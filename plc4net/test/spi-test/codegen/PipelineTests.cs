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
using org.apache.plc4net.tools.codegen.model;
using org.apache.plc4net.tools.codegen.model.fields;
using org.apache.plc4net.tools.codegen.output;
using Xunit;

namespace org.apache.plc4net.spi.test.codegen
{
    /// <summary>
    /// End-to-end for the pure-.NET mspec toolchain: text → type-model IR
    /// (<see cref="MspecModelBuilder"/>) → C# (<see cref="CSharpGenerator"/>).
    /// The generated Modbus model is compiled for real by the modbus project
    /// (and exercised against the test vectors in
    /// <see cref="ModbusGeneratedRoundTripTests"/>); this covers the shapes
    /// the IR and the emitter have to get right.
    /// </summary>
    public class PipelineTests
    {
        // ── IR ──────────────────────────────────────────────────

        [Fact]
        public void A_plain_type_has_its_fields_in_order()
        {
            var model = MspecModelBuilder.Build(@"
[type Item
    [simple uint 8  referenceType]
    [simple uint 16 fileNumber]
]
");
            var item = Assert.Single(model.Types);
            Assert.Equal("Item", item.Name);
            Assert.Equal(new[] { "referenceType", "fileNumber" }, item.Fields.Select(f => f.Name));
            var f0 = Assert.IsType<SimpleField>(item.Fields[0]);
            var t0 = Assert.IsType<SimpleTypeReference>(f0.Type);
            Assert.Equal(8, t0.SizeInBits);
        }

        [Fact]
        public void A_discriminated_type_lifts_its_cases_to_top_level_children()
        {
            var model = MspecModelBuilder.Build(@"
[discriminatedType Pdu(bit response)
    [discriminator uint 8 code]
    [typeSwitch code
        ['0x01' PduA [simple uint 16 a]]
        ['0x02' PduB [simple uint 16 b]]
    ]
]
");
            var parent = model.FindType("Pdu");
            Assert.True(parent.IsDiscriminatedParent);

            var a = model.FindType("PduA");
            Assert.Equal("Pdu", a.ParentName);
            Assert.True(a.IsDiscriminatedChild);
            Assert.Single(a.DiscriminatorValues);

            Assert.Equal(new[] { "PduA", "PduB" }, parent.TypeSwitch.CaseNames);
        }

        [Fact]
        public void An_enum_carries_its_base_type_and_values()
        {
            var model = MspecModelBuilder.Build("[enum uint 8 E ['1' A] ['2' B]]");
            var e = Assert.Single(model.Enums);
            Assert.Equal(SimpleTypeReference.Base.UInt, e.BaseType.BaseType);
            Assert.Equal(8, e.BaseType.SizeInBits);
            Assert.Equal(new[] { "A", "B" }, e.Values.Select(v => v.Name));
        }

        [Fact]
        public void An_enum_typed_field_resolves_to_an_enum_reference()
        {
            var model = MspecModelBuilder.Build(@"
[enum uint 8 ErrorCode ['1' BAD]]
[type Err [simple ErrorCode code]]
");
            var field = model.FindType("Err").Fields[0];
            Assert.IsType<EnumTypeReference>(field.Type);
        }

        // ── generation ──────────────────────────────────────────

        [Fact]
        public void Generates_a_file_per_type_and_enum()
        {
            var model = MspecModelBuilder.Build(@"
[enum uint 8 E ['1' A]]
[type T [simple uint 8 x]]
");
            var files = new CSharpGenerator(model, "demo", "demo.readwrite").Generate();
            Assert.Contains("model/T.cs", files.Keys);
            Assert.Contains("model/E.cs", files.Keys);
        }

        [Fact]
        public void A_generated_type_has_parse_serialize_and_length()
        {
            var model = MspecModelBuilder.Build("[type T [simple uint 16 startingAddress]]");
            var code = new CSharpGenerator(model, "demo", "demo.readwrite").Generate()["model/T.cs"];

            Assert.Contains("public static T StaticParse(ReadBuffer readBuffer)", code);
            Assert.Contains("readBuffer.ReadUshort(\"startingAddress\", 16)", code);
            Assert.Contains("public void Serialize(WriteBuffer writeBuffer)", code);
            Assert.Contains("writeBuffer.WriteUshort(\"startingAddress\", 16, StartingAddress)", code);
            Assert.Contains("public int GetLengthInBits()", code);
        }

        [Fact]
        public void An_implicit_field_is_read_and_discarded_but_written_from_its_formula()
        {
            var model = MspecModelBuilder.Build(@"
[type T
    [implicit uint 8 byteCount 'COUNT(value)']
    [array    byte   value count 'byteCount']
]
");
            var code = new CSharpGenerator(model, "demo", "demo.readwrite").Generate()["model/T.cs"];

            // parse: byteCount read into a local, then used as the array length
            Assert.Contains("var byteCount = readBuffer.ReadByte(\"byteCount\", 8);", code);
            Assert.Contains("readBuffer.ReadByteArray(\"value\", (int) (byteCount) * 8)", code);
            // serialize: byteCount computed from the array, not stored
            Assert.Contains("writeBuffer.WriteByte(\"byteCount\", 8, (byte) (Value.Length));", code);
            Assert.DoesNotContain("public byte ByteCount", code);
        }

        [Fact]
        public void An_optional_field_is_nullable_and_conditional()
        {
            var model = MspecModelBuilder.Build(@"
[type T
    [simple   uint 16 dataLength]
    [optional uint 8  extra 'dataLength >= 12']
]
");
            var code = new CSharpGenerator(model, "demo", "demo.readwrite").Generate()["model/T.cs"];

            Assert.Contains("public byte? Extra { get; }", code);
            // parse: only read when the condition holds
            Assert.Contains("byte? extra = null;", code);
            Assert.Contains("if ((dataLength >= 12))", code);
            // serialize / length: guarded by the null check
            Assert.Contains("if (Extra != null)", code);
            Assert.Contains("(Extra != null ? 8 : 0)", code);
        }

        [Fact]
        public void A_padding_field_loops_the_repeat_count()
        {
            var model = MspecModelBuilder.Build(
                "[type T [simple uint 8 n] [padding uint 8 pad '0x00' 'n % 2']]");
            var code = new CSharpGenerator(model, "demo", "demo.readwrite").Generate()["model/T.cs"];

            Assert.Contains("var _timesPadding = (int) ((n % 2));", code);
            Assert.Contains("while (_timesPadding-- > 0)", code);
        }

        [Fact]
        public void The_real_modbus_mspec_generates_without_throwing()
        {
            var repoRoot = RepoPaths.FindRepoRoot();
            if (repoRoot == null)
            {
                return;
            }
            var mspec = Path.Combine(repoRoot,
                "protocols", "modbus", "src", "main", "resources",
                "protocols", "modbus", "modbus.mspec");

            var model = MspecModelBuilder.BuildFile(mspec);
            var files = new CSharpGenerator(
                model, "modbus", "org.apache.plc4net.drivers.modbus.readwrite").Generate();

            Assert.True(model.Types.Count > 40);
            Assert.Contains("model/ModbusADU.cs", files.Keys);
            Assert.Contains("model/ModbusPDU.cs", files.Keys);
            Assert.All(files.Values, src => Assert.Contains("DO NOT EDIT", src));
        }
    }
}
