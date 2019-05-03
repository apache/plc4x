package org.apache.plc4x.codegen;

import com.fasterxml.jackson.annotation.JsonTypeInfo;

@JsonTypeInfo(
    use = JsonTypeInfo.Id.NAME,
    property = "_type")
public abstract class Node {

    public abstract <T> T accept(NodeVisitor<T> visitor);
}
