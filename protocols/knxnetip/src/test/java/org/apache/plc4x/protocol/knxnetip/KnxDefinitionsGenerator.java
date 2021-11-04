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
package org.apache.plc4x.protocol.knxnetip;

import org.apache.commons.io.FileUtils;
import org.apache.plc4x.protocol.knxnetip.handlers.ManufacturerIdsHandler;
import org.apache.plc4x.protocol.knxnetip.handlers.ProductDescriptionHandler;
import org.json.JSONArray;
import org.json.JSONObject;
import org.xml.sax.SAXException;

import javax.xml.parsers.ParserConfigurationException;
import javax.xml.parsers.SAXParser;
import javax.xml.parsers.SAXParserFactory;
import java.io.File;
import java.io.IOException;
import java.util.*;

public class KnxDefinitionsGenerator {

    public static final String HEADER = "/*\n" +
        " * Licensed to the Apache Software Foundation (ASF) under one\n" +
        " * or more contributor license agreements.  See the NOTICE file\n" +
        " * distributed with this work for additional information\n" +
        " * regarding copyright ownership.  The ASF licenses this file\n" +
        " * to you under the Apache License, Version 2.0 (the\n" +
        " * \"License\"); you may not use this file except in compliance\n" +
        " * with the License.  You may obtain a copy of the License at\n" +
        " *\n" +
        " *   http://www.apache.org/licenses/LICENSE-2.0\n" +
        " *\n" +
        " * Unless required by applicable law or agreed to in writing,\n" +
        " * software distributed under the License is distributed on an\n" +
        " * \"AS IS\" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY\n" +
        " * KIND, either express or implied.  See the License for the\n" +
        " * specific language governing permissions and limitations\n" +
        " * under the License.\n" +
        " */\n\n" +
        "\n" +
        "[enum uint 16 DeviceInformation(uint 16 'deviceDescriptor', string 'name', uint 16 'comObjectTableAddress')\n";

    public static final String FOOTER = "]\n";

    public static void main(String[] args) throws IOException, SAXException, ParserConfigurationException {
        int deviceCounter = 1;
        Set<String> processedIds = new HashSet<>();
        List<Integer> manufacturerIds;
        Map<String, Integer> comObjectTableStartAddresses = new TreeMap<>();

        File contentDir = new File(args[0]);
        if (!contentDir.exists()) {
            contentDir.mkdirs();
        }
        // Get an index of all manufacturers that have their product data online
        File manufacturersFile = new File(contentDir, "manufacturers.xml");

        // Parse the content.
        SAXParserFactory factory = SAXParserFactory.newInstance();
        SAXParser saxParser = factory.newSAXParser();
        ManufacturerIdsHandler manufacturerIdsHandler = new ManufacturerIdsHandler();
        saxParser.parse(manufacturersFile, manufacturerIdsHandler);

        // Get the final list of manufacturer ids from the handler.
        manufacturerIds = manufacturerIdsHandler.getManufacturerIds();

        // For each manufacturer id, get the catalog / product index
        // This is a JSON file (Yeah ... I know ... why stick to one
        // format?)
        for (Integer manufacturerId : manufacturerIds) {
            File manufacturerDirectory = new File(contentDir, String.format("M-%04X", manufacturerId));
            if (!manufacturerDirectory.exists()) {
                continue;
            }

            // Read the catalog file for this manufacturer
            File indexFile = new File(manufacturerDirectory, "catalog.json");
            String fileContent = FileUtils.readFileToString(indexFile, "UTF-8");
            JSONObject jsonObject = new JSONObject(fileContent);

            // Iterate over all the products in the catalog and fetch each one.
            // Each request will result in a little ZIP file, that contains some
            // boilerplate files, but also an XML file with the product descriptor.
            JSONArray entries = jsonObject.getJSONArray("Entries");
            for (Object entry : entries) {
                if (!(entry instanceof JSONObject)) {
                    System.out.println("Unexpected entry type");
                    continue;
                }
                JSONArray applicationIdentifier = ((JSONObject) entry).getJSONArray("ApplicationIdentifier");
                int applicationId = applicationIdentifier.getInt(2) << 8 | applicationIdentifier.getInt(3);
                int applicationVersion = applicationIdentifier.getInt(4) & 0xFF;

                String productCode = String.format("M-%04X_A-%04X-%02X", manufacturerId, applicationId, applicationVersion);

                // Make sure we only output one id once.
                if (processedIds.contains(productCode)) {
                    continue;
                } else {
                    processedIds.add(productCode);
                }

                // Check If we've already got that file (There are no updates, just new versions)
                File[] files = manufacturerDirectory.listFiles((dir, name) -> name.startsWith(productCode) && name.endsWith(".xml"));

                // If we've found a file, get the information from it.
                if (files.length > 0) {
                    ProductDescriptionHandler handler = new ProductDescriptionHandler();
                    saxParser.parse(files[0], handler);

                    // Unfortunately the enum gets soo big, that we can only generate it for the devices we need it for.
                    // These are the System 7 devices.
                    if(handler.getMaskVersion().startsWith("070")) {
                        //String cleanedName = handler.getName().replaceAll("\n", "").replaceAll("\r", "").replaceAll("\"", "inch");
                        if (handler.getComObjectTableAddress() != null) {
                            System.out.printf("    ['%4d' DEV%04X%04X%02X               ['0x%04X'                       ]]\n",
                                deviceCounter++, manufacturerId, applicationId, applicationVersion,
                                handler.getComObjectTableAddress());
                        }
                    }
                }
            }
        }
    }
}
