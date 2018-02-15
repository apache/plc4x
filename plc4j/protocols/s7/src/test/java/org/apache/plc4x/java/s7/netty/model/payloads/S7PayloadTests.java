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

import static org.assertj.core.api.Assertions.assertThat;

import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportSize;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;
import org.apache.plc4x.test.FastTests;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import java.util.ArrayList;

public class S7PayloadTests {

    @Test
    @Category(FastTests.class)
    public void varPayloadItem() {
        DataTransportErrorCode returnCode = DataTransportErrorCode.NOT_FOUND;
        DataTransportSize dataTransportSize = DataTransportSize.INTEGER;
        byte[] data = {(byte)0xFF};

        VarPayloadItem  varPayloadItem = new VarPayloadItem(returnCode, dataTransportSize, data);
        assertThat(varPayloadItem.getReturnCode()).isEqualTo(DataTransportErrorCode.NOT_FOUND).withFailMessage("Unexpected data transport error code");
        assertThat(varPayloadItem.getDataTransportSize()).isEqualTo(DataTransportSize.INTEGER).withFailMessage("Unexpected data transport size");
        assertThat(varPayloadItem.getData()[0]).isEqualTo((byte) 0xFF).withFailMessage("Unexpected user data");
    }

    @Test
    @Category(FastTests.class)
    public void varPayload() {
        ParameterType parameterType = ParameterType.DOWNLOAD_ENDED;
        ArrayList<VarPayloadItem> payloadItems = new ArrayList<>();
        byte[] data = {(byte)0xFF};
        
        payloadItems.add(new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BIT, data));

        VarPayload  varPayload = new VarPayload(parameterType, payloadItems);
        assertThat(varPayload.getType()).isEqualTo(ParameterType.DOWNLOAD_ENDED).withFailMessage("Unexpected parameter type");
        assertThat(varPayload.getPayloadItems()).containsAll(payloadItems).withFailMessage("Unexpected pay load items");
    }

}