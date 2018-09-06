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
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcClientDatatype;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;

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
 * <p>
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
 */
public class PlcFunctions {

    private PlcFunctions() {
        throw new IllegalStateException("Utility class!");
    }

    /**
     * Create a new Edgent {@link Supplier} to read data from the
     * plc device.
     * <p>
     * Every call to the returned {@link Supplier#get()} reads a
     * new data value from the plc device address and connection
     * associated with the {@code PlcConnectionAdapter}.
     * <p>
     *
     * @param adapter    the @{link PlcConnectionAdapter}
     * @param addressStr the plc device address string
     * @return the {@code Supplier<T>}
     * <p>
     * // TODO: No need to import the Topology module for just this comment ... I think
     * //see org.apache.edgent.topology.Topology#poll(Supplier, long, java.util.concurrent.TimeUnit)
     */
    /*public static Supplier<byte[]> rawSupplier(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newSupplier(byte[].class, PlcClientDatatype.RAW, addressStr);
    }

    public static Supplier<Object> objectSupplier(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newSupplier(Object.class, PlcClientDatatype.OBJECT, addressStr);
    }*/

    public static Supplier<Boolean> booleanSupplier(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newSupplier(Boolean.class, PlcClientDatatype.BOOLEAN, addressStr);
    }

    public static Supplier<Byte> byteSupplier(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newSupplier(Byte.class, PlcClientDatatype.BYTE, addressStr);
    }

    public static Supplier<Short> shortSupplier(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newSupplier(Short.class, PlcClientDatatype.SHORT, addressStr);
    }

    public static Supplier<Integer> integerSupplier(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newSupplier(Integer.class, PlcClientDatatype.INTEGER, addressStr);
    }

    public static Supplier<Long> longSupplier(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newSupplier(Long.class, PlcClientDatatype.LONG, addressStr);
    }

    public static Supplier<Float> floatSupplier(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newSupplier(Float.class, PlcClientDatatype.FLOAT, addressStr);
    }

    public static Supplier<Double> doubleSupplier(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newSupplier(Double.class, PlcClientDatatype.DOUBLE, addressStr);
    }

    public static Supplier<String> stringSupplier(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newSupplier(String.class, PlcClientDatatype.STRING, addressStr);
    }

    public static Supplier<LocalTime> timeSupplier(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newSupplier(LocalTime.class, PlcClientDatatype.TIME, addressStr);
    }

    public static Supplier<LocalDate> dateSupplier(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newSupplier(LocalDate.class, PlcClientDatatype.DATE, addressStr);
    }

    public static Supplier<LocalDateTime> dateTimeSupplier(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newSupplier(LocalDateTime.class, PlcClientDatatype.DATE_TIME, addressStr);
    }

    public static Supplier<PlcReadResponse> batchSupplier(PlcConnectionAdapter adapter, PlcReadRequest readRequest) {
        return adapter.newSupplier(readRequest);
    }

    /**
     * Create a new Edgent {@link Consumer} to write data to the
     * plc device.
     * <p>
     * Every call to the returned {@link Consumer#accept(Object)}
     * writes the value to the the device address and connection
     * associated with the {@code PlcConnectionAdapter}.
     *
     * @param adapter    the @{link PlcConnectionAdapter}
     * @param addressStr the plc device address string
     * @return the {@code Consumer<T>}
     * <p>
     * // TODO: No need to import the Topology module for just this comment ... I think
     * //see org.apache.edgent.topology.TStream#sink(Consumer)
     */
    /*public static Consumer<byte[]> rawConsumer(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newJsonConsumer(byte[].class, PlcClientDatatype.RAW, addressStr);
    }

    public static Consumer<Object> objectConsumer(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newJsonConsumer(Object.class, PlcClientDatatype.OBJECT, addressStr);
    }*/

    public static Consumer<Boolean> booleanConsumer(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newJsonConsumer(Boolean.class, PlcClientDatatype.BOOLEAN, addressStr);
    }

    public static Consumer<Byte> byteConsumer(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newJsonConsumer(Byte.class, PlcClientDatatype.BYTE, addressStr);
    }

    public static Consumer<Short> shortConsumer(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newJsonConsumer(Short.class, PlcClientDatatype.SHORT, addressStr);
    }

    public static Consumer<Integer> integerConsumer(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newJsonConsumer(Integer.class, PlcClientDatatype.INTEGER, addressStr);
    }

    public static Consumer<Long> longConsumer(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newJsonConsumer(Long.class, PlcClientDatatype.LONG, addressStr);
    }

    public static Consumer<Float> floatConsumer(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newJsonConsumer(Float.class, PlcClientDatatype.FLOAT, addressStr);
    }

    public static Consumer<Double> doubleConsumer(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newJsonConsumer(Double.class, PlcClientDatatype.DOUBLE, addressStr);
    }

    public static Consumer<String> stringConsumer(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newJsonConsumer(String.class, PlcClientDatatype.STRING, addressStr);
    }

    public static Consumer<LocalTime> timeConsumer(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newJsonConsumer(LocalTime.class, PlcClientDatatype.TIME, addressStr);
    }

    public static Consumer<LocalDate> dateConsumer(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newJsonConsumer(LocalDate.class, PlcClientDatatype.DATE, addressStr);
    }

    public static Consumer<LocalDateTime> dateTimeConsumer(PlcConnectionAdapter adapter, String addressStr) {
        return adapter.newJsonConsumer(LocalDateTime.class, PlcClientDatatype.DATE_TIME, addressStr);
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
     * @param adapter   the @{link PlcConnectionAdapter}
     * @param addressFn {@code Function} the returns a device {@code PlcField} from a {@code JsonObject}
     * @param valueFn   {@code Function} the returns a {@code Value} from a {@code JsonObject}
     * @return the {@code Consumer<JsonObject>}
     * <p>
     * // TODO: No need to import the Topology module for just this comment ... I think
     * //see org.apache.edgent.topology.TStream#sink(Consumer)
     */
    public static Consumer<JsonObject> booleanConsumer(PlcConnectionAdapter adapter, Function<JsonObject, String> addressFn, Function<JsonObject, Boolean> valueFn) {
        return adapter.newJsonConsumer(PlcClientDatatype.BOOLEAN, addressFn, valueFn);
    }

    public static Consumer<JsonObject> byteConsumer(PlcConnectionAdapter adapter, Function<JsonObject, String> addressFn, Function<JsonObject, Byte> valueFn) {
        return adapter.newJsonConsumer(PlcClientDatatype.BYTE, addressFn, valueFn);
    }

    public static Consumer<JsonObject> shortConsumer(PlcConnectionAdapter adapter, Function<JsonObject, String> addressFn, Function<JsonObject, Short> valueFn) {
        return adapter.newJsonConsumer(PlcClientDatatype.SHORT, addressFn, valueFn);
    }

    public static Consumer<JsonObject> integerConsumer(PlcConnectionAdapter adapter, Function<JsonObject, String> addressFn, Function<JsonObject, Integer> valueFn) {
        return adapter.newJsonConsumer(PlcClientDatatype.INTEGER, addressFn, valueFn);
    }

    public static Consumer<JsonObject> longConsumer(PlcConnectionAdapter adapter, Function<JsonObject, String> addressFn, Function<JsonObject, Long> valueFn) {
        return adapter.newJsonConsumer(PlcClientDatatype.LONG, addressFn, valueFn);
    }

    public static Consumer<JsonObject> floatConsumer(PlcConnectionAdapter adapter, Function<JsonObject, String> addressFn, Function<JsonObject, Float> valueFn) {
        return adapter.newJsonConsumer(PlcClientDatatype.FLOAT, addressFn, valueFn);
    }

    public static Consumer<JsonObject> doubleConsumer(PlcConnectionAdapter adapter, Function<JsonObject, String> addressFn, Function<JsonObject, Double> valueFn) {
        return adapter.newJsonConsumer(PlcClientDatatype.DOUBLE, addressFn, valueFn);
    }

    public static Consumer<JsonObject> stringConsumer(PlcConnectionAdapter adapter, Function<JsonObject, String> addressFn, Function<JsonObject, String> valueFn) {
        return adapter.newJsonConsumer(PlcClientDatatype.STRING, addressFn, valueFn);
    }

    public static Consumer<JsonObject> timeConsumer(PlcConnectionAdapter adapter, Function<JsonObject, String> addressFn, Function<JsonObject, LocalTime> valueFn) {
        return adapter.newJsonConsumer(PlcClientDatatype.TIME, addressFn, valueFn);
    }

    public static Consumer<JsonObject> dateConsumer(PlcConnectionAdapter adapter, Function<JsonObject, String> addressFn, Function<JsonObject, LocalDate> valueFn) {
        return adapter.newJsonConsumer(PlcClientDatatype.DATE, addressFn, valueFn);
    }

    public static Consumer<JsonObject> dateTimeConsumer(PlcConnectionAdapter adapter, Function<JsonObject, String> addressFn, Function<JsonObject, LocalDateTime> valueFn) {
        return adapter.newJsonConsumer(PlcClientDatatype.DATE_TIME, addressFn, valueFn);
    }

}
