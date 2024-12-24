# SIMPLE SECURE REVERSE PROXY

## Introduction

Or at least i hope it is simple.
ssrp is a reverse proxy specifically made to be stupidly simple to proxy based on a header containing a value in a list of values (comma separated).

so if a request has `x-groups: one,two,three` and we run with `-g one` the request will be allowed, however if we use `-g four` in this example the request will be denied.

## usage

```bash
> ssrp -h
Usage of ./ssrp:
  -g, --group strings   Allowed groups (can be specified multiple times)
  -i, --insecure        Ignore SSL certificate errors
  -l, --listen string   Address to listen on (default ":3000")
  -t, --target string   Target URL to proxy to (default "localhost:8080")
```

## schematic overview

```mermaid
sequenceDiagram
    participant Client
    participant Ingress as ingress-nginx
    participant OAuth2 as oauth2-proxy
    participant Azure as Azure Entra
    participant SSRP
    participant App as Application

    Client->>Ingress: Request
    Ingress->>OAuth2: Forward request
    alt User not authenticated
        OAuth2->>Azure: Redirect for authentication
        Azure->>Client: Login page
        Client->>Azure: Provide credentials
        Azure->>OAuth2: Auth token
        OAuth2->>Client: Set auth cookie
    end
    OAuth2->>SSRP: Forward authenticated request
    SSRP->>App: Route request
    App->>SSRP: Response
    SSRP->>OAuth2: Forward response
    OAuth2->>Ingress: Forward response
    Ingress->>Client: Deliver response
```
