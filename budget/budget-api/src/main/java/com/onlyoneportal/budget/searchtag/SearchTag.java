package com.onlyoneportal.budget.searchtag;

import java.io.Serializable;
import java.util.Objects;

public record SearchTag(String key, String value) implements Serializable {

    public static final String DEFAULT_KEY = "unknown";

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        SearchTag searchTag = (SearchTag) o;
        return Objects.equals(key, searchTag.key) && Objects.equals(value, searchTag.value);
    }

}
