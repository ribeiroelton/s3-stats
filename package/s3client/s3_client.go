package s3client

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/elribeiro/s3-stats-tool/internal/comparedate"
	log "github.com/sirupsen/logrus"
)

type S3AwsClientApi interface {
	ListBuckets(ctx context.Context,
		params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)

	ListObjectsV2(context.Context,
		*s3.ListObjectsV2Input, ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)

	GetBucketReplication(ctx context.Context,
		params *s3.GetBucketReplicationInput, optFns ...func(*s3.Options)) (*s3.GetBucketReplicationOutput, error)

	GetBucketLifecycleConfiguration(ctx context.Context, params *s3.GetBucketLifecycleConfigurationInput,
		optFns ...func(*s3.Options)) (*s3.GetBucketLifecycleConfigurationOutput, error)

	GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput,
		optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error)
}

type S3Client struct {
	Api S3AwsClientApi
}

type AllBucketsInput struct {
	FilterBucketName string
}
type Bucket struct {
	Name         string
	CreationDate time.Time
}

type AllBucketsOutput struct {
	Buckets []Bucket
}

type ObjectStatsInput struct {
	BucketName string
	Prefix     string
}

type ObjectStatsOutput struct {
	TotalFiles                 int64
	SizeInKB                   int64
	MostRecentFileModifiedDate time.Time
}

type StorageClass struct {
	SizeInKB     int64
	StorageClass string
}

type BucketReplicationInfoInput struct {
	BucketName string
}

type BucketReplicationRule struct {
	DestinationBucket  string
	DestinationAccount string
	StorageClass       string
	ID                 string
	Priority           int32
	Status             string
}

type BucketReplicationInfoOutput struct {
	ReplicationRules []BucketReplicationRule
}

type BucketLifeCycleInfoInput struct {
	BucketName string
}

type BucketLifeCycleRule struct {
	ID     string
	Status string
}

type BucketLifeCycleInfoOutput struct {
	LifeCycleRules []BucketLifeCycleRule
}

var awsConfig aws.Config

func NewS3Client() *S3Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("Error while creating S3 Client: ", err)
	}
	awsConfig = cfg
	return &S3Client{Api: s3.NewFromConfig(cfg)}
}

func (s3c S3Client) GetAllBuckets(c context.Context, params *AllBucketsInput) (AllBucketsOutput, error) {
	lbo, err := s3c.Api.ListBuckets(c, nil)
	if err != nil {
		log.Error("Error while listing buckets: ", err)
		return AllBucketsOutput{}, err
	}

	var bs AllBucketsOutput

	for _, b := range lbo.Buckets {
		if params.FilterBucketName != "" {
			if !strings.Contains(*b.Name, params.FilterBucketName) {
				continue
			}
		}
		nb := Bucket{
			Name:         *b.Name,
			CreationDate: *b.CreationDate,
		}
		bs.Buckets = append(bs.Buckets, nb)
	}

	return bs, nil
}

func (s3c S3Client) GetObjectStats(c context.Context,
	params *ObjectStatsInput) (*ObjectStatsOutput, error) {
	if params.BucketName == "" {
		return nil, errors.New("Bucket name is required")
	}

	bs := ObjectStatsOutput{}

	p := s3.ListObjectsV2Input{
		Bucket: &params.BucketName,
		Prefix: &params.Prefix,
	}

	loc, err := s3c.Api.GetBucketLocation(c, &s3.GetBucketLocationInput{Bucket: &params.BucketName})
	if err != nil {
		log.Error("Bucket location not found ", err)
		return nil, err
	}
	if loc.LocationConstraint == "" {
		loc.LocationConstraint = "us-east-1"
	}

	pg := s3.NewListObjectsV2Paginator(s3c.Api, &p)

	for pg.HasMorePages() {
		loo, err := pg.NextPage(c, func(o *s3.Options) { o.Region = string(loc.LocationConstraint) })
		if err != nil {
			log.Error("Error while listing objects: ", err)
			return nil, err
		}

		for _, o := range loo.Contents {
			bs.TotalFiles += 1
			bs.SizeInKB += (o.Size / 1024)
			bs.MostRecentFileModifiedDate = *comparedate.GetMostRecentDate(&bs.MostRecentFileModifiedDate, o.LastModified)
		}
	}

	return &bs, nil
}

func (s3c S3Client) GetBucketReplicationInfo(c context.Context,
	params *BucketReplicationInfoInput) (BucketReplicationInfoOutput, error) {
	if params.BucketName == "" {
		return BucketReplicationInfoOutput{}, errors.New("Bucket name is required")
	}

	p := &s3.GetBucketReplicationInput{Bucket: &params.BucketName}

	loc, err := s3c.Api.GetBucketLocation(c, &s3.GetBucketLocationInput{Bucket: &params.BucketName})
	if err != nil {
		log.Error("Bucket location not found ", err)
		return BucketReplicationInfoOutput{}, err
	}
	if loc.LocationConstraint == "" {
		loc.LocationConstraint = "us-east-1"
	}

	rep, err := s3c.Api.GetBucketReplication(c, p, func(o *s3.Options) { o.Region = string(loc.LocationConstraint) })
	if err != nil {
		var re *awshttp.ResponseError
		if errors.As(err, &re) {
			if re.Response.StatusCode == 404 {
				log.Debugf("Replication rules not found for bucket %v", *p.Bucket)
				return BucketReplicationInfoOutput{}, nil
			} else if re.Response.StatusCode == 403 {
				log.Debugf("User has no ownership of bucket %v", *p.Bucket)
			}
		}
		return BucketReplicationInfoOutput{}, err
	}

	bri := BucketReplicationInfoOutput{}

	for _, rr := range rep.ReplicationConfiguration.Rules {
		var acc string
		var sc string
		dv := "Same as Source"

		if rr.Destination.Account == nil {
			acc = dv
		} else {
			acc = *rr.Destination.Account
		}

		if rr.Destination.StorageClass == "" {
			sc = dv
		} else {
			sc = string(rr.Destination.StorageClass)
		}

		bri.ReplicationRules = append(bri.ReplicationRules, BucketReplicationRule{
			DestinationAccount: acc,
			DestinationBucket:  *rr.Destination.Bucket,
			ID:                 *rr.ID,
			Priority:           rr.Priority,
			Status:             string(rr.Status),
			StorageClass:       sc,
		})
	}

	return bri, err
}

func (s3c S3Client) GetBucketLifecycleInfo(c context.Context,
	params *BucketLifeCycleInfoInput) (BucketLifeCycleInfoOutput, error) {
	if params.BucketName == "" {
		return BucketLifeCycleInfoOutput{}, errors.New("Bucket name is required")
	}

	p := &s3.GetBucketLifecycleConfigurationInput{Bucket: &params.BucketName}

	loc, err := s3c.Api.GetBucketLocation(c, &s3.GetBucketLocationInput{Bucket: &params.BucketName})
	if err != nil {
		log.Error("Bucket location not found ", err)
		return BucketLifeCycleInfoOutput{}, err
	}
	if loc.LocationConstraint == "" {
		loc.LocationConstraint = "us-east-1"
	}

	lc, err := s3c.Api.GetBucketLifecycleConfiguration(c, p, func(o *s3.Options) { o.Region = string(loc.LocationConstraint) })
	if err != nil {
		var re *awshttp.ResponseError
		if errors.As(err, &re) {
			if re.Response.StatusCode == 404 {
				log.Debugf("Lifecycle rules not found for bucket %v", *p.Bucket)
				return BucketLifeCycleInfoOutput{}, nil
			} else if re.Response.StatusCode == 403 {
				log.Debugf("User has no ownership of bucket %v", *p.Bucket)

			}
		}
		log.Errorf("Error while getting lifecycle info: %v", err)
		return BucketLifeCycleInfoOutput{}, err
	}

	lfr := []BucketLifeCycleRule{}

	for _, r := range lc.Rules {
		lfr = append(lfr, BucketLifeCycleRule{
			ID:     *r.ID,
			Status: string(r.Status),
		})
	}

	return BucketLifeCycleInfoOutput{LifeCycleRules: lfr}, nil
}
