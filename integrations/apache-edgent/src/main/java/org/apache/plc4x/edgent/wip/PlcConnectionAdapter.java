package org.apache.plc4x.edgent.wip;

import java.util.Calendar;

import org.apache.edgent.function.Consumer;
import org.apache.edgent.function.Function;
import org.apache.edgent.function.Supplier;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcException;
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
 * PlcConnectionAdapter encapsulates a plc4x {@link PlcConnection}.
 * <p>
 * The idea here is to use PlcConnectionAdapter to enable our Edgent Supplier/Consumer
 * instances to be isolated from the details of / variability of 
 * PlcConnection mgmt and such.
 * <p>
 * A PlcConnectionAdapter is subject to the constraints of the underlying
 * PlcConnection for the device.
 * <p>
 * A single PlcConnectionAdaptor represents a single underlying PlcConnection/connection
 * to a plc device.
 * <p>
 * Multiple PlcConnectionAdaptor instances, hence multiple PlcConnection instances,
 * can be created for a single plc device subject to the underlying device's protocol's
 * support for multiple connections from a single client.
 * <p 
 * A single PlcConnectionAdapter can be used by multiple threads concurrently
 * (e.g., used by multiple PlcFunctions Consumers for {@code Topology.poll()} and/or 
 * multiple Suppliers for {@code TStream.sink()}). 
 * 
 * @see PlcFunctions
 */
public class PlcConnectionAdapter implements AutoCloseable{

  private final static Logger logger = LoggerFactory.getLogger(PlcConnectionAdapter.class);
  
  private String plcConnectionUrl;
  private PlcConnection plcConnection;
  
  /*
   * NOTES:
   * - if we get to the point of the application needing some feedback (possibly control)
   *   of read or write errors, my thinking is to enhance the PlcConnectionAdapter
   *   to enable the app to register an error callback handler or such.
   */
  
  public PlcConnectionAdapter(PlcConnection plcConnection) {
    this.plcConnection = plcConnection;
  }
  
  public PlcConnectionAdapter(String plcConnectionUrl) {
    this.plcConnectionUrl = plcConnectionUrl;
  }
  
  PlcConnection getConnection() throws PlcException {
    synchronized(this) {
      if (plcConnection == null) {
        plcConnection = new PlcDriverManager().getConnection(plcConnectionUrl);
      }
      return plcConnection;
    }
  }

  @Override
  public void close() throws Exception {
    // only close a connection this instance created/connected
    if (plcConnectionUrl != null) {
      if (plcConnection != null)
        plcConnection.close();
    }
  }

  <T> Supplier<T> newSupplier(Class<T> datatype, String addressStr) {
    PlcConnectionAdapter.checkDatatype(datatype);
    return new Supplier<T>() {
      private static final long serialVersionUID = 1L;

      @Override
      public T get() {
        PlcConnection connection = null;
        Address address = null;
        try {
          connection = getConnection();
          address = connection.parseAddress(addressStr);
          PlcReader reader = connection.getReader().get();
          PlcReadRequest<T> readRequest = PlcConnectionAdapter.newPlcReadRequest(datatype, address);
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

  <T> Consumer<T> newConsumer(Class<T> datatype, String addressStr) {
    PlcConnectionAdapter.checkDatatype(datatype);
    return new Consumer<T>() {
      private static final long serialVersionUID = 1L;

      @Override
      public void accept(T arg0) {
        PlcConnection connection = null;
        Address address = null;
        try {
          connection = getConnection();
          address = connection.parseAddress(addressStr);
          PlcWriter writer = connection.getWriter().get();
          PlcWriteRequest<T> writeReq = PlcConnectionAdapter.newPlcWriteRequest(address, arg0);
          writer.write(writeReq).get();
        }
        catch (Exception e) {
          logger.error("writing to plc device {} {} failed", connection, address, e);
        }
      }
      
    };
  }
  
  <T> Consumer<JsonObject> newConsumer(Class<T> datatype, Function<JsonObject,String> addressFn, Function<JsonObject,T> valueFn) {
    PlcConnectionAdapter.checkDatatype(datatype);
    return new Consumer<JsonObject>() {
      private static final long serialVersionUID = 1L;

      @Override
      public void accept(JsonObject jo) {
        PlcConnection connection = null;
        Address address = null;
        try {
          connection = getConnection();
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

  static void checkDatatype(Class<?> cls) {
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
  static <T> PlcWriteRequest<T> newPlcWriteRequest(Address address, T value) {
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
    else if (Calendar.class.isAssignableFrom(cls))
      return (PlcWriteRequest<T>) new CalendarPlcWriteRequest(address, (Calendar)value);
    else
      throw new IllegalArgumentException("Not a legal plc data type: "+cls.getSimpleName());
  }

  @SuppressWarnings("unchecked")
  static <T> PlcReadRequest<T> newPlcReadRequest(Class<T> datatype, Address address) {
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

}
