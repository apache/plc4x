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

using System.IO;
using System.Linq;
using org.apache.plc4net.tools.codegen;
using Xunit;

namespace org.apache.plc4net.test.codegen
{
    /// <summary>
    /// Verifies that the .NET ParserSerializerTestsuite runner can load and
    /// decode the same XML test data that plc4j and plc4go consume.
    /// </summary>
    public class TestsuiteRunnerTests
    {
        [Fact]
        public void Can_decode_hex_bytes()
        {
            var testCase = new ParserSerializerTestcase
            {
                Name = "test",
                RawHex = "010203ff"
            };

            var bytes = testCase.GetRawBytes();

            Assert.Equal(new byte[] { 0x01, 0x02, 0x03, 0xFF }, bytes);
        }

        [Fact]
        public void Whitespace_in_hex_is_ignored()
        {
            var testCase = new ParserSerializerTestcase
            {
                Name = "test",
                RawHex = "01 02\n03\r\n04"
            };

            var bytes = testCase.GetRawBytes();

            Assert.Equal(new byte[] { 0x01, 0x02, 0x03, 0x04 }, bytes);
        }

        [Fact]
        public void Empty_raw_produces_empty_bytes()
        {
            var empty = new ParserSerializerTestcase { Name = "test", RawHex = "" };
            var nullHex = new ParserSerializerTestcase { Name = "test", RawHex = null };

            Assert.Empty(empty.GetRawBytes());
            Assert.Empty(nullHex.GetRawBytes());
        }

        [Fact]
        public void Can_load_modbus_testsuite()
        {
            // Same test data consumed by plc4j and plc4go. Modbus, not KNX: the
            // generator tool has no driver-specific test resources of its own,
            // so this reads the one shared testsuite XML that already exists
            // on develop independently of any plc4net driver PR.
            // Fail loud rather than passing vacuously outside the repo layout.
            var repoRoot = RepoPaths.FindRepoRoot();
            Assert.NotNull(repoRoot);

            var path = Path.Combine(repoRoot!,
                "protocols", "modbus", "src", "test", "resources",
                "protocols", "modbus", "tcp", "ParserSerializerTestsuite.xml");

            var suite = ParserSerializerTestsuiteRunner.Load(path);

            Assert.Equal("Modbus", suite.Name);
            Assert.Equal("modbus", suite.ProtocolName);
            Assert.Equal("BIG_ENDIAN", suite.ByteOrder);
            Assert.True(suite.TestCases.Any());
        }
    }
}
