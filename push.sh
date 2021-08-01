#!/bin/bash

mosquitto_pub -i satelite -t afm/v1/settings/123 -m "{ \"model\":\"x\",\"firmware\":\"1.0.0\"}" -h localhost -p 1883 -q 0
# --cafile ssl/psikick-ca.pem --cert ssl/node_cert.pem --key ssl/node_pk.pem --tls-version tlsv1.2 --insecure