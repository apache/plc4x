package org.apache.plc4x.codegen;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

@JsonTypeName("alias")
public class AliasNode extends Node {

    @JsonProperty("asname")
    private String asname;

    @JsonProperty("name")
    private String name;

    public String getAsname() {
        return asname;
    }

    public void setAsname(String asname) {
        this.asname = asname;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
