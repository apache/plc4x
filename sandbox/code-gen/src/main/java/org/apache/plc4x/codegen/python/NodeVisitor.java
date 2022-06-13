/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.codegen.python;

public interface NodeVisitor<T> {

    T visit(AugAssignNode augAssignNode);

    T visit(AddNode addNode);

    T visit(EqNode eqNode);

    T visit(IfNode ifNode);

    T visit(WhileNode whileNode);

    T visit(NotEqNode notEqNode);

    T visit(NumNode numNode);

    T visit(PassNode passNode);

    T visit(ReturnNode returnNode);

    T visit(StoreNode storeNode);

    T visit(StrNode strNode);

    T visit(TupleNode tupleNode);

    T visit(NameNode nameNode);

    T visit(NameConstantNode nameConstantNode);

    T visit(MultNode multNode);

    T visit(ModuleNode moduleNode);

    T visit(LoadNode loadNode);

    T visit(ListNode listNode);

    T visit(KeywordNode keywordNode);

    T visit(ImportFromNode importFromNode);

    T visit(FunctionDefNode functionDefNode);

    T visit(ExprNode exprNode);

    T visit(CompareNode compareNode);

    T visit(AliasNode aliasNode);

    T visit(AnnAssignerNode annAssignerNode);

    T visit(ArgNode argNode);

    T visit(ClassDefNode classDefNode);

    T visit(CallNode callNode);

    T visit(BinOpNode binOpNode);

    T visit(AttributeNode attributeNode);

    T visit(AssignNode assignNode);

    T visit(ArgumentsNode argumentsNode);
}
