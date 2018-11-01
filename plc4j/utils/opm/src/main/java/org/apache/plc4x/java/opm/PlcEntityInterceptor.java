/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.opm;

import net.bytebuddy.implementation.bind.annotation.*;
import org.apache.commons.configuration2.Configuration;
import org.apache.commons.configuration2.SystemConfiguration;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Arrays;
import java.util.concurrent.Callable;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

/**
 * Interceptor for dynamic functionality of @{@link PlcEntity}.
 * Basically, its {@link #intercept(Object, Method, Callable, String, PlcDriverManager)} method is called for each
 * invocation of a method on a connected @{@link PlcEntity} and does then the dynamic part.
 *
 * For those not too familiar with the JVM's dispatch on can roughly imagine the intercept method being a "regular"
 * method on the "proxied" entity and all parameters of the intercept method could then be access to local fields.
 *
 * @author julian
 */
public class PlcEntityInterceptor {

    private static final Logger LOGGER = LoggerFactory.getLogger(PlcEntityInterceptor.class);

    private static final Configuration CONF = new SystemConfiguration();
    private static final long READ_TIMEOUT = CONF.getLong("org.apache.plc4x.java.opm.entity_manager.read_timeout", 1_000);

    private PlcEntityInterceptor() {
        throw new UnsupportedOperationException("This class is not to be instantiated");
    }

    /**
     * Basic Intersector for all methods on the proxy object.
     * It checks if the invoked method is a getter and if so, only retrieves the requested field, forwarding to
     * the {@link #fetchValueForGetter(Method, PlcDriverManager,String)} method.
     * <p>
     * If the field is no getter, then all fields are refreshed by calling {@link #refetchAllFields(Object, PlcDriverManager, String)}
     * and then, the method is invoked.
     *
     * @param proxy    Object to intercept
     * @param method   Method that was intercepted
     * @param callable Callable to call the method after fetching the values
     * @param address  Address of the plc (injected from private field)
     * @param driverManager DriverManager instance to use (injected from private field)
     * @return possible result of the original methods invocation
     * @throws OPMException Problems with plc / proxying
     */
    @SuppressWarnings("unused")
    @RuntimeType
    public static Object intercept(@This Object proxy, @Origin Method method, @SuperCall Callable<?> callable,
           @FieldValue(PlcEntityManager.PLC_ADDRESS_FIELD_NAME) String address,
           @FieldValue(PlcEntityManager.DRIVER_MANAGER_FIELD_NAME) PlcDriverManager driverManager) throws OPMException {
        LOGGER.trace("Invoked method {} on connected PlcEntity {}", method.getName(), method.getDeclaringClass().getName());

        // If "detached" (i.e. _driverManager is null) simply forward the call
        if (driverManager == null) {
            LOGGER.trace("Entity not connected, simply fowarding call");
            try {
                return callable.call();
            } catch (Exception e) {
                throw new OPMException("Exception during forwarding call", e);
            }
        }

        if (method.getName().startsWith("get")) {
            if (method.getParameterCount() > 0) {
                throw new OPMException("Only getter with no arguments are supported");
            }
            // Fetch single value
            LOGGER.trace("Invoked method {} is getter, trying to find annotated field and return requested value",
                method.getName());
            return fetchValueForGetter(method, driverManager, address);
        }

        if (method.getName().startsWith("is") && (method.getReturnType() == boolean.class || method.getReturnType() == Boolean.class)) {
            if (method.getParameterCount() > 0) {
                throw new OPMException("Only getter with no arguments are supported");
            }
            // Fetch single value
            LOGGER.trace("Invoked method {} is boolean flag method, trying to find annotated field and return requested value",
                method.getName());
            return fetchValueForIsGetter(method, driverManager, address);
        }

        // Fetch all values, than invoke method
        try {
            LOGGER.trace("Invoked method is no getter, refetch all fields and invoke method {} then", method.getName());
            refetchAllFields(proxy, driverManager, address);
            return callable.call();
        } catch (Exception e) {
            throw new OPMException("Unable to forward invocation " + method.getName() + " on connected PlcEntity", e);
        }
    }

    /**
     * Renews all values of all Fields that are annotated with {@link PlcEntity}.
     *
     * @param proxy Object to refresh the fields on.
     * @param driverManager
     * @throws OPMException on various errors.
     */
    @SuppressWarnings("squid:S1141") // Nested try blocks readability is okay, move to other method makes it imho worse
    static void refetchAllFields(Object proxy, PlcDriverManager driverManager, String address) throws OPMException {
        // Don't log o here as this would cause a second request against a plc so don't touch it, or if you log be aware of that
        Class<?> entityClass = proxy.getClass().getSuperclass();
        PlcEntity plcEntity = entityClass.getAnnotation(PlcEntity.class);
        if (plcEntity == null) {
            throw new OPMException("Non PlcEntity supplied");
        }

        try (PlcConnection connection = driverManager.getConnection(address)) {
            // Catch the exception, if no reader present (see below)
            // Build the query
            PlcReadRequest.Builder requestBuilder = connection.readRequestBuilder();

            Arrays.stream(entityClass.getDeclaredFields())
                .filter(field -> field.isAnnotationPresent(PlcField.class))
                .forEach(field ->
                    requestBuilder.addItem(
                        field.getDeclaringClass().getName() + "." + field.getName(),
                        field.getAnnotation(PlcField.class).value()
                    )
                );

            PlcReadRequest request = requestBuilder.build();

            PlcReadResponse response = getPlcReadResponse(request);

            // Fill all requested fields
            for (String fieldName : response.getFieldNames()) {
                LOGGER.trace("Value for field {}  is {}", fieldName, response.getObject(fieldName));
                String clazzFieldName = StringUtils.substringAfterLast(fieldName, ".");
                try {
                    setField(entityClass, proxy, response, clazzFieldName, fieldName);
                } catch (NoSuchFieldException | IllegalAccessException e) {
                    throw new PlcRuntimeException(e);
                }
            }
        } catch (PlcConnectionException e) {
            throw new OPMException("Problem during processing", e);
        } catch (Exception e) {
            throw new OPMException("Unexpected error during processing", e);
        }
    }

    private static Object fetchValueForIsGetter(Method m, PlcDriverManager driverManager, String address) throws OPMException {
        return fetchValueForGetter(m, 2, driverManager, address);
    }

    private static Object fetchValueForGetter(Method m, PlcDriverManager driverManager, String address) throws OPMException {
        return fetchValueForGetter(m, 3, driverManager, address);
    }

    private static Object fetchValueForGetter(Method m, int prefixLength, PlcDriverManager driverManager, String address) throws OPMException {
        String s = m.getName().substring(prefixLength);
        // First char to lower
        String variable = s.substring(0, 1).toLowerCase().concat(s.substring(1));
        LOGGER.trace("Looking for field with name {} after invokation of getter {}", variable, m.getName());
        PlcField annotation;
        Field field;
        try {
            field = m.getDeclaringClass().getDeclaredField(variable);
            annotation = field.getDeclaredAnnotation(PlcField.class);
        } catch (NoSuchFieldException e) {
            throw new OPMException("Unable to identify field annotated field for call to " + m.getName(), e);
        }
        try (PlcConnection connection = driverManager.getConnection(address)) {
            // Catch the exception, if no reader present (see below)

            // Use Fully qualified Name as field index
            String fqn = field.getDeclaringClass().getName() + "." + field.getName();

            PlcReadRequest request = connection.readRequestBuilder()
                .addItem(fqn, annotation.value())
                .build();

            PlcReadResponse response = getPlcReadResponse(request);

            return getTyped(m.getReturnType(), response, fqn);
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
     * @param o               Object to set the value on
     * @param response        Response to fetch the response from
     * @param targetFieldName Name of the field in the object
     * @param sourceFieldName Name of the field in the response
     * @throws NoSuchFieldException
     * @throws IllegalAccessException
     */
    static void setField(Class<?> clazz, Object o, PlcReadResponse response, String targetFieldName, String sourceFieldName) throws NoSuchFieldException, IllegalAccessException {
        LOGGER.debug("setField on clazz: {}, Object: {}, response: {}, targetFieldName: {}, sourceFieldName:{} ", clazz, o, response, targetFieldName, sourceFieldName);
        Field field = clazz.getDeclaredField(targetFieldName);
        field.setAccessible(true);
        try {
            field.set(o, getTyped(field.getType(), response, sourceFieldName));
        } catch (ClassCastException e) {
            throw new PlcRuntimeException(String.format("Unable to assign return value %s to field %s with type %s",
                response.getObject(sourceFieldName), targetFieldName, field.getType()), e);
        }
    }

    @SuppressWarnings("squid:S3776") // Cognitive Complexity not too high, as highly structured
    private static Object getTyped(Class<?> clazz, PlcReadResponse response, String sourceFieldName) {
        LOGGER.debug("getTyped clazz: {}, response: {}, fieldName: {}", clazz, response, sourceFieldName);
        if (response.getResponseCode(sourceFieldName) != PlcResponseCode.OK) {
            throw new PlcRuntimeException(String.format("Unable to read specified field %s, response code was %s",
                sourceFieldName, response));
        }
        if (clazz.isPrimitive()) {
            if (clazz == boolean.class) {
                return response.getBoolean(sourceFieldName);
            } else if (clazz == byte.class) {
                return response.getByte(sourceFieldName);
            } else if (clazz == short.class) {
                return response.getShort(sourceFieldName);
            } else if (clazz == int.class) {
                return response.getInteger(sourceFieldName);
            } else if (clazz == long.class) {
                return response.getLong(sourceFieldName);
            }
        }

        if (clazz == Boolean.class) {
            return response.getBoolean(sourceFieldName);
        } else if (clazz == Byte.class) {
            return response.getByte(sourceFieldName);
        } else if (clazz == Short.class) {
            return response.getShort(sourceFieldName);
        } else if (clazz == Integer.class) {
            return response.getInteger(sourceFieldName);
        } else if (clazz == Long.class) {
            return response.getLong(sourceFieldName);
        } else if (clazz == BigInteger.class) {
            return response.getBigInteger(sourceFieldName);
        } else if (clazz == Float.class) {
            return response.getFloat(sourceFieldName);
        } else if (clazz == Double.class) {
            return response.getDouble(sourceFieldName);
        } else if (clazz == BigDecimal.class) {
            return response.getBigDecimal(sourceFieldName);
        } else if (clazz == String.class) {
            return response.getString(sourceFieldName);
        } else if (clazz == LocalTime.class) {
            return response.getTime(sourceFieldName);
        } else if (clazz == LocalDate.class) {
            return response.getDate(sourceFieldName);
        } else if (clazz == LocalDateTime.class) {
            return response.getDateTime(sourceFieldName);
        } else if (clazz == byte[].class) {
            return ArrayUtils.toPrimitive(response.getByteArray(sourceFieldName));
        } else if (clazz == Byte[].class) {
            return response.getByteArray(sourceFieldName);
        }

        // Fallback
        Object responseObject = response.getObject(sourceFieldName);
        if (clazz.isAssignableFrom(responseObject.getClass())) {
            return responseObject;
        }

        // If nothing matched, throw
        throw new ClassCastException("Unable to return response item " + responseObject + "(" + responseObject.getClass() + ") as instance of " + clazz);
    }

    /**
     * Fetch the request and do appropriate error handling
     *
     * @param request the request to get the exception from
     * @return the response from the exception.
     * @throws OPMException on {@link InterruptedException} or {@link ExecutionException} or {@link TimeoutException}
     */
    static PlcReadResponse getPlcReadResponse(PlcReadRequest request) throws OPMException {
        try {
            return request.execute().get(READ_TIMEOUT, TimeUnit.MILLISECONDS);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new OPMException("Exception during execution", e);
        } catch (ExecutionException e) {
            throw new OPMException("Exception during execution", e);
        } catch (TimeoutException e) {
            throw new OPMException("Timeout during fetching values", e);
        }
    }
}
