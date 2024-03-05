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
package org.apache.plc4x.java.spi.configuration;

import org.apache.commons.lang3.ClassUtils;
import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.configuration.annotations.ComplexConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.ParameterConverter;
import org.apache.plc4x.java.spi.configuration.annotations.Required;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

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

    private static final Logger LOGGER = LoggerFactory.getLogger(ConfigurationFactory.class);

    public <T extends PlcConnectionConfiguration> T createConfiguration(Class<T> pClazz, String protocolCode, String transportCode,
                                                                        String transportConfig, String paramString) {

        // Get a map of all parameters in the connection string.
        Map<String, List<String>> paramStringValues = splitQuery(paramString);
        return createConfiguration(pClazz, protocolCode, transportCode, transportConfig, paramStringValues);
    }

    public <T extends PlcTransportConfiguration> T createTransportConfiguration(Class<T> pClazz, String protocolCode,
                                                                                String transportCode, String transportConfig,
                                                                                String paramString) {
        // Get a map of all parameters in the connection string.
        Map<String, List<String>> paramStringValues = splitQuery(paramString);
        // Filter out the properties, that don't have the current transport code as prefix.
        String prefix = transportCode + ".";
        Map<String, List<String>> filteredParamStringValues = new HashMap<>();
        for (String paramName : paramStringValues.keySet()) {
            if(paramName.startsWith(prefix)) {
                filteredParamStringValues.put(paramName.substring(prefix.length()), paramStringValues.get(paramName));
            }
        }
        return createConfiguration(pClazz, protocolCode, transportCode, transportConfig, filteredParamStringValues);
    }

    public <T> T createConfiguration(Class<T> pClazz, String protocolCode, String transportCode,
                                                           String transportConfig, Map<String, List<String>> paramStringValues) {
        // Get a map of all configuration parameter fields.
        // - Get a list of all fields in the given class.
        Map<String, Field> fields = Arrays.stream(FieldUtils.getAllFields(pClazz))
            // - Filter out only the ones annotated with the ConfigurationParameter annotation.
            .filter(field -> (field.getAnnotation(ConfigurationParameter.class) != null) || (field.getAnnotation(ComplexConfigurationParameter.class) != null))
            // - Create a map with the field-name as key and the field itself as value.
            .collect(Collectors.toMap(
                ConfigurationFactory::getConfigurationName,
                Function.identity()
            ));

        // Get a list of all required configuration parameters.
        // TODO: Check for the complex-fields with required annotations.
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
            // Add the other codes to the param strings, so we can make them accessible from Configurations.
            // Like this:
            //    @ConfigurationParameter("protocolCode")
            //    private String protocolCode;
            paramStringValues = new HashMap<>(paramStringValues);
            List<String> previousValue;
            previousValue = paramStringValues.put("protocol-code", List.of(protocolCode));
            if (previousValue != null) {
                LOGGER.warn("protocolCode with value {} overridden by", protocolCode);
            }
            previousValue = paramStringValues.put("transport-code", List.of(transportCode));
            if (previousValue != null) {
                LOGGER.warn("transportCode with value {} overridden by", transportCode);
            }
            previousValue = paramStringValues.put("transport-config", List.of(transportConfig));
            if (previousValue != null) {
                LOGGER.warn("transportConfig with value {} overridden by", transportConfig);
            }

            // Iterate over all fields and set the values to either the values specified
            // in the param string or to defaults configured by annotations.
            for (Map.Entry<String, Field> entry : fields.entrySet()) {
                final String configName = entry.getKey();
                final Field field = fields.get(configName);
                if (field.getAnnotation(ComplexConfigurationParameter.class) != null) {
                    // Filter out only the parameters with the given prefix.
                    String prefix = field.getAnnotation(ComplexConfigurationParameter.class).prefix() + ".";
                    Map<String, List<String>> filteredParamStringValues = new HashMap<>();
                    for (String paramName : paramStringValues.keySet()) {
                        if(paramName.startsWith(prefix)) {
                            filteredParamStringValues.put(paramName.substring(prefix.length()), paramStringValues.get(paramName));
                        }
                    }
                    Class<PlcConfiguration> configType = (Class<PlcConfiguration>) field.getType();
                    PlcConfiguration configValue = createConfiguration(configType, protocolCode, transportCode, transportConfig, filteredParamStringValues);
                    FieldUtils.writeField(instance, field.getName(), configValue, true);
                } else if (paramStringValues.containsKey(configName)) {
                    String stringValue = paramStringValues.get(configName).get(0);
                    // As the arguments might be URL encoded, be sure it's decoded.
                    stringValue = URLDecoder.decode(stringValue, StandardCharsets.UTF_8);
                    FieldUtils.writeField(instance, field.getName(), toFieldValue(field, stringValue), true);
                    missingFieldNames.remove(configName);
                } else {
                    Object defaultValue = getDefaultValueFromAnnotation(field);
                    // TODO: Check if the default values type matches.
                    if (defaultValue != null) {
                        FieldUtils.writeField(instance, field.getName(), defaultValue, true);
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

    public static <T> T configure(PlcConfiguration configuration, T obj) {
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
                        try {
                            ((HasConfiguration) obj).setConfiguration(configuration);
                        } catch (Throwable t) {
                            LOGGER.error("Error setting the configuration", t);
                            throw new PlcRuntimeException("Error setting the configuration", t);
                        }
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
    public static String getConfigurationName(Field field) {
        if (field.getAnnotation(ComplexConfigurationParameter.class) != null) {
            return field.getAnnotation(ComplexConfigurationParameter.class).prefix();
        } else if (StringUtils.isBlank(field.getAnnotation(ConfigurationParameter.class).value())) {
            return field.getName();
        } else {
            return field.getAnnotation(ConfigurationParameter.class).value();
        }
    }

    /**
     * Convert the string value from the parameter string into the type the field requires.
     *
     * @param field       field that should be set
     * @param valueString string representation of the value
     * @return parsed value of the field in the type the field requires
     */
    private static Object toFieldValue(Field field, String valueString) {
        if (field == null) {
            throw new IllegalArgumentException("Field not defined");
        }

        if (field.getAnnotation(ParameterConverter.class) != null) {
            Class<? extends ConfigurationParameterConverter<?>> converterClass = field.getAnnotation(ParameterConverter.class).value();

            try {
                ConfigurationParameterConverter<?> converter = converterClass.getDeclaredConstructor().newInstance();
                if (converter.getType().isAssignableFrom(field.getType())) {
                    return converter.convert(valueString);
                }
            } catch (InstantiationException | IllegalAccessException | InvocationTargetException |
                     NoSuchMethodException e) {
                throw new IllegalArgumentException("Could not initialize parameter converter", e);
            }
            throw new IllegalArgumentException("Unsupported field type " + field.getType() + " for converter " + converterClass);
        }

        if (field.getType() == String.class) {
            return valueString;
        }
        if ((field.getType() == boolean.class) || (field.getType() == Boolean.class)) {
            return Boolean.parseBoolean(valueString);
        }
        if ((field.getType() == byte.class) || (field.getType() == Byte.class)) {
            return Byte.parseByte(valueString);
        }
        if ((field.getType() == short.class) || (field.getType() == Short.class)) {
            return Short.parseShort(valueString);
        }
        if ((field.getType() == int.class) || (field.getType() == Integer.class)) {
            return Integer.parseInt(valueString);
        }
        if ((field.getType() == long.class) || (field.getType() == Long.class)) {
            return Long.parseLong(valueString);
        }
        if ((field.getType() == float.class) || (field.getType() == Float.class)) {
            return Float.parseFloat(valueString);
        }
        if ((field.getType() == double.class) || (field.getType() == Double.class)) {
            return Double.parseDouble(valueString);
        }
        if (field.getType().isEnum()) {
            return parseEnumValue(field, valueString);
        }
        throw new IllegalArgumentException("Unsupported property type " + field.getType().getName());
    }

    public static Object getDefaultValueFromAnnotation(Field field) {
        IntDefaultValue intDefaultValue = field.getAnnotation(IntDefaultValue.class);
        if (intDefaultValue != null) {
            return intDefaultValue.value();
        }
        LongDefaultValue longDefaultValue = field.getAnnotation(LongDefaultValue.class);
        if(longDefaultValue != null) {
            return longDefaultValue.value();
        }
        BooleanDefaultValue booleanDefaultValue = field.getAnnotation(BooleanDefaultValue.class);
        if (booleanDefaultValue != null) {
            return booleanDefaultValue.value();
        }
        FloatDefaultValue floatDefaultValue = field.getAnnotation(FloatDefaultValue.class);
        if (floatDefaultValue != null) {
            return floatDefaultValue.value();
        }
        DoubleDefaultValue doubleDefaultValue = field.getAnnotation(DoubleDefaultValue.class);
        if (doubleDefaultValue != null) {
            return doubleDefaultValue.value();
        }
        StringDefaultValue stringDefaultValue = field.getAnnotation(StringDefaultValue.class);
        if (stringDefaultValue != null) {
            if (field.getType().isEnum()) {
                return parseEnumValue(field, stringDefaultValue.value());
            }
            return stringDefaultValue.value();
        }
        return null;
    }

    private static Enum<?> parseEnumValue(Field field, String valueString) {
        return Enum.valueOf((Class<Enum>) field.getType(), valueString);
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
