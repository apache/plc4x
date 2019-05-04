package org.apache.plc4x.codegen.version2;

import java.util.List;

/**
 * Defines a File in Java
 */
public class FileNode implements Node {

    private final ClassDefinition mainClass;
    private final List<ClassDefinition> innerClasses;

    public FileNode(ClassDefinition mainClass, List<ClassDefinition> innerClasses) {
        this.mainClass = mainClass;
        this.innerClasses = innerClasses;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generateFile(mainClass, innerClasses);
    }
}
