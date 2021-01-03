package org.apache.plc4x.java.utils.connectionpool;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.junit.jupiter.api.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.*;

/**
 * TODO write comment
 *
 * @author julian
 * Created by julian on 24.02.20
 */
class CachedDriverManagerMT {

    private static final Logger logger = LoggerFactory.getLogger(CachedDriverManagerMT.class);

    public static final String PLC_IP = "s7://192.168.167.210/1/1";
//    public static final String PLC_IP = "s7://127.0.0.1/1/1";

    @Test
    void queryPlc() throws InterruptedException {
        ScheduledExecutorService pool = Executors.newScheduledThreadPool(4);

        CachedDriverManager driverManager = new CachedDriverManager(PLC_IP, () -> {
            try {
                PlcConnection connection = new PlcDriverManager().getConnection(PLC_IP);
                // Kill it every second
                pool.schedule(() -> {
                    try {
                        System.out.println("Close...");
                        connection.close();
                    } catch (Exception e) {
                        e.printStackTrace();
                    }
                }, 1, TimeUnit.SECONDS);
                return connection;
            } catch (PlcConnectionException e) {
                throw new RuntimeException("", e);
            }
        });

        for (int i = 1; i <= 100_000; i++) {
            pool.submit(() -> {
                try (PlcConnection conn = driverManager.getConnection(PLC_IP)) {
                    PlcReadResponse response = conn.readRequestBuilder().addItem("asdf", "%DB444:14.0:BOOL").build().execute().get(500, TimeUnit.MILLISECONDS);
                } catch (Exception e) {
                    // Intentionally do nothing...
                }
            });
            Thread.sleep(1);
        }

        pool.shutdown();
        pool.awaitTermination(1, TimeUnit.DAYS);
    }

    @Test
    void queryPlcWithPool() throws InterruptedException {
        ScheduledExecutorService pool = Executors.newScheduledThreadPool(4);

        PooledDriverManager driverManager = new PooledDriverManager();

        for (int i = 1; i <= 100_000; i++) {
            pool.submit(() -> {
                try (PlcConnection conn = driverManager.getConnection(PLC_IP)) {
                    PlcReadResponse response = conn.readRequestBuilder().addItem("asdf", "%DB444:14.0:BOOL").build().execute().get(500, TimeUnit.MILLISECONDS);
                } catch (Exception e) {
                    // Intentionally do nothing...
                }
            });
            Thread.sleep(1);
        }

        pool.shutdown();
        pool.awaitTermination(1, TimeUnit.DAYS);
    }

    @Test
    void borrowAndDontReturn() throws InterruptedException {
        ScheduledExecutorService pool = Executors.newScheduledThreadPool(4);

        PooledDriverManager driverManager = new PooledDriverManager();

        for (int i = 1; i <= 100_000; i++) {
            pool.submit(() -> {
                try {
                    PlcConnection conn = driverManager.getConnection(PLC_IP);
                    System.out.println("Successfully got a Connection...");
                    PlcReadResponse response = conn.readRequestBuilder().addItem("asdf", "%DB444:14.0:BOOL").build().execute().get(500, TimeUnit.MILLISECONDS);
                    System.out.println("Response: " + response.getBoolean("asdf"));
                } catch (PlcConnectionException | InterruptedException | ExecutionException | TimeoutException e) {
                    // Intentionally do Nothing...
                }
            });
            Thread.sleep(1);
        }

        pool.shutdown();
        pool.awaitTermination(1, TimeUnit.DAYS);
    }
}