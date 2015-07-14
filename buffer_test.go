package buffer

import (
	"fmt"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/garyburd/redigo/redis"
)

var r *redis.Pool

func init() {
	r = redis.NewPool(dial, 0)
}

func dial() (redis.Conn, error) {
	return redis.Dial("tcp", ":6379")
}

func setup() *Buffer {
	conn := r.Get()

	conn.Do("SELECT", 42)
	conn.Do("FLUSHDB")
	conn.Flush()

	return New(r)
}

func TestIncrs(t *testing.T) {
	b := setup()

	a := "a"
	b.Incr(a)
	b.Incr(a)
	b.Incr(a)

	c := "c"
	b.Incr(c)

	assert.Equal(t, 2, b.Length())
	assert.Equal(t, 3, b.GetIncr(a))
	assert.Equal(t, 1, b.GetIncr(c))
}

func TestSets(t *testing.T) {
	b := setup()

	a := "a"
	c := "c"

	b.Set(a, "hello")
	b.Set(a, "hello")
	b.Set(c, "world")

	assert.Equal(t, 2, b.Length())

	str := string(b.GetSet(a))
	assert.Equal(t, "hello", str)

	str = string(b.GetSet(c))
	assert.Equal(t, "world", str)
}

func TestExpire(t *testing.T) {
	b := setup()
	key := "a"

	b.Expire(key, 10)

	assert.Equal(t, 10, b.GetExpire(key))
}

func TestClear(t *testing.T) {
	b := setup()
	b.Set("a", "hello")
	b.Clear()

	assert.Equal(t, 0, b.Length())
}

func TestSAdds(t *testing.T) {
	b := setup()

	k := "a"
	v := "hello"

	b.SAdd(k, v)

	assert.Equal(t, 1, b.Length())

	s := b.sadds[k]
	assert.T(t, s.Has(v))
}

func TestHsets(t *testing.T) {
	b := setup()

	k := "a"
	f := "b"
	v := "hello"

	b.Hset(k, f, v)

	assert.Equal(t, 1, b.Length())

	_, ok := b.hsets[k][f]
	assert.T(t, ok)
}

func TestLength(t *testing.T) {
	b := setup()

	b.Set("a", "hello")

	assert.Equal(t, 1, b.Length())
}

func TestFlush(t *testing.T) {
	b := setup()

	b.Set("key", "value")
	b.Incr("incr")
	b.Hset("hash", "field", "val")
	b.Expire("expire", 1)
	b.SAdd("set", "val")

	flushed := false
	b.Flush(func(err error) {
		assert.Equal(t, nil, err)
		flushed = true
	})
	assert.Equal(t, true, flushed)
}

func TestRedis(t *testing.T) {
	b := setup()

	for i := 0; i < 10; i++ {
		b.Incr("key")
	}

	b.Flush(func(error) {})
	res, _ := redis.Int64(r.Get().Do("GET", "key"))
	fmt.Println(res)
	assert.Equal(t, int64(10), res)
}
