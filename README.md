# gowork
Concurrency patterns for golang

These concurrency patterns make use of closures to execute your tasks

## MutliTask

MultiTask can be used when you want to run several tasks in parallel and wait for them all to complete

## Batch

Batch can be used when you have a task you want to do for every item in a list but you want to limit 
the number of items processed at any given time.
