package encode

import (
	"context"
)

// RedisDB provides the service to retreive
// and manage the encodings from the db
type s3DB struct {
}

// IsEncodingAvailable will check if encoding is present inside the db
func (s3DB *s3DB) IsEncodingAvailable(ctx context.Context, encodingHash string) (bool, error) {
	return false, nil
}

// GetEncoding will get the encoding from the db
func (s3DB *s3DB) GetEncoding(ctx context.Context, encodingHash string) ([]byte, error) {
	return nil, nil
}

// GetMultipleEncodings will get the list of encodings from the db
func (s3DB *s3DB) GetMultipleEncodings(ctx context.Context, encodingHashList ...string) ([][]byte, error) {
	return nil, nil
}

// InsertEncoding will insert the encoding in the db
func (s3DB *s3DB) InsertEncoding(ctx context.Context, encodingHash string, byteStream []byte) error {
	return nil
}
