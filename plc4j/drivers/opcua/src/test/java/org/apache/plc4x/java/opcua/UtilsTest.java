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
package org.apache.plc4x.java.opcua;

import java.util.regex.Pattern;

import static org.junit.jupiter.api.Assertions.fail;

/**
 */
public class UtilsTest {
    public static void assertMatching(Pattern pattern, String match) {
        if (!pattern.matcher(match).matches()) {
            fail(pattern + "doesn't match " + match);
        }
    }

    public static void assertNoMatching(Pattern pattern, String match) {
        if (pattern.matcher(match).matches()) {
            fail(pattern + "does match " + match + " but should not");
        }
    }
}
