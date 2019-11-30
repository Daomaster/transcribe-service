package storage

// storage interface
type Storage interface {
	// upload to the storage service then output the storage url
	Upload(filename string) (string, error)
}
