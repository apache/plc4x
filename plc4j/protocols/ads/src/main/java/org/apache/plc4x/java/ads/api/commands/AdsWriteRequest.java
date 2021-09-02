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

import org.apache.plc4x.java.ads.api.commands.types.Data;
import org.apache.plc4x.java.ads.api.commands.types.IndexGroup;
import org.apache.plc4x.java.ads.api.commands.types.IndexOffset;
import org.apache.plc4x.java.ads.api.commands.types.Length;
import org.apache.plc4x.java.ads.api.generic.AdsData;
import org.apache.plc4x.java.ads.api.generic.AmsHeader;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.api.util.LengthSupplier;

import static java.util.Objects.requireNonNull;

/**
 * With ADS Write data can be written to an ADS device. The data are addressed by the Index Group and the Index Offset.
 */
@AdsCommandType(Command.ADS_WRITE)
public class AdsWriteRequest extends AdsAbstractRequest {

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

    private final transient LengthSupplier lengthSupplier;

    private AdsWriteRequest(AmsHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length, Data data) {
        super(amsHeader);
        this.indexGroup = requireNonNull(indexGroup);
        this.indexOffset = requireNonNull(indexOffset);
        this.length = requireNonNull(length);
        this.data = requireNonNull(data);
        this.lengthSupplier = null;
    }

    private AdsWriteRequest(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, IndexGroup indexGroup, IndexOffset indexOffset, Data data) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.indexGroup = requireNonNull(indexGroup);
        this.indexOffset = requireNonNull(indexOffset);
        this.length = null;
        this.data = requireNonNull(data);
        this.lengthSupplier = data;
    }

    public static AdsWriteRequest of(AmsHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length, Data data) {
        return new AdsWriteRequest(amsHeader, indexGroup, indexOffset, length, data);
    }

    public static AdsWriteRequest of(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, IndexGroup indexGroup, IndexOffset indexOffset, Data data) {
        return new AdsWriteRequest(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, indexGroup, indexOffset, data);
    }

    public IndexGroup getIndexGroup() {
        return indexGroup;
    }

    public IndexOffset getIndexOffset() {
        return indexOffset;
    }

    public Length getLength() {
        return lengthSupplier == null ? length : Length.of(lengthSupplier);
    }

    public Data getData() {
        return data;
    }

    @Override
    public AdsData getAdsData() {
        return buildADSData(indexGroup, indexOffset, getLength(), data);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof AdsWriteRequest))
            return false;
        if (!super.equals(o))
            return false;

        AdsWriteRequest that = (AdsWriteRequest) o;

        if (!indexGroup.equals(that.indexGroup))
            return false;
        if (!indexOffset.equals(that.indexOffset))
            return false;
        if (!getLength().equals(that.getLength()))
            return false;

        return data.equals(that.data);
    }

    @Override
    public int hashCode() {
        int result = super.hashCode();
        result = 31 * result + indexGroup.hashCode();
        result = 31 * result + indexOffset.hashCode();
        result = 31 * result + getLength().hashCode();
        result = 31 * result + data.hashCode();
        return result;
    }

    @Override
    public String toString() {
        return "AdsWriteRequest{" +
            "indexGroup=" + indexGroup +
            ", indexOffset=" + indexOffset +
            ", length=" + getLength() +
            ", data=" + data +
            "} " + super.toString();
    }

}
