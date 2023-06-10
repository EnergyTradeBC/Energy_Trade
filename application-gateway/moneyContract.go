package main

import (
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Initializes the money smart contract passing 3 settings (name, symbol and decimals).
// Can be used only by the org that will behave as the central banker
func initializeMoneyContract(contract *client.Contract, quantity string) {
	fmt.Printf("\n--> Submit Transaction: initializeMoneyContract, initializes the money contract with 'name', 'symbol' and 'decimals' \n")

	_, err := contract.Submit("Initialize",
		client.WithArguments(moneyName, moneySymbol, moneyDecimals))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Evaluate a transaction requesting the client ID of the caller.
func getAccountID(contract *client.Contract) string {
	fmt.Printf("\n--> Evaluate Transaction: getAccountID, retrieves the client ID of the caller \n")

	evaluateResult, err := contract.EvaluateTransaction("ClientAccountID")

	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	// Rimane però codificato BASE-64, va bene?
	result := string(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)

	return result
}

// Evaluate a transaction requesting the money balance of the caller.
func getAccountBalance(contract *client.Contract) string {
	fmt.Printf("\n--> Evaluate Transaction: getAccountID, retrieves the client ID of the caller \n")

	evaluateResult, err := contract.EvaluateTransaction("ClientAccountBalance")

	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := string(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)

	return result
}

// Creates <amount> of new tokens and then transfers them to the recipient (clientID of the recipient).
// Can be used only by the org that will behave as the central banker
func mintAndTransfer(contract *client.Contract, recipientID string, amount string) {
	fmt.Printf("\n--> Submit Transaction: mintAndTransfer, creates new tokens and transfers them to the recipient \n")

	_, err := contract.Submit("Mint",
		client.WithArguments(recipientID, amount))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Transfers <amount> of tokens to the new recipient (clientID of the recipient).
func transferMoney(contract *client.Contract, recipientID string, amount string) {
	fmt.Printf("\n--> Submit Transaction: transferMoney, transfer amount of tokens from the caller to a new recipient \n")

	_, err := contract.Submit("Transfer",
		client.WithArguments(recipientID, amount))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}
