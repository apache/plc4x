package org.apache.plc4x.utils.maven.site.asciidoctor.test.processors;

import org.asciidoctor.ast.AbstractBlock;
import org.asciidoctor.ast.Block;
import org.asciidoctor.extension.BlockMacroProcessor;

import java.util.Arrays;
import java.util.Map;

public class GistBlockMacroProcessor extends BlockMacroProcessor {

    public GistBlockMacroProcessor(String macroName, Map<String, Object> config) {
        super(macroName, config);
    }

    @Override
    public Block process(AbstractBlock parent, String target,
                          Map<String, Object> attributes) {

      String content = "<div class=\"content\">\n" +
              "<script src=\"https://gist.github.com/"+target+".js\"></script>\n" +
              "</div>";

      return createBlock(parent, "pass", Arrays.asList(content), attributes,
              this.getConfig());
    }

  }
