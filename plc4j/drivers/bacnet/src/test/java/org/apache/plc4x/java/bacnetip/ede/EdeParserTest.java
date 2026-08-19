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
package org.apache.plc4x.java.bacnetip.ede;

import com.sun.management.UnixOperatingSystemMXBean;
import org.apache.plc4x.java.bacnetip.ede.model.Datapoint;
import org.apache.plc4x.java.bacnetip.ede.model.EdeModel;
import org.apache.plc4x.java.bacnetip.tag.BacNetIpTag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.io.File;
import java.lang.management.ManagementFactory;
import java.lang.management.OperatingSystemMXBean;
import java.nio.file.Files;
import java.nio.file.Path;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.junit.jupiter.api.Assumptions.assumeTrue;

/**
 * The parser used to hand its {@link java.io.FileReader} to commons-csv and never close either one,
 * so every parsed EDE file leaked a file handle for the lifetime of the JVM. On Windows the leaked
 * handle also keeps the file locked, which is how it first surfaced.
 */
class EdeParserTest {

    @TempDir
    Path tempDir;

    @Test
    void parsesADatapoint() throws Exception {
        File edeFile = writeEdeFile("edeDataText.csv");

        EdeModel model = new EdeParser().parseFile(edeFile);

        Datapoint datapoint = model.getDatapoint(new BacNetIpTag(12L, 2, 7));
        assertNotNull(datapoint, "the single data row is picked up");
        assertEquals("Temperature", datapoint.getObjectName());
    }

    /**
     * Windows refuses to delete a file while a handle on it is open, so this fails there if the
     * reader is left open. On POSIX the delete succeeds either way - {@link #doesNotLeakFileHandles}
     * is what catches the leak on those platforms.
     */
    @Test
    void releasesTheFileAfterParsing() throws Exception {
        File edeFile = writeEdeFile("edeDataText.csv");

        new EdeParser().parseFile(edeFile);

        Files.delete(edeFile.toPath());
        assertFalse(edeFile.exists());
    }

    /**
     * The delete check above is a no-op on POSIX, so count the descriptors instead.
     */
    @Test
    void doesNotLeakFileHandles() throws Exception {
        OperatingSystemMXBean osBean = ManagementFactory.getOperatingSystemMXBean();
        assumeTrue(osBean instanceof UnixOperatingSystemMXBean, "open descriptors are not countable here");
        UnixOperatingSystemMXBean unixBean = (UnixOperatingSystemMXBean) osBean;

        File edeFile = writeEdeFile("edeDataText.csv");
        EdeParser parser = new EdeParser();
        // One parse up front so that anything opened lazily on the first pass is not counted.
        parser.parseFile(edeFile);

        long before = unixBean.getOpenFileDescriptorCount();
        for (int i = 0; i < 100; i++) {
            parser.parseFile(edeFile);
        }
        long leaked = unixBean.getOpenFileDescriptorCount() - before;

        // Unrelated activity can open a descriptor or two, but not a hundred of them.
        assertTrue(leaked < 10, "leaked " + leaked + " file descriptors over 100 parses");
    }

    /**
     * Five rows of preamble, the layout version in the second column, one more row to skip and the
     * column names - see {@link EdeParser} and the version 2 layout.
     */
    private File writeEdeFile(String name) throws Exception {
        Path file = tempDir.resolve(name);
        Files.writeString(file, String.join("\n",
            "#Data Exchange File",
            "PROJECT_NAME;test",
            "VERSION_OF_REFERENCEFILE;1",
            "TIMESTAMP_OF_LAST_CHANGE;2026-01-01",
            "AUTHOR_OF_LAST_CHANGE;plc4x",
            "VERSION_OF_LAYOUT;2",
            "#mandatory;mandatory;mandatory;mandatory;mandatory",
            "keyname;device obj.-instance;object-name;object-type;object-instance;description;"
                + "default value;min. value;max. value;commandable;supports COV;hi limit;lo limit;"
                + "state-text-reference;unit-code;vendor-specific-address",
            "temp-1;12;Temperature;2;7;Outside temperature;0;-50;150;N;Y;100;-40;;62;"
        ) + "\n");
        return file.toFile();
    }
}
