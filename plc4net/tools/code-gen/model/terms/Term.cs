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
using System.Globalization;
using System.Linq;

namespace org.apache.plc4net.tools.codegen.model.terms;

/// <summary>A node in an MSpec expression tree.</summary>
public abstract class Term
{
}

public sealed class IntegerLiteral(long value, string? text = null) : Term
{
    public long Value { get; } = value;
    public string Text { get; } = text ?? value.ToString(CultureInfo.InvariantCulture);
    public override string ToString() => Text;
}

public sealed class HexadecimalLiteral(string text) : Term
{
    public string Text { get; } = text;
    public override string ToString() => Text;
}

public sealed class FloatLiteral(double value) : Term
{
    public double Value { get; } = value;
    public override string ToString() => Value.ToString("R", CultureInfo.InvariantCulture);
}

public sealed class StringLiteral(string value) : Term
{
    public string Value { get; } = value;
    public override string ToString() =>
        $"\"{Value.Replace("\\", "\\\\").Replace("\"", "\\\"")}\"";
}

public sealed class BooleanLiteral(bool value) : Term
{
    public bool Value { get; } = value;
    public override string ToString() => Value ? "true" : "false";
}

public sealed class NullLiteral : Term
{
    public override string ToString() => "null";
}

public sealed class VariableLiteral(
    string name,
    IReadOnlyList<Term>? args = null,
    IReadOnlyList<Term>? index = null,
    VariableLiteral? child = null) : Term
{
    public string Name { get; } = name;
    public IReadOnlyList<Term>? Args { get; } = args;
    public IReadOnlyList<Term> Index { get; } = index ?? [];
    public VariableLiteral? Child { get; } = child;
    public bool IsCall => Args != null;

    public override string ToString()
    {
        string value = Name;
        if (Args != null)
        {
            value += $"({string.Join(", ", Args)})";
        }
        value += string.Concat(Index.Select(item => $"[{item}]"));
        return Child == null ? value : $"{value}.{Child}";
    }
}

public sealed class BinaryExpression(Term left, string op, Term right) : Term
{
    public Term Left { get; } = left;
    public string Operator { get; } = op;
    public Term Right { get; } = right;
    public override string ToString() => $"({Left} {Operator} {Right})";
}

public sealed class UnaryExpression(string op, Term operand) : Term
{
    public string Operator { get; } = op;
    public Term Operand { get; } = operand;
    public override string ToString() => $"{Operator}{Operand}";
}

public sealed class IndexExpression(Term target, IReadOnlyList<Term> indexes) : Term
{
    public Term Target { get; } = target;
    public IReadOnlyList<Term> Indexes { get; } = indexes;
    public override string ToString() =>
        (Target is UnaryExpression ? $"({Target})" : Target.ToString()) +
        string.Concat(Indexes.Select(item => $"[{item}]"));
}

public sealed class TernaryExpression(Term condition, Term whenTrue, Term whenFalse) : Term
{
    public Term Condition { get; } = condition;
    public Term WhenTrue { get; } = whenTrue;
    public Term WhenFalse { get; } = whenFalse;
    public override string ToString() => $"({Condition} ? {WhenTrue} : {WhenFalse})";
}
