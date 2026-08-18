/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
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
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionFactory;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Array;
import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.lang.reflect.ParameterizedType;
import java.lang.reflect.Type;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.temporal.ChronoUnit;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collection;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.Callable;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

/**
 * Interceptor for dynamic functionality of @{@link PlcEntity}.
 * Basically, its {@link #interceptGetter(Object, Method, Callable, String, PlcConnectionFactory, AliasRegistry, Map, Map)} method is called for each
 * invocation of a method on a connected @{@link PlcEntity} and does then the dynamic part.
 * <p>
 * For those not too familiar with the JVM's dispatch on can roughly imagine the intercept method being a "regular"
 * method on the "proxied" entity and all parameters of the intercept method could then be access to local tags.
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
     * Basic Interceptor for all methods on the proxy object.
     * It checks if the invoked method is a getter and if so, only retrieves the requested tag, forwarding to
     * the {@link #fetchAndSetValueForGetter(Object, Method, PlcConnectionFactory, String, AliasRegistry, Map)} method.
     * <p>
     * If the tag is no getter, then all tags are refreshed by calling {@link #readAllFields(Object, PlcConnectionFactory, String, AliasRegistry, Map)}
     * and then, the method is invoked.
     *
     * @param proxy             Object to intercept
     * @param method            Method that was intercepted
     * @param callable          Callable to call the method after fetching the values
     * @param address           Address of the plc (injected from private tag)
     * @param connectionFactory PlcConnectionFactory instance to use (injected from private tag)
     * @return possible result of the original methods invocation
     * @throws OPMException Problems with plc / proxying
     */
    @SuppressWarnings({"unused", "squid:S00107"})
    @RuntimeType
    public static Object interceptGetter(@This Object proxy, @Origin Method method, @SuperCall Callable<?> callable,
                                         @FieldValue(PlcEntityManager.PLC_ADDRESS_FIELD_NAME) String address,
                                         @FieldValue(PlcEntityManager.CONNECTION_FACTORY_FIELD_NAME) PlcConnectionFactory connectionFactory,
                                         @FieldValue(PlcEntityManager.ALIAS_REGISTRY) AliasRegistry registry,
                                         @FieldValue(PlcEntityManager.LAST_FETCHED) Map<String, Instant> lastFetched,
                                         @FieldValue(PlcEntityManager.LAST_WRITTEN) Map<String, Instant> lastWritten) throws OPMException {
        LOGGER.trace("Invoked method {} on connected PlcEntity {}", method.getName(), method.getDeclaringClass().getName());

        // If "detached" (i.e. _driverManager is null) simply forward the call
        if (connectionFactory == null) {
            LOGGER.trace("Entity not connected, simply forwarding call");
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
            LOGGER.trace("Invoked method {} is getter, trying to find annotated tag and return requested value",
                method.getName());

            fetchAndSetValueForGetter(proxy, method, connectionFactory, address, registry, lastFetched);
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
            LOGGER.trace("Invoked method {} is boolean flag method, trying to find annotated tag and return requested value",
                method.getName());
            fetchAndSetValueForIsGetter(proxy, method, connectionFactory, address, registry, lastFetched);
            try {
                return callable.call();
            } catch (Exception e) {
                throw new OPMException("Unable to forward invocation " + method.getName() + " on connected PlcEntity", e);
            }
        }

        // Fetch all values then invoke method
        try {
            LOGGER.trace("Invoked method is no getter, refetch all tags and invoke method {} then", method.getName());
            readAllFields(proxy, connectionFactory, address, registry, lastFetched);
            Object call = callable.call();
            // We write back
            // cdutz: Disabled this, as it seemed to make no real sense to me, as we're writing back values that we just read without any chance of them being changed.
            //writeAllFields(proxy, connectionFactory, address, registry, lastWritten);
            return call;
        } catch (Exception e) {
            throw new OPMException("Unable to forward invocation " + method.getName() + " on connected PlcEntity", e);
        }
    }

    @SuppressWarnings({"unused", "squid:S00107"})
    @RuntimeType
    public static Object interceptSetter(@This Object proxy, @Origin Method method, @SuperCall Callable<?> callable,
                                         @FieldValue(PlcEntityManager.PLC_ADDRESS_FIELD_NAME) String address,
                                         @FieldValue(PlcEntityManager.CONNECTION_FACTORY_FIELD_NAME) PlcConnectionFactory connectionFactory,
                                         @FieldValue(PlcEntityManager.ALIAS_REGISTRY) AliasRegistry registry,
                                         @FieldValue(PlcEntityManager.LAST_FETCHED) Map<String, Instant> lastFetched,
                                         @Argument(0) Object argument) throws OPMException {
        LOGGER.trace("Invoked method {} on connected PlcEntity {}", method.getName(), method.getDeclaringClass().getName());

        // If "detached" (i.e. _driverManager is null) simply forward the call
        if (connectionFactory == null) {
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
            LOGGER.trace("Invoked method {} is setter, trying to find annotated tag and return requested value",
                method.getName());

            return setValueForSetter(proxy, method, callable, connectionFactory, address, registry, lastFetched, argument);
        }

        // Fetch all values then invoke method
        try {
            LOGGER.trace("Invoked method is no getter, refetch all tags and invoke method {} then", method.getName());
            readAllFields(proxy, connectionFactory, address, registry, lastFetched);
            return callable.call();
        } catch (Exception e) {
            throw new OPMException("Unable to forward invocation " + method.getName() + " on connected PlcEntity", e);
        }
    }

    /**
     * Reads all values of all tags that are annotated with {@link PlcEntity}.
     *
     * @param proxy         Object to refresh the tags on.
     * @param connectionFactory Connection Manager to use
     * @param registry      AliasRegistry to use
     * @param lastFetched   instants when which property was last fetched
     * @throws OPMException on various errors.
     */
    @SuppressWarnings("squid:S1141") // Nested try blocks readability is okay, move to other method makes it imho worse
    static void readAllFields(Object proxy, PlcConnectionFactory connectionFactory, String address, AliasRegistry registry, Map<String, Instant> lastFetched) throws OPMException {
        // Don't log o here as this would cause a second request against a plc so don't touch it, or if you log be aware of that
        Class<?> entityClass = proxy.getClass().getSuperclass();
        LOGGER.trace("Re-fetching all tags on proxy object of class {}", entityClass);
        PlcEntity plcEntity = entityClass.getAnnotation(PlcEntity.class);
        if (plcEntity == null) {
            throw new OPMException("Non PlcEntity supplied");
        }

        // Check if all tags are valid
        for (Field field : entityClass.getDeclaredFields()) {
            if (field.isAnnotationPresent(PlcTag.class)) {
                OpmUtils.getOrResolveAddress(registry, field.getAnnotation(PlcTag.class).value());
            }
        }
        try (PlcConnection connection = connectionFactory.getConnection(address)) {
            // Catch the exception, if no reader present (see below)
            // Build the query
            PlcReadRequest.Builder requestBuilder = connection.readRequestBuilder();

            Arrays.stream(entityClass.getDeclaredFields())
                .filter(field -> field.isAnnotationPresent(PlcTag.class))
                .filter(field -> needsToBeSynced(lastFetched, field))
                .forEach(field ->
                    requestBuilder.addTagAddress(
                        getFqn(field),
                        OpmUtils.getOrResolveAddress(registry, field.getAnnotation(PlcTag.class).value())
                    )
                );

            PlcReadRequest request = requestBuilder.build();

            LOGGER.trace("Request for re-fetch of {} was build and is {}", entityClass, request);

            PlcReadResponse response = getPlcReadResponse(request);

            // Fill all requested tags
            for (String fieldName : response.getTagNames()) {
                // Fill into Cache
                lastFetched.put(fieldName, Instant.now());

                LOGGER.trace("Value for tag {}  is {}", fieldName, response.getObject(fieldName));
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

    static void writeAllFields(Object proxy, PlcConnectionFactory connectionFactory, String address, AliasRegistry registry, Map<String, Instant> lastWritten) throws OPMException {
        // Don't log o here as this would cause a second request against a plc so don't touch it, or if you log be aware of that
        Class<?> entityClass = proxy.getClass().getSuperclass();
        LOGGER.trace("Writing all tags on proxy object of class {}", entityClass);
        PlcEntity plcEntity = entityClass.getAnnotation(PlcEntity.class);
        if (plcEntity == null) {
            throw new OPMException("Non PlcEntity supplied");
        }

        // Check if all tags are valid
        for (Field field : entityClass.getDeclaredFields()) {
            if (field.isAnnotationPresent(PlcTag.class)) {
                OpmUtils.getOrResolveAddress(registry, field.getAnnotation(PlcTag.class).value());
            }
        }
        try (PlcConnection connection = connectionFactory.getConnection(address)) {
            // Catch the exception, if no reader present (see below)
            // Build the query
            PlcWriteRequest.Builder requestBuilder = connection.writeRequestBuilder();

            Arrays.stream(entityClass.getDeclaredFields())
                .filter(field -> field.isAnnotationPresent(PlcTag.class))
                .filter(field -> needsToBeSynced(lastWritten, field))
                .forEach(field ->
                    requestBuilder.addTagAddress(
                        getFqn(field),
                        OpmUtils.getOrResolveAddress(registry, field.getAnnotation(PlcTag.class).value()),
                        getFromField(field, proxy)
                    )
                );

            PlcWriteRequest request = requestBuilder.build();

            LOGGER.trace("Request for write of {} was build and is {}", entityClass, request);

            PlcWriteResponse response = getPlcWriteResponse(request);

            // Fill all requested tags
            for (String fieldName : response.getTagNames()) {
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
     * Checks if a tags needs to be re-fetched/re-written, i.e., the cached values are too old.
     */
    private static boolean needsToBeSynced(Map<String, Instant> lastSynced, Field field) {
        Objects.requireNonNull(field);
        long cacheDurationMillis = field.getAnnotation(PlcTag.class).cacheDurationMillis();
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

    private static void fetchAndSetValueForIsGetter(Object proxy, Method m, PlcConnectionFactory connectionFactory, String address, AliasRegistry registry, Map<String, Instant> lastFetched) throws OPMException {
        fetchAndSetValueForGetter(proxy, m, 2, connectionFactory, address, registry, lastFetched);
    }

    private static void fetchAndSetValueForGetter(Object proxy, Method m, PlcConnectionFactory connectionFactory, String address, AliasRegistry registry, Map<String, Instant> lastFetched) throws OPMException {
        fetchAndSetValueForGetter(proxy, m, 3, connectionFactory, address, registry, lastFetched);
    }

    private static void fetchAndSetValueForGetter(Object proxy, Method m, int prefixLength, PlcConnectionFactory connectionFactory,
                                                  String address, AliasRegistry registry, Map<String, Instant> lastFetched) throws OPMException {
        String s = m.getName().substring(prefixLength);
        // First char to lower
        String variable = s.substring(0, 1).toLowerCase().concat(s.substring(1));
        LOGGER.trace("Looking for tag with name {} after invocation of getter {}", variable, m.getName());
        PlcTag annotation;
        Field field;
        try {
            field = m.getDeclaringClass().getDeclaredField(variable);
            annotation = field.getDeclaredAnnotation(PlcTag.class);
        } catch (NoSuchFieldException e) {
            throw new OPMException("Unable to identify tag with name '" + variable + "' for call to '" + m.getName() + "'", e);
        }

        // Use Fully qualified Name as tag index
        String fqn = getFqn(field);

        // Check if cache is still active
        if (!needsToBeSynced(lastFetched, field)) {
            return;
        }
        try (PlcConnection connection = connectionFactory.getConnection(address)) {
            // Catch the exception, if no reader present (see below)

            PlcReadRequest request = connection.readRequestBuilder()
                .addTagAddress(fqn, OpmUtils.getOrResolveAddress(registry, annotation.value()))
                .build();

            PlcReadResponse response = getPlcReadResponse(request);

            // Fill into Cache
            lastFetched.put(field.getName(), Instant.now());

            Object value = getTyped(m.getReturnType(), m.getGenericReturnType(), response, fqn);
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

    private static Object setValueForSetter(Object proxy, Method m, Callable<?> callable, PlcConnectionFactory connectionFactory,
                                            String address, AliasRegistry registry, Map<String, Instant> lastFetched, Object object) throws OPMException {
        String s = m.getName().substring(3);
        // First char to lower
        String variable = s.substring(0, 1).toLowerCase().concat(s.substring(1));
        LOGGER.trace("Looking for tag with name {} after invokation of getter {}", variable, m.getName());
        PlcTag annotation;
        Field field;
        try {
            field = m.getDeclaringClass().getDeclaredField(variable);
            annotation = field.getDeclaredAnnotation(PlcTag.class);
        } catch (NoSuchFieldException e) {
            throw new OPMException("Unable to identify tag with name '" + variable + "' for call to '" + m.getName() + "'", e);
        }

        // Use Fully qualified Name as tag index
        String fqn = getFqn(field);

        try (PlcConnection connection = connectionFactory.getConnection(address)) {
            // Catch the exception, if no reader present (see below)

            PlcWriteRequest.Builder builder = connection.writeRequestBuilder();
            if (object instanceof Collection) {
                builder.addTagAddress(fqn, OpmUtils.getOrResolveAddress(registry, annotation.value()), ((Collection<?>) object).toArray());
            } else {
                builder.addTagAddress(fqn, OpmUtils.getOrResolveAddress(registry, annotation.value()), object);
            }
            PlcWriteRequest request = builder.build();

            PlcWriteResponse response = getPlcWriteResponse(request);

            // Fill into Cache
            lastFetched.put(field.getName(), Instant.now());

            LOGGER.debug("getTyped clazz: {}, response: {}, tagName: {}", m.getParameters()[0].getType(), response, fqn);
            if (response.getResponseCode(fqn) != PlcResponseCode.OK) {
                throw new PlcRuntimeException(String.format("Unable to read specified tag '%s', response code was '%s'",
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
     * Tries to set a response Item to a tag in the given object.
     * This is one by looking for a tag in the class and a response item
     * which is equal to the given tagName parameter.
     *
     * @param o               Object to set the value on
     * @param response        Response to fetch the response from
     * @param targetTagName Name of the tag in the object
     * @param sourceTagName Name of the tag in the response
     * @throws NoSuchFieldException   If a tag is not present in entity
     * @throws IllegalAccessException If a tag in the entity cannot be accessed
     */
    static void setField(Class<?> clazz, Object o, PlcReadResponse response, String targetTagName, String sourceTagName) throws NoSuchFieldException, IllegalAccessException {
        LOGGER.debug("setField on clazz: {}, Object: {}, response: {}, targetFieldName: {}, sourceFieldName:{} ", clazz, o, response, targetTagName, sourceTagName);
        Field field = clazz.getDeclaredField(targetTagName);
        field.setAccessible(true);
        try {
            field.set(o, getTyped(field.getType(), field.getGenericType(), response, sourceTagName));
        } catch (ClassCastException e) {
            throw new PlcRuntimeException(String.format("Unable to assign return value %s to tag %s with type %s",
                response.getObject(sourceTagName), targetTagName, field.getType()), e);
        }
    }

    /** The element type of a collection type, or Object when it isn't known. */
    private static Class<?> elementType(Type genericType) {
        if (genericType instanceof ParameterizedType parameterized) {
            Type[] arguments = parameterized.getActualTypeArguments();
            if (arguments.length == 1 && arguments[0] instanceof Class<?> elementClass) {
                return elementClass;
            }
        }
        return Object.class;
    }

    /** Reads all values of a tag and returns them as an array of the given component type. */
    private static Object toArray(Class<?> componentType, PlcReadResponse response, String sourceFieldName) {
        List<?> values = new ArrayList<>(getAllTyped(componentType, response, sourceFieldName));
        Object array = Array.newInstance(componentType, values.size());
        for (int i = 0; i < values.size(); i++) {
            // Array.set unboxes for primitive component types.
            Array.set(array, i, values.get(i));
        }
        return array;
    }

    /** Reads all values of a tag into a collection assignable to the declared collection type. */
    private static Object toCollection(Class<?> clazz, Class<?> elementType, PlcReadResponse response,
                                       String sourceFieldName) {
        Collection<?> values = getAllTyped(elementType, response, sourceFieldName);
        if (clazz.isAssignableFrom(ArrayList.class)) {
            return new ArrayList<>(values);
        }
        if (clazz.isAssignableFrom(LinkedHashSet.class)) {
            return new LinkedHashSet<>(values);
        }
        throw new ClassCastException("Unable to return the values of tag " + sourceFieldName
            + " as an instance of " + clazz + ". Use List, Collection, Set or an array.");
    }

    /**
     * The typed accessor matching the element type, mirroring the scalar conversions below.
     * Unknown types fall back to the plain objects of the response.
     */
    private static Collection<?> getAllTyped(Class<?> elementType, PlcReadResponse response, String sourceFieldName) {
        if (elementType == boolean.class || elementType == Boolean.class) {
            return response.getAllBooleans(sourceFieldName);
        } else if (elementType == byte.class || elementType == Byte.class) {
            return response.getAllBytes(sourceFieldName);
        } else if (elementType == short.class || elementType == Short.class) {
            return response.getAllShorts(sourceFieldName);
        } else if (elementType == int.class || elementType == Integer.class) {
            return response.getAllIntegers(sourceFieldName);
        } else if (elementType == long.class || elementType == Long.class) {
            return response.getAllLongs(sourceFieldName);
        } else if (elementType == float.class || elementType == Float.class) {
            return response.getAllFloats(sourceFieldName);
        } else if (elementType == double.class || elementType == Double.class) {
            return response.getAllDoubles(sourceFieldName);
        } else if (elementType == BigInteger.class) {
            return response.getAllBigIntegers(sourceFieldName);
        } else if (elementType == BigDecimal.class) {
            return response.getAllBigDecimals(sourceFieldName);
        } else if (elementType == String.class) {
            return response.getAllStrings(sourceFieldName);
        } else if (elementType == LocalTime.class) {
            return response.getAllTimes(sourceFieldName);
        } else if (elementType == LocalDate.class) {
            return response.getAllDates(sourceFieldName);
        } else if (elementType == LocalDateTime.class) {
            return response.getAllDateTimes(sourceFieldName);
        }
        return response.getAllObjects(sourceFieldName);
    }

    static Object getTyped(Class<?> clazz, PlcReadResponse response, String sourceFieldName) {
        return getTyped(clazz, null, response, sourceFieldName);
    }

    /**
     * Converts the value of one tag of a response into the type of the annotated field.
     *
     * @param clazz       the declared type of the field or the return type of the getter
     * @param genericType the generic type of the same, when known. It carries the element type of a
     *                    collection field ({@code List<Integer>}), which the raw class does not -
     *                    see GH-1947. May be null, in which case collection elements come back as
     *                    plain objects.
     */
    @SuppressWarnings({"squid:S3776", "squid:MethodCyclomaticComplexity"})
    // Cognitive Complexity not too high, as highly structured
    static Object getTyped(Class<?> clazz, Type genericType, PlcReadResponse response, String sourceFieldName) {
        LOGGER.debug("getTyped clazz: {}, response: {}, tagName: {}", clazz, response, sourceFieldName);
        if (response.getResponseCode(sourceFieldName) != PlcResponseCode.OK) {
            throw new PlcRuntimeException(String.format("Unable to read specified tag '%s', response code was '%s'",
                sourceFieldName, response.getResponseCode(sourceFieldName)));
        }
        // A tag reading several values (e.g. "input-register:1:UINT[2]") has to be assigned to an
        // array or collection field, with its elements converted to the field's element type -
        // handing out the raw PlcValues makes every element access throw a ClassCastException.
        if (clazz.isArray()) {
            return toArray(clazz.getComponentType(), response, sourceFieldName);
        }
        if (Collection.class.isAssignableFrom(clazz)) {
            return toCollection(clazz, elementType(genericType), response, sourceFieldName);
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
