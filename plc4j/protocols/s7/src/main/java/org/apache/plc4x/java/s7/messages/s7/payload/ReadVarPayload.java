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
package org.apache.plc4x.java.s7.messages.s7.payload;

import org.apache.plc4x.java.exceptions.PlcException;
import org.apache.plc4x.java.s7.messages.s7.params.ReadVarParameter;
import org.apache.plc4x.java.s7.messages.s7.types.DataTransportSize;

import java.io.DataOutputStream;

public class ReadVarPayload implements S7Payload {

    private final ReadVarParameter parameter;
    private final DataTransportSize dataTransportSize;
    private final byte[] data;

    public ReadVarPayload(ReadVarParameter parameter, DataTransportSize dataTransportSize, byte[] data) {
        this.parameter = parameter;
        this.dataTransportSize = dataTransportSize;
        this.data = data;
    }

    public ReadVarParameter getParameter() {
        return parameter;
    }

    public DataTransportSize getDataTransportSize() {
        return dataTransportSize;
    }

    public byte[] getData() {
        return data;
    }

    @Override
    public void serialize(DataOutputStream dos) throws PlcException {

    }

    @Override
    public short getLength() {
        return 0;
    }

}
