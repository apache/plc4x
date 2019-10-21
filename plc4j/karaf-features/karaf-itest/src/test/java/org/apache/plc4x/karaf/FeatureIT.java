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

import org.apache.karaf.itests.KarafTestSupport;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.ops4j.pax.exam.Configuration;
import org.ops4j.pax.exam.Option;
import org.ops4j.pax.exam.junit.PaxExam;
import org.ops4j.pax.exam.karaf.options.KarafDistributionOption;
import org.ops4j.pax.exam.karaf.options.LogLevelOption;
import org.ops4j.pax.exam.options.MavenArtifactUrlReference;
import org.ops4j.pax.exam.spi.reactors.ExamReactorStrategy;
import org.ops4j.pax.exam.spi.reactors.PerMethod;
import org.osgi.framework.Bundle;

import java.util.stream.Stream;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.ops4j.pax.exam.CoreOptions.maven;
import static org.ops4j.pax.exam.karaf.options.KarafDistributionOption.logLevel;

@RunWith(PaxExam.class)
@ExamReactorStrategy(PerMethod.class)
public class FeatureIT extends KarafTestSupport {

    @Configuration
    public Option[] config() {
        Option[] options = new Option[]{
            logLevel(LogLevelOption.LogLevel.DEBUG),
            KarafDistributionOption.editConfigurationFilePut("etc/system.properties", "plc4x.version", System.getProperty("plc4x.version")),
        };
        return Stream.of(super.config(), options).flatMap(Stream::of).toArray(Option[]::new);
    }

    /**
     * Checks:
     * - the feature can be installed
     * - the Driver Bundle is there and active
     * - the Bundle contains the DS for PlcDriver
     * - the feature is installed
     */
    @Test
    public void installAndAssertFeature() throws Exception {
        // Assemble Maven URL
        final MavenArtifactUrlReference featureRepo = maven()
            .groupId("org.apache.plc4x")
            .artifactId("driver-s7-feature")
            .version(System.getProperty("plc4x.version"))
            .type("xml")
            .classifier("features");
        System.out.println("Installing feature repo " + featureRepo.getURL());

        // Install the feature-repo
        addFeaturesRepository(featureRepo.getURL());

        // Install the feature
        installAndAssertFeature("driver-s7-feature");

        // Print Bundles and fetch result
        String bundles = executeCommand("bundle:list -t 0");

        // Find that line
        // 84 │ Active │  80 │ 0.5.0.SNAPSHOT │ PLC4J: Driver: S7
        assertContains("PLC4J: Driver: S7", bundles);
        System.out.println(bundles);

        // Find Bundle for more detailed check
        final Bundle bundle = findBundleByName("org.apache.plc4x.plc4j-driver-s7");

        // Bundle has to be ACTIVE
        assertNotNull(bundle);
        assertEquals(Bundle.ACTIVE, bundle.getState());

        // Check declarative service is present
        String services = executeCommand("services -p " + bundle.getBundleId());
        System.out.println("Services: " + services);
        assertContains("component.name = org.apache.plc4x.java.s7.S7PlcDriver", services);
        assertContains("objectClass = [org.apache.plc4x.java.spi.PlcDriver]", services);

        // Just for Debugging...
        String features = executeCommand("feature:list -i");
        System.out.print(features);
        assertContains("driver-s7-feature", features);
    }

}