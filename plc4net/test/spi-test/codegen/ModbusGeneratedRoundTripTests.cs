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
using org.apache.plc4net.drivers.modbus.readwrite.model;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.tools.codegen;
using Xunit;

namespace org.apache.plc4net.spi.test.codegen
{
    /// <summary>
    /// Runs the shared <c>ParserSerializerTestsuite.xml</c> Modbus TCP vectors -
    /// the same data plc4j and plc4go validate against - through the C# model
    /// that <see cref="output.CSharpGenerator"/> produced from modbus.mspec.
    /// A vector passes when its bytes parse and serialize back byte-identical.
    /// </summary>
    /// <remarks>
    /// This is the end-to-end proof for the pure-.NET generator: mspec -> IR ->
    /// C# -> a driver that actually round-trips real frames.
    /// </remarks>
    public class ModbusGeneratedRoundTripTests
    {
        private static readonly Assembly ModbusModel = typeof(ModbusADU).Assembly;
        private const string ModelNamespace = "org.apache.plc4net.drivers.modbus.readwrite.model";

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
        public void The_discriminated_type_resolves_to_the_expected_concrete_class()
        {
            var suite = LoadSuite();
            if (suite == null)
            {
                return;
            }

            var request = suite.TestCases.Single(t => t.Name == "Read Input Registers Request");
            var parsed = Parse(request);

            // ModbusADU -> ModbusTcpADU, pdu -> ModbusPDUReadInputRegistersRequest
            Assert.Equal("ModbusTcpADU", parsed.GetType().Name);
            var pdu = parsed.GetType().GetProperty("Pdu")!.GetValue(parsed);
            Assert.Equal("ModbusPDUReadInputRegistersRequest", pdu!.GetType().Name);
            Assert.Equal((ushort)2258, pdu.GetType().GetProperty("StartingAddress")!.GetValue(pdu));
        }

        // ── helpers ─────────────────────────────────────────────

        private static object Parse(ParserSerializerTestcase tc)
        {
            var rootType = ModbusModel.GetType($"{ModelNamespace}.{tc.RootType}")
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
                "protocols", "modbus", "src", "test", "resources",
                "protocols", "modbus", "tcp", "ParserSerializerTestsuite.xml");
            return ParserSerializerTestsuiteRunner.Load(path);
        }

        private static string ToHex(byte[] bytes) =>
            string.Concat(bytes.Select(b => b.ToString("x2", CultureInfo.InvariantCulture)));
    }
}
