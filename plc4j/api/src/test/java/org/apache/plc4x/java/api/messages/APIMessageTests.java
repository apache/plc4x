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

package org.apache.plc4x.java.api.messages;

import static org.assertj.core.api.Assertions.assertThat;

import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.api.messages.items.ReadResponseItem;
import org.apache.plc4x.java.api.messages.items.WriteRequestItem;
import org.apache.plc4x.java.api.messages.items.WriteResponseItem;
import org.apache.plc4x.java.api.messages.mock.MockAddress;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcWriteRequest;
import org.apache.plc4x.java.api.types.ResponseCode;
import org.apache.plc4x.test.FastTests;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import java.util.*;


public class APIMessageTests {

    @Test
    @Category(FastTests.class)
    public void readRequestItemSize() {
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, address, 1);
        assertThat(readRequestItem.getAddress()).isEqualTo(address).withFailMessage("Unexpected address");
        assertThat(readRequestItem.getDatatype()).isEqualTo(Byte.class).withFailMessage("Unexpected data type");
        assertThat(readRequestItem.getSize()).isEqualTo(1).withFailMessage("Unexpected size");
    }

    @Test
    @Category(FastTests.class)
    public void readRequestItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, address);
        assertThat(readRequestItem.getAddress()).isEqualTo(address).withFailMessage("Unexpected address");
        assertThat(readRequestItem.getDatatype()).isEqualTo(Byte.class).withFailMessage("Unexpected data type");
        assertThat(readRequestItem.getSize()).isEqualTo(1).withFailMessage("Unexpected size");
    }

    @Test
    @Category(FastTests.class)
    public void readResponseItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, address, 1);
        ReadResponseItem<Byte> readResponseItem = new ReadResponseItem<>(readRequestItem, ResponseCode.OK, Collections.emptyList());
        assertThat(readResponseItem.getResponseCode()).isEqualTo(ResponseCode.OK).withFailMessage("Unexpected response code");
        assertThat(readResponseItem.getValues()).isEmpty();
        assertThat(readResponseItem.getRequestItem()).isEqualTo(readRequestItem).withFailMessage("Unexpected read request item");
    }

    @Test
    @Category(FastTests.class)
    public void writeRequestItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        WriteRequestItem<Byte> writeRequestItem = new WriteRequestItem<>(Byte.class, address, (byte) 0x45);

        assertThat(writeRequestItem.getAddress()).isEqualTo(address).withFailMessage("Unexpected address");
        assertThat(writeRequestItem.getDatatype()).isEqualTo(Byte.class).withFailMessage("Unexpected data type");
        assertThat(writeRequestItem.getValues().get(0)).isEqualTo((byte) 0x45).withFailMessage("Unexpected value");
    }

    @Test
    @Category(FastTests.class)
    public void writeRequestItems() {
        MockAddress address = new MockAddress("mock:/DATA");
        Byte data[] = {(byte) 0x23, (byte) 0x84};
        WriteRequestItem<Byte> writeRequestItem = new WriteRequestItem<>(Byte.class, address, data);
        assertThat(writeRequestItem.getAddress()).isEqualTo(address).withFailMessage("Unexpected address");
        assertThat(writeRequestItem.getDatatype()).isEqualTo(Byte.class).withFailMessage("Unexpected data type");
        assertThat(writeRequestItem.getValues().get(0)).isEqualTo((byte) 0x23).withFailMessage("Unexpected value");
        assertThat(writeRequestItem.getValues().get(1)).isEqualTo((byte) 0x84).withFailMessage("Unexpected value");
    }

    @Test
    @Category(FastTests.class)
    public void writeResponseItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        WriteRequestItem<Byte> writeRequestItem = new WriteRequestItem<>(Byte.class, address, (byte) 0x3B);
        WriteResponseItem<Byte> writeResponseItem = new WriteResponseItem<>(writeRequestItem, ResponseCode.OK);
        assertThat(writeResponseItem.getResponseCode()).isEqualTo(ResponseCode.OK).withFailMessage("Unexpected response code");
        assertThat(writeResponseItem.getRequestItem()).isEqualTo(writeRequestItem).withFailMessage("Unexpected response item");
    }

    @Test
    @Category(FastTests.class)
    public void plcReadRequestEmpty() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        assertThat(plcReadRequest.getRequestItems()).isEmpty();
        assertThat(plcReadRequest.getNumberOfItems()).isEqualTo(0).withFailMessage("Expected request items to be zero");
    }

    @Test
    @Category(FastTests.class)
    public void plcReadRequestAddress() {
        MockAddress address = new MockAddress("mock:/DATA");
        PlcReadRequest plcReadRequest = new TypeSafePlcReadRequest<>(Byte.class, address);
        assertThat(plcReadRequest.getRequestItems()).hasSize(1).withFailMessage("Expected one request item");
        assertThat(plcReadRequest.getNumberOfItems()).isEqualTo(1).withFailMessage("Expected one request item");
    }

    @Test
    @Category(FastTests.class)
    public void plcReadRequestSize() {
        MockAddress address = new MockAddress("mock:/DATA");
        PlcReadRequest plcReadRequest = PlcReadRequest.builder().addItem(Byte.class, address, (byte) 1).build(Byte.class);
        assertThat(plcReadRequest.getRequestItems()).hasSize(1).withFailMessage("Expected one request item");
        assertThat(plcReadRequest.getNumberOfItems()).isEqualTo(1).withFailMessage("Expected one request item");
    }

    @Test
    @Category(FastTests.class)
    public void plcReadRequestAddItem() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        assertThat(plcReadRequest.getRequestItems()).isEmpty();
        assertThat(plcReadRequest.getNumberOfItems()).isEqualTo(0).withFailMessage("Expected request items to be zero");
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, address, (byte) 1);
        plcReadRequest.addItem(readRequestItem);
        assertThat(plcReadRequest.getRequestItems()).hasSize(1).withFailMessage("Expected one request item");
        assertThat(plcReadRequest.getNumberOfItems()).isEqualTo(1).withFailMessage("Expected one request item");
    }

    @Test
    @Category(FastTests.class)
    public void plcReadResponse() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        List<ReadResponseItem<?>> responseItems = new ArrayList<>();
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, address, 1);
        ReadResponseItem<Byte> readResponseItem = new ReadResponseItem<>(readRequestItem, ResponseCode.OK, Collections.emptyList());
        responseItems.add(readResponseItem);
        PlcReadResponse plcReadResponse = new PlcReadResponse(plcReadRequest, responseItems);
        assertThat(plcReadResponse.getRequest().getNumberOfItems()).isEqualTo(0).withFailMessage("Unexpected number of response items");
        assertThat(plcReadResponse.getRequest()).isEqualTo(plcReadRequest).withFailMessage("Unexpected read request");
        assertThat(plcReadResponse.getResponseItems()).hasSize(1).withFailMessage("Unexpected number of response items");
        assertThat(plcReadResponse.getResponseItems()).containsAll(responseItems).withFailMessage("Unexpected items in response items");
    }

    @Test
    @Category(FastTests.class)
    public void plcWriteRequestEmpty() {
        PlcWriteRequest plcWriteRequest = new PlcWriteRequest();
        assertThat(plcWriteRequest.getRequestItems()).isEmpty();
        assertThat(plcWriteRequest.getNumberOfItems()).isEqualTo(0).withFailMessage("Expected request items to be zero");
    }

    @Test
    @Category(FastTests.class)
    public void plcWriteRequestObject() {
        MockAddress address = new MockAddress("mock:/DATA");
        PlcWriteRequest plcWriteRequest = new TypeSafePlcWriteRequest<>(Byte.class, address, (byte) 0x33);
        assertThat(plcWriteRequest.getRequestItems()).hasSize(1).withFailMessage("Expected no request item");
        assertThat(plcWriteRequest.getNumberOfItems()).isEqualTo(1).withFailMessage("Expected one request item");
        List values = plcWriteRequest.getRequestItems().get(0).getValues();
        assertThat((byte) values.get(0)).isEqualTo((byte) 0x33).withFailMessage("Expected value 0x33");
    }

    @Test
    @Category(FastTests.class)
    public void plcWriteRequestObjects() {
        MockAddress address = new MockAddress("mock:/DATA");
        Byte[] data = {(byte) 0x22, (byte) 0x66};
        PlcWriteRequest plcWriteRequest = new TypeSafePlcWriteRequest<>(Byte.class, address, data);
        assertThat(plcWriteRequest.getRequestItems()).hasSize(1).withFailMessage("Expected one request item");
        assertThat(plcWriteRequest.getNumberOfItems()).isEqualTo(1).withFailMessage("Expected one request item");
        List values = plcWriteRequest.getRequestItems().get(0).getValues();
        assertThat((Byte) values.get(0)).isEqualTo((byte) 0x22).withFailMessage("Expected value 0x22");
        assertThat((Byte) values.get(1)).isEqualTo((byte) 0x66).withFailMessage("Expected value 0x66");
    }

    @Test
    @Category(FastTests.class)
    public void plcWriteResponse() {
        PlcWriteRequest plcWriteRequest = new PlcWriteRequest();
        List<WriteResponseItem<?>> responseItems = new ArrayList<>();
        MockAddress address = new MockAddress("mock:/DATA");
        WriteRequestItem<Byte> writeRequestItem = new WriteRequestItem<>(Byte.class, address, (byte) 1);
        WriteResponseItem<Byte> writeResponseItem = new WriteResponseItem<>(writeRequestItem, ResponseCode.OK);
        responseItems.add(writeResponseItem);
        PlcWriteResponse plcReadResponse = new PlcWriteResponse(plcWriteRequest, responseItems);
        assertThat(plcReadResponse.getRequest().getNumberOfItems()).isEqualTo(0).withFailMessage("Unexpected number of response items");
        assertThat(plcReadResponse.getRequest()).isEqualTo(plcWriteRequest).withFailMessage("Unexpected read request");
        assertThat(plcReadResponse.getResponseItems()).hasSize(1).withFailMessage("Unexpected number of response items");
        assertThat(plcReadResponse.getResponseItems()).containsAll(responseItems).withFailMessage("Unexpected items in response items");
    }

    @Test
    @Category(FastTests.class)
    public void bulkPlcWriteResponseGetValue() {
        PlcWriteRequest plcWriteRequest = new PlcWriteRequest();
        List<WriteResponseItem<?>> responseItems = new ArrayList<>();
        MockAddress address = new MockAddress("mock:/DATA");
        WriteRequestItem<Byte> writeRequestItem1 = new WriteRequestItem<>(Byte.class, address, (byte) 1);
        WriteRequestItem<Byte> writeRequestItem2 = new WriteRequestItem<>(Byte.class, address, (byte) 1);
        WriteResponseItem<Byte> writeResponseItem1 = new WriteResponseItem<>(writeRequestItem1, ResponseCode.OK);
        WriteResponseItem<Byte> writeResponseItem2 = new WriteResponseItem<>(writeRequestItem2, ResponseCode.OK);
        responseItems.add(writeResponseItem1);
        responseItems.add(writeResponseItem2);
        PlcWriteResponse plcWriteResponse = new PlcWriteResponse(plcWriteRequest, responseItems);
        Optional<WriteResponseItem<Byte>> responseValue1 = plcWriteResponse.getValue(writeRequestItem1);
        Optional<WriteResponseItem<Byte>> responseValue2 = plcWriteResponse.getValue(writeRequestItem2);
        assertThat(responseValue1).isEqualTo(Optional.of(writeResponseItem1)).withFailMessage("Unexpected items in response items");
        assertThat(responseValue2).isEqualTo(Optional.of(writeResponseItem2)).withFailMessage("Unexpected items in response items");
    }

    @Test
    @Category(FastTests.class)
    public void nonExistingItemPlcWriteResponseGetValue() {
        PlcWriteRequest plcWriteRequest = new PlcWriteRequest();
        List<WriteResponseItem<?>> responseItems = new ArrayList<>();
        MockAddress address = new MockAddress("mock:/DATA");
        WriteRequestItem<Byte> nonExistingWriteRequestItem = new WriteRequestItem<>(Byte.class, address, (byte) 1);
        PlcWriteResponse plcWriteResponse = new PlcWriteResponse(plcWriteRequest, responseItems);
        Optional<WriteResponseItem<Byte>> responseValue1 = plcWriteResponse.getValue(nonExistingWriteRequestItem);
        assertThat(responseValue1).isEqualTo(Optional.empty()).withFailMessage("Unexpected items in response items");
    }

    @Test
    @Category(FastTests.class)
    public void bulkPlcReadResponseGetValue() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        List<ReadResponseItem<?>> responseItems = new ArrayList<>();
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem1 = new ReadRequestItem<>(Byte.class, address, 1);
        ReadRequestItem<Byte> readRequestItem2 = new ReadRequestItem<>(Byte.class, address, 1);
        ReadResponseItem<Byte> readResponseItem1 = new ReadResponseItem<>(readRequestItem1, ResponseCode.OK, Collections.emptyList());
        ReadResponseItem<Byte> readResponseItem2 = new ReadResponseItem<>(readRequestItem2, ResponseCode.OK, Collections.emptyList());
        responseItems.add(readResponseItem1);
        responseItems.add(readResponseItem2);
        PlcReadResponse plcReadResponse = new PlcReadResponse(plcReadRequest, responseItems);
        Optional<ReadResponseItem<Byte>> responseValue1 = plcReadResponse.getValue(readRequestItem1);
        Optional<ReadResponseItem<Byte>> responseValue2 = plcReadResponse.getValue(readRequestItem2);
        assertThat(responseValue1).isEqualTo(Optional.of(readResponseItem1)).withFailMessage("Unexpected items in response items");
        assertThat(responseValue2).isEqualTo(Optional.of(readResponseItem2)).withFailMessage("Unexpected items in response items");
    }

    @Test
    @Category(FastTests.class)
    public void nonExistingItemPlcReadResponseGetValue() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        List<ReadResponseItem<?>> responseItems = new ArrayList<>();
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> nonExistingReadRequestItem = new ReadRequestItem<>(Byte.class, address, 1);
        PlcReadResponse plcReadResponse = new PlcReadResponse(plcReadRequest, responseItems);
        Optional<ReadResponseItem<Byte>> responseValue1 = plcReadResponse.getValue(nonExistingReadRequestItem);
        assertThat(responseValue1).isEqualTo(Optional.empty()).withFailMessage("Unexpected items in response items");
    }

}
