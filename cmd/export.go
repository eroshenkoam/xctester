package cmd

import (
	"github.com/eroshennkoam/xcresults/xcrun"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var output string
var clean bool

var exportCmd = &cobra.Command{
	Use:   "export [path-to-xcresult]",
	Short: "Export Xcode 11+ xcresult to allure format",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := filepath.Abs(args[0])
		log.Println("Input directory:", input)
		if _, err := os.Stat(input); os.IsNotExist(err) {
			log.Fatalln("Directory [", input, "] does not exists")
		}
		log.Println("Output directory:", output)
		if _, err := os.Stat(output); !os.IsNotExist(err) && !clean {
			log.Fatalln("Output directory already exists, use [--clean] option for deleting it. Skipping...")
			return
		} else {
			log.Println("Cleaning output directory...")
			if err := os.RemoveAll(output); err != nil {
				log.Fatalln("Can not delete output directory...")
			}
		}
		if _, err := os.Stat(output); os.IsNotExist(err) {
			log.Println("Creating output directory")
			if err := os.Mkdir(output, 0700); err != nil {
				log.Fatalln("Can not create output directory...")
			}
		}
		xcrun.Export(input, output)
	},
}

func init() {
	exportCmd.Flags().StringVarP(&output, "output", "o", "allure-results", "output directory (default is 'allure-results')")
	exportCmd.Flags().BoolVarP(&clean, "clean", "c", false, "clean output directory if exists")
	RootCmd.AddCommand(exportCmd)
}
