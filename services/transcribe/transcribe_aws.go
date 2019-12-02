package transcribe

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"io/ioutil"
	"net/http"
	"time"
)

type awsTranscribeService struct {
	transcriber *transcribeservice.TranscribeService
}

// initialize the aws transcription service
func InitAWSTranscribeService() {
	// create a new aws session
	sess := session.Must(session.NewSession())

	// create the aws transcribe client
	svc := transcribeservice.New(sess)

	// create the service client
	var client awsTranscribeService
	client.transcriber = svc

	Client = &client
}

// function takes a s3 bucket location then return the transcription result
func (a *awsTranscribeService) Transcribe(id string, storagePath string) (string, error) {
	// prepare the transcription input
	input := transcribeservice.StartTranscriptionJobInput{
		LanguageCode:         aws.String("en-US"),
		Media:                &transcribeservice.Media{MediaFileUri: aws.String(storagePath)},
		MediaFormat:          aws.String(transcribeservice.MediaFormatMp4),
		TranscriptionJobName: aws.String(id),
	}

	// start the transcription job
	result, err := a.transcriber.StartTranscriptionJob(&input)
	if err != nil {
		return "", err
	}

	// check if the result is in progress
	status := *result.TranscriptionJob.TranscriptionJobStatus
	if status != transcribeservice.TranscriptionJobStatusInProgress {
		if status == transcribeservice.TranscriptionJobStatusFailed {
			// failed then return as error
			return "", errors.New(*result.TranscriptionJob.FailureReason)
		}

		// check if the job has completed
		if status == transcribeservice.TranscriptionJobStatusCompleted {
			// succeed so return the result
			return GetJsonFromS3URL(*result.TranscriptionJob.Transcript.TranscriptFileUri)
		}
	}

	// job still in progress check it in an interval
	jobResult, err := a.GetTranscriptionResult(*input.TranscriptionJobName)

	return jobResult, err
}

// function to check the status and get the transcription result if completed
func (a *awsTranscribeService) GetTranscriptionResult(jobName string) (string, error) {
	// prepare the input for getting the transcription job
	input := transcribeservice.GetTranscriptionJobInput{TranscriptionJobName: aws.String(jobName)}

	// get the state of the transcription job
	result, err := a.transcriber.GetTranscriptionJob(&input)
	if err != nil {
		return "", err
	}

	// check if the job status has failed
	status := *result.TranscriptionJob.TranscriptionJobStatus
	if status == transcribeservice.TranscriptionJobStatusFailed {
		// failed then return as error
		return "", errors.New(*result.TranscriptionJob.FailureReason)
	}

	// check if the job has completed
	if status == transcribeservice.TranscriptionJobStatusCompleted {
		// succeed so return the result
		return GetJsonFromS3URL(*result.TranscriptionJob.Transcript.TranscriptFileUri)
	}

	// still in progress then wait for 1 sec then go recursive
	time.Sleep(time.Second)

	r, err := a.GetTranscriptionResult(jobName)

	return r, err
}

// helper function to grab the json content from the s3 bucket provided by transcription service
func GetJsonFromS3URL(url string) (string, error) {
	// have a http client and add time out to 1 minute (just in case)
	client := &http.Client{
		Timeout: time.Minute,
	}

	// GET from the url
	res, err := client.Get(url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	
	// make sure it got a 200 response
	if res.StatusCode != http.StatusOK {
		return "", errors.New("response from transcription result s3 bucket is not 200")
	}

	// transform body to string
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}