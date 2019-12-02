package storage

import "io"

var StorageClient Storage

// storage interface
type Storage interface {
	// upload to the storage service then output the storage url
	Upload(id string, filename string, input io.Reader) (string, error)
}
