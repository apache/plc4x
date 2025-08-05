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

package org.apache.plc4x.java.spi.metadata;

import static org.apache.plc4x.java.spi.generation.WithReaderWriterArgs.WithRenderAsList;

import java.util.Collections;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.Map;
import java.util.Objects;
import java.util.Set;
import org.apache.plc4x.java.api.metadata.Metadata;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class DefaultMetadata implements Metadata, Serializable {

    private final Metadata parent;
    private final Map<Key<?>, Object> values;

    DefaultMetadata(Map<Key<?>, Object> values) {
        this(values, Metadata.EMPTY);
    }

    public DefaultMetadata(Map<Key<?>, Object> values, Metadata parent) {
        this.values = new LinkedHashMap<>(values);
        this.parent = Objects.requireNonNull(parent, "Parent metadata must not be null");
    }

    @Override
    public Set<Key<?>> keys() {
        Set<Key<?>> keys = new LinkedHashSet<>(values.keySet());
        keys.addAll(parent.keys());
        return Collections.unmodifiableSet(keys);
    }

    @Override
    public Map<Key<?>, Object> entries() {
        Map<Key<?>, Object> copy = new LinkedHashMap<>(parent.entries());
        copy.putAll(values);
        return Map.copyOf(copy);
    }

    @Override
    public Object getValue(Key<?> key) {
        Object value = values.get(key);
        if (value == null) {
            return parent.getValue(key);
        }
        return value;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
            for (Key<?> metadataKey : keys()) {
            writeBuffer.pushContext("entry", WithRenderAsList(false));
            writeBuffer.writeString("key", metadataKey.getKey().length(), metadataKey.getKey());
            String value = "" + getValue(metadataKey);
            writeBuffer.writeString("value", value.length(), value);
            writeBuffer.popContext("entry");
        }
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof Metadata)) {
            return false;
        }
        Metadata that = (Metadata) o;
        return Objects.equals(entries(), that.entries());
    }

    @Override
    public int hashCode() {
        return Objects.hash(entries());
    }

    public static class Builder {
        private final Logger logger = LoggerFactory.getLogger(Builder.class);

        private final Map<Key<?>, Object> values = new LinkedHashMap<>();
        private final Metadata parent;

        public Builder() {
            this(DefaultMetadata.EMPTY);
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
            return new DefaultMetadata(values, parent);
        }
    }

}
