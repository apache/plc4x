package org.apache.plc4x.java.opcua.connection;

import com.google.common.collect.Sets;
import org.eclipse.milo.opcua.stack.core.util.SelfSignedCertificateBuilder;
import org.eclipse.milo.opcua.stack.core.util.SelfSignedCertificateGenerator;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.InputStream;
import java.io.OutputStream;
import java.net.Inet4Address;
import java.net.InetAddress;
import java.net.NetworkInterface;
import java.net.SocketException;
import java.net.UnknownHostException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.Key;
import java.security.KeyPair;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.cert.X509Certificate;
import java.util.Collections;
import java.util.Enumeration;
import java.util.HashSet;
import java.util.Iterator;
import java.util.Set;
import java.util.regex.Pattern;

/**
 * Copied from Eclipse Milo Project:
 * https://github.com/eclipse/milo/blob/master/milo-examples/server-examples/src/main/java/org/eclipse/milo/examples/server/KeyStoreLoader.java
 */
public class KeyStoreLoader {

    private static final Pattern IP_ADDR_PATTERN = Pattern.compile(
        "^(([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\.){3}([01]?\\d\\d?|2[0-4]\\d|25[0-5])$");

    private static final String CLIENT_ALIAS = "client-cert";
    private static final char[] PASSWORD = "plc4x".toCharArray();

    private final Logger logger = LoggerFactory.getLogger(getClass());

    private X509Certificate clientCertificate;
    private KeyPair clientKeyPair;

//    KeyStoreLoader load(Path baseDir) throws Exception {
//        return load(baseDir, true);
//    }
//
//    KeyStoreLoader load(Path baseDir, boolean resolve) throws Exception {
//        return load(baseDir, null, resolve);
//    }
//
//    KeyStoreLoader load(Path file, String password) throws KeyStoreException, IOException, CertificateException, NoSuchAlgorithmException, UnrecoverableKeyException {
//        KeyStore keyStore = KeyStore.getInstance("PKCS12");
//
//        Path serverKeyStore = file;
//
//        logger.info("Loading KeyStore at {}", serverKeyStore);
//
//        try (InputStream in = Files.newInputStream(serverKeyStore)) {
//            keyStore.load(in, password.toCharArray());
//        }
//
//        Key serverPrivateKey = keyStore.getKey(CLIENT_ALIAS, password.toCharArray());
//        if (serverPrivateKey instanceof PrivateKey) {
//            clientCertificate = (X509Certificate) keyStore.getCertificate(CLIENT_ALIAS);
//            PublicKey serverPublicKey = clientCertificate.getPublicKey();
//            clientKeyPair = new KeyPair(serverPublicKey, (PrivateKey) serverPrivateKey);
//        }
//
//        return this;
//    }

    KeyStoreLoader load(Path baseDir) throws Exception {
        KeyStore keyStore = KeyStore.getInstance("PKCS12");

        Path serverKeyStore = baseDir;

        logger.info("Loading KeyStore at {}", serverKeyStore);

        if (!Files.exists(serverKeyStore)) {
            keyStore.load(null, PASSWORD);

            KeyPair keyPair = SelfSignedCertificateGenerator.generateRsaKeyPair(2048);

            SelfSignedCertificateBuilder builder = new SelfSignedCertificateBuilder(keyPair)
                .setCommonName("Eclipse Milo Example Client")
                .setOrganization("digitalpetri")
                .setOrganizationalUnit("dev")
                .setLocalityName("Folsom")
                .setStateName("CA")
                .setCountryCode("US")
//                .setApplicationUri("urn:eclipse:milo:examples:client")
                .setApplicationUri("urn:plc4x-client")
                .addDnsName("localhost")
                .addIpAddress("127.0.0.1");

            // Get as many hostnames and IP addresses as we can listed in the certificate.
            for (String hostname : getHostnames("0.0.0.0", true)) {
                if (IP_ADDR_PATTERN.matcher(hostname).matches()) {
                    builder.addIpAddress(hostname);
                } else {
                    builder.addDnsName(hostname);
                }
            }

            X509Certificate certificate = builder.build();

            keyStore.setKeyEntry(CLIENT_ALIAS, keyPair.getPrivate(), PASSWORD, new X509Certificate[]{certificate});
            try (OutputStream out = Files.newOutputStream(serverKeyStore)) {
                keyStore.store(out, PASSWORD);
            }
        } else {
            try (InputStream in = Files.newInputStream(serverKeyStore)) {
                keyStore.load(in, PASSWORD);
            }
        }

        Key serverPrivateKey = keyStore.getKey(CLIENT_ALIAS, PASSWORD);
        if (serverPrivateKey instanceof PrivateKey) {
            clientCertificate = (X509Certificate) keyStore.getCertificate(CLIENT_ALIAS);
            PublicKey serverPublicKey = clientCertificate.getPublicKey();
            clientKeyPair = new KeyPair(serverPublicKey, (PrivateKey) serverPrivateKey);
        }

        return this;
    }

    X509Certificate getClientCertificate() {
        return clientCertificate;
    }

    KeyPair getClientKeyPair() {
        return clientKeyPair;
    }

    public Set<String> getHostnames(String address, boolean includeLoopback) {
        HashSet hostnames = Sets.newHashSet();

        try {
            InetAddress inetAddress = InetAddress.getByName(address);
            if (inetAddress.isAnyLocalAddress()) {
                try {
                    Enumeration<NetworkInterface> nis = NetworkInterface.getNetworkInterfaces();
                    Iterator var5 = Collections.list(nis).iterator();

                    while(var5.hasNext()) {
                        NetworkInterface ni = (NetworkInterface)var5.next();
                        Collections.list(ni.getInetAddresses()).forEach((ia) -> {
                            if (ia instanceof Inet4Address) {
                                boolean loopback = ia.isLoopbackAddress();
                                if (!loopback || includeLoopback) {
                                    hostnames.add(ia.getHostName());
                                    hostnames.add(ia.getHostAddress());
                                    hostnames.add(ia.getCanonicalHostName());
                                }
                            }

                        });
                    }
                } catch (SocketException var7) {
                    logger.warn("Failed to NetworkInterfaces for bind address: {}", address, var7);
                }
            } else {
                boolean loopback = inetAddress.isLoopbackAddress();
                if (!loopback || includeLoopback) {
                    hostnames.add(inetAddress.getHostName());
                    hostnames.add(inetAddress.getHostAddress());
                    hostnames.add(inetAddress.getCanonicalHostName());
                }
            }
        } catch (UnknownHostException var8) {
            logger.warn("Failed to get InetAddress for bind address: {}", address, var8);
        }

        return hostnames;
    }

}
