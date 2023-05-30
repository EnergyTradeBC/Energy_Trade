# createChannelGenesisBlock CHANNEL_NAME
createChannelGenesisBlock() {
  	export PATH=${PWD}/../bin:$PATH
	export FABRIC_CFG_PATH=${PWD}/configtx
	which configtxgen
	if [ "$?" -ne 0 ]; then
		fatalln "configtxgen tool not found."
	fi
	set -x
	configtxgen -profile TwoOrgsApplicationGenesis -outputBlock ./channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
	res=$?
	{ set +x; } 2>/dev/null
  	verifyResult $res "Failed to generate channel configuration transaction..."
}

# createChannel CHANNEL_NAME DELAY MAX_RETRY
createChannel() {
  	export PATH=${PWD}/../bin:$PATH
	export FABRIC_CFG_PATH=$PWD/../config/
	#actAsOrg 1
	# Poll in case the raft leader is not set yet
	local rc=1
	local COUNTER=1
	while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ] ; do
		sleep $DELAY
		set -x
		osnadmin channel join --channelID $CHANNEL_NAME --config-block ./channel-artifacts/${CHANNEL_NAME}.block -o localhost:7053 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY" >&log.txt
		res=$?
		{ set +x; } 2>/dev/null
		let rc=$res
		COUNTER=$(expr $COUNTER + 1)
	done
	cat log.txt
	verifyResult $res "Channel creation failed"
}

updateConfig() {
	infoln "Generating config tx to add Org${ORG}"
    ${CONTAINER_CLI} exec cli ./scripts/org-scripts/updateChannelConfig.sh $ORG $CHANNEL_NAME $CLI_DELAY $CLI_TIMEOUT $VERBOSE
    if [ $? -ne 0 ]; then
        fatalln "ERROR !!!! Unable to create config tx"
    fi
}

signConfig() {
	infoln "Signing config tx to add Org${ORG}"
    ${CONTAINER_CLI} exec cli ./scripts/org-scripts/signChannelConfig.sh $ORG $AS_ORG $CHANNEL_NAME $CLI_DELAY $CLI_TIMEOUT $VERBOSE
    if [ $? -ne 0 ]; then
        fatalln "ERROR !!!! Unable to sign config tx"
    fi
}

submitConfig() {
	infoln "Submitting config tx to add Org${ORG}"
    ${CONTAINER_CLI} exec cli ./scripts/org-scripts/submitChannelConfig.sh $ORG $AS_ORG $CHANNEL_NAME $CLI_DELAY $CLI_TIMEOUT $VERBOSE
    if [ $? -ne 0 ]; then
        fatalln "ERROR !!!! Unable to sumbit config tx"
    fi
}

# joinChannel CHANNEL_NAME DELAY MAX_RETRY ORG
joinChannel() {
  	export PATH=${PWD}/../bin:$PATH
  	export FABRIC_CFG_PATH=$PWD/../config/
  	BLOCKFILE="./channel-artifacts/${CHANNEL_NAME}.block"
  	actAsOrg $ORG
	local rc=1
	local COUNTER=1
	## Sometimes Join takes time, hence retry
	while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ] ; do
    sleep $DELAY
    set -x
    peer channel join -b $BLOCKFILE >&log.txt
    res=$?
    { set +x; } 2>/dev/null
		let rc=$res
		COUNTER=$(expr $COUNTER + 1)
	done
	cat log.txt
	verifyResult $res "After $MAX_RETRY attempts, peer0.org${ORG} has failed to join channel '$CHANNEL_NAME' "
}

# setAnchorPeer CHANNEL_NAME ORG
setAnchorPeer() {
  	${CONTAINER_CLI} exec cli ./scripts/setAnchorPeer.sh $ORG $CHANNEL_NAME 
}