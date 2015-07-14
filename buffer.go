package buffer

import (
	"github.com/garyburd/redigo/redis"
	"gopkg.in/fatih/set.v0"
)

// Buffer.
type Buffer struct {
	pool    *redis.Pool
	incrs   map[string]int
	sets    map[string]string
	sadds   map[string]*set.Set
	hsets   map[string]map[string]string
	expires map[string]int
}

// Create new Buffer instance.
func New(pool *redis.Pool) *Buffer {
	b := Buffer{pool: pool}
	b.Clear()
	return &b
}

// Increment value for key.
func (b *Buffer) Incr(key string) {
	b.incrs[key] += 1
}

// Expire a key after seconds.
func (b *Buffer) Expire(key string, seconds int) {
	b.expires[key] = seconds
}

// Set value to store for key.
func (b *Buffer) Set(key string, value string) {
	b.sets[key] = value
}

// Set hash key value for key.
func (b *Buffer) Hset(key, field, val string) {
	m, ok := b.hsets[key]
	if !ok {
		m = map[string]string{}
	}
	m[field] = val
	b.hsets[key] = m
}

// Add value to set.
func (b *Buffer) SAdd(key string, value string) {
	s, ok := b.sadds[key]
	if !ok {
		s = set.New()
	}
	s.Add(value)
	b.sadds[key] = s
}

// Get stored incr value for key.
func (b *Buffer) GetIncr(key string) int {
	return b.incrs[key]
}

// Get stored expired value for key.
func (b *Buffer) GetExpire(key string) int {
	return b.expires[key]
}

// Get stored set value for key.
func (b *Buffer) GetSet(key string) string {
	return b.sets[key]
}

// Get total count of all keys.
func (b *Buffer) Length() int {
	return len(b.incrs) +
		len(b.sets) +
		len(b.sadds) +
		len(b.hsets)
}

// Clear stored values.
func (b *Buffer) Clear() {
	b.incrs = map[string]int{}
	b.sets = map[string]string{}
	b.hsets = map[string]map[string]string{}
	b.expires = map[string]int{}
	b.sadds = map[string]*set.Set{}
}

type ErrFunc func(error)

// Flush calls to redis.
func (b *Buffer) Flush(fn ErrFunc) {
	db := b.pool.Get()
	defer db.Close()

	for k, v := range b.incrs {
		fn(db.Send("INCRBY", k, v))
	}

	args := []interface{}{}
	for k, v := range b.sets {
		args = append(args, k, v)
	}
	fn(db.Send("MSET", args...))

	for k, v := range b.expires {
		fn(db.Send("EXPIRE", k, v))
	}

	for k, s := range b.sadds {
		args := []interface{}{k}
		for _, e := range s.List() {
			args = append(args, e)
		}
		fn(db.Send("SADD", args...))
	}

	for k, m := range b.hsets {
		args := []interface{}{k}
		for k, v := range m {
			args = append(args, k, v)
		}
		fn(db.Send("HMSET", args...))
	}

	fn(db.Flush())
}
