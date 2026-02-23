package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	aws_dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/mrflick72/budget/budget-api/adapter/budget/dynamodb"
	"github.com/mrflick72/budget/budget-api/adapter/tags/rest"
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
)

func NewBudgetExpenseRepository() expense.BudgetExpenseRepository {
	return &dynamodb.DynamoDbBudgetExpenseRepository{
		TableName:               "budget-expenses",
		Client:                  aws_dynamodb.NewFromConfig(aws.Config{}), // You should provide your AWS configuration here
		BudgetExpenseIdProvider: &dynamodb.DynamoDbBudgetExpenseIdProvider{},
	}	
}

func NewBudgetExpenseActionsFacade() expense.BudgetExpenseActions {
	budgetExpenseRepository := NewBudgetExpenseRepository()
	searchTagRepository := rest.NewSearchTagRepository()
	createBudgetExpense := &expense.CreateBudgetExpense{
		Repository: budgetExpenseRepository,
	}
	updateBudgetExpense := &expense.UpdateBudgetExpense{
		Repository: budgetExpenseRepository,
	}
	findSpentBudget := &expense.FindSpentBudget{
		BudgetExpenseRepository: budgetExpenseRepository,
		SearchTagRepository:     searchTagRepository,
	}
	deleteBudgetExpense := &expense.DeleteBudgetExpense{
		Repository: budgetExpenseRepository,
	}

	return &expense.BudgetExpenseActionsFacade{
		CreateBudgetExpenseAction: createBudgetExpense,
		UpdateBudgetExpenseAction: updateBudgetExpense,
		FindSpentBudgetAction:     findSpentBudget,
		DeleteBudgetExpenseAction: deleteBudgetExpense,
	}
}

