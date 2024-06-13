# AnimeFrameBot

## Prerequisite

Go 1.22 or later is required.

## File Structure
- cmd/apiserver: The main API server program  
    `main_test.go` contains integration testing.
- internal: feature implementation  
    `*_test.go` contains unit testing.
    `http.go` contains handler implementation.

## Running

```sh
go run ./cmd/apiserver
```

This will start the API server at `http://localhost:8763`, if you want to change the port, you can change `":8763"` in `func main` of `cmd/apiserver/main.go`.

### Running tests
Hint: The following commands starts in `AnimeFrameBot/api-server` directory.


#### Testing Prerequisite
```sh
go install gotest.tools/gotestsum@latest
go install github.com/go-gremlins/gremlins/cmd/gremlins@v0.5.0
```

#### Unit Testing
frame:
```
cd ./internal/frame
go test . -v
gotestsum --format testname
```

upload:
```
cd ./internal/upload
go test . -v
gotestsum --format testname
```

#### Integration Testing
```
cd ./cmd/apiserver
go test . -v
gotestsum --format testname
```

#### Fuzz Testing
```
cd ./internal/frame
go test -fuzz=Fuzz -fuzztime 30s
```

#### Code Coverage
Per-file coverage:
```
go test -cover -covermode=count -coverpkg=./internal/... ./... -coverprofile=c.out
go tool cover -html=c.out
```

Per-function coverage:
```
go tool cover -func=c.out
```

#### Mutation Testing
```
# run mutation test
gremlins unleash --coverpkg "./internal/..." -i --timeout-coefficient 5 --invert-logical [--output mutation.json] ./internal

# to list lived mutations from file:
jq -r '.files[] | .file_name as $file | .mutations[] | select(.status == "LIVED") | "\($file):\(.line):\(.column) \(.type)"' mutation.json

# or pipe directly:
gremlins unleash --coverpkg "./internal/..." -i --timeout-coefficient 5 --invert-logical ./internal | grep --color=always "LIVED"
```
