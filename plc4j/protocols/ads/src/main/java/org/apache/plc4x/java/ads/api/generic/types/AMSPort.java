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
import org.apache.plc4x.java.ads.api.util.ByteValue;

import java.nio.ByteBuffer;
import java.util.regex.Pattern;

/**
 * The ADS devices in a TwinCAT message router are uniquely identified by a number referred to as the ADS-PortNr. For ADS devices this has a fixed specification, whereas pure ADS client applications (e.g. a visualisation system) are allocated a variable ADS port number when they first access the message router.
 *
 * @see <a href="https://infosys.beckhoff.com/content/1033/tcadscommon/html/tcadscommon_identadsdevice.htm?id=3991659524769593444">ADS device identification</a>
 */
public class AMSPort extends ByteValue {

    public static final Pattern AMS_PORT_PATTERN = Pattern.compile("\\d+");

    public static final int NUM_BYTES = 2;

    private AMSPort(byte... value) {
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

    public static AMSPort of(ByteBuf byteBuf) {
        return of(byteBuf.readUnsignedShortLE());
    }

    @Override
    public String toString() {
        return Integer.toString(getBytes()[1] << 8 | getBytes()[0] & 0xff);
    }

    public static class ReservedPorts {

        private ReservedPorts() {
            // Container class
        }

        /**
         * Logger (only NT-Log)
         */
        public static final AMSPort logger = AMSPort.of(100);
        /**
         * Eventlogger
         */
        public static final AMSPort eventLogger = AMSPort.of(110);
        /**
         * IO
         */
        public static final AMSPort io = AMSPort.of(300);
        /**
         * additional Task 1
         */
        public static final AMSPort additionalTask1 = AMSPort.of(301);
        /**
         * additional Task 2
         */
        public static final AMSPort additionalTask2 = AMSPort.of(302);

        /**
         * NC
         */
        public static final AMSPort nc = AMSPort.of(500);
        /**
         * PLC RuntimeSystem 1
         */
        public static final AMSPort plcRuntimeSystem1 = AMSPort.of(801);
        /**
         * PLC RuntimeSystem 2
         */
        public static final AMSPort plcRuntimeSystem2 = AMSPort.of(811);
        /**
         * PLC RuntimeSystem 3
         */
        public static final AMSPort plcRuntimeSystem3 = AMSPort.of(821);
        /**
         * PLC RuntimeSystem 4
         */
        public static final AMSPort plcRuntimeSystem4 = AMSPort.of(831);

        /**
         * Camshaft controller
         */
        public static final AMSPort camShaftController = AMSPort.of(900);

        /**
         * System Service
         */
        public static final AMSPort systemService = AMSPort.of(900);
        /**
         * Scope
         */
        public static final AMSPort scope = AMSPort.of(900);
    }
}
