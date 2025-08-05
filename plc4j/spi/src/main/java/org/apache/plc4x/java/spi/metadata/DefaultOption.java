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

package org.apache.plc4x.java.spi.metadata;

import org.apache.plc4x.java.api.metadata.Option;
import org.apache.plc4x.java.api.types.OptionType;

import java.util.Optional;

public class DefaultOption implements Option {

    private final String key;
    private final OptionType type;
    private final String description;
    private final boolean required;
    private final Object defaultValue;
    private final String since;

    public DefaultOption(String key, OptionType type, String description, boolean required, Object defaultValue, String since) {
        this.key = key;
        this.type = type;
        this.description = description;
        this.required = required;
        this.defaultValue = defaultValue;
        this.since = since;
    }

    @Override
    public String getKey() {
        return key;
    }

    @Override
    public OptionType getType() {
        return type;
    }

    @Override
    public String getDescription() {
        return description;
    }

    @Override
    public boolean isRequired() {
        return required;
    }

    @Override
    public Optional<Object> getDefaultValue() {
        return Optional.ofNullable(defaultValue);
    }

    @Override
    public Optional<String> getSince() {
        return Optional.ofNullable(since);
    }

}
