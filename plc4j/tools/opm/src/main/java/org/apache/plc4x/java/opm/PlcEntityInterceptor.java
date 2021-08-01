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
import org.apache.commons.lang3.Validate;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.temporal.ChronoUnit;
import java.util.Arrays;
import java.util.Map;
import java.util.concurrent.Callable;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

/**
 * Interceptor for dynamic functionality of @{@link PlcEntity}.
 * Basically, its {@link #interceptGetter(Object, Method, Callable, String, PlcDriverManager, AliasRegistry, Map, Map)} method is called for each
 * invocation of a method on a connected @{@link PlcEntity} and does then the dynamic part.
 * <p>
 * For those not too familiar with the JVM's dispatch on can roughly imagine the intercept method being a "regular"
 * method on the "proxied" entity and all parameters of the intercept method could then be access to local fields.
 *
 * @author julian
 */
@SuppressWarnings({"common-java:DuplicatedBlocks", "Duplicates"})
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
     * the {@link #fetchAndSetValueForGetter(Object, Method, PlcDriverManager, String, AliasRegistry, Map)} method.
     * <p>
     * If the field is no getter, then all fields are refreshed by calling {@link #refetchAllFields(Object, PlcDriverManager, String, AliasRegistry, Map)}
     * and then, the method is invoked.
     *
     * @param proxy         Object to intercept
     * @param method        Method that was intercepted
     * @param callable      Callable to call the method after fetching the values
     * @param address       Address of the plc (injected from private field)
     * @param driverManager DriverManager instance to use (injected from private field)
     * @return possible result of the original methods invocation
     * @throws OPMException Problems with plc / proxying
     */
    @SuppressWarnings({"unused", "squid:S00107"})
    @RuntimeType
    public static Object interceptGetter(@This Object proxy, @Origin Method method, @SuperCall Callable<?> callable,
                                         @FieldValue(PlcEntityManager.PLC_ADDRESS_FIELD_NAME) String address,
                                         @FieldValue(PlcEntityManager.DRIVER_MANAGER_FIELD_NAME) PlcDriverManager driverManager,
                                         @FieldValue(PlcEntityManager.ALIAS_REGISTRY) AliasRegistry registry,
                                         @FieldValue(PlcEntityManager.LAST_FETCHED) Map<String, Instant> lastFetched,
                                         @FieldValue(PlcEntityManager.LAST_WRITTEN) Map<String, Instant> lastWritten) throws OPMException {
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

            fetchAndSetValueForGetter(proxy, method, driverManager, address, registry, lastFetched);
            try {
                return callable.call();
            } catch (Exception e) {
                throw new OPMException("Unable to forward invocation " + method.getName() + " on connected PlcEntity", e);
            }
        }

        if (method.getName().startsWith("is") && (method.getReturnType() == boolean.class || method.getReturnType() == Boolean.class)) {
            if (method.getParameterCount() > 0) {
                throw new OPMException("Only getter with no arguments are supported");
            }
            // Fetch single value
            LOGGER.trace("Invoked method {} is boolean flag method, trying to find annotated field and return requested value",
                method.getName());
            fetchAndSetValueForIsGetter(proxy, method, driverManager, address, registry, lastFetched);
            try {
                return callable.call();
            } catch (Exception e) {
                throw new OPMException("Unable to forward invocation " + method.getName() + " on connected PlcEntity", e);
            }
        }

        // Fetch all values, than invoke method
        try {
            LOGGER.trace("Invoked method is no getter, refetch all fields and invoke method {} then", method.getName());
            refetchAllFields(proxy, driverManager, address, registry, lastFetched);
            Object call = callable.call();
            // We write back
            writeAllFields(proxy, driverManager, address, registry, lastWritten);
            return call;
        } catch (Exception e) {
            throw new OPMException("Unable to forward invocation " + method.getName() + " on connected PlcEntity", e);
        }
    }

    @SuppressWarnings({"unused", "squid:S00107"})
    @RuntimeType
    public static Object interceptSetter(@This Object proxy, @Origin Method method, @SuperCall Callable<?> callable,
                                         @FieldValue(PlcEntityManager.PLC_ADDRESS_FIELD_NAME) String address,
                                         @FieldValue(PlcEntityManager.DRIVER_MANAGER_FIELD_NAME) PlcDriverManager driverManager,
                                         @FieldValue(PlcEntityManager.ALIAS_REGISTRY) AliasRegistry registry,
                                         @FieldValue(PlcEntityManager.LAST_FETCHED) Map<String, Instant> lastFetched,
                                         @Argument(0) Object argument) throws OPMException {
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

        if (method.getName().startsWith("set")) {
            if (method.getParameterCount() != 1) {
                throw new OPMException("Only setter with one arguments are supported");
            }
            // Set single value
            LOGGER.trace("Invoked method {} is setter, trying to find annotated field and return requested value",
                method.getName());

            return setValueForSetter(proxy, method, callable, driverManager, address, registry, lastFetched, argument);
        }

        // Fetch all values, than invoke method
        try {
            LOGGER.trace("Invoked method is no getter, refetch all fields and invoke method {} then", method.getName());
            refetchAllFields(proxy, driverManager, address, registry, lastFetched);
            return callable.call();
        } catch (Exception e) {
            throw new OPMException("Unable to forward invocation " + method.getName() + " on connected PlcEntity", e);
        }
    }

    /**
     * Renews all values of all Fields that are annotated with {@link PlcEntity}.
     *
     * @param proxy         Object to refresh the fields on.
     * @param driverManager Driver Manager to use
     * @param registry      AliasRegistry to use
     * @param lastFetched
     * @throws OPMException on various errors.
     */
    @SuppressWarnings("squid:S1141") // Nested try blocks readability is okay, move to other method makes it imho worse
    static void refetchAllFields(Object proxy, PlcDriverManager driverManager, String address, AliasRegistry registry, Map<String, Instant> lastFetched) throws OPMException {
        // Don't log o here as this would cause a second request against a plc so don't touch it, or if you log be aware of that
        Class<?> entityClass = proxy.getClass().getSuperclass();
        LOGGER.trace("Refetching all fields on proxy object of class {}", entityClass);
        PlcEntity plcEntity = entityClass.getAnnotation(PlcEntity.class);
        if (plcEntity == null) {
            throw new OPMException("Non PlcEntity supplied");
        }

        // Check if all fields are valid
        for (Field field : entityClass.getDeclaredFields()) {
            if (field.isAnnotationPresent(PlcField.class)) {
                OpmUtils.getOrResolveAddress(registry, field.getAnnotation(PlcField.class).value());
            }
        }
        try (PlcConnection connection = driverManager.getConnection(address)) {
            // Catch the exception, if no reader present (see below)
            // Build the query
            PlcReadRequest.Builder requestBuilder = connection.readRequestBuilder();

            Arrays.stream(entityClass.getDeclaredFields())
                .filter(field -> field.isAnnotationPresent(PlcField.class))
                .filter(field -> needsToBeSynced(lastFetched, field))
                .forEach(field ->
                    requestBuilder.addItem(
                        getFqn(field),
                        OpmUtils.getOrResolveAddress(registry, field.getAnnotation(PlcField.class).value())
                    )
                );

            PlcReadRequest request = requestBuilder.build();

            LOGGER.trace("Request for refetch of {} was build and is {}", entityClass, request);

            PlcReadResponse response = getPlcReadResponse(request);

            // Fill all requested fields
            for (String fieldName : response.getFieldNames()) {
                // Fill into Cache
                lastFetched.put(fieldName, Instant.now());

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

    static void writeAllFields(Object proxy, PlcDriverManager driverManager, String address, AliasRegistry registry, Map<String, Instant> lastWritten) throws OPMException {
        // Don't log o here as this would cause a second request against a plc so don't touch it, or if you log be aware of that
        Class<?> entityClass = proxy.getClass().getSuperclass();
        LOGGER.trace("Writing all fields on proxy object of class {}", entityClass);
        PlcEntity plcEntity = entityClass.getAnnotation(PlcEntity.class);
        if (plcEntity == null) {
            throw new OPMException("Non PlcEntity supplied");
        }

        // Check if all fields are valid
        for (Field field : entityClass.getDeclaredFields()) {
            if (field.isAnnotationPresent(PlcField.class)) {
                OpmUtils.getOrResolveAddress(registry, field.getAnnotation(PlcField.class).value());
            }
        }
        try (PlcConnection connection = driverManager.getConnection(address)) {
            // Catch the exception, if no reader present (see below)
            // Build the query
            PlcWriteRequest.Builder requestBuilder = connection.writeRequestBuilder();

            Arrays.stream(entityClass.getDeclaredFields())
                .filter(field -> field.isAnnotationPresent(PlcField.class))
                .filter(field -> needsToBeSynced(lastWritten, field))
                .forEach(field ->
                    requestBuilder.addItem(
                        getFqn(field),
                        OpmUtils.getOrResolveAddress(registry, field.getAnnotation(PlcField.class).value()),
                        getFromField(field, proxy)
                    )
                );

            PlcWriteRequest request = requestBuilder.build();

            LOGGER.trace("Request for write of {} was build and is {}", entityClass, request);

            PlcWriteResponse response = getPlcWriteResponse(request);

            // Fill all requested fields
            for (String fieldName : response.getFieldNames()) {
                // Fill into Cache
                lastWritten.put(fieldName, Instant.now());
            }
        } catch (PlcConnectionException e) {
            throw new OPMException("Problem during processing", e);
        } catch (Exception e) {
            throw new OPMException("Unexpected error during processing", e);
        }
    }

    private static Object getFromField(Field field, Object object) {
        try {
            field.setAccessible(true);
            return field.get(object);
        } catch (IllegalAccessException e) {
            throw new PlcRuntimeException(e);
        }
    }

    private static String getFqn(Field field) {
        return field.getDeclaringClass().getName() + "." + field.getName();
    }

    /**
     * Checks if a field needs to be refetched/rewritten, i.e., the cached values are too old.
     */
    private static boolean needsToBeSynced(Map<String, Instant> lastSynced, Field field) {
        Validate.notNull(field);
        long cacheDurationMillis = field.getAnnotation(PlcField.class).cacheDurationMillis();
        if (cacheDurationMillis < 0) {
            return true;
        }
        String fqn = getFqn(field);
        if (lastSynced.containsKey(fqn)) {
            Instant last = lastSynced.get(fqn);
            return Instant.now().minus(cacheDurationMillis, ChronoUnit.MILLIS).isAfter(last);
        }
        return true;
    }

    private static void fetchAndSetValueForIsGetter(Object proxy, Method m, PlcDriverManager driverManager, String address, AliasRegistry registry, Map<String, Instant> lastFetched) throws OPMException {
        fetchAndSetValueForGetter(proxy, m, 2, driverManager, address, registry, lastFetched);
    }

    private static void fetchAndSetValueForGetter(Object proxy, Method m, PlcDriverManager driverManager, String address, AliasRegistry registry, Map<String, Instant> lastFetched) throws OPMException {
        fetchAndSetValueForGetter(proxy, m, 3, driverManager, address, registry, lastFetched);
    }

    private static void fetchAndSetValueForGetter(Object proxy, Method m, int prefixLength, PlcDriverManager driverManager,
                                                  String address, AliasRegistry registry, Map<String, Instant> lastFetched) throws OPMException {
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
            throw new OPMException("Unable to identify field with name '" + variable + "' for call to '" + m.getName() + "'", e);
        }

        // Use Fully qualified Name as field index
        String fqn = getFqn(field);

        // Check if cache is still active
        if (!needsToBeSynced(lastFetched, field)) {
            return;
        }
        try (PlcConnection connection = driverManager.getConnection(address)) {
            // Catch the exception, if no reader present (see below)

            PlcReadRequest request = connection.readRequestBuilder()
                .addItem(fqn, OpmUtils.getOrResolveAddress(registry, annotation.value()))
                .build();

            PlcReadResponse response = getPlcReadResponse(request);

            // Fill into Cache
            lastFetched.put(field.getName(), Instant.now());

            Object value = getTyped(m.getReturnType(), response, fqn);
            setForField(field, proxy, value);
        } catch (ClassCastException e) {
            throw new OPMException("Unable to return response as suitable type", e);
        } catch (Exception e) {
            throw new OPMException("Problem during processing", e);
        }
    }

    private static void setForField(Field field, Object proxy, Object value) {
        try {
            field.setAccessible(true);
            field.set(proxy, value);
        } catch (IllegalAccessException e) {
            throw new PlcRuntimeException(e);
        }
    }

    private static Object setValueForSetter(Object proxy, Method m, Callable<?> callable, PlcDriverManager driverManager,
                                            String address, AliasRegistry registry, Map<String, Instant> lastFetched, Object object) throws OPMException {
        String s = m.getName().substring(3);
        // First char to lower
        String variable = s.substring(0, 1).toLowerCase().concat(s.substring(1));
        LOGGER.trace("Looking for field with name {} after invokation of getter {}", variable, m.getName());
        PlcField annotation;
        Field field;
        try {
            field = m.getDeclaringClass().getDeclaredField(variable);
            annotation = field.getDeclaredAnnotation(PlcField.class);
        } catch (NoSuchFieldException e) {
            throw new OPMException("Unable to identify field with name '" + variable + "' for call to '" + m.getName() + "'", e);
        }

        // Use Fully qualified Name as field index
        String fqn = getFqn(field);

        try (PlcConnection connection = driverManager.getConnection(address)) {
            // Catch the exception, if no reader present (see below)

            PlcWriteRequest request = connection.writeRequestBuilder()
                .addItem(fqn, OpmUtils.getOrResolveAddress(registry, annotation.value()), object)
                .build();

            PlcWriteResponse response = getPlcWriteResponse(request);

            // Fill into Cache
            lastFetched.put(field.getName(), Instant.now());

            LOGGER.debug("getTyped clazz: {}, response: {}, fieldName: {}", m.getParameters()[0].getType(), response, fqn);
            if (response.getResponseCode(fqn) != PlcResponseCode.OK) {
                throw new PlcRuntimeException(String.format("Unable to read specified field '%s', response code was '%s'",
                    fqn, response.getResponseCode(fqn)));
            }
            callable.call();
            return null;
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
     * @throws NoSuchFieldException   If a field is not present in entity
     * @throws IllegalAccessException If a field in the entity cannot be accessed
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

    @SuppressWarnings({"squid:S3776", "squid:MethodCyclomaticComplexity"})
    // Cognitive Complexity not too high, as highly structured
    static Object getTyped(Class<?> clazz, PlcReadResponse response, String sourceFieldName) {
        LOGGER.debug("getTyped clazz: {}, response: {}, fieldName: {}", clazz, response, sourceFieldName);
        if (response.getResponseCode(sourceFieldName) != PlcResponseCode.OK) {
            throw new PlcRuntimeException(String.format("Unable to read specified field '%s', response code was '%s'",
                sourceFieldName, response.getResponseCode(sourceFieldName)));
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
        return getFromFuture(request);
    }

    /**
     * Fetch the request and do appropriate error handling
     *
     * @param request the request to get the exception from
     * @return the response from the exception.
     * @throws OPMException on {@link InterruptedException} or {@link ExecutionException} or {@link TimeoutException}
     */
    public static PlcWriteResponse getPlcWriteResponse(PlcWriteRequest request) throws OPMException {
        return getFromFuture(request);
    }

    @SuppressWarnings("unchecked")
    private static <REQ extends PlcRequest, RES extends PlcResponse> RES getFromFuture(REQ request) throws OPMException {
        try {
            return (RES) request.execute().get(READ_TIMEOUT, TimeUnit.MILLISECONDS);
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
