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
using org.apache.plc4net.tools.codegen.output;
using Xunit;

namespace org.apache.plc4net.test.codegen;

public class CSharpExpressionRendererTests
{
    [Theory]
    [InlineData("0xFFFFFFFFFFFFFFFF", "0xFFFFFFFFFFFFFFFF")]
    [InlineData("pdu.lengthInBytes + 1", "(pdu.GetLengthInBytes() + 1)")]
    [InlineData("COUNT(value)", "value.Count")]
    [InlineData("2 ^ exponent", "System.Math.Pow(2, exponent)")]
    [InlineData("ARRAY_SIZE_IN_BYTES(items)", "StaticHelper.ArraySizeInBytes(items)")]
    [InlineData("STATIC_CALL(\"rtuCrcCheck\", address, pdu)", "StaticHelper.RtuCrcCheck(address, pdu)")]
    [InlineData("CAST(parameter, \"S7Parameter\").items[0]", "((S7Parameter) parameter).Items[0]")]
    [InlineData("CAST(parameter, S7Parameter).items[0]", "((S7Parameter) parameter).Items[0]")]
    [InlineData("(-value)[0]", "(-value)[0]")]
    [InlineData("(!flag)[0]", "(!flag)[0]")]
    public void RendersExpressionAsCSharp(string expression, string expected)
    {
        Assert.Equal(expected, Render(expression));
    }

    [Fact]
    public void UsesScopeForReferencesCountsAndMembers()
    {
        var renderer = new CSharpExpressionRenderer(new TestScope());

        Assert.Equal("(_value.Start + _value.Size)",
            renderer.Render(MspecExpressionParser.Parse("start + size")));
        Assert.Equal("_value.Bytes.Length",
            renderer.Render(MspecExpressionParser.Parse("COUNT(bytes)")));
        Assert.Equal("_value.Value.GetCode()",
            renderer.Render(MspecExpressionParser.Parse("value.code")));
    }

    [Fact]
    public void SupportsQualifiedUnquotedCastTypesWithoutResolvingThemAsValues()
    {
        var renderer = new CSharpExpressionRenderer(new TestScope());

        Assert.Equal("((Protocol.Message) _value.Value)",
            renderer.Render(MspecExpressionParser.Parse("CAST(value, Protocol.Message)")));
    }

    [Fact]
    public void RendersTheMspecEmptyStringCompatibilitySentinel()
    {
        Assert.Equal("\"\"", Render(MspecReader.EmptyStringSentinel));
    }

    private static string Render(string expression) =>
        new CSharpExpressionRenderer().Render(MspecExpressionParser.Parse(expression));

    private sealed class TestScope : IExpressionScope
    {
        public string ResolveReference(string name) => $"_value.{char.ToUpperInvariant(name[0])}{name[1..]}";

        public string? ResolveCount(string name) =>
            name == "bytes" ? "_value.Bytes.Length" : null;

        public string? ResolveMember(string name) =>
            name == "code" ? "GetCode()" : null;
    }
}
