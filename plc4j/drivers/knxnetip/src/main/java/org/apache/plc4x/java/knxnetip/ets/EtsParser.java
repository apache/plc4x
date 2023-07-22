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
package org.apache.plc4x.java.knxnetip.ets;

import net.lingala.zip4j.ZipFile;
import net.lingala.zip4j.exception.ZipException;
import net.lingala.zip4j.model.FileHeader;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.knxnetip.ets.filehandlers.Ets5FileHandler;
import org.apache.plc4x.java.knxnetip.ets.filehandlers.Ets6FileHandler;
import org.apache.plc4x.java.knxnetip.ets.filehandlers.EtsFileHandler;
import org.apache.plc4x.java.knxnetip.ets.model.AddressType;
import org.apache.plc4x.java.knxnetip.ets.model.EtsModel;
import org.apache.plc4x.java.knxnetip.ets.model.Function;
import org.apache.plc4x.java.knxnetip.ets.model.GroupAddress;
import org.apache.plc4x.java.knxnetip.readwrite.KnxDatapointType;
import org.w3c.dom.Attr;
import org.w3c.dom.Document;
import org.w3c.dom.Element;
import org.w3c.dom.NodeList;
import org.xml.sax.SAXException;

import javax.xml.XMLConstants;
import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.parsers.ParserConfigurationException;
import javax.xml.xpath.*;
import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.HashMap;
import java.util.Map;
import java.util.Scanner;
import java.util.TreeMap;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class EtsParser {

    public EtsModel parse(File knxprojFile, String password) {
        try {
            DocumentBuilderFactory factory = DocumentBuilderFactory.newInstance();
            factory.setAttribute(XMLConstants.ACCESS_EXTERNAL_DTD, "");
            factory.setAttribute(XMLConstants.ACCESS_EXTERNAL_SCHEMA, "");
            DocumentBuilder builder = factory.newDocumentBuilder();
            XPathFactory xPathFactory = XPathFactory.newInstance();
            XPath xPath = xPathFactory.newXPath();

            try (ZipFile zipFile = new ZipFile(knxprojFile)) {
                Path tempDir = Files.createTempDirectory(null);

                final FileHeader knxMasterFileHeader = zipFile.getFileHeader("knx_master.xml");
                if(knxMasterFileHeader == null) {
                    throw new PlcRuntimeException("Error accessing knx_master.xml file.");
                }
                zipFile.extractFile("knx_master.xml", tempDir.toFile().getAbsolutePath());
                File knxMasterFile = new File(tempDir.toFile(), "knx_master.xml");
                // If the file contains: <KNX xmlns="http://knx.org/xml/project/21"> it's an ETS6 file
                // In all other cases, we'll treat it as ETS5
                String etsSchemaVersion = null;
                try(Scanner scanner = new Scanner(knxMasterFile)) {
                    while (scanner.hasNextLine()) {
                        final String curLine = scanner.nextLine();
                        if(curLine.contains("http://knx.org/xml/project/")) {
                            etsSchemaVersion = curLine.substring(curLine.indexOf("http://knx.org/xml/project/") + "http://knx.org/xml/project/".length());
                            etsSchemaVersion = etsSchemaVersion.substring(0, etsSchemaVersion.indexOf("\""));
                            break;
                        }
                    }
                }
                EtsFileHandler fileHandler = ("21".equals(etsSchemaVersion)) ? new Ets6FileHandler() : new Ets5FileHandler();

                ////////////////////////////////////////////////////////////////////////////////
                // File containing the information on the type of encoding used for group addresses.
                ////////////////////////////////////////////////////////////////////////////////
                Document projectHeaderDoc;
                Document projectDoc;
                String projectNumber = this.getProjectNumber(zipFile);
                FileHeader projectFileHeader = zipFile.getFileHeader(projectNumber + "/project.xml");
                if (projectFileHeader == null) {
                    // This is the case of a knxproj file which is password-protected
                    final FileHeader encryptedProjectFileHeader = zipFile.getFileHeader(projectNumber + ".zip");
                    if(encryptedProjectFileHeader == null) {
                        throw new PlcRuntimeException(String.format("Error accessing project header file. Project file '%s/project.xml' and '%s.zip' don't exist.", projectNumber, projectNumber));
                    }
                    // Dump the encrypted zip to a temp file.
                    zipFile.extractFile(projectNumber + ".zip", tempDir.toFile().getAbsolutePath());
                    File tempFile = new File(tempDir.toFile(), projectNumber + ".zip");

                    // Unzip the archive inside the archive.
                    try (ZipFile projectZipFile = fileHandler.getProjectFiles(tempFile, password)) {
                        final FileHeader compressedProjectFileHeader = projectZipFile.getFileHeader("project.xml");
                        if (compressedProjectFileHeader == null) {
                            throw new PlcRuntimeException(String.format("Error accessing project header file: Project file 'project.xml' inside '%s.zip'.", projectNumber));
                        }
                        projectHeaderDoc = builder.parse(projectZipFile.getInputStream(compressedProjectFileHeader));
                        FileHeader projectFileFileHeader = projectZipFile.getFileHeader("0.xml");
                        if (projectFileFileHeader == null) {
                            throw new PlcRuntimeException("Error accessing project file.");
                        }
                        projectDoc = builder.parse(projectZipFile.getInputStream(projectFileFileHeader));
                    }
                } else {
                    projectHeaderDoc = builder.parse(zipFile.getInputStream(projectFileHeader));
                    FileHeader projectFileFileHeader = zipFile.getFileHeader(projectNumber + "/0.xml");
                    if (projectFileFileHeader == null) {
                        throw new PlcRuntimeException("Error accessing project file.");
                    }
                    projectDoc = builder.parse(zipFile.getInputStream(projectFileFileHeader));
                }
                final XPathExpression xpathGroupAddressStyle = xPath.compile("/KNX/Project/ProjectInformation/@GroupAddressStyle");
                Attr groupAddressStyle = (Attr) xpathGroupAddressStyle.evaluate(projectHeaderDoc, XPathConstants.NODE);
                byte groupAddressStyleCode = getGroupAddressLevel(groupAddressStyle.getTextContent());

                ////////////////////////////////////////////////////////////////////////////////
                // General information on the type of encoding and the value ranges.
                ////////////////////////////////////////////////////////////////////////////////
                FileHeader knxMasterDataFileFileHeader = zipFile.getFileHeader("knx_master.xml");
                if (knxMasterDataFileFileHeader == null) {
                    throw new PlcRuntimeException("Error accessing KNX master file.");
                }
                Document knxMasterDoc = builder.parse(zipFile.getInputStream(knxMasterDataFileFileHeader));
                final XPathExpression xpathDatapointSubtype = xPath.compile("//DatapointSubtype");
                NodeList datapointSubtypeNodes = (NodeList) xpathDatapointSubtype.evaluate(knxMasterDoc, XPathConstants.NODESET);

                // Build an index of the internal data-types.
                Map<String, KnxDatapointType> knxDatapointTypeMap = new TreeMap<>();
                for (KnxDatapointType value : KnxDatapointType.values()) {
                    knxDatapointTypeMap.put(value.getDatapointMainType().getNumber() + "#" + value.getNumber(), value);
                }

                Map<String, AddressType> addressTypes = new HashMap<>();
                for (int i = 0; i < datapointSubtypeNodes.getLength(); i++) {
                    final Element datapointSubtypeNode = (Element) datapointSubtypeNodes.item(i);
                    final String id = datapointSubtypeNode.getAttribute("Id");
                    final int subType = Integer.parseInt(datapointSubtypeNode.getAttribute("Number"));
                    final int mainType = Integer.parseInt(
                        ((Element) datapointSubtypeNode.getParentNode().getParentNode()).getAttribute("Number"));
                    final String name = datapointSubtypeNode.getAttribute("Text");
                    addressTypes.put(id, new AddressType(id, mainType, subType, name));
                }

                ////////////////////////////////////////////////////////////////////////////////
                // File containing all the information about group addresses used, their names, types etc.
                ////////////////////////////////////////////////////////////////////////////////
                final Map<String, String> topologyNames = new HashMap<>();
                final XPathExpression topology = xPath.compile("//Topology");
                final Element topologyElement = (Element) topology.evaluate(projectDoc, XPathConstants.NODE);
                final NodeList areas = topologyElement.getElementsByTagName("Area");
                for (int a = 0; a < areas.getLength(); a++) {
                    final Element areaNode = (Element) areas.item(a);
                    final String curAreaAddress = areaNode.getAttribute("Address");
                    topologyNames.put(curAreaAddress, areaNode.getAttribute("Name"));

                    final NodeList lines = areaNode.getElementsByTagName("Line");
                    for (int l = 0; l < lines.getLength(); l++) {
                        final Element lineNode = (Element) lines.item(l);
                        final String curLineAddress = curAreaAddress + "/" + lineNode.getAttribute("Address");
                        topologyNames.put(curLineAddress, lineNode.getAttribute("Name"));
                    }
                }

                final Map<String, Function> groupAddressRefs = new HashMap<>();
                final XPathExpression xpathGroupAddressRef = xPath.compile("//GroupAddressRef");
                NodeList groupAddressRefNodes = (NodeList) xpathGroupAddressRef.evaluate(projectDoc, XPathConstants.NODESET);
                for (int i = 0; i < groupAddressRefNodes.getLength(); i++) {
                    final Element groupAddressRefNode = (Element) groupAddressRefNodes.item(i);
                    final String refId = groupAddressRefNode.getAttribute("RefId");
                    final Element functionNode = (Element) groupAddressRefNode.getParentNode();
                    final String functionName = functionNode.getAttribute("Name");
                    // Function Type information is stored in knx_master.xml (//FunctionType[@id='functionTypeId']
                    final String functionTypeId = functionNode.getAttribute("Type");
                    final Element spaceNode = (Element) functionNode.getParentNode();
                    final String spaceName = spaceNode.getAttribute("Name");

                    final Function function = new Function(refId, functionName, functionTypeId, spaceName);
                    groupAddressRefs.put(refId, function);
                }

                final XPathExpression xpathGroupAddresses = xPath.compile("//GroupAddress");
                NodeList groupAddressNodes = (NodeList) xpathGroupAddresses.evaluate(projectDoc, XPathConstants.NODESET);
                Map<String, GroupAddress> groupAddresses = new HashMap<>();
                for (int i = 0; i < groupAddressNodes.getLength(); i++) {
                    final Element groupAddressNode = (Element) groupAddressNodes.item(i);

                    final String id = groupAddressNode.getAttribute("Id");
                    final Function function = groupAddressRefs.get(id);

                    final int addressInt = Integer.parseInt(groupAddressNode.getAttribute("Address"));
                    final String knxGroupAddress = EtsModel.parseGroupAddress(groupAddressStyleCode, addressInt);

                    final String name = groupAddressNode.getAttribute("Name");

                    final String typeString = groupAddressNode.getAttribute("DatapointType");
                    final AddressType addressType = addressTypes.get(typeString);

                    if (addressType != null) {
                        // Lookup the driver internal data-type.
                        final KnxDatapointType datapointType = knxDatapointTypeMap.get(
                            addressType.getMainType() + "#" + addressType.getSubType());

                        GroupAddress groupAddress = new GroupAddress(knxGroupAddress, name, datapointType, function);
                        groupAddresses.put(knxGroupAddress, groupAddress);
                    }
                }
                return new EtsModel(groupAddressStyleCode, groupAddresses, topologyNames);
            }
        } catch (IOException | ParserConfigurationException | SAXException | XPathExpressionException e) {
            // Zip and Xml Stuff
            throw new PlcRuntimeException(e);
        }
    }

    private byte getGroupAddressLevel(String knxprojValue) {
        if ("ThreeLevel".equals(knxprojValue)) {
            return (byte) 3;
        } else if ("TwoLevel".equals(knxprojValue)) {
            return (byte) 2;
        } else if ("Free".equals(knxprojValue)) {
            return (byte) 1;
        }
        throw new PlcRuntimeException("Unsupported GroupAddressStyle=" + knxprojValue);
    }

    private String getProjectNumber(ZipFile zipFile) throws ZipException
    {
        // The following is defined within the KNX spec: Project Scheme 2.1 public documentation
        //
        // P-iiii/project.xml Created by user; contains the global data for project iiii (internal project ID, formatted as 4 hex digits).
        Pattern pattern = Pattern.compile( "^P-[\\dA-F]{4}", Pattern.CASE_INSENSITIVE);
        for (FileHeader fileHeader : zipFile.getFileHeaders ( ) ) {
            Matcher matcher = pattern.matcher(fileHeader.getFileName());
            if (matcher.find()) {
                return matcher.group();
            }
        }
        
        throw new PlcRuntimeException("Error determining project number.");
    }

}
