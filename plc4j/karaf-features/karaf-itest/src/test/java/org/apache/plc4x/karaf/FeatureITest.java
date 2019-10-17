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
import static org.ops4j.pax.exam.CoreOptions.mavenBundle;
import static org.ops4j.pax.exam.karaf.options.KarafDistributionOption.features;
import static org.ops4j.pax.exam.karaf.options.KarafDistributionOption.logLevel;

@RunWith(PaxExam.class)
@ExamReactorStrategy(PerMethod.class)
public class FeatureITest extends KarafTestSupport {

    @Configuration
    public Option[] config() {
        Option[] options = new Option[]{
            logLevel(LogLevelOption.LogLevel.DEBUG),
            KarafDistributionOption.editConfigurationFilePut("etc/system.properties", "plc4x.version", System.getProperty("plc4x.version")),
//            features("scr"),
//            mavenBundle()
//                .groupId("org.ops4j.pax.exam.samples")
//                .artifactId("pax-exam-sample8-ds")
//                .version("4.13.1")
//                .start()
        };
        return Stream.of(super.config(), options).flatMap(Stream::of).toArray(Option[]::new);

//
//        String httpPort = Integer.toString(getAvailablePort(Integer.parseInt(MIN_HTTP_PORT), Integer.parseInt(MAX_HTTP_PORT)));
//        String rmiRegistryPort = Integer.toString(getAvailablePort(Integer.parseInt(MIN_RMI_REG_PORT), Integer.parseInt(MAX_RMI_REG_PORT)));
//        String rmiServerPort = Integer.toString(getAvailablePort(Integer.parseInt(MIN_RMI_SERVER_PORT), Integer.parseInt(MAX_RMI_SERVER_PORT)));
//        String sshPort = Integer.toString(getAvailablePort(Integer.parseInt(MIN_SSH_PORT), Integer.parseInt(MAX_SSH_PORT)));
//        String localRepository = System.getProperty("org.ops4j.pax.url.mvn.localRepository");
//        if (localRepository == null) {
//            localRepository = "";
//        }
//
//        if (JavaVersionUtil.getMajorVersion() >= 9) {
//            return new Option[]{
//                //KarafDistributionOption.debugConfiguration("8889", true),
//                karafDistributionConfiguration().frameworkUrl(karafUrl).name("Apache Karaf").unpackDirectory(new File("target/exam")),
//                // enable JMX RBAC security, thanks to the KarafMBeanServerBuilder
//                configureSecurity().disableKarafMBeanServerBuilder(),
//                // configureConsole().ignoreLocalConsole(),
//                keepRuntimeFolder(),
//                logLevel(LogLevelOption.LogLevel.INFO),
//                mavenBundle().groupId("org.awaitility").artifactId("awaitility").versionAsInProject(),
//                mavenBundle().groupId("org.apache.servicemix.bundles").artifactId("org.apache.servicemix.bundles.hamcrest").versionAsInProject(),
//                mavenBundle().groupId("org.apache.karaf.itests").artifactId("common").versionAsInProject(),
//                editConfigurationFilePut("etc/org.ops4j.pax.web.cfg", "org.osgi.service.http.port", httpPort),
//                editConfigurationFilePut("etc/org.apache.karaf.management.cfg", "rmiRegistryPort", rmiRegistryPort),
//                editConfigurationFilePut("etc/org.apache.karaf.management.cfg", "rmiServerPort", rmiServerPort),
//                editConfigurationFilePut("etc/org.apache.karaf.shell.cfg", "sshPort", sshPort),
//                editConfigurationFilePut("etc/org.ops4j.pax.url.mvn.cfg", "org.ops4j.pax.url.mvn.localRepository", localRepository),
//                new VMOption("--add-reads=java.xml=java.logging"),
//                new VMOption("--add-exports=java.base/org.apache.karaf.specs.locator=java.xml,ALL-UNNAMED"),
//                new VMOption("--patch-module"),
//                new VMOption("java.base=lib/endorsed/org.apache.karaf.specs.locator-"
//                    + System.getProperty("karaf.version") + ".jar"),
//                new VMOption("--patch-module"),
//                new VMOption("java.xml=lib/endorsed/org.apache.karaf.specs.java.xml-"
//                    + System.getProperty("karaf.version") + ".jar"),
//                new VMOption("--add-opens"),
//                new VMOption("java.base/java.security=ALL-UNNAMED"),
//                new VMOption("--add-opens"),
//                new VMOption("java.base/java.net=ALL-UNNAMED"),
//                new VMOption("--add-opens"),
//                new VMOption("java.base/java.lang=ALL-UNNAMED"),
//                new VMOption("--add-opens"),
//                new VMOption("java.base/java.util=ALL-UNNAMED"),
//                new VMOption("--add-opens"),
//                new VMOption("java.naming/javax.naming.spi=ALL-UNNAMED"),
//                new VMOption("--add-opens"),
//                new VMOption("java.rmi/sun.rmi.transport.tcp=ALL-UNNAMED"),
//                new VMOption("--add-exports=java.base/sun.net.www.protocol.http=ALL-UNNAMED"),
//                new VMOption("--add-exports=java.base/sun.net.www.protocol.https=ALL-UNNAMED"),
//                new VMOption("--add-exports=java.base/sun.net.www.protocol.jar=ALL-UNNAMED"),
//                new VMOption("--add-exports=jdk.naming.rmi/com.sun.jndi.url.rmi=ALL-UNNAMED"),
//                new VMOption("-classpath"),
//                new VMOption("lib/jdk9plus/*" + File.pathSeparator + "lib/boot/*")
//            };
//        } else {
//            return new Option[]{
//                //KarafDistributionOption.debugConfiguration("8889", true),
//                karafDistributionConfiguration().frameworkUrl(karafUrl).name("Apache Karaf").unpackDirectory(new File("target/exam")),
//                // enable JMX RBAC security, thanks to the KarafMBeanServerBuilder
//                configureSecurity().disableKarafMBeanServerBuilder(),
//                // configureConsole().ignoreLocalConsole(),
//                keepRuntimeFolder(),
//                logLevel(LogLevelOption.LogLevel.INFO),
//                mavenBundle().groupId("org.awaitility").artifactId("awaitility").versionAsInProject(),
//                mavenBundle().groupId("org.apache.servicemix.bundles").artifactId("org.apache.servicemix.bundles.hamcrest").versionAsInProject(),
//                mavenBundle().groupId("org.apache.karaf.itests").artifactId("common").versionAsInProject(),
//                editConfigurationFilePut("etc/org.ops4j.pax.web.cfg", "org.osgi.service.http.port", httpPort),
//                editConfigurationFilePut("etc/org.apache.karaf.management.cfg", "rmiRegistryPort", rmiRegistryPort),
//                editConfigurationFilePut("etc/org.apache.karaf.management.cfg", "rmiServerPort", rmiServerPort),
//                editConfigurationFilePut("etc/org.apache.karaf.shell.cfg", "sshPort", sshPort),
//                editConfigurationFilePut("etc/org.ops4j.pax.url.mvn.cfg", "org.ops4j.pax.url.mvn.localRepository", localRepository)
//            };
//        }
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