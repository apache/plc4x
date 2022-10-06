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
package org.apache.plc4x.java.opcuaserver;

import java.io.*;
import java.net.InetAddress;
import java.security.*;
import java.security.cert.X509Certificate;
import java.util.Arrays;
import java.util.LinkedHashSet;
import java.util.Set;
import java.util.concurrent.CompletableFuture;

import java.nio.file.Path;
import java.nio.file.FileSystems;


import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;

import org.apache.plc4x.java.opcuaserver.context.CertificateKeyPair;
import org.apache.plc4x.java.opcuaserver.context.CertificateGenerator;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.eclipse.milo.opcua.sdk.server.OpcUaServer;
import org.eclipse.milo.opcua.sdk.server.api.config.OpcUaServerConfig;
import org.eclipse.milo.opcua.sdk.server.identity.CompositeValidator;
import org.eclipse.milo.opcua.sdk.server.identity.UsernameIdentityValidator;
import org.eclipse.milo.opcua.sdk.server.identity.X509IdentityValidator;
import org.eclipse.milo.opcua.stack.core.StatusCodes;
import org.eclipse.milo.opcua.stack.core.UaRuntimeException;
import org.eclipse.milo.opcua.stack.core.security.DefaultCertificateManager;
import org.eclipse.milo.opcua.stack.core.security.DefaultTrustListManager;
import org.eclipse.milo.opcua.stack.core.security.SecurityPolicy;
import org.eclipse.milo.opcua.stack.core.transport.TransportProfile;
import org.eclipse.milo.opcua.stack.core.types.builtin.DateTime;
import org.eclipse.milo.opcua.stack.core.types.builtin.LocalizedText;
import org.eclipse.milo.opcua.stack.core.types.enumerated.MessageSecurityMode;
import org.eclipse.milo.opcua.stack.core.types.structured.BuildInfo;
import org.eclipse.milo.opcua.stack.core.util.CertificateUtil;
import org.eclipse.milo.opcua.stack.server.EndpointConfiguration;
import org.eclipse.milo.opcua.stack.server.security.DefaultServerCertificateValidator;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import static org.eclipse.milo.opcua.sdk.server.api.config.OpcUaServerConfig.USER_TOKEN_POLICY_ANONYMOUS;
import static org.eclipse.milo.opcua.sdk.server.api.config.OpcUaServerConfig.USER_TOKEN_POLICY_USERNAME;
import static org.eclipse.milo.opcua.sdk.server.api.config.OpcUaServerConfig.USER_TOKEN_POLICY_X509;

import org.apache.commons.cli.*;

import org.apache.plc4x.java.opcuaserver.backend.Plc4xNamespace;
import org.apache.plc4x.java.opcuaserver.configuration.*;

public class OPCUAServer {

    private final Logger logger = LoggerFactory.getLogger(getClass());

    private Configuration config;
    private PasswordConfiguration passwordConfig;
    private CommandLine cmd = null;

    static {
        // Required for SecurityPolicy.Aes256_Sha256_RsaPss
        Security.addProvider(new BouncyCastleProvider());
    }

    protected String[] setPasswords() {
        Console cnsl = System.console();
        String[] ret = new String[3];

        System.out.println("Please enter password for certificate:- ");
        ret[0] = String.valueOf(cnsl.readPassword());

        System.out.println("Please enter a username for the OPC UA server admin account:- ");
        ret[1] = String.valueOf(cnsl.readLine());

        System.out.println("Please enter a password for the OPC UA server admin account:- ");
        ret[2] = String.valueOf(cnsl.readPassword());

        return ret;
    }

    private void setPasswordWrapper() {
        String[] ret;
        if (cmd.hasOption("test")) {
            ret = new String[] {"password", "admin", "password"};
        } else {
            ret = setPasswords();
        }
        try {
            passwordConfig.setSecurityPassword(ret[0]);
            passwordConfig.createUser(ret[1], ret[2], "admin-group");
        } catch (IOException e) {
            logger.error("Unable to save config file, please check folder permissions. " + e);
            System.exit(1);
        }
    }

    private void readPasswordConfig() {
        //Read Config File
        ObjectMapper mapper = new ObjectMapper(new YAMLFactory());
        mapper.findAndRegisterModules();
        try {
            Path path = FileSystems.getDefault().getPath(config.getDir()).resolve("security/.jibberish");
            File file = path.toFile();
            if (file.isFile() && !cmd.hasOption("set-passwords")) {
                passwordConfig = mapper.readValue(file, PasswordConfiguration.class);
                passwordConfig.setPasswordConfigFile(path);
            } else if (file.isFile() && cmd.hasOption("set-passwords")) {
                passwordConfig = mapper.readValue(file, PasswordConfiguration.class);
                passwordConfig.setPasswordConfigFile(path);
                setPasswordWrapper();
            } else {
                if (cmd.hasOption("interactive") || cmd.hasOption("set-passwords")) {
                    file.getParentFile().mkdirs();
                    passwordConfig = new PasswordConfiguration();
                    passwordConfig.setVersion("0.8");
                    passwordConfig.setPasswordConfigFile(path);
                    setPasswordWrapper();
                } else {
                    logger.info("Please re-run with the -i switch to setup the config file");
                    System.exit(1);
                }
            }
        } catch (IOException e) {
            logger.info("Error parsing password file " + e);
        }
    }

    private void readCommandLineArgs(String[] args) {
        Options options = new Options();

        Option input = new Option("c", "configfile", true, "configuration file");
        input.setRequired(true);
        options.addOption(input);

        Option setPassword = new Option("s", "set-passwords", false, "Reset passwords");
        setPassword.setRequired(false);
        options.addOption(setPassword);

        Option interactive = new Option("i", "interactive", false, "Interactively get asked to setup the config file from the console");
        interactive.setRequired(false);
        options.addOption(interactive);

        Option test = new Option("t", "test", false, "Used for testing the OPC UA Server");
        test.setRequired(false);
        options.addOption(test);

        CommandLineParser parser = new DefaultParser();
        HelpFormatter formatter = new HelpFormatter();
        cmd = null;

        try {
            cmd = parser.parse(options, args);
        } catch (ParseException e) {
            logger.info(e.getMessage());
            formatter.printHelp("Plc4x OPC UA Server", options);
            System.exit(1);
        }

        String configFile = cmd.getOptionValue("configfile");
        logger.info("Reading configuration file: {}", configFile);

        //Read Config File
        ObjectMapper mapper = new ObjectMapper(new YAMLFactory());
        mapper.findAndRegisterModules();
        try {
            config = mapper.readValue(new File(configFile), Configuration.class);
            config.setConfigFile(configFile);
            //Checking if the security directory has been configured.
            if (config.getDir() == null) {
                throw new IOException("Please set the dir in the config file");
            }

            readPasswordConfig();
        } catch (IOException e) {
            logger.info("Error parsing config file " + e);
            System.exit(1);
        }
    }

    public static void main(String[] args) throws Exception {
        OPCUAServer serverInit = new OPCUAServer(args);
        serverInit.getServer().startup().get();
        CompletableFuture<Void> future = new CompletableFuture<>();
        future.get();
    }

    private final OpcUaServer server;
    private final Plc4xNamespace plc4xNamespace;
    private final String certificateFileName = "plc4x-opcuaserver.pfx";

    public OPCUAServer(String[] args) throws Exception {

        readCommandLineArgs(args);

        File securityTempDir = new File(config.getDir(), "security");
        if (!securityTempDir.exists() && !securityTempDir.mkdirs()) {
            logger.error("Unable to create directory please confirm folder permissions on " + securityTempDir.toString());
            System.exit(1);
        }
        logger.info("Security Directory is: {}", securityTempDir.getAbsolutePath()); //

        File pkiDir = FileSystems.getDefault().getPath(config.getDir()).resolve("pki").toFile();
        DefaultTrustListManager trustListManager = new DefaultTrustListManager(pkiDir);
        logger.info("Certificate directory is: {}, Please move certificates from the reject dir to the trusted directory to allow encrypted access", pkiDir.getAbsolutePath());

        DefaultServerCertificateValidator certificateValidator =  new DefaultServerCertificateValidator(trustListManager);

        UsernameIdentityValidator identityValidator = new UsernameIdentityValidator(
            true,
            authChallenge -> {
                boolean check = passwordConfig.checkPassword(authChallenge.getUsername(), authChallenge.getPassword());
                if (!check) {
                    logger.info("Invalid password for user:- " + authChallenge.getUsername());
                }
                return check;
            }
        );

        KeyStore keyStore = KeyStore.getInstance("PKCS12");

        File serverKeyStore = securityTempDir.toPath().resolve(certificateFileName).toFile();

        X509IdentityValidator x509IdentityValidator = new X509IdentityValidator(c -> true);

        CertificateKeyPair certificate = null;
        if (!serverKeyStore.exists()) {
            if (!cmd.hasOption("interactive")) {
                logger.info("Please re-run with the -i switch to setup the security certificate key store");
                System.exit(1);
            }
            certificate = CertificateGenerator.generateCertificate();
            logger.info("Creating new KeyStore at {}", serverKeyStore);
            keyStore.load(null, passwordConfig.getSecurityPassword().toCharArray());
            keyStore.setKeyEntry("plc4x-certificate-alias", certificate.getKeyPair().getPrivate(), passwordConfig.getSecurityPassword().toCharArray(), new X509Certificate[] { certificate.getCertificate() });
            keyStore.store(new FileOutputStream(serverKeyStore), passwordConfig.getSecurityPassword().toCharArray());
        } else {
            logger.info("Loading KeyStore at {}", serverKeyStore);
            keyStore.load(new FileInputStream(serverKeyStore), passwordConfig.getSecurityPassword().toCharArray());
            String alias = keyStore.aliases().nextElement();
            KeyPair kp = new KeyPair(keyStore.getCertificate(alias).getPublicKey(),
                (PrivateKey) keyStore.getKey(alias, passwordConfig.getSecurityPassword().toCharArray()));
            certificate = new CertificateKeyPair(kp,(X509Certificate) keyStore.getCertificate(alias));
        }

        String applicationUri = CertificateUtil
            .getSanUri(certificate.getCertificate())
            .orElseThrow(() -> new UaRuntimeException(
                StatusCodes.Bad_ConfigurationError,
                "certificate is missing the application URI"));

        Set<EndpointConfiguration> endpointConfigurations = new LinkedHashSet<>();

        String hostname = InetAddress.getLocalHost().getHostName();

        EndpointConfiguration.Builder builder = EndpointConfiguration.newBuilder()
            .setBindAddress("0.0.0.0")
            .setHostname(hostname)
            .setPath("/plc4x")
            .setCertificate(certificate.getCertificate())
            .setBindPort(config.getTcpPort())
            .setSecurityMode(MessageSecurityMode.None)
            .addTokenPolicies(
                USER_TOKEN_POLICY_ANONYMOUS,
                USER_TOKEN_POLICY_USERNAME,
                USER_TOKEN_POLICY_X509);

        endpointConfigurations.add(
            builder.copy()
                .setSecurityPolicy(SecurityPolicy.Basic256Sha256)
                .setSecurityMode(MessageSecurityMode.SignAndEncrypt)
                .build()
        );

        endpointConfigurations.add(
            builder.copy()
                .setHostname("127.0.0.1")
                .setSecurityPolicy(SecurityPolicy.Basic256Sha256)
                .setSecurityMode(MessageSecurityMode.SignAndEncrypt)
                .build()
        );

        EndpointConfiguration.Builder discoveryBuilder = builder.copy()
            .setPath("/discovery")
            .setSecurityPolicy(SecurityPolicy.None)
            .setSecurityMode(MessageSecurityMode.None);

        endpointConfigurations.add(discoveryBuilder.build());

        EndpointConfiguration.Builder discoveryLocalBuilder = builder.copy()
            .setPath("/discovery")
            .setHostname("127.0.0.1")
            .setSecurityPolicy(SecurityPolicy.None)
            .setSecurityMode(MessageSecurityMode.None);

        endpointConfigurations.add(discoveryLocalBuilder.build());

        EndpointConfiguration.Builder discoveryLocalPlc4xBuilder = builder.copy()
            .setPath("/plc4x/discovery")
            .setHostname("127.0.0.1")
            .setSecurityPolicy(SecurityPolicy.None)
            .setSecurityMode(MessageSecurityMode.None);

        endpointConfigurations.add(discoveryLocalPlc4xBuilder.build());

        if (!config.getDisableInsecureEndpoint()) {
            EndpointConfiguration.Builder noSecurityBuilder = builder.copy()
                .setSecurityPolicy(SecurityPolicy.None)
                .setTransportProfile(TransportProfile.TCP_UASC_UABINARY);
            endpointConfigurations.add(noSecurityBuilder.build());
        }

        //Always add an unsecured endpoint to localhost, this is a work around for Milo throwing an exception if it isn't here.
        EndpointConfiguration.Builder noSecurityBuilder = builder.copy()
            .setSecurityPolicy(SecurityPolicy.None)
            .setHostname("127.0.0.1")
            .setTransportProfile(TransportProfile.TCP_UASC_UABINARY)
            .setSecurityMode(MessageSecurityMode.None);
        endpointConfigurations.add(noSecurityBuilder.build());

        DefaultCertificateManager certificateManager = new DefaultCertificateManager(
            certificate.getKeyPair(),
            Arrays.stream(keyStore.getCertificateChain(keyStore.getCertificateAlias(certificate.getCertificate())))// Added so that existing certificates are loaded on startup
                .map(X509Certificate.class::cast)
                .toArray(X509Certificate[]::new)
        );

        OpcUaServerConfig serverConfig = OpcUaServerConfig.builder()
            .setApplicationUri(applicationUri)
            .setApplicationName(LocalizedText.english(applicationUri))
            .setEndpoints(endpointConfigurations)
            .setBuildInfo(
                new BuildInfo(
                    "urn:eclipse:milo:plc4x:server",
                    "org.apache.plc4x",
                    config.getName(),
                    OpcUaServer.SDK_VERSION,
                    "", DateTime.now()))
            .setCertificateManager(certificateManager)
            .setTrustListManager(trustListManager)
            .setCertificateValidator(certificateValidator)
            .setIdentityValidator(new CompositeValidator(identityValidator, x509IdentityValidator))
            .setProductUri("urn:eclipse:milo:plc4x:server")
            .build();

        server = new OpcUaServer(serverConfig);

        plc4xNamespace = new Plc4xNamespace(server, config);
        plc4xNamespace.startup();
    }

    public OpcUaServer getServer() {
        return server;
    }

    public CompletableFuture<OpcUaServer> startup() {
        return server.startup();
    }

    public CompletableFuture<OpcUaServer> shutdown() {
        plc4xNamespace.shutdown();
        return server.shutdown();
    }

}
