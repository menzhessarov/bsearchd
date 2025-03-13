## bsearchd
bsearchd is an API service that provides a single endpoint for retrieving an index of a value.
The values are stored in memory of the service and are loaded on start-up.

### How to run

1. Clone the project
2. Make sure you've installed Go with the version that is shown in `go.mod` file
3. Run `go mod tidy` to install the project dependencies
4. Run `make run` to run the web server.
5. [Optional] to run test use this command `make test`

Use this command to interact with the API. Change *100* to any number.
```shell
  curl localhost:8080/values/100
```

### Configuration

You can change these variables in `.env` file in order to configure the service:

**LOG_LEVEL** - sets the logging level for the whole service, valid values are `debug`, `info`, `error`.

**HTTP_PORT** - sets the port on which the web server will run.

**INPUT_FILE** - sets the path of the file that contains the values which are then stored in memory by the service.

**CONFORMATION** - sets the leeway (as %) for returnable value i.e. how close to the queried value the result can be.  