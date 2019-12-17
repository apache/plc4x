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

import org.apache.plc4x.java.api.value.PlcSimpleValue;

public class PlcInteger extends PlcSimpleValue<Integer> {

    PlcInteger(Integer value) {
        super(value, true);
    }

    PlcInteger(int value) {
        super(value, false);
    }

    @Override public boolean isString() {
        return true;
    }

    @Override public String getString() {
        return value.toString();
    }

    @Override public boolean getBoolean() {
        // We like C
        return !(value == 0);
    }

    @Override public double getDouble() {
        return (double)value;
    }

    @Override public float getFloat() {
        return (float)value;
    }

    @Override public long getLong() {
        return value;
    }

    @Override public int getInteger() {
        return value;
    }

    @Override public String toString() {
        return String.valueOf(value);
    }
}
