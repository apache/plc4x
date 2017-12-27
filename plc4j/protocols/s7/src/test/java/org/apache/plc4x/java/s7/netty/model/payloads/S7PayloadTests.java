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

import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportSize;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class S7PayloadTests {

    @Test
    @Tag("fast")
    void varPayloadItem() {
        DataTransportErrorCode returnCode = DataTransportErrorCode.NOT_FOUND;
        DataTransportSize dataTransportSize = DataTransportSize.INTEGER;
        byte[] data = {(byte)0xFF};

        VarPayloadItem  varPayloadItem = new VarPayloadItem(returnCode, dataTransportSize, data);
        assertTrue(varPayloadItem.getReturnCode() == DataTransportErrorCode.NOT_FOUND, "Unexpected data transport error code");
        assertTrue(varPayloadItem.getDataTransportSize() == DataTransportSize.INTEGER, "Unexpected data transport size");
        assertTrue(varPayloadItem.getData()[0] == (byte) 0xFF, "Unexpected user data");
    }

    @Test
    @Tag("fast")
    void varPayload() {
        ParameterType parameterType = ParameterType.DOWNLOAD_ENDED;
        ArrayList<VarPayloadItem> payloadItems = new ArrayList<>();
        byte[] data = {(byte)0xFF};
        
        payloadItems.add(new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BIT, data));

        VarPayload  varPayload = new VarPayload(parameterType, payloadItems);
        assertTrue(varPayload.getType() == ParameterType.DOWNLOAD_ENDED, "Unexpected parameter type");
        assertTrue(varPayload.getPayloadItems().containsAll(payloadItems) , "Unexpected pay load items");
    }

}