#!/bin/sh

url=http://$GATEWAY_URL/metrics
echo "Call URL: $url"

for i in `seq 1 500000`
do 
   echo "call $i"
   curl -s $url > /dev/null
done