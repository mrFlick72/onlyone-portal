package com.onlyoneportal.budget.expense.repository;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.model.BudgetExpenseId;
import com.onlyoneportal.budget.infrastructure.dynamodb.BudgetDynamoDbIdFactory;
import com.onlyoneportal.budget.infrastructure.dynamodb.DynamoDbAttributeValueFactory;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserName;
import com.onlyoneportal.budget.user.UserRepository;
import software.amazon.awssdk.services.dynamodb.DynamoDbClient;
import software.amazon.awssdk.services.dynamodb.model.AttributeValue;
import software.amazon.awssdk.services.dynamodb.model.DeleteItemRequest;
import software.amazon.awssdk.services.dynamodb.model.PutItemRequest;
import software.amazon.awssdk.services.dynamodb.model.QueryRequest;

import java.time.LocalDate;
import java.util.*;
import java.util.stream.Collectors;

public class DynamoDbBudgetExpenseRepository implements BudgetExpenseRepository {
    private final String tableName;
    private final DynamoDbClient dynamoClient;
    private final BudgetDynamoDbIdFactory<BudgetExpenseId, BudgetExpense> idFactory;
    private final UserRepository userRepository;
    private final DynamoDbAttributeValueFactory attributeValueFactory;

    public DynamoDbBudgetExpenseRepository(String tableName,
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
    public Optional<BudgetExpense> findFor(BudgetExpenseId budgetExpenseId) {
        HashMap<String, AttributeValue> itemKeyCondition = itemKeyConditionFor(budgetExpenseId);
        return dynamoClient.query(
                        QueryRequest.builder()
                                .tableName(tableName)
                                .keyConditionExpression("pk =:pk AND range_key =:range_key")
                                .expressionAttributeValues(itemKeyCondition)
                                .build()
                ).items()
                .stream()
                .findFirst()
                .map(this::fromDynamoDbToModel);
    }

    private HashMap<String, AttributeValue> itemKeyConditionFor(BudgetExpenseId budgetExpenseId) {
        HashMap<String, AttributeValue> itemKeyCondition = new HashMap<>();
        itemKeyCondition.put(":pk", attributeValueFactory.stringAttributeFor(idFactory.partitionKeyFrom(budgetExpenseId)));
        itemKeyCondition.put(":range_key", attributeValueFactory.stringAttributeFor(idFactory.rangeKeyFrom(budgetExpenseId)));
        return itemKeyCondition;
    }


    @Override
    public List<BudgetExpense> findByDateRange(UserName userName, Date star, Date end, String... searchTags) {
        Set<String> pks = primaryKeysFor(star, end);
        List<String> searchTagsList = Arrays.stream(searchTags).toList();
        return pks.stream()
                .flatMap(pk -> {
                    Map<String, AttributeValue> itemKeyCondition = monthlyItemKeyConditionFor(pk,star, end);
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
                .filter(budgetExpense -> haveToBeIncluded(searchTagsList, budgetExpense))
                .collect(Collectors.toList());
    }

    private static boolean haveToBeIncluded(List<String> searchTagsList, BudgetExpense budgetExpense) {
        return searchTagsList.isEmpty() || searchTagsList.contains(budgetExpense.tag());

    }

    private Set<String> primaryKeysFor(Date star, Date end) {
        Set<String> primaryKeys = new TreeSet<>();
        LocalDate current = star.getLocalDate();
        while (!end.getLocalDate().equals(current)) {
            primaryKeys.add(idFactory.partitionKeyFrom(new Date(current), userRepository.currentLoggedUserName()));
            current= current.plusDays(1);
        }

        return primaryKeys;
    }

    private Map<String, AttributeValue> monthlyItemKeyConditionFor(String pk, Date star, Date end) {
        HashMap<String, AttributeValue> itemKeyCondition = new HashMap<>();
        itemKeyCondition.put(":pk", attributeValueFactory.stringAttributeFor(pk));
        itemKeyCondition.put(":star", attributeValueFactory.stringAttributeFor(star.isoFormattedDate()));
        itemKeyCondition.put(":end", attributeValueFactory.stringAttributeFor(end.isoFormattedDate()));
        return itemKeyCondition;
    }

    @Override
    public BudgetExpense save(BudgetExpense budgetExpense) {
        BudgetExpense budgetExpenseToSave = new BudgetExpense(
                idFactory.budgetIdFrom(budgetExpense),
                budgetExpense.userName(),
                budgetExpense.date(),
                budgetExpense.amount(),
                budgetExpense.note(),
                budgetExpense.tag()
        );
        dynamoClient.putItem(
                PutItemRequest.builder()
                        .tableName(tableName)
                        .item(putItemPayloadFor(budgetExpenseToSave))
                        .build()
        );

        return budgetExpenseToSave;
    }

    @Override
    public void delete(BudgetExpenseId idBudgetExpense) {
        HashMap<String, AttributeValue> itemKeyCondition = keysFor(idBudgetExpense);
        dynamoClient.deleteItem(
                DeleteItemRequest.builder()
                        .tableName(tableName)
                        .key(itemKeyCondition)
                        .build()
        );
    }

    private HashMap<String, AttributeValue> keysFor(BudgetExpenseId idBudgetExpense) {
        HashMap<String, AttributeValue> itemKeyCondition = new HashMap<>();
        itemKeyCondition.put("pk", attributeValueFactory.stringAttributeFor(idFactory.partitionKeyFrom(idBudgetExpense)));
        itemKeyCondition.put("range_key", attributeValueFactory.stringAttributeFor(idFactory.rangeKeyFrom(idBudgetExpense)));
        return itemKeyCondition;
    }

    private BudgetExpense fromDynamoDbToModel(Map<String, AttributeValue> item) {
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