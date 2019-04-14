package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

const (
	Support = "--support_datastore_emulator=true"
)

func main() {

	err := run()
	if err != nil {
		log.Printf("%+v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {

	projectId, err := ready("app.yaml")
	if err != nil {
		return xerrors.Errorf("error ready(): %w", err)
	}
	fmt.Printf("ProjectID:%s\n", projectId)

	err = runDevApp(projectId)
	if err != nil {
		return xerrors.Errorf("error runDevApp(): %w", err)
	}
	return nil
}

func ready(f string) (string, error) {

	_, err := os.Stat(f)
	if err != nil {
		return "", xerrors.Errorf("error devapp is %s path: %w", f, err)
	}

	projectId, err := getProjectID(f)
	if err != nil {
		return "", xerrors.Errorf("error getProjectID(): %w", err)
	}
	return projectId, nil
}

func runDevApp(p string) error {

	emu := []string{"beta", "emulators", "datastore", "env-init"}
	envCmd := exec.Command("gcloud", emu...)

	out, err := envCmd.Output()
	if err != nil {
		return xerrors.Errorf("error emulator env-init command: %w", err)
	}

	setLine := strings.Split(string(out), "\n")
	envLine := make([]string, 0)

	for _, elm := range setLine {
		if elm == "" {
			continue
		}
		s := strings.Split(elm, " ")
		eq := s[1]
		envLine = append(envLine, eq)
	}

	devPath, err := exec.LookPath("dev_appserver.py")
	if err != nil {
		return xerrors.Errorf("error dev_appserver.py LookPath: %w", err)
	}

	devCmd := exec.Command("python", devPath, Support, "-A="+p, ".")
	devCmd.Env = append(os.Environ(), envLine...)
	devCmd.Stdout = os.Stdout
	devCmd.Stderr = os.Stderr
	err = devCmd.Run()

	if err != nil {
		return xerrors.Errorf("error dev_appserver.py: %w", err)
	}

	return nil
}

type AppYaml struct {
	Runtime      string `yaml:"runtime"`
	Main         string `yaml:"main"`
	EnvVariables Env    `yaml:"env_variables"`
}

type Env struct {
	ProjectID string `yaml:"DATASTORE_PROJECT_ID"`
}

func getProjectID(f string) (string, error) {

	b, err := ioutil.ReadFile(f)
	if err != nil {
		return "", xerrors.Errorf("error ioutil.ReadFile(): %w", err)
	}
	var app AppYaml
	err = yaml.Unmarshal(b, &app)
	if err != nil {
		return "", xerrors.Errorf("error yaml.Unmarshal(): %w", err)
	}

	id := app.EnvVariables.ProjectID
	if id == "" {
		return "", xerrors.Errorf("error Project id is empty")
	}

	return id, nil
}
