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

import org.apache.commons.csv.CSVFormat;
import org.apache.commons.csv.CSVParser;
import org.apache.commons.csv.CSVRecord;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.bacnetip.ede.layouts.EdeLayout;
import org.apache.plc4x.java.bacnetip.ede.layouts.EdeLayoutFactory;
import org.apache.plc4x.java.bacnetip.ede.model.Datapoint;
import org.apache.plc4x.java.bacnetip.ede.model.EdeModel;
import org.apache.plc4x.java.bacnetip.field.BacNetIpField;

import java.io.*;
import java.nio.file.FileVisitResult;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.SimpleFileVisitor;
import java.nio.file.attribute.BasicFileAttributes;
import java.util.*;

public class EdeParser {

    public EdeModel parseDirectory(File edeDirectory) {
        Map<BacNetIpField, Datapoint> datapoints = new HashMap<>();
        List<File> edeFiles = findAllEdeFiles(edeDirectory);
        try {
            for (File edeFile : edeFiles) {
                datapoints.putAll(parseFileDatapoints(edeFile));
            }
        } catch (Exception e) {
            throw new PlcRuntimeException("Error parsing EDE files", e);
        }
        return new EdeModel(datapoints);
    }

    public EdeModel parseFile(File edeFile) {
        if(!edeFile.exists()) {
            throw new PlcRuntimeException("EDE File at " + edeFile.getPath() + " doesn't exist.");
        }
        return new EdeModel(parseFileDatapoints(edeFile));
    }

    private List<File> findAllEdeFiles(File curDir) {
        List<File> edeFiles = new LinkedList<>();
        try {
            Files.walkFileTree(curDir.toPath(), new SimpleFileVisitor<Path>() {
                @Override
                public FileVisitResult visitFile(Path path, BasicFileAttributes attrs) {
                    final File file = path.toFile();
                    if(file.isFile()) {
                        // If the name starts with "edeDataText" this is probably an EDE file.
                        if (file.getName().startsWith("edeDataText")) {
                            String suffix = file.getName().substring("edeDataText".length());
                            // If there's a matching "edeStateText" file, we've found a match.
                            File stateFile = new File(file.getParentFile(), "edeStateText" + suffix);
                            if (stateFile.exists() && stateFile.isFile()) {
                                edeFiles.add(path.toFile());
                            }
                        }
                    } else if (file.isDirectory()) {
                        edeFiles.addAll(findAllEdeFiles(file));
                    }
                    return FileVisitResult.CONTINUE;
                }
            });
        } catch (IOException e) {
            throw new PlcRuntimeException("Error scanning EDE directories", e);
        }
        return edeFiles;
    }

    private Map<BacNetIpField, Datapoint> parseFileDatapoints(File edeFile) {
        try {
            Reader in = new FileReader(edeFile);
            final CSVParser parser = CSVFormat.newFormat(';').parse(in);
            final Iterator<CSVRecord> iterator = parser.iterator();

            // Skip the header.
            for(int i = 0; i < 5; i++) {
                if(!iterator.hasNext()) {
                    throw new PlcRuntimeException("Invalid EDE file format");
                }
                iterator.next();
            }
            final CSVRecord layoutVersionRow = iterator.next();

            EdeLayout edeLayout = EdeLayoutFactory.getLayoutForVersion(safeCastInteger(layoutVersionRow, 1));
            if(edeLayout == null) {
                throw new PlcRuntimeException("rted EDE file layout version " + layoutVersionRow.get(1));
            }
            // Just skip the next row.
            iterator.next();
            // Get the column names.
            final CSVRecord columnNames = iterator.next();

            Map<BacNetIpField, Datapoint> datapoints = new HashMap<>();
            // Process the content.
            iterator.forEachRemaining(record -> {
                Long deviceInstance = safeCastLong(record, edeLayout.getDeviceInstancePos());
                Integer objectType = safeCastInteger(record, edeLayout.getObjectTypePos());
                Long objectInstance = safeCastLong(record, edeLayout.getObjectInstancePos());
                if((deviceInstance == null) || (objectType == null) || (objectInstance == null)) {
                    return;
                }
                BacNetIpField address = new BacNetIpField(deviceInstance, objectType, objectInstance);
                String keyName = safeString(record, edeLayout.getKeyNamePos());
                String objectName = safeString(record, edeLayout.getObjectNamePos());
                String description = safeString(record, edeLayout.getDescriptionPos());
                Double defaultValue = safeCastDouble(record, edeLayout.getDefaultValuePos());
                Double minValue = safeCastDouble(record, edeLayout.getMinValuePos());
                Double maxValue = safeCastDouble(record, edeLayout.getMaxValuePos());
                Boolean commandable = safeCastBoolean(record, edeLayout.getCommandablePos());
                Boolean supportsCov = safeCastBoolean(record, edeLayout.getSupportsCovPos());
                Double hiLimit = safeCastDouble(record, edeLayout.getHiLimitPos());
                Double lowLimit = safeCastDouble(record, edeLayout.getLowLimitPos());
                String stateTextReference = safeString(record, edeLayout.getStateTextReferencePos());
                Integer unitCode = safeCastInteger(record, edeLayout.getUnitCodePos());
                Integer vendorSpecificAddress = safeCastInteger(record, edeLayout.getVendorSpecificAddressPos());
                Integer notificationClass = safeCastInteger(record, edeLayout.getNotificationClassPos());
                Datapoint datapoint = new Datapoint(address, keyName, objectName, description, defaultValue, minValue,
                    maxValue, commandable, supportsCov, hiLimit, lowLimit, stateTextReference, unitCode,
                    vendorSpecificAddress, notificationClass);
                datapoints.put(address, datapoint);
            });
            return datapoints;
        } catch (IOException e) {
            throw new PlcRuntimeException("Error parsing EDE file", e);
        }
    }

    private String safeString(CSVRecord record, int pos) {
        if((pos == -1) || (record.get(pos) == null) || record.get(pos).isEmpty()) {
            return null;
        }
        return record.get(pos);
    }

    private Double safeCastDouble(CSVRecord record, int pos) {
        if(pos == -1) {
            return null;
        }
        try {
            return Double.valueOf(record.get(pos).replace(",", "."));
        } catch (Exception e) {
            return null;
        }
    }

    private Long safeCastLong(CSVRecord record, int pos) {
        if(pos == -1) {
            return null;
        }
        try {
            return Long.valueOf(record.get(pos));
        } catch (Exception e) {
            return null;
        }
    }

    private Integer safeCastInteger(CSVRecord record, int pos) {
        if(pos == -1) {
            return null;
        }
        try {
            return Integer.valueOf(record.get(pos));
        } catch (Exception e) {
            return null;
        }
    }

    private Boolean safeCastBoolean(CSVRecord record, int pos) {
        if(pos == -1) {
            return null;
        }
        try {
            switch (record.get(pos)) {
                case "Y":
                    return true;
                case "N":
                    return false;
                default:
                    return null;
            }
        } catch (Exception e) {
            return null;
        }
    }

}
