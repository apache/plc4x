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

/**
 * Fluent builder for constructing {@link CanFrame} instances with full validation.
 * <p>
 * Validates all CAN 2.0 constraints on {@link #build()}:
 * <ul>
 *   <li>Standard identifiers must be in range 0–0x7FF (11 bits)</li>
 *   <li>Extended identifiers must be in range 0–0x1FFFFFFF (29 bits)</li>
 *   <li>Data payload must not exceed 8 bytes</li>
 *   <li>RTR frames should have no data payload</li>
 * </ul>
 * <p>
 * Example usage:
 * <pre>{@code
 * CanFrame frame = CanFrame.builder()
 *     .identifier(0x123)
 *     .data(new byte[]{0x01, 0x02, 0x03})
 *     .build();
 * }</pre>
 */
public class CanFrameBuilder {

    private int identifier;
    private byte[] data = new byte[0];
    private boolean extended = false;
    private boolean rtr = false;
    private boolean identifierSet = false;

    /**
     * Creates a new builder with default values.
     */
    CanFrameBuilder() {
        // Package-private constructor — use CanFrame.builder()
    }

    /**
     * Sets the CAN arbitration identifier.
     *
     * @param identifier the CAN ID (validated on build)
     * @return this builder for chaining
     */
    public CanFrameBuilder identifier(int identifier) {
        this.identifier = identifier;
        this.identifierSet = true;
        return this;
    }

    /**
     * Sets the frame data payload.
     *
     * @param data the payload bytes (0–8 bytes, validated on build)
     * @return this builder for chaining
     * @throws IllegalArgumentException if data is null
     */
    public CanFrameBuilder data(byte[] data) {
        if (data == null) {
            throw new IllegalArgumentException("Data must not be null");
        }
        this.data = data;
        return this;
    }

    /**
     * Sets whether this frame uses an extended (29-bit) identifier.
     *
     * @param extended true for extended, false for standard (11-bit)
     * @return this builder for chaining
     */
    public CanFrameBuilder extended(boolean extended) {
        this.extended = extended;
        return this;
    }

    /**
     * Sets the Remote Transmission Request flag.
     *
     * @param rtr true if this is an RTR frame
     * @return this builder for chaining
     */
    public CanFrameBuilder rtr(boolean rtr) {
        this.rtr = rtr;
        return this;
    }

    /**
     * Builds the {@link CanFrame}, validating all CAN 2.0 constraints.
     *
     * @return the constructed CAN frame
     * @throws IllegalArgumentException if any constraint is violated
     */
    public CanFrame build() {
        if (!identifierSet) {
            throw new IllegalArgumentException("Identifier must be set");
        }

        // Validate identifier range based on type
        if (identifier < 0) {
            throw new IllegalArgumentException(
                    "CAN identifier must not be negative, got: " + identifier);
        }
        if (extended) {
            if (identifier > CanFrame.MAX_EXTENDED_ID) {
                throw new IllegalArgumentException(String.format(
                        "Extended CAN identifier must be in range 0–0x%X, got: 0x%X",
                        CanFrame.MAX_EXTENDED_ID, identifier));
            }
        } else {
            if (identifier > CanFrame.MAX_STANDARD_ID) {
                throw new IllegalArgumentException(String.format(
                        "Standard CAN identifier must be in range 0–0x%X, got: 0x%X",
                        CanFrame.MAX_STANDARD_ID, identifier));
            }
        }

        // Validate data length
        if (data.length > CanFrame.MAX_DATA_LENGTH) {
            throw new IllegalArgumentException(String.format(
                    "CAN frame data must not exceed %d bytes, got: %d",
                    CanFrame.MAX_DATA_LENGTH, data.length));
        }

        // RTR frames should have empty data
        if (rtr && data.length > 0) {
            throw new IllegalArgumentException(
                    "RTR frames must not contain data payload");
        }

        return new CanFrame(identifier, data, extended, rtr);
    }
}
