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
package org.apache.plc4x.java.umas.readwrite.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.util.Collections;
import java.util.List;
import java.util.Objects;
import java.util.regex.Pattern;

/**
 * UMAS tag identified by a symbolic name.
 * Supports IEC 61131-3 conventions:
 * <ul>
 *   <li>{@code g_r32} -- simple global variable</li>
 *   <li>{@code g_plant.meta.r32} -- nested struct member access</li>
 *   <li>{@code g_arrInt[3]} -- array element access</li>
 *   <li>{@code g_plant.items[2].value} -- mixed struct/array access</li>
 * </ul>
 */
public class SymbolicUmasTag implements UmasTag {

    private static final Pattern SYMBOLIC_ADDRESS_PATTERN =
        Pattern.compile("^([a-zA-Z_]\\w*)(\\[\\d+])*(\\.([a-zA-Z_]\\w*)(\\[\\d+])*)*$");

    private final String symbolicAddress;
    private final PlcValueType dataType;
    private final List<ArrayInfo> arrayInfo;

    public SymbolicUmasTag(String symbolicAddress, PlcValueType dataType, List<ArrayInfo> arrayInfo) {
        this.symbolicAddress = Objects.requireNonNull(symbolicAddress);
        this.dataType = dataType;
        this.arrayInfo = Objects.requireNonNull(arrayInfo);
    }

    public static SymbolicUmasTag of(String address) {
        if (!matches(address)) {
            throw new PlcInvalidTagException(address, SYMBOLIC_ADDRESS_PATTERN, "{symbolic-address}");
        }
        return new SymbolicUmasTag(address, null, Collections.emptyList());
    }

    public static boolean matches(String address) {
        return SYMBOLIC_ADDRESS_PATTERN.matcher(address).matches();
    }

    public String getSymbolicAddress() {
        return symbolicAddress;
    }

    @Override
    public String getAddressString() {
        return symbolicAddress;
    }

    @Override
    public PlcValueType getPlcValueType() {
        // Return null (not PlcValueType.NULL) when type is unknown so the
        // DefaultPlcValueHandler preserves the original PlcValue on writes
        // instead of discarding it as PlcNull.
        return dataType;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return arrayInfo;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof SymbolicUmasTag that)) return false;
        return Objects.equals(symbolicAddress, that.symbolicAddress);
    }

    @Override
    public int hashCode() {
        return Objects.hash(symbolicAddress);
    }

    @Override
    public String toString() {
        return "SymbolicUmasTag{symbolicAddress='" + symbolicAddress + "'}";
    }

}
