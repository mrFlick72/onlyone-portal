package com.onlyoneportal.budget.infrastructure.strings;

import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Base64;

public class Sha256Utils {
    private final static Base64.Encoder ENCODER = Base64.getEncoder();

    public static String sha256For(String input) {
        MessageDigest digest;
        try {
            digest = MessageDigest.getInstance("SHA-256");
        } catch (NoSuchAlgorithmException e) {
            throw new RuntimeException(e);
        }
        byte[] hash = digest.digest(input.getBytes(StandardCharsets.UTF_8));

        return ENCODER.encodeToString(hash);
    }
}
