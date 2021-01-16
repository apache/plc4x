/*
 * Copyright (c) 2019 the Eclipse Milo Authors
 *
 * This program and the accompanying materials are made
 * available under the terms of the Eclipse Public License 2.0
 * which is available at https://www.eclipse.org/legal/epl-2.0/
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package org.apache.plc4x.java.opcuaserver;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.security.Key;
import java.security.KeyPair;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.cert.X509Certificate;
import java.util.Arrays;
import java.util.Set;
import java.util.UUID;
import java.util.regex.Pattern;

import com.google.common.collect.Sets;
import org.apache.plc4x.java.opcuaserver.configuration.Configuration;
import org.apache.plc4x.java.opcuaserver.configuration.PasswordConfiguration;
import org.eclipse.milo.opcua.sdk.server.util.HostnameUtil;
import org.eclipse.milo.opcua.stack.core.util.SelfSignedCertificateBuilder;
import org.eclipse.milo.opcua.stack.core.util.SelfSignedCertificateGenerator;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

class KeyStoreLoader {

    private static final Pattern IP_ADDR_PATTERN = Pattern.compile(
        "^(([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\.){3}([01]?\\d\\d?|2[0-4]\\d|25[0-5])$");

    private final Logger logger = LoggerFactory.getLogger(getClass());

    private final String certificateFileName = "plc4x-opcuaserver.pfx";

    private X509Certificate[] serverCertificateChain;
    private X509Certificate serverCertificate;
    private KeyPair serverKeyPair;

    private Configuration config;
    private PasswordConfiguration passwordConfig;

    public KeyStoreLoader(Configuration config, PasswordConfiguration passwordConfig, Boolean interactive) {
        this.config = config;
        this.passwordConfig = passwordConfig;

        File securityTempDir = new File(config.getDir(), "security");
        if (!securityTempDir.exists() && !securityTempDir.mkdirs()) {
            logger.error("Unable to create directory please confirm folder permissions on " + securityTempDir.toString());
            System.exit(1);
        }
        logger.info("security dir: {}", securityTempDir.getAbsolutePath());

        try {
            load(securityTempDir, interactive);
        } catch (Exception e) {
            logger.error("Error loading the key store " + e);
            System.exit(1);
        }
    }

    public KeyStoreLoader load(File baseDir, boolean interactive) throws Exception {
        KeyStore keyStore = KeyStore.getInstance("PKCS12");

        File serverKeyStore = baseDir.toPath().resolve(certificateFileName).toFile();

        if (!serverKeyStore.exists()) {
            if (!interactive) {
                logger.info("Please re-run with the -i switch to setup the security certificate key store");
                System.exit(1);
            }

            logger.info("Creating keystore at {}", serverKeyStore);
            keyStore.load(null, passwordConfig.getSecurityPassword().toCharArray());

            logger.info("Creating self signed certiciate {}", serverKeyStore);
            KeyPair keyPair = SelfSignedCertificateGenerator.generateRsaKeyPair(2048);

            String applicationUri = "urn:eclipse:milo:plc4x:server" + UUID.randomUUID();

            SelfSignedCertificateBuilder builder = new SelfSignedCertificateBuilder(keyPair)
                .setCommonName(applicationUri)
                .setOrganization("org.apache")
                .setOrganizationalUnit("plc4x")
                .setLocalityName("Wakefield")
                .setStateName("MA")
                .setCountryCode("US")
                .setApplicationUri(applicationUri);

            // Get as many hostnames and IP addresses as we can listed in the certificate.
            Set<String> hostnames = Sets.union(
                Sets.newHashSet(HostnameUtil.getHostname()),
                HostnameUtil.getHostnames("0.0.0.0", false)
            );

            logger.info("using IP address/hostnames {}", hostnames.toString());

            for (String hostname : hostnames) {
                if (IP_ADDR_PATTERN.matcher(hostname).matches()) {
                    builder.addIpAddress(hostname);
                } else {
                    builder.addDnsName(hostname);
                }
            }

            X509Certificate certificate = builder.build();

            keyStore.setKeyEntry(config.getName(), keyPair.getPrivate(), passwordConfig.getSecurityPassword().toCharArray(), new X509Certificate[]{certificate});
            keyStore.store(new FileOutputStream(serverKeyStore), passwordConfig.getSecurityPassword().toCharArray());

            logger.info("Self signed certificate created. Replace {} and update config file passwords if not using a signed certificate.", serverKeyStore);

        } else {
            logger.info("Loading KeyStore at {}", serverKeyStore);
            keyStore.load(new FileInputStream(serverKeyStore), passwordConfig.getSecurityPassword().toCharArray());
        }

        Key serverPrivateKey = keyStore.getKey(config.getName(), passwordConfig.getSecurityPassword().toCharArray());
        if (serverPrivateKey instanceof PrivateKey) {
            serverCertificate = (X509Certificate) keyStore.getCertificate(config.getName());

            serverCertificateChain = Arrays.stream(keyStore.getCertificateChain(config.getName()))
                .map(X509Certificate.class::cast)
                .toArray(X509Certificate[]::new);

            PublicKey serverPublicKey = serverCertificate.getPublicKey();
            serverKeyPair = new KeyPair(serverPublicKey, (PrivateKey) serverPrivateKey);
        }

        return this;
    }

    X509Certificate getServerCertificate() {
        return serverCertificate;
    }

    public X509Certificate[] getServerCertificateChain() {
        return serverCertificateChain;
    }

    KeyPair getServerKeyPair() {
        return serverKeyPair;
    }

}
