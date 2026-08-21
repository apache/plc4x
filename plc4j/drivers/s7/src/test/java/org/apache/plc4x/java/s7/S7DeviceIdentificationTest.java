/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.s7.S7ScriptedConnectionHarness.ScriptedS7Transport;
import org.apache.plc4x.java.s7.userdata.S7SzlService.S7DeviceIdentification;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * A device inventory needs more than the article number the connect-time probe keeps, so the
 * connection can read the three identification SZLs on demand. Devices answer different subsets
 * of them — notably S7-1200/1500 usually refuse the protection-status list — so whatever did
 * come back has to survive the ones that didn't.
 */
class S7DeviceIdentificationTest {

    @Test
    void readsOrderCodeAndFirmwareFromTheModuleIdentificationList() throws Exception {
        S7DeviceIdentification id = identificationFrom(new SzlScript(true, true, true));
        assertEquals("6ES7 511-1AK02-0AB0", id.orderCode());
        assertEquals("V 2.6.0", id.firmwareVersion());
    }

    @Test
    void readsNamesFromTheComponentIdentificationList() throws Exception {
        S7DeviceIdentification id = identificationFrom(new SzlScript(true, true, true));
        assertEquals("PLC_1", id.moduleName());
        assertEquals("S C-J2UD57132018", id.serialNumber());
    }

    @Test
    void readsProtectionLevelFromTheProtectionStatusList() throws Exception {
        S7DeviceIdentification id = identificationFrom(new SzlScript(true, true, true));
        assertEquals(Integer.valueOf(2), id.protectionLevel());
        assertEquals(Boolean.TRUE, id.passwordProtected());
    }

    @Test
    void keepsIdentificationWhenTheProtectionListIsRefused() throws Exception {
        // What an S7-1200/1500 does: answers 0x0011 and 0x001C, rejects 0x0232.
        S7DeviceIdentification id = identificationFrom(new SzlScript(true, true, false));
        assertEquals("V 2.6.0", id.firmwareVersion());
        assertEquals("PLC_1", id.moduleName());
        assertNull(id.protectionLevel());
    }

    @Test
    void keepsFirmwareWhenTheComponentListIsRefused() throws Exception {
        S7DeviceIdentification id = identificationFrom(new SzlScript(true, false, true));
        assertEquals("V 2.6.0", id.firmwareVersion());
        assertNull(id.moduleName());
    }

    @Test
    void returnsAnEmptyIdentificationWhenEverySzlIsRefused() throws Exception {
        S7DeviceIdentification id = identificationFrom(new SzlScript(false, false, false));
        assertEquals(S7DeviceIdentification.EMPTY, id);
    }

    /**
     * A list that doesn't fit in one PDU comes back as a chain: the first chunk sets
     * lastDataUnit and can stop in the middle of a record, and the continuation carries no
     * header of its own. Unless the driver echoes the sequence number and appends the chunks,
     * the straddling record is lost - and the CPU refuses every later SZL read on the
     * connection because the sequence was left open.
     */
    @Test
    void assemblesAChainedSzlResponse() throws Exception {
        ScriptedS7Transport transport = new ScriptedS7Transport();
        S7CotpConnection connection = S7ScriptedConnectionHarness.newConnectedConnection(transport);
        transport.resetCounters();

        byte[] whole = moduleIdentificationData();
        int cut = 8 + 14;   // stop halfway through the first record
        byte[] firstChunk = Arrays.copyOfRange(whole, 0, cut);
        byte[] rest = Arrays.copyOfRange(whole, cut, whole.length);

        CompletableFuture<S7DeviceIdentification> future = connection.readDeviceIdentification();
        answerWith(transport, future, request -> {
            int tpduReference = S7ScriptedConnectionHarness.tpduReferenceOfFrame(request);
            if (S7ScriptedConnectionHarness.szlIdOfRequest(request) != 0x0011) {
                return S7ScriptedConnectionHarness.szlErrorResponse(tpduReference, 0xD401);
            }
            // A continuation is a ReadSZL request echoing the sequence number we handed out.
            return S7ScriptedConnectionHarness.sequenceNumberOfRequest(request) == SEQUENCE
                ? S7ScriptedConnectionHarness.szlResponse(tpduReference, rest, SEQUENCE, 0)
                : S7ScriptedConnectionHarness.szlResponse(tpduReference, firstChunk, SEQUENCE, 1);
        });

        S7DeviceIdentification id = future.get(20, TimeUnit.SECONDS);
        connection.close();
        assertEquals("V 2.6.0", id.firmwareVersion());
        assertEquals("6ES7 511-1AK02-0AB0", id.orderCode());
    }

    /**
     * A device that keeps claiming there is more data must not spin forever - the driver caps
     * the chain and returns what it has.
     */
    @Test
    void stopsChainingAfterTheContinuationLimit() throws Exception {
        ScriptedS7Transport transport = new ScriptedS7Transport();
        S7CotpConnection connection = S7ScriptedConnectionHarness.newConnectedConnection(transport);
        transport.resetCounters();

        CompletableFuture<S7DeviceIdentification> future = connection.readDeviceIdentification();
        answerWith(transport, future, request -> {
            int tpduReference = S7ScriptedConnectionHarness.tpduReferenceOfFrame(request);
            if (S7ScriptedConnectionHarness.szlIdOfRequest(request) != 0x0011) {
                return S7ScriptedConnectionHarness.szlErrorResponse(tpduReference, 0xD401);
            }
            // Always "more data follows".
            return S7ScriptedConnectionHarness.szlResponse(
                tpduReference, moduleIdentificationData(), SEQUENCE, 1);
        });

        S7DeviceIdentification id = future.get(20, TimeUnit.SECONDS);
        connection.close();
        // It gave up, but kept what the chunks contained rather than failing outright.
        assertEquals("6ES7 511-1AK02-0AB0", id.orderCode());
        assertTrue(transport.writeCount() <= 20,
            "expected the chain to be capped, but the driver sent " + transport.writeCount()
                + " requests");
    }

    /** Sequence number the scripted device hands out for its chained reads. */
    private static final int SEQUENCE = 3;

    /** Answers every frame the driver writes until {@code future} completes. */
    private static void answerWith(ScriptedS7Transport transport, CompletableFuture<?> future,
                                   Responder responder) {
        Thread thread = new Thread(() -> {
            int answered = 0;
            while (!future.isDone()) {
                try {
                    if (transport.writeCount() > answered) {
                        byte[] request = transport.writtenFrames().get(answered);
                        answered++;
                        transport.deliver(responder.replyTo(request));
                        transport.runDataListener();
                    }
                    Thread.sleep(5);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    return;
                } catch (Exception e) {
                    throw new IllegalStateException("could not answer SZL request", e);
                }
            }
        }, "szl-responder");
        thread.setDaemon(true);
        thread.start();
    }

    private interface Responder {
        byte[] replyTo(byte[] request) throws Exception;
    }

    /**
     * Runs {@code readDeviceIdentification()} against the scripted transport, answering each
     * SZL request the driver sends according to {@code script}.
     */
    private static S7DeviceIdentification identificationFrom(SzlScript script) throws Exception {
        ScriptedS7Transport transport = new ScriptedS7Transport();
        S7CotpConnection connection = S7ScriptedConnectionHarness.newConnectedConnection(transport);
        transport.resetCounters();

        CompletableFuture<S7DeviceIdentification> future = connection.readDeviceIdentification();

        Thread responder = new Thread(() -> {
            int answered = 0;
            while (!future.isDone()) {
                try {
                    if (transport.writeCount() > answered) {
                        byte[] request = transport.writtenFrames().get(answered);
                        answered++;
                        transport.deliver(script.replyTo(request));
                        transport.runDataListener();
                    }
                    Thread.sleep(5);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    return;
                } catch (Exception e) {
                    throw new IllegalStateException("could not answer SZL request", e);
                }
            }
        }, "szl-responder");
        responder.setDaemon(true);
        responder.start();

        S7DeviceIdentification id = future.get(20, TimeUnit.SECONDS);
        connection.close();
        return id;
    }

    /** Decides, per SZL-ID, whether the simulated device answers or rejects the request. */
    private record SzlScript(boolean answerModule, boolean answerComponent, boolean answerProtection) {

        byte[] replyTo(byte[] request) throws Exception {
            int szlId = S7ScriptedConnectionHarness.szlIdOfRequest(request);
            int tpduReference = S7ScriptedConnectionHarness.tpduReferenceOfFrame(request);
            return switch (szlId) {
                case 0x0011 -> answerModule
                    ? S7ScriptedConnectionHarness.szlResponse(tpduReference, moduleIdentificationData())
                    : S7ScriptedConnectionHarness.szlErrorResponse(tpduReference, 0xD401);
                case 0x001C -> answerComponent
                    ? S7ScriptedConnectionHarness.szlResponse(tpduReference, componentIdentificationData())
                    : S7ScriptedConnectionHarness.szlErrorResponse(tpduReference, 0xD401);
                case 0x0232 -> answerProtection
                    ? S7ScriptedConnectionHarness.szlResponse(tpduReference, protectionStatusData())
                    : S7ScriptedConnectionHarness.szlErrorResponse(tpduReference, 0xD401);
                default -> throw new IllegalStateException(
                    "unexpected SZL-ID 0x" + Integer.toHexString(szlId));
            };
        }
    }

    private static byte[] moduleIdentificationData() {
        return szlData(0x0011, 0x0000, 28,
            moduleRecord(1, "6ES7 511-1AK02-0AB0", 5, 1, 2),
            moduleRecord(6, "6ES7 511-1AK02-0AB0", 1, 0, 0),
            moduleRecord(7, "6ES7 511-1AK02-0AB0", 2, 6, 0));
    }

    private static byte[] componentIdentificationData() {
        return szlData(0x001C, 0x0000, 34,
            componentRecord(1, "PLANT-LINE-3"),
            componentRecord(2, "PLC_1"),
            componentRecord(5, "S C-J2UD57132018"));
    }

    private static byte[] protectionStatusData() {
        // sch_schal=0, sch_par=2 (password configured), sch_rel=2, bart_sch=2 (RUN-P), anl_sch=0
        return szlData(0x0232, 0x0004, 12, words(0x0004, 0, 2, 2, 2, 0));
    }

    private static byte[] moduleRecord(int index, String mlfb, int v1, int v2, int v3) {
        byte[] record = new byte[28];
        record[0] = (byte) (index >> 8);
        record[1] = (byte) index;
        byte[] mlfbBytes = mlfb.getBytes(StandardCharsets.US_ASCII);
        System.arraycopy(mlfbBytes, 0, record, 2, Math.min(mlfbBytes.length, 20));
        record[25] = (byte) v1;
        record[26] = (byte) v2;
        record[27] = (byte) v3;
        return record;
    }

    private static byte[] componentRecord(int index, String value) {
        byte[] record = new byte[34];
        record[0] = (byte) (index >> 8);
        record[1] = (byte) index;
        byte[] valueBytes = value.getBytes(StandardCharsets.US_ASCII);
        System.arraycopy(valueBytes, 0, record, 2, Math.min(valueBytes.length, 32));
        return record;
    }

    private static byte[] szlData(int szlId, int szlIndex, int recordLength, byte[]... records) {
        byte[] data = new byte[8 + records.length * recordLength];
        byte[] header = words(szlId, szlIndex, recordLength, records.length);
        System.arraycopy(header, 0, data, 0, 8);
        for (int i = 0; i < records.length; i++) {
            System.arraycopy(records[i], 0, data, 8 + i * recordLength,
                Math.min(records[i].length, recordLength));
        }
        return data;
    }

    private static byte[] words(int... values) {
        byte[] out = new byte[values.length * 2];
        for (int i = 0; i < values.length; i++) {
            out[i * 2] = (byte) (values[i] >> 8);
            out[i * 2 + 1] = (byte) values[i];
        }
        return out;
    }
}
