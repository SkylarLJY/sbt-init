package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
)

func getVersion() (string, error) {
	// get Scala version
	var out bytes.Buffer
	cmd := exec.Command("scala", "--version")

	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}
	re := regexp.MustCompile(`\d+\.\d+\.\d+`)
	res := re.FindString(out.String())
	if res == "" {
		return "", errors.New("No Scala version found")
	}
	return res, nil

}

func createDirs() error {
	// mkdir -p src/{main,test}/{java,resources,scala}
	// mkdir lib project target
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	for _, parent := range []string{"main", "test"} {
		for _, child := range []string{"java", "resource", "scala"} {
			err = os.MkdirAll(path.Join(root, "src", parent, child), 0777)
			if err != nil {
				return err
			}
		}
	}

	for _, dir := range []string{"lib", "project", "target"} {
		err := os.Mkdir(path.Join(root, dir), 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

func createBuildFile(name string, version string) error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	f, err := os.Create(path.Join(root, "build.sbt"))
	if err != nil {
		return err
	}
	content := fmt.Sprintf("name := \"%s\"\nversion := \"1.0\"\nscalaVersion := \"%s\"", name, version)
	_, err = f.Write([]byte(content))
	return err
}

func main() {
	projName := flag.String("proj", "", "name of the Scala project")
	flag.Parse()
	if *projName == "" {
		fmt.Println("Please enter a name for the project")
		os.Exit(0)
	}
	version, err := getVersion()
	if err != nil {
		log.Fatal(err)
	}
	// run commands to create dirs
	if err = createDirs(); err != nil {
		log.Fatal(err)
	}

	// create build.sbt
	if err = createBuildFile(*projName, version); err != nil {
		log.Fatal(err)
	}
}
