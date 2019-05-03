package org.apache.plc4x.codegen;

import com.fasterxml.jackson.annotation.JsonProperty;

public abstract class LineEntryNode extends Node {

    @JsonProperty("col_offset")
    private long colOffset;

    @JsonProperty("lineno")
    private long lineno;

    public long getColOffset() {
        return colOffset;
    }

    public void setColOffset(long colOffset) {
        this.colOffset = colOffset;
    }

    public long getLineno() {
        return lineno;
    }

    public void setLineno(long lineno) {
        this.lineno = lineno;
    }
}
