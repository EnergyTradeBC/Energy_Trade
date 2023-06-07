package main

import (
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func getAllEnergyAssets(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: GetAllEnergyAssets, function returns all the current energy assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetAllAssets")

	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

// Evaluate a transaction by assetID to query ledger state.
func readEnergyAssetByID(contract *client.Contract) {
	fmt.Printf("\n--> Evaluate Transaction: ReadAsset, function returns asset attributes\n")

	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", assetID)

	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)

	// return result  ==> ritorno una stringa o per esempio una struct specifica?
	// 					  se invece il risultato è negativo, ovvero non esiste un asset con quell'ID, come si comporta?
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func createEnergyAsset(contract *client.Contract, quantity string) {
	fmt.Printf("\n--> Submit Transaction: CreateEnergyAsset, creates new energy asset with asset_ID and quantity arguments \n")

	// _, err := contract.SubmitTransaction("CreateAsset", asset_ID, quantity)
	// Specifico che l'endorsement deve essere fornito solo dall'organizzazione stessa che sta creando l'asset (nessun altro deve confermare)
	_, err := contract.Submit("CreateAsset",
		client.WithArguments(assetID, quantity),
		client.WithEndorsingOrganizations(orgMSP))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Submit transaction asynchronously, blocking until the transaction has been sent to the orderer, and allowing
// this thread to process the chaincode response (e.g. update a UI) without waiting for the commit notification
func transferAssetAsync(contract *client.Contract, newOwner_ID string, transfer_quantity string) {
	fmt.Printf("\n--> Async Submit Transaction: TransferAsset, transfer part or the entire energy asset to a new owner")

	submitResult, commit, err := contract.SubmitAsync("TransferAsset",
		client.WithArguments(assetID, newOwner_ID, transfer_quantity))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))
	}

	fmt.Printf("\n*** Successfully submitted transaction to transfer ownership from %s to Mark. \n", string(submitResult))
	fmt.Println("*** Waiting for transaction commit.")

	if commitStatus, err := commit.Status(); err != nil {
		panic(fmt.Errorf("failed to get commit status: %w", err))
	} else if !commitStatus.Successful {
		panic(fmt.Errorf("transaction %s failed to commit with status: %d", commitStatus.TransactionID, int32(commitStatus.Code)))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func deleteEnergyAsset(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: DeleteEnergyAsset, deletes an energy asset using its asset_ID \n")

	// _, err := contract.SubmitTransaction("DeleteAsset", asset_ID)
	// Specifico che l'endorsement è richiesto solo all'org che vuole creare la transazione
	_, err := contract.Submit("DeleteAsset",
		client.WithArguments(assetID),
		client.WithEndorsingOrganizations(orgMSP))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}
