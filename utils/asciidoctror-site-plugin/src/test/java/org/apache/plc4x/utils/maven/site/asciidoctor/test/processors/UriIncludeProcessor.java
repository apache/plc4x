package org.apache.plc4x.utils.maven.site.asciidoctor.test.processors;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.MalformedURLException;
import java.net.URL;
import java.util.Map;

import org.asciidoctor.ast.DocumentRuby;
import org.asciidoctor.extension.IncludeProcessor;
import org.asciidoctor.extension.PreprocessorReader;

public class UriIncludeProcessor extends IncludeProcessor {

    public UriIncludeProcessor(Map<String, Object> config) {
        super(config);
        System.out.println(this.getClass().getSimpleName() + "(" 
              + this.getClass().getSuperclass().getSimpleName() + ") initialized");
    }

    @Override
    public boolean handles(String target) {
        return target.startsWith("http://") || target.startsWith("https://");
    }

    @Override
    public void process(DocumentRuby document, PreprocessorReader reader, String target,
            Map<String, Object> attributes) {

        System.out.println("Processing "+ this.getClass().getSimpleName());

        StringBuilder content = readContent(target);
        reader.push_include(content.toString(), target, target, 1, attributes);

    }

    private StringBuilder readContent(String target) {

        StringBuilder content = new StringBuilder();

        try {

            URL url = new URL(target);
            InputStream openStream = url.openStream();

            BufferedReader bufferedReader = new BufferedReader(
                    new InputStreamReader(openStream));

            String line = null;
            while ((line = bufferedReader.readLine()) != null) {
                content.append(line + "\n");
            }

            bufferedReader.close();

        } catch (MalformedURLException e) {
            throw new IllegalArgumentException(e);
        } catch (IOException e) {
            throw new IllegalArgumentException(e);
        }
        return content;
    }

  }
