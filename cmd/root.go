/*
Copyright © 2023 sfk@live.cn
*/
package cmd

import (
	_ "embed"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ares",
	Short: "A Go Project.",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
