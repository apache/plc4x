package org.apache.plc4x.codegen.version2;

import java.util.List;

public class PhpGenerator implements Generator {

    private final CodeWriter writer;

    public PhpGenerator(CodeWriter writer) {
        this.writer = writer;
    }

    @Override public void generate(ConstantNode constantNode) {
        this.writer.write(constantNode.getValue().toString());
    }

    @Override public void generateDeclarationWithInitializer(DeclarationStatement declarationStatement) {
        this.writer.startLine("");
        this.writer.write("$");
        this.writer.write(declarationStatement.getParameterExpression().getName());
        this.writer.write(" = ");
        declarationStatement.getInitializer().write(this);
        this.writer.write("\n");
    }

    @Override public void generateDeclaration(DeclarationStatement declarationStatement) {

    }

    @Override public void generate(ParameterExpression parameterExpression) {
        // not needed
    }

    @Override public void generate(Primitive primitive) {
        // Not needed, php is typeless
    }

    @Override public void generate(IfStatement ifStatement) {

    }

    @Override public void writeBlock(Block statements) {

    }

    @Override public void generate(BinaryExpression binaryExpression) {

    }

    @Override public void generate(AssignementExpression assignementExpression) {

    }

    @Override public void generateStaticCall(Method method, List<Node> constantNode) {

    }

    @Override public void generateCall(Node target, Method method, List<Node> constantNode) {

    }

    @Override public void generate(NewExpression newExpression) {

    }

    @Override public void generate(MethodDefinition methodDefinition) {

    }

    @Override public void generateReturn(Expression value) {

    }

    @Override public void generateClass(String namespace, String className, List<FieldDeclaration> fields, List<ConstructorDeclaration> constructors, List<MethodDefinition> methods, List<ClassDefinition> innerClasses, boolean mainClass) {

    }

    @Override public void generateFieldDeclaration(TypeNode type, String name) {

    }

    @Override public void generateFieldReference(TypeNode type, String name) {

    }

    @Override public void generateConstructor(String className, List<ParameterExpression> parameters, Block body) {

    }

    @Override public void generateFile(ClassDefinition mainClass, List<ClassDefinition> innerClasses) {

    }
}
