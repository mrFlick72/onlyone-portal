package kafka

import (
	"context"
	"encoding/json"

	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	CreateAction = "CREATE"
	UpdateAction = "UPDATE"
	DeleteAction = "DELETE"
)

type KafkaBudgetExpenseEventPublisher struct {
	Topic  string
	Client *kgo.Client
	logger *logging.Logger
}

// BudgetExpenseEvent is the wire contract published to Kafka.
//
// The amount and date are rendered as strings: the domain Money and Date value
// objects keep their state in unexported fields with no JSON representation of
// their own, so marshaling the raw domain entity would drop them silently.
type BudgetExpenseEvent struct {
	Action  string            `json:"action"`
	Payload BudgetExpenseBody `json:"payload"`
}

type BudgetExpenseBody struct {
	Id       string           `json:"id"`
	UserName string           `json:"userName"`
	Date     string           `json:"date"`
	Amount   string           `json:"amount"`
	Note     string           `json:"note"`
	Tags     []tags.SearchTag `json:"tags"`
}

func NewKafkaBudgetExpenseEventPublisher(topic string, client *kgo.Client) expense.BudgetExpenseEventPublisher {
	return &KafkaBudgetExpenseEventPublisher{
		Topic:  topic,
		Client: client,
		logger: logging.GetLoggerInstanceForComponentByType(&KafkaBudgetExpenseEventPublisher{}),
	}
}

func (eventPublisher *KafkaBudgetExpenseEventPublisher) CreateBudgetExpense(ctx context.Context, budgetExpense expense.BudgetExpense) error {
	return eventPublisher.publishEventFor(ctx, CreateAction, budgetExpense)
}

func (eventPublisher *KafkaBudgetExpenseEventPublisher) UpdateBudgetExpense(ctx context.Context, budgetExpense expense.BudgetExpense) error {
	return eventPublisher.publishEventFor(ctx, UpdateAction, budgetExpense)
}

func (eventPublisher *KafkaBudgetExpenseEventPublisher) DeleteBudgetExpense(ctx context.Context, budgetExpense expense.BudgetExpense) error {
	return eventPublisher.publishEventFor(ctx, DeleteAction, budgetExpense)
}

func (eventPublisher *KafkaBudgetExpenseEventPublisher) publishEventFor(ctx context.Context, action string, budgetExpense expense.BudgetExpense) error {
	encoded, err := encodeBudgetExpenseEvent(action, budgetExpense)
	if err != nil {
		eventPublisher.logger.LogErrorfFor("Kafka event encoding fails: %v", err)
		return err
	}

	record := &kgo.Record{Topic: eventPublisher.Topic, Value: encoded}
	result := eventPublisher.Client.ProduceSync(ctx, record)

	if err := result.FirstErr(); err != nil {
		eventPublisher.logger.LogErrorfFor("Kafka Messages publishing fails: %v", err)
		return err
	}
	return nil
}

func encodeBudgetExpenseEvent(action string, budgetExpense expense.BudgetExpense) ([]byte, error) {
	return json.Marshal(newBudgetExpenseEvent(action, budgetExpense))
}

func newBudgetExpenseEvent(action string, budgetExpense expense.BudgetExpense) BudgetExpenseEvent {
	return BudgetExpenseEvent{
		Action: action,
		Payload: BudgetExpenseBody{
			Id:       budgetExpense.Id,
			UserName: budgetExpense.UserName,
			Date:     budgetExpense.Date.GetFormattedDate(),
			Amount:   budgetExpense.Amount.StringifyAmount(),
			Note:     budgetExpense.Note,
			Tags:     budgetExpense.Tags,
		},
	}
}
