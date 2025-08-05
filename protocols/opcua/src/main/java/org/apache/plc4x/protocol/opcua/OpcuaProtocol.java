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
package org.apache.plc4x.protocol.opcua;

import org.apache.plc4x.plugins.codegenerator.language.mspec.parser.MessageFormatParser;
import org.apache.plc4x.plugins.codegenerator.language.mspec.protocol.ProtocolHelpers;
import org.apache.plc4x.plugins.codegenerator.language.mspec.protocol.ValidatableTypeContext;
import org.apache.plc4x.plugins.codegenerator.protocol.Protocol;
import org.apache.plc4x.plugins.codegenerator.protocol.TypeContext;
import org.apache.plc4x.plugins.codegenerator.types.exceptions.GenerationException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class OpcuaProtocol implements Protocol, ProtocolHelpers {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaProtocol.class);

    @Override
    public String getName() {
        return "opcua";
    }

    @Override
    public TypeContext getTypeContext() throws GenerationException {
        LOGGER.info("Parsing: opc-manual.mspec");
        ValidatableTypeContext typeContext;

        typeContext = new MessageFormatParser().parse(getMspecStream("opc-manual"));

        LOGGER.info("Parsing: opc-services.mspec");
        typeContext = new MessageFormatParser().parse(getMspecStream("opc-services"), typeContext);

        LOGGER.info("Parsing: opc-attribute.mspec");
        typeContext = new MessageFormatParser().parse(getMspecStream("opc-attribute"), typeContext);

        LOGGER.info("Parsing: opc-status.mspec");
        typeContext = new MessageFormatParser().parse(getMspecStream("opc-status"), typeContext);

        LOGGER.info("Parsing: opc-types.mspec");
        typeContext = new MessageFormatParser().parse(getMspecStream("opc-types"), typeContext);

        typeContext.validate();

        return typeContext;
    }
}
