package storage

import "io"

var StorageClient Storage

// storage interface
type Storage interface {
	// upload to the storage service then output the storage url
	Upload(filename string, input io.Reader) (string, error)
}
