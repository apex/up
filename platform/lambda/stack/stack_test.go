package stack

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/tj/assert"
)

func TestResourcesCompleted(t *testing.T) {
	resources := []*cloudformation.StackResource{
		{
			LogicalResourceId:    aws.String("DnsZoneSomethingComRecordApiSomethingCom"),
			PhysicalResourceId:   aws.String("api.something.com"),
			ResourceStatus:       aws.String("CREATE_IN_PROGRESS"),
			ResourceStatusReason: aws.String("Resource creation Initiated"),
			ResourceType:         aws.String("AWS::Route53::RecordSet"),
			StackId:              aws.String("arn:aws:cloudformation:us-west-2:foobarbaz:stack/app/ad3af570-8511-11e7-8832-50d5ca789e4a"),
			StackName:            aws.String("app"),
		},
		{
			LogicalResourceId:  aws.String("ApiProxyMethod"),
			PhysicalResourceId: aws.String("app-ApiProx-33K7PKBL7HNI"),
			ResourceStatus:     aws.String("CREATE_COMPLETE"),
			ResourceType:       aws.String("AWS::ApiGateway::Method"),
			StackId:            aws.String("arn:aws:cloudformation:us-west-2:foobarbaz:stack/app/ad3af570-8511-11e7-8832-50d5ca789e4a"),
			StackName:          aws.String("app"),
		},
		{
			LogicalResourceId: aws.String("Another"),
			ResourceStatus:    aws.String("CREATE_COMPLETE"),
			ResourceType:      aws.String("AWS::ApiGateway::Method"),
			StackId:           aws.String("arn:aws:cloudformation:us-west-2:foobarbaz:stack/app/ad3af570-8511-11e7-8832-50d5ca789e4a"),
			StackName:         aws.String("app"),
		},
	}

	states := map[string]Status{
		"DnsZoneSomethingComRecordApiSomethingCom": CreateComplete,
		"app-ApiProx-33K7PKBL7HNI":                 CreateComplete,
	}

	c := resourcesCompleted(resources, states)
	assert.Len(t, c, 1)
}
