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
package org.apache.plc4x.java.df1.util;

import org.apache.plc4x.java.df1.DF1Symbol;
import org.apache.plc4x.java.df1.DF1SymbolMessageFrame;
import org.apache.plc4x.java.df1.DF1UnprotectedReadRequest;
import org.apache.plc4x.java.df1.DF1UnprotectedReadResponse;
import org.apache.plc4x.java.utils.ParseException;
import org.apache.plc4x.java.utils.ReadBuffer;
import org.apache.plc4x.java.utils.WriteBuffer;

public class DF1Utils {

    public static short CRCCheck(Object... args) {
        DF1Symbol symbol = (DF1Symbol) args[0];
        if(symbol instanceof DF1SymbolMessageFrame) {
            DF1SymbolMessageFrame messageFrame = (DF1SymbolMessageFrame) symbol;

            short destinationAddress = messageFrame.getDestinationAddress();
            short sourceAddress = messageFrame.getSourceAddress();
            short commandDiscriminatorValues = (short) messageFrame.getCommand().getDiscriminatorValues()[0];
            short status = messageFrame.getCommand().getStatus();
            int   counter = messageFrame.getCommand().getTransactionCounter();
            if(messageFrame.getCommand() instanceof DF1UnprotectedReadRequest) {
                DF1UnprotectedReadRequest unprotectedReadRequestCommand = (DF1UnprotectedReadRequest) messageFrame.getCommand();
                try {
                    WriteBuffer writeBuffer = new WriteBuffer(10, false);
                    writeBuffer.writeUnsignedShort(8, destinationAddress);
                    writeBuffer.writeUnsignedShort(8, sourceAddress);
                    writeBuffer.writeUnsignedShort(8, commandDiscriminatorValues);
                    writeBuffer.writeUnsignedShort(8, status);
                    writeBuffer.writeUnsignedInt(16, (short) counter);
                    writeBuffer.writeUnsignedInt(16, (short) unprotectedReadRequestCommand.getAddress());
                    writeBuffer.writeUnsignedShort(8, (byte) unprotectedReadRequestCommand.getSize());
                    writeBuffer.writeUnsignedShort(8, (byte) 0x03);

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
            } else if(messageFrame.getCommand() instanceof DF1UnprotectedReadResponse) {
                DF1UnprotectedReadResponse unprotectedReadResponseCommand = (DF1UnprotectedReadResponse) messageFrame.getCommand();
                try {
                    WriteBuffer writeBuffer = new WriteBuffer(10, false);
                    writeBuffer.writeUnsignedShort(8, destinationAddress);
                    writeBuffer.writeUnsignedShort(8, sourceAddress);
                    writeBuffer.writeUnsignedShort(8, commandDiscriminatorValues);
                    writeBuffer.writeUnsignedShort(8, status);
                    writeBuffer.writeUnsignedInt(16, (short) counter);
                    for (short data : unprotectedReadResponseCommand.getData()) {
                        writeBuffer.writeUnsignedShort(8,  data);
                    }
                    writeBuffer.writeUnsignedShort(8, (byte) 0x03);

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

    public static boolean dataTerminate(ReadBuffer io) {
        try {
            if ((io.peekByte(0) == (byte) 0x10) && (io.peekByte(1) == (byte) 0x03)) {
                return true;
            }
        } catch (ParseException e) {
            // Just ignore and return false.
        }
        return false;
    }

}
