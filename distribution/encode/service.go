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
		Encodings: make([]bool, len(recipe.Keys)),
	}

	for i, encodingHash := range recipe.Keys {
		declaration.Encodings[i], _ = s.db.IsEncodingAvailable(ctx, encodingHash)
	}
	return declaration, nil
}

//GetAvailableBlocksFromDB gets available blocks from db and constructs a declaration
func (s *Service) GetAvailableBlocksFromDB(ctx context.Context, recipe encode.Recipe) (encode.Declaration, [][]byte, error) {
	blocksFromDB, err := s.db.GetMultipleEncodings(ctx, recipe.Keys...)
	declaration := encode.Declaration{
		Encodings: make([]bool, len(recipe.Keys)),
	}
	for i, v := range blocksFromDB {
		declaration.Encodings[i] = (v != nil)
	}
	return declaration, blocksFromDB, err
}

// InsertMissingEncodings will insert the encoding in the backend data store
func (s *Service) InsertMissingEncodings(ctx context.Context, recipe encode.Recipe, d encode.Declaration, byteStream []byte) error {
	for i, exists := range d.Encodings {
		if exists == false {
			startIndex, endIndex := encode.BlockIndices(i, len(byteStream))
			s.db.InsertEncoding(ctx, recipe.Keys[i], byteStream[startIndex:endIndex])
		}
	}
	return nil
}

// AssembleBlob will assemble the blob using the recipe and the byte streams
func (s *Service) AssembleBlob(ctx context.Context, r encode.Recipe, b encode.BlockResponse, dbBlocks [][]byte, lengthOfByteStream int) ([]byte, error) {
	blockResponse := make([]byte, lengthOfByteStream)

	if Debug == true {
		fmt.Println("Length of byte stream: ", lengthOfByteStream)
		fmt.Println("Length of recipe: ", len(r.Keys))
		fmt.Println("Length of Blocks from DB: ", len(dbBlocks))
	}

	for i, val := range dbBlocks {
		key := r.Keys[i]

		var block []byte
		if val == nil {
			block = b.Blocks[i]
		} else {
			block = val
			fmt.Println("Block fetched from db:", key)
		}

		_, endIndex := encode.BlockIndices(i, lengthOfByteStream)
		startIndex := endIndex - len(block)

		if Debug == true {
			fmt.Println("=================================")
			fmt.Println("Got block of length: ", len(block))
			fmt.Println("Start index: ", startIndex)
			fmt.Println("End index: ", endIndex)
			fmt.Println("=================================")
		}

		copy(blockResponse[startIndex:endIndex], block)
	}

	return blockResponse, nil
}
