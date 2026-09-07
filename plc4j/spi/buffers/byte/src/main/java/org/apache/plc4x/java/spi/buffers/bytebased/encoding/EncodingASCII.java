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
package org.apache.plc4x.java.spi.buffers.bytebased.encoding;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;

/**
 * In ASCII encoding, numbers are simply represented as their ASCII-encoded string value
 * Each value must have a bit-length that must be a multiple of 8.
 *
 * <p>Strings in ASCII fields are right-padded with spaces ({@code 0x20}) — the
 * Open-Protocol spec mandates this and the round-trip decode (in
 * {@link BaseStringEncoding}) already trims trailing spaces, so this matches.
 * The {@link BaseStringEncoding#encodeString(int, String, boolean)} default
 * leaves a zero-padded byte array, which other ASCII users could quietly
 * tolerate but which Open-Protocol explicitly forbids.</p>
 */
public class EncodingASCII extends BaseStringEncoding {

    public static final String NAME = "ASCII";

    private static final WithOption OPTION = WithOption.WithEncoding(NAME);

    public static WithOption optionEncodingASCII() {
        return OPTION;
    }

    @Override
    public String getName() {
        return NAME;
    }

    @Override
    protected int getBitsPerCharacter() {
        return 8;
    }

    @Override
    protected Charset getCharset() {
        return StandardCharsets.US_ASCII;
    }

    /**
     * String values get right-padded with spaces ({@code 0x20}) rather than
     * the default zero-pad. Numeric values still get the parent class's
     * '0'-left-pad treatment.
     */
    @Override
    public byte[] encodeString(int numBits, String value) throws BufferException {
        byte[] result = super.encodeString(numBits, value);
        int numBytes = numBits / 8;
        byte[] valueBytes = value == null ? new byte[0] : value.getBytes(getCharset());
        // The parent left the trailing slots (after the actual string bytes)
        // zero-initialized; replace those with ASCII spaces. We only touch
        // bytes the parent didn't populate — anything inside the string range
        // is left untouched.
        for (int i = valueBytes.length; i < numBytes; i++) {
            if (result[i] == 0) {
                result[i] = ' ';
            }
        }
        return result;
    }

}
