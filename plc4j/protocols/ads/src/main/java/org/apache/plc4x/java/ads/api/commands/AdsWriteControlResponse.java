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

import org.apache.plc4x.java.ads.api.commands.types.Result;
import org.apache.plc4x.java.ads.api.generic.AdsData;
import org.apache.plc4x.java.ads.api.generic.AmsHeader;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;

import static java.util.Objects.requireNonNull;

/**
 * Changes the ADS status and the device status of an ADS device.
 */
@AdsCommandType(Command.ADS_WRITE_CONTROL)
public class AdsWriteControlResponse extends AdsAbstractResponse {
    /**
     * 4 bytes	ADS error number
     */
    private final Result result;

    private AdsWriteControlResponse(AmsHeader amsHeader, Result result) {
        super(amsHeader);
        this.result = requireNonNull(result);
    }

    private AdsWriteControlResponse(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, Result result) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.result = requireNonNull(result);
    }

    public static AdsWriteControlResponse of(AmsHeader amsHeader, Result result) {
        return new AdsWriteControlResponse(amsHeader, result);
    }

    public static AdsWriteControlResponse of(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, Result result) {
        return new AdsWriteControlResponse(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, result);
    }

    public Result getResult() {
        return result;
    }

    @Override
    public AdsData getAdsData() {
        return buildADSData(result);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof AdsWriteControlResponse))
            return false;
        if (!super.equals(o))
            return false;

        AdsWriteControlResponse that = (AdsWriteControlResponse) o;

        return result.equals(that.result);
    }

    @Override
    public int hashCode() {
        int result1 = super.hashCode();
        result1 = 31 * result1 + result.hashCode();
        return result1;
    }

    @Override
    public String toString() {
        return "AdsWriteControlResponse{" +
            "result=" + result +
            "} " + super.toString();
    }
}
