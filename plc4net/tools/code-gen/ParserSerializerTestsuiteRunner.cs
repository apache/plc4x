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

using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Xml.Linq;

namespace org.apache.plc4net.tools.codegen
{
    /// <summary>
    /// .NET executor for PLC4X ParserSerializerTestsuite XML files.
    /// Mirrors the Java <c>ParserSerializerTestsuiteRunner</c> and Go
    /// <c>ParserSerializerTestRunner</c>: each reads the same XML test
    /// data and runs it against their own language's generated code.
    /// </summary>
    /// <remarks>
    /// The XML schema lives at
    /// https://plc4x.apache.org/schemas/parser-serializer-testsuite.xsd
    /// and is shared across all languages. This runner loads the XML,
    /// decodes the hex test vectors, and exposes them as a typed model
    /// so that a protocol-specific test harness can feed them into the
    /// generated parse/serialize code.
    /// </remarks>
    public class ParserSerializerTestsuiteRunner
    {
        /// <summary>
        /// Loads a testsuite XML file.
        /// </summary>
        public static ParserSerializerTestsuite Load(string filePath)
        {
            var doc = XDocument.Load(filePath);
            var ns = doc.Root?.GetDefaultNamespace() ?? XNamespace.None;
            var testNamespace = ns.NamespaceName.Contains("parser-serializer")
                ? ns
                : XNamespace.Get("https://plc4x.apache.org/schemas/parser-serializer-testsuite.xsd");

            var root = doc.Root;
            if (root == null)
            {
                throw new InvalidDataException("Empty testsuite file.");
            }

            var suite = new ParserSerializerTestsuite
            {
                Name = root.Element("name")?.Value,
                ProtocolName = root.Element("protocol-name")?.Value,
                OutputFlavor = root.Element("output-flavor")?.Value,
                ByteOrder = root.Attribute("byteOrder")?.Value ?? "BIG_ENDIAN"
            };

            foreach (var caseEl in root.Elements("testcase"))
            {
                var testCase = new ParserSerializerTestcase
                {
                    Name = caseEl.Element("name")?.Value ?? "unnamed",
                    RawHex = caseEl.Element("raw")?.Value,
                    RootType = caseEl.Element("root-type")?.Value,
                    // The expected parsed tree. Leaf elements carry dataType /
                    // bitLength attributes and the value as text; container
                    // elements name the (possibly discriminated) type.
                    ExpectedXml = caseEl.Element("xml")?.Elements().FirstOrDefault(),
                };

                // Collect expected values from the <parser-arguments> element.
                var args = caseEl.Element("parser-arguments");
                if (args != null)
                {
                    foreach (var arg in args.Elements())
                    {
                        testCase.ParserArguments[arg.Name.LocalName] = arg.Value;
                    }
                }

                suite.TestCases.Add(testCase);
            }

            return suite;
        }
    }

    /// <summary>
    /// The typed content of a ParserSerializerTestsuite.xml.
    /// </summary>
    public class ParserSerializerTestsuite
    {
        public string Name { get; set; }
        public string ProtocolName { get; set; }
        public string OutputFlavor { get; set; }
        public string ByteOrder { get; set; }

        public List<ParserSerializerTestcase> TestCases { get; set; }
            = new List<ParserSerializerTestcase>();
    }

    /// <summary>
    /// A single test case: raw hex bytes and the expected parsed field values.
    /// </summary>
    public class ParserSerializerTestcase
    {
        public string Name { get; set; }
        public string RawHex { get; set; }
        public string RootType { get; set; }

        /// <summary>The <c>&lt;xml&gt;</c> body's single child - the expected
        /// parsed tree.</summary>
        public XElement ExpectedXml { get; set; }

        /// <summary>Parser argument name → value.</summary>
        public Dictionary<string, string> ParserArguments { get; set; }
            = new Dictionary<string, string>();

        /// <summary>Decodes <see cref="RawHex"/> into bytes.</summary>
        public byte[] GetRawBytes()
        {
            if (string.IsNullOrEmpty(RawHex))
            {
                return System.Array.Empty<byte>();
            }
            var hex = RawHex.Replace(" ", "").Replace("\n", "").Replace("\r", "");
            if (hex.Length % 2 != 0)
                throw new System.FormatException(
                    $"Hex string has odd length ({hex.Length}); must be even. Raw: '{RawHex}'");
            var bytes = new byte[hex.Length / 2];
            for (var i = 0; i < bytes.Length; i++)
            {
                bytes[i] = System.Convert.ToByte(hex.Substring(i * 2, 2), 16);
            }
            return bytes;
        }
    }
}
