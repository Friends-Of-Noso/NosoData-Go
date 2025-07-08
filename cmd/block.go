/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Friends-Of-Noso/NosoData-Go/legacy"
	"github.com/Friends-Of-Noso/NosoData-Go/utils"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/spf13/cobra"
	"github.com/walles/moar/m"
)

var (
	block legacy.LegacyBlock

	// blockCmd represents the block command
	blockCmd = &cobra.Command{
		Use:   "block [flags] <BLOCK NUMBER>",
		Short: "Outputs the block in JSON",
		Args:  cobra.MinimumNArgs(1),
		// 	Long: `A longer description that spans multiple lines and likely contains examples
		// and usage of using your command. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		Example: `  # Display block 100000
  $ nosodata block --test-data <path to folder containing "100000.blk"> 100000

  # Display same block in JSON
  $ nosodata block --json --test-data <path to folder containing "100000.blk"> 100000`,
		Run: runBlock,
	}
)

func init() {
	rootCmd.AddCommand(blockCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// blockCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// blockCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runBlock(cmd *cobra.Command, args []string) {
	blockNumber := args[0]
	number, err := strconv.ParseInt(blockNumber, 10, 64)
	if err != nil {
		fmt.Printf("could not convert '%s' to an integer: %v", blockNumber, err)
		os.Exit(1)
	}

	blockFile := filepath.Join(
		testdata,
		fmt.Sprintf("%d.blk", number),
	)

	if err := block.ReadFromFile(blockFile); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	displayBlock(block, json)
}

func displayBlock(block legacy.LegacyBlock, jsonOutput bool) {
	buf := new(bytes.Buffer)

	options := m.ReaderOptions{}

	if jsonOutput {
		options.ShouldFormat = true
		options.Style = styles.Get("native")
		fmt.Fprintln(buf, block.AsJSON())
	} else {
		options.ShouldFormat = false
		fmt.Fprintln(buf, "Number:           ", block.Number)
		fmt.Fprintf(buf, "HASH:              '%s'\n", block.HASH)
		fmt.Fprintln(buf, "Time Start:       ", time.Unix(block.TimeStart, 0))
		fmt.Fprintln(buf, "Time End:         ", time.Unix(block.TimeEnd, 0))
		fmt.Fprintln(buf, "Time Total:       ", block.TimeTotal, "seconds")
		fmt.Fprintln(buf, "Time Last 20:     ", block.TimeLast20, "seconds")
		fmt.Fprintln(buf, "Transaction Count:", block.TransactionsCount)
		fmt.Fprintln(buf, "Difficulty:       ", block.Difficulty)
		fmt.Fprintf(buf, "Target Hash:      '%s'\n", block.TargetHash.GetString())
		fmt.Fprintf(buf, "Solution:         '%s'\n", block.Solution.GetString())
		fmt.Fprintf(buf, "Last Block Hash:  '%s'\n", block.LastBlockHash.GetString())
		fmt.Fprintf(buf, "Miner:            '%s'\n", block.Miner.GetString())
		fmt.Fprintln(buf, "Fee:             ", utils.ToNoso(block.Fee))
		fmt.Fprintln(buf, "Reward:          ", utils.ToNoso(block.Reward))

		if block.TransactionsCount > 0 {
			fmt.Fprintf(buf, "Transactions(%d):\n", block.TransactionsCount)
			for n := range block.TransactionsCount {
				fmt.Fprintf(buf, "  OrderID: '%s'\n", block.Transactions[n].OrderID.GetString())
				fmt.Fprintf(buf, "      TransferID:     '%s'\n", block.Transactions[n].TransferID.GetString())
				fmt.Fprintln(buf, "      Block:         ", block.Transactions[n].Block)
				fmt.Fprintln(buf, "      Order lines:   ", block.Transactions[n].OrderLinesCount)
				fmt.Fprintf(buf, "      Order type:     '%s'\n", block.Transactions[n].OrderType.GetString())
				fmt.Fprintln(buf, "      Timestamp:     ", time.Unix(block.Transactions[n].TimeStamp, 0))
				fmt.Fprintf(buf, "      Reference:      '%s'\n", block.Transactions[n].Reference.GetString())
				fmt.Fprintln(buf, "      Transfer Index:", block.Transactions[n].TransferIndex)
				fmt.Fprintf(buf, "      Sender:         '%s'\n", block.Transactions[n].Sender.GetString())
				fmt.Fprintf(buf, "      Address:        '%s'\n", block.Transactions[n].Address.GetString())
				fmt.Fprintf(buf, "      Receiver:       '%s'\n", block.Transactions[n].Receiver.GetString())
				fmt.Fprintln(buf, "      Fee:           ", utils.ToNoso(block.Transactions[n].AmountFee))
				fmt.Fprintln(buf, "      Value:         ", utils.ToNoso(block.Transactions[n].AmountTransfer))
				fmt.Fprintf(buf, "      Signature:      '%s'\n", block.Transactions[n].Signature.GetString())
			}
		} else {
			fmt.Fprintln(buf, "No transactions")
		}

		if block.ProofOfStakeRewardCount > 0 {
			fmt.Fprintf(buf, "PoS rewards(%d):\n", block.ProofOfStakeRewardCount)
			fmt.Fprintln(buf, "  Amount:", utils.ToNoso(block.ProofOfStakeRewardAmount))
			for n := range block.ProofOfStakeRewardCount {
				fmt.Fprintf(buf, "  Address: '%s'\n", block.ProofOfStakeRewardAddresses[n].GetString())
			}
		} else {
			fmt.Fprintln(buf, "No PoS rewards")
		}

		if block.MasterNodeRewardCount > 0 {
			fmt.Fprintf(buf, "MN rewards(%d):\n", block.MasterNodeRewardCount)
			fmt.Fprintln(buf, "  Amount:", utils.ToNoso(block.MasterNodeRewardAmount))
			for n := range block.MasterNodeRewardCount {
				fmt.Fprintf(buf, "  Address: '%s'\n", block.MasterNodeRewardAddresses[n].GetString())
			}
		} else {
			fmt.Fprintln(buf, "No MN rewards")
		}
	}

	reader, err := m.NewReaderFromStream(
		fmt.Sprintf("Block: %d", block.Number),
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
