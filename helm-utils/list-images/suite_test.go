package list_images

import (
	"fmt"
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

func TestListImagesWithValues(t *testing.T) {

	images, err := ListChartImages("/usr/local/bin/helm", "./test-0.1.0.tgz", []string{"-f", "../../test/values.yaml"})
	if err != nil {
		fmt.Println(err, "Error while listing chart images")
		t.Fail()
	}
	for _, image := range images {
		if image.Tag == "v1.2.3" {
			t.Fail()
		}
	}
}

func TestListImagesWithoutValues(t *testing.T) {

	images, err := ListChartImages("/usr/local/bin/helm", "./test-0.1.0.tgz", []string{})
	if err != nil {
		fmt.Println(err, "Error while listing chart images")
		t.Fail()
	}
	tagFound := false
	for _, image := range images {
		if image.Tag == "v1.2.3" {
			tagFound = true
		}
	}
	if !tagFound {
		t.Fail()
	}
}
