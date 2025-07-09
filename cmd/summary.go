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
	"github.com/Friends-Of-Noso/NosoData-Go/utils"
)

var (
	summary legacy.LegacySummary

	// summaryCmd represents the summary command
	summaryCmd = &cobra.Command{
		Use:   "summary",
		Short: "Outputs the summary in text or JSON",
		// 		Long: `A longer description that spans multiple lines and likely contains examples
		// and usage of using your command. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		Example: `  # Display summary in text format
  $ nosodata summary --test-data <path to folder containing "summary.psk">

  # Display summary in JSON format
  $ nosodata summary --json --test-data <path to folder containing "summary.psk">`,
		Run: runSummary,
	}
)

func init() {
	rootCmd.AddCommand(summaryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// summaryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// summaryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runSummary(cmd *cobra.Command, args []string) {
	summaryFile := filepath.Join(
		testdata,
		"sumary.psk",
	)

	if err := summary.ReadFromFile(summaryFile); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	displaySummary(json)
}

func displaySummary(jsonOutput bool) {
	buf := new(bytes.Buffer)

	options := m.ReaderOptions{}

	if jsonOutput {
		options.ShouldFormat = true
		options.Style = styles.Get("native")
		fmt.Fprintln(buf, summary.AsJSON())
	} else {
		options.ShouldFormat = false
		for i, a := range summary.Accounts {
			fmt.Fprintln(buf, "Position:", i)
			fmt.Fprintf(buf, "    Hash:           '%s'\n", a.Hash.GetString())
			fmt.Fprintf(buf, "    Custom:         '%s'\n", a.Custom.GetString())
			fmt.Fprintln(buf, "    Balance:       ", utils.ToNoso(a.Balance))
			fmt.Fprintln(buf, "    Score:         ", utils.ToNoso(a.Score))
			fmt.Fprintln(buf, "    Last Operation:", utils.ToNoso(a.LastOperation))
		}
	}

	reader, err := m.NewReaderFromStream(
		"Summary",
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
