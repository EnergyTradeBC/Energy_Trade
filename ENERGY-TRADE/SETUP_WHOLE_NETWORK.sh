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

./scripts/deployCC.sh packageChaincode -ccn chaincodetest -ccp ../asset-transfer-energy

./scripts/deployCC.sh installChaincode 1 -ccn chaincodetest
./scripts/deployCC.sh installChaincode 2 -ccn chaincodetest
./scripts/deployCC.sh installChaincode 3 -ccn chaincodetest
./scripts/deployCC.sh installChaincode 4 -ccn chaincodetest
./scripts/deployCC.sh installChaincode 5 -ccn chaincodetest

#./scripts/deployCC.sh checkCommitReadiness 1 -ccn chaincodetest -c channeltest

./scripts/deployCC.sh approveForMyOrg 1 -ccn chaincodetest -c channeltest
./scripts/deployCC.sh approveForMyOrg 2 -ccn chaincodetest -c channeltest
./scripts/deployCC.sh approveForMyOrg 3 -ccn chaincodetest -c channeltest
./scripts/deployCC.sh approveForMyOrg 4 -ccn chaincodetest -c channeltest
./scripts/deployCC.sh approveForMyOrg 5 -ccn chaincodetest -c channeltest

#./scripts/deployCC.sh checkCommitReadiness 1 -ccn chaincodetest -c channeltest

./scripts/deployCC.sh commitChaincodeDefinition 1 2 3 4 5 -ccn chaincodetest -c channeltest

# Export environment variables to operate and call the smart contract from the peer of org 1
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

set +x