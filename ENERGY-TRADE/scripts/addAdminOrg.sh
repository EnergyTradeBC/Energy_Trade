. scripts/envVar.sh
. scripts/utils.sh
. scripts/addOrgHelp.sh
. scripts/addOrgUtils.sh

if [[ $# -lt 1 ]] ; then
  	addOrgHelp
  	exit 0
else
  	MODE=$1
  	shift
fi

if [[ $# -lt 1 ]] ; then
  	addOrgHelp $MODE
  	exit 0
else
  	ORG=$1
  	shift
fi

if [[ $# -lt 1 ]] ; then
  	fatalln "Two ports required"
else
  	PORT_1=$1
  	shift
fi

if [[ $# -lt 1 ]] ; then
  	fatalln "Two ports required"
else
  	PORT_2=$1
  	shift
fi

: ${CONTAINER_CLI:="docker"}
: ${CONTAINER_CLI_COMPOSE:="${CONTAINER_CLI}-compose"}
infoln "Using ${CONTAINER_CLI} and ${CONTAINER_CLI_COMPOSE}"

# Using crpto vs CA. default is cryptogen
CRYPTO="cryptogen"
# timeout duration - the duration the CLI should wait for a response from
# another container before giving up
CLI_TIMEOUT=10
#default for delay
CLI_DELAY=3
# channel name defaults to "mychannel"
CHANNEL_NAME="mychannel"
# use this as the docker compose couch file
COMPOSE_FILE_COUCH_BASE=org${ORG}/compose/compose-couch-org${ORG}.yaml
COMPOSE_FILE_COUCH_ORG=org${ORG}/compose/${CONTAINER_CLI}/docker-compose-couch-org${ORG}.yaml
# use this as the default docker-compose yaml definition
COMPOSE_FILE_BASE=org${ORG}/compose/compose-org${ORG}.yaml
COMPOSE_FILE_ORG=org${ORG}/compose/${CONTAINER_CLI}/docker-compose-org${ORG}.yaml
# certificate authorities compose file
COMPOSE_FILE_CA_BASE=org${ORG}/compose/compose-ca-org${ORG}.yaml
COMPOSE_FILE_CA_ORG=org${ORG}/compose/${CONTAINER_CLI}/docker-compose-ca-org${ORG}.yaml
# database
DATABASE="leveldb"

# Get docker sock path from environment variable
SOCK="${DOCKER_HOST:-/var/run/docker.sock}"
DOCKER_SOCK="${SOCK##unix://}"

if [ "$MODE" == "files" ]; then
	infoln "Creating org${ORG} files"
  	addOrgFiles
elif [ "$MODE" == "add" ]; then
	infoln "Adding org${ORG} to the network"
  	addOrg
else
  	addOrgHelp
  	exit 1
fi