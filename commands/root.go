package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "xctester",
	Short: "Tool for parsing Xcode 11+ xcresult",
}

func Execute() {
	log.SetFlags(0)
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
