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
import org.apache.plc4x.java.ads.api.util.UnsignedShortLEByteValue;

@SuppressWarnings("unused") // Due to predefined AdsStates
public class AdsState extends UnsignedShortLEByteValue {

    public static final int NUM_BYTES = UnsignedShortLEByteValue.UNSIGNED_SHORT_LE_NUM_BYTES;

    private AdsState(byte... values) {
        super(values);
    }

    private AdsState(int value) {
        super(value);
    }

    private AdsState(String value) {
        super(value);
    }

    private AdsState(ByteBuf byteBuf) {
        super(byteBuf);
    }

    public static AdsState of(byte... values) {
        return new AdsState(values);
    }

    public static AdsState of(int value) {
        return new AdsState(value);
    }

    public static AdsState of(ByteBuf byteBuf) {
        return new AdsState(byteBuf);
    }

    public static AdsState of(String value) {
        return new AdsState(value);
    }

    // TODO: defined as enum so we might do this too similar to AdsReturnCode
    public static final class DefinedValues {
        public static final AdsState ADSSTATE_INVALID = AdsState.of(0);
        public static final AdsState ADSSTATE_IDLE = AdsState.of(1);
        public static final AdsState ADSSTATE_RESET = AdsState.of(2);
        public static final AdsState ADSSTATE_INIT = AdsState.of(3);
        public static final AdsState ADSSTATE_START = AdsState.of(4);
        public static final AdsState ADSSTATE_RUN = AdsState.of(5);
        public static final AdsState ADSSTATE_STOP = AdsState.of(6);
        public static final AdsState ADSSTATE_SAVECFG = AdsState.of(7);
        public static final AdsState ADSSTATE_LOADCFG = AdsState.of(8);
        public static final AdsState ADSSTATE_POWERFAILURE = AdsState.of(9);
        public static final AdsState ADSSTATE_POWERGOOD = AdsState.of(10);
        public static final AdsState ADSSTATE_ERROR = AdsState.of(11);
        public static final AdsState ADSSTATE_SHUTDOWN = AdsState.of(12);
        public static final AdsState ADSSTATE_SUSPEND = AdsState.of(13);
        public static final AdsState ADSSTATE_RESUME = AdsState.of(14);
        public static final AdsState ADSSTATE_CONFIG = AdsState.of(15);
        public static final AdsState ADSSTATE_RECONFIG = AdsState.of(16);
        public static final AdsState ADSSTATE_STOPPING = AdsState.of(17);
        public static final AdsState ADSSTATE_INCOMPATIBLE = AdsState.of(18);
        public static final AdsState ADSSTATE_EXCEPTION = AdsState.of(19);
        public static final AdsState ADSSTATE_MAXSTATES = AdsState.of(65535);
    }
}
