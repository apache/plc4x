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
package org.apache.plc4x.java.opcua.protocol;

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.opcua.OpcuaConnection;
import org.apache.plc4x.java.opcua.readwrite.Variant;
import org.apache.plc4x.java.opcua.readwrite.VariantDataValue;
import org.apache.plc4x.java.opcua.readwrite.VariantDiagnosticInfo;
import org.apache.plc4x.java.opcua.readwrite.VariantExpandedNodeId;
import org.apache.plc4x.java.opcua.readwrite.VariantInt32;
import org.apache.plc4x.java.opcua.readwrite.VariantNull;
import org.apache.plc4x.java.opcua.readwrite.VariantVariant;
import org.apache.plc4x.java.opcua.tag.OpcuaTag;
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.apache.plc4x.java.spi.values.PlcNull;
import org.apache.plc4x.java.spi.values.PlcTIME;
import org.junit.jupiter.api.Test;

import java.util.List;

import static org.assertj.core.api.Assertions.assertThat;

/**
 * Tests the mapping of OPC UA variants to PlcValues, in particular that a node which exists
 * but currently carries no value (Null variant, VariantType 0) comes back as a PlcNull
 * instead of a raw {@code null} (see GH-2511).
 */
public class OpcuaVariantToPlcValueTest {

    @Test
    public void nullVariantBecomesPlcNull() {
        Variant variant = new VariantNull(false, false, null, null);

        PlcValue value = OpcuaConnection.variantToPlcValue(OpcuaTag.of("ns=2;i=1"), variant);

        assertThat(value).isInstanceOf(PlcNull.class);
        assertThat(value.isNull()).isTrue();
    }

    @Test
    public void nullVariantIsNotAffectedByATypeSuffix() {
        // A tag with a type suffix must not have its "no value" re-interpreted as a zero-valued
        // TIME - it stays a PlcNull.
        Variant variant = new VariantNull(false, false, null, null);

        PlcValue value = OpcuaConnection.variantToPlcValue(OpcuaTag.of("ns=2;i=1;TIME"), variant);

        assertThat(value).isInstanceOf(PlcNull.class);
    }

    @Test
    public void nullVariantIsHandledWithoutATag() {
        PlcValue value = OpcuaConnection.variantToPlcValue(null, new VariantNull(false, false, null, null));

        assertThat(value).isInstanceOf(PlcNull.class);
    }

    @Test
    public void regularVariantIsStillMapped() {
        Variant variant = new VariantInt32(false, false, null, null, null, List.of(42));

        PlcValue value = OpcuaConnection.variantToPlcValue(OpcuaTag.of("ns=2;i=1"), variant);

        assertThat(value).isEqualTo(new PlcDINT(42));
    }

    @Test
    public void unsupportedVariantTypesStillReturnNull() {
        // These variant types have no mapping yet, so callers still have to deal with a null
        // result - that is what tells them apart from a Null variant, which maps to PlcNull.
        // If one of them ever gets a mapping, drop it from here.
        assertThat(OpcuaConnection.variantToPlcValue(OpcuaTag.of("ns=2;i=1"),
            new VariantDataValue(false, false, null, null, null, List.of()))).isNull();
        assertThat(OpcuaConnection.variantToPlcValue(OpcuaTag.of("ns=2;i=1"),
            new VariantDiagnosticInfo(false, false, null, null, null, List.of()))).isNull();
        assertThat(OpcuaConnection.variantToPlcValue(OpcuaTag.of("ns=2;i=1"),
            new VariantVariant(false, false, null, null, null, List.of()))).isNull();
        assertThat(OpcuaConnection.variantToPlcValue(OpcuaTag.of("ns=2;i=1"),
            new VariantExpandedNodeId(false, false, null, null, null, List.of()))).isNull();
    }

    @Test
    public void typeSuffixStillOverridesARegularVariant() {
        Variant variant = new VariantInt32(false, false, null, null, null, List.of(42));

        PlcValue value = OpcuaConnection.variantToPlcValue(OpcuaTag.of("ns=2;i=1;TIME"), variant);

        assertThat(value).isEqualTo(new PlcTIME(42L));
    }
}
