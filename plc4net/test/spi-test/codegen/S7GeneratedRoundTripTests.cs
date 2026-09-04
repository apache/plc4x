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
using System.Globalization;
using System.IO;
using System.Linq;
using System.Reflection;
using org.apache.plc4net.drivers.s7.readwrite.model;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.tools.codegen;
using Xunit;

namespace org.apache.plc4net.spi.test.codegen
{
    /// <summary>
    /// Runs the shared <c>ParserSerializerTestsuite.xml</c> S7 vectors - the
    /// same frames plc4j and plc4go validate - through the C# model
    /// <see cref="output.CSharpGenerator"/> produced from s7.mspec. A vector
    /// passes when its bytes parse and serialize back byte-identical.
    /// </summary>
    /// <remarks>
    /// S7 exercises the generator far harder than Modbus: nested discriminated
    /// types (TPKT -&gt; COTP -&gt; S7 message -&gt; parameter / payload),
    /// pre-typeSwitch implicit lengths a child has to see, a parameterised
    /// self-referential enum (<c>TransportSize</c>) and an <c>enum</c> field
    /// keyed on an attribute rather than the constant.
    /// </remarks>
    public class S7GeneratedRoundTripTests
    {
        private static readonly Assembly S7Model = typeof(TPKTPacket).Assembly;
        private const string ModelNamespace = "org.apache.plc4net.drivers.s7.readwrite.model";

        public static TheoryData<string> Vectors()
        {
            var data = new TheoryData<string>();
            var suite = LoadSuite();
            if (suite == null)
            {
                data.Add("(no checkout - skipped)");
                return data;
            }
            foreach (var tc in suite.TestCases)
            {
                data.Add(tc.Name);
            }
            return data;
        }

        [Theory]
        [MemberData(nameof(Vectors))]
        public void Vector_round_trips(string name)
        {
            var suite = LoadSuite();
            if (suite == null)
            {
                return; // outside a checkout
            }

            var tc = suite.TestCases.Single(t => t.Name == name);
            var expected = tc.GetRawBytes();

            var parsed = Parse(tc);
            Assert.NotNull(parsed);

            var writeBuffer = new WriteBuffer();
            parsed.GetType().GetMethod("Serialize")!.Invoke(parsed, new object[] { writeBuffer });
            var actual = writeBuffer.GetBytes();

            Assert.Equal(ToHex(expected), ToHex(actual));
        }

        [Fact]
        public void The_nested_discriminated_types_resolve_to_the_expected_classes()
        {
            var suite = LoadSuite();
            if (suite == null)
            {
                return;
            }

            var parsed = Parse(suite.TestCases.Single(t => t.Name == "COTP Connection Request"));

            Assert.Equal("TPKTPacket", parsed.GetType().Name);
            var cotp = parsed.GetType().GetProperty("Payload")!.GetValue(parsed);
            Assert.Equal("COTPPacketConnectionRequest", cotp!.GetType().Name);
        }

        // ── helpers ─────────────────────────────────────────────

        private static object Parse(ParserSerializerTestcase tc)
        {
            var rootType = S7Model.GetType($"{ModelNamespace}.{tc.RootType}")
                           ?? throw new InvalidOperationException($"generated type {tc.RootType} not found");
            var staticParse = rootType.GetMethod("StaticParse", BindingFlags.Public | BindingFlags.Static)!;

            var pars = staticParse.GetParameters();
            var argv = new object[pars.Length];
            argv[0] = new ReadBuffer(tc.GetRawBytes());
            for (var i = 1; i < pars.Length; i++)
            {
                var raw = tc.ParserArguments[pars[i].Name!];
                argv[i] = Coerce(raw, pars[i].ParameterType);
            }

            return staticParse.Invoke(null, argv);
        }

        private static object Coerce(string raw, Type target)
        {
            if (target.IsEnum)
            {
                return Enum.Parse(target, raw, ignoreCase: true);
            }
            if (target == typeof(bool))
            {
                return bool.Parse(raw);
            }
            return Convert.ChangeType(raw, target, CultureInfo.InvariantCulture);
        }

        private static ParserSerializerTestsuite LoadSuite()
        {
            var repoRoot = RepoPaths.FindRepoRoot();
            if (repoRoot == null)
            {
                return null;
            }
            var path = Path.Combine(repoRoot,
                "protocols", "s7", "src", "test", "resources",
                "protocols", "s7", "ParserSerializerTestsuite.xml");
            return ParserSerializerTestsuiteRunner.Load(path);
        }

        private static string ToHex(byte[] bytes) =>
            string.Concat(bytes.Select(b => b.ToString("x2", CultureInfo.InvariantCulture)));
    }
}
