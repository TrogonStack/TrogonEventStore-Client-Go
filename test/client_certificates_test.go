package test

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"io"
//	"path/filepath"
//	"testing"
//
//	"github.com/google/uuid"
//	"github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//)
//
//func TestClientCertificates(t *testing.T) {
//	endpoint := "localhost:2111,localhost:2112,localhost:2113"
//
//	t.Run("ValidClientCertificatesTlsRelativePath", func(t *testing.T) {
//		tlsCaFile := "../certs/ca/ca.crt"
//		userCertFile := "../certs/user-admin/user-admin.crt"
//		userKeyFile := "../certs/user-admin/user-admin.key"
//
//		config, err := trogoneventstore.ParseConnectionString(fmt.Sprintf("esdb://admin:changeit@%s?tls=true&tlscafile=%s&userCertFile=%s&userKeyFile=%s", endpoint, tlsCaFile, userCertFile, userKeyFile))
//		if err != nil {
//			t.Fatalf("Unexpected configuration error: %s", err.Error())
//		}
//
//		c, err := trogoneventstore.NewClient(config)
//		if err != nil {
//			t.Fatalf("Unexpected error: %s", err.Error())
//		}
//		defer c.Close()
//
//		numberOfEventsToRead := 1
//		numberOfEvents := uint64(numberOfEventsToRead)
//		opts := trogoneventstore.ReadAllOptions{
//			From:           trogoneventstore.Start{},
//			Direction:      trogoneventstore.Backwards,
//			ResolveLinkTos: true,
//		}
//
//		stream, err := c.ReadAll(context.Background(), opts, numberOfEvents)
//		require.NoError(t, err)
//		defer stream.Close()
//		evt, err := stream.Recv()
//		require.Nil(t, evt)
//		require.True(t, errors.Is(err, io.EOF))
//	})
//
//	t.Run("ValidClientCertificatesTlsWithAbsolutePath", func(t *testing.T) {
//		userCertFile, err := filepath.Abs("../certs/user-admin/user-admin.crt")
//		if err != nil {
//			t.Fatalf("Unexpected error: %s", err.Error())
//		}
//
//		userKeyFile, err := filepath.Abs("../certs/user-admin/user-admin.key")
//		if err != nil {
//			t.Fatalf("Unexpected error: %s", err.Error())
//		}
//
//		tlsCaFile, err := filepath.Abs("../certs/ca/ca.crt")
//		if err != nil {
//			t.Fatalf("Unexpected error: %s", err.Error())
//		}
//
//		config, err := trogoneventstore.ParseConnectionString(fmt.Sprintf("esdb://admin:changeit@%s?tls=true&tlscafile=%s&usercertfile=%s&userkeyfile=%s", endpoint, tlsCaFile, userCertFile, userKeyFile))
//		if err != nil {
//			t.Fatalf("Unexpected configuration error: %s", err.Error())
//		}
//
//		c, err := trogoneventstore.NewClient(config)
//		if err != nil {
//			t.Fatalf("Unexpected error: %s", err.Error())
//		}
//		defer c.Close()
//
//		numberOfEventsToRead := 1
//		numberOfEvents := uint64(numberOfEventsToRead)
//		opts := trogoneventstore.ReadAllOptions{
//			From:           trogoneventstore.Start{},
//			Direction:      trogoneventstore.Backwards,
//			ResolveLinkTos: true,
//		}
//
//		stream, err := c.ReadAll(context.Background(), opts, numberOfEvents)
//		require.NoError(t, err)
//		defer stream.Close()
//		evt, err := stream.Recv()
//		require.Nil(t, evt)
//		require.True(t, errors.Is(err, io.EOF))
//	})
//
//	t.Run("InvalidUserCertificates", func(t *testing.T) {
//		tlsCaFile := "../certs/ca/ca.crt"
//		userCertFile := "../certs/user-invalid/user-invalid.crt"
//		userKeyFile := "../certs/user-invalid/user-invalid.key"
//
//		config, err := trogoneventstore.ParseConnectionString(fmt.Sprintf("esdb://admin:changeit@%s?tls=true&tlscafile=%s&usercertfile=%s&userkeyfile=%s", endpoint, tlsCaFile, userCertFile, userKeyFile))
//		if err != nil {
//			t.Fatalf("Unexpected configuration error: %s", err.Error())
//		}
//
//		c, err := trogoneventstore.NewClient(config)
//		if err != nil {
//			t.Fatalf("Unexpected error: %s", err.Error())
//		}
//		defer c.Close()
//
//		testEvent := fixture.CreateTestEvent()
//
//		streamID := uuid.NewString()
//		opts := trogoneventstore.AppendToStreamOptions{
//			StreamState: trogoneventstore.Any{},
//		}
//
//		result, err := c.AppendToStream(context.Background(), streamID, opts, testEvent)
//		require.Nil(t, result)
//		require.Error(t, err)
//		require.Contains(t, err.Error(), "Unauthenticated")
//	})
//
//	t.Run("MissingCertificateFile", func(t *testing.T) {
//		tlsCaFile := "../certs/ca/ca.crt"
//		userCertFile := "../certs/user-admin/user-admin.crt"
//
//		config, err := trogoneventstore.ParseConnectionString(fmt.Sprintf("esdb://admin:changeit@%s?tls=true&tlscafile=%s&usercertfile=%s", endpoint, tlsCaFile, userCertFile))
//
//		_, err = trogoneventstore.NewClient(config)
//		serverError, ok := trogoneventstore.FromError(err)
//		require.False(t, ok)
//		require.NotNil(t, serverError)
//		assert.Contains(t, serverError.Error(), "both userCertFile and userKeyFile must be provided")
//	})
//}
