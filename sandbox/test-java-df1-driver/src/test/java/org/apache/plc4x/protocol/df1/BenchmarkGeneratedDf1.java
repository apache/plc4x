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

public class BenchmarkGeneratedDf1 {

    public static void main(String[] args) throws Exception {
        /*// Manually build a message
        ReadRequest readRequest = new ReadRequest(new DF1SymbolMessageFrameStart((short) 0x09, (short) 0x00, new DF1ReadRequest((short) 0x00, 0x01, 0x0B, (short) 0x02)), new DF1SymbolMessageFrameEnd());

        // Serialize the message
        WriteBuffer writeBuffer = new WriteBuffer(100, false);
        new ReadRequestIO().serialize(writeBuffer, readRequest);
        byte[] data = writeBuffer.getData();

        // Send the serialized message to the PLC via COM port
        SerialPort comPort = SerialPort.getCommPort("/dev/cu.usbserial-AL065SUZ");
        comPort.openPort();
        comPort.setComPortParameters(19200, 8, SerialPort.ONE_STOP_BIT, SerialPort.NO_PARITY);
        comPort.writeBytes(data, 20);

        // Give the PLC some time to respond.
        while (comPort.bytesAvailable() == 0) {
            Thread.sleep(22);
        }

        // Read the response
        byte[] readBytes = new byte[comPort.bytesAvailable()];
        comPort.readBytes(readBytes, readBytes.length);

        // Parse the ACK/NACK response
        ReadBuffer readBuffer = new ReadBuffer(readBytes);
        Result result = new ResultIO().parse(readBuffer);

        // Check if the response was an ACK
        if(result.getResult() instanceof DF1SymbolMessageFrameACK) {
            // The actual result is sent with a little delay.
            while (comPort.bytesAvailable() == 0) {
                Thread.sleep(22);
            }

            // Read the actual response data
            readBytes = new byte[comPort.bytesAvailable()];
            comPort.readBytes(readBytes, readBytes.length);

            // Parse the response
            readBuffer = new ReadBuffer(readBytes);
            ReadResponse readResponse = new ReadResponseIO().parse(readBuffer, (short) 0x02);

            // So something senseless ;-)
            System.out.println(readResponse);
        } else {
            System.out.println("Didn't get an ACK");
        }

        comPort.closePort();*/
    }

}
