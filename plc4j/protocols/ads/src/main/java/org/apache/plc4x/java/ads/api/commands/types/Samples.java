/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.ads.api.commands.types;

import org.apache.plc4x.java.ads.api.util.ByteValue;

import java.nio.ByteBuffer;

public class Samples extends ByteValue {

    public static final int NUM_BYTES = 4;

    Samples(byte... values) {
        super(values);
        assertLength(NUM_BYTES);
    }

    public static Samples of(long numberOfSamples) {
        checkUnsignedBounds(numberOfSamples, NUM_BYTES);
        return new Samples(ByteBuffer.allocate(NUM_BYTES)
            .put((byte) (numberOfSamples >> 24 & 0xff))
            .put((byte) (numberOfSamples >> 16 & 0xff))
            .put((byte) (numberOfSamples >> 8 & 0xff))
            .put((byte) (numberOfSamples & 0xff))
            .array());
    }

    public static Samples of(String numberOfSamples) {
        return of(Long.parseLong(numberOfSamples));
    }

    public static Samples of(byte... values) {
        return new Samples(values);
    }

    @Override
    public String toString() {
        return "" + (getBytes()[0] << 24 | getBytes()[1] << 16 | getBytes()[2] << 8 | getBytes()[3]);
    }
}
