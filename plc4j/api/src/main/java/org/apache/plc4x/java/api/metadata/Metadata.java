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

package org.apache.plc4x.java.api.metadata;

import java.util.Map;
import java.util.Objects;
import java.util.Set;

public interface Metadata {

    Metadata EMPTY = new Metadata() {
        @Override
        public Set<Key<?>> keys() {
            return Set.of();
        }

        @Override
        public Map<Key<?>, Object> entries() {
            return Map.of();
        }

        @Override
        public Object getValue(Key<?> key) {
            return null;
        }

    };

    class Key<T> {

        private final String key;
        private final Class<T> type;

        protected Key(String key, Class<T> type) {
            this.key = key;
            this.type = type;
        }

        public String getKey() {
            return key;
        }

        public boolean validate(Object value) {
            return type.isInstance(value);
        }

        public static <T> Key<T> of(String key, Class<T> type) {
            return new Key<>(key, type);
        }

        @Override
        public boolean equals(Object o) {
            if (this == o) {
                return true;
            }
            if (!(o instanceof Key)) {
                return false;
            }
            Key<?> key1 = (Key<?>) o;
            return Objects.equals(getKey(), key1.getKey()) && Objects.equals(type, key1.type);
        }

        @Override
        public int hashCode() {
            return Objects.hash(getKey(), type);
        }
    }

    Set<Key<?>> keys();
    Map<Key<?>, Object> entries();
    Object getValue(Key<?> key);
}
