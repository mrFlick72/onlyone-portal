

server:
  port: 3050

idp:
  jwks-endpoint: http://local.api.vauthenticator.com:9090/oauth2/jwks

user:
  required-role: USER_ROLE

logger:
  file-name: logs.log


cors:
  allowed:
    origins: "http://local.onlyone-portal.com:8070"

CONFIG_FILE_LOCATION