package com.onlyoneportal.budget.searchtag;

import com.onlyoneportal.budget.infrastructure.dynamodb.DynamoDbAttributeValueFactory;
import com.onlyoneportal.budget.user.UserRepository;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import software.amazon.awssdk.services.dynamodb.DynamoDbClient;

import static com.onlyoneportal.budget.support.DatabaseUtils.*;
import static com.onlyoneportal.budget.support.UserTestFixture.A_USER_NAME;
import static org.mockito.BDDMockito.given;

@ExtendWith(MockitoExtension.class)
class DynamoDBSearchTagRepositoryTest {

    private SearchTagRepository searchTagRepository;

    @Mock
    private UserRepository userRepository;

    @BeforeEach
    public void setUp() {
        DynamoDbClient client = dynamoClient();
        resetDatabase();
        searchTagRepository = new DynamoDBSearchTagRepository(SEARCH_TAG_TABLE_NAME, userRepository, client, new DynamoDbAttributeValueFactory());

        given(userRepository.currentLoggedUserName())
                .willReturn(A_USER_NAME);

        loadSearchTags().forEach(searchTagRepository::save);
    }


    @Test
    public void findAll() {
        Assertions.assertEquals(24, searchTagRepository.findAllSearchTag().size());
    }

    @Test
    public void findSearchTagBy() {
        Assertions.assertEquals(searchTagRepository.findSearchTagBy("super-market"), new SearchTag("super-market", "Spesa"));
    }
    @Test
    public void save() {
        searchTagRepository.save(new SearchTag("test", "Test"));
        SearchTag actual = searchTagRepository.findSearchTagBy("test");
        Assertions.assertEquals(new SearchTag("test", "Test"), actual);
    }
}