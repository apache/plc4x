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
package org.apache.plc4x.java.utils.cache;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.utils.cache.exceptions.PlcConnectionManagerClosedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashSet;
import java.util.Set;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * A connection manager that caches and reuses PLC connections.
 * <p>
 * This implementation provides:
 * - Connection pooling with automatic reuse
 * - Idle timeout (closes connections unused for too long)
 * - Max lease timeout (force-returns connections leased too long)
 * - Connection validation via ping (only for connections idle beyond threshold)
 * - Thread-safe access to cached connections
 * - Proper cleanup on shutdown
 * <p>
 * Usage Example:
 * <pre>
 * // Get the default connection manager (supports all drivers)
 * PlcConnectionManager connectionManager = PlcDriverManager.getDefault();
 *
 * CachedPlcConnectionManager manager = CachedPlcConnectionManager.getBuilder()
 *     .withConnectionManager(connectionManager)
 *     .withMaxIdleTime(5, TimeUnit.MINUTES)
 *     .withMaxLeaseTime(1, TimeUnit.MINUTES)
 *     .build();
 *
 * try (PlcConnection conn = manager.getConnection("ads:tcp://192.168.1.1")) {
 *     // Use connection - works with any protocol!
 * } // Automatically returned to cache
 *
 * manager.close(); // Shutdown all connections
 * </pre>
 * <p>
 * Thread Safety: This class is fully thread-safe. Multiple threads can call getConnection()
 * concurrently, and connections will be properly shared/reused.
 * <p>
 * Deadlock Prevention:
 * - Uses ConcurrentHashMap to avoid synchronization on the entire pool
 * - Uses per-connection {@link java.util.concurrent.locks.ReentrantLock}s (in ConnectionContainer)
 *   for fine-grained locking. ReentrantLock (rather than {@code synchronized}) keeps the cache
 *   virtual-thread-friendly: a virtual thread that blocks while holding or awaiting one of these
 *   locks parks and releases its carrier instead of pinning it.
 * - A connection's I/O (connect / ping / close) runs under that connection's own lock only, so it
 *   never blocks operations on other connections; the pool itself is never globally locked.
 * - Timeout tasks are cancelled before closing connections to avoid races.
 */
public class CachedPlcConnectionManager implements PlcConnectionManager, AutoCloseable {

    private static final Logger LOGGER = LoggerFactory.getLogger(CachedPlcConnectionManager.class);

    // Default configuration
    private static final long DEFAULT_MAX_IDLE_TIME_MS = TimeUnit.MINUTES.toMillis(5);
    private static final long DEFAULT_MAX_LEASE_TIME_MS = TimeUnit.MINUTES.toMillis(1);
    private static final long DEFAULT_MAX_WAIT_TIME_MS = TimeUnit.SECONDS.toMillis(30);
    private static final long DEFAULT_PING_TIMEOUT_MS = TimeUnit.SECONDS.toMillis(5);
    private static final long DEFAULT_IDLE_PING_THRESHOLD_MS = TimeUnit.SECONDS.toMillis(30);
    private static final long DEFAULT_CLOSE_TIMEOUT_MS = TimeUnit.SECONDS.toMillis(5);

    /**
     * Grace margin added to the caller-facing lease {@code get()} timeout, on top of the per-connection
     * max-wait. It ensures the container's own wait timeout (scheduled at exactly max-wait) fires and
     * marks a queued waiter "done" before this manager gives up — preventing a returned connection from
     * being handed to a waiter whose caller has already timed out. Purely an upper safety bound; the
     * {@code get()} returns as soon as the future completes.
     */
    private static final long LEASE_WAIT_GRACE_MS = TimeUnit.SECONDS.toMillis(1);

    private final PlcConnectionManager connectionManager;
    private final ScheduledExecutorService scheduler;
    private final long maxIdleTimeMs;
    private final long maxLeaseTimeMs;
    private final long maxWaitTimeMs;
    private final long pingTimeoutMs;
    private final long idlePingThresholdMs;
    private final long closeTimeoutMs;

    // Connection pool: connectionString -> ConnectionContainer
    // ConcurrentHashMap provides thread-safe access without global synchronization
    private final ConcurrentHashMap<String, ConnectionContainer> cachedConnections = new ConcurrentHashMap<>();

    // Shutdown flag to prevent use after close
    private volatile boolean closed = false;

    /**
     * Private constructor - use Builder to create instances.
     */
    private CachedPlcConnectionManager(PlcConnectionManager connectionManager, ScheduledExecutorService scheduler,
                                       long maxIdleTimeMs, long maxLeaseTimeMs, long maxWaitTimeMs,
                                       long pingTimeoutMs, long idlePingThresholdMs, long closeTimeoutMs) {
        this.connectionManager = connectionManager;
        this.scheduler = scheduler;
        this.maxIdleTimeMs = maxIdleTimeMs;
        this.maxLeaseTimeMs = maxLeaseTimeMs;
        this.maxWaitTimeMs = maxWaitTimeMs;
        this.pingTimeoutMs = pingTimeoutMs;
        this.idlePingThresholdMs = idlePingThresholdMs;
        this.closeTimeoutMs = closeTimeoutMs;

        LOGGER.info("Created CachedPlcConnectionManager with maxIdle={}ms, maxLease={}ms, maxWait={}ms, pingTimeout={}ms, idlePingThreshold={}ms, closeTimeout={}ms",
            maxIdleTimeMs, maxLeaseTimeMs, maxWaitTimeMs, pingTimeoutMs, idlePingThresholdMs, closeTimeoutMs);
    }

    @Override
    public PlcConnection getConnection(String connectionString) throws PlcConnectionException {
        return getConnection(connectionString, null);
    }

    @Override
    public PlcConnection getConnection(String connectionString, PlcAuthentication authentication) throws PlcConnectionException {
        if (closed) {
            throw new PlcConnectionManagerClosedException();
        }

        LOGGER.debug("Requesting connection for: {}", connectionString);

        // Get or create connection atomically with validation
        // compute() ensures all validation happens atomically for the first thread
        // The returned ConnectionContainer is stored in the map
        ConnectionContainer container = cachedConnections.computeIfAbsent(connectionString, (key) -> {
            LOGGER.debug("Creating new connection container for: {}", connectionString);
            return new ConnectionContainer(connectionString, () -> {
                if (authentication != null) {
                    return connectionManager.getConnection(connectionString, authentication);
                } else {
                    return connectionManager.getConnection(connectionString);
                }}, scheduler,
                maxIdleTimeMs,
                maxLeaseTimeMs,
                maxWaitTimeMs,
                pingTimeoutMs,
                idlePingThresholdMs,
                closeTimeoutMs);
        });

        // Lease the connection - THIS WILL BLOCK IF ALREADY LEASED
        Future<PlcConnection> leaseFuture = container.lease();
        try {
            // Wait a grace margin BEYOND the container's own max-wait timeout. The container already
            // guarantees the lease future completes within maxWaitTimeMs (it schedules its own wait
            // timeout), so letting our get() time out at exactly maxWaitTimeMs would race that: if our
            // get() won, we'd abandon a queue entry that is still "not done", and a concurrent return
            // could then hand the connection to that phantom waiter — wedging the connection. The grace
            // lets the container's timeout win and mark the entry done first; get() still returns the
            // instant the future completes, so a normal wait-timeout is not slowed down.
            return leaseFuture.get(maxWaitTimeMs + LEASE_WAIT_GRACE_MS, TimeUnit.MILLISECONDS);
        } catch (ExecutionException | InterruptedException | TimeoutException e) {
            throw new PlcConnectionException("Error acquiring lease for connection", e);
        }
    }

    @Override
    public void close() throws Exception {
        if (closed) {
            LOGGER.debug("Connection manager already closed");
            return;
        }

        closed = true;
        LOGGER.info("Closing CachedPlcConnectionManager with {} connections", cachedConnections.size());

        // Shut the scheduler down first to prevent new timeout tasks from starting
        // This is critical to prevent scheduler threads from trying to close connections
        // while we're also closing them
        scheduler.shutdownNow();

        // Close all connections
        // Iterate over a copy to avoid ConcurrentModificationException
        cachedConnections.values().forEach(container -> {
            try {
                LOGGER.debug("Closing connection '{}'", container.getConnectionString());
                container.close();
            } catch (Exception e) {
                if (LOGGER.isTraceEnabled()) {
                    LOGGER.error("Error closing connection '{}'", container.getConnectionString(), e);
                } else {
                    LOGGER.error("Error closing connection '{}': {}", container.getConnectionString(), e.getMessage());
                }
            }
        });

        // Clear the pool
        cachedConnections.clear();

        LOGGER.info("CachedPlcConnectionManager closed");
    }

    /**
     * Get the number of cached connections.
     *
     * @return Number of connections in the pool
     */
    public int getCachedConnectionCount() {
        return cachedConnections.size();
    }

    /**
     * Get the number of active leases across all connections.
     *
     * @return Total active lease count
     */
    public int getActiveLeaseCount() {
        return cachedConnections.values().stream()
            .mapToInt(value -> value.isLeased() ? 1 : 0)
            .sum();
    }

    /**
     * Get a snapshot of the connection strings currently held in the cache.
     * <p>
     * The returned set is a copy — mutating it does not affect the cache, and the cache may change
     * concurrently after this call. Useful for monitoring and for deciding which connections to
     * {@link #removeCachedConnection(String) evict}.
     *
     * @return a snapshot {@link Set} of the cached connection strings (never {@code null})
     */
    public Set<String> getCachedConnections() {
        return new HashSet<>(cachedConnections.keySet());
    }

    /**
     * Forcibly evict and close a cached connection by its connection string.
     * <p>
     * The connection is removed from the cache first (so the next {@link #getConnection(String)}
     * transparently establishes a fresh one) and then closed. If the connection is currently leased,
     * it is force-closed regardless — the lease holder's next operation will fail, mirroring the
     * max-lease-timeout behavior. This is the manual counterpart to the automatic idle/lease eviction
     * and is intended for management/operations tooling (e.g. dropping a connection whose device was
     * reconfigured or whose credentials changed).
     *
     * @param connectionString the connection string to evict.
     * @return {@code true} if a cached connection was found and evicted; {@code false} if none was cached.
     */
    public boolean removeCachedConnection(String connectionString) {
        // Remove from the map first so a concurrent getConnection() builds a fresh container rather
        // than racing on the one we are about to close.
        ConnectionContainer container = cachedConnections.remove(connectionString);
        if (container == null) {
            LOGGER.debug("removeCachedConnection: no cached connection for '{}'", connectionString);
            return false;
        }
        try {
            LOGGER.info("Forcibly evicting cached connection '{}'", connectionString);
            container.close();
        } catch (Exception e) {
            if (LOGGER.isTraceEnabled()) {
                LOGGER.error("Error closing evicted connection '{}'", connectionString, e);
            } else {
                LOGGER.error("Error closing evicted connection '{}': {}", connectionString, e.getMessage());
            }
        }
        return true;
    }

    /**
     * Create a new builder for CachedPlcConnectionManager.
     *
     * @return A new builder instance
     */
    public static Builder getBuilder() {
        return new Builder();
    }

    /**
     * Builder for CachedPlcConnectionManager.
     * <p>
     * Provides a fluent API for configuring the connection manager.
     */
    public static class Builder {
        private PlcConnectionManager connectionManager;
        private ScheduledExecutorService scheduler = defaultScheduler();
        private long maxIdleTimeMs = DEFAULT_MAX_IDLE_TIME_MS;
        private long maxLeaseTimeMs = DEFAULT_MAX_LEASE_TIME_MS;
        private long maxWaitTimeMs = DEFAULT_MAX_WAIT_TIME_MS;
        private long pingTimeoutMs = DEFAULT_PING_TIMEOUT_MS;
        private long idlePingThresholdMs = DEFAULT_IDLE_PING_THRESHOLD_MS;
        private long closeTimeoutMs = DEFAULT_CLOSE_TIMEOUT_MS;

        /**
         * Set the PLC connection manager to use for creating connections.
         * This allows the cache to handle connections for multiple PLC types/protocols.
         *
         * @param connectionManager The PLC connection manager (typically PlcDriverManager)
         * @return This builder
         */
        public Builder withConnectionManager(PlcConnectionManager connectionManager) {
            this.connectionManager = connectionManager;
            return this;
        }

        /**
         * Set the scheduler to use for the idle / lease / wait / validation timeout tasks.
         * <p>
         * By default, a small daemon {@link ScheduledExecutorService} pool is created (see
         * {@code defaultScheduler()}). Provide your own if you want to share a scheduler across
         * managers or control the scheduling threads. The supplied scheduler is shut down by
         * {@link CachedPlcConnectionManager#close()}.
         *
         * @param scheduler the scheduled executor service
         * @return This builder
         */
        public Builder withScheduler(ScheduledExecutorService scheduler) {
            this.scheduler = scheduler;
            return this;
        }

        /**
         * Set the maximum idle time before a connection is closed.
         * <p>
         * If a connection is not leased for this amount of time, it will be
         * automatically closed and removed from the pool.
         *
         * @param maxIdleTime The maximum idle time
         * @param timeUnit    The time unit
         * @return This builder
         */
        public Builder withMaxIdleTime(long maxIdleTime, TimeUnit timeUnit) {
            this.maxIdleTimeMs = timeUnit.toMillis(maxIdleTime);
            return this;
        }

        /**
         * Set the maximum lease time before a connection is force-returned.
         * <p>
         * If a connection is leased for longer than this time, it will be
         * automatically closed. This prevents connection leaks from clients
         * that forget to return connections.
         * <p>
         * Set to 0 to disable max lease timeout.
         *
         * @param maxLeaseTime The maximum lease time
         * @param timeUnit     The time unit
         * @return This builder
         */
        public Builder withMaxLeaseTime(long maxLeaseTime, TimeUnit timeUnit) {
            this.maxLeaseTimeMs = timeUnit.toMillis(maxLeaseTime);
            return this;
        }

        /**
         * Set the timeout for connection ping validation.
         * <p>
         * When a cached connection needs to be validated (because it has been
         * idle beyond the idle ping threshold), this timeout controls how long
         * to wait for the ping response.
         *
         * @param pingTimeout The ping timeout
         * @param timeUnit    The time unit
         * @return This builder
         */
        public Builder withPingTimeout(long pingTimeout, TimeUnit timeUnit) {
            this.pingTimeoutMs = timeUnit.toMillis(pingTimeout);
            return this;
        }

        /**
         * Set the maximum wait time for acquiring a lease.
         * <p>
         * If a thread is waiting for a connection lease and doesn't get it within
         * this time, a TimeoutException will be thrown. This prevents threads from
         * waiting indefinitely for busy connections.
         * <p>
         * Default: 30 seconds
         *
         * @param maxWaitTime The maximum wait time
         * @param timeUnit    The time unit
         * @return This builder
         */
        public Builder withMaxWaitTime(long maxWaitTime, TimeUnit timeUnit) {
            this.maxWaitTimeMs = timeUnit.toMillis(maxWaitTime);
            return this;
        }

        /**
         * Set the idle ping threshold.
         * <p>
         * Cached connections that haven't been used for this amount of time
         * will be pinged before being handed out to verify they're still alive.
         * This prevents unnecessary pings for recently used connections.
         * <p>
         * Default: 30 seconds
         *
         * @param idlePingThreshold The idle ping threshold
         * @param timeUnit          The time unit
         * @return This builder
         */
        public Builder withIdlePingThreshold(long idlePingThreshold, TimeUnit timeUnit) {
            this.idlePingThresholdMs = timeUnit.toMillis(idlePingThreshold);
            return this;
        }

        /**
         * Set the maximum time to wait for an underlying {@code connection.close()} to complete.
         * <p>
         * Closing a connection on a wedged socket can hang indefinitely while the per-connection lock
         * is held, which would freeze every operation on that connection (lease, return, eviction).
         * Bounding the close caps that worst case: after this time the close is abandoned to a
         * background daemon thread and the cache recovers (the next acquisition establishes a fresh
         * connection).
         * <p>
         * Set to 0 to disable the bound (legacy synchronous close — wait indefinitely).
         * Default: 5 seconds.
         *
         * @param closeTimeout The maximum close time
         * @param timeUnit     The time unit
         * @return This builder
         */
        public Builder withCloseTimeout(long closeTimeout, TimeUnit timeUnit) {
            this.closeTimeoutMs = timeUnit.toMillis(closeTimeout);
            return this;
        }

        /** Default number of daemon scheduler threads (see {@link #defaultScheduler()}). */
        private static final int DEFAULT_SCHEDULER_THREADS = 2;

        /**
         * Creates the default daemon scheduler used when the caller does not supply one via
         * {@link #withScheduler(ScheduledExecutorService)}.
         * <p>
         * A small <em>pool</em> (rather than a single thread) is used so that a timeout task which
         * blocks on a busy container's lock — while that container holds the lock during a slow
         * connect / ping / close — cannot stall the timeout tasks of all other containers. Every
         * task body either acquires the per-container lock (so same-container tasks still serialize)
         * or just completes a future, so concurrent execution is safe.
         */
        private static ScheduledExecutorService defaultScheduler() {
            AtomicInteger threadIndex = new AtomicInteger();
            return Executors.newScheduledThreadPool(DEFAULT_SCHEDULER_THREADS, runnable -> {
                Thread thread = new Thread(runnable,
                    "CachedPlcConnectionManager-Scheduler-" + threadIndex.incrementAndGet());
                thread.setDaemon(true);
                return thread;
            });
        }

        /**
         * Build the CachedPlcConnectionManager.
         *
         * @return A new CachedPlcConnectionManager instance
         * @throws IllegalStateException if connectionManager is not set
         */
        public CachedPlcConnectionManager build() {
            if (connectionManager == null) {
                throw new IllegalStateException("PlcConnectionManager must be set using withConnectionManager()");
            }
            return new CachedPlcConnectionManager(
                connectionManager,
                scheduler,
                maxIdleTimeMs,
                maxLeaseTimeMs,
                maxWaitTimeMs,
                pingTimeoutMs,
                idlePingThresholdMs,
                closeTimeoutMs
            );
        }
    }

}
