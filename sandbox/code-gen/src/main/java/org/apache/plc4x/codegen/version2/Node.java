package org.apache.plc4x.codegen.version2;

public interface Node {

    <T> T accept(NodeVisitor<T> visitor);

    void write(Generator writer);

}
