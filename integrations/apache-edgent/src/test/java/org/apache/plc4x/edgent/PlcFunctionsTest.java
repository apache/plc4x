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

import org.apache.edgent.function.Consumer;
import org.apache.edgent.function.Function;
import org.apache.edgent.function.Supplier;
import org.apache.plc4x.edgent.mock.MockAddress;
import org.apache.plc4x.edgent.mock.MockConnection;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import com.google.gson.JsonObject;

public class PlcFunctionsTest {
  
  // TODO figure out how to get these run via Eclipse (org.junit.jupiter.api?) and remove this
  // Ah... Junit 5... needs newer Eclipse (newer than neon 1.a)
  public static void main(String[] args) throws Exception {
    PlcFunctionsTest t = new PlcFunctionsTest();
    t.testSupplier();
    t.testNewConsumer1();
    t.testNewConsumer2();
    System.out.println("SUCCESS");
  }
  
  protected MockConnection getMockConnection() throws PlcConnectionException {
    return (MockConnection) new PlcDriverManager().getConnection("mock://some-cool-url");
  }

  /*
   * test PlcConnectionAdapter.newSupplier
   */
  @SuppressWarnings("unchecked")
  @Test
  @Tag("fast")
  public void testSupplier() throws Exception {
    String addressStr = "MyReadWriteAddress/0";
    MockAddress address = new MockAddress(addressStr);
    PlcConnectionAdapter adapter = new PlcConnectionAdapter(getMockConnection());
    MockConnection connection = (MockConnection) adapter.getConnection();

    Supplier supplier;
    
    supplier = PlcFunctions.booleanSupplier(adapter, addressStr);
    PlcConnectionAdapterTest.checkSupplier(connection, address, (Supplier<Boolean>)supplier, true, false);
    
    supplier = PlcFunctions.byteSupplier(adapter, addressStr);
    PlcConnectionAdapterTest.checkSupplier(connection, address, (Supplier<Byte>)supplier, (byte)0x1, (byte)0x2, (byte)0x3);

    supplier = PlcFunctions.shortSupplier(adapter, addressStr);
    PlcConnectionAdapterTest.checkSupplier(connection, address, (Supplier<Short>)supplier, (short)1, (short)2, (short)3);

    supplier = PlcFunctions.integerSupplier(adapter, addressStr);
    PlcConnectionAdapterTest.checkSupplier(connection, address, (Supplier<Integer>)supplier, 1000, 1001, 1002);
    
    supplier = PlcFunctions.floatSupplier(adapter, addressStr);
    PlcConnectionAdapterTest.checkSupplier(connection, address, (Supplier<Float>)supplier, 1000.5f, 1001.5f, 1002.5f);
    
    supplier = PlcFunctions.stringSupplier(adapter, addressStr);
    PlcConnectionAdapterTest.checkSupplier(connection, address, (Supplier<String>)supplier, "one", "two", "three");
    
    adapter.close();
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
    
    consumer = PlcFunctions.booleanConsumer(adapter, addressStr);
    PlcConnectionAdapterTest.checkConsumer(connection, address, consumer, true, false);
    
    consumer = PlcFunctions.byteConsumer(adapter, addressStr);
    PlcConnectionAdapterTest.checkConsumer(connection, address, consumer, (byte)0x1, (byte)0x2, (byte)0x3);
    
    consumer = PlcFunctions.shortConsumer(adapter, addressStr);
    PlcConnectionAdapterTest.checkConsumer(connection, address, consumer, (short)1, (short)2, (short)3);

    consumer = PlcFunctions.integerConsumer(adapter, addressStr);
    PlcConnectionAdapterTest.checkConsumer(connection, address, consumer, 1000, 1001, 1002);
    
    consumer = PlcFunctions.floatConsumer(adapter, addressStr);
    PlcConnectionAdapterTest.checkConsumer(connection, address, consumer, 1000.5f, 1001.5f, 1002.5f);
    
    consumer = PlcFunctions.stringConsumer(adapter, addressStr);
    PlcConnectionAdapterTest.checkConsumer(connection, address, consumer, "one", "two", "three");
    
    adapter.close();
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
    
    consumer = PlcFunctions.booleanConsumer(adapter, addressFn, t -> t.get("value").getAsBoolean());
    PlcConnectionAdapterTest.checkConsumerJson(connection, address, consumer, true, false);
    
    consumer = PlcFunctions.byteConsumer(adapter, addressFn, t -> t.get("value").getAsByte());
    PlcConnectionAdapterTest.checkConsumerJson(connection, address, consumer, (byte)0x1, (byte)0x2, (byte)0x3);

    consumer = PlcFunctions.shortConsumer(adapter, addressFn, t -> t.get("value").getAsShort());
    PlcConnectionAdapterTest.checkConsumerJson(connection, address, consumer, (short)1, (short)2, (short)3);

    consumer = PlcFunctions.integerConsumer(adapter, addressFn, t -> t.get("value").getAsInt());
    PlcConnectionAdapterTest.checkConsumerJson(connection, address, consumer, 1000, 1001, 1002);
    
    consumer = PlcFunctions.floatConsumer(adapter, addressFn, t -> t.get("value").getAsFloat());
    PlcConnectionAdapterTest.checkConsumerJson(connection, address, consumer, 1000.5f, 1001.5f, 1002.5f);
    
    consumer = PlcFunctions.stringConsumer(adapter, addressFn, t -> t.get("value").getAsString());
    PlcConnectionAdapterTest.checkConsumerJson(connection, address, consumer, "one", "two", "three");
    
    adapter.close();
  }

}
