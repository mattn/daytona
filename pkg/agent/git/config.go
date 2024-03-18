package git

import (
	"bytes"
	"os"
	"path"

	"github.com/daytonaio/daytona/pkg/serverapiclient"
	"gopkg.in/ini.v1"
)

func SetGitConfig(userData *serverapiclient.GitUserData) error {
	gitConfigFileName := path.Join(os.Getenv("HOME"), ".gitconfig")

	var gitConfigContent []byte
	gitConfigContent, err := os.ReadFile(gitConfigFileName)
	if err != nil {
		gitConfigContent = []byte{}
	}

	cfg, err := ini.Load(gitConfigContent)
	if err != nil {
		return err
	}

	if !cfg.HasSection("credential") {
		cfg.NewSection("credential")
	}

	cfg.Section("credential").NewKey("helper", "/usr/local/bin/daytona git-cred")

	if userData != nil {
		if !cfg.HasSection("user") {
			cfg.NewSection("user")
		}

		if userData.Name != nil {
			cfg.Section("user").NewKey("name", *userData.Name)
		}
		if userData.Email != nil {
			cfg.Section("user").NewKey("email", *userData.Email)
		}
	}

	var buf bytes.Buffer
	_, err = cfg.WriteTo(&buf)
	if err != nil {
		return err
	}

	err = os.WriteFile(gitConfigFileName, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
