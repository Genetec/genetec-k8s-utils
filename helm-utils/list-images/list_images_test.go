package list_images

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("List images", func() {
	str := ""
	dockerImages := []DockerImage{}
	JustBeforeEach(func() {
		buf := bytes.NewBufferString(str)
		dockerImages, _ = ParseImages(buf)
	})

	Context("Parsing images from a list", func() {
		BeforeEach(func() {
			str = `containers:
			- name: capstan
				image: "127.0.0.1:30000/capstan:0.1.0.0"
				image: "docker.io/capstan:0.1.0.0"
				image: docker.io/capstan:0.1.0.0
				image: docker.io/auth/capstan"
					 `
		})
		It("Has only three images as two are duplicate", func() {
			Expect(dockerImages).To(HaveLen(3))
		})
		It("Has no empty fields", func() {
			for _, di := range dockerImages {
				Expect(di.Repo).ToNot(BeEmpty())
				Expect(di.Registry).ToNot(BeEmpty())
				Expect(di.Image).ToNot(BeEmpty())
				Expect(di.Tag).ToNot(BeEmpty())
				Expect(di.ShaRef).To(BeEmpty())
			}
		})
		It("has good values for its fields", func() {
			img := dockerImages[0]
			Expect(img.Registry).To(Equal("docker.io"))
			Expect(img.Repo).To(Equal("auth/capstan"))
			Expect(img.Tag).To(Equal("latest"))
			Expect(img.PullReference(true)).To(Equal("docker.io/auth/capstan:latest"))
			Expect(img.PushReference(true)).To(Equal("docker.io/auth/capstan:latest"))

			img = dockerImages[1]
			Expect(img.Registry).To(Equal("127.0.0.1:30000"))
			Expect(img.Repo).To(Equal("capstan"))
			Expect(img.Tag).To(Equal("0.1.0.0"))
			Expect(img.PullReference(true)).To(Equal("127.0.0.1:30000/capstan:0.1.0.0"))
			Expect(img.PushReference(true)).To(Equal("127.0.0.1:30000/capstan:0.1.0.0"))
		})
	})
	Context("An image with a sha reference", func() {
		BeforeEach(func() {
			str = `containers:
    - name: capstan
      image: "docker.io/toto/busybox@sha256:c5439d7db88ab5423999530349d327b04279ad3161d7596d2126dfb5b02bfd1f"
         `
		})
		It("has a ShaRef and an appropriate Pull and Push Reference", func() {
			img := dockerImages[0]
			Expect(img.Repo).To(Equal("toto/busybox"))
			Expect(img.Image).To(Equal("busybox"))
			Expect(img.ShaRef).To(Equal("sha256:c5439d7db88ab5423999530349d327b04279ad3161d7596d2126dfb5b02bfd1f"))
			Expect(img.PullReference(true)).To(Equal("docker.io/toto/busybox@sha256:c5439d7db88ab5423999530349d327b04279ad3161d7596d2126dfb5b02bfd1f"))
			Expect(img.PushReference(true)).To(Equal("docker.io/toto/busybox:c5439d7db88ab5423999530349d327b04279ad3161d7596d2126dfb5b02bfd1f"))
		})
	})
	Context("Registry is excluded", func() {
		BeforeEach(func() {
			str = `containers:
    - name: capstan
      image: "docker.io/atoto/busybox@sha256:c5439d7db88ab5423999530349d327b04279ad3161d7596d2126dfb5b02bfd1f"
			image: "docker.io/btoto/busybox:v1.2.3
         `
		})
		It("has no registry in for its pull or push reference", func() {
			img := dockerImages[0]
			Expect(img.Repo).To(Equal("atoto/busybox"))
			Expect(img.ShaRef).To(Equal("sha256:c5439d7db88ab5423999530349d327b04279ad3161d7596d2126dfb5b02bfd1f"))
			Expect(img.PullReference(false)).To(Equal("atoto/busybox@sha256:c5439d7db88ab5423999530349d327b04279ad3161d7596d2126dfb5b02bfd1f"))
			Expect(img.PushReference(false)).To(Equal("atoto/busybox:c5439d7db88ab5423999530349d327b04279ad3161d7596d2126dfb5b02bfd1f"))

			img = dockerImages[1]
			Expect(img.Repo).To(Equal("btoto/busybox"))
			Expect(img.ShaRef).To(Equal(""))
			Expect(img.PullReference(false)).To(Equal("btoto/busybox:v1.2.3"))
			Expect(img.PushReference(false)).To(Equal("btoto/busybox:v1.2.3"))

		})
	})
	Context("Make DockerImage from string", func() {
		var img DockerImage
		BeforeEach(func() {
			str = `127.0.0.1:30000/capstan:0.1.0.0`
			img, _ = NewDockerImageFromString(str)
		})
		It("has no empty fields", func() {
			Expect(img.Repo).ToNot(BeEmpty())
			Expect(img.Registry).ToNot(BeEmpty())
			Expect(img.Image).ToNot(BeEmpty())
			Expect(img.Tag).ToNot(BeEmpty())
			Expect(img.ShaRef).To(BeEmpty())
		})
		It("has good values for its fields", func() {
			Expect(img.Registry).To(Equal("127.0.0.1:30000"))
			Expect(img.Repo).To(Equal("capstan"))
			Expect(img.Tag).To(Equal("0.1.0.0"))
			Expect(img.PullReference(true)).To(Equal("127.0.0.1:30000/capstan:0.1.0.0"))
			Expect(img.PushReference(true)).To(Equal("127.0.0.1:30000/capstan:0.1.0.0"))
		})
	})
	Context("Make DockerImage from sha256 string", func() {
		var img DockerImage
		BeforeEach(func() {
			str = `docker.io/toto/busybox@sha256:c5439d7db88ab5423999530349d327b04279ad3161d7596d2126dfb5b02bfd1f`
			img, _ = NewDockerImageFromString(str)
		})
		It("has a ShaRef and an appropriate Pull and Push Reference", func() {
			Expect(img.Repo).To(Equal("toto/busybox"))
			Expect(img.Image).To(Equal("busybox"))
			Expect(img.ShaRef).To(Equal("sha256:c5439d7db88ab5423999530349d327b04279ad3161d7596d2126dfb5b02bfd1f"))
			Expect(img.PullReference(true)).To(Equal("docker.io/toto/busybox@sha256:c5439d7db88ab5423999530349d327b04279ad3161d7596d2126dfb5b02bfd1f"))
			Expect(img.PushReference(true)).To(Equal("docker.io/toto/busybox:c5439d7db88ab5423999530349d327b04279ad3161d7596d2126dfb5b02bfd1f"))
		})
	})
})
