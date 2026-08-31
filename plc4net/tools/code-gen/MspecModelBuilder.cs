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
                        Type = BuildTypeReference(cf.type),
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
            _protocol.Types.Add(type);
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
                    ConstantValues = v.constantValueExpressions != null
                        ? v.constantValueExpressions.expression().Select(ParseExpression).ToList()
                        : (IReadOnlyList<Term>)System.Array.Empty<Term>(),
                });
            }
            _protocol.Enums.Add(e);
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
            _protocol.DataIos.Add(dio);
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
                    SerializeExpression = ParseExpression(im.serializeExpression),
                },
                _ when f.reservedField() is { } r => new ReservedField
                {
                    Name = "reserved",
                    Type = BuildSimpleType(r.type),
                    ReferenceValue = ParseExpression(r.expected),
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
                    LoopExpression = ParseExpression(a.loopExpression),
                },
                _ when f.checksumField() is { } ck => new ChecksumField
                {
                    Name = ck.name.GetText(),
                    Type = BuildSimpleType(ck.type),
                    ChecksumExpression = ParseExpression(ck.checksumExpression),
                },
                _ when f.virtualField() is { } vf => new VirtualField
                {
                    Name = vf.name.GetText(),
                    Type = BuildTypeReference(vf.type),
                    ValueExpression = ParseExpression(vf.valueExpression),
                },
                _ when f.optionalField() is { } of => new OptionalField
                {
                    Name = of.name.GetText(),
                    Type = BuildTypeReference(of.type),
                    Condition = of.condition != null ? ParseExpression(of.condition) : null,
                },
                _ when f.manualField() is { } mf => new ManualField
                {
                    Name = mf.name.GetText(),
                    Type = BuildTypeReference(mf.type),
                    ParseExpression = ParseExpression(mf.parseExpression),
                    SerializeExpression = ParseExpression(mf.serializeExpression),
                    LengthExpression = ParseExpression(mf.lengthExpression),
                },
                _ when f.paddingField() is { } pf => new PaddingField
                {
                    Name = pf.name.GetText(),
                    Type = BuildSimpleType(pf.type),
                    PaddingValue = ParseExpression(pf.paddingValue),
                    TimesPadding = ParseExpression(pf.timesPadding),
                },
                _ when f.typeSwitchField() is { } ts => BuildTypeSwitch(ts, owner.Name, _protocol.Types),
                _ => new UnsupportedField
                {
                    Name = "unsupported",
                    MspecKeyword = f.GetChild(0)?.GetText(),
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
                        : (IReadOnlyList<Term>)System.Array.Empty<Term>(),
                };
                FillArguments(cs.argumentList(), child.Arguments);
                foreach (var fd in cs.fieldDefinition())
                {
                    AddField(child, fd);
                }

                field.CaseNames.Add(child.Name);
                childSink.Add(child);
            }

            return field;
        }

        // ── type references ──────────────────────────────────────

        private TypeReference BuildTypeReference(MSpecParser.TypeReferenceContext ctx)
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
                Arguments = ctx.@params != null
                    ? ctx.@params.expression().Select(ParseExpression).ToList()
                    : (IReadOnlyList<Term>)System.Array.Empty<Term>(),
            };
        }

        private static SimpleTypeReference BuildSimpleType(MSpecParser.DataTypeContext dt)
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
                _ => (SimpleTypeReference.Base.UInt, size),
            };

            return new SimpleTypeReference { BaseType = baseType, SizeInBits = bits };
        }

        /// <summary>Promotes a <see cref="ComplexTypeReference"/> to an
        /// <see cref="EnumTypeReference"/> once every enum name is known.</summary>
        private void ResolveEnumReferences()
        {
            var enums = _protocol.Enums.ToDictionary(e => e.Name);

            EnumTypeReference Promote(ComplexTypeReference c) =>
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
                foreach (var f in t.Fields)
                {
                    Fix(() => f.Type, v => f.Type = v);
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

        private static Term ParseExpression(MSpecParser.ExpressionContext exprCtx)
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
                return new StringLiteral(text);
            }
        }

        private static Term ParseValueLiteral(MSpecParser.ValueLiteralContext ctx)
        {
            var text = SourceText(ctx).Trim();
            try
            {
                return MspecExpressionParser.Parse(text.Trim('\''));
            }
            catch (MspecParseException)
            {
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
                    Type = BuildTypeReference(a.type),
                });
            }
        }

        private static void FillAttributes(
            MSpecParser.AttributeListContext list, Dictionary<string, Term> target)
        {
            if (list == null)
            {
                return;
            }
            foreach (var attr in list.attribute())
            {
                target[attr.name.Text] = ParseExpression(attr.value);
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
