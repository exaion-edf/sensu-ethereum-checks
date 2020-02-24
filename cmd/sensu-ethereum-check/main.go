package main

import (
	"encoding/json"
	"fmt"
	"github.com/onrik/ethrpc"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

type Block struct {
	Number           string        `json:"number"`
}

var (
	ethUrl                    string
	minPairCrit               int
	minPairWarn               int
	minerAddr                 string
	maxBlocks                 int
	maxMinuteWithoutBlockWarn float64
	maxMinuteWithoutBlockCrit float64
)

func configureRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sensu-ethereum-check",
		Short: "The Sensu Go Ethereum Check plugin",
		RunE:  run,
	}

	cmd.Flags().StringVarP(&ethUrl,
		"rpc-url",
		"u",
		"http://127.0.0.1:8545",
		"Ethereum RPC URL")

	cmd.Flags().IntVarP(&minPairWarn,
		"warn-peers",
		"p",
		0,
		"Warning eth pairs amount")

	cmd.Flags().IntVarP(&minPairCrit,
		"crit-peers",
		"P",
		0,
		"Critical eth pairs amount")

	cmd.Flags().StringVarP(&minerAddr,
		"miner-addr",
		"a",
		"0x00",
		"Miner address")

	cmd.Flags().IntVarP(&maxBlocks,
		"max-blocks",
		"x",
		100,
		"Max blocks to check for address")

	cmd.Flags().Float64VarP(&maxMinuteWithoutBlockWarn,
		"warn-max-time-without-block",
		"b",
		10.0,
		"Warning max minutes without block")

	cmd.Flags().Float64VarP(&maxMinuteWithoutBlockCrit,
		"crit-max-time-without-block",
		"B",
		20.0,
		"Critcal max minutes without block")

	return cmd
}

func main() {
	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		_ = cmd.Help()
		return fmt.Errorf("invalid argument(s) received")
	}

	client := ethrpc.New(ethUrl)

	err := checkPeers(client)
	if err != nil {
		panic(err)
	}

	if minerAddr != "0x00" {
		err := checkNetworkConnection(client)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func checkPeers(client *ethrpc.EthRPC) error {
	peer, err := client.NetPeerCount()
	if err != nil {
		fmt.Printf("CRITICAL: %s not answering to RPC requests\n", client.URL())
		os.Exit(2)
	}

	if peer <= minPairCrit {
		fmt.Printf("CRITICAL: %d peers\n", peer)
		os.Exit(2)
	} else if peer <= minPairWarn {
		fmt.Printf("WARNING: %d peers\n", peer)
		os.Exit(1)
	}

	fmt.Printf("%d peers\n", peer)

	return nil
}

func checkNetworkConnection(client *ethrpc.EthRPC) error {
	var latest Block
	latestData, err := client.Call("eth_getBlockByNumber", "latest", true)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(latestData, &latest)
	if err != nil {
		panic(err)
	}

	n, err := strconv.ParseInt(latest.Number, 0, 64)
	if err != nil {
		panic(err)
	}
	nb := int(n)

	searching := true
	blocks := 0
	for searching && blocks < maxBlocks {
		b, err := client.EthGetBlockByNumber(nb, false)
		if err != nil {
			panic(err)
		}

		if b.Miner == minerAddr {
			t := time.Unix(int64(b.Timestamp), 0)
			diff := time.Now().Sub(t)

			if diff.Minutes() > maxMinuteWithoutBlockCrit {
				fmt.Printf("No new block seen since %f minutes\n", diff.Minutes())
				os.Exit(2)
			} else if diff.Minutes() > maxMinuteWithoutBlockWarn {
				fmt.Printf("No new block seen since %f minutes\n", diff.Minutes())
				os.Exit(1)
			}

			fmt.Printf("Last block seen %f minutes ago\n", diff.Minutes())
			return nil
		}

		blocks++
		nb--
	}

	if blocks >= maxBlocks {
		fmt.Printf("No block seen in last %d blocks\n", maxBlocks)
		os.Exit(2)
	}

	return nil
}
