package org.apache.plc4x.codegen.ast;

public interface Node {

    <T> T accept(NodeVisitor<T> visitor);

    void write(Generator writer);

}
