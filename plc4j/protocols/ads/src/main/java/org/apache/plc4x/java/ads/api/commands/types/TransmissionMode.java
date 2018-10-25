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

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.UnsignedIntLEByteValue;

@SuppressWarnings("unused") // Due to predefined TransmissionModes
public class TransmissionMode extends UnsignedIntLEByteValue {

    public static final int NUM_BYTES = UnsignedIntLEByteValue.UNSIGNED_INT_LE_NUM_BYTES;

    private TransmissionMode(byte... values) {
        super(values);
    }

    private TransmissionMode(long value) {
        super(value);
    }

    private TransmissionMode(String value) {
        super(value);
    }

    private TransmissionMode(ByteBuf byteBuf) {
        super(byteBuf);
    }

    public static TransmissionMode of(byte... values) {
        return new TransmissionMode(values);
    }

    public static TransmissionMode of(long value) {
        return new TransmissionMode(value);
    }

    public static TransmissionMode of(String value) {
        return new TransmissionMode(value);
    }

    public static TransmissionMode of(ByteBuf byteBuf) {
        return new TransmissionMode(byteBuf);
    }

    // TODO: defined as enum so we might do this too similar to AdsReturnCode
    public static final class DefinedValues {
        public static final TransmissionMode ADSTRANS_NOTRANS = TransmissionMode.of(0);
        public static final TransmissionMode ADSTRANS_CLIENTCYCLE = TransmissionMode.of(1);
        public static final TransmissionMode ADSTRANS_CLIENTONCHA = TransmissionMode.of(2);
        public static final TransmissionMode ADSTRANS_SERVERCYCLE = TransmissionMode.of(3);
        public static final TransmissionMode ADSTRANS_SERVERONCHA = TransmissionMode.of(4);
        public static final TransmissionMode ADSTRANS_SERVERCYCLE2 = TransmissionMode.of(5);
        public static final TransmissionMode ADSTRANS_SERVERONCHA2 = TransmissionMode.of(6);
        public static final TransmissionMode ADSTRANS_CLIENT1REQ = TransmissionMode.of(10);
        public static final TransmissionMode ADSTRANS_MAXMODES = TransmissionMode.of(Integer.MAX_VALUE);

        private DefinedValues() {
            // Container class
        }
    }
}
