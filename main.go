package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var (
	imgs    = flag.String("list", "", "Path to the list of images to tag.")
	pulltag = flag.String("pull-tag", "", "The image tag to pull.")
	tag     = flag.String("tag", "", "The tag to apply to the pulled images.")
	reg     = flag.String("registry", "", "The registry to pull the images from.")
)

func init() {
	flag.Parse()
}

// ImageSpecifier contains the parts of an image name, the registry, the image,
// and the tag.
type ImageSpecifier struct {
	Registry string
	Image    string
	Tag      string
}

// New returns a new ImageSpecifier
func New(reg, img, tag string) *ImageSpecifier {
	i := &ImageSpecifier{
		Registry: reg,
		Image:    img,
		Tag:      tag,
	}
	return i
}

func (i *ImageSpecifier) String() string {
	return fmt.Sprintf("%s/%s:%s", i.Registry, i.Image, i.Tag)
}

// ReadLines slurps in a file and returns the contents divided up into a slice
// of []byte.
func ReadLines(filename string) ([]string, error) {
	var lines []string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, err
}

// PullImage pulls a Docker image.
func PullImage(spec *ImageSpecifier) error {
	cmd := exec.Command(
		"docker",
		"pull",
		spec.String(),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// TagImage tags a docker image.
func TagImage(pullspec *ImageSpecifier, newspec *ImageSpecifier) error {
	cmd := exec.Command(
		"docker",
		"tag",
		"-f",
		pullspec.String(),
		newspec.String(),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// PushImage pushes a docker image.
func PushImage(pushspec *ImageSpecifier) error {
	cmd := exec.Command(
		"docker",
		"push",
		pushspec.String(),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	if *imgs == "" {
		fmt.Println("--list must be set")
		os.Exit(1)
	}
	if *pulltag == "" {
		fmt.Println("--pull-tag must be set")
		os.Exit(1)
	}
	if *tag == "" {
		fmt.Println("--tag must be set")
		os.Exit(1)
	}
	images, err := ReadLines(*imgs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, i := range images {
		var (
			pullimage = New(*reg, i, *pulltag)
			pushimage = New(*reg, i, *tag)
		)
		err = PullImage(pullimage)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = TagImage(pullimage, pushimage)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = PushImage(pushimage)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
