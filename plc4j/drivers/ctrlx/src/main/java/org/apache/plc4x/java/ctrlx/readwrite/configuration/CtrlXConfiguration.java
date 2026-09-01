/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.ctrlx.readwrite.configuration;

import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Secret;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.StringDefaultValue;

public class CtrlXConfiguration implements Configuration {

    /**
     * A key store holding the certificates this connection should trust, in place of the
     * certificate authorities the JVM ships with.
     *
     * <p>Named the same as the equivalent in the OPC UA driver and the TLS transport, so that
     * pinning a device's own certificate reads the same wherever it is done.</p>
     */
    @ConfigurationParameter("tls.trust-store")
    @Description("Key store of certificates to trust, instead of the JVM's public authorities")
    public String trustStoreFile;

    @Secret
    @ConfigurationParameter("tls.trust-store-password")
    @Description("Password of the trust store named by tls.trust-store")
    public String trustStorePassword;

    @ConfigurationParameter("tls.trust-store-type")
    @StringDefaultValue("PKCS12")
    @Description("Type of the trust store named by tls.trust-store")
    public String trustStoreType;

    /**
     * A single PEM certificate to trust, for the common case of one device carrying one
     * certificate and no key store to put it in.
     */
    @ConfigurationParameter("server-certificate-file")
    @Description("PEM certificate of the device to trust")
    public String serverCertificateFile;

    /**
     * Whether to go back to trusting the certificate this driver ships and to accepting it for any
     * host.
     *
     * <p>That certificate is in the jar, so it identifies nobody: anything holding it and its key
     * is trusted, and because it was the only anchor, a device with a properly issued certificate
     * could not be reached at all. It remains available for a device still on its factory
     * settings, and says what it is doing when used.</p>
     */
    @ConfigurationParameter("allow-factory-default-certificate")
    @BooleanDefaultValue(false)
    @Description("Trust the factory default certificate shipped with this driver, for any host")
    public boolean allowFactoryDefaultCertificate;

    /**
     * Whether to accept a certificate issued for some other name than the host connected to.
     */
    @ConfigurationParameter("ignore-common-name")
    @BooleanDefaultValue(false)
    @Description("Accept a device certificate issued for a different host than the one connected to")
    public boolean ignoreCommonName;

    /**
     * How many nodes a single browse will read.
     *
     * <p>A browse asks the device for the children of a node and then asks the same of each child,
     * so how much work it is depends entirely on the tree the device describes. Nothing in the
     * answer says how large that tree is, or that it ends.</p>
     */
    @ConfigurationParameter("browse-max-total-nodes")
    @IntDefaultValue(1000000)
    @Description("Largest number of nodes a single browse will read; 0 for no limit")
    public int browseMaxTotalNodes;

    /**
     * How deep a browse will follow children.
     *
     * <p>Each child is addressed by its parent's path with the child's name appended, so a device
     * naming a child that leads back to a node already passed produces a longer path every time
     * and never repeats one. Depth is what ends that, not a record of where the browse has
     * been.</p>
     */
    @ConfigurationParameter("browse-max-depth")
    @IntDefaultValue(64)
    @Description("How deep a browse will follow children; 0 for no limit")
    public int browseMaxDepth;

    public int getBrowseMaxTotalNodes() {
        return browseMaxTotalNodes;
    }

    public int getBrowseMaxDepth() {
        return browseMaxDepth;
    }

    public String getTrustStoreFile() {
        return trustStoreFile;
    }

    public String getTrustStorePassword() {
        return trustStorePassword;
    }

    public String getTrustStoreType() {
        return trustStoreType;
    }

    public String getServerCertificateFile() {
        return serverCertificateFile;
    }

    public boolean isAllowFactoryDefaultCertificate() {
        return allowFactoryDefaultCertificate;
    }

    public boolean isIgnoreCommonName() {
        return ignoreCommonName;
    }
}
