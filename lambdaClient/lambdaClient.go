package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	invoke "github.com/aws/aws-sdk-go/service/lambda"
)

/**
Simple lambda integration
 */
func main() {
	//region := os.Getenv("AWS_REGION")
	region := "us-east-1"
	sess, err := session.NewSession(&aws.Config{ // Use aws sdk to connect to dynamoDB
		Region: &region,
	})
	svc := invoke.New(sess)

	data := make(map[string]interface{})
	data["repository"] = map[string]interface{} {
		"name" : "jimmy",
	}
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	type Payload struct {
		// You can also include more objects in the structure like below,
		// but for my purposes body was all that was required
		// Method string `json:"httpMethod"`
		Body string `json:"body"`
	}
	p := Payload{
		// Method: "POST",
		Body: string(body),
	}
	payload, err := json.Marshal(p)
	// Result should be: {"body":"{\"name\":\"Jimmy\"}"}
	// This is the required format for the lambda request body.

	if err != nil {
		fmt.Println("Json Marshalling error")
	}
	fmt.Println(string(payload))

	input := &invoke.InvokeInput{
		FunctionName:   aws.String("github-webhook"),
		InvocationType: aws.String("RequestResponse"),
		LogType:        aws.String("Tail"),
		Payload:        body,
	}
	result, err := svc.Invoke(input)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err.Error())
	}

	var m map[string]interface{}
	json.Unmarshal(result.Payload, &m)
}


