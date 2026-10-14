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
using System.Linq;
using System.Reflection;
using org.apache.plc4net.drivers.s7.readwrite.model;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.tools.codegen;
using Xunit;

namespace org.apache.plc4net.drivers.s7.test
{
    /// <summary>
    /// Runs the ParserSerializerTestsuite vectors that plc4j and plc4go validate against
    /// through the generated S7 model. A vector passes when its bytes parse and serialize
    /// back byte-identical.
    /// </summary>
    public class GeneratedModelTests
    {
        private static readonly Assembly S7Model = typeof(TPKTPacket).Assembly;
        private static readonly Lazy<ParserSerializerTestsuite> Suite =
            new Lazy<ParserSerializerTestsuite>(() => ParserSerializerTestsuiteRunner.Load(RepoPaths.ParserSerializerTestsuite()));
        private const string ModelNamespace = "org.apache.plc4net.drivers.s7.readwrite.model";

        public static TheoryData<string> Vectors()
        {
            var data = new TheoryData<string>();
            foreach (var testCase in LoadSuite().TestCases)
            {
                data.Add(testCase.Name);
            }
            return data;
        }

        [Theory]
        [MemberData(nameof(Vectors))]
        public void Vector_round_trips_byte_identical(string name)
        {
            var testCase = LoadSuite().TestCases.Single(t => t.Name == name);
            var expected = testCase.GetRawBytes();

            var parsed = Parse(testCase);

            var writeBuffer = new WriteBuffer();
            parsed.GetType().GetMethod("Serialize")!.Invoke(parsed, new object[] { writeBuffer });

            Assert.Equal(ToHex(expected), ToHex(writeBuffer.GetBytes()));
        }

        [Theory]
        [MemberData(nameof(Vectors))]
        public void Vector_reports_its_own_length(string name)
        {
            var testCase = LoadSuite().TestCases.Single(t => t.Name == name);
            var expected = testCase.GetRawBytes();

            var parsed = Parse(testCase);

            // The lengths a message computes feed the implicit length fields of its parents.
            var length = (int)parsed.GetType().GetMethod("GetLengthInBytes")!.Invoke(parsed, null)!;
            Assert.Equal(expected.Length, length);
        }

        [Theory]
        [InlineData(0)]
        [InlineData(1)]
        [InlineData(6)]
        public void The_length_of_a_type_with_a_manual_field_is_the_number_of_bytes_it_serializes(int dataLength)
        {
            // AssociatedValueType holds the manual field valueLength, which the mspec sizes as '2':
            // the two bytes it takes on the wire.
            var value = new AssociatedValueType(
                DataTransportErrorCode.OK,
                DataTransportSize.BYTE_WORD_DWORD,
                (ushort)dataLength,
                Enumerable.Repeat((byte)0x5A, dataLength).ToList());

            var writeBuffer = new WriteBuffer();
            value.Serialize(writeBuffer);

            Assert.Equal(writeBuffer.GetBytes().Length, value.GetLengthInBytes());
            Assert.Equal(4 + dataLength, value.GetLengthInBytes());
        }

        [Fact]
        public void The_shared_suite_is_actually_loaded()
        {
            // Guards against a path or layout change turning the theory above into zero
            // cases, which xunit reports as a pass.
            Assert.True(LoadSuite().TestCases.Count >= 10, $"only {LoadSuite().TestCases.Count} vectors");
        }

        [Fact]
        public void The_nested_discriminated_types_resolve_to_the_expected_classes()
        {
            var parsed = Parse(LoadSuite().TestCases.Single(t => t.Name == "COTP Connection Request"));

            Assert.Equal("TPKTPacket", parsed.GetType().Name);
            var cotp = parsed.GetType().GetProperty("Payload")!.GetValue(parsed);
            Assert.Equal("COTPPacketConnectionRequest", cotp!.GetType().Name);
        }

        private static object Parse(ParserSerializerTestcase testCase)
        {
            var rootType = S7Model.GetType($"{ModelNamespace}.{testCase.RootType}")
                           ?? throw new InvalidOperationException($"generated type {testCase.RootType} not found");
            var staticParse = rootType.GetMethod("StaticParse", BindingFlags.Public | BindingFlags.Static)!;

            var parameters = staticParse.GetParameters();
            var arguments = new object[parameters.Length];
            arguments[0] = new ReadBuffer(testCase.GetRawBytes());
            for (var i = 1; i < parameters.Length; i++)
            {
                arguments[i] = Coerce(testCase.ParserArguments[parameters[i].Name!], parameters[i].ParameterType);
            }
            return staticParse.Invoke(null, arguments)!;
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
            => Suite.Value;

        private static string ToHex(byte[] bytes)
            => string.Concat(bytes.Select(b => b.ToString("x2", CultureInfo.InvariantCulture)));
    }
}
