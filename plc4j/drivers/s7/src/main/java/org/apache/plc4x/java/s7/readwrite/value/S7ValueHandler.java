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

package org.apache.plc4x.java.api.value;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.IEC61131ValueHandler;
import org.apache.plc4x.java.api.value.PlcValueHandler;

import java.math.BigInteger;
import java.math.BigDecimal;
import java.nio.charset.StandardCharsets;
import java.util.BitSet;
import java.util.LinkedList;
import java.util.List;


public abstract class S7ValueHandler extends IEC61131ValueHandler {

    @Override
    public PlcValue customDataType(PlcField field, Object value) {
        switch (field.getPlcDataType().toUpperCase()) {
            case "TIME":
                return PlcBOOL.of(value);
            case "DATE":
                return PlcBYTE.of(value);
            case "DATE_AND_TIME":
                return PlcSINT.of(value);
            default:
                throw PlcUnsuppportedDataTypeException("Data Type " + field.getPlcDataType())
                    + "Is not supported");
        }
    }
}
