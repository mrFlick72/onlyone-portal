package com.onlyoneportal.budget.infrastructure.dynamodb;

import java.util.UUID;

public class UUIDSaltGenerator implements SaltGenerator {
    @Override
    public String newSalt() {
        return UUID.randomUUID().toString();
    }
}
