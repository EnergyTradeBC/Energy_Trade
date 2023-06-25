package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type Asset struct {
	Asset_ID string `json:"Asset_ID"`
	Owner_ID string `json:"Owner_ID"`
	Quantity string `json:"Quantity"`
}

type EnergyEvent struct {
	New_OwnerID string  `json:"New_OwnerID"`
	Quantity    float32 `json:"Quantity"`
}

// Retrieves all the energy assets present on the ledger.
// func getAllEnergyAssets(contract *client.Contract) string {
// 	fmt.Println("\n--> Evaluate Transaction: GetAllEnergyAssets, function returns all the current energy assets on the ledger")

// 	evaluateResult, err := contract.EvaluateTransaction("GetAllAssets")

// 	if err != nil {
// 		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
// 	}
// 	result := formatJSON(evaluateResult)

// 	fmt.Printf("*** Result:%s\n", result)

// 	return result
// }

// Retrieves the asset, if existing, with the assetID and returns the quantity.
func readEnergyAssetByID(contract *client.Contract) string {
	fmt.Printf("\n--> Evaluate Transaction: ReadAsset, function returns asset attributes\n")

	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", assetID)

	if err != nil {
		// Used as signal that there is no asset with that ID
		return "failed"
	}
	result := formatJSON(evaluateResult)

	var asset Asset
	err = json.Unmarshal(evaluateResult, &asset)
	if err != nil {
		return "Error during unmarshal"
	}

	fmt.Printf("*** Result:%s\n", result)

	return asset.Quantity // se invece il risultato è negativo, ovvero non esiste un asset con quell'ID, come si comporta?
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
// Creates a new asset with <assetID> and quantity received as input
func createEnergyAsset(contract *client.Contract, quantity string) {
	fmt.Printf("\n--> Submit Transaction: CreateEnergyAsset, creates new energy asset with asset_ID and quantity arguments \n")

	// _, err := contract.SubmitTransaction("CreateAsset", asset_ID, quantity)
	// Specifico che l'endorsement deve essere fornito solo dall'organizzazione stessa che sta creando l'asset (nessun altro deve confermare)
	_, err := contract.Submit("CreateAsset",
		client.WithArguments(assetID, quantity))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Submit transaction asynchronously, blocking until the transaction has been sent to the orderer, and allowing
// this thread to process the chaincode response (e.g. update a UI) without waiting for the commit notification
// Transfers a specific quantity of an asset to a new owner with <newOwner_ID>
func transferEnergyAssetAsync(contract *client.Contract, newOwner_ID string, transfer_quantity string) {
	fmt.Printf("\n--> Async Submit Transaction: TransferAsset, transfer part or the entire energy asset to a new owner")

	submitResult, commit, err := contract.SubmitAsync("TransferAsset",
		client.WithArguments(assetID, newOwner_ID, transfer_quantity))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))
	}

	fmt.Printf("\n*** Successfully submitted transaction to transfer ownership from %s to %s. \n", string(submitResult), newOwner_ID)
	fmt.Println("*** Waiting for transaction commit.")

	if commitStatus, err := commit.Status(); err != nil {
		panic(fmt.Errorf("failed to get commit status: %w", err))
	} else if !commitStatus.Successful {
		panic(fmt.Errorf("transaction %s failed to commit with status: %d", commitStatus.TransactionID, int32(commitStatus.Code)))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
// Deletes the asse twith <assetID>
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
