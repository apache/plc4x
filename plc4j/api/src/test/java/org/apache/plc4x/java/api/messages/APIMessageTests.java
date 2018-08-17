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

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

import static org.hamcrest.collection.IsCollectionWithSize.hasSize;
import static org.hamcrest.collection.IsEmptyCollection.empty;
import static org.hamcrest.collection.IsIterableContainingInOrder.contains;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;


public class APIMessageTests {

    @Test
    @Category(FastTests.class)
    public void readRequestItemSize() {
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, address, 1);
        assertThat("Unexpected address", readRequestItem.getAddress(), equalTo(address));
        assertThat("Unexpected data type", readRequestItem.getDatatype(), equalTo(Byte.class));
        assertThat("Unexpected size", readRequestItem.getSize(), equalTo(1));
    }

    @Test
    @Category(FastTests.class)
    public void readRequestItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, address);
        assertThat("Unexpected address", readRequestItem.getAddress(), equalTo(address));
        assertThat("Unexpected data type", readRequestItem.getDatatype(), equalTo(Byte.class));
        assertThat("Unexpected size", readRequestItem.getSize(), equalTo(1));
    }

    @Test
    @Category(FastTests.class)
    public void readResponseItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, address, 1);
        ReadResponseItem<Byte> readResponseItem = new ReadResponseItem<>(readRequestItem, ResponseCode.OK);
        assertThat("Unexpected response code", readResponseItem.getResponseCode(), equalTo(ResponseCode.OK));
        assertThat(readResponseItem.getValues(), empty());
        assertThat("Unexpected read request item", readResponseItem.getRequestItem(), equalTo(readRequestItem));
    }

    @Test
    @Category(FastTests.class)
    public void writeRequestItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        WriteRequestItem<Byte> writeRequestItem = new WriteRequestItem<>(Byte.class, address, (byte) 0x45);

        assertThat("Unexpected address", writeRequestItem.getAddress(), equalTo(address));
        assertThat("Unexpected data type", writeRequestItem.getDatatype(), equalTo(Byte.class));
        assertThat("Unexpected value", writeRequestItem.getValues().get(0), equalTo((byte) 0x45));
    }

    @Test
    @Category(FastTests.class)
    public void writeRequestItems() {
        MockAddress address = new MockAddress("mock:/DATA");
        Byte data[] = {(byte) 0x23, (byte) 0x84};
        WriteRequestItem<Byte> writeRequestItem = new WriteRequestItem<>(Byte.class, address, data);
        assertThat("Unexpected address", writeRequestItem.getAddress(), equalTo(address));
        assertThat("Unexpected data type", writeRequestItem.getDatatype(), equalTo(Byte.class));
        assertThat("Unexpected value", writeRequestItem.getValues().get(0), equalTo((byte) 0x23));
        assertThat("Unexpected value", writeRequestItem.getValues().get(1), equalTo((byte) 0x84));
    }

    @Test
    @Category(FastTests.class)
    public void writeResponseItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        WriteRequestItem<Byte> writeRequestItem = new WriteRequestItem<>(Byte.class, address, (byte) 0x3B);
        WriteResponseItem<Byte> writeResponseItem = new WriteResponseItem<>(writeRequestItem, ResponseCode.OK);
        assertThat("Unexpected response code", writeResponseItem.getResponseCode(), equalTo(ResponseCode.OK));
        assertThat("Unexpected response item", writeResponseItem.getRequestItem(), equalTo(writeRequestItem));
    }

    @Test
    public void plcProprietaryRequest() {
        Object customMessage = new Object();
        PlcProprietaryRequest<Object> plcProprietaryRequest = new PlcProprietaryRequest<>(customMessage);
        assertThat("Unexpected request type", plcProprietaryRequest.getRequest(), equalTo(customMessage));
    }

    @Test
    public void plcProprietaryResponse() {
        Object customMessage = new Object();
        Object customMessageResponse = new Object();
        PlcProprietaryResponse<Object> plcProprietaryResponse = new PlcProprietaryResponse<>(new PlcProprietaryRequest<>(customMessage), customMessageResponse);
        assertThat("Unexpected response type", plcProprietaryResponse.getResponse(), equalTo(customMessageResponse));
    }

    @Test
    @Category(FastTests.class)
    public void plcReadRequestEmpty() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        assertThat(plcReadRequest.getRequestItems(), empty());
        assertThat("Expected request items to be zero", plcReadRequest.getNumberOfItems(), equalTo(0));
    }

    @Test
    @Category(FastTests.class)
    public void plcReadRequestAddress() {
        MockAddress address = new MockAddress("mock:/DATA");
        PlcReadRequest plcReadRequest = new TypeSafePlcReadRequest<>(Byte.class, address);
        assertThat("Expected one request item", plcReadRequest.getRequestItems(), hasSize(1));
        assertThat("Expected one request item", plcReadRequest.getNumberOfItems(), equalTo(1));
    }

    @Test
    @Category(FastTests.class)
    public void plcReadRequestSize() {
        MockAddress address = new MockAddress("mock:/DATA");
        PlcReadRequest plcReadRequest = PlcReadRequest.builder().addItem(Byte.class, address, (byte) 1).build(Byte.class);
        assertThat("Expected one request item", plcReadRequest.getRequestItems(), hasSize(1));
        assertThat("Expected one request item", plcReadRequest.getNumberOfItems(), equalTo(1));
    }

    @Test
    @Category(FastTests.class)
    public void plcReadRequestAddItem() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        assertThat(plcReadRequest.getRequestItems(), empty());
        assertThat("Expected request items to be zero", plcReadRequest.getNumberOfItems(), equalTo(0));
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, address, (byte) 1);
        plcReadRequest.addItem(readRequestItem);
        assertThat("Expected one request item", plcReadRequest.getRequestItems(), hasSize(1));
        assertThat("Expected one request item", plcReadRequest.getNumberOfItems(), equalTo(1));
    }

    @Test
    @Category(FastTests.class)
    public void plcReadResponse() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        List<ReadResponseItem<?>> responseItems = new ArrayList<>();
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, address, 1);
        ReadResponseItem<Byte> readResponseItem = new ReadResponseItem<>(readRequestItem, ResponseCode.OK);
        responseItems.add(readResponseItem);
        PlcReadResponse plcReadResponse = new PlcReadResponse(plcReadRequest, responseItems);
        assertThat("Unexpected number of response items", plcReadResponse.getRequest().getNumberOfItems(), equalTo(0));
        assertThat("Unexpected read request", plcReadResponse.getRequest(), equalTo(plcReadRequest));
        assertThat("Unexpected number of response items", plcReadResponse.getResponseItems(), hasSize(1));
        assertThat("Unexpected items in response items", plcReadResponse.getResponseItems(), contains(readResponseItem));
    }

    @Test
    @Category(FastTests.class)
    public void plcWriteRequestEmpty() {
        PlcWriteRequest plcWriteRequest = new PlcWriteRequest();
        assertThat(plcWriteRequest.getRequestItems(), empty());
        assertThat("Expected request items to be zero", plcWriteRequest.getNumberOfItems(), equalTo(0));
    }

    @Test
    @Category(FastTests.class)
    public void plcWriteRequestObject() {
        MockAddress address = new MockAddress("mock:/DATA");
        PlcWriteRequest plcWriteRequest = new TypeSafePlcWriteRequest<>(Byte.class, address, (byte) 0x33);
        assertThat("Expected one request item", plcWriteRequest.getRequestItems(), hasSize(1));
        assertThat("Expected one request item", plcWriteRequest.getNumberOfItems(), equalTo(1));
        List values = plcWriteRequest.getRequestItems().get(0).getValues();
        assertThat(values.get(0), equalTo((byte) 0x33));
    }

    @Test
    @Category(FastTests.class)
    public void plcWriteRequestObjects() {
        MockAddress address = new MockAddress("mock:/DATA");
        Byte[] data = {(byte) 0x22, (byte) 0x66};
        PlcWriteRequest plcWriteRequest = new TypeSafePlcWriteRequest<>(Byte.class, address, data);
        assertThat("Expected one request item", plcWriteRequest.getRequestItems(), hasSize(1));
        assertThat("Expected one request item", plcWriteRequest.getNumberOfItems(), equalTo(1));
        List values = plcWriteRequest.getRequestItems().get(0).getValues();
        assertThat(values.get(0), equalTo((byte) 0x22));
        assertThat(values.get(1), equalTo((byte) 0x66));
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
        assertThat("Unexpected number of response items", plcReadResponse.getRequest().getNumberOfItems(), equalTo(0));
        assertThat("Unexpected read request", plcReadResponse.getRequest(), equalTo(plcWriteRequest));
        assertThat("Unexpected number of response items", plcReadResponse.getResponseItems(), hasSize(1));
        assertThat("Unexpected items in response items", plcReadResponse.getResponseItems(), contains(writeResponseItem));
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
        assertThat("Unexpected items in response items", responseValue1, equalTo(Optional.of(writeResponseItem1)));
        assertThat("Unexpected items in response items", responseValue2, equalTo(Optional.of(writeResponseItem2)));
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
        assertThat("Unexpected items in response items", responseValue1, equalTo(Optional.empty()));
    }

    @Test
    @Category(FastTests.class)
    public void bulkPlcReadResponseGetValue() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        List<ReadResponseItem<?>> responseItems = new ArrayList<>();
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem1 = new ReadRequestItem<>(Byte.class, address, 1);
        ReadRequestItem<Byte> readRequestItem2 = new ReadRequestItem<>(Byte.class, address, 1);
        ReadResponseItem<Byte> readResponseItem1 = new ReadResponseItem<>(readRequestItem1, ResponseCode.OK);
        ReadResponseItem<Byte> readResponseItem2 = new ReadResponseItem<>(readRequestItem2, ResponseCode.OK);
        responseItems.add(readResponseItem1);
        responseItems.add(readResponseItem2);
        PlcReadResponse plcReadResponse = new PlcReadResponse(plcReadRequest, responseItems);
        Optional<ReadResponseItem<Byte>> responseValue1 = plcReadResponse.getValue(readRequestItem1);
        Optional<ReadResponseItem<Byte>> responseValue2 = plcReadResponse.getValue(readRequestItem2);
        assertThat("Unexpected items in response items", responseValue1, equalTo(Optional.of(readResponseItem1)));
        assertThat("Unexpected items in response items", responseValue2, equalTo(Optional.of(readResponseItem2)));
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
        assertThat("Unexpected items in response items", responseValue1, equalTo(Optional.empty()));
    }

}
