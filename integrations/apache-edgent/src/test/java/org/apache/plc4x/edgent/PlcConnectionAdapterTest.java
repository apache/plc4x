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
package org.apache.plc4x.edgent;

import com.google.gson.JsonObject;
import org.apache.edgent.function.Consumer;
import org.apache.edgent.function.Function;
import org.apache.edgent.function.Supplier;
import org.apache.plc4x.edgent.mock.MockConnection;
import org.apache.plc4x.edgent.mock.MockField;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.messages.items.DefaultLongFieldItem;
import org.apache.plc4x.test.FastTests;
import org.hamcrest.core.IsInstanceOf;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

import static org.hamcrest.Matchers.is;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNot.not;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.hamcrest.core.IsSame.sameInstance;
import static org.junit.Assert.assertThat;

public class PlcConnectionAdapterTest {

    private MockConnection getMockConnection() throws PlcConnectionException {
        return (MockConnection) new PlcDriverManager().getConnection("mock-for-edgent-integration://some-cool-url");
    }

    /*
     * Test the PlcConnectionAdapter(PlcConnection) ctor, getConnection() and close()
     */
    @Test
    @Category(FastTests.class)
    public void testCtor1() throws Exception {
        MockConnection mockConnection = getMockConnection();
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(mockConnection);
        assertThat(mockConnection, sameInstance(adapter.getConnection()));
        // and again... multiple adapter.getConnection() returns the same
        assertThat(mockConnection, sameInstance(adapter.getConnection()));
        adapter.close();
    }

    /*
     * Test the PlcConnectionAdapter(url) ctor, getConnection() and close()
     */
    @Test
    @Category(FastTests.class)
    public void testCtor2() throws Exception {
        MockConnection mockConnection = getMockConnection();
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(mockConnection.getUrl());
        MockConnection mockConnection2 = (MockConnection) adapter.getConnection();
        assertThat(mockConnection, not(sameInstance(mockConnection2)));
        assertThat(mockConnection.getUrl(), sameInstance(mockConnection2.getUrl()));
        // and again... multiple adapter.getConnection() returns the same
        assertThat(mockConnection2, sameInstance(adapter.getConnection()));
        adapter.close();
    }

    private <T> void checkRead(MockConnection connection, PlcReadRequest request, T value) throws InterruptedException, ExecutionException {
        // this is really a tests of our mock tooling but knowing it's behaving as expected
        // will help identify problems in the adapter/supplier/consumer
        PlcField plcField = request.getFields().get(0);
        // TODO: smart value conversion
        connection.setFieldItem(plcField, new DefaultLongFieldItem(0L));

        CompletableFuture<PlcReadResponse> cf = connection.read(request);

        assertThat(cf.isDone(), is(true));
        PlcReadResponse response = cf.get();
        // TODO: fixme
        // assertThat(value, equalTo(response.getResponseItems().get(0).getValues().get(0)));
    }


    private <T> void checkWrite(MockConnection connection, PlcWriteRequest request, T value) throws InterruptedException, ExecutionException {
        // this is really a tests of our mock tooling but knowing it's behaving as expected
        // will help identify problems in the adapter/supplier/consumer
        PlcField plcField = request.getFields().get(0);
        connection.setFieldItem(plcField, new DefaultLongFieldItem(0L));

        CompletableFuture<PlcWriteResponse> cf = connection.write(request);

        assertThat(cf.isDone(), is(true));
        PlcWriteResponse response = cf.get();
        assertThat(response, notNullValue());

        // TODO: fixme
        /*
        Object writtenData = connection.getDataValue(request.getRequestItems().get(0).getField());
        if (writtenData.getClass().isArray()) {
            writtenData = Array.get(writtenData, 0);
        }
        if (List.class.isAssignableFrom(writtenData.getClass())) {
            @SuppressWarnings("unchecked")
            List<Object> writtenDataList = (List<Object>) writtenData;
            writtenData = writtenDataList.get(0);
        }
        assertThat(value, equalTo(writtenData));*/
    }

    /*
     * Verify the adapter yields the appropriate PlcReadRequest for each type and that it works.
     */
    @Test
    @Category(FastTests.class)
    public void testNewPlcReadRequest() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockField address = new MockField(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        {
            PlcReadRequest readRequest = adapter.readRequestBuilder().addItem("test", addressStr).build();
            assertThat(readRequest.getNumberOfFields(), equalTo(1));
            assertThat(readRequest.getField("test"), notNullValue());
            assertThat(readRequest.getField("test"), IsInstanceOf.instanceOf(MockField.class));
            assertThat(((MockField) readRequest.getField("test")).getAddress(), equalTo(addressStr));
            checkRead(connection, readRequest, true);
            checkRead(connection, readRequest, false);
        }
        /*{
            TypeSafePlcReadRequest<Byte> request = PlcConnectionAdapter.newPlcReadRequest(Byte.class, address);
            PlcReadRequestItem<Byte> requestItem = request.getCheckedReadRequestItems().get(0);
            Class<Byte> dataType = requestItem.getDatatype();
            assertThat(dataType, equalTo(Byte.class));
            assertThat(address, sameInstance(requestItem.getField()));
            checkRead(connection, request, (byte) 0x13);
            checkRead(connection, request, (byte) 0x23);
        }
        {
            TypeSafePlcReadRequest<Short> request = PlcConnectionAdapter.newPlcReadRequest(Short.class, address);
            PlcReadRequestItem<Short> requestItem = request.getCheckedReadRequestItems().get(0);
            Class<Short> dataType = requestItem.getDatatype();
            assertThat(dataType, equalTo(Short.class));
            assertThat(address, sameInstance(requestItem.getField()));
            checkRead(connection, request, (short) 13);
            checkRead(connection, request, (short) 23);
        }
        {
            TypeSafePlcReadRequest<Integer> request = PlcConnectionAdapter.newPlcReadRequest(Integer.class, address);
            PlcReadRequestItem<Integer> requestItem = request.getCheckedReadRequestItems().get(0);
            Class<Integer> dataType = requestItem.getDatatype();
            assertThat(dataType, equalTo(Integer.class));
            assertThat(address, sameInstance(requestItem.getField()));
            checkRead(connection, request, 33);
            checkRead(connection, request, -133);
        }
        {
            TypeSafePlcReadRequest<Float> request = PlcConnectionAdapter.newPlcReadRequest(Float.class, address);
            PlcReadRequestItem<Float> requestItem = request.getCheckedReadRequestItems().get(0);
            Class<Float> dataType = requestItem.getDatatype();
            assertThat(dataType, equalTo(Float.class));
            assertThat(address, sameInstance(requestItem.getField()));
            checkRead(connection, request, 43.5f);
            checkRead(connection, request, -143.5f);
        }
        {
            TypeSafePlcReadRequest<String> request = PlcConnectionAdapter.newPlcReadRequest(String.class, address);
            PlcReadRequestItem<String> requestItem = request.getCheckedReadRequestItems().get(0);
            Class<String> dataType = requestItem.getDatatype();
            assertThat(dataType, equalTo(String.class));
            assertThat(address, sameInstance(requestItem.getField()));
            checkRead(connection, request, "ReadySetGo");
            checkRead(connection, request, "OneMoreTime");
        }
        {
            TypeSafePlcReadRequest<Calendar> request = PlcConnectionAdapter.newPlcReadRequest(Calendar.class, address);
            PlcReadRequestItem<Calendar> requestItem = request.getCheckedReadRequestItems().get(0);
            Class<Calendar> dataType = requestItem.getDatatype();
            assertThat(dataType, equalTo(Calendar.class));
            assertThat(address, sameInstance(requestItem.getField()));
            checkRead(connection, request, Calendar.getInstance());
        }*/
        adapter.close();
    }


    /*
     * Verify the adapter yields the appropriate PlcWriteRequest for each type and that it works.
     */
    @Test
    @Category(FastTests.class)
    public void testNewPlcWriteRequest() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockField address = new MockField(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        /*{
            TypeSafePlcWriteRequest<Boolean> request = PlcConnectionAdapter.newPlcWriteRequest(address, true);
            PlcWriteRequestItem<Boolean> requestItem = request.getCheckedRequestItems().get(0);
            Class<Boolean> dataType = requestItem.getDatatype();
            assertThat(dataType, typeCompatibleWith(Boolean.class));
            assertThat(address, sameInstance(requestItem.getField()));
            checkWrite(connection, request, true);
        }
        {
            TypeSafePlcWriteRequest<Byte> request = PlcConnectionAdapter.newPlcWriteRequest(address, (byte) 0x113);
            PlcWriteRequestItem<Byte> requestItem = request.getCheckedRequestItems().get(0);
            Class<Byte> dataType = requestItem.getDatatype();
            assertThat(dataType, typeCompatibleWith(Byte.class));
            assertThat(address, sameInstance(requestItem.getField()));
            checkWrite(connection, request, (byte) 0x113);
        }
        {
            TypeSafePlcWriteRequest<Short> request = PlcConnectionAdapter.newPlcWriteRequest(address, (short) 113);
            PlcWriteRequestItem<Short> requestItem = request.getCheckedRequestItems().get(0);
            Class<Short> dataType = requestItem.getDatatype();
            assertThat(dataType, typeCompatibleWith(Short.class));
            assertThat(address, sameInstance(requestItem.getField()));
            checkWrite(connection, request, (short) 113);
        }
        {
            TypeSafePlcWriteRequest<Integer> request = PlcConnectionAdapter.newPlcWriteRequest(address, 1033);
            PlcWriteRequestItem<Integer> requestItem = request.getCheckedRequestItems().get(0);
            Class<Integer> dataType = requestItem.getDatatype();
            assertThat(dataType, typeCompatibleWith(Integer.class));
            assertThat(address, sameInstance(requestItem.getField()));
            checkWrite(connection, request, 1033);
        }
        {
            TypeSafePlcWriteRequest<Float> request = PlcConnectionAdapter.newPlcWriteRequest(address, 1043.5f);
            PlcWriteRequestItem<Float> requestItem = request.getCheckedRequestItems().get(0);
            Class<Float> dataType = requestItem.getDatatype();
            assertThat(dataType, typeCompatibleWith(Float.class));
            assertThat(address, sameInstance(requestItem.getField()));
            checkWrite(connection, request, 1043.5f);
        }
        {
            TypeSafePlcWriteRequest<String> request = PlcConnectionAdapter.newPlcWriteRequest(address, "A written value");
            PlcWriteRequestItem<String> requestItem = request.getCheckedRequestItems().get(0);
            Class<String> dataType = requestItem.getDatatype();
            assertThat(dataType, typeCompatibleWith(String.class));
            assertThat(address, sameInstance(requestItem.getField()));
            checkWrite(connection, request, "A written value");
        }
        {
            Calendar calValue = Calendar.getInstance();
            TypeSafePlcWriteRequest<Calendar> request = PlcConnectionAdapter.newPlcWriteRequest(address, calValue);
            PlcWriteRequestItem<Calendar> requestItem = request.getCheckedRequestItems().get(0);
            Class<Calendar> dataType = requestItem.getDatatype();
            assertThat(dataType, typeCompatibleWith(Calendar.class));
            assertThat(address, sameInstance(requestItem.getField()));
            checkWrite(connection, request, calValue);
        }*/
        adapter.close();
    }

    /*
     * test PlcConnectionAdapter.newSupplier
     */
    @Test
    @Category(FastTests.class)
    public void testNewSupplier() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockField address = new MockField(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        /*{
            Supplier<Boolean> supplier = adapter.newSupplier(Boolean.class, addressStr);
            assertThat(supplier, not(sameInstance(adapter.newSupplier(Boolean.class, addressStr))));
            checkSupplier(connection, address, supplier, true, false);
        }
        {
            Supplier<Byte> supplier = adapter.newSupplier(Byte.class, addressStr);
            checkSupplier(connection, address, supplier, (byte) 0x1, (byte) 0x2, (byte) 0x3);
        }
        {
            Supplier<Short> supplier = adapter.newSupplier(Short.class, addressStr);
            checkSupplier(connection, address, supplier, (short) 1, (short) 2, (short) 3);
        }
        {
            Supplier<Integer> supplier = adapter.newSupplier(Integer.class, addressStr);
            checkSupplier(connection, address, supplier, 1000, 1001, 1002);
        }
        {
            Supplier<Float> supplier = adapter.newSupplier(Float.class, addressStr);
            checkSupplier(connection, address, supplier, 1000.5f, 1001.5f, 1002.5f);
        }
        {
            Supplier<?> supplier = adapter.newSupplier(String.class, addressStr);
            checkSupplier(connection, address, supplier, "one", "two", "three");
        }*/
        adapter.close();
    }

    /*
     * test PlcConnectionAdapter.newSupplier with read exception
     */
    @Test
    @Category(FastTests.class)
    public void testNewSupplierNeg() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockField address = new MockField(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        /*Supplier<String> supplier = adapter.newSupplier(String.class, addressStr);
        checkSupplier(2, connection, address, supplier, "one", "two", "three");*/

        adapter.close();
    }

    static <T> void checkSupplier(MockConnection connection, PlcField field, Supplier<T> supplier, Object... values) {
        checkSupplier(0, connection, field, supplier, values);
    }

    private static <T> void checkSupplier(int readFailureCountTrigger, MockConnection connection, PlcField field, Supplier<T> supplier, Object... values) {
        // verify that a read failure doesn't kill the consumer
        // it logs (not verified) but returns null (as designed) and keeps working for the subsequent reads
        connection.setReadException(readFailureCountTrigger, "This is a mock read exception");
        int readCount = 0;
        /*for (Object value : values) {
            connection.setDataValue(field, value);
            T readData = supplier.get();
            // System.out.println("checkSupplier"+(readFailureCountTrigger > 0 ? "NEG" : "")+": value:"+value+" readData:"+readData);
            if (readFailureCountTrigger <= 0)
                assertThat(value, equalTo(readData));
            else {
                if (++readCount != readFailureCountTrigger)
                    assertThat(value, equalTo(readData));
                else
                    assertThat(readData, nullValue());
            }
        }*/
    }

    /*
     * test PlcConnectionAdapter.newJsonConsumer(address)
     */
    @Test
    @Category(FastTests.class)
    public void testNewConsumer1() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockField address = new MockField(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        Consumer<?> consumer;

        /*consumer = adapter.newConsumer(Boolean.class, addressStr);
        assertThat(consumer, not(sameInstance(adapter.newConsumer(Boolean.class, addressStr))));
        checkConsumer(connection, address, consumer, true, false);

        consumer = adapter.newConsumer(Byte.class, addressStr);
        checkConsumer(connection, address, consumer, (byte) 0x1, (byte) 0x2, (byte) 0x3);

        consumer = adapter.newConsumer(Short.class, addressStr);
        checkConsumer(connection, address, consumer, (short) 1, (short) 2, (short) 3);

        consumer = adapter.newConsumer(Integer.class, addressStr);
        checkConsumer(connection, address, consumer, 1000, 1001, 1002);

        consumer = adapter.newConsumer(Float.class, addressStr);
        checkConsumer(connection, address, consumer, 1000.5f, 1001.5f, 1002.5f);

        consumer = adapter.newConsumer(String.class, addressStr);
        checkConsumer(connection, address, consumer, "one", "two", "three");*/

        adapter.close();
    }

    /*
     * test PlcConnectionAdapter.newJsonConsumer(address) with write exception
     */
    @Test
    @Category(FastTests.class)
    public void testNewConsumer1Neg() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockField address = new MockField(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        Consumer<?> consumer;

        /*consumer = adapter.newConsumer(String.class, addressStr);
        checkConsumer(2, connection, address, consumer, "one", "two", "three");*/

        adapter.close();
    }

    static <T> void checkConsumer(MockConnection connection, PlcField field, Consumer<T> consumer, Object... values) {
        checkConsumer(0, connection, field, consumer, values);
    }

    private static <T> void checkConsumer(int writeFailureCountTrigger, MockConnection connection, PlcField field, Consumer<T> consumer, Object... values) {
        // verify that a write failure doesn't kill the consumer
        // it logs (not verified) but keeps working for the subsequent writes
        connection.setWriteException(writeFailureCountTrigger, "This is a mock write exception");
        int writeCount = 0;
        Object previousValue = null;
        /*for (Object value : values) {
            @SuppressWarnings("unchecked")
            T tValue = (T) value;
            consumer.accept(tValue);
            Object writtenData = connection.getDataValue(field);
            if (List.class.isAssignableFrom(writtenData.getClass())) {
              @SuppressWarnings("unchecked")
              List<Object> writtenDataList = (List<Object>) writtenData;
                writtenData = writtenDataList.get(0);
            }
            if (writtenData.getClass().isArray()) {
                writtenData = Array.get(writtenData, 0);
            }
            // System.out.println("checkConsumer"+(writeFailureCountTrigger > 0 ? "NEG" : "")+": value:"+value+" writtenData:"+writtenData);
            if (writeFailureCountTrigger <= 0)
                assertThat(value, equalTo(writtenData));
            else {
                if (++writeCount != writeFailureCountTrigger)
                    assertThat(value, equalTo(writtenData));
                else
                    assertThat(previousValue, equalTo(writtenData));
            }
            previousValue = value;
        }*/
    }

    /*
     * test PlcConnectionAdapter.newJsonConsumer(addressFn, valueFn)
     */
    @Test
    @Category(FastTests.class)
    public void testNewConsumer2() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockField address = new MockField(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        Consumer<JsonObject> consumer;

        Function<JsonObject, String> addressFn = t -> t.get("address").getAsString();

        /*consumer = adapter.newJsonConsumer(Boolean.class, addressFn, t -> t.get("value").getAsBoolean());
        checkConsumerJson(connection, address, consumer, true, false);

        consumer = adapter.newJsonConsumer(Byte.class, addressFn, t -> t.get("value").getAsByte());
        checkConsumerJson(connection, address, consumer, (byte) 0x1, (byte) 0x2, (byte) 0x3);

        consumer = adapter.newJsonConsumer(Short.class, addressFn, t -> t.get("value").getAsShort());
        checkConsumerJson(connection, address, consumer, (short) 1, (short) 2, (short) 3);

        consumer = adapter.newJsonConsumer(Integer.class, addressFn, t -> t.get("value").getAsInt());
        checkConsumerJson(connection, address, consumer, 1000, 1001, 1002);

        consumer = adapter.newJsonConsumer(Float.class, addressFn, t -> t.get("value").getAsFloat());
        checkConsumerJson(connection, address, consumer, 1000.5f, 1001.5f, 1002.5f);

        consumer = adapter.newJsonConsumer(String.class, addressFn, t -> t.get("value").getAsString());
        checkConsumerJson(connection, address, consumer, "one", "two", "three");*/

        adapter.close();
    }

    /*
     * test PlcConnectionAdapter.newJsonConsumer(addressFn, valueFn) with write failure
     */
    @Test
    @Category(FastTests.class)
    public void testNewConsumer2Neg() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockField address = new MockField(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        Consumer<JsonObject> consumer;

        Function<JsonObject, String> addressFn = t -> t.get("address").getAsString();

        /*consumer = adapter.newJsonConsumer(String.class, addressFn, t -> t.get("value").getAsString());
        checkConsumerJson(2, connection, address, consumer, "one", "two", "three");*/

        adapter.close();
    }

    static <T> void checkConsumerJson(MockConnection connection, MockField address, Consumer<JsonObject> consumer, Object... values) {
        checkConsumerJson(0, connection, address, consumer, values);
    }

    private static <T> void checkConsumerJson(int writeFailureCountTrigger, MockConnection connection,
                                              MockField field, Consumer<JsonObject> consumer, Object... values) {
        if (writeFailureCountTrigger > 0) {
            connection.setWriteException(writeFailureCountTrigger, "This is a mock write exception");
        }
        int writeCount = 0;
        Object previousValue = null;
        for (Object value : values) {

            // build the JsonObject to consume
            JsonObject jo = new JsonObject();
            jo.addProperty("address", field.getAddress());
            if (value instanceof Boolean) {
                jo.addProperty("value", (Boolean) value);
            } else if (value instanceof Number) {
                jo.addProperty("value", (Number) value);
            } else if (value instanceof String) {
                jo.addProperty("value", (String) value);
            }

            consumer.accept(jo);

            /*Object writtenData = connection.getDataValue(field);
            if (writtenData.getClass().isArray()) {
                Object[] writtenDataArray = (Object[]) writtenData;
                writtenData = Array.get(writtenDataArray, 0);
            }
            if (List.class.isAssignableFrom(writtenData.getClass())) {
                @SuppressWarnings("unchecked")
                List<Object> writtenDataList = (List<Object>) writtenData;
                writtenData = writtenDataList.get(0);
            }
            // System.out.println("checkConsumerJson"+(writeFailureCountTrigger > 0 ? "NEG" : "")+": value:"+value+" writtenData:"+writtenData);
            if (writeFailureCountTrigger <= 0)
                assertThat(value, equalTo(writtenData));
            else {
                if (++writeCount != writeFailureCountTrigger)
                    assertThat(value, equalTo(writtenData));
                else
                    assertThat(previousValue, equalTo(writtenData));
            }*/
            previousValue = value;
        }
    }

}
