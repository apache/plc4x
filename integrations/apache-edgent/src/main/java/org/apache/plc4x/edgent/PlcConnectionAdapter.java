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
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcWriteRequest;
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
 * <p>
 * A single PlcConnectionAdapter can be used by multiple threads concurrently
 * (e.g., used by multiple PlcFunctions Consumers for {@code Topology.poll()} and/or
 * multiple Suppliers for {@code TStream.sink()}).
 *
 * @see PlcFunctions
 */
public class PlcConnectionAdapter implements AutoCloseable {

    private static final Logger logger = LoggerFactory.getLogger(PlcConnectionAdapter.class);

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
        synchronized (this) {
            if (plcConnection == null) {
                plcConnection = new PlcDriverManager().getConnection(plcConnectionUrl);
            }
            return plcConnection;
        }
    }

    @Override
    public void close() throws Exception {
        // only close a connection this instance created/connected
        if (plcConnectionUrl != null && plcConnection != null) {
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
                    TypeSafePlcReadRequest<T> readRequest = PlcConnectionAdapter.newPlcReadRequest(datatype, address);
                    return reader.read(readRequest).get().getResponseItem().get().getValues().get(0);
                } catch (Exception e) {
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
                    PlcWriteRequest writeReq = PlcConnectionAdapter.newPlcWriteRequest(address, arg0);
                    writer.write(writeReq).get();
                } catch (Exception e) {
                    logger.error("writing to plc device {} {} failed", connection, address, e);
                }
            }

        };
    }

    <T> Consumer<JsonObject> newConsumer(Class<T> datatype, Function<JsonObject, String> addressFn, Function<JsonObject, T> valueFn) {
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
                    PlcWriteRequest writeReq = newPlcWriteRequest(address, value);
                    writer.write(writeReq).get();
                } catch (Exception e) {
                    logger.error("writing to plc device {} {} failed", connection, address, e);
                }
            }

        };
    }

    static void checkDatatype(Class<?> cls) {
        if (cls == Boolean.class
            || cls == Byte.class
            || cls == Short.class
            || cls == Integer.class
            || cls == Float.class
            || cls == String.class
            || cls == Calendar.class)
            return;
        throw new IllegalArgumentException("Not a legal plc data type: " + cls.getSimpleName());
    }

    @SuppressWarnings("unchecked")
    static <T> TypeSafePlcWriteRequest<T> newPlcWriteRequest(Address address, T value) {
        Class<T> cls = (Class<T>) value.getClass();
        return new TypeSafePlcWriteRequest<>(cls, address, value);
    }

    static <T> TypeSafePlcReadRequest<T> newPlcReadRequest(Class<T> datatype, Address address) {
        return new TypeSafePlcReadRequest<>(datatype, address);
    }

}
