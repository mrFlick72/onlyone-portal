package kafka

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	testTopicPartitions        = 3
	testTopicReplicationFactor = 1
)

// These tests publish through a real broker. Bring one up first:
//
//	cd budget && docker compose up -d kafka
//
// The broker address is taken from KAFKA_BROKERS (default localhost:9092). When
// no broker is reachable the whole suite is skipped rather than failing, so the
// fast unit tests in this package still run on a machine without Kafka.

func kafkaBrokers() []string {
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		return []string{brokers}
	}
	return []string{"localhost:9092"}
}

func newKafkaClient(t *testing.T, opts ...kgo.Opt) *kgo.Client {
	t.Helper()

	client, err := kgo.NewClient(append([]kgo.Opt{kgo.SeedBrokers(kafkaBrokers()...)}, opts...)...)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx); err != nil {
		client.Close()
	}
	return client
}

// createTestTopic provisions a unique topic up front so the tests do not depend
// on the broker's auto-create setting (a real cluster commonly disables it). The
// topic is removed on test cleanup.
func createTestTopic(t *testing.T, client *kgo.Client) string {
	t.Helper()

	topic := "budget-api.expense.test." + uuid.NewString()
	admin := kadm.NewClient(client)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := admin.CreateTopic(ctx, testTopicPartitions, testTopicReplicationFactor, nil, topic)
	require.NoError(t, err)
	require.NoError(t, resp.Err)

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, _ = admin.DeleteTopics(ctx, topic)
	})
	return topic
}

func consumeAll(t *testing.T, topic string, expected int) []BudgetExpenseEvent {
	t.Helper()

	consumer := newKafkaClient(t,
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	defer consumer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	events := make([]BudgetExpenseEvent, 0, expected)
	for len(events) < expected {
		fetches := consumer.PollFetches(ctx)
		if err := ctx.Err(); err != nil {
			t.Fatalf("timed out waiting for %d records, got %d: %v", expected, len(events), err)
		}
		require.Empty(t, fetches.Errors())
		fetches.EachRecord(func(record *kgo.Record) {
			var event BudgetExpenseEvent
			require.NoError(t, json.Unmarshal(record.Value, &event))
			events = append(events, event)
		})
	}
	return events
}

func TestPublishBudgetExpenseEventsRoundTrip(t *testing.T) {
	producer := newKafkaClient(t)
	defer producer.Close()

	topic := createTestTopic(t, producer)
	publisher := NewKafkaBudgetExpenseEventPublisher(topic, producer)
	budgetExpense := aBudgetExpense()
	ctx := context.Background()

	require.NoError(t, publisher.CreateBudgetExpense(ctx, budgetExpense))
	require.NoError(t, publisher.UpdateBudgetExpense(ctx, budgetExpense))
	require.NoError(t, publisher.DeleteBudgetExpense(ctx, budgetExpense))

	events := consumeAll(t, topic, 3)

	actions := make([]string, 0, len(events))
	for _, event := range events {
		actions = append(actions, event.Action)
		assert.Equal(t, "10.50", event.Payload.Amount)
		assert.Equal(t, "01/02/2024", event.Payload.Date)
		assert.Equal(t, "testuser", event.Payload.UserName)
	}
	assert.ElementsMatch(t, []string{CreateAction, UpdateAction, DeleteAction}, actions)
}

func TestPublishBudgetExpenseEventToClosedClientReturnsError(t *testing.T) {
	producer := newKafkaClient(t)
	producer.Close()

	publisher := NewKafkaBudgetExpenseEventPublisher("budget-api.expense.test."+uuid.NewString(), producer)

	err := publisher.CreateBudgetExpense(context.Background(), aBudgetExpense())
	assert.Error(t, err)
}
