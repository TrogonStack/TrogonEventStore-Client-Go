package samples

import (
	"context"
	"fmt"

	"github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"
)

func ExcludeSystemEvents(db *trogoneventstore.Client) {
	// region exclude-system
	sub, err := db.SubscribeToAll(context.Background(), trogoneventstore.SubscribeToAllOptions{
		Filter: trogoneventstore.ExcludeSystemEventsFilter(),
	})

	if err != nil {
		panic(err)
	}

	defer sub.Close()

	for {
		event := sub.Recv()

		if event.EventAppeared != nil {
			streamId := event.EventAppeared.OriginalEvent().StreamID
			revision := event.EventAppeared.OriginalEvent().EventNumber

			fmt.Printf("received event %v@%v", revision, streamId)
		}

		if event.SubscriptionDropped != nil {
			break
		}
	}

	// endregion exclude-system
}

func EventTypePrefix(db *trogoneventstore.Client) {
	// region event-type-prefix
	sub, err := db.SubscribeToAll(context.Background(), trogoneventstore.SubscribeToAllOptions{
		Filter: &trogoneventstore.SubscriptionFilter{
			Type:     trogoneventstore.EventFilterType,
			Prefixes: []string{"customer-"},
		},
	})

	if err != nil {
		panic(err)
	}

	defer sub.Close()

	for {
		event := sub.Recv()

		if event.EventAppeared != nil {
			streamId := event.EventAppeared.OriginalEvent().StreamID
			revision := event.EventAppeared.OriginalEvent().EventNumber

			fmt.Printf("received event %v@%v", revision, streamId)
		}

		if event.SubscriptionDropped != nil {
			break
		}
	}

	// endregion event-type-prefix
}

func EventTypeRegex(db *trogoneventstore.Client) {
	// region event-type-regex
	sub, err := db.SubscribeToAll(context.Background(), trogoneventstore.SubscribeToAllOptions{
		Filter: &trogoneventstore.SubscriptionFilter{
			Type:  trogoneventstore.EventFilterType,
			Regex: "^user|^company",
		},
	})

	if err != nil {
		panic(err)
	}

	defer sub.Close()

	for {
		event := sub.Recv()

		if event.EventAppeared != nil {
			streamId := event.EventAppeared.OriginalEvent().StreamID
			revision := event.EventAppeared.OriginalEvent().EventNumber

			fmt.Printf("received event %v@%v", revision, streamId)
		}

		if event.SubscriptionDropped != nil {
			break
		}
	}

	// endregion event-type-regex
}

func StreamPrefix(db *trogoneventstore.Client) {
	// region stream-prefix
	sub, err := db.SubscribeToAll(context.Background(), trogoneventstore.SubscribeToAllOptions{
		Filter: &trogoneventstore.SubscriptionFilter{
			Type:     trogoneventstore.StreamFilterType,
			Prefixes: []string{"user-"},
		},
	})

	if err != nil {
		panic(err)
	}

	defer sub.Close()

	for {
		event := sub.Recv()

		if event.EventAppeared != nil {
			streamId := event.EventAppeared.OriginalEvent().StreamID
			revision := event.EventAppeared.OriginalEvent().EventNumber

			fmt.Printf("received event %v@%v", revision, streamId)
		}

		if event.SubscriptionDropped != nil {
			break
		}
	}

	// endregion stream-prefix
}

func StreamRegex(db *trogoneventstore.Client) {
	// region stream-regex
	sub, err := db.SubscribeToAll(context.Background(), trogoneventstore.SubscribeToAllOptions{
		Filter: &trogoneventstore.SubscriptionFilter{
			Type:  trogoneventstore.StreamFilterType,
			Regex: "^user|^company",
		},
	})

	if err != nil {
		panic(err)
	}

	defer sub.Close()

	for {
		event := sub.Recv()

		if event.EventAppeared != nil {
			streamId := event.EventAppeared.OriginalEvent().StreamID
			revision := event.EventAppeared.OriginalEvent().EventNumber

			fmt.Printf("received event %v@%v", revision, streamId)
		}

		if event.SubscriptionDropped != nil {
			break
		}
	}

	// endregion stream-regex
}

func CheckpointCallbackWithInterval(db *trogoneventstore.Client) {
	// region checkpoint-with-interval
	sub, err := db.SubscribeToAll(context.Background(), trogoneventstore.SubscribeToAllOptions{
		Filter: &trogoneventstore.SubscriptionFilter{
			Type:  trogoneventstore.EventFilterType,
			Regex: "/^[^\\$].*/",
		},
	})
	// endregion checkpoint-with-interval
	if err != nil {
		panic(err)
	}

	defer sub.Close()

	// region checkpoint
	for {
		event := sub.Recv()

		if event.EventAppeared != nil {
			streamId := event.EventAppeared.OriginalEvent().StreamID
			revision := event.EventAppeared.OriginalEvent().EventNumber

			fmt.Printf("received event %v@%v", revision, streamId)
		}

		if event.CheckPointReached != nil {
			// Save commit position to a persistent store as a checkpoint
			fmt.Printf("checkpoint taken at %v", event.CheckPointReached.Commit)
		}

		if event.SubscriptionDropped != nil {
			break
		}
	}

	// endregion checkpoint
}
