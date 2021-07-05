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
package org.apache.plc4x.codegen;

import org.apache.plc4x.codegen.ast.*;
import org.apache.plc4x.codegen.util.PojoFactory;
import org.dom4j.*;
import org.dom4j.io.SAXReader;
import org.xml.sax.SAXException;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Iterator;

/**
 * Util to create a POJO from a DFDL file.
 */
public class DFDLUtil {

    private static final Namespace XS_NAMESPACE = new Namespace("xs", "http://www.w3.org/2001/XMLSchema");
    private static final QName COMPLEX_TYPE = new QName("complexType", XS_NAMESPACE);
    private static final QName SEQUENCE = new QName("sequence", XS_NAMESPACE);
    private static final QName ELEMENT = new QName("element", XS_NAMESPACE);

    public void transform(File dfdlSchemaFile, File outputDir) throws IOException {
        if (!outputDir.exists()) {
            throw new FileNotFoundException(outputDir.getAbsolutePath() + " does not exist");
        }
        if (!outputDir.isDirectory()) {
            throw new FileNotFoundException(outputDir.getAbsolutePath() + " is not a directory");
        }

        final Document schema = parseDFDLSchema(dfdlSchemaFile);

        final Iterator<Element> iterator = getMainTypes(schema);

        final PojoFactory factory = new PojoFactory();
        while (iterator.hasNext()) {
            final Element element = iterator.next();
            final Iterator<Element> sequenceIterator = element.elementIterator(SEQUENCE);
            final Element sequence = getSequence(sequenceIterator);

            // Now make a POJO with all "elements" as fields
            final ArrayList<PojoFactory.Field> fields = new ArrayList<>();
            final Iterator<Element> elementIterator = sequence.elementIterator(ELEMENT);
            while (elementIterator.hasNext()) {
                final Element elem = elementIterator.next();
                fields.add(new PojoFactory.Field(Expressions.typeOf(elem.attributeValue("type")), elem.attributeValue("name")));
            }
            final PojoFactory.PojoDescription desc = new PojoFactory.PojoDescription(element.attributeValue("name"), fields);
            final ClassDeclaration classDeclaration = factory.create(desc);

            // Now generate the code for that
            final CodeWriter writer = new CodeWriter(4);
            final Generator generator = new JavaGenerator(writer);
            classDeclaration.write(generator);

            System.out.println("Class Definition");
            System.out.println("----------------------------");
            System.out.println(writer.getCode());
        }
    }

    private Element getSequence(Iterator<Element> sequenceIterator) {
        assert sequenceIterator.hasNext();
        final Element sequence = sequenceIterator.next();
        assert !sequenceIterator.hasNext();
        return sequence;
    }

    private Document parseDFDLSchema(File schemaFile) throws IOException {
        try {
            SAXReader reader = new SAXReader();
            reader.setFeature("http://apache.org/xml/features/disallow-doctype-decl", true);
            return reader.read(schemaFile);
        } catch (DocumentException | SAXException e) {
            // Do something
            throw new IOException("Unable to read DFDL Schema File", e);
        }
    }

    private Iterator<Element> getMainTypes(Document dfdlSchema) {
        Element rootElement = dfdlSchema.getRootElement();
        return rootElement.elementIterator(COMPLEX_TYPE);
    }
}
