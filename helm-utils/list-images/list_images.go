package list_images

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"sort"
)

func ListChartImages(helmPath string, chartPath string, args []string) (images []DockerImage, err error) {
	helmArgs := append([]string{"template", chartPath}, args...)
	cmd := exec.Command(helmPath, helmArgs...)
	var buf bytes.Buffer

	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		return nil, err
	}
	images, err = ParseImages(&buf)
	return images, err
}

func ParseImages(buf *bytes.Buffer) (images []DockerImage, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// We want to match the repo name, the image name, the tag name: docker.io/redis:latest -> docker.io | redis | latest
	regex := regexp.MustCompile(`(?m)^\s*image:\s*(?:"|')?(?:(?P<registry>.*?\..+?)\/)?(?:(?P<repo>(?:.*\/)?(?P<image>[^\s:]+)))(?:(?:\s*:\s*)(?P<tag>[^"'\s]+)*(?:"|')?\s*)?$`)
	results := regex.FindAllStringSubmatch(buf.String(), -1)
	names := regex.SubexpNames()

	tempMap := map[string]string{}
	dockerImages := map[DockerImage]bool{}

	for _, match := range results {
		for j, n := range match {
			tempMap[names[j]] = n
		}
		img, err := NewDockerImage(tempMap)
		if err != nil {
			return nil, err
		}
		dockerImages[img] = true
	}

	// Unique
	result := []DockerImage{}
	for image := range dockerImages {
		result = append(result, image)
	}

	// Ordered alphabetically
	sort.Slice(result, func(i, j int) bool {
		if result[i].Repo != result[j].Repo {
			return result[i].Repo < result[j].Repo
		}
		return result[i].Tag < result[j].Tag
	})

	return result, nil
}
