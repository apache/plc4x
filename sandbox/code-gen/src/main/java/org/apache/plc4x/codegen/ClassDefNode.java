package org.apache.plc4x.codegen;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

import java.util.ArrayList;
import java.util.List;

@JsonTypeName("ClassDef")
@JsonIgnoreProperties(ignoreUnknown = true)
public class ClassDefNode extends LineEntryNode {

    @JsonProperty("bases")
    private List<Node> bases = new ArrayList<>();

    @JsonProperty("body")
    private List<Node> body = new ArrayList<>();

    @JsonProperty("name")
    private String name;

    @JsonCreator
    public ClassDefNode() {
    }

    public List<Node> getBases() {
        return bases;
    }

    public void setBases(List<Node> bases) {
        this.bases = bases;
    }

    public List<Node> getBody() {
        return body;
    }

    public void setBody(List<Node> body) {
        this.body = body;
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
