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
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.*;

class APIMessageTests {

    @Test
    @Tag("fast")
    void readRequestItemSize() {
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, address, 1);
        assertTrue(readRequestItem.getAddress().equals(address), "Unexpected address");
        assertTrue(readRequestItem.getDatatype() == Byte.class, "Unexpected data type");
        assertTrue(readRequestItem.getSize() == 1, "Unexpected size");
    }

    @Test
    @Tag("fast")
    void readRequestItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, address);
        assertTrue(readRequestItem.getAddress().equals(address), "Unexpected address");
        assertTrue(readRequestItem.getDatatype() == Byte.class, "Unexpected data type");
        assertTrue(readRequestItem.getSize() == 1, "Unexpected size");
    }

    @Test
    @Tag("fast")
    void readResponseItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, address, 1);
        ReadResponseItem<Byte> readResponseItem = new ReadResponseItem<>(readRequestItem, ResponseCode.OK, Collections.emptyList());
        assertTrue(readResponseItem.getResponseCode() == ResponseCode.OK, "Unexpected response code");
        assertTrue(readResponseItem.getValues().isEmpty(), "List should be empty");
        assertTrue(readResponseItem.getRequestItem().equals(readRequestItem), "Unexpected read request item");
    }

    @Test
    @Tag("fast")
    void writeRequestItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        WriteRequestItem<Byte> writeRequestItem = new WriteRequestItem<>(Byte.class, address, (byte) 0x45);

        assertTrue(writeRequestItem.getAddress().equals(address), "Unexpected address");
        assertTrue(writeRequestItem.getDatatype() == Byte.class, "Unexpected data type");
        assertTrue(writeRequestItem.getValues().get(0) == 0x45, "Unexpected value");
    }

    @Test
    @Tag("fast")
    void writeRequestItems() {
        MockAddress address = new MockAddress("mock:/DATA");
        Byte data[] = {(byte) 0x23, (byte) 0x84};
        WriteRequestItem<Byte> writeRequestItem = new WriteRequestItem<>(Byte.class, address, data);
        assertTrue(writeRequestItem.getAddress().equals(address), "Unexpected address");
        assertTrue(writeRequestItem.getDatatype() == Byte.class, "Unexpected data type");
        assertTrue(writeRequestItem.getValues().get(0) == (byte) 0x23, "Unexpected value");
        assertTrue(writeRequestItem.getValues().get(1) == (byte) 0x84, "Unexpected value");
    }

    @Test
    @Tag("fast")
    void writeResponseItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        WriteRequestItem<Byte> writeRequestItem = new WriteRequestItem<>(Byte.class, address, (byte) 0x3B);
        WriteResponseItem<Byte> writeResponseItem = new WriteResponseItem<>(writeRequestItem, ResponseCode.OK);
        assertTrue(writeResponseItem.getResponseCode() == ResponseCode.OK, "Unexpected response code");
        assertTrue(writeResponseItem.getRequestItem().equals(writeRequestItem), "Unexpected response item");
    }

    @Test
    @Tag("fast")
    void plcReadRequestEmpty() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        assertTrue(plcReadRequest.getRequestItems().isEmpty(), "Request items not empty");
        assertTrue(plcReadRequest.getNumberOfItems() == 0, "Expected request items to be zero");
    }

    @Test
    @Tag("fast")
    void plcReadRequestAddress() {
        MockAddress address = new MockAddress("mock:/DATA");
        PlcReadRequest plcReadRequest = new TypeSafePlcReadRequest<>(Byte.class, address);
        assertTrue(plcReadRequest.getRequestItems().size() == 1, "Expected one request item");
        assertTrue(plcReadRequest.getNumberOfItems() == 1, "Expected one request item");
    }

    @Test
    @Tag("fast")
    void plcReadRequestSize() {
        MockAddress address = new MockAddress("mock:/DATA");
        PlcReadRequest plcReadRequest = PlcReadRequest.builder().addItem(Byte.class, address, (byte) 1).build(Byte.class);
        assertTrue(plcReadRequest.getRequestItems().size() == 1, "Expected one request item");
        assertTrue(plcReadRequest.getNumberOfItems() == 1, "Expected one request item");
    }

    @Test
    @Tag("fast")
    void plcReadRequestAddItem() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        assertTrue(plcReadRequest.getRequestItems().isEmpty(), "Request items not empty");
        assertTrue(plcReadRequest.getNumberOfItems() == 0, "Expected request items to be zero");
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, address, (byte) 1);
        plcReadRequest.addItem(readRequestItem);
        assertTrue(plcReadRequest.getRequestItems().size() == 1, "Expected one request item");
        assertTrue(plcReadRequest.getNumberOfItems() == 1, "Expected one request item");
    }

    @Test
    @Tag("fast")
    void plcReadResponse() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        List<ReadResponseItem<?>> responseItems = new ArrayList<>();
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, address, 1);
        ReadResponseItem<Byte> readResponseItem = new ReadResponseItem<>(readRequestItem, ResponseCode.OK, Collections.emptyList());
        responseItems.add(readResponseItem);
        PlcReadResponse plcReadResponse = new PlcReadResponse(plcReadRequest, responseItems);
        assertTrue(plcReadResponse.getRequest().getNumberOfItems() == 0, "Unexpected number of response items");
        assertTrue(plcReadResponse.getRequest().equals(plcReadRequest), "Unexpected read request");
        assertTrue(plcReadResponse.getResponseItems().size() == 1, "Unexpected number of response items");
        assertTrue(plcReadResponse.getResponseItems().containsAll(responseItems), "Unexpected items in response items");
    }

    @Test
    @Tag("fast")
    void plcWriteRequestEmpty() {
        PlcWriteRequest plcWriteRequest = new PlcWriteRequest();
        assertTrue(plcWriteRequest.getRequestItems().isEmpty(), "Request items not empty");
        assertTrue(plcWriteRequest.getNumberOfItems() == 0, "Expected request items to be zero");
    }

    @Test
    @Tag("fast")
    void plcWriteRequestObject() {
        MockAddress address = new MockAddress("mock:/DATA");
        PlcWriteRequest plcWriteRequest = new TypeSafePlcWriteRequest<>(Byte.class, address, (byte) 0x33);
        assertTrue(plcWriteRequest.getRequestItems().size() == 1, "Expected no request item");
        assertTrue(plcWriteRequest.getNumberOfItems() == 1, "Expected one request item");
        List values = plcWriteRequest.getRequestItems().get(0).getValues();
        assertTrue((byte) values.get(0) == (byte) 0x33, "Expected value 0x33");
    }

    @Test
    @Tag("fast")
    void plcWriteRequestObjects() {
        MockAddress address = new MockAddress("mock:/DATA");
        Byte[] data = {(byte) 0x22, (byte) 0x66};
        PlcWriteRequest plcWriteRequest = new TypeSafePlcWriteRequest<>(Byte.class, address, data);
        assertTrue(plcWriteRequest.getRequestItems().size() == 1, "Expected one request item");
        assertTrue(plcWriteRequest.getNumberOfItems() == 1, "Expected one request item");
        List values = plcWriteRequest.getRequestItems().get(0).getValues();
        assertTrue((Byte) values.get(0) == (byte) 0x22, "Expected value 0x22");
        assertTrue((Byte) values.get(1) == (byte) 0x66, "Expected value 0x66");
    }

    @Test
    @Tag("fast")
    void plcWriteResponse() {
        PlcWriteRequest plcWriteRequest = new PlcWriteRequest();
        List<WriteResponseItem<?>> responseItems = new ArrayList<>();
        MockAddress address = new MockAddress("mock:/DATA");
        WriteRequestItem<Byte> writeRequestItem = new WriteRequestItem<>(Byte.class, address, (byte) 1);
        WriteResponseItem<Byte> writeResponseItem = new WriteResponseItem<>(writeRequestItem, ResponseCode.OK);
        responseItems.add(writeResponseItem);
        PlcWriteResponse plcReadResponse = new PlcWriteResponse(plcWriteRequest, responseItems);
        assertTrue(plcReadResponse.getRequest().getNumberOfItems() == 0, "Unexpected number of response items");
        assertTrue(plcReadResponse.getRequest().equals(plcWriteRequest), "Unexpected read request");
        assertTrue(plcReadResponse.getResponseItems().size() == 1, "Unexpected number of response items");
        assertTrue(plcReadResponse.getResponseItems().containsAll(responseItems), "Unexpected items in response items");
    }

    @Test
    @Tag("fast")
    void bulkPlcWriteResponseGetValue() {
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
        assertEquals(Optional.of(writeResponseItem1), responseValue1, "Unexpected items in response items");
        assertEquals(Optional.of(writeResponseItem2), responseValue2, "Unexpected items in response items");
    }

    @Test
    @Tag("fast")
    void nonExistingItemPlcWriteResponseGetValue() {
        PlcWriteRequest plcWriteRequest = new PlcWriteRequest();
        List<WriteResponseItem<?>> responseItems = new ArrayList<>();
        MockAddress address = new MockAddress("mock:/DATA");
        WriteRequestItem<Byte> nonExistingWriteRequestItem = new WriteRequestItem<>(Byte.class, address, (byte) 1);
        PlcWriteResponse plcWriteResponse = new PlcWriteResponse(plcWriteRequest, responseItems);
        Optional<WriteResponseItem<Byte>> responseValue1 = plcWriteResponse.getValue(nonExistingWriteRequestItem);
        assertEquals(Optional.empty(), responseValue1, "Unexpected items in response items");
    }

    @Test
    @Tag("fast")
    void bulkPlcReadResponseGetValue() {
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
        assertEquals(Optional.of(readResponseItem1), responseValue1, "Unexpected items in response items");
        assertEquals(Optional.of(readResponseItem2), responseValue2, "Unexpected items in response items");
    }

    @Test
    @Tag("fast")
    void nonExistingItemPlcReadResponseGetValue() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        List<ReadResponseItem<?>> responseItems = new ArrayList<>();
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem<Byte> nonExistingReadRequestItem = new ReadRequestItem<>(Byte.class, address, 1);
        PlcReadResponse plcReadResponse = new PlcReadResponse(plcReadRequest, responseItems);
        Optional<ReadResponseItem<Byte>> responseValue1 = plcReadResponse.getValue(nonExistingReadRequestItem);
        assertEquals(Optional.empty(), responseValue1, "Unexpected items in response items");
    }

    @Disabled("FIXME: timing issues")
    @Test
    void accessReadResponseFromRequest() throws Exception {
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, null, (byte) 0x0);
        try {
            readRequestItem.getResponseItem().get(10, TimeUnit.MILLISECONDS);
            fail("We should no get anything because of the short Timeout");
        } catch (Exception ignore) {
        }
        ReadResponseItem<Byte> readResponseItem = new ReadResponseItem<>(readRequestItem, ResponseCode.OK, Collections.singletonList((byte) 0x0));
        ReadResponseItem<Byte> byteReadResponseItem = readRequestItem.getResponseItem().get(10, TimeUnit.MILLISECONDS);
        assertEquals(readResponseItem, byteReadResponseItem);
    }

    @Disabled("FIXME: timing issues")
    @Test
    void accessWriteResponseFromRequest() throws Exception {
        WriteRequestItem<Byte> writeRequestItem = new WriteRequestItem<>(Byte.class, null, (byte) 0x0);
        try {
            writeRequestItem.getResponseItem().get(10, TimeUnit.MILLISECONDS);
            fail("We should no get anything because of the short Timeout");
        } catch (Exception ignore) {
        }
        WriteResponseItem<Byte> writeResponseItem = new WriteResponseItem<>(writeRequestItem, ResponseCode.OK);
        WriteResponseItem<Byte> byteWriteResponseItem = writeRequestItem.getResponseItem().get(10, TimeUnit.MILLISECONDS);
        assertEquals(writeResponseItem, byteWriteResponseItem);
    }

    @Disabled("FIXME: timing issues")
    @Test
    void accessResponseFromRequestMultiThreaded() throws Exception {
        ReadRequestItem<Byte> readRequestItem = new ReadRequestItem<>(Byte.class, null, (byte) 0x0);
        CountDownLatch successLatch = new CountDownLatch(10);
        CountDownLatch failLatch = new CountDownLatch(10);
        for (int i = 0; i < 10; i++) {
            new Thread(() -> {
                try {
                    readRequestItem.getResponseItem().get(10, TimeUnit.SECONDS);
                    successLatch.countDown();
                } catch (Exception ignore) {
                    failLatch.countDown();
                }
            }).start();
        }
        assertEquals(10, successLatch.getCount());

        ReadResponseItem<Byte> readResponseItem = new ReadResponseItem<>(readRequestItem, ResponseCode.OK, Collections.singletonList((byte) 0x0));
        ReadResponseItem<Byte> byteReadResponseItem = readRequestItem.getResponseItem().get(10, TimeUnit.MILLISECONDS);
        assertEquals(readResponseItem, byteReadResponseItem);
        assertEquals(0, successLatch.getCount());
        assertEquals(10, failLatch.getCount());
    }

}
