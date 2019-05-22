package awsapi

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

const (
	// AsyncInvokeEvent indicates invocation type = Event
	AsyncInvokeEvent = "Event"
	// AsyncInvokeRequest indicates invocation type = RequestResponse
	AsyncInvokeRequest = "RequestResponse"

	// ActionInvoke indicates action to invoke function
	ActionInvoke = "lambda:InvokeFunction"
	// PrincipalEvents indicates principal of cloudwatchevents to invoke lambda
	PrincipalEvents = "events.amazonaws.com"
)

// InvokeLambda invoke a lambda function
func (awsapi *AWSAPI) InvokeLambda(lambdasvc lambdaiface.LambdaAPI,
	funcArn string, payload []byte, invokeType string) (*lambda.InvokeOutput, error) {

	input := lambda.InvokeInput{
		FunctionName:   aws.String(funcArn),
		Payload:        payload,
		InvocationType: aws.String(invokeType),
	}
	output, err := lambdasvc.Invoke(&input)
	if err != nil {
		log.Printf("AWSError: fail to invoke lambda %s input=(%s) invocation=(%s), err: %v",
			funcArn, string(payload), invokeType, err)
		return nil, err
	}
	log.Printf("Lambda Function invoked: arn %s input=(%s) invocation=(%s), output=(%+v)",
		funcArn, string(payload), invokeType, output)
	return output, nil
}

// LambdaAddPermission helps to call lambda.AddPermission api
func (awsapi *AWSAPI) LambdaAddPermission(lambdasvc lambdaiface.LambdaAPI,
	action, funcName, principal, source, stateID string) error {
	input := &lambda.AddPermissionInput{
		Action:       aws.String(action),
		FunctionName: aws.String(funcName),
		Principal:    aws.String(principal),
		SourceArn:    aws.String(source),
		StatementId:  aws.String(stateID),
	}
	output, err := lambdasvc.AddPermission(input)
	if err != nil {
		log.Printf("AWSError: fail to add permission to lambda %s "+
			"w/ (source=%s, action=%s, principal=%s, stateID=%s)", funcName, source,
			action, principal, stateID)
		return err
	}

	/* statement example:
	 * `{"Sid":"randstateid12345","Effect":"Allow",
	 * "Principal":{"Service":"events.amazonaws.com"},
	 * "Action":"lambda:InvokeFunction","Resource":"arn:aws:lambda:
	 * ap-northeast-1:054657590879:function:test_gui_sqs_receiver",
	 * "Condition":{"ArnLike":{"AWS:SourceArn":"arn:aws:events:ap-northeast-1:
	 * 054657590879:rule/lambdaname-randid"}}}) `
	 */
	log.Printf("Succeed to add permission to lambda %s", *output.Statement)
	return nil
}

// NewLambdaAPI returns a Lambda API
func NewLambdaAPI(sess AWSSession) lambdaiface.LambdaAPI {
	lambdasess, err := sess.NewSession()
	if err != nil {
		log.Printf("AWSError: Cannot create aws session, err %#v", err)
		return nil
	}
	return lambda.New(lambdasess)
}
