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
package org.apache.plc4x.java.ads.api.generic.types;

import org.apache.plc4x.java.ads.api.util.ByteValue;

import java.nio.ByteBuffer;
import java.util.regex.Pattern;

public class AMSPort extends ByteValue {

    public static final Pattern AMS_PORT_PATTERN = Pattern.compile("\\d+");

    public static final int NUM_BYTES = 2;

    AMSPort(byte... value) {
        super(value);
        assertLength(NUM_BYTES);
    }

    public static AMSPort of(byte... values) {
        return new AMSPort(values);
    }

    public static AMSPort of(int port) {
        checkUnsignedBounds(port, NUM_BYTES);
        return new AMSPort(ByteBuffer.allocate(NUM_BYTES)
            // LE
            .put((byte) (port & 0xff))
            .put((byte) (port >> 8 & 0xff))
            .array());
    }

    public static AMSPort of(String port) {
        if (!AMS_PORT_PATTERN.matcher(port).matches()) {
            throw new IllegalArgumentException(port + " must match " + AMS_PORT_PATTERN);
        }
        return of(Integer.parseInt(port));
    }

    @Override
    public String toString() {
        return "" + (getBytes()[1] << 8 | getBytes()[0]);
    }
}
