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
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.xml.sax.SAXException;

import javax.xml.XMLConstants;
import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.parsers.ParserConfigurationException;
import javax.xml.transform.Transformer;
import javax.xml.transform.TransformerException;
import javax.xml.transform.TransformerFactory;
import javax.xml.transform.dom.DOMSource;
import javax.xml.transform.stream.StreamResult;
import javax.xml.xpath.*;
import java.io.File;
import java.io.FilterInputStream;
import java.io.InputStream;
import java.nio.file.StandardCopyOption;
import java.io.IOException;
import java.io.StringWriter;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Comparator;
import java.util.HashMap;
import java.util.Map;
import java.util.Scanner;
import java.util.TreeMap;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Stream;

public class EtsParser {

    private static final Logger LOGGER = LoggerFactory.getLogger(EtsParser.class);

    /**
     * How much a single entry of a knxproj may expand to.
     *
     * <p>A zip entry says how big it will be once inflated, and it can say anything. Compression
     * ratios of a thousand to one are ordinary for repetitive data, so a small file can name a
     * large one - and the XML here is not merely inflated but parsed into a document, which costs
     * several times the text again.</p>
     *
     * <p>The declared size is checked before anything is read, and the bytes are counted as they
     * arrive, because the declaration is written by whoever made the file. The largest knx_master
     * seen in practice is a few tens of megabytes, so this leaves room to spare while keeping a
     * crafted file from deciding how much memory we use.</p>
     */
    private static final long MAX_UNCOMPRESSED_ENTRY_SIZE = 128L * 1024 * 1024;

    private static boolean isEts6Schema(String schemaVersion) {
        if (schemaVersion == null) {
            return false;
        }
        try {
            return Integer.parseInt(schemaVersion) >= 21;
        } catch (NumberFormatException e) {
            return false;
        }
    }

    /**
     * Reads an entry, refusing one that says it is larger than the budget and stopping if it turns
     * out to be larger than it said.
     */
    private static InputStream boundedEntryStream(ZipFile zipFile, FileHeader header)
            throws IOException {
        long declared = header.getUncompressedSize();
        if (declared > MAX_UNCOMPRESSED_ENTRY_SIZE) {
            throw new PlcRuntimeException(String.format(
                "Refusing to read '%s': it declares %d bytes once expanded, over the %d byte limit",
                header.getFileName(), declared, MAX_UNCOMPRESSED_ENTRY_SIZE));
        }
        return new BoundedInputStream(zipFile.getInputStream(header), MAX_UNCOMPRESSED_ENTRY_SIZE,
            header.getFileName());
    }

    /**
     * Writes an entry out under the same budget. Not zip4j's extractFile, which would have written
     * whatever the entry expanded to before anybody could object.
     */
    private static void extractBounded(ZipFile zipFile, FileHeader header, Path target)
            throws IOException {
        try (InputStream in = boundedEntryStream(zipFile, header)) {
            Files.copy(in, target, StandardCopyOption.REPLACE_EXISTING);
        }
    }

    /** Counts what passes through and stops once it is more than was allowed. */
    static final class BoundedInputStream extends FilterInputStream {

        private final long limit;
        private final String entryName;
        private long read;

        BoundedInputStream(InputStream in, long limit, String entryName) {
            super(in);
            this.limit = limit;
            this.entryName = entryName;
        }

        private void count(long justRead) throws IOException {
            if (justRead < 0) {
                return;
            }
            read += justRead;
            if (read > limit) {
                throw new IOException(String.format(
                    "'%s' expanded past the %d byte limit, whatever its header claimed",
                    entryName, limit));
            }
        }

        @Override
        public int read() throws IOException {
            int b = super.read();
            count(b < 0 ? 0 : 1);
            return b;
        }

        @Override
        public int read(byte[] b, int off, int len) throws IOException {
            int n = super.read(b, off, len);
            count(n);
            return n;
        }
    }

    public EtsModel parse(File knxprojFile, String password) {
        try {
            DocumentBuilderFactory factory = DocumentBuilderFactory.newInstance();
            factory.setAttribute(XMLConstants.ACCESS_EXTERNAL_DTD, "");
            factory.setAttribute(XMLConstants.ACCESS_EXTERNAL_SCHEMA, "");
            factory.setFeature("http://apache.org/xml/features/disallow-doctype-decl", true);
            factory.setFeature(XMLConstants.FEATURE_SECURE_PROCESSING, true);
            DocumentBuilder builder = factory.newDocumentBuilder();
            XPathFactory xPathFactory = XPathFactory.newInstance();
            XPath xPath = xPathFactory.newXPath();

            // Unpacked beside the process rather than held in memory, because zip4j needs a file
            // to open the archive that password-protected projects keep inside the archive.
            Path tempDir = Files.createTempDirectory(null);
            try (ZipFile zipFile = new ZipFile(knxprojFile)) {

                final FileHeader knxMasterFileHeader = zipFile.getFileHeader("knx_master.xml");
                if(knxMasterFileHeader == null) {
                    throw new PlcRuntimeException("Error accessing knx_master.xml file.");
                }
                File knxMasterFile = new File(tempDir.toFile(), "knx_master.xml");
                extractBounded(zipFile, knxMasterFileHeader, knxMasterFile.toPath());
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
                // Schema 21 (ETS6) onwards — 22, 23, ... — all use the same
                // PBKDF2-SHA256 derivation with a fixed salt. Anything older
                // uses the raw password.
                EtsFileHandler fileHandler = isEts6Schema(etsSchemaVersion)
                    ? new Ets6FileHandler()
                    : new Ets5FileHandler();

                ////////////////////////////////////////////////////////////////////////////////
                // File containing the information on the type of encoding used for group addresses.
                ////////////////////////////////////////////////////////////////////////////////
                Document projectHeaderDoc;
                Document projectDoc;
                String projectNumber = this.getProjectNumber(zipFile);
                // If the project is not password protected, it contains a file called `project.xml`
                // If it is a password protected one, there is a zip file inside, {project-number}.zip,
                // which is password protected, which then contains the project.xml.
                FileHeader projectFileHeader = zipFile.getFileHeader(projectNumber + "/project.xml");
                if (projectFileHeader == null) {
                    // This is the case of a knxproj file which is password-protected
                    final FileHeader encryptedProjectFileHeader = zipFile.getFileHeader(projectNumber + ".zip");
                    if(encryptedProjectFileHeader == null) {
                        throw new PlcRuntimeException(String.format("Error accessing project header file. Project file '%s/project.xml' and '%s.zip' don't exist.", projectNumber, projectNumber));
                    }
                    // Dump the encrypted zip to a temp file.
                    File tempFile = new File(tempDir.toFile(), projectNumber + ".zip");
                    extractBounded(zipFile, encryptedProjectFileHeader, tempFile.toPath());

                    // Unzip the archive inside the archive.
                    //noinspection CommentedOutCode
                    try (ZipFile projectZipFile = fileHandler.getProjectFiles(tempFile, password)) {
                        final FileHeader compressedProjectFileHeader = projectZipFile.getFileHeader("project.xml");
                        if (compressedProjectFileHeader == null) {
                            throw new PlcRuntimeException(String.format("Error accessing project header file: Project file 'project.xml' inside '%s.zip'.", projectNumber));
                        }
                        projectHeaderDoc = builder.parse(boundedEntryStream(projectZipFile, compressedProjectFileHeader));

                        // TODO: Leave this in, as it helps debug problems, wen a new ETS version comes out.
                        // Print the document content to the console.
                        //System.out.println("Project Header Doc:");
                        //printDocument(projectHeaderDoc);

                        FileHeader projectFileFileHeader = projectZipFile.getFileHeader("0.xml");
                        if (projectFileFileHeader == null) {
                            throw new PlcRuntimeException("Error accessing project file.");
                        }
                        projectDoc = builder.parse(boundedEntryStream(projectZipFile, projectFileFileHeader));

                        // TODO: Leave this in, as it helps debug problems, wen a new ETS version comes out.
                        //System.out.println("Project Doc:");
                        //printDocument(projectDoc);
                    }
                } else {
                    projectHeaderDoc = builder.parse(boundedEntryStream(zipFile, projectFileHeader));
                    FileHeader projectFileFileHeader = zipFile.getFileHeader(projectNumber + "/0.xml");
                    if (projectFileFileHeader == null) {
                        throw new PlcRuntimeException("Error accessing project file.");
                    }
                    projectDoc = builder.parse(boundedEntryStream(zipFile, projectFileFileHeader));
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
                Document knxMasterDoc = builder.parse(boundedEntryStream(zipFile, knxMasterDataFileFileHeader));
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
            } finally {
                // What we unpacked goes again whether we got a model out of it or not. A read that
                // failed is the case that repeats - an operator pointing the driver at the wrong
                // file does it until they find the right one.
                deleteRecursively(tempDir);
            }
        } catch (IOException | ParserConfigurationException | SAXException | XPathExpressionException e) {
            // Zip and Xml Stuff
            throw new PlcRuntimeException(e);
        }
    }

    /**
     * Removes a directory and everything below it, reporting what it could not remove rather than
     * failing the read over it - by the time this runs the model is either built or lost already,
     * and neither outcome is improved by a second error.
     */
    private static void deleteRecursively(Path directory) {
        if (directory == null || !Files.exists(directory)) {
            return;
        }
        try (Stream<Path> paths = Files.walk(directory)) {
            // Deepest first, since a directory only goes once it is empty.
            paths.sorted(Comparator.reverseOrder()).forEach(path -> {
                try {
                    Files.delete(path);
                } catch (IOException e) {
                    LOGGER.warn("Could not remove {} while cleaning up after reading a knxproj", path, e);
                }
            });
        } catch (IOException e) {
            LOGGER.warn("Could not clean up {} after reading a knxproj", directory, e);
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

    @SuppressWarnings("All")
    public void printDocument(Document doc) {
        try {
            DOMSource domSource = new DOMSource(doc);
            StringWriter writer = new StringWriter();
            StreamResult result = new StreamResult(writer);
            TransformerFactory tf = TransformerFactory.newInstance();
            Transformer transformer = tf.newTransformer();
            transformer.transform(domSource, result);
            System.out.println(writer.toString());
        } catch(TransformerException ex) {
            ex.printStackTrace();
        }
    }

}
