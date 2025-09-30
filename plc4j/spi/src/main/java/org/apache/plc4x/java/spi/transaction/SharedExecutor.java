package org.apache.plc4x.java.spi.transaction;

import org.apache.commons.lang3.concurrent.BasicThreadFactory;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

/**
 * A new class created to centrally manage executors, which would otherwise be created and never removed for each connection request.
 * That is, given a connection, ThreadPools were created for each request through it; these only died at the end of the connection, so in
 * the case of a cached connection, the ThreadPools accumulated. This way, there is always only one shared ThreadPool: nothing changes for the normal connection
 * while avoiding the problem with the cached connection.
 *
 * NB: the number of threads here is set to 10 (in the original, non-shared connection, there were 4) for safety.
 * If necessary (e.g. massive several connections to independent PLCs), it can be increased.
 */
public class SharedExecutor {
    /**
     * Usato da opcua driver (e forse altri.. ma non s7)
     */
    private static final ExecutorService tmExecutor = Executors.newFixedThreadPool(
            10,
            new BasicThreadFactory.Builder()
                    .namingPattern("plc4x-tm-thread-%d")
                    .daemon(true)
                    .priority(Thread.MAX_PRIORITY)
                    .build()
    );

    /**
     * Usato solo da s7 driver
     */
    private static final ExecutorService appExecutor = Executors.newFixedThreadPool(
            10,
            new BasicThreadFactory.Builder()
                    .namingPattern("plc4x-app-thread-%d")
                    .daemon(true)
                    .priority(Thread.MAX_PRIORITY)
                    .build()
    );

    public static ExecutorService getTmExecutor() {
        return tmExecutor;
    }

    public static ExecutorService getAppExecutor() {
        return appExecutor;
    }

    public static void shutdown() {
        tmExecutor.shutdownNow();
        appExecutor.shutdownNow();
    }
}
