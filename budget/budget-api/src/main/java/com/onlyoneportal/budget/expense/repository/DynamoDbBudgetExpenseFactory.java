package com.onlyoneportal.budget.expense.repository;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.model.BudgetExpenseId;
import com.onlyoneportal.budget.infrastructure.dynamodb.BudgetDynamoDbIdFactory;
import com.onlyoneportal.budget.infrastructure.dynamodb.DynamoDbAttributeValueFactory;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserName;
import com.onlyoneportal.budget.user.UserRepository;
import software.amazon.awssdk.services.dynamodb.model.AttributeValue;

import java.time.LocalDate;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;

public class DynamoDbBudgetExpenseFactory  {
    private final BudgetDynamoDbIdFactory<BudgetExpenseId, BudgetExpense> idFactory;
    private final UserRepository userRepository;
    private final DynamoDbAttributeValueFactory attributeValueFactory;

    public DynamoDbBudgetExpenseFactory(BudgetDynamoDbIdFactory idFactory,
                                        UserRepository userRepository,
                                        DynamoDbAttributeValueFactory attributeValueFactory) {
        this.idFactory = idFactory;
        this.userRepository = userRepository;
        this.attributeValueFactory = attributeValueFactory;
    }


    public HashMap<String, AttributeValue> itemKeyConditionFor(BudgetExpenseId budgetExpenseId) {
        HashMap<String, AttributeValue> itemKeyCondition = new HashMap<>();
        itemKeyCondition.put(":pk", attributeValueFactory.stringAttributeFor(idFactory.partitionKeyFrom(budgetExpenseId)));
        itemKeyCondition.put(":range_key", attributeValueFactory.stringAttributeFor(idFactory.rangeKeyFrom(budgetExpenseId)));
        return itemKeyCondition;
    }

    public Set<String> primaryKeysFor(Date star, Date end) {
        Set<String> primaryKeys = new TreeSet<>();
        LocalDate current = star.getLocalDate();
        while (!end.getLocalDate().equals(current)) {
            primaryKeys.add(idFactory.partitionKeyFrom(new Date(current), userRepository.currentLoggedUserName()));
            current= current.plusDays(1);
        }

        return primaryKeys;
    }

    public Map<String, AttributeValue> monthlyItemKeyConditionFor(String pk, Date star, Date end) {
        HashMap<String, AttributeValue> itemKeyCondition = new HashMap<>();
        itemKeyCondition.put(":pk", attributeValueFactory.stringAttributeFor(pk));
        itemKeyCondition.put(":star", attributeValueFactory.stringAttributeFor(star.isoFormattedDate()));
        itemKeyCondition.put(":end", attributeValueFactory.stringAttributeFor(end.isoFormattedDate()));
        return itemKeyCondition;
    }


    public HashMap<String, AttributeValue> keysFor(BudgetExpenseId idBudgetExpense) {
        HashMap<String, AttributeValue> itemKeyCondition = new HashMap<>();
        itemKeyCondition.put("pk", attributeValueFactory.stringAttributeFor(idFactory.partitionKeyFrom(idBudgetExpense)));
        itemKeyCondition.put("range_key", attributeValueFactory.stringAttributeFor(idFactory.rangeKeyFrom(idBudgetExpense)));
        return itemKeyCondition;
    }

    public BudgetExpense fromDynamoDbToModel(Map<String, AttributeValue> item) {
        return new BudgetExpense(
                new BudgetExpenseId(item.get("budget_id").s()),
                new UserName(item.get("user_name").s()),
                Date.isoDateFor(item.get("transaction_date").s()),
                Money.moneyFor(item.get("amount").s()),
                item.get("note").s(),
                item.get("tag").s()
        );
    }

    private Map<String, AttributeValue> putItemPayloadFor(BudgetExpense budgetExpense) {
        Map<String, AttributeValue> payload = new HashMap<>();

        payload.put("pk", attributeValueFactory.stringAttributeFor(idFactory.partitionKeyFrom(budgetExpense.id())));
        payload.put("range_key", attributeValueFactory.stringAttributeFor(idFactory.rangeKeyFrom(budgetExpense.id())));

        payload.put("budget_id", attributeValueFactory.stringAttributeFor(budgetExpense.id().content()));
        payload.put("user_name", attributeValueFactory.stringAttributeFor(userRepository.currentLoggedUserName().content()));
        payload.put("transaction_date", attributeValueFactory.stringAttributeFor(budgetExpense.date().isoFormattedDate()));
        payload.put("amount", attributeValueFactory.stringAttributeFor(budgetExpense.amount().stringifyAmount()));
        payload.put("note", attributeValueFactory.stringAttributeFor(budgetExpense.note()));
        payload.put("tag", attributeValueFactory.stringAttributeFor(budgetExpense.tag()));

        return payload;
    }

}