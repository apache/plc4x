package org.apache.plc4x.utils.maven.site.asciidoctor.test.processors;

import org.asciidoctor.ast.Document;
import org.asciidoctor.extension.Preprocessor;
import org.asciidoctor.extension.PreprocessorReader;

import java.util.Map;

public class FailingPreprocessor extends Preprocessor {

    public FailingPreprocessor(Map<String, Object> config) {
        super(config);
        System.out.println(this.getClass().getSimpleName() + "("
                + this.getClass().getSuperclass().getSimpleName() + ") initialized");
    }

    @Override
    public PreprocessorReader process(Document document, PreprocessorReader reader) {
        System.out.println("Processing "+ this.getClass().getSimpleName());
        System.out.println("Processing: blocks found: " + document.getBlocks().size());
        throw new RuntimeException("That's all folks");
    }

}
