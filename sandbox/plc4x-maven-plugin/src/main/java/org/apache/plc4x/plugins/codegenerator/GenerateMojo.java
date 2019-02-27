package org.apache.plc4x.plugins.codegenerator;


import org.apache.maven.artifact.DependencyResolutionRequiredException;
import org.apache.maven.plugin.AbstractMojo;
import org.apache.maven.plugin.MojoExecutionException;
import org.apache.maven.plugins.annotations.LifecyclePhase;
import org.apache.maven.plugins.annotations.Mojo;
import org.apache.maven.plugins.annotations.Parameter;
import org.apache.maven.plugins.annotations.ResolutionScope;
import org.apache.maven.project.MavenProject;

import java.io.File;

/**
 * Generate the types, serializer and parser classes based on a DFDL shema.
 */
@Mojo(name = "generate-driver",
    defaultPhase = LifecyclePhase.GENERATE_SOURCES,
    requiresDependencyResolution = ResolutionScope.RUNTIME)
public class GenerateMojo extends AbstractMojo {

    @Parameter(defaultValue = "${project}", readonly = true, required = true)
    private MavenProject project;

    /**
     * Where to generate the code.
     */
    @Parameter(defaultValue = "${project.build.directory}/generated-sources/plc4x/", property = "outputDir", required = true)
    private File outputDirectory;

    /**
     * The path to the DFDL schema file in the modules classpath (Ideally as runtime dependency)
     */
    @Parameter(property = "dfdl-schema", required = true)
    private File dfdlSchemaFile;

    public void execute()
        throws MojoExecutionException {

        // Make sure the output directory exists.
        if (!outputDirectory.exists()) {
            if(!outputDirectory.mkdirs()) {
                throw new MojoExecutionException("Could not generate output directory " + outputDirectory.getAbsolutePath());
            }
        }

        try {
            for (String runtimeClasspathElement : project.getRuntimeClasspathElements()) {
                getLog().info("Inspecting " + runtimeClasspathElement);
            }
        } catch (DependencyResolutionRequiredException e) {
            throw new MojoExecutionException("Error resolving dependencies", e);
        }
    }

}
