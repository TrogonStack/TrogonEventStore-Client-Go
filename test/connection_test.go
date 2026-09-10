package test

import (
	"context"
	"testing"
	"time"

	"github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestConnectionSuite(t *testing.T) {
	suite.Run(t, new(ConnectionTestSuite))
}

type ConnectionTestSuite struct {
	suite.Suite
	fixture *ClientFixture
}

func (s *ConnectionTestSuite) SetupSuite() {
	s.fixture = NewSecureSingleNodeClientFixture(s.T())
}

func (s *ConnectionTestSuite) TearDownSuite() {
	s.fixture.Close(s.T())
}

func (s *ConnectionTestSuite) TestCloseConnection() {
	client := s.fixture.Client()
	testEvent := s.fixture.CreateTestEvent()

	streamID := uuid.New()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := trogoneventstore.AppendToStreamOptions{
		StreamState: trogoneventstore.NoStream{},
	}

	_, err := client.AppendToStream(ctx, streamID.String(), opts, testEvent)
	s.NoError(err)

	// Close the client and attempt another append to verify connection is closed.
	client.Close()
	opts.StreamState = trogoneventstore.Any{}
	_, err = client.AppendToStream(ctx, streamID.String(), opts, testEvent)

	serverError, ok := trogoneventstore.FromError(err)
	assert.False(s.T(), ok)
	assert.Equal(s.T(), trogoneventstore.ErrorCodeConnectionClosed, serverError.Code())
}
