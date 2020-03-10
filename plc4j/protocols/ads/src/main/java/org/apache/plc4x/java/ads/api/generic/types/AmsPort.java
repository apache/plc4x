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

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.UnsignedShortLEByteValue;

import java.util.regex.Pattern;

/**
 * The ADS devices in a TwinCAT message router are uniquely identified by a number referred to as the ADS-PortNr. For ADS devices this has a fixed specification, whereas pure ADS client applications (e.g. a visualisation system) are allocated a variable ADS port number when they first access the message router.
 *
 * @see <a href="https://infosys.beckhoff.com/content/1033/tcadscommon/html/tcadscommon_identadsdevice.htm?id=3991659524769593444">ADS device identification</a>
 */
@SuppressWarnings("unused") // Due to predefined ports
public class AmsPort extends UnsignedShortLEByteValue {

    public static final Pattern AMS_PORT_PATTERN = Pattern.compile("\\d+");

    public static final int NUM_BYTES = UnsignedShortLEByteValue.UNSIGNED_SHORT_LE_NUM_BYTES;

    private AmsPort(byte... value) {
        super(value);
    }

    private AmsPort(int value) {
        super(value);
    }

    public static AmsPort of(byte... values) {
        return new AmsPort(values);
    }

    public static AmsPort of(int port) {
        return new AmsPort(port);
    }

    public static AmsPort of(String port) {
        if (!AMS_PORT_PATTERN.matcher(port).matches()) {
            throw new IllegalArgumentException(port + " must match " + AMS_PORT_PATTERN);
        }
        return of(Integer.parseInt(port));
    }

    public static AmsPort of(ByteBuf byteBuf) {
        return of(byteBuf.readUnsignedShortLE());
    }

    @Override
    public String toString() {
        return String.valueOf(getAsInt());
    }

    public static class ReservedPorts {

        /**
         * Logger (only NT-Log)
         */
        public static final AmsPort logger = AmsPort.of(100);
        /**
         * Eventlogger
         */
        public static final AmsPort eventLogger = AmsPort.of(110);
        /**
         * IO
         */
        public static final AmsPort io = AmsPort.of(300);
        /**
         * additional Task 1
         */
        public static final AmsPort additionalTask1 = AmsPort.of(301);
        /**
         * additional Task 2
         */
        public static final AmsPort additionalTask2 = AmsPort.of(302);

        /**
         * NC
         */
        public static final AmsPort nc = AmsPort.of(500);
        /**
         * PLC RuntimeSystem 1
         */
        public static final AmsPort plcRuntimeSystem1 = AmsPort.of(801);
        /**
         * PLC RuntimeSystem 2
         */
        public static final AmsPort plcRuntimeSystem2 = AmsPort.of(811);
        /**
         * PLC RuntimeSystem 3
         */
        public static final AmsPort plcRuntimeSystem3 = AmsPort.of(821);
        /**
         * PLC RuntimeSystem 4
         */
        public static final AmsPort plcRuntimeSystem4 = AmsPort.of(831);

        /**
         * Camshaft controller
         */
        public static final AmsPort camShaftController = AmsPort.of(900);

        /**
         * System Service
         */
        public static final AmsPort systemService = AmsPort.of(900);
        /**
         * Scope
         */
        public static final AmsPort scope = AmsPort.of(900);

        private ReservedPorts() {
            // Container class
        }
    }
}
