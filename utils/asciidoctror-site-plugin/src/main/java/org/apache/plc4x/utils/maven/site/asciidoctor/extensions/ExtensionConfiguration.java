package org.apache.plc4x.utils.maven.site.asciidoctor.extensions;

import org.apache.maven.plugins.annotations.Parameter;
import org.apache.plc4x.utils.maven.site.asciidoctor.AsciidoctorMaven;

/**
 * Holds a processor's configuration parameters in the pom.xml
 * 
 * @author abelsromero
 */
public class ExtensionConfiguration {

    public static final String PREFIX = AsciidoctorMaven.PREFIX + "extension.";
    
    /**
     * Fully qualified name of the processor
     */
    @Parameter(property = PREFIX + "className", required = true)
    private String className;

    /**
     * Optional. Block name in case of setting a Block, BlockMacro or
     * InlineMacro processor
     */
    @Parameter(property = PREFIX + "blockName")
    private String blockName;

    public ExtensionConfiguration() {
    }

    public String getClassName() {
        return className;
    }

    public void setClassName(String className) {
        this.className = className;
    }

    public String getBlockName() {
        return blockName;
    }

    public void setBlockName(String blockName) {
        this.blockName = blockName;
    }

}
