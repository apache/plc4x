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

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.opcua.protocol.model.OpcuaIdentifierType;

import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
/**
 * @author Matthias Milan Stlrljic
 * Created by Matthias Milan Stlrljic on 10.05.2019
 */
public class OpcuaField implements PlcField {
    //TODO: Add a correct regex definition of all the different sub-types of the identifiers --> requires perhaps individual type definitions
    public static final Pattern ADDRESS_PATTERN = Pattern.compile("^ns=(?<namespace>\\d+);((?<identifierType>[isgb])=((?<identifier>[\\w.\\-/=%_]+))?)");

    private final OpcuaIdentifierType identifierType;

    private final int namespace;

    private final String identifier;

    protected OpcuaField(int namespace, OpcuaIdentifierType identifierType, String identifier) {
        this.namespace = namespace;
        this.identifier = identifier;
        this.identifierType = identifierType;
        if (this.identifier == null || this.namespace < 0) {
            throw new IllegalArgumentException("Identifier can not be null or Namespace can not be lower then 0.");
        }
    }

    private OpcuaField(Integer namespace, String identifier, OpcuaIdentifierType identifierType) {
        this.identifier = Objects.requireNonNull(identifier);
        this.identifierType = Objects.requireNonNull(identifierType);
        this.namespace = namespace != null ? namespace : 0;
        if (this.namespace < 0) {
            throw new IllegalArgumentException("namespace must be greater then zero. Was " + this.namespace);
        }
    }

    public static OpcuaField of(String address) {
        Matcher matcher = ADDRESS_PATTERN.matcher(address);
        if (!matcher.matches()) {
            throw new PlcInvalidFieldException(address, ADDRESS_PATTERN, "{address}");
        }
        String identifier = matcher.group("identifier");

        String identifierTypeString = matcher.group("identifierType");
        OpcuaIdentifierType identifierType = OpcuaIdentifierType.fromString(identifierTypeString);

        String namespaceString = matcher.group("namespace");
        Integer namespace = namespaceString != null ? Integer.valueOf(namespaceString) : 0;

        return new OpcuaField(namespace, identifier, identifierType);
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
            "identifierType=" + identifierType.getText() +
            "identifier=" + identifier +
            '}';
    }
}
