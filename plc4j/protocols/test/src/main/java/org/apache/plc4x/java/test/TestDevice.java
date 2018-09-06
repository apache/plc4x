/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.test;

import java.util.*;

/**
 * Test device storing its state in memory.
 * Values are stored in a HashMap.
 */
class TestDevice {
    private final Random random = new Random();
    private final String name;
    private final Map<TestAddress, Object> state = new HashMap<>();

    TestDevice(String name) {
        this.name = name;
    }

    @SuppressWarnings("unchecked")
    <T> Optional<T> get(Class<? extends T> type, TestAddress address) {
        Objects.requireNonNull(address);
        if (address.equals(TestAddress.RANDOM)) {
            return Optional.of(randomValue(type));
        } else {
            return Optional.ofNullable((T) state.get(address));
        }
    }

    public void set(TestAddress address, Object value) {
        state.put(address, value);
    }

    @SuppressWarnings("unchecked")
    private <T> T randomValue(Class<T> type) {
        Object result = null;

        // TODO: implement for further data types

        if (type == Integer.class)
            result = random.nextInt();

        if (type == Byte.class) {
            byte[] bytes = new byte[1];
            random.nextBytes(bytes);
            result = bytes[0];
        }

        if (type == Short.class) {
            result = random.nextInt(1 << 16);
        }

        return (T) result;
    }

    @Override
    public String toString() {
        return name;
    }
}
