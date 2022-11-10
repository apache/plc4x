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
package org.apache.plc4x.java.opm;

import org.apache.commons.lang3.Validate;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Utility methods for usage with OPM.
 */
final class OpmUtils {

    private static final String ADDRESS = "address";
    static final Pattern pattern = Pattern.compile("^\\$\\{(?<" + ADDRESS + ">.*)}$");

    private OpmUtils() {
        // Util class
    }

    static <T> PlcEntity getPlcEntityAndCheckPreconditions(Class<T> clazz) {
        PlcEntity annotation = clazz.getAnnotation(PlcEntity.class);
        if (annotation == null) {
            throw new IllegalArgumentException("Given Class is no Plc Entity, i.e., not annotated with @PlcEntity");
        }
        // Check if default constructor exists
        try {
            clazz.getConstructor();
        } catch (NoSuchMethodException e) {
            throw new IllegalArgumentException("Cannot use PlcEntity without default constructor", e);
        }
        return annotation;
    }

    static String getOrResolveAddress(AliasRegistry registry, String addressString) {
        if (!isValidExpression(addressString)) {
            throw new IllegalArgumentException("Invalid Syntax, either use tag address (no starting $) " +
                "or an alias with Syntax ${xxx}. But given was '" + addressString + "'");
        }
        if (!isAlias(addressString)) {
            return addressString;
        }
        String alias = getAlias(addressString);
        if (registry.canResolve(alias)) {
            return registry.resolve(alias);
        } else {
            throw new IllegalArgumentException("Unable to resolve Alias '" + alias + "' in Schema Registry");
        }

    }

    /**
     * Checks whether a given String is a valid OPM Expression, this means
     * either an Address or an alias ${xxx}.
     */
    static boolean isValidExpression(String s) {
        Validate.notNull(s);
        return !s.startsWith("$") || pattern.matcher(s).matches();
    }

    static boolean isAlias(String s) {
        Validate.notNull(s);
        return s.startsWith("$") && pattern.matcher(s).matches();
    }

    static String getAlias(String s) {
        Matcher matcher = pattern.matcher(s);
        if (!matcher.matches()) {
            throw new IllegalArgumentException("Invalid Syntax, no Alias found in String '" + s + "'. Syntax is ${xxx}");
        }
        return matcher.group(ADDRESS);
    }

}
