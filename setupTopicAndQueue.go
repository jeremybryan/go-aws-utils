package main

import (
    "flag"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/sns"
    "github.com/aws/aws-sdk-go/service/sqs"
    "goutilities.com/awsutil/utility"
    "strings"
    "fmt"
)

func main() {
    profilePtr := flag.String("profile", "default", "a string")
    regionPtr := flag.String("region", "us-east-1", "a string")
    queuePtr := flag.String("queue", "infra-event-queue", "a string")
    topicPtr := flag.String("topic", "infrastructure-event", "a string")

    flag.Parse()

    requiredQueueName := *queuePtr
    requiredTopic := *topicPtr
    queueURL := ""
    topicArn := ""
    protocolName := "sqs"
    regionType := ""

    //Set the region type (gov or commercial)
    if strings.Contains(*regionPtr, "gov") {
        regionType = "gov"
    }

    fmt.Printf("Proceeding with \n queue=%s\n profile=%s\n topic=%s\n region=%s\n",
        *queuePtr, *profilePtr, *topicPtr, *regionPtr)

    //Get a session for interfacing with AWS
    sess := utility.GetSession(*profilePtr, *regionPtr)

    // Create an SQS and SNS service client.
    snsSvc := sns.New(sess)
    sqsSvc := sqs.New(sess)

    //First let's check to see if the required topic exists
    if !utility.TopicExists(requiredTopic, snsSvc) {
        fmt.Println("Required topic doesn't exist, creating it now.")
        //create it
        //we could grab the ARN from here
        topicArn = utility.CreateTopic(requiredTopic, snsSvc)
    }

    //Get the ARN for the topic
    if topicArn == "" {
        topicArn = utility.RetrieveTopicArn(requiredTopic, snsSvc)
    }

    // Establish queue
    //First, get a list of all queues and see if our queue exists
    queueURL = utility.RetrieveQueueURL(sqsSvc, requiredQueueName)

    //Need to create the queue if it doesn't exist
    if queueURL == "" {
        result, err := sqsSvc.CreateQueue(&sqs.CreateQueueInput{
            QueueName: aws.String(requiredQueueName),
            Attributes: map[string]*string{
                "DelaySeconds":           aws.String("15"),
                "MessageRetentionPeriod": aws.String("86400"),
            },
        })
        if err != nil {
            fmt.Println("Error", err)
            return
        }
        fmt.Println("Success", *result.QueueUrl)
        queueURL = *result.QueueUrl
    }

    //Now that we know the queue exists...we need to register it to listen to the topic

    // No way to retrieve the queue ARN through the SDK, manual string replace to generate the ARN
    queueARN := utility.ConvertQueueURLToARN(queueURL, regionType)

    fmt.Println("Topic URN", topicArn)
    fmt.Println("Protocol Name", protocolName)
    fmt.Println("Queue ARN", queueARN)

    if topicArn != "" {
        subscribeQueueInput := sns.SubscribeInput{
            TopicArn: &topicArn,
            Protocol: &protocolName,
            Endpoint: &queueARN,
        }

        createSubRes, err := snsSvc.Subscribe(&subscribeQueueInput)

        if err != nil {
            fmt.Println(err.Error())
        }

        if createSubRes != nil {
            fmt.Println(*createSubRes.SubscriptionArn)
        }
    }

    policyContent := "{\"Version\": \"2012-10-17\",  \"Id\": \"" + queueARN + "/SQSDefaultPolicy\",  \"Statement\": [    {     \"Sid\": \"Sid1580665629194\",      \"Effect\": \"Allow\",      \"Principal\": {        \"AWS\": \"*\"      },      \"Action\": \"SQS:SendMessage\",      \"Resource\": \"" + queueARN + "\",      \"Condition\": {        \"ArnEquals\": {         \"aws:SourceArn\": \"" + topicArn + "\"        }      }    }  ]}"

    attr := make(map[string]*string, 1)
    attr["Policy"] = &policyContent

    setQueueAttrInput := sqs.SetQueueAttributesInput{
        QueueUrl:   &queueURL,
        Attributes: attr,
    }

    var _, err = sqsSvc.SetQueueAttributes(&setQueueAttrInput)

    if err != nil {
        fmt.Println(err.Error())
    }
}
