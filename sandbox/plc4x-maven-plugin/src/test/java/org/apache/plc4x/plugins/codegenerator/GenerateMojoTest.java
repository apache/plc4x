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

import org.apache.maven.plugin.testing.MojoRule;
import org.apache.maven.plugin.testing.WithoutMojo;
import org.apache.maven.plugin.testing.resources.TestResources;
import org.apache.maven.project.MavenProject;
import org.junit.Rule;
import org.junit.Test;

import java.io.File;

import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;

public class GenerateMojoTest {

    @Rule
    public MojoRule rule = new MojoRule();

    @Rule
    public TestResources testResources = new TestResources();

    /**
     * @throws Exception if any
     */
    @Test
    public void testSomething() throws Exception {
        File baseDir = testResources.getBasedir("simple-embedded-schema");
        MavenProject project = rule.readMavenProject(baseDir);

        GenerateMojo generateMojo = (GenerateMojo) rule.lookupConfiguredMojo(project, "generate-driver");
        assertNotNull(generateMojo);
        generateMojo.execute();

        File outputDirectory = (File) rule.getVariableValueFromObject(generateMojo, "outputDir");
        assertNotNull(outputDirectory);
        assertTrue(outputDirectory.exists());
    }

    /**
     * Do not need the MojoRule.
     */
    @WithoutMojo
    @Test
    public void testSomethingWhichDoesNotNeedTheMojoAndProbablyShouldBeExtractedIntoANewClassOfItsOwn() {
        assertTrue(true);
    }

}

