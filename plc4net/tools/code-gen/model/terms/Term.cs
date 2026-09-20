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

namespace org.apache.plc4net.tools.codegen.model.terms
{
    /// <summary>
    /// A node in an mspec expression tree.
    /// </summary>
    /// <remarks>
    /// mspec expressions appear inside single quotes - array lengths
    /// (<c>count 'byteCount'</c>), implicit-field values
    /// (<c>'pdu.lengthInBytes + 1'</c>), discriminator literals, virtual-field
    /// formulas. <see cref="MspecExpressionParser"/> turns the quoted text into
    /// one of these; a language renderer (for C#,
    /// <see cref="output.CSharpExpressionRenderer"/>) turns it back into source.
    ///
    /// The shape mirrors the Java <c>Term</c> / <c>Literal</c> hierarchy in
    /// plc4x-code-generation so the same mspec means the same thing in both.
    /// <c>ToString()</c> renders a language-neutral, mostly C#-shaped form and
    /// is used by tests and diagnostics, not by the generator.
    /// </remarks>
    public abstract class Term
    {
    }

    /// <summary>
    /// An integer literal. <see cref="Text"/> preserves the spelling as it
    /// appeared in the mspec (e.g. <c>0x0E</c>) so generated code and
    /// diagnostics read like the source; <see cref="Value"/> is the number.
    /// </summary>
    public class IntegerLiteral : Term
    {
        public long Value { get; }

        /// <summary>The literal exactly as written, or a decimal rendering.</summary>
        public string Text { get; }

        public bool IsHex => Text.StartsWith("0x") || Text.StartsWith("0X");

        public IntegerLiteral(long value, string text = null)
        {
            Value = value;
            Text = text ?? value.ToString(CultureInfo.InvariantCulture);
        }

        public override string ToString() => Text;
    }

    /// <summary>A floating-point literal (e.g. <c>'1.5'</c>).</summary>
    public class FloatLiteral : Term
    {
        public double Value { get; }
        public FloatLiteral(double value) { Value = value; }
        public override string ToString() => Value.ToString("R", CultureInfo.InvariantCulture);
    }

    /// <summary>A string literal. The stored value has no surrounding quotes.</summary>
    public class StringLiteral : Term
    {
        public string Value { get; }
        public StringLiteral(string value) { Value = value; }
        public override string ToString() => "\"" + Value + "\"";
    }

    /// <summary>A boolean literal.</summary>
    public class BooleanLiteral : Term
    {
        public bool Value { get; }
        public BooleanLiteral(bool value) { Value = value; }
        public override string ToString() => Value ? "true" : "false";
    }

    /// <summary>The <c>null</c> literal.</summary>
    public class NullLiteral : Term
    {
        public override string ToString() => "null";
    }

    /// <summary>
    /// A name reference, optionally a call (<c>args != null</c>), optionally
    /// indexed (<c>index</c>), optionally with a member access chain
    /// (<c>child</c>). Covers the grammar's <c>identifierSegment</c>:
    /// <c>name (args)? [index]? (. child)?</c>.
    /// </summary>
    /// <example>
    /// <c>pdu.lengthInBytes</c>, <c>COUNT(value)</c>,
    /// <c>STATIC_CALL("rtuCrcCheck", address, pdu)</c>, <c>items[0].fileNumber</c>.
    /// </example>
    public class VariableLiteral : Term
    {
        public string Name { get; }

        /// <summary>Call arguments, or null when this is a plain reference.
        /// An empty list means a call with no arguments (<c>foo()</c>).</summary>
        public IReadOnlyList<Term> Args { get; }

        /// <summary>Index expressions (<c>a[i][j]</c>), or an empty list.</summary>
        public IReadOnlyList<Term> Index { get; }

        /// <summary>The next segment after a <c>.</c>, or null.</summary>
        public VariableLiteral Child { get; }

        public VariableLiteral(
            string name,
            IReadOnlyList<Term> args = null,
            IReadOnlyList<Term> index = null,
            VariableLiteral child = null)
        {
            Name = name;
            Args = args;
            Index = index ?? System.Array.Empty<Term>();
            Child = child;
        }

        public bool IsCall => Args != null;

        public override string ToString()
        {
            var s = Name;
            if (Args != null)
            {
                s += "(" + string.Join(", ", Args.Select(a => a.ToString())) + ")";
            }
            s = Index.Aggregate(s, (acc, i) => acc + "[" + i + "]");
            if (Child != null)
            {
                s += "." + Child;
            }
            return s;
        }
    }

    /// <summary>A binary operation <c>left op right</c>.</summary>
    public class BinaryExpression : Term
    {
        public Term Left { get; }
        public Term Right { get; }
        public string Operator { get; }

        public BinaryExpression(Term left, string op, Term right)
        {
            Left = left;
            Operator = op;
            Right = right;
        }

        public override string ToString() => $"({Left} {Operator} {Right})";
    }

    /// <summary>A unary operation <c>op operand</c> (<c>-x</c>, <c>!x</c>).</summary>
    public class UnaryExpression : Term
    {
        public Term Operand { get; }
        public string Operator { get; }

        public UnaryExpression(string op, Term operand)
        {
            Operator = op;
            Operand = operand;
        }

        public override string ToString() => $"{Operator}{Operand}";
    }

    /// <summary>A ternary <c>condition ? whenTrue : whenFalse</c>.</summary>
    public class TernaryExpression : Term
    {
        public Term Condition { get; }
        public Term WhenTrue { get; }
        public Term WhenFalse { get; }

        public TernaryExpression(Term condition, Term whenTrue, Term whenFalse)
        {
            Condition = condition;
            WhenTrue = whenTrue;
            WhenFalse = whenFalse;
        }

        public override string ToString() => $"({Condition} ? {WhenTrue} : {WhenFalse})";
    }
}
