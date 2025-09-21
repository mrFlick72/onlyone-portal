package com.onlyoneportal.budget.user;

import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;

public class SpringSecurityUserRepository implements UserRepository {

    @Override
    public UserName currentLoggedUserName() {
        JwtAuthenticationToken principal = (JwtAuthenticationToken) SecurityContextHolder.getContext().getAuthentication();
        return new UserName(String.valueOf(principal.getToken().getClaims().get("user_name")));
    }
}
