# go-telegram-oidc
Telegram Bot written on Go with OIDC authentication implementation

## environment variables
- TG_TOKEN: <'your token for telegram-bot'>
- CLIENT_ID: <'OICD server client id'>
- CLIENT_SECRET: <'OICD server client secret'>
- REDIRECT_HOST: <'where OICD server will redirect user after authentication, your tg-service must expose public endpoint "/auth"'>
- AUTH_URL: <'OICD server auth url'>
- TOKEN_URL: <'OICD server token url'>

## see example configuration with [keycloak](https://github.com/keycloak/keycloak) as OICD provider in [compose.yaml](https://github.com/un1uckyyy/go-telegram-oidc/blob/main/compose.yaml)
