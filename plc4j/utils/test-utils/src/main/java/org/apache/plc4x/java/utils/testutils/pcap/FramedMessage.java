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
 * A single protocol message extracted from a pcap capture, with its raw bytes and direction metadata.
 *
 * @param rawBytes   the complete protocol message bytes (ready for parsing)
 * @param isResponse true if this packet was sent from the server (i.e. srcPort matches the protocol port)
 * @param index      sequential index of this message in the capture (used for test naming)
 */
public record FramedMessage(byte[] rawBytes, boolean isResponse, int index) {
}
