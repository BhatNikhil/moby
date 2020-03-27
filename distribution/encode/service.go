package encode

import (
	"context"

	"github.com/docker/distribution/encode"
)

//Service class embodices the functionality to
type Service struct {
	db DB
}

//NewService creates a new struct for Service
func NewService(database DB) Service {
	return Service{
		db: database,
	}
}

// GetDeclaration will get the 'declaration' which indicates which
// encodings reffered by the recipe are already held by the service
// for the recipe
func (s *Service) GetDeclaration(ctx context.Context, recipe encode.Recipe) (encode.Declaration, error) {
	declaration := encode.Declaration{
		Encodings: make([]bool, len(recipe.Recipe)),
	}

	for i, encodingHash := range recipe.Recipe {
		declaration.Encodings[i], _ = s.db.IsEncodingAvailable(ctx, encodingHash)
	}
	return declaration, nil
}

// InsertEncoding will insert the encoding in the backend data store
func (s *Service) InsertEncoding(ctx context.Context, encodingHash string, byteStream []byte) error {
	return s.db.InsertEncoding(ctx, encodingHash, byteStream)
}

// AssembleBlob will assemble the blob using the recipe and the byte streams
func (s *Service) AssembleBlob(ctx context.Context, r encode.Recipe, b encode.BlockResponse, d encode.Declaration, lengthOfByteStream int) ([]byte, error) {
	blockResponse := make([]byte, lengthOfByteStream)
	for i, val := range d.Encodings {
		key := r.Recipe[i]

		var block []byte
		if val == true {
			block, _ = s.db.GetEncoding(ctx, key)
		} else {
			block = b.Blocks[i]
		}
		copy(blockResponse[i*encode.ShiftOfWindow:i*encode.ShiftOfWindow+len(block)], block)
	}

	return blockResponse, nil
}
