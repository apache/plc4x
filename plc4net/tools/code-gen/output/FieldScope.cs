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
using org.apache.plc4net.tools.codegen.model;
using org.apache.plc4net.tools.codegen.model.fields;

namespace org.apache.plc4net.tools.codegen.output
{
    /// <summary>
    /// Resolves the names in an mspec expression against one complex type's
    /// fields and parser arguments. During parse, a field name is a local
    /// holding the value just read; during serialize, it is the property. An
    /// enum member is <c>Enum.Member</c>; an unknown name is left alone.
    /// </summary>
    internal sealed class FieldScope : IExpressionScope
    {
        private readonly Dictionary<string, TypeReference> _byName = new Dictionary<string, TypeReference>();
        private readonly HashSet<string> _enumNames;
        // constant name -> "EnumName.CONSTANT" (first enum wins on a clash)
        private readonly Dictionary<string, string> _enumConstants;
        // every enum attribute-accessor name, e.g. `sizeInBits`, `code` - a
        // member access on one of these is a Get… extension-method call.
        private readonly HashSet<string> _enumAccessors;
        private readonly bool _serialize;

        public FieldScope(Protocol protocol, ComplexTypeDefinition type, bool serialize = false)
        {
            _serialize = serialize;
            _hoistedLocals = new HashSet<string>();
            _enumNames = new HashSet<string>(protocol.Enums.Select(e => e.Name));
            _enumAccessors = new HashSet<string>(
                protocol.Enums.SelectMany(e => e.Arguments.Select(a => a.Name)));
            _enumConstants = new Dictionary<string, string>();
            foreach (var e in protocol.Enums)
            {
                foreach (var v in e.Values)
                {
                    _enumConstants.TryAdd(v.Name, $"{e.Name}.{v.Name}");
                }
            }

            var owner = type.ParentName != null ? protocol.FindType(type.ParentName) : type;
            if (owner != null)
            {
                foreach (var a in owner.Arguments)
                {
                    _byName[a.Name] = a.Type;
                }
                foreach (var f in owner.Fields.Where(f => f.Name != null))
                {
                    _byName[f.Name] = f.Type;
                }
            }
            foreach (var f in type.Fields.Where(f => f.Name != null))
            {
                _byName[f.Name] = f.Type;
            }
        }

        private FieldScope(FieldScope src, bool serialize)
        {
            _serialize = serialize;
            _enumNames = src._enumNames;
            _enumAccessors = src._enumAccessors;
            _enumConstants = src._enumConstants;
            _byName = src._byName;
            _hoistedLocals = src._hoistedLocals;
        }

        public FieldScope ForSerialize() => new FieldScope(this, serialize: true);

        /// <summary>The C# for <c>_value</c> - the field's own value, set by the
        /// generator around a <c>manual</c> field's serialize expression.</summary>
        public string CurrentValue { get; set; }

        /// <summary>Names the generator has hoisted as a plain local (a
        /// recomputed implicit field in a child's serialize / length) - these
        /// resolve to the camelCase local, not a Pascal property. Shared with
        /// the serialize view so a mark on either is seen by both.</summary>
        private readonly HashSet<string> _hoistedLocals;

        public void MarkLocal(string name) => _hoistedLocals.Add(name);

        public void ClearLocals() => _hoistedLocals.Clear();

        public string ResolveReference(string name)
        {
            switch (name)
            {
                // mspec context variables available inside manual / expression fields
                case "readBuffer":
                case "writeBuffer":
                    return name;
                case "_lastItem":
                    // the generator gives a type that needs it a `_lastItem`
                    // parameter, threaded from the enclosing array loop.
                    return "_lastItem";
                case "curPos":
                case "_curPos":
                    // mspec curPos is bytes consumed since this type's parse
                    // began; the generator captures _startPos (bits) at entry.
                    return _serialize
                        ? "((writeBuffer.GetPos() - _startPos) / 8)"
                        : "((readBuffer.GetPos() - _startPos) / 8)";
                case "_value":
                    return CurrentValue ?? "_value";
                case "_type":
                    return "this";
                // the length of the type being defined, used bare inside an
                // implicit / virtual field's expression
                case "lengthInBytes":
                    return "GetLengthInBytes()";
                case "lengthInBits":
                    return "GetLengthInBits()";
            }

            if (_enumNames.Contains(name))
            {
                return name;
            }
            if (_hoistedLocals.Contains(name))
            {
                return CSharpGenerator.Camel(name);
            }
            if (_byName.ContainsKey(name))
            {
                return _serialize ? CSharpGenerator.Pascal(name) : CSharpGenerator.Camel(name);
            }
            // A bare enum constant used in a discriminator / expression, e.g.
            // `BYTE_WORD_DWORD` for a `DataTransportSize`-typed const field.
            if (_enumConstants.TryGetValue(name, out var qualified))
            {
                return qualified;
            }
            return null;
        }

        /// <summary>Resolves only bare enum-constant names - for rendering an
        /// enum's own value expressions and its parameter tables, where there
        /// are no fields in scope.</summary>
        public static IExpressionScope EnumConstantsOnly(Protocol protocol)
            => new EnumScope(protocol);

        private sealed class EnumScope : IExpressionScope
        {
            private readonly Dictionary<string, string> _constants = new Dictionary<string, string>();

            public EnumScope(Protocol protocol)
            {
                foreach (var e in protocol.Enums)
                {
                    foreach (var v in e.Values)
                    {
                        _constants.TryAdd(v.Name, $"{e.Name}.{v.Name}");
                    }
                }
            }

            public string ResolveReference(string name)
                => _constants.TryGetValue(name, out var q) ? q : null;
        }

        public string ResolveMember(string name)
            => _enumAccessors.Contains(name)
                ? $"Get{CSharpGenerator.Pascal(name)}()"
                : null;

        public string ResolveCount(string name)
        {
            var reference = ResolveReference(name) ?? name;
            if (_byName.TryGetValue(name, out var type)
                && type is SimpleTypeReference { IsByteBased: true })
            {
                return reference + ".Length";
            }
            // A byte-based array field is byte[]; anything else is List<T>.
            return reference + ".Count";
        }
    }
}
