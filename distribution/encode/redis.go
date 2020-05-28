package encode

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/docker/distribution/encode"
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
	conn := r.pool.Get()
	defer conn.Close()

	key := getDBIdentifier(encodingHash)
	return redis.Bool(conn.Do("EXISTS", key))
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
func (r *RedisDB) GetMultipleEncodings(ctx context.Context, encodingHashList ...string) (map[string][]byte, error) {
	conn := r.pool.Get()
	defer conn.Close()

	now := time.Now()

	//fmt.Println("LEn of encodings:", len(encodingHashList))
	blocksFromDB := make(map[string][]byte, len(encodingHashList))
	BatchSize := 5000
	for i := 0; i < len(encodingHashList); i = i + BatchSize {
		startIndex := i
		endIndex := i + BatchSize
		if endIndex > len(encodingHashList) {
			endIndex = len(encodingHashList)
		}
		//fmt.Println("Start iNdex:", startIndex)
		//fmt.Println("End iNdex:", endIndex)

		batchKeys := encodingHashList[startIndex:endIndex]
		dbKeys := make([]interface{}, len(batchKeys))
		for j := range batchKeys {
			dbKeys[j] = getDBIdentifier(batchKeys[j])
		}
		rawBlocksFromDB, _ := redis.Values(conn.Do("MGET", dbKeys...))

		for k, rawEncoding := range rawBlocksFromDB {
			switch encoding := rawEncoding.(type) {
			case []byte:
				blocksFromDB[batchKeys[k]] = encoding
			default:
				blocksFromDB[batchKeys[k]] = nil
			}
		}
	}
	encode.PerfLog(fmt.Sprintf("Time to MGET all encodings is %s", time.Since(now)))

	return blocksFromDB, nil
}

// InsertEncoding will insert the encoding in the db
func (r *RedisDB) InsertEncoding(ctx context.Context, encodingHash string, byteStream []byte) error {
	conn := r.pool.Get()
	defer conn.Close()

	key := getDBIdentifier(encodingHash)
	_, err := conn.Do("SET", key, byteStream)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// InsertEncodings will insert a collection of encodings in the db
func (r *RedisDB) InsertEncodings(ctx context.Context, keys []string, blocks [][]byte) error {
	conn := r.pool.Get()
	defer conn.Close()

	if len(keys) == 0 {
		return nil
	}
	values := make([]interface{}, 2*len(keys))
	for i := range keys {
		values[2*i] = getDBIdentifier(keys[i])
		values[2*i+1] = blocks[i]
	}

	_, err := conn.Do("MSET", values...)
	if err != nil {
		fmt.Println(err)
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
