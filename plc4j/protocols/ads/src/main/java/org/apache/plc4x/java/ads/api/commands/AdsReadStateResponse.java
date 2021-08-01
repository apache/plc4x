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

import org.apache.plc4x.java.ads.api.commands.types.AdsState;
import org.apache.plc4x.java.ads.api.commands.types.DeviceState;
import org.apache.plc4x.java.ads.api.commands.types.Result;
import org.apache.plc4x.java.ads.api.generic.AdsData;
import org.apache.plc4x.java.ads.api.generic.AmsHeader;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;

import static java.util.Objects.requireNonNull;

/**
 * Reads the ADS status and the device status of an ADS device.
 */
@AdsCommandType(Command.ADS_READ_STATE)
public class AdsReadStateResponse extends AdsAbstractResponse {

    /**
     * 4 bytes	ADS error number
     */
    private final Result result;

    /**
     * 2 bytes	New ADS status (see data type ADSSTATE of the ADS-DLL).
     */
    private final AdsState adsState;

    /**
     * 2 bytes	New device status.
     */
    private final DeviceState deviceState;

    private AdsReadStateResponse(AmsHeader amsHeader, Result result, AdsState adsState, DeviceState deviceState) {
        super(amsHeader);
        this.result = requireNonNull(result);
        this.adsState = requireNonNull(adsState);
        this.deviceState = requireNonNull(deviceState);
    }

    private AdsReadStateResponse(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, Result result, AdsState adsState, DeviceState deviceState) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.result = requireNonNull(result);
        this.adsState = requireNonNull(adsState);
        this.deviceState = requireNonNull(deviceState);
    }

    public static AdsReadStateResponse of(AmsHeader amsHeader, Result result, AdsState adsState, DeviceState deviceState) {
        return new AdsReadStateResponse(amsHeader, result, adsState, deviceState);
    }

    public static AdsReadStateResponse of(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, Result result, AdsState adsState, DeviceState deviceState) {
        return new AdsReadStateResponse(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, result, adsState, deviceState);
    }

    public Result getResult() {
        return result;
    }

    public AdsState getAdsState() {
        return adsState;
    }

    public DeviceState getDeviceState() {
        return deviceState;
    }

    @Override
    public AdsData getAdsData() {
        return buildADSData(result, adsState, deviceState);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof AdsReadStateResponse)) return false;
        if (!super.equals(o)) return false;

        AdsReadStateResponse that = (AdsReadStateResponse) o;

        if (!result.equals(that.result)) return false;
        if (!adsState.equals(that.adsState)) return false;
        return deviceState.equals(that.deviceState);
    }

    @Override
    public int hashCode() {
        int result1 = super.hashCode();
        result1 = 31 * result1 + result.hashCode();
        result1 = 31 * result1 + adsState.hashCode();
        result1 = 31 * result1 + deviceState.hashCode();
        return result1;
    }

    @Override
    public String toString() {
        return "AdsReadStateResponse{" +
            "result=" + result +
            ", adsState=" + adsState +
            ", deviceState=" + deviceState +
            "} " + super.toString();
    }
}
