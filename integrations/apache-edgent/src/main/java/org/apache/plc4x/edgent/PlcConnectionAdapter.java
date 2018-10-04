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
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcClientDatatype;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;

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

    private static final String FIELD_NAME = "default";

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

    public PlcReadRequest.Builder readRequestBuilder() throws PlcException {
        return getConnection().readRequestBuilder().orElseThrow(
            () -> new PlcException("This connection doesn't support reading"));
    }

    Supplier<PlcReadResponse> newSupplier(PlcReadRequest readRequest) {
        return new Supplier<PlcReadResponse>() {
            private static final long serialVersionUID = 1L;

            @Override
            public PlcReadResponse get() {
                PlcConnection connection = null;
                try {
                    connection = getConnection();
                    return readRequest.execute().get();
                } catch (Exception e) {
                    logger.error("reading from plc device {} {} failed", connection, readRequest, e);
                    return null;
                }
            }
        };
    }

    <T> Supplier<T> newSupplier(Class<T> genericDatatype, PlcClientDatatype clientDatatype, String fieldQuery) {
        // satisfy sonar's "Reduce number of anonymous class lines" code smell
        return new MySupplier<>(genericDatatype, clientDatatype, fieldQuery);
    }

    private class MySupplier<T> implements Supplier<T> {

        private static final long serialVersionUID = 1L;

        private Class<T> genericDatatype;
        private PlcClientDatatype clientDatatype;
        private String fieldQuery;

        MySupplier(Class<T> genericDatatype, PlcClientDatatype clientDatatype, String fieldQuery) {
            this.genericDatatype = genericDatatype;
            this.clientDatatype = clientDatatype;
            this.fieldQuery = fieldQuery;
        }

        @Override
        @SuppressWarnings("unchecked")
        public T get() {
            PlcConnection connection = null;
            PlcField field = null;
            try {
                connection = getConnection();
                PlcReadRequest readRequest = connection.readRequestBuilder().orElseThrow(() -> new PlcException("This connection doesn't support reading")).addItem(FIELD_NAME, fieldQuery).build();
                PlcReadResponse readResponse = readRequest.execute().get();
                Object value = null;
                switch (clientDatatype) {
                    case BYTE:
                        value = readResponse.getByte(FIELD_NAME);
                        break;
                    case SHORT:
                        value = readResponse.getShort(FIELD_NAME);
                        break;
                    case INTEGER:
                        value = readResponse.getInteger(FIELD_NAME);
                        break;
                    case LONG:
                        value = readResponse.getLong(FIELD_NAME);
                        break;
                    case FLOAT:
                        value = readResponse.getFloat(FIELD_NAME);
                        break;
                    case DOUBLE:
                        value = readResponse.getDouble(FIELD_NAME);
                        break;
                    case STRING:
                        value = readResponse.getString(FIELD_NAME);
                        break;
                    case TIME:
                        value = readResponse.getTime(FIELD_NAME);
                        break;
                    case DATE:
                        value = readResponse.getDate(FIELD_NAME);
                        break;
                    case DATE_TIME:
                        value = readResponse.getDateTime(FIELD_NAME);
                        break;
                }
                if (value != null) {
                    if (genericDatatype.isAssignableFrom(value.getClass())) {
                        return (T) value;
                    } else {
                        logger.error("types don't match {} should be of type {}", value.getClass(), genericDatatype);
                    }
                }
            } catch (Exception e) {
                logger.error("reading from plc device {} {} failed", connection, field, e);
            }
            return null;
        }
    }

    <T> Consumer<T> newJsonConsumer(Class<T> genericDatatype, PlcClientDatatype clientDatatype, String fieldQuery) {
        return new ObjectConsumer<>(genericDatatype, clientDatatype, fieldQuery);
    }

    <T> Consumer<JsonObject> newJsonConsumer(PlcClientDatatype clientDatatype, Function<JsonObject, String> fieldQueryFn, Function<JsonObject, T> fieldValueFn) {
        return new JsonConsumer<>(clientDatatype, fieldQueryFn, fieldValueFn);
    }

    private abstract class BaseConsumer<T> implements Consumer<T> {

        protected void write(PlcClientDatatype clientDatatype, String fieldQuery, Object fieldValue) {
            PlcConnection connection = null;
            try {
                connection = getConnection();
                PlcWriteRequest.Builder builder = connection.writeRequestBuilder().orElseThrow(() -> new PlcException("This connection doesn't support writing"));
                PlcWriteRequest writeRequest = builder.build();
                addItem(builder, clientDatatype, fieldQuery, fieldValue);
                writeRequest.execute().get();
            } catch (Exception e) {
                logger.error("writing to plc device {} {} failed", connection, fieldQuery, e);
            }

        }

        private void addItem(PlcWriteRequest.Builder builder,
                             PlcClientDatatype clientDatatype, String fieldQuery, Object fieldValue) {
            switch (clientDatatype) {
                case BYTE:
                    if (fieldValue instanceof Byte) {
                        builder.addItem(FIELD_NAME, fieldQuery, (Byte) fieldValue);
                    }
                    break;
                case SHORT:
                    if (fieldValue instanceof Short) {
                        builder.addItem(FIELD_NAME, fieldQuery, (Short) fieldValue);
                    }
                    break;
                case INTEGER:
                    if (fieldValue instanceof Integer) {
                        builder.addItem(FIELD_NAME, fieldQuery, (Integer) fieldValue);
                    }
                    break;
                case LONG:
                    if (fieldValue instanceof Long) {
                        builder.addItem(FIELD_NAME, fieldQuery, (Long) fieldValue);
                    }
                    break;
                case FLOAT:
                    if (fieldValue instanceof Float) {
                        builder.addItem(FIELD_NAME, fieldQuery, (Float) fieldValue);
                    }
                    break;
                case DOUBLE:
                    if (fieldValue instanceof Double) {
                        builder.addItem(FIELD_NAME, fieldQuery, (Double) fieldValue);
                    }
                    break;
                case STRING:
                    if (fieldValue instanceof String) {
                        builder.addItem(FIELD_NAME, fieldQuery, (String) fieldValue);
                    }
                    break;
                case TIME:
                    if (fieldValue instanceof LocalTime) {
                        builder.addItem(FIELD_NAME, fieldQuery, (LocalTime) fieldValue);
                    }
                    break;
                case DATE:
                    if (fieldValue instanceof LocalDate) {
                        builder.addItem(FIELD_NAME, fieldQuery, (LocalDate) fieldValue);
                    }
                    break;
                case DATE_TIME:
                    if (fieldValue instanceof LocalDateTime) {
                        builder.addItem(FIELD_NAME, fieldQuery, (LocalDateTime) fieldValue);
                    }
                    break;
            }
        }
    }

    private class ObjectConsumer<T> extends BaseConsumer<T> {
        private static final long serialVersionUID = 1L;

        private PlcClientDatatype clientDatatype;
        private String fieldQuery;

        ObjectConsumer(Class<T> genericDatatype, PlcClientDatatype clientDatatype, String fieldQuery) {
            this.clientDatatype = clientDatatype;
            this.fieldQuery = fieldQuery;
        }

        @Override
        public void accept(Object fieldValue) {
            write(clientDatatype, fieldQuery, fieldValue);
        }
    }

    private class JsonConsumer<T> extends BaseConsumer<JsonObject> {
        private static final long serialVersionUID = 1L;

        private PlcClientDatatype clientDatatype;
        private Function<JsonObject, String> fieldQueryFn;
        private Function<JsonObject, T> fieldValueFn;

        JsonConsumer(PlcClientDatatype clientDatatype, Function<JsonObject, String> fieldQueryFn, Function<JsonObject, T> fieldValueFn) {
            this.clientDatatype = clientDatatype;
            this.fieldQueryFn = fieldQueryFn;
            this.fieldValueFn = fieldValueFn;
        }

        @Override
        public void accept(JsonObject jsonObject) {
            String fieldQuery = fieldQueryFn.apply(jsonObject);
            Object fieldValue = fieldValueFn.apply(jsonObject);
            write(clientDatatype, fieldQuery, fieldValue);
        }
    }

}
