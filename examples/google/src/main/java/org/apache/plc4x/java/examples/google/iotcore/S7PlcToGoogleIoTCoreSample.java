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
package org.apache.plc4x.java.examples.google.iotcore;

import io.jsonwebtoken.JwtBuilder;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.SignatureAlgorithm;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.eclipse.paho.client.mqttv3.*;
import org.eclipse.paho.client.mqttv3.persist.MemoryPersistence;
import org.joda.time.DateTime;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.ByteBuffer;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.KeyFactory;
import java.security.spec.PKCS8EncodedKeySpec;

// [START iot_mqtt_includes]
// [END iot_mqtt_includes]

public class S7PlcToGoogleIoTCoreSample {

    private static final Logger logger = LoggerFactory.getLogger(S7PlcToGoogleIoTCoreSample.class);

    // [START iot_mqtt_jwt]

    /**
     * Create a Cloud IoT Core JWT for the given project id, signed with the given RSA key.
     */
    private static String createJwtRsa(String projectId, String privateKeyFile) throws Exception {
        DateTime now = new DateTime();
        // Create a JWT to authenticate this device. The device will be disconnected after the token
        // expires, and will have to reconnect with a new token. The audience field should always be set
        // to the GCP project id.
        JwtBuilder jwtBuilder =
            Jwts.builder()
                .setIssuedAt(now.toDate())
                .setExpiration(now.plusMinutes(20).toDate())
                .setAudience(projectId);

        byte[] keyBytes = Files.readAllBytes(Paths.get(privateKeyFile));
        PKCS8EncodedKeySpec spec = new PKCS8EncodedKeySpec(keyBytes);
        KeyFactory kf = KeyFactory.getInstance("RSA");

        return jwtBuilder.signWith(SignatureAlgorithm.RS256, kf.generatePrivate(spec)).compact();
    }

    /**
     * Create a Cloud IoT Core JWT for the given project id, signed with the given ES key.
     */
    private static String createJwtEs(String projectId, String privateKeyFile) throws Exception {
        DateTime now = new DateTime();
        // Create a JWT to authenticate this device. The device will be disconnected after the token
        // expires, and will have to reconnect with a new token. The audience field should always be set
        // to the GCP project id.
        JwtBuilder jwtBuilder =
            Jwts.builder()
                .setIssuedAt(now.toDate())
                .setExpiration(now.plusMinutes(20).toDate())
                .setAudience(projectId);

        byte[] keyBytes = Files.readAllBytes(Paths.get(privateKeyFile));
        PKCS8EncodedKeySpec spec = new PKCS8EncodedKeySpec(keyBytes);
        KeyFactory kf = KeyFactory.getInstance("EC");

        return jwtBuilder.signWith(SignatureAlgorithm.ES256, kf.generatePrivate(spec)).compact();
    }
    // [END iot_mqtt_jwt]

    /**
     * Attaches the callback used when configuration changes occur.
     */
    private static void attachCallback(MqttClient client, String deviceId) throws MqttException {
        // [START iot_mqtt_configcallback]
        MqttCallback mCallback = new MqttCallback() {
            @Override
            public void connectionLost(Throwable cause) {
                // Do nothing...
            }

            @Override
            public void messageArrived(String topic, MqttMessage message) {
                String payload = new String(message.getPayload());
                System.out.println("Payload : " + payload);
                // TODO: Insert your parsing / handling of the configuration message here.
            }

            @Override
            public void deliveryComplete(IMqttDeliveryToken token) {
                // Do nothing;
            }
        };

        String configTopic = String.format("/devices/%s/config", deviceId);
        client.subscribe(configTopic, 1);

        client.setCallback(mCallback);
    }
    // [END iot_mqtt_configcallback]


    private static void setConnectPassword(MqttExampleOptions options, MqttConnectOptions connectOptions) throws Exception {
        switch (options.algorithm) {
            case "RS256":
                connectOptions.setPassword(
                    createJwtRsa(options.projectId, options.privateKeyFile).toCharArray());
                break;
            case "ES256":
                connectOptions.setPassword(
                    createJwtEs(options.projectId, options.privateKeyFile).toCharArray());
                break;
            default:
                throw new IllegalArgumentException(
                    "Invalid algorithm " + options.algorithm
                        + ". Should be one of 'RS256' or 'ES256'.");
        }
    }

    /**
     * Example code do demonstrate sending events from an S7 device to Microsoft Azure IoT Hub
     *
     * @param args Expected: [plc4x connection string, plc4x address, IoT-Hub connection string].
     */
    public static void main(String[] args) throws Exception {

        // [START iot_mqtt_configuremqtt]
        MqttExampleOptions options = MqttExampleOptions.fromFlags(args);
        if (options == null) {
            // Could not parse.
            System.exit(1);
        }

        // Build the connection string for Google's Cloud IoT Core MQTT server. Only SSL
        // connections are accepted. For server authentication, the JVM's root certificates
        // are used.
        final String mqttServerAddress =
            String.format("ssl://%s:%s", options.mqttBridgeHostname, options.mqttBridgePort);

        // Create our MQTT client. The mqttClientId is a unique string that identifies this device. For
        // Google Cloud IoT Core, it must be in the format below.
        final String mqttClientId =
            String.format("projects/%s/locations/%s/registries/%s/devices/%s",
                options.projectId, options.cloudRegion, options.registryId, options.deviceId);

        MqttConnectOptions connectOptions = new MqttConnectOptions();
        // Note that the Google Cloud IoT Core only supports MQTT 3.1.1, and Paho requires that we
        // explictly set this. If you don't set MQTT version, the server will immediately close its
        // connection to your device.
        connectOptions.setMqttVersion(MqttConnectOptions.MQTT_VERSION_3_1_1);

        // With Google Cloud IoT Core, the username field is ignored, however it must be set for the
        // Paho client library to send the password field. The password field is used to transmit a JWT
        // to authorize the device.
        connectOptions.setUserName("unused");

        DateTime iat = new DateTime();
        setConnectPassword(options, connectOptions);
        // [END iot_mqtt_configuremqtt]

        // [START iot_mqtt_publish]
        // Create a client, and connect to the Google MQTT bridge.
        MqttClient client = new MqttClient(mqttServerAddress, mqttClientId, new MemoryPersistence());

        // Both connect and publish operations may fail. If they do, allow retries but with an
        // exponential backoff time period.
        long initialConnectIntervalMillis = 500L;
        long maxConnectIntervalMillis = 6000L;
        long maxConnectRetryTimeElapsedMillis = 900000L;
        float intervalMultiplier = 1.5f;

        long retryIntervalMs = initialConnectIntervalMillis;
        long totalRetryTimeMs = 0;

        while (!client.isConnected() && totalRetryTimeMs < maxConnectRetryTimeElapsedMillis) {
            try {
                client.connect(connectOptions);
            } catch (MqttException e) {
                int reason = e.getReasonCode();

                // If the connection is lost or if the server cannot be connected, allow retries, but with
                // exponential backoff.
                System.out.println("An error occurred: " + e.getMessage());
                if (reason == MqttException.REASON_CODE_CONNECTION_LOST
                    || reason == MqttException.REASON_CODE_SERVER_CONNECT_ERROR) {
                    System.out.println("Retrying in " + retryIntervalMs / 1000.0 + " seconds.");
                    Thread.sleep(retryIntervalMs);
                    totalRetryTimeMs += retryIntervalMs;
                    retryIntervalMs *= intervalMultiplier;
                    if (retryIntervalMs > maxConnectIntervalMillis) {
                        retryIntervalMs = maxConnectIntervalMillis;
                    }
                } else {
                    throw e;
                }
            }
        }

        attachCallback(client, options.deviceId);

        // Publish to the events or state topic based on the flag.
        String subTopic = options.messageType.equals("event") ? "events" : options.messageType;

        // The MQTT topic that this device will publish telemetry data to. The MQTT topic name is
        // required to be in the format below. Note that this is not the same as the device registry's
        // Cloud Pub/Sub topic.
        String mqttTopic = String.format("/devices/%s/%s", options.deviceId, subTopic);

        // Connect to Plc
        logger.info("Connecting to Plc");
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection("s7://10.10.64.20/1/1")) {
            logger.info("Connected");

            PlcReadRequest readRequest = plcConnection.readRequestBuilder().get().addItem("outputs", "OUTPUTS/0").build();

            while (!Thread.currentThread().isInterrupted()) {

                PlcReadResponse plcReadResponse = readRequest.execute().get();

                // Refresh the connection credentials before the JWT expires.
                // [START iot_mqtt_jwt_refresh]
                long secsSinceRefresh = ((new DateTime()).getMillis() - iat.getMillis()) / 1000;
                if (secsSinceRefresh > (options.tokenExpMins * 60)) {
                    System.out.format("\tRefreshing token after: %d seconds\n", secsSinceRefresh);
                    iat = new DateTime();
                    setConnectPassword(options, connectOptions);
                    client.disconnect();
                    client.connect();
                    attachCallback(client, options.deviceId);
                }
                // [END iot_mqtt_jwt_refresh]

                // Send data to cloud
                for (String fieldName : plcReadResponse.getFieldNames()) {
                    Long l = plcReadResponse.getLong(fieldName);
                    byte[] array = ByteBuffer.allocate(8).putLong(l).array();
                    String result = Long.toBinaryString(l);
                    System.out.println("Outputs: " + result);
                    // Publish "array" to the MQTT topic. qos=1 means at least once delivery. Cloud IoT Core
                    // also supports qos=0 for at most once delivery.
                    MqttMessage message = new MqttMessage(array);
                    message.setQos(1);
                    client.publish(mqttTopic, message);
                    if (options.messageType.equals("event")) {
                        // Send telemetry events every second
                        Thread.sleep(1000);
                    } else {
                        // Note: Update Device state less frequently than with telemetry events
                        Thread.sleep(5000);
                    }
                }
            }
        }

        System.out.println("Sent all messages. Goodbye!");
        // [END iot_mqtt_publish]
    }
}
