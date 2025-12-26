package com.onlyoneportal.budget.searchtag;

import java.util.List;

import org.springframework.core.ParameterizedTypeReference;
import org.springframework.http.HttpMethod;
import org.springframework.http.RequestEntity;
import org.springframework.web.client.RestTemplate;

public class RestTagsRepository implements SearchTagRepository {

    private final RestTemplate restTemplate;
    private final String baseUrl;

    public RestTagsRepository(RestTemplate restTemplate, String baseUrl) {
        this.restTemplate = restTemplate;
        this.baseUrl = baseUrl;
    }

    @Override
    public SearchTag findSearchTagBy(String key) {
        return findAllSearchTag().stream()
                .filter(tag -> tag.key().equals(key))
                .findFirst()
                .orElse(null);
    }

    @Override
    public List<SearchTag> findAllSearchTag() {

        return restTemplate.exchange(baseUrl + "/api/tags",
                HttpMethod.GET, RequestEntity.EMPTY, new ParameterizedTypeReference<List<SearchTag>>() {
                }).getBody();
    }

    @Override
    public void save(SearchTag searchTag) {
        restTemplate.put(baseUrl + "/api/tags", searchTag);
    }

}
