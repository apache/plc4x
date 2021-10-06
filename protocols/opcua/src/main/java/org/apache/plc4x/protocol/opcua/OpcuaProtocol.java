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
package org.apache.plc4x.protocol.opcua;

import org.apache.plc4x.plugins.codegenerator.language.mspec.parser.MessageFormatParser;
import org.apache.plc4x.plugins.codegenerator.protocol.Protocol;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.exceptions.GenerationException;

import java.io.InputStream;
import java.util.LinkedHashMap;
import java.util.Map;

public class OpcuaProtocol implements Protocol {

    @Override
    public String getName() {
        return "opcua";
    }

    @Override
    public Map<String, TypeDefinition> getTypeDefinitions() throws GenerationException {
        System.out.println("Parsing: opc-manual.mspec");
        InputStream manualInputStream = OpcuaProtocol.class.getResourceAsStream(
            "/protocols/opcua/opc-manual.mspec");
        if(manualInputStream == null) {
            throw new GenerationException("Error loading message-format schema for protocol '" + getName() + "'");
        }
        Map<String, TypeDefinition> typeDefinitionMap =
            new LinkedHashMap<>(new MessageFormatParser().parse(manualInputStream));

        System.out.println("Parsing: opc-services.mspec");
        InputStream servicesInputStream = OpcuaProtocol.class.getResourceAsStream(
            "/protocols/opcua/opc-services.mspec");
        if(servicesInputStream == null) {
            throw new GenerationException("Error loading message-format schema for protocol '" + getName() + "'");
        }
        typeDefinitionMap.putAll(new MessageFormatParser().parse(servicesInputStream));

        System.out.println("Parsing: opc-status.mspec");
        InputStream statusInputStream = OpcuaProtocol.class.getResourceAsStream(
            "/protocols/opcua/opc-status.mspec");
        if(statusInputStream == null) {
            throw new GenerationException("Error loading message-format schema for protocol '" + getName() + "'");
        }
        typeDefinitionMap.putAll(new MessageFormatParser().parse(statusInputStream));

        System.out.println("Parsing: opc-types.mspec");
        InputStream typesInputStream = OpcuaProtocol.class.getResourceAsStream(
            "/protocols/opcua/opc-types.mspec");
        if(typesInputStream == null) {
            throw new GenerationException("Error loading message-format schema for protocol '" + getName() + "'");
        }
        typeDefinitionMap.putAll(new MessageFormatParser().parse(typesInputStream));
        return typeDefinitionMap;
    }
}
