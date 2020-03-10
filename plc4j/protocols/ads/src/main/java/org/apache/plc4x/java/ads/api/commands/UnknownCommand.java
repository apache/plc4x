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

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.generic.AdsData;
import org.apache.plc4x.java.ads.api.generic.AmsHeader;
import org.apache.plc4x.java.ads.api.generic.AmsPacket;
import org.apache.plc4x.java.ads.api.generic.types.*;

import static java.util.Objects.requireNonNull;

/**
 * Unknown ADS Package
 */
@AdsCommandType(Command.UNKNOWN)
public class UnknownCommand extends AmsPacket {

    private final transient ByteBuf remainingBytes;

    private UnknownCommand(AmsHeader amsHeader, ByteBuf remainingBytes) {
        super(amsHeader);
        this.remainingBytes = requireNonNull(remainingBytes);
    }

    private UnknownCommand(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, State stateId, Invoke invokeId, ByteBuf remainingBytes) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, stateId, invokeId);
        this.remainingBytes = requireNonNull(remainingBytes);
    }

    @Override
    public AdsData getAdsData() {
        return () -> remainingBytes;
    }

    public static UnknownCommand of(AmsHeader amsHeader, ByteBuf remainingBytes) {
        return new UnknownCommand(amsHeader, remainingBytes);
    }

    public static UnknownCommand of(AmsNetId targetAmsNetId, AmsPort targetAmsPort, AmsNetId sourceAmsNetId, AmsPort sourceAmsPort, State stateId, Invoke invokeId, ByteBuf remainingBytes) {
        return new UnknownCommand(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, stateId, invokeId, remainingBytes);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof UnknownCommand)) {
            return false;
        }
        if (!super.equals(o)) {
            return false;
        }

        UnknownCommand that = (UnknownCommand) o;

        return remainingBytes.equals(that.remainingBytes);
    }

    @Override
    public int hashCode() {
        int result = super.hashCode();
        result = 31 * result + remainingBytes.hashCode();
        return result;
    }

    @Override
    public String toString() {
        return "UnknownCommand";
    }
}
