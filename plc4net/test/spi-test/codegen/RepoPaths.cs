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

namespace org.apache.plc4net.spi.test.codegen
{
    /// <summary>
    /// Resolves files in the repository checkout that integration-style tests
    /// need — the real modbus.mspec, the KNX test-suite XML. The test host's
    /// working directory is the test output directory, several levels below
    /// the checkout root, so a hard-coded walk-up depth silently misses the
    /// root when the layout changes. Walking up until a directory containing
    /// pom.xml and plc4net/ marks the root is robust to that depth.
    /// </summary>
    internal static class RepoPaths
    {
        /// <summary>
        /// Root of the PLC4X checkout, or null when the tests run somewhere
        /// that is not inside a checkout (e.g. a decompiled source tree).
        /// </summary>
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
            return null;
        }
    }
}
