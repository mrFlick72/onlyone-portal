package com.onlyoneportal.budget.expense.converter;

import com.onlyoneportal.budget.expense.endpoint.BudgetSearchCriteriaRepresentation;
import org.springframework.core.convert.converter.Converter;

import java.util.Arrays;
import java.util.HashMap;
import java.util.Optional;
import java.util.function.Function;

import static java.util.Arrays.asList;

public class StringToBudgetSearchCriteriaRepresentationConverter implements Converter<String, BudgetSearchCriteriaRepresentation> {

    @Override
    public BudgetSearchCriteriaRepresentation convert(String representation) {

        Function<String, BudgetSearchCriteriaRepresentation> converter =
                (r) -> Arrays.stream(r.split(";"))
                        .map(searchCriteria -> {
                            HashMap<String, String> hashMap = new HashMap<>();
                            String[] split = searchCriteria.split("=");
                            if (split.length == 2) {
                                hashMap.put(split[0], split[1]);
                            }
                            return hashMap;
                        })
                        .reduce((hashMap, hashMap2) -> {
                            hashMap.putAll(hashMap2);
                            return hashMap;
                        }).map(hashMap -> new BudgetSearchCriteriaRepresentation(Integer.parseInt(hashMap.getOrDefault("month", "1")),
                                Integer.parseInt(hashMap.getOrDefault("year", "1970")),
                                Optional.ofNullable(hashMap.get("searchTag"))
                                        .map(tag -> asList(tag.split(",")))
                                        .orElse(asList())))
                        .orElse(BudgetSearchCriteriaRepresentation.empty());


        return Optional.ofNullable(representation)
                .map(converter)
                .orElse(BudgetSearchCriteriaRepresentation.empty());
    }
}