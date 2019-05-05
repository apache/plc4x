package org.apache.plc4x.codegen.ast;

import java.util.List;
import java.util.Set;

public interface Generator {

    /**
     * Do preliminary stuff.
     * @param root
     * @return
     */
    Node prepare(Node root);

    void generate(ConstantNode constantNode);

    void generateDeclarationWithInitializer(DeclarationStatement declarationStatement);

    void generateDeclaration(DeclarationStatement declarationStatement);

    void generate(ParameterExpression parameterExpression);

    void generate(Primitive primitive);

    void generate(IfStatement ifStatement);

    void writeBlock(Block statements);

    void generate(BinaryExpression binaryExpression);

    void generate(AssignementExpression assignementExpression);

    void generateStaticCall(Method method, List<Node> constantNode);

    void generateCall(Node target, Method method, List<Node> constantNode);

    void generate(NewExpression newExpression);

    void generate(MethodDefinition methodDefinition);

    void generateReturn(Expression value);

    void generateClass(String namespace, String className, List<FieldDeclaration> fields, List<ConstructorDeclaration> constructors, List<MethodDefinition> methods, List<ClassDefinition> innerClasses, boolean mainClass);

    void generateFieldDeclaration(Set<Modifier> modifiers, TypeNode type, String name, Expression initializer);

    void generateFieldReference(TypeNode type, String name);

    void generateConstructor(Set<Modifier> modifiers, String className, List<ParameterExpression> parameters, Block body);

    void generateFile(ClassDefinition mainClass, List<ClassDefinition> innerClasses);

    void generateType(String typeString);

    void generateComment(String comment);

    void generateNoOp();
}
