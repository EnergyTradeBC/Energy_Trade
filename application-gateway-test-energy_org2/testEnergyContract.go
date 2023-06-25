/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	mspID        = "Org2MSP"
	cryptoPath   = "../ENERGY-TRADE/organizations/peerOrganizations/org2.example.com"
	certPath     = cryptoPath + "/users/User1@org2.example.com/msp/signcerts/User1@org2.example.com-cert.pem"
	keyPath      = cryptoPath + "/users/User1@org2.example.com/msp/keystore/"
	tlsCertPath  = cryptoPath + "/peers/peer0.org2.example.com/tls/ca.crt"
	peerEndpoint = "localhost:9051"
	gatewayPeer  = "peer0.org2.example.com"
	assetID      = "energy_2"
)

func main() {

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
	}

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

	// Override default values for chaincode and channel name as they may differ in testing contexts.
	chaincodeName := "chaincodetest"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "channeltest"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	// createEnergyAssetAsync(contract, "10")
	getAllEnergyAssets(contract)
	// readEnergyAssetByID(contract)
	// transferEnergyAssetAsync(contract, "pippo", "5")
	// getAllEnergyAssets(contract)

	time.Sleep(30 * time.Second)
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

// Evaluate a transaction to query ledger state.

// GENERA ERRORE QUANDO SI FA UNA QUERY E IL LEDGER STATE E' VUOTO (ERRORE DI FORMATTAZIONE JSON)
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
		client.WithArguments(assetID, quantity))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func createEnergyAssetAsync(contract *client.Contract, quantity string) {
	fmt.Printf("\n--> Async Submit Transaction: Create Asset, create a new asset\n")

	submitResult, commit, err := contract.SubmitAsync("CreateAsset",
		client.WithArguments(assetID, quantity))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))
	}

	fmt.Printf("\n*** Successfully submitted transaction to create asset %s. \n", assetID)
	fmt.Printf("\n*** SubmitResult content: %s.\n", submitResult)
	fmt.Println("*** Waiting for transaction commit.")

	if commitStatus, err := commit.Status(); err != nil {
		panic(fmt.Errorf("failed to get commit status: %w", err))
	} else if !commitStatus.Successful {
		panic(fmt.Errorf("transaction %s failed to commit with status: %d", commitStatus.TransactionID, int32(commitStatus.Code)))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Submit transaction asynchronously, blocking until the transaction has been sent to the orderer, and allowing
// this thread to process the chaincode response (e.g. update a UI) without waiting for the commit notification
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
		client.WithEndorsingOrganizations(mspID))

	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
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
