#!/bin/bash

set -e
set +x

curl -f -X GET https://username:password@subway-broker.cfapps.io/v2/catalog

id=$(uuid)
curl -f -X PUT https://username:password@subway-broker.cfapps.io/v2/service_instances/$id -d '{"plan_id": "1", "service_id":"1"}'
curl -f -X PUT https://username:password@subway-broker.cfapps.io/v2/service_instances/$id/service_bindings/$id -d '{"plan_id": "1", "service_id":"1", "app_guid": "x"}'
curl -f -X DELETE https://username:password@subway-broker.cfapps.io/v2/service_instances/$id/service_bindings/$id
curl -f -X DELETE https://username:password@subway-broker.cfapps.io/v2/service_instances/$id

id=$(uuid)
curl -f -X PUT https://username:password@subway-broker.cfapps.io/v2/service_instances/$id -d '{"plan_id": "1", "service_id":"1"}'
curl -f -X PUT https://username:password@subway-broker.cfapps.io/v2/service_instances/$id/service_bindings/$id -d '{"plan_id": "1", "service_id":"1", "app_guid": "x"}'

id=$(uuid)
curl -f -X PUT https://username:password@subway-broker.cfapps.io/v2/service_instances/$id -d '{"plan_id": "1", "service_id":"1"}'
curl -f -X PUT https://username:password@subway-broker.cfapps.io/v2/service_instances/$id/service_bindings/$id -d '{"plan_id": "1", "service_id":"1", "app_guid": "x"}'
