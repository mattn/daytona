// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/daytonaio/daytona/cmd/daytona/config"
	workspace "github.com/daytonaio/daytona/pkg/cmd/workspace"
	view "github.com/daytonaio/daytona/pkg/views/prune"
	view_util "github.com/daytonaio/daytona/pkg/views/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var forceFlag bool

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prunes all Daytona data from the current device",
	Long:  "Prunes all Daytona data from the current device - including all workspaces, configuration files, and SSH files. This command is irreversible.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var confirmCheck bool
		var serverStoppedCheck bool
		var switchProfileCheck bool

		c, err := config.GetConfig()
		if err != nil {
			log.Fatal(err)
		}

		if c.ActiveProfileId != "default" {
			view.SwitchProfilePrompt(&switchProfileCheck)
			if !switchProfileCheck {
				view_util.RenderInfoMessage("Operation cancelled.")
				return
			}
			c.ActiveProfileId = "default"
			err = c.Save()
			if err != nil {
				log.Fatal(err)
			}
		}

		if !forceFlag {
			view.ConfirmPrompt(&confirmCheck)
			if !confirmCheck {
				view_util.RenderInfoMessage("Operation cancelled.")
				return
			}
		}

		err = workspace.DeleteAllWorkspaces()
		if err != nil {
			log.Fatal(err)
		}

		err = config.UnlinkSshFiles()
		if err != nil {
			log.Fatal(err)
		}

		err = config.DeleteAutocompletionData()
		if err != nil {
			log.Fatal(err)
		}

		view.ServerStoppedPrompt(&serverStoppedCheck)
		if !serverStoppedCheck {
			view_util.RenderInfoMessage("Operation cancelled.")
			return
		}

		err = config.DeleteConfigDir()
		if err != nil {
			log.Fatal(err)
		}

		view_util.RenderInfoMessage("All Daytona data has been successfully cleared from the device.\nYou may now stop the Daytona Server and delete the binary.")
	},
}

func init() {
	pruneCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force prune without prompt")
}
