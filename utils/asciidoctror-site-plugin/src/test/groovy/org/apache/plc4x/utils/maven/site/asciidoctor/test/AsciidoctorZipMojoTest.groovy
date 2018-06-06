package org.apache.plc4x.utils.maven.site.asciidoctor.test

import org.apache.commons.io.FileUtils
import org.apache.plc4x.utils.maven.site.asciidoctor.AsciidoctorZipMojo
import org.apache.plc4x.utils.maven.site.asciidoctor.test.plexus.MockPlexusContainer
import spock.lang.Specification

import java.util.zip.ZipFile

/**
 *
 */
class AsciidoctorZipMojoTest extends Specification {

    /**
     * Intercept Asciidoctor mojo constructor to mock and inject required
     * plexus objects
     */
    def setupSpec() {
        MockPlexusContainer mockPlexusContainer = new MockPlexusContainer()
        def oldConstructor = AsciidoctorZipMojo.constructors[0]

        AsciidoctorZipMojo.metaClass.constructor = {
            def mojo = oldConstructor.newInstance()
            mockPlexusContainer.initializeContext(mojo)
            return mojo
        }
    }

    def "zip it"() {
        given: 'an empty output directory'
            def outputDir = new File('target/asciidoctor-zip-output')
            outputDir.deleteDir()
            outputDir.mkdirs()

            def zip = new File('target/asciidoctor-zip.zip')
            zip.delete()

        when: 'zip mojo is called'
            def srcDir = new File('target/test-classes/src/asciidoctor-zip')
            srcDir.mkdirs()

            new File(srcDir, "sample.asciidoc").withWriter {
                it << '''
                Title
                =====
                test
                '''.stripIndent()
            }

            def mojo = new AsciidoctorZipMojo()
            mojo.backend = 'html'
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.zipDestination = zip
            mojo.zip = true
            mojo.execute()

        then: 'a zip is created'
            mojo.zipDestination.exists()

            def entries = new ZipFile(mojo.zipDestination).entries()
           entries.hasMoreElements()
           def entryName = entries.nextElement().name
           entryName == 'asciidoctor-zip/target/asciidoctor-zip-output/sample.html' || entryName == 'asciidoctor-zip/target\\asciidoctor-zip-output\\sample.html'
    }

    def 'should replicate source structure in zip-standard paths'() {
        setup:
            File srcDir = new File('src/test/resources/src/asciidoctor/relative-path-treatment')
            // Create random folder to avoid mix files when something goes wrong
            String outputPath = "target/asciidoctor-output-relative-${UUID.randomUUID().toString().split('-') [0]}/"
            File outputDir = new File(outputPath)
            
            File zip = new File("$outputPath/asciidoctor-zip.zip")

        when:
            AsciidoctorZipMojo mojo = new AsciidoctorZipMojo()
            mojo.backend = 'html5'
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.preserveDirectories = true
            mojo.relativeBaseDir = true
            mojo.sourceHighlighter = 'coderay'
            mojo.zipDestination = zip
            mojo.zip = true
            mojo.execute()

        then:
            ZipFile zipfile = new ZipFile(zip)
            def entries = (zipfile).entries()*.getName().collect() {
                // Protection to avoid errors on diferent OS
                it.replaceAll('\\\\', '/') - outputPath
            }
            zipfile.close()
            // Paths are adapted for the test are do not match the real paths inside the zip
            def expected = [
                'asciidoctor-zip/HelloWorld.groovy',
                'asciidoctor-zip/HelloWorld.html',
                'asciidoctor-zip/level-1-1/asciidoctor-icon.jpg',
                'asciidoctor-zip/level-1-1/HelloWorld2.groovy',
                'asciidoctor-zip/level-1-1/HelloWorld2.html',
                'asciidoctor-zip/level-1-1/HelloWorld22.html',
                'asciidoctor-zip/level-1-1/level-2-1/HelloWorld3.groovy',
                'asciidoctor-zip/level-1-1/level-2-1/HelloWorld3.html',
                'asciidoctor-zip/level-1-1/level-2-2/HelloWorld3.groovy',
                'asciidoctor-zip/level-1-1/level-2-2/HelloWorld3.html',
                'asciidoctor-zip/level-1-1/level-2-2/level-3-1/HelloWorld4.groovy',
                'asciidoctor-zip/level-1-1/level-2-2/level-3-1/HelloWorld4.html'
            ]
            expected.containsAll(entries)
            entries.containsAll(expected)

        cleanup:
            FileUtils.deleteDirectory(outputDir)
    }

    def 'should not replicate source structure in zip-standard paths'() {
        setup:
            File srcDir = new File('src/test/resources/src/asciidoctor/relative-path-treatment')
            String outputPath = "target/asciidoctor-output-relative-${UUID.randomUUID().toString().split('-') [0]}/"
            File outputDir = new File(outputPath)

            File zip = new File("$outputPath/asciidoctor-zip.zip")

        when:
            AsciidoctorZipMojo mojo = new AsciidoctorZipMojo()
            mojo.backend = 'html5'
            mojo.sourceDirectory = srcDir
            mojo.outputDirectory = outputDir
            mojo.sourceHighlighter = 'coderay'
            mojo.zipDestination = zip
            mojo.zip = true
            mojo.execute()

        then:
            ZipFile zipfile = new ZipFile(zip)
            def entries = (zipfile).entries()*.getName().collect() {
                // Protection to avoid errors on diferent OS
                it.replaceAll('\\\\', '/') - outputPath
            }
            zipfile.close()
            // Paths are adapted for the test are do not match the real paths inside the zip
            def expected = [
                'asciidoctor-zip/HelloWorld.groovy',
                'asciidoctor-zip/HelloWorld.html',
                'asciidoctor-zip/HelloWorld2.html',
                'asciidoctor-zip/HelloWorld22.html',
                'asciidoctor-zip/HelloWorld3.html',
                'asciidoctor-zip/HelloWorld4.html',
                'asciidoctor-zip/level-1-1/asciidoctor-icon.jpg',
                'asciidoctor-zip/level-1-1/HelloWorld2.groovy',
                'asciidoctor-zip/level-1-1/level-2-1/HelloWorld3.groovy',
                'asciidoctor-zip/level-1-1/level-2-2/HelloWorld3.groovy',
                'asciidoctor-zip/level-1-1/level-2-2/level-3-1/HelloWorld4.groovy'
            ]
            expected.containsAll(entries)
            entries.containsAll(expected)

        cleanup:
            FileUtils.deleteDirectory(outputDir)

    }

}
