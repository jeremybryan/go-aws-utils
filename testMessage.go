package main

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"goutilities.com/snsutil/utility"
)
func main() {
	requiredTopic := "infrastructure-event"
	sess := utility.GetSession("tapestry")

	snsSvc := sns.New(sess)

	topicArn := utility.RetrieveTopicArn(requiredTopic, snsSvc)
	utility.SendTestMessage("Hi There...", topicArn, snsSvc)
}
