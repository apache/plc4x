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

import org.apache.plc4x.java.base.messages.items.FieldItem;

import java.util.*;

/**
 * Test device storing its state in memory.
 * Values are stored in a HashMap.
 */
class TestDevice {

    private final Random random = new Random();
    private final String name;
    private final Map<TestField, FieldItem> state = new HashMap<>();

    TestDevice(String name) {
        this.name = name;
    }

    Optional<FieldItem> get(TestField field) {
        Objects.requireNonNull(field);
        switch(field.getType()) {
            case STATE:
                return Optional.ofNullable(state.get(field));
            case RANDOM:
                return Optional.of(randomValue(field.getDataType()));
            case STDOUT:
                return Optional.empty();
        }
        throw new IllegalArgumentException("Unsupported field type: " + field.getType().name());
    }

    void set(TestField field, FieldItem value) {
        Objects.requireNonNull(field);
        switch (field.getType()) {
            case STATE:
                state.put(field, value);
                return;
            case STDOUT:
                System.out.printf("TEST PLC STDOUT [%s]: %s%n", field.getName(), Objects.toString(value.getValues()[0]));
                return;
            case RANDOM:
                System.out.printf("TEST PLC RANDOM [%s]: %s%n", field.getName(), Objects.toString(value.getValues()[0]));
                return;
        }
        throw new IllegalArgumentException("Unsupported field type: " + field.getType().name());
    }

    @SuppressWarnings("unchecked")
    private FieldItem randomValue(Class<?> type) {
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

        return new TestFieldItem(new Object[]{result});
    }

    @Override
    public String toString() {
        return name;
    }

}
