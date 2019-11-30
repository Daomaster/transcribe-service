package transcribe

// transcription service interface
type Transcriber interface {
	// transcription takes a storage url then output the json string result
	Transcribe(storagePath string) (string, error)
}
