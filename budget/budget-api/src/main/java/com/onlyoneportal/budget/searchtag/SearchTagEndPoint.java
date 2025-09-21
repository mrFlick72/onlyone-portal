package com.onlyoneportal.budget.searchtag;

import com.onlyoneportal.budget.infrastructure.dynamodb.SaltGenerator;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
public class SearchTagEndPoint {

    private final SearchTagRepository searchTagRepository;

    public SearchTagEndPoint(SearchTagRepository searchTagRepository, SaltGenerator saltGenerator) {
        this.searchTagRepository = searchTagRepository;
    }

    @GetMapping("/budget-expense/search-tag")
    public ResponseEntity findAllSearchTag() {
        return ResponseEntity.ok(searchTagRepository.findAllSearchTag());
    }

    @PutMapping("/budget-expense/search-tag")
    public ResponseEntity newSearchTag(@RequestBody SearchTag searchTag) {
        searchTagRepository.save(searchTag);
        return ResponseEntity.noContent().build();
    }

}