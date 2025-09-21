package com.onlyoneportal.budget.searchtag;

import com.onlyoneportal.budget.user.UserName;
import com.onlyoneportal.budget.user.UserRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.redis.connection.lettuce.LettuceConnectionFactory;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.data.redis.serializer.JdkSerializationRedisSerializer;
import org.springframework.data.redis.serializer.StringRedisSerializer;

import java.util.List;

import static com.onlyoneportal.budget.infrastructure.strings.Sha256Utils.sha256For;
import static java.util.Arrays.asList;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.BDDMockito.given;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CachedSearchTagRepositoryTest {
    String cacheKey = "search_tag_cache_anemail@test.com_key";
    String findAllCacheKey = "search_tag_cache_anemail@test.com";

    public static final UserName USER_NAME = new UserName("anemail@test.com");
    @Mock
    private SearchTagRepository repository;

    @Mock
    private UserRepository userRepository;

    private RedisTemplate<String, SearchTag> redisTemplate;

    private CachedSearchTagRepository underTest;

    @BeforeEach
    void setup() {
        reset(repository);
        reset(userRepository);

        setUpRedisTemplate();
        underTest = new CachedSearchTagRepository("search_tag_cache", 1000, redisTemplate, userRepository, repository);

        given(userRepository.currentLoggedUserName())
                .willReturn(USER_NAME);

        redisTemplate.opsForHash().delete(cacheKey,sha256For(cacheKey));
        redisTemplate.opsForHash().delete(findAllCacheKey,sha256For(findAllCacheKey));
    }

    private void setUpRedisTemplate() {
        LettuceConnectionFactory connectionFactory = new LettuceConnectionFactory("localhost", 6379);
        connectionFactory.afterPropertiesSet();
        redisTemplate = new RedisTemplate<>();
        redisTemplate.setConnectionFactory(connectionFactory);
        redisTemplate.setKeySerializer(new StringRedisSerializer());
        redisTemplate.setHashKeySerializer(new StringRedisSerializer());
        redisTemplate.setValueSerializer(new JdkSerializationRedisSerializer());
        redisTemplate.afterPropertiesSet();
    }


    @Test
    void whenASearchTagIsStored() {
        SearchTag searchTag = new SearchTag("key", "value");

        underTest.save(searchTag);

        SearchTag actual = (SearchTag) redisTemplate.opsForHash().get(cacheKey, sha256For(cacheKey));

        assertEquals(searchTag, actual);
        verify(repository).save(searchTag);
    }

    @Test
    void whenASearchTagIsRetrievedFromTheCache() {
        SearchTag searchTag = new SearchTag("key", "value");

        redisTemplate.opsForHash().put(cacheKey, sha256For(cacheKey), searchTag);

        SearchTag actual = underTest.findSearchTagBy("key");

        SearchTag expectedFromTheCache = (SearchTag) redisTemplate.opsForHash().get(cacheKey, sha256For(cacheKey));

        assertEquals(expectedFromTheCache, actual);
        verifyNoMoreInteractions(repository);
    }

    @Test
    void whenASearchTagIsRetrievedFromTheDatabase() {
        String searchTagKey = "key";
        SearchTag searchTag = new SearchTag(searchTagKey, "value");

        given(repository.findSearchTagBy(searchTagKey))
                .willReturn(searchTag);

        SearchTag actual = underTest.findSearchTagBy(searchTagKey);

        SearchTag expectedFromTheCache = (SearchTag) redisTemplate.opsForHash().get(cacheKey, sha256For(cacheKey));

        assertEquals(expectedFromTheCache, actual);
    }

    @Test
    void whenFindAllFromTheCache() {
        SearchTag searchTag = new SearchTag("key", "value");

        redisTemplate.opsForHash().put(findAllCacheKey, sha256For(findAllCacheKey), asList(searchTag));

        List<SearchTag> actual = underTest.findAllSearchTag();
        List<SearchTag> expectedFromTheCache = (List<SearchTag>) redisTemplate.opsForHash().get(findAllCacheKey, sha256For(findAllCacheKey));

        assertNotNull(actual);
        assertNotNull(expectedFromTheCache);
        assertEquals(expectedFromTheCache, actual);
        verifyNoMoreInteractions(repository);
    }

    @Test
    void whenFindAllIsRetrievedFromTheDatabase() {
        SearchTag searchTag = new SearchTag("key", "value");

        given(repository.findAllSearchTag())
                .willReturn(asList(searchTag));

        List<SearchTag> actual = underTest.findAllSearchTag();

        List<SearchTag> expectedFromTheCache = (List<SearchTag>) redisTemplate.opsForHash().get(findAllCacheKey, sha256For(findAllCacheKey));

        assertNotNull(actual);
        assertNotNull(expectedFromTheCache);
        assertEquals(expectedFromTheCache, actual);
    }
    @Test
    void whenASaveDoEvictFindAllCache() {
        SearchTag searchTag = new SearchTag("key", "value");
        SearchTag anotherSearchTag = new SearchTag("key2", "value2");

        given(repository.findAllSearchTag())
                .willReturn(asList(searchTag, anotherSearchTag));

        List<SearchTag> actual = underTest.findAllSearchTag();

        List<SearchTag> expectedFromTheCache = (List<SearchTag>) redisTemplate.opsForHash().get(findAllCacheKey, sha256For(findAllCacheKey));

        assertNotNull(actual);
        assertNotNull(expectedFromTheCache);
        assertEquals(expectedFromTheCache, actual);

        underTest.save(searchTag);
        expectedFromTheCache = (List<SearchTag>) redisTemplate.opsForHash().get(findAllCacheKey, sha256For(findAllCacheKey));
        assertNull(expectedFromTheCache);
    } 

}