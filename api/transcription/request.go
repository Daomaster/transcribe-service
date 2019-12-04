package transcription

import "mime/multipart"

type CreateTranscriptionRequest struct {
	VideoFile *multipart.FileHeader `form:"file" binding:"required"`
}

type GetTranscriptionByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}