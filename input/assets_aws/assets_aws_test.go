package assets_aws

import (
    "github.com/aws/aws-sdk-go-v2/service/ec2/types"
    "github.com/stretchr/testify/assert"
    "testing"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/elastic-agent-libs/mapstr"
	"github.com/elastic/inputrunner/mocks"
	"github.com/golang/mock/gomock"
)

func TestAssetsAWS_publishAWSAsset_IncludesRequiredFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	publisher := mocks.NewMockPublisher(ctrl)
	expectedAsset := mapstr.M{
		"cloud.provider":   "aws",
		"cloud.region":     "eu-west-1",
		"cloud.account.id": "1234",
		"asset.type":       "aws.ec2.instance",
		"asset.id":         "i-1234",
		"asset.ean":        "aws.ec2.instance:i-1234",
	}
	publisher.EXPECT().Publish(beat.Event{Fields: expectedAsset})
	publishAWSAsset(publisher, "eu-west-1", "1234", "aws.ec2.instance", "i-1234", nil, nil, nil, nil)
}

func TestAssetsAWS_publishAWSAsset_IncludesTagsInMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	publisher := mocks.NewMockPublisher(ctrl)
	expectedAsset := mapstr.M{
		"cloud.provider":   "aws",
		"cloud.region":     "eu-west-1",
		"cloud.account.id": "1234",
		"asset.type":       "aws.ec2.instance",
		"asset.id":         "i-1234",
		"asset.ean":        "aws.ec2.instance:i-1234",
		"asset.metadata": mapstr.M{
			"tags": map[string]string{
				"tag1": "a",
				"tag2": "b",
			},
		},
	}
	publisher.EXPECT().Publish(beat.Event{Fields: expectedAsset})

	tags := map[string]string{"tag1": "a", "tag2": "b"}
	publishAWSAsset(publisher, "eu-west-1", "1234", "aws.ec2.instance", "i-1234", nil, nil, tags, nil)
}

func TestAssetsAWS_flattenEC2Tags(t *testing.T) {
    tag1, tag2, a, b := "tag1", "tag2", "a", "b"
    tags := []types.Tag{{Key: &tag1, Value: &a}, {Key: &tag2, Value: &b}}
    flat := flattenEC2Tags(tags)
    expected := map[string]string{"tag1": "a", "tag2": "b"}
    assert.Equal(t, expected, flat)
}