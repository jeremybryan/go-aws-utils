package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"goutilities.com/awsutil/utility"
	"net/http"
	"time"
        "os"
)

func main() {
	//infra-event-queue
	queuePtr := flag.String("queue", "", "a string")
	profilePtr := flag.String("profile", "default", "a string")
	regionPtr := flag.String("region", "us-east-1", "a string")
	endpointPtr := flag.String("endpoint", "", "a string")
	flag.Parse()

	if *queuePtr == "" {
		fmt.Println("Queue name is a required parameter, set it and retry")
		os.Exit(-1)
	}

	fmt.Printf("Proceeding with \n queue=%s\n profile=%s\n and region=%s.\n",
		*queuePtr, *profilePtr, *regionPtr,)

	if *endpointPtr == "" {
		fmt.Println("HTTP endpoint has not been set, no action will be taken")
	} else {
		fmt.Printf("Using HTTP endpoint %s as the action endpoint", *endpointPtr)
	}

	// Initialize the AWS session
	sess := utility.GetSession(*profilePtr, *regionPtr)

	// Create new services for SQS and SNS
	sqsSvc := sqs.New(sess)

	requiredQueueName := *queuePtr

	queueURL := utility.RetrieveQueueURL(sqsSvc, requiredQueueName)
	if queueURL == "" {
		fmt.Printf("The specified queue %s was not found, exiting now\n", requiredQueueName)
		os.Exit(-1)
	}
	go checkMessages(*sqsSvc, queueURL, *endpointPtr)

	_, _ = fmt.Scanln()
}

func checkMessages(sqsSvc sqs.SQS, queueURL string, endpoint string) {
	for ; ; {
		retrieveMessageRequest := sqs.ReceiveMessageInput{
			QueueUrl: &queueURL,
		}

		retrieveMessageResponse, _ := sqsSvc.ReceiveMessage(&retrieveMessageRequest)

		if len(retrieveMessageResponse.Messages) > 0 {

			processedReceiptHandles := make([]*sqs.DeleteMessageBatchRequestEntry, len(retrieveMessageResponse.Messages))

			for i, mess := range retrieveMessageResponse.Messages {
				fmt.Println(mess.String())

				if endpoint != "" {
					callEndpoint(endpoint)
				}

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

func callEndpoint(endpoint string)  {
	resp, err := http.Get(endpoint)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Response status:", resp.Status)
}
