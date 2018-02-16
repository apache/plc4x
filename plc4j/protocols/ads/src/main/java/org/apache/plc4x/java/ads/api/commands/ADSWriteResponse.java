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

import org.apache.plc4x.java.ads.api.commands.types.Result;
import org.apache.plc4x.java.ads.api.generic.ADSData;
import org.apache.plc4x.java.ads.api.generic.AMSHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPHeader;
import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;

/**
 * With ADS Write data can be written to an ADS device.
 */
@ADSCommandType(Command.ADS_Write)
public class ADSWriteResponse extends ADSAbstractResponse {

    /**
     * 4 bytes	ADS error number
     */
    private final Result result;

    protected ADSWriteResponse(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, Result result) {
        super(amstcpHeader, amsHeader);
        this.result = result;
    }

    protected ADSWriteResponse(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, Result result) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.result = result;
    }

    public static ADSWriteResponse of(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, Result result) {
        return new ADSWriteResponse(amstcpHeader, amsHeader, result);
    }

    public static ADSWriteResponse of(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, Result result) {
        return new ADSWriteResponse(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, result);
    }

    @Override
    public ADSData getAdsData() {
        return buildADSData(result);
    }

    public Result getResult() {
        return result;
    }

    @Override
    public String toString() {
        return "ADSWriteResponse{" +
            "result=" + result +
            "} " + super.toString();
    }
}
