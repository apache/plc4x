package org.apache.plc4x.codegen.ast;

import java.util.List;

public class ClassDefinition implements Node {

    private final String namespace;
    private final String className;

    private final List<FieldDeclaration> fields;
    private final List<ConstructorDeclaration> constructors;
    private final List<MethodDefinition> methods;
    private final List<ClassDefinition> innerClasses;

    public ClassDefinition(String namespace, String className, List<FieldDeclaration> fields, List<ConstructorDeclaration> constructors, List<MethodDefinition> methods, List<ClassDefinition> innerClasses) {
        this.namespace = namespace;
        this.className = className;
        this.fields = fields;
        this.constructors = constructors;
        this.methods = methods;
        this.innerClasses = innerClasses;
    }

    public String getNamespace() {
        return namespace;
    }

    public String getClassName() {
        return className;
    }

    public List<FieldDeclaration> getFields() {
        return fields;
    }

    public List<ConstructorDeclaration> getConstructors() {
        return constructors;
    }

    public List<MethodDefinition> getMethods() {
        return methods;
    }

    public List<ClassDefinition> getInnerClasses() {
        return innerClasses;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generateClass(this.namespace, this.className, this.fields, this.constructors, this.methods, innerClasses, true);
    }
}
