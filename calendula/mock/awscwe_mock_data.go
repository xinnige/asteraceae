package mock

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	cwe "github.com/aws/aws-sdk-go/service/cloudwatchevents"
	cweiface "github.com/aws/aws-sdk-go/service/cloudwatchevents/cloudwatcheventsiface"
)

// CWEAPI mocks s3iface.S3API
type CWEAPI struct {
	cweiface.CloudWatchEventsAPI
	listRulesOutput    *cwe.ListRuleNamesByTargetOutput
	listRulesError     error
	listIndex          *int
	listNexts          []string
	listErrors         []string
	describeRuleOutput *cwe.DescribeRuleOutput
	describeRuleError  error
	putRuleOutput      *cwe.PutRuleOutput
	putRuleError       error
	putTargetsOutput   *cwe.PutTargetsOutput
	putTargetsError    error
}

// ListRuleNamesByTarget mocks cloudwatchevents.ListRuleNamesByTarget
func (cweapi *CWEAPI) ListRuleNamesByTarget(
	*cwe.ListRuleNamesByTargetInput) (*cwe.ListRuleNamesByTargetOutput, error) {
	if cweapi.listRulesError != nil {
		return cweapi.listRulesOutput, cweapi.listRulesError
	}
	// returns error if defined not ""
	if *cweapi.listIndex < len(cweapi.listErrors) {
		if len(cweapi.listErrors[*cweapi.listIndex]) != 0 {
			return cweapi.listRulesOutput,
				fmt.Errorf("%s", cweapi.listErrors[*cweapi.listIndex])
		}
	}

	if *cweapi.listIndex < len(cweapi.listNexts) {
		cweapi.listRulesOutput.NextToken = &cweapi.listNexts[*cweapi.listIndex]
		*cweapi.listIndex = (*cweapi.listIndex + 1) % len(cweapi.listNexts)
	} else {
		cweapi.listRulesOutput.NextToken = nil
	}
	return cweapi.listRulesOutput, nil
}

// DescribeRule mocks cloudwatchevents.DescribeRule
func (cweapi *CWEAPI) DescribeRule(
	*cwe.DescribeRuleInput) (*cwe.DescribeRuleOutput, error) {
	return cweapi.describeRuleOutput, cweapi.describeRuleError
}

// PutRule mocks cloudwatchevents.PutRule
func (cweapi *CWEAPI) PutRule(*cwe.PutRuleInput) (*cwe.PutRuleOutput, error) {
	return cweapi.putRuleOutput, cweapi.putRuleError
}

// PutTargets mocks cloudwatchevents.PutRule
func (cweapi *CWEAPI) PutTargets(
	*cwe.PutTargetsInput) (*cwe.PutTargetsOutput, error) {
	return cweapi.putTargetsOutput, cweapi.putTargetsError
}

func mockListRuleOutput(names []string) *cwe.ListRuleNamesByTargetOutput {
	return &cwe.ListRuleNamesByTargetOutput{
		RuleNames: aws.StringSlice(names),
	}
}

func mockDescribeRuleOutput(
	arn, name, schedule, state string) *cwe.DescribeRuleOutput {
	return &cwe.DescribeRuleOutput{
		Arn:                aws.String(arn),
		Name:               aws.String(name),
		ScheduleExpression: aws.String(schedule),
		State:              aws.String(state),
	}
}

func mockPutRuleOutput(arn string) *cwe.PutRuleOutput {
	return &cwe.PutRuleOutput{RuleArn: aws.String(arn)}
}

func mockPutTargetsOutputError(err string) *cwe.PutTargetsOutput {
	return &cwe.PutTargetsOutput{
		FailedEntries: []*cwe.PutTargetsResultEntry{&cwe.PutTargetsResultEntry{
			ErrorCode:    aws.String("fake-code"),
			ErrorMessage: aws.String(err),
			TargetId:     aws.String("fake-id"),
		}},
		FailedEntryCount: aws.Int64(1),
	}
}

// CWEAPIListRules mocks a CloudWatchEventsAPI API w/ ListRules
func CWEAPIListRules(
	names []string, nexts []string) cweiface.CloudWatchEventsAPI {
	return &CWEAPI{
		listRulesOutput: mockListRuleOutput(names),
		listIndex:       aws.Int(0),
		listNexts:       nexts,
	}
}

// CWEAPIListRulesError mocks a CloudWatchEventsAPI API w/ ListRules error
func CWEAPIListRulesError() cweiface.CloudWatchEventsAPI {
	return &CWEAPI{
		listRulesError: fmt.Errorf("FakeListRulesError"),
	}
}

// CWEAPIListRulesError mocks a CloudWatchEventsAPI API w/ ListRules error
func CWEAPIListRulesLoopError() cweiface.CloudWatchEventsAPI {
	return &CWEAPI{
		// listRulesError: fmt.Errorf("FakeListRulesError"),
		listRulesOutput: mockListRuleOutput([]string{"rule01"}),
		listIndex:       aws.Int(0),
		listNexts:       []string{"123", ""},
		listErrors:      []string{"", "FakeListRulesError"},
	}
}

// CWEAPIDescribeRule mocks a CloudWatchEventsAPI API w/ DescribeRule
func CWEAPIDescribeRule(
	arn, name, schedule, state string) cweiface.CloudWatchEventsAPI {
	return &CWEAPI{
		describeRuleOutput: mockDescribeRuleOutput(arn, name, schedule, state),
	}
}

// CWEAPIDescribeRuleError mocks a CloudWatchEventsAPI API w/ DescribeRule Error
func CWEAPIDescribeRuleError() cweiface.CloudWatchEventsAPI {
	return &CWEAPI{
		describeRuleError: fmt.Errorf("FakeDescribeRuleError"),
	}
}

// CWEAPIPutRule mocks a CloudWatchEventsAPI API w/ PutRule
func CWEAPIPutRule(arn string) cweiface.CloudWatchEventsAPI {
	return &CWEAPI{
		putRuleOutput: mockPutRuleOutput(arn),
	}
}

// CWEAPIPutRuleError mocks a CloudWatchEventsAPI API w/ PutRule Error
func CWEAPIPutRuleError() cweiface.CloudWatchEventsAPI {
	return &CWEAPI{
		putRuleError: fmt.Errorf("FakePutRuleError"),
	}
}

// CWEAPIPutTargets mocks a CloudWatchEventsAPI API w/ PutTargets
func CWEAPIPutTargets() cweiface.CloudWatchEventsAPI {
	return &CWEAPI{
		putTargetsOutput: &cwe.PutTargetsOutput{},
	}
}

// CWEAPIPutTargetsError mocks a CloudWatchEventsAPI API w/ PutTargets Error
func CWEAPIPutTargetsError() cweiface.CloudWatchEventsAPI {
	return &CWEAPI{
		putTargetsOutput: mockPutTargetsOutputError("fake-error"),
		putTargetsError:  fmt.Errorf("FakePutTargetsError"),
	}
}

// CWEAPIAll mocks a CloudWatchEventsAPI generally
func CWEAPIAll(names []string,
	name, arn, schedule, state string) cweiface.CloudWatchEventsAPI {
	return &CWEAPI{
		putTargetsOutput:   &cwe.PutTargetsOutput{},
		listRulesOutput:    mockListRuleOutput(names),
		listIndex:          aws.Int(0),
		putRuleOutput:      mockPutRuleOutput(arn),
		describeRuleOutput: mockDescribeRuleOutput(arn, name, schedule, state),
	}
}

// CWEAPIAllListError mocks a CloudWatchEventsAPI generally w/ ListRuleNames error
func CWEAPIAllListError(
	name, arn, schedule, state, err string) cweiface.CloudWatchEventsAPI {
	return &CWEAPI{
		putTargetsOutput:   &cwe.PutTargetsOutput{},
		putRuleOutput:      mockPutRuleOutput(arn),
		describeRuleOutput: mockDescribeRuleOutput(arn, name, schedule, state),
		listRulesError:     fmt.Errorf(err),
		listIndex:          aws.Int(0),
	}
}

// CWEAPIAllDescribeError mocks a CloudWatchEventsAPI w/ DescribeRule error
func CWEAPIAllDescribeError(
	names []string, arn, err string) cweiface.CloudWatchEventsAPI {
	return &CWEAPI{
		putTargetsOutput:  &cwe.PutTargetsOutput{},
		putRuleOutput:     mockPutRuleOutput(arn),
		listRulesOutput:   mockListRuleOutput(names),
		describeRuleError: fmt.Errorf(err),
		listIndex:         aws.Int(0),
	}
}

// CWEAPIAllPutRuleError mocks a CloudWatchEventsAPI  w/ PutRule error
func CWEAPIAllPutRuleError(names []string,
	name, arn, schedule, state, err string) cweiface.CloudWatchEventsAPI {
	return &CWEAPI{
		putTargetsOutput:   &cwe.PutTargetsOutput{},
		listRulesOutput:    mockListRuleOutput(names),
		describeRuleOutput: mockDescribeRuleOutput(arn, name, schedule, state),
		putRuleError:       fmt.Errorf(err),
		listIndex:          aws.Int(0),
	}
}

// CWEAPIAllPutTargetsError mocks a CloudWatchEventsAPI w/ PutTargets error
func CWEAPIAllPutTargetsError(names []string,
	name, arn, schedule, state, err string) cweiface.CloudWatchEventsAPI {
	return &CWEAPI{
		listRulesOutput:    mockListRuleOutput(names),
		putRuleOutput:      mockPutRuleOutput(arn),
		describeRuleOutput: mockDescribeRuleOutput(arn, name, schedule, state),
		putTargetsError:    fmt.Errorf(err),
		listIndex:          aws.Int(0),
	}
}
