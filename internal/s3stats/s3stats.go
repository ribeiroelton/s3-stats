package s3stats

import (
	"context"
	"sync"

	"time"

	"com.github.elribeiro.s3-analytics/package/s3client"
	log "github.com/sirupsen/logrus"
)

type S3ClientApi interface {
	GetAllBuckets(c context.Context, params *s3client.AllBucketsInput) (s3client.AllBucketsOutput, error)

	GetObjectStats(c context.Context, params *s3client.ObjectStatsInput) (*s3client.ObjectStatsOutput, error)

	GetBucketReplicationInfo(c context.Context,
		params *s3client.BucketReplicationInfoInput) (s3client.BucketReplicationInfoOutput, error)

	GetBucketLifecycleInfo(c context.Context,
		params *s3client.BucketLifeCycleInfoInput) (s3client.BucketLifeCycleInfoOutput, error)
}

type S3Stats struct {
	Api S3ClientApi
}

type GenerateBucketStatsInput struct {
	GetReplicationRules bool
	GetLifecycleRules   bool
	FilterObjectPrefix  string
	FilterBucketName    string
	NumberOfThreads     int
}

type GetBucketStatsInput struct {
	GetReplicationRules bool
	GetLifecycleRules   bool
	FilterPrefix        string
}
type BucketStats struct {
	Name                       string            `json:"name"`
	CreationDate               time.Time         `json:"creation_date"`
	TotalFiles                 int64             `json:"total_files"`
	SizeInKB                   int64             `json:"size_in_kb"`
	MostRecentFile             string            `json:"most_recent_file"`
	MostRecentFileModifiedDate time.Time         `json:"most_recent_file_modified_date"`
	ReplicationRules           []ReplicationRule `json:"replication_rules"`
	LifecycleRules             []LifecycleRule   `json:"lifecyle_rules"`
}

type GenerateBucketStatsOutput struct {
	BucketsStats []BucketStats
}

type ReplicationRule struct {
	DestinationBucket  string `json:"destination_bucket"`
	DestinationAccount string `json:"destination_account"`
	StorageClass       string `json:"storage_class"`
	ID                 string `json:"id"`
	Priority           int32  `json:"priority"`
	Status             string `json:"status"`
}

type LifecycleRule struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func NewS3Stats() *S3Stats {
	return &S3Stats{
		Api: s3client.NewS3Client(),
	}
}

var (
	lock = sync.Mutex{}
	bsl  = []BucketStats{}
	wg   = sync.WaitGroup{}
)

func (s3s S3Stats) GenerateBucketStats(params *GenerateBucketStatsInput) (GenerateBucketStatsOutput, error) {
	if params.NumberOfThreads == 0 {
		params.NumberOfThreads = 1
	}

	bl, err := s3s.Api.GetAllBuckets(context.TODO(), &s3client.AllBucketsInput{FilterBucketName: params.FilterBucketName})
	if err != nil {
		log.Error("Error while getting bucket list: ", err)
		return GenerateBucketStatsOutput{}, err
	}
	if len(bl.Buckets) == 0 {
		log.Error("No buckets found for the provided account and region")
	} else {
		log.Infof("Got %v buckets", len(bl.Buckets))
	}

	inputChannel := make(chan s3client.Bucket, len(bl.Buckets))

	log.Infof("Creating %v workers for concurrent running", params.NumberOfThreads)
	p := GetBucketStatsInput{
		GetReplicationRules: params.GetReplicationRules,
		GetLifecycleRules:   params.GetLifecycleRules,
		FilterPrefix:        params.FilterObjectPrefix,
	}
	for i := 0; i < params.NumberOfThreads; i++ {
		go s3s.getBucketStats(inputChannel, &p)
	}
	wg.Add(params.NumberOfThreads)

	for _, b := range bl.Buckets {
		log.Infof("Queeing %v for processing", b.Name)
		inputChannel <- b
	}

	close(inputChannel)
	wg.Wait()

	return GenerateBucketStatsOutput{BucketsStats: bsl}, nil
}

func (s3s S3Stats) getBucketStats(inputChannel chan s3client.Bucket, params *GetBucketStatsInput) {
	for b := range inputChannel {

		log.Infof("Getting stats for bucket %v", b.Name)
		bs, err := s3s.Api.GetObjectStats(context.TODO(), &s3client.ObjectStatsInput{BucketName: b.Name, Prefix: params.FilterPrefix})
		if err != nil {
			log.Error("Error while getting object stats: ", err)
			wg.Done()
			return
		}

		var repRules []ReplicationRule

		if params.GetReplicationRules {
			log.Infof("Getting replication info for bucket %v", b.Name)
			bri, err := s3s.Api.GetBucketReplicationInfo(context.TODO(), &s3client.BucketReplicationInfoInput{BucketName: b.Name})
			if err != nil {
				log.Error("Error while getting replication info stats: ", err)
				wg.Done()
				return
			}

			for _, rr := range bri.ReplicationRules {
				repRules = append(repRules, ReplicationRule{
					DestinationAccount: rr.DestinationAccount,
					DestinationBucket:  rr.DestinationBucket,
					StorageClass:       rr.StorageClass,
					ID:                 rr.ID,
					Priority:           rr.Priority,
					Status:             rr.Status,
				})
			}
		}

		var lcRules []LifecycleRule

		if params.GetLifecycleRules {
			log.Infof("Getting lifecycle info for bucket %v", b.Name)
			lcr, err := s3s.Api.GetBucketLifecycleInfo(context.TODO(), &s3client.BucketLifeCycleInfoInput{BucketName: b.Name})
			if err != nil {
				log.Error("Error while getting lifecycle info stats: ", err)
				wg.Done()
				return
			}
			for _, lc := range lcr.LifeCycleRules {
				lcRules = append(lcRules, LifecycleRule{
					ID:     lc.ID,
					Status: lc.Status,
				})
			}

		}

		log.Infof("Generating ouput data for bucket %v", b.Name)
		lock.Lock()
		bsl = append(bsl, BucketStats{
			Name:                       b.Name,
			CreationDate:               b.CreationDate,
			TotalFiles:                 bs.TotalFiles,
			SizeInKB:                   bs.SizeInKB,
			MostRecentFileModifiedDate: bs.MostRecentFileModifiedDate,
			ReplicationRules:           repRules,
			LifecycleRules:             lcRules,
		})
		lock.Unlock()
	}
	wg.Done()
}
