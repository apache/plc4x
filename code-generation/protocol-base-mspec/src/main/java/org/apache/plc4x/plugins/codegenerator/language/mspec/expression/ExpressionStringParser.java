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
package org.apache.plc4x.plugins.codegenerator.language.mspec.expression;

import org.antlr.v4.runtime.CharStreams;
import org.antlr.v4.runtime.CommonTokenStream;
import org.antlr.v4.runtime.tree.ParseTreeWalker;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.io.IOException;
import java.io.InputStream;

public class ExpressionStringParser {

    public Term parse(InputStream source) {
        ExpressionLexer lexer;
        try {
            lexer = new ExpressionLexer(CharStreams.fromStream(source));
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        ExpressionStringListener listener = new ExpressionStringListener();
        new ParseTreeWalker().walk(listener, new ExpressionParser(new CommonTokenStream(lexer)).expressionString());
        return listener.getRoot();
    }

}
