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
 * Describes how to detect individual protocol message boundaries within a TCP byte stream.
 *
 * <p>The framing logic reads a length field at a fixed offset within each message header,
 * then computes the total message size as {@code lengthFieldValue + lengthAdjustment}.
 *
 * @param lengthFieldOffset byte offset from the start of the message to the length field
 * @param lengthFieldSize   size of the length field in bytes (typically 2 or 4)
 * @param bigEndian         whether the length field is big-endian encoded
 * @param lengthAdjustment  added to the length field value to get total message size
 */
public record FramingSpec(int lengthFieldOffset, int lengthFieldSize, boolean bigEndian, int lengthAdjustment) {

    /**
     * Validates that the framing spec parameters are sensible.
     */
    public FramingSpec {
        if (lengthFieldOffset < 0) {
            throw new IllegalArgumentException("lengthFieldOffset must be >= 0, got " + lengthFieldOffset);
        }
        if (lengthFieldSize != 1 && lengthFieldSize != 2 && lengthFieldSize != 4) {
            throw new IllegalArgumentException("lengthFieldSize must be 1, 2, or 4, got " + lengthFieldSize);
        }
    }
}
