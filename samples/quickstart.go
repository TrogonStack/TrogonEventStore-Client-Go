package samples

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"
)

func Run() {
	// region createClient
	settings, err := trogoneventstore.ParseConnectionString("{connectionString}")

	if err != nil {
		panic(err)
	}

	db, err := trogoneventstore.NewClient(settings)

	// endregion createClient
	if err != nil {
		panic(err)
	}

	// region createEvent
	testEvent := TestEvent{
		Id:            uuid.NewString(),
		ImportantData: "I wrote my first event!",
	}

	data, err := json.Marshal(testEvent)

	if err != nil {
		panic(err)
	}

	eventData := trogoneventstore.EventData{
		ContentType: trogoneventstore.ContentTypeJson,
		EventType:   "TestEvent",
		Data:        data,
	}
	// endregion createEvent

	// region appendEvents
	_, err = db.AppendToStream(context.Background(), "some-stream", trogoneventstore.AppendToStreamOptions{}, eventData)
	// endregion appendEvents

	if err != nil {
		panic(err)
	}

	// region readStream
	stream, err := db.ReadStream(context.Background(), "some-stream", trogoneventstore.ReadStreamOptions{}, 10)

	if err != nil {
		panic(err)
	}

	defer stream.Close()

	for {
		event, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			panic(err)
		}

		// Doing something productive with the event
		fmt.Println(event)
	}
	// endregion readStream
}
