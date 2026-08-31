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

using System.IO;
using Antlr4.Runtime;
using org.apache.plc4net.tools.codegen.grammar;

namespace org.apache.plc4net.tools.codegen
{
    /// <summary>
    /// Parses an mspec protocol description file into an ANTLR parse tree.
    /// This is the first stage of the pure-.NET mspec toolchain: the grammar
    /// layer that turns text into a parse tree. The type model layer (the
    /// next stage) walks the tree with a listener.
    /// </summary>
    public static class MspecReader
    {
        /// <summary>
        /// Parses the content of an mspec file and returns the parse tree root.
        /// </summary>
        public static MSpecParser.FileContext Read(string mspecContent, bool strict = true)
        {
            var inputStream = new AntlrInputStream(mspecContent);
            var lexer = new MSpecLexer(inputStream);
            var tokenStream = new CommonTokenStream(lexer);
            var parser = new MSpecParser(tokenStream);

            // Collect errors rather than printing to stderr.
            var errorListener = new CollectingErrorListener();
            parser.RemoveErrorListeners();
            parser.AddErrorListener(errorListener);

            var tree = parser.file();

            if (errorListener.HasErrors && strict)
            {
                throw new MspecParseException(
                    $"Failed to parse mspec: {string.Join("; ", errorListener.Errors)}");
            }

            // Non-strict: ANTLR's error recovery has already skipped the bad
            // tokens and produced a best-effort tree. The Java MessageFormatParser
            // behaves the same on knx-master-data.mspec (it logs the same
            // "mismatched input" lines to stderr and carries on).
            LastErrors = errorListener.Errors;
            return tree;
        }

        /// <summary>Syntax errors from the most recent non-strict
        /// <see cref="Read(string, bool)"/>, for a caller that wants to report them.</summary>
        public static System.Collections.Generic.IReadOnlyList<string> LastErrors { get; private set; }
            = System.Array.Empty<string>();

        /// <summary>
        /// Parses an mspec file at the given path.
        /// </summary>
        public static MSpecParser.FileContext ReadFile(string filePath, bool strict = true)
        {
            return Read(File.ReadAllText(filePath), strict);
        }
    }

    /// <summary>
    /// Collects syntax errors rather than printing them.
    /// </summary>
    internal class CollectingErrorListener : BaseErrorListener
    {
        private readonly System.Collections.Generic.List<string> _errors
            = new System.Collections.Generic.List<string>();

        public System.Collections.Generic.IReadOnlyList<string> Errors => _errors;

        public bool HasErrors => _errors.Count > 0;

        public override void SyntaxError(TextWriter output, IRecognizer recognizer,
            IToken offendingSymbol, int line, int charPositionInLine,
            string msg, RecognitionException e)
        {
            _errors.Add($"line {line}:{charPositionInLine} {msg}");
        }
    }

    /// <summary>
    /// Raised when an mspec file fails to parse.
    /// </summary>
    public class MspecParseException : System.Exception
    {
        public MspecParseException(string message) : base(message) { }
        public MspecParseException(string message, System.Exception inner) : base(message, inner) { }
    }
}
