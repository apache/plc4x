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

package org.apache.plc4x.java.spi.parser;

import org.apache.commons.beanutils.PropertyUtils;
import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import java.lang.reflect.Field;
import java.lang.reflect.InvocationTargetException;
import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.UnknownHostException;
import java.util.AbstractMap;
import java.util.Arrays;
import java.util.Collections;
import java.util.Iterator;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Properties;
import java.util.function.Function;
import java.util.stream.Collectors;

import static java.util.stream.Collectors.mapping;
import static java.util.stream.Collectors.toList;

/**
 * A query contains for our cases mostly of three parts
 * - protocol identifier
 * - connection address (ip/port), serial port, ...
 * - path parameters
 */
public class ConnectionParser {

    private final String connectionString;
    private URI uri;

    public ConnectionParser(String protocolCode, String connectionString) throws PlcConnectionException {
        this.connectionString = connectionString;
        try {
            this.uri = new URI(connectionString);
        } catch (URISyntaxException e) {
            throw new PlcConnectionException("Unable to parse connection string '" + connectionString + "'", e);
        }
        if (!protocolCode.equals(uri.getScheme())) {
            throw new PlcConnectionException("The given Connection String does not match the expected Protocol '" + protocolCode + "'");
        }
    }

    public SocketAddress getSocketAddress() {
        return this.getSocketAddress(-1);
    }

    /**
     * Convenvience Method, as its sometimes allowed to omit port in the URI String, as its
     * default for some protocols.
     * Of course only makes sense for TCP based Protocols
     *
     * @param defaultPort Default Port
     * @return Valid InetSocketAddress
     */
    public SocketAddress getSocketAddress(int defaultPort) {
        try {
            String hostString = uri.getHost();
            int port = uri.getPort();
            if (port == -1) {
                if (defaultPort == -1) {
                    throw new PlcRuntimeException("No port given in URI String and no default Port given!");
                } else {
                    port = defaultPort;
                }
            }
            return new InetSocketAddress(InetAddress.getByName(hostString), port);
        } catch (UnknownHostException e) {
            throw new PlcRuntimeException("Unable to resolve Host in connection  string '" + connectionString + "'", e);
        }
    }

    /**
     * Returns all Properties as Properties Object
     * All Keys are translated to Upper Case
     */
    public Properties getProperties() {
        Properties properties = new Properties();
        splitQuery(uri).entrySet()
            .forEach(entry -> properties.setProperty(entry.getKey().toUpperCase(), entry.getValue().get(0)));

        return properties;
    }

    private static Object getDefaultValueFromAnnotation(Field field) {
        IntDefaultValue intDefaultValue = field.getAnnotation(IntDefaultValue.class);
        if (intDefaultValue != null) {
            return intDefaultValue.value();
        }
        BooleanDefaultValue booleanDefaultValue = field.getAnnotation(BooleanDefaultValue.class);
        if (booleanDefaultValue != null) {
            return booleanDefaultValue.value();
        }
        DoubleDefaultValue doubleDefaultValue = field.getAnnotation(DoubleDefaultValue.class);
        if (doubleDefaultValue != null) {
            return doubleDefaultValue.value();
        }
        StringDefaultValue stringDefaultValue = field.getAnnotation(StringDefaultValue.class);
        if (stringDefaultValue != null) {
            return stringDefaultValue.value();
        }
        return null;
    }

    // TODO Respect Path Params
    public <T> T createConfiguration(Class<T> pClazz) {
        Map<String, Field> fieldMap = Arrays.stream(FieldUtils.getAllFields(pClazz))
            .filter(field -> field.getAnnotation(ConfigurationParameter.class) != null)
            .collect(Collectors.toMap(
                field -> getConfigurationName(field, field.getName()),
                Function.identity()
            ));

        T instance;
        try {
            instance = pClazz.newInstance();
        } catch (InstantiationException | IllegalAccessException e) {
            throw new IllegalArgumentException("Unable to Instantiate Configuration Class", e);
        }
        try {
            Map<String, List<String>> stringListMap = splitQuery(uri);

            // TODO notify on unmatched parameters

            Iterator<Map.Entry<String, Field>> iterator = fieldMap.entrySet().iterator();
            for (Iterator<Map.Entry<String, Field>> iter = iterator; iter.hasNext(); ) {
                Map.Entry<String, Field> entry = iter.next();
                // TODO field name also from annotation
                if (stringListMap.containsKey(entry.getKey())) {
                    final Field field = fieldMap.get(entry.getKey());
                    field.setAccessible(true);
                    if (field.getType().isAssignableFrom(String.class)) {
                        field.set(instance, stringListMap.get(entry.getKey()).get(0));
                    } else if (field.getType().isAssignableFrom(int.class)) {
                        field.setInt(instance, Integer.parseInt(stringListMap.get(entry.getKey()).get(0)));
                    }
                    iter.remove();
                } else {
                    // TODO Implement other types
                    // Check IntDefaultValue
                    Object defaultValue = getDefaultValueFromAnnotation(fieldMap.get(entry.getKey()));
                    if (defaultValue != null) {
                        try {
                            PropertyUtils.setSimpleProperty(instance, fieldMap.get(entry.getKey()).getName(), defaultValue);
                        } catch (InvocationTargetException | NoSuchMethodException e) {
                            throw new IllegalStateException(String.format("Unable to inject Configuration into field '%s' on Configuration %s", entry.getKey(), pClazz.getSimpleName()), e);
                        }
                    }
                    iter.remove();
                    continue;

                }
            }

            // TODO refactor
            List<String> missingFields = fieldMap.entrySet().stream()
                .filter(entry -> entry.getValue().getAnnotation(Required.class) != null)
                .map(entry -> entry.getValue().getAnnotation(ConfigurationParameter.class) != null ?
                    (getConfigurationName(entry.getValue(), entry.getKey())) : entry.getKey())
                .collect(toList());

            if (missingFields.size() > 0) {
                throw new IllegalArgumentException("Missing required fields: " + missingFields);
            }
        } catch (IllegalAccessException e) {
            throw new IllegalArgumentException("Unable to access all fields from Configuration Class '" + pClazz.getSimpleName() + "'", e);
        }
        return instance;
    }

    /** Extracts the Name from the configuration if given, uses given "name" otherwise */
    private String getConfigurationName(Field field, String name) {
        if (StringUtils.isBlank(field.getAnnotation(ConfigurationParameter.class).value())) {
            return name;
        } else {
            return field.getAnnotation(ConfigurationParameter.class).value();
        }
    }

    /**
     * https://stackoverflow.com/questions/13592236/parse-a-uri-string-into-name-value-collection/13592567#13592567
     */
    public static Map<String, List<String>> splitQuery(URI url) {
        if (StringUtils.isBlank(url.getQuery())) {
            return Collections.emptyMap();
        }
        return Arrays.stream(url.getQuery().split("&"))
            .map(ConnectionParser::splitQueryParameter)
            .collect(Collectors.groupingBy(AbstractMap.SimpleImmutableEntry::getKey, LinkedHashMap::new, mapping(Map.Entry::getValue, toList())));
    }

    public static AbstractMap.SimpleImmutableEntry<String, String> splitQueryParameter(String it) {
        final int idx = it.indexOf("=");
        final String key = idx > 0 ? it.substring(0, idx) : it;
        final String value = idx > 0 && it.length() > idx + 1 ? it.substring(idx + 1) : null;
        return new AbstractMap.SimpleImmutableEntry<>(key, value);
    }

}
