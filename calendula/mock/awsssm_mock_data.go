package mock

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// SSMAPI mocks s3iface.S3API
type SSMAPI struct {
	ssmiface.SSMAPI
	putparamOutput   *ssm.PutParameterOutput
	addtagsOutput    *ssm.AddTagsToResourceOutput
	getparamsOutput  *ssm.GetParametersByPathOutput
	listtagsOutput   *ssm.ListTagsForResourceOutput
	delparamsOutput  *ssm.DeleteParametersOutput
	getContents      []*ssm.GetParametersByPathOutput
	listtagsContents []*ssm.ListTagsForResourceOutput
	putErr           error
	addtagsErr       error
	getErr           error
	listtagsErr      error
	delparamsErr     error
	getIndex         *int
	listtagsIndex    *int
}

// PutParameter mocks SSMAPI.PutParameter
func (m SSMAPI) PutParameter(
	*ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
	return m.putparamOutput, m.putErr
}

// AddTagsToResource mocks SSMAPI.PutParameter
func (m SSMAPI) AddTagsToResource(
	*ssm.AddTagsToResourceInput) (*ssm.AddTagsToResourceOutput, error) {
	return m.addtagsOutput, m.addtagsErr
}

// GetParametersByPath mocks SSMAPI.PutParameter
func (m SSMAPI) GetParametersByPath(*ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error) {
	output := m.getparamsOutput
	if *m.getIndex < len(m.getContents) {
		output = m.getContents[*m.getIndex]
		*m.getIndex = (*m.getIndex + 1) % len(m.getContents)
	}
	return output, m.getErr
}

// DeleteParameters mocks SSMAPI.DeleteParameters
func (m SSMAPI) DeleteParameters(
	*ssm.DeleteParametersInput) (*ssm.DeleteParametersOutput, error) {
	return m.delparamsOutput, m.delparamsErr
}

// ListTagsForResource mocks SSMAPI.ListTagsForResource
func (m SSMAPI) ListTagsForResource(
	*ssm.ListTagsForResourceInput) (*ssm.ListTagsForResourceOutput, error) {
	output := m.listtagsOutput
	if *m.listtagsIndex < len(m.listtagsContents) {
		output = m.listtagsContents[*m.listtagsIndex]
		*m.listtagsIndex = (*m.listtagsIndex + 1) % len(m.listtagsContents)
	}
	return output, m.listtagsErr
}

// SSMPutParameter returns a pointer to SSMAPI w/ ssm.PutParameter
func SSMPutParameter(version int64) *SSMAPI {
	return &SSMAPI{
		addtagsOutput:  &ssm.AddTagsToResourceOutput{},
		putparamOutput: &ssm.PutParameterOutput{Version: aws.Int64(version)},
	}
}

// SSMPutParameterError returns a *SSMAPI w/ ssm.AddTagsToResource error
func SSMPutParameterError() *SSMAPI {
	return &SSMAPI{
		putErr: fmt.Errorf("FakeSSMPutError"),
	}
}

// SSMAddTags returns a pointer to SSMAPI w/ ssm.PutParameter
func SSMAddTags() *SSMAPI {
	return &SSMAPI{
		putparamOutput: &ssm.PutParameterOutput{Version: aws.Int64(32)},
		addtagsOutput:  &ssm.AddTagsToResourceOutput{},
	}
}

// SSMAddTagsError returns a pointer to SSMAPI w/ ssm.AddTagsToResource error
func SSMAddTagsError() *SSMAPI {
	return &SSMAPI{
		addtagsErr: fmt.Errorf("FakeSSMAddTagsError"),
	}
}

// SSMGetParams returns a pointer to SSMAPI w/ ssm.GetParametersByPath
func SSMGetParams(names, values [][]string, nexts []*string) *SSMAPI {
	paramList := make([]*ssm.GetParametersByPathOutput, len(names))
	for idx := range names {
		p := make([]*ssm.Parameter, len(names[idx]))
		for i := range names[idx] {
			p[i] = &ssm.Parameter{
				Name:  &names[idx][i],
				Value: &values[idx][i],
			}
		}
		paramList[idx] = &ssm.GetParametersByPathOutput{
			NextToken:  nexts[idx],
			Parameters: p,
		}
	}
	return &SSMAPI{
		listtagsOutput:  &ssm.ListTagsForResourceOutput{},
		getparamsOutput: &ssm.GetParametersByPathOutput{},
		getIndex:        aws.Int(0),
		getContents:     paramList,
	}
}

// SSMGetParamsError returns *SSMAPI w/ ssm.GetParametersByPath error
func SSMGetParamsError() *SSMAPI {
	return &SSMAPI{
		getErr:   fmt.Errorf("FakeSSMAGetParamsError"),
		getIndex: aws.Int(0),
	}
}

// SSMListTags returns a pointer to SSMAPI w/ ssm.ListTagsForResource
func SSMListTags(names, values [][]string) *SSMAPI {
	tagList := make([]*ssm.ListTagsForResourceOutput, len(names))
	for idx := range names {
		tlist := make([]*ssm.Tag, len(names[idx]))
		for i := range names[idx] {
			tlist[i] = &ssm.Tag{
				Key:   &names[idx][i],
				Value: &values[idx][i],
			}
		}
		tagList[idx] = &ssm.ListTagsForResourceOutput{
			TagList: tlist,
		}
	}
	return &SSMAPI{
		listtagsOutput:   &ssm.ListTagsForResourceOutput{},
		getparamsOutput:  &ssm.GetParametersByPathOutput{},
		getIndex:         aws.Int(0),
		listtagsIndex:    aws.Int(0),
		listtagsContents: tagList,
	}
}

// SSMListTagsError returns *SSMAPI w/ ssm.ListTagsForResource error
func SSMListTagsError() *SSMAPI {
	return &SSMAPI{
		listtagsErr:   fmt.Errorf("FakeSSMListTagsError"),
		listtagsIndex: aws.Int(0),
	}
}

// SSMDeleteParams returns *SSMAPI w/ ssm.DeleteParameters
func SSMDeleteParams(deleted, invalids []string) *SSMAPI {
	return &SSMAPI{
		delparamsOutput: &ssm.DeleteParametersOutput{
			DeletedParameters: aws.StringSlice(deleted),
			InvalidParameters: aws.StringSlice(invalids),
		},
	}
}

// SSMDeleteParamsError returns *SSMAPI w/ ssm.DeleteParameters error
func SSMDeleteParamsError() *SSMAPI {
	return &SSMAPI{
		delparamsErr: fmt.Errorf("FakeSSMDeleteParamsError"),
	}
}
