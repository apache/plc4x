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

import java.util.Collections;
import java.util.HashMap;
import java.util.LinkedHashSet;
import java.util.Map;
import java.util.Set;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Metadata {

    public static class Key<T> {

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

    }

    public final static Metadata EMPTY = new Metadata(Collections.emptyMap());

    private final Metadata parent;
    private final Map<Key<?>, Object> values;

    Metadata(Map<Key<?>, Object> values) {
        this(values, EMPTY);
    }

    public Metadata(Map<Key<?>, Object> values, Metadata parent) {
        this.parent = parent;
        this.values = values;
    }

    public Set<Key<?>> keys() {
        Set<Key<?>> keys = new LinkedHashSet<>(values.keySet());
        keys.addAll(parent.keys());
        return Collections.unmodifiableSet(keys);
    }

    public Object getValue(Key<?> key) {
        Object value = values.get(key);
        if (value == null) {
            return parent.getValue(key);
        }
        return value;
    }

    public static class Builder {
        private final Logger logger = LoggerFactory.getLogger(Builder.class);

        private final Map<Key<?>, Object> values = new HashMap<>();
        private final Metadata parent;

        public Builder() {
            this(Metadata.EMPTY);
        }

        public Builder(Metadata parent) {
            this.parent = parent;
        }

        public <T> Builder put(Key<T> key, T value) {
            if (!key.validate(value)) {
                logger.debug("Ignore metadata value {}, it does not match constraints imposed by key {}", value, key);
                return this;
            }

            values.put(key, value);
            return this;
        }

        public Metadata build() {
            return new Metadata(values, parent);
        }
    }

}
