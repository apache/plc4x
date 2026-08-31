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
using System.Collections.Generic;
using System.Linq;
using org.apache.plc4net.tools.codegen.model;
using org.apache.plc4net.tools.codegen.model.fields;
using org.apache.plc4net.tools.codegen.model.terms;

namespace org.apache.plc4net.tools.codegen.output
{
    /// <summary>
    /// Emits C# for an mspec <see cref="Protocol"/>: one file per complex type
    /// (class + <c>StaticParse</c> + <c>Serialize</c> + length) and per enum,
    /// plus a constants class. This is the pure-.NET replacement for the Java
    /// freemarker <c>model-template</c> / <c>io-template</c>.
    /// </summary>
    public sealed class CSharpGenerator
    {
        private readonly Protocol _protocol;
        private readonly string _namespace;
        private readonly string _protocolName;

        public CSharpGenerator(Protocol protocol, string protocolName, string namespaceName)
        {
            _protocol = protocol;
            _protocolName = protocolName;
            _namespace = namespaceName;
        }

        /// <summary>Generates every file; key is a relative path.</summary>
        public Dictionary<string, string> Generate()
        {
            var files = new Dictionary<string, string>();

            foreach (var e in _protocol.Enums)
            {
                files[$"model/{e.Name}.cs"] = EmitEnum(e);
            }
            foreach (var t in _protocol.Types)
            {
                files[$"model/{t.Name}.cs"] = EmitComplexType(t);
            }
            foreach (var dio in _protocol.DataIos)
            {
                files[$"model/{dio.Name}.cs"] = EmitDataIo(dio);
            }
            if (_protocol.Constants.Count > 0)
            {
                files[$"model/{Pascal(_protocolName)}Constants.cs"] = EmitConstants();
            }

            files[$"model/{Pascal(_protocolName)}StaticHelper.cs"] = EmitStaticHelper();
            return files;
        }

        // ── complex types ───────────────────────────────────────

        private string EmitComplexType(ComplexTypeDefinition type)
        {
            var scope = new FieldScope(_protocol, type);
            var render = new CSharpExpressionRenderer(scope)
            {
                StaticHelperClass = $"{Pascal(_protocolName)}StaticHelper",
            };

            var c = new CodeWriter();
            Header(c);
            c.Line("using System;");
            c.Line("using System.Linq;");
            c.Line("using org.apache.plc4net.spi.drivers;");
            c.Line("using org.apache.plc4net.spi.generation;");
            c.Line();
            c.Line($"namespace {_namespace}.model");
            c.Line("{");
            c.Indent();

            var parent = type.IsDiscriminatedChild ? $" : {type.ParentName}" : " : IMessage";
            var mods = type.IsDiscriminatedParent ? "abstract partial" : "partial";
            c.Line($"public {mods} class {type.Name}{parent}");
            c.Line("{");
            c.Indent();

            EmitConstFields(c, type, render);
            EmitDiscriminatorAccessors(c, type, render);
            EmitProperties(c, type);
            EmitConstructor(c, type);
            EmitStaticParse(c, type, render, scope);
            EmitSerialize(c, type, render, scope);
            EmitLength(c, type, render, scope);

            c.Outdent();
            c.Line("}");
            c.Outdent();
            c.Line("}");
            return c.ToString();
        }

        private void EmitConstFields(CodeWriter c, ComplexTypeDefinition type, CSharpExpressionRenderer r)
        {
            foreach (var f in type.ConstFields)
            {
                c.Line($"public const {CSharpTypeMapper.CSharpType(f.Type)} {Pascal(f.Name)} = {r.Render(f.ReferenceValue)};");
            }
            if (type.ConstFields.Any())
            {
                c.Line();
            }
        }

        private void EmitDiscriminatorAccessors(CodeWriter c, ComplexTypeDefinition type, CSharpExpressionRenderer r)
        {
            if (type.IsDiscriminatedParent)
            {
                foreach (var d in type.Fields.OfType<DiscriminatorField>())
                {
                    c.Line($"public abstract {CSharpTypeMapper.CSharpType(d.Type)} {Pascal(d.Name)} {{ get; }}");
                }
                // A parser-arg discriminator with no matching field becomes an
                // abstract accessor too, so children can pin their value.
                foreach (var name in ParserArgDiscriminatorNames(type))
                {
                    var arg = type.Arguments.First(a => a.Name == name);
                    c.Line($"public abstract {CSharpTypeMapper.CSharpType(arg.Type)} {Pascal(name)} {{ get; }}");
                }
                c.Line();
            }

            if (type.IsDiscriminatedChild && _protocol.FindType(type.ParentName) is { } p)
            {
                var values = type.DiscriminatorValues ?? Array.Empty<Term>();
                var discriminators = DiscriminatorList(p, r);
                var emitted = false;
                for (var i = 0; i < discriminators.Count; i++)
                {
                    var d = discriminators[i];
                    if (d.AccessorName == null)
                    {
                        continue; // property field or dotted path - no accessor
                    }
                    var pinned = values.Count > i && values[i] != null;
                    // An unpinned discriminator (the case listed fewer values,
                    // or a null slot) takes the type's default, matching plc4j:
                    // ModbusPDUError leaves functionFlag / response unspecified
                    // and its accessors return 0 / false.
                    var rhs = pinned
                        ? RenderDiscriminatorValue(values[i], d.Type, r)
                        : $"default({CSharpTypeMapper.CSharpType(d.Type)})";
                    c.Line($"public override {CSharpTypeMapper.CSharpType(d.Type)} {Pascal(d.AccessorName)} => {rhs};");
                    emitted = true;
                }
                if (emitted)
                {
                    c.Line();
                }
            }
        }

        private void EmitProperties(CodeWriter c, ComplexTypeDefinition type)
        {
            foreach (var f in type.PropertyFields)
            {
                c.Line($"public {PropertyType(f)} {Pascal(f.Name)} {{ get; }}");
            }
            if (type.PropertyFields.Any())
            {
                c.Line();
            }
        }

        private void EmitConstructor(CodeWriter c, ComplexTypeDefinition type)
        {
            var all = AllValueFields(type);
            if (all.Count == 0)
            {
                return;
            }

            var own = type.PropertyFields.ToList();
            var inherited = ParentPrefixFields(type).Concat(ParentSuffixFields(type)).ToList();

            var pars = string.Join(", ", all.Select(f => $"{PropertyType(f)} {Camel(f.Name)}"));
            var mod = type.IsDiscriminatedParent ? "protected" : "public";
            var baseCall = inherited.Count > 0
                ? " : base(" + string.Join(", ", inherited.Select(f => Camel(f.Name))) + ")"
                : "";

            c.Line($"{mod} {type.Name}({pars}){baseCall}");
            c.Line("{");
            c.Indent();
            foreach (var f in own)
            {
                c.Line($"{Pascal(f.Name)} = {Camel(f.Name)};");
            }
            c.Outdent();
            c.Line("}");
            c.Line();
        }

        // ── StaticParse ─────────────────────────────────────────

        private void EmitStaticParse(
            CodeWriter c, ComplexTypeDefinition type, CSharpExpressionRenderer r, FieldScope scope)
        {
            var args = ParseSignatureArgs(type);
            // `new` only when the child's StaticParse has the same signature as
            // the base's (no extra context parameters) - otherwise it neither
            // hides anything (CS0109) nor needs to.
            var hide = type.IsDiscriminatedChild && ParentContextFields(type).Count == 0
                ? "new "
                : "";
            c.Line($"public static {hide}{type.Name} StaticParse(ReadBuffer readBuffer{Comma(args)})");
            c.Line("{");
            c.Indent();

            // `curPos` is "bytes since this type's parse began". A child that
            // needs it receives _startPos from the parent (the byte count is
            // measured from where the parent started, before the discriminators);
            // everyone else captures it here.
            if (NeedsStartPos(type) && !(type.IsDiscriminatedChild && UsesCurPos(type)))
            {
                c.Line("var _startPos = readBuffer.GetPos();");
            }

            foreach (var f in FieldsToEmit(type))
            {
                EmitFieldParse(c, type, f, r, scope);
            }

            if (!type.IsDiscriminatedParent)
            {
                // prefix (passed in), own + suffix (read above), in wire order.
                var ctorArgs = string.Join(", ", AllValueFields(type).Select(f => Camel(f.Name)));
                c.Line($"return new {type.Name}({ctorArgs});");
            }

            c.Outdent();
            c.Line("}");
            c.Line();
        }

        private void EmitFieldParse(
            CodeWriter c, ComplexTypeDefinition type, Field f, CSharpExpressionRenderer r, FieldScope scope)
        {
            switch (f)
            {
                case DiscriminatorField d:
                    c.Line($"var {Camel(d.Name)} = {CSharpTypeMapper.ReadCall(d.Type, d.Name)};");
                    break;

                case SimpleField:
                    c.Line($"var {Camel(f.Name)} = {ReadFieldValue(f, r)};");
                    break;

                case EnumField ef:
                    c.Line($"var {Camel(ef.Name)} = {EnumFieldRead(ef, r)};");
                    break;

                case ConstField cf:
                    c.Line($"var {Camel(cf.Name)} = {CSharpTypeMapper.ReadCall(cf.Type, cf.Name)};");
                    c.Line($"if (!Equals({Camel(cf.Name)}, {WithType(r.Render(cf.ReferenceValue), cf.Type)}))");
                    c.Indent();
                    c.Line($"throw new ParseException($\"Expected constant {{{Pascal(cf.Name)}}} for '{cf.Name}' but got {{{Camel(cf.Name)}}}\");");
                    c.Outdent();
                    break;

                case ImplicitField imf:
                    c.Line($"var {Camel(imf.Name)} = {CSharpTypeMapper.ReadCall(imf.Type, imf.Name)};");
                    break;

                case ReservedField rf:
                    c.Line("{");
                    c.Indent();
                    c.Line($"var reserved = {CSharpTypeMapper.ReadCall(rf.Type, "reserved")};");
                    c.Line($"if (!Equals(reserved, {WithType(r.Render(rf.ReferenceValue), rf.Type)})) {{ /* mspec reserved: value differs from the spec default */ }}");
                    c.Outdent();
                    c.Line("}");
                    break;

                case ArrayField af:
                    EmitArrayParse(c, af, r, scope);
                    break;

                case ChecksumField ck:
                    c.Line($"var {Camel(ck.Name)} = {CSharpTypeMapper.ReadCall(ck.Type, ck.Name)};");
                    break;

                case OptionalField opt:
                    c.Line($"{PropertyType(opt)} {Camel(opt.Name)} = null;");
                    if (opt.Condition != null)
                    {
                        c.Line($"if ({r.Render(opt.Condition)})");
                        c.Line("{");
                        c.Indent();
                    }
                    c.Line($"{Camel(opt.Name)} = {ReadFieldValue(opt, r)};");
                    if (opt.Condition != null)
                    {
                        c.Outdent();
                        c.Line("}");
                    }
                    break;

                case ManualField mf:
                    c.Line($"var {Camel(mf.Name)} = ({CSharpTypeMapper.CSharpType(mf.Type)}) ({r.Render(mf.ParseExpression)});");
                    break;

                case PaddingField pad:
                    c.Line("{");
                    c.Indent();
                    c.Line($"var _timesPadding = (int) ({r.Render(pad.TimesPadding)});");
                    c.Line("while (_timesPadding-- > 0)");
                    c.Line("{");
                    c.Indent();
                    c.Line($"{CSharpTypeMapper.ReadCall(pad.Type, "padding")};");
                    c.Outdent();
                    c.Line("}");
                    c.Outdent();
                    c.Line("}");
                    break;

                case TypeSwitchField ts:
                    EmitTypeSwitchParse(c, type, ts, r);
                    break;

                case VirtualField:
                    break;

                case UnsupportedField uf:
                    c.Line($"// TODO unsupported mspec field '{uf.MspecKeyword}': {uf.RawText}");
                    break;
            }
        }

        private void EmitArrayParse(CodeWriter c, ArrayField af, CSharpExpressionRenderer r, FieldScope scope)
        {
            var element = ElementType(af);
            var (elemType, isByteArray) = CSharpTypeMapper.ArrayType(element);
            var name = Camel(af.Name);

            if (isByteArray && af.LoopType == ArrayField.Loop.Count)
            {
                c.Line($"var {name} = readBuffer.ReadByteArray(\"{af.Name}\", (int) ({r.Render(af.LoopExpression)}) * 8);");
                return;
            }

            c.Line($"var {name} = new {elemType}();");
            // Underscore-prefixed locals: mspec field names never start with one,
            // so a loop's bookkeeping cannot collide with an implicit / simple
            // field of the same base name (e.g. an `itemsCount` implicit beside
            // an `items count 'itemsCount'` array).
            // `_lastItem` is only knowable ahead of time for a counted loop.
            var lastItem = af.LoopType == ArrayField.Loop.Count && ElementUsesLastItem(af)
                ? $"_{name}I == _{name}Cnt - 1"
                : null;
            if (af.LoopType == ArrayField.Loop.Length)
            {
                c.Line($"var _{name}End = readBuffer.GetPos() + (int) ({r.Render(af.LoopExpression)}) * 8;");
                c.Line($"while (readBuffer.GetPos() < _{name}End)");
            }
            else
            {
                c.Line($"var _{name}Cnt = (int) ({r.Render(af.LoopExpression)});");
                c.Line($"for (var _{name}I = 0; _{name}I < _{name}Cnt; _{name}I++)");
            }
            c.Line("{");
            c.Indent();
            var read = CSharpTypeMapper.ReadCall(element, af.Name, ElementArgs(element, r, lastItem));
            c.Line($"{name}.Add({read});");
            c.Outdent();
            c.Line("}");
        }

        private void EmitTypeSwitchParse(
            CodeWriter c, ComplexTypeDefinition parent, TypeSwitchField ts, CSharpExpressionRenderer r)
        {
            var discriminators = DiscriminatorList(parent, r);

            foreach (var childName in ts.CaseNames)
            {
                var child = _protocol.FindType(childName);
                var values = child.DiscriminatorValues ?? Array.Empty<Term>();
                var conds = new List<string>();
                for (var i = 0; i < discriminators.Count && i < values.Count; i++)
                {
                    if (values[i] == null)
                    {
                        continue;
                    }
                    var d = discriminators[i];
                    conds.Add($"Equals({d.DispatchExpr}, {RenderDiscriminatorValue(values[i], d.Type, r)})");
                }
                var test = conds.Count == 0 ? "true" : string.Join(" && ", conds);
                var childArgs = ChildParseArgNames(child)
                    .Concat(ParentContextFields(child).Select(f => Camel(f.Name)));
                if (UsesCurPos(child))
                {
                    childArgs = childArgs.Append("_startPos");
                }
                c.Line($"if ({test})");
                c.Line("{");
                c.Indent();
                c.Line($"return {childName}.StaticParse(readBuffer{Comma(childArgs)});");
                c.Outdent();
                c.Line("}");
            }
            c.Line($"throw new ParseException(\"No matching subtype found for {parent.Name}\");");
        }

        // ── Serialize ───────────────────────────────────────────

        private void EmitSerialize(
            CodeWriter c, ComplexTypeDefinition type, CSharpExpressionRenderer r, FieldScope scope)
        {
            if (type.IsDiscriminatedChild)
            {
                c.Line($"protected override void SerializeChild(WriteBuffer writeBuffer)");
            }
            else if (UsesLastItem(type))
            {
                c.Line("public void Serialize(WriteBuffer writeBuffer) => Serialize(writeBuffer, false);");
                c.Line();
                c.Line("public void Serialize(WriteBuffer writeBuffer, bool _lastItem)");
            }
            else
            {
                c.Line("public void Serialize(WriteBuffer writeBuffer)");
            }
            c.Line("{");
            c.Indent();

            var serScope = scope.ForSerialize();
            var serRender = new CSharpExpressionRenderer(serScope) { StaticHelperClass = r.StaticHelperClass };

            HoistParentImplicits(c, type, serRender, serScope);

            foreach (var f in FieldsToEmit(type))
            {
                // `_value` inside a manual field's serialize expression is that
                // field's own value.
                serScope.CurrentValue = f is ManualField ? Pascal(f.Name) : null;
                EmitFieldSerialize(c, type, f, serRender);
            }
            serScope.CurrentValue = null;
            serScope.ClearLocals();
            if (type.IsDiscriminatedParent)
            {
                c.Line("SerializeChild(writeBuffer);");
            }

            c.Outdent();
            c.Line("}");
            c.Line();

            if (type.IsDiscriminatedParent)
            {
                c.Line("protected abstract void SerializeChild(WriteBuffer writeBuffer);");
                c.Line();
            }
        }

        private void EmitFieldSerialize(
            CodeWriter c, ComplexTypeDefinition type, Field f, CSharpExpressionRenderer r)
        {
            switch (f)
            {
                case DiscriminatorField d:
                    c.Line($"{CSharpTypeMapper.WriteCall(d.Type, d.Name, Pascal(d.Name))};");
                    break;

                case SimpleField:
                    c.Line($"{CSharpTypeMapper.WriteCall(f.Type, f.Name, Pascal(f.Name))};");
                    break;

                case EnumField ef:
                    c.Line($"{EnumFieldWrite(ef)};");
                    break;

                case ConstField cf:
                    c.Line($"{CSharpTypeMapper.WriteCall(cf.Type, cf.Name, Pascal(cf.Name))};");
                    break;

                case ImplicitField imf:
                    var val = $"({CSharpTypeMapper.CSharpType(imf.Type)}) ({r.Render(imf.SerializeExpression)})";
                    c.Line($"{CSharpTypeMapper.WriteCall(imf.Type, imf.Name, val)};");
                    break;

                case ReservedField rf:
                    c.Line($"{CSharpTypeMapper.WriteCall(rf.Type, "reserved", WithType(r.Render(rf.ReferenceValue), rf.Type))};");
                    break;

                case ArrayField af:
                    EmitArraySerialize(c, af, r);
                    break;

                case ChecksumField ck:
                    var cv = $"({CSharpTypeMapper.CSharpType(ck.Type)}) ({r.Render(ck.ChecksumExpression)})";
                    c.Line($"{CSharpTypeMapper.WriteCall(ck.Type, ck.Name, cv)};");
                    break;

                case OptionalField opt:
                    // .Value for a nullable value type; the reference itself
                    // for a nullable class.
                    var optRef = opt.Type is ComplexTypeReference
                        ? Pascal(opt.Name)
                        : $"{Pascal(opt.Name)}.Value";
                    c.Line($"if ({Pascal(opt.Name)} != null)");
                    c.Line("{");
                    c.Indent();
                    c.Line($"{CSharpTypeMapper.WriteCall(opt.Type, opt.Name, optRef)};");
                    c.Outdent();
                    c.Line("}");
                    break;

                case ManualField mf:
                    c.Line($"{r.Render(mf.SerializeExpression)};");
                    break;

                case PaddingField pad:
                    c.Line("{");
                    c.Indent();
                    c.Line($"var _timesPadding = (int) ({r.Render(pad.TimesPadding)});");
                    c.Line("while (_timesPadding-- > 0)");
                    c.Line("{");
                    c.Indent();
                    c.Line($"{CSharpTypeMapper.WriteCall(pad.Type, "padding", WithType(r.Render(pad.PaddingValue), pad.Type))};");
                    c.Outdent();
                    c.Line("}");
                    c.Outdent();
                    c.Line("}");
                    break;

                case TypeSwitchField:
                case VirtualField:
                case UnsupportedField:
                    break;
            }
        }

        private void EmitArraySerialize(CodeWriter c, ArrayField af, CSharpExpressionRenderer r)
        {
            var element = ElementType(af);
            var (_, isByteArray) = CSharpTypeMapper.ArrayType(element);
            var prop = Pascal(af.Name);

            if (isByteArray)
            {
                c.Line($"writeBuffer.WriteByteArray(\"{af.Name}\", {prop});");
                return;
            }

            if (ElementUsesLastItem(af))
            {
                c.Line($"for (var _i = 0; _i < {prop}.Count; _i++)");
                c.Line("{");
                c.Indent();
                c.Line($"{prop}[_i].Serialize(writeBuffer, _i == {prop}.Count - 1);");
                c.Outdent();
                c.Line("}");
                return;
            }

            c.Line($"foreach (var _e in {prop})");
            c.Line("{");
            c.Indent();
            c.Line($"{CSharpTypeMapper.WriteCall(element, af.Name, "_e")};");
            c.Outdent();
            c.Line("}");
        }

        // ── length ──────────────────────────────────────────────

        private void EmitLength(
            CodeWriter c, ComplexTypeDefinition type, CSharpExpressionRenderer r, FieldScope scope)
        {
            var serScope = scope.ForSerialize();
            var serRender = new CSharpExpressionRenderer(serScope) { StaticHelperClass = r.StaticHelperClass };

            var childMethod = type.IsDiscriminatedChild;
            var lastItem = !childMethod && UsesLastItem(type);
            if (lastItem)
            {
                c.Line("public int GetLengthInBits() => GetLengthInBits(false);");
                c.Line();
                c.Line("public int GetLengthInBits(bool _lastItem)");
            }
            else
            {
                c.Line(childMethod
                    ? "protected override int GetLengthInBitsChild()"
                    : "public int GetLengthInBits()");
            }
            c.Line("{");
            c.Indent();
            c.Line("var lengthInBits = 0;");
            HoistParentImplicits(c, type, serRender, serScope);
            foreach (var f in FieldsToEmit(type))
            {
                var term = LengthContribution(f, serRender);
                if (term != null)
                {
                    c.Line($"lengthInBits += {term};");
                }
            }
            serScope.ClearLocals();
            if (type.IsDiscriminatedParent)
            {
                c.Line("lengthInBits += GetLengthInBitsChild();");
            }
            c.Line("return lengthInBits;");
            c.Outdent();
            c.Line("}");
            c.Line();

            if (type.IsDiscriminatedParent)
            {
                c.Line("protected abstract int GetLengthInBitsChild();");
                c.Line();
            }
            if (!childMethod)
            {
                c.Line($"public int GetLengthInBytes() => GetLengthInBits({(lastItem ? "false" : "")}) / 8;");
            }
        }

        /// <summary>
        /// A discriminated child's serialize / length can reference the
        /// parent's pre-typeSwitch implicit fields, which the parent (not the
        /// child) writes. Recompute the ones it names as locals at the top so
        /// the expressions resolve.
        /// </summary>
        private void HoistParentImplicits(
            CodeWriter c, ComplexTypeDefinition type, CSharpExpressionRenderer r, FieldScope scope)
        {
            if (!type.IsDiscriminatedChild || _protocol.FindType(type.ParentName) is not { } p)
            {
                return;
            }
            var referenced = new HashSet<string>();
            foreach (var f in FieldsToEmit(type))
            {
                foreach (var t in SerializeTerms(f))
                {
                    CollectNames(t, referenced);
                }
            }
            foreach (var imf in p.PrefixImplicitFields.Where(i => referenced.Contains(i.Name)))
            {
                scope.MarkLocal(imf.Name);
                c.Line($"var {Camel(imf.Name)} = ({CSharpTypeMapper.CSharpType(imf.Type)}) ({r.Render(imf.SerializeExpression)});");
            }
        }

        /// <summary>The expression trees a field contributes to
        /// <c>Serialize</c> / <c>GetLengthInBits</c>. Excludes the parse-only
        /// ones: an optional's condition, an array's loop expression, a manual
        /// field's parse expression, and a complex field's constructor
        /// arguments (those are passed to <c>StaticParse</c>, never re-evaluated
        /// on the way out).</summary>
        private static IEnumerable<Term> SerializeTerms(Field f)
        {
            switch (f)
            {
                case ImplicitField i: yield return i.SerializeExpression; break;
                case ChecksumField ck: yield return ck.ChecksumExpression; break;
                case ReservedField r: yield return r.ReferenceValue; break;
                case VirtualField v: yield return v.ValueExpression; break;
                case ManualField m:
                    yield return m.SerializeExpression;
                    yield return m.LengthExpression;
                    break;
                case PaddingField p:
                    yield return p.PaddingValue;
                    yield return p.TimesPadding;
                    break;
            }
        }

        private string LengthContribution(Field f, CSharpExpressionRenderer r)
        {
            switch (f)
            {
                case SimpleField or ConstField or DiscriminatorField or ImplicitField
                    or ReservedField or ChecksumField or EnumField:
                    var bits = CSharpTypeMapper.FixedBitLength(f.Type);
                    if (bits >= 0)
                    {
                        return bits.ToString();
                    }
                    if (f.Type is ComplexTypeReference)
                    {
                        return $"{Pascal(f.Name)}.GetLengthInBits()";
                    }
                    return null; // vstring / unmodelled - length not computed

                case VirtualField:
                    return null;

                case ArrayField af:
                    var element = ElementType(af);
                    var (_, isByteArray) = CSharpTypeMapper.ArrayType(element);
                    var prop = Pascal(af.Name);
                    if (isByteArray)
                    {
                        return $"({prop}.Length * 8)";
                    }
                    var eb = CSharpTypeMapper.FixedBitLength(element);
                    if (eb >= 0)
                    {
                        return $"({prop}.Count * {eb})";
                    }
                    return ElementUsesLastItem(af)
                        ? $"System.Linq.Enumerable.Range(0, {prop}.Count).Sum(_i => {prop}[_i].GetLengthInBits(_i == {prop}.Count - 1))"
                        : $"{prop}.Sum(_e => _e.GetLengthInBits())";

                case OptionalField opt:
                    var name = Pascal(opt.Name);
                    var ob = CSharpTypeMapper.FixedBitLength(opt.Type);
                    if (ob >= 0)
                    {
                        return $"({name} != null ? {ob} : 0)";
                    }
                    return opt.Type is ComplexTypeReference
                        ? $"({name}?.GetLengthInBits() ?? 0)"
                        : null;

                case ManualField mf:
                    return $"(({r.Render(mf.LengthExpression)}) * 8)";

                case PaddingField pad:
                    var pb = CSharpTypeMapper.FixedBitLength(pad.Type);
                    return $"(({r.Render(pad.TimesPadding)}) * {(pb > 0 ? pb : 8)})";

                default:
                    return null;
            }
        }

        // ── enums ───────────────────────────────────────────────

        private string EmitEnum(EnumTypeDefinition e)
        {
            var r = new CSharpExpressionRenderer(FieldScope.EnumConstantsOnly(_protocol));
            var baseType = e.BaseType != null
                ? CSharpTypeMapper.SimpleCSharpType(e.BaseType)
                : "int";

            var c = new CodeWriter();
            Header(c);
            c.Line($"namespace {_namespace}.model");
            c.Line("{");
            c.Indent();
            c.Line($"public enum {e.Name} : {baseType}");
            c.Line("{");
            c.Indent();
            foreach (var v in e.Values)
            {
                var val = v.Value != null ? r.Render(v.Value) : null;
                c.Line(val != null ? $"{v.Name} = {val}," : $"{v.Name},");
            }
            c.Outdent();
            c.Line("}");

            if (e.Arguments.Count > 0)
            {
                c.Line();
                c.Line($"public static class {e.Name}Extensions");
                c.Line("{");
                c.Indent();
                for (var ai = 0; ai < e.Arguments.Count; ai++)
                {
                    var arg = e.Arguments[ai];
                    var rows = e.Values.Where(v => v.ConstantValues.Count > ai).ToList();
                    var nullable = IsValueType(arg.Type)
                        && rows.Any(v => v.ConstantValues[ai] is NullLiteral);
                    var retType = CSharpTypeMapper.CSharpType(arg.Type) + (nullable ? "?" : "");
                    // `Get` prefix: keeps the method name off the enum type name
                    // it may return (`DataTransportSize`), and matches the
                    // `x.attr` -> `x.GetAttr()` the expression renderer emits.
                    c.Line($"public static {retType} Get{Pascal(arg.Name)}(this {e.Name} value) => value switch");
                    c.Line("{");
                    c.Indent();
                    // dedup on the enum's underlying value: S7's TransportSize
                    // has COUNTER and DATE_AND_TIME both at 0x1C.
                    EmitEnumArms(c, rows, NumericKey, v => $"{e.Name}.{v.Name}",
                        v => RenderEnumParamValue(v.ConstantValues[ai], arg.Type, r));
                    c.Line("_ => default,");
                    c.Outdent();
                    c.Line("};");
                }
                foreach (var key in EnumFieldKeys(e.Name))
                {
                    var ki = e.Arguments.FindIndex(a => a.Name == key);
                    if (ki < 0)
                    {
                        continue;
                    }
                    var keyType = e.Arguments[ki].Type;
                    var rows = e.Values.Where(v => v.ConstantValues.Count > ki).ToList();
                    c.Line($"public static {e.Name} FirstEnumForField{Pascal(key)}({CSharpTypeMapper.CSharpType(keyType)} {Camel(key)}) => {Camel(key)} switch");
                    c.Line("{");
                    c.Indent();
                    // "first" lookup: the earliest constant wins each key value.
                    string PatternFor(EnumValue v) => RenderEnumParamValue(v.ConstantValues[ki], keyType, r);
                    EmitEnumArms(c, rows, PatternFor, PatternFor, v => $"{e.Name}.{v.Name}");
                    c.Line("_ => default,");
                    c.Outdent();
                    c.Line("};");
                }
                c.Outdent();
                c.Line("}");
            }

            c.Outdent();
            c.Line("}");
            return c.ToString();
        }

        /// <summary>Emits <c>pattern =&gt; result,</c> switch arms, skipping any
        /// whose <paramref name="dedupKey"/> repeats an earlier arm - two enum
        /// constants can share an underlying value, and an enum can map several
        /// constants to the same attribute value.</summary>
        private static void EmitEnumArms(
            CodeWriter c, IEnumerable<EnumValue> rows,
            System.Func<EnumValue, string> dedupKey,
            System.Func<EnumValue, string> pattern,
            System.Func<EnumValue, string> result)
        {
            var seen = new HashSet<string>();
            foreach (var v in rows)
            {
                if (seen.Add(dedupKey(v)))
                {
                    c.Line($"{pattern(v)} => {result(v)},");
                }
            }
        }

        private static string NumericKey(EnumValue v) =>
            v.Value is IntegerLiteral i ? i.Value.ToString() : v.Value?.ToString() ?? v.Name;

        /// <summary>Renders one entry of an enum's parameter table, coerced to
        /// the parameter's declared type: a bare word is a string literal for a
        /// <c>vstring</c> parameter, a member of the target enum for an
        /// enum-typed one, and a char literal for a one-character integer.</summary>
        private string RenderEnumParamValue(Term term, TypeReference target, CSharpExpressionRenderer r)
        {
            switch (target)
            {
                case EnumTypeReference et:
                    if (term is NullLiteral)
                    {
                        return "null";
                    }
                    if (term is VariableLiteral { Child: null, IsCall: false } ev)
                    {
                        return $"{et.Name}.{ev.Name}";
                    }
                    return r.Render(term);

                case SimpleTypeReference { BaseType: SimpleTypeReference.Base.String
                    or SimpleTypeReference.Base.VString }:
                    if (term is NullLiteral)
                    {
                        return "null";
                    }
                    if (term is VariableLiteral { Child: null, IsCall: false } sv)
                    {
                        return "\"" + sv.Name + "\"";
                    }
                    return r.Render(term);

                case SimpleTypeReference { IsIntegerLike: true } it:
                    if (term is VariableLiteral { Child: null, IsCall: false, Name.Length: 1 } cv)
                    {
                        return $"({CSharpTypeMapper.SimpleCSharpType(it)}) '{cv.Name}'";
                    }
                    return r.Render(term);

                default:
                    return r.Render(term);
            }
        }

        private static bool IsValueType(TypeReference t) =>
            t is EnumTypeReference
            || t is SimpleTypeReference s
               && s.BaseType is not (SimpleTypeReference.Base.String or SimpleTypeReference.Base.VString);

        /// <summary>The distinct <c>KeyAccessor</c>s of every <c>enum</c> field
        /// in the protocol that is typed as <paramref name="enumName"/> - each
        /// needs a <c>FirstEnumForField…</c> reverse lookup.</summary>
        private IEnumerable<string> EnumFieldKeys(string enumName) =>
            _protocol.Types
                .SelectMany(t => t.Fields)
                .OfType<EnumField>()
                .Where(f => f.KeyAccessor != null
                            && (f.Type as EnumTypeReference)?.Name == enumName)
                .Select(f => f.KeyAccessor)
                .Distinct();

        private string EmitConstants()
        {
            var r = new CSharpExpressionRenderer(FieldScope.EnumConstantsOnly(_protocol));
            var c = new CodeWriter();
            Header(c);
            c.Line($"namespace {_namespace}.model");
            c.Line("{");
            c.Indent();
            c.Line($"public static class {Pascal(_protocolName)}Constants");
            c.Line("{");
            c.Indent();
            foreach (var k in _protocol.Constants)
            {
                c.Line($"public const {CSharpTypeMapper.CSharpType(k.Type)} {Pascal(k.Name)} = {r.Render(k.Value)};");
            }
            c.Outdent();
            c.Line("}");
            c.Outdent();
            c.Line("}");
            return c.ToString();
        }

        // ── dataIo (parses to / from an IPlcValue) ──────────────

        /// <summary>IEC value wrappers whose plc4net type is exactly
        /// <c>Plc&lt;caseName&gt;</c> and which round-trip through a single
        /// scalar accessor.</summary>
        private static readonly HashSet<string> ScalarPlcValueCases = new HashSet<string>
        {
            "BOOL", "BYTE", "WORD", "DWORD", "LWORD",
            "SINT", "USINT", "INT", "UINT", "DINT", "UDINT", "LINT", "ULINT",
            "REAL", "LREAL", "CHAR", "WCHAR", "STRING", "WSTRING",
        };

        /// <summary>The TIA date / time <c>DataItem</c> cases - each maps to a
        /// <c>Plc{caseName}</c> built through a factory / segment call, and the
        /// S7 date helpers (<c>parseS5Time</c>, <c>parseTiaDate</c>,
        /// <c>parseSiemensYear</c>, BCD). Only the multi-value
        /// <c>numberOfValues &gt; 1</c> list cases stay a throwing stub now.</summary>
        private static readonly HashSet<string> TemporalPlcValueCases = new HashSet<string>
        {
            "TIME", "LTIME", "DATE", "TIME_OF_DAY", "LTIME_OF_DAY",
            "DATE_AND_TIME", "DATE_AND_LTIME",
        };

        private string EmitDataIo(DataIoTypeDefinition dio)
        {
            var helper = $"{Pascal(_protocolName)}StaticHelper";

            var c = new CodeWriter();
            Header(c);
            c.Line("using System;");
            c.Line("using org.apache.plc4net.api.value;");
            c.Line("using org.apache.plc4net.spi.generation;");
            c.Line("using org.apache.plc4net.spi.model.values;");
            c.Line();
            c.Line($"namespace {_namespace}.model");
            c.Line("{");
            c.Indent();
            c.Line("/// <summary>");
            c.Line($"/// mspec <c>[dataIo {dio.Name}]</c> - reads and writes one");
            c.Line("/// <see cref=\"IPlcValue\"/> whose wire layout the parser arguments pick.");
            c.Line("/// </summary>");
            c.Line($"public static class {dio.Name}");
            c.Line("{");
            c.Indent();

            var argDecls = dio.Arguments
                .Select(a => $"{CSharpTypeMapper.CSharpType(a.Type)} {Camel(a.Name)}")
                .ToList();
            var argNames = dio.Arguments.Select(a => Camel(a.Name)).ToList();

            c.Line($"public static IPlcValue StaticParse(ReadBuffer readBuffer{Comma(argDecls)})");
            c.Line("{");
            c.Indent();
            for (var i = 0; i < dio.Cases.Count; i++)
            {
                EmitDataIoCase(c, dio, dio.Cases[i], helper, DataIoMode.Parse, i == 0);
            }
            c.Line("return new PlcNULL();");
            c.Outdent();
            c.Line("}");
            c.Line();

            c.Line($"public static void StaticSerialize(WriteBuffer writeBuffer, IPlcValue _value{Comma(argDecls)})");
            c.Line("{");
            c.Indent();
            for (var i = 0; i < dio.Cases.Count; i++)
            {
                EmitDataIoCase(c, dio, dio.Cases[i], helper, DataIoMode.Serialize, i == 0);
            }
            c.Outdent();
            c.Line("}");
            c.Line();

            c.Line($"public static int GetLengthInBytes(IPlcValue _value{Comma(argDecls)}) =>");
            c.Indent();
            c.Line($"(GetLengthInBits(_value{Comma(argNames)}) + 7) / 8;");
            c.Outdent();
            c.Line();
            c.Line($"public static int GetLengthInBits(IPlcValue _value{Comma(argDecls)})");
            c.Line("{");
            c.Indent();
            c.Line("var lengthInBits = 0;");
            for (var i = 0; i < dio.Cases.Count; i++)
            {
                EmitDataIoCase(c, dio, dio.Cases[i], helper, DataIoMode.Length, i == 0);
            }
            c.Line("return lengthInBits;");
            c.Outdent();
            c.Line("}");

            c.Outdent();
            c.Line("}");
            c.Outdent();
            c.Line("}");
            return c.ToString();
        }

        private enum DataIoMode { Parse, Serialize, Length }

        private void EmitDataIoCase(
            CodeWriter c, DataIoTypeDefinition dio, ComplexTypeDefinition cs, string helper, DataIoMode mode, bool first)
        {
            var serialize = mode != DataIoMode.Parse;
            var synthetic = new ComplexTypeDefinition { Name = dio.Name };
            synthetic.Arguments.AddRange(dio.Arguments);
            synthetic.Fields.AddRange(cs.Fields);
            var scope = new FieldScope(_protocol, synthetic, serialize);
            // A dataIo's parser arguments are method parameters in every method,
            // never properties - keep them camelCase even while serializing.
            foreach (var a in dio.Arguments)
            {
                scope.MarkLocal(a.Name);
            }
            var r = new CSharpExpressionRenderer(scope) { StaticHelperClass = helper };

            c.Line($"{(first ? "if" : "else if")} ({DataIoCaseTest(dio, cs, r)})");
            c.Line("{");
            c.Indent();
            c.Line($"// {cs.Name}");

            var props = cs.PropertyFields.ToList();
            var valueField = props.Count == 1 ? props[0] : null;
            var listField = valueField as ArrayField;
            var temporal = TemporalPlcValueCases.Contains(cs.Name)
                           && props.All(f => f is not ArrayField);
            var scalar = valueField is not (null or ArrayField) && ScalarPlcValueCases.Contains(cs.Name);

            if (!temporal && !scalar
                && !(listField?.Type is SimpleTypeReference && listField.LoopType == ArrayField.Loop.Count))
            {
                c.Line($"throw new NotImplementedException(\"{dio.Name} '{cs.Name}' is not a shape the generator emits yet (design.md GAP-8)\");");
                c.Outdent();
                c.Line("}");
                return;
            }

            switch (mode)
            {
                case DataIoMode.Parse when temporal:
                    EmitDataIoTemporalParse(c, cs, r);
                    break;
                case DataIoMode.Parse when listField != null:
                    EmitDataIoListParse(c, listField, r);
                    break;
                case DataIoMode.Parse:
                    foreach (var f in cs.Fields)
                    {
                        EmitDataIoFieldParse(c, f, r);
                    }
                    c.Line($"return {DataIoReturnValue(cs, valueField)};");
                    break;

                case DataIoMode.Serialize when temporal:
                    EmitDataIoTemporalSerialize(c, cs, helper);
                    break;
                case DataIoMode.Serialize when listField != null:
                    EmitDataIoListSerialize(c, listField);
                    break;
                case DataIoMode.Serialize:
                    foreach (var f in cs.Fields)
                    {
                        EmitDataIoFieldSerialize(c, cs, f, valueField, r);
                    }
                    break;

                case DataIoMode.Length when listField != null:
                    c.Line($"lengthInBits += _value.GetLength() * {CSharpTypeMapper.FixedBitLength(listField.Type)};");
                    break;
                case DataIoMode.Length:
                    foreach (var f in cs.Fields)
                    {
                        var bits = DataIoFieldLength(f, r);
                        if (bits != null)
                        {
                            c.Line($"lengthInBits += {bits};");
                        }
                    }
                    break;
            }

            c.Outdent();
            c.Line("}");
        }

        /// <summary>Parses one TIA date / time case to its <c>Plc…</c> wrapper.
        /// Mirrors plc4j's <c>DataItem</c>: a factory / segment call, with the S7
        /// helpers doing S5TIME, the Siemens epoch and BCD.</summary>
        private void EmitDataIoTemporalParse(CodeWriter c, ComplexTypeDefinition cs, CSharpExpressionRenderer r)
        {
            var field = cs.PropertyFields.FirstOrDefault();
            switch (cs.Name)
            {
                case "TIME":
                    c.Line(field is ManualField mf
                        ? $"var milliseconds = {r.Render(RewriteTypeEncoding(mf.ParseExpression, "UTF8"))};"
                        : "var milliseconds = readBuffer.ReadUint(\"milliseconds\", 32);");
                    c.Line("return PlcTIME.OfMilliseconds(milliseconds);");
                    break;
                case "LTIME":
                    c.Line("var nanoseconds = readBuffer.ReadUlong(\"nanoseconds\", 64);");
                    c.Line("return PlcLTIME.OfNanoseconds(nanoseconds);");
                    break;
                case "DATE":
                    c.Line($"var daysSinceEpoch = {r.Render(RewriteTypeEncoding(((ManualField) field).ParseExpression, "UTF8"))};");
                    c.Line("return PlcDATE.OfDaysSinceEpoch(daysSinceEpoch);");
                    break;
                case "TIME_OF_DAY":
                    c.Line("var millisecondsSinceMidnight = readBuffer.ReadUint(\"millisecondsSinceMidnight\", 32);");
                    c.Line("return PlcTIME_OF_DAY.OfMillisecondsSinceMidnight(millisecondsSinceMidnight);");
                    break;
                case "LTIME_OF_DAY":
                    c.Line("var nanosecondsSinceMidnight = readBuffer.ReadUlong(\"nanosecondsSinceMidnight\", 64);");
                    c.Line("return PlcLTIME_OF_DAY.OfNanosecondsSinceMidnight(nanosecondsSinceMidnight);");
                    break;
                case "DATE_AND_TIME":
                    c.Line($"var year = {r.StaticHelperClass}.ParseSiemensYear(readBuffer);");
                    foreach (var f in new[] { "month", "day", "hour", "minutes", "seconds" })
                    {
                        c.Line($"var {f} = {r.StaticHelperClass}.BcdToBin(readBuffer.ReadByte(\"{f}\", 8));");
                    }
                    c.Line($"var millisecondsOfSecond = {r.StaticHelperClass}.BcdToBin12(readBuffer.ReadUshort(\"millisecondsOfSecond\", 12));");
                    c.Line("readBuffer.ReadByte(\"dayOfWeek\", 4);");
                    c.Line("return PlcDATE_AND_TIME.OfSegments(year, month, day, hour, minutes, seconds, millisecondsOfSecond * 1000000);");
                    break;
                case "DATE_AND_LTIME" when cs.PropertyFields.Count() == 1:
                    c.Line("var nanosecondsSinceEpoch = readBuffer.ReadUlong(\"nanosecondsSinceEpoch\", 64);");
                    c.Line("return PlcDATE_AND_LTIME.OfNanosecondsSinceEpoch(nanosecondsSinceEpoch);");
                    break;
                case "DATE_AND_LTIME": // DTL - plain-binary segments, widths from the mspec
                    foreach (var f in cs.Fields.OfType<SimpleField>())
                    {
                        var read = DataIoSimpleRead(f);
                        // dayOfWeek is recomputed from the DateTime on serialize.
                        c.Line(f.Name == "dayOfWeek" ? $"{read};" : $"var {Camel(f.Name)} = {read};");
                    }
                    c.Line("return PlcDATE_AND_LTIME.OfSegments(year, month, day, hour, minutes, seconds, nanosecondsOfSecond);");
                    break;
            }
        }

        private void EmitDataIoTemporalSerialize(CodeWriter c, ComplexTypeDefinition cs, string helper)
        {
            var field = cs.PropertyFields.FirstOrDefault();
            // S7 DT / DTL number the week Sunday=1 .. Saturday=7 (per s7.mspec's
            // own DTL comment and the Siemens DTL doc). This is a DELIBERATE
            // departure from plc4j, whose DataItem writes ISO-8601
            // Monday=1 .. Sunday=7 (getDayOfWeek().getValue()) - the two never
            // agree on any weekday. Both discard the field on parse (the
            // DateTime carries its own weekday), so it only shows up when a
            // device consumes the serialized frame, or when diffing plc4net vs
            // plc4j serializer output.
            const string dow = "(byte) ((int) _dt.DayOfWeek + 1)";
            switch (cs.Name)
            {
                case "TIME" when field is ManualField:
                    c.Line($"{helper}.SerializeS5Time(writeBuffer, _value);");
                    break;
                case "TIME":
                    c.Line("writeBuffer.WriteUint(\"milliseconds\", 32, (uint) _value.GetDuration().TotalMilliseconds);");
                    break;
                case "LTIME":
                    c.Line("writeBuffer.WriteUlong(\"nanoseconds\", 64, ((PlcLTIME) _value).GetNanoseconds());");
                    break;
                case "DATE":
                    c.Line($"{helper}.SerializeTiaDate(writeBuffer, _value);");
                    break;
                case "TIME_OF_DAY":
                    c.Line("writeBuffer.WriteUint(\"millisecondsSinceMidnight\", 32, (uint) _value.GetTime().ToTimeSpan().TotalMilliseconds);");
                    break;
                case "LTIME_OF_DAY":
                    c.Line("writeBuffer.WriteUlong(\"nanosecondsSinceMidnight\", 64, ((PlcLTIME_OF_DAY) _value).GetNanosecondsSinceMidnight());");
                    break;
                case "DATE_AND_TIME":
                    c.Line("var _dt = _value.GetDateTime();");
                    c.Line($"{helper}.SerializeSiemensYear(writeBuffer, _value);");
                    c.Line($"writeBuffer.WriteByte(\"month\", 8, {helper}.BinToBcd(_dt.Month));");
                    c.Line($"writeBuffer.WriteByte(\"day\", 8, {helper}.BinToBcd(_dt.Day));");
                    c.Line($"writeBuffer.WriteByte(\"hour\", 8, {helper}.BinToBcd(_dt.Hour));");
                    c.Line($"writeBuffer.WriteByte(\"minutes\", 8, {helper}.BinToBcd(_dt.Minute));");
                    c.Line($"writeBuffer.WriteByte(\"seconds\", 8, {helper}.BinToBcd(_dt.Second));");
                    c.Line($"writeBuffer.WriteUshort(\"millisecondsOfSecond\", 12, {helper}.BinToBcd12(_dt.Millisecond));");
                    c.Line($"writeBuffer.WriteByte(\"dayOfWeek\", 4, {dow});");
                    break;
                case "DATE_AND_LTIME" when cs.PropertyFields.Count() == 1:
                    c.Line("writeBuffer.WriteUlong(\"nanosecondsSinceEpoch\", 64, ((PlcDATE_AND_LTIME) _value).GetNanosecondsSinceEpoch());");
                    break;
                case "DATE_AND_LTIME":
                    c.Line("var _dtl = (PlcDATE_AND_LTIME) _value;");
                    c.Line("var _dt = _dtl.GetDateTime();");
                    c.Line("writeBuffer.WriteUshort(\"year\", 16, (ushort) _dt.Year);");
                    c.Line("writeBuffer.WriteByte(\"month\", 8, (byte) _dt.Month);");
                    c.Line("writeBuffer.WriteByte(\"day\", 8, (byte) _dt.Day);");
                    c.Line($"writeBuffer.WriteByte(\"dayOfWeek\", 8, {dow});");
                    c.Line("writeBuffer.WriteByte(\"hour\", 8, (byte) _dt.Hour);");
                    c.Line("writeBuffer.WriteByte(\"minutes\", 8, (byte) _dt.Minute);");
                    c.Line("writeBuffer.WriteByte(\"seconds\", 8, (byte) _dt.Second);");
                    c.Line("writeBuffer.WriteUint(\"nanosecondsOfSecond\", 32, _dtl.GetNanosecondsOfSecond());");
                    break;
            }
        }

        /// <summary>A <c>numberOfValues &gt; 1</c> dataIo case (Modbus's
        /// <c>DataItem</c> List cases) - a counted array of one primitive,
        /// wrapped element-by-element into a <see cref="PlcList"/>.</summary>
        private void EmitDataIoListParse(CodeWriter c, ArrayField af, CSharpExpressionRenderer r)
        {
            var element = (SimpleTypeReference) af.Type;
            var wrap = DataIoPlcWrapper(element);
            c.Line("var value = new System.Collections.Generic.List<IPlcValue>();");
            c.Line($"var _valueCnt = (int) ({r.Render(af.LoopExpression)});");
            c.Line("for (var _i = 0; _i < _valueCnt; _i++)");
            c.Line("{");
            c.Indent();
            if (element.BaseType is SimpleTypeReference.Base.String)
            {
                c.Line($"var _e = readBuffer.ReadString(\"value\", {element.SizeInBits}, {CSharpTypeMapper.EncodingExpr(element.Encoding)});");
                c.Line($"value.Add(new {wrap}(_e.Length > 0 ? _e[0] : '\\0'));");
            }
            else
            {
                c.Line($"value.Add(new {wrap}({CSharpTypeMapper.ReadCall(element, "value")}));");
            }

            c.Outdent();
            c.Line("}");
            c.Line("return new PlcList(value);");
        }

        private void EmitDataIoListSerialize(CodeWriter c, ArrayField af)
        {
            var element = (SimpleTypeReference) af.Type;
            c.Line("foreach (var _e in _value.GetList())");
            c.Line("{");
            c.Indent();
            var item = element.BaseType is SimpleTypeReference.Base.String
                ? "_e.GetString()"
                : $"({CSharpTypeMapper.SimpleCSharpType(element)}) _e.{PlcValueGetter(element)}";
            c.Line(element.BaseType is SimpleTypeReference.Base.String
                ? $"writeBuffer.WriteString(\"value\", {element.SizeInBits}, \"{element.Encoding ?? "UTF8"}\", {item});"
                : $"{CSharpTypeMapper.WriteCall(element, "value", item)};");
            c.Outdent();
            c.Line("}");
        }

        /// <summary>The <c>Plc…</c> wrapper for one element of a dataIo list,
        /// keyed on the wire type (<c>bit</c> → <c>PlcBOOL</c>, <c>int 16</c> →
        /// <c>PlcINT</c>, <c>string 16</c> → <c>PlcWCHAR</c>, …).</summary>
        private static string DataIoPlcWrapper(SimpleTypeReference s) => s.BaseType switch
        {
            SimpleTypeReference.Base.Bit => "PlcBOOL",
            SimpleTypeReference.Base.Byte => "PlcUSINT",
            SimpleTypeReference.Base.UInt => s.SizeInBits <= 8 ? "PlcUSINT"
                : s.SizeInBits <= 16 ? "PlcUINT"
                : s.SizeInBits <= 32 ? "PlcUDINT" : "PlcULINT",
            SimpleTypeReference.Base.Int => s.SizeInBits <= 8 ? "PlcSINT"
                : s.SizeInBits <= 16 ? "PlcINT"
                : s.SizeInBits <= 32 ? "PlcDINT" : "PlcLINT",
            SimpleTypeReference.Base.Float or SimpleTypeReference.Base.UFloat =>
                s.SizeInBits <= 32 ? "PlcREAL" : "PlcLREAL",
            SimpleTypeReference.Base.String or SimpleTypeReference.Base.VString =>
                s.SizeInBits <= 8 ? "PlcCHAR" : "PlcWCHAR",
            _ => "PlcNULL",
        };

        private string DataIoCaseTest(
            DataIoTypeDefinition dio, ComplexTypeDefinition cs, CSharpExpressionRenderer r)
        {
            var discs = dio.TypeSwitch.Discriminators.OfType<VariableLiteral>().ToList();
            var values = cs.DiscriminatorValues ?? Array.Empty<Term>();
            var conds = new List<string>();
            for (var i = 0; i < discs.Count && i < values.Count; i++)
            {
                if (values[i] == null)
                {
                    continue;
                }
                var argType = dio.Arguments.FirstOrDefault(a => a.Name == discs[i].Name)?.Type;
                conds.Add($"Equals({Camel(discs[i].Name)}, {RenderDiscriminatorValue(values[i], argType, r)})");
            }
            return conds.Count == 0 ? "true" : string.Join(" && ", conds);
        }

        private void EmitDataIoFieldParse(CodeWriter c, Field f, CSharpExpressionRenderer r)
        {
            switch (f)
            {
                case ReservedField rf:
                    c.Line("{");
                    c.Indent();
                    c.Line($"var reserved = {CSharpTypeMapper.ReadCall(rf.Type, "reserved")};");
                    c.Line($"if (!Equals(reserved, {WithType(r.Render(rf.ReferenceValue), rf.Type)})) {{ /* mspec reserved: value differs from the spec default */ }}");
                    c.Outdent();
                    c.Line("}");
                    break;

                case ManualField mf:
                    var parse = RewriteTypeEncoding(mf.ParseExpression, FieldEncoding(mf));
                    c.Line($"var {Camel(mf.Name)} = ({CSharpTypeMapper.CSharpType(mf.Type)}) ({r.Render(parse)});");
                    break;

                case SimpleField sf:
                    c.Line($"var {Camel(sf.Name)} = {DataIoSimpleRead(sf)};");
                    break;
            }
        }

        private void EmitDataIoFieldSerialize(
            CodeWriter c, ComplexTypeDefinition cs, Field f, Field valueField, CSharpExpressionRenderer r)
        {
            switch (f)
            {
                case ReservedField rf:
                    c.Line($"{CSharpTypeMapper.WriteCall(rf.Type, "reserved", WithType(r.Render(rf.ReferenceValue), rf.Type))};");
                    break;

                case ManualField mf:
                    c.Line($"{r.Render(RewriteTypeEncoding(mf.SerializeExpression, FieldEncoding(mf)))};");
                    break;

                case SimpleField sf when sf.Type is SimpleTypeReference st:
                    var getter = st.BaseType is SimpleTypeReference.Base.String or SimpleTypeReference.Base.VString
                        ? "_value.GetString()"
                        : $"({CSharpTypeMapper.SimpleCSharpType(st)}) _value.{PlcValueGetter(st)}";
                    if (st.BaseType is SimpleTypeReference.Base.String)
                    {
                        c.Line($"writeBuffer.WriteString(\"{sf.Name}\", {st.SizeInBits}, \"{FieldEncoding(sf)}\", {getter});");
                    }
                    else
                    {
                        c.Line($"{CSharpTypeMapper.WriteCall(sf.Type, sf.Name, getter)};");
                    }
                    break;
            }
        }

        private string DataIoFieldLength(Field f, CSharpExpressionRenderer r)
        {
            switch (f)
            {
                case ReservedField or SimpleField:
                    var bits = CSharpTypeMapper.FixedBitLength(f.Type);
                    return bits >= 0 ? bits.ToString() : null;
                case ManualField mf:
                    // A dataIo manual length expression is already in bits
                    // (plc4j's DataItem: `(stringLength * 8) + 16`).
                    return $"({r.Render(RewriteTypeEncoding(mf.LengthExpression, FieldEncoding(mf)))})";
                default:
                    return null;
            }
        }

        /// <summary>The <c>return new Plc…(value)</c> for a scalar
        /// <c>DataItem</c> case. The case name is the plc4net value type.</summary>
        private string DataIoReturnValue(ComplexTypeDefinition cs, Field valueField)
        {
            var v = Camel(valueField.Name);
            switch (cs.Name)
            {
                case "CHAR":
                case "WCHAR":
                    // read as a 1- or 2-byte string; Plc(W)CHAR wraps a char.
                    return $"new Plc{cs.Name}({v}.Length > 0 ? {v}[0] : '\\0')";
                default:
                    // BOOL / BYTE / … / STRING / WSTRING - Plc<caseName>(value).
                    return $"new Plc{cs.Name}({v})";
            }
        }

        private string DataIoSimpleRead(SimpleField sf)
        {
            if (sf.Type is SimpleTypeReference { BaseType: SimpleTypeReference.Base.String } s)
            {
                return $"readBuffer.ReadString(\"{sf.Name}\", {s.SizeInBits}, {CSharpTypeMapper.EncodingExpr(FieldEncoding(sf))})";
            }
            return CSharpTypeMapper.ReadCall(sf.Type, sf.Name);
        }

        /// <summary>The <c>IPlcValue</c> accessor that yields a value of this
        /// wire type, for a dataIo serialize.</summary>
        private static string PlcValueGetter(SimpleTypeReference s) => s.BaseType switch
        {
            SimpleTypeReference.Base.Bit => "GetBool()",
            SimpleTypeReference.Base.Byte => "GetByte()",
            SimpleTypeReference.Base.UInt => s.SizeInBits <= 8 ? "GetByte()"
                : s.SizeInBits <= 16 ? "GetUshort()"
                : s.SizeInBits <= 32 ? "GetUint()" : "GetUlong()",
            SimpleTypeReference.Base.Int => s.SizeInBits <= 8 ? "GetSbyte()"
                : s.SizeInBits <= 16 ? "GetShort()"
                : s.SizeInBits <= 32 ? "GetInt()" : "GetLong()",
            SimpleTypeReference.Base.Float or SimpleTypeReference.Base.UFloat =>
                s.SizeInBits <= 32 ? "GetFloat()" : "GetDouble()",
            _ => "GetString()",
        };

        /// <summary>The string encoding for a dataIo value field: the
        /// <c>encoding</c> / <c>stringEncoding</c> attribute (s7.mspec spells it
        /// the first way, modbus.mspec the second), or the type's own resolved
        /// encoding, defaulting to UTF-8.</summary>
        private static string FieldEncoding(Field f)
        {
            if (f.Type is SimpleTypeReference { Encoding: { } typeEncoding })
            {
                return typeEncoding;
            }
            foreach (var key in new[] { "encoding", "stringEncoding" })
            {
                if (f.Attributes.TryGetValue(key, out var t) && t is StringLiteral s)
                {
                    return s.Value;
                }
            }
            return "UTF8";
        }

        /// <summary>Inlines <c>_type.encoding</c> (a dataIo case's own encoding
        /// attribute) as a string literal, the way plc4j's <c>DataItem</c> does -
        /// a dataIo has no instance for <c>_type</c> to resolve against.</summary>
        private static Term RewriteTypeEncoding(Term t, string encoding)
        {
            switch (t)
            {
                case VariableLiteral v when v.Name == "_type"
                        && v.Child is { Name: "encoding", Child: null }:
                    return new StringLiteral(encoding);
                case VariableLiteral v:
                    return new VariableLiteral(
                        v.Name,
                        v.Args?.Select(a => RewriteTypeEncoding(a, encoding)).ToList(),
                        v.Index?.Select(i => RewriteTypeEncoding(i, encoding)).ToList(),
                        v.Child);
                case BinaryExpression b:
                    return new BinaryExpression(
                        RewriteTypeEncoding(b.Left, encoding), b.Operator, RewriteTypeEncoding(b.Right, encoding));
                case UnaryExpression u:
                    return new UnaryExpression(u.Operator, RewriteTypeEncoding(u.Operand, encoding));
                case TernaryExpression x:
                    return new TernaryExpression(
                        RewriteTypeEncoding(x.Condition, encoding),
                        RewriteTypeEncoding(x.WhenTrue, encoding),
                        RewriteTypeEncoding(x.WhenFalse, encoding));
                default:
                    return t;
            }
        }

        // ── static helper (STATIC_CALL targets + array sizing) ──

        private string EmitStaticHelper()
        {
            // name -> the checksum/field C# type the call has to return
            var staticCalls = new SortedDictionary<string, string>(StringComparer.Ordinal);
            CollectStaticCalls(staticCalls);

            var cls = $"{Pascal(_protocolName)}StaticHelper";
            var c = new CodeWriter();
            Header(c);
            c.Line("using System;");
            c.Line("using System.Collections.Generic;");
            c.Line("using System.Linq;");
            c.Line("using org.apache.plc4net.spi.drivers;");
            c.Line();
            c.Line($"namespace {_namespace}.model");
            c.Line("{");
            c.Indent();
            c.Line("/// <summary>");
            c.Line($"/// The array-size helper, and the <c>STATIC_CALL</c> targets in");
            c.Line($"/// {_protocolName}.mspec. The generated bodies throw; supply the real");
            c.Line("/// implementation in a sibling non-generated <c>partial</c> file.");
            c.Line("/// </summary>");
            c.Line($"public static partial class {cls}");
            c.Line("{");
            c.Indent();
            c.Line("/// <summary>Serialized byte length of a sequence of messages.</summary>");
            c.Line("public static int ArraySizeInBytes<T>(IEnumerable<T> items) where T : IMessage");
            c.Indent();
            c.Line("=> items?.Sum(i => i.GetLengthInBytes()) ?? 0;");
            c.Outdent();
            foreach (var (name, retType) in staticCalls)
            {
                c.Line();
                c.Line($"public static {retType} {name}(params object[] args)");
                c.Indent();
                c.Line($"=> throw new NotImplementedException(\"{_protocolName}.mspec STATIC_CALL '{name}' has no implementation yet\");");
                c.Outdent();
            }
            c.Outdent();
            c.Line("}");
            c.Outdent();
            c.Line("}");
            return c.ToString();
        }

        private void CollectStaticCalls(IDictionary<string, string> sink)
        {
            foreach (var type in _protocol.Types)
            {
                foreach (var f in type.Fields)
                {
                    var retType = f.Type != null ? CSharpTypeMapper.CSharpType(f.Type) : "object";
                    foreach (var term in FieldTerms(f))
                    {
                        WalkStaticCalls(term, retType, sink);
                    }
                }
            }
        }

        private static void WalkStaticCalls(Term t, string retType, IDictionary<string, string> sink)
        {
            switch (t)
            {
                case VariableLiteral { Name: "STATIC_CALL" } v
                    when v.Args?.Count >= 1 && v.Args[0] is StringLiteral s:
                    sink[Pascal(s.Value)] = retType;
                    break;
                case VariableLiteral v:
                    foreach (var a in v.Args ?? Enumerable.Empty<Term>())
                    {
                        WalkStaticCalls(a, retType, sink);
                    }
                    if (v.Child != null)
                    {
                        WalkStaticCalls(v.Child, retType, sink);
                    }
                    break;
                case BinaryExpression b:
                    WalkStaticCalls(b.Left, retType, sink);
                    WalkStaticCalls(b.Right, retType, sink);
                    break;
                case UnaryExpression u:
                    WalkStaticCalls(u.Operand, retType, sink);
                    break;
                case TernaryExpression tt:
                    WalkStaticCalls(tt.Condition, retType, sink);
                    WalkStaticCalls(tt.WhenTrue, retType, sink);
                    WalkStaticCalls(tt.WhenFalse, retType, sink);
                    break;
            }
        }

        private static IEnumerable<Term> FieldTerms(Field f)
        {
            switch (f)
            {
                case ImplicitField i: yield return i.SerializeExpression; break;
                case ChecksumField ck: yield return ck.ChecksumExpression; break;
                case ArrayField a: yield return a.LoopExpression; break;
                case VirtualField v: yield return v.ValueExpression; break;
                case ConstField c2: yield return c2.ReferenceValue; break;
                case ReservedField r: yield return r.ReferenceValue; break;
                case OptionalField o: yield return o.Condition; break;
                case ManualField m:
                    yield return m.ParseExpression;
                    yield return m.SerializeExpression;
                    yield return m.LengthExpression;
                    break;
                case PaddingField p:
                    yield return p.PaddingValue;
                    yield return p.TimesPadding;
                    break;
            }
        }

        // ── helpers ─────────────────────────────────────────────

        private string ReadFieldValue(Field f, CSharpExpressionRenderer r) =>
            CSharpTypeMapper.ReadCall(f.Type, f.Name, ComplexArgs(f.Type, r));

        /// <summary>Reads an <c>enum</c> field. With a key accessor
        /// (<c>[enum TransportSize transportSize code]</c>) the wire value is
        /// the key, mapped back through the generated reverse lookup.</summary>
        private string EnumFieldRead(EnumField ef, CSharpExpressionRenderer r)
        {
            if (ef.Type is EnumTypeReference et && ef.KeyAccessor != null)
            {
                var baseRead = CSharpTypeMapper.ReadCall(et.BaseType, ef.Name);
                return $"{et.Name}Extensions.FirstEnumForField{Pascal(ef.KeyAccessor)}({baseRead})";
            }
            return ReadFieldValue(ef, r);
        }

        private string EnumFieldWrite(EnumField ef)
        {
            if (ef.Type is EnumTypeReference et && ef.KeyAccessor != null)
            {
                var baseCs = CSharpTypeMapper.SimpleCSharpType(et.BaseType);
                var val = $"({baseCs}) {Pascal(ef.Name)}.Get{Pascal(ef.KeyAccessor)}()";
                return CSharpTypeMapper.WriteCall(et.BaseType, ef.Name, val);
            }
            return CSharpTypeMapper.WriteCall(ef.Type, ef.Name, Pascal(ef.Name));
        }

        /// <summary><see cref="ComplexArgs"/> plus the <c>_lastItem</c> flag when
        /// the array loop knows it.</summary>
        private string ElementArgs(TypeReference type, CSharpExpressionRenderer r, string lastItem)
        {
            var args = ComplexArgs(type, r);
            if (lastItem == null)
            {
                return args;
            }
            return args.Length == 0 ? lastItem : $"{args}, {lastItem}";
        }

        private string ComplexArgs(TypeReference type, CSharpExpressionRenderer r)
        {
            if (type is not ComplexTypeReference c)
            {
                return "";
            }
            // Cast each argument to the callee's declared parameter type: an
            // mspec count / position expression is `int`, but a `uint 8`
            // parameter is `byte`.
            var callee = _protocol.FindType(c.Name);
            return string.Join(", ", c.Arguments.Select((a, i) =>
            {
                var rendered = r.Render(a);
                var pType = callee != null && i < callee.Arguments.Count
                    ? callee.Arguments[i].Type
                    : null;
                return pType != null ? WithType(rendered, pType) : rendered;
            }));
        }

        private static TypeReference ElementType(ArrayField af) => af.Type;

        private string PropertyType(Field f)
        {
            if (f is ArrayField af)
            {
                var (t, _) = CSharpTypeMapper.ArrayType(af.Type);
                return t;
            }
            var t2 = CSharpTypeMapper.CSharpType(f.Type);
            // An optional field is absent when its condition is false.
            return f is OptionalField ? t2 + "?" : t2;
        }

        /// <summary>One switch discriminator, resolved against the parent.</summary>
        private sealed class Discriminator
        {
            /// <summary>The C# accessor name a child would override, or null
            /// when the value already lives in a property or is a dotted path
            /// (both cases: no accessor is generated).</summary>
            public string AccessorName { get; init; }

            public TypeReference Type { get; init; }

            /// <summary>The C# that reads the value during dispatch - a local
            /// (a discriminator or simple field the parent already read), a
            /// parser-argument name, or a rendered dotted path.</summary>
            public string DispatchExpr { get; init; }
        }

        private List<Discriminator> DiscriminatorList(
            ComplexTypeDefinition parent, CSharpExpressionRenderer render = null)
        {
            var ts = parent.TypeSwitch;
            var result = new List<Discriminator>();
            if (ts == null)
            {
                return result;
            }

            foreach (var term in ts.Discriminators.OfType<VariableLiteral>())
            {
                if (term.Child != null || term.IsCall)
                {
                    // e.g. parameter.parameterType - used for dispatch, never an accessor.
                    var r = render ?? new CSharpExpressionRenderer(new FieldScope(_protocol, parent));
                    result.Add(new Discriminator
                    {
                        AccessorName = null,
                        Type = new SimpleTypeReference { BaseType = SimpleTypeReference.Base.UInt, SizeInBits = 8 },
                        DispatchExpr = r.Render(term),
                    });
                    continue;
                }

                var name = term.Name;
                var field = parent.Fields.FirstOrDefault(x => x.Name == name);
                var arg = parent.Arguments.FirstOrDefault(a => a.Name == name);
                var type = field?.Type ?? arg?.Type
                    ?? new SimpleTypeReference { BaseType = SimpleTypeReference.Base.UInt, SizeInBits = 8 };

                result.Add(new Discriminator
                {
                    // A property field carries the value itself; only a bare
                    // `discriminator` / parser-arg gets an accessor.
                    AccessorName = field is { IsProperty: true } ? null : name,
                    Type = type,
                    DispatchExpr = Camel(name),
                });
            }
            return result;
        }

        private IEnumerable<string> ParserArgDiscriminatorNames(ComplexTypeDefinition parent) =>
            DiscriminatorList(parent)
                .Where(d => d.AccessorName != null
                            && parent.Arguments.Any(a => a.Name == d.AccessorName)
                            && parent.Fields.All(f => f.Name != d.AccessorName))
                .Select(d => d.AccessorName);

        private IReadOnlyList<Field> ParentPrefixFields(ComplexTypeDefinition child) =>
            child.ParentName != null && _protocol.FindType(child.ParentName) is { } p
                ? p.PrefixPropertyFields.ToList()
                : Array.Empty<Field>();

        private IReadOnlyList<Field> ParentSuffixFields(ComplexTypeDefinition child) =>
            child.ParentName != null && _protocol.FindType(child.ParentName) is { } p
                ? p.SuffixPropertyFields.ToList()
                : Array.Empty<Field>();

        /// <summary>The parent's pre-typeSwitch value-carrying fields that this
        /// child actually names in one of its own or trailing-field
        /// expressions - the child's <c>StaticParse</c> takes these as extra
        /// parameters so those expressions resolve. Filtered to what is
        /// referenced so an unrelated protocol (Modbus) gets no spurious
        /// parameters.</summary>
        private IReadOnlyList<Field> ParentContextFields(ComplexTypeDefinition child)
        {
            if (child.ParentName == null || _protocol.FindType(child.ParentName) is not { } p)
            {
                return Array.Empty<Field>();
            }
            var referenced = new HashSet<string>();
            foreach (var f in child.Fields.Concat(p.FieldsAfterSwitch))
            {
                foreach (var t in AllFieldTerms(f))
                {
                    CollectNames(t, referenced);
                }
            }
            // Prefix property fields are always needed - the child's constructor
            // takes them (see AllValueFields). Prefix discriminators / implicits
            // / consts come along only when an expression names them.
            return p.PrefixContextFields
                .Where(f => f.IsProperty || referenced.Contains(f.Name))
                .ToList();
        }

        private static IEnumerable<Term> AllFieldTerms(Field f)
        {
            foreach (var t in FieldTerms(f))
            {
                yield return t;
            }
            if (f.Type is ComplexTypeReference c)
            {
                foreach (var a in c.Arguments)
                {
                    yield return a;
                }
            }
        }

        private static void CollectNames(Term t, HashSet<string> sink)
        {
            switch (t)
            {
                case VariableLiteral v:
                    if (v.Name != null)
                    {
                        sink.Add(v.Name);
                    }
                    foreach (var a in v.Args ?? Enumerable.Empty<Term>())
                    {
                        CollectNames(a, sink);
                    }
                    foreach (var ix in v.Index ?? Enumerable.Empty<Term>())
                    {
                        CollectNames(ix, sink);
                    }
                    CollectNames(v.Child, sink);
                    break;
                case BinaryExpression b:
                    CollectNames(b.Left, sink);
                    CollectNames(b.Right, sink);
                    break;
                case UnaryExpression u:
                    CollectNames(u.Operand, sink);
                    break;
                case TernaryExpression tt:
                    CollectNames(tt.Condition, sink);
                    CollectNames(tt.WhenTrue, sink);
                    CollectNames(tt.WhenFalse, sink);
                    break;
            }
        }

        /// <summary>Constructor parameters and property fields for a type, in
        /// wire order: inherited prefix, own, inherited suffix.</summary>
        private List<Field> AllValueFields(ComplexTypeDefinition type) =>
            ParentPrefixFields(type)
                .Concat(type.PropertyFields)
                .Concat(ParentSuffixFields(type))
                .ToList();

        /// <summary>The fields this type's own parse / serialize / length
        /// covers. A discriminated parent stops at the typeSwitch; a child
        /// picks up the parent's trailing fields after its own body.</summary>
        private IReadOnlyList<Field> FieldsToEmit(ComplexTypeDefinition type)
        {
            if (type.IsDiscriminatedParent)
            {
                return type.FieldsThroughSwitch.ToList();
            }
            if (type.IsDiscriminatedChild && _protocol.FindType(type.ParentName) is { } p)
            {
                return type.Fields.Concat(p.FieldsAfterSwitch).ToList();
            }
            return type.Fields;
        }

        private string RenderDiscriminatorValue(Term value, TypeReference dType, CSharpExpressionRenderer r)
        {
            if (dType is EnumTypeReference e && value is VariableLiteral v && v.Child == null && !v.IsCall)
            {
                return $"{e.Name}.{v.Name}";
            }
            return WithType(r.Render(value), dType);
        }

        private static string WithType(string rendered, TypeReference type)
        {
            if (type is SimpleTypeReference s && s.IsIntegerLike && s.BaseType != SimpleTypeReference.Base.Bit)
            {
                return $"({CSharpTypeMapper.SimpleCSharpType(s)}) ({rendered})";
            }
            return rendered;
        }

        private List<string> ParseSignatureArgs(ComplexTypeDefinition type)
        {
            // A child inherits the parent's parser arguments, and also takes
            // the parent's prefix fields (read by the parent, passed down so
            // the child can build a complete object).
            var owner = type.IsDiscriminatedChild ? _protocol.FindType(type.ParentName) : type;
            var args = owner.Arguments
                .Select(a => $"{CSharpTypeMapper.CSharpType(a.Type)} {Camel(a.Name)}")
                .ToList();
            args.AddRange(ParentContextFields(type)
                .Select(f => $"{PropertyType(f)} {Camel(f.Name)}"));
            if (type.IsDiscriminatedChild && UsesCurPos(type))
            {
                args.Add("int _startPos");
            }
            if (UsesLastItem(type))
            {
                args.Add("bool _lastItem = false");
            }
            return args;
        }

        /// <summary>Whether an expression in this type reads <c>_lastItem</c> -
        /// true only for a type used as an array element whose last element is
        /// laid out differently (S7's <c>S7VarPayloadDataItem</c> padding).</summary>
        private bool UsesLastItem(ComplexTypeDefinition type)
        {
            var names = new HashSet<string>();
            foreach (var f in FieldsToEmit(type))
            {
                foreach (var t in AllFieldTerms(f))
                {
                    CollectNames(t, names);
                }
            }
            return names.Contains("_lastItem");
        }

        /// <summary>Whether <paramref name="af"/>'s element type threads
        /// <c>_lastItem</c> - the loop then has to pass the per-iteration flag
        /// and use an index-aware serialize / length.</summary>
        private bool ElementUsesLastItem(ArrayField af) =>
            ElementType(af) is ComplexTypeReference c
            && _protocol.FindType(c.Name) is { } et
            && UsesLastItem(et);

        /// <summary>Whether this type's own parse emits a field whose expression
        /// reads <c>curPos</c> - it then needs <c>_startPos</c> in scope.</summary>
        private bool UsesCurPos(ComplexTypeDefinition type)
        {
            var names = new HashSet<string>();
            foreach (var f in FieldsToEmit(type))
            {
                foreach (var t in AllFieldTerms(f))
                {
                    CollectNames(t, names);
                }
            }
            return names.Contains("curPos") || names.Contains("_curPos");
        }

        /// <summary>A discriminated parent that itself does not read
        /// <c>curPos</c> still captures <c>_startPos</c> when a child needs it,
        /// to pass down.</summary>
        private bool NeedsStartPos(ComplexTypeDefinition type) =>
            UsesCurPos(type)
            || type.IsDiscriminatedParent
               && type.TypeSwitch.CaseNames
                   .Select(_protocol.FindType)
                   .Where(child => child != null)
                   .Any(UsesCurPos);

        private IEnumerable<string> ChildParseArgNames(ComplexTypeDefinition child)
        {
            var parent = _protocol.FindType(child.ParentName);
            return parent.Arguments.Select(a => Camel(a.Name));
        }

        private static string Comma(IEnumerable<string> parts)
        {
            var list = parts.ToList();
            return list.Count == 0 ? "" : ", " + string.Join(", ", list);
        }

        private static string Comma(string s) => string.IsNullOrEmpty(s) ? "" : ", " + s;

        private void Header(CodeWriter c)
        {
            foreach (var line in LicenseHeader)
            {
                c.Line(line);
            }
            c.Line();
            c.Line("// Code generated from the mspec by plc4net-code-gen. DO NOT EDIT.");
            c.Line();
        }

        private static readonly string[] LicenseHeader =
        {
            "//",
            "// Licensed to the Apache Software Foundation (ASF) under one",
            "// or more contributor license agreements.  See the NOTICE file",
            "// distributed with this work for additional information",
            "// regarding copyright ownership.  The ASF licenses this file",
            "// to you under the Apache License, Version 2.0 (the",
            "// \"License\"); you may not use this file except in compliance",
            "// with the License.  You may obtain a copy of the License at",
            "//",
            "//      https://www.apache.org/licenses/LICENSE-2.0",
            "//",
            "// Unless required by applicable law or agreed to in writing,",
            "// software distributed under the License is distributed on an",
            "// \"AS IS\" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY",
            "// KIND, either express or implied.  See the License for the",
            "// specific language governing permissions and limitations",
            "// under the License.",
            "//",
        };

        internal static string Pascal(string s) =>
            string.IsNullOrEmpty(s) ? s : char.ToUpperInvariant(s[0]) + s.Substring(1);

        internal static string Camel(string s) =>
            string.IsNullOrEmpty(s) ? s : char.ToLowerInvariant(s[0]) + s.Substring(1);
    }
}
