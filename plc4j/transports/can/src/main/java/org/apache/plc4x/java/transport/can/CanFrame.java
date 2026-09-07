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
package org.apache.plc4x.java.transport.can;

import java.util.Arrays;
import java.util.Objects;

/**
 * Immutable value object representing a CAN (Controller Area Network) bus frame.
 * <p>
 * A CAN frame is the fundamental unit of communication on a CAN bus. It consists of:
 * <ul>
 *   <li>An arbitration identifier (11-bit standard or 29-bit extended)</li>
 *   <li>A data payload of 0–8 bytes (CAN 2.0 standard)</li>
 *   <li>Control flags: extended identifier flag and Remote Transmission Request (RTR) flag</li>
 * </ul>
 * <p>
 * Use {@link CanFrameBuilder} to construct instances with full validation.
 *
 * @see CanFrameBuilder
 */
public class CanFrame {

    /** Maximum CAN identifier value for standard (11-bit) frames. */
    public static final int MAX_STANDARD_ID = 0x7FF;

    /** Maximum CAN identifier value for extended (29-bit) frames. */
    public static final int MAX_EXTENDED_ID = 0x1FFFFFFF;

    /** Maximum data payload length in bytes (CAN 2.0). */
    public static final int MAX_DATA_LENGTH = 8;

    private final int identifier;
    private final byte[] data;
    private final boolean extended;
    private final boolean rtr;

    /**
     * Constructs a CAN frame. Use {@link #builder()} for validated construction.
     *
     * @param identifier the CAN arbitration identifier
     * @param data       the frame payload (defensively copied)
     * @param extended   true for 29-bit extended identifier, false for 11-bit standard
     * @param rtr        true for Remote Transmission Request frame
     */
    CanFrame(int identifier, byte[] data, boolean extended, boolean rtr) {
        this.identifier = identifier;
        this.data = data != null ? Arrays.copyOf(data, data.length) : new byte[0];
        this.extended = extended;
        this.rtr = rtr;
    }

    /**
     * Creates a new builder for constructing CAN frames.
     *
     * @return a new {@link CanFrameBuilder} instance
     */
    public static CanFrameBuilder builder() {
        return new CanFrameBuilder();
    }

    /**
     * Returns the CAN arbitration identifier.
     *
     * @return the identifier (11-bit if standard, 29-bit if extended)
     */
    public int getIdentifier() {
        return identifier;
    }

    /**
     * Returns a defensive copy of the frame's data payload.
     *
     * @return a copy of the data bytes
     */
    public byte[] getData() {
        return Arrays.copyOf(data, data.length);
    }

    /**
     * Returns the length of the data payload in bytes.
     *
     * @return data length (0–8)
     */
    public int getDataLength() {
        return data.length;
    }

    /**
     * Returns whether this frame uses an extended (29-bit) identifier.
     *
     * @return true if extended, false if standard (11-bit)
     */
    public boolean isExtended() {
        return extended;
    }

    /**
     * Returns whether this is a Remote Transmission Request frame.
     *
     * @return true if RTR
     */
    public boolean isRtr() {
        return rtr;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        CanFrame canFrame = (CanFrame) o;
        return identifier == canFrame.identifier
                && extended == canFrame.extended
                && rtr == canFrame.rtr
                && Arrays.equals(data, canFrame.data);
    }

    @Override
    public int hashCode() {
        int result = Objects.hash(identifier, extended, rtr);
        result = 31 * result + Arrays.hashCode(data);
        return result;
    }

    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder("CanFrame{");
        sb.append("id=0x").append(String.format(extended ? "%08X" : "%03X", identifier));
        if (extended) sb.append(" (EXT)");
        if (rtr) sb.append(" (RTR)");
        sb.append(", data=[");
        for (int i = 0; i < data.length; i++) {
            if (i > 0) sb.append(' ');
            sb.append(String.format("%02X", data[i] & 0xFF));
        }
        sb.append("], len=").append(data.length);
        sb.append('}');
        return sb.toString();
    }
}
