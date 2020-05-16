# websockets

websockets for real-time communications with the frontend

## Installation

```bash
go get -v github.com/cuttle-ai/websockets
```

## Usage

Navigate into the project directory and run the following command

```bash
go run main.go
```

### Environment Variables

| Enivironment Variable           | Description                                                                                     |
| ------------------------------- | ----------------------------------------------------------------------------------------------- |
| **PORT**                        | Port on to which application server listens to. Default value is 8080                           |
| **RESPONSE_TIMEOUT**            | Timeout for the server to write response. Default value is 100ms                                |
| **REQUEST_BODY_READ_TIMEOUT**   | Timeout for reading the request body send to the server. Default value is 20ms                  |
| **RESPONSE_BODY_WRITE_TIMEOUT** | Timeout for writing the response body. Default value is 20ms                                    |
| **PRODUCTION**                  | Flag to denote whether the server is running in production. Default value is `false`            |
| **SKIP_VAULT**                  | Skip loading the configurations from vault server. Default value is `false`.                    |
| **IS_TEST**                     | Denoting the run is test. This will load the test configuration from vault                      |
| **MAX_REQUESTS**                | Maximum no. of concurrent requests supported by the server. Default value is 1000               |
| **REQUEST_CLEAN_UP_CHECK**      | Time interval after which error request app context cleanup has to be done. Default value is 2m |

## Author

Melvin Davis<melvinodsa@gmail.com>
