package main

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"goutilities.com/awsutil/utility"
)
/**
Simple test class for sending event to mimic infrastructure events to the
topic
 */
func main() {
	requiredTopic := "infrastructure-event"
	sess := utility.GetSession("default", "us-east-1")

	snsSvc := sns.New(sess)

	topicArn := utility.RetrieveTopicArn(requiredTopic, snsSvc)
	utility.SendTestMessage("Hi There...", topicArn, snsSvc)
}
