package suite01

import (
	"testing"

	"github.com/goharbor/harbor/src/testing/apitests/api-testing/envs"
)

// TestRun : Start to run the case
func TestRun(t *testing.T) {
	// Initialize env
	if err := envs.ConcourseCIEnv.Load(); err != nil {
		t.Fatal(err.Error())
	}

	suite := ConcourseCiSuite01{}
	report := suite.Run(&envs.ConcourseCIEnv)
	report.Print()
	if report.IsFail() {
		t.Fail()
	}
}
