package encode

import (
	"context"
	"fmt"

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

	fmt.Println("Length of byte stream: ", lengthOfByteStream)
	fmt.Println("Length of recipe: ", len(r.Recipe))
	fmt.Println("Length of declaration: ", len(d.Encodings))

	for i, val := range d.Encodings {
		key := r.Recipe[i]

		var block []byte
		if val == true {
			block, _ = s.db.GetEncoding(ctx, key)
		} else {
			block = b.Blocks[i]
		}
		startIndex := i * encode.ShiftOfWindow
		endIndex := i*encode.ShiftOfWindow + len(block)

		fmt.Println("=================================")
		fmt.Println("Got block of length: ", len(block))
		fmt.Println("Start index: ", startIndex)
		fmt.Println("End index: ", endIndex)
		fmt.Println("=================================")

		copy(blockResponse[startIndex:endIndex], block)
	}

	return blockResponse, nil
}
