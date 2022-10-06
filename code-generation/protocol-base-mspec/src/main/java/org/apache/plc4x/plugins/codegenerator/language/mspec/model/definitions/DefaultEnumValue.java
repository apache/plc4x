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
package org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions;

import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;

import java.util.Map;
import java.util.Objects;
import java.util.Optional;

public class DefaultEnumValue implements EnumValue {

    private final String value;
    private final String name;
    private final Map<String, String> constants;

    public DefaultEnumValue(String value, String name, Map<String, String> constants) {
        this.value = Objects.requireNonNull(value);
        this.name = Objects.requireNonNull(name);
        this.constants = constants;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public String getValue() {
        return value;
    }

    @Override
    public Optional<String> getConstant(String name) {
        if (constants == null) {
            return Optional.empty();
        }
        return Optional.ofNullable(constants.get(name));
    }

    @Override
    public String toString() {
        return "DefaultEnumValue{" +
            "value='" + value + '\'' +
            ", name='" + name + '\'' +
            ", constants=" + constants +
            '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DefaultEnumValue that = (DefaultEnumValue) o;
        return Objects.equals(value, that.value) && Objects.equals(name, that.name) && Objects.equals(constants, that.constants);
    }

    @Override
    public int hashCode() {
        return Objects.hash(value, name, constants);
    }
}
