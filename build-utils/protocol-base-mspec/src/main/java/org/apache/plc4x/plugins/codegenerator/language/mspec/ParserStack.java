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
package org.apache.plc4x.plugins.codegenerator.language.mspec;

import java.util.Map;
import java.util.Stack;

/**
 * Helper to hold information about parsing of rules.
 */
public class ParserStack {

    private final Stack<ParserCall> lines = new Stack<>();

    public void push(int line, int character, Map<String, Object> context) {
        lines.push(new ParserCall(line, character, context));
    }

    public void pop() {
        lines.pop();
    }

    public void clear() {
        lines.clear();
    }

    @Override
    public String toString() {
        final StringBuilder sb = new StringBuilder("ParserStack:");
        for (int index = 0, linesSize = lines.size(); index < linesSize; index++) {
            ParserCall call = lines.get(index);
            sb.append(index).append(". ").append(call).append("\n");
        }
        return sb.toString();
    }

    private class ParserCall {
        private final int line;
        private final int character;
        private final Map<String, Object> context;

        private ParserCall(int line, int character, Map<String, Object> context) {
            this.line = line;
            this.character = character;
            this.context = context;
        }

        @Override
        public String toString() {
            return "line: " + line + ":" + character + (context.isEmpty() ? "" : ", context: " + context);
        }
    }
}
