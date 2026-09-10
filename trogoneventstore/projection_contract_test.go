package trogoneventstore

import (
	"context"
	"net"
	"reflect"
	"sync"
	"testing"

	"github.com/TrogonStack/TrogonEventStore-Client-Go/protos/trogoneventstore/protocols/v1/projections"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type projectionCaptureServer struct {
	projections.UnimplementedProjectionsServer
	created chan *projections.CreateReq
	updated chan *projections.UpdateReq
}

func (s *projectionCaptureServer) Create(_ context.Context, request *projections.CreateReq) (*projections.CreateResp, error) {
	s.created <- request
	return &projections.CreateResp{}, nil
}

func (s *projectionCaptureServer) Update(_ context.Context, request *projections.UpdateReq) (*projections.UpdateResp, error) {
	s.updated <- request
	return &projections.UpdateResp{}, nil
}

func newProjectionCaptureClient(t *testing.T) (*ProjectionClient, *projectionCaptureServer) {
	t.Helper()

	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	capture := &projectionCaptureServer{
		created: make(chan *projections.CreateReq, 1),
		updated: make(chan *projections.UpdateReq, 1),
	}
	projections.RegisterProjectionsServer(server, capture)
	go func() {
		_ = server.Serve(listener)
	}()

	connection, err := grpc.NewClient(
		"passthrough:///projection-capture",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}

	closed := new(int32)
	connectionRequests := make(chan msg)
	go func() {
		for request := range connectionRequests {
			get, ok := request.(getConnection)
			if ok {
				get.channel <- connectionHandle{id: uuid.New(), connection: connection}
			}
		}
	}()

	inner := &Client{
		grpcClient: &grpcClient{
			channel:   connectionRequests,
			closeFlag: closed,
			once:      new(sync.Once),
		},
		config: &Configuration{DisableTLS: true},
	}
	t.Cleanup(func() {
		inner.grpcClient.close()
		_ = connection.Close()
		server.Stop()
		_ = listener.Close()
	})

	return NewProjectionClientFromExistingClient(inner), capture
}

func TestProjectionContractMatchesServer(t *testing.T) {
	fields := (&projections.CreateReq_Options{}).ProtoReflect().Descriptor().Fields()

	if fields.ByName("annotations") == nil {
		t.Fatal("projection create contract must expose annotations")
	}
	if fields.ByName("engine_version") != nil {
		t.Fatal("projection create contract must not expose an unsupported engine version")
	}

	optionsType := reflect.TypeOf(CreateProjectionOptions{})
	annotations, ok := optionsType.FieldByName("Annotations")
	if !ok || annotations.Type != reflect.TypeOf(map[string]string{}) {
		t.Fatal("projection options must expose string annotations")
	}
	if _, ok := optionsType.FieldByName("Metadata"); ok {
		t.Fatal("projection options must use the server contract name")
	}
}

func TestProjectionClientSendsAnnotations(t *testing.T) {
	client, capture := newProjectionCaptureClient(t)

	createAnnotations := map[string]string{"deployment": "blue", "owner": "inventory"}
	if err := client.Create(context.Background(), "inventory", "fromAll()", CreateProjectionOptions{
		Annotations: createAnnotations,
	}); err != nil {
		t.Fatal(err)
	}
	if got := (<-capture.created).GetOptions().GetAnnotations(); !reflect.DeepEqual(got, createAnnotations) {
		t.Fatalf("create annotations = %v, want %v", got, createAnnotations)
	}

	updateAnnotations := map[string]string{"deployment": "green"}
	if err := client.Update(context.Background(), "inventory", "fromAll()", UpdateProjectionOptions{
		Annotations: updateAnnotations,
	}); err != nil {
		t.Fatal(err)
	}
	if got := (<-capture.updated).GetOptions().GetAnnotations(); !reflect.DeepEqual(got, updateAnnotations) {
		t.Fatalf("update annotations = %v, want %v", got, updateAnnotations)
	}
}
