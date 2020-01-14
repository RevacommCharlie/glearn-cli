package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

const withConfigFixture = "../../fixtures/test-block-with-config"

func Test_PreviewDetectsConfig(t *testing.T) {
	createdConfig, _ := doesConfigExistOrCreate(withConfigFixture, "", false)
	if createdConfig {
		t.Errorf("Created a config when one existed")
	}
}

const withNoConfigFixture = "../../fixtures/test-block-no-config"

func Test_PreviewBuildsAutoConfig(t *testing.T) {
	createdConfig, _ := doesConfigExistOrCreate(withNoConfigFixture, "", false)
	if createdConfig == false {
		t.Errorf("Should of created a config file")
	}

	b, err := ioutil.ReadFile(withNoConfigFixture + "/autoconfig.yaml")
	if err != nil {
		fmt.Print(err)
	}

	config := string(b)

	if !strings.Contains(config, "Title: Unit 1") {
		t.Errorf("Autoconfig should have a unit title of Unit 1")
	}

	if !strings.Contains(config, "Path: /units/test.md") {
		t.Errorf("Autoconfig should have a lesson with a path of /units/test.md")
	}
}

const withNoUnitsDirFixture = "../../fixtures/test-block-no-units-dir"

func Test_PreviewBuildsAutoConfigDeclaredUnitsDir(t *testing.T) {
	createdConfig, _ := doesConfigExistOrCreate(withNoUnitsDirFixture, "foo", false)
	if createdConfig == false {
		t.Errorf("Should of created a config file")
	}

	b, err := ioutil.ReadFile(withNoUnitsDirFixture + "/autoconfig.yaml")
	if err != nil {
		fmt.Print(err)
	}

	config := string(b)

	if !strings.Contains(config, "Title: Foo") {
		t.Errorf("Autoconfig should have a unit title of Foo")
	}

	if !strings.Contains(config, "Path: /foo/test.md") {
		t.Errorf("Autoconfig should have a lesson with a path of /foo/test.md")
	}
}

func Test_AutoConfigAddsInFileTypesOrVisibility(t *testing.T) {
	createdConfig, _ := doesConfigExistOrCreate(withNoConfigFixture, "", false)
	if createdConfig == false {
		t.Errorf("Should of created a config file")
	}

	b, err := ioutil.ReadFile(withNoConfigFixture + "/autoconfig.yaml")
	if err != nil {
		fmt.Print(err)
	}

	config := string(b)

	if !strings.Contains(config, "Type: Checkpoint") {
		t.Errorf("Autoconfig should have a content file of type checkpoint")
	}

	if !strings.Contains(config, "Path: /units/checkpoint.md") {
		t.Errorf("Autoconfig should have a checkpoint with a path of /units/checkpoint.md")
	}

	if !strings.Contains(config, "Type: Instructor") {
		t.Errorf("Autoconfig should have a content file of type Instructor")
	}

	if !strings.Contains(config, "Path: /units/instructor.md") {
		t.Errorf("Autoconfig should have a content file with a path of /units/instructor.md")
	}

	if !strings.Contains(config, "Type: Resource") {
		t.Errorf("Autoconfig should have a content file of type resource")
	}

	if !strings.Contains(config, "Path: /units/resource.md") {
		t.Errorf("Autoconfig should have a content file with a path of /units/resource.md")
	}

	if !strings.Contains(config, "DefaultVisibility: hidden") {
		t.Errorf("Autoconfig should have a content file of with a DefaultVisibility of hidden")
	}

	if !strings.Contains(config, "Path: /units/hidden.file.md") {
		t.Errorf("Autoconfig should have a content file with a path of /units/hidden.file.md")
	}
}
