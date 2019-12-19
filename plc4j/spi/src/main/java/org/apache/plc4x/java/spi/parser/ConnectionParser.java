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

import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.reflect.FieldUtils;

import java.lang.reflect.Field;
import java.net.URI;
import java.net.URISyntaxException;
import java.util.AbstractMap;
import java.util.Arrays;
import java.util.Collections;
import java.util.Iterator;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.function.Function;
import java.util.stream.Collectors;

import static java.util.stream.Collectors.mapping;
import static java.util.stream.Collectors.toList;

public class ConnectionParser {

    // TODO Respect Path Params
    public static <T> T parse(String string, Class<T> pClazz) {
        Map<String, Field> fieldMap = Arrays.stream(FieldUtils.getAllFields(pClazz))
            .filter(field -> field.getAnnotation(ConfigurationParameter.class) != null)
            .collect(Collectors.toMap(Field::getName, Function.identity()));

        T instance;
        try {
            instance = pClazz.newInstance();
        } catch (InstantiationException | IllegalAccessException e) {
            throw new IllegalStateException();
        }
        try {
            URI url = new URI(string);
            Map<String, List<String>> stringListMap = splitQuery(url);

            // TODO notify on umatched parameters

            Iterator<Map.Entry<String, Field>> iterator = fieldMap.entrySet().iterator();
            for (Iterator<Map.Entry<String, Field>> iter = iterator; iter.hasNext(); ) {
                Map.Entry<String, Field> entry = iter.next();
                // TODO field name also from annotation
                if (stringListMap.containsKey(entry.getKey())) {
                    fieldMap.get(entry.getKey()).setAccessible(true);
                    fieldMap.get(entry.getKey()).setInt(instance, Integer.parseInt(stringListMap.get(entry.getKey()).get(0)));
                    iter.remove();
                } else {
                    // TODO Implement other types
                    IntDefaultValue intDefaultValue = fieldMap.get(entry.getKey()).getAnnotation(IntDefaultValue.class);
                    if (intDefaultValue != null) {
                        fieldMap.get(entry.getKey()).setAccessible(true);
                        fieldMap.get(entry.getKey()).setInt(instance, intDefaultValue.value());
                        iter.remove();
                    }
                }
            }

            // TODO refactor
            List<String> missingFields = fieldMap.entrySet().stream()
                .filter(entry -> entry.getValue().getAnnotation(Required.class) != null)
                .map(entry -> entry.getValue().getAnnotation(ConfigurationParameter.class) != null ?
                    // In Memory of S. Ruehl
                    (StringUtils.isBlank(entry.getValue().getAnnotation(ConfigurationParameter.class).value()) ? entry.getKey() : entry.getValue().getAnnotation(ConfigurationParameter.class).value()) : entry.getKey())
                .collect(toList());

            if (missingFields.size() > 0) {
                throw new IllegalArgumentException("Missing required fields: " + missingFields);
            }
        } catch (URISyntaxException | IllegalAccessException e) {
            e.printStackTrace();
        }
        return instance;
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
