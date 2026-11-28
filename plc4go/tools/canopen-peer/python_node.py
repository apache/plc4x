# Licensed to the Apache Software Foundation (ASF) under one or more contributor license
# agreements. See the NOTICE file distributed with this work for additional information regarding
# copyright ownership. The ASF licenses this file to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance with the License. You may obtain a
# copy of the License at https://www.apache.org/licenses/LICENSE-2.0. Unless required by
# applicable law or agreed to in writing, software distributed under the License is distributed on
# an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#
# Test peer only: a CANopen node from the python canopen library (LocalNode: SDO server, NMT
# slave, heartbeat producer) for internal/canopen's peer_linux_test.go.
#
#   python3 python_node.py <interface> <node id>
#
# Its object dictionary:
#
#   0x1000/0  UNSIGNED32      ro  0x00020192
#   0x1001/0  UNSIGNED8       ro  0
#   0x1008/0  VISIBLE_STRING  ro  "plc4x canopen python peer" (segmented)
#   0x1017/0  UNSIGNED16      rw  0, the heartbeat producer time
#   0x1018/1  UNSIGNED32      ro  0x0000AB42, the vendor id
#   0x2000/0  VISIBLE_STRING  rw  ""
#   0x2001/0  INTEGER32       rw  -7
#   0x2002/0  REAL32          rw  1.5
#
# While operational the node sends 0x2001 as its TPDO 1 (0x180 + node id) every 100 ms.
import sys
import threading
import time

import canopen
from canopen import objectdictionary as od


def variable(index, name, data_type, access, default, subindex=0):
    v = od.ODVariable(name, index, subindex)
    v.data_type = data_type
    v.access_type = access
    v.default = default
    return v


def dictionary():
    d = od.ObjectDictionary()
    d.add_object(variable(0x1000, "Device type", od.UNSIGNED32, "ro", 0x00020192))
    d.add_object(variable(0x1001, "Error register", od.UNSIGNED8, "ro", 0))
    d.add_object(variable(0x1008, "Manufacturer device name", od.VISIBLE_STRING, "ro", "plc4x canopen python peer"))
    d.add_object(variable(0x1017, "Producer heartbeat time", od.UNSIGNED16, "rw", 0))
    identity = od.ODRecord("Identity", 0x1018)
    identity.add_member(variable(0x1018, "Highest sub-index", od.UNSIGNED8, "ro", 1, 0))
    identity.add_member(variable(0x1018, "Vendor-ID", od.UNSIGNED32, "ro", 0x0000AB42, 1))
    d.add_object(identity)
    d.add_object(variable(0x2000, "Text", od.VISIBLE_STRING, "rw", ""))
    d.add_object(variable(0x2001, "Number", od.INTEGER32, "rw", -7))
    d.add_object(variable(0x2002, "Real", od.REAL32, "rw", 1.5))
    return d


def main():
    interface, node_id = sys.argv[1], int(sys.argv[2], 0)
    network = canopen.Network()
    network.connect(channel=interface, interface="socketcan")
    node = network.create_node(node_id, dictionary())
    node.nmt.state = "PRE-OPERATIONAL"
    network.send_message(0x700 + node_id, [0x00])  # boot-up
    print("python canopen node", node_id, "on", interface, flush=True)

    def tpdo():
        while True:
            time.sleep(0.1)
            if node.nmt.state == "OPERATIONAL":
                value = node.sdo[0x2001].raw
                network.send_message(0x180 + node_id, value.to_bytes(4, "little", signed=True))

    threading.Thread(target=tpdo, daemon=True).start()
    try:
        while True:
            time.sleep(1)
    finally:
        network.disconnect()


if __name__ == "__main__":
    main()
