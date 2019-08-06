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
package org.apache.plc4x.protocol.df1;

import org.apache.plc4x.java.df1.DF1ReadRequest;
import org.apache.plc4x.java.df1.DF1Symbol;
import org.apache.plc4x.java.df1.DF1SymbolMessageFrameStart;
import org.apache.plc4x.java.utils.ParseException;
import org.apache.plc4x.java.utils.WriteBuffer;

import java.nio.ByteBuffer;

public class DF1Utils {

    public static short CRCCheck(Object... args) {
        DF1Symbol symbol = (DF1Symbol) args[1];
        short messageType = (short) args[0];
        if(symbol instanceof DF1SymbolMessageFrameStart) {
            DF1SymbolMessageFrameStart messageFrameStart = (DF1SymbolMessageFrameStart) symbol;

            short destinationAddress = messageFrameStart.getDestinationAddress();
            short sourceAddress = messageFrameStart.getSourceAddress();
            short commandDiscriminatorValues = (short) messageFrameStart.getCommand().getDiscriminatorValues()[0];
            short status = messageFrameStart.getCommand().getStatus();
            int   counter = messageFrameStart.getCommand().getTransactionCounter();
            if(messageFrameStart.getCommand() instanceof DF1ReadRequest) {
                DF1ReadRequest readRequestCommand = (DF1ReadRequest) messageFrameStart.getCommand();

                try {
                    WriteBuffer writeBuffer = new WriteBuffer(10, false);
                    writeBuffer.writeUnsignedShort(8, destinationAddress);
                    writeBuffer.writeUnsignedShort(8, sourceAddress);
                    writeBuffer.writeUnsignedShort(8, commandDiscriminatorValues);
                    writeBuffer.writeUnsignedShort(8, status);
                    writeBuffer.writeUnsignedInt(16, (short) counter);
                    writeBuffer.writeUnsignedInt(16, (short) readRequestCommand.getAddress());
                    writeBuffer.writeUnsignedShort(8, (byte) readRequestCommand.getSize());
                    writeBuffer.writeUnsignedShort(8, (byte) messageType);

                    byte[] data = writeBuffer.getData();

                    int tmp = 0;
                    int crcL, crcR;

                    for (int newByte : data) {
                        crcL = tmp >> 8;
                        crcR = tmp & 0xFF;
                        tmp = (crcL << 8) + (newByte ^ crcR);
                        for (int j=0; j<8; j++)
                            if (tmp % 2 == 1) {     // check if LSB shifted out is 1 or 0
                                tmp = tmp >> 1;
                                tmp = tmp ^ 0xA001;
                            } else {
                                tmp = tmp >> 1;
                            }
                    }

                    return (short) tmp;
                } catch (ParseException e) {
                    throw new RuntimeException("Something wen't wrong during the CRC check", e);
                }
            }
        }

        return 0;
    }

}
