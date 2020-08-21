 This project contains a set of utilities for interfacing with AWS notification and queue services. Using the AWS CLI
 and SDK a you can configure a AWS Simple Notification Service (SNS) topic and AWS Simple Queue Service (SQS) queue and 
 then listen for events on the queue. The original insipiration for this was to use this as part of a proof-of-concept 
 demonstrating monitoring of infrastructure change events (via CloudTrail and/or CloudWatch).
 
 ### Installing the AWS SDK
 ```
$ go get github.com/aws/aws-sdk-go
```

 ### Set up a Topic and Queue 
 ##### Build 
 ```
$ go build setupTopicAndQueue.go
```
 ##### Running
 ``` 
 $  ./setupTopicAndQueue -profile=foo -region=us-gov-west-1
 ```
 Parameter Options
 
 | Parameter Name | Description | Default Value |
|---|---|---|
| profile | sets the AWS Cli profile to be used for accessing AWS SDK  | default |
| region | sets the region to operate one | us-east-1|
| topic | defines the SNS topic name where events should be sent | infrastructure-event |
| queue | sets the SQS queue name to be monitoring | infra-event-queue |


 ### Set up and run a Queue Listener 
 ##### Build 
  ```
 $ go build queueListener.go
 ```

 ##### Running
 ``` 
 ./queueListener -queue=infra-event-queue -profile=foo -region=us-gov-west-1 -getEndpoint=http://google.com -putEndpoint=http://foo.io
 ```

##### Docker
Building the docker image:
```
docker build . -t listener:latest
```
 ##### Running in Docker 
 The AWS SDK uses the AWS Cli credentials to interface with AWS thus they need to be provided. We 
 accomplish this by attaching a read only volume to the container.
```
 docker run -v $HOME/.aws/credentials:/root/.aws/credentials:ro --env QUEUE=infra-event-queue --env REGION=default listener:latest
```
  
  Parameter Options
 
 | Parameter Name | Description | Default Value |
|---|---|---|
| profile | sets the AWS Cli profile to be used for accessing AWS SDK  | default |
| region | sets the region to operate one | none (required input) |
| queue | sets the SQS queue name to be monitoring | none (required input) |
| getEndpoint | establishes url for http GET call to be made on event received | none (optional input) |
| putEndpoint | establishes url for http PUT call to be made on event received | none (optional input) |
 
To Do:
Add threshold 
Add parameter time interval 
JSON configuration 
Interrogate aws credential to obtain region if no region is provided 
 Reference: https://dev.to/jeastham1993/how-to-use-amazon-sqs-and-sns-for-inter-service-communication-part-2-2pna

##### GitHub Actions

##### Subdirectories 
The following directories in contain other, not directly related, projects/applications.
These were written either to gain an understanding or to test some aspect of the queue listener
app.

- httpTesters
- lambdaClient
- testMessage
- topicAndQueueSetup