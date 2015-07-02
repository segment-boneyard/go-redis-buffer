# go-redis-buffer

Buffer that aggregates and flushes redis queries.

## Example

```go
buf := buffer.New(redis)
buf.Incr("hello-world")
buf.Flush(func(err error){
    fmt.Errorf("redis error: %s, err")
})
```

## License

MIT
