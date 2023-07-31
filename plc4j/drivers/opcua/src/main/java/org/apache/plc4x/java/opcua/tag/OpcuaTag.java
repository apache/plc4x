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

import org.apache.commons.lang3.EnumUtils;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.opcua.readwrite.OpcuaDataType;
import org.apache.plc4x.java.opcua.readwrite.OpcuaIdentifierType;

import java.time.Duration;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class OpcuaTag implements PlcSubscriptionTag {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("^ns=(?<namespace>\\d+);(?<identifierType>[isgb])=(?<identifier>[^;]+)?(;(?<datatype>[a-zA-Z_]+))?");

    private final OpcuaIdentifierType identifierType;

    private final int namespace;

    private final String identifier;

    private final OpcuaDataType dataType;

    private OpcuaTag(Integer namespace, String identifier, OpcuaIdentifierType identifierType, OpcuaDataType dataType) {
        this.identifier = Objects.requireNonNull(identifier);
        this.identifierType = Objects.requireNonNull(identifierType);
        this.namespace = namespace != null ? namespace : 0;
        if (this.namespace < 0) {
            throw new IllegalArgumentException("namespace must be greater then zero. Was " + this.namespace);
        }
        this.dataType = dataType;
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

        return new OpcuaTag(namespace, identifier, identifierType, dataType);
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

    @Override
    public String getAddressString() {
        String address = String.format("ns=%d;%s=%s", namespace, identifierType.getValue(), identifier);
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
        return namespace == that.namespace && identifier.equals(that.identifier) && identifierType == that.identifierType;
    }

    @Override
    public int hashCode() {
        return Objects.hash(namespace);
    }

    @Override
    public String toString() {
        return "OpcuaTag{" +
            "namespace=" + namespace +
            "identifierType=" + identifierType.getValue() +
            "identifier=" + identifier +
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
