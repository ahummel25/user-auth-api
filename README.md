# Simple User Auth API built with Go and the [Serverless framework](https://www.serverless.com/framework/docs/getting-started)

Reference the `serverless.yml` file in the root of the project. This is the main file used for deploying the serverless application.

## Install Go
- https://go.dev/doc/install

## Install Serverless
- https://www.serverless.com/framework/docs/getting-started

## Setup

```sh
$ git clone git@github.com:ahummel25/user-auth-api.git
$ cd user-auth-api
$ go mod install
$ go mod tidy
```

## Development

```sh
$ go run server.go

connect to http://localhost:8080/graphiql for GraphQL playground
```

## Test

```sh
$ go test ./...
```

## Build/Deploy (Must have AWS credentials configured)

```sh
$ make deploy
```

## License

MIT © [Andrew Hummel](https://andrewhummel.dev)