#!/usr/bin/env python3
"""
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.

Generate a test PCAP file with sample Ethernet frames.
This script creates a minimal PCAP file for testing the pcap-replay transport.
"""

import struct
import time

def write_pcap_header(f):
    """Write PCAP global header."""
    magic = 0xa1b2c3d4  # Magic number (microsecond precision)
    version_major = 2
    version_minor = 4
    thiszone = 0  # GMT to local correction
    sigfigs = 0   # Accuracy of timestamps
    snaplen = 65535  # Max length of captured packets
    network = 1   # Data link type (1 = Ethernet)

    f.write(struct.pack('<IHHIIII', magic, version_major, version_minor,
                        thiszone, sigfigs, snaplen, network))

def write_packet(f, data, timestamp=None):
    """Write a packet to the PCAP file."""
    if timestamp is None:
        timestamp = time.time()

    ts_sec = int(timestamp)
    ts_usec = int((timestamp - ts_sec) * 1000000)
    incl_len = len(data)
    orig_len = len(data)

    # Write packet header
    f.write(struct.pack('<IIII', ts_sec, ts_usec, incl_len, orig_len))
    # Write packet data
    f.write(data)

def create_ethernet_frame(dst_mac, src_mac, ethertype, payload):
    """Create an Ethernet frame."""
    frame = bytearray()

    # Destination MAC (6 bytes)
    frame.extend(bytes.fromhex(dst_mac.replace(':', '')))

    # Source MAC (6 bytes)
    frame.extend(bytes.fromhex(src_mac.replace(':', '')))

    # EtherType (2 bytes)
    frame.extend(struct.pack('>H', ethertype))

    # Payload
    frame.extend(payload)

    # Padding to minimum frame size (60 bytes without FCS)
    while len(frame) < 60:
        frame.append(0)

    return bytes(frame)

def main():
    """Generate test.pcap file."""

    # MAC addresses matching the test
    local_mac = '00:11:22:33:44:55'
    remote_mac = 'AA:BB:CC:DD:EE:FF'

    # Custom EtherType for testing (using a test protocol ID)
    ethertype = 0x88B5  # Example custom protocol

    # Create test payloads
    payloads = [
        b'\x01\x02\x03\x04',  # Simple test data
        b'\x05\x06\x07\x08\x09\x0A',  # Another test packet
        b'Hello World!',  # Text payload
        bytes(range(32)),  # Sequential bytes
    ]

    output_file = 'test.pcap'

    with open(output_file, 'wb') as f:
        # Write PCAP header
        write_pcap_header(f)

        base_time = time.time()

        # Write packets with different directions
        for i, payload in enumerate(payloads):
            timestamp = base_time + (i * 0.01)  # 10ms apart

            # Alternate between incoming and outgoing
            if i % 2 == 0:
                # Incoming: from remote to local
                frame = create_ethernet_frame(local_mac, remote_mac, ethertype, payload)
            else:
                # Outgoing: from local to remote
                frame = create_ethernet_frame(remote_mac, local_mac, ethertype, payload)

            write_packet(f, frame, timestamp)

    print(f'Created {output_file} with {len(payloads)} packets')
    print(f'Local MAC: {local_mac}')
    print(f'Remote MAC: {remote_mac}')
    print(f'EtherType: 0x{ethertype:04X}')

if __name__ == '__main__':
    main()