/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.ads.api.commands.types;

import io.netty.buffer.ByteBuf;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.ads.api.util.ByteValue;
import org.apache.plc4x.java.api.exceptions.PlcProtocolPayloadTooBigException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import java.nio.charset.Charset;
import java.util.Arrays;

import static java.util.Objects.requireNonNull;

public class Device extends ByteValue {

    public static final int NUM_BYTES = 16;

    private Device(byte... values) {
        super(values);
        assertLength(NUM_BYTES);
    }

    public static Device of(byte... values) {
        return new Device(values);
    }

    public static Device of(ByteBuf byteBuf) {
        byte[] values = new byte[NUM_BYTES];
        byteBuf.readBytes(values);
        return of(values);
    }

    public static Device of(String value) {
        return of(value, Charset.defaultCharset());
    }

    public static Device of(String value, Charset charset) {
        requireNonNull(value);
        requireNonNull(charset);
        byte[] bytes = value.getBytes(charset);
        if (bytes.length > NUM_BYTES) {
            throw new PlcRuntimeException(new PlcProtocolPayloadTooBigException("ADS/AMS", NUM_BYTES, bytes.length, value));
        }
        return new Device(Arrays.copyOf(bytes, NUM_BYTES));
    }

    @Override
    public long getCalculatedLength() {
        return NUM_BYTES;
    }

    public String getAsString() {
        return getAsString(Charset.defaultCharset());
    }

    public String getAsString(Charset charset) {
        int nullTermination = ArrayUtils.indexOf(value, (byte) 0);
        byte[] reducedValue = Arrays.copyOfRange(value, 0, nullTermination);
        return new String(reducedValue, charset);
    }

    @Override
    public String toString() {
        // TODO: this might break some outputs like surefire if this id can contain non printable characters
        return "Device{" + getAsString() + "} " + super.toString();
    }
}
