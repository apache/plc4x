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
package org.apache.plc4x.java.opcua.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.opcua.readwrite.AttributeId;
import org.apache.plc4x.java.opcua.readwrite.OpcuaDataType;
import org.apache.plc4x.java.opcua.readwrite.OpcuaIdentifierType;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class OpcuaTagCoverageTest {

    @Test
    void parsesNumericIdentifier() {
        OpcuaTag tag = OpcuaTag.of("ns=2;i=10846");
        assertThat(tag.getIdentifier()).isEqualTo("10846");
        assertThat(tag.getIdentifierType()).isEqualTo(OpcuaIdentifierType.NUMBER_IDENTIFIER);
        assertThat(tag.getNamespace()).isEqualTo(2);
        assertThat(tag.getAttributeId()).isEqualTo(AttributeId.Value);
        assertThat(tag.getDataType()).isEqualTo(OpcuaDataType.NULL);
        assertThat(tag.getConfig()).isEmpty();
    }

    @Test
    void parsesStringIdentifierWithDataType() {
        OpcuaTag tag = OpcuaTag.of("ns=2;s=foo.bar;DINT");
        assertThat(tag.getIdentifier()).isEqualTo("foo.bar");
        assertThat(tag.getIdentifierType()).isEqualTo(OpcuaIdentifierType.STRING_IDENTIFIER);
        assertThat(tag.getDataType()).isEqualTo(OpcuaDataType.DINT);
        assertThat(tag.getPlcValueType()).isEqualTo(PlcValueType.DINT);
    }

    @Test
    void parsesGuidAndBinaryIdentifiers() {
        OpcuaTag guid = OpcuaTag.of("ns=2;g=09087e75-8e5e-499b-954f-f2a8624db28a");
        assertThat(guid.getIdentifierType()).isEqualTo(OpcuaIdentifierType.GUID_IDENTIFIER);
        OpcuaTag bin = OpcuaTag.of("ns=2;b=asvaewavarahreb==");
        assertThat(bin.getIdentifierType()).isEqualTo(OpcuaIdentifierType.BINARY_IDENTIFIER);
    }

    @Test
    void parsesAttributeIdByNameAndNumber() {
        assertThat(OpcuaTag.of("ns=2;i=1;a=NodeId").getAttributeId()).isEqualTo(AttributeId.NodeId);
        // 13 is the numeric value of AttributeId.Value.
        OpcuaTag byNumber = OpcuaTag.of("ns=2;i=1;a=13");
        assertThat(byNumber.getAttributeId()).isEqualTo(AttributeId.Value);
    }

    @Test
    void parsesTrailingConfigSegment() {
        // The address regex's identifier group ([^;]+) is greedy, so a trailing
        // |k=v config only parses when a ;DATATYPE segment separates them.
        OpcuaTag tag = OpcuaTag.of("ns=2;i=1;DINT|sampling-interval=500,queue-size=10");
        assertThat(tag.getConfig())
            .containsEntry("sampling-interval", "500")
            .containsEntry("queue-size", "10");
    }

    @Test
    void matchesStaticHelperAgreesWithOf() {
        String valid = "ns=2;i=10846";
        assertThat(OpcuaTag.matches(valid)).isTrue();
        assertThat(OpcuaTag.matches("definitely-not-an-opcua-address")).isFalse();
    }

    @Test
    void invalidAddressThrows() {
        assertThatThrownBy(() -> OpcuaTag.of("not-a-tag"))
            .isInstanceOf(PlcInvalidTagException.class);
    }

    @Test
    void unknownDataTypeThrows() {
        assertThatThrownBy(() -> OpcuaTag.of("ns=2;i=1;NOT_A_TYPE"))
            .isInstanceOf(PlcUnsupportedDataTypeException.class);
    }

    @Test
    void getAddressStringRoundtripsWithoutDataType() {
        // No data type → no trailing ;TYPE
        OpcuaTag tag = OpcuaTag.of("ns=2;i=10846");
        // The toString suffixes ;NULL because the default datatype is NULL,
        // but getAddressString omits attribute-id=Value (the default).
        assertThat(tag.getAddressString()).contains("ns=2;i=10846");
    }

    @Test
    void getAddressStringIncludesNonDefaultAttribute() {
        OpcuaTag tag = OpcuaTag.of("ns=2;i=1;a=NodeId");
        assertThat(tag.getAddressString()).contains(";a=NodeId");
    }

    @Test
    void getTagReturnsSelfShape() {
        OpcuaTag tag = OpcuaTag.of("ns=2;i=10846");
        assertThat(tag.getTag()).isInstanceOf(OpcuaTag.class);
        assertThat(((OpcuaTag) tag.getTag()).getIdentifier()).isEqualTo("10846");
    }

    @Test
    void equalsAndHashCodeRespectAddressIdentity() {
        OpcuaTag a = OpcuaTag.of("ns=2;i=10846");
        OpcuaTag b = OpcuaTag.of("ns=2;i=10846");
        OpcuaTag c = OpcuaTag.of("ns=2;i=99999");
        assertThat(a).isEqualTo(a).isEqualTo(b).isNotEqualTo(c).isNotEqualTo("string").isNotEqualTo(null);
        assertThat(a.hashCode()).isEqualTo(b.hashCode());
    }

    @Test
    void toStringRendersAllFields() {
        String s = OpcuaTag.of("ns=2;s=foo;DINT").toString();
        assertThat(s).contains("namespace=2").contains("identifier=foo");
    }

    @Test
    void subscriptionAccessorsReturnDefaults() {
        OpcuaTag tag = OpcuaTag.of("ns=2;i=1");
        assertThat(tag.getPlcSubscriptionType()).isNull();
        assertThat(tag.getDuration()).isEmpty();
        // Default array info comes from the PlcSubscriptionTag interface default impl.
        assertThat(tag.getArrayInfo()).isNotNull();
    }
}
