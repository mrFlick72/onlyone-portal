package com.onlyoneportal.budget.security;

import java.io.IOException;

import org.springframework.http.HttpRequest;
import org.springframework.http.client.ClientHttpRequestExecution;
import org.springframework.http.client.ClientHttpRequestInterceptor;
import org.springframework.http.client.ClientHttpResponse;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.server.resource.authentication.BearerTokenAuthentication;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;

public class BearerTokenInterceptor implements ClientHttpRequestInterceptor {

  @Override
  public ClientHttpResponse intercept(
      HttpRequest request, byte[] body, ClientHttpRequestExecution execution)
      throws IOException {

    Authentication auth = SecurityContextHolder.getContext().getAuthentication();

    if (auth instanceof JwtAuthenticationToken jwtAuth) {
      request.getHeaders().setBearerAuth(jwtAuth.getToken().getTokenValue());
    } else if (auth instanceof BearerTokenAuthentication bearer) {
      request.getHeaders().setBearerAuth(bearer.getToken().getTokenValue());
    }

    return execution.execute(request, body);
  }
}