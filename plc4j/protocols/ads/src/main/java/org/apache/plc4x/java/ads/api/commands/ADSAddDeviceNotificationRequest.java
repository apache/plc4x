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
package org.apache.plc4x.java.ads.api.commands;

import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.commands.types.Length;
import org.apache.plc4x.java.ads.api.generic.ADSData;
import org.apache.plc4x.java.ads.api.generic.AMSHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPPaket;
import org.apache.plc4x.java.ads.api.generic.types.*;
import org.apache.plc4x.java.ads.api.util.ByteValue;

/**
 * A notification is created in an ADS device.
 * <p>
 * Note: We recommend to announce not more then 550 notifications per device.
 * You can increase the payload by organizing the data in structures.
 */
public class ADSAddDeviceNotificationRequest extends AMSTCPPaket {

    /**
     * 4 bytes	Index Group of the data, which should be sent per notification.
     */
    private final IndexGroup indexGroup;
    /**
     * 4 bytes	Index Offset of the data, which should be sent per notification.
     */
    private final IndexOffset indexOffset;
    /**
     * 4 bytes	Length of data in bytes, which should be sent per notification.
     */
    private final Length length;
    /**
     * 4 bytes	See description of the structure ADSTRANSMODE at the ADS-DLL.
     */
    private final TransmissionMode transmissionMode;
    /**
     * 4 bytes	At the latest after this time, the ADS Device Notification is called. The unit is 1ms.
     */
    private final MaxDelay maxDelay;
    /**
     * 4 bytes	The ADS server checks if the value changes in this time slice. The unit is 1ms
     */
    private final CycleTime cycleTime;
    /**
     * 16bytes	Must be set to 0
     */
    private final Reserved reserved = Reserved.INSTANCE;

    public ADSAddDeviceNotificationRequest(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length, TransmissionMode transmissionMode, MaxDelay maxDelay, CycleTime cycleTime) {
        super(amstcpHeader, amsHeader);
        this.indexGroup = indexGroup;
        this.indexOffset = indexOffset;
        this.length = length;
        this.transmissionMode = transmissionMode;
        this.maxDelay = maxDelay;
        this.cycleTime = cycleTime;
    }

    public ADSAddDeviceNotificationRequest(AMSHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length, TransmissionMode transmissionMode, MaxDelay maxDelay, CycleTime cycleTime) {
        super(amsHeader);
        this.indexGroup = indexGroup;
        this.indexOffset = indexOffset;
        this.length = length;
        this.transmissionMode = transmissionMode;
        this.maxDelay = maxDelay;
        this.cycleTime = cycleTime;
    }

    public ADSAddDeviceNotificationRequest(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, IndexGroup indexGroup, IndexOffset indexOffset, Length length, TransmissionMode transmissionMode, MaxDelay maxDelay, CycleTime cycleTime) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.indexGroup = indexGroup;
        this.indexOffset = indexOffset;
        this.length = length;
        this.transmissionMode = transmissionMode;
        this.maxDelay = maxDelay;
        this.cycleTime = cycleTime;
    }

    @Override
    public ADSData getAdsData() {
        return buildADSData(indexGroup, indexOffset, length, transmissionMode, maxDelay, cycleTime, reserved);
    }

    @Override
    public Command getCommandId() {
        return Command.ADS_Add_Device_Notification;
    }

    @Override
    public State getStateId() {
        return State.ADS_REQUEST_TCP;
    }

    public static class Reserved extends ByteValue {

        private static final Reserved INSTANCE = new Reserved();

        private Reserved() {
            super((byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00);
            assertLength(16);
        }
    }
}
