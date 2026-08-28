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
package org.apache.plc4x.java.umas.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.spi.drivers.model.AddressConstraints;
import org.apache.plc4x.java.spi.drivers.model.ArrayNotationParser;
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

    /**
     * A range has to be contiguous to be one read. Only the last dimension of the trailing
     * selection may span more than one element: "a[1].b[2..5]" is one run of a single
     * sub-structure, while "a[1..3].b" would be member b of three separate elements and
     * "a[1..3][2]" a strided slice - neither of which any single request can fetch.
     *
     * <p>An interior range is already refused by the symbolic path pattern, which accepts only
     * bare indices between the dots.
     */
    private static final AddressConstraints CONSTRAINTS =
        AddressConstraints.UNCONSTRAINED.withOnlyTrailingDimensionMayBeRange(true);

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
        // A trailing bracket run is the selection; everything before it is the symbolic path.
        // Splitting first is what lets the last segment carry a range - the per-segment group in
        // SYMBOLIC_ADDRESS_PATTERN would otherwise swallow a bare trailing index.
        String expression = ArrayNotationParser.expressionPart(address);
        List<ArrayInfo> selection = expression.isEmpty()
            ? Collections.emptyList()
            : ArrayNotationParser.parse(expression, address, CONSTRAINTS);
        return new SymbolicUmasTag(address, null, selection);
    }

    public static boolean matches(String address) {
        return SYMBOLIC_ADDRESS_PATTERN.matcher(ArrayNotationParser.addressPart(address)).matches();
    }

    /**
     * The selection the address states, or an empty list where it states none. Derived from the
     * address rather than from the constructor, because a tag built directly - as the driver does
     * when it browses the symbol table - carries the variable's declared shape in its arrayInfo,
     * not the user's selection.
     *
     * <p>Distinct from {@link #getArrayInfo()}, which describes the shape of the value the caller
     * receives.
     */
    public List<ArrayInfo> getSelection() {
        String expression = ArrayNotationParser.expressionPart(symbolicAddress);
        return expression.isEmpty()
            ? Collections.emptyList()
            : ArrayNotationParser.parse(expression, symbolicAddress, CONSTRAINTS);
    }

    /**
     * The declared lower bound the address states for its trailing dimension, or {@code null}
     * where it states none. The device's own declaration is authoritative; this is the user's
     * statement of intent, to be checked against it when the symbol is resolved.
     */
    public Integer getDeclaredBase() {
        String expression = ArrayNotationParser.expressionPart(symbolicAddress);
        if (expression.isEmpty() || !expression.contains(";")) {
            return null;
        }
        List<ArrayInfo> dimensions = ArrayNotationParser.parse(expression, symbolicAddress);
        return dimensions.get(dimensions.size() - 1).getBase();
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
        // Returning null (not PlcValueType.NULL) when type is unknown so the
        // DefaultPlcValueHandler preserves the original PlcValue on writes
        // instead of discarding it as PlcNull.
        return dataType;
    }

    @Override
    /**
     * The shape of the value the caller receives: empty for a scalar, one entry per dimension
     * for an array. A bare index selects one element and so reports empty; a range reports its
     * dimensions. Where the address states no selection, the driver fills this in from the
     * device's declared dimensions so a bare array address reports the whole array.
     */
    public List<ArrayInfo> getArrayInfo() {
        if (ArrayNotationParser.selectsSingleElement(
                ArrayNotationParser.expressionPart(symbolicAddress))) {
            return Collections.emptyList();
        }
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
