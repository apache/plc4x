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
import org.apache.plc4x.java.ads.api.commands.types.ReadLength;
import org.apache.plc4x.java.ads.api.commands.types.WriteLength;
import org.apache.plc4x.java.ads.api.generic.ADSData;
import org.apache.plc4x.java.ads.api.generic.AMSHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPHeader;
import org.apache.plc4x.java.ads.api.generic.AMSTCPPaket;
import org.apache.plc4x.java.ads.api.generic.types.*;
import org.apache.plc4x.java.ads.api.util.ByteValue;

/**
 * With ADS Read Write data will be written to an ADS device. Additionally, data can be read from the ADS device.
 * <p>
 * The data which can be read are addressed by the Index Group and the Index Offset
 */
public class ADSReadWriteRequest extends AMSTCPPaket {

    /**
     * 4 bytes	Index Group, in which the data should be written.
     */
    private final IndexGroup indexGroup;
    /**
     * 4 bytes	Index Offset, in which the data should be written
     */
    private final IndexOffset indexOffset;
    /**
     * 4 bytes	Length of data in bytes, which should be read.
     */
    private final ReadLength readLength;
    /**
     * 4 bytes	Length of data in bytes, which should be written
     */
    private final WriteLength writeLength;
    /**
     * n bytes	Data which are written in the ADS device.
     */
    private final Data data;

    public ADSReadWriteRequest(AMSTCPHeader amstcpHeader, AMSHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, ReadLength readLength, WriteLength writeLength, Data data) {
        super(amstcpHeader, amsHeader);
        this.indexGroup = indexGroup;
        this.indexOffset = indexOffset;
        this.readLength = readLength;
        this.writeLength = writeLength;
        this.data = data;
    }

    public ADSReadWriteRequest(AMSHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, ReadLength readLength, WriteLength writeLength, Data data) {
        super(amsHeader);
        this.indexGroup = indexGroup;
        this.indexOffset = indexOffset;
        this.readLength = readLength;
        this.writeLength = writeLength;
        this.data = data;
    }

    public ADSReadWriteRequest(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId, org.apache.plc4x.java.ads.api.generic.types.Data nData, IndexGroup indexGroup, IndexOffset indexOffset, ReadLength readLength, WriteLength writeLength, Data data) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, nData);
        this.indexGroup = indexGroup;
        this.indexOffset = indexOffset;
        this.readLength = readLength;
        this.writeLength = writeLength;
        this.data = data;
    }

    @Override
    public ADSData getAdsData() {
        return ADSData.EMPTY;
    }

    public static class Data extends ByteValue {
        public Data(byte... value) {
            super(value);
        }
    }

    @Override
    public Command getCommandId() {
        return Command.ADS_Read_Write;
    }

    @Override
    public State getStateId() {
        return State.ADS_REQUEST_TCP;
    }
}
