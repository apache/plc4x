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

package org.apache.plc4x.java.opm;

import net.bytebuddy.ByteBuddy;
import net.bytebuddy.implementation.MethodDelegation;
import net.bytebuddy.implementation.bind.annotation.*;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.concurrent.Callable;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

import static net.bytebuddy.matcher.ElementMatchers.any;

/**
 * Plc4x equivalent of Jpas EntityManager for implementing Object-Plc-Mapping.
 * This means that calls to a plc can be done by using plain POJOs with Annotations.
 * <p>
 * First, the necessary annotations are {@link PlcEntity} and {@link PlcField}.
 * For a class to be usable as PlcEntity it needs
 * <ul>
 * <li>be non-final (as proxiing has to be used in case of {@link #connect(Class)}</li>
 * <li>a public no args constructor for instanciation</li>
 * <li>Needs to be annotated with {@link PlcEntity} and has a valid value which is the connection string</li>
 * </ul>
 * <p>
 * Basically, the {@link PlcEntityManager} has to operation "modes" represented by the methods {@link #read(Class)} and
 * {@link #connect(Class)}.
 * <p>
 * For a field to get Values from the Plc Injected it needs to be annotated with the {@link PlcField} annotation.
 * The value has to be the plc fields string (which is inserted in the {@link PlcReadRequest}).
 * The connection string is taken from the value of the {@link PlcEntity} annotation on the class.
 * <p>
 * The {@link #read(Class)} method has no direkt equivalent in JPA (as far as I know) as it only returns a "detached"
 * entity. This means it fetches all values from the plc that are annotated wiht the {@link PlcField} annotations.
 * <p>
 * The {@link #connect(Class)} method is more JPA-like as it returns a "connected" entity. This means, that each
 * time one of the getters on the returned entity is called a call is made to the plc (and the field value is changed
 * for this specific field).
 * Furthermore, if a method which is no getter is called, then all {@link PlcField}s are refreshed before doing the call.
 * Thus, all operations on fields that are annotated with {@link PlcField} are always done against the "live" values
 * from the PLC.
 * <p>
 * // TODO Add detach method
 */
public class PlcEntityManager {

    private static final Logger LOGGER = LoggerFactory.getLogger(PlcEntityManager.class);

    private final PlcDriverManager driverManager;

    public PlcEntityManager() {
        this.driverManager = new PlcDriverManager();
    }

    public PlcEntityManager(PlcDriverManager driverManager) {
        this.driverManager = driverManager;
    }

    public <T> T read(Class<T> clazz) throws OPMException {
        PlcEntity annotation = getPlcEntityAndCheckPreconditions(clazz);
        String source = annotation.value();

        try (PlcConnection connection = driverManager.getConnection(source)) {

            if (!connection.readRequestBuilder().isPresent()) {
                throw new OPMException("Unable to get Reader for connection with url '" + source + "'");
            }

            PlcReadRequest.Builder requestBuilder = connection.readRequestBuilder().get();

            // Do the necessary queries for all fields
            // HashMap<ReadRequestItem<?>, Field> requestItems = new HashMap<>();
            for (Field field : clazz.getDeclaredFields()) {
                PlcField fieldAnnotation = field.getAnnotation(PlcField.class);
                if (fieldAnnotation == null) {
                    // Ignore that field
                    continue;
                }
                // Create the suitable Request
                String query = fieldAnnotation.value();
                requestBuilder.addItem(field.getName(), query);
            }

            // Build the request
            PlcReadRequest request;
            try {
                request = requestBuilder.build();
            } catch (PlcInvalidFieldException e) {
                throw new OPMException("Unable to parse one field request", e);
            }

            // Perform the request
            PlcReadResponse response;
            try {
                // TODO: make configurable.
                response = request.execute().get(1_000, TimeUnit.MILLISECONDS);
            } catch (InterruptedException | ExecutionException e) {
                throw new OPMException("Request fetching not able", e);
            } catch (TimeoutException e) {
                throw new OPMException("Timeout during fetching values", e);
            }

            // Construct the Object
            T instance = clazz.getConstructor().newInstance();

            // Fill all requested fields
            for (String fieldName : response.getFieldNames()) {
                setField(clazz, instance, response, fieldName);
            }
            return instance;
        } catch (PlcConnectionException e) {
            throw new OPMException("Unable to get connection with url '" + source + "'", e);
        } catch (Exception e) {
            throw new OPMException("Unable to fetch PlcEntity " + clazz.getName(), e);
        }
    }

    /**
     * Returns a connected proxy.
     *
     * @param clazz clazz to be connected.
     * @param <T>   type of param {@code clazz}.
     * @return a connected entity.
     * @throws OPMException when proxy can't be build.
     */
    public <T> T connect(Class<T> clazz) throws OPMException {
        getPlcEntityAndCheckPreconditions(clazz);
        try {
            // Use Byte Buddy to generate a subclassed proxy that delegates all PlcField Methods
            // to the intercept method
            return new ByteBuddy()
                .subclass(clazz)
                .method(any()).intercept(MethodDelegation.to(this))
                .make()
                .load(Thread.currentThread().getContextClassLoader())
                .getLoaded()
                .getConstructor()
                .newInstance();
        } catch (NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            throw new OPMException("Unable to instantiate Proxy", e);
        }
    }

    private <T> PlcEntity getPlcEntityAndCheckPreconditions(Class<T> clazz) {
        PlcEntity annotation = clazz.getAnnotation(PlcEntity.class);
        if (annotation == null) {
            throw new IllegalArgumentException("Given Class is no Plc Entity, i.e., not annotated with @PlcEntity");
        }
        // Check if default constructor exists
        try {
            clazz.getConstructor();
        } catch (NoSuchMethodException e) {
            throw new IllegalArgumentException("Cannot use PlcEntity without default constructor", e);
        }
        return annotation;
    }

    //------------------------------------------------------------------------------------------------
    //
    //  Methods for interception for the proxy object
    //
    //------------------------------------------------------------------------------------------------

    /**
     * Basic Intersector for all methods on the proxy object.
     * It checks if the invoked method is a getter and if so, only retrieves the requested field, forwarding to
     * the {@link #fetchValueForGetter(Object, Method)} method.
     * <p>
     * If the field is no getter, then all fields are refreshed by calling {@link #refetchAllFields(Object)}
     * and then, the method is invoked.
     *
     * @param o    Object to intercept
     * @param m    Method that was intercepted
     * @param c    Callable to call the method after fetching the values
     * @param that Reference to the proxy object
     * @return possible result of the original methods invocation
     * @throws OPMException Problems with plc / proxying
     */
    @RuntimeType
    // TODO: avoid using single letter params.
    public Object intercept(@This Object o, @Origin Method m, @SuperCall Callable<?> c, @Super Object that) throws OPMException {
        LOGGER.trace("Invoked method {} on connected PlcEntity {}", m.getName(), that);

        if (m.getName().startsWith("get")) {
            // Fetch single value
            LOGGER.trace("Invoked method {} is getter, trying to find annotated field and return requested value",
                m.getName());
            return fetchValueForGetter(that, m);
        }

        if (m.getName().startsWith("is") && (m.getReturnType() == boolean.class || m.getReturnType() == Boolean.class)) {
            // Fetch single value
            LOGGER.trace("Invoked method {} is boolean flag method, trying to find annotated field and return requested value",
                m.getName());
            return fetchValueForGetter(that, m);
        }

        // Fetch all values, than invoke method
        try {
            LOGGER.trace("Invoked method is no getter, refetch all fields and invoke method {} then", m.getName());
            refetchAllFields(o);
            return c.call();
        } catch (Exception e) {
            throw new OPMException("Unbale to forward invokation " + m.getName() + " on connected PlcEntity", e);
        }
    }

    /**
     * Renews all values of all Fields that are annotated with {@link PlcEntity}.
     *
     * @param o Object to refresh he fields on.
     * @throws OPMException
     */
    private void refetchAllFields(Object o) throws OPMException {
        Class<?> superclass = o.getClass().getSuperclass();
        PlcEntity plcEntity = superclass.getAnnotation(PlcEntity.class);

        try (PlcConnection connection = driverManager.getConnection(plcEntity.value())) {
            // Catch the exception, if no reader present (see below)
            // Build the query
            PlcReadRequest.Builder builder = connection.readRequestBuilder().get();
            for (Field field : superclass.getDeclaredFields()) {
                // Check if the field has an annotation
                PlcField plcField = field.getDeclaredAnnotation(PlcField.class);
                if (plcField != null) {
                    LOGGER.trace("Adding field " + field.getName() + " to request as " + plcField.value());
                    builder.addItem(field.getName(), plcField.value());
                }

            }
            PlcReadRequest request = builder.build();

            PlcReadResponse response = getPlcReadResponse(request);

            // Fill all requested fields
            for (String fieldName : response.getFieldNames()) {
                LOGGER.trace("Value for field " + fieldName + " is " + response.getObject(fieldName));
                try {
                    setField(o.getClass().getSuperclass(), o, response, fieldName);
                } catch (NoSuchFieldException | IllegalAccessException e) {
                    throw new PlcRuntimeException(e);
                }
            }
        } catch (Exception e) {
            throw new OPMException("Problem during processing", e);
        }
    }

    // TODO: why isn't o used?
    private Object fetchValueForGetter(Object o, Method m) throws OPMException {
        String s = m.getName().substring(3);
        // First char to lower
        String variable = s.substring(0, 1).toLowerCase().concat(s.substring(1));
        LOGGER.trace("Looking for field with name {} after invokation of getter {}", variable, m.getName());
        PlcField annotation;
        try {
            annotation = m.getDeclaringClass().getDeclaredField(variable).getDeclaredAnnotation(PlcField.class);
        } catch (NoSuchFieldException e) {
            throw new OPMException("Unable to identify field annotated field for call to " + m.getName(), e);
        }
        PlcEntity plcEntity = m.getDeclaringClass().getAnnotation(PlcEntity.class);
        try (PlcConnection connection = driverManager.getConnection(plcEntity.value())) {
            // Catch the exception, if no reader present (see below)

            // Assume to do the query here...
            PlcReadRequest request = connection.readRequestBuilder().orElseThrow(IllegalStateException::new)
                .addItem(m.getName(), annotation.value())
                .build();

            PlcReadResponse response;
            response = getPlcReadResponse(request);

            return getTyped(m.getReturnType(), response, m.getName());
        } catch (ClassCastException e) {
            throw new OPMException("Unable to return response as suitable type", e);
        } catch (Exception e) {
            throw new OPMException("Problem during processing", e);
        }
    }

    /**
     * Tries to set a response Item to a field in the given object.
     * This is one by looking for a field in the class and a response item
     * which is equal to the given fieldName parameter.
     *
     * @param o         Object to set the value on
     * @param response  Response to fetch the response from
     * @param fieldName Name of the field in the object and the response
     * @throws NoSuchFieldException
     * @throws IllegalAccessException
     */
    private void setField(Class<?> clazz, Object o, PlcReadResponse response, String fieldName) throws NoSuchFieldException, IllegalAccessException {
        Field field = clazz.getDeclaredField(fieldName);
        field.setAccessible(true);
        try {
            field.set(o, getTyped(field.getType(), response, fieldName));
        } catch (ClassCastException e) {
            // TODO should we simply fail here?
            LOGGER.warn("Unable to assign return value {} to field {} with type {}", response.getObject(fieldName), fieldName, field.getType(), e);
        }
    }

    private Object getTyped(Class<?> clazz, PlcReadResponse response, String fieldName) {
        if (clazz.isPrimitive()) {
            if (clazz == byte.class) {
                return response.getByte(fieldName);
            } else if (clazz == int.class) {
                return response.getInteger(fieldName);
            } else if (clazz == long.class) {
                return response.getLong(fieldName);
            } else if (clazz == short.class) {
                return response.getShort(fieldName);
            }
            // TODO this should fail on Short, Integer, ... because it always gets a Long
        } else if (clazz.isAssignableFrom(response.getObject(fieldName).getClass())) {
            return response.getObject(fieldName);
        }
        // If nothing matched, throw
        throw new ClassCastException("Unable to return response item " + response.getObject(fieldName) + " as instance of " + clazz);
    }

    /**
     * Fetch the request and do appropriate error handling
     *
     * @param request
     * @return
     * @throws OPMException
     */
    private PlcReadResponse getPlcReadResponse(PlcReadRequest request) throws OPMException {
        PlcReadResponse response;
        try {
            response = request.execute().get();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new OPMException("Exception during execution", e);
        } catch (ExecutionException e) {
            throw new OPMException("Exception during execution", e);
        }
        return response;
    }

}
