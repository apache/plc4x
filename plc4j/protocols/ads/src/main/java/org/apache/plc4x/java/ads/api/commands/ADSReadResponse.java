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

import org.apache.plc4x.java.ads.api.commands.types.Data;
import org.apache.plc4x.java.ads.api.commands.types.Length;
import org.apache.plc4x.java.ads.api.commands.types.Result;
import org.apache.plc4x.java.ads.api.generic.ADSData;
import org.apache.plc4x.java.ads.api.generic.AMSHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPHeader;
import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.api.util.LengthSupplier;

import static java.util.Objects.requireNonNull;

/**
 * With ADS Read data can be read from an ADS device
 */
@ADSCommandType(Command.ADS_Read)
public class ADSReadResponse extends ADSAbstractResponse {

    /**
     * 4 bytes	ADS error number
     */
    private final Result result;
    /**
     * 4 bytes	Length of data which are supplied back.
     */
    private final Length length;
    /**
     * n bytes	Data which are supplied back.
     */
    private final Data data;

    ////
    // Used when fields should be calculated. TODO: check if we better work with a subclass.
    private final LengthSupplier lengthSupplier;
    private final boolean calculated;
    //
    ///

    private ADSReadResponse(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, Result result, Length length, Data data) {
        super(amstcpHeader, amsHeader);
        this.result = requireNonNull(result);
        this.length = requireNonNull(length);
        this.data = requireNonNull(data);
        this.lengthSupplier = null;
        this.calculated = false;
    }

    private ADSReadResponse(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, Result result, Data data) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.result = requireNonNull(result);
        this.length = null;
        this.data = requireNonNull(data);
        this.lengthSupplier = data;
        this.calculated = true;
    }

    public static ADSReadResponse of(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, Result result, Length length, Data data) {
        return new ADSReadResponse(amstcpHeader, amsHeader, result, length, data);
    }

    public static ADSReadResponse of(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, Result result, Data data) {
        return new ADSReadResponse(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, result, data);
    }

    @Override
    public ADSData getAdsData() {
        return buildADSData(result, (calculated ? Length.of(lengthSupplier.getCalculatedLength()) : length), data);
    }

    public Result getResult() {
        return result;
    }

    public Length getLength() {
        return (calculated ? Length.of(lengthSupplier.getCalculatedLength()) : length);
    }

    public Data getData() {
        return data;
    }

    @Override
    public String toString() {
        return "ADSReadResponse{" +
            "result=" + result +
            ", length=" + (calculated ? Length.of(lengthSupplier.getCalculatedLength()) : length) +
            ", data=" + data +
            "} " + super.toString();
    }
}
