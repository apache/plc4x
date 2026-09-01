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

import net.bytebuddy.ByteBuddy;
import net.bytebuddy.TypeCache;
import net.bytebuddy.description.modifier.Visibility;
import net.bytebuddy.implementation.MethodDelegation;
import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnectionFactory;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.time.Instant;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.Callable;

import static net.bytebuddy.matcher.ElementMatchers.isDeclaredBy;
import static net.bytebuddy.matcher.ElementMatchers.not;

/**
 * PLC4X equivalent of JPAs EntityManager for implementing Object-Plc-Mapping.
 * This means that calls to a plc can be done by using plain POJOs with Annotations.
 * <p>
 * First, the necessary annotations are {@link PlcEntity} and {@link PlcTag}.
 * For a class to be usable as PlcEntity it needs
 * <ul>
 * <li>be non-final (as proxying has to be used in case of {@link #connect(Class, String)}</li>
 * <li>a public no args constructor for instantiation</li>
 * <li>Needs to be annotated with {@link PlcEntity} and has a valid value which is the connection string</li>
 * </ul>
 * <p>
 * Basically, the {@link PlcEntityManager} has to operation "modes" represented by the methods {@link #read(Class, String)} and
 * {@link #connect(Class, String)}.
 * <p>
 * For a tag to get Values from the Plc Injected it needs to be annotated with the {@link PlcTag} annotation.
 * The value has to be the plc tags string (which is inserted in the {@link PlcReadRequest}).
 * The connection string is taken from the value of the {@link PlcEntity} annotation on the class.
 * <p>
 * The {@link #read(Class, String)} method has no direkt equivalent in JPA (as far as I know) as it only returns a "detached"
 * entity. This means it fetches all values from the plc that are annotated with the {@link PlcTag} annotations.
 * <p>
 * The {@link #connect(Class, String)} method is more JPA-like as it returns a "connected" entity. This means, that each
 * time one of the getters on the returned entity is called a call is made to the plc (and the tag value is changed
 * for this specific tag).
 * Furthermore, if a method which is no getter is called, then all {@link PlcTag}s are refreshed before doing the call.
 * Thus, all operations on tags that are annotated with {@link PlcTag} are always done against the "live" values
 * from the PLC.
 * <p>
 * A connected @{@link PlcEntity} can be disconnected calling {@link #disconnect(Object)}, then it behaves like the
 * regular Pojo it was before.
 * <p>
 * All invocations on the getters are forwarded to the
 * {@link PlcEntityInterceptor#interceptGetter(Object, Method, Callable, String, PlcConnectionFactory, AliasRegistry, Map, Map)}
 * method.
 */
public class PlcEntityManager {

    private static final Logger LOGGER = LoggerFactory.getLogger(PlcEntityManager.class);

    public static final String PLC_ADDRESS_FIELD_NAME = "_plcAddress";
    static final String CONNECTION_FACTORY_FIELD_NAME = "_connectionFactory";
    static final String ALIAS_REGISTRY = "_aliasRegistry";
    public static final String LAST_FETCHED = "_lastFetched";
    public static final String LAST_WRITTEN = "_lastWritten";

    /**
     * Generated proxy classes, keyed by class loader and entity class. Soft references so the
     * cache can't keep class loaders alive on its own.
     */
    private static final TypeCache<Class<?>> PROXY_CACHE =
        new TypeCache.WithInlineExpunction<>(TypeCache.Sort.SOFT);

    private final PlcConnectionFactory connectionFactory;
    private final SimpleAliasRegistry registry;

    public PlcEntityManager() {
        this(new DefaultPlcDriverManager());
    }

    public PlcEntityManager(PlcConnectionFactory connectionFactory) {
        this(connectionFactory, new SimpleAliasRegistry());
    }

    public PlcEntityManager(PlcConnectionFactory connectionFactory, SimpleAliasRegistry registry) {
        this.connectionFactory = connectionFactory;
        this.registry = registry;
    }

    public <T> T read(Class<T> clazz, String address) throws OPMException {
        T connect = connect(clazz, address);
        disconnect(connect);
        return connect;
    }

    public <T> T write(Class<T> clazz, String address, T object) throws OPMException {
        T merge = merge(clazz, address, object);
        disconnect(merge);
        return merge;
    }

    /**
     * Returns a connected proxy.
     *
     * @param clazz clazz to be connected.
     * @param <T>   type of param {@code clazz}.
     * @return a connected entity.
     * @throws OPMException when proxy can't be built.
     */
    public <T> T connect(Class<T> clazz, String address) throws OPMException {
        return connect(clazz, address, null);
    }


    /**
     * Returns a connected proxy.
     *
     * @param clazz clazz to be connected.
     * @param <T>   type of param {@code clazz}.
     * @return a connected entity.
     * @throws OPMException when proxy can't be built.
     */
    public <T> T merge(Class<T> clazz, String address, T instance) throws OPMException {
        return connect(clazz, address, instance);
    }

    private <T> T connect(Class<T> clazz, String address, T existingInstance) throws OPMException {
        OpmUtils.getPlcEntityAndCheckPreconditions(clazz);
        try {
            T instance = proxyClassFor(clazz)
                .getConstructor()
                .newInstance();
            // Set connection value into the private field
            FieldUtils.writeDeclaredField(instance, PLC_ADDRESS_FIELD_NAME, address, true);
            FieldUtils.writeDeclaredField(instance, CONNECTION_FACTORY_FIELD_NAME, connectionFactory, true);
            FieldUtils.writeDeclaredField(instance, ALIAS_REGISTRY, registry, true);
            Map<String, Instant> lastFetched = new HashMap<>();
            FieldUtils.writeDeclaredField(instance, LAST_FETCHED, lastFetched, true);
            Map<String, Instant> lastWritten = new HashMap<>();
            FieldUtils.writeDeclaredField(instance, LAST_WRITTEN, lastWritten, true);

            // Initially fetch all values
            if (existingInstance == null) {
                PlcEntityInterceptor.readAllFields(instance, connectionFactory, address, registry, lastFetched);
            } else {
                // Copy all field values from the existing instance to the new one.
                FieldUtils.getAllFieldsList(clazz).stream()
                    .peek(field -> field.setAccessible(true))
                    .forEach(field -> setValueToField(field, instance, getValueFromField(field, existingInstance)));

                PlcEntityInterceptor.writeAllFields(instance, connectionFactory, address, registry, lastWritten);
            }

            return instance;
        } catch (NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException | IllegalAccessError e) {
            throw new OPMException("Unable to instantiate Proxy", e);
        }
    }

    /**
     * Returns the generated proxy class for the given entity class, generating it on first use.
     * <p>
     * The generated type depends only on the entity class, so it is cached and reused - see
     * GH-1935. Generating and loading a fresh subclass on every {@link #connect(Class, String)}
     * (and therefore on every {@link #read(Class, String)} / {@link #write(Class, String, Object)})
     * leaks one class per call, which fills up metaspace in any polling application.
     */
    @SuppressWarnings("unchecked")
    private <T> Class<? extends T> proxyClassFor(Class<T> clazz) throws OPMException {
        ClassLoader classLoader = Thread.currentThread().getContextClassLoader();
        try {
            // The cache is keyed by class loader and entity class, so the same entity loaded by
            // two class loaders correctly gets a proxy per loader. The monitor makes concurrent
            // callers wait for one generation instead of racing to generate the same class.
            return (Class<? extends T>) PROXY_CACHE.findOrInsert(classLoader, clazz,
                () -> new ByteBuddy()
                    .subclass(clazz)
                    .defineField(PLC_ADDRESS_FIELD_NAME, String.class, Visibility.PRIVATE)
                    .defineField(CONNECTION_FACTORY_FIELD_NAME, PlcConnectionFactory.class, Visibility.PRIVATE)
                    .defineField(ALIAS_REGISTRY, AliasRegistry.class, Visibility.PRIVATE)
                    .defineField(LAST_FETCHED, Map.class, Visibility.PRIVATE)
                    .defineField(LAST_WRITTEN, Map.class, Visibility.PRIVATE)
                    .method(not(isDeclaredBy(Object.class))).intercept(MethodDelegation.to(PlcEntityInterceptor.class))
                    .make()
                    .load(classLoader)
                    .getLoaded(),
                PROXY_CACHE);
        } catch (RuntimeException e) {
            throw new OPMException("Unable to generate Proxy for " + clazz.getName(), e);
        }
    }

    private Object getValueFromField(Field field, Object object) {
        try {
            return field.get(object);
        } catch (IllegalAccessException e) {
            throw new PlcRuntimeException(e);
        }
    }

    private void setValueToField(Field field, Object object, Object value) {
        try {
            field.set(object, value);
        } catch (IllegalAccessException e) {
            throw new PlcRuntimeException(e);
        }
    }

    /**
     * Disconnects the given instance.
     *
     * @param entity Instance of a PlcEntity.
     * @throws OPMException Is thrown when the plc is already disconnected or no entity.
     */
    public void disconnect(Object entity) throws OPMException {
        // Check if this is an entity
        PlcEntity annotation = entity.getClass().getSuperclass().getAnnotation(PlcEntity.class);
        if (annotation == null) {
            throw new OPMException("Unable to disconnect Object, is no entity!");
        }
        try {
            Object manager = FieldUtils.readDeclaredField(entity, CONNECTION_FACTORY_FIELD_NAME, true);
            if (manager == null) {
                throw new OPMException("Instance is already disconnected!");
            }
            FieldUtils.writeDeclaredField(entity, CONNECTION_FACTORY_FIELD_NAME, null, true);
        } catch (IllegalAccessException e) {
            throw new OPMException("Unable to fetch driverManager instance on entity instance", e);
        }
    }

}
