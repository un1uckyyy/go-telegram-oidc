# go-telegram-oidc
Telegram Bot written on Go with OIDC authentication implementation

## environment variables
- TG_TOKEN: <'your token for telegram-bot'>
- CLIENT_ID: <'OIDC server client id'>
- CLIENT_SECRET: <'OIDC server client secret'>
- REDIRECT_HOST: <'where OIDC server will redirect user after authentication, your tg-service must expose public endpoint "/auth"'>
- AUTH_URL: <'OIDC server auth url'>
- TOKEN_URL: <'OIDC server token url'>

## see example configuration with [keycloak](https://github.com/keycloak/keycloak) as oidc provider in [compose.yaml](https://github.com/un1uckyyy/go-telegram-oidc/blob/main/compose.yaml)
