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

import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

class TagTagConfigParserTest {

    @Test
    public void testConfigOptions() {
        parse("aaa:123{unit-id: 10}", "unit-id", "10");
        parse("aaa:123{unit-id: -10}", "unit-id", "-10");
        parse("aaa:123{unit-id: 10.0}", "unit-id", "10.0");
        parse("aaa:123{unit-id: -10.0}", "unit-id", "-10.0");
        parse("aaa:123{unit-id: '10.0'}", "unit-id", "10.0");
        parse("aaa:123{unit-id: \"10.0\"}", "unit-id", "10.0");
        parse("aaa:123{unit-id: true}", "unit-id", "true");
        parse("aaa:123{unit-id: false}", "unit-id", "false");
        parse("aaa:123{val1: 1, val2: 2}", "val1", "1", "val2", "2");
    }

    private void parse(String address, String key, String value) {
        Map<String, String> config = TagConfigParser.parse(address);
        verify(config, key, value);
    }

    private void parse(String address, String key1, String value1, String key2, String value2) {
        Map<String, String> config = TagConfigParser.parse(address);
        verify(config, key1, value1);
        verify(config, key2, value2);
    }

    private void verify(Map<String, String> config, String key, String value) {
        assertTrue(config.containsKey(key), "Config does not contain key: " + key);
        assertEquals(value, config.get(key), "Value mismatch");
    }

}