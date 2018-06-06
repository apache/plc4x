package org.apache.plc4x.utils.maven.site.asciidoctor.test.site

import org.apache.maven.doxia.sink.Sink
import org.apache.maven.project.MavenProject
import org.apache.plc4x.utils.maven.site.asciidoctor.site.AsciidoctorDoxiaParser
import org.codehaus.plexus.util.xml.Xpp3DomBuilder
import spock.lang.Specification

class AsciidoctorDoxiaParserTest extends Specification {

    public static final String TEST_DOCS_PATH = 'src/test/resources/src/asciidoctor'

    def "should render html"() {
        given:
        final File srcAsciidoc = new File("$TEST_DOCS_PATH/sample.asciidoc")
        final Sink sink = createSinkMock()

        AsciidoctorDoxiaParser parser = new AsciidoctorDoxiaParser()
        parser.@project = createMavenProjectMock()

        when:
        parser.parse(new FileReader(srcAsciidoc), sink)

        then:
        String outputText = sink.sinkedText
        outputText.contains '<h1>Document Title</h1>'
        outputText.contains '<div class="ulist">'
        outputText.contains '<div class="listingblock">'
        outputText.contains "require 'asciidoctor'"
        // icons as text
        outputText.contains '<div class="title">Note</div>'
    }

    def "should render html with an attribute"() {
        given:
        final File srcAsciidoc = new File("$TEST_DOCS_PATH/sample.asciidoc")
        Reader reader = new FileReader(srcAsciidoc)
        final Sink sink = createSinkMock()

        AsciidoctorDoxiaParser parser = new AsciidoctorDoxiaParser()
        parser.@project = createMavenProjectMock('''
                    <configuration>
                        <asciidoc>
                            <attributes>
                                <icons>font</icons>
                            </attributes>
                        </asciidoc>
                    </configuration>''')

        when:
        parser.parse(reader, sink)

        then:
        String outputText = sink.sinkedText
        // :icons: font
        outputText.contains '<i class="fa icon-note" title="Note"></i>'
    }

    def "should render html with baseDir option"() {
        given:
        final File srcAsciidoc = new File("$TEST_DOCS_PATH/main-document.adoc")
        final Sink sink = createSinkMock()

        AsciidoctorDoxiaParser parser = new AsciidoctorDoxiaParser()
        parser.@project = createMavenProjectMock("""
                     <configuration>
                        <asciidoc>
                            <baseDir>${new File(srcAsciidoc.parent).absolutePath}</baseDir>
                        </asciidoc>
                     </configuration>""")

        when:
        parser.parse(new FileReader(srcAsciidoc), sink)

        then: 'include works'
        String outputText = sink.sinkedText
        outputText.contains '<h1>Include test</h1>'
        outputText.contains 'println "HelloWorld from Groovy on ${new Date()}"'
    }

    def "should render html with relative baseDir option"() {
        given:
        final File srcAsciidoc = new File("$TEST_DOCS_PATH/main-document.adoc")
        final Sink sink = createSinkMock()

        AsciidoctorDoxiaParser parser = new AsciidoctorDoxiaParser()
        parser.@project = createMavenProjectMock("""
                     <configuration>
                        <asciidoc>
                            <baseDir>${TEST_DOCS_PATH}</baseDir>
                        </asciidoc>
                     </configuration>""")

        when:
        parser.parse(new FileReader(srcAsciidoc), sink)

        then: 'include works'
        String outputText = sink.sinkedText
        outputText.contains '<h1>Include test</h1>'
        outputText.contains 'println "HelloWorld from Groovy on ${new Date()}"'
    }

    def "should render html with templateDir option"() {
        given:
        final File srcAsciidoc = new File("$TEST_DOCS_PATH/sample.asciidoc")
        final Sink sink = createSinkMock()

        AsciidoctorDoxiaParser parser = new AsciidoctorDoxiaParser()
        parser.@project = createMavenProjectMock("""
                     <configuration>
                        <asciidoc>
                            <templateDir>${TEST_DOCS_PATH}/templates</templateDir>
                        </asciidoc>
                     </configuration>""")

        when:
        parser.parse(new FileReader(srcAsciidoc), sink)

        then:
        String outputText = sink.sinkedText
        outputText.contains '<h1>Document Title</h1>'
        outputText.contains '<p class="custom-template ">'
    }

    def "should render html with attributes and baseDir option"() {
        given:
        final File srcAsciidoc = new File("$TEST_DOCS_PATH/main-document.adoc")
        final Sink sink = createSinkMock()

        AsciidoctorDoxiaParser parser = new AsciidoctorDoxiaParser()
        parser.@project = createMavenProjectMock("""
                    <configuration>
                        <asciidoc>
                            <baseDir>${new File(srcAsciidoc.parent).absolutePath}</baseDir>
                            <attributes>
                                <sectnums></sectnums>
                                <icons>font</icons>
                                <my-label>Hello World!!</my-label>
                            </attributes>
                        </asciidoc>
                    </configuration>""")

        when:
        parser.parse(new FileReader(srcAsciidoc), sink)

        then:
        String outputText = sink.sinkedText
        outputText.contains '<h1>Include test</h1>'
        outputText.contains '<h2 id="code">1. Code</h2>'
        outputText.contains '<h2 id="optional_section">2. Optional section</h2>'
        outputText.contains 'println "HelloWorld from Groovy on ${new Date()}"'
        outputText.contains 'Hello World!!'
        outputText.contains '<i class="fa icon-tip" title="Tip"></i>'
    }

    def "should process empty self-closing XML attributes"() {
        given:
        final File srcAsciidoc = new File("$TEST_DOCS_PATH/sample.asciidoc")
        final Sink sink = createSinkMock()

        AsciidoctorDoxiaParser parser = new AsciidoctorDoxiaParser()
        parser.@project = createMavenProjectMock("""
                     <configuration>
                       <asciidoc>
                         <attributes>
                           <sectnums/>
                         </attributes>
                       </asciidoc>
                     </configuration>""")

        when:
        parser.parse(new FileReader(srcAsciidoc), sink)

        then:
        String outputText = sink.sinkedText
        outputText.contains '<h2 id="id_section_a">1. Section A</h2>'
        outputText.contains '<h3 id="id_section_a_subsection">1.1. Section A Subsection</h3>'
    }

    def "should process empty value XML attributes"() {
        given:
        final File srcAsciidoc = new File("$TEST_DOCS_PATH/sample.asciidoc")
        final Sink sink = createSinkMock()

        AsciidoctorDoxiaParser parser = new AsciidoctorDoxiaParser()
        parser.@project = createMavenProjectMock("""
                     <configuration>
                       <asciidoc>
                         <attributes>
                           <sectnums></sectnums>
                         </attributes>
                       </asciidoc>
                     </configuration>""")

        when:
        parser.parse(new FileReader(srcAsciidoc), sink)

        then:
        String outputText = sink.sinkedText
        outputText.contains '<h2 id="id_section_a">1. Section A</h2>'
        outputText.contains '<h3 id="id_section_a_subsection">1.1. Section A Subsection</h3>'
    }

    private MavenProject createMavenProjectMock(final String configuration = null) {
        [getGoalConfiguration: { pluginGroupId, pluginArtifactId, executionId, goalId ->
            configuration ? Xpp3DomBuilder.build(new StringReader(configuration)) : null
        },
         getBasedir          : {
             new File('.')
         }] as MavenProject
    }

    /**
     * Creates a {@link Sink} mock that allows retrieving a text previously sinked.
     */
    private Sink createSinkMock() {
        String text
        return [rawText : { t ->
            text = t
        }, getSinkedText: {
            text
        }] as MySink
    }

    interface MySink extends Sink {
        String getSinkedText()
    }

}
