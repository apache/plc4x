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
package org.apache.plc4x.edgent.wip;

import java.util.Calendar;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

import org.apache.edgent.function.Consumer;
import org.apache.edgent.function.Function;
import org.apache.edgent.function.Supplier;
import org.apache.plc4x.edgent.mock.MockAddress;
import org.apache.plc4x.edgent.mock.MockConnection;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.BooleanPlcReadRequest;
import org.apache.plc4x.java.api.messages.BooleanPlcWriteRequest;
import org.apache.plc4x.java.api.messages.BytePlcReadRequest;
import org.apache.plc4x.java.api.messages.BytePlcWriteRequest;
import org.apache.plc4x.java.api.messages.CalendarPlcReadRequest;
import org.apache.plc4x.java.api.messages.CalendarPlcWriteRequest;
import org.apache.plc4x.java.api.messages.FloatPlcReadRequest;
import org.apache.plc4x.java.api.messages.FloatPlcWriteRequest;
import org.apache.plc4x.java.api.messages.IntegerPlcReadRequest;
import org.apache.plc4x.java.api.messages.IntegerPlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.messages.StringPlcReadRequest;
import org.apache.plc4x.java.api.messages.StringPlcWriteRequest;
import org.apache.plc4x.java.api.model.Address;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import com.google.gson.JsonObject;

public class PlcConcectionAdapterTest {
  
  // TODO figure out how to get these run via Eclipse (org.junit.jupiter.api?) and remove this
  // Ah... Junit 5... needs newer Eclipse (newer than neon 1.a)
  public static void main(String[] args) throws Exception {
    PlcConcectionAdapterTest t = new PlcConcectionAdapterTest();
    t.testCtor1();
    t.testCtor2();
    t.testCheckDatatype();
    t.testNewPlcReadRequest();
    t.testNewPlcWriteRequest();
    t.testNewSupplier();
    t.testNewSupplierNeg();
    t.testNewConsumer1();
    t.testNewConsumer1Neg();
    t.testNewConsumer2();
    t.testNewConsumer2Neg();
    System.out.println("SUCCESS");
  }
  
  protected MockConnection getMockConnection() throws PlcConnectionException {
    return (MockConnection) new PlcDriverManager().getConnection("mock://some-cool-url");
  }
  
  /*
   * Test the PlcConnectionAdapter(PlcConnection) ctor, getConnection() and close()
   */
  @Test
  @Tag("fast")
  public void testCtor1() throws Exception {
    MockConnection mockConnection = getMockConnection();
    PlcConnectionAdapter adapter = new PlcConnectionAdapter(mockConnection);
    Assertions.assertSame(mockConnection, adapter.getConnection());
    // and again... multiple adapter.getConnection() returns the same
    Assertions.assertSame(mockConnection, adapter.getConnection());
    adapter.close();
  }
  
  /*
   * Test the PlcConnectionAdapter(url) ctor, getConnection() and close()
   */
  @Test
  @Tag("fast")
  public void testCtor2() throws Exception {
    MockConnection mockConnection = getMockConnection();
    PlcConnectionAdapter adapter = new PlcConnectionAdapter(mockConnection.getUrl());
    MockConnection mockConnection2 = (MockConnection) adapter.getConnection();
    Assertions.assertNotSame(mockConnection, mockConnection2);
    Assertions.assertSame(mockConnection.getUrl(), mockConnection2.getUrl());
    // and again... multiple adapter.getConnection() returns the same
    Assertions.assertSame(mockConnection2, adapter.getConnection());
    adapter.close();
  }
 
  @Test
  @Tag("fast")
  public void testCheckDatatype() throws Exception {
    PlcConnectionAdapter.checkDatatype(Boolean.class);
    PlcConnectionAdapter.checkDatatype(Byte.class);
    PlcConnectionAdapter.checkDatatype(Integer.class);
    PlcConnectionAdapter.checkDatatype(Float.class);
    PlcConnectionAdapter.checkDatatype(String.class);
    PlcConnectionAdapter.checkDatatype(Calendar.class);
    Assertions.assertThrows(IllegalArgumentException.class,
        () -> PlcConnectionAdapter.checkDatatype(Short.class));
    Assertions.assertThrows(IllegalArgumentException.class,
        () -> PlcConnectionAdapter.checkDatatype(Long.class));
    Assertions.assertThrows(IllegalArgumentException.class,
        () -> PlcConnectionAdapter.checkDatatype(Double.class));
  }
  
  private <T> void checkRead(MockConnection connection, PlcReadRequest<T> request, T value) throws InterruptedException, ExecutionException {
    // this is really a tests of our mock tooling but knowing it's behaving as expected
    // will help identify problems in the adapter/supplier/consumer
    connection.setDataValue(request.getAddress(), value);
    
    CompletableFuture<PlcReadResponse<T>> cf = connection.read(request);
    
    Assertions.assertTrue(cf.isDone());
    PlcReadResponse<T> response = cf.get();
    Assertions.assertEquals(value, response.getValue());
  }
  
  private <T> void checkWrite(MockConnection connection, PlcWriteRequest<T> request, T value) throws InterruptedException, ExecutionException {
    // this is really a tests of our mock tooling but knowing it's behaving as expected
    // will help identify problems in the adapter/supplier/consumer
    connection.setDataValue(request.getAddress(), value);
    
    CompletableFuture<PlcWriteResponse<T>> cf = connection.write(request);
    
    Assertions.assertTrue(cf.isDone());
    PlcWriteResponse<T> response = cf.get();
    Assertions.assertNotNull(response);
    Assertions.assertEquals(value, connection.getDataValue(request.getAddress()));
  }
  
  /*
   * Verify the adapter yields the appropriate PlcReadRequest for each type and that it works.
   */
  @SuppressWarnings("unchecked")
  @Test
  @Tag("fast")
  public void testNewPlcReadRequest() throws Exception {
    String addressStr = "MyReadWriteAddress/0";
    MockAddress address = new MockAddress(addressStr);
    PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
    MockConnection connection = (MockConnection) adapter.getConnection();

    PlcReadRequest<?> request;
    
    request = PlcConnectionAdapter.newPlcReadRequest(Boolean.class, address);
    Assertions.assertTrue(request instanceof BooleanPlcReadRequest, "class:"+request.getClass());
    Assertions.assertSame(address, request.getAddress());
    checkRead(connection, (PlcReadRequest<Boolean>)request, true);
    checkRead(connection, (PlcReadRequest<Boolean>)request, false);
    
    request = PlcConnectionAdapter.newPlcReadRequest(Byte.class, address);
    Assertions.assertTrue(request instanceof BytePlcReadRequest, "class:"+request.getClass());
    Assertions.assertSame(address, request.getAddress());
    checkRead(connection, (PlcReadRequest<Byte>)request, (byte)0x13);
    checkRead(connection, (PlcReadRequest<Byte>)request, (byte)0x23);
    
    request = PlcConnectionAdapter.newPlcReadRequest(Integer.class, address);
    Assertions.assertTrue(request instanceof IntegerPlcReadRequest, "class:"+request.getClass());
    Assertions.assertSame(address, request.getAddress());
    checkRead(connection, (PlcReadRequest<Integer>)request, 33);
    checkRead(connection, (PlcReadRequest<Integer>)request, -133);
    
    request = PlcConnectionAdapter.newPlcReadRequest(Float.class, address);
    Assertions.assertTrue(request instanceof FloatPlcReadRequest, "class:"+request.getClass());
    Assertions.assertSame(address, request.getAddress());
    checkRead(connection, (PlcReadRequest<Float>)request, 43.5f);
    checkRead(connection, (PlcReadRequest<Float>)request, -143.5f);
    
    request = PlcConnectionAdapter.newPlcReadRequest(String.class, address);
    Assertions.assertTrue(request instanceof StringPlcReadRequest, "class:"+request.getClass());
    Assertions.assertSame(address, request.getAddress());
    checkRead(connection, (PlcReadRequest<String>)request, "ReadySetGo");
    checkRead(connection, (PlcReadRequest<String>)request, "OneMoreTime");
    
    request = PlcConnectionAdapter.newPlcReadRequest(Calendar.class, address);
    Assertions.assertTrue(request instanceof CalendarPlcReadRequest, "class:"+request.getClass());
    Assertions.assertSame(address, request.getAddress());
    checkRead(connection, (PlcReadRequest<Calendar>)request, Calendar.getInstance());
    
    adapter.close();
  }
  
  
  /*
   * Verify the adapter yields the appropriate PlcWriteRequest for each type and that it works.
   */
  @SuppressWarnings("unchecked")
  @Test
  @Tag("fast")
  public void testNewPlcWriteRequest() throws Exception {
    String addressStr = "MyReadWriteAddress/0";
    MockAddress address = new MockAddress(addressStr);
    PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
    MockConnection connection = (MockConnection) adapter.getConnection();

    PlcWriteRequest<?> request;
    
    request = PlcConnectionAdapter.newPlcWriteRequest(address, true);
    Assertions.assertTrue(request instanceof BooleanPlcWriteRequest, "class:"+request.getClass());
    Assertions.assertSame(address, request.getAddress());
    checkWrite(connection, (PlcWriteRequest<Boolean>)request, true);
    
    request = PlcConnectionAdapter.newPlcWriteRequest(address, (byte)0x113);
    Assertions.assertTrue(request instanceof BytePlcWriteRequest, "class:"+request.getClass());
    Assertions.assertSame(address, request.getAddress());
    checkWrite(connection, (PlcWriteRequest<Byte>)request, (byte)0x113);
    
    request = PlcConnectionAdapter.newPlcWriteRequest(address, 1033);
    Assertions.assertTrue(request instanceof IntegerPlcWriteRequest, "class:"+request.getClass());
    Assertions.assertSame(address, request.getAddress());
    checkWrite(connection, (PlcWriteRequest<Integer>)request, 1033);
    
    request = PlcConnectionAdapter.newPlcWriteRequest(address, 1043.5f);
    Assertions.assertTrue(request instanceof FloatPlcWriteRequest, "class:"+request.getClass());
    Assertions.assertSame(address, request.getAddress());
    checkWrite(connection, (PlcWriteRequest<Float>)request, 1043.5f);
    
    request = PlcConnectionAdapter.newPlcWriteRequest(address, "A written value");
    Assertions.assertTrue(request instanceof StringPlcWriteRequest, "class:"+request.getClass());
    Assertions.assertSame(address, request.getAddress());
    checkWrite(connection, (PlcWriteRequest<String>)request, "A written value");
    
    Calendar calValue = Calendar.getInstance();
    request = PlcConnectionAdapter.newPlcWriteRequest(address, calValue);
    Assertions.assertTrue(request instanceof CalendarPlcWriteRequest, "class:"+request.getClass());
    Assertions.assertSame(address, request.getAddress());
    checkWrite(connection, (PlcWriteRequest<Calendar>)request, calValue);
    
    adapter.close();
  }

  /*
   * test PlcConnectionAdapter.newSupplier
   */
  @SuppressWarnings("unchecked")
  @Test
  @Tag("fast")
  public void testNewSupplier() throws Exception {
    String addressStr = "MyReadWriteAddress/0";
    MockAddress address = new MockAddress(addressStr);
    PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
    MockConnection connection = (MockConnection) adapter.getConnection();

    Supplier<?> supplier;
    
    supplier = adapter.newSupplier(Boolean.class, addressStr);
    Assertions.assertNotSame(supplier, adapter.newSupplier(Boolean.class, addressStr));
    checkSupplier(connection, address, (Supplier<Boolean>)supplier, true, false);
    
    supplier = adapter.newSupplier(Byte.class, addressStr);
    checkSupplier(connection, address, (Supplier<Byte>)supplier, (byte)0x1, (byte)0x2, (byte)0x3);
    
    supplier = adapter.newSupplier(Integer.class, addressStr);
    checkSupplier(connection, address, (Supplier<Integer>)supplier, 1000, 1001, 1002);
    
    supplier = adapter.newSupplier(Float.class, addressStr);
    checkSupplier(connection, address, (Supplier<Float>)supplier, 1000.5f, 1001.5f, 1002.5f);
    
    supplier = adapter.newSupplier(String.class, addressStr);
    checkSupplier(connection, address, (Supplier<String>)supplier, "one", "two", "three");
    
    adapter.close();
  }
  
  /*
   * test PlcConnectionAdapter.newSupplier with read exception
   */
  @SuppressWarnings("unchecked")
  @Test
  @Tag("fast")
  public void testNewSupplierNeg() throws Exception {
    String addressStr = "MyReadWriteAddress/0";
    MockAddress address = new MockAddress(addressStr);
    PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
    MockConnection connection = (MockConnection) adapter.getConnection();

    Supplier<?> supplier;
    
    supplier = adapter.newSupplier(String.class, addressStr);
    checkSupplier(2, connection, address, (Supplier<String>)supplier, "one", "two", "three");
    
    adapter.close();
  }
  
  private <T> void checkSupplier(MockConnection connection, Address address, Supplier<T> supplier, Object ... values) throws Exception {
    checkSupplier(0, connection, address, supplier, values);
  }
  private <T> void checkSupplier(int readFailureCountTrigger, MockConnection connection, Address address, Supplier<T> supplier, Object ... values) throws Exception {
    // verify that a read failure doesn't kill the consumer
    // it logs (not verified) but returns null (as designed) and keeps working for the subsequent reads
    connection.setReadException(readFailureCountTrigger, "This is a mock read exception");
    int readCount = 0;
    for (Object value : values) {
      connection.setDataValue(address, value);
      T readData = supplier.get();
      // System.out.println("checkSupplier"+(readFailureCountTrigger > 0 ? "NEG" : "")+": value:"+value+" readData:"+readData);
      if (readFailureCountTrigger <= 0)
        Assertions.assertEquals(value, readData);
      else {
        if (++readCount != readFailureCountTrigger)
          Assertions.assertEquals(value, readData);
        else
          Assertions.assertNull(readData);
      }
    }
  }

  /*
   * test PlcConnectionAdapter.newConsumer(address)
   */
  @SuppressWarnings("unchecked")
  @Test
  @Tag("fast")
  public void testNewConsumer1() throws Exception {
    String addressStr = "MyReadWriteAddress/0";
    MockAddress address = new MockAddress(addressStr);
    PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
    MockConnection connection = (MockConnection) adapter.getConnection();

    Consumer<?> consumer;
    
    consumer = adapter.newConsumer(Boolean.class, addressStr);
    Assertions.assertNotSame(consumer, adapter.newConsumer(Boolean.class, addressStr));
    checkConsumer(connection, address, (Consumer<Boolean>)consumer, true, false);
    
    consumer = adapter.newConsumer(Byte.class, addressStr);
    checkConsumer(connection, address, (Consumer<Byte>)consumer, (byte)0x1, (byte)0x2, (byte)0x3);
    
    consumer = adapter.newConsumer(Integer.class, addressStr);
    checkConsumer(connection, address, (Consumer<Integer>)consumer, 1000, 1001, 1002);
    
    consumer = adapter.newConsumer(Float.class, addressStr);
    checkConsumer(connection, address, (Consumer<Float>)consumer, 1000.5f, 1001.5f, 1002.5f);
    
    consumer = adapter.newConsumer(String.class, addressStr);
    checkConsumer(connection, address, (Consumer<String>)consumer, "one", "two", "three");
    
    adapter.close();
  }

  /*
   * test PlcConnectionAdapter.newConsumer(address) with write exception
   */
  @SuppressWarnings("unchecked")
  @Test
  @Tag("fast")
  public void testNewConsumer1Neg() throws Exception {
    String addressStr = "MyReadWriteAddress/0";
    MockAddress address = new MockAddress(addressStr);
    PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
    MockConnection connection = (MockConnection) adapter.getConnection();

    Consumer<?> consumer;
    
    consumer = adapter.newConsumer(String.class, addressStr);
    checkConsumer(2, connection, address, (Consumer<String>)consumer, "one", "two", "three");
    
    adapter.close();
  }
  
  private <T> void checkConsumer(MockConnection connection, Address address, Consumer<T> consumer, Object ... values) throws Exception {
    checkConsumer(0, connection, address, consumer, values);
  }
  @SuppressWarnings("unchecked")
  private <T> void checkConsumer(int writeFailureCountTrigger, MockConnection connection, Address address, Consumer<T> consumer, Object ... values) throws Exception {
    // verify that a write failure doesn't kill the consumer
    // it logs (not verified) but keeps working for the subsequent writes
    connection.setWriteException(writeFailureCountTrigger, "This is a mock write exception");
    int writeCount = 0;
    Object previousValue = null;
    for (Object value : values) {
      consumer.accept((T)value);
      T writtenData = (T) connection.getDataValue(address);
      // System.out.println("checkConsumer"+(writeFailureCountTrigger > 0 ? "NEG" : "")+": value:"+value+" writtenData:"+writtenData);
      if (writeFailureCountTrigger <= 0)
        Assertions.assertEquals(value, writtenData);
      else { 
        if (++writeCount != writeFailureCountTrigger)
          Assertions.assertEquals(value, writtenData);
        else
          Assertions.assertEquals(previousValue, writtenData);
      }
      previousValue = value;
    }
  }

  /*
   * test PlcConnectionAdapter.newConsumer(addressFn, valueFn)
   */
  @Test
  @Tag("fast")
  public void testNewConsumer2() throws Exception {
    String addressStr = "MyReadWriteAddress/0";
    MockAddress address = new MockAddress(addressStr);
    PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
    MockConnection connection = (MockConnection) adapter.getConnection();

    Consumer<JsonObject> consumer;
    
    Function<JsonObject,String> addressFn = t -> t.get("address").getAsString(); 
    
    consumer = adapter.newConsumer(Boolean.class, addressFn, t -> t.get("value").getAsBoolean());
    checkConsumerJson(connection, address, consumer, true, false);
    
    consumer = adapter.newConsumer(Byte.class, addressFn, t -> t.get("value").getAsByte());
    checkConsumerJson(connection, address, consumer, (byte)0x1, (byte)0x2, (byte)0x3);
    
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
  @Tag("fast")
  public void testNewConsumer2Neg() throws Exception {
    String addressStr = "MyReadWriteAddress/0";
    MockAddress address = new MockAddress(addressStr);
    PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
    MockConnection connection = (MockConnection) adapter.getConnection();

    Consumer<JsonObject> consumer;
    
    Function<JsonObject,String> addressFn = t -> t.get("address").getAsString(); 
    
    consumer = adapter.newConsumer(String.class, addressFn, t -> t.get("value").getAsString());
    checkConsumerJson(2, connection, address, consumer, "one", "two", "three");
    
    adapter.close();
  }
  
  private <T> void checkConsumerJson(MockConnection connection, MockAddress address, Consumer<JsonObject> consumer, Object ... values) throws Exception {
    checkConsumerJson(0, connection, address, consumer, values);
  }
  private <T> void checkConsumerJson(int writeFailureCountTrigger, MockConnection connection, MockAddress address, Consumer<JsonObject> consumer, Object ... values) throws Exception {
    if (writeFailureCountTrigger > 0)
      connection.setWriteException(writeFailureCountTrigger, "This is a mock write exception");
    int writeCount = 0;
    Object previousValue = null;
    for (Object value : values) {
      
      // build the JsonObject to consume
      JsonObject jo = new JsonObject();
      jo.addProperty("address", address.getAddress());
      if (value instanceof Boolean)
        jo.addProperty("value", (Boolean)value);
      else if (value instanceof Number)
        jo.addProperty("value", (Number)value);
      else if (value instanceof String)
        jo.addProperty("value", (String)value);
      
      consumer.accept(jo);
      
      @SuppressWarnings("unchecked")
      T writtenData = (T) connection.getDataValue(address);
      // System.out.println("checkConsumerJson"+(writeFailureCountTrigger > 0 ? "NEG" : "")+": value:"+value+" writtenData:"+writtenData);
      if (writeFailureCountTrigger <= 0)
        Assertions.assertEquals(value, writtenData);
      else { 
        if (++writeCount != writeFailureCountTrigger)
          Assertions.assertEquals(value, writtenData);
        else
          Assertions.assertEquals(previousValue, writtenData);
      }
      previousValue = value;
    }
  }

}
