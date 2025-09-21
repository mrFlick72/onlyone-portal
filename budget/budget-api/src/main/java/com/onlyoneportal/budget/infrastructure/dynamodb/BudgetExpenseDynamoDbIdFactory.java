package com.onlyoneportal.budget.infrastructure.dynamodb;

import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.model.BudgetExpenseId;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserName;

import java.util.Base64;
import java.util.Optional;

public class BudgetExpenseDynamoDbIdFactory implements BudgetDynamoDbIdFactory<BudgetExpenseId, BudgetExpense> {

    private final static Base64.Encoder ENCODER = Base64.getEncoder();

    private final SaltGenerator saltGenerator;

    public BudgetExpenseDynamoDbIdFactory(SaltGenerator saltGenerator) {
        this.saltGenerator = saltGenerator;
    }


    @Override
    public BudgetExpenseId budgetIdFrom(BudgetExpense budgetExpense) {
        return Optional.ofNullable(budgetExpense.id())
                .orElseGet(() -> new BudgetExpenseId(String.format("%s-%s", partitionKeyFrom(budgetExpense.date(), budgetExpense.userName()), rangeKeyFrom(budgetExpense))));
    }

    @Override
    public String partitionKeyFrom(Date date, UserName userName) {
        int budgetExpenseYear = date.getLocalDate().getYear();
        int budgetExpenseMonth = date.getLocalDate().getMonthValue();
        String budgetExpenseUser = userName.content();

        String partitionKey = String.format("%s_%s_%s", budgetExpenseYear, budgetExpenseMonth, budgetExpenseUser);

        return ENCODER.encodeToString(partitionKey.getBytes());
    }

    @Override
    public String partitionKeyFrom(BudgetExpenseId id) {
        return id.content().split("-")[0];
    }

    @Override
    public String rangeKeyFrom(BudgetExpenseId id) {
        return id.content().split("-")[1];
    }


    private String rangeKeyFrom(BudgetExpense budgetExpense) {
        int dayOfMonth = budgetExpense.date().getLocalDate().getDayOfMonth();
        String salt = saltGenerator.newSalt();

        String rangeKey = String.format("%s_%s", dayOfMonth, salt);

        return ENCODER.encodeToString(rangeKey.getBytes());
    }

}