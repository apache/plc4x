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
package org.apache.plc4x.java.spi.config;

import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Secret;

import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.util.Collections;
import java.util.LinkedHashSet;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;

/**
 * The parameter names a configuration class declares as carrying secrets.
 *
 * <p>This is what lets redaction be driven by the declaration rather than by a list of words
 * matched against names - see {@link Secret} for why that distinction matters.</p>
 */
public final class SecretParameters {

    /** Cached per class: the answer cannot change for the lifetime of the class. */
    private static final Map<Class<?>, Set<String>> CACHE = new ConcurrentHashMap<>();

    private SecretParameters() {
        // Utility class.
    }

    /**
     * The {@link ConfigurationParameter} names of every {@link Secret} field on the given class and
     * its supertypes.
     *
     * <p>A {@code @Secret} field without a {@code @ConfigurationParameter} contributes nothing here
     * - it has no name to match in a connection string - but is still covered by the
     * {@code toString()} rule.</p>
     *
     * @param configurationClass the class to inspect; {@code null} yields an empty set
     * @return an unmodifiable set of parameter names, never {@code null}
     */
    public static Set<String> namesFor(Class<?> configurationClass) {
        if (configurationClass == null) {
            return Collections.emptySet();
        }
        return CACHE.computeIfAbsent(configurationClass, SecretParameters::collect);
    }

    /** Every field marked as a secret on the given class and its supertypes, parameter or not. */
    public static Set<Field> fieldsOf(Class<?> configurationClass) {
        Set<Field> fields = new LinkedHashSet<>();
        Class<?> current = configurationClass;
        while ((current != null) && (current != Object.class)) {
            for (Field field : current.getDeclaredFields()) {
                if (!Modifier.isStatic(field.getModifiers()) && field.isAnnotationPresent(Secret.class)) {
                    fields.add(field);
                }
            }
            current = current.getSuperclass();
        }
        return fields;
    }

    private static Set<String> collect(Class<?> configurationClass) {
        Set<String> names = new LinkedHashSet<>();
        for (Field field : fieldsOf(configurationClass)) {
            ConfigurationParameter parameter = field.getAnnotation(ConfigurationParameter.class);
            if (parameter == null) {
                continue;
            }
            // A parameter with no explicit name is addressed by its field name, the same rule
            // ConfigurationFactory applies when it resolves one.
            names.add(parameter.value().isEmpty() ? field.getName() : parameter.value());
        }
        return Collections.unmodifiableSet(names);
    }
}
