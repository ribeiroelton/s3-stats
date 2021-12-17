package params_test

import (
	"testing"

	params "github.com/elribeiro/s3-stats-tool/internal/params"
)

func TestParamsInput(t *testing.T) {
	result := struct {
		wantFilterName      string
		wantNumberOfThreads int
		wantGetReplication  bool
		wantGetLifecycle    bool
		wantWriteToFile     bool
		wantObjectFilter    string
		wantBucketFilter    string
	}{
		wantFilterName:      "",
		wantNumberOfThreads: 2,
		wantGetReplication:  false,
		wantGetLifecycle:    false,
		wantWriteToFile:     false,
		wantObjectFilter:    "",
		wantBucketFilter:    "",
	}

	params := params.ParamsInput()

	if params.FilterBucketName != result.wantFilterName {
		t.Errorf("Expecting %v, got %v", result.wantFilterName, params.FilterBucketName)
	}

	if params.NumberOfThreads != result.wantNumberOfThreads {
		t.Errorf("Expecting %v, got %v", result.wantNumberOfThreads, params.NumberOfThreads)
	}

	if params.GetReplicationRules != result.wantGetReplication {
		t.Errorf("Expecting %v, got %v", result.wantGetReplication, params.GetReplicationRules)
	}

	if params.GetLifecycleRules != result.wantGetLifecycle {
		t.Errorf("Expecting %v, got %v", result.wantGetLifecycle, params.GetLifecycleRules)
	}

	if params.WriteToFile != result.wantWriteToFile {
		t.Errorf("Expecting %v, got %v", result.wantWriteToFile, params.WriteToFile)
	}

	if params.FilterObjectPrefix != result.wantObjectFilter {
		t.Errorf("Expecting %v, got %v", result.wantObjectFilter, params.FilterObjectPrefix)
	}

	if params.FilterBucketName != result.wantBucketFilter {
		t.Errorf("Expecting %v, got %v", result.wantBucketFilter, params.FilterBucketName)
	}

}
