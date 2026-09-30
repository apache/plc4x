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
using Antlr4.Runtime;
using Antlr4.Runtime.Misc;
using org.apache.plc4net.tools.codegen.grammar;
using org.apache.plc4net.tools.codegen.model;
using org.apache.plc4net.tools.codegen.model.fields;
using org.apache.plc4net.tools.codegen.model.terms;

namespace org.apache.plc4net.tools.codegen
{
    /// <summary>
    /// Turns a parsed mspec file into the type-model IR
    /// (<see cref="Protocol"/>). Walks the ANTLR parse tree with the generated
    /// typed context classes; expressions inside ticks are handed to
    /// <see cref="MspecExpressionParser"/>.
    /// </summary>
    public sealed class MspecModelBuilder
    {
        private readonly Protocol _protocol = new Protocol();

        /// <summary>Set by <c>BuildFiles</c> so diagnostics and errors can
        /// name the mspec file they come from.</summary>
        private string? _currentFile;

        public static Protocol Build(string mspecContent)
            => new MspecModelBuilder().Run(MspecReader.Read(mspecContent));

        public static Protocol BuildFile(string path)
            => BuildFiles(path);

        /// <summary>
        /// Builds one model from several mspec files - a protocol whose types
        /// are split across files (knxnetip.mspec references
        /// <c>KnxPropertyDataType</c> from device-info.mspec) is compiled the
        /// way the Java plugin does it: every <c>*.mspec</c> in the directory
        /// as one unit.
        /// </summary>
        public static Protocol BuildFiles(params string[] paths) => BuildFiles(true, paths);

        public static Protocol BuildFiles(bool strict, params string[] paths)
        {
            var builder = new MspecModelBuilder();
            foreach (var path in paths)
            {
                builder._currentFile = path;
                builder.RunInto(MspecReader.ReadFile(path, strict));
            }
            builder.ResolveEnumReferences();
            return builder._protocol;
        }

        private Protocol Run(MSpecParser.FileContext file)
        {
            RunInto(file);
            ResolveEnumReferences();
            return _protocol;
        }

        private void RunInto(MSpecParser.FileContext file)
        {
            var constants = file.contantsDefinition();
            if (constants != null)
            {
                foreach (var cf in constants.constField())
                {
                    _protocol.Constants.Add(new ConstantDeclaration
                    {
                        Name = cf.name.GetText(),
                        // A const's type token is mandatory in the grammar.
                        Type = BuildTypeReference(cf.type)!,
                        Value = ParseValueLiteral(cf.expected),
                    });
                }
            }

            foreach (var def in file.complexTypeDefinition())
            {
                BuildComplexTypeDefinition(def.complexType());
            }
        }

        // ── type definitions ─────────────────────────────────────

        private void BuildComplexTypeDefinition(MSpecParser.ComplexTypeContext ctx)
        {
            if (ctx.ENUM() != null)
            {
                BuildEnum(ctx);
                return;
            }
            if (ctx.DATAIO() != null)
            {
                BuildDataIo(ctx);
                return;
            }

            var type = new ComplexTypeDefinition { Name = ctx.name.GetText() };
            FillArguments(ctx.argumentList(), type.Arguments);
            FillAttributes(ctx.attributeList(), type.Attributes);
            foreach (var fd in ctx.fieldDefinition())
            {
                AddField(type, fd);
            }
            AddTypeDefinition(type);
        }

        private void BuildEnum(MSpecParser.ComplexTypeContext ctx)
        {
            var e = new EnumTypeDefinition
            {
                Name = ctx.name.GetText(),
                BaseType = ctx.dataType() != null ? BuildSimpleType(ctx.dataType()) : null,
            };
            FillArguments(ctx.argumentList(), e.Arguments);
            FillAttributes(ctx.attributeList(), e.Attributes);

            foreach (var v in ctx.enumValueDefinition())
            {
                e.Values.Add(new EnumValue
                {
                    Name = v.name.Text,
                    Value = v.valueExpression != null ? ParseExpression(v.valueExpression) : null,
                    // Every element here comes from a real expression node in
                    // the constant-value list, never from an absent optional.
                    ConstantValues = v.constantValueExpressions != null
                        ? v.constantValueExpressions.expression().Select(e => ParseExpression(e)!).ToList()
                        : (IReadOnlyList<Term>)System.Array.Empty<Term>(),
                });
            }
            AddTypeDefinition(e);
        }

        private void BuildDataIo(MSpecParser.ComplexTypeContext ctx)
        {
            var dio = new DataIoTypeDefinition { Name = ctx.name.GetText() };
            FillArguments(ctx.argumentList(), dio.Arguments);
            FillAttributes(ctx.attributeList(), dio.Attributes);

            var tsCtx = ctx.dataIoDefinition()?.typeSwitchField();
            if (tsCtx != null)
            {
                // dataIo case sub-types are not lifted to top-level types: the
                // names repeat (CHAR / STRING / TIME each appear twice) and a
                // dataIo parses to an IPlcValue, not a generated class. They are
                // kept on the DataIoTypeDefinition for the emitter.
                dio.TypeSwitch = BuildTypeSwitch(tsCtx, dio.Name, dio.Cases);
            }
            AddTypeDefinition(dio);
        }

        // ── fields ───────────────────────────────────────────────

        private void AddField(ComplexTypeDefinition owner, MSpecParser.FieldDefinitionContext fd)
        {
            var f = fd.field();
            Field field = f switch
            {
                _ when f.simpleField() is { } s => new SimpleField
                {
                    Name = s.name.GetText(),
                    Type = BuildTypeReference(s.type),
                },
                _ when f.constField() is { } c => new ConstField
                {
                    Name = c.name.GetText(),
                    Type = BuildTypeReference(c.type),
                    ReferenceValue = ParseValueLiteral(c.expected),
                },
                _ when f.implicitField() is { } im => new ImplicitField
                {
                    Name = im.name.GetText(),
                    Type = BuildSimpleType(im.type),
                    // The serialize expression is mandatory in the grammar.
                    SerializeExpression = ParseExpression(im.serializeExpression)!,
                },
                _ when f.reservedField() is { } r => new ReservedField
                {
                    Name = "reserved",
                    Type = BuildSimpleType(r.type),
                    // The reserved value literal is mandatory in the grammar.
                    ReferenceValue = ParseExpression(r.expected)!,
                },
                _ when f.discriminatorField() is { } d => new DiscriminatorField
                {
                    Name = d.name.GetText(),
                    Type = BuildTypeReference(d.type),
                },
                _ when f.enumField() is { } en => new EnumField
                {
                    Name = en.name.GetText(),
                    Type = BuildTypeReference(en.type),
                    KeyAccessor = en.fieldName?.GetText(),
                },
                _ when f.arrayField() is { } a => new ArrayField
                {
                    Name = a.name.GetText(),
                    Type = BuildTypeReference(a.type),
                    LoopType = a.loopType.Text.Trim('\'') switch
                    {
                        "length" => ArrayField.Loop.Length,
                        "terminated" => ArrayField.Loop.Terminated,
                        _ => ArrayField.Loop.Count,
                    },
                    // The loop (count / length / terminator) expression is
                    // mandatory in the grammar.
                    LoopExpression = ParseExpression(a.loopExpression)!,
                },
                _ when f.checksumField() is { } ck => new ChecksumField
                {
                    Name = ck.name.GetText(),
                    Type = BuildSimpleType(ck.type),
                    ChecksumExpression = ParseExpression(ck.checksumExpression)!,
                },
                _ when f.virtualField() is { } vf => new VirtualField
                {
                    Name = vf.name.GetText(),
                    Type = BuildTypeReference(vf.type),
                    ValueExpression = ParseExpression(vf.valueExpression)!,
                },
                _ when f.optionalField() is { } of => new OptionalField
                {
                    Name = of.name.GetText(),
                    Type = BuildTypeReference(of.type),
                    Condition = of.condition != null ? ParseExpression(of.condition) : null,
                },
                // All three manual expressions are mandatory in the grammar.
                _ when f.manualField() is { } mf => new ManualField
                {
                    Name = mf.name.GetText(),
                    Type = BuildTypeReference(mf.type),
                    ParseExpression = ParseExpression(mf.parseExpression)!,
                    SerializeExpression = ParseExpression(mf.serializeExpression)!,
                    LengthExpression = ParseExpression(mf.lengthExpression)!,
                },
                // Both padding expressions are mandatory in the grammar.
                _ when f.paddingField() is { } pf => new PaddingField
                {
                    Name = pf.name.GetText(),
                    Type = BuildSimpleType(pf.type),
                    PaddingValue = ParseExpression(pf.paddingValue)!,
                    TimesPadding = ParseExpression(pf.timesPadding)!,
                },
                _ when f.typeSwitchField() is { } ts => BuildTypeSwitch(ts, owner.Name, _protocol.Types),
                // The field context always has at least the keyword token as
                // its first child.
                _ => new UnsupportedField
                {
                    Name = "unsupported",
                    MspecKeyword = f.GetChild(0)?.GetText() ?? f.GetType().Name,
                    RawText = SourceText(f),
                },
            };

            FillAttributes(fd.attributeList(), field.Attributes);
            ApplyStringEncoding(field);
            owner.Fields.Add(field);
        }

        private TypeSwitchField BuildTypeSwitch(
            MSpecParser.TypeSwitchFieldContext ts, string parentName,
            List<ComplexTypeDefinition> childSink)
        {
            var field = new TypeSwitchField
            {
                Name = "typeSwitch",
                Discriminators = ts.multipleVariableLiterals().variableLiteral()
                    .Select(v =>
                    {
                        try { return MspecExpressionParser.Parse(SourceText(v)); }
                        catch (MspecParseException) { return new VariableLiteral(SourceText(v)); }
                    })
                    .ToList(),
            };

            foreach (var cs in ts.caseStatement())
            {
                var child = new ComplexTypeDefinition
                {
                    Name = cs.name.Text,
                    ParentName = parentName,
                    DiscriminatorValues = cs.discriminatorValues != null
                        ? cs.discriminatorValues.expression().Select(ParseExpression).ToList()
                        : (IReadOnlyList<Term?>)System.Array.Empty<Term>(),
                };
                FillArguments(cs.argumentList(), child.Arguments);
                foreach (var fd in cs.fieldDefinition())
                {
                    AddField(child, fd);
                }

                field.CaseNames.Add(child.Name);
                if (ReferenceEquals(childSink, _protocol.Types))
                {
                    // A lifted child becomes its own output file; a name
                    // shared with any other type would silently collide.
                    AddTypeDefinition(child);
                }
                else
                {
                    // dataIo cases stay local and may repeat names.
                    childSink.Add(child);
                }
            }

            return field;
        }

        // ── type references ──────────────────────────────────────

        private TypeReference? BuildTypeReference(MSpecParser.TypeReferenceContext ctx)
        {
            if (ctx == null)
            {
                return null;
            }
            if (ctx.simpleTypeReference != null)
            {
                return BuildSimpleType(ctx.simpleTypeReference);
            }

            return new ComplexTypeReference
            {
                Name = ctx.complexTypeReference.Text,
                // Every element here is a real parser-argument expression,
                // never from an absent optional.
                Arguments = ctx.@params != null
                    ? ctx.@params.expression().Select(e => ParseExpression(e)!).ToList()
                    : (IReadOnlyList<Term>)System.Array.Empty<Term>(),
            };
        }

        /// <summary>Registers a type / enum / dataIo, rejecting a name that
        /// already exists anywhere in the model: output files are keyed by
        /// type name, so a second definition would silently win the file.</summary>
        private void AddTypeDefinition(TypeDefinition definition)
        {
            var clash = _protocol.Types.FirstOrDefault(t => t.Name == definition.Name)
                ?? (TypeDefinition?)_protocol.Enums.FirstOrDefault(e => e.Name == definition.Name)
                ?? _protocol.DataIos.FirstOrDefault(d => d.Name == definition.Name);
            if (clash != null)
            {
                throw new System.InvalidOperationException(
                    $"duplicate type name '{definition.Name}' in {_currentFile ?? "<unknown file>"}; "
                    + "output is keyed by type name, so the second definition would overwrite the first");
            }
            switch (definition)
            {
                case EnumTypeDefinition e:
                    _protocol.Enums.Add(e);
                    break;
                case DataIoTypeDefinition d:
                    _protocol.DataIos.Add(d);
                    break;
                default:
                    _protocol.Types.Add((ComplexTypeDefinition)definition);
                    break;
            }
        }

        private SimpleTypeReference BuildSimpleType(MSpecParser.DataTypeContext dt)
        {
            var baseText = dt.@base.Text;
            var size = dt.size != null
                ? int.Parse(dt.size.Text, System.Globalization.CultureInfo.InvariantCulture)
                : 0;

            var (baseType, bits) = baseText switch
            {
                "bit" => (SimpleTypeReference.Base.Bit, 1),
                "byte" => (SimpleTypeReference.Base.Byte, 8),
                "uint" => (SimpleTypeReference.Base.UInt, size),
                "int" => (SimpleTypeReference.Base.Int, size),
                "float" => (SimpleTypeReference.Base.Float, size),
                "ufloat" => (SimpleTypeReference.Base.UFloat, size),
                "string" => (SimpleTypeReference.Base.String, size),
                "vstring" => (SimpleTypeReference.Base.VString, size),
                "time" => (SimpleTypeReference.Base.Time, 32),
                "date" => (SimpleTypeReference.Base.Date, 32),
                "dateTime" => (SimpleTypeReference.Base.DateTime, 64),
                _ => throw new System.NotSupportedException(
                    $"base type '{baseText}' is not modelled (vint/vuint have no C# mapping yet)"),
            };

            return new SimpleTypeReference
            {
                BaseType = baseType,
                SizeInBits = bits,
                // `vstring 'expr'` - the length in bits is a run-time expression
                // (ADS: `vstring 'stringLength * 8' value`).
                LengthExpression = baseType == SimpleTypeReference.Base.VString && dt.length != null
                    ? ParseExpression(dt.length)
                    : null,
            };
        }

        /// <summary>Promotes a <see cref="ComplexTypeReference"/> to an
        /// <see cref="EnumTypeReference"/> once every enum name is known.</summary>
        private void ResolveEnumReferences()
        {
            var enums = _protocol.Enums.ToDictionary(e => e.Name);

            EnumTypeReference? Promote(ComplexTypeReference c) =>
                enums.TryGetValue(c.Name, out var e)
                    ? new EnumTypeReference
                    {
                        Name = e.Name,
                        BaseType = e.BaseType
                            ?? new SimpleTypeReference { BaseType = SimpleTypeReference.Base.UInt, SizeInBits = 8 },
                    }
                    : null;

            void Fix(System.Func<TypeReference> get, System.Action<TypeReference> set)
            {
                if (get() is ComplexTypeReference c && Promote(c) is { } er)
                {
                    set(er);
                }
            }

            foreach (var t in _protocol.Types)
            {
                foreach (var a in t.Arguments)
                {
                    Fix(() => a.Type, v => a.Type = v);
                }
                // Field.Type is nullable (null for typeSwitch), unlike
                // Argument.Type above, so it is handled separately rather
                // than widening Fix's own delegate types for one caller.
                foreach (var f in t.Fields)
                {
                    if (f.Type is ComplexTypeReference fc && Promote(fc) is { } fer)
                    {
                        f.Type = fer;
                    }
                }
            }
            foreach (var en in _protocol.Enums)
            {
                foreach (var a in en.Arguments)
                {
                    Fix(() => a.Type, v => a.Type = v);
                }
            }
            foreach (var dio in _protocol.DataIos)
            {
                foreach (var a in dio.Arguments)
                {
                    Fix(() => a.Type, v => a.Type = v);
                }
            }
            foreach (var c in _protocol.Constants)
            {
                Fix(() => c.Type, v => c.Type = v);
            }
        }

        // ── expressions, attributes, arguments ───────────────────

        private Term? ParseExpression(MSpecParser.ExpressionContext exprCtx)
        {
            if (exprCtx == null)
            {
                return null;
            }

            var raw = SourceText(exprCtx).Trim();
            if (raw == "*")
            {
                return new VariableLiteral("*");
            }

            var text = raw.Length >= 2 && raw[0] == '\'' && raw[^1] == '\''
                ? raw.Substring(1, raw.Length - 2)
                : raw;

            try
            {
                return MspecExpressionParser.Parse(text);
            }
            catch (MspecParseException)
            {
                // Expression positions (count/length/serialize/...) never
                // legitimately hold prose, so a failure here is a broken
                // mspec; record it rather than silently emitting a string
                // literal that breaks generated code far from the cause.
                _protocol.Diagnostics.Add(
                    $"{_currentFile ?? "<model>"}: expression '{text}' did not parse; kept as a string literal");
                return new StringLiteral(text);
            }
        }

        private Term ParseValueLiteral(MSpecParser.ValueLiteralContext ctx)
        {
            var text = SourceText(ctx).Trim();
            try
            {
                return MspecExpressionParser.Parse(text.Trim('\''));
            }
            catch (MspecParseException)
            {
                // Deliberate: value-table cells legitimately hold quoted
                // strings that are not expressions (KNX's 'DPST-1-1').
                return new StringLiteral(text);
            }
        }

        private void FillArguments(MSpecParser.ArgumentListContext list, List<Argument> target)
        {
            if (list == null)
            {
                return;
            }
            foreach (var a in list.argument())
            {
                target.Add(new Argument
                {
                    Name = a.name.GetText(),
                    // A parser argument's type token is mandatory in the grammar.
                    Type = BuildTypeReference(a.type)!,
                });
            }
        }

        private void FillAttributes(
            MSpecParser.AttributeListContext list, Dictionary<string, Term> target)
        {
            if (list == null)
            {
                return;
            }
            foreach (var attr in list.attribute())
            {
                // An attribute's value is mandatory in the grammar (name=value).
                target[attr.name.Text] = ParseExpression(attr.value)!;
            }
        }

        private static void ApplyStringEncoding(Field field)
        {
            if (field.Type is SimpleTypeReference s
                && field.Attributes.TryGetValue("stringEncoding", out var enc)
                && enc is StringLiteral sl)
            {
                s.Encoding = sl.Value;
            }
        }

        private static string SourceText(ParserRuleContext ctx) =>
            ctx.Start.InputStream.GetText(Interval.Of(ctx.Start.StartIndex, ctx.Stop.StopIndex));
    }
}
