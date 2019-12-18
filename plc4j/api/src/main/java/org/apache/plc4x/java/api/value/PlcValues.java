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

package org.apache.plc4x.java.api.value;

import org.apache.plc4x.java.api.exceptions.PlcIncompatibleDatatypeException;
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

    public static PlcValue of(String s) {
        return new PlcString(s);
    }

    public static PlcValue of(Integer i) {
        return new PlcInteger(i);
    }

    public static PlcValue of(int i) {
        return new PlcInteger(i);
    }

    public static PlcValue of(Boolean b) {
        return new PlcBoolean(b);
    }

    public static PlcValue of(boolean b) {
        return new PlcBoolean(b);
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
        try {
            String simpleName = o.getClass().getSimpleName();
            Class<?> clazz = o.getClass();
            switch (simpleName) {
                case "ArrayList":
                    simpleName = "List";
                    clazz = List.class;
            }
            Constructor<?> constructor = Class.forName(PlcValues.class.getPackage().getName() + ".Plc" + simpleName).getDeclaredConstructor(clazz);
            constructor.setAccessible(true);
            return ((PlcValue) constructor.newInstance(o));
        } catch (InstantiationException | IllegalAccessException | InvocationTargetException | NoSuchMethodException | ClassNotFoundException e) {
            LOGGER.warn("Cannot wrap", e);
            throw new PlcIncompatibleDatatypeException(o.getClass());
        }
    }


}
