/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.bacnetip.ede;

import org.apache.commons.csv.CSVFormat;
import org.apache.commons.csv.CSVParser;
import org.apache.commons.csv.CSVRecord;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.bacnetip.ede.model.Datapoint;
import org.apache.plc4x.java.bacnetip.ede.model.EdeModel;
import org.apache.plc4x.java.bacnetip.field.BacNetIpField;

import java.io.*;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;

public class EdeParser {

    public EdeModel parse(File edeFile) {
        if(!edeFile.exists()) {
            throw new PlcRuntimeException("EDE File at " + edeFile.getPath() + " doesn't exist.");
        }
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
            // Check if the version of the layout is 3
            if(!layoutVersionRow.get(1).equals("3")) {
                throw new PlcRuntimeException("Unsupported EDE file layout version " + layoutVersionRow.get(2) +
                    ". Currently only version 3 is supported.");
            }
            // Just skip the next row.
            iterator.next();
            // Get the column names.
            final CSVRecord columnNames = iterator.next();


            Map<BacNetIpField, Datapoint> datapoints = new HashMap<>();
            // Process the content.
            iterator.forEachRemaining(record -> {
                BacNetIpField address = new BacNetIpField(Long.parseLong(
                    record.get(1)), Integer.parseInt(record.get(3)), Long.parseLong(record.get(4)));
                String keyName = (record.get(0).length() == 0) ? null : record.get(0);
                String objectName = (record.get(2).length() == 0) ? null : record.get(2);
                String description = (record.get(5).length() == 0) ? null : record.get(5);
                Double defaultValue = (record.get(6).length() == 0) ? null : Double.valueOf(record.get(6).replace(",", "."));
                Double minValue = (record.get(7).length() == 0) ? null : Double.valueOf(record.get(7).replace(",", "."));
                Double maxValue = (record.get(8).length() == 0) ? null : Double.valueOf(record.get(8).replace(",", "."));
                Boolean commandable = (record.get(9).length() == 0) ? null :
                    record.get(9).equalsIgnoreCase("Y") ? Boolean.TRUE : record.get(9).equalsIgnoreCase("N") ? Boolean.TRUE : null;
                Boolean supportsCov = (record.get(10).length() == 0) ? null :
                    record.get(10).equalsIgnoreCase("Y") ? Boolean.TRUE : record.get(10).equalsIgnoreCase("N") ? Boolean.TRUE : null;
                Double hiLimit = (record.get(11).length() == 0) ? null : Double.valueOf(record.get(11).replace(",", "."));
                Double lowLimit = (record.get(12).length() == 0) ? null : Double.valueOf(record.get(12).replace(",", "."));
                String stateTextReference = (record.get(13).length() == 0) ? null : record.get(13);
                Integer unitCode = (record.get(14).length() == 0) ? null : Integer.valueOf(record.get(14));
                Integer vendorSpecificAddress = (record.get(15).length() == 0) ? null : Integer.valueOf(record.get(15));
                Datapoint datapoint = new Datapoint(address, keyName, objectName, description, defaultValue, minValue,
                    maxValue, commandable, supportsCov, hiLimit, lowLimit, stateTextReference, unitCode,
                    vendorSpecificAddress);
                datapoints.put(address, datapoint);
            });
            return new EdeModel(datapoints);
        } catch (IOException e) {
            throw new PlcRuntimeException("Error parsing EDE file", e);
        }

    }

    public static void main(String[] args) throws Exception {
        EdeModel model = new EdeParser().parse(new File("/Users/christofer.dutz/Projects/Apache/PLC4X-Documents/BacNET/Merck/EDE-Files/F135/edeDataText3029.csv"));
        System.out.println(model);
    }

}
