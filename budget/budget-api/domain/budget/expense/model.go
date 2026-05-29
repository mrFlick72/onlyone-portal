package expense

import (
	"slices"

	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
)

type SpentBudget struct {
	BudgetExpenseList []BudgetExpense
	SearchTags        map[tags.SearchTagKey]tags.SearchTagValue
}


func NewSpentBudget(
	budgetExpenseList []BudgetExpense,
	searchTags []tags.SearchTag) *SpentBudget {

	return &SpentBudget{
		BudgetExpenseList: budgetExpenseList,
		SearchTags:        adaptSearchTagFormListToMap(searchTags),
	}
}

func adaptSearchTagFormListToMap(searchTags []tags.SearchTag) map[tags.SearchTagKey]tags.SearchTagValue {
	result := make(map[tags.SearchTagKey]tags.SearchTagValue)

	for _, searchTag := range searchTags {
		result[searchTag.Key] = searchTag.Value
	}

	return result
}

func (spentBudget *SpentBudget) Total() money.Money {
	total := money.Zero()
	for _, budgetExpense := range spentBudget.BudgetExpenseList {
		total = total.Plus(budgetExpense.Amount)
	}
	return total
}


func (spentBudget *SpentBudget) TotalForSearchTags() map[tags.SearchTag]money.Money {
	result := make(map[tags.SearchTag]money.Money)
	for _, budgetExpense := range spentBudget.BudgetExpenseList {
		searchTag := spentBudget.findSearchTagFor(budgetExpense.Tag.Key)
		if searchTag != nil {
			currentTotal, exists := result[*searchTag]
			if !exists {
				currentTotal = money.Zero()
			}
			result[*searchTag] = currentTotal.Plus(budgetExpense.Amount)
		}
	}
	return result
}

func (spentBudget *SpentBudget) findSearchTagFor(searchTagKey string) *tags.SearchTag {
	value, ok := spentBudget.SearchTags[searchTagKey]
	if !ok {
		return nil
	}
	return &tags.SearchTag{
		Key:   searchTagKey,
		Value: value,
	}
}

func (spentBudget *SpentBudget) DailyBudgetExpenseList() []DailyBudgetExpense {
	result := make(map[date.Date][]BudgetExpense)
	for _, budgetExpense := range spentBudget.BudgetExpenseList {
		key := budgetExpense.Date
		value, exists := result[key]
		if !exists {
			value = []BudgetExpense{}
		}
		result[key] = append(value, budgetExpense)
	}
	resultList := make([]DailyBudgetExpense, 0)
	for key, value := range result {
		total := money.Zero()
		for _, budgetExpense := range value {
			total = total.Plus(budgetExpense.Amount)
		}
		resultList = append(resultList, DailyBudgetExpense{
			BudgetExpenseList: &value,
			Date:              key,
			Total:             total,
		})
	}
		slices.SortFunc(resultList, func(a, b DailyBudgetExpense) int {
			return a.Date.GetTime().Compare(b.Date.GetTime())
		})
	return resultList
}

type BudgetExpense struct {
	Id       BudgetExpenseId
	UserName UserName
	Date     date.Date
	Amount   money.Money
	Note     string
	Tag      tags.SearchTag
}

type BudgetExpenseId = string
type UserName = string

type BudgetExpenseIdProvider interface {
	GenerateIdFor(budgetExpense *BudgetExpense) BudgetExpenseId
}


type DailyBudgetExpense struct {
	BudgetExpenseList *[]BudgetExpense
	Date              date.Date
	Total             money.Money
}
