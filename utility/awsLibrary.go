package utility

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/defaults"
    "github.com/aws/aws-sdk-go/aws/endpoints"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sns"
    "github.com/aws/aws-sdk-go/service/sqs"
    "strings"

    "fmt"
    "os"
)

func GetSession() *session.Session {
    // Initialize a session that the SDK will use to load
    // credentials from the shared credentials file. (~/.aws/credentials).
    sess := session.Must(
        session.NewSession(&aws.Config{
            // Use the SDK's SharedCredentialsProvider directly instead of the
            // SDK's default credential chain. This ensures that the
            // application can call Config.Credentials.Expire. This  is counter
            // to the SDK's default credentials chain, which  will never reread
            // the shared credentials file.
            Credentials: credentials.NewCredentials(&credentials.SharedCredentialsProvider{
                Filename: defaults.SharedCredentialsFilename(),
                Profile:  "default",
            }),
            Region: aws.String(endpoints.UsEast1RegionID),
        }),
    )

    return sess
}

func CreateTopic(topic string, svc *sns.SNS) string {
    result, err := svc.CreateTopic(&sns.CreateTopicInput{
        Name: aws.String(topic),
    })
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    fmt.Println("Topic has been created ARN is:: " + *result.TopicArn)
    return *result.TopicArn
}

func TopicExists(requiredTopic string, svc *sns.SNS) bool {
    resultT, errT := svc.ListTopics(nil)
    if errT != nil {
        fmt.Println(errT.Error())
        os.Exit(1)
    }

    for _, t := range resultT.Topics {
        if strings.Contains(*t.TopicArn, requiredTopic) {
            fmt.Println("Topic exists, return ARN")
            return true
        }
    }
    fmt.Println("Topic not found.")
    return false
}

func RetrieveTopicArn(requiredTopic string, svc *sns.SNS) string {
    resultT, errT := svc.ListTopics(nil)
    if errT != nil {
        fmt.Println(errT.Error())
        os.Exit(1)
    }

    var arn = ""
    for _, t := range resultT.Topics {
        if strings.Contains(*t.TopicArn, requiredTopic) {
            fmt.Println("Topic exists, return ARN")
            arn = *t.TopicArn
            break
        }
    }
    return arn
}

func SendTestMessage(message, arn string, svc *sns.SNS) *sns.PublishInput {
    input := &sns.PublishInput{
        Message:  aws.String(message),
        TopicArn: aws.String(arn),
    }

    result, err := svc.Publish(input)

    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    fmt.Println(*result.MessageId)
    return input
}

func ConvertQueueURLToARN(inputURL string) string {
    // Awfully bad string replace code to convert a SQS queue URL to an ARN
    queueARN := strings.Replace(strings.Replace(strings.Replace(inputURL, "https://sqs.", "arn:aws:sqs:", -1), ".amazonaws.com/", ":", -1), "/", ":", -1)
    return queueARN
}

func RetrieveQueueURL(sqsSvc *sqs.SQS, requiredQueueName string) string {
    listQueuesRequest := sqs.ListQueuesInput{}
    listQueueResults, _ := sqsSvc.ListQueues(&listQueuesRequest)
    queueURL := ""
    for _, t := range listQueueResults.QueueUrls {
        fmt.Println(*t)
        // If one of the returned queue URL's contains the required name we need
        // then break the loop
        if strings.Contains(*t, requiredQueueName) {
            fmt.Println("Queue exists.")
            queueURL = *t
            break
        }
    }
    return queueURL
}

