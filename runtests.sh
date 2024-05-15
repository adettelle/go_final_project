#!/bin/bash
echo 'running go test -run ^TestApp$ ./tests'
go test -run ^TestApp$ ./tests

echo 'running go test -run ^TestDB$ ./tests'
go test -run ^TestDB$ ./tests

echo 'running go test -run ^TestAddTask$ ./tests'
go test -run ^TestAddTask$ ./tests

echo 'running go test -run ^TestNextDate$ ./tests'
go test -run ^TestNextDate$ ./tests

echo 'running go test -run ^TestTasks$ ./tests'
go test -run ^TestTasks$ ./tests

echo 'go test -run ^TestEditTask$ ./tests'
go test -run ^TestEditTask$ ./tests

echo 'go test -run ^TestDone$ ./tests'
go test -run ^TestDone$ ./tests

echo 'go test -run ^TestDelTask$ ./tests'
go test -run ^TestDelTask$ ./tests