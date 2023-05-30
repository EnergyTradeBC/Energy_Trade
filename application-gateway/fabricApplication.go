/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type NetworkConfig struct {
	ChannelName         string
	ChaincodeMoneyName  string
	ChaincodeEnergyName string
}

// const (
// 	mspID        = "Org1MSP"
// 	cryptoPath   = "../../test-network/organizations/peerOrganizations/org1.example.com"
// 	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem"
// 	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
// 	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
// 	peerEndpoint = "localhost:7051"
// 	gatewayPeer  = "peer0.org1.example.com"
// )

func main() {
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	channelName, chaincodeMoneyName, chaincodeEnergyName := readNetworkConfig()

	if cc_money_name := os.Getenv("CHAINCODE_MONEY_NAME"); cc_money_name != "" {
		chaincodeMoneyName = cc_money_name
	}

	if cc_energy_name := os.Getenv("CHAINCODE_ENERGY_NAME"); cc_energy_name != "" {
		chaincodeEnergyName = cc_energy_name
	}

	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract_money := network.GetContract(chaincodeMoneyName)
	contract_energy := network.GetContract(chaincodeEnergyName)

	for true {

	}

}

// UTILITY FUNCTIONS

// Reads from a json the channel name and the name of the chaincodes committed to that channel
func readNetworkConfig() (string, string, string) {
	content, err := ioutil.ReadFile("./network_config.json")
	if err != nil {
		panic(fmt.Errorf("Error when opening file: %w", err))
	}

	var payload NetworkConfig
	err = json.Unmarshal(content, &payload)
	if err != nil {
		panic(fmt.Errorf("Error during Unmarshal(): %w", err))
	}

	return payload.ChannelName, payload.ChaincodeMoneyName, payload.ChaincodeEnergyName
}

// // newGrpcConnection creates a gRPC connection to the Gateway server.
// func newGrpcConnection() *grpc.ClientConn {
// 	certificate, err := loadCertificate(tlsCertPath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	certPool := x509.NewCertPool()
// 	certPool.AddCert(certificate)
// 	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

// 	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
// 	if err != nil {
// 		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
// 	}

// 	return connection
// }

// // newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
// func newIdentity() *identity.X509Identity {
// 	certificate, err := loadCertificate(certPath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	id, err := identity.NewX509Identity(mspID, certificate)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return id
// }

// func loadCertificate(filename string) (*x509.Certificate, error) {
// 	certificatePEM, err := os.ReadFile(filename)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read certificate file: %w", err)
// 	}
// 	return identity.CertificateFromPEM(certificatePEM)
// }

// // newSign creates a function that generates a digital signature from a message digest using a private key.
// func newSign() identity.Sign {
// 	files, err := os.ReadDir(keyPath)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to read private key directory: %w", err))
// 	}
// 	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))

// 	if err != nil {
// 		panic(fmt.Errorf("failed to read private key file: %w", err))
// 	}

// 	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
// 	if err != nil {
// 		panic(err)
// 	}

// 	sign, err := identity.NewPrivateKeySign(privateKey)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return sign
// }

// FUNCTIONS TO MANAGE THE ENERGY CONTRACT

// Evaluate a transaction to query ledger state.
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
func readEnergyAssetByID(contract *client.Contract, asset_ID string) {
	fmt.Printf("\n--> Evaluate Transaction: ReadAsset, function returns asset attributes\n")

	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", asset_ID)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func createEnergyAsset(contract *client.Contract, asset_ID string, quantity string) {
	fmt.Printf("\n--> Submit Transaction: CreateEnergyAsset, creates new energy asset with asset_ID and quantity arguments \n")

	_, err := contract.SubmitTransaction("CreateAsset", asset_ID, quantity)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Submit transaction asynchronously, blocking until the transaction has been sent to the orderer, and allowing
// this thread to process the chaincode response (e.g. update a UI) without waiting for the commit notification
func transferAssetAsync(contract *client.Contract, asset_ID string, newOwner_ID string, transfer_quantity string) {
	fmt.Printf("\n--> Async Submit Transaction: TransferAsset, transfer part or the entire energy asset to a new owner")

	submitResult, commit, err := contract.SubmitAsync("TransferAsset", client.WithArguments(asset_ID, newOwner_ID, transfer_quantity))
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

// Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
