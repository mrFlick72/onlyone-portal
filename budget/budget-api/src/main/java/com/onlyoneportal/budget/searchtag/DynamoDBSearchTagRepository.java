package com.onlyoneportal.budget.searchtag;

import com.onlyoneportal.budget.infrastructure.dynamodb.DynamoDbAttributeValueFactory;
import com.onlyoneportal.budget.user.UserRepository;
import software.amazon.awssdk.services.dynamodb.DynamoDbClient;
import software.amazon.awssdk.services.dynamodb.model.AttributeValue;
import software.amazon.awssdk.services.dynamodb.model.PutItemRequest;
import software.amazon.awssdk.services.dynamodb.model.QueryRequest;
import software.amazon.awssdk.services.dynamodb.model.QueryResponse;

import java.util.*;
import java.util.stream.Collectors;

public class DynamoDBSearchTagRepository implements SearchTagRepository {
    private final String tableName;
    private final UserRepository userRepository;
    private final DynamoDbClient client;
    private final DynamoDbAttributeValueFactory attributeValueFactory;

    public DynamoDBSearchTagRepository(String tableName,
                                       UserRepository userRepository,
                                       DynamoDbClient client,
                                       DynamoDbAttributeValueFactory attributeValueFactory) {
        this.tableName = tableName;
        this.userRepository = userRepository;
        this.client = client;
        this.attributeValueFactory = attributeValueFactory;
    }

    @Override
    public SearchTag findSearchTagBy(String key) {
        HashMap<String, AttributeValue> itemKeyCondition = itemKeyConditionFor(key);
        QueryResponse query = client
                .query(
                        QueryRequest.builder()
                                .tableName(this.tableName)
                                .keyConditionExpression("user_name =:user_name AND search_tag_key =:search_tag_key")
                                .expressionAttributeValues(itemKeyCondition)
                                .build()
                );

        return query.items().stream()
                .map(this::fromDynamoDbToModel)
                .findFirst()
                .orElse(null);
    }

    @Override
    public List<SearchTag> findAllSearchTag() {
        QueryResponse query = client.query(
                QueryRequest.builder()
                        .tableName(tableName)
                        .keyConditionExpression("user_name =:user_name")
                        .expressionAttributeValues(findAllItemKeyCondition())
                        .build()
        );
        return query.items().stream()
                .map(this::fromDynamoDbToModel)
                .collect(Collectors.toList());
    }

    @Override
    public void save(SearchTag searchTag) {
        client.putItem(
                PutItemRequest.builder()
                        .tableName(this.tableName)
                        .item(putItemPayloadFor(searchTag))
                        .build()
        );
    }

    private HashMap<String, AttributeValue> itemKeyConditionFor(String key) {
        HashMap<String, AttributeValue> itemKeyCondition = new HashMap<>();
        itemKeyCondition.put(":user_name", attributeValueFactory.stringAttributeFor(userRepository.currentLoggedUserName().content()));
        itemKeyCondition.put(":search_tag_key", attributeValueFactory.stringAttributeFor(key));
        return itemKeyCondition;
    }

    private HashMap<String, AttributeValue> findAllItemKeyCondition() {
        HashMap<String, AttributeValue> itemKeyCondition = new HashMap<>();
        itemKeyCondition.put(":user_name", attributeValueFactory.stringAttributeFor(userRepository.currentLoggedUserName().content()));
        return itemKeyCondition;
    }

    private SearchTag fromDynamoDbToModel(Map<String, AttributeValue> attributes) {
        return new SearchTag(attributes.get("search_tag_key").s(), attributes.get("search_tag_value").s());
    }

    private HashMap<String, AttributeValue> putItemPayloadFor(SearchTag searchTag) {
        HashMap<String, AttributeValue> attributes = new HashMap<>();
        attributes.put("user_name", attributeValueFactory.stringAttributeFor(userRepository.currentLoggedUserName().content()));
        attributes.put("search_tag_key", attributeValueFactory.stringAttributeFor(keyFrom(searchTag)));
        attributes.put("search_tag_value", attributeValueFactory.stringAttributeFor(searchTag.value()));
        return attributes;
    }

    private static String keyFrom(SearchTag searchTag) {
        return Optional.ofNullable(searchTag.key())
                .map(id -> id.isEmpty() ? UUID.randomUUID().toString() : id)
                .orElse(UUID.randomUUID().toString());
    }
}