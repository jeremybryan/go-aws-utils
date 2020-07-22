package main

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"goutilities.com/awsutil/utility"
)
func main() {
	requiredTopic := "infrastructure-event"
	sess := utility.GetSession("default", "us-east-1")

	snsSvc := sns.New(sess)

	topicArn := utility.RetrieveTopicArn(requiredTopic, snsSvc)
	utility.SendTestMessage("Hi There...", topicArn, snsSvc)
}
