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

import java.util.List;
import java.util.Map;
import java.util.Set;

/**
 * Base Type of all Types.
 */
public interface PlcValue {

    // Simple Types

    boolean isSimple();

    boolean isNullable();

    boolean isNull();

    // Generic (\o/ Sebastian)

    boolean is(Class<?> clazz);

    boolean isConvertibleTo(Class<?> clazz);

    <T> T get(Class<T> clazz);

    // Boolean

    boolean isBoolean();

    boolean getBoolean();

    // Integer

    boolean isLong();

    long getLong();

    boolean isInteger();

    int getInteger();

    // Floating Point

    boolean isDouble();

    double getDouble();

    boolean isFloat();

    float getFloat();

    // String

    boolean isString();

    String getString();

    // Raw Access

    byte[] getRaw();

    // List Methods

    boolean isList();

    int length();

    PlcValue getIndex(int i);

    List<? extends PlcValue> getList();

    // Struct Methods

    boolean isStruct();

    Set<String> getKeys();

    boolean hasKey(String key);

    PlcValue getValue(String key);

    Map<String, ? extends PlcValue> getStruct();

}
