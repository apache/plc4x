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

import org.apache.plc4x.java.df1.readwrite.DF1Command;
import org.apache.plc4x.java.df1.readwrite.DF1UnprotectedReadRequest;
import org.apache.plc4x.java.df1.readwrite.DF1UnprotectedReadResponse;
import org.apache.plc4x.java.spi.generation.*;

public class DF1Utils {

    public static int crcCheck(Object... args) {
        short destinationAddress = (short) args[0];
        short sourceAddress = (short) args[1];
        DF1Command command = (DF1Command) args[2];
        short status = command.getStatus();
        int transactionCounter = command.getTransactionCounter();
        if(command instanceof DF1UnprotectedReadRequest) {
            try {
                DF1UnprotectedReadRequest unprotectedReadRequest = (DF1UnprotectedReadRequest) command;
                WriteBufferByteBased writeBuffer = new WriteBufferByteBased(10, false);
                writeBuffer.writeUnsignedShort(8, destinationAddress);
                writeBuffer.writeUnsignedShort(8, sourceAddress);
                writeBuffer.writeUnsignedShort(8, command.getCommandCode());
                writeBuffer.writeUnsignedShort(8, status);
                writeBuffer.writeUnsignedInt(16, (short) transactionCounter);
                writeBuffer.writeUnsignedInt(16, (short) unprotectedReadRequest.getAddress());
                writeBuffer.writeUnsignedShort(8, (byte) unprotectedReadRequest.getSize());
                writeBuffer.writeUnsignedShort(8, (byte) 0x03);

                byte[] data = writeBuffer.getData();
                return calculateCRC(data) & 0xFFFF;

            } catch (ParseException e) {
                throw new RuntimeException("Something went wrong during the CRC check", e);
            }
        } else if(command instanceof DF1UnprotectedReadResponse) {
            DF1UnprotectedReadResponse unprotectedReadResponseCommand = (DF1UnprotectedReadResponse) command;
            try {
                // TODO: size has to be dependent on actual size requested
                WriteBufferByteBased writeBuffer = new WriteBufferByteBased(10, false);
                writeBuffer.writeUnsignedShort(8, destinationAddress);
                writeBuffer.writeUnsignedShort(8, sourceAddress);
                writeBuffer.writeUnsignedShort(8, command.getCommandCode());
                writeBuffer.writeUnsignedShort(8, status);
                writeBuffer.writeUnsignedInt(16, (short) transactionCounter);
                boolean escape10 = false;
                for (short data : unprotectedReadResponseCommand.getData()) {
//                    if (escape10 == true){
//                        if (data == 0x10) {
//
//                        }
//                    } else{
//
//                    }
                    writeBuffer.writeUnsignedShort(8,  data);
                }
                writeBuffer.writeUnsignedShort(8, (byte) 0x03);

                byte[] data = writeBuffer.getData();
                return calculateCRC(data) & 0xFFFF;

            } catch (ParseException e) {
                throw new RuntimeException("Something went wrong during the CRC check", e);
            }
        }

        return 0;
    }

    public static boolean dataTerminate(ReadBuffer io) {
        ReadBufferByteBased rbbb = (ReadBufferByteBased)io;
        try {
            // The byte sequence 0x10 followed by 0x03 indicates the end of the message,
            // so if we would read this, we abort the loop and stop reading data.
            if ((rbbb.peekByte(0) == (byte) 0x10) && (rbbb.peekByte(1) == (byte) 0x03)) {
                return true;
            }
        } catch (ParseException e) {
            // Just ignore and return false.
        }
        return false;
    }

    public static short readData(ReadBuffer io) {
        ReadBufferByteBased rbbb = (ReadBufferByteBased)io;
        try {
            // If we read a 0x10, this has to be followed by another 0x10, which is how
            // this value is escaped in DF1, so if we encounter two 0x10, we simply ignore the first.
            if ((rbbb.peekByte(0) == (byte) 0x10) && (rbbb.peekByte(1) == 0x10)) {
                io.readByte(8);
            }
            return io.readUnsignedShort(8);
        } catch (ParseException e) {
            throw new RuntimeException("Error parsing data", e);
        }
    }

    public static void writeData(WriteBuffer io, short data) {
        try {
            // If a value is 0x10, this has to be duplicated in order to be escaped.
            if ((data == (short) 0x10)) {
                io.writeUnsignedShort(8, (short) 0x10);
            }
            io.writeUnsignedShort(8, data);
        } catch (ParseException e) {
            throw new RuntimeException("Error parsing data", e);
        }
    }

    public static int dataLength(short[] data) {
        int i = 0;
        for(short dataByte : data) {
            // If a value is 0x10, this has to be duplicated which increases the message size by one.
            if(dataByte == 0x10) {
                i++;
            }
            i++;
        }
        return i;
    }

    private static short calculateCRC(byte[] crcBytes) {
        int tmp = 0;
        int crcL, crcR;

        for (int newByte : crcBytes) {
            crcL = tmp >> 8;
            crcR = tmp & 0xFF;
            tmp = (crcL << 8) + (newByte ^ crcR);
            for (int j = 0; j < 8; j++)
                if (tmp % 2 == 1) {     // check if LSB shifted out is 1 or 0
                    tmp = tmp >> 1;
                    tmp = tmp ^ 0xA001;
                } else {
                    tmp = tmp >> 1;
                }
        }
        return Short.reverseBytes((short) tmp);
    }
}
