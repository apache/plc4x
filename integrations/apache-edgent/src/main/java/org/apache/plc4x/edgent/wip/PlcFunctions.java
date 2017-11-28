package org.apache.plc4x.edgent.wip;

import java.util.Calendar;

import org.apache.edgent.function.Consumer;
import org.apache.edgent.function.Function;
import org.apache.edgent.function.Supplier;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.connection.PlcWriter;
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
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.StringPlcReadRequest;
import org.apache.plc4x.java.api.messages.StringPlcWriteRequest;
import org.apache.plc4x.java.api.model.Address;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.gson.JsonObject;

/**
 * WIP - A plc4x Apache Edgent Supplier and Consumer connector factory.
 * <p>
 * TODO:
 * Are there cases where a single logical poll would want to read from 
 * multiple addrs/data (of different types) from a device and bundle the values
 * into a single TStream tuple (e.g., a JsonObject)?  How would we support that?
 * Is there a similar need for writing to multiple addrs/values on a device?
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

  private final static Logger logger = LoggerFactory.getLogger(PlcFunctions.class);

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
   * @param addressStr the plc device address string
   * @param datatype Class<T>
   * @param address memory address on the device
   * @return the {@code Supplier<T>}
   * 
   * @throws IllegalArgumentException if {@code datatype} isn't a plc supported data type
   * 
   * @see org.apache.edgent.topology.Topology#poll(Supplier, long, java.util.concurrent.TimeUnit)
   */
  private static <T> Supplier<T> supplier(PlcConnectionAdapter adapter, String addressStr, Class<T> datatype) {
    checkDatatype(datatype);
    return new Supplier<T>() {
      private static final long serialVersionUID = 1L;

      @Override
      public T get() {
        PlcConnection connection = null;
        Address address = null;
        try {
          connection = adapter.getConnection();
          address = connection.parseAddress(addressStr);
          PlcReader reader = connection.getReader().get();
          PlcReadRequest<T> readRequest = newPlcReadRequest(datatype, address);
          T value = reader.read(readRequest).get().getValue();
          return value;
        }
        catch (Exception e) {
          logger.error("reading from plc device {} {} failed", connection, address, e);
          return null;
        }
      }
      
    };
  }

  public static Supplier<Boolean> booleanSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return supplier(adapter, addressStr, Boolean.class);
  }
  public static Supplier<Byte> byteSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return supplier(adapter, addressStr, Byte.class);
  }
  public static Supplier<Integer> integerSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return supplier(adapter, addressStr, Integer.class);
  }
  public static Supplier<Float> floatSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return supplier(adapter, addressStr, Float.class);
  }
  public static Supplier<String> stringSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return supplier(adapter, addressStr, String.class);
  }
  public static Supplier<Calendar> calendarSupplier(PlcConnectionAdapter adapter, String addressStr) {
    return supplier(adapter, addressStr, Calendar.class);
  }

  private static void checkDatatype(Class<?> cls) {
    if (cls == Boolean.class
        || cls == Byte.class
        || cls == Integer.class
        || cls == Float.class
        || cls == String.class
        || cls == Calendar.class)
      return;
    throw new IllegalArgumentException("Not a legal plc data type: "+cls.getSimpleName());
  }
  
  @SuppressWarnings("unchecked")
  private static <T> PlcReadRequest<T> newPlcReadRequest(Class<T> datatype, Address address) {
    if (datatype == Boolean.class)
      return (PlcReadRequest<T>) new BooleanPlcReadRequest(address);
    else if (datatype == Byte.class)
      return (PlcReadRequest<T>) new BytePlcReadRequest(address);
    else if (datatype == Integer.class)
      return (PlcReadRequest<T>) new IntegerPlcReadRequest(address);
    else if (datatype == Float.class)
      return (PlcReadRequest<T>) new FloatPlcReadRequest(address);
    else if (datatype == String.class)
      return (PlcReadRequest<T>) new StringPlcReadRequest(address);
    else if (datatype == Calendar.class)
      return (PlcReadRequest<T>) new CalendarPlcReadRequest(address);
    else
      throw new IllegalArgumentException("Not a legal plc data type: "+datatype.getSimpleName());
  }
  
  @SuppressWarnings("unchecked")
  private static <T> PlcWriteRequest<T> newPlcWriteRequest(Address address, T value) {
    Class<?> cls = value.getClass();
    if (cls == Boolean.class)
      return (PlcWriteRequest<T>) new BooleanPlcWriteRequest(address, (Boolean)value);
    else if (cls == Byte.class)
      return (PlcWriteRequest<T>) new BytePlcWriteRequest(address, (Byte)value);
    else if (cls == Integer.class)
      return (PlcWriteRequest<T>) new IntegerPlcWriteRequest(address, (Integer)value);
    else if (cls == Float.class)
      return (PlcWriteRequest<T>) new FloatPlcWriteRequest(address, (Float)value);
    else if (cls == String.class)
      return (PlcWriteRequest<T>) new StringPlcWriteRequest(address, (String)value);
    else if (cls == Calendar.class)
      return (PlcWriteRequest<T>) new CalendarPlcWriteRequest(address, (Calendar)value);
    else
      throw new IllegalArgumentException("Not a legal plc data type: "+cls.getSimpleName());
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
  private static <T> Consumer<T> consumer(PlcConnectionAdapter adapter, String addressStr, Class<T> datatype) {
    checkDatatype(datatype);
    return new Consumer<T>() {
      private static final long serialVersionUID = 1L;

      @Override
      public void accept(T arg0) {
        PlcConnection connection = null;
        Address address = null;
        try {
          connection = adapter.getConnection();
          address = connection.parseAddress(addressStr);
          PlcWriter writer = connection.getWriter().get();
          PlcWriteRequest<T> writeReq = newPlcWriteRequest(address, arg0);
          writer.write(writeReq).get();
        }
        catch (Exception e) {
          logger.error("writing to plc device {} {} failed", connection, address, e);
        }
      }
      
    };
  }

  public static Consumer<Boolean> booleanConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return consumer(adapter, addressStr, Boolean.class);
  }
  public static Consumer<Byte> byteConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return consumer(adapter, addressStr, Byte.class);
  }
  public static Consumer<Integer> integerConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return consumer(adapter, addressStr, Integer.class);
  }
  public static Consumer<Float> floatConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return consumer(adapter, addressStr, Float.class);
  }
  public static Consumer<String> stringConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return consumer(adapter, addressStr, String.class);
  }
  public static Consumer<Calendar> calendarConsumer(PlcConnectionAdapter adapter, String addressStr) {
    return consumer(adapter, addressStr, Calendar.class);
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
    return new Consumer<JsonObject>() {
      private static final long serialVersionUID = 1L;

      @Override
      public void accept(JsonObject jo) {
        PlcConnection connection = null;
        Address address = null;
        try {
          connection = adapter.getConnection();
          String addressStr = addressFn.apply(jo);
          address = connection.parseAddress(addressStr);
          T value = valueFn.apply(jo);
          PlcWriter writer = connection.getWriter().get();
          PlcWriteRequest<T> writeReq = newPlcWriteRequest(address, value);
          writer.write(writeReq).get();
        }
        catch (Exception e) {
          logger.error("writing to plc device {} {} failed", connection, address, e);
        }
      }
      
    };
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
