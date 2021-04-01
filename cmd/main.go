package main

import (
	params "com.github.elribeiro.s3-analytics/internal/params"
	"com.github.elribeiro.s3-analytics/internal/report"
	"com.github.elribeiro.s3-analytics/internal/s3stats"
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
