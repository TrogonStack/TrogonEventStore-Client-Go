package samples

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"
)

func createPersistentSubscription(client *trogoneventstore.Client) {
	// #region create-persistent-subscription-to-stream
	err := client.CreatePersistentSubscription(context.Background(), "test-stream", "subscription-group", trogoneventstore.PersistentStreamSubscriptionOptions{})

	if err != nil {
		panic(err)
	}
	// #endregion create-persistent-subscription-to-stream
}

func connectToPersistentSubscriptionToStream(client *trogoneventstore.Client) {
	// #region subscribe-to-persistent-subscription-to-stream
	sub, err := client.SubscribeToPersistentSubscription(context.Background(), "test-stream", "subscription-group", trogoneventstore.SubscribeToPersistentSubscriptionOptions{})

	if err != nil {
		panic(err)
	}

	for {
		event := sub.Recv()

		if event.EventAppeared != nil {
			sub.Ack(event.EventAppeared.Event)
		}

		if event.SubscriptionDropped != nil {
			break
		}
	}
	// #endregion subscribe-to-persistent-subscription-to-stream
}

func connectToPersistentSubscriptionToAll(client *trogoneventstore.Client) {
	// #region subscribe-to-persistent-subscription-to-all
	sub, err := client.SubscribeToPersistentSubscriptionToAll(context.Background(), "subscription-group", trogoneventstore.SubscribeToPersistentSubscriptionOptions{})

	if err != nil {
		panic(err)
	}

	for {
		event := sub.Recv()

		if event.EventAppeared != nil {
			sub.Ack(event.EventAppeared.Event)
		}

		if event.SubscriptionDropped != nil {
			break
		}
	}
	// #endregion subscribe-to-persistent-subscription-to-all
}

func createPersistentSubscriptionToAll(client *trogoneventstore.Client) {
	// #region create-persistent-subscription-to-all
	options := trogoneventstore.PersistentAllSubscriptionOptions{
		Filter: &trogoneventstore.SubscriptionFilter{
			Type:     trogoneventstore.StreamFilterType,
			Prefixes: []string{"test"},
		},
	}

	err := client.CreatePersistentSubscriptionToAll(context.Background(), "subscription-group", options)

	if err != nil {
		panic(err)
	}
	// #endregion create-persistent-subscription-to-all
}

func connectToPersistentSubscriptionWithManualAcks(client *trogoneventstore.Client) {
	// #region subscribe-to-persistent-subscription-with-manual-acks
	sub, err := client.SubscribeToPersistentSubscription(context.Background(), "test-stream", "subscription-group", trogoneventstore.SubscribeToPersistentSubscriptionOptions{})

	if err != nil {
		panic(err)
	}

	for {
		event := sub.Recv()

		if event.EventAppeared != nil {
			sub.Ack(event.EventAppeared.Event)
		}

		if event.SubscriptionDropped != nil {
			break
		}
	}
	// #endregion subscribe-to-persistent-subscription-with-manual-acks
}

func updatePersistentSubscription(client *trogoneventstore.Client) {
	// #region update-persistent-subscription
	options := trogoneventstore.PersistentStreamSubscriptionOptions{
		Settings: &trogoneventstore.PersistentSubscriptionSettings{
			ResolveLinkTos:       true,
			CheckpointLowerBound: 20,
		},
	}

	err := client.UpdatePersistentSubscription(context.Background(), "test-stream", "subscription-group", options)

	if err != nil {
		panic(err)
	}
	// #endregion update-persistent-subscription
}

func deletePersistentSubscription(client *trogoneventstore.Client) {
	// #region delete-persistent-subscription
	err := client.DeletePersistentSubscription(context.Background(), "test-stream", "subscription-group", trogoneventstore.DeletePersistentSubscriptionOptions{})

	if err != nil {
		panic(err)
	}
	// #endregion delete-persistent-subscription
}

func deletePersistentSubscriptionToAll(client *trogoneventstore.Client) {
	// #region delete-persistent-subscription-all
	err := client.DeletePersistentSubscriptionToAll(context.Background(), "test-stream", trogoneventstore.DeletePersistentSubscriptionOptions{})

	if err != nil {
		panic(err)
	}
	// #endregion delete-persistent-subscription-all
}

func getPersistentSubscriptionToStreamInfo(client *trogoneventstore.Client) {
	// #region get-persistent-subscription-to-stream-info
	info, err := client.GetPersistentSubscriptionInfo(context.Background(), "test-stream", "subscription-group", trogoneventstore.GetPersistentSubscriptionOptions{})

	if err != nil {
		panic(err)
	}

	log.Printf("groupName: %s eventsource: %s status: %s", info.GroupName, info.EventSource, info.Status)
	// #endregion get-persistent-subscription-to-stream-info
}

func getPersistentSubscriptionToAllInfo(client *trogoneventstore.Client) {
	// #region get-persistent-subscription-to-all-info
	info, err := client.GetPersistentSubscriptionInfoToAll(context.Background(), "subscription-group", trogoneventstore.GetPersistentSubscriptionOptions{})

	if err != nil {
		panic(err)
	}

	log.Printf("groupName: %s eventsource: %s status: %s", info.GroupName, info.EventSource, info.Status)
	// #endregion get-persistent-subscription-to-all-info
}

func replayParkedToStream(client *trogoneventstore.Client) {
	// #region replay-parked-of-persistent-subscription-to-stream
	err := client.ReplayParkedMessages(context.Background(), "test-stream", "subscription-group", trogoneventstore.ReplayParkedMessagesOptions{
		StopAt: 10,
	})

	if err != nil {
		panic(err)
	}
	// #endregion replay-parked-of-persistent-subscription-to-stream
}

func replayParkedToAll(client *trogoneventstore.Client) {
	// #region replay-parked-of-persistent-subscription-to-all
	err := client.ReplayParkedMessagesToAll(context.Background(), "subscription-group", trogoneventstore.ReplayParkedMessagesOptions{
		StopAt: 10,
	})

	if err != nil {
		panic(err)
	}
	// #endregion replay-parked-of-persistent-subscription-to-all
}

func listPersistentSubscriptionsToStream(client *trogoneventstore.Client) {
	// #region list-persistent-subscriptions-to-stream
	subs, err := client.ListPersistentSubscriptionsForStream(context.Background(), "test-stream", trogoneventstore.ListPersistentSubscriptionsOptions{})

	if err != nil {
		panic(err)
	}

	var entries []string

	for i := range subs {
		entries = append(
			entries,
			fmt.Sprintf(
				"groupName: %s eventSource: %s status: %s",
				subs[i].GroupName,
				subs[i].EventSource,
				subs[i].Status,
			),
		)
	}

	log.Printf("subscriptions to stream: [ %s ]", strings.Join(entries, ","))
	// #endregion list-persistent-subscriptions-to-stream
}

func listPersistentSubscriptionsToAll(client *trogoneventstore.Client) {
	// #region list-persistent-subscriptions-to-all
	subs, err := client.ListPersistentSubscriptionsToAll(context.Background(), trogoneventstore.ListPersistentSubscriptionsOptions{})

	if err != nil {
		panic(err)
	}

	var entries []string

	for i := range subs {
		entries = append(
			entries,
			fmt.Sprintf(
				"groupName: %s eventSource: %s status: %s",
				subs[i].GroupName,
				subs[i].EventSource,
				subs[i].Status,
			),
		)
	}

	log.Printf("subscriptions to stream: [ %s ]", strings.Join(entries, ","))
	// #endregion list-persistent-subscriptions-to-all
}

func restartPersistentSubscriptionSubsystem(client *trogoneventstore.Client) {
	// #region restart-persistent-subscription-subsystem
	err := client.RestartPersistentSubscriptionSubsystem(context.Background(), trogoneventstore.RestartPersistentSubscriptionSubsystemOptions{})

	if err != nil {
		panic(err)
	}
	// #endregion restart-persistent-subscription-subsystem
}
