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

import java.util.Calendar;

import org.apache.edgent.function.Consumer;
import org.apache.edgent.function.Function;
import org.apache.edgent.function.Supplier;

import com.google.gson.JsonObject;

/**
 * WIP - A plc4x Apache Edgent {@link Supplier} and {@link Consumer} connector factory.
 * <p>
 * TODO:
 * Are there cases where a single logical poll would want to read from 
 * multiple addrs/data (of different types) from a device and bundle the values
 * into a single TStream tuple (e.g., a JsonObject)?  How would we support that?
 * Is there a similar need for writing to multiple addrs/values on a device?
 * Ah... NOTE: plc4c "batch" requests are coming and will help to address this.
 * 
 * <p>
 * Sample use to read plc device data into an Edgent TStream:
 * <pre>{@code
 * PlcConnection plcConnection = new PlcDriverManager().getConnection("s7://192.168.0.1/0/0");
 * plcConnection.connect();
 * PlcConnectionAdapter adapter = new PlcConnectionAdapter(plcConnection));
 * 
 * DirectProvider dp = new DirectProvider();
 * Topology top = dp.newTopology();
 * TStream<Byte> stream = top.poll(PlcFunctions.byteSupplier(adapter, "INPUTS/0", 1, TimeUnit.SECONDS);
 * stream.print();
 * dp.submit(top);
 * }</pre>
 * 
 * Sample use to write Edgent TStream data to a plc device:
 * <pre>{@code
 * PlcConnection plcConnection = new PlcDriverManager().getConnection("s7://192.168.0.1/0/0");
 * plcConnection.connect();
 * PlcConnectionAdapter adapter = new PlcConnectionAdapter(plcConnection);
 * 
 * DirectProvider dp = new DirectProvider();
 * Topology top = dp.newTopology();
 * 
 * TStream<Byte> stream = ...
 * stream.print();
 * TSink<Byte> sink = stream.sink(PlcFunctions.byteConsumer(adapter, "OUTPUTS/0"));
 * 
 * dp.submit(top);
 * }</pre>
 * 
 */
public class PlcFunctions {

  /**
   * Create a new Edgent {@link Supplier} to read data from the 
   * plc device.
   * <p>
   * Every call to the returned {@link Supplier#get()} reads a
   * new data value from the plc device address and connection
   * associated with the {@code PlcConnectionAdapter}.
   * <p>
   * 
   * @param adapter the @{link PlcConnectionAdapter}
   * @param addressStr the plc device address string
   * @return the {@code Supplier<T>}
   * 
   * @see org.apache.edgent.topology.Topology#poll(Supplier, long, java.util.concurrent.TimeUnit)
   */
  public static Supplier<Boolean> booleanSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return adapter.newSupplier(Boolean.class, addressStr);
  }
  public static Supplier<Byte> byteSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return adapter.newSupplier(Byte.class, addressStr);
  }
  public static Supplier<Short> shortSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return adapter.newSupplier(Short.class, addressStr);
  }
  public static Supplier<Integer> integerSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return adapter.newSupplier(Integer.class, addressStr);
  }
  public static Supplier<Float> floatSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return adapter.newSupplier(Float.class, addressStr);
  }
  public static Supplier<String> stringSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return adapter.newSupplier(String.class, addressStr);
  }
  public static Supplier<Calendar> calendarSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return adapter.newSupplier(Calendar.class, addressStr);
  }

  /**
   * Create a new Edgent {@link Consumer} to write data to the 
   * plc device.
   * <p>
   * Every call to the returned {@link Consumer#accept(Object)}
   * writes the value to the the device address and connection
   * associated with the {@code PlcConnectionAdapter}.
   * 
   * @param adapter the @{link PlcConnectionAdapter}
   * @param addressStr the plc device address string
   * @return the {@code Consumer<T>}
   * 
   * @see org.apache.edgent.topology.TStream#sink(Consumer)
   */
  public static Consumer<Boolean> booleanConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return adapter.newConsumer(Boolean.class, addressStr);
  }
  public static Consumer<Byte> byteConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return adapter.newConsumer(Byte.class, addressStr);
  }
  public static Consumer<Short> shortConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return adapter.newConsumer(Short.class, addressStr);
  }
  public static Consumer<Integer> integerConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return adapter.newConsumer(Integer.class, addressStr);
  }
  public static Consumer<Float> floatConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return adapter.newConsumer(Float.class, addressStr);
  }
  public static Consumer<String> stringConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return adapter.newConsumer(String.class, addressStr);
  }
  public static Consumer<Calendar> calendarConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return adapter.newConsumer(Calendar.class, addressStr);
  }

  /**
   * Create a new Edgent {@link Consumer} to write data to the 
   * plc device.
   * <p>
   * TODO: Is it premature to supply this?
   * <p>
   * Every call to the returned {@link Consumer#accept(Object)}
   * <ul>
   * <li>calls {@code addressFn} to get the device address string</li>
   * <li>calls {@code valueFn} to get the {@code T} to write</li>
   * <li>writes the value to the device address using the connection
   * associated with the {@code PlcConnectionAdapter}.</li>
   * </ul>
   * 
   * @param adapter the @{link PlcConnectionAdapter}
   * @param addressFn {@code Function} the returns a device {@code Address} from a {@code JsonObject}
   * @param valueFn {@code Function} the returns a {@code Value} from a {@code JsonObject}
   * @return the {@code Consumer<JsonObject>}
   * 
   * @see org.apache.edgent.topology.TStream#sink(Consumer)
   */
  public static Consumer<JsonObject> booleanConsumer(PlcConnectionAdapter adapter, Function<JsonObject,String> addressFn, Function<JsonObject,Boolean> valueFn) {
    return adapter.newConsumer(Boolean.class, addressFn, valueFn);
  }
  public static Consumer<JsonObject> byteConsumer(PlcConnectionAdapter adapter, Function<JsonObject,String> addressFn, Function<JsonObject,Byte> valueFn) {
    return adapter.newConsumer(Byte.class, addressFn, valueFn);
  }
  public static Consumer<JsonObject> shortConsumer(PlcConnectionAdapter adapter, Function<JsonObject,String> addressFn, Function<JsonObject,Short> valueFn) {
    return adapter.newConsumer(Short.class, addressFn, valueFn);
  }
  public static Consumer<JsonObject> integerConsumer(PlcConnectionAdapter adapter, Function<JsonObject,String> addressFn, Function<JsonObject,Integer> valueFn) {
    return adapter.newConsumer(Integer.class, addressFn, valueFn);
  }
  public static Consumer<JsonObject> floatConsumer(PlcConnectionAdapter adapter, Function<JsonObject,String> addressFn, Function<JsonObject,Float> valueFn) {
    return adapter.newConsumer(Float.class, addressFn, valueFn);
  }
  public static Consumer<JsonObject> stringConsumer(PlcConnectionAdapter adapter, Function<JsonObject,String> addressFn, Function<JsonObject,String> valueFn) {
    return adapter.newConsumer(String.class, addressFn, valueFn);
  }
  public static Consumer<JsonObject> calendarConsumer(PlcConnectionAdapter adapter, Function<JsonObject,String> addressFn, Function<JsonObject,Calendar> valueFn) {
    return adapter.newConsumer(Calendar.class, addressFn, valueFn);
  }
  
}
