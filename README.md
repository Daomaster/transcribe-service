# Transcription Service
### Description
This is the transcription service for the STT exercise,
it contains two part backend service (Golang, MySQL) and [simple web app (Angular)][l3].
The application is hosted in AWS EC2 with `docker-compose`. 
And under the hood it is using the AWS transcribe service.
The application allow user to upload a mp4 file and get the transcript from AWS transcribe service.

### API
Note: The postman collection and environment is in the /utility directory, 
and in the collection the `POST /api/auth/login` request will automatically embed jwt in the postman environment after get the response.

```sh
POST /api/users    
POST /api/auth/login
GET /api/auth/refresh  
POST /api/transcription
GET /api/transcription
GET /api/transcription/:id
```

### Docker Images
[transcription-service-ui][l1]

[transcription-service][l2]


[l1]: <https://hub.docker.com/repository/docker/daomaster/transcribe-service-ui>
[l2]: <https://hub.docker.com/repository/docker/daomaster/transcribe-service>
[l3]: <https://github.com/Daomaster/transcribe-service-ui>

### Run Locally
The whole service be can run locally, but it needs to provide a few environment variables (.env works also)
```bash
AWS_REGION="region string of the aws services"
AWS_ACCESS_KEY_ID="aws credential key id"
AWS_SECRET_ACCESS_KEY="aws credential key"
BUCKET_NAME="aws bucket name"
MYSQL_PASSWORD="mysql password"
MYSQL_USER="mysql user name"
```

Once the envs are met, copy the docker-compose.yml to desired directory,
then run
```bash
docker-compose up -d
```

Note: If only want to run backend service, then you can just create a `config.json` file in the project directory

Example as below, but `AWS_REGION` `AWS_ACCESS_KEY_ID` `AWS_SECRET_ACCESS_KEY` still needs to be environment variables
```json
{
  "MYSQL_HOSTNAME": "localhost",
  "MYSQL_ROOT_PWD": "test123",
  "MYSQL_USER": "root",
  "AWS_BUCKET_NAME": "fm-demo-stt"
}
```

#### If you have any questions feel free to reach out to me
