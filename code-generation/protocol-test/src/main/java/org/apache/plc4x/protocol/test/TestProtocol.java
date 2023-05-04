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
package org.apache.plc4x.protocol.test;

import org.apache.plc4x.plugins.codegenerator.language.mspec.parser.MessageFormatParser;
import org.apache.plc4x.plugins.codegenerator.protocol.Protocol;
import org.apache.plc4x.plugins.codegenerator.protocol.TypeContext;
import org.apache.plc4x.plugins.codegenerator.types.exceptions.GenerationException;

import java.io.InputStream;

public class TestProtocol implements Protocol {

    @Override
    public String getName() {
        return "test";
    }

    @Override
    public TypeContext getTypeContext() throws GenerationException {
        InputStream schemaInputStream = TestProtocol.class.getResourceAsStream("/protocols/test/test.mspec");
        if (schemaInputStream == null) {
            throw new GenerationException("Error loading message-format schema for protocol '" + getName() + "'");
        }
        TypeContext typeContext = new MessageFormatParser().parse(schemaInputStream);
        if (typeContext.getUnresolvedTypeReferences().size() > 0) {
            throw new GenerationException("Unresolved types left: " + typeContext.getUnresolvedTypeReferences());
        }
        return typeContext;
    }

}
