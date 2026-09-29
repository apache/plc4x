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

using System;
using System.Collections.Generic;
using System.IO;
using Antlr4.Runtime;
using org.apache.plc4net.tools.codegen.grammar;
using Xunit;

namespace org.apache.plc4net.test.codegen;

public class ParserSmokeTests
{
    public static TheoryData<string> MSpecFixtures => new()
    {
        "mspec.example",
        "mspec.example2",
        "mspec.example3"
    };

    public static TheoryData<string> Expressions => new()
    {
        "payload.lengthInBytes+4",
        "CAST(parameter,S7ParameterUserData).items(hurz)[0]",
        "enabled && (length >= 4 ? true : false)",
        "0x2a << 1"
    };

    [Theory]
    [MemberData(nameof(MSpecFixtures))]
    public void ParsesExistingMSpecFixture(string fixture)
    {
        string input = File.ReadAllText(Path.Combine(AppContext.BaseDirectory, "Fixtures", fixture));

        AssertMSpecParses(input);
    }

    [Theory]
    [MemberData(nameof(Expressions))]
    public void ParsesRepresentativeExpression(string input)
    {
        var lexer = new ExpressionLexer(CharStreams.fromString(input));
        var lexerErrors = new CollectingErrorListener<int>();
        lexer.RemoveErrorListeners();
        lexer.AddErrorListener(lexerErrors);
        var tokens = new CommonTokenStream(lexer);
        var parser = new ExpressionParser(tokens);
        var parserErrors = new CollectingErrorListener<IToken>();
        parser.RemoveErrorListeners();
        parser.AddErrorListener(parserErrors);

        parser.expressionString();

        Assert.Empty(lexerErrors.Errors);
        Assert.Empty(parserErrors.Errors);
        Assert.Equal(TokenConstants.EOF, parser.CurrentToken.Type);
    }

    [Fact]
    public void ParsesIndentedEmptyLineThroughCompatibilityPredicate()
    {
        const string input = "[type Sample\n    [simple uint 8 value]\n    \n]\n";
        var lexer = new MSpecLexer(CharStreams.fromString(input));
        var tokens = new CommonTokenStream(lexer);

        tokens.Fill();

        Assert.Contains(tokens.GetTokens(), token =>
            token.Type == MSpecLexer.EmptyLine && token.Text == "    \n");
        AssertMSpecParses(input);
    }

    private static void AssertMSpecParses(string input)
    {
        var lexer = new MSpecLexer(CharStreams.fromString(input));
        var lexerErrors = new CollectingErrorListener<int>();
        lexer.RemoveErrorListeners();
        lexer.AddErrorListener(lexerErrors);
        var tokens = new CommonTokenStream(lexer);
        var parser = new MSpecParser(tokens);
        var parserErrors = new CollectingErrorListener<IToken>();
        parser.RemoveErrorListeners();
        parser.AddErrorListener(parserErrors);

        parser.file();

        Assert.Empty(lexerErrors.Errors);
        Assert.Empty(parserErrors.Errors);
        Assert.Equal(TokenConstants.EOF, parser.CurrentToken.Type);
    }

    private sealed class CollectingErrorListener<TSymbol> : IAntlrErrorListener<TSymbol>
    {
        public List<string> Errors { get; } = [];

        public void SyntaxError(
            TextWriter output,
            IRecognizer recognizer,
            TSymbol offendingSymbol,
            int line,
            int charPositionInLine,
            string msg,
            RecognitionException e)
        {
            Errors.Add($"{line}:{charPositionInLine}: {msg}");
        }
    }
}