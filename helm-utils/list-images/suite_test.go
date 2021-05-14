package list_images

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestListImages(t *testing.T) {
	RegisterFailHandler(Fail)

	toFile := os.Getenv("TEST_REPORT_PATH")
	if toFile == "" {
		toFile = "go-report.xml"
	}
	junitReporter := reporters.NewJUnitReporter(toFile)

	RunSpecsWithDefaultAndCustomReporters(t, "List Images Suite", []Reporter{junitReporter})
}
