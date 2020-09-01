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

package org.apache.plc4x.plugins.codegenerator.language.mspec.parser;

import org.antlr.v4.runtime.CharStream;
import org.antlr.v4.runtime.CharStreams;
import org.antlr.v4.runtime.CommonTokenStream;
import org.antlr.v4.runtime.tree.ParseTreeWalker;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecLexer;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecParser;
import org.apache.plc4x.plugins.codegenerator.language.mspec.ParserStack;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.URISyntaxException;
import java.net.URL;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Map;

public class MessageFormatParser {

    private final ParserStack parserStack;

    public MessageFormatParser() {
        this(new ParserStack());
    }

    public MessageFormatParser(ParserStack stack) {
        this.parserStack = stack;
    }

    public Map<String, TypeDefinition> parse(InputStream source) {
        try {
            return parse(CharStreams.fromStream(source));
        } catch (IOException e) {
            throw new RuntimeException("Could not read source stream", e);
        }
    }

    public Map<String, TypeDefinition> parse(InputStream source, String name) {
        try {
            return parse(CharStreams.fromReader(new InputStreamReader(source), name));
        } catch (IOException e) {
            throw new RuntimeException("Could not read source stream", e);
        }
    }

    public Map<String, TypeDefinition> parse(Path source) {
        try {
            return parse(CharStreams.fromPath(source));
        } catch (IOException e) {
            throw new RuntimeException("Could not read source path", e);
        }
    }

    public Map<String, TypeDefinition> parse(URL source) {
        try {
            return parse(Paths.get(source.toURI()));
        } catch (URISyntaxException e) {
            throw new RuntimeException("Invalid source URI", e);
        }
    }

    public Map<String, TypeDefinition> parse(File source) {
        try {
            return parse(CharStreams.fromFileName(source.getAbsolutePath()));
        } catch (IOException e) {
            throw new RuntimeException("Could not open file", e);
        }
    }

    private Map<String, TypeDefinition> parse(CharStream input) {
        MSpecLexer lexer = new MSpecLexer(input);
        MessageFormatListener listener = new MessageFormatListener(parserStack, MSpecParser.ruleNames);
        new ParseTreeWalker().walk(listener, new MSpecParser(new CommonTokenStream(lexer)).file());
        return listener.getTypes();
    }

}
