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
package org.apache.plc4x.karaf;

import org.apache.plc4x.java.spi.PlcDriver;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.ops4j.pax.exam.Configuration;
import org.ops4j.pax.exam.Option;
import org.ops4j.pax.exam.junit.PaxExam;
import org.ops4j.pax.exam.options.MavenArtifactUrlReference;
import org.ops4j.pax.exam.options.MavenUrlReference;
import org.ops4j.pax.exam.spi.reactors.ExamReactorStrategy;
import org.ops4j.pax.exam.spi.reactors.PerClass;
import org.ops4j.pax.exam.spi.reactors.PerMethod;

import javax.inject.Inject;

import java.io.File;

import static org.junit.Assert.assertNotNull;
import static org.ops4j.pax.exam.CoreOptions.bundle;
import static org.ops4j.pax.exam.CoreOptions.junitBundles;
import static org.ops4j.pax.exam.CoreOptions.maven;
import static org.ops4j.pax.exam.CoreOptions.mavenBundle;
import static org.ops4j.pax.exam.CoreOptions.options;
import static org.ops4j.pax.exam.karaf.options.KarafDistributionOption.features;
import static org.ops4j.pax.exam.karaf.options.KarafDistributionOption.karafDistributionConfiguration;
import static org.ops4j.pax.exam.karaf.options.KarafDistributionOption.keepRuntimeFolder;

@RunWith(PaxExam.class)
@ExamReactorStrategy(PerMethod.class)
public class FeatureTest {

    @Inject
    private PlcDriver plcDriver;

    @Configuration
    public Option[] config() {
        MavenArtifactUrlReference karafUrl = maven()
            .groupId("org.apache.karaf")
            .artifactId("apache-karaf")
            .version("3.0.0")
            .type("tar.gz");
        MavenUrlReference karafStandardRepo = maven()
            .groupId("org.apache.karaf.features")
            .artifactId("standard")
            .classifier("features")
            .version("4.2.1")
            .type("xml")
            .versionAsInProject();
        return new Option[] {
            // KarafDistributionOption.debugConfiguration("5005", true),
            karafDistributionConfiguration()
                .frameworkUrl(karafUrl)
                .unpackDirectory(new File("target/exam"))
                .useDeployFolder(false),
            keepRuntimeFolder(),
            features(karafStandardRepo , "scr"),
            mavenBundle()
                .groupId("org.ops4j.pax.exam.samples")
                .artifactId("pax-exam-sample8-ds")
                .version("4.2.1")
                //.versionAsInProject()
            .start(),
        };
    }

    @Test
    public void getHelloService() {
        assertNotNull(plcDriver);
        // assertEquals("Hello Pax!", helloService.getMessage());
    }

}
