package com.onlyoneportal.budget.searchtag;


import java.util.List;

public interface SearchTagRepository {

    SearchTag findSearchTagBy(String key);
    List<SearchTag> findAllSearchTag();
    void save(SearchTag searchTag);
}
