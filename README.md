# go-greetings

A small go app that offers a rest api to store greetings in and retrieve them from a db.

## Docker image

When using the docker image specify the postgres connection string in the environment variable `GOGREETING_DATASOURCE`.

## Endpoints

It has some very simple url based endpoints that you can use directly from the browser.

### remember

Save a greeting for john:

```
localhost:8080/remember/john/greet_john
```

### greet

Return the greeting for john

```
localhost:8080/greet/john
```