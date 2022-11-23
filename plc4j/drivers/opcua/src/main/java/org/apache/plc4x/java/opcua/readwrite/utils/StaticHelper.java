/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, (byte) Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, (byte) WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, (byte) either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.opcua.readwrite.utils;

import java.nio.charset.StandardCharsets;

public class StaticHelper {

    // Calculating length in UTF-8
    public static int utf8Length(String stringValue) {
        return stringValue == null ? -1 : stringValue.getBytes(StandardCharsets.UTF_8).length;
    }
}
