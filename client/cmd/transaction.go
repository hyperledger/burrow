// Copyright 2015, 2016 Eris Industries (UK) Ltd.
// This file is part of Eris-RT

// Eris-RT is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// Eris-RT is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Eris-RT.  If not, see <http://www.gnu.org/licenses/>.

package commands

import (
	"strings"
	"os"

	"github.com/spf13/cobra"

	log "github.com/eris-ltd/eris-logger"

	"github.com/eris-ltd/eris-db/client/transaction"
)

var TransactionCmd = &cobra.Command{
	Use:   "tx",
	Short: "eris-client tx formulates and signs a transaction to a chain",
	Long:  `eris-client tx formulates and signs a transaction to a chain.
`,
	Run:   func(cmd *cobra.Command, args []string) { cmd.Help() },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 
		if !strings.HasPrefix(clientDo.NodeAddrFlag, "http://") {
			clientDo.NodeAddrFlag = "http://" + clientDo.NodeAddrFlag
		}
		if !strings.HasSuffix(clientDo.NodeAddrFlag, "/") {
			clientDo.NodeAddrFlag += "/"
		}

		if !strings.HasPrefix(clientDo.SignAddrFlag, "http://") {
			clientDo.SignAddrFlag = "http://" + clientDo.SignAddrFlag
		}
	},
}

func buildTransactionCommand() {
	// Transaction command has subcommands send, name, call, bond,
	// unbond, rebond, permissions. Dupeout transaction is not accessible through the command line.

	addTransactionPersistentFlags()

	// SendTx
	var sendCmd = &cobra.Command{
		Use:   "send",
		Short: "eris-client tx send --amt <amt> --to <addr>",
		Long:  "eris-client tx send --amt <amt> --to <addr>",
		Run:   func(cmd *cobra.Command, args []string) {
			transaction.Send(clientDo)
		},
		PersistentPreRun: assertChainIdSet,
	}
	sendCmd.Flags().StringVarP(&clientDo.AmtFlag, "amt", "a", "", "specify an amount")
	sendCmd.Flags().StringVarP(&clientDo.ToFlag, "to", "t", "", "specify an address to send to")

	// NameTx
	var nameCmd = &cobra.Command{
		Use:   "name",
		Short: "eris-client tx name --amt <amt> --name <name> --data <data>",
		Long:  "eris-client tx name --amt <amt> --name <name> --data <data>",
		Run:   func(cmd *cobra.Command, args []string) {
			// transaction.Name(clientDo)
		},
		PersistentPreRun: assertChainIdSet,
	}
	nameCmd.Flags().StringVarP(&clientDo.AmtFlag, "amt", "a", "", "specify an amount")
	nameCmd.Flags().StringVarP(&clientDo.NameFlag, "name", "n", "", "specify a name")
	nameCmd.Flags().StringVarP(&clientDo.DataFlag, "data", "d", "", "specify some data")
	nameCmd.Flags().StringVarP(&clientDo.DataFileFlag, "data-file", "", "", "specify a file with some data")
	nameCmd.Flags().StringVarP(&clientDo.FeeFlag, "fee", "f", "", "specify the fee to send")

	// CallTx
	var callCmd = &cobra.Command{
		Use:   "call",
		Short: "eris-client tx call --amt <amt> --fee <fee> --gas <gas> --to <contract addr> --data <data>",
		Long:  "eris-client tx call --amt <amt> --fee <fee> --gas <gas> --to <contract addr> --data <data>",
		Run:   func(cmd *cobra.Command, args []string) {
			// transaction.Call(clientDo)
		},
		PersistentPreRun: assertChainIdSet,
	}
	callCmd.Flags().StringVarP(&clientDo.AmtFlag, "amt", "a", "", "specify an amount")
	callCmd.Flags().StringVarP(&clientDo.ToFlag, "to", "t", "", "specify an address to send to")
	callCmd.Flags().StringVarP(&clientDo.DataFlag, "data", "d", "", "specify some data")
	callCmd.Flags().StringVarP(&clientDo.FeeFlag, "fee", "f", "", "specify the fee to send")
	callCmd.Flags().StringVarP(&clientDo.GasFlag, "gas", "g", "", "specify the gas limit for a CallTx")

	// BondTx
	var bondCmd = &cobra.Command{
		Use:   "bond",
		Short: "eris-client tx bond --pubkey <pubkey> --amt <amt> --unbond-to <address>",
		Long:  "eris-client tx bond --pubkey <pubkey> --amt <amt> --unbond-to <address>",
		Run:   func(cmd *cobra.Command, args []string) {
			// transaction.Bond(clientDo)
		},
		PersistentPreRun: assertChainIdSet,
	}
	bondCmd.Flags().StringVarP(&clientDo.AmtFlag, "amt", "a", "", "specify an amount")
	bondCmd.Flags().StringVarP(&clientDo.UnbondtoFlag, "to", "t", "", "specify an address to unbond to")

	// UnbondTx
	var unbondCmd = &cobra.Command{
		Use:   "unbond",
		Short: "eris-client tx unbond --addr <address> --height <block_height>",
		Long:  "eris-client tx unbond --addr <address> --height <block_height>",
		Run:   func(cmd *cobra.Command, args []string) {
			// transaction.Unbond(clientDo)
		},
		PersistentPreRun: assertChainIdSet,
	}
	unbondCmd.Flags().StringVarP(&clientDo.AddrFlag, "addr", "a", "", "specify an address")
	unbondCmd.Flags().StringVarP(&clientDo.HeightFlag, "height", "n", "", "specify a height to unbond at")

	// RebondTx
	var rebondCmd = &cobra.Command{
		Use:   "rebond",
		Short: "eris-client tx rebond --addr <address> --height <block_height>",
		Long:  "eris-client tx rebond --addr <address> --height <block_height>",
		Run:   func(cmd *cobra.Command, args []string) {
			// transaction.Rebond(clientDo)
		},
		PersistentPreRun: assertChainIdSet,
	}
	rebondCmd.Flags().StringVarP(&clientDo.AddrFlag, "addr", "a", "", "specify an address")
	rebondCmd.Flags().StringVarP(&clientDo.HeightFlag, "height", "n", "", "specify a height to unbond at")

	// PermissionsTx
	var permissionsCmd = &cobra.Command{
		Use:   "permission",
		Short: "eris-client tx perm <function name> <args ...>",
		Long:  "eris-client tx perm <function name> <args ...>",
		Run:   func(cmd *cobra.Command, args []string) {
			// transaction.Permsissions(clientDo)
		},
		PersistentPreRun: assertChainIdSet,
	}
	permissionsCmd.Flags().StringVarP(&clientDo.AddrFlag, "addr", "a", "", "specify an address")
	permissionsCmd.Flags().StringVarP(&clientDo.HeightFlag, "height", "n", "", "specify a height to unbond at")

	TransactionCmd.AddCommand(sendCmd, nameCmd, callCmd, bondCmd, unbondCmd, rebondCmd, permissionsCmd)
}

func addTransactionPersistentFlags() {
	TransactionCmd.PersistentFlags().StringVarP(&clientDo.SignAddrFlag, "sign-addr", "", defaultKeyDaemonAddress(), "set eris-keys daemon address (default respects $ERIS_CLIENT_SIGN_ADDRESS)")
	TransactionCmd.PersistentFlags().StringVarP(&clientDo.NodeAddrFlag, "node-addr", "", defaultNodeRpcAddress(), "set the eris-db node rpc server address (default respects $ERIS_CLIENT_NODE_ADDRESS)")
	TransactionCmd.PersistentFlags().StringVarP(&clientDo.PubkeyFlag, "pubkey", "", defaultPublicKey(), "specify the public key to sign with (defaults to $ERIS_CLIENT_PUBLIC_KEY)")
	TransactionCmd.PersistentFlags().StringVarP(&clientDo.AddrFlag, "addr", "", defaultAddress(), "specify the account address (from which the public key can be fetch from eris-keys) (default respects $ERIS_CLIENT_ADDRESS)")
	TransactionCmd.PersistentFlags().StringVarP(&clientDo.ChainidFlag, "chain-id", "", defaultChainId(), "specify the chainID (default respects $CHAIN_ID)")
	TransactionCmd.PersistentFlags().StringVarP(&clientDo.NonceFlag, "nonce", "", "", "specify the nonce to use for the transaction (should equal the sender account's nonce + 1)")

	// TransactionCmd.PersistentFlags().BoolVarP(&clientDo.SignFlag, "sign", "s", false, "sign the transaction using the eris-keys daemon")
	TransactionCmd.PersistentFlags().BoolVarP(&clientDo.BroadcastFlag, "broadcast", "b", false, "broadcast the transaction to the blockchain")
	TransactionCmd.PersistentFlags().BoolVarP(&clientDo.WaitFlag, "wait", "w", false, "wait for the transaction to be committed in a block")
}

//------------------------------------------------------------------------------
// Defaults

func defaultChainId() string {
	return setDefaultString("CHAIN_ID", "")
}

func defaultKeyDaemonAddress() string {
	return setDefaultString("ERIS_CLIENT_SIGN_ADDRESS", "http://localhost:4767")
}

func defaultNodeRpcAddress() string {
	return setDefaultString("ERIS_CLIENT_NODE_ADDRESS", "http://localhost:46657")
}

func defaultPublicKey() string {
	return setDefaultString("ERIS_CLIENT_PUBLIC_KEY", "")
}

func defaultAddress() string {
	return setDefaultString("ERIS_CLIENT_ADDRESS", "")
}

//------------------------------------------------------------------------------
// Helper functions

func assertChainIdSet(cmd *cobra.Command, args []string) {
	if clientDo.ChainidFlag == "" {
		log.Fatalf(`Cannot run "eris-client tx" without chain id set in --chain-id or $CHAIN_ID.`)
		os.Exit(1)
	}
}