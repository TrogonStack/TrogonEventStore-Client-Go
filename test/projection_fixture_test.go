package test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"
)

type ProjectFixture struct {
	*ClientFixture
	projectionClient *trogoneventstore.ProjectionClient
}

func NewProjectFixture(t *testing.T, clientFixture *ClientFixture) *ProjectFixture {
	t.Helper()

	projectionClient := trogoneventstore.NewProjectionClientFromExistingClient(clientFixture.Client())

	return &ProjectFixture{
		projectionClient: projectionClient,
		ClientFixture:    clientFixture,
	}
}

func (f *ProjectFixture) WaitUntilProjectionStatusIs(t *testing.T, timeout time.Duration, name string, targetStatus string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("Timeout waiting for status %s", targetStatus)
		default:
			status, err := f.projectionClient.GetStatus(context.Background(), name, trogoneventstore.GenericProjectionOptions{})
			if err != nil {
				t.Logf("Status check error: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			if strings.Contains(status.Status, targetStatus) {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (f *ProjectFixture) WaitUntilProjectionStateReady(t *testing.T, duration time.Duration, name string) {
	done := make(chan *state)
	go func() {
		for {
			result, err := f.projectionClient.GetState(context.Background(), name, trogoneventstore.GetStateProjectionOptions{})

			if err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			var state state
			content, err := result.MarshalJSON()

			if err != nil {
				t.Errorf("error when unmarshalling protobuf types to json: %v", err)
			}

			err = json.Unmarshal(content, &state)
			if err != nil {
				t.Errorf("error when unmarshalling projection internal state type: %v", err)
			}

			done <- &state
			return
		}
	}()

	select {
	case <-done:
		return
	case <-time.After(duration):
		t.Errorf("unable to get projection '%s' internal state in a timely manner", name)
	}
}

func (f *ProjectFixture) WaitUntilProjectionResultReady(t *testing.T, duration time.Duration, name string) {
	done := make(chan *state)
	go func() {
		for {
			result, err := f.projectionClient.GetResult(context.Background(), name, trogoneventstore.GetResultProjectionOptions{})

			if err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			var state state
			content, err := result.MarshalJSON()

			if err != nil {
				t.Errorf("error when unmarshalling protobuf types to json: %v", err)
			}

			err = json.Unmarshal(content, &state)
			if err != nil {
				t.Errorf("error when unmarshalling projection internal result type: %v", err)
			}

			done <- &state
			return
		}
	}()

	select {
	case <-done:
		return
	case <-time.After(duration):
		t.Errorf("unable to get projection '%s' internal state in a timely manner", name)
	}
}

func (f *ProjectFixture) Close(t *testing.T) {
	if f.projectionClient != nil {
		if err := f.projectionClient.Close(); err != nil {
			t.Error("Error closing client:", err)
		}
	}

	f.ClientFixture.Close(t)
}

type state struct {
	Foo struct {
		Baz struct {
			Count float64 `json:"count"`
		} `json:"baz"`
	} `json:"foo"`
}
