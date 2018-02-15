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
import org.apache.plc4x.edgent.mock.MockAddress;
import org.apache.plc4x.edgent.mock.MockConnection;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.api.messages.items.WriteRequestItem;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadResponse;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcWriteRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcWriteResponse;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.test.FastTests;
import org.assertj.core.api.Assertions;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import java.lang.reflect.Array;
import java.util.Calendar;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

import static org.assertj.core.api.Assertions.assertThatThrownBy;

public class PlcConnectionAdapterTest {

    protected MockConnection getMockConnection() throws PlcConnectionException {
        return (MockConnection) new PlcDriverManager().getConnection("mock://some-cool-url");
    }

    /*
     * Test the PlcConnectionAdapter(PlcConnection) ctor, getConnection() and close()
     */
    @Test
    @Category(FastTests.class)
    public void testCtor1() throws Exception {
        MockConnection mockConnection = getMockConnection();
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(mockConnection);
        Assertions.assertThat(mockConnection).isSameAs(adapter.getConnection());
        // and again... multiple adapter.getConnection() returns the same
        Assertions.assertThat(mockConnection).isSameAs(adapter.getConnection());
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
        Assertions.assertThat(mockConnection).isNotSameAs(mockConnection2);
        Assertions.assertThat(mockConnection.getUrl()).isSameAs(mockConnection2.getUrl());
        // and again... multiple adapter.getConnection() returns the same
        Assertions.assertThat(mockConnection2).isSameAs(adapter.getConnection());
        adapter.close();
    }

    @Test
    @Category(FastTests.class)
    public void testCheckDatatype() throws Exception {
        PlcConnectionAdapter.checkDatatype(Boolean.class);
        PlcConnectionAdapter.checkDatatype(Byte.class);
        PlcConnectionAdapter.checkDatatype(Short.class);
        PlcConnectionAdapter.checkDatatype(Integer.class);
        PlcConnectionAdapter.checkDatatype(Float.class);
        PlcConnectionAdapter.checkDatatype(String.class);
        PlcConnectionAdapter.checkDatatype(Calendar.class);
        assertThatThrownBy(() ->
            PlcConnectionAdapter.checkDatatype(Long.class))
            .isInstanceOf(IllegalArgumentException.class);
        assertThatThrownBy(() ->
            PlcConnectionAdapter.checkDatatype(Double.class))
            .isInstanceOf(IllegalArgumentException.class);
    }

    private <T> void checkRead(MockConnection connection, TypeSafePlcReadRequest<T> request, T value) throws InterruptedException, ExecutionException {
        // this is really a tests of our mock tooling but knowing it's behaving as expected
        // will help identify problems in the adapter/supplier/consumer
        connection.setDataValue(request.getCheckedReadRequestItems().get(0).getAddress(), value);

        CompletableFuture<TypeSafePlcReadResponse<T>> cf = connection.read(request);

        Assertions.assertThat(cf.isDone()).isTrue();
        TypeSafePlcReadResponse<T> response = cf.get();
        Assertions.assertThat(value).isEqualTo(response.getResponseItems().get(0).getValues().get(0));
    }

    @SuppressWarnings("unchecked")
    private <T> void checkWrite(MockConnection connection, TypeSafePlcWriteRequest<T> request, T value) throws InterruptedException, ExecutionException {
        // this is really a tests of our mock tooling but knowing it's behaving as expected
        // will help identify problems in the adapter/supplier/consumer
        connection.setDataValue(request.getRequestItems().get(0).getAddress(), value);

        CompletableFuture<TypeSafePlcWriteResponse<T>> cf = connection.write(request);

        Assertions.assertThat(cf.isDone()).isTrue();
        PlcWriteResponse response = cf.get();
        Assertions.assertThat(response).isNotNull();
        T writtenData = (T) connection.getDataValue(request.getRequestItems().get(0).getAddress());
        if (writtenData.getClass().isArray()) {
            writtenData = (T) Array.get(writtenData, 0);
        }
        if (List.class.isAssignableFrom(writtenData.getClass())) {
            writtenData = (T) ((List) writtenData).get(0);
        }
        Assertions.assertThat(value).isEqualTo(writtenData);
    }

    /*
     * Verify the adapter yields the appropriate PlcReadRequest for each type and that it works.
     */
    @SuppressWarnings("unchecked")
    @Test
    @Category(FastTests.class)
    public void testNewPlcReadRequest() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockAddress address = new MockAddress(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        {
            TypeSafePlcReadRequest<Boolean> request = PlcConnectionAdapter.newPlcReadRequest(Boolean.class, address);
            ReadRequestItem<Boolean> requestItem = request.getCheckedReadRequestItems().get(0);
            Class<Boolean> dataType = requestItem.getDatatype();
            Assertions.assertThat(dataType).isEqualTo(Boolean.class);
            Assertions.assertThat(address).isSameAs(requestItem.getAddress());
            checkRead(connection, request, true);
            checkRead(connection, request, false);
        }
        {
            TypeSafePlcReadRequest<Byte> request = PlcConnectionAdapter.newPlcReadRequest(Byte.class, address);
            ReadRequestItem<Byte> requestItem = request.getCheckedReadRequestItems().get(0);
            Class<Byte> dataType = requestItem.getDatatype();
            Assertions.assertThat(dataType).isEqualTo(Byte.class);
            Assertions.assertThat(address).isSameAs(requestItem.getAddress());
            checkRead(connection, request, (byte) 0x13);
            checkRead(connection, request, (byte) 0x23);
        }
        {
            TypeSafePlcReadRequest<Short> request = PlcConnectionAdapter.newPlcReadRequest(Short.class, address);
            ReadRequestItem<Short> requestItem = request.getCheckedReadRequestItems().get(0);
            Class<Short> dataType = requestItem.getDatatype();
            Assertions.assertThat(dataType).isEqualTo(Short.class);
            Assertions.assertThat(address).isSameAs(requestItem.getAddress());
            checkRead(connection, request, (short) 13);
            checkRead(connection, request, (short) 23);
        }
        {
            TypeSafePlcReadRequest<Integer> request = PlcConnectionAdapter.newPlcReadRequest(Integer.class, address);
            ReadRequestItem<Integer> requestItem = request.getCheckedReadRequestItems().get(0);
            Class<Integer> dataType = requestItem.getDatatype();
            Assertions.assertThat(dataType).isEqualTo(Integer.class);
            Assertions.assertThat(address).isSameAs(requestItem.getAddress());
            checkRead(connection, request, 33);
            checkRead(connection, request, -133);
        }
        {
            TypeSafePlcReadRequest<Float> request = PlcConnectionAdapter.newPlcReadRequest(Float.class, address);
            ReadRequestItem<Float> requestItem = request.getCheckedReadRequestItems().get(0);
            Class<Float> dataType = requestItem.getDatatype();
            Assertions.assertThat(dataType).isEqualTo(Float.class);
            Assertions.assertThat(address).isSameAs(requestItem.getAddress());
            checkRead(connection, request, 43.5f);
            checkRead(connection, request, -143.5f);
        }
        {
            TypeSafePlcReadRequest<String> request = PlcConnectionAdapter.newPlcReadRequest(String.class, address);
            ReadRequestItem<String> requestItem = request.getCheckedReadRequestItems().get(0);
            Class<String> dataType = requestItem.getDatatype();
            Assertions.assertThat(dataType).isEqualTo(String.class);
            Assertions.assertThat(address).isSameAs(requestItem.getAddress());
            checkRead(connection, request, "ReadySetGo");
            checkRead(connection, request, "OneMoreTime");
        }
        {
            TypeSafePlcReadRequest<Calendar> request = PlcConnectionAdapter.newPlcReadRequest(Calendar.class, address);
            ReadRequestItem<Calendar> requestItem = request.getCheckedReadRequestItems().get(0);
            Class<Calendar> dataType = requestItem.getDatatype();
            Assertions.assertThat(dataType).isEqualTo(Calendar.class);
            Assertions.assertThat(address).isSameAs(requestItem.getAddress());
            checkRead(connection, request, Calendar.getInstance());
        }
        adapter.close();
    }


    /*
     * Verify the adapter yields the appropriate PlcWriteRequest for each type and that it works.
     */
    @SuppressWarnings("unchecked")
    @Test
    @Category(FastTests.class)
    public void testNewPlcWriteRequest() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockAddress address = new MockAddress(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        {
            TypeSafePlcWriteRequest<Boolean> request = PlcConnectionAdapter.newPlcWriteRequest(address, true);
            WriteRequestItem<Boolean> requestItem = request.getCheckedRequestItems().get(0);
            Class<Boolean> dataType = requestItem.getDatatype();
            Assertions.assertThat(Boolean.class).isAssignableFrom(dataType);
            Assertions.assertThat(address).isSameAs(requestItem.getAddress());
            checkWrite(connection, request, true);
        }
        {
            TypeSafePlcWriteRequest<Byte> request = PlcConnectionAdapter.newPlcWriteRequest(address, (byte) 0x113);
            WriteRequestItem<Byte> requestItem = request.getCheckedRequestItems().get(0);
            Class<Byte> dataType = requestItem.getDatatype();
            Assertions.assertThat(Byte.class).isAssignableFrom(dataType);
            Assertions.assertThat(address).isSameAs(requestItem.getAddress());
            checkWrite(connection, request, (byte) 0x113);
        }
        {
            TypeSafePlcWriteRequest<Short> request = PlcConnectionAdapter.newPlcWriteRequest(address, (short) 113);
            WriteRequestItem<Short> requestItem = request.getCheckedRequestItems().get(0);
            Class<Short> dataType = requestItem.getDatatype();
            Assertions.assertThat(Short.class).isAssignableFrom(dataType);
            Assertions.assertThat(address).isSameAs(requestItem.getAddress());
            checkWrite(connection, request, (short) 113);
        }
        {
            TypeSafePlcWriteRequest<Integer> request = PlcConnectionAdapter.newPlcWriteRequest(address, 1033);
            WriteRequestItem<Integer> requestItem = request.getCheckedRequestItems().get(0);
            Class<Integer> dataType = requestItem.getDatatype();
            Assertions.assertThat(Integer.class).isAssignableFrom(dataType);
            Assertions.assertThat(address).isSameAs(requestItem.getAddress());
            checkWrite(connection, request, 1033);
        }
        {
            TypeSafePlcWriteRequest<Float> request = PlcConnectionAdapter.newPlcWriteRequest(address, 1043.5f);
            WriteRequestItem<Float> requestItem = request.getCheckedRequestItems().get(0);
            Class<Float> dataType = requestItem.getDatatype();
            Assertions.assertThat(Float.class).isAssignableFrom(dataType);
            Assertions.assertThat(address).isSameAs(requestItem.getAddress());
            checkWrite(connection, request, 1043.5f);
        }
        {
            TypeSafePlcWriteRequest<String> request = PlcConnectionAdapter.newPlcWriteRequest(address, "A written value");
            WriteRequestItem<String> requestItem = request.getCheckedRequestItems().get(0);
            Class<String> dataType = requestItem.getDatatype();
            Assertions.assertThat(String.class).isAssignableFrom(dataType);
            Assertions.assertThat(address).isSameAs(requestItem.getAddress());
            checkWrite(connection, request, "A written value");
        }
        {
            Calendar calValue = Calendar.getInstance();
            TypeSafePlcWriteRequest<Calendar> request = PlcConnectionAdapter.newPlcWriteRequest(address, calValue);
            WriteRequestItem<Calendar> requestItem = request.getCheckedRequestItems().get(0);
            Class<Calendar> dataType = requestItem.getDatatype();
            Assertions.assertThat(Calendar.class).isAssignableFrom(dataType);
            Assertions.assertThat(address).isSameAs(requestItem.getAddress());
            checkWrite(connection, request, calValue);
        }
        adapter.close();
    }

    /*
     * test PlcConnectionAdapter.newSupplier
     */
    @SuppressWarnings("unchecked")
    @Test
    @Category(FastTests.class)
    public void testNewSupplier() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockAddress address = new MockAddress(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        {
            Supplier<Boolean> supplier = adapter.newSupplier(Boolean.class, addressStr);
            Assertions.assertThat(supplier).isNotSameAs(adapter.newSupplier(Boolean.class, addressStr));
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
        }
        adapter.close();
    }

    /*
     * test PlcConnectionAdapter.newSupplier with read exception
     */
    @SuppressWarnings("unchecked")
    @Test
    @Category(FastTests.class)
    public void testNewSupplierNeg() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockAddress address = new MockAddress(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        Supplier<String> supplier = adapter.newSupplier(String.class, addressStr);
        checkSupplier(2, connection, address, supplier, "one", "two", "three");

        adapter.close();
    }

    static <T> void checkSupplier(MockConnection connection, Address address, Supplier<T> supplier, Object... values) throws Exception {
        checkSupplier(0, connection, address, supplier, values);
    }

    private static <T> void checkSupplier(int readFailureCountTrigger, MockConnection connection, Address address, Supplier<T> supplier, Object... values) throws Exception {
        // verify that a read failure doesn't kill the consumer
        // it logs (not verified) but returns null (as designed) and keeps working for the subsequent reads
        connection.setReadException(readFailureCountTrigger, "This is a mock read exception");
        int readCount = 0;
        for (Object value : values) {
            connection.setDataValue(address, value);
            T readData = supplier.get();
            // System.out.println("checkSupplier"+(readFailureCountTrigger > 0 ? "NEG" : "")+": value:"+value+" readData:"+readData);
            if (readFailureCountTrigger <= 0)
                Assertions.assertThat(value).isEqualTo(readData);
            else {
                if (++readCount != readFailureCountTrigger)
                    Assertions.assertThat(value).isEqualTo(readData);
                else
                    Assertions.assertThat(readData).isNull();
            }
        }
    }

    /*
     * test PlcConnectionAdapter.newConsumer(address)
     */
    @SuppressWarnings("unchecked")
    @Test
    @Category(FastTests.class)
    public void testNewConsumer1() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockAddress address = new MockAddress(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        Consumer<?> consumer;

        consumer = adapter.newConsumer(Boolean.class, addressStr);
        Assertions.assertThat(consumer).isNotSameAs(adapter.newConsumer(Boolean.class, addressStr));
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
        checkConsumer(connection, address, consumer, "one", "two", "three");

        adapter.close();
    }

    /*
     * test PlcConnectionAdapter.newConsumer(address) with write exception
     */
    @SuppressWarnings("unchecked")
    @Test
    @Category(FastTests.class)
    public void testNewConsumer1Neg() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockAddress address = new MockAddress(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        Consumer<?> consumer;

        consumer = adapter.newConsumer(String.class, addressStr);
        checkConsumer(2, connection, address, consumer, "one", "two", "three");

        adapter.close();
    }

    static <T> void checkConsumer(MockConnection connection, Address address, Consumer<T> consumer, Object... values) throws Exception {
        checkConsumer(0, connection, address, consumer, values);
    }

    @SuppressWarnings("unchecked")
    private static <T> void checkConsumer(int writeFailureCountTrigger, MockConnection connection, Address address, Consumer<T> consumer, Object... values) throws Exception {
        // verify that a write failure doesn't kill the consumer
        // it logs (not verified) but keeps working for the subsequent writes
        connection.setWriteException(writeFailureCountTrigger, "This is a mock write exception");
        int writeCount = 0;
        Object previousValue = null;
        for (Object value : values) {
            consumer.accept((T) value);
            T writtenData = (T) connection.getDataValue(address);
            if (List.class.isAssignableFrom(writtenData.getClass())) {
                writtenData = (T) ((List) writtenData).get(0);
            }
            if (writtenData.getClass().isArray()) {
                writtenData = (T) Array.get(writtenData, 0);
            }
            // System.out.println("checkConsumer"+(writeFailureCountTrigger > 0 ? "NEG" : "")+": value:"+value+" writtenData:"+writtenData);
            if (writeFailureCountTrigger <= 0)
                Assertions.assertThat(value).isEqualTo(writtenData);
            else {
                if (++writeCount != writeFailureCountTrigger)
                    Assertions.assertThat(value).isEqualTo(writtenData);
                else
                    Assertions.assertThat(previousValue).isEqualTo(writtenData);
            }
            previousValue = value;
        }
    }

    /*
     * test PlcConnectionAdapter.newConsumer(addressFn, valueFn)
     */
    @Test
    @Category(FastTests.class)
    public void testNewConsumer2() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockAddress address = new MockAddress(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        Consumer<JsonObject> consumer;

        Function<JsonObject, String> addressFn = t -> t.get("address").getAsString();

        consumer = adapter.newConsumer(Boolean.class, addressFn, t -> t.get("value").getAsBoolean());
        checkConsumerJson(connection, address, consumer, true, false);

        consumer = adapter.newConsumer(Byte.class, addressFn, t -> t.get("value").getAsByte());
        checkConsumerJson(connection, address, consumer, (byte) 0x1, (byte) 0x2, (byte) 0x3);

        consumer = adapter.newConsumer(Short.class, addressFn, t -> t.get("value").getAsShort());
        checkConsumerJson(connection, address, consumer, (short) 1, (short) 2, (short) 3);

        consumer = adapter.newConsumer(Integer.class, addressFn, t -> t.get("value").getAsInt());
        checkConsumerJson(connection, address, consumer, 1000, 1001, 1002);

        consumer = adapter.newConsumer(Float.class, addressFn, t -> t.get("value").getAsFloat());
        checkConsumerJson(connection, address, consumer, 1000.5f, 1001.5f, 1002.5f);

        consumer = adapter.newConsumer(String.class, addressFn, t -> t.get("value").getAsString());
        checkConsumerJson(connection, address, consumer, "one", "two", "three");

        adapter.close();
    }

    /*
     * test PlcConnectionAdapter.newConsumer(addressFn, valueFn) with write failure
     */
    @Test
    @Category(FastTests.class)
    public void testNewConsumer2Neg() throws Exception {
        String addressStr = "MyReadWriteAddress/0";
        MockAddress address = new MockAddress(addressStr);
        PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
        MockConnection connection = (MockConnection) adapter.getConnection();

        Consumer<JsonObject> consumer;

        Function<JsonObject, String> addressFn = t -> t.get("address").getAsString();

        consumer = adapter.newConsumer(String.class, addressFn, t -> t.get("value").getAsString());
        checkConsumerJson(2, connection, address, consumer, "one", "two", "three");

        adapter.close();
    }

    static <T> void checkConsumerJson(MockConnection connection, MockAddress address, Consumer<JsonObject> consumer, Object... values) throws Exception {
        checkConsumerJson(0, connection, address, consumer, values);
    }

    private static <T> void checkConsumerJson(int writeFailureCountTrigger, MockConnection connection, MockAddress address, Consumer<JsonObject> consumer, Object... values) throws Exception {
        if (writeFailureCountTrigger > 0)
            connection.setWriteException(writeFailureCountTrigger, "This is a mock write exception");
        int writeCount = 0;
        Object previousValue = null;
        for (Object value : values) {

            // build the JsonObject to consume
            JsonObject jo = new JsonObject();
            jo.addProperty("address", address.getAddress());
            if (value instanceof Boolean)
                jo.addProperty("value", (Boolean) value);
            else if (value instanceof Number)
                jo.addProperty("value", (Number) value);
            else if (value instanceof String)
                jo.addProperty("value", (String) value);

            consumer.accept(jo);

            @SuppressWarnings("unchecked")
            T writtenData = (T) connection.getDataValue(address);
            if (writtenData.getClass().isArray()) {
                writtenData = (T) Array.get(writtenData, 0);
            }
            if (List.class.isAssignableFrom(writtenData.getClass())) {
                writtenData = (T) ((List) writtenData).get(0);
            }
            // System.out.println("checkConsumerJson"+(writeFailureCountTrigger > 0 ? "NEG" : "")+": value:"+value+" writtenData:"+writtenData);
            if (writeFailureCountTrigger <= 0)
                Assertions.assertThat(value).isEqualTo(writtenData);
            else {
                if (++writeCount != writeFailureCountTrigger)
                    Assertions.assertThat(value).isEqualTo(writtenData);
                else
                    Assertions.assertThat(previousValue).isEqualTo(writtenData);
            }
            previousValue = value;
        }
    }

}
