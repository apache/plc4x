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
package org.apache.plc4x.java.s7.netty.model.payloads;

import org.apache.plc4x.java.s7.netty.model.params.items.ReadVarRequestItem;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;

public abstract class ReadVarPayload implements S7Payload {

    private final ReadVarRequestItem readVarRequestItem;
    private final DataTransportErrorCode dataTransportErrorCode;

    public ReadVarPayload(ReadVarRequestItem readVarRequestItem, DataTransportErrorCode dataTransportErrorCode) {
        this.readVarRequestItem = readVarRequestItem;
        this.dataTransportErrorCode = dataTransportErrorCode;
    }

    @Override
    public ParameterType getType() {
        return ParameterType.READ_VAR;
    }

    public ReadVarRequestItem getReadVarRequestItem() {
        return readVarRequestItem;
    }

    public DataTransportErrorCode getDataTransportErrorCode() {
        return dataTransportErrorCode;
    }

}
