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

import org.apache.commons.compress.archivers.zip.ZipArchiveEntry;
import org.apache.commons.compress.archivers.zip.ZipArchiveInputStream;
import org.apache.commons.io.FileUtils;
import org.apache.http.HttpEntity;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.util.EntityUtils;
import org.apache.plc4x.protocol.knxnetip.handlers.ManufacturerIdsHandler;
import org.json.JSONArray;
import org.json.JSONObject;
import org.xml.sax.SAXException;

import javax.xml.parsers.ParserConfigurationException;
import javax.xml.parsers.SAXParser;
import javax.xml.parsers.SAXParserFactory;
import java.io.File;
import java.io.IOException;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;

public class KnxSpecificationStoreUpdater {

    /**
     * Description of the KNX Foundation Webservice: http://onlinecatalog.knx.org/Download/help
     * @param args
     */
    public static void main(String[] args) {
        List<Integer> manufacturerIds;
        Map<String, Integer> comObjectTableStartAddresses = new TreeMap<>();

        File contentDir = new File("content");
        if (!contentDir.exists()) {
            contentDir.mkdirs();
        }
        try (CloseableHttpClient httpClient = HttpClients.createDefault()) {
            // Get an index of all manufacturers that have their product data online
            HttpGet getManufacturersRequest = new HttpGet("http://onlinecatalog.knx.org/Download/Manufacturers");
            try (CloseableHttpResponse response = httpClient.execute(getManufacturersRequest)) {
                HttpEntity entity = response.getEntity();

                // Copy to a local file.
                File manufacturersFile = new File(contentDir, "manufacturers.xml");
                FileUtils.copyInputStreamToFile(entity.getContent(), manufacturersFile);

                // Parse the content.
                SAXParserFactory factory = SAXParserFactory.newInstance();
                SAXParser saxParser = factory.newSAXParser();
                ManufacturerIdsHandler manufacturerIdsHandler = new ManufacturerIdsHandler();
                saxParser.parse(manufacturersFile, manufacturerIdsHandler);

                manufacturerIds = manufacturerIdsHandler.getManufacturerIds();
            } catch (SAXException | ParserConfigurationException e) {
                e.printStackTrace();
                return;
            }

            // For each manufacturer id, get the catalog / product index
            // This is a JSON file (Yeah ... I know ... why stick to one
            // format?)
            for (Integer manufacturerId : manufacturerIds) {
                File manufacturerDirectory = new File(contentDir, String.format("M-%04X", manufacturerId));
                if (!manufacturerDirectory.exists()) {
                    manufacturerDirectory.mkdirs();
                }
                boolean indexNeedsUpdate = true;
                // If the file exists, check if there's a newer version.
                File indexFile = new File(manufacturerDirectory, "catalog.json");
                JSONObject jsonObject = null;
                if (indexFile.exists()) {
                    String fileContent = FileUtils.readFileToString(indexFile, "UTF-8");
                    jsonObject = new JSONObject(fileContent);
                    int version = jsonObject.getInt("Version");
                    HttpGet checkNewerVersionRequest = new HttpGet(String.format("http://onlinecatalog.knx.org/Download/HasNewerIndexFile/%d/%d", manufacturerId, version));
                    try (CloseableHttpResponse response = httpClient.execute(checkNewerVersionRequest)) {
                        // This is an XML response just containing either:
                        // <boolean xmlns="http://schemas.microsoft.com/2003/10/Serialization/">false</boolean>
                        // or:
                        // <boolean xmlns="http://schemas.microsoft.com/2003/10/Serialization/">true</boolean>
                        if (EntityUtils.toString(response.getEntity()).contains("false")) {
                            indexNeedsUpdate = false;
                        }
                    }
                }
                if (indexNeedsUpdate) {
                    HttpGet getProductIndexRequest = new HttpGet("http://onlinecatalog.knx.org/Download/Index/" + manufacturerId);
                    try (CloseableHttpResponse response = httpClient.execute(getProductIndexRequest)) {
                        if (response.getStatusLine().getStatusCode() != 200) {
                            System.out.println("Got an unexpected status code " + response.getStatusLine().getStatusCode());
                        }
                        HttpEntity entity = response.getEntity();

                        // Save the content locally.
                        FileUtils.writeStringToFile(indexFile, EntityUtils.toString(entity).substring(3), "UTF-8");

                        String fileContent = FileUtils.readFileToString(indexFile, "UTF-8");
                        jsonObject = new JSONObject(fileContent);
                    }
                }

                // Iterate over all the products in the catalog and fetch each one.
                // Each request will result in a little ZIP file, that contains some
                // boilerplate files, but also an XML file with the product descriptor.
                JSONArray entries = jsonObject.getJSONArray("Entries");
                for (Object entry : entries) {
                    if (!(entry instanceof JSONObject)) {
                        System.out.println("Unexpected entry type");
                        continue;
                    }
                    String productId = ((JSONObject) entry).getString("Id");
                    JSONArray applicationIdentifier = ((JSONObject) entry).getJSONArray("ApplicationIdentifier");
                    int applicationId = applicationIdentifier.getInt(2) << 8 | applicationIdentifier.getInt(3);
                    int applicationVersion = applicationIdentifier.getInt(4) & 0xFF;

                    if(applicationId == 0 && applicationVersion == 0) {
                        System.out.println("SKIPPED");
                        continue;
                    }

                    String productCode = String.format("M-%04X_A-%04X-%02X", manufacturerId, applicationId, applicationVersion);

                    System.out.print("Fetching product: " + productCode + ": ");

                    // Check If we've already got that file (There are no updates, just new versions)
                    File[] files = manufacturerDirectory.listFiles((dir, name) -> name.startsWith(productCode));
                    // If we've already got the file, skip loading it
                    if(files.length > 0) {
                        System.out.println("SKIPPED");
                        continue;
                    }

                    HttpPost downloadProduct = new HttpPost("http://onlinecatalog.knx.org/Download/DownloadProduct");
                    downloadProduct.setHeader("Content-type", "application/json");
                    downloadProduct.setEntity(new StringEntity(String.format("{\"CatalogIds\":[\"%s\"],\"LanguageIds\":[\"en-US\"]}", productId)));
                    try (CloseableHttpResponse downloadProductResponse = httpClient.execute(downloadProduct)) {
                        if (downloadProductResponse.getStatusLine().getStatusCode() != 200) {
                            if (EntityUtils.toString(downloadProductResponse.getEntity()).contains("NotInAnyMarket")) {
                                // This product id is not available in the online product list
                                // Create an empty file indicating it's not in the catalog so we won't try to fetch it again.
                                File dummy = new File(manufacturerDirectory, productCode + ".failed");
                                dummy.createNewFile();
                                System.out.println("FAILED");
                            } else {
                                System.out.println("Got an unexpected status code " + downloadProductResponse.getStatusLine().getStatusCode());
                            }
                        } else {
                            String expectedPrefix = String.format("M-%04X/%s", manufacturerId, productCode);
                            try (ZipArchiveInputStream zis = new ZipArchiveInputStream(downloadProductResponse.getEntity().getContent())) {
                                ZipArchiveEntry zipFileEntry;
                                while ((zipFileEntry = zis.getNextZipEntry()) != null) {
                                    String fileName = zipFileEntry.getName();
                                    if (fileName.startsWith(expectedPrefix)) {
                                        File productFile = new File(manufacturerDirectory, fileName.substring(fileName.indexOf('/') + 1));
                                        FileUtils.copyInputStreamToFile(zis, productFile);
                                        System.out.println("SUCCESS");
                                        break;
                                    }
                                }
                            }
                        }
                    }
                }
            }
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

}
