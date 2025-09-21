package com.onlyoneportal.budget.revenue;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.infrastructure.dynamodb.BudgetDynamoDbIdFactory;
import com.onlyoneportal.budget.infrastructure.dynamodb.DynamoDbAttributeValueFactory;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserRepository;
import software.amazon.awssdk.services.dynamodb.DynamoDbClient;
import software.amazon.awssdk.services.dynamodb.model.AttributeValue;
import software.amazon.awssdk.services.dynamodb.model.DeleteItemRequest;
import software.amazon.awssdk.services.dynamodb.model.PutItemRequest;
import software.amazon.awssdk.services.dynamodb.model.QueryRequest;

import java.time.LocalDate;
import java.util.*;
import java.util.stream.Collectors;

public class DynamoDbBudgetRevenueRepository implements BudgetRevenueRepository {

    private final String tableName;
    private final DynamoDbClient dynamoClient;
    private final BudgetDynamoDbIdFactory<BudgetRevenueId, BudgetRevenue> idFactory;
    private final UserRepository userRepository;
    private final DynamoDbAttributeValueFactory attributeValueFactory;

    public DynamoDbBudgetRevenueRepository(String tableName,
                                           DynamoDbClient dynamoClient,
                                           BudgetDynamoDbIdFactory idFactory,
                                           UserRepository userRepository,
                                           DynamoDbAttributeValueFactory attributeValueFactory) {
        this.tableName = tableName;
        this.dynamoClient = dynamoClient;
        this.idFactory = idFactory;
        this.userRepository = userRepository;
        this.attributeValueFactory = attributeValueFactory;
    }

    @Override
    public List<BudgetRevenue> findByDateRange(String user, Date star, Date end) {
        Set<String> pks = primaryKeysFor(star, end);
        return pks.stream()
                .flatMap(pk -> {
                    Map<String, AttributeValue> itemKeyCondition = yearlyItemKeyConditionFor(pk, star, end);
                    return dynamoClient.query(
                            QueryRequest.builder()
                                    .tableName(tableName)
                                    .keyConditionExpression("pk =:pk")
                                    .filterExpression("transaction_date >= :star AND transaction_date <= :end")
                                    .expressionAttributeValues(itemKeyCondition)
                                    .build()
                    ).items().stream();
                })
                .map(this::fromDynamoDbToModel)
                .sorted(Comparator.comparing(o -> o.registrationDate().getLocalDate()))
                .collect(Collectors.toList());
    }

    private Set<String> primaryKeysFor(Date star, Date end) {
        Set<String> primaryKeys = new TreeSet<>();
        LocalDate current = star.getLocalDate();
        while (!current.isAfter(end.getLocalDate())) {
            primaryKeys.add(idFactory.partitionKeyFrom(new Date(current), userRepository.currentLoggedUserName()));
            current = current.plusYears(1);
        }

        return primaryKeys;
    }

    private Map<String, AttributeValue> yearlyItemKeyConditionFor(String pk, Date star, Date end) {
        HashMap<String, AttributeValue> itemKeyCondition = new HashMap<>();
        itemKeyCondition.put(":pk", attributeValueFactory.stringAttributeFor(pk));
        itemKeyCondition.put(":star", attributeValueFactory.stringAttributeFor(star.isoFormattedDate()));
        itemKeyCondition.put(":end", attributeValueFactory.stringAttributeFor(end.isoFormattedDate()));
        return itemKeyCondition;
    }

    private BudgetRevenue fromDynamoDbToModel(Map<String, AttributeValue> item) {
        return new BudgetRevenue(
                new BudgetRevenueId(item.get("budget_id").s()),
                item.get("user_name").s(),
                Date.isoDateFor(item.get("transaction_date").s()),
                Money.moneyFor(item.get("amount").s()),
                item.get("note").s()
        );
    }

    @Override
    public BudgetRevenue save(BudgetRevenue budgetRevenue) {
        BudgetRevenue budgetRevenueToSave = new BudgetRevenue(
                idFactory.budgetIdFrom(budgetRevenue),
                budgetRevenue.userName(),
                budgetRevenue.registrationDate(),
                budgetRevenue.amount(),
                budgetRevenue.note()
        );

        dynamoClient.putItem(
                PutItemRequest.builder()
                        .tableName(tableName)
                        .item(putItemPayloadFor(budgetRevenueToSave))
                        .build()
        );

        return budgetRevenueToSave;
    }

    private Map<String, AttributeValue> putItemPayloadFor(BudgetRevenue budgetExpense) {
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

    @Override
    public void update(BudgetRevenue budgetRevenue) {
        save(budgetRevenue);
    }

    @Override
    public void delete(BudgetRevenueId idBudgetRevenue) {
        HashMap<String, AttributeValue> itemKeyCondition = keysFor(idBudgetRevenue);
        dynamoClient.deleteItem(
                DeleteItemRequest.builder()
                        .tableName(tableName)
                        .key(itemKeyCondition)
                        .build()
        );
    }

    private HashMap<String, AttributeValue> keysFor(BudgetRevenueId idBudgetRevenue) {
        HashMap<String, AttributeValue> itemKeyCondition = new HashMap<>();
        itemKeyCondition.put("pk", attributeValueFactory.stringAttributeFor(idFactory.partitionKeyFrom(idBudgetRevenue)));
        itemKeyCondition.put("range_key", attributeValueFactory.stringAttributeFor(idFactory.rangeKeyFrom(idBudgetRevenue)));
        return itemKeyCondition;
    }

}
