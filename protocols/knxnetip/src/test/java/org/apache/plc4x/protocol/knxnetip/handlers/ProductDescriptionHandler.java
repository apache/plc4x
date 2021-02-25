package org.apache.plc4x.protocol.knxnetip.handlers;

import org.xml.sax.Attributes;
import org.xml.sax.SAXException;
import org.xml.sax.helpers.DefaultHandler;

import java.util.HashMap;
import java.util.Map;

public class ProductDescriptionHandler extends DefaultHandler {

    private Map<String, Integer> addresses = new HashMap<>();
    private String maskVersion = null;
    private String name = null;
    private Integer comObjectTableAddress = null;

    @Override
    public void startElement(String uri, String localName, String qName, Attributes attributes) throws SAXException {
        if(qName.equalsIgnoreCase("ApplicationProgram")) {
            String maskVersionString = attributes.getValue("MaskVersion");
            maskVersion = maskVersionString.substring(maskVersionString.indexOf('-') + 1);
            name = attributes.getValue("Name");
        } else if(qName.equalsIgnoreCase("AbsoluteSegment")) {
            String id = attributes.getValue("Id");
            Integer address = Integer.parseInt(attributes.getValue("Address"));
            addresses.put(id, address);
        } else if(qName.equalsIgnoreCase("ComObjectTable")) {
            String codeSegment = attributes.getValue("CodeSegment");
            comObjectTableAddress = addresses.get(codeSegment);
        }
    }

    public String getMaskVersion() {
        return maskVersion;
    }

    public String getName() {
        return name;
    }

    public Integer getComObjectTableAddress() {
        return comObjectTableAddress;
    }

}
