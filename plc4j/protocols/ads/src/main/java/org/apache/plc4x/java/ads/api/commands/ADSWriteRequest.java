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
import org.apache.plc4x.java.ads.api.util.LengthSupplier;

import static java.util.Objects.requireNonNull;

/**
 * With ADS Write data can be written to an ADS device. The data are addressed by the Index Group and the Index Offset.
 */
@ADSCommandType(Command.ADS_Write)
public class ADSWriteRequest extends ADSAbstractRequest {

    /**
     * 4 bytes	Index Group in which the data should be written
     */
    private final IndexGroup indexGroup;
    /**
     * 4 bytes	Index Offset, in which the data should be written
     */
    private final IndexOffset indexOffset;
    /**
     * 4 bytes	Length of data in bytes which are written
     */
    private final Length length;
    /**
     * n bytes	Data which are written in the ADS device.
     */
    private final Data data;

    ////
    // Used when fields should be calculated. TODO: check if we better work with a subclass.
    private final LengthSupplier lengthSupplier;
    private final boolean calculated;
    //
    ///

    private ADSWriteRequest(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length, Data data) {
        super(amstcpHeader, amsHeader);
        this.indexGroup = requireNonNull(indexGroup);
        this.indexOffset = requireNonNull(indexOffset);
        this.length = requireNonNull(length);
        this.data = requireNonNull(data);
        this.lengthSupplier = null;
        this.calculated = false;
    }

    private ADSWriteRequest(AMSHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length, Data data) {
        super(amsHeader);
        this.indexGroup = requireNonNull(indexGroup);
        this.indexOffset = requireNonNull(indexOffset);
        this.length = requireNonNull(length);
        this.data = requireNonNull(data);
        this.lengthSupplier = null;
        this.calculated = false;
    }

    private ADSWriteRequest(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, IndexGroup indexGroup, IndexOffset indexOffset, Data data) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.indexGroup = requireNonNull(indexGroup);
        this.indexOffset = requireNonNull(indexOffset);
        this.length = null;
        this.data = requireNonNull(data);
        this.lengthSupplier = data;
        this.calculated = true;
    }

    public static ADSWriteRequest of(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length, Data data) {
        return new ADSWriteRequest(amstcpHeader, amsHeader, indexGroup, indexOffset, length, data);
    }

    public static ADSWriteRequest of(AMSHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length, Data data) {
        return new ADSWriteRequest(amsHeader, indexGroup, indexOffset, length, data);
    }

    public static ADSWriteRequest of(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, IndexGroup indexGroup, IndexOffset indexOffset, Data data) {
        return new ADSWriteRequest(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, indexGroup, indexOffset, data);
    }

    @Override
    public ADSData getAdsData() {
        return buildADSData(indexGroup, indexOffset, (calculated ? WriteLength.of(lengthSupplier.getCalculatedLength()) : length), data);
    }

    @Override
    public String toString() {
        return "ADSWriteRequest{" +
            "indexGroup=" + indexGroup +
            ", indexOffset=" + indexOffset +
            ", length=" + (calculated ? WriteLength.of(lengthSupplier.getCalculatedLength()) : length) +
            ", data=" + data +
            "} " + super.toString();
    }

}
