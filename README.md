# goworktidy

Adds `go mod edit` commands to modules defined in `go.work` to your subprojects,
allowing `go mod tidy` to run.

## The Problem

Running `go mod tidy` in a workspace with modules of different names, or with
unpublished modules, will fail.

## The Manual Solution

You can run `go mod edit -replace module=../path` in each subproject to allow
`go mod tidy` to run.

## The Automated Solution

`goworktidy` will run identify all modules in your `go.work`, then identify
which subprojects are missing replace declarations, and add them.
