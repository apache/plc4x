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
package org.apache.plc4x.java.ads.api.generic.calculated;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.generic.AMSHeader;
import org.apache.plc4x.java.ads.api.generic.types.*;

import static org.apache.plc4x.java.ads.api.util.ByteReadableUtils.buildByteBuff;

public class CalculatedAMSHeader extends AMSHeader {

    protected final CommandIdSupplier commandIdSupplier;

    protected final StateIdSupplier stateIdSupplier;

    protected final LengthSupplier lengthSupplier;

    protected CalculatedAMSHeader(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, CommandIdSupplier commandIdSupplier, StateIdSupplier stateIdSupplier, LengthSupplier lengthSupplier, Invoke invokeId) {
        super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, null, null, null, AMSError.NONE, invokeId);
        this.commandIdSupplier = commandIdSupplier;
        this.stateIdSupplier = stateIdSupplier;
        this.lengthSupplier = lengthSupplier;
    }

    public static CalculatedAMSHeader of(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, CommandIdSupplier commandIdSupplier, StateIdSupplier stateIdSupplier, LengthSupplier lengthSupplier, Invoke invokeId) {
        return new CalculatedAMSHeader(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, commandIdSupplier, stateIdSupplier, lengthSupplier, invokeId);
    }

    @FunctionalInterface
    public interface CommandIdSupplier {
        Command getCommandId();
    }

    @FunctionalInterface
    public interface StateIdSupplier {
        State getStateId();
    }

    public interface LengthSupplier {
        DataLength getLength();
    }

    @Override
    public ByteBuf getByteBuf() {
        return buildByteBuff(
            targetAmsNetId,
            targetAmsPort,
            sourceAmsNetId,
            sourceAmsPort,
            commandIdSupplier.getCommandId(),
            stateIdSupplier.getStateId(),
            lengthSupplier.getLength(),
            code,
            invokeId);
    }
}
