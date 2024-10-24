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

package org.apache.plc4x.java.opcua.tag;

import org.apache.plc4x.java.opcua.readwrite.StatusCode;

public final class OpcuaQualityStatus {

    private static final long STATUS_MASK = 0xC0000000L;
    private static final long STATUS_GOOD = 0x00000000L;
    private static final long STATUS_UNCERTAIN = 0x40000000L;
    private static final long STATUS_BAD = 0x80000000L;

    private final StatusCode statusCode;

    public OpcuaQualityStatus(StatusCode statusCode) {
        this.statusCode = statusCode;
    }

    public boolean isGood() {
        return (statusCode.getStatusCode() & STATUS_MASK) == STATUS_GOOD;
    }

    public boolean isBad() {
        return (statusCode.getStatusCode() & STATUS_MASK) == STATUS_BAD;
    }

    public boolean isUncertain() {
        return (statusCode.getStatusCode() & STATUS_MASK) == STATUS_UNCERTAIN;
    }

    @Override
    public String toString() {
        if (isGood()) {
            return "good";
        } else if (isBad()) {
            return "bad";
        }
        return "uncertain";
    }
}
