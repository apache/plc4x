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
import org.apache.plc4x.java.ads.api.generic.ADSData;
import org.apache.plc4x.java.ads.api.generic.AMSHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPHeader;
import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;

/**
 * Reads the name and the version number of the ADS device.
 */
@ADSCommandType(Command.ADS_Read_Device_Info)
public class ADSReadDeviceInfoResponse extends ADSAbstractResponse {
    /**
     * 4 bytes	ADS error number.
     */
    private final Result result;
    /**
     * Version	1 byte	Major version number
     */
    private final MajorVersion majorVersion;
    /**
     * Version	1 byte	Minor version number
     */
    private final MinorVersion minorVersion;
    /**
     * Build	2 bytes	Build number
     */
    private final Version version;
    /**
     * Name	16 bytes	Name of ADS device
     */
    private final Device device;

    public ADSReadDeviceInfoResponse(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, Result result, MajorVersion majorVersion, MinorVersion minorVersion, Version version, Device device) {
        super(amstcpHeader, amsHeader);
        this.result = result;
        this.majorVersion = majorVersion;
        this.minorVersion = minorVersion;
        this.version = version;
        this.device = device;
    }

    public ADSReadDeviceInfoResponse(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, Result result, MajorVersion majorVersion, MinorVersion minorVersion, Version version, Device device) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.result = result;
        this.majorVersion = majorVersion;
        this.minorVersion = minorVersion;
        this.version = version;
        this.device = device;
    }

    @Override
    public ADSData getAdsData() {
        return buildADSData(result, majorVersion, minorVersion, version, device);
    }

    @Override
    public String toString() {
        return "ADSReadDeviceInfoResponse{" +
            "result=" + result +
            ", majorVersion=" + majorVersion +
            ", minorVersion=" + minorVersion +
            ", version=" + version +
            ", device=" + device +
            "} " + super.toString();
    }
}
