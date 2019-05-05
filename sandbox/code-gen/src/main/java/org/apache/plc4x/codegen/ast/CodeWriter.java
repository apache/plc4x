package org.apache.plc4x.codegen.ast;

import java.util.stream.IntStream;

/**
 * Helper class to print code.
 */
public class CodeWriter {

    private StringBuffer buffer = new StringBuffer();
    private int tabSize;

    private int intendationLvl = 0;

    public CodeWriter(int tabSize) {
        this.tabSize = tabSize;
    }

    public void startBlock() {
        this.intendationLvl += tabSize;
    }

    public void endBlock() {
        if (intendationLvl < tabSize) {
            throw new RuntimeException("Closing a Block which is not open!");
        }
        this.intendationLvl -= tabSize;
    }

    public void write(String s) {
        buffer.append(s);
    }

    public void startLine(String s) {
        writeIntendation();
        write(s);
    }

    public void endLine() {
        buffer.append("\n");
    }

    public void writeLine(String s) {
        writeIntendation();
        buffer.append(s);
        buffer.append("\n");
    }

    private void writeIntendation() {
        // Write the intendation
        IntStream.range(0, intendationLvl).forEach(i -> buffer.append(" "));
    }

    public String getCode() {
        return buffer.toString();
    }
}
