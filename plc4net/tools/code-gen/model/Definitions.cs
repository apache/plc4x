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
using System.Linq;
using org.apache.plc4net.tools.codegen.model.terms;

namespace org.apache.plc4net.tools.codegen.model
{
    /// <summary>A top-level named type in an mspec file.</summary>
    public abstract class TypeDefinition
    {
        public string Name { get; set; }

        /// <summary>Attributes on the type header, e.g.
        /// <c>byteOrder='"BIG_ENDIAN"'</c>.</summary>
        public Dictionary<string, Term> Attributes { get; } = new Dictionary<string, Term>();

        public Term GetAttribute(string name) =>
            Attributes.TryGetValue(name, out var v) ? v : null;
    }

    /// <summary>A parser argument in a type header: <c>(bit response)</c>.</summary>
    public sealed class Argument
    {
        public string Name { get; set; }
        public TypeReference Type { get; set; }
    }

    /// <summary>
    /// A structured type - a <c>[type]</c> or a <c>[discriminatedType]</c>.
    /// A discriminated parent is the one whose fields include a
    /// <see cref="fields.TypeSwitchField"/>; its children are also
    /// <see cref="ComplexTypeDefinition"/> instances, linked by
    /// <see cref="ParentName"/> / <see cref="DiscriminatorValues"/>.
    /// </summary>
    public sealed class ComplexTypeDefinition : TypeDefinition
    {
        public List<Argument> Arguments { get; } = new List<Argument>();

        public List<fields.Field> Fields { get; } = new List<fields.Field>();

        // ── discriminated-type wiring ────────────────────────────

        /// <summary>Set on a child type: the parent it extends.</summary>
        public string ParentName { get; set; }

        /// <summary>Set on a child type: the case values, positionally matched
        /// against the parent's <see cref="fields.TypeSwitchField.Discriminators"/>.
        /// A null entry is a wildcard (the case listed fewer values).</summary>
        public IReadOnlyList<Term> DiscriminatorValues { get; set; }

        public bool IsDiscriminatedParent =>
            Fields.Any(f => f is fields.TypeSwitchField);

        public bool IsDiscriminatedChild => ParentName != null;

        public fields.TypeSwitchField TypeSwitch =>
            Fields.OfType<fields.TypeSwitchField>().FirstOrDefault();

        /// <summary>The fields that hold a value and therefore get a property
        /// and a constructor parameter (everything except const / implicit /
        /// reserved / virtual / typeSwitch and the abstract markers).</summary>
        public IEnumerable<fields.Field> PropertyFields =>
            Fields.Where(f => f.IsProperty);

        /// <summary>
        /// On a discriminated parent: the property fields declared before the
        /// <c>typeSwitch</c>. A child inherits these - the parent reads and
        /// writes them, and passes them to the child's <c>StaticParse</c> and
        /// on to <c>base(...)</c>.
        /// </summary>
        public IEnumerable<fields.Field> PrefixPropertyFields =>
            Fields.TakeWhile(f => f is not fields.TypeSwitchField).Where(f => f.IsProperty);

        /// <summary>
        /// On a discriminated parent: the property fields declared after the
        /// <c>typeSwitch</c> - a shared trailer (S7's COTPPacket parameters /
        /// payload). The child reads and writes these, after its own body.
        /// </summary>
        public IEnumerable<fields.Field> SuffixPropertyFields =>
            Fields.SkipWhile(f => f is not fields.TypeSwitchField).Skip(1).Where(f => f.IsProperty);

        /// <summary>
        /// On a discriminated parent: the pre-<c>typeSwitch</c> fields that
        /// carry a value a child might reference in an expression - property
        /// fields plus discriminators, implicit and const fields. The parent
        /// reads them and passes them to the child's <c>StaticParse</c>.
        /// </summary>
        public IEnumerable<fields.Field> PrefixContextFields =>
            Fields.TakeWhile(f => f is not fields.TypeSwitchField)
                  .Where(f => f.Name != null
                      && f is not fields.ReservedField
                      && f is not fields.VirtualField);

        /// <summary>Pre-<c>typeSwitch</c> implicit fields - a child recomputes
        /// these at the top of its serialize / length.</summary>
        public IEnumerable<fields.ImplicitField> PrefixImplicitFields =>
            Fields.TakeWhile(f => f is not fields.TypeSwitchField).OfType<fields.ImplicitField>();

        /// <summary>All fields up to and including the <c>typeSwitch</c> - what
        /// a discriminated parent's own <c>StaticParse</c> / <c>Serialize</c>
        /// covers.</summary>
        public IEnumerable<fields.Field> FieldsThroughSwitch
        {
            get
            {
                foreach (var f in Fields)
                {
                    yield return f;
                    if (f is fields.TypeSwitchField)
                    {
                        yield break;
                    }
                }
            }
        }

        /// <summary>The fields after the <c>typeSwitch</c> in declaration
        /// order (a superset of <see cref="SuffixPropertyFields"/> - includes
        /// reserved / const trailers too).</summary>
        public IEnumerable<fields.Field> FieldsAfterSwitch =>
            Fields.SkipWhile(f => f is not fields.TypeSwitchField).Skip(1);

        public IEnumerable<fields.ConstField> ConstFields => Fields.OfType<fields.ConstField>();
    }

    /// <summary>An <c>[enum]</c> type.</summary>
    public sealed class EnumTypeDefinition : TypeDefinition
    {
        /// <summary>The wire type: <c>uint 8</c> for
        /// <c>[enum uint 8 ModbusErrorCode ...]</c>. Null for a bare
        /// <c>[enum Name ...]</c> (rare; treated as an int).</summary>
        public SimpleTypeReference BaseType { get; set; }

        /// <summary>Per-constant parameters: <c>(uint 8 dataTypeSize)</c>.</summary>
        public List<Argument> Arguments { get; } = new List<Argument>();

        public List<EnumValue> Values { get; } = new List<EnumValue>();
    }

    /// <summary>One constant in an enum.</summary>
    public sealed class EnumValue
    {
        public string Name { get; set; }

        /// <summary>The wire value, e.g. <c>0x02</c>.</summary>
        public Term Value { get; set; }

        /// <summary>Values for the enum's <see cref="EnumTypeDefinition.Arguments"/>,
        /// positionally: <c>['1' BOOL ['2']]</c> gives <c>dataTypeSize = 2</c>.</summary>
        public IReadOnlyList<Term> ConstantValues { get; set; } = System.Array.Empty<Term>();
    }

    /// <summary>
    /// A <c>[dataIo]</c> type - parses to / serializes from an
    /// <c>IPlcValue</c> rather than a generated class. The <c>typeSwitch</c>
    /// picks the wire layout from the parser arguments (S7's <c>DataItem</c>
    /// keys on <c>dataProtocolId</c>); each case yields one <c>PlcValue</c>.
    /// </summary>
    public sealed class DataIoTypeDefinition : TypeDefinition
    {
        public List<Argument> Arguments { get; } = new List<Argument>();
        public fields.TypeSwitchField TypeSwitch { get; set; }

        /// <summary>The <c>typeSwitch</c> cases, in declaration order. Unlike a
        /// discriminated type's children these are not lifted to top-level
        /// types - the case names repeat (<c>CHAR</c>, <c>STRING</c>, <c>TIME</c>
        /// each appear twice) and a case is a value shape, not a class.</summary>
        public List<ComplexTypeDefinition> Cases { get; } = new List<ComplexTypeDefinition>();
    }
}
