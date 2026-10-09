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

namespace org.apache.plc4net.tools.codegen.grammar
{
    /// <summary>
    /// The mspec lexer's <c>EmptyLine</c> rule starts with the semantic predicate
    /// <c>{getCharPositionInLine() == 0}?</c> in the grammar under
    /// <c>code-generation/protocol-base-mspec</c>. ANTLR copies predicate text into
    /// the generated code unchanged, so this part of the (partial) lexer supplies
    /// that method - the Java runtime's name for the C# runtime's <c>Column</c> -
    /// and the grammar can be used exactly as it is upstream.
    /// </summary>
    public partial class MSpecLexer
    {
        private int getCharPositionInLine() => Column;
    }
}