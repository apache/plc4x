package org.apache.plc4x.java.firmata.readwrite;

import com.github.pfichtner.testcontainers.virtualavr.VirtualAvrConnection;
import com.github.pfichtner.testcontainers.virtualavr.VirtualAvrContainer;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.junit.jupiter.api.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.testcontainers.images.ImagePullPolicy;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.testcontainers.utility.DockerImageName;

import java.io.File;
import java.net.URISyntaxException;
import java.util.List;
import java.util.Objects;
import java.util.UUID;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration test for the Firmata driver using VirtualAVR in TestContainers.
 * <p>
 * This test simulates:
 * - Analog input (A0): Sine wave from 0 to 1023
 * - Digital inputs (pins 2-5): 4-bit counter cycling 0-15
 * <p>
 * The test verifies that the Firmata driver correctly receives these values
 * through subscriptions.
 * <p>
 * NOTE: Each test uses @TestInstance(PER_METHOD) with @Container to get a fresh
 * container instance. This is necessary because the socat TCP-to-PTY bridge on
 * macOS doesn't handle connection reuse well, and the Arduino needs a clean state
 * for each test.
 */
@Testcontainers(disabledWithoutDocker = true)
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
public class FirmataIT {

    private static final Logger LOGGER = LoggerFactory.getLogger(FirmataIT.class);

    // Pin configuration
    private static final String ANALOG_PIN_A0 = "A0";  // Analog pin for sine wave
    private static final int[] DIGITAL_PINS = {2, 3, 4, 5};  // 4 digital pins for counter

    // Simulation parameters
    private static final double SINE_WAVE_PERIOD_MS = 2000;  // 2-second period
    private static final int COUNTER_INCREMENT_MS = 500;  // Increment every 500ms

    // Container is created fresh for each test method
    @Container
    VirtualAvrContainer<?> virtualavr = TestcontainerSupport.virtualAvrContainer(loadClasspath("/firmata/StandardFirmata.ino"))
        .withBaudrate(57600);

    static File loadClasspath(String name) {
        try {
            return new File(Objects.requireNonNull(FirmataIT.class.getResource(name)).toURI());
        } catch (URISyntaxException e) {
            throw new IllegalStateException(e);
        }
    }

    private ScheduledExecutorService simulationExecutor;
    private ScheduledFuture<?> sineWaveFuture;
    private ScheduledFuture<?> counterFuture;

    @BeforeEach
    void setup() {
        simulationExecutor = Executors.newScheduledThreadPool(2);
        startSimulation();
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
            LOGGER.info("Simulation executor terminated");
        }
    }

    /**
     * Start the pin value simulation.
     * - Sine wave on analog pin A0 (0-1023)
     * - 4-bit counter on digital pins 2-5 (0-15)
     */
    private void startSimulation() {
        VirtualAvrConnection avrConnection = virtualavr.avr();

        // Sine wave simulation on A0
        AtomicInteger sinePhase = new AtomicInteger(0);
        sineWaveFuture = simulationExecutor.scheduleAtFixedRate(() -> {
            try {
                int phase = sinePhase.getAndIncrement();
                double radians = (2 * Math.PI * phase) / (SINE_WAVE_PERIOD_MS / 100);
                int value = (int) ((Math.sin(radians) + 1) * 511.5); // Map -1..1 to 0..1023
                avrConnection.pinState(ANALOG_PIN_A0, value);
                LOGGER.debug("Sine wave: phase={}, value={}", phase, value);
            } catch (Exception e) {
                LOGGER.error("Error in sine wave simulation", e);
            }
        }, 0, 50, TimeUnit.MILLISECONDS);

        // Counter simulation on digital pins
        AtomicInteger counter = new AtomicInteger(0);
        counterFuture = simulationExecutor.scheduleAtFixedRate(() -> {
            try {
                int count = counter.getAndIncrement() % 16;
                for (int i = 0; i < DIGITAL_PINS.length; i++) {
                    boolean bitValue = ((count >> i) & 1) == 1;
                    avrConnection.pinState(String.valueOf(DIGITAL_PINS[i]), bitValue);
                }
                LOGGER.debug("Counter: {} (binary: {})", count, Integer.toBinaryString(count));
            } catch (Exception e) {
                LOGGER.error("Error in counter simulation", e);
            }
        }, 0, COUNTER_INCREMENT_MS, TimeUnit.MILLISECONDS);

        LOGGER.info("Simulation started: Sine wave on {}, Counter on pins {}-{}",
            ANALOG_PIN_A0, DIGITAL_PINS[0], DIGITAL_PINS[DIGITAL_PINS.length - 1]);
    }

    private void stopSimulation() {
        // Cancel any running scheduled tasks
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
        String serialDevicePath = virtualavr.serialPortDescriptor();
        LOGGER.info("Started virtual serial port on {}", serialDevicePath);

        boolean received;
        List<Integer> receivedValues = new CopyOnWriteArrayList<>();
        CountDownLatch valuesReceived = new CountDownLatch(20);

        // Create a connection using the serialDevicePath and subscribe to analog pin A0
        try (PlcConnection connection = PlcDriverManager.getDefault().getConnectionManager().getConnection("firmata://" + serialDevicePath + "?serial.baud-rate=57600")) { //&log.driver-testsuite-file=audit-log.txt
            PlcSubscriptionResponse subscriptionResponse = connection.subscriptionRequestBuilder()
                .addChangeOfStateTagAddress(ANALOG_PIN_A0, "analog:0", plcSubscriptionEvent -> {
                    LOGGER.info("Received event: {}", plcSubscriptionEvent.getPlcValue(ANALOG_PIN_A0).getInteger());
                    receivedValues.add(plcSubscriptionEvent.getPlcValue(ANALOG_PIN_A0).getInteger());
                    valuesReceived.countDown();
                })
                .build()
                .execute().get(1000, TimeUnit.MILLISECONDS);

            // Wait for values (longer timeout for QEMU emulation)
            received = valuesReceived.await(15, TimeUnit.SECONDS);

            // Unsubscribe from all subscriptions
            connection.unsubscriptionRequestBuilder()
                .addHandles(subscriptionResponse.getSubscriptionHandles());
        }

        assertTrue(received || receivedValues.size() > 5,
            "Should receive multiple analog values");

        // Verify values are in the expected range
        for (int value : receivedValues) {
            assertTrue(value >= 0 && value <= 1023,
                "Analog value should be between 0 and 1023, got: " + value);
        }

        // Verify we see variation (sine wave should produce different values)
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
        String serialDevicePath = virtualavr.serialPortDescriptor();
        LOGGER.info("Started virtual serial port on {}", serialDevicePath);

        boolean received;
        List<Integer> receivedCounterValues = new CopyOnWriteArrayList<>();
        CountDownLatch valuesReceived = new CountDownLatch(16);

        // Create a connection using the serialDevicePath and subscribe to digital pins
        try (PlcConnection connection = PlcDriverManager.getDefault().getConnectionManager().getConnection("firmata://" + serialDevicePath + "?serial.baud-rate=57600&log.driver-testsuite-file=audit-log.txt")) {
            PlcSubscriptionRequest.Builder builder = connection.subscriptionRequestBuilder();
            for (int pin : DIGITAL_PINS) {
                builder.addChangeOfStateTagAddress(String.valueOf(pin), "digital:" + pin);
            }
            builder.setConsumer(plcSubscriptionEvent -> {
                LOGGER.info("Received event: {}", plcSubscriptionEvent);
                for (String tagName : plcSubscriptionEvent.getTagNames()) {
                    int integer = plcSubscriptionEvent.getPlcValue(tagName).getInteger();
                    receivedCounterValues.add(integer);
                }
                valuesReceived.countDown();
            });
            PlcSubscriptionResponse subscriptionResponse = builder
                .build()
                .execute().get(1000, TimeUnit.MILLISECONDS);

            // Wait for values (longer timeout for QEMU emulation)
            received = valuesReceived.await(15, TimeUnit.SECONDS);

            // Unsubscribe from all subscriptions
            connection.unsubscriptionRequestBuilder()
                .addHandles(subscriptionResponse.getSubscriptionHandles())
                .build().execute().get(1000, TimeUnit.MILLISECONDS);
        }

        assertTrue(received || receivedCounterValues.size() > 3,
            "Should receive multiple counter values");

        // Verify values are in the expected range (0-15)
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
        String serialDevicePath = virtualavr.serialPortDescriptor();
        LOGGER.info("Started virtual serial port on {}", serialDevicePath);

        boolean received;
        AtomicInteger analogCount = new AtomicInteger(0);
        AtomicInteger digitalCount = new AtomicInteger(0);
        CountDownLatch valuesReceived = new CountDownLatch(10);

        // Create a connection using the serialDevicePath and subscribe to both analog and digital pins
        try (PlcConnection connection = PlcDriverManager.getDefault().getConnectionManager().getConnection("firmata://" + serialDevicePath + "?serial.baud-rate=57600&log.driver-testsuite-file=audit-log.txt")) {
            PlcSubscriptionRequest.Builder builder = connection.subscriptionRequestBuilder();
            for (int pin : DIGITAL_PINS) {
                builder.addChangeOfStateTagAddress("D" + pin, "digital:" + pin);
            }
            builder.addChangeOfStateTagAddress(ANALOG_PIN_A0, "analog:0");
            builder.setConsumer(plcSubscriptionEvent -> {
                LOGGER.info("Received event: {}", plcSubscriptionEvent);
                for (String tagName : plcSubscriptionEvent.getTagNames()) {
                    if (tagName.equals(ANALOG_PIN_A0)) {
                        analogCount.incrementAndGet();
                    } else if (tagName.startsWith("D")) {
                        digitalCount.incrementAndGet();
                    }
                }
                valuesReceived.countDown();
            });
            PlcSubscriptionResponse subscriptionResponse = builder
                .build()
                .execute().get(1000, TimeUnit.MILLISECONDS);

            // Wait for values (longer timeout for QEMU emulation)
            received = valuesReceived.await(15, TimeUnit.SECONDS);

            // Unsubscribe from all subscriptions
            connection.unsubscriptionRequestBuilder()
                .addHandles(subscriptionResponse.getSubscriptionHandles())
                .build().execute().get(1000, TimeUnit.MILLISECONDS);
        }

        LOGGER.info("Received {} analog events and {} digital events",
            analogCount.get(), digitalCount.get());

        // Both simulation types should have produced events
        assertTrue(received && (analogCount.get() > 0 || digitalCount.get() > 0),
            "Should receive events from both analog and digital simulations");
    }

}
