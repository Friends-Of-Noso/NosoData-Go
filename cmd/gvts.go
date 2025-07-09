package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/spf13/cobra"
	"github.com/walles/moar/m"

	"github.com/Friends-Of-Noso/NosoData-Go/legacy"
)

var (
	gvts legacy.LegacyGVT

	// gvtsCmd represents the gvts command
	gvtsCmd = &cobra.Command{
		Use:   "gvts [flags]",
		Short: "Outputs the gvts in text or JSON",
		// 	Long: `A longer description that spans multiple lines and likely contains examples
		// and usage of using your command. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		Example: `  # Display gvts in text format
  $ nosodata gvts --test-data <path to folder containing "gvts.psk">

  # Display gvts in JSON format
  $ nosodata gvts --json --test-data <path to folder containing "100000.blk"> 100000`,
		Run: runGVTS,
	}
)

func init() {
	rootCmd.AddCommand(gvtsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gvtsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gvtsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runGVTS(cmd *cobra.Command, args []string) {
	gvtsFile := filepath.Join(
		testdata,
		"gvts.psk",
	)

	if err := gvts.ReadFromFile(gvtsFile); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	displayGVTS(json)
}

func displayGVTS(jsonOutput bool) {
	buf := new(bytes.Buffer)

	options := m.ReaderOptions{}

	if jsonOutput {
		options.ShouldFormat = true
		options.Style = styles.Get("native")
		fmt.Fprintln(buf, gvts.AsJSON())
	} else {
		options.ShouldFormat = false
		for i, e := range gvts.Entries {
			fmt.Fprintln(buf, "Position:", i)
			fmt.Fprintf(buf, "    Number:  '%s'\n", e.Number.GetString())
			fmt.Fprintf(buf, "    Owner:   '%s'\n", e.Owner.GetString())
			fmt.Fprintf(buf, "    Hash:    '%s'\n", e.Hash.GetString())
			fmt.Fprintln(buf, "    Control:", e.Control)
		}
	}

	reader, err := m.NewReaderFromStream(
		"GVTS",
		buf,
		formatters.TTY,
		options,
	)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	pager := m.NewPager(reader)
	pager.WrapLongLines = true

	err = pager.Page()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
