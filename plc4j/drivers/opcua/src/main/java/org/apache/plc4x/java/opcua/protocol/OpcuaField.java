/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.opcua.protocol;

import org.apache.commons.lang3.EnumUtils;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionField;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.opcua.readwrite.types.OpcuaIdentifierType;
import org.apache.plc4x.java.opcua.readwrite.types.OpcuaDataType;

import java.time.Duration;
import java.util.Objects;
import java.util.Optional;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 */
public class OpcuaField implements PlcSubscriptionField {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("^ns=(?<namespace>\\d+);(?<identifierType>[isgb])=((?<identifier>[^;]+))?(;(?<datatype>[a-zA-Z_]+))?");

    private final OpcuaIdentifierType identifierType;

    private final int namespace;

    private final String identifier;

    private final OpcuaDataType dataType;

    protected OpcuaField(int namespace, OpcuaIdentifierType identifierType, String identifier, OpcuaDataType dataType) {
        this.namespace = namespace;
        this.identifier = identifier;
        this.identifierType = identifierType;
        if (this.identifier == null || this.namespace < 0) {
            throw new IllegalArgumentException("Identifier can not be null or Namespace can not be lower then 0.");
        }
        this.dataType = dataType;
    }

    private OpcuaField(Integer namespace, String identifier, OpcuaIdentifierType identifierType, OpcuaDataType dataType) {
        this.identifier = Objects.requireNonNull(identifier);
        this.identifierType = Objects.requireNonNull(identifierType);
        this.namespace = namespace != null ? namespace : 0;
        if (this.namespace < 0) {
            throw new IllegalArgumentException("namespace must be greater then zero. Was " + this.namespace);
        }
        this.dataType = dataType;
    }

    public static OpcuaField of(String address) {
        Matcher matcher = ADDRESS_PATTERN.matcher(address);
        if (!matcher.matches()) {
            throw new PlcInvalidFieldException(address, ADDRESS_PATTERN, "{address}");
        }
        String identifier = matcher.group("identifier");

        String identifierTypeString = matcher.group("identifierType");
        OpcuaIdentifierType identifierType = OpcuaIdentifierType.enumForValue(identifierTypeString);

        String namespaceString = matcher.group("namespace");
        Integer namespace = namespaceString != null ? Integer.valueOf(namespaceString) : 0;

        String dataTypeString = matcher.group("datatype") != null ? matcher.group("datatype").toUpperCase() : "NULL";
        if (!EnumUtils.isValidEnum(OpcuaDataType.class, dataTypeString)) {
            throw new PlcUnsupportedDataTypeException("Datatype " + dataTypeString + " is unsupported by this protocol");
        }
        OpcuaDataType dataType = OpcuaDataType.valueOf(dataTypeString);

        return new OpcuaField(namespace, identifier, identifierType, dataType);
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
    public String getPlcDataType() {
        return dataType.name();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof OpcuaField)) {
            return false;
        }
        OpcuaField that = (OpcuaField) o;
        return namespace == that.namespace && identifier.equals(that.identifier) && identifierType == that.identifierType;
    }

    @Override
    public int hashCode() {
        return Objects.hash(namespace);
    }

    @Override
    public String toString() {
        return "OpcuaField{" +
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
