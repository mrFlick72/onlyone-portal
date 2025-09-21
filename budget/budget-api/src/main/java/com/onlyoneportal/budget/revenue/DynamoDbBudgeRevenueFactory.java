package com.onlyoneportal.budget.revenue;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.infrastructure.dynamodb.BudgetDynamoDbIdFactory;
import com.onlyoneportal.budget.infrastructure.dynamodb.DynamoDbAttributeValueFactory;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserRepository;
import software.amazon.awssdk.services.dynamodb.model.AttributeValue;

import java.time.LocalDate;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;

public class DynamoDbBudgeRevenueFactory {
    private final BudgetDynamoDbIdFactory<BudgetRevenueId, BudgetRevenue> idFactory;
    private final UserRepository userRepository;
    private final DynamoDbAttributeValueFactory attributeValueFactory;

    public DynamoDbBudgeRevenueFactory(BudgetDynamoDbIdFactory<BudgetRevenueId, BudgetRevenue> idFactory,
                                       UserRepository userRepository, DynamoDbAttributeValueFactory attributeValueFactory) {
        this.idFactory = idFactory;
        this.userRepository = userRepository;
        this.attributeValueFactory = attributeValueFactory;
    }

    public Set<String> primaryKeysFor(Date star, Date end) {
        Set<String> primaryKeys = new TreeSet<>();
        LocalDate current = star.getLocalDate();
        while (!current.isAfter(end.getLocalDate())) {
            primaryKeys.add(idFactory.partitionKeyFrom(new Date(current), userRepository.currentLoggedUserName()));
            current = current.plusYears(1);
        }

        return primaryKeys;
    }

    public Map<String, AttributeValue> yearlyItemKeyConditionFor(String pk, Date star, Date end) {
        HashMap<String, AttributeValue> itemKeyCondition = new HashMap<>();
        itemKeyCondition.put(":pk", attributeValueFactory.stringAttributeFor(pk));
        itemKeyCondition.put(":star", attributeValueFactory.stringAttributeFor(star.isoFormattedDate()));
        itemKeyCondition.put(":end", attributeValueFactory.stringAttributeFor(end.isoFormattedDate()));
        return itemKeyCondition;
    }

    public BudgetRevenue fromDynamoDbToModel(Map<String, AttributeValue> item) {
        return new BudgetRevenue(
                new BudgetRevenueId(item.get("budget_id").s()),
                item.get("user_name").s(),
                Date.isoDateFor(item.get("transaction_date").s()),
                Money.moneyFor(item.get("amount").s()),
                item.get("note").s()
        );
    }


    public Map<String, AttributeValue> putItemPayloadFor(BudgetRevenue budgetExpense) {
        Map<String, AttributeValue> payload = new HashMap<>();

        payload.put("pk", attributeValueFactory.stringAttributeFor(idFactory.partitionKeyFrom(budgetExpense.id())));
        payload.put("range_key", attributeValueFactory.stringAttributeFor(idFactory.rangeKeyFrom(budgetExpense.id())));

        payload.put("budget_id", attributeValueFactory.stringAttributeFor(budgetExpense.id().content()));
        payload.put("user_name", attributeValueFactory.stringAttributeFor(userRepository.currentLoggedUserName().content()));
        payload.put("transaction_date", attributeValueFactory.stringAttributeFor(budgetExpense.registrationDate().isoFormattedDate()));
        payload.put("amount", attributeValueFactory.stringAttributeFor(budgetExpense.amount().stringifyAmount()));
        payload.put("note", attributeValueFactory.stringAttributeFor(budgetExpense.note()));

        return payload;
    }


    public HashMap<String, AttributeValue> keysFor(BudgetRevenueId idBudgetRevenue) {
        HashMap<String, AttributeValue> itemKeyCondition = new HashMap<>();
        itemKeyCondition.put("pk", attributeValueFactory.stringAttributeFor(idFactory.partitionKeyFrom(idBudgetRevenue)));
        itemKeyCondition.put("range_key", attributeValueFactory.stringAttributeFor(idFactory.rangeKeyFrom(idBudgetRevenue)));
        return itemKeyCondition;
    }

}