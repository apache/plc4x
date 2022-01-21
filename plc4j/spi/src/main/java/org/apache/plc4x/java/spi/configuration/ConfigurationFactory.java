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
package org.apache.plc4x.java.spi.configuration;

import org.apache.commons.lang3.ClassUtils;
import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.spi.configuration.annotations.*;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.*;

import java.lang.reflect.Field;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.ParameterizedType;
import java.lang.reflect.Type;
import java.net.URLDecoder;
import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.function.Function;
import java.util.stream.Collectors;

import static java.util.stream.Collectors.mapping;
import static java.util.stream.Collectors.toList;

/**
 * A query contains for our cases mostly of three parts
 * - protocol identifier
 * - (optional) transport identifier
 * - connection address (ip/port), serial port, path (depending on the type of transport) ...
 * - (optional) path parameters
 */
public class ConfigurationFactory {

    // TODO Respect Path Params
    public <T extends Configuration> T createConfiguration(Class<T> pClazz, String configurationString) {
        // Get a map of all configuration parameter fields.
        // - Get a list of all fields in the given class.
        Map<String, Field> fields = Arrays.stream(FieldUtils.getAllFields(pClazz))
            // - Filter out only the ones annotated with the ConfigurationParameter annotation.
            .filter(field -> field.getAnnotation(ConfigurationParameter.class) != null)
            // - Create a map with the field-name as key and the field itself as value.
            .collect(Collectors.toMap(
                ConfigurationFactory::getConfigurationName,
                Function.identity()
            ));

        // Get a list of all required configuration parameters.
        List<String> missingFieldNames = fields.values().stream()
            .filter(field -> field.getAnnotation(Required.class) != null)
            .map(ConfigurationFactory::getConfigurationName)
            .collect(toList());

        // Create a new instance of the configuration object.
        T instance;
        try {
            instance = pClazz.getDeclaredConstructor().newInstance();
        } catch (InvocationTargetException | InstantiationException |
            IllegalAccessException | NoSuchMethodException e) {
            throw new IllegalArgumentException("Unable to Instantiate Configuration Class", e);
        }

        // Process the parameters passed in with the connection string.
        try {
            // Get a map of all parameters in the connection string.
            Map<String, List<String>> paramStringValues = splitQuery(configurationString);

            // Iterate over all fields and set the values to either the values specified
            // in the param string or to defaults configured by annotations.
            for (Map.Entry<String, Field> entry : fields.entrySet()) {
                final String configName = entry.getKey();
                final Field field = fields.get(configName);
                if (paramStringValues.containsKey(configName)) {
                    String stringValue = paramStringValues.get(configName).get(0);
                    // As the arguments might be URL encoded, be sure it's decoded.
                    stringValue = URLDecoder.decode(stringValue, StandardCharsets.UTF_8);
                    setFieldValue(instance, field, stringValue);
                    missingFieldNames.remove(configName);
                } else {
                    boolean success = setFieldDefaultValueFromAnnotation(instance, field);
                    if (success) {
                        missingFieldNames.remove(configName);
                    }
                }
            }

            // If in the end still some required parameters are missing, output an error.
            if (!missingFieldNames.isEmpty()) {
                throw new IllegalArgumentException("Missing required fields: " + missingFieldNames);
            }
        } catch (IllegalAccessException e) {
            throw new IllegalArgumentException("Unable to access all fields from Configuration Class '" + pClazz.getSimpleName() + "'", e);
        }
        return instance;
    }

    public static <T> T configure(Configuration configuration, T obj) {
        // Check if in this object is configurable at all.
        if (ClassUtils.isAssignable(obj.getClass(), HasConfiguration.class)) {
            // Check if the type declared by the HasConfiguration interface is
            // compatible with the given configuration type.
            Optional<ParameterizedType> typeOptional = Arrays.stream(obj.getClass().getGenericInterfaces())
                // Check if the interface has a type parameter
                .filter(ParameterizedType.class::isInstance)
                .map(ParameterizedType.class::cast)
                .filter(type -> type.getRawType().equals(HasConfiguration.class))
                .findAny();
            if (typeOptional.isPresent()) {
                final ParameterizedType parameterizedType = typeOptional.get();
                final Type configType = parameterizedType.getActualTypeArguments()[0];
                if (configType instanceof Class) {
                    Class<?> configClass = (Class<?>) configType;
                    if (configClass.isAssignableFrom(configuration.getClass())) {
                        ((HasConfiguration) obj).setConfiguration(configuration);
                    }
                }
            }

        }
        return obj;
    }

    /**
     * Get the configuration parameter name for configuration parameters.
     * If an explicit name is provided in the annotation, use that else use the name of the field itself.
     *
     * @param field name of the field.
     * @return name of the configuration (either from the annotation or from the field itself)
     */
    private static String getConfigurationName(Field field) {
        if (StringUtils.isBlank(field.getAnnotation(ConfigurationParameter.class).value())) {
            return field.getName();
        } else {
            return field.getAnnotation(ConfigurationParameter.class).value();
        }
    }

    /**
     * Set the instance field value.
     *
     * @param instance    instance that should be set
     * @param field       field that should be set
     * @param valueString string representation of the value
     */
    private static void setFieldValue(Object instance, Field field, String valueString) throws IllegalAccessException {
        if (field == null) {
            throw new IllegalArgumentException("Field not defined");
        }

        try {
            // Make the field accessible
            field.setAccessible(true);

            // If a ParameterConverter is provided, use this for the conversion instead.
            if (field.getAnnotation(ParameterConverter.class) != null) {
                Class<? extends ConfigurationParameterConverter<?>> converterClass = field.getAnnotation(ParameterConverter.class).value();
                try {
                    ConfigurationParameterConverter<?> converter = converterClass.getDeclaredConstructor().newInstance();
                    if (converter.getType().isAssignableFrom(field.getType())) {
                        Object value = converter.convert(valueString);
                        field.set(instance, value);
                        return;
                    }
                    throw new IllegalArgumentException("Unsupported field type " + field.getType() + " for converter " + converterClass);
                } catch (InstantiationException | IllegalAccessException | InvocationTargetException | NoSuchMethodException e) {
                    throw new IllegalArgumentException("Could not initialize parameter converter", e);
                }
            }
            // Use the default value conversions.
            else {
                if (field.getType() == String.class) {
                    field.set(instance, valueString);
                } else if ((field.getType() == boolean.class) || (field.getType() == Boolean.class)) {
                    field.setBoolean(instance, Boolean.parseBoolean(valueString));
                } else if ((field.getType() == byte.class) || (field.getType() == Byte.class)) {
                    field.setByte(instance, Byte.parseByte(valueString));
                } else if ((field.getType() == short.class) || (field.getType() == Short.class)) {
                    field.setShort(instance, Short.parseShort(valueString));
                } else if ((field.getType() == int.class) || (field.getType() == Integer.class)) {
                    field.setInt(instance, Integer.parseInt(valueString));
                } else if ((field.getType() == long.class) || (field.getType() == Long.class)) {
                    field.setLong(instance, Long.parseLong(valueString));
                } else if ((field.getType() == float.class) || (field.getType() == Float.class)) {
                    field.setFloat(instance, Float.parseFloat(valueString));
                } else if ((field.getType() == double.class) || (field.getType() == Double.class)) {
                    field.setDouble(instance, Double.parseDouble(valueString));
                } else {
                    throw new IllegalArgumentException("Unsupported property type " + field.getType().getName());
                }
            }
        } finally {
            // Ensure we undo the accessible-making
            field.setAccessible(false);
        }
    }

    /**
     * Set the instance field default value from annotation.
     *
     * @param instance    instance that should be set
     * @param field       field that should be set
     * @return field is set
     */
    private static boolean setFieldDefaultValueFromAnnotation(Object instance, Field field) throws IllegalAccessException {
        try {
            // Make the field accessible
            field.setAccessible(true);

            // Go through the possible default-value annotations and use the first one that's there.
            IntDefaultValue intDefaultValue = field.getAnnotation(IntDefaultValue.class);
            if (intDefaultValue != null) {
                field.setInt(instance, intDefaultValue.value());
                return true;
            }
            BooleanDefaultValue booleanDefaultValue = field.getAnnotation(BooleanDefaultValue.class);
            if (booleanDefaultValue != null) {
                field.setBoolean(instance, booleanDefaultValue.value());
                return true;
            }
            FloatDefaultValue floatDefaultValue = field.getAnnotation(FloatDefaultValue.class);
            if (floatDefaultValue != null) {
                field.setFloat(instance, floatDefaultValue.value());
                return true;
            }
            DoubleDefaultValue doubleDefaultValue = field.getAnnotation(DoubleDefaultValue.class);
            if (doubleDefaultValue != null) {
                field.setDouble(instance, doubleDefaultValue.value());
                return true;
            }
            StringDefaultValue stringDefaultValue = field.getAnnotation(StringDefaultValue.class);
            if (stringDefaultValue != null) {
                field.set(instance, stringDefaultValue.value());
                return true;
            }
        } finally {
            // Ensure we undo the accessible-making
            field.setAccessible(false);
        }
        return false;
    }

    /**
     * https://stackoverflow.com/questions/13592236/parse-a-uri-string-into-name-value-collection/13592567#13592567
     */
    private static Map<String, List<String>> splitQuery(String paramString) {
        if (StringUtils.isBlank(paramString)) {
            return Collections.emptyMap();
        }
        // Split the individual parameters which are separated by "&" characters.
        return Arrays.stream(paramString.split("&"))
            // Split the stream of parameters (Which can be "key=value", "key=" or just "key" into tuples of
            // "key & value"
            .map(ConfigurationFactory::splitQueryParameter)
            // Build a map of "key & List<value>" where the values of elements with equal "key" are
            // added to a list of values.
            .collect(Collectors.groupingBy(AbstractMap.SimpleImmutableEntry::getKey, LinkedHashMap::new,
                mapping(Map.Entry::getValue, toList())));
    }

    private static AbstractMap.SimpleImmutableEntry<String, String> splitQueryParameter(String it) {
        final int idx = it.indexOf('=');
        final String key = idx > 0 ? it.substring(0, idx) : it;
        final String value = idx > 0 && it.length() > idx + 1 ? it.substring(idx + 1) : null;
        return new AbstractMap.SimpleImmutableEntry<>(key, value);
    }

}
