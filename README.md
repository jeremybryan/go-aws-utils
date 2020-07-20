 This project contains a set of utilities for interfacing with AWS notification and queue services. Using the AWS CLI
 and SDK a you can configure a AWS Simple Notification Service (SNS) topic and AWS Simple Queue Service (SQS) queue and 
 then listen for events on the queue. The original insipiration for this was to use this as part of a proof-of-concept 
 demonstrating monitoring of infrastructure change events (via CloudTrail and/or CloudWatch).
 
 ### Set up a Topic and Queue 
 ##### Build 
 ```
$ go build setupTopicAndQueue.go
```
 
 ##### Running
 ``` 
 $  ./setupTopicAndQueue -profile=foo -region=us-gov-west-1
 ```
 Available Parameters
 
 ### Set up and run a Queue Listener 
 ##### Build 
  ```
 $ go build queueListener.go
 ```

 ##### Running
 ``` 
 ./queueListener -queue=infra-event-queue -profile=foo -region=us-gov-west-1 -endpoint=http://google.com
 ```
  
 
 
 
 Reference: https://dev.to/jeastham1993/how-to-use-amazon-sqs-and-sns-for-inter-service-communication-part-2-2pna
