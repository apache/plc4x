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
import org.apache.plc4x.java.ads.api.commands.types.Length;
import org.apache.plc4x.java.ads.api.commands.types.Result;
import org.apache.plc4x.java.ads.api.generic.AdsData;
import org.apache.plc4x.java.ads.api.generic.AmsHeader;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Command;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.api.util.LengthSupplier;

import static java.util.Objects.requireNonNull;

/**
 * With ADS Read Write data will be written to an ADS device. Additionally, data can be read from the ADS device.
 */
@AdsCommandType(Command.ADS_READ_WRITE)
public class AdsReadWriteResponse extends AdsAbstractResponse {

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

    private final transient LengthSupplier lengthSupplier;

    private AdsReadWriteResponse(AmsHeader amsHeader, Result result, Length length, Data data) {
        super(amsHeader);
        this.result = requireNonNull(result);
        this.length = requireNonNull(length);
        this.data = requireNonNull(data);
        this.lengthSupplier = null;
    }

    private AdsReadWriteResponse(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, Result result, Data data) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId);
        this.result = requireNonNull(result);
        this.length = null;
        this.data = requireNonNull(data);
        this.lengthSupplier = data;
    }

    public static AdsReadWriteResponse of(AmsHeader amsHeader, Result result, Length length, Data data) {
        return new AdsReadWriteResponse(amsHeader, result, length, data);
    }

    public static AdsReadWriteResponse of(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, Invoke invokeId, Result result, Data data) {
        return new AdsReadWriteResponse(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId, result, data);
    }

    public Result getResult() {
        return result;
    }

    public Length getLength() {
        return lengthSupplier == null ? length : Length.of(lengthSupplier);
    }

    public Data getData() {
        return data;
    }

    @Override
    public AdsData getAdsData() {
        return buildADSData(result, getLength(), data);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof AdsReadWriteResponse))
            return false;
        if (!super.equals(o))
            return false;

        AdsReadWriteResponse that = (AdsReadWriteResponse) o;

        if (!result.equals(that.result))
            return false;
        if (!getLength().equals(that.getLength()))
            return false;
        return data.equals(that.data);
    }

    @Override
    public int hashCode() {
        int result1 = super.hashCode();
        result1 = 31 * result1 + result.hashCode();
        result1 = 31 * result1 + getLength().hashCode();
        result1 = 31 * result1 + data.hashCode();
        return result1;
    }

    @Override
    public String toString() {
        return "AdsReadWriteResponse{" +
            "result=" + result +
            ", length=" + getLength() +
            ", data=" + data +
            "} " + super.toString();
    }
}
