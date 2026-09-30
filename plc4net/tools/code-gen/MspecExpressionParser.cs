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
using System.Globalization;
using System.IO;
using System.Linq;
using Antlr4.Runtime;
using org.apache.plc4net.tools.codegen.grammar;
using org.apache.plc4net.tools.codegen.model.terms;

namespace org.apache.plc4net.tools.codegen;

/// <summary>Converts an expression parsed by the shared grammar into its semantic tree.</summary>
public static class MspecExpressionParser
{
    public static Term Parse(string expression)
    {
        ArgumentNullException.ThrowIfNull(expression);
        var lexer = new ExpressionLexer(CharStreams.fromString(expression));
        var lexerErrors = new ErrorCollector<int>();
        lexer.RemoveErrorListeners();
        lexer.AddErrorListener(lexerErrors);
        var parser = new ExpressionParser(new CommonTokenStream(lexer));
        var parserErrors = new ErrorCollector<IToken>();
        parser.RemoveErrorListeners();
        parser.AddErrorListener(parserErrors);

        ExpressionParser.ExpressionStringContext parsed = parser.expressionString();
        string[] errors = lexerErrors.Errors.Concat(parserErrors.Errors).ToArray();
        if (errors.Length != 0)
        {
            throw new MspecParseException(
                $"Failed to parse expression '{expression}': {string.Join("; ", errors)}");
        }
        try
        {
            return Convert(parsed.expression());
        }
        catch (Exception exception) when (exception is FormatException or OverflowException)
        {
            throw new MspecParseException(
                $"Failed to parse expression '{expression}': numeric literal is outside the supported range.",
                exception);
        }
    }

    private static Term Convert(ExpressionParser.ExpressionContext context) => context switch
    {
        ExpressionParser.NumberExpressionContext value => Number(value.Number().GetText()),
        ExpressionParser.HexExpressionContext value =>
            new HexadecimalLiteral(value.HexExpression().GetText()),
        ExpressionParser.BoolExpressionContext value =>
            new BooleanLiteral(bool.Parse(value.Bool().GetText())),
        ExpressionParser.NullExpressionContext => new NullLiteral(),
        ExpressionParser.StringExpressionContext value => WithIndexes(
            new StringLiteral(Unquote(value.String().GetText())), value.indexes()),
        ExpressionParser.IdentifierExpressionContext value => Identifier(value.identifierSegment()),
        ExpressionParser.UnaryMinusExpressionContext value =>
            new UnaryExpression("-", Convert(value.expression())),
        ExpressionParser.NotExpressionContext value =>
            new UnaryExpression("!", Convert(value.expression())),
        ExpressionParser.ExpressionExpressionContext value =>
            WithIndexes(Convert(value.expression()), value.indexes()),
        ExpressionParser.PowerExpressionContext value => Binary(value, "^"),
        ExpressionParser.MultExpressionContext value => Binary(value, value.op.Text),
        ExpressionParser.AddExpressionContext value => Binary(value, value.op.Text),
        ExpressionParser.BitShiftExpressionContext value => Binary(value, value.op.Text),
        ExpressionParser.CompExpressionContext value => Binary(value, value.op.Text),
        ExpressionParser.EqExpressionContext value => Binary(value, value.op.Text),
        ExpressionParser.AndExpressionContext value => Binary(value, "&&"),
        ExpressionParser.BitAndExpressionContext value => Binary(value, "&"),
        ExpressionParser.OrExpressionContext value => Binary(value, "||"),
        ExpressionParser.BitOrExpressionContext value => Binary(value, "|"),
        ExpressionParser.IfExpressionContext value => new TernaryExpression(
            Convert(value.expression(0)), Convert(value.expression(1)), Convert(value.expression(2))),
        _ => throw new MspecParseException($"Unsupported expression: {context.GetText()}")
    };

    private static Term Number(string text)
    {
        if (!text.Contains('.'))
        {
            return new IntegerLiteral(long.Parse(text, CultureInfo.InvariantCulture), text);
        }
        double value = double.Parse(text, CultureInfo.InvariantCulture);
        return double.IsFinite(value)
            ? new FloatLiteral(value)
            : throw new OverflowException();
    }

    private static BinaryExpression Binary(ParserRuleContext context, string op)
    {
        ExpressionParser.ExpressionContext[] operands =
            context.GetRuleContexts<ExpressionParser.ExpressionContext>();
        return new BinaryExpression(Convert(operands[0]), op, Convert(operands[1]));
    }

    private static Term WithIndexes(Term target, ExpressionParser.IndexesContext? indexes) =>
        indexes == null
            ? target
            : new IndexExpression(target, indexes.expression().Select(Convert).ToArray());

    private static VariableLiteral Identifier(ExpressionParser.IdentifierSegmentContext context)
    {
        ExpressionParser.IdentifierSegmentArgumentsContext? arguments =
            context.identifierSegmentArguments();
        ExpressionParser.IdentifierSegmentIndexesContext? indexes =
            context.identifierSegmentIndexes();
        ExpressionParser.IdentifierSegmentRestContext? rest = context.identifierSegmentRest();
        return new VariableLiteral(
            context.Identifier().GetText(),
            arguments?.arguments().expression().Select(Convert).ToArray(),
            indexes?.indexes().expression().Select(Convert).ToArray(),
            rest == null ? null : Identifier(rest.identifierSegment()));
    }

    private static string Unquote(string text)
    {
        char delimiter = text[0];
        var value = new System.Text.StringBuilder(text.Length - 2);
        for (int i = 1; i < text.Length - 1; i++)
        {
            if (text[i] == '\\' && i + 1 < text.Length - 1 &&
                (text[i + 1] == delimiter || text[i + 1] == '\\'))
            {
                i++;
            }
            value.Append(text[i]);
        }
        return value.ToString();
    }

    private sealed class ErrorCollector<TSymbol> : IAntlrErrorListener<TSymbol>
    {
        public List<string> Errors { get; } = [];

        public void SyntaxError(TextWriter output, IRecognizer recognizer, TSymbol offendingSymbol,
            int line, int charPositionInLine, string msg, RecognitionException e) =>
            Errors.Add($"line {line}:{charPositionInLine} {msg}");
    }
}

public sealed class MspecParseException : Exception
{
    public MspecParseException(string message) : base(message)
    {
    }

    public MspecParseException(string message, Exception innerException) : base(message, innerException)
    {
    }
}
