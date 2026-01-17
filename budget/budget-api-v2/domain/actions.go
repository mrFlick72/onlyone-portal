package domain

import (
	"context"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type CreateBudgetExpense struct {
	repository BudgetExpenseRepository
}

func (action *CreateBudgetExpense) Execute(ctx *context.Context, budgetExpense *BudgetExpense) error {
	user, _ := security.GetCurrentUser(ctx)
	budgetExpense.UserName = *user.UserName
	return action.repository.Save(ctx, budgetExpense)
}

type FindBudgetExpense struct {
	repository BudgetExpenseRepository
}

/*
	public void updateWithoutAttachment(BudgetExpense budgetExpense) {
	    budgetExpenseRepository.findFor(budgetExpense.id())
	            .ifPresent(foundBudgetExpense -> {
	                BudgetExpense updatedBudgetExpense = new BudgetExpense(budgetExpense.id(),
	                        budgetExpense.userName(),
	                        budgetExpense.date(),
	                        budgetExpense.amount(), budgetExpense.note(),
	                        budgetExpense.tag()
	                );

	                budgetExpenseRepository.save(updatedBudgetExpense);
	            });
	}
*/
type UpdateBudgetExpense struct {
	repository BudgetExpenseRepository
}

func (action *UpdateBudgetExpense) Execute(ctx *context.Context, budgetExpense *BudgetExpense) error {
	existingBudgetExpense, err := action.repository.FindFor(ctx, budgetExpense.Id)
	if err != nil {
		return err
	}
	if existingBudgetExpense != nil {
		return action.repository.Save(ctx, budgetExpense)
	}	
	return nil
}

type DeleteBudgetExpense struct {
	repository BudgetExpenseRepository
}
