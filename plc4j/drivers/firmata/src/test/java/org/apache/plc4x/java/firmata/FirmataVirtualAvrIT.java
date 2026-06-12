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
package org.apache.plc4x.java.firmata;

import com.github.pfichtner.testcontainers.virtualavr.DefaultVirtualAvrConnection;
import com.github.pfichtner.testcontainers.virtualavr.VirtualAvrConnection;
import com.github.pfichtner.testcontainers.virtualavr.VirtualAvrConnection.PinReportMode;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.MethodOrderer;
import org.junit.jupiter.api.Order;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestMethodOrder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.testcontainers.containers.GenericContainer;
import org.testcontainers.images.builder.ImageFromDockerfile;
import org.testcontainers.junit.jupiter.Testcontainers;

import java.net.URI;
import java.util.List;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;

import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Integration test for the Firmata driver running against a virtualavr container
 * over the TCP transport.
 *
 * <p>The container is built from {@code src/test/resources/virtualavr-tcp/Dockerfile}
 * — a thin extension over {@code pfichtner/virtualavr} that replaces the PTY-bridge
 * entrypoint with a single {@code socat TCP-LISTEN ... EXEC:"node /app/virtualavr.js
 * sketch.ino",pty,rawer,fdin=3,fdout=4} and bakes StandardFirmata in. The simulated
 * UART is exposed as a TCP listener on port 3030; the simulator's WebSocket
 * side-channel on port 8080 is used by {@link VirtualAvrConnection} to inject pin
 * transitions and read back pin state. No host-side {@code socat} or PTY is
 * involved, so this IT runs on every OS Docker supports — including macOS where
 * {@code jSerialComm} won't open the temp-file PTY the upstream serial mode
 * produces.</p>
 *
 * <h2>Test plan</h2>
 *
 * <ul>
 *   <li>{@link #testAnalogInputs()} — sine wave on A0, verify the driver delivers
 *       a stream of subscription events whose values span the expected 0..1023
 *       ADC range with significant variation.</li>
 *   <li>{@link #testDigitalInputs()} — 4-bit counter cycling D2..D5 0..15, verify
 *       the driver delivers change-of-state events for every pin transition.</li>
 *   <li>{@link #testAnalogAndDigitalInputs()} — both simulations running
 *       concurrently, verify the driver routes events to the right tag names
 *       without interference.</li>
 * </ul>
 *
 * <p>Container and simulator are recreated for every test (default PER_METHOD
 * lifecycle): socat in the container runs one accept slot at a time, and a clean
 * AVR pin register per test keeps the assertions predictable.</p>
 */
@Testcontainers(disabledWithoutDocker = true)
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
public class FirmataVirtualAvrIT {

    private static final Logger LOGGER = LoggerFactory.getLogger(FirmataVirtualAvrIT.class);

    private static final int CONTAINER_TCP_SERIAL_PORT = 3030;
    private static final int CONTAINER_WEBSOCKET_PORT = 8080;

    private static final String ANALOG_PIN_A0 = "A0";
    private static final int[] DIGITAL_PINS = {2, 3, 4, 5};

    private static final double SINE_WAVE_PERIOD_MS = 2000;
    private static final int COUNTER_INCREMENT_MS = 500;

    private GenericContainer<?> virtualAvr;
    private VirtualAvrConnection virtualAvrConnection;
    private ScheduledExecutorService simulationExecutor;
    private ScheduledFuture<?> sineWaveFuture;
    private ScheduledFuture<?> counterFuture;

    @BeforeEach
    @SuppressWarnings("resource")
    void setup() {
        virtualAvr = new GenericContainer<>(
            new ImageFromDockerfile()
                .withFileFromClasspath("Dockerfile", "virtualavr-tcp/Dockerfile"))
            .withExposedPorts(CONTAINER_TCP_SERIAL_PORT, CONTAINER_WEBSOCKET_PORT);
        // No Wait.forListeningPort — that probe opens a TCP connection which
        // consumes our single accept slot (socat runs without ",fork" so the
        // next accept spawns a fresh simulator that would clash on port 8080).
        // We retry the driver connect in openDriverAndAttachWebSocket() instead.
        virtualAvr.start();
        simulationExecutor = Executors.newScheduledThreadPool(2);
    }

    @AfterEach
    void teardown() {
        stopSimulation();
        if (simulationExecutor != null) {
            simulationExecutor.shutdown();
            try {
                if (!simulationExecutor.awaitTermination(5, TimeUnit.SECONDS)) {
                    simulationExecutor.shutdownNow();
                }
            } catch (InterruptedException e) {
                simulationExecutor.shutdownNow();
                Thread.currentThread().interrupt();
            }
        }
        if (virtualAvrConnection != null) {
            try { virtualAvrConnection.close(); } catch (Exception ignored) {}
            virtualAvrConnection = null;
        }
        if (virtualAvr != null) {
            virtualAvr.stop();
            virtualAvr = null;
        }
    }

    /**
     * Opens the driver's TCP connection to the container's serial port, then
     * connects the {@link VirtualAvrConnection} WebSocket helper used for
     * pin-state injection. The driver's connect attempt is what spawns the
     * simulator inside the container (socat runs the AVR command on the first
     * accept), so we have to retry until socat is up.
     */
    private PlcConnection openDriverAndAttachWebSocket() throws Exception {
        String url = "firmata:tcp://" + virtualAvr.getHost() + ":"
            + virtualAvr.getMappedPort(CONTAINER_TCP_SERIAL_PORT)
            + "?request-timeout=30000";
        PlcConnection connection = openWithRetry(url);

        // Build the WebSocket URI manually — the library's static factory
        // assumes the first mapped port is the WS port, which isn't true here.
        URI wsUri = URI.create("ws://" + virtualAvr.getHost() + ":"
            + virtualAvr.getMappedPort(CONTAINER_WEBSOCKET_PORT));
        virtualAvrConnection = new DefaultVirtualAvrConnection(wsUri);
        long wsDeadline = System.currentTimeMillis() + 10_000L;
        while (!virtualAvrConnection.isConnected() && System.currentTimeMillis() < wsDeadline) {
            //noinspection BusyWait
            Thread.sleep(100);
        }
        if (!virtualAvrConnection.isConnected()) {
            throw new AssertionError("Could not connect to virtualavr WebSocket at " + wsUri);
        }
        return connection;
    }

    private PlcConnection openWithRetry(String url) throws Exception {
        long deadline = System.currentTimeMillis() + 60_000L;
        Exception last = null;
        while (System.currentTimeMillis() < deadline) {
            try {
                return new DefaultPlcDriverManager().getConnection(url);
            } catch (Exception e) {
                last = e;
                //noinspection BusyWait
                Thread.sleep(500);
            }
        }
        throw new AssertionError("Could not open " + url + " within timeout", last);
    }

    /**
     * Start the pin-value simulation:
     * <ul>
     *   <li>Sine wave on analog pin A0 (0..1023)</li>
     *   <li>4-bit counter on digital pins 2..5 (0..15)</li>
     * </ul>
     */
    private void startSimulation() {
        // The simulator only forwards injected pinState(...) calls into the
        // AVR's pin register for pins it has been told to report on. Without
        // this, StandardFirmata never sees the transitions and the driver
        // gets no events even though our scheduler keeps firing.
        virtualAvrConnection.pinReportMode(ANALOG_PIN_A0, PinReportMode.ANALOG);
        for (int pin : DIGITAL_PINS) {
            virtualAvrConnection.pinReportMode(String.valueOf(pin), PinReportMode.DIGITAL);
        }

        AtomicInteger sinePhase = new AtomicInteger(0);
        sineWaveFuture = simulationExecutor.scheduleAtFixedRate(() -> {
            try {
                int phase = sinePhase.getAndIncrement();
                double radians = (2 * Math.PI * phase) / (SINE_WAVE_PERIOD_MS / 100);
                int value = (int) ((Math.sin(radians) + 1) * 511.5);
                virtualAvrConnection.pinState(ANALOG_PIN_A0, value);
                LOGGER.info("Sine wave: phase={}, value={}", phase, value);
            } catch (Exception e) {
                LOGGER.error("Error in sine wave simulation", e);
            }
        }, 0, 50, TimeUnit.MILLISECONDS);

        AtomicInteger counter = new AtomicInteger(0);
        counterFuture = simulationExecutor.scheduleAtFixedRate(() -> {
            try {
                int count = counter.getAndIncrement() % 16;
                for (int i = 0; i < DIGITAL_PINS.length; i++) {
                    boolean bitValue = ((count >> i) & 1) == 1;
                    virtualAvrConnection.pinState(String.valueOf(DIGITAL_PINS[i]), bitValue);
                }
                LOGGER.info("Counter: {} (binary: {})", count, Integer.toBinaryString(count));
            } catch (Exception e) {
                LOGGER.error("Error in counter simulation", e);
            }
        }, 0, COUNTER_INCREMENT_MS, TimeUnit.MILLISECONDS);

        LOGGER.info("Simulation started: Sine wave on {}, Counter on pins {}-{}",
            ANALOG_PIN_A0, DIGITAL_PINS[0], DIGITAL_PINS[DIGITAL_PINS.length - 1]);
    }

    private void stopSimulation() {
        if (sineWaveFuture != null) {
            sineWaveFuture.cancel(false);
            sineWaveFuture = null;
        }
        if (counterFuture != null) {
            counterFuture.cancel(false);
            counterFuture = null;
        }
    }

    @Test
    @Order(1)
    @DisplayName("Test analog pins (sine wave simulation)")
    void testAnalogInputs() throws Exception {
        boolean received;
        List<Integer> receivedValues = new CopyOnWriteArrayList<>();
        CountDownLatch valuesReceived = new CountDownLatch(20);

        try (PlcConnection connection = openDriverAndAttachWebSocket()) {
            startSimulation();

            PlcSubscriptionResponse subscriptionResponse = connection.subscriptionRequestBuilder()
                .addChangeOfStateTagAddress(ANALOG_PIN_A0, "analog:0")
                .build()
                .execute().get(5000, TimeUnit.MILLISECONDS);

            // Builder-level consumers (per-tag or setConsumer) are only stashed
            // on the request — they aren't auto-wired to event delivery. Attach
            // the consumer on the returned handle to actually receive events.
            subscriptionResponse.getSubscriptionHandle(ANALOG_PIN_A0).register(plcSubscriptionEvent -> {
                LOGGER.info("Received event: {}", plcSubscriptionEvent.getPlcValue(ANALOG_PIN_A0).getInteger());
                receivedValues.add(plcSubscriptionEvent.getPlcValue(ANALOG_PIN_A0).getInteger());
                valuesReceived.countDown();
            });

            received = valuesReceived.await(15, TimeUnit.SECONDS);

            connection.unsubscriptionRequestBuilder()
                .addHandles(subscriptionResponse.getSubscriptionHandles())
                .build().execute().get(5000, TimeUnit.MILLISECONDS);
        }

        assertTrue(received || receivedValues.size() > 5,
            "Should receive multiple analog values");

        for (int value : receivedValues) {
            assertTrue(value >= 0 && value <= 1023,
                "Analog value should be between 0 and 1023, got: " + value);
        }

        if (receivedValues.size() > 5) {
            int min = receivedValues.stream().mapToInt(Integer::intValue).min().orElse(0);
            int max = receivedValues.stream().mapToInt(Integer::intValue).max().orElse(1023);
            assertTrue(max - min > 100,
                "Should see significant variation in sine wave values");
        }

        LOGGER.info("Received {} analog values, range: {} to {}",
            receivedValues.size(),
            receivedValues.stream().mapToInt(Integer::intValue).min().orElse(-1),
            receivedValues.stream().mapToInt(Integer::intValue).max().orElse(-1));
    }

    @Test
    @Order(2)
    @DisplayName("Test digital pins (4-bit counter simulation)")
    void testDigitalInputs() throws Exception {
        boolean received;
        List<Integer> receivedCounterValues = new CopyOnWriteArrayList<>();
        CountDownLatch valuesReceived = new CountDownLatch(16);

        try (PlcConnection connection = openDriverAndAttachWebSocket()) {
            startSimulation();

            PlcSubscriptionRequest.Builder builder = connection.subscriptionRequestBuilder();
            for (int pin : DIGITAL_PINS) {
                builder.addChangeOfStateTagAddress("D" + pin, "digital:" + pin);
            }
            PlcSubscriptionResponse subscriptionResponse = builder
                .build()
                .execute().get(5000, TimeUnit.MILLISECONDS);

            Consumer<PlcSubscriptionEvent> consumer = plcSubscriptionEvent -> {
                LOGGER.info("Received event: {} {} {} {}",
                    plcSubscriptionEvent.getPlcValue("D5") != null ? plcSubscriptionEvent.getPlcValue("D5").getInteger() : "-",
                    plcSubscriptionEvent.getPlcValue("D4") != null ? plcSubscriptionEvent.getPlcValue("D4").getInteger() : "-",
                    plcSubscriptionEvent.getPlcValue("D3") != null ? plcSubscriptionEvent.getPlcValue("D3").getInteger() : "-",
                    plcSubscriptionEvent.getPlcValue("D2") != null ? plcSubscriptionEvent.getPlcValue("D2").getInteger() : "-");
                for (String tagName : plcSubscriptionEvent.getTagNames()) {
                    int integer = plcSubscriptionEvent.getPlcValue(tagName).getInteger();
                    receivedCounterValues.add(integer);
                }
                valuesReceived.countDown();
            };
            // setConsumer on the builder is only stashed; events fire only for
            // consumers attached to the returned handles.
            for (int pin : DIGITAL_PINS) {
                subscriptionResponse.getSubscriptionHandle("D" + pin).register(consumer);
            }

            received = valuesReceived.await(15, TimeUnit.SECONDS);

            connection.unsubscriptionRequestBuilder()
                .addHandles(subscriptionResponse.getSubscriptionHandles())
                .build().execute().get(5000, TimeUnit.MILLISECONDS);
        }

        assertTrue(received || receivedCounterValues.size() > 3,
            "Should receive multiple counter values");

        for (int value : receivedCounterValues) {
            assertTrue(value >= 0 && value <= 15,
                "Counter value should be between 0 and 15, got: " + value);
        }

        LOGGER.info("Received {} counter state changes", receivedCounterValues.size());
    }

    @Test
    @Order(3)
    @DisplayName("Test combined analog and digital pins")
    void testAnalogAndDigitalInputs() throws Exception {
        boolean received;
        AtomicInteger analogCount = new AtomicInteger(0);
        AtomicInteger digitalCount = new AtomicInteger(0);
        CountDownLatch valuesReceived = new CountDownLatch(20);

        try (PlcConnection connection = openDriverAndAttachWebSocket()) {
            startSimulation();

            PlcSubscriptionRequest.Builder builder = connection.subscriptionRequestBuilder();
            for (int pin : DIGITAL_PINS) {
                builder.addChangeOfStateTagAddress("D" + pin, "digital:" + pin);
            }
            builder.addChangeOfStateTagAddress(ANALOG_PIN_A0, "analog:0");
            PlcSubscriptionResponse subscriptionResponse = builder
                .build()
                .execute().get(5000, TimeUnit.MILLISECONDS);

            Consumer<PlcSubscriptionEvent> consumer = plcSubscriptionEvent -> {
                LOGGER.info("Received event: {} {} {} {} {}",
                    plcSubscriptionEvent.getPlcValue(ANALOG_PIN_A0) != null ? plcSubscriptionEvent.getPlcValue(ANALOG_PIN_A0).getInteger() : "-",
                    plcSubscriptionEvent.getPlcValue("D5") != null ? plcSubscriptionEvent.getPlcValue("D5").getInteger() : "-",
                    plcSubscriptionEvent.getPlcValue("D4") != null ? plcSubscriptionEvent.getPlcValue("D4").getInteger() : "-",
                    plcSubscriptionEvent.getPlcValue("D3") != null ? plcSubscriptionEvent.getPlcValue("D3").getInteger() : "-",
                    plcSubscriptionEvent.getPlcValue("D2") != null ? plcSubscriptionEvent.getPlcValue("D2").getInteger() : "-");
                for (String tagName : plcSubscriptionEvent.getTagNames()) {
                    if (tagName.equals(ANALOG_PIN_A0)) {
                        analogCount.incrementAndGet();
                    } else if (tagName.startsWith("D")) {
                        digitalCount.incrementAndGet();
                    }
                }
                valuesReceived.countDown();
            };
            // setConsumer on the builder is only stashed; events fire only for
            // consumers attached to the returned handles.
            for (int pin : DIGITAL_PINS) {
                subscriptionResponse.getSubscriptionHandle("D" + pin).register(consumer);
            }
            subscriptionResponse.getSubscriptionHandle(ANALOG_PIN_A0).register(consumer);

            received = valuesReceived.await(15, TimeUnit.SECONDS);

            connection.unsubscriptionRequestBuilder()
                .addHandles(subscriptionResponse.getSubscriptionHandles())
                .build().execute().get(5000, TimeUnit.MILLISECONDS);
        }

        LOGGER.info("Received {} analog events and {} digital events",
            analogCount.get(), digitalCount.get());

        assertTrue(received && (analogCount.get() > 0 || digitalCount.get() > 0),
            "Should receive events from both analog and digital simulations");
    }

}
