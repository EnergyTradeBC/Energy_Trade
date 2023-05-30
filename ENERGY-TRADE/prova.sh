set -x
./network.sh up

./scripts/createChannel.sh genesis -c channeltest
./scripts/createChannel.sh create -c channeltest

./scripts/createChannel.sh join 1 -c channeltest
./scripts/createChannel.sh join 2 -c channeltest

./scripts/addOrg.sh files 3 -c channeltest
./scripts/addOrg.sh add 3 -c channeltest
./scripts/createChannel.sh updateConfig 3 -c channeltest
./scripts/createChannel.sh signConfig 3 1 -c channeltest
./scripts/createChannel.sh submitConfig 3 2 -c channeltest
./scripts/createChannel.sh join 3 -c channeltest

./scripts/addOrg.sh files 4 -c channeltest
./scripts/addOrg.sh add 4 -c channeltest
./scripts/createChannel.sh updateConfig 4 -c channeltest
./scripts/createChannel.sh signConfig 4 1 -c channeltest
./scripts/createChannel.sh signConfig 4 2 -c channeltest
./scripts/createChannel.sh submitConfig 4 3 -c channeltest
./scripts/createChannel.sh join 4 -c channeltest

./scripts/addOrg.sh files 5 -c channeltest
./scripts/addOrg.sh add 5 -c channeltest
./scripts/createChannel.sh updateConfig 5 -c channeltest
./scripts/createChannel.sh signConfig 5 1 -c channeltest
./scripts/createChannel.sh signConfig 5 2 -c channeltest
./scripts/createChannel.sh signConfig 5 3 -c channeltest
./scripts/createChannel.sh submitConfig 5 4 -c channeltest
./scripts/createChannel.sh join 5 -c channeltest

./scripts/deployCC.sh packageChaincode -ccn moneyAsset -ccp ../asset-transfer-money

./scripts/deployCC.sh installChaincode 1 -ccn moneyAsset
./scripts/deployCC.sh installChaincode 2 -ccn moneyAsset
./scripts/deployCC.sh installChaincode 3 -ccn moneyAsset
./scripts/deployCC.sh installChaincode 4 -ccn moneyAsset
./scripts/deployCC.sh installChaincode 5 -ccn moneyAsset

#./scripts/deployCC.sh checkCommitReadiness 1 -ccn moneyAsset -c channeltest

./scripts/deployCC.sh approveForMyOrg 1 -ccn moneyAsset -c channeltest
./scripts/deployCC.sh approveForMyOrg 2 -ccn moneyAsset -c channeltest
./scripts/deployCC.sh approveForMyOrg 3 -ccn moneyAsset -c channeltest
./scripts/deployCC.sh approveForMyOrg 4 -ccn moneyAsset -c channeltest
./scripts/deployCC.sh approveForMyOrg 5 -ccn moneyAsset -c channeltest

#./scripts/deployCC.sh checkCommitReadiness 1 -ccn moneyAsset -c channeltest

./scripts/deployCC.sh commitChaincodeDefinition 1 2 3 4 5 -ccn moneyAsset -c channeltest -ccs 1

# Export environment variables to operate and call the smart contract from the peer of org 5 (MINTER)

export CORE_PEER_LOCALMSPID="Org5MSP"
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org5.example.com/users/Admin@org5.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org5.example.com/tlsca/tlsca.org5.example.com-cert.pem
export CORE_PEER_ADDRESS=localhost:13051
export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls 
--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" 
--peerAddresses localhost:7051 
--tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" 
--peerAddresses localhost:9051 
--tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
--peerAddresses localhost:11051 
--tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt" 
--peerAddresses localhost:12051
--tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org4.example.com/peers/peer0.org4.example.com/tls/ca.crt"
--peerAddresses localhost:13051
--tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org5.example.com/peers/peer0.org5.example.com/tls/ca.crt")

export PATH="/home/bchain/EnergyTrade/bin:$PATH"
export FABRIC_CFG_PATH=/home/bchain/EnergyTrade/config/

set +x