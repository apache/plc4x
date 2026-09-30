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

namespace org.apache.plc4net.tools.codegen.output;

/// <summary>Resolves MSpec names in the context of generated C# code.</summary>
public interface IExpressionScope
{
    string? ResolveReference(string name);

    string? ResolveCount(string name) => null;

    string? ResolveMember(string name) => null;
}

/// <summary>Renders an MSpec expression tree as a C# expression.</summary>
public sealed class CSharpExpressionRenderer(IExpressionScope? scope = null)
{
    private static readonly IReadOnlyDictionary<string, string> MemberAccessors =
        new Dictionary<string, string>
        {
            ["lengthInBytes"] = "GetLengthInBytes()",
            ["lengthInBits"] = "GetLengthInBits()"
        };

    public string StaticHelperClass { get; set; } = "StaticHelper";

    public string Render(Term? term) => term switch
    {
        null => "null",
        IntegerLiteral literal => literal.Text,
        HexadecimalLiteral literal => literal.Text,
        FloatLiteral literal => RenderFloat(literal),
        BooleanLiteral literal => literal.Value ? "true" : "false",
        NullLiteral => "null",
        StringLiteral literal => RenderString(literal.Value),
        UnaryExpression expression => RenderUnary(expression),
        IndexExpression expression => RenderIndex(expression),
        TernaryExpression expression =>
            $"({Render(expression.Condition)} ? {Render(expression.WhenTrue)} : {Render(expression.WhenFalse)})",
        BinaryExpression expression => RenderBinary(expression),
        VariableLiteral variable => RenderVariable(variable),
        _ => throw new NotSupportedException($"Cannot render term {term.GetType().Name}")
    };

    private static string RenderFloat(FloatLiteral literal)
    {
        string text = literal.Value.ToString("R", CultureInfo.InvariantCulture);
        return text.IndexOfAny(['.', 'e', 'E']) < 0 ? text + ".0" : text;
    }

    private static string RenderString(string value) =>
        $"\"{value.Replace("\\", "\\\\").Replace("\"", "\\\"")}\"";

    private string RenderUnary(UnaryExpression expression)
    {
        string operand = Render(expression.Operand);
        return expression.Operator +
               (expression.Operand is UnaryExpression ? $"({operand})" : operand);
    }

    private string RenderIndex(IndexExpression expression)
    {
        string target = Render(expression.Target);
        if (expression.Target is UnaryExpression)
        {
            target = $"({target})";
        }
        return target + string.Concat(expression.Indexes.Select(index => $"[{Render(index)}]"));
    }

    private string RenderBinary(BinaryExpression expression) =>
        expression.Operator == "^"
            ? $"System.Math.Pow({Render(expression.Left)}, {Render(expression.Right)})"
            : $"({Render(expression.Left)} {expression.Operator} {Render(expression.Right)})";

    private string RenderVariable(VariableLiteral variable)
    {
        if (variable.Name == MspecReader.EmptyStringSentinel &&
            !variable.IsCall && variable.Index.Count == 0 && variable.Child == null)
        {
            return "\"\"";
        }

        string head;
        if (variable.IsCall && RenderBuiltin(variable) is { } builtin)
        {
            head = builtin;
        }
        else if (variable.IsCall)
        {
            head = (Resolve(variable.Name) ?? variable.Name) +
                   $"({RenderArguments(variable.Args)})";
        }
        else
        {
            head = Resolve(variable.Name) ?? variable.Name;
        }

        head += RenderIndexes(variable.Index);
        for (VariableLiteral? child = variable.Child; child != null; child = child.Child)
        {
            head += "." + RenderMember(child);
        }
        return head;
    }

    private string RenderMember(VariableLiteral member)
    {
        string value;
        if (member.IsCall)
        {
            value = $"{Pascal(member.Name)}({RenderArguments(member.Args)})";
        }
        else if (MemberAccessors.TryGetValue(member.Name, out string? accessor))
        {
            value = accessor;
        }
        else
        {
            value = scope?.ResolveMember(member.Name) ?? Pascal(member.Name);
        }
        return value + RenderIndexes(member.Index);
    }

    private string RenderIndexes(IReadOnlyList<Term> indexes) =>
        string.Concat(indexes.Select(index => $"[{Render(index)}]"));

    private string RenderArguments(IReadOnlyList<Term>? arguments) =>
        string.Join(", ", (arguments ?? []).Select(Render));

    private string? RenderBuiltin(VariableLiteral variable)
    {
        if (variable.Args is not { } arguments)
        {
            return null;
        }
        return variable.Name switch
        {
            "COUNT" when arguments.Count == 1 => RenderCount(arguments[0]),
            "CEIL" when arguments.Count == 1 =>
                $"(int) System.Math.Ceiling((double) ({Render(arguments[0])}))",
            "ARRAY_SIZE_IN_BYTES" when arguments.Count == 1 =>
                $"{StaticHelperClass}.ArraySizeInBytes({Render(arguments[0])})",
            "STATIC_CALL" when arguments.Count >= 1 && arguments[0] is StringLiteral function =>
                $"{StaticHelperClass}.{Pascal(function.Value)}({string.Join(", ", arguments.Skip(1).Select(Render))})",
            "CAST" when arguments.Count == 2 && RenderCastType(arguments[1]) is { } type =>
                $"(({type}) {Render(arguments[0])})",
            _ => null
        };
    }

    private string RenderCount(Term argument)
    {
        if (argument is VariableLiteral { Child: null, Args: null } variable &&
            scope?.ResolveCount(variable.Name) is { } count)
        {
            return count;
        }
        return Render(argument) + ".Count";
    }

    private static string? RenderCastType(Term term) => term switch
    {
        StringLiteral literal => literal.Value,
        VariableLiteral variable => RenderTypeName(variable),
        _ => null
    };

    private static string? RenderTypeName(VariableLiteral variable)
    {
        var parts = new List<string>();
        for (VariableLiteral? part = variable; part != null; part = part.Child)
        {
            if (part.IsCall || part.Index.Count != 0)
            {
                return null;
            }
            parts.Add(part.Name);
        }
        return string.Join(".", parts);
    }

    private string? Resolve(string name) => scope?.ResolveReference(name);

    private static string Pascal(string name) => string.IsNullOrEmpty(name)
        ? name
        : char.ToUpperInvariant(name[0]) + name[1..];
}
