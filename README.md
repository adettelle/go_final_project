# ToDo List

#### It is a simple web application to keep daily tasks.

With this application you can:
- Add task;
- Get list of tasks;
- Delete task;
- Change task's parameters;
- Mark task as completed.

#### This project includes ability to sign in.

If you mark a task as completed, it is transferred to the next date in accordance with the rule.

## Running the project

To run the project locally use the following command:
```
go run ./cmd/main.go
```
Settings are passed with the following environment variables:

- `TODO_PASSWORD` - required. Api authentication password.
- `TODO_DBFILE` - optional. SQLite DB file location.
- `TODO_PORT` - optional. Rest api port.

#### Running with docker 
``` bash
docker build -t go_final_project:v1.0.0 .
```

``` bash
docker run -p 7540:7540 -e TODO_PASSWORD=1234567 golang-go_final_project:v1.0.0
```

## Running tests

``` bash
# Starting tests requires running application.
go test -v ./tests
```