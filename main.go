package main

import (
	_ "embed"
	"fmt"
	"freakbot/app/cmd"
	"freakbot/app/util"
	"freakbot/app/util/mylog"
	"os"

	"github.com/spf13/cobra"
	"go.szostok.io/version/extension"
)

func main() {
	mylog.Preinit()

	fmt.Fprintln(os.Stderr, util.Banner)

	rootCmd := &cobra.Command{Use: "freakbot"}
	rootCmd.AddCommand(cmd.Run)
	rootCmd.AddCommand(extension.NewVersionCobraCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
		return
	}
}
