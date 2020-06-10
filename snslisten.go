package main

import (
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sns"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/defaults"
    "github.com/aws/aws-sdk-go/aws/endpoints"

    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("You must supply a topic name")
        fmt.Println("Usage: go run awsmonitor.go TOPIC")
        os.Exit(1)
    }

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

    svc := sns.New(sess)
 
    if os.Args[1] != "none" {
       result, err := svc.CreateTopic(&sns.CreateTopicInput{
          Name: aws.String(os.Args[1]),
       })
       if err != nil {
         fmt.Println(err.Error())
          os.Exit(1)
       }

       fmt.Println(*result.TopicArn)
    } else {
      fmt.Println("Not creating a new topic")
    } 


    resultT, errT := svc.ListTopics(nil)
    if errT != nil {
        fmt.Println(errT.Error())
        os.Exit(1)
    }

  
    var arn = ""
    for _, t := range resultT.Topics {
        fmt.Println(*t.TopicArn)
        arn = *t.TopicArn
    }

    //publish test message
    input := &sns.PublishInput{
        Message:  aws.String("Hello world!"),
        TopicArn: aws.String(arn),
    }

    result, err := svc.Publish(input)
    if err != nil {
        fmt.Println("Publish error:", err)
        return
    }

    fmt.Println(result)
}
