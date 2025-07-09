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
	pso legacy.LegacyPSO

	// psosCmd represents the psos command
	psosCmd = &cobra.Command{
		Use:   "psos",
		Short: "Outputs the psos in text or JSON",
		// 	Long: `A longer description that spans multiple lines and likely contains examples
		// and usage of using your command. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		Example: `  # Display psos in text format
  $ nosodata psos --test-data <path to folder containing "psos.dat">

  # Display psos in JSON format
  $ nosodata psos --json --test-data <path to folder containing "psos.dat">`,
		Run: runPSOS,
	}
)

func init() {
	rootCmd.AddCommand(psosCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// psoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// psoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runPSOS(cmd *cobra.Command, args []string) {
	psoFile := filepath.Join(
		testdata,
		"psos.dat",
	)

	if err := pso.ReadFromFile(psoFile); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	displayPSO(json)
}

func displayPSO(jsonOutput bool) {
	buf := new(bytes.Buffer)

	options := m.ReaderOptions{}

	if jsonOutput {
		options.ShouldFormat = true
		options.Style = styles.Get("native")
		fmt.Fprintln(buf, pso.AsJSON())
	} else {
		options.ShouldFormat = false
		fmt.Fprintln(buf, "Block:", pso.Block)
		fmt.Fprintf(buf, "  MN Locks(%d):\n", pso.MNLockCount)
		for i, mli := range pso.MNLocks {
			fmt.Fprintln(buf, "  Position:", i)
			fmt.Fprintf(buf, "      Address: '%s'\n", mli.Address.GetString())
			fmt.Fprintln(buf, "       Expire:", mli.Expire, "seconds")
		}
		fmt.Fprintf(buf, "  PSO Count(%d):\n", pso.PSOCount)
	}

	reader, err := m.NewReaderFromStream(
		"PSO",
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
