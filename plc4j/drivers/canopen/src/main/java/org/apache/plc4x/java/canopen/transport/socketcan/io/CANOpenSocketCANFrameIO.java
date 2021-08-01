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
package org.apache.plc4x.java.canopen.transport.socketcan.io;

import org.apache.plc4x.java.canopen.transport.CANOpenFrame;
import org.apache.plc4x.java.canopen.transport.socketcan.CANOpenSocketCANFrame;
import org.apache.plc4x.java.canopen.helper.HeaderParser;
import org.apache.plc4x.java.canopen.readwrite.CANOpenPayload;
import org.apache.plc4x.java.canopen.readwrite.io.CANOpenPayloadIO;
import org.apache.plc4x.java.canopen.readwrite.types.CANOpenService;
import org.apache.plc4x.java.spi.generation.MessageIO;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class CANOpenSocketCANFrameIO implements MessageIO<CANOpenFrame, CANOpenFrame> {

    private static final Logger LOGGER = LoggerFactory.getLogger(CANOpenSocketCANFrameIO.class);

    @Override
    public CANOpenSocketCANFrame parse(ReadBuffer io, Object... args) throws ParseException {
        return CANOpenSocketCANFrameIO.staticParse(io);
    }

    @Override
    public void serialize(WriteBuffer io, CANOpenFrame value, Object... args) throws ParseException {
        CANOpenSocketCANFrameIO.staticSerialize(io, (CANOpenSocketCANFrame) value);
    }

    public static CANOpenSocketCANFrame staticParse(ReadBuffer io) throws ParseException {
        int startPos = io.getPos();
        int curPos;

        // Simple Field (rawId)
        int rawId = io.readInt(32);

        // Virtual field (Just declare a local variable so we can access it in the parser)
        int identifier = (int) (HeaderParser.readIdentifier(rawId));

        CANOpenService service = serviceId(identifier);
        int nodeId = Math.abs(service.getMin() - identifier);

        // Implicit Field (size) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
        short size = io.readUnsignedShort(8);

        // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
        {
            short reserved = io.readUnsignedShort(8);
            if(reserved != (short) 0x0) {
                LOGGER.info("Expected constant value " + 0x0 + " but got " + reserved + " for reserved field.");
            }
        }

        // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
        {
            short reserved = io.readUnsignedShort(8);
            if(reserved != (short) 0x0) {
                LOGGER.info("Expected constant value " + 0x0 + " but got " + reserved + " for reserved field.");
            }
        }

        // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
        {
            short reserved = io.readUnsignedShort(8);
            if(reserved != (short) 0x0) {
                LOGGER.info("Expected constant value " + 0x0 + " but got " + reserved + " for reserved field.");
            }
        }

        // Array field (data)
        // Count array
        if(size > Integer.MAX_VALUE) {
            throw new ParseException("Array count of " + (size) + " exceeds the maximum allowed count of " + Integer.MAX_VALUE);
        }

        final CANOpenPayload payload = CANOpenPayloadIO.staticParse(io, service);

        // Padding Field (padding)
        {
            int _timesPadding = (int) ((8) - payload.getLengthInBytes());
            while ((io.hasMore(8)) && (_timesPadding-- > 0)) {
                // Just read the padding data and ignore it
                io.readUnsignedShort(8);
            }
        }

        // Create the instance
        return new CANOpenSocketCANFrame(nodeId, service, payload);
    }

    public static void staticSerialize(WriteBuffer io, CANOpenSocketCANFrame _value) throws ParseException {
        int startPos = io.getPos();

        // Simple Field (service)
        int nodeId = _value.getNodeId();
        int service = _value.getService().getMin();
        io.writeInt(32, (service + nodeId));

        // Implicit Field (size) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
        final CANOpenPayload payload = _value.getPayload();
        short size = (short) (payload == null ? 0 : payload.getLengthInBytes());
        io.writeUnsignedShort(8, ((Number) (size)).shortValue());

        // Reserved Field (reserved)
        io.writeUnsignedShort(8, ((Number) (short) 0x0).shortValue());

        // Reserved Field (reserved)
        io.writeUnsignedShort(8, ((Number) (short) 0x0).shortValue());

        // Reserved Field (reserved)
        io.writeUnsignedShort(8, ((Number) (short) 0x0).shortValue());

        // Array Field (data)
        if(_value.getPayload() != null) {
            payload.getMessageIO().serialize(io, payload);
        }

        // Padding Field (padding)
        {
            int _timesPadding = (int) ((8) - size);
            while (_timesPadding-- > 0) {
                short _paddingValue = (short) (0);
                io.writeUnsignedShort(8, ((Number) (_paddingValue)).shortValue());
            }
        }
    }

    public static CANOpenService serviceId(int cobId) {
        // form 32 bit socketcan identifier
        CANOpenService service = CANOpenService.enumForValue((byte) (cobId >> 7));
        if (service == null) {
            for (CANOpenService val : CANOpenService.values()) {
                if (val.getMin() > cobId && val.getMax() < cobId) {
                    return val;
                }
            }
        }
        return service;
    }

}
