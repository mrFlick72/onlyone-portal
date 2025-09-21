package com.onlyoneportal.budget.config;

import com.onlyoneportal.budget.expense.repository.BudgetExpenseRepository;
import com.onlyoneportal.budget.expense.repository.DynamoDbBudgetExpenseRepository;
import com.onlyoneportal.budget.revenue.BudgetRevenueRepository;
import com.onlyoneportal.budget.revenue.DynamoDbBudgetRevenueRepository;
import com.onlyoneportal.budget.infrastructure.dynamodb.*;
import com.onlyoneportal.budget.searchtag.CachedSearchTagRepository;
import com.onlyoneportal.budget.searchtag.DynamoDBSearchTagRepository;
import com.onlyoneportal.budget.searchtag.SearchTagRepository;
import com.onlyoneportal.budget.user.SpringSecurityUserRepository;
import com.onlyoneportal.budget.user.UserRepository;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.web.client.RestTemplateBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.web.client.RestTemplate;
import software.amazon.awssdk.services.dynamodb.DynamoDbClient;

@Configuration(proxyBeanMethods = false)
public class RepositoryConfiguration {

    @Bean
    public UserRepository userRepository() {
        return new SpringSecurityUserRepository();
    }

    @Bean
    public BudgetRevenueRepository budgetRevenueRepository(DynamoDbClient dynamoDbClient,
                                                           @Value("${budget-api.dynamo-db.budget-revenue.table-name}") String tableName,
                                                           UserRepository userRepository, SaltGenerator saltGenerator) {

        return new DynamoDbBudgetRevenueRepository(tableName, dynamoDbClient,
                new BudgetRevenueDynamoDbIdFactory(saltGenerator),
                userRepository, new DynamoDbAttributeValueFactory());
    }

    @Bean
    public BudgetExpenseRepository budgetExpenseRepository(DynamoDbClient dynamoDbClient,
                                                           @Value("${budget-api.dynamo-db.budget-expense.table-name}") String tableName,
                                                           UserRepository userRepository, SaltGenerator saltGenerator) {

        return new DynamoDbBudgetExpenseRepository(tableName, dynamoDbClient,
                new BudgetExpenseDynamoDbIdFactory(saltGenerator),
                userRepository, new DynamoDbAttributeValueFactory());
    }

    @Bean
    public SearchTagRepository searchTagRepository(DynamoDbClient dynamoDbClient,
                                                   RedisTemplate redisTemplate,
                                                   @Value("${budget-api.dynamo-db.search-tags.cache-name}") String cacheName,
                                                   @Value("${budget-api.dynamo-db.search-tags.cache-ttl}") Integer cacheTtl,
                                                   @Value("${budget-api.dynamo-db.search-tags.table-name}") String tableName,
                                                   UserRepository userRepository) {
        DynamoDBSearchTagRepository repository = new DynamoDBSearchTagRepository(tableName, userRepository, dynamoDbClient, new DynamoDbAttributeValueFactory());
        return new CachedSearchTagRepository(cacheName, cacheTtl, redisTemplate, userRepository, repository);
    }

    @Bean
    public RestTemplate repositoryServiceRestTemplate() {
        return new RestTemplateBuilder().build();
    }

    @Bean
    public SaltGenerator saltGenerator() {
        return new UUIDSaltGenerator();
    }
}
