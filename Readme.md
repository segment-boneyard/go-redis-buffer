# go-redis-buffer

A buffer that for making batched updates to Redis.

Our Redis instances often ended up being the bottleneck for a write-heavy
workload. Fortunately most of the operations can applied in any order, so we
end up batching them up.

Instead of sending 10 `INCR` commands, we send a single `INCRBY 10` command.
Similarly, we'll send along only the latest `SET` and `MSET` commands.

Note this buffer is *not* concurrency safe. It should be run within a single
goroutine.

## Example

```go
buf := buffer.New(redis)
buf.Incr("hello-world")
buf.Flush(func(err error){
    if err != nil {
      fmt.Errorf("redis error: %s, err")
    }
    // clear when we're done
    buf.clear()
})
```

## Supported Operations

 - *INCR*
 - *SET*
 - *INCRYBY*
 - *HSET*
 - *EXPIRE*

## License

MIT
