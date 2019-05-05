package org.apache.plc4x.codegen;

import org.dom4j.Document;
import org.dom4j.DocumentException;
import org.dom4j.Element;
import org.dom4j.Namespace;
import org.dom4j.QName;
import org.dom4j.io.SAXReader;

import java.io.File;
import java.util.Iterator;

/**
 * Util to create a POJO from a DFDL file.
 */
public class DFDLUtil {

    private static final Namespace xsNamespace = new Namespace("xs", "http://www.w3.org/2001/XMLSchema");
    private static final QName complexType = new QName("complexType", xsNamespace);

    public void transform(File dfdlSchemaFile, File outputDir) {
        assert outputDir.exists();
        assert outputDir.isDirectory();


        final Document schema = parseDFDLSchema(dfdlSchemaFile);

        final Iterator<Element> iterator = getMainTypes(schema);

        while (iterator.hasNext()) {
            final Element element = iterator.next();
            System.out.println(element);
        }
    }

    private Document parseDFDLSchema(File schemaFile) {
        try {
            SAXReader reader = new SAXReader();
            return reader.read(schemaFile);
        } catch (DocumentException e) {
            // Do something
            throw new RuntimeException("Unable to read DFDL Schema File", e);
        }
    }

    private Iterator<Element> getMainTypes(Document dfdlSchema) {
        Element rootElement = dfdlSchema.getRootElement();
        return rootElement.elementIterator(complexType);
    }
}
