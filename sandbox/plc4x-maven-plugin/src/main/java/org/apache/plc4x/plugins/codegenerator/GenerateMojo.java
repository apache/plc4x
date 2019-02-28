/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
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
    private File outputDir;

    /**
     * The path to the DFDL schema file in the modules classpath (Ideally as runtime dependency)
     */
    @Parameter(property = "dfdl-schema", required = true)
    private String dfdlSchema;

    public void execute()
        throws MojoExecutionException {

        // Make sure the output directory exists.
        if (!outputDir.exists()) {
            if (!outputDir.mkdirs()) {
                throw new MojoExecutionException("Could not generate output directory " + outputDir.getAbsolutePath());
            }
        }

        try {
            File schemaFile = new File(dfdlSchema);
            if(!schemaFile.exists()) {
                for (String runtimeClasspathElement : project.getRuntimeClasspathElements()) {
                    File classPathElement = new File(runtimeClasspathElement);
                    if (classPathElement.exists()) {
                        // This is usually the "target/classes" directory.
                        if (classPathElement.isDirectory()) {
                            schemaFile = new File(classPathElement, dfdlSchema);
                            if (schemaFile.exists()) {
                                break;
                            }
                        }
                        // These are usually jar files.
                        else {
                            // TODO: Implement ...
                        }
                    }
                }
            }
        } catch (DependencyResolutionRequiredException e) {
            throw new MojoExecutionException("Error resolving dependencies", e);
        }
    }

}
