/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.ads.api.commands;

import org.apache.plc4x.java.ads.api.commands.types.AdsStampHeader;
import org.apache.plc4x.java.ads.api.commands.types.Length;
import org.apache.plc4x.java.ads.api.commands.types.Stamps;
import org.apache.plc4x.java.ads.api.generic.AdsData;
import org.apache.plc4x.java.ads.api.generic.AmsHeader;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.api.util.ByteReadable;
import org.apache.plc4x.java.ads.api.util.LengthSupplier;

import java.util.List;

import static java.util.Objects.requireNonNull;

/**
 * Data will carry forward independently from an ADS device to a Client
 * <p>
 * The data which are transfered at the Device Notification  are multiple nested into one another. The Notification Stream contains an array with elements of type AdsStampHeader. This array again contains elements of type AdsNotificationSample.
 */
@AdsCommandType(Command.ADS_DEVICE_NOTIFICATION)
public class AdsDeviceNotificationRequest extends AdsAbstractRequest {

    /**
     * 4 bytes	Size of data in byte.
     */
    private final Length length;
    /**
     * 4 bytes	Number of elements of type AdsStampHeader.
     */
    private final Stamps stamps;
    /**
     * n bytes	Array with elements of type AdsStampHeader.
     */
    private final List<AdsStampHeader> adsStampHeaders;

    private final transient LengthSupplier lengthSupplier;

    private AdsDeviceNotificationRequest(AmsHeader amsHeader, Length length, Stamps stamps, List<AdsStampHeader> adsStampHeaders) {
        super(amsHeader);
        this.length = requireNonNull(length);
        this.stamps = requireNonNull(stamps);
        this.adsStampHeaders = requireNonNull(adsStampHeaders);
        lengthSupplier = null;
    }

    private AdsDeviceNotificationRequest(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, Stamps stamps, List<AdsStampHeader> adsStampHeaders) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.length = null;
        this.stamps = requireNonNull(stamps);
        this.adsStampHeaders = requireNonNull(adsStampHeaders);
        this.lengthSupplier = () -> {
            long aggregateLength = 0;
            for (LengthSupplier supplier : adsStampHeaders) {
                aggregateLength += supplier.getCalculatedLength();
            }
            return aggregateLength + Stamps.NUM_BYTES;
        };
    }

    public static AdsDeviceNotificationRequest of(AmsHeader amsHeader, Length length, Stamps stamps, List<AdsStampHeader> adsStampHeaders) {
        return new AdsDeviceNotificationRequest(amsHeader, length, stamps, adsStampHeaders);
    }

    public static AdsDeviceNotificationRequest of(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, Stamps stamps, List<AdsStampHeader> adsStampHeaders) {
        return new AdsDeviceNotificationRequest(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, stamps, adsStampHeaders);
    }

    public Length getLength() {
        return lengthSupplier == null ? length : Length.of(lengthSupplier);
    }

    public Stamps getStamps() {
        return stamps;
    }

    public List<AdsStampHeader> getAdsStampHeaders() {
        return adsStampHeaders;
    }

    public LengthSupplier getLengthSupplier() {
        return lengthSupplier;
    }

    @Override
    public AdsData getAdsData() {
        return buildADSData(getLength(), stamps, buildADSData(adsStampHeaders.toArray(new ByteReadable[adsStampHeaders.size()])));
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof AdsDeviceNotificationRequest))
            return false;
        if (!super.equals(o))
            return false;

        AdsDeviceNotificationRequest that = (AdsDeviceNotificationRequest) o;

        if (!getLength().equals(that.getLength()))
            return false;
        if (!stamps.equals(that.stamps))
            return false;

        return adsStampHeaders.equals(that.adsStampHeaders);
    }

    @Override
    public int hashCode() {
        int result = super.hashCode();
        result = 31 * result + getLength().hashCode();
        result = 31 * result + stamps.hashCode();
        result = 31 * result + adsStampHeaders.hashCode();
        return result;
    }

    @Override
    public String toString() {
        return "AdsDeviceNotificationRequest{" +
            "length=" + getLength() +
            ", stamps=" + stamps +
            ", adsStampHeaders=" + adsStampHeaders +
            "} " + super.toString();
    }
}
