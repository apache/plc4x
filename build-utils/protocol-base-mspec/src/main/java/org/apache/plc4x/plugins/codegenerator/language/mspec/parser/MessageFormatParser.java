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
import org.antlr.v4.runtime.tree.ParseTree;
import org.antlr.v4.runtime.tree.ParseTreeWalker;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecLexer;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecParser;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;

import java.io.IOException;
import java.io.InputStream;
import java.util.Map;
import java.util.function.Function;
import java.util.function.Supplier;

public class MessageFormatParser {

    // TODO: can be used a instance once thread safety is asserted
    private Function<CharStream, MSpecLexer> lexerBuilder;

    // TODO: can be used a instance once thread safety is asserted
    private Function<MSpecLexer, CommonTokenStream> commonTokenStreamBuilder;

    // TODO: can be used a instance once thread safety is asserted
    private Function<CommonTokenStream, MSpecParser> commonTokenStreamMSpecParserBuilder;

    // TODO: can be used a instance once thread safety is asserted
    private Supplier<ParseTreeWalker> voidParseTreeWalkerBuilder;

    // TODO: can be used a instance once thread safety is asserted
    private Supplier<MessageFormatListener> messageFormatListenerBuilder;

    public MessageFormatParser() {
        lexerBuilder = MSpecLexer::new;
        commonTokenStreamBuilder = CommonTokenStream::new;
        commonTokenStreamMSpecParserBuilder = MSpecParser::new;
        voidParseTreeWalkerBuilder = ParseTreeWalker::new;
        messageFormatListenerBuilder = MessageFormatListener::new;
    }

    MessageFormatParser(Function<CharStream, MSpecLexer> lexerBuilder, Function<MSpecLexer, CommonTokenStream> commonTokenStreamBuilder, Function<CommonTokenStream, MSpecParser> commonTokenStreamMSpecParserBuilder, Supplier<ParseTreeWalker> voidParseTreeWalkerBuilder, Supplier<MessageFormatListener> messageFormatListenerBuilder) {
        this.lexerBuilder = lexerBuilder;
        this.commonTokenStreamBuilder = commonTokenStreamBuilder;
        this.commonTokenStreamMSpecParserBuilder = commonTokenStreamMSpecParserBuilder;
        this.voidParseTreeWalkerBuilder = voidParseTreeWalkerBuilder;
        this.messageFormatListenerBuilder = messageFormatListenerBuilder;
    }

    public Map<String, TypeDefinition> parse(InputStream source) {
        MSpecLexer lexer;
        try {
            lexer = lexerBuilder.apply(CharStreams.fromStream(source));
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        CommonTokenStream tokens = commonTokenStreamBuilder.apply(lexer);
        MSpecParser parser = commonTokenStreamMSpecParserBuilder.apply(tokens);
        ParseTree tree = parser.file();
        ParseTreeWalker walker = voidParseTreeWalkerBuilder.get();
        MessageFormatListener listener = messageFormatListenerBuilder.get();
        walker.walk(listener, tree);
        return listener.getTypes();
    }

}
