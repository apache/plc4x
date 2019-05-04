package org.apache.plc4x.codegen.version2;

import java.util.List;

public class ClassDefinition implements Node {

    private final String namespace;
    private final String className;

    private final List<FieldDeclaration> fields;
    private final List<ConstructorDeclaration> constructors;
    private final List<MethodDefinition> methods;

    public ClassDefinition(String namespace, String className, List<FieldDeclaration> fields, List<ConstructorDeclaration> constructors, List<MethodDefinition> methods) {
        this.namespace = namespace;
        this.className = className;
        this.fields = fields;
        this.constructors = constructors;
        this.methods = methods;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generateClass(this.namespace, this.className, this.fields, this.constructors, this.methods);
    }
}
