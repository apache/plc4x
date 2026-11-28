#!/bin/sh
# Licensed to the Apache Software Foundation (ASF) under the Apache License, Version 2.0.
#
# Runs the CANopen peer tests against CANopenNode and python canopen, inside the canopen-peer
# image:
#
#   docker build -t canopen-peer plc4go/tools/canopen-peer
#   docker run --rm --privileged -v "$PWD":/src -w /src/plc4go canopen-peer sh tools/canopen-peer/run-test.sh
#
# Both nodes sit on one virtual CAN bus (vcan0): CANopenNode as node 0x21, python canopen as node
# 0x22. The bus traffic and each node's log are printed after the test.
set -e
ip link add dev vcan0 type vcan
ip link set up vcan0
mkdir -p /tmp/canopend
# canopend reports a storage error (EMCY 0x5000) and refuses NMT start when its storage files are
# missing; a file holding "-" says "no stored data, use the defaults".
for f in od_comm.persist mainline.persist; do echo - > /tmp/canopend/$f; done
candump -L vcan0 > /tmp/candump.log 2>&1 &
(cd /tmp/canopend && canopend vcan0 -i 0x21 > /tmp/canopend.log 2>&1 &)
/opt/canopen/bin/python3 /src/plc4go/tools/canopen-peer/python_node.py vcan0 0x22 > /tmp/python_node.log 2>&1 &
sleep 2
status=0
CANOPEN_TEST_INTERFACE=vcan0 CANOPEN_TEST_CANOPENNODE=0x21 CANOPEN_TEST_PYTHON=0x22 \
  go test -timeout 180s -count=1 -v -run 'TestSocketCAN|TestPeer' ./internal/canopen/ || status=$?
echo "---- canopend log"
grep -v '^\s*$' /tmp/canopend.log | tail -20
echo "---- python canopen log"
grep -v '^\s*$' /tmp/python_node.log | tail -20
echo "---- bus (last 40 frames)"
tail -40 /tmp/candump.log
exit $status
