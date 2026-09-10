package test

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"

	"github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"
)

func TestReadAllSuite(t *testing.T) {
	suite.Run(t, new(ReadAllTestSuite))
}

type ReadAllTestSuite struct {
	suite.Suite
	fixture *ClientFixture
}

func (s *ReadAllTestSuite) SetupTest() {
	s.fixture = NewInsecureClientFixture(s.T())
}

func (s *ReadAllTestSuite) TestReadAllEventsForwardsFromZeroPosition() {
	fixture := s.fixture
	client := fixture.Client()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	numberOfEvents := uint64(10)

	opts := trogoneventstore.ReadAllOptions{
		Direction:      trogoneventstore.Forwards,
		From:           trogoneventstore.Start{},
		ResolveLinkTos: true,
	}
	stream, err := client.ReadAll(ctx, opts, numberOfEvents)
	if err != nil {
		s.T().Fatalf("Unexpected failure %+v", err)
	}

	defer stream.Close()

	events, err := fixture.CollectEvents(stream)
	if err != nil {
		s.T().Fatalf("Unexpected failure %+v", err)
	}

	s.Equal(numberOfEvents, uint64(len(events)), "Expected the correct number of messages to be returned")
}
func (s *ReadAllTestSuite) TestReadAllEventsForwardsFromNonZeroPosition() {
	fixture := s.fixture
	client := fixture.Client()

	streamId := fixture.NewStreamId()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	initialEvent := fixture.CreateTestEvent()
	result, err := fixture.client.AppendToStream(ctx, streamId, trogoneventstore.AppendToStreamOptions{}, initialEvent)

	fixture.CreateTestEvents(streamId, 1000)

	s.NoError(err)
	opts := trogoneventstore.ReadAllOptions{
		From:           trogoneventstore.Position{Commit: result.CommitPosition, Prepare: result.PreparePosition},
		ResolveLinkTos: true,
	}

	stream, err := client.ReadAll(ctx, opts, 1000)
	if err != nil {
		s.T().Fatalf("Unexpected failure %+v", err)
	}

	defer stream.Close()

	events, err := fixture.CollectEvents(stream)
	s.NoError(err)

	s.Equal(events[0].OriginalEvent().Position.Commit, result.CommitPosition)
	s.Equal(events[0].OriginalEvent().EventID, initialEvent.EventID)
	s.True(events[1].OriginalEvent().Position.Commit > result.CommitPosition)
}
func (s *ReadAllTestSuite) TestReadAllEventsBackwardsFromZeroPosition() {
	fixture := s.fixture
	client := fixture.Client()

	streamId := fixture.NewStreamId()
	testEvents := fixture.CreateTestEvents(streamId, 1000)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	opts := trogoneventstore.ReadAllOptions{
		From:           trogoneventstore.End{},
		Direction:      trogoneventstore.Backwards,
		ResolveLinkTos: true,
	}

	stream, err := client.ReadAll(ctx, opts, 1000)
	if err != nil {
		s.T().Fatalf("Unexpected failure %+v", err)
	}

	defer stream.Close()

	collectedEvents, err := fixture.CollectEvents(stream)
	if err != nil {
		s.T().Fatalf("Unexpected failure %+v", err)
	}

	s.NotEmpty(collectedEvents)
	s.Len(collectedEvents, 1000)

	found := false
	for _, collectedEvent := range collectedEvents {
		for _, testEvent := range testEvents {
			if collectedEvent.OriginalEvent().EventID == testEvent.EventID {
				found = true
				break
			}
		}
	}
	s.True(found, "Expected to find at least one event from the testEvents")
}
func (s *ReadAllTestSuite) TestReadAllEventsBackwardsFromNonZeroPosition() {
	fixture := s.fixture
	client := fixture.Client()

	streamId := fixture.NewStreamId()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	testEvents := fixture.CreateTestEvents(streamId, 1000)

	initialEvent := fixture.CreateTestEvent()
	result, err := fixture.client.AppendToStream(ctx, streamId, trogoneventstore.AppendToStreamOptions{}, initialEvent)
	s.NoError(err)

	fixture.CreateTestEvents(streamId, 10)

	opts := trogoneventstore.ReadAllOptions{
		From:           trogoneventstore.Position{Commit: result.CommitPosition, Prepare: result.PreparePosition},
		Direction:      trogoneventstore.Backwards,
		ResolveLinkTos: true,
	}

	stream, err := client.ReadAll(ctx, opts, 1000)
	if err != nil {
		s.T().Fatalf("Unexpected failure %+v", err)
	}

	defer stream.Close()

	collectedEvents, err := fixture.CollectEvents(stream)
	if err != nil {
		s.T().Fatalf("Unexpected failure %+v", err)
	}

	found := false
	for _, collectedEvent := range collectedEvents {
		for _, testEvent := range testEvents {
			if collectedEvent.OriginalEvent().EventID == testEvent.EventID {
				found = true
				break
			}
		}
	}
	s.True(found, "Expected to find at least one event from the testEvents")
}

func (s *ReadAllTestSuite) TestReadAllEventsWithCredentialsOverride() {
	fixture := s.fixture
	client := fixture.Client()
	streamId := fixture.NewStreamId()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	result, err := fixture.client.AppendToStream(ctx, streamId, trogoneventstore.AppendToStreamOptions{}, fixture.CreateTestEvent())
	s.NoError(err)

	fixture.CreateTestEvents(streamId, 10)

	opts := trogoneventstore.ReadAllOptions{
		Authenticated: &trogoneventstore.Credentials{
			Login:    "admin",
			Password: "changeit",
		},
		From:           trogoneventstore.Position{Commit: result.CommitPosition, Prepare: result.PreparePosition},
		Direction:      trogoneventstore.Forwards,
		ResolveLinkTos: false,
	}

	stream, err := client.ReadAll(ctx, opts, 10)
	if err != nil {
		s.T().Fatalf("Unexpected failure %+v", err)
	}

	defer stream.Close()

	_, err = fixture.CollectEvents(stream)
	if err != nil {
		s.T().Fatalf("Unexpected failure %+v", err)
	}
}
