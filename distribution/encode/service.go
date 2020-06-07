package encode

import (
	"context"
	"crypto/md5"
	"encoding/base64"
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

// InsertMissingEncodings will insert the encoding in the backend data store
func (s *Service) InsertMissingEncodings(ctx context.Context, blockKeysInDB []string, byteStream []byte) error {
	var keys []string
	var blocks [][]byte

	for i, blockKey := range blockKeysInDB {
		if blockKey == "0" { // block key not in db
			startIndex, endIndex := encode.BlockIndices(i, len(byteStream))
			block := byteStream[startIndex:endIndex]
			keyAsBytes := md5.Sum(block) //TODO: Move this code to the same place sevrer library and import it here
			keyBase64 := base64.StdEncoding.EncodeToString(keyAsBytes[:])
			keys = append(keys, keyBase64[0:len(keyBase64)-2])
			blocks = append(blocks, block)
		}
	}

	s.db.InsertEncodings(ctx, keys, blocks)

	return nil
}

// AssembleBlob will assemble the blob using the recipe and the byte streams
func (s *Service) AssembleBlob(ctx context.Context, b encode.BlockResponse, blockKeys []string, lengthOfByteStream int) ([]byte, error) {
	blockResponse := make([]byte, lengthOfByteStream)

	if Debug == true {
		fmt.Println("Length of byte stream: ", lengthOfByteStream)
	}

	var blockKeysFromDB []string
	for _, v := range blockKeys {
		if v != "0" {
			blockKeysFromDB = append(blockKeysFromDB, v)
		}
	}
	blocksFromDB, _ := s.db.GetMultipleEncodings(ctx, blockKeysFromDB...)

	for i, val := range b.Blocks {
		var block []byte
		if val != nil {
			block = b.Blocks[i]
		} else {
			block, _ = blocksFromDB[blockKeys[i]]
		}

		startIndex, endIndex := encode.BlockIndices(i, lengthOfByteStream)

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
