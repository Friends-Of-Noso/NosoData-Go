package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	testdata string
	json     bool
	rootCmd  = &cobra.Command{
		Use:     "nosodata",
		Version: "0.0.9",
		Short:   "A tool to inspect Noso's blockchain files",
		// 		Long: `A longer description that spans multiple lines and likely contains
		// examples and usage of using your application. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.NosoData-Go.yaml)")
	rootCmd.PersistentFlags().StringVarP(&testdata, "test-data", "t", "", "data folder containing Noso blockchain files")
	err := rootCmd.MarkPersistentFlagRequired("test-data")
	if err != nil {
		fmt.Printf("could not make 'test-data' required: %v\n", err)
		os.Exit(1)
	}
	err = rootCmd.MarkPersistentFlagDirname("test-data")
	if err != nil {
		fmt.Printf("could not make 'test-data' required: %v\n", err)
		os.Exit(1)
	}

	rootCmd.PersistentFlags().BoolVarP(&json, "json", "j", false, "output JSON")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
