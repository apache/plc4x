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
package org.apache.plc4x.java.ads.api.generic;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.generic.types.*;
import org.apache.plc4x.java.ads.api.util.ByteReadable;
import org.apache.plc4x.java.ads.api.util.LengthSupplier;

import static java.util.Objects.requireNonNull;
import static org.apache.plc4x.java.ads.api.util.ByteReadableUtils.buildByteBuff;

/**
 * AMS Header	32 bytes	The AMS/TCP-Header contains the addresses of the transmitter and receiver. In addition the AMS error code , the ADS command Id and some other information.
 */
public class AMSHeader implements ByteReadable {
    /**
     * This is the AMSNetId of the station, for which the packet is intended. Remarks see below.
     */
    private final AMSNetId targetAmsNetId;
    /**
     * This is the AMSPort of the station, for which the packet is intended.
     */
    private final AMSPort targetAmsPort;
    /**
     * This contains the AMSNetId of the station, from which the packet was sent.
     */
    private final AMSNetId sourceAmsNetId;
    /**
     * This contains the AMSPort of the station, from which the packet was sent.
     */
    private final AMSPort sourceAmsPort;
    /**
     * 2 bytes.
     */
    private final Command commandId;
    /**
     * 2 bytes.
     */
    private final State stateFlags;
    /**
     * 4 bytes	Size of the data range. The unit is byte.
     */
    private final DataLength dataLength;

    /**
     * 4 bytes	AMS error number. See ADS Return Codes.
     */
    private final AMSError code;

    /**
     * 4 bytes	Free usable 32 bit array. Usually this array serves to send an Id. This Id makes is possible to assign a received response to a request, which was sent before.
     */
    private final Invoke invokeId;

    private final LengthSupplier dataLengthSupplier;

    private AMSHeader(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Command commandId, State stateFlags, DataLength dataLength, AMSError code, Invoke invokeId) {
        this.targetAmsNetId = requireNonNull(targetAmsNetId);
        this.targetAmsPort = requireNonNull(targetAmsPort);
        this.sourceAmsNetId = requireNonNull(sourceAmsNetId);
        this.sourceAmsPort = requireNonNull(sourceAmsPort);
        this.commandId = requireNonNull(commandId);
        this.stateFlags = requireNonNull(stateFlags);
        this.dataLength = requireNonNull(dataLength);
        this.code = requireNonNull(code);
        this.invokeId = requireNonNull(invokeId);
        dataLengthSupplier = null;

    }

    private AMSHeader(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Command commandId, State stateFlags, LengthSupplier dataLengthSupplier, AMSError code, Invoke invokeId) {
        this.targetAmsNetId = requireNonNull(targetAmsNetId);
        this.targetAmsPort = requireNonNull(targetAmsPort);
        this.sourceAmsNetId = requireNonNull(sourceAmsNetId);
        this.sourceAmsPort = requireNonNull(sourceAmsPort);
        this.commandId = requireNonNull(commandId);
        this.stateFlags = requireNonNull(stateFlags);
        this.dataLength = null;
        this.code = requireNonNull(code);
        this.invokeId = requireNonNull(invokeId);
        this.dataLengthSupplier = requireNonNull(dataLengthSupplier);
    }

    public static AMSHeader of(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Command commandId, State stateFlags, DataLength dataLength, AMSError code, Invoke invokeId) {
        return new AMSHeader(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, commandId, stateFlags, dataLength, code, invokeId);
    }

    public static AMSHeader of(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Command commandId, State stateFlags, LengthSupplier dataLengthSupplier, AMSError code, Invoke invokeId) {
        return new AMSHeader(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, commandId, stateFlags, dataLengthSupplier, code, invokeId);
    }

    @Override
    public ByteBuf getByteBuf() {
        return buildByteBuff(
            targetAmsNetId,
            targetAmsPort,
            sourceAmsNetId,
            sourceAmsPort,
            commandId,
            stateFlags,
            getDataLength(),
            code,
            invokeId);
    }

    public AMSNetId getTargetAmsNetId() {
        return targetAmsNetId;
    }

    public AMSPort getTargetAmsPort() {
        return targetAmsPort;
    }

    public AMSNetId getSourceAmsNetId() {
        return sourceAmsNetId;
    }

    public AMSPort getSourceAmsPort() {
        return sourceAmsPort;
    }

    public Command getCommandId() {
        return commandId;
    }

    public State getStateFlags() {
        return stateFlags;
    }

    public DataLength getDataLength() {
        return dataLengthSupplier == null ? dataLength : DataLength.of(dataLengthSupplier);
    }

    public AMSError getCode() {
        return code;
    }

    public Invoke getInvokeId() {
        return invokeId;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof AMSHeader))
            return false;

        AMSHeader amsHeader = (AMSHeader) o;

        if (!targetAmsNetId.equals(amsHeader.targetAmsNetId))
            return false;
        if (!targetAmsPort.equals(amsHeader.targetAmsPort))
            return false;
        if (!sourceAmsNetId.equals(amsHeader.sourceAmsNetId))
            return false;
        if (!sourceAmsPort.equals(amsHeader.sourceAmsPort))
            return false;
        if (commandId != amsHeader.commandId)
            return false;
        if (!stateFlags.equals(amsHeader.stateFlags))
            return false;
        if (!code.equals(amsHeader.code))
            return false;
        if (!invokeId.equals(amsHeader.invokeId))
            return false;
        
        return getDataLength().equals(((AMSHeader) o).getDataLength());
    }

    @Override
    public int hashCode() {
        int result = targetAmsNetId.hashCode();
        result = 31 * result + targetAmsPort.hashCode();
        result = 31 * result + sourceAmsNetId.hashCode();
        result = 31 * result + sourceAmsPort.hashCode();
        result = 31 * result + commandId.hashCode();
        result = 31 * result + stateFlags.hashCode();
        result = 31 * result + code.hashCode();
        result = 31 * result + invokeId.hashCode();
        result = 31 * result + getDataLength().hashCode();
        return result;
    }

    @Override
    public String toString() {
        return "AMSHeader{" +
            "targetAmsNetId=" + targetAmsNetId +
            ", targetAmsPort=" + targetAmsPort +
            ", sourceAmsNetId=" + sourceAmsNetId +
            ", sourceAmsPort=" + sourceAmsPort +
            ", commandId=" + commandId +
            ", stateFlags=" + stateFlags +
            ", dataLength=" + getDataLength() +
            ", code=" + code +
            ", invokeId=" + invokeId +
            '}';
    }
}
