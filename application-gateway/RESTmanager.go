package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MONEY STRUCTS

type MoneyTransactions_struct struct {
	Transactions []float64
	Price        float32
	Timestamp    string
}

type GET_moneyTransactions struct {
	TransactionsArray []*MoneyTransactions_struct
}

// ENERGY STRUCTS

type EnergyTransactions_struct struct {
	Transactions []float64
	Timestamp    string
}

type GET_energyTransactions struct {
	TransactionsArray []*EnergyTransactions_struct
}

type GET_remainingAssets_struct struct {
	AssetArray []*RemainingAsset
}

var moneyTransactionsArray GET_moneyTransactions
var energyTransactionsArray GET_energyTransactions
var remainingAssetsArray GET_remainingAssets_struct

// MONEY SETTINGS

// getMoneySettings responds with struct which contains the settings for the money management
func getMoneySettings(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, moneySettings)
}

// postMoneySettings retrieves the settings from the mobile application and write them into the money settings structure
func postMoneySettings(c *gin.Context) {
	// Call BindJSON to bind the received JSON to the money settings struct
	if err := c.BindJSON(&moneySettings); err != nil {
		return
		// Can we do it directly or it is necessary to creare first a transient variable and then overwrite the struct?
	}

	c.IndentedJSON(http.StatusCreated, moneySettings)
}

// MONEY ASSET and TRANSFER

// getCurrentBalance responds with the total actual money balance
func getCurrentBalance(c *gin.Context) {
	currentBalance := getAccountBalance(contract_money)

	c.IndentedJSON(http.StatusOK, currentBalance)
}

// getMoneyTransactions responds with a struct that contains an array of money transactions (sent and received)
// with timestamp and the prize, which is referred to the final prize of the auction that causes the transfer
func getMoneyTransactions(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, moneyTransactionsArray)
}

// ENERGY ASSET and TRANSFER

// getCurrentEnergyAsset responds with the current energy asset
func getCurrentEnergyAsset(c *gin.Context) {
	currentAsset := readEnergyAssetByID(contract_energy)

	c.IndentedJSON(http.StatusOK, currentAsset)
}

// getEnergyTransactions responds with a struct that contains an array of energy transactions (sent and received) with timestamp
func getEnergyTransactions(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, energyTransactionsArray)
}

// getRemainingEnergy responds with a struct that contains an array of remaining energy with timestamp
func getRemainingEnergy(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, remainingAssetsArray)
}
