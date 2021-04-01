package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"com.github.elribeiro.s3-analytics/internal/s3stats"
	log "github.com/sirupsen/logrus"
)

type Report struct {
	BucketStats s3stats.GenerateBucketStatsOutput
	WriteToFile bool
}

func OutputData(params *Report) {
	jsonString, err := json.MarshalIndent(params.BucketStats, "", " ")
	if err != nil {
		log.Fatal("Error while marshalling input data, details: ", err)
	}

	if params.WriteToFile {
		fileName := "s3stats-" + time.Now().Format("2006-01-02") + ".json"
		ioutil.WriteFile(fileName, jsonString, os.ModePerm)
	} else {
		fmt.Println(string(jsonString))
	}
}
