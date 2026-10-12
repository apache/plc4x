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
using System.IO;

namespace org.apache.plc4net.drivers.modbus.test
{
    /// <summary>
    /// Locates the shared Modbus test data in the repository checkout. The test
    /// host starts several directories below the checkout root, so the root is
    /// found by walking up to the directory that holds both <c>pom.xml</c> and
    /// <c>plc4net/</c>. A missing root or file fails the test: silently skipping
    /// would make a layout change look like a pass.
    /// </summary>
    internal static class RepoPaths
    {
        public static string FindRepoRoot()
        {
            var dir = new DirectoryInfo(Directory.GetCurrentDirectory());
            while (dir != null)
            {
                if (File.Exists(Path.Combine(dir.FullName, "pom.xml")) &&
                    Directory.Exists(Path.Combine(dir.FullName, "plc4net")))
                {
                    return dir.FullName;
                }
                dir = dir.Parent;
            }
            throw new InvalidOperationException(
                "The Modbus tests need the PLC4X checkout: no parent of " +
                $"'{Directory.GetCurrentDirectory()}' holds both pom.xml and plc4net/.");
        }

        /// <summary>
        /// Path of a ParserSerializerTestsuite.xml shared with plc4j and plc4go,
        /// e.g. <c>protocol</c> = "tcp", "rtu" or "ascii".
        /// </summary>
        public static string ParserSerializerTestsuite(string protocol)
        {
            var path = Path.Combine(FindRepoRoot(),
                "protocols", "modbus", "src", "test", "resources",
                "protocols", "modbus", protocol, "ParserSerializerTestsuite.xml");
            if (!File.Exists(path))
            {
                throw new FileNotFoundException(
                    $"The shared Modbus {protocol} test suite is missing.", path);
            }
            return path;
        }
    }
}