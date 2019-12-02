package transcribe

var Client Transcribe

// transcription service interface
type Transcribe interface {
	// transcription takes a storage url then output the json string result
	Transcribe(id string, storagePath string) (string, error)
}
