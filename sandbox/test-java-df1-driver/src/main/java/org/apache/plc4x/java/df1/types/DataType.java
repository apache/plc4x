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

import org.apache.plc4x.java.api.types.PlcValueType;

public enum DataType {
    BIT((short)1, PlcValueType.BOOL),
    BIT_STRING((short)-1, PlcValueType.List),
    BYTE_STRING((short)-1, PlcValueType.BYTE),
    INTEGER((short)2, PlcValueType.INT),
    TIMER((short)-1, PlcValueType.TIME),
    COUNTER((short)-1, PlcValueType.LINT),
    GENERAL_COUNT_STRUCTURE((short)-1, PlcValueType.Struct),
    FLOAT((short)-1, PlcValueType.REAL),
    ARRAY((short)-1, PlcValueType.List),
    ADDRESS((short)-1, PlcValueType.STRING),
    BINARY_CODED_DECIMAL((short)-1, PlcValueType.LINT);

    private final short length;
    private final PlcValueType plcValueType;

    DataType(short length, PlcValueType plcValueType) {
        this.length = length;
        this.plcValueType = plcValueType;
    }

    public short getLength() {
        return length;
    }

    public PlcValueType getPlcValueType() {
        return plcValueType;
    }

}
