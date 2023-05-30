#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

# This script brings up a Hyperledger Fabric network for testing smart contracts
# and applications. The test network consists of two organizations with one
# peer each, and a single node Raft ordering service. Users can also use this
# script to create a channel deploy a chaincode on the channel
#
# prepending $PWD/../bin to PATH to ensure we are picking up the correct binaries
# this may be commented out to resolve installed version of tools if desired
#
# However using PWD in the path has the side effect that location that
# this script is run from is critical. To ease this, get the directory
# this script is actually in and infer location from there. (putting first)

ROOTDIR=$(cd "$(dirname "$0")" && pwd)
export PATH=${ROOTDIR}/../bin:${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx
export VERBOSE=false

# push to the required directory & set a trap to go back if needed
pushd ${ROOTDIR} > /dev/null
trap "popd > /dev/null" EXIT

. scripts/utils.sh
. networkUtils.sh

: ${CONTAINER_CLI:="docker"}
: ${CONTAINER_CLI_COMPOSE:="${CONTAINER_CLI}-compose"}
infoln "Using ${CONTAINER_CLI} and ${CONTAINER_CLI_COMPOSE}"

# Using crpto vs CA. default is cryptogen
CRYPTO="cryptogen"
# timeout duration - the duration the CLI should wait for a response from
# another container before giving up
MAX_RETRY=5
# default for delay between commands
CLI_DELAY=3
# channel name defaults to "mychannel"
CHANNEL_NAME="mychannel"
# chaincode name defaults to "NA"
CC_NAME="NA"
# chaincode path defaults to "NA"
CC_SRC_PATH="NA"
# endorsement policy defaults to "NA". This would allow chaincodes to use the majority default policy.
CC_END_POLICY="NA"
# collection configuration defaults to "NA"
CC_COLL_CONFIG="NA"
# chaincode init function defaults to "NA"
CC_INIT_FCN="NA"
# use this as the default docker-compose yaml definition
COMPOSE_FILE_BASE=compose-test-net.yaml
# docker-compose.yaml file if you are using couchdb
COMPOSE_FILE_COUCH=compose-couch.yaml
# certificate authorities compose file
COMPOSE_FILE_CA=compose-ca.yaml
#
# chaincode language defaults to "NA"
CC_SRC_LANGUAGE="NA"
# default to running the docker commands for the CCAAS
CCAAS_DOCKER_RUN=true
# Chaincode version
CC_VERSION="1.0"
# Chaincode definition sequence
CC_SEQUENCE=1
# default database
DATABASE="leveldb"

# Get docker sock path from environment variable
SOCK="${DOCKER_HOST:-/var/run/docker.sock}"
DOCKER_SOCK="${SOCK##unix://}"

# Parse commandline args

## Parse mode
if [[ $# -lt 1 ]] ; then
	printHelp
	exit 0
else
	MODE=$1
	shift
fi

# parse a createChannel subcommand if used
if [[ $# -ge 1 ]] ; then
	key="$1"
fi

# parse flags

while [[ $# -ge 1 ]] ; do
	key="$1"
	shift
	case $key in
	-h )
		printHelp $MODE
		exit 0
		;;
	-ca )
		CRYPTO="Certificate Authorities"
		;;
	-s )
		if [[ $# -ge 1 ]] ; then
			DATABASE="$1"
			shift
		else
			fatalln "No database name provided after flag '-s'"
		fi
		;;
	-verbose )
		VERBOSE=true
		;;
	* )
		errorln "Unknown flag: $key"
		printHelp
		exit 1
		;;
	esac
	shift
done

# Are we generating crypto material with this command?
if [ ! -d "organizations/peerOrganizations" ]; then
	CRYPTO_MODE="with crypto from '${CRYPTO}'"
else
	CRYPTO_MODE=""
fi

# Determine mode of operation and printing out what we asked for
if [ "$MODE" == "up" ]; then
	infoln "Starting network"
	up
elif [ "$MODE" == "down" ]; then
	infoln "Stopping network"
	down
elif [ "$MODE" == "restart" ]; then
	infoln "Restarting network"
	down
	up
else
	printHelp
	exit 1
fi
