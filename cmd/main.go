package main

import (
	params "github.com/elribeiro/s3-stats-tool/internal/params"
	"github.com/elribeiro/s3-stats-tool/internal/report"
	"github.com/elribeiro/s3-stats-tool/internal/s3stats"
	log "github.com/sirupsen/logrus"
)

func main() {
	params := params.ParamsInput()

	s3s := s3stats.NewS3Stats()

	bs, err := s3s.GenerateBucketStats(&s3stats.GenerateBucketStatsInput{
		GetReplicationRules: params.GetReplicationRules,
		GetLifecycleRules:   params.GetLifecycleRules,
		NumberOfThreads:     params.NumberOfThreads,
		FilterObjectPrefix:  params.FilterObjectPrefix,
		FilterBucketName:    params.FilterBucketName})

	if err != nil {
		log.Fatal("Error: ", err)
	}

	report.OutputData(&report.Report{BucketStats: bs, WriteToFile: params.WriteToFile})
}
