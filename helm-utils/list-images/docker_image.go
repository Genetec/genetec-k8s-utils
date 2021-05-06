package list_images

import (
	"fmt"
	"strings"

	ms "github.com/mitchellh/mapstructure"
)

type DockerImage struct {
	Registry string `mapstructure:"registry"`
	Repo     string `mapstructure:"repo"`
	Image    string `mapstructure:"image"`
	Tag      string `mapstructure:"tag"`
	ShaRef   string `mapstructure:"sharef"`
}

func NewDockerImage(imageMap map[string]string) (DockerImage, error) {
	var img DockerImage
	err := ms.Decode(imageMap, &img)
	if err != nil {
		return DockerImage{}, err
	}
	if img.Registry == "" {
		img.Registry = "docker.io"
	}
	if img.Tag == "" {
		img.Tag = "latest"
	}

	needle := "@sha256"
	shaIndex := strings.Index(img.Repo, needle)
	img.Repo = strings.ReplaceAll(img.Repo, `"`, "")
	img.Repo = strings.ReplaceAll(img.Repo, `'`, "")
	if shaIndex != -1 {
		// Skipping the "@"
		img.ShaRef = img.Repo[shaIndex+1:] + ":" + img.Tag
		img.Repo = img.Repo[:shaIndex]
	}

	return img, nil
}

func (image DockerImage) PullReference() string {
	if image.ShaRef != "" {
		return fmt.Sprintf("%s/%s@%s", image.Registry, image.Repo, image.ShaRef)
	}
	return fmt.Sprintf("%s/%s:%s", image.Registry, image.Repo, image.Tag)
}

func (image DockerImage) PushReference() string {
	return fmt.Sprintf("%s/%s:%s", image.Registry, image.Repo, image.Tag)
}

func (image *DockerImage) RepoAddress() string {
	return fmt.Sprintf("%s/%s", image.Registry, image.Repo)
}
