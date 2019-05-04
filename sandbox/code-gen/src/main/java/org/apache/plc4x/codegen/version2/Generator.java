package org.apache.plc4x.codegen.version2;

import java.util.List;

public interface Generator {

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

    void generateFieldDeclaration(TypeNode type, String name);

    void generateFieldReference(TypeNode type, String name);

    void generateConstructor(String className, List<ParameterExpression> parameters, Block body);

    void generateFile(ClassDefinition mainClass, List<ClassDefinition> innerClasses);
}
