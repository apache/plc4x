/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.codegen;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.plc4x.codegen.python.*;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.File;
import java.io.IOException;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertTrue;

class ModuleNodeTest {

    private ObjectMapper mapper;

    @BeforeEach
    void setUp() {
        mapper = new ObjectMapper();
        mapper.registerSubtypes(ModuleNode.class);
        mapper.registerSubtypes(ClassDefNode.class);
        mapper.registerSubtypes(PassNode.class);
        mapper.registerSubtypes(FunctionDefNode.class);
        mapper.registerSubtypes(ArgumentsNode.class);
        mapper.registerSubtypes(ExprNode.class);
        mapper.registerSubtypes(CallNode.class);
        mapper.registerSubtypes(StrNode.class);
        mapper.registerSubtypes(NameNode.class);
        mapper.registerSubtypes(LoadNode.class);
        mapper.registerSubtypes(AnnAssignerNode.class);
        mapper.registerSubtypes(StoreNode.class);
        mapper.registerSubtypes(NumNode.class);
        mapper.registerSubtypes(AugAssignNode.class);
        mapper.registerSubtypes(AddNode.class);
        mapper.registerSubtypes(ImportFromNode.class);
        mapper.registerSubtypes(AliasNode.class);
        mapper.registerSubtypes(ArgNode.class);
        mapper.registerSubtypes(AssignNode.class);
        mapper.registerSubtypes(AttributeNode.class);
        mapper.registerSubtypes(WhileNode.class);
        mapper.registerSubtypes(CompareNode.class);
        mapper.registerSubtypes(NotEqNode.class);
        mapper.registerSubtypes(ReturnNode.class);
        mapper.registerSubtypes(IfNode.class);
        mapper.registerSubtypes(KeywordNode.class);
        mapper.registerSubtypes(BinOpNode.class);
        mapper.registerSubtypes(ListNode.class);
        mapper.registerSubtypes(MultNode.class);
        mapper.registerSubtypes(NameConstantNode.class);
        mapper.registerSubtypes(TupleNode.class);
        mapper.registerSubtypes(EqNode.class);
    }

    @Test
    void serialize() throws IOException {
        final ModuleNode node = new ModuleNode();
        node.setBody(Collections.singletonList(new ClassDefNode()));

        final String s = mapper.writeValueAsString(node);

        System.out.println(s);

        // Reread
        final Node root = mapper.readValue(s, Node.class);

        System.out.println(root);
    }

    @Test
    void deserializeExample() throws IOException {
        final Node root = mapper.readValue(new File("src/main/resources/example.json"), Node.class);

        assertTrue(root instanceof ModuleNode);
        assertTrue(((ModuleNode) root).getBody().get(0) instanceof ClassDefNode);

        // Print the node again
        System.out.println(mapper.writeValueAsString(root));
    }

    @Test
    void deserializeExample2() throws IOException {
        final Node root = mapper.readValue(new File("src/main/resources/example2.json"), Node.class);

        assertTrue(root instanceof ModuleNode);

        // Print the node again
        System.out.println(mapper.writeValueAsString(root));
    }

    @Test
    void visitMinimalTree() throws IOException {
        final Node root = mapper.readValue(new File("src/main/resources/example2.json"), Node.class);

        CodePrinter printer = new CodePrinter(4);

        NodeVisitor<String> toString = new AbstractNodeVisitor<String>() {


            @Override
            public String visit(CompareNode compareNode) {
                assert compareNode.getComparators().size() == 1;
                return compareNode.getLeft().accept(this) + " " + compareNode.getOps().get(0).accept(this) + " " + compareNode.getComparators().get(0).accept(this);
            }

            @Override
            public String visit(EqNode eqNode) {
                return "==";
            }

            @Override
            public String visit(NotEqNode notEqNode) {
                return "!=";
            }

            @Override
            public String visit(NumNode numNode) {
                return Double.toString(numNode.getN());
            }

            @Override
            public String visit(StrNode strNode) {
                return strNode.getS();
            }

            @Override
            public String visit(NameNode nameNode) {
                return nameNode.getId();
            }

            @Override
            public String visit(AttributeNode attributeNode) {
                return attributeNode.getAttr();
            }

            @Override
            public String visit(CallNode callNode) {
                final String function = callNode.getFunc().accept(this);

                String s = function + "(";
                for (Node arg : callNode.getArgs()) {
                    s += arg.accept(this) + ", ";
                }
                for (Node keyword : callNode.getKeywords()) {
                    s += keyword.accept(this);
                }
                s += ")";
                return s;
            }

            @Override
            public String visit(KeywordNode keywordNode) {
                return keywordNode.getArg() + " = " + keywordNode.getValue().accept(this);
            }

            @Override
            public String visit(AssignNode assignNode) {
                assert assignNode.getTargets().size() == 1;
                printer.writeLine(String.format("%s = %s;", assignNode.getTargets().get(0).accept(this), assignNode.getValue().accept(this))); ;
                return super.visit(assignNode);
            }

            @Override
            public String visit(IfNode ifNode) {
                printer.writeLine(String.format("If (%s) {", ifNode.getTest().accept(this)));
                printer.startBlock();
                for (Node node : ifNode.getBody()) {
                    node.accept(this);
                }
                printer.endBlock();
                printer.writeLine("}");
                return null;
            }
        };


        root.accept(toString);

        System.out.println(printer.getCode());
    }
}