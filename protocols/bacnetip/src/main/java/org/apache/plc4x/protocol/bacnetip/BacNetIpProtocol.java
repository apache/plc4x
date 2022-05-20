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
package org.apache.plc4x.protocol.bacnetip;

import org.apache.plc4x.plugins.codegenerator.language.mspec.parser.MessageFormatParser;
import org.apache.plc4x.plugins.codegenerator.protocol.Protocol;
import org.apache.plc4x.plugins.codegenerator.protocol.TypeContext;
import org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.exceptions.GenerationException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.InputStream;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.function.Consumer;

public class BacNetIpProtocol implements Protocol {
    private static final Logger LOGGER = LoggerFactory.getLogger(BacNetIpProtocol.class);

    @Override
    public String getName() {
        return "bacnetip";
    }

    @Override
    public TypeContext getTypeContext() throws GenerationException {
        LOGGER.info("Parsing: bacnet-vendorids.mspec");
        InputStream bacnetVendorIdsSchemaInputStream = BacNetIpProtocol.class.getResourceAsStream(
            "/protocols/bacnetip/bacnet-vendorids.mspec");
        if (bacnetVendorIdsSchemaInputStream == null) {
            throw new GenerationException("Error loading vendorId schema for protocol '" + getName() + "'");
        }
        Map<String, TypeDefinition> typeDefinitionMap = new LinkedHashMap<>();
        TypeContext typeContext;

        typeContext = new MessageFormatParser().parse(bacnetVendorIdsSchemaInputStream);
        typeDefinitionMap.putAll(typeContext.getTypeDefinitions());

        LOGGER.info("Parsing: bacnet-private-enums.mspec");
        InputStream bacnetPrivateEnumsSchemaInputStream = BacNetIpProtocol.class.getResourceAsStream(
            "/protocols/bacnetip/bacnet-private-enums.mspec");
        if (bacnetPrivateEnumsSchemaInputStream == null) {
            throw new GenerationException("Error loading private enum schema for protocol '" + getName() + "'");
        }
        typeContext = new MessageFormatParser().parse(bacnetPrivateEnumsSchemaInputStream, typeDefinitionMap, typeContext.getUnresolvedTypeReferences());
        typeDefinitionMap.putAll(typeContext.getTypeDefinitions());

        LOGGER.info("Parsing: bacnet-enums.mspec");
        InputStream bacnetEnumsSchemaInputStream = BacNetIpProtocol.class.getResourceAsStream(
            "/protocols/bacnetip/bacnet-enums.mspec");
        if (bacnetEnumsSchemaInputStream == null) {
            throw new GenerationException("Error loading enum schema for protocol '" + getName() + "'");
        }
        typeContext = new MessageFormatParser().parse(bacnetEnumsSchemaInputStream, typeDefinitionMap, typeContext.getUnresolvedTypeReferences());
        typeDefinitionMap.putAll(typeContext.getTypeDefinitions());

        LOGGER.info("Parsing: bacnetip.mspec");
        InputStream bacnetipSchemaInputStream = BacNetIpProtocol.class.getResourceAsStream(
            "/protocols/bacnetip/bacnetip.mspec");
        if (bacnetipSchemaInputStream == null) {
            throw new GenerationException("Error loading schema for protocol '" + getName() + "'");
        }
        typeContext = new MessageFormatParser().parse(bacnetipSchemaInputStream, typeDefinitionMap, typeContext.getUnresolvedTypeReferences());
        typeDefinitionMap.putAll(typeContext.getTypeDefinitions());

        if (typeContext.getUnresolvedTypeReferences().size() > 0) {
            throw new GenerationException("Unresolved types left: " + typeContext.getUnresolvedTypeReferences());
        }

        return new TypeContext() {
            @Override
            public Map<String, TypeDefinition> getTypeDefinitions() {
                return typeDefinitionMap;
            }

            @Override
            public Map<String, List<Consumer<TypeDefinition>>> getUnresolvedTypeReferences() {
                return Collections.emptyMap();
            }
        };
    }

}
