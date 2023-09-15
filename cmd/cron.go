/*
Copyright © 2022 sfk <sfk@live.cn>
*/
package cmd

import (
	"strings"

	"github.com/isfk/xiaohui/internal/pkg/global"
	"github.com/isfk/xiaohui/internal/server/cron"
	"github.com/spf13/cobra"
)

// cronCmd represents the cron command
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Start `cron` server.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		global.Author = strings.Trim(Author, "\n")
		global.Branch = strings.Trim(Branch, "\n")
		global.Version = strings.Trim(Version, "\n")
		global.Date = strings.Trim(Date, "\n")
		cron.Start()
	},
}

func init() {
	rootCmd.AddCommand(cronCmd)
}
