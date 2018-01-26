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

import org.apache.plc4x.java.ads.model.generic.ADSData;
import org.apache.plc4x.java.ads.model.generic.AMSHeader;
import org.apache.plc4x.java.ads.model.generic.AMSTCPHeader;
import org.apache.plc4x.java.ads.model.generic.AMSTCPPaket;
import org.apache.plc4x.java.ads.model.util.ByteValue;

/**
 * With ADS Write data can be written to an ADS device.
 */
public class ADSWriteResponse extends AMSTCPPaket {

    /**
     * 4 bytes	ADS error number
     */
    private final ADSReadResponse.Result result;

    public ADSWriteResponse(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, ADSReadResponse.Result result) {
        super(amstcpHeader, amsHeader);
        this.result = result;
    }

    @Override
    public ADSData getAdsData() {
        return buildADSData(result);
    }

    public static class Result extends ByteValue {
        public Result(byte... value) {
            super(value);
            assertLength(4);
        }
    }

}
