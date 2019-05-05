package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

import java.util.ArrayList;
import java.util.List;

@JsonTypeName("ImportFrom")
public class ImportFromNode extends LineEntryNode {

    @JsonProperty("level")
    private long level;

    @JsonProperty("module")
    private String module;

    @JsonProperty("names")
    private List<Node> names = new ArrayList<>();

    public long getLevel() {
        return level;
    }

    public void setLevel(long level) {
        this.level = level;
    }

    public String getModule() {
        return module;
    }

    public void setModule(String module) {
        this.module = module;
    }

    public List<Node> getNames() {
        return names;
    }

    public void setNames(List<Node> names) {
        this.names = names;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
