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
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcBrowseRequest;
import org.apache.plc4x.java.api.messages.PlcBrowseRequestInterceptor;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;
import org.apache.plc4x.java.api.messages.PlcPingRequest;
import org.apache.plc4x.java.api.messages.PlcPingResponse;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.s7.configuration.S7Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Collection;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicReference;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Supplier;

/**
 * Highly-available S7 connection wrapping two underlying {@link S7CotpConnection} instances —
 * one to a primary host and one to a secondary host. Routes every operation to whichever
 * inner connection is currently "active", and transparently fails over to the other when
 * the active one drops.
 *
 * <h3>Use cases</h3>
 * <ul>
 *   <li><b>S7-400H redundant CPU pair</b> — two CPUs hot-syncing on separate chassis. The
 *       wrapper talks to both, switches when one dies. The classic SPI1 use case.</li>
 *   <li><b>Network-level redundancy on a single CPU</b> — e.g. an S7-300 with both an
 *       on-board PN port and a CP 343-1 module exposed at separate IPs. Same data on
 *       both interfaces, but the connection survives a cable / NIC / switch failure.</li>
 * </ul>
 *
 * <h3>Behaviour</h3>
 * <ul>
 *   <li>{@code connect()} attempts both inner connections in parallel; succeeds if at
 *       least one comes up. Active defaults to the first that connected, with a primary
 *       preference.</li>
 *   <li>Each request is dispatched to the active inner. If the operation fails because
 *       the inner's transport is no longer connected, the wrapper swaps to the other
 *       inner and retries once.</li>
 *   <li>A scheduled heartbeat task periodically reconnects dropped inners and switches
 *       back to the primary when it recovers (configurable preference).</li>
 *   <li>{@code close()} closes both inners. The heartbeat is stopped first.</li>
 * </ul>
 *
 * <h3>Subscriptions and failover</h3>
 * Push subscriptions (alarm, cyclic) are PLC-side state keyed to a TCP connection — when
 * a TCP drops, the PLC drops its corresponding subscription state. This wrapper does
 * <em>not</em> currently re-subscribe automatically on failover. Callers wanting transparent
 * subscription failover should listen for connection-lost events and re-subscribe in
 * application code; the wrapper itself only ensures reads/writes/pings keep working.
 */
public class S7HCotpConnection extends ConnectionBase<S7Configuration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(S7HCotpConnection.class);

    /**
     * Heartbeat tick interval and wrapper-level failover timeout, both in milliseconds.
     * Read from {@link S7Configuration#getHaHeartbeatInterval()} /
     * {@link S7Configuration#getHaFailoverTimeout()} at connection construction so each
     * connection can be tuned independently via URL params {@code ?ha-heartbeat-interval=}
     * and {@code ?ha-failover-timeout=}.
     */
    private final long heartbeatIntervalMs;
    private final long failoverTimeoutMs;

    /**
     * Factories that build fresh inner connections from the configured URL. Used both for
     * the initial connect and to <em>rebuild</em> an inner whose transport has died — the
     * underlying TCP socket can't be reopened once the kernel has closed it (cable pull,
     * peer reset), so the heartbeat task drops the corpse and constructs a new one.
     */
    private final Supplier<S7CotpConnection> primaryFactory;
    private final Supplier<S7CotpConnection> secondaryFactory;
    /** Current inner connections. Replaced by the heartbeat when an inner is rebuilt. */
    private final AtomicReference<S7CotpConnection> primaryRef = new AtomicReference<>();
    private final AtomicReference<S7CotpConnection> secondaryRef = new AtomicReference<>();
    /** Whichever ref ({@link #primaryRef} or {@link #secondaryRef}) is currently active. */
    private final AtomicReference<AtomicReference<S7CotpConnection>> activeSlot = new AtomicReference<>();
    /**
     * Inners we've recently given up on (after a wrapper-level timeout or operation
     * failure) and that the heartbeat task should re-validate before we trust them again.
     * Keyed by slot reference (primary/secondary), not by the inner instance, so it
     * survives rebuilds.
     */
    private final java.util.Set<AtomicReference<S7CotpConnection>> stale = java.util.concurrent.ConcurrentHashMap.newKeySet();
    private final ScheduledExecutorService heartbeat;

    public S7HCotpConnection(S7Configuration configuration,
                             TransportInstance<?> primaryTransport,
                             AuditLog auditLog,
                             Supplier<S7CotpConnection> primaryFactory,
                             Supplier<S7CotpConnection> secondaryFactory) {
        // The base class wants a transport. We never actually use it (request lifecycle is
        // delegated to inner connections), but ConnectionBase needs a non-null reference.
        super(configuration, primaryTransport, auditLog);
        this.heartbeatIntervalMs = configuration.getHaHeartbeatInterval();
        this.failoverTimeoutMs = configuration.getHaFailoverTimeout();
        this.primaryFactory = primaryFactory;
        this.secondaryFactory = secondaryFactory;
        this.primaryRef.set(primaryFactory.get());
        this.secondaryRef.set(secondaryFactory.get());
        this.heartbeat = Executors.newSingleThreadScheduledExecutor(r -> {
            Thread t = new Thread(r, "S7H-Heartbeat");
            t.setDaemon(true);
            return t;
        });
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        // Reuse the same tag handler as the standard S7 connection — both inners use it.
        return new org.apache.plc4x.java.s7.tag.S7PlcTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new DefaultPlcValueHandler();
    }

    @Override
    @SuppressWarnings("resource")
    public boolean isConnected() {
        return current(primaryRef).isConnected() || current(secondaryRef).isConnected();
    }

    @Override
    @SuppressWarnings("resource")
    public PlcConnectionMetadata getMetadata() {
        AtomicReference<S7CotpConnection> slot = activeSlot.get();
        S7CotpConnection a = slot != null ? slot.get() : current(primaryRef);
        PlcConnectionMetadata inner = a.getMetadata();
        // Subscriptions (alarm push, cyclic services) are PLC-side state per TCP connection.
        // The wrapper can't transparently re-establish them on the survivor when an inner
        // dies, so we hide the capability rather than let callers set up subscriptions
        // that silently break on failover. Application code that wants subscriptions on a
        // dual-path setup should connect to one host directly, listen for connection-lost
        // events, and re-subscribe in app code.
        return new PlcConnectionMetadata() {
            @Override public boolean isReadSupported()      { return inner.isReadSupported(); }
            @Override public boolean isWriteSupported()     { return inner.isWriteSupported(); }
            @Override public boolean isSubscribeSupported() { return false; }
            @Override public boolean isBrowseSupported()    { return inner.isBrowseSupported(); }
        };
    }

    @Override
    protected void onConnect() throws PlcConnectionException {
        // Connect both inners in parallel; succeed if at least one comes up. Failure of a
        // single inner is logged but doesn't abort the wrapper — the user might be testing
        // with one of the IPs unreachable on purpose.
        CompletableFuture<Boolean> primaryFuture = connectAsync(current(primaryRef), "primary");
        CompletableFuture<Boolean> secondaryFuture = connectAsync(current(secondaryRef), "secondary");
        CompletableFuture.allOf(primaryFuture, secondaryFuture).join();

        boolean primaryUp = primaryFuture.getNow(false);
        boolean secondaryUp = secondaryFuture.getNow(false);
        if (!primaryUp && !secondaryUp) {
            throw new PlcConnectionException(
                "Both primary and secondary S7H endpoints failed to connect");
        }
        activeSlot.set(primaryUp ? primaryRef : secondaryRef);
        LOGGER.info("S7H connection up — primary={}, secondary={}, active={}",
            primaryUp ? "OK" : "DOWN", secondaryUp ? "OK" : "DOWN", primaryUp ? "primary" : "secondary");

        heartbeat.scheduleAtFixedRate(this::supervise, heartbeatIntervalMs, heartbeatIntervalMs, TimeUnit.MILLISECONDS);
    }

    @Override
    public void close() throws Exception {
        heartbeat.shutdownNow();
        Exception primaryEx = closeQuietly(current(primaryRef), "primary");
        Exception secondaryEx = closeQuietly(current(secondaryRef), "secondary");
        super.close();
        if (primaryEx != null) throw primaryEx;
        if (secondaryEx != null) throw secondaryEx;
    }

    // ------------------------------------------------------------------------------------
    // Operation delegation. Each onX hook routes to the active inner; on transport-level
    // failure we swap and retry once.
    // ------------------------------------------------------------------------------------

    @Override
    protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
        return withFailover(c -> c.read(readRequest));
    }

    @Override
    protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
        return withFailover(c -> c.write(writeRequest));
    }

    @Override
    protected CompletableFuture<PlcPingResponse> onPing(PlcPingRequest pingRequest) {
        return withFailover(c -> c.ping().thenApply(r -> (PlcPingResponse) r));
    }

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowse(PlcBrowseRequest browseRequest) {
        return withFailover(c -> c.browse(browseRequest));
    }

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowseWithInterceptor(PlcBrowseRequest browseRequest, PlcBrowseRequestInterceptor interceptor) {
        return withFailover(c -> c.browseWithInterceptor(browseRequest, interceptor));
    }

    @Override
    protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest subscriptionRequest) {
        // Refuse subscriptions outright on the dual-path connection. PLC-side subscription
        // state is keyed to a TCP connection, so a transparent fail-over would silently
        // drop alarm pushes and cyclic streams — better to surface the limitation now via
        // an UNSUPPORTED response than to let the app discover it after a real failover.
        // Apps that want subscriptions in a dual-path setup should connect to one host
        // directly, listen for connection-lost events, and re-subscribe in app code.
        return CompletableFuture.failedFuture(new PlcRuntimeException(
            "Subscribe not supported on S7H dual-path connections — see S7HCotpConnection "
            + "javadoc. Connect to a single host instead."));
    }

    @Override
    protected CompletableFuture<PlcUnsubscriptionResponse> onUnsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        return CompletableFuture.failedFuture(new PlcRuntimeException(
            "Unsubscribe not supported on S7H dual-path connections"));
    }

    @Override
    protected PlcConsumerRegistration onRegisterConsumer(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        throw new PlcRuntimeException("Register consumer not supported on S7H dual-path connections");
    }

    @Override
    protected void onUnregisterConsumer(PlcConsumerRegistration registration) {
        // No-op: nothing was registered. Throwing here would be hostile if app code calls
        // unregister defensively in cleanup paths.
    }

    // ------------------------------------------------------------------------------------
    // Internals
    // ------------------------------------------------------------------------------------

    /** Read the inner currently held in a slot. Slots get replaced when an inner is rebuilt. */
    private static S7CotpConnection current(AtomicReference<S7CotpConnection> slot) {
        return slot.get();
    }

    /** Pick a slot to try next: prefer active if it's not stale; else fall back to alt. */
    @SuppressWarnings("resource")
    private AtomicReference<S7CotpConnection> chooseSlot() {
        AtomicReference<S7CotpConnection> a = activeSlot.get();
        if (a != null && !stale.contains(a) && current(a).isConnected()) {
            return a;
        }
        AtomicReference<S7CotpConnection> alt = (a == primaryRef) ? secondaryRef : primaryRef;
        if (alt != null && !stale.contains(alt) && current(alt).isConnected()) {
            activeSlot.set(alt);
            return alt;
        }
        // Both stale or disconnected — last resort: use whichever still claims connected.
        if (a != null && current(a).isConnected()) return a;
        if (alt != null && current(alt).isConnected()) return alt;
        return null;
    }

    /**
     * Invoke {@code op} on the active inner with a wrapper-level timeout
     * ({@link #failoverTimeoutMs}). If the inner doesn't respond in that window — or
     * fails outright — mark its slot stale and retry once on the other slot. The stale
     * flag gets cleared by the heartbeat task once it successfully validates (ping or
     * full rebuild), so a transient drop heals without staying out of rotation forever.
     */
    private <R> CompletableFuture<R> withFailover(Function<S7CotpConnection, CompletableFuture<R>> op) {
        AtomicReference<S7CotpConnection> targetSlot = chooseSlot();
        if (targetSlot == null) {
            return CompletableFuture.failedFuture(
                new PlcRuntimeException("Neither primary nor secondary S7H endpoint is reachable"));
        }
        S7CotpConnection target = current(targetSlot);
        return tryWithTimeout(target, op).handle((r, t) -> {
            if (t == null) {
                return CompletableFuture.completedFuture(r);
            }
            stale.add(targetSlot);
            AtomicReference<S7CotpConnection> altSlot = (targetSlot == primaryRef) ? secondaryRef : primaryRef;
            S7CotpConnection alt = current(altSlot);
            if (!stale.contains(altSlot) && alt.isConnected()) {
                LOGGER.info("S7H operation failed on {} ({}), retrying on {}",
                    label(targetSlot), describe(t), label(altSlot));
                activeSlot.set(altSlot);
                return tryWithTimeout(alt, op);
            }
            LOGGER.info("S7H operation failed on {} and no healthy alternate available: {}",
                label(targetSlot), describe(t));
            return CompletableFuture.<R>failedFuture(t);
        }).thenCompose(Function.identity());
    }

    private <R> CompletableFuture<R> tryWithTimeout(S7CotpConnection inner,
                                                    Function<S7CotpConnection, CompletableFuture<R>> op) {
        return op.apply(inner).orTimeout(failoverTimeoutMs, TimeUnit.MILLISECONDS);
    }

    private String label(AtomicReference<S7CotpConnection> slot) {
        return slot == primaryRef ? "primary" : "secondary";
    }

    /**
     * Pretty-print a Throwable for log lines: prefer the message; fall back to the simple
     * class name when {@link Throwable#getMessage()} is null (e.g. {@link
     * java.util.concurrent.TimeoutException} routinely is).
     */
    private static String describe(Throwable t) {
        if (t == null) return "null";
        String msg = t.getMessage();
        return msg != null ? msg : t.getClass().getSimpleName();
    }

    /**
     * Heartbeat tick — for each slot, validate its inner by pinging the PLC; if the ping
     * fails or the underlying transport has died, <em>rebuild</em> the inner from its
     * factory and reconnect. Successfully-validated inners have their stale flag cleared
     * so {@link #withFailover} considers them again. Finally, if primary is healthy and
     * we're not on it, switch active back.
     *
     * <p>Rebuild-on-failure is necessary because the SPI3 transport doesn't reopen after
     * the kernel closes the socket — the only way back is a fresh transport instance.
     */
    @SuppressWarnings("resource")
    private void supervise() {
        try {
            validate(primaryRef, primaryFactory, "primary");
            validate(secondaryRef, secondaryFactory, "secondary");
            // Prefer primary when it's healthy.
            if (!stale.contains(primaryRef) && current(primaryRef).isConnected()
                && activeSlot.get() != primaryRef) {
                LOGGER.info("S7H heartbeat: primary recovered, switching active back to primary");
                activeSlot.set(primaryRef);
            }
        } catch (Throwable t) {
            LOGGER.debug("S7H heartbeat tick failed", t);
        }
    }

    /**
     * Probe a slot's inner by pinging the PLC and rebuild it if needed. Logs the
     * <em>transitions</em>:
     * <ul>
     *   <li>healthy → broken: one INFO {@code "X disrupted"} line.</li>
     *   <li>broken → healthy: one INFO {@code "X recovered"} line.</li>
     *   <li>still broken / still healthy: silent at INFO (DEBUG details available).</li>
     * </ul>
     * The {@link #stale} set doubles as the "is it currently broken?" predicate, so first
     * insert/remove on each transition gives us clean once-per-edge logging.
     */
    private void validate(AtomicReference<S7CotpConnection> slot,
                          Supplier<S7CotpConnection> factory, String label) {
        S7CotpConnection inner = slot.get();
        String reason;
        boolean nowHealthy;
        if (!inner.isConnected()) {
            reason = "transport dead";
            nowHealthy = rebuild(slot, factory, label);
        } else {
            try {
                inner.ping().get(failoverTimeoutMs, TimeUnit.MILLISECONDS);
                reason = null;
                nowHealthy = true;
            } catch (Exception e) {
                reason = "ping failed: " + describe(e);
                nowHealthy = rebuild(slot, factory, label);
            }
        }

        if (nowHealthy) {
            if (stale.remove(slot)) {
                LOGGER.info("S7H heartbeat: {} recovered", label);
            }
        } else if (stale.add(slot)) {
            LOGGER.info("S7H heartbeat: {} disrupted ({}), rebuild attempt failed", label, reason);
        } else {
            LOGGER.debug("S7H heartbeat: {} still down ({}), rebuild attempt failed", label, reason);
        }
    }

    /**
     * Replace the inner held in {@code slot} with a fresh one built by {@code factory},
     * closing the old one first. Returns {@code true} iff the new inner connected.
     */
    private boolean rebuild(AtomicReference<S7CotpConnection> slot,
                            Supplier<S7CotpConnection> factory, String label) {
        S7CotpConnection old = slot.get();
        try { old.close(); } catch (Exception ignored) { /* best-effort */ }
        S7CotpConnection fresh;
        try {
            fresh = factory.get();
        } catch (Exception e) {
            LOGGER.debug("S7H heartbeat: factory for {} threw: {}", label, e.getMessage());
            return false;
        }
        try {
            fresh.connect();
        } catch (Exception e) {
            LOGGER.debug("S7H heartbeat: reconnect of fresh {} failed: {}", label, e.getMessage());
            return false;
        }
        slot.set(fresh);
        return true;
    }

    private static CompletableFuture<Boolean> connectAsync(S7CotpConnection inner, String label) {
        return CompletableFuture.supplyAsync(() -> {
            try {
                inner.connect();
                return true;
            } catch (Exception e) {
                LOGGER.debug("S7H {} connect failed: {}", label, e.getMessage());
                return false;
            }
        });
    }

    private static Exception closeQuietly(S7CotpConnection inner, String label) {
        try {
            inner.close();
            return null;
        } catch (Exception e) {
            LOGGER.debug("S7H close on {} failed: {}", label, e.getMessage());
            return e;
        }
    }
}
