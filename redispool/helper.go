package redispool

import (
	"github.com/garyburd/redigo/redis"
)

// reply helper func
var (
	Bool       = redis.Bool
	ByteSlices = redis.ByteSlices
	Bytes      = redis.Bytes
	Float64    = redis.Float64
	Int        = redis.Int
	Int64      = redis.Int64
	Int64Map   = redis.Int64Map
	IntMap     = redis.IntMap
	Ints       = redis.Ints
	MultiBulk  = redis.MultiBulk
	Scan       = redis.Scan
	ScanSlice  = redis.ScanSlice
	ScanStruct = redis.ScanStruct
	String     = redis.String
	StringMap  = redis.StringMap
	Strings    = redis.Strings
	Uint64     = redis.Uint64
	Values     = redis.Values
)

// if not exists in redis, err should be ErrNil
var (
	ErrNil = redis.ErrNil
)
