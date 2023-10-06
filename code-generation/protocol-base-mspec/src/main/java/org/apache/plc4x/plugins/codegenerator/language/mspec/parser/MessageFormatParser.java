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
package org.apache.plc4x.plugins.codegenerator.language.mspec.parser;

import org.antlr.v4.runtime.CharStreams;
import org.antlr.v4.runtime.CommonTokenStream;
import org.antlr.v4.runtime.tree.ParseTreeWalker;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecLexer;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecParser;
import org.apache.plc4x.plugins.codegenerator.language.mspec.protocol.ValidatableTypeContext;
import org.apache.plc4x.plugins.codegenerator.protocol.TypeContext;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.io.InputStream;
import java.util.List;
import java.util.Map;
import java.util.function.Consumer;

public class MessageFormatParser {

    private static final Logger LOGGER = LoggerFactory.getLogger(MessageFormatParser.class);

    public ValidatableTypeContext parse(InputStream source) {
        return parse(source, null);
    }

    public ValidatableTypeContext parse(InputStream source, TypeContext existingTypeContext) {
        LOGGER.debug("Parsing: {}", source);
        MSpecLexer lexer;
        try {
            lexer = new MSpecLexer(CharStreams.fromStream(source));
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        MessageFormatListener listener;
        if (existingTypeContext == null) {
            listener = new MessageFormatListener();
        } else {
            if (LOGGER.isDebugEnabled()) {
                Map<String, TypeDefinition> exitingTypeDefinitions = existingTypeContext.getTypeDefinitions();
                if (exitingTypeDefinitions != null) {
                    LOGGER.debug("Continue with {} exitingTypeDefinitions", exitingTypeDefinitions.size());
                }
                Map<String, List<Consumer<TypeDefinition>>> unresolvedTypeReferences = existingTypeContext.getUnresolvedTypeReferences();
                if (unresolvedTypeReferences != null) {
                    LOGGER.debug("Continue with {} unresolvedTypeReferences", unresolvedTypeReferences.size());
                }
            }

            listener = new MessageFormatListener(existingTypeContext);
        }

        new ParseTreeWalker().walk(listener, new MSpecParser(new CommonTokenStream(lexer)).file());
        LOGGER.info("Checking for open consumers");
        listener.typeDefinitionConsumers.forEach((key, value) -> LOGGER.warn("{} has {} open consumers", key, value.size()));
        return new ValidatableTypeContext() {

            @Override
            public Map<String, TypeDefinition> getTypeDefinitions() {
                return listener.types;
            }

            @Override
            public Map<String, List<Consumer<TypeDefinition>>> getUnresolvedTypeReferences() {
                return listener.typeDefinitionConsumers;
            }
        };
    }

}
