package org.apache.plc4x.codegen;

import org.apache.plc4x.codegen.ast.ClassDefinition;
import org.apache.plc4x.codegen.ast.CodeWriter;
import org.apache.plc4x.codegen.ast.Generator;
import org.apache.plc4x.codegen.ast.JavaGenerator;
import org.apache.plc4x.codegen.ast.PojoFactory;
import org.apache.plc4x.codegen.ast.Primitive;
import org.apache.plc4x.codegen.ast.PythonGenerator;
import org.apache.plc4x.codegen.ast.TypeNode;
import org.dom4j.Document;
import org.dom4j.DocumentException;
import org.dom4j.Element;
import org.dom4j.Namespace;
import org.dom4j.QName;
import org.dom4j.io.SAXReader;

import java.io.File;
import java.sql.SQLOutput;
import java.util.ArrayList;
import java.util.Iterator;

/**
 * Util to create a POJO from a DFDL file.
 */
public class DFDLUtil {

    private static final Namespace xsNamespace = new Namespace("xs", "http://www.w3.org/2001/XMLSchema");
    private static final QName complexType = new QName("complexType", xsNamespace);
    private static final QName sequence = new QName("sequence", xsNamespace);
    private static final QName ELEMENT = new QName("element", xsNamespace);

    public void transform(File dfdlSchemaFile, File outputDir) {
        assert outputDir.exists();
        assert outputDir.isDirectory();

        final Document schema = parseDFDLSchema(dfdlSchemaFile);

        final Iterator<Element> iterator = getMainTypes(schema);

        final PojoFactory factory = new PojoFactory();
        while (iterator.hasNext()) {
            final Element element = iterator.next();
            final Iterator<Element> sequenceIterator = element.elementIterator(sequence);
            final Element sequence = getSequence(sequenceIterator);

            // Now make a POJO with all "elements" as fields
            final ArrayList<PojoFactory.Field> fields = new ArrayList<>();
            final Iterator<Element> elementIterator = sequence.elementIterator(ELEMENT);
            while (elementIterator.hasNext()) {
                final Element elem = elementIterator.next();
                fields.add(new PojoFactory.Field(new TypeNode(elem.attributeValue("type")), elem.attributeValue("name")));
            }
            final PojoFactory.PojoDescription desc = new PojoFactory.PojoDescription(element.attributeValue("name"), fields);
            final ClassDefinition classDefinition = factory.create(desc);

            // Now generate the code for that
            final CodeWriter writer = new CodeWriter(4);
            final Generator generator = new JavaGenerator(writer);
            classDefinition.write(generator);

            System.out.println("Class Definition");
            System.out.println("----------------------------");
            System.out.println(writer.getCode());
        }
    }

    private Element getSequence(Iterator<Element> sequenceIterator) {
        assert sequenceIterator.hasNext();
        final Element sequence = sequenceIterator.next();
        assert sequenceIterator.hasNext() == false;
        return sequence;
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
