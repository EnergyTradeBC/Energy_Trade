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
AS_ORG="$2"
CHANNEL_NAME="$3"
DELAY="$4"
TIMEOUT="$5"
VERBOSE="$6"
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

infoln "Signing config transaction"
signConfigtxAsPeerOrg ${AS_ORG} org${ORG}_update_in_envelope.pb