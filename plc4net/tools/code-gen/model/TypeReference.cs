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
using org.apache.plc4net.tools.codegen.model.terms;

namespace org.apache.plc4net.tools.codegen.model
{
    /// <summary>
    /// The type of a field or argument as written in the mspec: either a
    /// <c>dataType</c> (<c>uint 16</c>, <c>bit</c>, <c>string 8</c>) or a
    /// reference to another named type, optionally with parser arguments
    /// (<c>ModbusPDU('response')</c>).
    /// </summary>
    public abstract class TypeReference
    {
    }

    /// <summary>An mspec <c>dataType</c> - a primitive on the wire.</summary>
    public sealed class SimpleTypeReference : TypeReference
    {
        public enum Base
        {
            Bit,
            Byte,
            UInt,
            Int,
            Float,
            UFloat,
            String,
            VString,
            Time,
            Date,
            DateTime,
        }

        public Base BaseType { get; set; }

        /// <summary>Width in bits. 1 for <c>bit</c>, 8 for <c>byte</c>, the
        /// declared size for <c>uint N</c> etc. Zero for a <c>vstring</c> whose
        /// length is an expression (see <see cref="LengthExpression"/>).</summary>
        public int SizeInBits { get; set; }

        /// <summary>For <c>vstring</c>: the run-time length expression.</summary>
        public Term LengthExpression { get; set; }

        /// <summary>The <c>stringEncoding</c> attribute on the field, if any
        /// (e.g. <c>"UTF8"</c>, <c>"UTF16BE"</c>).</summary>
        public string Encoding { get; set; }

        public bool IsIntegerLike =>
            BaseType is Base.Bit or Base.Byte or Base.UInt or Base.Int;

        public bool IsSigned => BaseType == Base.Int;

        public bool IsByteBased => BaseType == Base.Byte && SizeInBits == 8;

        public override string ToString() =>
            BaseType == Base.Bit ? "bit"
            : BaseType == Base.Byte ? "byte"
            : $"{BaseType.ToString().ToLowerInvariant()} {SizeInBits}";
    }

    /// <summary>A reference to another <c>[type]</c> / <c>[discriminatedType]</c>.</summary>
    public sealed class ComplexTypeReference : TypeReference
    {
        public string Name { get; set; }

        /// <summary>Parser arguments passed to the referenced type, e.g. the
        /// <c>'response'</c> in <c>ModbusPDU('response')</c>.</summary>
        public IReadOnlyList<Term> Arguments { get; set; } = System.Array.Empty<Term>();

        public override string ToString() =>
            Arguments.Count == 0 ? Name : $"{Name}({string.Join(", ", Arguments)})";
    }

    /// <summary>
    /// A reference resolved to an <c>[enum]</c> type. mspec does not spell
    /// these differently from a complex reference - the builder promotes a
    /// <see cref="ComplexTypeReference"/> to this once the enum types are known.
    /// </summary>
    public sealed class EnumTypeReference : TypeReference
    {
        public string Name { get; set; }

        /// <summary>The enum's underlying wire type (<c>uint 8</c> etc.).</summary>
        public SimpleTypeReference BaseType { get; set; }

        public override string ToString() => Name;
    }
}
