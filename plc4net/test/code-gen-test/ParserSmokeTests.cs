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
using org.apache.plc4net.tools.codegen;
using org.apache.plc4net.tools.codegen.grammar;
using org.apache.plc4net.tools.codegen.model.terms;
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

    [Theory]
    [InlineData("42", "42")]
    [InlineData("0x0E", "0x0E")]
    [InlineData("true", "true")]
    [InlineData("false", "false")]
    [InlineData("null", "null")]
    [InlineData("1.5", "1.5")]
    [InlineData("1 + 2 * 3", "(1 + (2 * 3))")]
    [InlineData("(1 + 2) * 3", "((1 + 2) * 3)")]
    [InlineData("8 / 2 % 3", "((8 / 2) % 3)")]
    [InlineData("2 ^ 3 ^ 2", "(2 ^ (3 ^ 2))")]
    [InlineData("0x2a << 1", "(0x2a << 1)")]
    [InlineData("a >= b && c != d", "((a >= b) && (c != d))")]
    [InlineData("a & b | c", "((a & b) | c)")]
    [InlineData("a || b", "(a || b)")]
    [InlineData("!flag", "!flag")]
    [InlineData("-1", "-1")]
    [InlineData("a ? b : c", "(a ? b : c)")]
    public void BuildsExpressionTree(string input, string expected)
    {
        Assert.Equal(expected, MspecExpressionParser.Parse(input).ToString());
    }

    [Theory]
    [InlineData("0x0E")]
    [InlineData("0xFFFFFFFFFFFFFFFF")]
    [InlineData("0x10000000000000000")]
    public void PreservesHexadecimalLiteralTextWithoutSignedConversion(string input)
    {
        var literal = Assert.IsType<HexadecimalLiteral>(MspecExpressionParser.Parse(input));

        Assert.Equal(input, literal.Text);
    }

    [Fact]
    public void BuildsCallsIndexesAndMemberChains()
    {
        var call = Assert.IsType<VariableLiteral>(
            MspecExpressionParser.Parse("CAST(parameter, Type).items[0].value"));

        Assert.Equal("CAST", call.Name);
        Assert.NotNull(call.Args);
        Assert.Equal(2, call.Args.Count);
        Assert.Equal("items", call.Child?.Name);
        Assert.Equal("0", Assert.Single(call.Child!.Index).ToString());
        Assert.Equal("value", call.Child.Child?.Name);
    }

    [Theory]
    [InlineData("\"a\\\"b\"", "a\"b", "\"a\\\"b\"")]
    [InlineData("'a\\'b'", "a'b", "\"a'b\"")]
    [InlineData("\"C:\\\\temp\"", "C:\\temp", "\"C:\\\\temp\"")]
    public void UnquotesAndRendersStringLiterals(string input, string value, string rendered)
    {
        var literal = Assert.IsType<StringLiteral>(MspecExpressionParser.Parse(input));

        Assert.Equal(value, literal.Value);
        Assert.Equal(rendered, literal.ToString());
    }

    [Fact]
    public void DistinguishesEmptyCallsAndMultipleIndexes()
    {
        var value = Assert.IsType<VariableLiteral>(MspecExpressionParser.Parse("values()[row][column]"));

        Assert.NotNull(value.Args);
        Assert.Empty(value.Args);
        Assert.Equal(2, value.Index.Count);
    }

    [Theory]
    [InlineData("\"abc\"[0]", "\"abc\"[0]")]
    [InlineData("(left + right)[offset]", "(left + right)[offset]")]
    [InlineData("(-value)[0]", "(-value)[0]")]
    [InlineData("(!flag)[0]", "(!flag)[0]")]
    public void PreservesIndexesOnAnyExpression(string input, string expected)
    {
        Assert.Equal(expected, MspecExpressionParser.Parse(input).ToString());
    }

    [Theory]
    [InlineData("pdu.lengthInBytes + 1")]
    [InlineData("COUNT(events) + 6")]
    [InlineData("(COUNT(fifoValue) * 2) / 2")]
    [InlineData("ARRAY_SIZE_IN_BYTES(items)")]
    [InlineData("STATIC_CALL(\"rtuCrcCheck\", address, pdu)")]
    [InlineData("STATIC_CALL(\"asciiLrcCheck\", address, pdu)")]
    public void BuildsExpressionsUsedByModbus(string input)
    {
        Assert.NotNull(MspecExpressionParser.Parse(input));
    }

    [Fact]
    public void RejectsInvalidExpressionWithSourceLocation()
    {
        MspecParseException error =
            Assert.Throws<MspecParseException>(() => MspecExpressionParser.Parse("1 +"));

        Assert.Contains("'1 +'", error.Message);
        Assert.Contains("line 1:", error.Message);
    }

    [Theory]
    [InlineData(false)]
    [InlineData(true)]
    public void RejectsOutOfRangeNumbersAsParseErrors(bool floatingPoint)
    {
        string input = floatingPoint ? new string('9', 400) + ".0" : new string('9', 30);
        MspecParseException error =
            Assert.Throws<MspecParseException>(() => MspecExpressionParser.Parse(input));

        Assert.Contains(input, error.Message);
        Assert.IsType<OverflowException>(error.InnerException);
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
