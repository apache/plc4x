package org.apache.plc4x.utils.maven.site.asciidoctor.test.processors;

import java.util.Map;

import org.asciidoctor.ast.Document;
import org.asciidoctor.extension.DocinfoProcessor;

public class MetaDocinfoProcessor extends DocinfoProcessor {

    public MetaDocinfoProcessor(Map<String, Object> config) {
        super(config);
        System.out.println(this.getClass().getSimpleName() + "(" 
                + this.getClass().getSuperclass().getSimpleName() + ") initialized");
    }

    @Override
    public String process(Document document) {
        System.out.println("Processing "+ this.getClass().getSimpleName());
        System.out.println("Processing: blocks found: " + document.getBlocks().size());
        return "<meta name=\"author\" content=\"asciidoctor\">";
    }

}
