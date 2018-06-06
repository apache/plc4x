package org.apache.plc4x.utils.maven.site.asciidoctor.test.processors;

import java.util.HashMap;
import java.util.Map;

import org.asciidoctor.ast.Document;
import org.asciidoctor.extension.Treeprocessor;

public class DummyTreeprocessor extends Treeprocessor {

    public DummyTreeprocessor() {
        super(new HashMap<String, Object>());
    }

    public DummyTreeprocessor(Map<String, Object> config) {
        super(config);
        System.out.println(this.getClass().getSimpleName() + "(" 
                + this.getClass().getSuperclass().getSimpleName() + ") initialized");
    }

    @Override
    public Document process(Document document) {
        System.out.println("Processing "+ this.getClass().getSimpleName());
        System.out.println("Processing: blocks found: " + document.getBlocks().size());
        return document;
    }

}
