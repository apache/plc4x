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
package org.apache.plc4x.protocol.cbus;

import org.apache.plc4x.plugins.codegenerator.language.mspec.parser.MessageFormatParser;
import org.apache.plc4x.plugins.codegenerator.language.mspec.protocol.ProtocolHelpers;
import org.apache.plc4x.plugins.codegenerator.language.mspec.protocol.ValidatableTypeContext;
import org.apache.plc4x.plugins.codegenerator.protocol.Protocol;
import org.apache.plc4x.plugins.codegenerator.protocol.TypeContext;
import org.apache.plc4x.plugins.codegenerator.types.exceptions.GenerationException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class CBusProtocol implements Protocol, ProtocolHelpers {

    private static final Logger LOGGER = LoggerFactory.getLogger(CBusProtocol.class);

    @Override
    public String getName() {
        return "c-bus";
    }

    @Override
    public TypeContext getTypeContext() throws GenerationException {
        ValidatableTypeContext typeContext;
        String mspecName;

        mspecName = getName() + "_lighting_application";
        LOGGER.info("Parsing: {}.mspec", mspecName);
        typeContext = new MessageFormatParser().parse(getMspecStream(mspecName));

        mspecName = getName() + "_security_application";
        LOGGER.info("Parsing: {}.mspec", mspecName);
        typeContext = new MessageFormatParser().parse(getMspecStream(mspecName), typeContext);

        mspecName = getName() + "_metering_application";
        LOGGER.info("Parsing: {}.mspec", mspecName);
        typeContext = new MessageFormatParser().parse(getMspecStream(mspecName), typeContext);

        mspecName = getName() + "_trigger_control_application";
        LOGGER.info("Parsing: {}.mspec", mspecName);
        typeContext = new MessageFormatParser().parse(getMspecStream(mspecName), typeContext);

        mspecName = getName() + "_enable_control_application";
        LOGGER.info("Parsing: {}.mspec", mspecName);
        typeContext = new MessageFormatParser().parse(getMspecStream(mspecName), typeContext);

        mspecName = getName() + "_temperature_broadcast_application";
        LOGGER.info("Parsing: {}.mspec", mspecName);
        typeContext = new MessageFormatParser().parse(getMspecStream(mspecName), typeContext);

        mspecName = getName() + "_access_control_application";
        LOGGER.info("Parsing: {}.mspec", mspecName);
        typeContext = new MessageFormatParser().parse(getMspecStream(mspecName), typeContext);

        mspecName = getName() + "_media_transport_control_application";
        LOGGER.info("Parsing: {}.mspec", mspecName);
        typeContext = new MessageFormatParser().parse(getMspecStream(mspecName), typeContext);

        mspecName = getName() + "_clock_and_timekeeping_application";
        LOGGER.info("Parsing: {}.mspec", mspecName);
        typeContext = new MessageFormatParser().parse(getMspecStream(mspecName), typeContext);

        mspecName = getName() + "_telephony_application";
        LOGGER.info("Parsing: {}.mspec", mspecName);
        typeContext = new MessageFormatParser().parse(getMspecStream(mspecName), typeContext);

        mspecName = getName() + "_air_conditioning_application";
        LOGGER.info("Parsing: {}.mspec", mspecName);
        typeContext = new MessageFormatParser().parse(getMspecStream(mspecName), typeContext);

        mspecName = getName() + "_measurement_application";
        LOGGER.info("Parsing: {}.mspec", mspecName);
        typeContext = new MessageFormatParser().parse(getMspecStream(mspecName), typeContext);

        mspecName = getName() + "_error_reporting_application";
        LOGGER.info("Parsing: {}.mspec", mspecName);
        typeContext = new MessageFormatParser().parse(getMspecStream(mspecName), typeContext);

        LOGGER.info("Parsing: c-bus.mspec");
        typeContext = new MessageFormatParser().parse(getMspecStream(), typeContext);

        typeContext.validate();

        return typeContext;
    }

}
