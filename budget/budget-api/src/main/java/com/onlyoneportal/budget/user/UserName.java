package com.onlyoneportal.budget.user;

import java.util.Objects;

public record UserName(String content) {

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        UserName userName = (UserName) o;
        return Objects.equals(content, userName.content);
    }
}
