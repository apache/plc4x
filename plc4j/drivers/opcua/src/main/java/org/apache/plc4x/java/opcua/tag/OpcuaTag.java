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

import java.nio.charset.StandardCharsets;
import java.util.Map;
import java.util.Map.Entry;
import org.apache.commons.lang3.EnumUtils;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.opcua.readwrite.AttributeId;
import org.apache.plc4x.java.opcua.readwrite.OpcuaDataType;
import org.apache.plc4x.java.opcua.readwrite.OpcuaIdentifierType;

import java.time.Duration;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
public class OpcuaTag implements PlcSubscriptionTag {

    // Inline tag-config pattern that the old SPI's {@code TagConfigParser} used
    // to append; kept here so the address-string syntax stays compatible.
    private static final String TAG_CONFIG_PATTERN = "(\\|(?<config>(?:(?:[a-zA-Z\\-_]+=[a-zA-Z0-9\\-_]+)(?:,(?:[a-zA-Z\\-_]+=[a-zA-Z0-9\\-_]+))*)))?";
    private static final String OPC_UTA_TAG_ADDRESS = "^ns=(?<namespace>\\d+);(?<identifierType>[isgb])=(?<identifier>[^;]+)?(;a=(?<attributeId>[^;]+))?(;(?<datatype>[a-zA-Z_]+))?";
    public static final Pattern ADDRESS_PATTERN = Pattern.compile(OPC_UTA_TAG_ADDRESS + TAG_CONFIG_PATTERN + "$");

    private final OpcuaIdentifierType identifierType;

    private final int namespace;

    private final String identifier;

    private final AttributeId attributeId;

    private final OpcuaDataType dataType;

    private final Map<String, String> config;

    private OpcuaTag(Integer namespace, String identifier, OpcuaIdentifierType identifierType, AttributeId attributeId, OpcuaDataType dataType, Map<String, String> config) {
        this.identifier = Objects.requireNonNull(identifier);
        this.identifierType = Objects.requireNonNull(identifierType);
        this.namespace = namespace != null ? namespace : 0;
        if (this.namespace < 0) {
            throw new IllegalArgumentException("namespace must be greater then zero. Was " + this.namespace);
        }
        this.attributeId = attributeId;
        this.dataType = dataType;
        this.config = config;
    }

    public static OpcuaTag of(String address) {
        Matcher matcher = ADDRESS_PATTERN.matcher(address);
        if (!matcher.matches()) {
            throw new PlcInvalidTagException(address, ADDRESS_PATTERN, "{address}");
        }
        String identifier = matcher.group("identifier");

        String identifierTypeString = matcher.group("identifierType");
        OpcuaIdentifierType identifierType = OpcuaIdentifierType.enumForValue(identifierTypeString);

        String namespaceString = matcher.group("namespace");
        Integer namespace = namespaceString != null ? Integer.parseInt(namespaceString) : 0;

        String dataTypeString = matcher.group("datatype") != null ? matcher.group("datatype").toUpperCase() : "NULL";
        if (!EnumUtils.isValidEnum(OpcuaDataType.class, dataTypeString)) {
            throw new PlcUnsupportedDataTypeException("Datatype " + dataTypeString + " is unsupported by this protocol");
        }
        OpcuaDataType dataType = OpcuaDataType.valueOf(dataTypeString);

        String attributeElement = matcher.group("attributeId");
        AttributeId attributeId = AttributeId.Value;
        if (attributeElement != null) {
            if (attributeElement.matches("\\d+")) {
                attributeId = AttributeId.enumForValue(Long.parseLong(attributeElement));
            } else {
                attributeId = AttributeId.valueOf(attributeElement);
            }
        }
        return new OpcuaTag(namespace, identifier, identifierType, attributeId, dataType, parseConfig(matcher.group("config")));
    }

    /** Parses the tag's config tail ({@code |k=v,k=v}) into a map. */
    private static Map<String, String> parseConfig(String config) {
        if (config == null || config.isEmpty()) {
            return java.util.Collections.emptyMap();
        }
        Map<String, String> result = new java.util.HashMap<>();
        for (String entry : config.split(",")) {
            int eq = entry.indexOf('=');
            if (eq > 0) {
                result.put(entry.substring(0, eq), entry.substring(eq + 1));
            }
        }
        return result;
    }

    @Override
    public PlcTag getTag() {
        return new OpcuaTag(namespace, identifier, identifierType, attributeId, dataType, config);
    }

    public static boolean matches(String address) {
        return ADDRESS_PATTERN.matcher(address).matches();
    }

    public int getNamespace() {
        return namespace;
    }

    public String getIdentifier() {
        return identifier;
    }

    public OpcuaIdentifierType getIdentifierType() {
        return identifierType;
    }

    public OpcuaDataType getDataType() {
        return dataType;
    }

    public AttributeId getAttributeId() {
        return attributeId;
    }

    public Map<String, String> getConfig() {
        return config;
    }

    @Override
    public String getAddressString() {
        String address = String.format("ns=%d;%s=%s", namespace, identifierType.getValue(), identifier);
        if (attributeId != AttributeId.Value) {
            address += ";a=" + attributeId.name();
        }
        if (dataType != null) {
            address += ";" + dataType.name();
        }
        return address;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.valueOf(dataType.name());
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return PlcSubscriptionTag.super.getArrayInfo();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof OpcuaTag)) {
            return false;
        }
        OpcuaTag that = (OpcuaTag) o;
        return namespace == that.namespace &&
            identifier.equals(that.identifier) &&
            identifierType == that.identifierType &&
            attributeId == that.attributeId &&
            config.equals(that.config);
    }

    @Override
    public int hashCode() {
        return Objects.hash(namespace, identifier, identifierType, attributeId, config);
    }

    @Override
    public String toString() {
        return "OpcuaTag{" +
            " namespace=" + namespace +
            " identifierType=" + identifierType.getValue() +
            " identifier=" + identifier +
            " attributeId=" + attributeId.name() +
            " config=" + config +
            '}';
    }

    @Override
    public PlcSubscriptionType getPlcSubscriptionType() {
        return null;
    }

    @Override
    public Optional<Duration> getDuration() {
        return Optional.empty();
    }

}
