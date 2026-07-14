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
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.LinkedList;
import java.util.Queue;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Future;
import java.util.concurrent.RejectedExecutionException;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.locks.ReentrantLock;

/**
 * Internal container that tracks a single connection's state in the cache.
 * <p>
 * This class manages:
 * - The underlying connection instance
 * - Lease tracking (who's using it, how long)
 * - Idle timeout (close if unused too long)
 * - Max lease timeout (force return if leased too long)
 * <p>
 * Thread Safety: This class uses a ReentrantLock to ensure thread-safe access to
 * the connection state. The lock prevents race conditions during lease/return operations.
 * Using ReentrantLock provides better performance under contention compared to synchronized.
 * <p>
 * Deadlock Prevention: Always acquires locks in the same order and releases them promptly.
 * Timer tasks are canceled before attempting to close connections to prevent deadlocks.
 * Every blocking I/O performed under the lock is time-bounded — connect (see {@link #connectBounded()},
 * bounded by {@code maxWaitTimeMs}), ping ({@code pingTimeoutMs}) and close ({@code closeTimeoutMs}) — so
 * a hung driver can never hold the lock indefinitely and starve other callers of the same connection.
 */
class ConnectionContainer {

    private static final Logger LOGGER = LoggerFactory.getLogger(ConnectionContainer.class);

    private final String connectionString;
    private final ThrowingSupplier<PlcConnection> connectionFactory;
    private final ScheduledExecutorService scheduler;
    private final long maxIdleTimeMs;
    private final long maxLeaseTimeMs;
    private final long maxWaitTimeMs;
    private final long pingTimeoutMs;
    private final long idlePingThresholdMs;
    private final long closeTimeoutMs;
    private final Queue<QueueEntry> queue;
    private final ReentrantLock lock;

    // Tracks subscriptions and event listeners for state recovery after reconnection
    private final ConnectionStateTracker stateTracker;

    private PlcConnection connection;
    // When true, the next lease MUST validate (ping) the connection regardless of how recently it
    // was used. Set when a lease is force-returned (e.g. by a max-lease timeout) because the
    // connection may be dead. In all other cases validation is decided purely from idle time at
    // lease() time, so it is deterministic and needs no scheduled task.
    private boolean forceValidationOnNextLease;
    // Tracks the last time a lease was successfully returned (i.e., the connection was used).
    // Used to determine if a ping validation is actually needed: if a successful operation
    // completed recently, the connection is known to be alive and pinging is unnecessary.
    private long lastSuccessfulReturnTimeMs;
    // Volatile so the lease state can be read without acquiring the lock (see isLeased()): a
    // monitoring call must never block behind a container that is mid-connect/ping/close under
    // the lock. All WRITES still happen under the lock; only the monitoring READ is lock-free.
    private volatile LeasedPlcConnection leasedConnection;

    // Monotonically increasing lease identity counter. Each new lease gets the next ID.
    // Used to detect and ignore stale callbacks from old leases that fire after a new
    // lease has been created (e.g., timeout-triggered return followed by stale invalidation).
    private long currentLeaseId;

    // Scheduled timeout tasks (must be canceled to prevent leaks)
    private volatile ScheduledFuture<?> idleTimeoutTask;
    private volatile ScheduledFuture<?> maxLeaseTimeoutTask;

    /**
     * Convenience constructor with an <em>unbounded</em> close (legacy behavior): the underlying
     * {@code connection.close()} runs synchronously with no timeout. Retained for internal/test use.
     */
    public ConnectionContainer(String connectionString, ThrowingSupplier<PlcConnection> connectionFactory,
                               ScheduledExecutorService scheduler,
                               long maxIdleTimeMs, long maxLeaseTimeMs, long maxWaitTimeMs,
                               long pingTimeoutMs, long idlePingThresholdMs) {
        this(connectionString, connectionFactory, scheduler, maxIdleTimeMs, maxLeaseTimeMs, maxWaitTimeMs,
            pingTimeoutMs, idlePingThresholdMs, 0);
    }

    public ConnectionContainer(String connectionString, ThrowingSupplier<PlcConnection> connectionFactory,
                               ScheduledExecutorService scheduler,
                               long maxIdleTimeMs, long maxLeaseTimeMs, long maxWaitTimeMs,
                               long pingTimeoutMs, long idlePingThresholdMs, long closeTimeoutMs) {
        this.connectionString = connectionString;
        this.connectionFactory = connectionFactory;
        this.scheduler = scheduler;
        this.maxIdleTimeMs = maxIdleTimeMs;
        this.maxLeaseTimeMs = maxLeaseTimeMs;
        this.maxWaitTimeMs = maxWaitTimeMs;
        this.pingTimeoutMs = pingTimeoutMs;
        this.idlePingThresholdMs = idlePingThresholdMs;
        this.closeTimeoutMs = closeTimeoutMs;
        this.queue = new LinkedList<>();
        this.forceValidationOnNextLease = false;
        this.lastSuccessfulReturnTimeMs = 0;
        this.currentLeaseId = 0;
        this.lock = new ReentrantLock(true); // Fair lock to prevent thread starvation
        this.stateTracker = new ConnectionStateTracker(connectionString);
    }

    /**
     * Lease this connection to a caller.
     * <p>
     * BLOCKS if the connection is already leased by another thread.
     * Only one client can have access to the connection at a time.
     * <p>
     * Increments the active lease count and cancels the idle timeout.
     * If maxLeaseTimeMs > 0, schedules a task to force-return the connection.
     *
     * @return A LeasedPlcConnection wrapper
     * @throws IllegalStateException if the connection is closed
     */
    public Future<PlcConnection> lease() {
        lock.lock();
        try {
            CompletableFuture<PlcConnection> connectionFuture = new CompletableFuture<>();

            // If the connection is currently idle, return the connection immediately.
            if (leasedConnection == null) {
                // Try to get a new connection if we haven't got one yet.
                if(connection == null) {
                    try {
                        connection = connectBounded();
                    } catch (PlcConnectionException e) {
                        if (LOGGER.isTraceEnabled()) {
                            LOGGER.warn("Exception while getting connection for lease", e);
                        } else {
                            LOGGER.warn("Exception while getting connection for lease: {}", e.getMessage());
                        }
                        connectionFuture.completeExceptionally(e);
                        return connectionFuture;
                    }
                }
                // If the connection was idle in the pool for too long, we need to validate it first.
                // If validation fails, we need to close the existing connection and get a new one.
                // Skip validation if a successful operation completed recently — the connection
                // is known to be alive and pinging would just add unnecessary blocking latency.
                //
                // This decision is computed entirely here, at lease() time, from the idle duration
                // (plus the force flag). It does NOT depend on any scheduled task having fired, so
                // it is fully deterministic. A zero idlePingThreshold means "always validate", which
                // the (now - lastReturn) >= 0 comparison already yields.
                else if (forceValidationOnNextLease
                    || lastSuccessfulReturnTimeMs == 0
                    || (System.currentTimeMillis() - lastSuccessfulReturnTimeMs) >= idlePingThresholdMs) {
                    if (!validate()) {
                        LOGGER.debug("Connection failed validation, closing and getting new connection");
                        closeBounded(connection);
                        try {
                            connection = connectBounded();
                            // Restore tracked state (subscriptions, event listeners) on new connection
                            if (stateTracker.hasStateToRestore()) {
                                stateTracker.restoreState(connection);
                            }
                        } catch (PlcConnectionException e) {
                            if (LOGGER.isTraceEnabled()) {
                                LOGGER.warn("Exception while getting connection for lease", e);
                            } else {
                                LOGGER.warn("Exception while getting connection for lease: {}", e.getMessage());
                            }
                            connectionFuture.completeExceptionally(e);
                            return connectionFuture;
                        }
                    }
                }

                // Cancel idle timeout since the connection is now in use
                cancelIdleTimeout();

                // Create the lease first (increments currentLeaseId), then schedule
                // the timeout so it captures the correct lease identity
                leasedConnection = createLeasedConnection();

                // Schedule max lease timeout if configured
                if (maxLeaseTimeMs > 0) {
                    scheduleMaxLeaseTimeout();
                }

                connectionFuture.complete(leasedConnection);
            }
            // Otherwise add the future to the queue for completion as soon as the connection is returned.
            else {
                // Schedule a task that times out the wait if it's taking too long.
                ScheduledFuture<?> waitTimeoutTask = scheduler.schedule(
                    () -> connectionFuture.completeExceptionally(new TimeoutException("Connection wait timed out")),
                    maxWaitTimeMs, TimeUnit.MILLISECONDS);
                QueueEntry queueEntry = new QueueEntry(connectionFuture, waitTimeoutTask);
                queue.add(queueEntry);
            }

            return connectionFuture;
        } finally {
            lock.unlock();
        }
    }

    private LeasedPlcConnection createLeasedConnection() {
        currentLeaseId++;
        long leaseId = currentLeaseId;
        return new LeasedPlcConnection(connectionString, connection, stateTracker, (conn, invalidated, lid) -> {
            if(invalidated) {
                invalidateLease(lid);
            } else {
                returnLease(lid, false);
            }
        }, leaseId);
    }

    /**
     * Return a lease to this connection.
     * <p>
     * Stops the timers for the max lease timeout and starts the idle timeout if no other client was waiting.
     * If another client is waiting for access to this connection, return a new lease to that client and start
     * the max lease timeout.
     *
     * @param leaseId The lease identity of the returning connection. If this doesn't match
     *                the current active lease, the return is silently ignored (stale callback).
     * @param forceValidation if true, always mark the connection as requiring validation before reuse.
     *                        This is set to true when the max lease timeout fires, ensuring that
     *                        connections returned due to timeout are validated (pinged) before being
     *                        handed out again — preventing dead connections from being recycled.
     */
    private void returnLease(long leaseId, boolean forceValidation) {
        lock.lock();
        try {
            // Verify this callback is from the current active lease — ignore stale returns
            if (leaseId != currentLeaseId) {
                LOGGER.debug("Ignoring stale return from lease {}, current is {}: {}", leaseId, currentLeaseId, connectionString);
                return;
            }

            // Stop the timeout handler.
            cancelMaxLeaseTimeout();
            // Track the time of this successful return — used to skip unnecessary
            // ping validation when the connection was recently used successfully.
            if (!forceValidation) {
                lastSuccessfulReturnTimeMs = System.currentTimeMillis();
            }
            LOGGER.trace("Returned connection {}", connectionString);

            // If the queue is empty, clear the leased connection and start the idle timeout.
            if(queue.isEmpty()) {
                leasedConnection = null;
                scheduleIdleTimeout();
                // A force-return means the connection may be dead: validate on the next lease no
                // matter how recently it was used. Otherwise leave it to the idle-time check in
                // lease() — no scheduled task needed.
                forceValidationOnNextLease = forceValidation;
                LOGGER.trace("Connection idle {}", connectionString);
                return;
            }

            // If the connection was previously invalidated, get a new connection.
            if (connection == null) {
                try {
                    connection = connectBounded();
                    // Restore tracked state (subscriptions, event listeners) on new connection
                    if (stateTracker.hasStateToRestore()) {
                        stateTracker.restoreState(connection);
                    }
                } catch (PlcConnectionException e) {
                    // If creating the connection fails, close all waiting client's futures exceptionally.
                    do {
                        QueueEntry queueEntry = queue.poll();
                        if (queueEntry != null) {
                            queueEntry.future.completeExceptionally(e);
                        }
                    } while(!queue.isEmpty());
                    return;
                }
            }

            // If there were other clients waiting in the queue, take the one waiting the longest
            // and check if it's already timed out. If it is, go through the queue until we find one
            // that hasn't timed out yet. Then return a new lease to that client.
            while (!queue.isEmpty()) {
                QueueEntry queueEntry = queue.poll();
                if ((queueEntry != null) && (!queueEntry.future.isDone())) {
                    leasedConnection = createLeasedConnection();
                    // Schedule max lease timeout if configured
                    if (maxLeaseTimeMs > 0) {
                        scheduleMaxLeaseTimeout();
                    }
                    queueEntry.waitTimeoutTask.cancel(false);
                    queueEntry.future.complete(leasedConnection);
                    return;
                }
            }

            // All queued entries had already timed out — go idle
            leasedConnection = null;
            scheduleIdleTimeout();
            // See the queue-empty branch above: force validation only on a force-return; otherwise
            // the idle-time check in lease() decides deterministically.
            forceValidationOnNextLease = forceValidation;
            LOGGER.trace("Connection idle after all waiters timed out {}", connectionString);
        } finally {
            lock.unlock();
        }
    }

    /**
     * Invalidate this connection due to an error.
     * <p>
     * This immediately closes the connection, marks it as closed, and invokes the invalidation callback
     * to allow the owner to remove it from any tracking structures.
     * Used when an exception occurs during connection use.
     * <p>
     * If the leaseId doesn't match the current active lease, the invalidation is silently ignored.
     * This prevents stale callbacks from old leases from destroying active connections.
     *
     * @param leaseId The lease identity of the invalidating connection
     */
    private void invalidateLease(long leaseId) {
        lock.lock();
        try {
            // Verify this callback is from the current active lease — ignore stale invalidations
            if (leaseId != currentLeaseId) {
                LOGGER.debug("Ignoring stale invalidation from lease {}, current is {}: {}", leaseId, currentLeaseId, connectionString);
                return;
            }

            LOGGER.warn("Invalidating connection due to error: {}", connectionString);
            // Clear the leased connection so the container is available for new leases
            leasedConnection = null;
            // Close the underlying connection
            closeInternal();

            // Fail all queued waiters immediately so they don't leave stale entries
            // that could cause a phantom lease in returnLease()
            PlcConnectionException invalidationError = new PlcConnectionException("Connection was invalidated: " + connectionString);
            while (!queue.isEmpty()) {
                QueueEntry queueEntry = queue.poll();
                if (queueEntry != null) {
                    queueEntry.waitTimeoutTask.cancel(false);
                    queueEntry.future.completeExceptionally(invalidationError);
                }
            }
        } finally {
            lock.unlock();
        }
    }

    /**
     * Validate this connection to ensure it's still usable.
     * <p>
     * This method checks if the connection needs validation based on idle time,
     * and performs a ping if necessary. Returns true if the connection is valid
     * and can be used, false otherwise.
     * <p>
     * This method handles all validation logic internally, removing that
     * responsibility from the CachedPlcConnectionManager.
     *
     * @return true if the connection is valid, false if it should be discarded
     */
    public boolean validate() {
        // Validate by checking connection status and pinging
        if ((connection == null) || !connection.isConnected()) {
            LOGGER.debug("Connection reports not connected: {}", connectionString);
            return false;
        }

        // Perform ping with timeout
        try {
            java.util.concurrent.CompletableFuture<?> pingFuture = connection.ping();
            pingFuture.get(pingTimeoutMs, java.util.concurrent.TimeUnit.MILLISECONDS);
            LOGGER.trace("Connection ping successful '{}'", connectionString);
            return true;

        } catch (java.util.concurrent.TimeoutException e) {
            LOGGER.warn("Connection ping timeout '{}'", connectionString);
            return false;

        } catch (Exception e) {
            if (LOGGER.isTraceEnabled()) {
                LOGGER.warn("Connection ping failed '{}'", connectionString, e);
            } else {
                LOGGER.warn("Connection ping failed '{}': {}", connectionString, e.getMessage());
            }
            return false;
        }
    }

    /**
     * Schedule a task to close the connection if it remains idle.
     */
    private void scheduleIdleTimeout() {
        // Cancel any existing idle timeout
        cancelIdleTimeout();

        try {
            idleTimeoutTask = scheduler.schedule(this::onIdleTimeout, maxIdleTimeMs, TimeUnit.MILLISECONDS);
            LOGGER.trace("Scheduled idle timeout for {} in {}ms", connectionString, maxIdleTimeMs);
        } catch (RejectedExecutionException e) {
            // Scheduler was already shut down (manager shutting down)
            LOGGER.debug("Cannot schedule idle timeout, scheduler shut down: {}", connectionString);
        }
    }

    /**
     * Fired by the scheduler when the idle timeout elapses. Acquires the lock and re-checks that the
     * connection is still idle before closing it: a lease may have been granted between the time this
     * task was scheduled and the time it actually runs (the scheduler cannot interrupt an
     * already-started task), so closing unconditionally would race with an active lease.
     */
    private void onIdleTimeout() {
        lock.lock();
        try {
            if (leasedConnection == null) {
                closeInternal();
            }
        } finally {
            lock.unlock();
        }
    }

    /**
     * Schedule a task to force-return the connection if leased too long.
     * The client using it can no longer use it from now on.
     */
    private void scheduleMaxLeaseTimeout() {
        // Cancel any existing max lease timeout
        cancelMaxLeaseTimeout();

        // Capture the lease ID at scheduling time so the timeout handler always
        // references the correct lease even if a new lease is created before it fires
        long scheduledLeaseId = currentLeaseId;
        try {
            maxLeaseTimeoutTask = scheduler.schedule(() -> {
                LOGGER.warn("Connection lease timeout reached, forcing close: {}", connectionString);
                // Force validation on timeout-returned connections — the connection may be
                // dead due to a network outage that caused the operation to hang
                returnLease(scheduledLeaseId, true);
            }, maxLeaseTimeMs, TimeUnit.MILLISECONDS);
            LOGGER.trace("Scheduled max lease timeout for {} in {}ms", connectionString, maxLeaseTimeMs);
        } catch (RejectedExecutionException e) {
            // Scheduler was already shut down (manager shutting down)
            LOGGER.debug("Cannot schedule max lease timeout, scheduler shut down: {}", connectionString);
        }
    }

    /**
     * Cancel the idle timeout task if running.
     */
    private void cancelIdleTimeout() {
        if (idleTimeoutTask != null) {
            idleTimeoutTask.cancel(false);
            idleTimeoutTask = null;
            LOGGER.trace("Cancelled idle timeout for {}", connectionString);
        }
    }

    /**
     * Cancel the max lease timeout task if running.
     */
    private void cancelMaxLeaseTimeout() {
        if (maxLeaseTimeoutTask != null) {
            maxLeaseTimeoutTask.cancel(false);
            maxLeaseTimeoutTask = null;
            LOGGER.trace("Cancelled max lease timeout for {}", connectionString);
        }
    }


    /**
     * Close this connection container and the underlying connection.
     * <p>
     * Cancels all timer tasks and closes the connection. After calling this,
     * the container should not be used anymore.
     */
    public void close() {
        lock.lock();
        try {
            if (leasedConnection != null) {
                leasedConnection.close();
            }
            closeInternal();
        } finally {
            lock.unlock();
        }
    }

    /**
     * Internal close method (must be called with lock held).
     * <p>
     * This prevents duplicate close operations and ensures timers are canceled
     * before closing the connection to avoid deadlocks.
     * <p>
     * Note: This method is reentrant-safe because it's called from methods that
     * already hold the ReentrantLock (invalidateLease, close).
     */
    private void closeInternal() {
        // Cancel all timeout tasks BEFORE closing connection
        // This is critical to prevent races where a timeout task
        // tries to acquire the lock while we're closing
        cancelIdleTimeout();
        cancelMaxLeaseTimeout();

        // Drop the reference FIRST so the next lease reconnects even if the close hangs and gets
        // abandoned by the close timeout below.
        PlcConnection toClose = connection;
        connection = null;
        closeBounded(toClose);
    }

    /**
     * Establishes a connection without ever blocking the caller — and therefore the container lock —
     * for longer than {@code maxWaitTimeMs}. The underlying driver {@code connect()} can otherwise hang
     * indefinitely (e.g. a wedged TCP handshake, or a driver awaiting a {@code CompletableFuture} that
     * never completes). Because this runs under the container lock, an unbounded connect would hold the
     * lock forever and freeze every other operation on this connection — including waiters, who block on
     * {@code lock.lock()} in {@link #lease()} before they ever obtain a future and so never see their own
     * {@code maxWaitTimeMs} timeout. Bounding the connect turns "stuck forever" into "fails after
     * {@code maxWaitTimeMs}, dead connect abandoned to a background daemon thread", honouring the wait
     * timeout the caller was promised. A non-positive {@code maxWaitTimeMs} disables the bound (legacy
     * synchronous connect).
     *
     * @return the freshly established connection (never {@code null})
     * @throws PlcConnectionException if the connect fails or does not complete within {@code maxWaitTimeMs}
     */
    private PlcConnection connectBounded() throws PlcConnectionException {
        if (maxWaitTimeMs <= 0) {
            return connectionFactory.get();
        }
        CompletableFuture<PlcConnection> result = new CompletableFuture<>();
        Thread connector = new Thread(() -> {
            try {
                PlcConnection established = connectionFactory.get();
                // If the caller already gave up (timed out below), it will never use this connection —
                // close it here instead of leaking it. complete() returns false iff we lost that race.
                if (!result.complete(established)) {
                    closeQuietly(established);
                }
            } catch (Throwable t) {
                result.completeExceptionally(t);
            }
        }, "cache-connect-" + connectionString);
        connector.setDaemon(true);
        connector.start();
        try {
            return result.get(maxWaitTimeMs, TimeUnit.MILLISECONDS);
        } catch (TimeoutException e) {
            // Abandon the hanging connect. If the connector produced a connection in the tiny race
            // window between get() timing out and this line, completeExceptionally() fails — recover
            // and close that now-orphaned connection so it is not leaked.
            if (!result.completeExceptionally(e)) {
                result.thenAccept(this::closeQuietly);
            }
            connector.interrupt(); // best-effort; the underlying connect may not honor interruption
            LOGGER.warn("Connection establishment did not complete within {}ms; abandoning it to the "
                + "background and failing this attempt: {}", maxWaitTimeMs, connectionString);
            throw new PlcConnectionException("Connection establishment did not complete within "
                + maxWaitTimeMs + "ms: " + connectionString, e);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PlcConnectionException("Interrupted while establishing connection: " + connectionString, e);
        } catch (ExecutionException e) {
            Throwable cause = e.getCause();
            if (cause instanceof PlcConnectionException) {
                throw (PlcConnectionException) cause;
            }
            throw new PlcConnectionException("Error establishing connection: " + connectionString, cause);
        }
    }

    /**
     * Closes {@code conn} without ever blocking the caller — and therefore the container lock — for
     * longer than {@code closeTimeoutMs}. A {@code close()} on a wedged socket can otherwise hang
     * indefinitely while the lock is held, which would freeze every operation on this connection
     * (lease, return, eviction). Bounding it turns "stuck forever" into "recovers after
     * {@code closeTimeoutMs}, dead close abandoned to a background daemon thread". A non-positive
     * {@code closeTimeoutMs} disables the bound (legacy synchronous close).
     */
    private void closeBounded(PlcConnection conn) {
        if (conn == null) {
            return;
        }
        if (closeTimeoutMs <= 0) {
            closeQuietly(conn);
            return;
        }
        Thread closer = new Thread(() -> closeQuietly(conn), "cache-close-" + connectionString);
        closer.setDaemon(true);
        closer.start();
        try {
            closer.join(closeTimeoutMs);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
        if (closer.isAlive()) {
            LOGGER.warn("Connection close did not complete within {}ms; abandoning it to the background "
                + "and continuing: {}", closeTimeoutMs, connectionString);
            closer.interrupt(); // best-effort; the underlying close may not honor interruption
        }
    }

    /**
     * Closes {@code conn}, swallowing and logging any error (a close failure must never propagate out
     * of the cache's own bookkeeping).
     */
    private void closeQuietly(PlcConnection conn) {
        try {
            conn.close();
            LOGGER.debug("Closed connection '{}'", connectionString);
        } catch (Exception e) {
            if (LOGGER.isTraceEnabled()) {
                LOGGER.error("Error closing connection '{}'", connectionString, e);
            } else {
                LOGGER.error("Error closing connection '{}': {}", connectionString, e.getMessage());
            }
        }
    }

    /**
     * Check if this connection is currently leased.
     *
     * @return true if leased, false otherwise
     */
    public boolean isLeased() {
        // Lock-free read of the volatile lease state: a monitoring call (e.g. getActiveLeaseCount)
        // must not block behind a container that is mid-connect/ping/close under the lock.
        return leasedConnection != null;
    }

    /**
     * Check if this connection is closed.
     *
     * @return true if closed, false otherwise
     */
    public boolean isClosed() {
        lock.lock();
        try {
            return connection == null;
        } finally {
            lock.unlock();
        }
    }

    /**
     * Get the underlying connection (for testing).
     *
     * @return The wrapped connection
     */
    PlcConnection getConnection() {
        lock.lock();
        try {
            return connection;
        } finally {
            lock.unlock();
        }
    }

    public String getConnectionString() {
        return connectionString;
    }

    /**
     * Get the state tracker for this connection container.
     * <p>
     * The state tracker maintains subscriptions and event listeners that should
     * be restored when the connection is recreated.
     *
     * @return The connection state tracker
     */
    ConnectionStateTracker getStateTracker() {
        return stateTracker;
    }

    @FunctionalInterface
    public interface ThrowingSupplier<T> {
        T get() throws PlcConnectionException;
    }

    private record QueueEntry(CompletableFuture<PlcConnection> future, ScheduledFuture<?> waitTimeoutTask) {}

}
