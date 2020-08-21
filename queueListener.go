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

/**
This will listen for events from the queue and then, if configured, it will make calls
to endpoints to propagate the eventing.
 */
func main() {
	//infra-event-queue
	queuePtr := flag.String("queue", "", "a string")
	profilePtr := flag.String("profile", "default", "a string")
	regionPtr := flag.String("region", "", "a string")
	getEndpointPtr := flag.String("getEndpoint", "", "a string")
	putEndpointPtr := flag.String("putEndpoint", "", "a string")
	thresholdPtr := flag.Int("threshold", 10, "an int")
	flag.Parse()

	if *queuePtr == "" {
		fmt.Println("Queue name is a required parameter, set it and retry")
		os.Exit(-1)
	}

	if *regionPtr == "" {
		fmt.Println("Region is a required parameter, set it and retry")
		os.Exit(-1)
	}

	fmt.Printf("Proceeding with \n queue=%s\n profile=%s\n region=%s\n",
		*queuePtr, *profilePtr, *regionPtr)

	if *getEndpointPtr == "" {
		fmt.Println("HTTP GET endpoint has not been set, no action will be taken")
	} else {
		fmt.Printf("Using HTTP GET endpoint %s as the action endpoint.\n", *getEndpointPtr)
	}

	if *putEndpointPtr == "" {
		fmt.Println("HTTP POST endpoint has not been set, no action will be taken")
	} else {
		fmt.Printf("Using HTTP POST endpoint %s as the action endpoint.\n", *putEndpointPtr)
	}

	// Initialize the AWS session
	sess := utility.GetSession(*profilePtr, *regionPtr)

	if sess == nil {
		fmt.Println("Session was not obtained, exiting.")
	}

	// Create new services for SQS and SNS
	sqsSvc := sqs.New(sess)

	if sqsSvc == nil {
		fmt.Println("SQS Session was not obtained, exiting.")
		os.Exit(-1)
	}

	requiredQueueName := *queuePtr

	queueURL := utility.RetrieveQueueURL(sqsSvc, requiredQueueName)
	if queueURL == "" {
		fmt.Printf("The specified queue %s was not found, exiting now\n", requiredQueueName)
		os.Exit(-1)
	}
	go checkMessages(*sqsSvc, queueURL, *getEndpointPtr, *putEndpointPtr, *thresholdPtr)

	_, _ = fmt.Scanln()
}

func checkMessages(sqsSvc sqs.SQS, queueURL string, getEndpoint string, putEndpoint string, threshold int) {
	var count int = 0
	for ; ; {
	   retrieveMessageRequest := sqs.ReceiveMessageInput{
   	      QueueUrl: &queueURL,
	   }

	   retrieveMessageResponse, _ := sqsSvc.ReceiveMessage(&retrieveMessageRequest)

	   if len(retrieveMessageResponse.Messages) > 0 {
              count++
        
              if count >= threshold {
                 fmt.Println("Notification threshold has been reached, call endpoints.")
	         if getEndpoint != "" {
	            callGetEndpoint(getEndpoint)
	         }

	        if putEndpoint != "" {
	           callPutEndpoint(putEndpoint)
	        }

	        fmt.Println("Resetting threshold count")
	        count = 0
	      } 
	      cleanupMessages(sqsSvc, retrieveMessageResponse, queueURL)
	   }

	   if len(retrieveMessageResponse.Messages) == 0 {
	      fmt.Println(":(  I have no messages")
   	   }

	   fmt.Printf("%v+\n", time.Now())
	   time.Sleep(time.Minute)
     }
}

func cleanupMessages(sqsSvc sqs.SQS, retrieveMessageResponse *sqs.ReceiveMessageOutput, queueURL string) {
	fmt.Println("Cleaning up messages from the queue.")
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
			fmt.Println(":(  I have no messages, will check again momentarily")
		}

		fmt.Printf("%v+\n", time.Now())
		time.Sleep(time.Second * 30)
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
