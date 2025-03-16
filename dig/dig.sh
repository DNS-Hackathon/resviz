#!/bin/bash
dig @1.1.1.1 A $1 +trace +all | jc --dig | jq
