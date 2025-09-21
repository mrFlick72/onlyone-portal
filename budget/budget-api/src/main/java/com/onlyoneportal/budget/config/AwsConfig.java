package com.onlyoneportal.budget.config;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import software.amazon.awssdk.auth.credentials.AwsCredentialsProvider;
import software.amazon.awssdk.auth.credentials.EnvironmentVariableCredentialsProvider;
import software.amazon.awssdk.services.dynamodb.DynamoDbClient;
import software.amazon.awssdk.services.dynamodb.DynamoDbClientBuilder;

import java.net.URI;

@Configuration(proxyBeanMethods = false)
public class AwsConfig {

    private Logger logger = LoggerFactory.getLogger(AwsConfig.class);


    @Bean("awsCredentialsProvider")
    public AwsCredentialsProvider iamUserAwsCredentialsProvider() {
        return EnvironmentVariableCredentialsProvider.create();
    }

    @Bean
    public DynamoDbClient dynamoDbClient(@Value("${aws.dynamodb.endpointOverride:}") String endpointOverride,
                                         AwsCredentialsProvider awsCredentialsProvider) {
        DynamoDbClientBuilder dynamoDbClientBuilder = DynamoDbClient.builder()
                .credentialsProvider(awsCredentialsProvider);

        if(!endpointOverride.isEmpty()){
            logger.info("endpointOverride: " + endpointOverride);
            dynamoDbClientBuilder.endpointOverride(URI.create(endpointOverride));
        }
        return dynamoDbClientBuilder.build();

    }
}
