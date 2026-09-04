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
using org.apache.plc4net.tools.codegen.model.terms;

namespace org.apache.plc4net.tools.codegen.model
{
    /// <summary>
    /// The whole type model for one mspec file: the complex types (including
    /// the children lifted out of every <c>typeSwitch</c>), the enums, the
    /// <c>dataIo</c> types, and the top-level <c>constants</c> block.
    /// </summary>
    public sealed class Protocol
    {
        public List<ComplexTypeDefinition> Types { get; } = new List<ComplexTypeDefinition>();

        public List<EnumTypeDefinition> Enums { get; } = new List<EnumTypeDefinition>();

        public List<DataIoTypeDefinition> DataIos { get; } = new List<DataIoTypeDefinition>();

        public List<ConstantDeclaration> Constants { get; } = new List<ConstantDeclaration>();

        public ComplexTypeDefinition FindType(string name)
            => Types.FirstOrDefault(t => t.Name == name);

        public EnumTypeDefinition FindEnum(string name)
            => Enums.FirstOrDefault(e => e.Name == name);

        /// <summary>The concrete children of a discriminated parent, in
        /// declaration order.</summary>
        public IEnumerable<ComplexTypeDefinition> ChildrenOf(string parentName)
            => Types.Where(t => t.ParentName == parentName);
    }

    /// <summary>One entry in the top-level <c>[constants]</c> block.</summary>
    public sealed class ConstantDeclaration
    {
        public string Name { get; set; }
        public TypeReference Type { get; set; }
        public Term Value { get; set; }
    }
}
