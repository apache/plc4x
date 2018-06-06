package org.apache.plc4x.utils.maven.site.asciidoctor.test.plexus

import org.apache.maven.plugin.logging.Log
import org.apache.maven.plugin.logging.SystemStreamLog
import org.apache.maven.project.MavenProject
import org.apache.maven.shared.filtering.DefaultMavenFileFilter
import org.apache.maven.shared.filtering.DefaultMavenResourcesFiltering
import org.apache.plc4x.utils.maven.site.asciidoctor.AsciidoctorMojo
import org.sonatype.plexus.build.incremental.DefaultBuildContext

/**
 * Initializes a mocks components required to simulate a Plexus container for the tests.
 *
 * Note: This is a temporal solution. At some point 'proper' testing should be introduced, see:
 *  - https://vzurczak.wordpress.com/2014/07/23/write-unit-tests-for-a-maven-plug-in/
 *
 * @author abelsromero
 */
class MockPlexusContainer {

    class FakeMavenLogger {
        @Delegate
        Log logger = new SystemStreamLog()
    }

    void initializeContext(AsciidoctorMojo mojo) {

        mojo.@project = [
                getBasedir: {
                    return new File('.')
                }] as MavenProject

        mojo.@buildContext = new DefaultBuildContext()

        def logger = new FakeMavenLogger() as org.codehaus.plexus.logging.Logger

        DefaultMavenFileFilter mavenFileFilter = new DefaultMavenFileFilter()
        mavenFileFilter.@buildContext = mojo.@buildContext
        mavenFileFilter.enableLogging(logger)

        DefaultMavenResourcesFiltering resourceFilter = new DefaultMavenResourcesFiltering()
        resourceFilter.@mavenFileFilter = mavenFileFilter
        resourceFilter.@buildContext = mojo.@buildContext
        resourceFilter.initialize()
        resourceFilter.enableLogging(logger)
        mojo.encoding = "UTF-8"
        mojo.@outputResourcesFiltering = resourceFilter

    }

}
