#!/bin/bash

# imports  
. scripts/envVar.sh
. scripts/utils.sh
. scripts/actAsOrg.sh
. scripts/createChannelHelp.sh
. scripts/createChannelUtils.sh

if [ ! -d "channel-artifacts" ]; then
	mkdir channel-artifacts
fi

if [[ $# -lt 1 ]] ; then
  	createChannelHelp
  	exit 0
else
  	MODE=$1
  	shift
fi

if [ "$MODE" == "updateConfig" ]; then
	if [[ $# -lt 1 ]] ; then
		createChannelHelp $MODE
		exit 0
	else
		ORG=$1
		shift
	fi
fi

if [ "$MODE" == "signConfig" ]; then
	if [[ $# -lt 2 ]] ; then
		createChannelHelp $MODE
		exit 0
	else
		ORG=$1
		shift
		AS_ORG=$1
		shift
	fi
fi

if [ "$MODE" == "submitConfig" ]; then
	if [[ $# -lt 2 ]] ; then
		createChannelHelp $MODE
		exit 0
	else
		ORG=$1
		shift
		AS_ORG=$1
		shift
	fi
fi

if [ "$MODE" == "join" ]; then
	if [[ $# -lt 1 ]] ; then
		createChannelHelp $MODE
		exit 0
	else
		ORG=$1
		shift
	fi
fi

CONTAINER_CLI="docker"
CONTAINER_CLI_COMPOSE="${CONTAINER_CLI}-compose"
infoln "Using ${CONTAINER_CLI} and ${CONTAINER_CLI_COMPOSE}"

CHANNEL_NAME="cer"
DELAY="3"
MAX_RETRY="5"
VERBOSE="false"

while [[ $# -ge 1 ]] ; do
  	key="$1"
	shift
  	case $key in
  	-h)
    	createChannelHelp $MODE
    	exit 0
    	;;
  	-c)
		if [[ $# -ge 1 ]] ; then
    		CHANNEL_NAME="$1"
			shift
		else
			fatalln "No channel name provided after flag '-c'"
		fi
    	;;
	-d)
		if [[ $# -ge 1 ]] ; then
    		DELAY="$1"
    		shift
		else
			fatalln "No delay value provided after flag '-d'"
		fi
    	;;
  	-r)
		if [[ $# -ge 1 ]] ; then
    		MAX_RETRY="$1"
    		shift
		else
			fatalln "No max retry value provided after flag '-r'"
		fi
    	;;
  	-verbose)
    	VERBOSE=true
    	;;
  	*)
    	errorln "Unknown flag: $key"
    	createChannelHelp
    	exit 1
    	;;
  	esac
done

# Determine mode of operation and printing out what we asked for
if [ "$MODE" == "genesis" ]; then
	infoln "Generating channel genesis block '${CHANNEL_NAME}.block'"
  	createChannelGenesisBlock
elif [ "$MODE" == "create" ]; then
  	infoln "Creating channel ${CHANNEL_NAME}"
	createChannel
  	successln "Channel '$CHANNEL_NAME' created"
elif [ "$MODE" == "updateConfig" ]; then
  	infoln "Updating config transaction to add org${ORG} to channel ${CHANNEL_NAME}"
  	updateConfig
elif [ "$MODE" == "signConfig" ]; then
  	infoln "Signing config from org${AS_ORG} to add org${ORG} to channel ${CHANNEL_NAME}"
  	signConfig
elif [ "$MODE" == "submitConfig" ]; then
  	infoln "Sumbitting config from org${AS_ORG} to add org${ORG} to channel ${CHANNEL_NAME}"
  	submitConfig
elif [ "$MODE" == "join" ]; then
  	infoln "Joining ${ORG} peer to channel ${CHANNEL_NAME}"
  	joinChannel
elif [ "$MODE" == "anchor" ]; then
  	infoln "Setting anchor peer for org1..."
  	setAnchorPeer
else
  	createChannelHelp
  	exit 1
fi