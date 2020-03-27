package encode

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/garyburd/redigo/redis"
)

// RedisDB provides the service to retreive
// and manage the encodings from the db
type RedisDB struct {
	pool *redis.Pool
}

// NewRedisDB will generate a new DB object satisfying the ENcodeDB ifc
func NewRedisDB() RedisDB {
	return RedisDB{
		pool: configureRedis(),
	}
}

func getDBIdentifier(encodingHash string) string {
	return "recipe:" + encodingHash
}

// IsEncodingAvailable will check if encoding is present inside the db
func (r *RedisDB) IsEncodingAvailable(ctx context.Context, encodingHash string) (bool, error) {
	return false, nil
}

// GetEncoding will get the encoding from the db
func (r *RedisDB) GetEncoding(ctx context.Context, encodingHash string) ([]byte, error) {
	conn := r.pool.Get()
	defer conn.Close()

	key := getDBIdentifier(encodingHash)
	rawEncoding, err := conn.Do("GET", key)

	switch encoding := rawEncoding.(type) {
	case []byte:
		return encoding, err

	default:
		return nil, err
	}
}

// GetMultipleEncodings will get the list of encodings from the db
func (r *RedisDB) GetMultipleEncodings(ctx context.Context, encodingHashList ...string) ([][]byte, error) {
	return nil, nil
}

// InsertEncoding will insert the encoding in the db
func (r *RedisDB) InsertEncoding(ctx context.Context, encodingHash string, byteStream []byte) error {
	conn := r.pool.Get()
	defer conn.Close()

	key := getDBIdentifier(encodingHash)
	ok, err := conn.Do("PUT", key, byteStream)

	if ok != "OK" {
		return fmt.Errorf("Not Okay")
	}
	return err
}

func configureRedis() *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				log.Printf("ERROR: fail init redis pool: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}

	return pool
}
