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
using org.apache.plc4net.tools.codegen;
using Xunit;

namespace org.apache.plc4net.spi.test.codegen
{
    /// <summary>
    /// Verifies that the pure-.NET mspec parser can parse the Modbus protocol
    /// description, which is the simplest of the 24 protocol descriptions and
    /// therefore the first target for the generator pipeline.
    /// </summary>
    public class MspecParserTests
    {
        private static string GetModbusMspecPath()
        {
            // The mspec file lives in the repo root under protocols/modbus/...
            // FindRepoRoot walks up from the test output directory; null means
            // the tests are not running inside a checkout.
            var repoRoot = RepoPaths.FindRepoRoot();
            if (repoRoot == null) return null;
            return Path.Combine(repoRoot,
                "protocols", "modbus", "src", "main", "resources",
                "protocols", "modbus", "modbus.mspec");
        }

        [Fact]
        public void Modbus_mspec_parses_without_errors()
        {
            var path = GetModbusMspecPath();
            if (path == null)
            {
                // Skip gracefully when run outside the repo layout.
                return;
            }

            // Deliberately no File.Exists guard here: inside a checkout the
            // file must exist, and a missing one is a real error worth failing
            // on, not a reason to pass vacuously.
            var tree = MspecReader.ReadFile(path);

            Assert.NotNull(tree);
        }

        [Fact]
        public void A_minimal_mspec_parses()
        {
            var mspec = @"
[type Byte
    [simple uint 8 value]
]
";
            var tree = MspecReader.Read(mspec);
            Assert.NotNull(tree);
        }

        [Fact]
        public void Malformed_mspec_throws()
        {
            var ex = Assert.Throws<MspecParseException>(
                () => MspecReader.Read("[type Broken <<<"));

            Assert.Contains("Failed to parse", ex.Message);
        }
    }
}
