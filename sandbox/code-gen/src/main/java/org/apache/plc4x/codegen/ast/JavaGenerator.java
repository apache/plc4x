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
package org.apache.plc4x.codegen.ast;

import java.util.List;
import java.util.Set;

public class JavaGenerator implements Generator {

    public static final String PUBLIC_ = "public ";
    public static final String STATIC_ = "static ";
    private final CodeWriter writer;

    public JavaGenerator(CodeWriter writer) {
        this.writer = writer;
    }

    @Override
    public Node prepare(Node root) {
        return root;
    }

    @Override
    public void generate(ConstantExpression constantExpression) {
        writer.write(constantExpression.getValue().toString());
    }

    @Override
    public void generateDeclarationWithInitializer(DeclarationStatement declarationStatement) {
        declarationStatement.getParameterExpression().getType().write(this);
        writer.write(" ");
        declarationStatement.getParameterExpression().write(this);
        writer.write(" = ");
        declarationStatement.getInitializer().write(this);
    }

    @Override
    public void generateDeclaration(DeclarationStatement declarationStatement) {
        declarationStatement.getParameterExpression().getType().write(this);
        writer.write(" ");
        declarationStatement.getParameterExpression().write(this);
    }

    @Override
    public void generate(ParameterExpression parameterExpression) {
        writer.write(parameterExpression.getName());
    }

    @Override
    public void generatePrimitive(Primitive.DataType primitive) {
        switch (primitive) {
            case STRING:
                writer.write("String");
                break;
            case DOUBLE:
                writer.write("Double");
                break;
            case LONG:
                writer.write("Long");
                break;
            case INTEGER:
                writer.write("Integer");
                break;
            case BOOLEAN:
                writer.write("Boolean");
                break;
            case VOID:
                writer.write("Void");
                break;
            default:
                throw new UnsupportedOperationException("The primitive type " + primitive + " is not implemented!");
        }
    }

    @Override
    public void generate(IfStatement ifStatement) {
        writer.startLine("if (");
        ifStatement.getConditions().get(0).write(this);
        writer.write(") {\n");
        writeBlock(ifStatement.getBlocks().get(0));
        // For each remaining condition to an else if
        for (int i = 1; i < ifStatement.getConditions().size(); i++) {
            writer.startLine("} else if (");
            ifStatement.getConditions().get(i).write(this);
            writer.write(") {\n");
            writeBlock(ifStatement.getBlocks().get(i));
        }
        if (ifStatement.getBlocks().size() == ifStatement.getConditions().size() + 1) {
            writer.writeLine("} else {");
            writeBlock(ifStatement.getBlocks().get(ifStatement.getBlocks().size() - 1));
        }
        writer.writeLine("}");
    }

    @Override
    public void writeBlock(Block statements) {
        if (statements == null) {
            return;
        }
        writer.startBlock();
        for (Node statement : statements.getStatements()) {
            // Dont to the wrapping for If Statements
            if (statement instanceof IfStatement) {
                statement.write(this);
            } else {
                writer.startLine("");
                statement.write(this);
                writer.write(";");
                writer.endLine();
            }
        }
        writer.endBlock();
    }

    @Override
    public void generate(BinaryExpression binaryExpression) {
        binaryExpression.getLeft().write(this);
        writer.write(" ");
        writer.write(getOperator(binaryExpression.getOp()));
        writer.write(" ");
        binaryExpression.getRight().write(this);
    }

    @Override
    public void generate(AssignementExpression assignementExpression) {
        assignementExpression.getTarget().write(this);
        writer.write(" = ");
        assignementExpression.getValue().write(this);
    }

    @Override
    public void generateStaticCall(Method method, List<Node> arguments) {
        writer.write(method.getType().getTypeString());
        writer.write(".");
        writer.write(method.getName());
        writer.write("(");
        generateArgumentList(arguments);
        writer.write(")");
    }

    private void generateArgumentList(List<Node> arguments) {
        for (int i = 0; i < arguments.size(); i++) {
            arguments.get(i).write(this);
            if (i < arguments.size() - 1) {
                writer.write(", ");
            }
        }
    }

    @Override
    public void generateCall(Node target, Method method, List<Node> arguments) {
        target.write(this);
        writer.write(".");
        writer.write(method.getName());
        writer.write("(");
        generateArgumentList(arguments);
        writer.write(")");
    }

    @Override
    public void generate(NewExpression newExpression) {
        writer.write("new ");
        newExpression.getType().write(this);
        writer.write("(");
        generateArgumentList(newExpression.getArguments());
        writer.write(")");
    }

    @Override
    public void generate(MethodDefinition methodDefinition) {
        writer.startLine(PUBLIC_);
        if (methodDefinition.getModifiers().contains(Modifier.STATIC)) {
            writer.write(STATIC_);
        }
        // Special handling of VOID is necessary, to avoid having to return null in the end.
        if (methodDefinition.getResultType() instanceof Primitive &&
            ((Primitive) methodDefinition.getResultType()).getType() == Primitive.DataType.VOID) {
            writer.write("void");
        } else {
            methodDefinition.getResultType().write(this);
        }
        writer.write(" ");
        writer.write(methodDefinition.getName());
        writer.write("(");
        for (int i = 0; i < methodDefinition.getParameters().size(); i++) {
            methodDefinition.getParameters().get(i).getType().write(this);
            writer.write(" ");
            methodDefinition.getParameters().get(i).write(this);
            if (i < methodDefinition.getParameters().size() - 1) {
                writer.write(", ");
            }
        }
        writer.write(") {");
        writer.endLine();
        methodDefinition.getBody().write(this);
        writer.writeLine("}");
    }

    @Override
    public void generateReturn(Expression value) {
        writer.write("return ");
        value.write(this);
    }

    @Override
    public void generateClass(String namespace, String className, List<FieldDeclaration> fields, List<ConstructorDeclaration> constructors, List<MethodDefinition> methods, List<ClassDeclaration> innerClasses, boolean mainClass) {
        // Add static?!
        // Own File?
        writer.startLine(PUBLIC_);
        if (!mainClass) {
            writer.write(STATIC_);
        }
        writer.write("class ");
        writer.write(className);
        // TODO extends / implements
        writer.write(" {");
        writer.endLine();
        writer.startBlock();

        writer.writeLine("");
        // Fields
        for (FieldDeclaration field : fields) {
            field.write(this);
            writer.writeLine("");
        }

        // Constructors
        if (constructors != null) {
            for (ConstructorDeclaration constructor : constructors) {
                this.generateConstructor(constructor.getModifiers(), className, constructor.getParameters(), constructor.getBody());
                writer.writeLine("");
            }
        }

        // Methods
        for (MethodDefinition method : methods) {
            method.write(this);
            writer.writeLine("");
        }

        // If there are inner classes, implement them
        if (innerClasses != null) {
            for (ClassDeclaration innerClass : innerClasses) {
                this.generateClass(innerClass.getNamespace(), innerClass.getClassName(), innerClass.getFields(), innerClass.getConstructors(), innerClass.getMethods(), innerClass.getInnerClasses(), false);
            }
        }

        writer.endBlock();
        writer.writeLine("}");
    }

    @Override
    public void generateFieldDeclaration(Set<Modifier> modifiers, TypeDefinition type, String name, Expression initializer) {
        if (modifiers.contains(Modifier.PRIVATE)) {
            writer.startLine("private ");
        } else {
            writer.startLine(PUBLIC_);
        }
        if (modifiers.contains(Modifier.STATIC)) {
            writer.write(STATIC_);
        }
        if (modifiers.contains(Modifier.FINAL)) {
            writer.write("final ");
        }
        type.write(this);
        writer.write(" ");
        writer.write(name);
        // If it has an initializer, then do it here...
        if (initializer != null) {
            writer.write( " = ");
            initializer.write(this);
        }
        writer.write(";");
        writer.endLine();
    }

    @Override
    public void generateFieldReference(TypeDefinition type, String name) {
        writer.write("this.");
        writer.write(name);
    }

    @Override
    public void generateConstructor(Set<Modifier> modifiers, String className, List<ParameterExpression> parameters, Block body) {
        if (modifiers.contains(Modifier.PRIVATE)) {
            writer.startLine("private ");
        } else {
            writer.startLine(PUBLIC_);
        }
        writer.write(className);
        writer.write("(");
        for (int i = 0; i < parameters.size(); i++) {
            parameters.get(i).getType().write(this);
            writer.write(" ");
            parameters.get(i).write(this);
            if (i < parameters.size() - 1) {
                writer.write(", ");
            }
        }
        writer.write(") {");
        writer.endLine();
        body.write(this);
        writer.writeLine("}");
    }

    @Override
    public void generateFile(ClassDeclaration mainClass, List<ClassDeclaration> innerClasses) {
        generateClass(mainClass.getNamespace(), mainClass.getClassName(), mainClass.getFields(), mainClass.getConstructors(), mainClass.getMethods(), innerClasses, true);
    }

    @Override
    public void generateType(String typeString) {
        writer.write(typeString);
    }

    @Override
    public void generateComment(String comment) {
        writer.writeLine("// " + comment);
    }

    @Override
    public void generateNoOp() {
        writer.write(";");
    }

    private String getOperator(BinaryExpression.Operation op) {
        switch (op) {
            case EQ:
                return "==";
            case PLUS:
                return "+";
        }
        throw new UnsupportedOperationException("The Operator " + op + " is currently not implemented!");
    }
}
