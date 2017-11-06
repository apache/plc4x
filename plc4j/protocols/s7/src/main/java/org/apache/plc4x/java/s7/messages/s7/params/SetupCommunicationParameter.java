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
package org.apache.plc4x.java.s7.messages.s7.params;

import org.apache.plc4x.java.exception.PlcException;
import org.apache.plc4x.java.exception.PlcIoException;
import org.apache.plc4x.java.s7.messages.s7.types.Function;

import java.io.DataOutputStream;
import java.io.IOException;

public class SetupCommunicationParameter implements S7Parameter {

    private final short maxAmqCaller;
    private final short maxAmqCallee;
    private final short pduLength;

    public SetupCommunicationParameter(short maxAmqCaller, short maxAmqCallee, short pduLength) {
        this.maxAmqCaller = maxAmqCaller;
        this.maxAmqCallee = maxAmqCallee;
        this.pduLength = pduLength;
    }

    public short getMaxAmqCaller() {
        return maxAmqCaller;
    }

    public short getMaxAmqCallee() {
        return maxAmqCallee;
    }

    public short getPduLength() {
        return pduLength;
    }

    @Override
    public short getLength() {
        return 8;
    }

    @Override
    public void serialize(DataOutputStream dos) throws PlcException {
        try {
            dos.writeByte(Function.SETUP_COMMUNICATION.getCode());
            dos.writeByte((byte) 0x00);         // Reserved ...
            dos.writeShort(maxAmqCaller);
            dos.writeShort(maxAmqCallee);
            dos.writeShort(pduLength);
        } catch (IOException e) {
            throw new PlcIoException("Error serializing message", e);
        }
    }

}
