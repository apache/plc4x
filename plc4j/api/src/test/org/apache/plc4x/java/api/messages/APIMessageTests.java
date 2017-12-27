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
import org.apache.plc4x.java.api.types.ResponseCode;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertTrue;

class APIMessageTests {

    @Test
    @Tag("fast")
    void readRequestItemSize() {
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem readRequestItem = new ReadRequestItem(Byte.class, address, 1);
        assertTrue(readRequestItem.getAddress().equals(address), "Unexpected address");
        assertTrue(readRequestItem.getDatatype() == Byte.class, "Unexpected data type");
        assertTrue(readRequestItem.getSize() == 1, "Unexpected size");
    }

    @Test
    @Tag("fast")
    void readRequestItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem readRequestItem = new ReadRequestItem(Byte.class, address);
        assertTrue(readRequestItem.getAddress().equals(address), "Unexpected address");
        assertTrue(readRequestItem.getDatatype() == Byte.class, "Unexpected data type");
        assertTrue(readRequestItem.getSize() == 1, "Unexpected size");
    }

    @Test
    @Tag("fast")
    void readResponseItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem readRequestItem = new ReadRequestItem(Byte.class, address, 1);
        ReadResponseItem readResponseItem = new  ReadResponseItem(readRequestItem, ResponseCode.OK, Collections.emptyList());
        assertTrue(readResponseItem.getResponseCode() ==  ResponseCode.OK, "Unexpected response code");
        assertTrue(readResponseItem.getValues().isEmpty(), "List should be empty");
        assertTrue(readResponseItem.getRequestItem().equals(readRequestItem), "Unexpected read request item");
    }

    @Test
    @Tag("fast")
    void writeRequestItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        WriteRequestItem writeRequestItem = new WriteRequestItem(Byte.class, address, (byte) 0x45);
        assertTrue(writeRequestItem.getAddress().equals(address), "Unexpected address");
        assertTrue(writeRequestItem.getDatatype() == Byte.class, "Unexpected data type");
        assertTrue((Byte) writeRequestItem.getValues()[0] == 0x45, "Unexpected value");
    }

    @Test
    @Tag("fast")
    void writeRequestItems() {
        MockAddress address = new MockAddress("mock:/DATA");
        Byte data[] = { (byte) 0x23, (byte) 0x84 };
        WriteRequestItem writeRequestItem = new WriteRequestItem(Byte.class, address, data);
        assertTrue(writeRequestItem.getAddress().equals(address), "Unexpected address");
        assertTrue(writeRequestItem.getDatatype() == Byte.class, "Unexpected data type");
        assertTrue((Byte) writeRequestItem.getValues()[0] == (byte) 0x23, "Unexpected value");
        assertTrue((Byte) writeRequestItem.getValues()[1] == (byte) 0x84, "Unexpected value");
    }

    @Test
    @Tag("fast")
    void writeResponseItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        WriteRequestItem writeRequestItem = new WriteRequestItem(Byte.class, address, (byte) 0x3B);
        WriteResponseItem writeResponseItem = new  WriteResponseItem(writeRequestItem, ResponseCode.OK);
        assertTrue(writeResponseItem.getResponseCode() ==  ResponseCode.OK, "Unexpected response code");
        assertTrue(writeResponseItem.getRequestItem().equals(writeRequestItem),  "Unexpected response item");
    }

    @Test
    @Tag("fast")
    void plcReadRequestEmpty() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        assertTrue(plcReadRequest.getReadRequestItems().isEmpty(), "Request items not empty");
        assertTrue(plcReadRequest.getNumItems() == 0, "Expected request items to be zero");
    }

    @Test
    @Tag("fast")
    void plcReadRequestAddress() {
        MockAddress address = new MockAddress("mock:/DATA");
        PlcReadRequest plcReadRequest = new PlcReadRequest(Byte.class, address);
        assertTrue(plcReadRequest.getReadRequestItems().size() == 1, "Expected one request item");
        assertTrue(plcReadRequest.getNumItems() == 1, "Expected one request item");
    }

    @Test
    @Tag("fast")
    void plcReadRequestSize() {
        MockAddress address = new MockAddress("mock:/DATA");
        PlcReadRequest plcReadRequest = new PlcReadRequest(Byte.class, address, 1);
        assertTrue(plcReadRequest.getReadRequestItems().size() == 1, "Expected one request item");
        assertTrue(plcReadRequest.getNumItems() == 1, "Expected one request item");
    }

    @Test
    @Tag("fast")
    void plcReadRequestAddItem() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        assertTrue(plcReadRequest.getReadRequestItems().isEmpty(), "Request items not empty");
        assertTrue(plcReadRequest.getNumItems() == 0, "Expected request items to be zero");
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem readRequestItem = new ReadRequestItem(Byte.class, address, 1);
        plcReadRequest.addItem(readRequestItem);
        assertTrue(plcReadRequest.getReadRequestItems().size() == 1, "Expected one request item");
        assertTrue(plcReadRequest.getNumItems() == 1, "Expected one request item");
    }

    @Test
    @Tag("fast")
    void plcReadResponse() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        ArrayList<ReadResponseItem> responseItems  = new ArrayList<>();
        MockAddress address = new MockAddress("mock:/DATA");
        ReadRequestItem readRequestItem = new ReadRequestItem(Byte.class, address, 1);
        ReadResponseItem readResponseItem = new  ReadResponseItem(readRequestItem, ResponseCode.OK, Collections.emptyList());
        responseItems.add(readResponseItem);
        PlcReadResponse plcReadResponse = new PlcReadResponse(plcReadRequest, responseItems);
        assertTrue(plcReadResponse.getRequest().getNumItems() == 0, "Unexpected number of response items");
        assertTrue(plcReadResponse.getRequest().equals(plcReadRequest), "Unexpected read request");
        assertTrue(plcReadResponse.getResponseItems().size() == 1, "Unexpected number of response items");
        assertTrue(plcReadResponse.getResponseItems().containsAll(responseItems), "Unexpected items in response items");
    }

    @Test
    @Tag("fast")
    void plcWriteRequestItem() {
        MockAddress address = new MockAddress("mock:/DATA");
        WriteRequestItem writeRequestItem = new WriteRequestItem(Byte.class, address, (byte) 0x45);

        assertTrue(writeRequestItem.getAddress().equals(address), "Unexpected address");
        assertTrue(writeRequestItem.getDatatype() == Byte.class, "Unexpected data type");
        assertTrue((Byte) writeRequestItem.getValues()[0] == 0x45, "Unexpected value");
    }

    @Test
    @Tag("fast")
    void plcWriteRequestEmpty() {
        PlcWriteRequest plcWriteRequest = new PlcWriteRequest();
        assertTrue(plcWriteRequest.getRequestItems().isEmpty(), "Request items not empty");
        assertTrue(plcWriteRequest.getNumItems() == 0, "Expected request items to be zero");
    }

    @Test
    @Tag("fast")
    void plcWriteRequestObject() {
        MockAddress address = new MockAddress("mock:/DATA");
        PlcWriteRequest plcWriteRequest = new PlcWriteRequest(Byte.class, address, (byte) 0x33);
        assertTrue(plcWriteRequest.getRequestItems().size() == 1, "Expected no request item");
        assertTrue(plcWriteRequest.getNumItems() == 1, "Expected one request item");
        Object[] values = plcWriteRequest.getRequestItems().get(0).getValues();
        assertTrue((byte)values[0] == (byte) 0x33, "Expected value 0x33");
    }

    @Test
    @Tag("fast")
    void plcWriteRequestObjects() {
        MockAddress address = new MockAddress("mock:/DATA");
        Byte[] data = {(byte)0x22, (byte)0x66};
        PlcWriteRequest plcWriteRequest = new PlcWriteRequest(Byte.class, address, data);
        assertTrue(plcWriteRequest.getRequestItems().size() == 1, "Expected one request item");
        assertTrue(plcWriteRequest.getNumItems() == 1, "Expected one request item");
        Byte[] values = (Byte[])plcWriteRequest.getRequestItems().get(0).getValues();
        assertTrue(values[0] == (byte) 0x22, "Expected value 0x22");
        assertTrue(values[1] == (byte) 0x66, "Expected value 0x66");
    }

    @Test
    @Tag("fast")
    void plcWriteResponse() {
        PlcWriteRequest plcWriteRequest = new PlcWriteRequest();
        ArrayList<WriteResponseItem> responseItems  = new ArrayList<>();
        MockAddress address = new MockAddress("mock:/DATA");
        WriteRequestItem writeRequestItem = new WriteRequestItem(Byte.class, address, 1);
        WriteResponseItem writeResponseItem = new WriteResponseItem(writeRequestItem, ResponseCode.OK);
        responseItems.add(writeResponseItem);
        PlcWriteResponse plcReadResponse = new PlcWriteResponse(plcWriteRequest, responseItems);
        assertTrue(plcReadResponse.getRequest().getNumItems() == 0, "Unexpected number of response items");
        assertTrue(plcReadResponse.getRequest().equals(plcWriteRequest), "Unexpected read request");
        assertTrue(plcReadResponse.getResponseItems().size() == 1, "Unexpected number of response items");
        assertTrue(plcReadResponse.getResponseItems().containsAll(responseItems), "Unexpected items in response items");
    }
}
