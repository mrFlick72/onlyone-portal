package expense

import (
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	tagRep "github.com/mrflick72/budget/budget-api/web/tags"
)

func RepresentationModelToDomainModel(budgetExpenseRepresentation BudgetExpenseRepresentation) (*expense.BudgetExpense, error) {
	date, err := date.DateFor(budgetExpenseRepresentation.Date)
	if err != nil {
		return nil, err
	}
	money, err := money.MoneyFor(budgetExpenseRepresentation.Amount)
	if err != nil {
		return nil, err
	}
	return &expense.BudgetExpense{
		Id:       expense.BudgetExpenseId(budgetExpenseRepresentation.Id),
		UserName: "", // This will be set in the service layer using the current logged user
		Date:     *date,
		Amount:   money,
		Note:     budgetExpenseRepresentation.Note,
		Tag:      []tags.SearchTag{{Key: budgetExpenseRepresentation.Tag.Key, Value: budgetExpenseRepresentation.Tag.Value}},
	}, nil
}


func SpentBudgetDomainToRepresentationModel(spentBudget *expense.SpentBudget) *SpentBudgetRepresentation {

	dailyBudgetExpenseRepresentations := make([]DailyBudgetExpenseRepresentation, 0)
	for _, dailyBudgetExpense := range spentBudget.DailyBudgetExpenseList() {
		dailyBudgetExpenseRepresentation := DailyBudgetExpenseRepresentation{
			BudgetExpenseRepresentationList: budgetExpenseRepresentationList(&dailyBudgetExpense),
			Date:                            dailyBudgetExpense.Date.GetFormattedDate(),
			Total:                           dailyBudgetExpense.Total.StringifyAmount(),
		}
		dailyBudgetExpenseRepresentations = append(dailyBudgetExpenseRepresentations, dailyBudgetExpenseRepresentation)
	
	}

	totalBySearchTagDetails := make([]TotalBySearchTagDetail, 0)
	for totalBySearchTag, amount := range spentBudget.TotalForSearchTags() {
		totalBySearchTagDetails = append(totalBySearchTagDetails, TotalBySearchTagDetail{
			SearchTagKey:   totalBySearchTag.Key,
			SearchTagValue: totalBySearchTag.Value,
			Total:          amount.StringifyAmount(),
		})

	}

	result := &SpentBudgetRepresentation{
		Total:                            spentBudget.Total().StringifyAmount(),
		DailyBudgetExpenseRepresentation: dailyBudgetExpenseRepresentations,
		TotalDetailList:                  totalBySearchTagDetails,
	}
	return result
}


func budgetExpenseRepresentationList(dailyBudgetExpense *expense.DailyBudgetExpense) []BudgetExpenseRepresentation {
	budgetExpenseRepresentationList := make([]BudgetExpenseRepresentation, 0)
	for _, budgetExpense := range *dailyBudgetExpense.BudgetExpenseList {
		budgetExpenseRepresentationList = append(budgetExpenseRepresentationList, BudgetExpenseRepresentation{
			Id:       string(budgetExpense.Id),
			Date:     budgetExpense.Date.GetFormattedDate(),
			Amount:   budgetExpense.Amount.StringifyAmount(),
			Note:     budgetExpense.Note,
			Tag: tagRep.SearchTagRepresentation{
				Key:   budgetExpense.Tag[0].Key,
				Value: budgetExpense.Tag[0].Value,
			},
		})
	}
	return budgetExpenseRepresentationList
}
