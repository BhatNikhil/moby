package encode

import "context"

//DB is an interface to the data structure which holds the encoding
type DB interface {
	// IsEncodingAvailable will check if an encoding is avaialble in the db
	IsEncodingAvailable(ctx context.Context, encodingHash string) (bool, error)

	// GetEncoding gets encoding from the db
	GetEncoding(ctx context.Context, encodingHash string) ([]byte, error)

	//InsertEncoding will insert the encoding in the db
	InsertEncoding(ctx context.Context, encodingHash string, byteStream []byte) error

	//InsertEncodings will insert a list of encodings in the db
	InsertEncodings(ctx context.Context, encodings []string, blocks [][]byte) error

	//GetMultipleEncodings will get a list of encodings corresponding to the list provided
	GetMultipleEncodings(ctx context.Context, encodingHashList ...string) (map[string][]byte, error)
}
