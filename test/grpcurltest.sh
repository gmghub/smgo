#!/bin/bash

srv="127.0.0.1:50051"
grpcurlcmd="$(which grpcurl) -plaintext -proto ./v1/SmgoService.proto"

cd ./api

echo
echo "*** Stat tests"
# $grpcurlcmd -rpc-header 'calendar-userid: grpcuser' -d '{"startat": "2022-08-28T00:00:00.000000000Z"}' $srv event.EventService.ListDay
# $grpcurlcmd -rpc-header 'calendar-userid: grpcuser' -d '{"startat": "2022-08-28T00:00:00.000000000Z"}' $srv event.EventService.ListWeek
 $grpcurlcmd -d '{"statinterval": 1, "statperiod":5}' $srv smgo.SmgoService.GetSysStat

#$grpcurlcmd -v -d '{"startat": "2022-08-28T00:00:00.000000000Z"}' $srv event.EventService.ListMonth

exit 0
