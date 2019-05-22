package awsapi

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	cwe "github.com/aws/aws-sdk-go/service/cloudwatchevents"
	cweiface "github.com/aws/aws-sdk-go/service/cloudwatchevents/cloudwatcheventsiface"
)

const (
	size = 100
)

// CloudWatchEventRule defines necessary properties of a cloudwatch event rule
type CloudWatchEventRule struct {
	Arn                string `json:"arn"`
	Description        string `json:"description"`
	EventPattern       string `json:"event_pattern"`
	ManagedBy          string `json:"managed_by"`
	Name               string `json:"name"`
	RoleArn            string `json:"role_arn"`
	ScheduleExpression string `json:"schedule_expression"`
	State              string `json:"state"`
}

// NewEventRule parse cloudwatchevents.DescribeRuleOutput to CloudWatchEventRule
func NewEventRule(arn, desc, event, managedby,
	name, role, schedule, state *string) *CloudWatchEventRule {
	rule := &CloudWatchEventRule{}
	if arn != nil {
		rule.Arn = aws.StringValue(arn)
	}
	if desc != nil {
		rule.Description = aws.StringValue(desc)
	}
	if event != nil {
		rule.EventPattern = aws.StringValue(event)
	}
	if managedby != nil {
		rule.ManagedBy = aws.StringValue(managedby)
	}
	if name != nil {
		rule.Name = aws.StringValue(name)
	}
	if role != nil {
		rule.RoleArn = aws.StringValue(role)
	}
	if name != nil {
		rule.ScheduleExpression = aws.StringValue(schedule)
	}
	if state != nil {
		rule.State = aws.StringValue(state)
	}
	return rule

}

// NewCloudWatchEventsAPI returns a CloudWatch API
func NewCloudWatchEventsAPI(sess AWSSession) cweiface.CloudWatchEventsAPI {
	cwesess, err := sess.NewSession()
	if err != nil {
		log.Printf("AWSError: Cannot create aws session, err %#v", err)
		return nil
	}
	return cwe.New(cwesess)
}

// ListRuleNamesByTarget returns a list of rule names related to a target arn
func (awsapi *AWSAPI) ListRuleNamesByTarget(svc cweiface.CloudWatchEventsAPI,
	arn string) ([]string, error) {
	ruleNames := make([]string, 0)
	input := &cwe.ListRuleNamesByTargetInput{
		TargetArn: aws.String(arn),
		Limit:     aws.Int64(size),
	}
	output, err := svc.ListRuleNamesByTarget(input)
	if err != nil {
		log.Printf("AWSError: Cannot list rules of %s, err %v", arn, err)
		return nil, err
	}
	ruleNames = append(ruleNames, aws.StringValueSlice(output.RuleNames)...)
	for aws.StringValue(output.NextToken) != "" {
		input.NextToken = output.NextToken
		output, err := svc.ListRuleNamesByTarget(input)
		if err != nil {
			log.Printf("AWSError: Cannot list rules of %s, err %v", arn, err)
			return ruleNames, err
		}
		ruleNames = append(ruleNames, aws.StringValueSlice(output.RuleNames)...)
	}
	return ruleNames, nil
}

// DescribeRule returns an expression of a rule scheule searchd by name
func (awsapi *AWSAPI) DescribeRule(svc cweiface.CloudWatchEventsAPI,
	name string) (*CloudWatchEventRule, error) {
	input := &cwe.DescribeRuleInput{
		Name: aws.String(name),
	}
	output, err := svc.DescribeRule(input)
	if err != nil {
		log.Printf("AWSError: Get schedule expression of rule %s, err %v", name, err)
		return nil, err
	}
	return NewEventRule(output.Arn, output.Description, output.EventPattern,
		output.ManagedBy, output.Name, output.RoleArn, output.ScheduleExpression,
		output.State), nil
}

// PutRule helps to call PutRule api
func (awsapi *AWSAPI) PutRule(svc cweiface.CloudWatchEventsAPI,
	rule *CloudWatchEventRule) (string, error) {
	input := &cwe.PutRuleInput{
		Name:               aws.String(rule.Name),
		ScheduleExpression: aws.String(rule.ScheduleExpression),
		State:              aws.String(rule.State),
	}
	if len(rule.RoleArn) != 0 {
		input.RoleArn = aws.String(rule.RoleArn)
	}
	if len(rule.Description) != 0 {
		input.Description = aws.String(rule.Description)
	}
	if len(rule.EventPattern) != 0 {
		input.EventPattern = aws.String(rule.EventPattern)
	}

	output, err := svc.PutRule(input)
	if err != nil {
		log.Printf("AWSError: Cannot put rule %+v, err %v", input, err)
		return "", err
	}
	return aws.StringValue(output.RuleArn), nil
}

// PutTargets helps to call PutTargets api
func (awsapi *AWSAPI) PutTargets(svc cweiface.CloudWatchEventsAPI,
	rule, targetID, targetArn, targetInput string) error {
	target := &cwe.Target{
		Arn:   aws.String(targetArn),
		Id:    aws.String(targetID),
		Input: aws.String(targetInput)}
	input := &cwe.PutTargetsInput{
		Rule:    aws.String(rule),
		Targets: []*cwe.Target{target},
	}
	output, err := svc.PutTargets(input)
	if err != nil {
		log.Printf("AWSError: Cannot put targets (%s w/ %v), output=(%v), err %+v", rule, target, output, err)
		return err
	}
	return nil
}
