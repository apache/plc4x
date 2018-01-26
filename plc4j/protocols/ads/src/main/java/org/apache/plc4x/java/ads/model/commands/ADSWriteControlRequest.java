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

import org.apache.plc4x.java.ads.model.commands.types.ADSState;
import org.apache.plc4x.java.ads.model.commands.types.Data;
import org.apache.plc4x.java.ads.model.commands.types.DeviceState;
import org.apache.plc4x.java.ads.model.commands.types.Length;
import org.apache.plc4x.java.ads.model.generic.ADSData;
import org.apache.plc4x.java.ads.model.generic.AMSHeader;
import org.apache.plc4x.java.ads.model.generic.AMSTCPHeader;
import org.apache.plc4x.java.ads.model.generic.AMSTCPPaket;

/**
 * Changes the ADS status and the device status of an ADS device.
 * Additionally it is possible to send data to the ADS device to transfer further information.
 * These data were not analysed from the current ADS devices (PLC, NC, ...)
 */
public class ADSWriteControlRequest extends AMSTCPPaket {

    /**
     * 2 bytes	New ADS status (see data type ADSSTATE of the ADS-DLL).
     */
    private final ADSState adsState;
    /**
     * 2 bytes	New device status.
     */
    private final DeviceState deviceState;
    /**
     * 4 bytes	Length of data in byte.
     */
    private final Length length;
    /**
     * n bytes	Additional data which are sent to the ADS device
     */
    private final Data data;

    public ADSWriteControlRequest(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, ADSState adsState, DeviceState deviceState, Length length, Data data) {
        super(amstcpHeader, amsHeader);
        this.adsState = adsState;
        this.deviceState = deviceState;
        this.length = length;
        this.data = data;
    }

    @Override
    public ADSData getAdsData() {
        return buildADSData(adsState, deviceState, length, data);
    }

}
