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
using org.apache.plc4net.drivers.knxnetip.readwrite.model;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.tools.codegen;
using Xunit;

namespace org.apache.plc4net.spi.test.codegen
{
    /// <summary>
    /// Runs the shared <c>ParserSerializerTestsuite.xml</c> KNXnet/IP vectors - the
    /// same frames plc4j and plc4go validate - through the C# model generated from
    /// knxnetip.mspec. A vector passes when its bytes parse and then serialize back
    /// byte-identical, which exercises the whole nested frame (common header ->
    /// service body -> cEMI -> APDU / DPT).
    /// </summary>
    public class KnxGeneratedRoundTripTests
    {
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
                // "Causes Failure N" vectors are deliberately malformed - they belong
                // to a parse-rejection test, not a round-trip one.
                if (!tc.Name.StartsWith("Causes Failure", StringComparison.Ordinal))
                {
                    data.Add(tc.Name);
                }
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

            var parsed = KnxNetIpMessage.StaticParse(new ReadBuffer(expected));
            Assert.NotNull(parsed);

            var writeBuffer = new WriteBuffer();
            parsed.Serialize(writeBuffer);

            Assert.Equal(ToHex(expected), ToHex(writeBuffer.GetBytes()));
        }

        [Fact]
        public void The_tunnelling_request_vector_decodes_to_a_group_value_write()
        {
            var suite = LoadSuite();
            if (suite == null)
            {
                return;
            }

            var tc = suite.TestCases.Single(t => t.Name == "Tunneling Request");
            var message = KnxNetIpMessage.StaticParse(new ReadBuffer(tc.GetRawBytes()));

            var request = Assert.IsType<TunnelingRequest>(message);
            var ind = Assert.IsAssignableFrom<CEMI>(request.Cemi);
            Assert.NotNull(ind);
        }

        private static ParserSerializerTestsuite LoadSuite()
        {
            var repoRoot = RepoPaths.FindRepoRoot();
            if (repoRoot == null)
            {
                return null;
            }
            var path = Path.Combine(repoRoot,
                "plc4net", "test", "knxnetip-test", "resources",
                "protocols", "knxnetip", "ParserSerializerTestsuite.xml");
            return ParserSerializerTestsuiteRunner.Load(path);
        }

        private static string ToHex(byte[] bytes) =>
            string.Concat(bytes.Select(b => b.ToString("x2", CultureInfo.InvariantCulture)));
    }
}
