## OAuth2.0 server

During the invention of another one wheel I produced this. Client management, user management and oauth server - all in one.

### What you can find here:

* OAuth2.0 proxy server (*Resource Owner Password Credentials Grant https://www.rfc-editor.org/rfc/rfc6749#section-4.3)
* Client management
* User management

### API

Client:

* **POST** /clients - create client

User:

* **POST** /users - create user
* **PATCH** /users - update user

Token:

* **POST** /oauth/token - create user token
* **POST** /oauth/token/refresh - refresh user token
* **ANY** /oauth/token/validate - validate user token (provides user_id in `X-User-Id` response header as well)