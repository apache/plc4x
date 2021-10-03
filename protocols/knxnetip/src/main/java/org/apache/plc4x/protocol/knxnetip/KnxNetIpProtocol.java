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
package org.apache.plc4x.protocol.knxnetip;

import org.apache.plc4x.plugins.codegenerator.language.mspec.parser.MessageFormatParser;
import org.apache.plc4x.plugins.codegenerator.protocol.Protocol;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.exceptions.GenerationException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.InputStream;
import java.util.LinkedHashMap;
import java.util.Map;

public class KnxNetIpProtocol implements Protocol {

    private static final Logger LOGGER = LoggerFactory.getLogger(KnxNetIpProtocol.class);

    @Override
    public String getName() {
        return "knxnetip";
    }

    @Override
    public Map<String, TypeDefinition> getTypeDefinitions() throws GenerationException {
        LOGGER.info("Parsing: knxnetip.mspec");
        InputStream schemaInputStream = KnxNetIpProtocol.class.getResourceAsStream(
            "/protocols/knxnetip/knxnetip.mspec");
        if (schemaInputStream == null) {
            throw new GenerationException("Error loading message-format schema for protocol '" + getName() + "'");
        }
        Map<String, TypeDefinition> typeDefinitionMap =
            new LinkedHashMap<>(new MessageFormatParser().parse(schemaInputStream));

        LOGGER.info("Parsing: knx-master-data.mspec");
        InputStream masterDataInputStream = KnxNetIpProtocol.class.getResourceAsStream(
            "/protocols/knxnetip/knx-master-data.mspec");
        if (masterDataInputStream == null) {
            throw new GenerationException("Error loading knx-master-data schema for protocol '" + getName() + "'");
        }
        typeDefinitionMap.putAll(new MessageFormatParser().parse(masterDataInputStream));

        LOGGER.info("Parsing: device-info.mspec");
        InputStream deviceDataInputStream = KnxNetIpProtocol.class.getResourceAsStream(
            "/protocols/knxnetip/device-info.mspec");
        if (deviceDataInputStream == null) {
            throw new GenerationException("Error loading device-info schema for protocol '" + getName() + "'");
        }
        typeDefinitionMap.putAll(new MessageFormatParser().parse(deviceDataInputStream));

        return typeDefinitionMap;
    }

}
