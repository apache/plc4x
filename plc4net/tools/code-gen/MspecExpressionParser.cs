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
using System.Linq;
using Antlr4.Runtime;
using Antlr4.Runtime.Tree;
using org.apache.plc4net.tools.codegen.grammar;
using org.apache.plc4net.tools.codegen.model.terms;

namespace org.apache.plc4net.tools.codegen
{
    /// <summary>
    /// Parses an mspec expression - the text inside the single quotes of a
    /// <c>count</c>, <c>length</c>, <c>implicit</c>, <c>virtual</c>,
    /// <c>discriminator</c> or <c>const</c> field - into a
    /// <see cref="Term"/> tree.
    /// </summary>
    /// <remarks>
    /// Uses the ANTLR grammar <c>Expression.g4</c> (the same grammar plc4j and
    /// plc4go parse) via the checked-in C# lexer/parser. ANTLR only generated a
    /// listener, so the tree is walked here by pattern-matching the labelled
    /// alternative context classes (<c>#addExpression</c> -&gt;
    /// <c>AddExpressionContext</c>, etc.) - no visitor and therefore no need to
    /// re-run the ANTLR tool.
    /// </remarks>
    public static class MspecExpressionParser
    {
        /// <summary>
        /// Parses <paramref name="expression"/> (without the surrounding
        /// single quotes) and returns its tree.
        /// </summary>
        /// <exception cref="MspecParseException">the text is not a valid expression.</exception>
        public static Term Parse(string expression)
        {
            if (expression == null)
            {
                throw new ArgumentNullException(nameof(expression));
            }

            var lexer = new ExpressionLexer(new AntlrInputStream(expression));
            var parser = new ExpressionParser(new CommonTokenStream(lexer));

            // A lexer error (unrecognised character) leaves a gap the parser
            // then reports as a syntax error, so collecting on the parser alone
            // is enough and avoids a second listener type for the lexer's int
            // token symbols.
            var errors = new CollectingErrorListener();
            parser.RemoveErrorListeners();
            parser.AddErrorListener(errors);

            var tree = parser.expressionString();

            if (errors.HasErrors)
            {
                throw new MspecParseException(
                    $"Failed to parse expression '{expression}': {string.Join("; ", errors.Errors)}");
            }

            return Visit(tree.expression());
        }

        // ── expression alternatives ──────────────────────────────

        private static Term Visit(ExpressionParser.ExpressionContext ctx)
        {
            switch (ctx)
            {
                case ExpressionParser.UnaryMinusExpressionContext c:
                    return new UnaryExpression("-", Visit(c.expression()));
                case ExpressionParser.NotExpressionContext c:
                    return new UnaryExpression("!", Visit(c.expression()));

                case ExpressionParser.PowerExpressionContext c:
                    return Bin(c.expression(0), "^", c.expression(1));
                case ExpressionParser.MultExpressionContext c:
                    return Bin(c.expression(0), c.op.Text, c.expression(1));
                case ExpressionParser.AddExpressionContext c:
                    return Bin(c.expression(0), c.op.Text, c.expression(1));
                case ExpressionParser.BitShiftExpressionContext c:
                    return Bin(c.expression(0), c.op.Text, c.expression(1));
                case ExpressionParser.CompExpressionContext c:
                    return Bin(c.expression(0), c.op.Text, c.expression(1));
                case ExpressionParser.EqExpressionContext c:
                    return Bin(c.expression(0), c.op.Text, c.expression(1));
                case ExpressionParser.AndExpressionContext c:
                    return Bin(c.expression(0), "&&", c.expression(1));
                case ExpressionParser.BitAndExpressionContext c:
                    return Bin(c.expression(0), "&", c.expression(1));
                case ExpressionParser.OrExpressionContext c:
                    return Bin(c.expression(0), "||", c.expression(1));
                case ExpressionParser.BitOrExpressionContext c:
                    return Bin(c.expression(0), "|", c.expression(1));

                case ExpressionParser.IfExpressionContext c:
                    return new TernaryExpression(
                        Visit(c.expression(0)), Visit(c.expression(1)), Visit(c.expression(2)));

                case ExpressionParser.NumberExpressionContext c:
                    return Number(c.Number().GetText());
                case ExpressionParser.HexExpressionContext c:
                    var hex = c.HexExpression().GetText();
                    return new IntegerLiteral(
                        long.Parse(hex.Substring(2), NumberStyles.HexNumber, CultureInfo.InvariantCulture),
                        hex);
                case ExpressionParser.BoolExpressionContext c:
                    return new BooleanLiteral(c.Bool().GetText() == "true");
                case ExpressionParser.NullExpressionContext _:
                    return new NullLiteral();

                case ExpressionParser.StringExpressionContext c:
                    return Unquote(c.String().GetText());

                case ExpressionParser.IdentifierExpressionContext c:
                    return VisitSegment(c.identifierSegment());

                case ExpressionParser.ExpressionExpressionContext c:
                    // '(' expression ')' indexes? - the parens only grouped;
                    // a trailing index would apply to the grouped value, which
                    // no mspec in tree does, so it is not modelled.
                    return Visit(c.expression());

                default:
                    throw new MspecParseException(
                        $"Unhandled expression node {ctx.GetType().Name}: '{ctx.GetText()}'");
            }
        }

        private static Term Bin(
            ExpressionParser.ExpressionContext l, string op, ExpressionParser.ExpressionContext r)
            => new BinaryExpression(Visit(l), op, Visit(r));

        // ── identifierSegment: name (args)? [index]? (. child)? ───

        private static VariableLiteral VisitSegment(ExpressionParser.IdentifierSegmentContext ctx)
        {
            var name = ctx.name.Text;

            IReadOnlyList<Term> args = null;
            if (ctx.args != null)
            {
                args = ctx.args.arguments().expression().Select(Visit).ToList();
            }

            IReadOnlyList<Term> index = Array.Empty<Term>();
            if (ctx.index != null)
            {
                index = ctx.index.indexes().expression().Select(Visit).ToList();
            }

            VariableLiteral child = null;
            if (ctx.rest != null)
            {
                child = VisitSegment(ctx.rest.identifierSegment());
            }

            return new VariableLiteral(name, args, index, child);
        }

        // ── literals ─────────────────────────────────────────────

        private static Term Number(string text)
        {
            if (text.Contains('.'))
            {
                return new FloatLiteral(double.Parse(text, CultureInfo.InvariantCulture));
            }
            return new IntegerLiteral(long.Parse(text, CultureInfo.InvariantCulture), text);
        }

        private static StringLiteral Unquote(string text)
        {
            // The lexer keeps the quotes and any backslash escapes.
            if (text.Length >= 2 &&
                (text[0] == '"' || text[0] == '\'') && text[^1] == text[0])
            {
                text = text.Substring(1, text.Length - 2);
            }
            return new StringLiteral(text.Replace("\\\"", "\"").Replace("\\'", "'"));
        }
    }
}
