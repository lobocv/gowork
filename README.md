# gowork
Concurrency patterns for golang

These concurrency patterns make use of closures to execute your tasks

## MutliTask

MultiTask can be used when you want to run several tasks in parallel and wait for them all to complete

## Batch

Batch can be used when you have a task you want to do for every item in a list but you want to limit 
the number of items processed at any given time.

## BufferedBatch

BufferedBatch can be used when you have a task you want to run for every X number of items in a slice. For example
if you have a slice [1, 2, 3, 4, 5, 6, 7, 8] and a buffer size of 3, it would run jobs for [1,2,3], [4, 5, 6], [7, 8]  