package comparedate_test

import (
	"testing"
	"time"

	"github.com/elribeiro/s3-stats-tool/internal/comparedate"
)

func TestGetMostRecentDate(t *testing.T) {
	cur := time.Date(2020, time.April, 10, 22, 40, 20, 22, time.UTC)
	new := time.Date(2020, time.April, 10, 22, 40, 20, 23, time.UTC)

	r := comparedate.GetMostRecentDate(&cur, &new)

	if r.Equal(cur) {
		t.Errorf("Expecting %v, got %v", new, r)
	}

	r = comparedate.GetMostRecentDate(&new, &cur)

	if !r.Equal(new) {
		t.Errorf("Expecting %v, got %v", new, r)
	}

}
