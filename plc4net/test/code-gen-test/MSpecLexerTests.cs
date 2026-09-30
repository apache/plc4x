//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

using System.Linq;
using Antlr4.Runtime;
using org.apache.plc4net.tools.codegen.grammar;
using Xunit;

namespace org.apache.plc4net.test.codegen
{
    /// <summary>
    /// Pins the column-zero predicate of the mspec lexer's <c>EmptyLine</c> rule.
    /// The grammar is the upstream one and spells the predicate the way the ANTLR
    /// Java runtime does, <c>getCharPositionInLine()</c>; <c>MSpecLexer.Predicates.cs</c>
    /// supplies that method for the C# lexer. A line break that starts a line must
    /// lex as <c>EmptyLine</c>, one that ends a line of text as <c>NEWLINE</c>.
    /// </summary>
    public class MSpecLexerTests
    {
        private static int[] LineBreakTokenTypes(string text)
        {
            var lexer = new MSpecLexer(CharStreams.fromString(text));
            return lexer.GetAllTokens()
                .Select(token => token.Type)
                .Where(type => type == MSpecLexer.EmptyLine || type == MSpecLexer.NEWLINE)
                .ToArray();
        }

        [Fact]
        public void Line_break_at_column_zero_is_an_EmptyLine()
        {
            Assert.Equal(new[] { MSpecLexer.EmptyLine }, LineBreakTokenTypes("\n"));
        }

        [Fact]
        public void Whitespace_only_line_is_an_EmptyLine()
        {
            Assert.Equal(new[] { MSpecLexer.EmptyLine }, LineBreakTokenTypes("  \t\n"));
        }

        [Fact]
        public void Line_break_after_text_is_a_NEWLINE()
        {
            Assert.Equal(new[] { MSpecLexer.NEWLINE }, LineBreakTokenTypes("x\n"));
        }

        [Fact]
        public void Crlf_line_breaks_follow_the_same_rule()
        {
            Assert.Equal(new[] { MSpecLexer.EmptyLine, MSpecLexer.NEWLINE }, LineBreakTokenTypes("\r\nx\r\n"));
        }

        [Fact]
        public void Blank_line_between_two_lines_of_text()
        {
            // "a\n" ends a line of text; the second "\n" starts, and ends, a line of its own.
            Assert.Equal(
                new[] { MSpecLexer.NEWLINE, MSpecLexer.EmptyLine, MSpecLexer.NEWLINE },
                LineBreakTokenTypes("a\n\nb\n"));
        }
    }
}