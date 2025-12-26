package com.onlyoneportal.budget.searchtag;

import java.util.List;

import org.springframework.core.ParameterizedTypeReference;
import org.springframework.web.client.RestTemplate;

import com.fasterxml.jackson.core.type.TypeReference;

public class RestTagsRepository implements SearchTagRepository {

    private final RestTemplate restTemplate;
    private final String baseUrl;

    public RestTagsRepository(RestTemplate restTemplate, String baseUrl) {
        this.restTemplate = restTemplate;
        this.baseUrl = baseUrl;
    }

    @Override
    public SearchTag findSearchTagBy(String key) {
        restTemplate.getForObject(baseUrl + "/api/tags/" + key, SearchTag.class);
        return null;
    }

    @Override
    public List<SearchTag> findAllSearchTag() {
    restTemplate.getForEntity(baseUrl + "/api/tags", new ParameterizedTypeReference<List<SearchTag>>(){})
            return null;
    }

    @Override
    public void save(SearchTag searchTag) {
        throw new UnsupportedOperationException("Not supported yet.");
    }

}
