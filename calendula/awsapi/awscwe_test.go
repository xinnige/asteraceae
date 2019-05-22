package awsapi

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xinnige/asteraceae/calendula/mock"
)

func TestNewCWEAPI(t *testing.T) {
	sess := &AWSServiceSession{}
	cweapi := NewCloudWatchEventsAPI(sess)
	assert.NotNil(t, cweapi)
}

func TestNewCWEAPIError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSession := mock.NewMockAWSSession(mockCtrl)
	mockSession.EXPECT().NewSession().Return(
		nil, errors.New("FakeNewSessionError")).Times(1)

	cweapi := NewCloudWatchEventsAPI(mockSession)
	assert.Nil(t, cweapi)
}

func TestNewEventRule(t *testing.T) {
	rule := NewEventRule(aws.String("arn"), aws.String("description"),
		aws.String("event"), aws.String("managedby"), aws.String("fakename"),
		aws.String("role"), aws.String("schedule"), aws.String("state"))
	assert.Equal(t, "fakename", rule.Name)
}

func TestListRuleNames(t *testing.T) {
	awsapi := AWSAPI{}
	cwesvc := mock.CWEAPIListRules([]string{"test-01", "test-02"}, []string{"123", ""})

	output, err := awsapi.ListRuleNamesByTarget(cwesvc, "fake-arn")
	assert.Equal(t, []string{"test-01", "test-02", "test-01", "test-02"}, output)
	assert.Nil(t, err)

	cwesvc = mock.CWEAPIListRulesError()
	output, err = awsapi.ListRuleNamesByTarget(cwesvc, "fake-arn")
	assert.Nil(t, output)
	assert.Equal(t, "FakeListRulesError", err.Error())

	cwesvc = mock.CWEAPIListRulesLoopError()
	output, err = awsapi.ListRuleNamesByTarget(cwesvc, "fake-arn")
	assert.Equal(t, []string{"rule01"}, output)
	assert.Equal(t, "FakeListRulesError", err.Error())
}

func TestDescribeRule(t *testing.T) {
	awsapi := AWSAPI{}
	cwesvc := mock.CWEAPIDescribeRule("fake-arn", "fake-name", "cron(0 23 * * ? *)", "Enabled")

	output, err := awsapi.DescribeRule(cwesvc, "fake-arn")
	assert.Equal(t, "cron(0 23 * * ? *)", output.ScheduleExpression)
	assert.Nil(t, err)

	cwesvc = mock.CWEAPIDescribeRuleError()
	output, err = awsapi.DescribeRule(cwesvc, "fake-arn")
	assert.Equal(t, "FakeDescribeRuleError", err.Error())
	assert.Nil(t, output)
}

func TestPutRule(t *testing.T) {
	awsapi := AWSAPI{}
	cwesvc := mock.CWEAPIPutRule("fake-arn")

	rule := &CloudWatchEventRule{
		Name: "fake-name", ScheduleExpression: "fake-cron", State: "Enabled",
		Description: "fake-description", EventPattern: "fake-event", RoleArn: "fake-role",
	}
	output, err := awsapi.PutRule(cwesvc, rule)
	assert.Equal(t, "fake-arn", output)
	assert.Nil(t, err)

	cwesvc = mock.CWEAPIPutRuleError()
	output, err = awsapi.PutRule(cwesvc, rule)
	assert.Equal(t, "FakePutRuleError", err.Error())
	assert.Equal(t, "", output)
}

func TestPutTargets(t *testing.T) {
	awsapi := AWSAPI{}
	cwesvc := mock.CWEAPIPutTargets()

	err := awsapi.PutTargets(cwesvc, "rule-arn", "target-id", "target-arn", "{target-nput}")
	assert.Nil(t, err)

	cwesvc = mock.CWEAPIPutTargetsError()
	err = awsapi.PutTargets(cwesvc, "rule-arn", "target-id", "target-arn", "{target-nput}")
	assert.Equal(t, "FakePutTargetsError", err.Error())
}
