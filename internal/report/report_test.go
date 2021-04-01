package report_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"com.github.elribeiro.s3-analytics/internal/report"
	"com.github.elribeiro.s3-analytics/internal/s3stats"
)

func TestOutputData(t *testing.T) {
	lf := s3stats.LifecycleRule{
		ID:     "id1",
		Status: "Enabled",
	}
	rr := s3stats.ReplicationRule{
		DestinationBucket:  "dest1",
		DestinationAccount: "123456789123",
		StorageClass:       "STANDARD",
		ID:                 "id1",
		Priority:           1,
		Status:             "Enabled",
	}
	bs := s3stats.BucketStats{
		Name:                       "bucket1",
		CreationDate:               time.Date(2020, time.April, 10, 22, 40, 20, 11, time.UTC),
		TotalFiles:                 101,
		SizeInKB:                   2048,
		MostRecentFileModifiedDate: time.Date(2020, time.April, 10, 23, 40, 20, 11, time.UTC),
		ReplicationRules:           []s3stats.ReplicationRule{rr},
		LifecycleRules:             []s3stats.LifecycleRule{lf},
	}
	bso := s3stats.GenerateBucketStatsOutput{
		BucketsStats: []s3stats.BucketStats{bs},
	}

	report.OutputData(&report.Report{BucketStats: bso})
	report.OutputData(&report.Report{BucketStats: bso, WriteToFile: true})

	fileName := "s3stats-" + time.Now().Format("2006-01-02") + ".json"

	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Errorf("Got an error while reading output file, details %v", err)
	}
	t.Log(string(f))

	err = os.Remove(fileName)
	if err != nil {
		t.Errorf("Error while file cleanup, details: %v ", err)
	}
}
