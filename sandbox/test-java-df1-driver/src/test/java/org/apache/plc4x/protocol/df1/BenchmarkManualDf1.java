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

import com.fazecast.jSerialComm.SerialPort;

public class BenchmarkManualDf1 {

    public static void main(String[] args) throws Exception {
//        byte[] rData = Hex.decodeHex("0300006702f080320100000001005600000407120a10060001032b84000160120a10020001032b840001a0120a10010001032b840001a9120a10050001032b84000150120a10020001032b84000198120a10040001032b84000140120a10020001032b84000190");
//        long start = System.currentTimeMillis();
//        int numRunsParse = 2000000;
//        TPKTPacketIO tpktPacketIO = new TPKTPacketIO();
//
//
//        CommPortIdentifier id = CommPortIdentifier.getPortIdentifier("");
//
//
//        // Benchmark the parsing code
//        TPKTPacket packet = null;
//        for(int i = 0; i < numRunsParse; i++) {
//            ReadBuffer rBuf = new ReadBuffer(rData);
//            packet = tpktPacketIO.parse(rBuf);
//        }
//        long endParsing = System.currentTimeMillis();
//
//        System.out.println("Parsed " + numRunsParse + " packets in " + (endParsing - start) + "ms");
//        System.out.println("That's " + ((float) (endParsing - start) / numRunsParse) + "ms per packet");
//
//        // Benchmark the serializing code
//        int numRunsSerialize = 2000000;
//        byte[] oData = null;
//        for(int i = 0; i < numRunsSerialize; i++) {
//            WriteBuffer wBuf = new WriteBuffer(packet.getLengthInBytes());
//            tpktPacketIO.serialize(wBuf, packet);
//            oData = wBuf.getData();
//        }
//        long endSerializing = System.currentTimeMillis();
//
//        System.out.println("Serialized " + numRunsSerialize + " packets in " + (endSerializing - endParsing) + "ms");
//        System.out.println("That's " + ((float) (endSerializing - endParsing) / numRunsSerialize) + "ms per packet");
//        if(!Arrays.equals(rData, oData)) {
//            for(int i = 0; i < rData.length; i++) {
//                if(rData[i] != oData[i]) {
//                    System.out.println("Difference in byte " + i);
//                }
//            }
//            System.out.println("Not equals");
//        } else {
//            System.out.println("Bytes equal");
//        }


//        byte[] rData = {0x10, 0x02, 0x00, 0x09, 0x41, 0x00, 0x01, 0x00, 0x1F, 0x1F, 0x10, 0x03, 0x1A, 0x2B};
//
//        DF1SymbolIO df1SymbolIO = new DF1SymbolIO();
//        DF1Symbol packet;
//        ReadBuffer rBuf = new ReadBuffer(rData);
//        int statusWord = (rData[7]<<8) + rData[6];
//        DF1Command messageCommand = new DF1Command((short)rData[5]); //,(short)statusWord);
//        DF1SymbolMessageFrameStart messageStart = new DF1SymbolMessageFrameStart((short)rData[3],(short)rData[2], messageCommand);
//        packet = df1SymbolIO.parse(rBuf, (short) (rData.length-12), messageStart);
//
//        System.out.println("x: " + packet);
//        System.exit(0);



        SerialPort comPort = SerialPort.getCommPorts()[0];
        comPort.openPort();

        comPort.setComPortParameters(19200, 8, SerialPort.ONE_STOP_BIT, SerialPort.NO_PARITY);
//        comPort.setComPortTimeouts(SerialPort.TIMEOUT_NONBLOCKING, 0, 0);

        System.out.print(comPort.getSystemPortName() + " | ");
        //System.out.print(comPort.getPortDescription() + " | ");
        System.out.print(comPort.getDescriptivePortName() + " | Baud rate: ");
        System.out.println(comPort.getBaudRate());
//        System.out.println(comPort.getReadTimeout());
//        System.out.println(comPort.getWriteTimeout());



//        DF1SymbolIO df1message = new DF1SymbolIO();


        byte[] c_STX = {0x02};
        byte[] c_SOH = {0x01};
        byte[] c_ETX = {0x03};
        byte[] c_EOT = {0x04};
        byte[] c_ENQ = {0x05};
        byte[] c_ACK = {0x06};
        byte[] c_DLE = {0x10};
        byte[] c_NAK = {0x0f};

        byte[] c_DST = {0x09};
        byte[] c_SRC = {0x00};
        byte[] c_CMD = {0x06};
        byte[] c_FNC = {0x03};
        byte[] c_STS = {0x00};
        byte[] c_TNS = {0x01, 0x00};
        byte[] c_ADR = {0x11, 0x00};
        byte[] c_SZE = {0x02};
//        byte[] c_BCC = {-30};

//        byte[] c_STN = {0x11};
//        byte[] c_DST = {0x01};
//        byte[] c_SRC = {0x00};
//        byte[] c_SZE = {0x0c};
//        byte[] c_TNS = {0x41, 0x00};
//        byte[] c_ADR = {0x12, 0x00};
//        byte[] c_CRC = {-49, 0x40};


//        byte[] message = {0x10, 0x01, 0x11, 0x10, 0x02, 0x09, 0x00, 0x01, 0x00, 0x41, 0x00, 0x12, 0x00, 0x0c, 0x10, 0x03}; // halfduplex msg
//        comPort.writeBytes(message, message.length);

//        byte[] message = {0x10, 0x05, 0x11, -17};         // halfduplex poll
//        comPort.writeBytes(message, message.length);


//        comPort.writeBytes(c_DLE, 1);     // fullduplex msg Seite 235
//        comPort.writeBytes(c_STX, 1);
//        comPort.writeBytes(c_DST, 1);
//        comPort.writeBytes(c_SRC, 1);
//        comPort.writeBytes(c_CMD, 1);
//        comPort.writeBytes(c_STS, 1);
//        comPort.writeBytes(c_TNS, 2);
//        comPort.writeBytes(c_FNC, 1);
//        comPort.writeBytes(c_ADR, 2);
//        comPort.writeBytes(c_SZE, 1);
//        comPort.writeBytes(c_DLE, 1);
//        comPort.writeBytes(c_ETX, 1);

        byte[] msg = {0x10, 0x02,   0x09, 0x00,   0x01, 0x00,   0x01, 0x00,   0x0b, 0x00,   0x02,   0x10, 0x03};
        comPort.writeBytes(msg, 13);


//        int[] crcmsg = {c_DST[0], c_SRC[0], c_CMD[0], c_STS[0], c_TNS[0], c_TNS[1], c_ADR[0], c_ADR[1], c_SZE[0], c_ETX[0]}; // fullduplex CRC
//        int[] crcmsg = {c_DST[0], c_SRC[0], c_CMD[0], c_STS[0], c_TNS[0], c_TNS[1], c_FNC[0], c_ETX[0]};                       // diagnostic status request
//        int[] crcmsg = {0x11, 0x02, 0x09, 0x00, 0x01, 0x00, 0x41, 0x00, 0x12, 0x00, 0x0c, 0x03}; // halfduplex CRC
        int[] crcmsg = { 0x09, 0x00,   0x01, 0x00,   0x01, 0x00,   0x0b, 0x00,   0x02,   0x03};

        int[] c_CRC = CalcCRC(crcmsg);
        byte[] crc1 = {(byte)c_CRC[0]};
        byte[] crc2 = {(byte)c_CRC[1]};
        System.out.println("crc1: " + Integer.toHexString(crc1[0]));
        System.out.println("crc2: " + Integer.toHexString(crc2[0]));
        comPort.writeBytes(crc1, 1);
        comPort.writeBytes(crc2, 1);


        while (comPort.bytesAvailable() == 0) {
            Thread.sleep(22); }

        byte[] readBuffer = new byte[comPort.bytesAvailable()];
        int numRead = comPort.readBytes(readBuffer, readBuffer.length);
        System.out.println("Read " + numRead + " bytes.");

        for (byte c_RCV : readBuffer) {
            System.out.print(Integer.toHexString(c_RCV) + " | "); }
        System.out.println("");

//        if (numRead > 1) {
//            if (readBuffer[1] != 0x15) {
//                comPort.writeBytes(c_DLE, 1);
//                comPort.writeBytes(c_ACK, 1);
//            }
//        }

        while (comPort.bytesAvailable() == 0) {
            Thread.sleep(22); }


        byte[] readBuffer2 = new byte[comPort.bytesAvailable()];
        int numRead2 = comPort.readBytes(readBuffer2, readBuffer2.length);
        System.out.println("Read " + numRead2 + " bytes.");

        for (byte c_RCV2 : readBuffer2) {
            System.out.print(Integer.toHexString(c_RCV2) + " | "); }
        System.out.println("");

        comPort.closePort();

    }

        private static int[] CalcCRC(int[] crcBytes) {
            int tmp = 0;
            int crcL, crcR;

            for (int newByte : crcBytes ) {
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

//            return ((tmp & 0xFF) << 8) + (tmp >> 8);  // returns lowbyte|highbyte as one number, change return to non-array

            int[] tmparray = {(tmp & 0xFF), (tmp >> 8)};
            return tmparray;
        }

        private static int CalcBCC(int[] crcBytes) {
            int tmp = 0;
            int j = 0;

            for (int newByte : crcBytes) {
                tmp = tmp + newByte;
                if (newByte == 0x10) {
                    j = ++j; }
            }
            tmp = tmp - ((j/2) * 0x10);  // get rid of double DLE
            return ((~tmp) & 0b11111111) + 1;
        }
}
