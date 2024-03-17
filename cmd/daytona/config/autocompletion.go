// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"runtime"
	"strings"
)

const completionScriptNameRoot = "daytona.completion_script."

func DeleteAutocompletionData() error {
	shellName := getCurrentShell()
	if shellName == "" {
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	completionScriptPath := homeDir
	runCommandFilePath := homeDir

	switch shellName {
	case "bash":
		completionScriptPath += "/." + completionScriptNameRoot + "bash"
		runCommandFilePath += "/.bashrc"
	case "zsh":
		completionScriptPath += "/." + completionScriptNameRoot + "zsh"
		runCommandFilePath += "/.zshrc"
	case "fish":
		completionScriptPath += "/." + completionScriptNameRoot + "fish"
		runCommandFilePath += "/.config/fish/config.fish"
	case "powershell":
		completionScriptPath += "/" + completionScriptNameRoot + "ps1"
		runCommandFilePath += "/Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1"
	default:
		return nil
	}

	// Remove the line that sources the completion script
	err = removeLineFromFile(runCommandFilePath, completionScriptNameRoot)
	if err != nil {
		return err
	}

	// Remove the completion script if it exists
	_, err = os.Stat(completionScriptPath)
	if os.IsNotExist(err) {
		return nil
	}

	return os.Remove(completionScriptPath)
}

func getCurrentShell() string {
	var shell string

	if runtime.GOOS == "windows" {
		shell = os.Getenv("PSModulePath")
		if shell == "" {
			return ""
		}

		if strings.Contains(strings.ToLower(shell), "powershell") {
			return "powershell"
		}
	}

	shell = os.Getenv("SHELL")
	if shell == "" {
		return ""
	}

	shellParts := strings.Split(shell, "/")
	shellName := shellParts[len(shellParts)-1]

	// Normalize shell name
	switch shellName {
	case "bash", "zsh", "fish":
		return shellName
	default:
		return ""
	}
}

func removeLineFromFile(filePath string, lineText string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	// Remove the line containing the specified text
	var newLines []string
	for _, line := range lines {
		if !strings.Contains(line, lineText) {
			newLines = append(newLines, line)
		}
	}

	// Join the lines back together
	newContent := strings.Join(newLines, "\n")

	err = os.WriteFile(filePath, []byte(newContent), 0600)
	if err != nil {
		return err
	}

	return nil
}
