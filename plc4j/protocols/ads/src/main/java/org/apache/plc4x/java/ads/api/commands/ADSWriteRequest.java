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
import org.apache.plc4x.java.ads.api.commands.types.IndexGroup;
import org.apache.plc4x.java.ads.api.commands.types.IndexOffset;
import org.apache.plc4x.java.ads.api.commands.types.Length;
import org.apache.plc4x.java.ads.api.generic.ADSData;
import org.apache.plc4x.java.ads.api.generic.AMSHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPPaket;
import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;

/**
 * With ADS Write data can be written to an ADS device. The data are addressed by the Index Group and the Index Offset.
 */
@ADSCommandType(Command.ADS_Write)
public class ADSWriteRequest extends AMSTCPPaket {

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

    public ADSWriteRequest(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length, Data data) {
        super(amstcpHeader, amsHeader);
        this.indexGroup = indexGroup;
        this.indexOffset = indexOffset;
        this.length = length;
        this.data = data;
    }

    public ADSWriteRequest(AMSHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length, Data data) {
        super(amsHeader);
        this.indexGroup = indexGroup;
        this.indexOffset = indexOffset;
        this.length = length;
        this.data = data;
    }

    public ADSWriteRequest(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, IndexGroup indexGroup, IndexOffset indexOffset, Length length, Data data) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.indexGroup = indexGroup;
        this.indexOffset = indexOffset;
        this.length = length;
        this.data = data;
    }

    @Override
    public ADSData getAdsData() {
        return buildADSData(indexGroup, indexOffset, length, data);
    }
}
