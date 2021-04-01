package params

import (
	"flag"
)

type Params struct {
	GetReplicationRules bool
	GetLifecycleRules   bool
	FilterObjectPrefix  string
	FilterBucketName    string
	NumberOfThreads     int
	WriteToFile         bool
}

const (
	numberOfThreadsMsg = `
		Integer to define the number of threads to run concurrently.
		Each thread will process one bucket at a time 
		WATCH OUT: setting a high number may impact in high cost and computing resources usage
		`

	getReplicationRulesMsg = `
		Boolean to define if this job will collect replication rules as well
		 (default false)
	`

	getLifecycleRulesMsg = `
		Boolean to define if this job will collect lifecycle rules as well
		 (default false)
	`

	filterPrefixMsg = `
		String to filter only objects that has a specific prefix
		 (default no filter) 
	`

	filterBucketNameMsg = `
		String to filter only buckets that contains the specified value
		 (default no filter)
	`

	writeToFileMsg = `
		Bool to indicate if output will be to a file named st3stats-date.json,
		where date is the current date. If not set, will output to console in json format
		 (default false)
		`
)

func ParamsInput() *Params {

	numberOfThreads := flag.Int("t", 2, numberOfThreadsMsg)
	getReplicationRules := flag.Bool("r", false, getReplicationRulesMsg)
	getLifecycleRules := flag.Bool("l", false, getLifecycleRulesMsg)
	filterObjectPrefix := flag.String("fo", "", filterPrefixMsg)
	filterBucketName := flag.String("fb", "", filterBucketNameMsg)
	writeToFile := flag.Bool("o", false, writeToFileMsg)

	flag.Parse()

	return &Params{
		GetReplicationRules: *getReplicationRules,
		GetLifecycleRules:   *getLifecycleRules,
		FilterObjectPrefix:  *filterObjectPrefix,
		FilterBucketName:    *filterBucketName,
		NumberOfThreads:     *numberOfThreads,
		WriteToFile:         *writeToFile,
	}
}
