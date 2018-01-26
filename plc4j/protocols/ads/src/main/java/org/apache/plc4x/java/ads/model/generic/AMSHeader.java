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
package org.apache.plc4x.java.ads.model.generic;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.plc4x.java.ads.model.util.ByteReadable;
import org.apache.plc4x.java.ads.model.util.ByteValue;

import java.nio.ByteBuffer;

/**
 * AMS Header	32 bytes	The AMS/TCP-Header contains the addresses of the transmitter and receiver. In addition the AMS error code , the ADS command Id and some other information.
 */
public class AMSHeader implements ByteReadable {
    /**
     * This is the AMSNetId of the station, for which the packet is intended. Remarks see below.
     */
    private final AMSNetId targetAmsNetId;
    /**
     * This is the AMSPort of the station, for which the packet is intended.
     */
    private final AMSPort targetAmsPort;
    /**
     * This contains the AMSNetId of the station, from which the packet was sent.
     */
    private final AMSNetId sourceAmsNetId;
    /**
     * This contains the AMSPort of the station, from which the packet was sent.
     */
    private final AMSPort sourceAmsPort;

    private final Command commandId;

    private final State stateFlags;

    private final DataLength dataLength;

    private final Error code;

    private final Invoke invokeId;

    private final Data nData;

    public AMSHeader(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Command commandId, State stateFlags, DataLength dataLength, Error code, Invoke invokeId, Data nData) {
        this.targetAmsNetId = targetAmsNetId;
        this.targetAmsPort = targetAmsPort;
        this.sourceAmsNetId = sourceAmsNetId;
        this.sourceAmsPort = sourceAmsPort;
        this.commandId = commandId;
        this.stateFlags = stateFlags;
        this.dataLength = dataLength;
        this.code = code;
        this.invokeId = invokeId;
        this.nData = nData;
    }

    /**
     * The AMSNetId consists of 6 bytes and addresses the transmitter or receiver. One possible AMSNetId would be e.g.. 172.16.17.10.1.1. The storage arrangement in this example is as follows:
     * <p>
     * _____0     1     2     3     4     5
     * __+-----------------------------------+
     * 0 | 127 |  16 |  17 |  10 |   1 |   1 |
     * __+-----------------------------------+
     * <p>
     * <p>
     * The AMSNetId is purely logical and has usually no relation to the IP address. The AMSNetId is configurated at the target system. At the PC for this the TwinCAT System Control is used. If you use other hardware, see the considering documentation for notes about settings of the AMS NetId.
     */
    public static class AMSNetId extends ByteValue {

        public AMSNetId(int octed1, int octed2, int octed3, int octed4, int octed5, int octed6) {
            super((byte) octed1, (byte) octed2, (byte) octed3, (byte) octed4, (byte) octed5, (byte) octed6);
        }

        public AMSNetId(byte... value) {
            super(value);
            assertLength(6);
        }
    }

    public static class AMSPort extends ByteValue {

        public AMSPort(int port) {
            super(ByteBuffer.allocate(2).putInt(port).array());
        }

        public AMSPort(byte... value) {
            super(value);
            assertLength(2);
        }
    }

    /**
     * 2 bytes	see below.
     */
    enum Command implements ByteReadable {
        Invalid(0x0000),
        ADS_Read_Device_Info(0x0001),
        ADS_Read(0x0002),
        ADS_Write(0x0003),
        ADS_Read_State(0x0004),
        ADS_Write_Control(0x0005),
        ADS_Add_Device_Notification(0x0006),
        ADS_Delete_Device_Notification(0x0007),
        ADS_Device_Notification(0x0008),
        ADS_Read_Write(0x0009),
        /**
         * Other commands are not defined or are used internally. Therefore the Command Id  is only allowed to contain the above enumerated values!
         */
        UNKNOWN(0xffff);
        final byte[] value;

        Command(int value) {
            this.value = ByteBuffer.allocate(4).putInt(value).array();
        }

        @Override
        public byte[] getBytes() {
            return value;
        }

        public ByteBuf getByteBuf() {
            return Unpooled.buffer().writeBytes(value);
        }
    }

    /**
     * 2 bytes	see below.
     * <p>
     * State Flags
     * Flag     Description
     * 0x0001	0: Request / 1: Response
     * 0x0004	ADS command
     * The first bit marks, whether it´s a request or response. The third bit must be set to 1, to exchange data with ADS commands. The other bits are not defined or were used for other internal purposes.
     * <p>
     * Therefore the other bits must be set to 0!
     * <p>
     * Flag     Description
     * 0x000x	TCP Protocol
     * 0x004x	UDP Protocol
     * Bit number 7 marks, if it should be transfered with TCP or UDP.
     */
    enum State implements ByteReadable {
        ADS_REQUEST_TCP(0x0004),
        ADS_RESPONSE_TCP(0x0005),
        ADS_REQUEST_UDP(0x0044),
        ADS_RESPONSE_UDP(0x0045);
        final byte[] value;

        State(int value) {
            this.value = ByteBuffer.allocate(4).putInt(value).array();
        }

        @Override
        public byte[] getBytes() {
            return value;
        }

        public ByteBuf getByteBuf() {
            return Unpooled.buffer().writeBytes(value);
        }
    }

    /**
     * 4 bytes	Size of the data range. The unit is byte.
     */
    static class DataLength extends ByteValue {
        public DataLength(int length) {
            super(ByteBuffer.allocate(4).putInt(length).array());
        }

        public DataLength(byte... value) {
            super(value);
            assertLength(4);
        }
    }

    /**
     * 4 bytes	AMS error number. See ADS Return Codes.
     */
    static class Error extends ByteValue {
        public Error(byte... value) {
            super(value);
            assertLength(4);
        }
    }

    /**
     * 4 bytes	Free usable 32 bit array. Usually this array serves to send an Id. This Id makes is possible to assign a received response to a request, which was sent before.
     */
    static class Invoke extends ByteValue {
        public Invoke(byte... value) {
            super(value);
            assertLength(4);
        }
    }

    /**
     * bytes	Data range. The data range contains the parameter of the considering ADS commands.
     */
    static class Data extends ByteValue {
        public Data(byte... value) {
            super(value);
        }
    }

    @Override
    public ByteBuf getByteBuf() {
        return AMSTCPPaket.buildByteBuff(
            targetAmsNetId,
            targetAmsPort,
            sourceAmsNetId,
            sourceAmsPort,
            commandId,
            stateFlags,
            dataLength,
            code,
            invokeId,
            nData);
    }
}
