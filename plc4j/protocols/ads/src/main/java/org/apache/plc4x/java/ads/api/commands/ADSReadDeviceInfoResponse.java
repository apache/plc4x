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
import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;

import static java.util.Objects.requireNonNull;

/**
 * Reads the name and the version number of the ADS device.
 */
@ADSCommandType(Command.ADS_READ_DEVICE_INFO)
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

    private ADSReadDeviceInfoResponse(AMSHeader amsHeader, Result result, MajorVersion majorVersion, MinorVersion minorVersion, Version version, Device device) {
        super(amsHeader);
        this.result = requireNonNull(result);
        this.majorVersion = requireNonNull(majorVersion);
        this.minorVersion = requireNonNull(minorVersion);
        this.version = requireNonNull(version);
        this.device = requireNonNull(device);
    }

    private ADSReadDeviceInfoResponse(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, Result result, MajorVersion majorVersion, MinorVersion minorVersion, Version version, Device device) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.result = requireNonNull(result);
        this.majorVersion = requireNonNull(majorVersion);
        this.minorVersion = requireNonNull(minorVersion);
        this.version = requireNonNull(version);
        this.device = requireNonNull(device);
    }

    public static ADSReadDeviceInfoResponse of(AMSHeader amsHeader, Result result, MajorVersion majorVersion, MinorVersion minorVersion, Version version, Device device) {
        return new ADSReadDeviceInfoResponse(amsHeader, result, majorVersion, minorVersion, version, device);
    }

    public static ADSReadDeviceInfoResponse of(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, Result result, MajorVersion majorVersion, MinorVersion minorVersion, Version version, Device device) {
        return new ADSReadDeviceInfoResponse(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, result, majorVersion, minorVersion, version, device);
    }

    public Result getResult() {
        return result;
    }

    public MajorVersion getMajorVersion() {
        return majorVersion;
    }

    public MinorVersion getMinorVersion() {
        return minorVersion;
    }

    public Version getVersion() {
        return version;
    }

    public Device getDevice() {
        return device;
    }

    @Override
    public ADSData getAdsData() {
        return buildADSData(result, majorVersion, minorVersion, version, device);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof ADSReadDeviceInfoResponse))
            return false;
        if (!super.equals(o))
            return false;

        ADSReadDeviceInfoResponse that = (ADSReadDeviceInfoResponse) o;

        if (!result.equals(that.result))
            return false;
        if (!majorVersion.equals(that.majorVersion))
            return false;
        if (!minorVersion.equals(that.minorVersion))
            return false;
        if (!version.equals(that.version))
            return false;

        return device.equals(that.device);
    }

    @Override
    public int hashCode() {
        int result1 = super.hashCode();
        result1 = 31 * result1 + result.hashCode();
        result1 = 31 * result1 + majorVersion.hashCode();
        result1 = 31 * result1 + minorVersion.hashCode();
        result1 = 31 * result1 + version.hashCode();
        result1 = 31 * result1 + device.hashCode();
        return result1;
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
