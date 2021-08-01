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

import org.apache.plc4x.java.ads.api.commands.types.IndexGroup;
import org.apache.plc4x.java.ads.api.commands.types.IndexOffset;
import org.apache.plc4x.java.ads.api.commands.types.Length;
import org.apache.plc4x.java.ads.api.generic.AdsData;
import org.apache.plc4x.java.ads.api.generic.AmsHeader;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;

import static java.util.Objects.requireNonNull;

/**
 * With ADS Read data can be read from an ADS device.  The data are addressed by the Index Group and the Index Offset
 */
@AdsCommandType(Command.ADS_READ)
public class AdsReadRequest extends AdsAbstractRequest {

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

    private AdsReadRequest(AmsHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length) {
        super(amsHeader);
        this.indexGroup = requireNonNull(indexGroup);
        this.indexOffset = requireNonNull(indexOffset);
        this.length = requireNonNull(length);
    }

    private AdsReadRequest(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, IndexGroup indexGroup, IndexOffset indexOffset, Length length) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.indexGroup = requireNonNull(indexGroup);
        this.indexOffset = requireNonNull(indexOffset);
        this.length = requireNonNull(length);
    }

    public static AdsReadRequest of(AmsHeader amsHeader, IndexGroup indexGroup, IndexOffset indexOffset, Length length) {
        return new AdsReadRequest(amsHeader, indexGroup, indexOffset, length);
    }

    public static AdsReadRequest of(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, IndexGroup indexGroup, IndexOffset indexOffset, Length length) {
        return new AdsReadRequest(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, indexGroup, indexOffset, length);
    }

    public IndexGroup getIndexGroup() {
        return indexGroup;
    }

    public IndexOffset getIndexOffset() {
        return indexOffset;
    }

    public Length getLength() {
        return length;
    }

    @Override
    public AdsData getAdsData() {
        return buildADSData(indexGroup, indexOffset, length);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof AdsReadRequest))
            return false;
        if (!super.equals(o))
            return false;

        AdsReadRequest that = (AdsReadRequest) o;

        if (!indexGroup.equals(that.indexGroup))
            return false;
        if (!indexOffset.equals(that.indexOffset))
            return false;

        return length.equals(that.length);
    }

    @Override
    public int hashCode() {
        int result = super.hashCode();
        result = 31 * result + indexGroup.hashCode();
        result = 31 * result + indexOffset.hashCode();
        result = 31 * result + length.hashCode();
        return result;
    }

    @Override
    public String toString() {
        return "AdsReadRequest{" +
            "indexGroup=" + indexGroup +
            ", indexOffset=" + indexOffset +
            ", length=" + length +
            "} " + super.toString();
    }
}
