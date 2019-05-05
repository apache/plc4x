package org.apache.plc4x.codegen.ast;

import java.util.Arrays;
import java.util.Collections;
import java.util.HashSet;
import java.util.Set;

public class FieldDeclaration implements Node {

    private final Set<Modifier> modifiers;
    private final TypeNode type;
    private final String name;

    private final Expression initializer;

    public FieldDeclaration(Set<Modifier> modifiers, TypeNode type, String name, Expression initializer) {
        this.modifiers = modifiers;
        this.type = type;
        this.name = name;
        this.initializer = initializer;
    }

    public FieldDeclaration(TypeNode type, String name) {
        this(Collections.emptySet(), type, name, null);
    }

    public FieldDeclaration(TypeNode type, String name, Modifier... modifiers) {
        this(new HashSet<>(Arrays.asList(modifiers)), type, name, null);
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generateFieldDeclaration(modifiers, type, name, initializer);
    }
}
