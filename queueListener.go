package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"goutilities.com/snsutil/utility"
	"time"
)

func main() {
	// Initialize the AWS session
	sess := utility.GetSession()

	// Create new services for SQS and SNS
	sqsSvc := sqs.New(sess)

	requiredQueueName := "infra-event-queue"

	queueURL := utility.RetrieveQueueURL(sqsSvc, requiredQueueName)

	go checkMessages(*sqsSvc, queueURL)

	_, _ = fmt.Scanln()
}

func checkMessages(sqsSvc sqs.SQS, queueURL string) {
	for ; ; {
		retrieveMessageRequest := sqs.ReceiveMessageInput{
			QueueUrl: &queueURL,
		}

		retrieveMessageResponse, _ := sqsSvc.ReceiveMessage(&retrieveMessageRequest)

		if len(retrieveMessageResponse.Messages) > 0 {

			processedReceiptHandles := make([]*sqs.DeleteMessageBatchRequestEntry, len(retrieveMessageResponse.Messages))

			for i, mess := range retrieveMessageResponse.Messages {
				fmt.Println(mess.String())

				processedReceiptHandles[i] = &sqs.DeleteMessageBatchRequestEntry{
					Id: mess.MessageId,
					ReceiptHandle: mess.ReceiptHandle,
				}
			}

			deleteMessageRequest := sqs.DeleteMessageBatchInput{
				QueueUrl: &queueURL,
				Entries: processedReceiptHandles,
			}

			_,err := sqsSvc.DeleteMessageBatch(&deleteMessageRequest)

			if err != nil {
				fmt.Println(err.Error())
			}
		}

		if len(retrieveMessageResponse.Messages) == 0 {
			fmt.Println(":(  I have no messages")
		}

		fmt.Printf("%v+\n", time.Now())
		time.Sleep(time.Minute)
	}
}