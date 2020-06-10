package main

import (
    "fmt"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sqs"
)

// Usage:
// go run sqs_receive_message.go
func main() {
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    svc := sqs.New(sess)

    // URL to our queue
    qURL := "https://sqs.us-east-1.amazonaws.com/085141918894/TEST_QUEUE_NAME"

    result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
        AttributeNames: []*string{
            aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
        },
        MessageAttributeNames: []*string{
            aws.String(sqs.QueueAttributeNameAll),
        },
        QueueUrl:            &qURL,
        MaxNumberOfMessages: aws.Int64(10),
        VisibilityTimeout:   aws.Int64(60), // 60 seconds
        WaitTimeSeconds:     aws.Int64(0),

    })
    if err != nil {
        fmt.Println("Error", err)
        return
    }
    if len(result.Messages) == 0 {
        fmt.Println("Received no messages")
        return
    }

    fmt.Printf("Success: %+v\n", result.Messages)

    resultDelete, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
        QueueUrl:      &qURL,
        ReceiptHandle: result.Messages[0].ReceiptHandle,
    })

    if err != nil {
        fmt.Println("Delete Error", err)
        return
    }

    fmt.Println("Message Deleted", resultDelete)
}
// snippet-end:[sqs.go.receive_message]
