package org.apache.plc4x.utils.maven.site.asciidoctor.test

import groovy.io.FileType
import org.apache.commons.io.FileUtils
import org.apache.maven.model.Resource
import org.apache.plc4x.utils.maven.site.asciidoctor.AsciidoctorMojo
import org.apache.plc4x.utils.maven.site.asciidoctor.io.AsciidoctorFileScanner
import org.apache.plc4x.utils.maven.site.asciidoctor.test.plexus.MockPlexusContainer
import spock.lang.Specification

/**
 *
 */
class AsciidoctorMojoTest extends Specification {

    static final String DEFAULT_SOURCE_DIRECTORY = 'target/test-classes/src/asciidoctor'
    static final String MULTIPLE_RESOURCES_OUTPUT = 'target/asciidoctor-output/multiple-resources'

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

    def "renders docbook"() {
        setup:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output')

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'docbook'
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.sourceDocumentName = 'sample.asciidoc'
            mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('sample.xml')

            File sampleOutput = new File('sample.xml', outputDir)
            sampleOutput.length() > 0
    }

    def "renders html"() {
        setup:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output')

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'sample.asciidoc'
            mojo.resources = [new Resource(directory: '.', excludes: ['**/**'])]
            mojo.outputDirectory = outputDir
            mojo.headerFooter = true
            mojo.sourceHighlighter = 'coderay'
            // IMPORTANT Maven can only assign string values or null, so we have to emulate the value precisely in the test!
            // Believe it or not, null is the equivalent of writing <toc/> in the XML configuration
            mojo.attributes['toc'] = null
            mojo.attributes['linkcss!'] = ''
            mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('sample.html')

            File sampleOutput = new File('sample.html', outputDir)
            sampleOutput.length() > 0
            String text = sampleOutput.getText()
            text.contains('<body class="article">')
            text.contains('id="toc"')
            text.contains('Asciidoctor default stylesheet')
            !text.contains('<link rel="stylesheet" href="./asciidoctor.css">')
            text.contains('<pre class="CodeRay highlight">')
    }

    def "docinfo file should be ignored html"() {
        setup:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output-docinfo')

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'sample-with-doc-info.asciidoc'
            mojo.outputDirectory = outputDir
            mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('sample-with-doc-info.html')
            !outputDir.list().toList().contains('sample-with-doc-info-docinfo.html')
            !outputDir.list().toList().contains('sample-with-doc-info-docinfo-footer.html')
            !outputDir.list().toList().contains('sample-with-doc-info-docinfo.xml')
            !outputDir.list().toList().contains('sample-with-doc-info-docinfo-footer.xml')

            File sampleOutput = new File('sample-with-doc-info.html', outputDir)
            sampleOutput.length() > 0
            String text = sampleOutput.getText()
            text.contains('This is the sample-with-doc-info file.')
            text.contains('This is the docinfo html file.')
            text.contains('This is the docinfo html footer.')
    }

    def "should honor doctype set in document"() {
        setup:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output')
            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'book.adoc'
            mojo.outputDirectory = outputDir
            mojo.headerFooter = true
            mojo.attributes['linkcss'] = ''
            mojo.attributes['copycss!'] = ''
            mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('book.html')
            File sampleOutput = new File('book.html', outputDir)
            sampleOutput.length() > 0
            String text = sampleOutput.getText()
            text.contains('<body class="book">')
    }

    def "asciidoc file extension can be changed"() {
        given: 'an empty output directory'
            def outputDir = new File('target/asciidoctor-output')
            outputDir.delete()

        when: 'asciidoctor mojo is called with extension foo and bar and it exists a sample1.foo and a sample2.bar'
            def srcDir = new File(DEFAULT_SOURCE_DIRECTORY)

            outputDir.mkdirs()

            new File(srcDir, 'sample1.foo').withWriter {
                it << '''
                    Document Title
                    ==============

                    foo
                    '''.stripIndent()
            }
            new File(srcDir, 'sample2.bar').withWriter {
                it << '''
                    Document Title
                    ==============

                    bar
                    '''.stripIndent()
            }

            def mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.sourceDocumentExtensions = [ 'foo', 'bar' ]
            mojo.execute()

        then: 'sample1.html and sample2.html exist and contain the extension of the original file'
            def outputs = outputDir.list().toList()
            outputs.size() >= 2
            outputs.contains('sample1.html')
            outputs.contains('sample2.html')

            new File(outputDir, 'sample1.html').text.contains('foo')
            new File(outputDir, 'sample2.html').text.contains('bar')
    }

    def "header footer is enabled by default"() {
        when:
          AsciidoctorMojo mojo = new AsciidoctorMojo()
        then:
          mojo.headerFooter == true
    }

    def "should require library"() {
        setup:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output')

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.requires = ['time'] as List
            mojo.backend = 'html'
            mojo.outputDirectory = outputDir
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'sample.asciidoc'
            mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('sample.html')
            assert "constant".equals(org.asciidoctor.internal.JRubyRuntimeContext.get().evalScriptlet('defined? ::DateTime').toString())
    }

    def "embedding resources"() {
        setup:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output')

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.attributes["icons"] = "font"
            mojo.embedAssets = true
            mojo.imagesDir = ''
            mojo.outputDirectory = outputDir
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'sample-embedded.adoc'
            mojo.backend = 'html'
            mojo.execute()
        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('sample-embedded.html')

            File sampleOutput = new File(outputDir, 'sample-embedded.html')
            sampleOutput.length() > 0
            String text = sampleOutput.getText()
            text.contains('Asciidoctor default stylesheet')
            text.contains('data:image/png;base64,iVBORw0KGgo')
            text.contains('font-awesome.min.css')
            text.contains('i class="fa icon-tip"')
    }

    def "override output file"() {
        setup:
        File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
        File outputDir = new File('target/asciidoctor-output')
        File outputFile = new File( 'custom_output_file.html')

        if (!outputDir.exists())
            outputDir.mkdir()
        when:
        AsciidoctorMojo mojo = new AsciidoctorMojo()
        mojo.attributes["icons"] = "font"
        mojo.embedAssets = true
        mojo.imagesDir = ''
        mojo.outputDirectory = outputDir
        mojo.outputFile = outputFile
        mojo.sourceDirectory = srcDir
        mojo.sourceDocumentName = 'sample-embedded.adoc'
        mojo.backend = 'html'
        mojo.execute()
        then:
        outputDir.list().toList().isEmpty() == false
        outputDir.list().toList().contains('custom_output_file.html')

        File sampleOutput = new File(outputDir, 'custom_output_file.html')
        sampleOutput.length() > 0
        String text = sampleOutput.getText()
        text.contains('Asciidoctor default stylesheet')
        text.contains('data:image/png;base64,iVBORw0KGgo')
        text.contains('font-awesome.min.css')
        text.contains('i class="fa icon-tip"')
    }

    def "override output file with absolute path"() {
        setup:
        File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
        File outputDir = new File('target/asciidoctor-output')
        File outputFile = new File(System.getProperty('java.io.tmpdir'), 'custom_output_file.html')

        if (!outputDir.exists())
            outputDir.mkdir()
        when:
        AsciidoctorMojo mojo = new AsciidoctorMojo()
        mojo.attributes["icons"] = "font"
        mojo.embedAssets = true
        mojo.imagesDir = ''
        mojo.outputDirectory = outputDir
        mojo.outputFile = outputFile
        mojo.sourceDirectory = srcDir
        mojo.sourceDocumentName = 'sample-embedded.adoc'
        mojo.backend = 'html'
        mojo.execute()
        then:
        outputDir.list().toList().isEmpty() == false
        outputDir.list().toList().contains('custom_output_file.html')

        File sampleOutput = new File(outputDir, 'custom_output_file.html')
        sampleOutput.length() > 0
        String text = sampleOutput.getText()
        text.contains('Asciidoctor default stylesheet')
        text.contains('data:image/png;base64,iVBORw0KGgo')
        text.contains('font-awesome.min.css')
        text.contains('i class="fa icon-tip"')
    }

    def "missing-attribute skip"() {
        given:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output')

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.outputDirectory = outputDir
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'attribute-missing.adoc'
            mojo.backend = 'html'
            mojo.attributeMissing = 'skip'
            mojo.execute()
        then:
            File sampleOutput = new File(outputDir, 'attribute-missing.html')
            String text = sampleOutput.getText()
            text.contains('Here is a line that has an attribute that is {missing}!')
    }

    def "missing-attribute drop"() {
        given:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output')

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.outputDirectory = outputDir
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'attribute-missing.adoc'
            mojo.backend = 'html'
            mojo.attributeMissing = 'drop'
            mojo.execute()
        then:
            File sampleOutput = new File(outputDir, 'attribute-missing.html')
            String text = sampleOutput.getText()
            text.contains('Here is a line that has an attribute that is !')
            !text.contains('{name}')
    }

    def "missing-attribute drop-line"() {
        given:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output')

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.outputDirectory = outputDir
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'attribute-missing.adoc'
            mojo.backend = 'html'
            mojo.attributeMissing = 'drop-line'
            mojo.execute()
        then:
            File sampleOutput = new File(outputDir, 'attribute-missing.html')
            String text = sampleOutput.getText()
            !text.contains('Here is a line that has an attribute that is')
            !text.contains('{set: name!}')
    }

    def "undefined-attribute drop"() {
        given:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output')

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.outputDirectory = outputDir
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'attribute-undefined.adoc'
            mojo.backend = 'html'
            mojo.attributeUndefined = 'drop'
            mojo.execute()
        then:
            File sampleOutput = new File(outputDir, 'attribute-undefined.html')
            String text = sampleOutput.getText()
            text.contains('Here is a line that has an attribute that is !')
            !text.contains('{set: name!}')
    }

    def "undefined-attribute drop-line"() {
        given:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output')

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.outputDirectory = outputDir
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'attribute-undefined.adoc'
            mojo.backend = 'html'
            mojo.attributeMissing = 'drop-line'
            mojo.execute()
        then:
            File sampleOutput = new File(outputDir, 'attribute-undefined.html')
            String text = sampleOutput.getText()
            !text.contains('Here is a line that has an attribute that is')
            !text.contains('{set: name!}')
    }

    // Test for Issue 62
    def 'setting_boolean_values'() {
        given:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output-issue-62')

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.outputDirectory = outputDir
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'sample.asciidoc'
            mojo.backend = 'html'
            // IMPORTANT Maven can only assign string values or null, so we have to emulate the value precisely in the test!
            // Believe it or not, null is the equivalent of writing <toc/> in the XML configuration
            mojo.attributes.put('toc2', 'true')
            mojo.execute()
        then:
            File sampleOutput = new File(outputDir, 'sample.html')
            String text = sampleOutput.getText()
            text.contains('class="toc2"')

    }

    // Test for Issue 62 (unset)
    def 'unsetting_boolean_values'() {
        given:
        File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
        File outputDir = new File('target/asciidoctor-output-issue-62-unset')

        if (!outputDir.exists())
            outputDir.mkdir()
        when:
        AsciidoctorMojo mojo = new AsciidoctorMojo()
        mojo.outputDirectory = outputDir
        mojo.sourceDirectory = srcDir
        mojo.sourceDocumentName = 'sample.asciidoc'
        mojo.backend = 'html'
        // IMPORTANT Maven can only assign string values or null, so we have to emulate the value precisely in the test!
        // Believe it or not, null is the equivalent of writing <toc/> in the XML configuration
        mojo.attributes.put('toc2', 'false')
        mojo.execute()
        then:
        File sampleOutput = new File(outputDir, 'sample.html')
        String text = sampleOutput.getText()
        !text.contains('class="toc2"')
    }

    def 'test_imageDir_properly_passed'() {
        given:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output-imageDir')

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.outputDirectory = outputDir
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'imageDir.adoc'
            mojo.backend = 'html'
            mojo.imagesDir = 'images-dir'
            mojo.execute()
        then:
            File sampleOutput = new File(outputDir, 'imageDir.html')
            String text = sampleOutput.getText()
            text.contains('<img src="images-dir/my-cool-image.jpg" alt="my cool image">')
    }

    def 'includes_test'() {
        given:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output-include-test')

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.sourceDocumentName = new File('main-document.adoc')
            mojo.backend = 'html'
            mojo.execute()
        then:
            File mainDocumentOutput = new File(outputDir, 'main-document.html')
            String text = mainDocumentOutput.getText()
            text.contains('This is the parent document')
            text.contains('This is an included file.')
    }

    def 'skip'() {
        given:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output-skip-test')
            if (outputDir.exists())
                outputDir.delete()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.sourceDocumentName = new File('main-document.adoc')
            mojo.backend = 'html'
            mojo.skip = true
            mojo.execute()
        then:
            !outputDir.exists()
    }

    def 'issue-78'() {
        given:
        File srcDir = new File('target/test-classes/src/asciidoctor/issue-78')
        File outputDir = new File('target/asciidoctor-output-issue-78')

        if (!outputDir.exists())
            outputDir.mkdir()
        when:
        AsciidoctorMojo mojo = new AsciidoctorMojo()
        mojo.sourceDirectory = srcDir
        mojo.outputDirectory = outputDir
        mojo.sourceDocumentName = new File('main.adoc')
        mojo.doctype = 'book'
        mojo.embedAssets = true
        // IMPORTANT Maven can only assign string values or null, so we have to emulate the value precisely in the test!
        // Believe it or not, null is the equivalent of writing <toc/> in the XML configuration
        mojo.attributes['toc'] = 'true'
        mojo.backend = 'html'
        mojo.execute()
        then:
        File mainDocumentOutput = new File(outputDir, 'main.html')
        File imageFile = new File(outputDir, 'images/halliburton_lab.jpg')
        imageFile.exists()
        String text = mainDocumentOutput.getText()
        text.contains("<p>Here&#8217;s an image:</p>")
        text.contains('<img src="data:image/jpg;base64,/9j/4AAQSkZJRgABAQEASABIAAD/4gzESUNDX1BST0ZJTEUAAQEAAA')
    }

    def 'code highlighting should be rendered when set in the document header'() {
        setup:
        File srcDir = new File('src/test/resources/src/asciidoctor')
        File outputDir = new File('target/asciidoctor-output-sourceHighlighting/header')
        String documentName = 'sample-with-source-highlighting'

        when:
        AsciidoctorMojo mojo = new AsciidoctorMojo()
        mojo.sourceDirectory = srcDir
        mojo.outputDirectory = outputDir
        mojo.sourceDocumentName = new File("${documentName}.adoc")
        mojo.resources = [new Resource(directory: '.', excludes: ['**/**'])]
        mojo.backend = 'html5'
        mojo.execute()

        then:
        File mainDocumentOutput = new File(outputDir, "${documentName}.html")
        String text = mainDocumentOutput.getText()
        text.contains('<pre class="CodeRay highlight">')
    }

    /**
     * Tests CodeRay source code highlighting options.
     */
    def 'code highlighting - coderay'() {
        setup:
            File srcDir = new File('src/test/resources/src/asciidoctor')
            File outputDir = new File('target/asciidoctor-output-sourceHighlighting/coderay')

        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.sourceHighlighter = 'coderay'
            mojo.sourceDocumentName = new File('main-document.adoc')
            mojo.backend = 'html'
            mojo.execute()

        then:
            File mainDocumentOutput = new File(outputDir, 'main-document.html')
            String text = mainDocumentOutput.getText()
            text.contains('<pre class="CodeRay highlight">')
    }

    /**
     * Tests Highlight.js source code highlighting options.
     */
    def 'code highlighting - highlight.js'() {
        setup:
            File srcDir = new File('src/test/resources/src/asciidoctor')
            File outputDir = new File('target/asciidoctor-output-sourceHighlighting/highlightjs')

        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.sourceHighlighter = 'highlight.js'
            mojo.sourceDocumentName = new File('main-document.adoc')
            mojo.backend = 'html'
            mojo.execute()

        then:
            File mainDocumentOutput = new File(outputDir, 'main-document.html')
            String text = mainDocumentOutput.getText()
            text.contains('highlight')
    }

    /**
     * Tests Prettify source code highlighting options.
     */
    def 'code highlighting - prettify'() {
        setup:
            File srcDir = new File('src/test/resources/src/asciidoctor')
            File outputDir = new File('target/asciidoctor-output-sourceHighlighting/prettify')

        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.sourceHighlighter = 'prettify'
            mojo.sourceDocumentName = new File('main-document.adoc')
            mojo.backend = 'html'
            mojo.execute()

        then:
            File mainDocumentOutput = new File(outputDir, 'main-document.html')
            String text = mainDocumentOutput.getText()
            text.contains('prettify')
    }

    /**
     * Tests behavior when source code highlighting with Pygments is specified.
     *
     * Test checks that an exception is not thrown.
     */
    def 'code highlighting - pygments'() {
        setup:
            File srcDir = new File('src/test/resources/src/asciidoctor')
            File outputDir = new File('target/asciidoctor-output-sourceHighlighting/pygments')

        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.sourceHighlighter = 'pygments'
            mojo.sourceDocumentName = new File('main-document.adoc')
            mojo.backend = 'html'
            mojo.execute()

        then:
            File mainDocumentOutput = new File(outputDir, 'main-document.html')
            String text = mainDocumentOutput.getText()
            text.contains('Pygments is not available.')
            text.contains('<pre class="pygments highlight">')
    }

    /**
     * Tests behaviour when an invalid source code highlighting option is set.
     *
     * Test checks that no additional CSS are added.
     */
    def 'code highlighting - nonExistent'() {
        setup:
            File srcDir = new File('src/test/resources/src/asciidoctor')
            File outputDir = new File('target/asciidoctor-output-sourceHighlighting/nonExistent')

        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.sourceHighlighter = 'nonExistent'
            mojo.sourceDocumentName = new File('main-document.adoc')
            mojo.backend = 'html'
            mojo.execute()

        then:
            File mainDocumentOutput = new File(outputDir, 'main-document.html')
            String text = mainDocumentOutput.getText()
            // No extra CSS is added other than Asciidoctor's default
            text.count('<style>') == 1
    }

    /**
     * Tests for relative folder structures treatment
     */
    static final FileFilter DIRECTORY_FILTER = {File f -> f.isDirectory()} as FileFilter
    static final String ASCIIDOC_REG_EXP_EXTENSION = '.*\\.a((sc(iidoc)?)|d(oc)?)$'

    /**
     * Validates that the folder structures under certain files are the same
     *
     * @param expected
     *         list of expected folders
     * @param actual
     *         list of actual folders (the ones to validate)
     * @param ignoreUnderscore
     *         tells where to ignore hidden Asciidoctor files (prefixed with underscore) in the expected parameter
     */
    private void assertEqualsStructure (File[] expected, File[] actual, Boolean ignoreUnderscore = true) {
        if (ignoreUnderscore)
            expected = expected.findAll { !it.name.startsWith('_')}

        assert expected.length == actual.length
        expected*.name.containsAll(actual*.name)
        actual*.name.containsAll(expected*.name)
        for (File actualFile in actual) {
            File expectedFile = expected.find { it.getName() == actualFile.getName() }
            assert expectedFile != null

            // check that at least the number of html files and asciidoc are the same in each folder
            File[] expectedChildren = expectedFile.listFiles(DIRECTORY_FILTER)
            if (ignoreUnderscore)
                expectedChildren = expectedChildren.findAll { !it.name.startsWith('_')}

            File[] htmls =  actualFile.listFiles({File f -> f.getName() ==~ /.+html/} as FileFilter)
            if (htmls) {
                File[] asciidocs =  expectedFile.listFiles({File f -> f.getName() ==~ ASCIIDOC_REG_EXP_EXTENSION} as FileFilter)
                assert htmls.length == asciidocs.length
            }
            File[] actualChildren =  actualFile.listFiles(DIRECTORY_FILTER)
            assertEqualsStructure(expectedChildren, actualChildren)
        }
    }


    /**
     * Tests the behaviour when:
     *  - simple paths are used
     *  - preserveDirectories = true
     *  - relativeBaseDir = true
     *
     *  Expected:
     *   - all documents are rendered in the same folder structure found in the sourceDirectory
     *   - all documents are correctly rendered with the import
     */
    def 'should replicate source structure-standard paths'() {
        setup:
            File srcDir = new File('src/test/resources/src/asciidoctor/relative-path-treatment')
            File outputDir = new File('target/asciidoctor-output-relative')

        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html5'
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.imagesDir = '.'
            mojo.preserveDirectories = true
            mojo.relativeBaseDir = true
            mojo.sourceHighlighter = 'prettify'
            mojo.attributes = ['icons':'font']
            mojo.execute()

        then:
            outputDir.list().toList().isEmpty() == false
            assertEqualsStructure(srcDir.listFiles(DIRECTORY_FILTER), outputDir.listFiles(DIRECTORY_FILTER))
            def asciidocs = []
            outputDir.eachFileRecurse(FileType.FILES) {
                if (it.getName() ==~ /.+html/) asciidocs << it
            }
            asciidocs.size() == 6
            // Checks that all imports are found in the respective baseDir
            for (File renderedFile in asciidocs) {
                assert renderedFile.text.contains('Unresolved directive') == false
            }
        cleanup:
            // Avoids false positives in other tests
            FileUtils.deleteDirectory(outputDir)
    }

    /**
     * Tests the behaviour when:
     *  - complex paths are used
     *  - preserveDirectories = true
     *  - relativeBaseDir = true
     *
     *  Expected:
     *   - all documents are rendered in the same folder structure found in the sourceDirectory
     *   - all documents are correctly rendered with the import
     */
    def 'should replicate source structure-complex paths'() {
        setup:
            File srcDir = new File('src/test/resources/src/asciidoctor/relative-path-treatment/../relative-path-treatment')
            File outputDir = new File('target/../target/asciidoctor-output-relative')

        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html5'
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.preserveDirectories = true
            mojo.relativeBaseDir = true
            mojo.sourceHighlighter = 'coderay'
            mojo.attributes = ['icons':'font']
            mojo.execute()

        then:
            outputDir.list().toList().isEmpty() == false
            outputDir.listFiles({File f -> f.getName().endsWith('html')} as FileFilter).length == 1
            assertEqualsStructure(srcDir.listFiles(DIRECTORY_FILTER), outputDir.listFiles(DIRECTORY_FILTER))
            def asciidocs = []
            outputDir.eachFileRecurse(FileType.FILES) {
                if (it.getName() ==~ /.+html/) asciidocs << it
            }
            asciidocs.size() == 6
            // Checks that all imports are found in the respective baseDir
            for (File renderedFile in asciidocs) {
                assert renderedFile.text.contains('Unresolved directive') == false
            }
        cleanup:
            // Avoid possible false positives in other tests
            FileUtils.deleteDirectory(outputDir)
    }

    /**
     * Tests the behaviour when:
     *  - complex paths are used
     *  - preserveDirectories = false
     *  - relativeBaseDir = false
     *
     *  Expected:
     *   - all documents are rendered in the same outputDirectory. 1 document is overwritten
     *   - all documents but 1 (in the root) are incorrectly rendered because they cannot find the imported file
     */
    def 'should not replicate source structure-complex paths'() {
        setup:
            File srcDir = new File('src/test/resources/src/asciidoctor/relative-path-treatment/../relative-path-treatment')
            File outputDir = new File('target/../target/asciidoctor-output-relative')

        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html5'
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.sourceHighlighter = 'coderay'
            mojo.execute()

        then:
            outputDir.list().toList().isEmpty() == false
            // 1 file is missing because 2 share the same name and 1 is overwritten in outputDirectory
            def asciidocs  = outputDir.listFiles({File f -> f.getName().endsWith('html')} as FileFilter)
			asciidocs.length == 5
            // folders are copied anyway
            assertEqualsStructure(srcDir.listFiles(DIRECTORY_FILTER), outputDir.listFiles(DIRECTORY_FILTER))
			for (File renderedFile in asciidocs) {
				if (renderedFile.getName() != 'HelloWorld.html') {
					assert renderedFile.text.contains('Unresolved directive')
				}
			}
        cleanup:
            // Avoid possible false positives in other tests
            FileUtils.deleteDirectory(outputDir)
    }

    /**
     * Tests the behaviour when:
     *  - simple paths are used
     *  - preserveDirectories = true
     *  - relativeBaseDir = false
     *
     *  Expected:
     *   - all documents are rendered in the same folder structure found in the sourceDirectory
     *   - all documents but 1 (in the root) are incorrectly rendered because they cannot find the imported file
     */
    def 'should replicate source structure-no baseDir rewrite'() {
        setup:
            File srcDir = new File('src/test/resources/src/asciidoctor/relative-path-treatment')
            File outputDir = new File('target/asciidoctor-output-relative')

        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html5'
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.imagesDir = '.'
            mojo.preserveDirectories = true
			mojo.baseDir = srcDir
            //mojo.relativeBaseDir = true
            mojo.sourceHighlighter = 'prettify'
            mojo.attributes = ['icons':'font']
            mojo.execute()

        then:
            outputDir.list().toList().isEmpty() == false
            assertEqualsStructure(srcDir.listFiles(DIRECTORY_FILTER), outputDir.listFiles(DIRECTORY_FILTER))
            def asciidocs = []
            outputDir.eachFileRecurse(FileType.FILES) {
                if (it.getName() ==~ /.+html/) asciidocs << it
            }
            asciidocs.size() == 6
            // Looks for import errors in all files but the one in the root folder
            for (File renderedFile in asciidocs) {
                if (renderedFile.getName() != 'HelloWorld.html') {
                    assert renderedFile.text.contains('Unresolved directive')
                }
            }

        cleanup:
            // Avoids false positives in other tests
            FileUtils.deleteDirectory(outputDir)
    }

    /**
     * Tests the behaviour when:
     *  - simple paths are used
     *  - preserveDirectories = false
     *  - relativeBaseDir = true
     *
     *  Expected: all documents are correctly rendered in the same folder
     */
    def 'should not replicate source structure-baseDir rewrite'() {
        setup:
            File srcDir = new File('src/test/resources/src/asciidoctor/relative-path-treatment')
            File outputDir = new File('target/asciidoctor-output-relative')

        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.imagesDir = '.'
            mojo.preserveDirectories = false
            mojo.relativeBaseDir = true
            mojo.sourceHighlighter = 'prettify'
            mojo.attributes = ['icons':'font']
            mojo.execute()

        then:
            assertEqualsStructure(srcDir.listFiles(DIRECTORY_FILTER), outputDir.listFiles(DIRECTORY_FILTER))
			// all files are rendered in the outputDirectory
            def asciidocs = outputDir.listFiles({File f -> f.getName().endsWith('html')} as FileFilter)
			// 1 file is missing because 2 share the same name and 1 is overwritten in outputDirectory
            asciidocs.length == 5
            // Checks that all imports are found correctly because baseDir is adapted for each file
            for (File renderedFile in asciidocs) {
                assert renderedFile.text.contains('Unresolved directive') == false
            }

        cleanup:
            // Avoids false positives in other tests
            FileUtils.deleteDirectory(outputDir)
    }

    def 'project-version test'()
    {
        given:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File( 'target/asciidoctor-output-project-version-test' )

            if (!outputDir.exists()) {
                outputDir.mkdir()
            }
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.sourceDocumentName = new File( 'project-version.adoc' )
            mojo.backend = 'html'
            mojo.attributes['project-version'] = "1.0-SNAPSHOT"
            mojo.execute()
        then:
            File mainDocumentOutput = new File( outputDir, 'project-version.html' )
            String text = mainDocumentOutput.getText()
            assert text =~ /[vV]ersion 1\.0-SNAPSHOT/
            text.contains( "This is the project version: 1.0-SNAPSHOT" )
    }

    def 'github files can be included'() {
        setup:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File("target/test-resources/github-include/${System.currentTimeMillis()}")
            String documentName = 'github-include.adoc'

        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDocumentName = documentName
            mojo.backend = 'html5'
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.sourceHighlighter = 'coderay'
            mojo.attributes = ['allow-uri-read':'true']
            mojo.resources = [[
                                  directory: '.',
                                  excludes : ['**/**']
                          ] as Resource]
            mojo.execute()

        then:
            outputDir.list().toList().isEmpty() == false
            (new File(outputDir, 'github-include.html').text.contains('modelVersion'))
    }

    def "command line attributes replace configurations"() {
        setup:
            File srcDir = new File(DEFAULT_SOURCE_DIRECTORY)
            File outputDir = new File('target/asciidoctor-output/command-line-options')

            if (!outputDir.exists())
                outputDir.mkdirs()
        when: 'set toc and sourceHighlighter as XML configuration and command line attributes'
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.sourceDocumentName = 'sample.asciidoc'
            mojo.outputDirectory = outputDir
            mojo.sourceHighlighter = 'coderay'
            mojo.attributes['toc'] = 'left'
            // replace some options
            mojo.attributesChain = 'toc=right source-highlighter=highlight.js'
            mojo.execute()
        then: 'command line options are applied instead of xml configuration'
            outputDir.list().toList().isEmpty() == false
            outputDir.list().toList().contains('sample.html')
            File sampleOutput = new File('sample.html', outputDir)
            sampleOutput.length() > 0
            String text = sampleOutput.getText()
            text.contains('<body class="article toc2 toc-right">')
            text.contains('<pre class="highlightjs highlight">')
    }

    def "should skip processing when source directory does no exist"() {
        setup:
            def originalOut = System.out
            def newOut = new ByteArrayOutputStream()
            System.setOut(new PrintStream(newOut))

            File outputDir = new File("$MULTIPLE_RESOURCES_OUTPUT/skipped-process/${System.currentTimeMillis()}")
            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = new File(UUID.randomUUID().toString())
            mojo.backend = 'html5'
            mojo.outputDirectory = outputDir
            mojo.execute()
        then:
            newOut.toString().contains("sourceDirectory ${mojo.sourceDirectory} does not exist. Skip processing")
            !outputDir.exists()
        cleanup:
            System.setOut(originalOut)
    }

    def "should only convert documents and not copy any resources when resources directory does no exist"() {
        setup:
            File outputDir = new File("$MULTIPLE_RESOURCES_OUTPUT/multi-sources/error-source-not-found/${System.currentTimeMillis()}")

            if (!outputDir.exists())
                outputDir.mkdir()
        when: 'resource directory does not exist but source AsciiDoc documents do'
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = new File(DEFAULT_SOURCE_DIRECTORY)
            mojo.resources = [
                    [
                            directory: UUID.randomUUID().toString()
                    ] as Resource
            ]
            mojo.backend = 'html5'
            mojo.outputDirectory = outputDir
            mojo.execute()
        then: 'only rendered (html) files are found in the target directory'
            def allFiles = outputDir.listFiles({File f -> f.isFile()} as FileFilter)
            def htmlFiles = FileUtils.listFiles(outputDir, ['html'] as String[], true)
            allFiles.size() == htmlFiles.size()
            outputDir.listFiles({File f -> f.isDirectory()} as FileFilter).size() == 0

    }

    def "should only render a single file and not copy any resource"() {
        setup:
            File outputDir = new File("$MULTIPLE_RESOURCES_OUTPUT/file-pattern/${System.currentTimeMillis()}")

            if (!outputDir.exists())
                outputDir.mkdir()
            else
                FileUtils.deleteDirectory(outputDir)
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = new File(DEFAULT_SOURCE_DIRECTORY)
            mojo.sourceDocumentName = 'attribute-missing.adoc'
            // excludes all, nothing at all is copied
            mojo.resources = [[
                                      directory: DEFAULT_SOURCE_DIRECTORY,
                                      excludes : ['**/**']
                              ] as Resource]
            mojo.backend = 'html5'
            mojo.outputDirectory = outputDir
            mojo.execute()
        then:
            def files = outputDir.listFiles({File f -> f.isFile()} as FileFilter)
            files.size() == 1
            files*.name.containsAll(['attribute-missing.html'])
    }

    def "should copy all resources (2 directories with filters) into output folder"() {
        setup:
            File outputDir = new File("$MULTIPLE_RESOURCES_OUTPUT/multi-sources/${System.currentTimeMillis()}")
            String relativeTestsPath = "$DEFAULT_SOURCE_DIRECTORY/relative-path-treatment"

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = new File("$DEFAULT_SOURCE_DIRECTORY/issue-78")
            mojo.resources = [[
                                      directory: "$DEFAULT_SOURCE_DIRECTORY/issue-78",
                                      includes : ['**/*.adoc']
                              ] as Resource, [
                                      directory: relativeTestsPath,
                                      excludes : ['**/*.jpg']
                              ] as Resource]
            mojo.preserveDirectories = true
            mojo.backend = 'html5'
            mojo.outputDirectory = outputDir
            mojo.execute()
        then:
            def files = outputDir.listFiles({File f -> f.isFile()} as FileFilter)
            // includes 2 rendered AsciiDoc documents and 3 resources
            files.size() == 5
            // from 'issue-78' directory
            // resource files obtained using the include
            files*.name.findAll({it.endsWith('html')}).containsAll(['main.html', 'image-test.html'])
            // 'images' folder is not copied because it's not included
            files*.name.findAll({it == 'images'}) ==  []
            // from 'relative-path-treatment' directory
            // all folders and files are created because only image files are excluded
            assertEqualsStructure(new File(relativeTestsPath).listFiles(DIRECTORY_FILTER), outputDir.listFiles(DIRECTORY_FILTER))
            // images are excluded but not the rest of files
            FileUtils.listFiles(outputDir, ['groovy'] as String[], true).size() == 5
            FileUtils.listFiles(outputDir, ["jpg"] as String[], true).size() == 0
    }

    def "should render GitHub README alone"() {
        setup:
            File outputDir = new File("$MULTIPLE_RESOURCES_OUTPUT/readme/${System.currentTimeMillis()}")

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = new File('.')
            mojo.sourceDocumentName = 'README.adoc'
            mojo.resources = [[
                                      directory: ".",
                                      excludes : ['**/**']
                              ] as Resource]
            mojo.backend = 'html5'
            mojo.outputDirectory = outputDir
            mojo.execute()
        then:
            def files = outputDir.listFiles({File f -> f.isFile()} as FileFilter)
            // includes only 1 rendered AsciiDoc document
            files.size() == 1
            files.first().text.contains('Asciidoctor Maven Plugin')
    }

    def "should not include files in hidden directories (prefixes with underscore)"() {
        setup:
            File outputDir = new File("target/asciidoctor-output/hidden/${System.currentTimeMillis()}")
            String relativeTestsPath = "$DEFAULT_SOURCE_DIRECTORY/relative-path-treatment"

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = new File(relativeTestsPath)
            mojo.preserveDirectories = true
            mojo.backend = 'html5'
            mojo.outputDirectory = outputDir
            mojo.execute()
        then:
            def hiddenDirectories = ['_this_is_ignored', 'level-1-1/level-2-2/_this_is_ignored']
            hiddenDirectories. each { directoryPath ->
                assert new File(relativeTestsPath, directoryPath).exists()
                assert !(new File(outputDir, directoryPath).exists())
            }
    }

    def "should not crash when enabling maven-resource filtering"() {
        setup:
            File outputDir = new File("$MULTIPLE_RESOURCES_OUTPUT/readme/${System.currentTimeMillis()}")

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = new File('.')
            mojo.sourceDocumentName = 'README.adoc'
            mojo.resources = [[
                                      directory: ".",
                                      excludes : ['**/**'],
                                      filtering: true
                              ] as Resource]
            mojo.backend = 'html5'
            mojo.outputDirectory = outputDir
            mojo.execute()
        then:
            def files = outputDir.listFiles({File f -> f.isFile()} as FileFilter)
            // includes only 1 rendered AsciiDoc document
            files.size() == 1
            files.first().text.contains('Asciidoctor Maven Plugin')
    }

    def "should exclude custom source document Extensions by default"() {
        setup:
            File outputDir = new File("$MULTIPLE_RESOURCES_OUTPUT/readme/${System.currentTimeMillis()}")

            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = new File(DEFAULT_SOURCE_DIRECTORY)
            mojo.sourceDocumentExtensions = ['ext']
            mojo.resources = [[
                                      directory: DEFAULT_SOURCE_DIRECTORY,
                                      includes: ['**/*.adoc'],
                                      excludes: ['**/**']
                              ] as Resource]
            mojo.backend = 'html5'
            mojo.outputDirectory = outputDir
            mojo.execute()
        then:
            def files = outputDir.listFiles({File f -> f.isFile()} as FileFilter)
            FileUtils.listFiles(outputDir, ['ext'] as String[], true).isEmpty()
            // includes only 1 rendered AsciiDoc document
            def file = new File(outputDir, 'sample.html')
            file.text.contains('Asciidoctor default stylesheet')
    }

    def "should show message when overwriting files without outputFile"() {
        setup:
            def originalOut = System.out
            def newOut = new ByteArrayOutputStream()
            System.setOut(new PrintStream(newOut))
            def originalErr = System.err
            def newErr = new ByteArrayOutputStream()
            System.setErr(new PrintStream(newErr))

            // srcDir contains 6 documents, 2 of them with the same name (HellowWorld3.adoc)
            File srcDir = new File("$DEFAULT_SOURCE_DIRECTORY/relative-path-treatment/")
            File outputDir = new File("target/asciidoctor-output/overlapping-outputFile/${System.currentTimeMillis()}")
            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = srcDir
            mojo.backend = 'html5'
            mojo.outputDirectory = outputDir
            mojo.execute()
        then:
            def asciidocs = []
            outputDir.eachFileRecurse(FileType.FILES) {
                if (it.getName() ==~ /.+html/) asciidocs << it
            }
            asciidocs.size() == 5
            newOut.toString().count("Rendered") == 6
            newOut.toString().count("Duplicated destination found") == 1
        cleanup:
            System.setOut(originalOut)
            System.setErr(originalErr)
    }

    def "should show message when overwriting files using outputFile"() {
        setup:
            def originalOut = System.out
            def newOut = new ByteArrayOutputStream()
            System.setOut(new PrintStream(newOut))
            def originalErr = System.err
            def newErr = new ByteArrayOutputStream()
            System.setErr(new PrintStream(newErr))

            File srcDir = new File("$DEFAULT_SOURCE_DIRECTORY/relative-path-treatment/")
            File outputDir = new File("target/asciidoctor-output/overlapping-outputFile/${System.currentTimeMillis()}")
            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.sourceDirectory = srcDir
            mojo.backend = 'html5'
            mojo.outputDirectory = outputDir
            mojo.outputFile = new File('single-output.html')
            mojo.execute()
        then:
            def asciidocs = []
            outputDir.eachFileRecurse(FileType.FILES) {
                if (it.getName() ==~ /.+html/) asciidocs << it
            }
            asciidocs.size() == 1
            newOut.toString().count("Rendered") == 6
            newOut.toString().count("Duplicated destination found") == 5
        cleanup:
            System.setOut(originalOut)
            System.setErr(originalErr)
    }

    def "should not show message when overwriting files using outputFile and preserveDirectories"() {
        setup:
            def originalOut = System.out
            def newOut = new ByteArrayOutputStream()
            System.setOut(new PrintStream(newOut))
            def originalErr = System.err
            def newErr = new ByteArrayOutputStream()
            System.setErr(new PrintStream(newErr))

            File srcDir = new File("$DEFAULT_SOURCE_DIRECTORY/relative-path-treatment/")
            File outputDir = new File("target/asciidoctor-output/overlapping-outputFile/${System.currentTimeMillis()}")
            if (!outputDir.exists())
                outputDir.mkdir()
        when:
            AsciidoctorMojo mojo = new AsciidoctorMojo()
            mojo.getLog().errorEnabled
            mojo.sourceDirectory = srcDir
            mojo.backend = 'html5'
            mojo.outputDirectory = outputDir
            mojo.preserveDirectories = true
            mojo.outputFile = new File('single-output.html')
            mojo.execute()
        then:
            def asciidocs = []
            outputDir.eachFileRecurse(FileType.FILES) {
                if (it.getName() ==~ /.+html/) asciidocs << it
            }
            asciidocs.size() == 5
            newOut.toString().count("Rendered") == 6
            newOut.toString().count("Duplicated destination found") == 1
        cleanup:
            System.setOut(originalOut)
            System.setErr(originalErr)
    }

}
