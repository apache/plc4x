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

import org.apache.plc4x.java.ads.api.commands.types.IndexGroup;
import org.apache.plc4x.java.ads.api.commands.types.IndexOffset;
import org.apache.plc4x.java.ads.api.commands.types.Length;
import org.apache.plc4x.java.ads.api.generic.ADSData;
import org.apache.plc4x.java.ads.api.generic.AMSHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPPaket;
import org.apache.plc4x.java.ads.api.generic.types.*;

/**
 * With ADS Read data can be read from an ADS device.  The data are addressed by the Index Group and the Index Offset
 */
public class ADSReadRequest extends AMSTCPPaket {

    /**
     * 4 bytes	Index Group of the data which should be read.
     */
    private final IndexGroup indexGroup;

    /**
     * 4 bytes	Index Offset of the data which should be read.
     */
    private final IndexOffset indexOffset;

    /**
     * 4 bytes	Length of the data (in bytes) which should be read.
     */
    private final Length length;

    public ADSReadRequest(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length) {
        super(amstcpHeader, amsHeader);
        this.indexGroup = indexGroup;
        this.indexOffset = indexOffset;
        this.length = length;
    }

    public ADSReadRequest(AMSHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length) {
        super(amsHeader);
        this.indexGroup = indexGroup;
        this.indexOffset = indexOffset;
        this.length = length;
    }

    public ADSReadRequest(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, IndexGroup indexGroup, IndexOffset indexOffset, Length length) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.indexGroup = indexGroup;
        this.indexOffset = indexOffset;
        this.length = length;
    }

    @Override
    public ADSData getAdsData() {
        return buildADSData(indexGroup, indexOffset, length);
    }

    @Override
    protected Command getCommandId() {
        return Command.ADS_Read;
    }

    @Override
    protected State getStateId() {
        return State.ADS_REQUEST_TCP;
    }
}
