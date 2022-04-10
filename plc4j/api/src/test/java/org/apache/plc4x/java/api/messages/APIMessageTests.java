/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.api.messages;

public class APIMessageTests {

    /*@Test
    @Category(FastTests.class)
    public void readRequestItemSize() {
        MockField field = new MockField("mock:/DATA");
        PlcReadRequestItem<Byte> readRequestItem = new PlcReadRequestItem<>(Byte.class, field, 1);
        assertThat("Unexpected field", readRequestItem.getField(), equalTo(field));
        assertThat("Unexpected data type", readRequestItem.getDatatype(), equalTo(Byte.class));
        assertThat("Unexpected size", readRequestItem.getSize(), equalTo(1));
    }

    @Test
    @Category(FastTests.class)
    public void readRequestItem() {
        MockField field = new MockField("mock:/DATA");
        PlcReadRequestItem<Byte> readRequestItem = new PlcReadRequestItem<>(Byte.class, field);
        assertThat("Unexpected field", readRequestItem.getField(), equalTo(field));
        assertThat("Unexpected data type", readRequestItem.getDatatype(), equalTo(Byte.class));
        assertThat("Unexpected size", readRequestItem.getSize(), equalTo(1));
    }

    @Test
    @Category(FastTests.class)
    public void readResponseItem() {
        MockField field = new MockField("mock:/DATA");
        PlcReadRequestItem<Byte> readRequestItem = new PlcReadRequestItem<>(Byte.class, field, 1);
        PlcReadResponseItem<Byte> readResponseItem = new PlcReadResponseItem<>(readRequestItem, PlcResponseCode.OK);
        assertThat("Unexpected response code", readResponseItem.getResponseCode(), equalTo(PlcResponseCode.OK));
        assertThat(readResponseItem.getValues(), empty());
        assertThat("Unexpected read request item", readResponseItem.getRequestItem(), equalTo(readRequestItem));
    }

    @Test
    @Category(FastTests.class)
    public void writeRequestItem() {
        MockField field = new MockField("mock:/DATA");
        PlcWriteRequestItem<Byte> writeRequestItem = new PlcWriteRequestItem<>(Byte.class, field, (byte) 0x45);

        assertThat("Unexpected field", writeRequestItem.getField(), equalTo(field));
        assertThat("Unexpected data type", writeRequestItem.getDatatype(), equalTo(Byte.class));
        assertThat("Unexpected value", writeRequestItem.getValues().get(0), equalTo((byte) 0x45));
    }

    @Test
    @Category(FastTests.class)
    public void writeRequestItems() {
        MockField field = new MockField("mock:/DATA");
        Byte data[] = {(byte) 0x23, (byte) 0x84};
        PlcWriteRequestItem<Byte> writeRequestItem = new PlcWriteRequestItem<>(Byte.class, field, data);
        assertThat("Unexpected field", writeRequestItem.getField(), equalTo(field));
        assertThat("Unexpected data type", writeRequestItem.getDatatype(), equalTo(Byte.class));
        assertThat("Unexpected value", writeRequestItem.getValues().get(0), equalTo((byte) 0x23));
        assertThat("Unexpected value", writeRequestItem.getValues().get(1), equalTo((byte) 0x84));
    }

    @Test
    @Category(FastTests.class)
    public void writeResponseItem() {
        MockField field = new MockField("mock:/DATA");
        PlcWriteRequestItem<Byte> writeRequestItem = new PlcWriteRequestItem<>(Byte.class, field, (byte) 0x3B);
        PlcWriteResponseItem<Byte> writeResponseItem = new PlcWriteResponseItem<>(writeRequestItem, PlcResponseCode.OK);
        assertThat("Unexpected response code", writeResponseItem.getResponseCode(), equalTo(PlcResponseCode.OK));
        assertThat("Unexpected response item", writeResponseItem.getRequestItem(), equalTo(writeRequestItem));
    }

    @Test
    public void plcProprietaryRequest() {
        Object customMessage = new Object();
        PlcProprietaryRequest<Object> plcProprietaryRequest = new PlcProprietaryRequest<>(customMessage);
        assertThat("Unexpected request type", plcProprietaryRequest.getProprietaryRequest(), equalTo(customMessage));
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
    public void plcReadRequestField() {
        MockField field = new MockField("mock:/DATA");
        PlcReadRequest plcReadRequest = new TypeSafePlcReadRequest<>(Byte.class, field);
        assertThat("Expected one request item", plcReadRequest.getRequestItems(), hasSize(1));
        assertThat("Expected one request item", plcReadRequest.getNumberOfItems(), equalTo(1));
    }

    @Test
    @Category(FastTests.class)
    public void plcReadRequestSize() {
        MockField field = new MockField("mock:/DATA");
        PlcReadRequest plcReadRequest = PlcReadRequest.builder().addItem(Byte.class, field, (byte) 1).build(Byte.class);
        assertThat("Expected one request item", plcReadRequest.getRequestItems(), hasSize(1));
        assertThat("Expected one request item", plcReadRequest.getNumberOfItems(), equalTo(1));
    }

    @Test
    @Category(FastTests.class)
    public void plcReadRequestAddItem() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        assertThat(plcReadRequest.getRequestItems(), empty());
        assertThat("Expected request items to be zero", plcReadRequest.getNumberOfItems(), equalTo(0));
        MockField field = new MockField("mock:/DATA");
        PlcReadRequestItem<Byte> readRequestItem = new PlcReadRequestItem<>(Byte.class, field, (byte) 1);
        plcReadRequest.addItem(readRequestItem);
        assertThat("Expected one request item", plcReadRequest.getRequestItems(), hasSize(1));
        assertThat("Expected one request item", plcReadRequest.getNumberOfItems(), equalTo(1));
    }

    @Test
    @Category(FastTests.class)
    public void plcReadResponse() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        List<PlcReadResponseItem<?>> responseItems = new ArrayList<>();
        MockField field = new MockField("mock:/DATA");
        PlcReadRequestItem<Byte> readRequestItem = new PlcReadRequestItem<>(Byte.class, field, 1);
        PlcReadResponseItem<Byte> readResponseItem = new PlcReadResponseItem<>(readRequestItem, PlcResponseCode.OK);
        responseItems.add(readResponseItem);
        PlcReadResponse plcReadResponse = new PlcReadResponse(plcReadRequest, responseItems);
        assertThat("Unexpected number of response items", plcReadResponse.getProprietaryRequest().getNumberOfItems(), equalTo(0));
        assertThat("Unexpected read request", plcReadResponse.getProprietaryRequest(), equalTo(plcReadRequest));
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
        MockField field = new MockField("mock:/DATA");
        PlcWriteRequest plcWriteRequest = new TypeSafePlcWriteRequest<>(Byte.class, field, (byte) 0x33);
        assertThat("Expected one request item", plcWriteRequest.getRequestItems(), hasSize(1));
        assertThat("Expected one request item", plcWriteRequest.getNumberOfItems(), equalTo(1));
        List values = plcWriteRequest.getRequestItems().get(0).getValues();
        assertThat(values.get(0), equalTo((byte) 0x33));
    }

    @Test
    @Category(FastTests.class)
    public void plcWriteRequestObjects() {
        MockField field = new MockField("mock:/DATA");
        Byte[] data = {(byte) 0x22, (byte) 0x66};
        PlcWriteRequest plcWriteRequest = new TypeSafePlcWriteRequest<>(Byte.class, field, data);
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
        List<PlcWriteResponseItem<?>> responseItems = new ArrayList<>();
        MockField field = new MockField("mock:/DATA");
        PlcWriteRequestItem<Byte> writeRequestItem = new PlcWriteRequestItem<>(Byte.class, field, (byte) 1);
        PlcWriteResponseItem<Byte> writeResponseItem = new PlcWriteResponseItem<>(writeRequestItem, PlcResponseCode.OK);
        responseItems.add(writeResponseItem);
        PlcWriteResponse plcReadResponse = new PlcWriteResponse(plcWriteRequest, responseItems);
        assertThat("Unexpected number of response items", plcReadResponse.getProprietaryRequest().getNumberOfItems(), equalTo(0));
        assertThat("Unexpected read request", plcReadResponse.getProprietaryRequest(), equalTo(plcWriteRequest));
        assertThat("Unexpected number of response items", plcReadResponse.getResponseItems(), hasSize(1));
        assertThat("Unexpected items in response items", plcReadResponse.getResponseItems(), contains(writeResponseItem));
    }

    @Test
    @Category(FastTests.class)
    public void bulkPlcWriteResponseGetValue() {
        PlcWriteRequest plcWriteRequest = new PlcWriteRequest();
        List<PlcWriteResponseItem<?>> responseItems = new ArrayList<>();
        MockField field = new MockField("mock:/DATA");
        PlcWriteRequestItem<Byte> writeRequestItem1 = new PlcWriteRequestItem<>(Byte.class, field, (byte) 1);
        PlcWriteRequestItem<Byte> writeRequestItem2 = new PlcWriteRequestItem<>(Byte.class, field, (byte) 1);
        PlcWriteResponseItem<Byte> writeResponseItem1 = new PlcWriteResponseItem<>(writeRequestItem1, PlcResponseCode.OK);
        PlcWriteResponseItem<Byte> writeResponseItem2 = new PlcWriteResponseItem<>(writeRequestItem2, PlcResponseCode.OK);
        responseItems.add(writeResponseItem1);
        responseItems.add(writeResponseItem2);
        PlcWriteResponse plcWriteResponse = new PlcWriteResponse(plcWriteRequest, responseItems);
        Optional<PlcWriteResponseItem<Byte>> responseValue1 = plcWriteResponse.getValue(writeRequestItem1);
        Optional<PlcWriteResponseItem<Byte>> responseValue2 = plcWriteResponse.getValue(writeRequestItem2);
        assertThat("Unexpected items in response items", responseValue1, equalTo(Optional.of(writeResponseItem1)));
        assertThat("Unexpected items in response items", responseValue2, equalTo(Optional.of(writeResponseItem2)));
    }

    @Test
    @Category(FastTests.class)
    public void nonExistingItemPlcWriteResponseGetValue() {
        PlcWriteRequest plcWriteRequest = new PlcWriteRequest();
        List<PlcWriteResponseItem<?>> responseItems = new ArrayList<>();
        MockField field = new MockField("mock:/DATA");
        PlcWriteRequestItem<Byte> nonExistingWriteRequestItem = new PlcWriteRequestItem<>(Byte.class, field, (byte) 1);
        PlcWriteResponse plcWriteResponse = new PlcWriteResponse(plcWriteRequest, responseItems);
        Optional<PlcWriteResponseItem<Byte>> responseValue1 = plcWriteResponse.getValue(nonExistingWriteRequestItem);
        assertThat("Unexpected items in response items", responseValue1, equalTo(Optional.empty()));
    }

    @Test
    @Category(FastTests.class)
    public void bulkPlcReadResponseGetValue() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        List<PlcReadResponseItem<?>> responseItems = new ArrayList<>();
        MockField field = new MockField("mock:/DATA");
        PlcReadRequestItem<Byte> readRequestItem1 = new PlcReadRequestItem<>(Byte.class, field, 1);
        PlcReadRequestItem<Byte> readRequestItem2 = new PlcReadRequestItem<>(Byte.class, field, 1);
        PlcReadResponseItem<Byte> readResponseItem1 = new PlcReadResponseItem<>(readRequestItem1, PlcResponseCode.OK);
        PlcReadResponseItem<Byte> readResponseItem2 = new PlcReadResponseItem<>(readRequestItem2, PlcResponseCode.OK);
        responseItems.add(readResponseItem1);
        responseItems.add(readResponseItem2);
        PlcReadResponse plcReadResponse = new PlcReadResponse(plcReadRequest, responseItems);
        Optional<PlcReadResponseItem<Byte>> responseValue1 = plcReadResponse.getValue(readRequestItem1);
        Optional<PlcReadResponseItem<Byte>> responseValue2 = plcReadResponse.getValue(readRequestItem2);
        assertThat("Unexpected items in response items", responseValue1, equalTo(Optional.of(readResponseItem1)));
        assertThat("Unexpected items in response items", responseValue2, equalTo(Optional.of(readResponseItem2)));
    }

    @Test
    @Category(FastTests.class)
    public void nonExistingItemPlcReadResponseGetValue() {
        PlcReadRequest plcReadRequest = new PlcReadRequest();
        List<PlcReadResponseItem<?>> responseItems = new ArrayList<>();
        MockField field = new MockField("mock:/DATA");
        PlcReadRequestItem<Byte> nonExistingReadRequestItem = new PlcReadRequestItem<>(Byte.class, field, 1);
        PlcReadResponse plcReadResponse = new PlcReadResponse(plcReadRequest, responseItems);
        Optional<PlcReadResponseItem<Byte>> responseValue1 = plcReadResponse.getValue(nonExistingReadRequestItem);
        assertThat("Unexpected items in response items", responseValue1, equalTo(Optional.empty()));
    }*/

}
