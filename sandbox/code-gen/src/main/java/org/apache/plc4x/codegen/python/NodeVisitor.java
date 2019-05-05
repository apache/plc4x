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
