---
title: "Stargate"
linkTitle: "Stargate"
weight: 1
description: |
  Accessing the K8ssandra Stargate interfaces.
---

[Stargate](https://stargate.io/) is an open-source data gateway providing common
API interfaces for backend databases. With K8ssandra, Stargate may be deployed
in front of the Apache Cassandra cluster providing CQL, REST, GraphQL, and document-based API
endpoints. Stargate itself may be scaled horizontally within the cluster as needed. 
This scaling is done independently from the data layer.

This guide provides information about accessing the various API endpoints
provided by Stargate.

While this document will help get you up and going quickly with Stargate, more detailed 
information about using Stargate can be found in the Stargate 
[docs](https://stargate.io/docs/stargate/1.0/quickstart/quickstart.html). 

## Tools

* HTTP client (cURL, Postman, etc.)
* Web Browser

## Prerequisites

1. [K8ssandra Cluster]({{< ref "getting-started#install-k8ssandra" >}})
1. [Ingress]({{< ref "ingress" >}}) configured to expose each of the Stargate services (Auth, REST, GraphQL)
1. DNS names configured for the exposed Stargate services, referred to as `STARGATE_AUTH_DOMAIN`, `STARGATE_REST_DOMAIN`, and `STARGATE_GRAPHQL_DOMAIN` below.

## Access Auth API

Before accessing any of the provided Stargate data APIs, an auth token must be generated and provided
to subsequent data API requests.  Use the auth API to generate a token.

The default port exposed by Stargate for the auth API is `8081`, these examples will assume that is the
port exposed by the cluster ingress configuration for access.

The authorization API can be accessed at: [http://STARGATE_AUTH_DOMAIN/v1/auth](http://STARGATE_AUTH_DOMAIN/v1/auth)

Detailed information about the Stargate auth API can be found in the Stargate [docs](https://stargate.io/docs/stargate/1.0/developers-guide/auth.html).

### Extracting Cassandra username/password Secrets

The auth API requires the Cassandra username and password to be provided to it.  Those values can be 
extracted from the K8ssandra cluster through the following commands (replace `k8ssandra-cluster` with the
name configured for your running cluster).

Extract and decode the username secret:

```
kubectl get secret k8ssandra-cluster-superuser -o jsonpath="{.data.username}" | base64 --decode
```

Extract and decode the password secret:

```
kubectl get secret k8ssandra-cluster-superuser -o jsonpath="{.data.password}" | base64 --decode
```

### Generating Auth Tokens

Next, use the extracted and decoded secrets to request a token from the Stargate auth API.

```
curl -L -X POST 'http://_STARGATE_DOMAIN_/v1/auth' -H 'Content-Type: application/json' --data-raw '{"username": "k8ssandra-cluster-superuser", "password": "1LI8TebjjHYrqUk9xYbJnbYJheX3Ckq250byd2ePDPXNtweaYgznmg"}'
```

This request will return a response similar to the following. The value given for `authToken` will be required when making requests to the Stargate data APIs.

```
{"authToken":"e4b34bbc-0ebc-4e2a-86ca-04793ca658a7"}
```

### Using Auth Tokens

Stargate supports authorization within the data APIs through a custom HTTP header `x-cassandra-token`, which must be populated with the token given by the auth API.

## Access Document Data API

The Stargate document APIs provide a way schemaless way to store and interact with data inside of Cassandra.
The first step is to [create a namespace](https://stargate.io/docs/stargate/1.0/quickstart/quick_start-document.html#_creating_schema). 
That can be done with a request to the `/v2/schemas/namespaces` API:

```
curl --location --request POST 'http://STARGATE_REST_DOMAIN/v2/schemas/namespaces' \
--header "x-cassandra-token: e4b34bbc-0ebc-4e2a-86ca-04793ca658a7" \
--header 'Content-Type: application/json' \
--data '{
    "name": "mynamespace"
}'
```

That will use the auth token previously generated to request the creation of a namespace called `mynamespace`. The 
server should return a response like:

```
{"name":"mynamespace"}
```

Additional information related to using the Document APIs can be found in the Stargate [docs](https://stargate.io/docs/stargate/1.0/quickstart/quick_start-document.html).

## Access REST Data API

The Stargate REST APIs provide a RESTful way to store and interact with data inside of Cassandra that should feel
familiar to developers. Unlike the document APIs, some understanding of Cassandra data modeling will be required. The first step is to [create a keyspace](https://stargate.io/docs/stargate/1.0/quickstart/quick_start-rest.html#_creating_schema). 
That can be done with a request to the `/v2/schemas/keyspaces` API:

```
curl --location --request POST 'http://STARGATE_REST_DOMAIN/v2/schemas/keyspaces' \
--header "x-cassandra-token: e4b34bbc-0ebc-4e2a-86ca-04793ca658a7" \
--header 'Content-Type: application/json' \
--data '{
    "name": "mykeyspace"
}'
```

That will use the auth token previously generated to request the creation of a keyspace called `mykeyspace`. 
The server should return a response like:

```
{"name":"mykeyspace"}
```

Additional information related to using the Document APIs can be found in the Stargate [docs](https://stargate.io/docs/stargate/1.0/quickstart/quick_start-rest.html).

## Access GraphQL Data API

The Stargate GraphQL APIs provide a way to store and interact with data inside of Cassandra using the powerful GraphQL
query language and tooling ecosystem. Like the REST APIs, this does require some additional Cassandra data modeling
understanding. Like the REST APIs, The first step to using the GraphQL APIs is to [create a keyspace](https://stargate.io/docs/stargate/1.0/quickstart/quick_start-graphql.html#_creating_schema).

The easiest way to get started with the GraphQL APIs is to use the built-in GraphQL playground described in the next section.

Additional information related to using the Document APIs can be found in the Stargate [docs](https://stargate.io/docs/stargate/1.0/quickstart/quick_start-graphql.html).

### Access GraphQL playground

Stargate's GraphQL service provides an interactive "playground" application that can be used to interact with the GraphQL APIs.

The playground application can be accessed at [http://STARGATE_GRAPHQL_DOMAIN/playground](http://STARGATE_GRAPHQL_DOMAIN/playground).

Detailed information related to using the GraphQL playground can be found in the Stargate [docs](https://stargate.io/docs/stargate/1.0/developers-guide/graphql-using.html#_using_the_graphql_playground).
