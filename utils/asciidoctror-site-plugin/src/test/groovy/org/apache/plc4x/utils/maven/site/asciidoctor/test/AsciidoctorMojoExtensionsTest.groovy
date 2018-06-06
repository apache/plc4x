package org.apache.plc4x.utils.maven.site.asciidoctor.test

import org.apache.maven.plugin.MojoExecutionException
import org.apache.plc4x.utils.maven.site.asciidoctor.AsciidoctorMojo
import org.apache.plc4x.utils.maven.site.asciidoctor.AsciidoctorZipMojo
import org.apache.plc4x.utils.maven.site.asciidoctor.extensions.ExtensionConfiguration
import org.apache.plc4x.utils.maven.site.asciidoctor.test.plexus.MockPlexusContainer
import org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.ChangeAttributeValuePreprocessor
import org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.FailingPreprocessor
import org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.GistBlockMacroProcessor
import org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.ManpageInlineMacroProcessor
import org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.MetaDocinfoProcessor
import org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.UriIncludeProcessor
import org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.YellBlockProcessor

import spock.lang.Specification
import spock.lang.Unroll

/**
 * Specific tests to validate usage of AsciidoctorJ extension in AsciidoctorMojo.
 * 
 * Most of the examples have been directly adapted from the ones found in AsciidoctorJ 
 * documentation (https://github.com/asciidoctor/asciidoctorj/blob/master/README.adoc)
 *
 * @author abelsromero
 */
class AsciidoctorMojoExtensionsTest extends Specification {

    /**
     * Intercept Asciidoctor mojo constructor to mock and inject required
     * plexus objects
     */
    def setupSpec() {
        MockPlexusContainer mockPlexusContainer = new MockPlexusContainer()
        def oldConstructor = AsciidoctorMojo.constructors[0]

        AsciidoctorMojo.metaClass.constructor = {
            def mojo = oldConstructor.newInstance()
            mockPlexusContainer.initializeContext(mojo)
            return mojo
        }
    }

    static final String SRC_DIR = 'target/test-classes/src/asciidoctor/'
    static final String OUTPUT_DIR = 'target/asciidoctor-output-extensions'

    def "fails because processor is not found in classpath"() {
        setup:
            File srcDir = new File(SRC_DIR)
            File outputDir = new File("${OUTPUT_DIR}/preprocessor/${System.currentTimeMillis()}")
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'processors-sample.adoc'
            mojo.outputDirectory = outputDir
            mojo.headerFooter = true
            mojo.attributes['toc'] = null
            mojo.attributes['linkcss!'] = ''
            mojo.extensions = [
                [className: 'non.existent.Processor'] as ExtensionConfiguration
            ]
            mojo.execute()
        then:
            outputDir.list().size() == 0
            def e = thrown(MojoExecutionException)
            e.message.contains(mojo.extensions.get(0).className)
            e.message.contains('not found in classpath')
    }
    
    // This test is added to keep track of possible changes in the extension's SPI
    def "plugin fails because processor throws an uncaught exception"() {
        setup:
            File srcDir = new File(SRC_DIR)
            File outputDir = new File("${OUTPUT_DIR}/preprocessor/${System.currentTimeMillis()}")
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'processors-sample.adoc'
            mojo.outputDirectory = outputDir
            mojo.headerFooter = true
            mojo.attributes['toc'] = null
            mojo.attributes['linkcss!'] = ''
            mojo.extensions = [
                [className: FailingPreprocessor.class.canonicalName] as ExtensionConfiguration
            ]
            mojo.execute()
        then:
            // since v 1.5.4 resources are copied before rendering, so some files remain
            outputDir.list().size() > 0
            thrown(RuntimeException)
    }
    

    /**
     * Redirects output to validate specific traces left in the processors  
     */
    @Unroll
    def "tests that a #processorType is registered, initialized and executed"() {

        setup:
            ByteArrayOutputStream systemOut = new ByteArrayOutputStream()
            System.out = new PrintStream(systemOut)
    
            File srcDir = new File(SRC_DIR)
            File outputDir = new File("${OUTPUT_DIR}/processors/${System.currentTimeMillis()}")
    
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'processors-sample.adoc'
            mojo.outputDirectory = outputDir
            mojo.headerFooter = true
            mojo.attributes['toc'] = null
            mojo.extensions = [
                [className: "org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.${processorClass}"] as ExtensionConfiguration
            ]
            mojo.execute()

        expect:
            systemOut.toString().contains(initializationMessage)
            systemOut.toString().contains(executionMessage)

            where:
            processorClass          | processorType           || initializationMessage                                        || executionMessage
            'ChangeAttributeValuePreprocessor'|'Preprocessor' || "ChangeAttributeValuePreprocessor(Preprocessor) initialized" || 'Processing ChangeAttributeValuePreprocessor'
            'DummyTreeprocessor'    |'Treeprocessor'          || "DummyTreeprocessor(Treeprocessor) initialized"              || 'Processing DummyTreeprocessor'
            'DummyPostprocessor'    |'Postprocessor'          || "DummyPostprocessor(Postprocessor) initialized"              || 'Processing DummyPostprocessor'
            'MetaDocinfoProcessor'  |'DocinfoProcessor'       || "MetaDocinfoProcessor(DocinfoProcessor) initialized"         || 'Processing MetaDocinfoProcessor'
            'UriIncludeProcessor'   |'IncludeProcessor'       || "UriIncludeProcessor(IncludeProcessor) initialized"          || 'Processing UriIncludeProcessor'
    }

    def "successfully renders html with a preprocessor"() {
        setup:
            File srcDir = new File(SRC_DIR)
            File outputDir = new File("${OUTPUT_DIR}/preprocessor/${System.currentTimeMillis()}")
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'processors-sample.adoc'
            mojo.outputDirectory = outputDir
            mojo.headerFooter = true
            mojo.attributes['toc'] = null
            mojo.extensions = [
                [className: ChangeAttributeValuePreprocessor.class.canonicalName] as ExtensionConfiguration
            ]
            mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('processors-sample.html')
    
            File sampleOutput = new File(outputDir, 'processors-sample.html')
            sampleOutput.length() > 0
            String text = sampleOutput.getText()
            text.count(ChangeAttributeValuePreprocessor.AUTHOR_NAME) == 2
    }

    def "successfully renders html with a blockprocessor"() {
        setup:
            File srcDir = new File(SRC_DIR)
            File outputDir = new File("${OUTPUT_DIR}/blockprocessor/${System.currentTimeMillis()}")
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'processors-sample.adoc'
            mojo.outputDirectory = outputDir
            mojo.headerFooter = true
            mojo.attributes['toc'] = null
            mojo.extensions = [
                [className: YellBlockProcessor.class.canonicalName, blockName:'yell'] as ExtensionConfiguration
            ]
            mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('processors-sample.html')
    
            File sampleOutput = new File(outputDir, 'processors-sample.html')
            sampleOutput.length() > 0
            sampleOutput.getText().contains('The time is now. Get a move on.'.toUpperCase())
    }
    
    def "successfully renders html and adds meta tag with a DocinfoProcessor"() {
        setup:
            File srcDir = new File(SRC_DIR)
            File outputDir = new File("${OUTPUT_DIR}/docinfoProcessor/${System.currentTimeMillis()}")
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'processors-sample.adoc'
            mojo.outputDirectory = outputDir
            mojo.headerFooter = true
            mojo.attributes['toc'] = null
            mojo.extensions = [
                [className: MetaDocinfoProcessor.class.canonicalName] as ExtensionConfiguration
            ]
            mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('processors-sample.html')
    
            File sampleOutput = new File(outputDir, 'processors-sample.html')
            sampleOutput.length() > 0
            sampleOutput.text.contains("<meta name=\"author\" content=\"asciidoctor\">")
    }

    def "successfully renders html and modifies output with a BlockMacroProcessor"() {
        setup:
            File srcDir = new File(SRC_DIR)
            File outputDir = new File("${OUTPUT_DIR}/blockMacroProcessor/${System.currentTimeMillis()}")
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'processors-sample.adoc'
            mojo.outputDirectory = outputDir
            mojo.headerFooter = true
            mojo.attributes['toc'] = null
            mojo.extensions = [
                [className: GistBlockMacroProcessor.class.canonicalName, blockName:'gist'] as ExtensionConfiguration
            ]
            mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('processors-sample.html')
    
            File sampleOutput = new File(outputDir, 'processors-sample.html')
            sampleOutput.length() > 0
            sampleOutput.text.contains("<script src=\"https://gist.github.com/123456.js\"></script>")
    }

    def "successfully renders html and modifies output with a InlineMacroProcessor"() {
        setup:
            File srcDir = new File(SRC_DIR)
            File outputDir = new File("${OUTPUT_DIR}/inlineMacroProcessor/${System.currentTimeMillis()}")
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'processors-sample.adoc'
            mojo.outputDirectory = outputDir
            mojo.headerFooter = true
            mojo.attributes['toc'] = null
            mojo.extensions = [
                [className: ManpageInlineMacroProcessor.class.canonicalName, blockName:'man'] as ExtensionConfiguration
            ]
            mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('processors-sample.html')
    
            File sampleOutput = new File(outputDir, 'processors-sample.html')
            sampleOutput.length() > 0
            sampleOutput.text.contains("<p>See <a href=\"gittutorial.html\">gittutorial</a> to get started.</p>")
    }
    
    def "successfully renders html and modifies output with an IncludeProcessor"() {
        setup:
            File srcDir = new File(SRC_DIR)
            File outputDir = new File("${OUTPUT_DIR}/includeProcessor/${System.currentTimeMillis()}")
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'processors-sample.adoc'
            mojo.outputDirectory = outputDir
            mojo.headerFooter = true
            mojo.attributes['toc'] = null
            mojo.extensions = [
                [className: UriIncludeProcessor.class.canonicalName] as ExtensionConfiguration
            ]
        mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('processors-sample.html')
    
            File sampleOutput = new File(outputDir, 'processors-sample.html')
            sampleOutput.length() > 0
            sampleOutput.text.contains("source 'https://rubygems.org'")
    }

    def "executes the same preprocessor twice"() {
        setup:
            ByteArrayOutputStream systemOut = new ByteArrayOutputStream()
            System.out = new PrintStream(systemOut)
            File srcDir = new File(SRC_DIR)
            File outputDir = new File("${OUTPUT_DIR}/preprocessor/${System.currentTimeMillis()}")
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'processors-sample.adoc'
            mojo.outputDirectory = outputDir
            mojo.headerFooter = true
            mojo.attributes['toc'] = null
            mojo.extensions = [
                [className: ChangeAttributeValuePreprocessor.class.canonicalName] as ExtensionConfiguration,
                [className: ChangeAttributeValuePreprocessor.class.canonicalName] as ExtensionConfiguration
            ]
            mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('processors-sample.html')
    
            File sampleOutput = new File(outputDir, 'processors-sample.html')
            sampleOutput.length() > 0
            String text = sampleOutput.getText()
            text.count(ChangeAttributeValuePreprocessor.AUTHOR_NAME) == 2
            
            systemOut.toString().count('Processing ChangeAttributeValuePreprocessor') == 2
    }


    // Adding a BlockMacroProcessor or BlockProcessor makes the conversion fail
    def "successfully renders html with Preprocessor, DocinfoProcessor, InlineMacroProcessor and IncludeProcessor"() {
        setup:
            File srcDir = new File(SRC_DIR)
            File outputDir = new File("${OUTPUT_DIR}/processors/${System.currentTimeMillis()}")
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'processors-sample.adoc'
            mojo.outputDirectory = outputDir
            mojo.headerFooter = true
            mojo.attributes['toc'] = null
            mojo.extensions = [
                // Preprocessor
                [className: 'org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.ChangeAttributeValuePreprocessor'] as ExtensionConfiguration,
                // DocinfoProcessor
                [className: 'org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.MetaDocinfoProcessor'] as ExtensionConfiguration,
                // InlineMacroProcessor
                [className: 'org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.ManpageInlineMacroProcessor', blockName:'man'] as ExtensionConfiguration,
                // IncludeProcessor
                [className: 'org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.UriIncludeProcessor'] as ExtensionConfiguration,
            ]
            mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('processors-sample.html')
    
            File sampleOutput = new File(outputDir, 'processors-sample.html')
            sampleOutput.length() > 0

            String text = sampleOutput.text
            text.count(ChangeAttributeValuePreprocessor.AUTHOR_NAME) == 2
            text.contains("<meta name=\"author\" content=\"asciidoctor\">")
            text.contains("<p>See <a href=\"gittutorial.html\">gittutorial</a> to get started.</p>")
            text.contains("source 'https://rubygems.org'")
    }
    
    def "renders html when using all types of extensions"() {
        setup:
            File srcDir = new File(SRC_DIR)
            File outputDir = new File("${OUTPUT_DIR}/processors/${System.currentTimeMillis()}")
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'processors-sample.adoc'
            mojo.outputDirectory = outputDir
            mojo.headerFooter = true
            mojo.attributes['toc'] = ''
            mojo.attributes['linkcss'] = ''
            mojo.attributes['copycss!'] = ''
            mojo.extensions = [
                // Preprocessor
                [className: 'org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.ChangeAttributeValuePreprocessor'] as ExtensionConfiguration,
                // DocinfoProcessor
                [className: 'org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.MetaDocinfoProcessor'] as ExtensionConfiguration,
                // InlineMacroProcessor
                [className: 'org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.ManpageInlineMacroProcessor', blockName:'man'] as ExtensionConfiguration,
                // IncludeProcessor
                [className: 'org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.UriIncludeProcessor'] as ExtensionConfiguration,
                // BlockMacroProcessor
                [className: 'org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.GistBlockMacroProcessor', blockName:'gist'] as ExtensionConfiguration,
                // BlockProcessor
                [className: 'org.apache.plc4x.utils.maven.site.asciidoctor.test.processors.YellBlockProcessor', blockName:'yell'] as ExtensionConfiguration
            ]
            mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('processors-sample.html')
    
            File sampleOutput = new File(outputDir, 'processors-sample.html')
            sampleOutput.length() > 0

            String text = sampleOutput.text
            text.contains('<meta name="author" content="asciidoctor">')
            text.contains('<script src="https://gist.github.com/123456.js"></script>')
            text.contains('<p>See <a href="gittutorial.html">gittutorial</a> to get started.</p>')
            text.contains('<p>THE TIME IS NOW. GET A MOVE ON.</p>')
    }

    /**
     *  Manual test to validate automatic extension registration.
     *  To execute, copy org.asciidoctor.extension.spi.ExtensionRegistry to 
     *  /src/test/resources/META-INF/services/ and execute
     */
    @spock.lang.Ignore
    def "property extension"() {
        setup:
            File srcDir = new File(SRC_DIR)
            File outputDir = new File("${OUTPUT_DIR}/preprocessor/${System.currentTimeMillis()}")
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'processors-sample.adoc'
            mojo.outputDirectory = outputDir
            mojo.headerFooter = true
            mojo.attributes['toc'] = null
            mojo.attributes['linkcss!'] = ''
            mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('processors-sample.html')
    
            File sampleOutput = new File(outputDir, 'processors-sample.html')
            sampleOutput.length() > 0
            String text = sampleOutput.getText()
            text.count(ChangeAttributeValuePreprocessor.AUTHOR_NAME) == 2
    }
    
}
