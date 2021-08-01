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
package org.apache.plc4x.plugins.codegenerator.language.mspec.parser;

import org.antlr.v4.runtime.CharStreams;
import org.antlr.v4.runtime.CommonTokenStream;
import org.antlr.v4.runtime.tree.ParseTreeWalker;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecLexer;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecParser;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;

import java.io.IOException;
import java.io.InputStream;
import java.util.Map;

public class MessageFormatParser {

    public Map<String, TypeDefinition> parse(InputStream source) {
        MSpecLexer lexer;
        try {
            lexer = new MSpecLexer(CharStreams.fromStream(source));
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        MessageFormatListener listener = new MessageFormatListener();
        new ParseTreeWalker().walk(listener, new MSpecParser(new CommonTokenStream(lexer)).file());
        return listener.getTypes();
    }

}
