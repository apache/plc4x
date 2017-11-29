package org.apache.plc4x.edgent.wip;

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
   * new {@code valuetype} from the plc device address and connection
   * associated with the {@code PlcConnectionAdapter}
   * and returns the valuetype's contained value of type {@code T}.
   * <p>
   * 
   * @param adapter the @{link PlcConnectionAdapter}
   * @param datatype Class<T>
   * @param addressStr the plc device address string
   * @param address memory address on the device
   * @return the {@code Supplier<T>}
   * 
   * @throws IllegalArgumentException if {@code datatype} isn't a plc supported data type
   * 
   * @see org.apache.edgent.topology.Topology#poll(Supplier, long, java.util.concurrent.TimeUnit)
   */
  private static <T> Supplier<T> supplier(PlcConnectionAdapter adapter, Class<T> datatype, String addressStr) {
    return adapter.newSupplier(datatype, addressStr);
  }

  public static Supplier<Boolean> booleanSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return supplier(adapter, Boolean.class, addressStr);
  }
  public static Supplier<Byte> byteSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return supplier(adapter, Byte.class, addressStr);
  }
  public static Supplier<Integer> integerSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return supplier(adapter, Integer.class, addressStr);
  }
  public static Supplier<Float> floatSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return supplier(adapter, Float.class, addressStr);
  }
  public static Supplier<String> stringSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return supplier(adapter, String.class, addressStr);
  }
  public static Supplier<Calendar> calendarSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return supplier(adapter, Calendar.class, addressStr);
  }

  /**
   * Create a new Edgent {@link Consumer} to write data to the 
   * plc device.
   * <p>
   * Every call to the returned {@link Consumer<T>#accept(T)} 
   * converts the {@code T} argument to a {@code valuetype} 
   * and writes the value to the
   * the device address and connection
   * associated with the {@code PlcStaticAddressAdapter}.
   * <p>
   * Every call to the returned {@link Consumer<T>#accept(T)}
   * <ul>
   * <li>converts the {@code T} to a value of type {@code valuetype}</li>
   * <li>writes the value to the device address using the connection
   * associated with the {@code PlcStaticAddressAdapter}.</li>
   * </ul>
   * 
   * @param adapter the @{link PlcStaticAddressAdapter}
   * @param addressStr the plc device address string
   * @param valuetype a {@link Value} type (e.g., IntegerValue)
   * @param address memory address on the device
   * @return the {@code Consumer<T>}
   * 
   * @see org.apache.edgent.topology.TStream#sink(Consumer)
   */
  private static <T> Consumer<T> consumer(PlcConnectionAdapter adapter, Class<T> datatype, String addressStr) {
    return adapter.newConsumer(datatype, addressStr);
  }

  public static Consumer<Boolean> booleanConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return consumer(adapter, Boolean.class, addressStr);
  }
  public static Consumer<Byte> byteConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return consumer(adapter, Byte.class, addressStr);
  }
  public static Consumer<Integer> integerConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return consumer(adapter, Integer.class, addressStr);
  }
  public static Consumer<Float> floatConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return consumer(adapter, Float.class, addressStr);
  }
  public static Consumer<String> stringConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return consumer(adapter, String.class, addressStr);
  }
  public static Consumer<Calendar> calendarConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return consumer(adapter, Calendar.class, addressStr);
  }

  /**
   * Create a new Edgent {@link Consumer} to write data to the 
   * plc device.
   * <p>
   * TODO: Is it premature to supply this?
   * <p>
   * Every call to the returned {@link Consumer<JsonObject>#accept(JsonObject)}
   * <ul>
   * <li>calls {@code addressFn} to get the device address string</li>
   * <li>calls {@code valueFn} to get the {@code T} to write</li>
   * <li>converts the {@code T} to a value of type {@code valuetype}</li>
   * <li>writes the value to the device address using the connection
   * associated with the {@code PlcConnectionAdapter}.</li>
   * </ul>
   * 
   * @param adapter the @{link PlcConnectionAdapter}
   * @param valuetype a {@link Value} type (e.g., IntegerValue)
   * @param addressFn {@code Function} the returns a device {@code Address} from a {@code JsonObject}
   * @param valueFn {@code Function} the returns a {@code Value} from a {@code JsonObject}
   * @return the {@code Consumer<JsonObject>}
   * 
   * @see org.apache.edgent.topology.TStream#sink(Consumer)
   */
  private static <T> Consumer<JsonObject> consumer(PlcConnectionAdapter adapter, Class<T> datatype, Function<JsonObject,String> addressFn, Function<JsonObject,T> valueFn) {
    return adapter.newConsumer(datatype, addressFn, valueFn);
  }

  public static Consumer<JsonObject> booleanConsumer(PlcConnectionAdapter adapter, Function<JsonObject,String> addressFn, Function<JsonObject,Boolean> valueFn) {
    return consumer(adapter, Boolean.class, addressFn, valueFn);
  }
  public static Consumer<JsonObject> byteConsumer(PlcConnectionAdapter adapter, Function<JsonObject,String> addressFn, Function<JsonObject,Byte> valueFn) {
    return consumer(adapter, Byte.class, addressFn, valueFn);
  }
  public static Consumer<JsonObject> integerConsumer(PlcConnectionAdapter adapter, Function<JsonObject,String> addressFn, Function<JsonObject,Integer> valueFn) {
    return consumer(adapter, Integer.class, addressFn, valueFn);
  }
  public static Consumer<JsonObject> floatConsumer(PlcConnectionAdapter adapter, Function<JsonObject,String> addressFn, Function<JsonObject,Float> valueFn) {
    return consumer(adapter, Float.class, addressFn, valueFn);
  }
  public static Consumer<JsonObject> stringConsumer(PlcConnectionAdapter adapter, Function<JsonObject,String> addressFn, Function<JsonObject,String> valueFn) {
    return consumer(adapter, String.class, addressFn, valueFn);
  }
  public static Consumer<JsonObject> calendarConsumer(PlcConnectionAdapter adapter, Function<JsonObject,String> addressFn, Function<JsonObject,Calendar> valueFn) {
    return consumer(adapter, Calendar.class, addressFn, valueFn);
  }
  
}
