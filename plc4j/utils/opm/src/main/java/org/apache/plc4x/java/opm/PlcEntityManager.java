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
import net.bytebuddy.description.modifier.Visibility;
import net.bytebuddy.implementation.MethodDelegation;
import org.apache.commons.configuration2.Configuration;
import org.apache.commons.configuration2.SystemConfiguration;
import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.Arrays;
import java.util.concurrent.Callable;

import static net.bytebuddy.matcher.ElementMatchers.any;

/**
 * Plc4x equivalent of Jpas EntityManager for implementing Object-Plc-Mapping.
 * This means that calls to a plc can be done by using plain POJOs with Annotations.
 * <p>
 * First, the necessary annotations are {@link PlcEntity} and {@link PlcField}.
 * For a class to be usable as PlcEntity it needs
 * <ul>
 * <li>be non-final (as proxiing has to be used in case of {@link #connect(Class, String)}</li>
 * <li>a public no args constructor for instanciation</li>
 * <li>Needs to be annotated with {@link PlcEntity} and has a valid value which is the connection string</li>
 * </ul>
 * <p>
 * Basically, the {@link PlcEntityManager} has to operation "modes" represented by the methods {@link #read(Class,String)} and
 * {@link #connect(Class,String)}.
 * <p>
 * For a field to get Values from the Plc Injected it needs to be annotated with the {@link PlcField} annotation.
 * The value has to be the plc fields string (which is inserted in the {@link PlcReadRequest}).
 * The connection string is taken from the value of the {@link PlcEntity} annotation on the class.
 * <p>
 * The {@link #read(Class,String)} method has no direkt equivalent in JPA (as far as I know) as it only returns a "detached"
 * entity. This means it fetches all values from the plc that are annotated wiht the {@link PlcField} annotations.
 * <p>
 * The {@link #connect(Class,String)} method is more JPA-like as it returns a "connected" entity. This means, that each
 * time one of the getters on the returned entity is called a call is made to the plc (and the field value is changed
 * for this specific field).
 * Furthermore, if a method which is no getter is called, then all {@link PlcField}s are refreshed before doing the call.
 * Thus, all operations on fields that are annotated with {@link PlcField} are always done against the "live" values
 * from the PLC.
 * <p>
 * All invocations on the getters are forwarded to the {@link PlcEntityInterceptor#intercept(Object, Method, Callable, Object)}
 * method.
 * // TODO Add detach method
 */
public class PlcEntityManager {

    private static final Logger LOGGER = LoggerFactory.getLogger(PlcEntityManager.class);

    public static final String PLC_ADDRESS_FIELD_NAME = "_plcAddress";
    static final String DRIVER_MANAGER_FIELD_NAME = "_driverManager";

    private final PlcDriverManager driverManager;

    public PlcEntityManager() {
        this.driverManager = new PlcDriverManager();
    }

    public PlcEntityManager(PlcDriverManager driverManager) {
        this.driverManager = driverManager;
    }

    public <T> T read(Class<T> clazz, String address) throws OPMException {
        PlcEntity annotation = OpmUtils.getPlcEntityAndCheckPreconditions(clazz);

        try (PlcConnection connection = driverManager.getConnection(address)) {
            if (!connection.getMetadata().canRead()) {
                throw new OPMException("Unable to get Reader for connection with url '" + address + "'");
            }

            PlcReadRequest.Builder requestBuilder = connection.readRequestBuilder();

            Arrays.stream(clazz.getDeclaredFields())
                .filter(field -> field.isAnnotationPresent(PlcField.class))
                .forEach(field ->
                    requestBuilder.addItem(
                        field.getDeclaringClass().getName() + "." + field.getName(),
                        field.getAnnotation(PlcField.class).value()
                    )
                );

            // Build the request
            PlcReadRequest request = requestBuilder.build();

            // Perform the request
            PlcReadResponse response = PlcEntityInterceptor.getPlcReadResponse(request);

            // Construct the Object
            T instance = clazz.getConstructor().newInstance();

            // Fill all requested fields
            for (String fieldName : response.getFieldNames()) {
                String targetFieldName = StringUtils.substringAfterLast(fieldName, ".");
                PlcEntityInterceptor.setField(clazz, instance, response, targetFieldName, fieldName);
            }
            return instance;
        } catch (PlcInvalidFieldException e) {
            throw new OPMException("Unable to parse one field request", e);
        } catch (PlcConnectionException e) {
            throw new OPMException("Unable to get connection with url '" + address + "'", e);
        } catch (InstantiationException | InvocationTargetException | NoSuchMethodException | NoSuchFieldException | IllegalAccessException e) {
            throw new OPMException("Unable to fetch PlcEntity " + clazz.getName(), e);
        } catch (Exception e) {
            throw new OPMException("Unknown Error", e);
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
    public <T> T connect(Class<T> clazz, String address) throws OPMException {
        OpmUtils.getPlcEntityAndCheckPreconditions(clazz);
        try {
            // Use Byte Buddy to generate a subclassed proxy that delegates all PlcField Methods
            // to the intercept method
            T instance = new ByteBuddy()
                .subclass(clazz)
                .defineField(PLC_ADDRESS_FIELD_NAME, String.class, Visibility.PRIVATE)
                .defineField(DRIVER_MANAGER_FIELD_NAME, PlcDriverManager.class, Visibility.PRIVATE)
                .method(any()).intercept(MethodDelegation.to(PlcEntityInterceptor.class))
                .make()
                .load(Thread.currentThread().getContextClassLoader())
                .getLoaded()
                .getConstructor()
                .newInstance();
            // Set connection value into the private field
            FieldUtils.writeDeclaredField(instance, PLC_ADDRESS_FIELD_NAME, address, true);
            FieldUtils.writeDeclaredField(instance, DRIVER_MANAGER_FIELD_NAME, driverManager, true);
            return instance;
        } catch (NoSuchMethodException | InvocationTargetException | InstantiationException | IllegalAccessException e) {
            throw new OPMException("Unable to instantiate Proxy", e);
        }
    }

}
