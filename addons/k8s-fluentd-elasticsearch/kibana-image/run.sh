#!/bin/bash
set -e

# Add kibana as command if needed
if [[ "$1" == -* ]]; then
	set -- kibana "$@"
fi

export ELASTICSEARCH_URL=${ELASTICSEARCH_URL:-"http://localhost:9200"}
echo ELASTICSEARCH_URL=${ELASTICSEARCH_URL}

#export KIBANA_BASE_URL=${KIBANA_BASE_URL:-"''"}
#echo "server.basePath: ${KIBANA_BASE_URL}"
#echo "server.basePath: ${KIBANA_BASE_URL}" >> /etc/kibana/kibana.yml

kibana -e ${ELASTICSEARCH_URL}
