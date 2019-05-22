package awsapi

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// NewSSMAPI returns a s3 client
func NewSSMAPI(sess AWSSession) ssmiface.SSMAPI {
	s3ess, err := sess.NewSession()
	if err != nil {
		log.Printf("AWSError: Cannot create aws session. err %#v", err)
		return nil
	}
	return ssm.New(s3ess)
}

// PutParameter helps to put parameter to ssm as Secure String
func (api *AWSAPI) PutParameter(svc ssmiface.SSMAPI, name, value,
	keyid string, tags map[string]string, overwrite bool, isArray bool) error {
	input := &ssm.PutParameterInput{
		Name:      aws.String(name),
		Overwrite: aws.Bool(overwrite),
		Value:     aws.String(value),
	}
	if keyid != "" {
		input.Type = aws.String(ssm.ParameterTypeSecureString)
		input.KeyId = aws.String(keyid)
	} else if isArray {
		input.Type = aws.String(ssm.ParameterTypeStringList)
	} else {
		input.Type = aws.String(ssm.ParameterTypeString)
	}
	if tags != nil {
		input.Tags = make([]*ssm.Tag, len(tags))
		idx := 0
		for key, val := range tags {
			input.Tags[idx] = &ssm.Tag{
				Key:   aws.String(key),
				Value: aws.String(val),
			}
			idx++
		}
	}
	result, err := svc.PutParameter(input)
	if err != nil {
		log.Printf("SSMError: fail to put parameter to ssm (name=%s, value=%s, "+
			"keyid=%s, tags=%v, overwrite=%t), err: %v\n",
			name, value, keyid, tags, overwrite, err)
		return err
	}
	log.Printf("SSM: succeed to put parameter to ssm (name=%s, value=%s, "+
		"keyid=%s, tags=%v, overwrite=%t), versionID %d\n",
		name, value, keyid, tags, overwrite, *result.Version)
	return nil
}

// AddTagsToResource helps to add tags to ssm resources
func (api *AWSAPI) AddTagsToResource(svc ssmiface.SSMAPI,
	resID, resType string, tags map[string]string) error {
	input := &ssm.AddTagsToResourceInput{
		ResourceId:   aws.String(resID),
		ResourceType: aws.String(resType),
	}
	input.Tags = make([]*ssm.Tag, len(tags))
	idx := 0
	for key, val := range tags {
		input.Tags[idx] = &ssm.Tag{
			Key:   aws.String(key),
			Value: aws.String(val),
		}
		idx++
	}
	_, err := svc.AddTagsToResource(input)
	if err != nil {
		log.Printf("SSMError: fail to add tags %v to %s\n", tags, resID)
		return err
	}
	log.Printf("SSM: Succeed to add tags %v to %s\n", tags, resID)
	return nil
}

// GetParametersByPathIter helps to get all parameters in a certain path
func (api *AWSAPI) GetParametersByPathIter(svc ssmiface.SSMAPI,
	name string, next *string, recursive bool, maxsize int64) ([]*ssm.Parameter, *string, error) {
	input := &ssm.GetParametersByPathInput{
		MaxResults:     aws.Int64(maxsize),
		Path:           aws.String(name),
		Recursive:      aws.Bool(recursive),
		WithDecryption: aws.Bool(true),
	}
	if next != nil {
		input.NextToken = next
	}
	result, err := svc.GetParametersByPath(input)
	if err != nil {
		log.Printf("SSMError: fail to get parameters by path %s, err: %v\n", name, err)
		return nil, nil, err
	}
	return result.Parameters, result.NextToken, nil
}

// GetParametersByPath helps to get all parameters in a certain path
func (api *AWSAPI) GetParametersByPath(svc ssmiface.SSMAPI,
	name string, recursive bool, maxsize int64) (map[string]string, error) {
	results := make([]*ssm.Parameter, 0)
	sparams, next, err := api.GetParametersByPathIter(svc, name, nil, recursive, maxsize)
	if err == nil && len(sparams) == 0 {
		return nil, fmt.Errorf("parameter %s not found", name)
	}
	for next != nil && err == nil {
		results = append(results, sparams...)
		sparams, next, err = api.GetParametersByPathIter(svc, name, next, recursive, maxsize)
	}
	if err != nil {
		return nil, err
	}
	results = append(results, sparams...)

	paramMap := make(map[string]string)
	for _, param := range results {
		paramMap[*param.Name] = *param.Value
	}
	return paramMap, nil
}

// ListTagsForResource helps to list tags of a certain resource
func (api *AWSAPI) ListTagsForResource(svc ssmiface.SSMAPI,
	resID, resType string) (map[string]string, error) {
	input := &ssm.ListTagsForResourceInput{
		ResourceId:   aws.String(resID),
		ResourceType: aws.String(resType),
	}
	result, err := svc.ListTagsForResource(input)
	if err != nil {
		log.Printf("SSMError: fail to list tags of %s: %v\n", resID, err)
		return nil, err
	}
	tags := make(map[string]string)
	for _, ssmTag := range result.TagList {
		tags[*ssmTag.Key] = *ssmTag.Value
	}
	return tags, nil
}

// DeleteParameters helps to delete a list of parameters at maxsize 10
func (api *AWSAPI) DeleteParameters(svc ssmiface.SSMAPI, names []string) error {
	input := &ssm.DeleteParametersInput{Names: aws.StringSlice(names)}
	output, err := svc.DeleteParameters(input)
	if err != nil {
		return err
	}
	log.Printf("DeleteParameters: deleted parameters %v\n", output.DeletedParameters)
	log.Printf("DeleteParameters: invalid parameters %v\n", output.InvalidParameters)
	return nil
}
