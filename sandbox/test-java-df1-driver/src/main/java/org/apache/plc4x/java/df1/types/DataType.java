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
package org.apache.plc4x.java.df1.types;

public enum DataType {
    BIT((short)1),
    BIT_STRING((short)-1),
    BYTE_STRING((short)-1),
    INTEGER((short)2),
    TIMER((short)-1),
    COUNTER((short)-1),
    GENERAL_COUNT_STRUCTURE((short)-1),
    FLOAT((short)-1),
    ARRAY((short)-1),
    ADDRESS((short)-1),
    BINARY_CODED_DECIMAL((short)-1);

    private final short length;

    DataType(short length) {
        this.length = length;
    }

    public short getLength() {
        return length;
    }

}
