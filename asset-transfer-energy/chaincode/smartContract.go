package chaincode

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type Asset struct {
	Asset_ID string  `json:"Asset_ID"`
	Owner_ID string  `json:"Owner_ID"`
	Quantity float32 `json:"Quantity"`
}

// submittingClientIdentity is an internal function that retrieves the ID of the submitting client identity
func submittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}

// verifyClientOrgMatchesPeerOrg is an internal function used to verify client org id matches peer org id.
func verifyClientOrgMatchesPeerOrg(ctx contractapi.TransactionContextInterface) error {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	peerMSPID, err := shim.GetMSPID()
	if err != nil {
		return fmt.Errorf("failed getting the peer's MSPID: %v", err)
	}

	if clientMSPID != peerMSPID {
		return fmt.Errorf("client from org %v is not authorized to write private data from an org %v peer", clientMSPID, peerMSPID)
	}

	return nil
}

// verifyClientIDMatchesOwnerID is an internal function used to verify client id matches owner of the asset id.
func verifyClientIDMatchesOwnerID(ctx contractapi.TransactionContextInterface, asset_ID string) error {
	// Retrieve the actual version of the asset to verify that the client who is submitting
	// the update is the owner of the asset
	assetJSON, err := ctx.GetStub().GetState(asset_ID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return fmt.Errorf("the asset %s does not exist", asset_ID)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return err
	}

	// Get ID of submitting client identity
	clientID, err := submittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	// control if the owner of the asset and the submitting client are the same
	if clientID != asset.Owner_ID {
		return fmt.Errorf("the asset owner identity %s is different from the client identity %s", asset.Owner_ID, clientID)
	}

	return nil
}

// InitLedger adds a base set of assets to the ledger
// func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
// 	assets := []Asset{
// 		{ID: "asset1", Color: "blue", Size: 5, Owner: "Tomoko", AppraisedValue: 300},
// 		{ID: "asset2", Color: "red", Size: 5, Owner: "Brad", AppraisedValue: 400},
// 		{ID: "asset3", Color: "green", Size: 10, Owner: "Jin Soo", AppraisedValue: 500},
// 		{ID: "asset4", Color: "yellow", Size: 10, Owner: "Max", AppraisedValue: 600},
// 		{ID: "asset5", Color: "black", Size: 15, Owner: "Adriana", AppraisedValue: 700},
// 		{ID: "asset6", Color: "white", Size: 15, Owner: "Michel", AppraisedValue: 800},
// 	}

// 	for _, asset := range assets {
// 		assetJSON, err := json.Marshal(asset)
// 		if err != nil {
// 			return err
// 		}

// 		err = ctx.GetStub().PutState(asset.ID, assetJSON)
// 		if err != nil {
// 			return fmt.Errorf("failed to put to world state. %v", err)
// 		}
// 	}

// 	return nil
// }

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, asset_ID string, quantity float32) error {

	// Verify that the client is submitting request to peer in their organization
	// This is to ensure that a client from another org doesn't attempt to read or
	// write private data from this peer.
	err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return fmt.Errorf("CreateAsset cannot be performed: Error %v", err)
	}

	exists, err := s.AssetExists(ctx, asset_ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", asset_ID)
	}

	// Get ID of submitting client identity
	clientID, err := submittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	asset := Asset{
		Asset_ID: asset_ID,
		Owner_ID: clientID,
		Quantity: quantity,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(asset_ID, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, asset_ID string) (*Asset, error) {

	// Verify that the client is submitting request to peer in their organization
	// This is to ensure that a client from another org doesn't attempt to read or
	// write private data from this peer.
	err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return nil, fmt.Errorf("ReadAsset cannot be performed: Error %v", err)
	}

	assetJSON, err := ctx.GetStub().GetState(asset_ID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", asset_ID)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, asset_ID string, quantity float32) error {

	// Verify that the client is submitting request to peer in their organization
	// This is to ensure that a client from another org doesn't attempt to read or
	// write private data from this peer.
	err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return fmt.Errorf("UpdateAsset cannot be performed: Error %v", err)
	}

	exists, err := s.AssetExists(ctx, asset_ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", asset_ID)
	}

	// Verify that the client is submitting the update is the owner of the asset
	// This is to ensure that a client cannot modify an asset if it is not the owner.
	err = verifyClientIDMatchesOwnerID(ctx, asset_ID)
	if err != nil {
		return fmt.Errorf("UpdateAsset cannot be performed: Error %v", err)
	}

	// Get ID of submitting client identity
	clientID, err := submittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	// overwriting original asset with new asset
	asset := Asset{
		Asset_ID: asset_ID,
		Owner_ID: clientID,
		Quantity: quantity,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(asset_ID, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, asset_ID string) error {

	// Verify that the client is submitting request to peer in their organization
	// This is to ensure that a client from another org doesn't attempt to read or
	// write private data from this peer.
	err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return fmt.Errorf("DeleteAsset cannot be performed: Error %v", err)
	}

	exists, err := s.AssetExists(ctx, asset_ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", asset_ID)
	}

	// Verify that the client is submitting the update is the owner of the asset
	// This is to ensure that a client cannot modify an asset if it is not the owner.
	err = verifyClientIDMatchesOwnerID(ctx, asset_ID)
	if err != nil {
		return fmt.Errorf("DeleteAsset cannot be performed: Error %v", err)
	}

	return ctx.GetStub().DelState(asset_ID)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, asset_ID string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(asset_ID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, asset_ID string, newOwner_ID string, transfer_quantity float32) (string, error) {

	// Verify that the client is submitting request to peer in their organization
	// This is to ensure that a client from another org doesn't attempt to read or
	// write private data from this peer.
	err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return "", fmt.Errorf("TransferAsset cannot be performed: Error %v", err)
	}

	asset, err := s.ReadAsset(ctx, asset_ID)
	if err != nil {
		return "", err
	}

	// Verify that the client is submitting the update is the owner of the asset
	// This is to ensure that a client cannot modify an asset if it is not the owner.
	err = verifyClientIDMatchesOwnerID(ctx, asset_ID)
	if err != nil {
		return "", fmt.Errorf("TransferAsset cannot be performed: Error %v", err)
	}

	asset_quantity := asset.Quantity

	if asset_quantity == transfer_quantity {

		// oldOwner := asset.Owner_ID
		asset.Owner_ID = newOwner_ID

		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return "", err
		}

		err = ctx.GetStub().PutState(asset_ID, assetJSON)
		if err != nil {
			return "", err
		}

		return "Transferred entire asset", nil

	} else if asset_quantity > transfer_quantity {

		remaining_quantity := asset_quantity - transfer_quantity

		err := s.UpdateAsset(ctx, asset_ID, remaining_quantity)
		if err != nil {
			return "", fmt.Errorf("Old asset cannot be updated: Error %v", err)
		}

		// L'ASSET ID DEVE ESSERE CAMBIATO, MA QUALE USIAMO PER L'ASSET CHE VIENE CREATO IN MODO DA INVIARLO AL COMPRATORE?
		// Create a new asset using the newOwner_ID (create an asset
		// for the buyer instead of transfering it)
		asset := Asset{
			Asset_ID: asset_ID,
			Owner_ID: newOwner_ID,
			Quantity: transfer_quantity,
		}
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return "", err
		}

		// L'ASSET ID DEVE ESSERE CAMBIATO, MA QUALE USIAMO PER L'ASSET CHE VIENE CREATO IN MODO DA INVIARLO AL COMPRATORE?
		return "Transferred part of the asset", ctx.GetStub().PutState(asset_ID, assetJSON)

	}

	return "", nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
