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


public class AbstractNodeVisitor<T> implements NodeVisitor<T> {

    @Override
    public T visit(AugAssignNode augAssignNode) {
        augAssignNode.getOp().accept(this);
        augAssignNode.getTarget().accept(this);
        augAssignNode.getValue().accept(this);
        return null;
    }

    @Override
    public T visit(AddNode addNode) {
        return null;
    }

    @Override
    public T visit(EqNode eqNode) {
        return null;
    }

    @Override
    public T visit(IfNode ifNode) {
        for (Node node : ifNode.getBody()) {
            node.accept(this);
        }
        for (Node node : ifNode.getOrelse()) {
            node.accept(this);
        }
        ifNode.getTest().accept(this);
        return null;
    }

    @Override
    public T visit(WhileNode whileNode) {
        for (Node node : whileNode.getBody()) {
            node.accept(this);
        }
        for (Node node : whileNode.getOrelse()) {
            node.accept(this);
        }
        whileNode.getTest().accept(this);
        return null;
    }

    @Override
    public T visit(NotEqNode notEqNode) {
        return null;
    }

    @Override
    public T visit(NumNode numNode) {
        return null;
    }

    @Override
    public T visit(PassNode passNode) {
        return null;
    }

    @Override
    public T visit(ReturnNode returnNode) {
        returnNode.getValue().accept(this);
        return null;
    }

    @Override
    public T visit(StoreNode storeNode) {
        return null;
    }

    @Override
    public T visit(StrNode strNode) {
        return null;
    }

    @Override
    public T visit(TupleNode tupleNode) {
        for (Node elt : tupleNode.getElts()) {
            elt.accept(this);
        }
        return null;
    }

    @Override
    public T visit(NameNode nameNode) {
        return null;
    }

    @Override
    public T visit(NameConstantNode nameConstantNode) {
        return null;
    }

    @Override
    public T visit(MultNode multNode) {
        return null;
    }

    @Override
    public T visit(ModuleNode moduleNode) {
        for (Node node : moduleNode.getBody()) {
            node.accept(this);
        }
        return null;
    }

    @Override
    public T visit(LoadNode loadNode) {
        return null;
    }

    @Override
    public T visit(ListNode listNode) {
        for (Node elt : listNode.getElts()) {
            elt.accept(this);
        }
        return null;
    }

    @Override
    public T visit(KeywordNode keywordNode) {
        keywordNode.getValue().accept(this);
        return null;
    }

    @Override
    public T visit(ImportFromNode importFromNode) {
        for (Node name : importFromNode.getNames()) {
            name.accept(this);
        }
        return null;
    }

    @Override
    public T visit(FunctionDefNode functionDefNode) {
        functionDefNode.getArgs().accept(this);
        for (Node node : functionDefNode.getBody()) {
            node.accept(this);
        }
        if (functionDefNode.getReturns() != null) {
            functionDefNode.getReturns().accept(this);
        }
        return null;
    }

    @Override
    public T visit(ExprNode exprNode) {
        return exprNode.getValue().accept(this);
    }

    @Override
    public T visit(CompareNode compareNode) {
        for (Node comparator : compareNode.getComparators()) {
            comparator.accept(this);
        }
        compareNode.getLeft().accept(this);
        for (Node op : compareNode.getOps()) {
            op.accept(this);
        }
        return null;
    }

    @Override
    public T visit(AliasNode aliasNode) {
        return null;
    }

    @Override
    public T visit(AnnAssignerNode annAssignerNode) {
        annAssignerNode.getAnnotation().accept(this);
        annAssignerNode.getTarget().accept(this);
        if (annAssignerNode.getValue() != null) {
            return annAssignerNode.getValue().accept(this);
        } else {
            return null;
        }
    }

    @Override
    public T visit(ArgNode argNode) {
        if (argNode.getAnnotation() != null) {
            return argNode.getAnnotation().accept(this);
        } else {
            return null;
        }
    }

    @Override
    public T visit(ClassDefNode classDefNode) {
        for (Node basis : classDefNode.getBases()) {
            basis.accept(this);
        }
        for (Node node : classDefNode.getBody()) {
            node.accept(this);
        }
        return null;
    }

    @Override
    public T visit(CallNode callNode) {
        for (Node arg : callNode.getArgs()) {
            arg.accept(this);
        }
        for (Node keyword : callNode.getKeywords()) {
            keyword.accept(this);
        }
        return callNode.getFunc().accept(this);
    }

    @Override
    public T visit(BinOpNode binOpNode) {
        binOpNode.getLeft().accept(this);
        binOpNode.getOp().accept(this);
        binOpNode.getRight().accept(this);
        return null;
    }

    @Override
    public T visit(AttributeNode attributeNode) {
        return attributeNode.getValue().accept(this);
    }

    @Override
    public T visit(AssignNode assignNode) {
        for (Node target : assignNode.getTargets()) {
            target.accept(this);
        }
        return assignNode.getValue().accept(this);
    }

    @Override
    public T visit(ArgumentsNode argumentsNode) {
        for (Node arg : argumentsNode.getArgs()) {
            arg.accept(this);
        }
        return null;
    }
}
