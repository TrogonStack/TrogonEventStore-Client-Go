package samples

import (
	"context"
	"time"

	"github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"
)

func SubscribeToStream(db *trogoneventstore.Client) {
	options := trogoneventstore.SubscribeToStreamOptions{}
	// region subscribe-to-stream
	stream, err := db.SubscribeToStream(context.Background(), "some-stream", trogoneventstore.SubscribeToStreamOptions{})

	if err != nil {
		panic(err)
	}

	defer stream.Close()

	for {
		event := stream.Recv()

		if event.EventAppeared != nil {
			// handles the event...
		}

		if event.SubscriptionDropped != nil {
			break
		}
	}
	// endregion subscribe-to-stream

	// region subscribe-to-stream-from-position
	db.SubscribeToStream(context.Background(), "some-stream", trogoneventstore.SubscribeToStreamOptions{
		From: trogoneventstore.Revision(20),
	})
	// endregion subscribe-to-stream-from-position

	// region subscribe-to-stream-live
	options = trogoneventstore.SubscribeToStreamOptions{
		From: trogoneventstore.End{},
	}

	db.SubscribeToStream(context.Background(), "some-stream", options)
	// endregion subscribe-to-stream-live

	// region subscribe-to-stream-resolving-linktos
	options = trogoneventstore.SubscribeToStreamOptions{
		From:           trogoneventstore.Start{},
		ResolveLinkTos: true,
	}

	db.SubscribeToStream(context.Background(), "$et-myEventType", options)
	// endregion subscribe-to-stream-resolving-linktos

	// region subscribe-to-stream-subscription-dropped
	options = trogoneventstore.SubscribeToStreamOptions{
		From: trogoneventstore.Start{},
	}

	for {

		stream, err := db.SubscribeToStream(context.Background(), "some-stream", options)

		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		for {
			event := stream.Recv()

			if event.SubscriptionDropped != nil {
				stream.Close()
				break
			}

			if event.EventAppeared != nil {
				// handles the event...
				options.From = trogoneventstore.Revision(event.EventAppeared.OriginalEvent().EventNumber)
			}
		}
	}
	// endregion subscribe-to-stream-subscription-dropped
}

func SubscribeToAll(db *trogoneventstore.Client) {
	options := trogoneventstore.SubscribeToAllOptions{}
	// region subscribe-to-all
	stream, err := db.SubscribeToAll(context.Background(), trogoneventstore.SubscribeToAllOptions{})

	if err != nil {
		panic(err)
	}

	defer stream.Close()

	for {
		event := stream.Recv()

		if event.EventAppeared != nil {
			// handles the event...
		}

		if event.SubscriptionDropped != nil {
			break
		}
	}
	// endregion subscribe-to-all

	// region subscribe-to-all-from-position
	db.SubscribeToAll(context.Background(), trogoneventstore.SubscribeToAllOptions{
		From: trogoneventstore.Position{
			Commit:  1_056,
			Prepare: 1_056,
		},
	})
	// endregion subscribe-to-all-from-position

	// region subscribe-to-all-live
	db.SubscribeToAll(context.Background(), trogoneventstore.SubscribeToAllOptions{
		From: trogoneventstore.End{},
	})
	// endregion subscribe-to-all-live

	// region subscribe-to-all-subscription-dropped
	options = trogoneventstore.SubscribeToAllOptions{
		From: trogoneventstore.Start{},
	}

	for {
		stream, err := db.SubscribeToAll(context.Background(), options)

		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		for {
			event := stream.Recv()

			if event.SubscriptionDropped != nil {
				stream.Close()
				break
			}

			if event.EventAppeared != nil {
				// handles the event...
				options.From = event.EventAppeared.OriginalEvent().Position
			}
		}
	}
	// endregion subscribe-to-all-subscription-dropped
}

func SubscribeToFiltered(db *trogoneventstore.Client) {
	// region stream-prefix-filtered-subscription
	db.SubscribeToAll(context.Background(), trogoneventstore.SubscribeToAllOptions{
		Filter: &trogoneventstore.SubscriptionFilter{
			Type:     trogoneventstore.StreamFilterType,
			Prefixes: []string{"test-"},
		},
	})
	// endregion stream-prefix-filtered-subscription
	// region stream-regex-filtered-subscription
	db.SubscribeToAll(context.Background(), trogoneventstore.SubscribeToAllOptions{
		Filter: &trogoneventstore.SubscriptionFilter{
			Type:  trogoneventstore.StreamFilterType,
			Regex: "/invoice-\\d\\d\\d/g",
		},
	})
	// endregion stream-regex-filtered-subscription
}

func SubscribeToAllOverridingUserCredentials(db *trogoneventstore.Client) {
	// region overriding-user-credentials
	db.SubscribeToAll(context.Background(), trogoneventstore.SubscribeToAllOptions{
		Authenticated: &trogoneventstore.Credentials{
			Login:    "admin",
			Password: "changeit",
		},
	})
	// endregion overriding-user-credentials
}
