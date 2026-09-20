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

using org.apache.plc4net.tools.codegen;
using org.apache.plc4net.tools.codegen.model.terms;
using org.apache.plc4net.tools.codegen.output;
using Xunit;

namespace org.apache.plc4net.spi.test.codegen
{
    /// <summary>
    /// The mspec expression engine: <see cref="MspecExpressionParser"/> (text
    /// -&gt; <see cref="Term"/> tree, via the <c>Expression.g4</c> grammar) and
    /// <see cref="CSharpExpressionRenderer"/> (tree -&gt; C#). Every field that
    /// carries a formula - array <c>count</c> / <c>length</c>, <c>implicit</c>,
    /// <c>virtual</c>, <c>const</c>, discriminator values - depends on both.
    /// </summary>
    public class ExpressionTests
    {
        // ── literals ─────────────────────────────────────────────

        [Theory]
        [InlineData("42", "42")]
        [InlineData("0", "0")]
        [InlineData("0x0E", "0x0E")]
        [InlineData("0xFF", "0xFF")]
        [InlineData("true", "true")]
        [InlineData("false", "false")]
        [InlineData("null", "null")]
        [InlineData("1.5", "1.5")]
        public void Literals_round_trip(string mspec, string expected)
        {
            Assert.Equal(expected, MspecExpressionParser.Parse(mspec).ToString());
        }

        [Fact]
        public void String_literal_drops_the_quotes()
        {
            var term = Assert.IsType<StringLiteral>(MspecExpressionParser.Parse("\"UTF8\""));
            Assert.Equal("UTF8", term.Value);
        }

        // ── operators and precedence ─────────────────────────────

        [Theory]
        [InlineData("1 + 2", "(1 + 2)")]
        [InlineData("1 + 2 * 3", "(1 + (2 * 3))")]
        [InlineData("(1 + 2) * 3", "((1 + 2) * 3)")]
        [InlineData("byteCount - 6", "(byteCount - 6)")]
        [InlineData("numberOfValues * 8", "(numberOfValues * 8)")]
        [InlineData("(COUNT(fifoValue) * 2) + 2", "((COUNT(fifoValue) * 2) + 2)")]
        [InlineData("a >= b && c != d", "((a >= b) && (c != d))")]
        [InlineData("!flag", "!flag")]
        [InlineData("-1", "-1")]
        [InlineData("a ? b : c", "(a ? b : c)")]
        public void Operators_parse_with_the_right_shape(string mspec, string expectedTree)
        {
            Assert.Equal(expectedTree, MspecExpressionParser.Parse(mspec).ToString());
        }

        // ── identifier segments: name (args)? [index]? (. child)? ─

        [Fact]
        public void Dotted_path_becomes_a_child_chain()
        {
            var v = Assert.IsType<VariableLiteral>(MspecExpressionParser.Parse("pdu.lengthInBytes"));
            Assert.Equal("pdu", v.Name);
            Assert.Null(v.Args);
            Assert.Equal("lengthInBytes", v.Child.Name);
        }

        [Fact]
        public void Call_captures_arguments()
        {
            var v = Assert.IsType<VariableLiteral>(MspecExpressionParser.Parse("COUNT(value)"));
            Assert.Equal("COUNT", v.Name);
            Assert.NotNull(v.Args);
            Assert.Single(v.Args);
            Assert.Equal("value", v.Args[0].ToString());
        }

        [Fact]
        public void Static_call_keeps_the_string_target_and_the_rest()
        {
            var v = Assert.IsType<VariableLiteral>(
                MspecExpressionParser.Parse("STATIC_CALL(\"rtuCrcCheck\", address, pdu)"));
            Assert.Equal("STATIC_CALL", v.Name);
            Assert.Equal(3, v.Args.Count);
            Assert.IsType<StringLiteral>(v.Args[0]);
        }

        [Fact]
        public void Indexing_is_captured()
        {
            var v = Assert.IsType<VariableLiteral>(MspecExpressionParser.Parse("items[0].fileNumber"));
            Assert.Single(v.Index);
            Assert.Equal("0", v.Index[0].ToString());
            Assert.Equal("fileNumber", v.Child.Name);
        }

        // ── C# rendering ─────────────────────────────────────────

        [Theory]
        [InlineData("pdu.lengthInBytes + 1", "(pdu.GetLengthInBytes() + 1)")]
        [InlineData("COUNT(value)", "value.Count")]
        [InlineData("COUNT(value) + 6", "(value.Count + 6)")]
        [InlineData("byteCount - 6", "(byteCount - 6)")]
        [InlineData("(COUNT(fifoValue) * 2) / 2", "((fifoValue.Count * 2) / 2)")]
        [InlineData("ARRAY_SIZE_IN_BYTES(items)", "StaticHelper.ArraySizeInBytes(items)")]
        [InlineData("STATIC_CALL(\"rtuCrcCheck\", address, pdu)", "StaticHelper.RtuCrcCheck(address, pdu)")]
        [InlineData("numberOfValues * 8", "(numberOfValues * 8)")]
        [InlineData("errorFlag ? 1 : 0", "(errorFlag ? 1 : 0)")]
        public void Renders_to_C_sharp(string mspec, string expectedCSharp)
        {
            var term = MspecExpressionParser.Parse(mspec);
            Assert.Equal(expectedCSharp, new CSharpExpressionRenderer().Render(term));
        }

        [Fact]
        public void The_scope_rewrites_bare_references()
        {
            var term = MspecExpressionParser.Parse("startingAddress + quantity");
            var renderer = new CSharpExpressionRenderer(new PrefixScope("_value."));
            Assert.Equal("(_value.startingAddress + _value.quantity)", renderer.Render(term));
        }

        [Fact]
        public void Power_operator_maps_to_Math_Pow_not_xor()
        {
            var term = MspecExpressionParser.Parse("2 ^ exponent");
            Assert.Equal("System.Math.Pow(2, exponent)", new CSharpExpressionRenderer().Render(term));
        }

        [Fact]
        public void Every_modbus_field_expression_parses()
        {
            // The formulas that appear verbatim in modbus.mspec.
            string[] expressions =
            {
                "pdu.lengthInBytes + 1",
                "COUNT(value)",
                "COUNT(events) + 6",
                "byteCount - 6",
                "(COUNT(fifoValue) * 2) + 2",
                "(COUNT(fifoValue) * 2) / 2",
                "ARRAY_SIZE_IN_BYTES(items)",
                "COUNT(data) + 1",
                "dataLength - 1",
                "COUNT(recordData) / 2",
                "recordLength * 2",
                "COUNT(objects)",
                "STATIC_CALL(\"rtuCrcCheck\", address, pdu)",
                "STATIC_CALL(\"asciiLrcCheck\", address, pdu)",
            };

            var renderer = new CSharpExpressionRenderer();
            foreach (var expr in expressions)
            {
                var term = MspecExpressionParser.Parse(expr);
                Assert.False(string.IsNullOrEmpty(renderer.Render(term)), expr);
            }
        }

        [Fact]
        public void A_broken_expression_reports_where()
        {
            var ex = Assert.Throws<MspecParseException>(() => MspecExpressionParser.Parse("1 +"));
            Assert.Contains("'1 +'", ex.Message);
        }

        private sealed class PrefixScope : IExpressionScope
        {
            private readonly string _prefix;
            public PrefixScope(string prefix) => _prefix = prefix;
            public string ResolveReference(string name) => _prefix + name;
        }
    }
}
