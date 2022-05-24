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
import org.apache.plc4x.plugins.codegenerator.language.mspec.protocol.ProtocolHelpers;
import org.apache.plc4x.plugins.codegenerator.language.mspec.protocol.ValidatableTypeContext;
import org.apache.plc4x.plugins.codegenerator.protocol.Protocol;
import org.apache.plc4x.plugins.codegenerator.protocol.TypeContext;
import org.apache.plc4x.plugins.codegenerator.types.exceptions.GenerationException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class BacNetIpProtocol implements Protocol, ProtocolHelpers {

    private static final Logger LOGGER = LoggerFactory.getLogger(BacNetIpProtocol.class);

    @Override
    public String getName() {
        return "bacnetip";
    }

    @Override
    public TypeContext getTypeContext() throws GenerationException {
        ValidatableTypeContext typeContext;

        LOGGER.info("Parsing: bacnet-tags.mspec");
        typeContext = new MessageFormatParser().parse(getMspecStream("bacnet-tags"));

        LOGGER.info("Parsing: bacnet-private-enums.mspec");
        typeContext = new MessageFormatParser().parse(getMspecStream("bacnet-private-enums"), typeContext);

        LOGGER.info("Parsing: bacnet-enums.mspec");
        typeContext = new MessageFormatParser().parse(getMspecStream("bacnet-enums"), typeContext);

        LOGGER.info("Parsing: bacnetip.mspec");
        typeContext = new MessageFormatParser().parse(getMspecStream(), typeContext);

        // TODO: those should work above bacnetip.mspec but somehow if we move them we get a concurrent modification exception... debug that.
        LOGGER.info("Parsing: bacnet-vendorids.mspec");
        typeContext = new MessageFormatParser().parse(getMspecStream("bacnet-vendorids"), typeContext);

        LOGGER.info("Parsing: bacnet-private-enums-tagged.mspec");
        typeContext = new MessageFormatParser().parse(getMspecStream("bacnet-private-enums-tagged"), typeContext);

        LOGGER.info("Parsing: bacnet-enums-tagged.mspec");
        typeContext = new MessageFormatParser().parse(getMspecStream("bacnet-enums-tagged"), typeContext);

        typeContext.validate();

        return typeContext;
    }

}
