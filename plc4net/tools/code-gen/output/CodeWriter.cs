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

using System.Text;

namespace org.apache.plc4net.tools.codegen.output
{
    /// <summary>A minimal indent-aware source writer.</summary>
    internal sealed class CodeWriter
    {
        private readonly StringBuilder _sb = new StringBuilder();
        private int _depth;

        public void Indent() => _depth++;

        public void Outdent() => _depth = _depth > 0 ? _depth - 1 : 0;

        public void Line(string text = "")
        {
            if (text.Length == 0)
            {
                _sb.Append('\n');
                return;
            }
            _sb.Append(new string(' ', _depth * 4)).Append(text).Append('\n');
        }

        public override string ToString() => _sb.ToString();
    }
}
