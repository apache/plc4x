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
package org.apache.plc4x.java.bacnetip;

import org.apache.commons.io.FileUtils;
import org.apache.plc4x.java.bacnetip.readwrite.BVLC;
import org.apache.plc4x.test.generator.ParserSerializerTestsuiteGenerator;

import java.io.File;
import java.io.IOException;
import java.net.URL;
import java.nio.file.FileSystems;
import java.nio.file.Path;

public class BACnetParserSerializerTestSuiteGenerator {

    public static void main(String... args) throws Exception {
        String pcapFile = DownloadAndCache("bacnet-stack-services.cap");
        String xmlTestSuiteFile = Path.of("protocols/bacnetip/src/test/resources/protocols/bacnet/ParserSerializerTestsuite.xml").toAbsolutePath().toString();
        ParserSerializerTestsuiteGenerator.main("-tBACnet/IP", "-pbacnetip", BVLC.class.getName(), pcapFile, xmlTestSuiteFile);
    }

    private static String DownloadAndCache(String file) throws IOException {
        String tempDirectory = FileUtils.getTempDirectoryPath();
        File pcapFile = FileSystems.getDefault().getPath(tempDirectory, RandomPackagesTest.class.getSimpleName(), file).toFile();
        FileUtils.createParentDirectories(pcapFile);
        if (!pcapFile.exists()) {
            URL source = new URL("https://kargs.net/captures/" + file);
            FileUtils.copyURLToFile(source, pcapFile);
        }
        return pcapFile.getAbsolutePath();
    }
}
