package s3client_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/elribeiro/s3-stats-tool/package/s3client"
)

type S3AwsClientMock struct{}

func (s3c S3AwsClientMock) ListBuckets(ctx context.Context, params *s3.ListBucketsInput,
	optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	timeDate := time.Date(2020, time.April, 10, 22, 40, 20, 11, time.UTC)

	buckets := []types.Bucket{
		{Name: aws.String("bucket1"), CreationDate: &timeDate},
		{Name: aws.String("bucket2"), CreationDate: &timeDate},
	}

	output := &s3.ListBucketsOutput{Buckets: buckets}

	return output, nil

}

func (s3c S3AwsClientMock) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input,
	optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	if *params.Bucket != "bucket1" {
		err := errors.New("Bucket Not Found")
		return nil, err
	}

	d1 := time.Date(2020, time.April, 10, 22, 40, 20, 11, time.UTC)
	d2 := time.Date(2020, time.April, 10, 22, 40, 20, 22, time.UTC)
	objects := []types.Object{
		{Key: aws.String("item1"), Size: 3131, LastModified: &d1},
		{Key: aws.String("item2"), Size: 3232, LastModified: &d2},
	}

	output := &s3.ListObjectsV2Output{Contents: objects}

	return output, nil

}

func (s3c S3AwsClientMock) GetBucketReplication(ctx context.Context, params *s3.GetBucketReplicationInput,
	optFns ...func(*s3.Options)) (*s3.GetBucketReplicationOutput, error) {
	if *params.Bucket != "bucket1" {
		err := errors.New("Bucket Not Found")
		return nil, err
	}

	d := &types.Destination{Bucket: aws.String("bucket1"), Account: aws.String("123456789123"), StorageClass: types.StorageClassDeepArchive}

	rr := []types.ReplicationRule{
		{Destination: d, Status: types.ReplicationRuleStatusEnabled, ID: aws.String("id1"), Priority: 1},
		{Destination: d, Status: types.ReplicationRuleStatusDisabled, ID: aws.String("id2"), Priority: 2},
	}

	rc := types.ReplicationConfiguration{Rules: rr}

	r := s3.GetBucketReplicationOutput{ReplicationConfiguration: &rc}

	return &r, nil

}

func (s3c S3AwsClientMock) GetBucketLifecycleConfiguration(ctx context.Context, params *s3.GetBucketLifecycleConfigurationInput,
	optFns ...func(*s3.Options)) (*s3.GetBucketLifecycleConfigurationOutput, error) {
	if *params.Bucket != "bucket1" {
		err := errors.New("Bucket Not Found")
		return nil, err
	}

	r := []types.LifecycleRule{
		{Status: types.ExpirationStatusEnabled, ID: aws.String("id1")},
		{Status: types.ExpirationStatusDisabled, ID: aws.String("id2")},
	}

	blc := s3.GetBucketLifecycleConfigurationOutput{Rules: r}

	return &blc, nil
}

func (s3c S3AwsClientMock) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput,
	optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return &s3.GetBucketLocationOutput{LocationConstraint: types.BucketLocationConstraintApEast1}, nil
}

func TestGetAllBuckets(t *testing.T) {
	result := &struct {
		wantTotalListSize    int
		wantFilteredListSize int
		wantName             string
		wantCreationDate     time.Time
	}{
		wantTotalListSize:    2,
		wantFilteredListSize: 1,
		wantName:             "bucket1",
		wantCreationDate:     time.Date(2020, time.April, 10, 22, 40, 20, 11, time.UTC),
	}

	thisTime := time.Now()
	t.Log("Starting Test: ", thisTime)
	api := S3AwsClientMock{}
	s3c := s3client.S3Client{Api: api}

	bsl, _ := s3c.GetAllBuckets(context.TODO(), &s3client.AllBucketsInput{})
	if len(bsl.Buckets) != result.wantTotalListSize {
		t.Errorf("Expecting %v, got %v ", result.wantTotalListSize, len(bsl.Buckets))
	}
	for _, b := range bsl.Buckets {
		if b.Name == "" {
			t.Errorf("Expecting bucket name, got empty")
		}
		if b.CreationDate != result.wantCreationDate {
			t.Errorf("Expecting %v, but got: %v", result.wantCreationDate, b.CreationDate)
		}
	}

	bsl, _ = s3c.GetAllBuckets(context.TODO(), &s3client.AllBucketsInput{FilterBucketName: "bucket1"})
	if len(bsl.Buckets) != result.wantFilteredListSize {
		t.Errorf("Expecting %v, got %v ", result.wantFilteredListSize, len(bsl.Buckets))
	}
}

func TestGetObjectStats(t *testing.T) {
	result := &struct {
		wantSize         int64
		wantCreationDate time.Time
		wantTotalFiles   int64
	}{
		wantSize:         6363 / 1024,
		wantCreationDate: time.Date(2020, time.April, 10, 22, 40, 20, 22, time.UTC),
		wantTotalFiles:   2,
	}

	thisTime := time.Now()
	t.Log("Starting Test: ", thisTime)
	api := S3AwsClientMock{}
	s3c := s3client.S3Client{Api: api}

	object, _ := s3c.GetObjectStats(context.TODO(), &s3client.ObjectStatsInput{BucketName: "bucket1"})
	if object.SizeInKB != result.wantSize {
		t.Errorf("Expecting %v , got %v ", result.wantSize, object.SizeInKB)
	}
	if object.MostRecentFileModifiedDate != result.wantCreationDate {
		t.Errorf("Expecting %v, got %v", result.wantCreationDate, object.MostRecentFileModifiedDate)
	}
	if object.TotalFiles != result.wantTotalFiles {
		t.Errorf("Expecting %v, got %v", result.wantTotalFiles, object.TotalFiles)
	}

	_, err := s3c.GetObjectStats(context.TODO(), &s3client.ObjectStatsInput{BucketName: "bucket2"})
	notFoundMsg := "Bucket Not Found"
	if err.Error() != notFoundMsg {
		t.Errorf("Expected %v, got %v", notFoundMsg, err.Error())
	}

	_, err = s3c.GetObjectStats(context.TODO(), &s3client.ObjectStatsInput{})
	notFoundMsg = "Bucket name is required"
	if err.Error() != notFoundMsg {
		t.Errorf("Expected %v, got %v", notFoundMsg, err.Error())
	}
}
func TestGetBucketReplicationInfo(t *testing.T) {
	result := &struct {
		wantListSize int
		wantDAccount string
		wantDBucket  string
		wantSClass   types.StorageClass
		wantStatus   types.ReplicationRuleStatus
		wantID       string
		wantPriority int32
	}{
		wantListSize: 2,
		wantDAccount: "123456789123",
		wantDBucket:  "bucket1",
		wantSClass:   types.StorageClassDeepArchive,
		wantStatus:   types.ReplicationRuleStatusEnabled,
		wantID:       "id1",
		wantPriority: 1,
	}
	thisTime := time.Now()
	t.Log("Starting Test: ", thisTime)
	api := S3AwsClientMock{}
	s3c := s3client.S3Client{Api: api}

	r, _ := s3c.GetBucketReplicationInfo(context.TODO(), &s3client.BucketReplicationInfoInput{BucketName: "bucket1"})

	if len(r.ReplicationRules) != result.wantListSize {
		t.Errorf("Expected %v, got %v", result.wantListSize, len(r.ReplicationRules))
	}
	if r.ReplicationRules[0].DestinationAccount != result.wantDAccount {
		t.Errorf("Expected %v, got %v", result.wantDAccount, r.ReplicationRules[0].DestinationAccount)
	}

	if r.ReplicationRules[0].DestinationBucket != result.wantDBucket {
		t.Errorf("Expected %v, got %v", result.wantDAccount, r.ReplicationRules[0].DestinationBucket)
	}

	if types.StorageClass(r.ReplicationRules[0].StorageClass) != result.wantSClass {
		t.Errorf("Expected %v, got %v", result.wantSClass, r.ReplicationRules[0].StorageClass)
	}

	if r.ReplicationRules[0].Priority != result.wantPriority {
		t.Errorf("Expected %v, got %v", result.wantPriority, r.ReplicationRules[0].Priority)
	}

	if r.ReplicationRules[0].ID != result.wantID {
		t.Errorf("Expected %v, got %v", result.wantID, r.ReplicationRules[0].ID)
	}

	if r.ReplicationRules[0].Status != string(result.wantStatus) {
		t.Errorf("Expected %v, got %v", result.wantStatus, r.ReplicationRules[0].Status)
	}

	_, err := s3c.GetBucketReplicationInfo(context.TODO(), &s3client.BucketReplicationInfoInput{BucketName: "bucket2"})
	notFoundMsg := "Bucket Not Found"
	if err.Error() != notFoundMsg {
		t.Errorf("Expected %v, got %v", notFoundMsg, err.Error())
	}

	_, err = s3c.GetBucketReplicationInfo(context.TODO(), &s3client.BucketReplicationInfoInput{})
	notFoundMsg = "Bucket name is required"
	if err.Error() != notFoundMsg {
		t.Errorf("Expected %v, got %v", notFoundMsg, err.Error())
	}

}

func TestGetBucketLifecycleInfo(t *testing.T) {
	result := &struct {
		wantListSize int
		wantStatus   types.ExpirationStatus
		wantID       string
	}{
		wantListSize: 2,
		wantStatus:   types.ExpirationStatusDisabled,
		wantID:       "id1",
	}
	thisTime := time.Now()
	t.Log("Starting Test: ", thisTime)
	api := S3AwsClientMock{}
	s3c := s3client.S3Client{Api: api}

	r, _ := s3c.GetBucketLifecycleInfo(context.TODO(), &s3client.BucketLifeCycleInfoInput{BucketName: "bucket1"})

	if len(r.LifeCycleRules) != result.wantListSize {
		t.Errorf("Expected %v, got %v", result.wantListSize, len(r.LifeCycleRules))
	}

	if r.LifeCycleRules[0].ID != result.wantID {
		t.Errorf("Expected %v, got %v", result.wantID, r.LifeCycleRules[0].ID)
	}

	if r.LifeCycleRules[1].Status != string(result.wantStatus) {
		t.Errorf("Expected %v, got %v", result.wantStatus, r.LifeCycleRules[1].Status)
	}

	_, err := s3c.GetBucketLifecycleInfo(context.TODO(), &s3client.BucketLifeCycleInfoInput{BucketName: "bucket2"})
	notFoundMsg := "Bucket Not Found"
	if err.Error() != notFoundMsg {
		t.Errorf("Expected %v, got %v", notFoundMsg, err.Error())
	}

	_, err = s3c.GetBucketLifecycleInfo(context.TODO(), &s3client.BucketLifeCycleInfoInput{})
	notFoundMsg = "Bucket name is required"
	if err.Error() != notFoundMsg {
		t.Errorf("Expected %v, got %v", notFoundMsg, err.Error())
	}

}
