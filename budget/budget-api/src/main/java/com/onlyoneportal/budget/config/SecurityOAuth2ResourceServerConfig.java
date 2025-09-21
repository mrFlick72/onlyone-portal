package com.onlyoneportal.budget.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.config.annotation.web.configurers.CsrfConfigurer;
import org.springframework.security.config.annotation.web.configurers.oauth2.server.resource.OAuth2ResourceServerConfigurer;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationConverter;
import org.springframework.security.oauth2.server.resource.web.DefaultBearerTokenResolver;
import org.springframework.security.web.SecurityFilterChain;

import java.util.*;

import static java.util.stream.Collectors.toList;
import static org.springframework.security.config.Customizer.*;

@Configuration(proxyBeanMethods = false)
public class SecurityOAuth2ResourceServerConfig {

    @Value("${granted-role.budget-api}")
    private List<String> grantedRole;

    @Bean
    public JwtAuthenticationConverter jwtAuthenticationConverter() {
        JwtAuthenticationConverter jwtAuthenticationConverter = new JwtAuthenticationConverter();
        jwtAuthenticationConverter.setJwtGrantedAuthoritiesConverter(jwt -> {
            List<String> authorities = (List<String>) Optional.ofNullable(jwt.getClaim("authorities")).orElse(new ArrayList<String>());
            List<String> scopes = (List<String>) Optional.ofNullable(jwt.getClaim("scope")).orElse(new ArrayList<String>());
            List<GrantedAuthority> grantedAuthorities = authorities.stream().map(SimpleGrantedAuthority::new).collect(toList());
            List<GrantedAuthority> grantedScopes = scopes.stream().map(SimpleGrantedAuthority::new).collect(toList());
            List<GrantedAuthority> roles = new ArrayList<>();
            roles.addAll(grantedAuthorities);
            roles.addAll(grantedScopes);
            return roles;
        });

        return jwtAuthenticationConverter;
    }


    @Bean
    public SecurityFilterChain defaultSecurityFilterChain(HttpSecurity http,
                                                          @Value("${spring.security.oauth2.resourceserver.jwt.jwk-set-uri}") String jwksUri,
                                                          JwtAuthenticationConverter jwtAuthenticationConverter) throws Exception {
        return http.csrf(AbstractHttpConfigurer::disable).cors(withDefaults())
                .authorizeHttpRequests(authz ->
                        authz.requestMatchers("/actuator/**").permitAll()
                                .anyRequest().hasAnyAuthority(grantedRole.toArray(new String[grantedRole.size()]))
                )
                .oauth2ResourceServer(resourceServerConfigurer ->
                        resourceServerConfigurer
                                .jwt(jwtConfigurer ->
                                        jwtConfigurer.jwkSetUri(jwksUri)
                                                .jwtAuthenticationConverter(jwtAuthenticationConverter))
                                .bearerTokenResolver(bearerTokenResolver())
                ).build();
    }


    private DefaultBearerTokenResolver bearerTokenResolver() {
        DefaultBearerTokenResolver bearerTokenResolver = new DefaultBearerTokenResolver();
        bearerTokenResolver.setAllowFormEncodedBodyParameter(true);
        bearerTokenResolver.setAllowUriQueryParameter(true);
        return bearerTokenResolver;
    }

}
