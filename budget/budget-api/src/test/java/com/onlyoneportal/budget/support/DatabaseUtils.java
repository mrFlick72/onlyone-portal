package com.onlyoneportal.budget.support;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.revenue.BudgetRevenue;
import com.onlyoneportal.budget.searchtag.SearchTag;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserName;
import software.amazon.awssdk.auth.credentials.AwsBasicCredentials;
import software.amazon.awssdk.auth.credentials.StaticCredentialsProvider;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.dynamodb.DynamoDbClient;
import software.amazon.awssdk.services.dynamodb.model.*;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.URI;
import java.util.List;
import java.util.stream.Collectors;

public class DatabaseUtils {


    public static final String SEARCH_TAG_TABLE_NAME = "BUDGET_SERVICE_SEARCH_TAGS_STAGING";
    public static final String BUDGET_EXPENSE_TABLE_NAME = "BUDGET_EXPENSE_TABLE_NAME_STAGING";
    public static final String BUDGET_REVENUE_TABLE_NAME = "BUDGET_REVENUE_TABLE_NAME_STAGING";

    private DatabaseUtils() {
    }

    public static DynamoDbClient dynamoClient() {
        return DynamoDbClient.builder()
                .credentialsProvider(
                        StaticCredentialsProvider.create(
                                AwsBasicCredentials.create("ACCESS_KEY_ID", "SECRET_ACCESS_KEY"))
                ).region(Region.US_EAST_1)
                .endpointOverride(URI.create("http://localhost:4566"))
                .build();
    }

    public static List<SearchTag> loadSearchTags() {
        try (InputStream in = DatabaseUtils.class.getClassLoader().getResourceAsStream("search_tag/findAll.csv");
             InputStreamReader inputStreamReader = new InputStreamReader(in);
             BufferedReader br = new BufferedReader(inputStreamReader);
        ) {
            return br.lines().map(line -> {
                        String[] split = line.split(",");
                        return new SearchTag(split[0], split[1]);
                    })
                    .collect(Collectors.toList());

        } catch (IOException e) {
            throw new RuntimeException(e.getMessage(), e);
        }
    }

    public static List<BudgetExpense> loadBudgetExpense() {
        try (InputStream in = DatabaseUtils.class.getClassLoader().getResourceAsStream("budget-expense/find-by-date-range-data-set.csv");
             InputStreamReader inputStreamReader = new InputStreamReader(in);
             BufferedReader br = new BufferedReader(inputStreamReader);
        ) {
            return br.lines().map(line -> {
                        String[] split = line.split(",");
                        return new BudgetExpense(
                                null,
                                new UserName(split[0]),
                                Date.dateFor(split[1]),
                                Money.moneyFor(split[2]),
                                split[3],
                                split[4]
                        );
                    })
                    .collect(Collectors.toList());

        } catch (IOException e) {
            throw new RuntimeException(e.getMessage(), e);
        }
    }

    public static List<BudgetRevenue> loadBudgetRevenue() {
        try (InputStream in = DatabaseUtils.class.getClassLoader().getResourceAsStream("budget_revenue/find-by-date-range-data-set.csv");
             InputStreamReader inputStreamReader = new InputStreamReader(in);
             BufferedReader br = new BufferedReader(inputStreamReader);
        ) {
            return br.lines().map(line -> {
                        String[] split = line.split(",");
                        return new BudgetRevenue(
                                null,
                                split[0],
                                Date.dateFor(split[1]),
                                Money.moneyFor(split[2]),
                                split[3]
                        );
                    })
                    .collect(Collectors.toList());

        } catch (IOException e) {
            throw new RuntimeException(e.getMessage(), e);
        }
    }


    public static void resetDatabase() {
        DynamoDbClient dynamoDbClient = dynamoClient();
        try {
            dynamoDbClient.deleteTable(DeleteTableRequest.builder()
                    .tableName(SEARCH_TAG_TABLE_NAME)
                    .build());
            dynamoDbClient.deleteTable(DeleteTableRequest.builder()
                    .tableName(BUDGET_EXPENSE_TABLE_NAME)
                    .build());
            dynamoDbClient.deleteTable(DeleteTableRequest.builder()
                    .tableName(BUDGET_REVENUE_TABLE_NAME)
                    .build());
        } catch (Exception e) {
        }
        try {
            dynamoDbClient.createTable(CreateTableRequest.builder()
                    .tableName(SEARCH_TAG_TABLE_NAME)
                    .keySchema(KeySchemaElement.builder()
                                    .attributeName("user_name")
                                    .keyType(KeyType.HASH)
                                    .build(),
                            KeySchemaElement.builder()
                                    .attributeName("search_tag_key")
                                    .keyType(KeyType.RANGE)
                                    .build())
                    .attributeDefinitions(AttributeDefinition.builder()
                                    .attributeName("user_name")
                                    .attributeType(ScalarAttributeType.S)
                                    .build(),
                            AttributeDefinition.builder()
                                    .attributeName("search_tag_key")
                                    .attributeType(ScalarAttributeType.S)
                                    .build())
                    .billingMode(BillingMode.PAY_PER_REQUEST)
                    .build());

            dynamoDbClient.createTable(CreateTableRequest.builder()
                    .tableName(BUDGET_EXPENSE_TABLE_NAME)
                    .keySchema(KeySchemaElement.builder()
                                    .attributeName("pk")
                                    .keyType(KeyType.HASH)
                                    .build(),
                            KeySchemaElement.builder()
                                    .attributeName("range_key")
                                    .keyType(KeyType.RANGE)
                                    .build())
                    .attributeDefinitions(AttributeDefinition.builder()
                                    .attributeName("pk")
                                    .attributeType(ScalarAttributeType.S)
                                    .build(),
                            AttributeDefinition.builder()
                                    .attributeName("range_key")
                                    .attributeType(ScalarAttributeType.S)
                                    .build())
                    .billingMode(BillingMode.PAY_PER_REQUEST)
                    .build());


            dynamoDbClient.createTable(CreateTableRequest.builder()
                    .tableName(BUDGET_REVENUE_TABLE_NAME)
                    .keySchema(KeySchemaElement.builder()
                                    .attributeName("pk")
                                    .keyType(KeyType.HASH)
                                    .build(),
                            KeySchemaElement.builder()
                                    .attributeName("range_key")
                                    .keyType(KeyType.RANGE)
                                    .build())
                    .attributeDefinitions(AttributeDefinition.builder()
                                    .attributeName("pk")
                                    .attributeType(ScalarAttributeType.S)
                                    .build(),
                            AttributeDefinition.builder()
                                    .attributeName("range_key")
                                    .attributeType(ScalarAttributeType.S)
                                    .build())
                    .billingMode(BillingMode.PAY_PER_REQUEST)
                    .build());
        } catch (Exception e) {
        }
    }

}
