package org.apache.plc4x.utils.maven.site.asciidoctor.extensions;

import org.apache.maven.plugin.MojoExecutionException;
import org.asciidoctor.extension.Processor;

public interface ExtensionRegistry {

    /**
     * Checks if {@code extensionClassName} belongs to a valid {@link Processor}
     * class and if it can be found in the classpath
     * 
     * @param extensionClassName
     *             fully qualified name of the class implementing the extension
     * @param blockName
     *            required when declaring
     * 
     */
    void register(String extensionClassName, String blockName) throws MojoExecutionException;

}