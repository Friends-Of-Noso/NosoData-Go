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
	wallet legacy.LegacyWallet

	// walletCmd represents the wallet command
	walletCmd = &cobra.Command{
		Use:   "wallet",
		Short: "Outputs the wallet in text or JSON",
		// 		Long: `A longer description that spans multiple lines and likely contains examples
		// and usage of using your command. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		Example: `  # Display wallet in text format
  $ nosodata wallet --test-data <path to folder containing "wallet.pkw">

  # Display wallet in JSON format
  $ nosodata wallet --json --test-data <path to folder containing "wallet.pkw">`,
		Run: runWallet,
	}
)

func init() {
	rootCmd.AddCommand(walletCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// walletCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// walletCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runWallet(cmd *cobra.Command, args []string) {
	walletFile := filepath.Join(
		testdata,
		"wallet.pkw",
	)

	if err := wallet.ReadFromFile(walletFile); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	displayWallet(json)
}

func displayWallet(jsonOutput bool) {
	buf := new(bytes.Buffer)

	options := m.ReaderOptions{}

	if jsonOutput {
		options.ShouldFormat = true
		options.Style = styles.Get("native")
		fmt.Fprintln(buf, wallet.AsJSON())
	} else {
		options.ShouldFormat = false
		for i, a := range wallet.Accounts {
			fmt.Fprintln(buf, "Position:", i)
			fmt.Fprintf(buf, "    Hash: '%s'\n", a.Hash.GetString())
			fmt.Fprintf(buf, "    Custom:         '%s'\n", a.Custom.GetString())
			fmt.Fprintf(buf, "    Pub key:        '%s'\n", a.PublicKey.GetString())
			fmt.Fprintf(buf, "    Priv key:       '%s'\n", a.PrivateKey.GetString())
			fmt.Fprintln(buf, "    Balance:       ", utils.ToNoso(a.Balance))
			fmt.Fprintln(buf, "    Pending:       ", utils.ToNoso(a.Pending))
			fmt.Fprintln(buf, "    Score:         ", utils.ToNoso(a.Score))
			fmt.Fprintln(buf, "    Last Operation:", utils.ToNoso(a.LastOperation))
		}
	}

	reader, err := m.NewReaderFromStream(
		"Wallet",
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
