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
package org.apache.plc4x.java.ads.api.commands;

import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.AdsData;
import org.apache.plc4x.java.ads.api.generic.AmsHeader;
import org.apache.plc4x.java.ads.api.generic.types.*;
import org.apache.plc4x.java.ads.api.util.ByteValue;

import static java.util.Objects.requireNonNull;

/**
 * A notification is created in an ADS device.
 * <p>
 * Note: We recommend to announce not more then 550 notifications per device.
 * You can increase the payload by organizing the data in structures.
 */
@AdsCommandType(Command.ADS_ADD_DEVICE_NOTIFICATION)
public class AdsAddDeviceNotificationRequest extends AdsAbstractRequest {

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
    private final Reserved reserved;

    private AdsAddDeviceNotificationRequest(AmsHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length, TransmissionMode transmissionMode, MaxDelay maxDelay, CycleTime cycleTime) {
        super(amsHeader);
        this.indexGroup = requireNonNull(indexGroup);
        this.indexOffset = requireNonNull(indexOffset);
        this.length = requireNonNull(length);
        this.transmissionMode = requireNonNull(transmissionMode);
        this.maxDelay = requireNonNull(maxDelay);
        this.cycleTime = requireNonNull(cycleTime);
        this.reserved = Reserved.INSTANCE;
    }

    private AdsAddDeviceNotificationRequest(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, IndexGroup indexGroup, IndexOffset indexOffset, Length length, TransmissionMode transmissionMode, MaxDelay maxDelay, CycleTime cycleTime) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, State.DEFAULT, invokeId);
        this.indexGroup = requireNonNull(indexGroup);
        this.indexOffset = requireNonNull(indexOffset);
        this.length = requireNonNull(length);
        this.transmissionMode = requireNonNull(transmissionMode);
        this.maxDelay = requireNonNull(maxDelay);
        this.cycleTime = requireNonNull(cycleTime);
        this.reserved = Reserved.INSTANCE;
    }

    @Override
    public AdsData getAdsData() {
        return buildADSData(indexGroup, indexOffset, length, transmissionMode, maxDelay, cycleTime, reserved);
    }

    public static AdsAddDeviceNotificationRequest of(AmsHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length, TransmissionMode transmissionMode, MaxDelay maxDelay, CycleTime cycleTime) {
        return new AdsAddDeviceNotificationRequest(amsHeader, indexGroup, indexOffset, length, transmissionMode, maxDelay, cycleTime);
    }

    public static AdsAddDeviceNotificationRequest of(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, IndexGroup indexGroup, IndexOffset indexOffset, Length length, TransmissionMode transmissionMode, MaxDelay maxDelay, CycleTime cycleTime) {
        return new AdsAddDeviceNotificationRequest(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, indexGroup, indexOffset, length, transmissionMode, maxDelay, cycleTime);
    }

    public IndexGroup getIndexGroup() {
        return indexGroup;
    }

    public IndexOffset getIndexOffset() {
        return indexOffset;
    }

    public Length getLength() {
        return length;
    }

    public TransmissionMode getTransmissionMode() {
        return transmissionMode;
    }

    public MaxDelay getMaxDelay() {
        return maxDelay;
    }

    public CycleTime getCycleTime() {
        return cycleTime;
    }

    public Reserved getReserved() {
        return reserved;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof AdsAddDeviceNotificationRequest))
            return false;
        if (!super.equals(o))
            return false;

        AdsAddDeviceNotificationRequest that = (AdsAddDeviceNotificationRequest) o;

        if (!indexGroup.equals(that.indexGroup))
            return false;
        if (!indexOffset.equals(that.indexOffset))
            return false;
        if (!length.equals(that.length))
            return false;
        if (!transmissionMode.equals(that.transmissionMode))
            return false;
        if (!maxDelay.equals(that.maxDelay))
            return false;
        if (!cycleTime.equals(that.cycleTime))
            return false;

        return reserved.equals(that.reserved);
    }

    @Override
    public int hashCode() {
        int result = super.hashCode();
        result = 31 * result + indexGroup.hashCode();
        result = 31 * result + indexOffset.hashCode();
        result = 31 * result + length.hashCode();
        result = 31 * result + transmissionMode.hashCode();
        result = 31 * result + maxDelay.hashCode();
        result = 31 * result + cycleTime.hashCode();
        result = 31 * result + reserved.hashCode();
        return result;
    }

    @Override
    public String toString() {
        return "AdsAddDeviceNotificationRequest{" +
            "indexGroup=" + indexGroup +
            ", indexOffset=" + indexOffset +
            ", length=" + length +
            ", transmissionMode=" + transmissionMode +
            ", maxDelay=" + maxDelay +
            ", cycleTime=" + cycleTime +
            ", reserved=" + reserved +
            "} " + super.toString();
    }

    public static class Reserved extends ByteValue {

        public static final int NUM_BYTES = 16;

        private static final Reserved INSTANCE = new Reserved();

        private Reserved() {
            super((byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00);
            assertLength(NUM_BYTES);
        }
    }
}
