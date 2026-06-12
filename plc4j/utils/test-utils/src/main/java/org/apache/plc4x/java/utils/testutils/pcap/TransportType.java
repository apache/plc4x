/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.utils.testutils.pcap;

/**
 * Transport layer type used to determine how protocol messages are extracted from pcap captures.
 *
 * <ul>
 *   <li>{@link #TCP} — messages are extracted via TCP stream reassembly using a {@link FramingSpec}</li>
 *   <li>{@link #UDP} — each UDP datagram payload matching the configured port is one complete message</li>
 *   <li>{@link #ETHERNET} — each Ethernet frame matching the configured EtherType is one message
 *       (including MAC headers, since root types like {@code Ethernet_Frame} parse those)</li>
 * </ul>
 */
public enum TransportType {
    TCP,
    UDP,
    ETHERNET
}
