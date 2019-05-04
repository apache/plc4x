package org.apache.plc4x.codegen.version2;

/**
 * Stub for the Type System.
 */
public abstract class TypeNode implements Node {

    private final String typeString;

    protected TypeNode(String typeString) {
        this.typeString = typeString;
    }

    String getTypeString() {
        return this.typeString;
    }
}
