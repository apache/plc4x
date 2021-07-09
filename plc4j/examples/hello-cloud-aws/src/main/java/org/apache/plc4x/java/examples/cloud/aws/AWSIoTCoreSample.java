package org.apache.plc4x.java.examples.cloud.aws;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import software.amazon.awssdk.crt.CRT;
import software.amazon.awssdk.crt.CrtResource;
import software.amazon.awssdk.crt.CrtRuntimeException;
import software.amazon.awssdk.crt.io.*;
import software.amazon.awssdk.crt.mqtt.*;
import software.amazon.awssdk.iot.iotjobs.model.RejectedError;

import java.nio.charset.StandardCharsets;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutionException;

public class AWSIoTCoreSample {

    private static final Logger logger = LoggerFactory.getLogger(AWSIoTCoreSample.class);

    public static void main(String[] args) throws Exception {
        CliOptions options = CliOptions.fromArgs(args);
        if (options == null) {
            CliOptions.printHelp();
            System.exit(1);
        }

        //AWS Connection Callback method implements
        MqttClientConnectionEvents callbacks = new MqttClientConnectionEvents() {
            @Override
            public void onConnectionInterrupted(int errorCode) {
                if (errorCode != 0) {
                    logger.warn("Connection interrupted: " + errorCode + ": " + CRT.awsErrorString(errorCode));
                }
            }

            @Override
            public void onConnectionResumed(boolean sessionPresent) {
                logger.info(("Connection resumed: " + (sessionPresent ? "existing session" : "clean session")));
            }
        };

        String userName = options.getUserName();

        // Can Custom auth with Get method.
        // https://docs.aws.amazon.com/ko_kr/iot/latest/developerguide/custom-auth.html
        if (options.getAuthParams() != null && options.getAuthParams().length > 0) {
            if (userName.length() > 0) {
                StringBuilder usernameBuilder = new StringBuilder();
                usernameBuilder.append(userName);
                usernameBuilder.append("?");
                for (int i = 0; i < options.getAuthParams().length; ++i) {
                    usernameBuilder.append(options.getAuthParams().length);
                    if (i + 1 < options.getAuthParams().length) {
                        usernameBuilder.append("&");
                    }
                }
                userName = usernameBuilder.toString();
            }
        }

        try (EventLoopGroup eventLoopGroup = new EventLoopGroup(1);
             HostResolver resolver = new HostResolver(eventLoopGroup);
             ClientBootstrap clientBootstrap = new ClientBootstrap(eventLoopGroup, resolver);
             TlsContextOptions tlsContextOptions = TlsContextOptions.createWithMtlsFromPath(options.getCert(), options.getKey())) {
            tlsContextOptions.overrideDefaultTrustStoreFromPath(null, options.getRootca());

            int port = 8883;

            try (TlsContext tlsContext = new TlsContext(tlsContextOptions);
                 MqttClient client = new MqttClient(clientBootstrap, tlsContext);
                 MqttConnectionConfig config = new MqttConnectionConfig()) {

                config.setMqttClient(client);
                config.setClientId(options.getClientId());
                config.setConnectionCallbacks(callbacks);
                config.setCleanSession(true);
                config.setEndpoint(options.getEndpoint());
                config.setPort(port);

                if (userName != null && userName.length() > 0) {
                    config.setLogin(userName, options.getPassword());
                }

                try (MqttClientConnection connection = new MqttClientConnection(config)) {

                    CompletableFuture<Boolean> connected = connection.connect();
                    try {
                        boolean sessionPresent = connected.get();
                        logger.info("Connected to " + (!sessionPresent ? "new" : "existing") + " session!");
                    } catch (Exception ex) {
                        throw new RuntimeException("Exception occurred during connect", ex);
                    }

                    // Subscribe
                    CountDownLatch countDownLatch = new CountDownLatch(options.getMessagesToPublish());

                    CompletableFuture<Integer> subscribed = connection.subscribe(options.getTopic(), QualityOfService.AT_LEAST_ONCE, (message) -> {
                        String payload = new String(message.getPayload(), StandardCharsets.UTF_8);
                        logger.info("MESSAGE: " + payload);
                        countDownLatch.countDown();
                    });
                    subscribed.get();

                    try (PlcConnection plcConnection = new PlcDriverManager().getConnection(options.getConnectionString())) {
                        PlcReadRequest.Builder builder = plcConnection.readRequestBuilder();
                        for (int i = 0; i < options.getFieldAddress().length; i++) {
                            builder.addItem("value-" + options.getFieldAddress()[i], options.getFieldAddress()[i]);
                        }
                        PlcReadRequest readRequest = builder.build();

                        int count = 0;
                        while (count++ < options.getMessagesToPublish()) {
                            PlcReadResponse plcReadResponse = readRequest.execute().get();
                            for (String fieldName : plcReadResponse.getFieldNames()) {
                                if (plcReadResponse.getResponseCode(fieldName) == PlcResponseCode.OK) {
                                    int numValues = plcReadResponse.getNumberOfValues(fieldName);
                                    for (int i = 0; i < numValues; i++) {
                                        CompletableFuture<Integer> published = connection.publish(new MqttMessage(options.getTopic(), plcReadResponse.getString(fieldName, i).getBytes(), QualityOfService.AT_LEAST_ONCE, false));
                                        published.get();
                                        Thread.sleep(500);
                                    }
                                } else {
                                    logger.error("Error[{}]: {}", fieldName, plcReadResponse.getResponseCode(fieldName).name());
                                }
                            }
                        }
                    }
                    countDownLatch.await();

                    CompletableFuture<Void> disconnected = connection.disconnect();
                    disconnected.get();
                }
            } catch (CrtRuntimeException | InterruptedException | ExecutionException ex) {
                logger.error("Exception encountered: " + ex.toString());
            }
        }
        CrtResource.waitForNoResources();
        logger.info("Complete!");
        System.exit(0);
    }

    static void onRejectedError(RejectedError error) {
        logger.error("Request rejected: " + error.code.toString() + ": " + error.message);
    }
}

