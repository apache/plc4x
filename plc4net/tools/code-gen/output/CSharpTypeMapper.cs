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
using org.apache.plc4net.tools.codegen.model;

namespace org.apache.plc4net.tools.codegen.output
{
    /// <summary>
    /// Maps an mspec <see cref="TypeReference"/> to the C# type, and to the
    /// <c>ReadBuffer</c> / <c>WriteBuffer</c> call that reads or writes it.
    /// </summary>
    /// <remarks>
    /// Targets plc4net's hand-written buffer API: <c>ReadUshort("name", 16)</c>
    /// with the width passed explicitly, no context stack and no byte-order
    /// argument (the buffers are MSB-first / big-endian, which is what every
    /// mspec-described protocol here uses).
    /// </remarks>
    public static class CSharpTypeMapper
    {
        /// <summary>The C# type used for a field or property of this reference.</summary>
        public static string CSharpType(TypeReference type) => type switch
        {
            EnumTypeReference e => e.Name,
            ComplexTypeReference c => c.Name,
            SimpleTypeReference s => SimpleCSharpType(s),
            _ => "object",
        };

        public static string SimpleCSharpType(SimpleTypeReference s) => s.BaseType switch
        {
            SimpleTypeReference.Base.Bit => "bool",
            SimpleTypeReference.Base.Byte => "byte",
            SimpleTypeReference.Base.UInt => s.SizeInBits <= 8 ? "byte"
                : s.SizeInBits <= 16 ? "ushort"
                : s.SizeInBits <= 32 ? "uint" : "ulong",
            SimpleTypeReference.Base.Int => s.SizeInBits <= 8 ? "sbyte"
                : s.SizeInBits <= 16 ? "short"
                : s.SizeInBits <= 32 ? "int" : "long",
            SimpleTypeReference.Base.Float or SimpleTypeReference.Base.UFloat =>
                s.SizeInBits <= 32 ? "float" : "double",
            SimpleTypeReference.Base.String or SimpleTypeReference.Base.VString => "string",
            SimpleTypeReference.Base.Time or SimpleTypeReference.Base.Date
                or SimpleTypeReference.Base.DateTime => "long",
            _ => "object",
        };

        /// <summary>
        /// The C# element type of an array of <paramref name="element"/>, and
        /// whether the array is <c>byte[]</c> (vs <c>List&lt;T&gt;</c>).
        /// </summary>
        public static (string Type, bool IsByteArray) ArrayType(TypeReference element)
        {
            if (element is SimpleTypeReference { IsByteBased: true })
            {
                return ("byte[]", true);
            }
            return ($"System.Collections.Generic.List<{CSharpType(element)}>", false);
        }

        /// <summary>Bit width of a fixed-size simple type; -1 if run-time sized.</summary>
        public static int FixedBitLength(TypeReference type) => type switch
        {
            SimpleTypeReference { BaseType: SimpleTypeReference.Base.VString } => -1,
            SimpleTypeReference s => s.SizeInBits,
            EnumTypeReference e => e.BaseType.SizeInBits,
            _ => -1,
        };

        /// <summary>
        /// A <c>ReadBuffer</c> call that reads one value of
        /// <paramref name="type"/>. <paramref name="logicalName"/> is the field
        /// name for diagnostics; <paramref name="complexArgs"/> is the rendered
        /// argument list for a complex type's <c>StaticParse</c>.
        /// </summary>
        public static string ReadCall(
            TypeReference type, string logicalName, string complexArgs = "")
        {
            switch (type)
            {
                case EnumTypeReference e:
                    return $"({e.Name}) {ReadCall(e.BaseType, logicalName)}";
                case ComplexTypeReference c:
                    return $"{c.Name}.StaticParse(readBuffer{Prefix(complexArgs)})";
                case SimpleTypeReference s:
                    return SimpleReadCall(s, logicalName);
                default:
                    throw new NotSupportedException($"No read for {type}");
            }
        }

        private static string SimpleReadCall(SimpleTypeReference s, string n)
        {
            var q = $"\"{n}\"";
            return s.BaseType switch
            {
                SimpleTypeReference.Base.Bit => $"readBuffer.ReadBit({q})",
                SimpleTypeReference.Base.Byte => $"readBuffer.ReadByte({q}, 8)",
                SimpleTypeReference.Base.UInt => s.SizeInBits <= 8 ? $"readBuffer.ReadByte({q}, {s.SizeInBits})"
                    : s.SizeInBits <= 16 ? $"readBuffer.ReadUshort({q}, {s.SizeInBits})"
                    : s.SizeInBits <= 32 ? $"readBuffer.ReadUint({q}, {s.SizeInBits})"
                    : $"readBuffer.ReadUlong({q}, {s.SizeInBits})",
                SimpleTypeReference.Base.Int => s.SizeInBits <= 8 ? $"readBuffer.ReadSbyte({q}, {s.SizeInBits})"
                    : s.SizeInBits <= 16 ? $"readBuffer.ReadShort({q}, {s.SizeInBits})"
                    : s.SizeInBits <= 32 ? $"readBuffer.ReadInt({q}, {s.SizeInBits})"
                    : $"readBuffer.ReadLong({q}, {s.SizeInBits})",
                SimpleTypeReference.Base.Float or SimpleTypeReference.Base.UFloat =>
                    s.SizeInBits <= 32 ? $"readBuffer.ReadFloat({q}, 32)" : $"readBuffer.ReadDouble({q}, 64)",
                SimpleTypeReference.Base.String =>
                    $"readBuffer.ReadString({q}, {s.SizeInBits}, {EncodingExpr(s.Encoding)})",
                _ => throw new NotSupportedException($"No read for {s}"),
            };
        }

        /// <summary>
        /// A <c>WriteBuffer</c> statement that writes <paramref name="valueExpr"/>
        /// as <paramref name="type"/>.
        /// </summary>
        public static string WriteCall(TypeReference type, string logicalName, string valueExpr)
        {
            var q = $"\"{logicalName}\"";
            switch (type)
            {
                case EnumTypeReference e:
                    return WriteCall(e.BaseType, logicalName,
                        $"({SimpleCSharpType(e.BaseType)}) {valueExpr}");
                case ComplexTypeReference _:
                    return $"{valueExpr}.Serialize(writeBuffer)";
                case SimpleTypeReference s:
                    return s.BaseType switch
                    {
                        SimpleTypeReference.Base.Bit => $"writeBuffer.WriteBit({q}, {valueExpr})",
                        SimpleTypeReference.Base.Byte => $"writeBuffer.WriteByte({q}, 8, {valueExpr})",
                        SimpleTypeReference.Base.UInt => s.SizeInBits <= 8 ? $"writeBuffer.WriteByte({q}, {s.SizeInBits}, {valueExpr})"
                            : s.SizeInBits <= 16 ? $"writeBuffer.WriteUshort({q}, {s.SizeInBits}, {valueExpr})"
                            : s.SizeInBits <= 32 ? $"writeBuffer.WriteUint({q}, {s.SizeInBits}, {valueExpr})"
                            : $"writeBuffer.WriteUlong({q}, {s.SizeInBits}, {valueExpr})",
                        SimpleTypeReference.Base.Int => s.SizeInBits <= 8 ? $"writeBuffer.WriteSbyte({q}, {s.SizeInBits}, {valueExpr})"
                            : s.SizeInBits <= 16 ? $"writeBuffer.WriteShort({q}, {s.SizeInBits}, {valueExpr})"
                            : s.SizeInBits <= 32 ? $"writeBuffer.WriteInt({q}, {s.SizeInBits}, {valueExpr})"
                            : $"writeBuffer.WriteLong({q}, {s.SizeInBits}, {valueExpr})",
                        SimpleTypeReference.Base.Float or SimpleTypeReference.Base.UFloat =>
                            s.SizeInBits <= 32 ? $"writeBuffer.WriteFloat({q}, 32, {valueExpr})"
                                : $"writeBuffer.WriteDouble({q}, 64, {valueExpr})",
                        SimpleTypeReference.Base.String =>
                            $"writeBuffer.WriteString({q}, {s.SizeInBits}, \"{s.Encoding ?? "UTF8"}\", {valueExpr})",
                        _ => throw new NotSupportedException($"No write for {s}"),
                    };
                default:
                    throw new NotSupportedException($"No write for {type}");
            }
        }

        /// <summary>The <c>System.Text.Encoding</c> expression for an mspec
        /// string encoding name (<c>"UTF8"</c>, <c>"UTF16BE"</c>, …).</summary>
        public static string EncodingExpr(string encoding) =>
            (encoding ?? "UTF8").Replace("-", "").ToUpperInvariant() switch
            {
                "UTF16" or "UTF16LE" => "System.Text.Encoding.Unicode",
                "UTF16BE" => "System.Text.Encoding.BigEndianUnicode",
                "ASCII" => "System.Text.Encoding.ASCII",
                _ => "System.Text.Encoding.UTF8",
            };

        private static string Prefix(string args) =>
            string.IsNullOrEmpty(args) ? "" : ", " + args;
    }
}
