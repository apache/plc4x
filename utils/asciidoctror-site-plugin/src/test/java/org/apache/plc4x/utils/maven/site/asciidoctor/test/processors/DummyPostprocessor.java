package org.apache.plc4x.utils.maven.site.asciidoctor.test.processors;

import java.util.Map;

import org.asciidoctor.ast.Document;
import org.asciidoctor.extension.Postprocessor;

public class DummyPostprocessor extends Postprocessor {

    public DummyPostprocessor(Map<String, Object> config) {
        super(config);
        System.out.println(this.getClass().getSimpleName() + "(" 
                + this.getClass().getSuperclass().getSimpleName() + ") initialized");
    }

    @Override
    public String process(Document document, String output) {
        System.out.println("Processing "+ this.getClass().getSimpleName());
        System.out.println("Processing: blocks found: " + document.getBlocks().size());
        System.out.println("Processing: output size: " + output.length());
        return output;
    }

}
