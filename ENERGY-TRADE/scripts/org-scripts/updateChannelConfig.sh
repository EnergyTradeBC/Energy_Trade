#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# This script is designed to be run in the cli container as the
# first step of the EYFN tutorial.  It creates and submits a
# configuration transaction to add the org to the test network
#

ORG="$1"
CHANNEL_NAME="$2"
DELAY="$3"
TIMEOUT="$4"
VERBOSE="$5"
: ${CHANNEL_NAME:="cer"}
: ${DELAY:="3"}
: ${TIMEOUT:="10"}
: ${VERBOSE:="false"}
COUNTER=1
MAX_RETRY=5


# imports
. scripts/envVar.sh
. scripts/configUpdate.sh
. scripts/utils.sh
. scripts/actAsOrg.sh

infoln "Creating config transaction to add org${ORG} to network"

# Fetch the config for the channel, writing it to config.json
fetchChannelConfig 1 ${CHANNEL_NAME} config.json

# Modify the configuration to append the new org
set -x
jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"Org'${ORG}'MSP":.[1]}}}}}' config.json ./organizations/peerOrganizations/org${ORG}.example.com/org${ORG}.json > modified_config.json
{ set +x; } 2>/dev/null

# Compute a config update, based on the differences between config.json and modified_config.json, write it as a transaction to org${ORG}_update_in_envelope.pb
createConfigUpdate ${CHANNEL_NAME} config.json modified_config.json org${ORG}_update_in_envelope.pb