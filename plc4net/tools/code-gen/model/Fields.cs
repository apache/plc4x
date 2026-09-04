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

namespace org.apache.plc4net.tools.codegen.model.fields
{
    /// <summary>Base for the mspec field keywords.</summary>
    public abstract class Field
    {
        public string Name { get; set; }

        /// <summary>The field's declared type. Null for <c>typeSwitch</c>.</summary>
        public TypeReference Type { get; set; }

        public Dictionary<string, Term> Attributes { get; } = new Dictionary<string, Term>();

        /// <summary>Whether the field holds a value the message carries (a
        /// property + constructor parameter). False for the fields that exist
        /// only for the wire format.</summary>
        public virtual bool IsProperty => false;

        /// <summary>The mspec keyword, for diagnostics and switch templates.</summary>
        public string Keyword => GetType().Name.Replace("Field", "").ToLowerInvariant();
    }

    /// <summary><c>[simple uint 16 startingAddress]</c></summary>
    public sealed class SimpleField : Field
    {
        public override bool IsProperty => true;
    }

    /// <summary><c>[const uint 16 protocolIdentifier 0x0000]</c> - written on
    /// serialize, checked on parse, exposed as a C# <c>const</c>.</summary>
    public sealed class ConstField : Field
    {
        public Term ReferenceValue { get; set; }
    }

    /// <summary><c>[implicit uint 16 length 'pdu.lengthInBytes + 1']</c> -
    /// computed on serialize, read-and-kept-as-a-local on parse.</summary>
    public sealed class ImplicitField : Field
    {
        public Term SerializeExpression { get; set; }
    }

    /// <summary><c>[reserved uint 15 '0x0000']</c> - no name, no property.</summary>
    public sealed class ReservedField : Field
    {
        public Term ReferenceValue { get; set; }
    }

    /// <summary><c>[virtual uint 16 foo 'a + b']</c> - never on the wire, a
    /// computed accessor.</summary>
    public sealed class VirtualField : Field
    {
        public Term ValueExpression { get; set; }
        public override bool IsProperty => false;
    }

    /// <summary>
    /// <c>[discriminator uint 7 functionFlag]</c> - read first on the parent so
    /// the type switch can pick the child. Exposed as an (abstract on the
    /// parent, override on each child) accessor, not a stored property.
    /// </summary>
    public sealed class DiscriminatorField : Field
    {
    }

    /// <summary><c>[enum TransportSize transportSize code]</c> - an enum-typed
    /// field whose wire value is keyed on one of the enum's attributes rather
    /// than the enum constant itself. <see cref="Field.Type"/> is an
    /// <see cref="EnumTypeReference"/>.</summary>
    public sealed class EnumField : Field
    {
        /// <summary>The enum attribute the wire value maps through - <c>code</c>
        /// in <c>[enum TransportSize transportSize code]</c>. Null for the bare
        /// two-token form, where the wire value is the enum constant itself.</summary>
        public string KeyAccessor { get; set; }

        public override bool IsProperty => true;
    }

    /// <summary><c>[array byte value count 'byteCount']</c></summary>
    public sealed class ArrayField : Field
    {
        public enum Loop { Count, Length, Terminated }

        public Loop LoopType { get; set; }

        /// <summary>The count / byte-length / terminator expression.</summary>
        public Term LoopExpression { get; set; }

        public override bool IsProperty => true;
    }

    /// <summary><c>[checksum uint 16 crc 'STATIC_CALL("rtuCrcCheck", ...)']</c></summary>
    public sealed class ChecksumField : Field
    {
        public Term ChecksumExpression { get; set; }
    }

    /// <summary>
    /// <c>[typeSwitch driverType, response ['MODBUS_TCP' ModbusTcpADU [...]]]</c>
    /// - the discriminated dispatch. The case sub-types are lifted to top-level
    /// <see cref="ComplexTypeDefinition"/>s by the builder; the names are kept
    /// here in declaration order for the parse/serialize switch.
    /// </summary>
    public sealed class TypeSwitchField : Field
    {
        /// <summary>The discriminator expressions to switch on, e.g.
        /// <c>errorFlag</c>, <c>functionFlag</c>, <c>response</c>.</summary>
        public IReadOnlyList<Term> Discriminators { get; set; } = System.Array.Empty<Term>();

        /// <summary>Child type names in declaration order.</summary>
        public List<string> CaseNames { get; } = new List<string>();
    }

    /// <summary>
    /// <c>[optional uint 16 errorCode 'condition']</c> - present on the wire
    /// only when the condition holds. The property is nullable.
    /// </summary>
    public sealed class OptionalField : Field
    {
        /// <summary>The presence condition, or null (always present - rare).</summary>
        public Term Condition { get; set; }

        public override bool IsProperty => true;
    }

    /// <summary>
    /// <c>[manual Type name 'parse' 'serialize' 'length']</c> - the codec is
    /// three hand-supplied expressions (usually <c>STATIC_CALL</c>s that take
    /// <c>readBuffer</c> / <c>writeBuffer</c> / <c>_value</c>).
    /// </summary>
    public sealed class ManualField : Field
    {
        public Term ParseExpression { get; set; }
        public Term SerializeExpression { get; set; }
        public Term LengthExpression { get; set; }

        public override bool IsProperty => true;
    }

    /// <summary>
    /// <c>[padding uint 8 pad '0x00' 'timesPadding']</c> - writes the padding
    /// value <c>timesPadding</c> times; on parse, reads and drops that many.
    /// </summary>
    public sealed class PaddingField : Field
    {
        public Term PaddingValue { get; set; }
        public Term TimesPadding { get; set; }
    }

    /// <summary>
    /// A field keyword the generator does not handle yet
    /// (<c>abstract</c>, <c>assert</c>, <c>manualArray</c>, <c>peek</c>,
    /// <c>unknown</c>, <c>validation</c>). Kept in the model so a type that
    /// uses one still round-trips through the builder and the gap is visible
    /// in the output.
    /// </summary>
    public sealed class UnsupportedField : Field
    {
        public string MspecKeyword { get; set; }
        public string RawText { get; set; }
    }
}
