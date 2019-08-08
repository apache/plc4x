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

import org.apache.plc4x.java.df1.*;
import org.apache.plc4x.java.utils.ParseException;
import org.apache.plc4x.java.utils.ReadBuffer;
import org.apache.plc4x.java.utils.WriteBuffer;

public class DF1Utils {

    public static short crcCheck(Object... args) {
        short destinationAddress = (short) args[0];
        short sourceAddress = (short) args[1];
        DF1Command command = (DF1Command) args[2];
        short commandDiscriminatorValues = (short) command.getDiscriminatorValues()[0];
        short status = command.getStatus();
        int transactionCounter = command.getTransactionCounter();
        if(command instanceof DF1UnprotectedReadRequest) {
            try {
                DF1UnprotectedReadRequest unprotectedReadRequest = (DF1UnprotectedReadRequest) command;
                WriteBuffer writeBuffer = new WriteBuffer(10, false);
                writeBuffer.writeUnsignedShort(8, destinationAddress);
                writeBuffer.writeUnsignedShort(8, sourceAddress);
                writeBuffer.writeUnsignedShort(8, commandDiscriminatorValues);
                writeBuffer.writeUnsignedShort(8, status);
                writeBuffer.writeUnsignedInt(16, (short) transactionCounter);
                writeBuffer.writeUnsignedInt(16, (short) unprotectedReadRequest.getAddress());
                writeBuffer.writeUnsignedShort(8, (byte) unprotectedReadRequest.getSize());
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
        } else if(command instanceof DF1UnprotectedReadResponse) {
            DF1UnprotectedReadResponse unprotectedReadResponseCommand = (DF1UnprotectedReadResponse) command;
            try {
                WriteBuffer writeBuffer = new WriteBuffer(10, false);
                writeBuffer.writeUnsignedShort(8, destinationAddress);
                writeBuffer.writeUnsignedShort(8, sourceAddress);
                writeBuffer.writeUnsignedShort(8, commandDiscriminatorValues);
                writeBuffer.writeUnsignedShort(8, status);
                writeBuffer.writeUnsignedInt(16, (short) transactionCounter);
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
