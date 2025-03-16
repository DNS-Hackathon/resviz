#!/bin/bash
dig @1.1.1.1 SOA $1 +trace +all | jc --dig | jq
