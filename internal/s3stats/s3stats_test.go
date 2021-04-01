package s3stats_test

import (
	"context"
	"testing"
	"time"

	"com.github.elribeiro.s3-analytics/internal/s3stats"
	"com.github.elribeiro.s3-analytics/package/s3client"
)

type S3ClientApiMock struct{}

func (s3c S3ClientApiMock) GetAllBuckets(c context.Context, params *s3client.AllBucketsInput) (s3client.AllBucketsOutput, error) {
	bs := []s3client.Bucket{
		{Name: string("bucket1"), CreationDate: time.Date(2020, time.April, 10, 22, 40, 20, 11, time.UTC)},
		{Name: string("bucket2"), CreationDate: time.Date(2020, time.April, 10, 22, 40, 20, 33, time.UTC)},
	}

	abo := s3client.AllBucketsOutput{
		Buckets: bs,
	}

	return abo, nil
}

func (s3c S3ClientApiMock) GetObjectStats(c context.Context, params *s3client.ObjectStatsInput) (*s3client.ObjectStatsOutput, error) {
	bs := s3client.ObjectStatsOutput{
		TotalFiles:                 10,
		SizeInKB:                   2500,
		MostRecentFileModifiedDate: time.Date(2020, time.April, 10, 22, 40, 20, 11, time.UTC),
	}

	return &bs, nil
}

func (s3c S3ClientApiMock) GetBucketReplicationInfo(c context.Context,
	params *s3client.BucketReplicationInfoInput) (s3client.BucketReplicationInfoOutput, error) {
	rr := []s3client.BucketReplicationRule{
		{
			DestinationBucket:  "bucket1",
			DestinationAccount: "123456789012",
			StorageClass:       "DEEP_ARCHIVE",
			Status:             "Enabled",
			Priority:           1,
			ID:                 "id1",
		},
		{
			DestinationBucket:  "bucket2",
			DestinationAccount: "123456789012",
			StorageClass:       "STANDARD_IA",
			Status:             "Enabled",
			Priority:           2,
			ID:                 "id2",
		}}

	r := s3client.BucketReplicationInfoOutput{ReplicationRules: rr}

	return r, nil
}

func (s3s S3ClientApiMock) GetBucketLifecycleInfo(c context.Context,
	params *s3client.BucketLifeCycleInfoInput) (s3client.BucketLifeCycleInfoOutput, error) {
	lcrs := []s3client.BucketLifeCycleRule{
		{
			Status: "Enabled",
			ID:     "id1",
		},
		{
			Status: "Disabled",
			ID:     "id2",
		}}

	lco := s3client.BucketLifeCycleInfoOutput{LifeCycleRules: lcrs}

	return lco, nil
}

func TestGenerateBucketStats(t *testing.T) {
	result := &struct {
		wantListSize         int
		wantFilteredListSize int
		wantCreationDate     time.Time
		wantRegion           string
		wantLifeCycleID      string
		wantLifeCycleSize    int
		wantReplicationID    string
		wantReplicationSize  int
	}{
		wantListSize:         2,
		wantFilteredListSize: 1,
		wantCreationDate:     time.Date(2020, time.April, 10, 22, 40, 20, 11, time.UTC),
		wantRegion:           "us-east-1",
		wantLifeCycleID:      "id1",
		wantLifeCycleSize:    2,
		wantReplicationID:    "id1",
		wantReplicationSize:  2,
	}
	thisTime := time.Now()
	t.Log("Starting Test: ", thisTime)
	api := S3ClientApiMock{}
	s3s := s3stats.S3Stats{Api: api}

	r, _ := s3s.GenerateBucketStats(&s3stats.GenerateBucketStatsInput{GetReplicationRules: true, GetLifecycleRules: true, NumberOfThreads: 10})
	if len(r.BucketsStats) != result.wantListSize {
		t.Errorf("Expecting %v, got %v", result.wantListSize, len(r.BucketsStats))
	}

	if len(r.BucketsStats[0].LifecycleRules) != result.wantLifeCycleSize {
		t.Errorf("Expecting %v, got %v", result.wantLifeCycleSize, len(r.BucketsStats[0].LifecycleRules))
	}

	if len(r.BucketsStats[0].ReplicationRules) != result.wantReplicationSize {
		t.Errorf("Expecting %v, got %v", result.wantReplicationSize, len(r.BucketsStats[0].ReplicationRules))
	}

	if r.BucketsStats[0].LifecycleRules[0].ID != result.wantLifeCycleID {
		t.Errorf("Expecting %v, got %v", result.wantLifeCycleID, r.BucketsStats[0].LifecycleRules[0].ID)
	}

	if r.BucketsStats[0].ReplicationRules[0].ID != result.wantReplicationID {
		t.Errorf("Expecting %v, got %v", result.wantReplicationID, r.BucketsStats[0].ReplicationRules[0].ID)
	}

}
