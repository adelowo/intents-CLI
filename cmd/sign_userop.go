package cmd

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/blndgs/model"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"

	"github.com/blndgs/intents-sdk/pkg/config"
	"github.com/blndgs/intents-sdk/pkg/ethclient"
	"github.com/blndgs/intents-sdk/pkg/userop"
	"github.com/blndgs/intents-sdk/utils"
)

// init initializes the signUserOp command and adds it to the root command.
func init() {
	utils.AddCommonFlags(SignUserOpCmd)
}

// SignUserOpCmd represents the command to sign user operations.
var SignUserOpCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a userOp with JSON input",
	Run: func(cmd *cobra.Command, args []string) {
		// Read configuration and initialize necessary components.
		nodeUrl, _, entrypointAddr, eoaSigner := config.ReadConf()
		userOp := utils.GetUserOps(cmd)

		zeroGas := utils.IsZeroGas(cmd)

		fmt.Println("is zero gas enabled: ", zeroGas)

		sender := userOp.Sender

		fmt.Println("sender address: ", sender)

		// Initialize an Ethereum client and retrieve nonce and chain ID.
		ethClient := ethclient.NewClient(nodeUrl)
		// get nonce
		nonce, err := ethClient.GetNonce(sender)
		if err != nil {
			panic(err)
		}

		fmt.Println("nonce: ", nonce)
		unsignedUserOp := utils.UpdateUserOp(userOp, nonce, zeroGas)
		fmt.Println("unsignedUserOp: ", unsignedUserOp.String())
		chainID, err := ethClient.GetChainID(sender)
		if err != nil {
			panic(err)
		}

		fmt.Printf("\nchain-id:%s\n", chainID)
		fmt.Printf("userOp:%s\n\n", unsignedUserOp.GetUserOpHash(entrypointAddr, chainID).String())

		// Sign the user operation and prepare it for sending.
		signUserOp(chainID, entrypointAddr, eoaSigner, userOp)
		// Print signature
		utils.PrintSignature(userOp)
	},
}

// signUserOp signs a user operation using the provided parameters and
// prepares it for sending. It utilizes the userop package for signing.
func signUserOp(chainID *big.Int, entryPointAddr common.Address, signer *signer.EOA, signedUserOp *model.UserOperation) {
	signedOps, err := userop.Sign(chainID, entryPointAddr, signer, signedUserOp)
	if err != nil {
		panic(err)
	}

	fmt.Printf("signed userOp:\n%s\n", signedOps)

	// Marshal signedOps into JSON
	jsonBytes, err := json.Marshal(signedOps)
	if err != nil {
		panic(fmt.Errorf("error marshaling signed operations to JSON: %v", err))
	}

	// Print JSON string
	fmt.Println("signed UserOps in JSON:", string(jsonBytes))
}
