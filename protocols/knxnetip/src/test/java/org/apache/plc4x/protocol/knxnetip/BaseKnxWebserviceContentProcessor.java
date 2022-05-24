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
import org.json.JSONArray;
import org.json.JSONObject;

import javax.xml.parsers.SAXParser;
import javax.xml.parsers.SAXParserFactory;
import java.io.File;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

/**
 * Base helper that crawls over all xml descriptions found in the local storage.
 */
public abstract class BaseKnxWebserviceContentProcessor {

    /**
     * Logic to crawl the local knx profile storage for further processing.
     * @param contentDir points to the directory maintained by the {@link KnxSpecificationStoreUpdater}
     * @throws Exception something went wrong
     */
    public void processDirectory(File contentDir) throws Exception {
        Set<String> processedIds = new HashSet<>();
        List<Integer> manufacturerIds;

        if (!contentDir.exists()) {
            throw new RuntimeException("The input directory doesn't exist.");
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
                if ((files != null) && (files.length > 0)) {
                    for (File file : files) {
                        processFile(file, manufacturerId, applicationId, applicationVersion);
                    }
                }
            }
        }
    }

    public abstract void processFile(File file, int manufacturerId, int applicationId, int applicationVersion) throws Exception;
}
