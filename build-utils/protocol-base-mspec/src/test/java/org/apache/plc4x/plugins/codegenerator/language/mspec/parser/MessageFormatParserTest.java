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

import org.antlr.v4.runtime.CommonTokenStream;
import org.antlr.v4.runtime.tree.ParseTreeWalker;
import org.apache.commons.io.IOUtils;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecLexer;
import org.apache.plc4x.plugins.codegenerator.language.mspec.MSpecParser;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.io.InputStream;
import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.mock;

// TODO: implement me
@ExtendWith(MockitoExtension.class)
class MessageFormatParserTest {

    //@InjectMocks
    MessageFormatParser SUT;

    @Mock
    MSpecLexer mSpecLexer;

    @Mock
    CommonTokenStream commonTokenStream;

    @Mock
    MSpecParser mSpecParser;

    @Mock
    ParseTreeWalker parseTreeWalker;

    @Mock
    MessageFormatListener messageFormatListener;

    @BeforeEach
    void setUp() {
        SUT = new MessageFormatParser(
            (_ignore) -> mSpecLexer,
            (_ignore) -> commonTokenStream,
            (_ignore) -> mSpecParser,
            () -> parseTreeWalker,
            () -> messageFormatListener
        );
    }

    @Test
    void parseNull() {
        assertThrows(NullPointerException.class, () -> SUT.parse(null));
    }

    @Disabled("mockito broken because of NPE in REAL code")
    @Test
    void parseSomething() {
        // TODO: seems like mockito is broken somehow...
        MSpecParser mSpecParser = mock(MSpecParser.class);
        // TODO: ... because this throws a NPE... (somehow call ends up in real code)
        mSpecParser.file();
        InputStream is = IOUtils.toInputStream("test", StandardCharsets.UTF_8);
        SUT.parse(is);
    }
}