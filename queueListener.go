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
	putEndpointPtr := flag.String("putEndpoint", "", "a string")
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

	if *putEndpointPtr == "" {
		fmt.Println("HTTP POST endpoint has not been set, no action will be taken")
	} else {
		fmt.Printf("Using HTTP POST endpoint %s as the action endpoint", *putEndpointPtr)
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
	go checkMessages(*sqsSvc, queueURL, *getEndpointPtr, *putEndpointPtr)

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

func callPutEndpoint(endpoint string)  {
	req, err := http.NewRequest("PUT", endpoint, nil)
	if err != nil {
		fmt.Errorf("Error creating HTTP PUT request %s", err.Error())
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Errorf("Error completing HTTP PUT %s", err.Error())
		return
	}
	defer resp.Body.Close()
	fmt.Printf("HTTP Post to %s has been completed.\n Response status: %s\n", endpoint, resp.Status)
}
