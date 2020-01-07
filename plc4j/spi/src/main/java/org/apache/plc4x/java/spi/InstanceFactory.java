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

package org.apache.plc4x.java.spi;

import org.apache.commons.lang3.ClassUtils;
import org.apache.plc4x.java.spi.parser.ConnectionParser;

import java.lang.reflect.ParameterizedType;
import java.util.Arrays;
import java.util.Optional;

/**
 * General Class which is used to create instances of Classes.
 */
public class InstanceFactory {

    private final ConnectionParser parser;

    public InstanceFactory(ConnectionParser parser) {
        this.parser = parser;
    }

    public InstanceFactory() {
        this(null);
    }

    public <T> T createInstance(Class<T> clazz) {
        try {
            T instance = clazz.newInstance();
            // Inject Configuration, if wanted
            if (ClassUtils.isAssignable(clazz, HasConfiguration.class)) {
                Optional<ParameterizedType> typeOptional = Arrays.stream(clazz.getGenericInterfaces())
                    .filter(type -> type instanceof ParameterizedType)
                    .map(type -> ((ParameterizedType) type))
                    .filter(type -> type.getRawType().equals(HasConfiguration.class))
                    .findAny();
                if (!typeOptional.isPresent()) {
                    throw new IllegalStateException("This should never happen!");
                }
                Class<?> configurationClass = (Class<?>) typeOptional.get().getActualTypeArguments()[0];
                // Try to get the Configuration
                Object configuration;
                if (parser != null) {
                    configuration = parser.createConfiguration(configurationClass);
                } else {
                    configuration = configurationClass.newInstance();
                }
                ((HasConfiguration) instance).setConfiguration(configuration);
                // System.out.println("The Configuration has to be of Type " + configurationType);
            }
            // Set all Properties
            // transport.setProperties(parser.getProperties());
            return instance;
        } catch (InstantiationException | IllegalAccessException e) {
            throw new IllegalStateException("Cannot Instantiate Transport '"
                + clazz.getSimpleName()
                + "'. Cannot access Default no Args Constructor.", e);
        }
    }
}
