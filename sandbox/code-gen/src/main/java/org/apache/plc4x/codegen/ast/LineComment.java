package org.apache.plc4x.codegen.ast;

/**
 * Comment over one Line
 */
public class LineComment implements Node {

    private final String comment;

    public LineComment(String comment) {
        this.comment = comment;
    }

    public String getComment() {
        return comment;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generateComment(comment);
    }
}
