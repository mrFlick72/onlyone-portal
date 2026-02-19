package web

import (
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
)

// todo to be tested
func RepresentationModelToDomainModel(budgetExpenseRepresentation BudgetExpenseRepresentation) *expense.BudgetExpense {
	date, _ := date.DateFor(budgetExpenseRepresentation.Date)
	money, _ := money.MoneyFor(budgetExpenseRepresentation.Amount)
	return &expense.BudgetExpense{
		Id:       expense.BudgetExpenseId(budgetExpenseRepresentation.Id),
		UserName: "", // This will be set in the service layer using the current logged user
		Date:     *date,
		Amount:   money,
		Note:     budgetExpenseRepresentation.Note,
		Tag:      budgetExpenseRepresentation.TagKey,
	}
}

//     public SpentBudgetRepresentation domainToRepresentationModel(SpentBudget spentBudget) {
//         return new SpentBudgetRepresentation(spentBudget.total().stringifyAmount(),
//                 spentBudget.dailyBudgetExpenseList().stream()
//                         .map(dailyBudgetExpense ->
//                                 new DailyBudgetExpenseRepresentation(budgetExpenseRepresentationList(dailyBudgetExpense),
//                                         dailyBudgetExpense.date().formattedDate(),
//                                         dailyBudgetExpense.total().stringifyAmount())).collect(toList()),
//                 spentBudget.totalForSearchTags().entrySet().stream()
//                         .map(total -> new TotalBySearchTagDetail(total.getKey().key(),
//                                 total.getKey().value(),
//                                 total.getValue().stringifyAmount()))
//                         .collect(toList()));
//     }


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

	result := &SpentBudgetRepresentation{
		Total:                            spentBudget.Total().StringifyAmount(),
		DailyBudgetExpenseRepresentation: dailyBudgetExpenseRepresentations,
		TotalDetailList:                  nil,
	}
	return result
}

    // private List<BudgetExpenseRepresentation> budgetExpenseRepresentationList(DailyBudgetExpense dailyBudgetExpense) {
    //     return dailyBudgetExpense.budgetExpenseList().stream()
    //             .map(budgetExpenseConverter::domainToRepresentationModel)
    //             .collect(toList());
    // }

func budgetExpenseRepresentationList(dailyBudgetExpense *expense.DailyBudgetExpense) []BudgetExpenseRepresentation {
	budgetExpenseRepresentationList := make([]BudgetExpenseRepresentation, 0)
	for _, budgetExpense := range *dailyBudgetExpense.BudgetExpenseList {
		budgetExpenseRepresentationList = append(budgetExpenseRepresentationList, BudgetExpenseRepresentation{
			Id:       string(budgetExpense.Id),
			Date:     budgetExpense.Date.GetFormattedDate(),
			Amount:   budgetExpense.Amount.StringifyAmount(),
			Note:     budgetExpense.Note,
			TagKey:   budgetExpense.Tag,
			TagValue: "", // This will be set in the service layer using the current logged user
		})
	}
	return budgetExpenseRepresentationList
}	
