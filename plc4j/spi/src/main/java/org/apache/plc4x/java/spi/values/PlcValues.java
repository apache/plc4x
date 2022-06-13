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
package org.apache.plc4x.java.spi.values;

import org.apache.plc4x.java.api.exceptions.PlcIncompatibleDatatypeException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Constructor;
import java.lang.reflect.InvocationTargetException;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.Map;

public class PlcValues {

    private static final Logger LOGGER = LoggerFactory.getLogger(PlcValues.class);

    private PlcValues() {
    }

    public static PlcValue of(List<PlcValue> list) {
        return new PlcList(list);
    }

    public static PlcValue of(PlcValue... items) {
        return new PlcList(Arrays.asList(items));
    }

    public static PlcValue of(String key, PlcValue value) {
        return new PlcStruct(Collections.singletonMap(key, value));
    }

    public static PlcValue of(Map<String, PlcValue> map) {
        return new PlcStruct(map);
    }

    public static PlcValue of(Object o) {
        if (o == null) {
            return new PlcNull();
        }
        try {
            String simpleName = o.getClass().getSimpleName();
            Class<?> clazz = o.getClass();
            if (o instanceof List) {
                simpleName = "List";
                clazz = List.class;
            } else if (clazz.isArray()) {
                simpleName = "List";
                clazz = List.class;
                Object[] objectArray = (Object[]) o;
                o = Arrays.asList(objectArray);
            }
            // If it's one of the LocalDate, LocalTime or LocalDateTime, cut off the "Local".
            if (simpleName.startsWith("Local")) {
                simpleName = simpleName.substring(5);
            }
            Constructor<?> constructor = Class.forName(PlcValues.class.getPackage().getName() + ".Plc" + simpleName).getDeclaredConstructor(clazz);
            return ((PlcValue) constructor.newInstance(o));
        } catch (InstantiationException | IllegalAccessException | InvocationTargetException | NoSuchMethodException | ClassNotFoundException e) {
            LOGGER.warn("Cannot wrap", e);
            throw new PlcIncompatibleDatatypeException(o.getClass());
        }
    }
}
