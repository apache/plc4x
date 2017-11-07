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

import org.apache.plc4x.java.exceptions.PlcException;
import org.apache.plc4x.java.exceptions.PlcIoException;
import org.apache.plc4x.java.exceptions.PlcProtocolException;
import org.apache.plc4x.java.s7.messages.s7.types.Function;

import java.io.DataInputStream;
import java.io.IOException;

public class S7ParameterFactory {

    public static S7Parameter parse(DataInputStream dis, boolean isResponse) throws PlcException {
        try {
            Function function = Function.valueOf(dis.readByte());
            if(function == null) {
                throw new PlcProtocolException("Could not find function");
            }
            switch (function) {
                case SETUP_COMMUNICATION:
                    dis.readByte(); // Reserved ...
                    short callingMaxAmq = dis.readShort();
                    short calledMaxAmq = dis.readShort();
                    short pduLength = dis.readShort();
                    return new SetupCommunicationParameter(callingMaxAmq, calledMaxAmq, pduLength);
                case READ_VAR:
                    byte numItems = dis.readByte();
                    if(isResponse) {
                        return new ReadVarParameter(numItems);
                    } /*else {
                        SpecificationType specificationType = SpecificationType.valueOf(dis.readByte());
                        byte addressSpecificationLenght = dis.readByte();
                        VariableAddressingMode variableAddressingMode = VariableAddressingMode.valueOf(dis.readByte());
                        TransportSize transportSize = TransportSize.valueOf(dis.readByte());
                        short length = dis.readShort();
                        short dbNumber = dis.readShort();
                        MemoryArea memoryArea = MemoryArea.valueOf(dis.readByte());
                        short byteAddress = dis.readShort(); // TODO: Is not 100% correct
                        byte bitAddress = dis.readByte(); // TODO: Is not 100% correct
                    }*/
                default:
                    throw new PlcProtocolException("Parameter not implemented yet " + function.name());
            }
        } catch (IOException e) {
            throw new PlcIoException("Error parsing message", e);
        }
    }

}
