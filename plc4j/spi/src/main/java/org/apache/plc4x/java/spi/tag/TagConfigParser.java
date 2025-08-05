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

package org.apache.plc4x.java.spi.tag;

import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Dedicated parser which parses config portion of tag address.
 * Tag config comes in curly braces and follows json syntax with key value pairs (fields are not quoted)
 */
public class TagConfigParser {

    public static final Pattern TAG_CONFIG_PATTERN = Pattern.compile("(\\{(?<config>.*?)})?$");

    protected static final Pattern KEY_VALUE_PATTERN = Pattern.compile("(?<parameter>[\\w\\-_]+):\\s*(?<value>-?\\d+.\\d+|-?\\d+|\"[^\"]*\"|'[^']*'|true|false),?");

    private TagConfigParser() {
        // Prevent this from being instantiated.
    }

    public static Map<String, String> parse(String tagAddress) {
        Map<String, String> params = new LinkedHashMap<>();
        Matcher matcher = TAG_CONFIG_PATTERN.matcher(tagAddress);
        if (matcher.find() && matcher.group("config") != null) {
            Matcher kv = KEY_VALUE_PATTERN.matcher(matcher.group("config"));
            while (kv.find()) {
                params.put(kv.group("parameter"), clean(kv.group("value")));
            }
        }
        return params;
    }

    private static String clean(String value) {
        if (value.startsWith("'") && value.endsWith("'")) {
            return value.substring(1, value.length() - 1);
        }
        if (value.startsWith("\"") && value.endsWith("\"")) {
            return value.substring(1, value.length() - 1);
        }
        return value;
    }

}
