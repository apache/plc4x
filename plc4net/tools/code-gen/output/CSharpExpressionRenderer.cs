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
using System.Globalization;
using System.Linq;
using org.apache.plc4net.tools.codegen.model.terms;

namespace org.apache.plc4net.tools.codegen.output
{
    /// <summary>
    /// Resolves a bare mspec identifier to the C# text that reads its value in
    /// the current context.
    /// </summary>
    /// <remarks>
    /// The same mspec name means different things depending on where it is
    /// rendered: while parsing it is a local holding a value just read off the
    /// wire; while serializing it is a property on the message; a parser
    /// argument is a method parameter either way. The generator supplies the
    /// mapping; the renderer only asks.
    /// </remarks>
    public interface IExpressionScope
    {
        /// <summary>The C# for reading <paramref name="name"/>, or
        /// <c>null</c> to fall back to the name unchanged.</summary>
        string ResolveReference(string name);

        /// <summary>The C# for "how many elements in the collection
        /// <paramref name="name"/>" - <c>.Length</c> for a <c>byte[]</c>,
        /// <c>.Count</c> for a <c>List&lt;T&gt;</c>. Return <c>null</c> to let
        /// the renderer default to <c>.Count</c>.</summary>
        string ResolveCount(string name) => null;

        /// <summary>The C# for a member access <c>x.<paramref name="name"/></c>
        /// where <c>name</c> is not a plain property - an enum attribute
        /// accessor such as <c>transportSize.sizeInBits</c> becomes a method
        /// call. Return <c>null</c> to fall back to the PascalCase member.</summary>
        string ResolveMember(string name) => null;
    }

    /// <summary>
    /// Renders a <see cref="Term"/> tree as a C# expression.
    /// </summary>
    /// <remarks>
    /// Operators map straight across except mspec <c>^</c>, which is
    /// exponentiation (C# <c>^</c> is xor). The ALL-CAPS mspec built-ins
    /// (<c>COUNT</c>, <c>CEIL</c>, <c>ARRAY_SIZE_IN_BYTES</c>,
    /// <c>STATIC_CALL</c>) are expanded here; anything else that looks like a
    /// call is passed through as a C# call.
    /// </remarks>
    public class CSharpExpressionRenderer
    {
        private readonly IExpressionScope _scope;

        /// <summary>Name of the generated per-protocol static helper class that
        /// carries <c>STATIC_CALL</c> targets and array-size helpers.</summary>
        public string StaticHelperClass { get; set; } = "StaticHelper";

        public CSharpExpressionRenderer(IExpressionScope scope = null)
        {
            _scope = scope;
        }

        /// <summary>mspec member accessors that are methods in the generated model.</summary>
        private static readonly IReadOnlyDictionary<string, string> MemberAccessors =
            new Dictionary<string, string>
            {
                ["lengthInBytes"] = "GetLengthInBytes()",
                ["lengthInBits"] = "GetLengthInBits()",
            };

        public string Render(Term term)
        {
            switch (term)
            {
                case null:
                    return "null";
                case IntegerLiteral i:
                    return i.Text;
                case FloatLiteral f:
                    var ft = f.Value.ToString("R", CultureInfo.InvariantCulture);
                    return ft.IndexOfAny(new[] { '.', 'e', 'E' }) < 0 ? ft + ".0" : ft;
                case BooleanLiteral b:
                    return b.Value ? "true" : "false";
                case NullLiteral _:
                    return "null";
                case StringLiteral s:
                    return "\"" + s.Value.Replace("\\", "\\\\").Replace("\"", "\\\"") + "\"";
                case UnaryExpression u:
                    return u.Operator + Wrap(u.Operand);
                case TernaryExpression t:
                    return $"({Render(t.Condition)} ? {Render(t.WhenTrue)} : {Render(t.WhenFalse)})";
                case BinaryExpression bin:
                    return RenderBinary(bin);
                case VariableLiteral v:
                    return RenderVariable(v);
                default:
                    throw new NotSupportedException(
                        $"Cannot render term {term.GetType().Name}");
            }
        }

        private string RenderBinary(BinaryExpression bin)
        {
            if (bin.Operator == "^")
            {
                // mspec '^' is exponentiation.
                return $"System.Math.Pow({Render(bin.Left)}, {Render(bin.Right)})";
            }
            return $"({Render(bin.Left)} {bin.Operator} {Render(bin.Right)})";
        }

        private string Wrap(Term t)
        {
            var rendered = Render(t);
            return t is BinaryExpression || t is TernaryExpression ? "(" + rendered + ")" : rendered;
        }

        private string RenderVariable(VariableLiteral v)
        {
            string head;
            if (v.IsCall && RenderBuiltin(v) is { } expanded)
            {
                head = expanded;
            }
            else if (v.IsCall)
            {
                head = (Resolve(v.Name) ?? v.Name) + "(" + RenderArgs(v.Args) + ")";
            }
            else
            {
                head = Resolve(v.Name) ?? v.Name;
            }

            head += RenderIndexes(v.Index);

            for (var child = v.Child; child != null; child = child.Child)
            {
                head += "." + RenderMember(child);
            }

            return head;
        }

        private string RenderMember(VariableLiteral member)
        {
            string text;
            if (member.IsCall)
            {
                text = Pascal(member.Name) + "(" + RenderArgs(member.Args) + ")";
            }
            else if (MemberAccessors.TryGetValue(member.Name, out var accessor))
            {
                text = accessor;
            }
            else if (_scope?.ResolveMember(member.Name) is { } resolved)
            {
                text = resolved;
            }
            else
            {
                text = Pascal(member.Name);
            }
            return text + RenderIndexes(member.Index);
        }

        private string RenderIndexes(IReadOnlyList<Term> index)
            => index == null || index.Count == 0
                ? string.Empty
                : string.Concat(index.Select(i => "[" + Render(i) + "]"));

        private string RenderArgs(IReadOnlyList<Term> args)
            => string.Join(", ", (args ?? Array.Empty<Term>()).Select(Render));

        private string RenderBuiltin(VariableLiteral v)
        {
            switch (v.Name)
            {
                case "COUNT" when v.Args.Count == 1:
                    // byte[] wants .Length, List<T> wants .Count - the scope
                    // knows which; .Count is the fallback.
                    if (v.Args[0] is VariableLiteral { Child: null, Args: null } arg
                        && _scope?.ResolveCount(arg.Name) is { } count)
                    {
                        return count;
                    }
                    return Render(v.Args[0]) + ".Count";
                case "CEIL" when v.Args.Count == 1:
                    return $"(int) System.Math.Ceiling((double) ({Render(v.Args[0])}))";
                case "ARRAY_SIZE_IN_BYTES" when v.Args.Count == 1:
                    return $"{StaticHelperClass}.ArraySizeInBytes({Render(v.Args[0])})";
                case "STATIC_CALL" when v.Args.Count >= 1 && v.Args[0] is StringLiteral fn:
                    var rest = v.Args.Skip(1).Select(Render);
                    return $"{StaticHelperClass}.{Pascal(fn.Value)}({string.Join(", ", rest)})";
                case "CAST" when v.Args.Count == 2 && v.Args[1] is StringLiteral castType:
                    // CAST(x, "T") -> ((T) x); a trailing member access stays.
                    return $"(({castType.Value}) {Render(v.Args[0])})";
                default:
                    return null;
            }
        }

        private string Resolve(string name) => _scope?.ResolveReference(name);

        private static string Pascal(string name)
            => string.IsNullOrEmpty(name)
                ? name
                : char.ToUpperInvariant(name[0]) + name.Substring(1);
    }
}
