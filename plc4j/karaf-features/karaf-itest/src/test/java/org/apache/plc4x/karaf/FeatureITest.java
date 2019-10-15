/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.apache.plc4x.karaf;

import org.apache.karaf.bundle.core.BundleService;
import org.apache.karaf.features.FeaturesService;
import org.apache.plc4x.java.spi.PlcDriver;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.ops4j.pax.exam.Configuration;
import org.ops4j.pax.exam.ConfigurationManager;
import org.ops4j.pax.exam.Option;
import org.ops4j.pax.exam.junit.PaxExam;
import org.ops4j.pax.exam.options.MavenArtifactUrlReference;
import org.ops4j.pax.exam.options.MavenUrlReference;
import org.osgi.framework.Bundle;
import org.osgi.framework.BundleContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.inject.Inject;
import java.io.File;
import java.util.Arrays;
import java.util.Optional;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;
import static org.ops4j.pax.exam.CoreOptions.bundle;
import static org.ops4j.pax.exam.CoreOptions.maven;
import static org.ops4j.pax.exam.CoreOptions.mavenBundle;
import static org.ops4j.pax.exam.karaf.options.KarafDistributionOption.configureConsole;
import static org.ops4j.pax.exam.karaf.options.KarafDistributionOption.features;
import static org.ops4j.pax.exam.karaf.options.KarafDistributionOption.karafDistributionConfiguration;
import static org.ops4j.pax.exam.karaf.options.KarafDistributionOption.keepRuntimeFolder;

@RunWith(PaxExam.class)
public class FeatureITest {

    private static Logger LOG = LoggerFactory.getLogger(FeatureITest.class);

    @Inject
    private FeaturesService featuresService;

    @Inject
    private BundleContext bundleContext;

    @Configuration
    public Option[] config() {
        MavenArtifactUrlReference karafUrl = maven()
            .groupId("org.apache.karaf")
            .artifactId("apache-karaf")
            .version(karafVersion())
            .type("zip");

        MavenUrlReference karafStandardRepo = maven()
            .groupId("org.apache.karaf.features")
            .artifactId("standard")
            .version(karafVersion())
            .classifier("features")
            .type("xml");

        final MavenArtifactUrlReference plc4xRepo = maven()
            .groupId("org.apache.plc4x")
            .artifactId("driver-s7-feature")
            .version("0.5.0-SNAPSHOT")
            .classifier("features")
            .type("xml");
        return new Option[] {
            // KarafDistributionOption.debugConfiguration("5005", true),
            karafDistributionConfiguration()
                .frameworkUrl(karafUrl)
                .unpackDirectory(new File("target", "exam"))
                .useDeployFolder(false),
            keepRuntimeFolder(),
            configureConsole().ignoreLocalConsole(),
            features(karafStandardRepo , "scr"),
            features(plc4xRepo, "driver-s7-feature"),
            mavenBundle()
                .groupId("org.apache.plc4x")
                .artifactId("plc4j-driver-s7")
                .version("0.5.0-SNAPSHOT")
                .start()
        };
    }

    public static String karafVersion() {
        ConfigurationManager cm = new ConfigurationManager();
        String karafVersion = cm.getProperty("pax.exam.karaf.version", "3.0.0");
        return karafVersion;
    }


    @Test
    public void checkFeatureInstalled() throws Exception {
        assertTrue(
            featuresService.isInstalled(featuresService.getFeature("driver-s7-feature"))
        );
    }

    @Test
    public void checkBundleStarted() throws Exception {
        for (Bundle bundle : bundleContext.getBundles()) {
            System.out.println(bundle.getSymbolicName());
        }

        // Try to find the bundle
        final Optional<Bundle> optionalBundle = Arrays.stream(bundleContext.getBundles())
            .filter(bundle -> "org.apache.plc4x.plc4j-driver-s7".equals(bundle.getSymbolicName()))
            .findFirst();

        // Ensure that the bundle is resolved
        assertTrue(optionalBundle.isPresent());

        // Check if the bundle is active
        assertEquals(Bundle.ACTIVE, optionalBundle.get().getState());
    }

}
