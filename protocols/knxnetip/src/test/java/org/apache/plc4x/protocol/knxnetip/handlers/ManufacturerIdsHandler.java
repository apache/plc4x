package org.apache.plc4x.protocol.knxnetip.handlers;

import org.xml.sax.Attributes;
import org.xml.sax.SAXException;
import org.xml.sax.helpers.DefaultHandler;

import java.util.ArrayList;
import java.util.List;

public class ManufacturerIdsHandler extends DefaultHandler {

    private boolean inElement = false;

    private List<Integer> manufacturerIds = new ArrayList<>();

    @Override
    public void startElement(String uri, String localName, String qName, Attributes attributes) throws SAXException {
        inElement = qName.equalsIgnoreCase("unsignedShort");
    }

    @Override
    public void endElement(String uri, String localName, String qName) throws SAXException {
        inElement = ! (inElement && qName.equalsIgnoreCase("unsignedShort"));
    }

    public void characters(char[] ch, int start, int length) throws SAXException {
        if(inElement) {
            String content = new String(ch, start, length);
            int manufacturerId = Integer.parseInt(content);
            manufacturerIds.add(manufacturerId);
        }
    }

    public List<Integer> getManufacturerIds() {
        return manufacturerIds;
    }

}
