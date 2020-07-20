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
	getEndpointPtr := flag.String("getEndpoint", "", "a string")
	postEndpointPtr := flag.String("postEndpoint", "", "a string")
	flag.Parse()

	if *queuePtr == "" {
		fmt.Println("Queue name is a required parameter, set it and retry")
		os.Exit(-1)
	}

	fmt.Printf("Proceeding with \n queue=%s\n profile=%s\n and region=%s.\n",
		*queuePtr, *profilePtr, *regionPtr,)

	if *getEndpointPtr == "" {
		fmt.Println("HTTP GET endpoint has not been set, no action will be taken")
	} else {
		fmt.Printf("Using HTTP GET endpoint %s as the action endpoint", *getEndpointPtr)
	}

	if *postEndpointPtr == "" {
		fmt.Println("HTTP POST endpoint has not been set, no action will be taken")
	} else {
		fmt.Printf("Using HTTP POST endpoint %s as the action endpoint", *postEndpointPtr)
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
	go checkMessages(*sqsSvc, queueURL, *getEndpointPtr, *postEndpointPtr)

	_, _ = fmt.Scanln()
}

func checkMessages(sqsSvc sqs.SQS, queueURL string, getEndpoint string, postEndpoint string) {
	for ; ; {
		retrieveMessageRequest := sqs.ReceiveMessageInput{
			QueueUrl: &queueURL,
		}

		retrieveMessageResponse, _ := sqsSvc.ReceiveMessage(&retrieveMessageRequest)

		if len(retrieveMessageResponse.Messages) > 0 {

			processedReceiptHandles := make([]*sqs.DeleteMessageBatchRequestEntry, len(retrieveMessageResponse.Messages))

			for i, mess := range retrieveMessageResponse.Messages {
				fmt.Println(mess.String())

				if getEndpoint != "" {
					callGetEndpoint(getEndpoint)
				}
				if postEndpoint != "" {
					callPostEndpoint(postEndpoint)
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

func callGetEndpoint(endpoint string)  {
	resp, err := http.Get(endpoint)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Response status:", resp.Status)
}

func callPostEndpoint(endpoint string)  {
	resp, err := http.Post(endpoint, "application/json", nil)
	if err != nil {
		fmt.Errorf("Error completing HTTP POST %s", err.Error())
	}
	defer resp.Body.Close()
	fmt.Printf("HTTP Post to %s has been completedl Response status: %s", endpoint, resp.Status)
}
