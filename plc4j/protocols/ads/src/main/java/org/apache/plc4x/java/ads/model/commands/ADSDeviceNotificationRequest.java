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
package org.apache.plc4x.java.ads.model.commands;

import org.apache.plc4x.java.ads.model.commands.types.AdsStampHeader;
import org.apache.plc4x.java.ads.model.commands.types.Length;
import org.apache.plc4x.java.ads.model.commands.types.Stamps;
import org.apache.plc4x.java.ads.model.generic.ADSData;
import org.apache.plc4x.java.ads.model.generic.AMSHeader;
import org.apache.plc4x.java.ads.model.generic.AMSTCPHeader;
import org.apache.plc4x.java.ads.model.generic.AMSTCPPaket;

/**
 * Data will carry forward independently from an ADS device to a Client
 * <p>
 * The data which are transfered at the Device Notification  are multiple nested into one another. The Notification Stream contains an array with elements of type AdsStampHeader. This array again contains elements of type AdsNotificationSample.
 */
public class ADSDeviceNotificationRequest extends AMSTCPPaket {

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
    private final AdsStampHeader adsStampHeader;

    public ADSDeviceNotificationRequest(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, Length length, Stamps stamps, AdsStampHeader adsStampHeader) {
        super(amstcpHeader, amsHeader);
        this.length = length;
        this.stamps = stamps;
        this.adsStampHeader = adsStampHeader;
    }

    @Override
    public ADSData getAdsData() {
        return buildADSData(length, stamps, adsStampHeader);
    }

}
