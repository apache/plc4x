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

package org.apache.plc4x.plugins.codegenerator;

import org.apache.maven.plugin.MojoExecutionException;
import org.apache.plc4x.codegen.ast.ClassDeclaration;
import org.apache.plc4x.plugins.codegenerator.model.EnumType;
import org.apache.plc4x.plugins.codegenerator.model.SimpleType;
import org.apache.plc4x.plugins.codegenerator.model.SimpleTypeVarLength;
import org.apache.plc4x.plugins.codegenerator.model.Type;
import org.dom4j.*;
import org.dom4j.io.SAXReader;

import java.io.InputStream;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class ProtocolModel {

    private Type rootType;
    private Map<QName, Type> types;

    public ProtocolModel() {
        types = new HashMap<>();
    }

    public void parseBaseSchema(InputStream schemaInputStream) throws MojoExecutionException {
        Document dfdlSpecification = parseDFDLSchema(schemaInputStream);

        // Get all namespaces that are declared.
        Map<String, Namespace> namespaces = new HashMap<>();
        for (Namespace declaredNamespace : dfdlSpecification.getRootElement().declaredNamespaces()) {
            namespaces.put(declaredNamespace.getPrefix(), declaredNamespace);
        }

        // Check that some of the expected namespaces are declared:
        if(!namespaces.containsKey("dfdl")) {
            throw new RuntimeException("No namespace declaration for prefix 'dfdl' defined.");
        }
        Namespace dfdlNamespace = namespaces.get("dfdl");

        if(!namespaces.containsKey("xs")) {
            throw new RuntimeException("No namespace declaration for prefix 'xs' defined.");
        }
        Namespace xsNamespace = namespaces.get("xs");

        // Try to find the target namespace.
        Namespace targetNamespace = null;
        String targetNamespaceUri = dfdlSpecification.getRootElement().attributeValue("targetNamespace");
        for (Namespace namespace : namespaces.values()) {
            if(namespace.getURI().equalsIgnoreCase(targetNamespaceUri)) {
                targetNamespace = namespace;
                break;
            }
        }
        if(targetNamespace == null) {
            throw new RuntimeException("Couldn't find targetNamespace");
        }

        List<Element> simpleTypes = dfdlSpecification.getRootElement().elements(new QName("simpleType", xsNamespace));
        for (Element simpleType : simpleTypes) {
            String typeName = simpleType.attributeValue("name");
            if(simpleType.attributeValue(new QName("length", dfdlNamespace)) == null) {
                types.put(new QName(typeName, targetNamespace), new SimpleTypeVarLength(typeName));
            } else {
                types.put(new QName(typeName, targetNamespace), new SimpleType(typeName));
            }
        }
    }

    public void parseSchema(InputStream schemaInputStream) throws MojoExecutionException {
        Document dfdlSpecification = parseDFDLSchema(schemaInputStream);

        // Get all namespaces that are declared.
        Map<String, Namespace> namespaces = new HashMap<>();
        for (Namespace declaredNamespace : dfdlSpecification.getRootElement().declaredNamespaces()) {
            namespaces.put(declaredNamespace.getPrefix(), declaredNamespace);
        }

        // Check that some of the expected namespaces are declared:
        if(!namespaces.containsKey("dfdl")) {
            throw new RuntimeException("No namespace declaration for prefix 'dfdl' defined.");
        }
        Namespace dfdlNamespace = namespaces.get("dfdl");

        if(!namespaces.containsKey("xs")) {
            throw new RuntimeException("No namespace declaration for prefix 'xs' defined.");
        }
        Namespace xsNamespace = namespaces.get("xs");

        if(!namespaces.containsKey("plc4x")) {
            throw new RuntimeException("No namespace declaration for prefix 'plc4x' defined.");
        }
        Namespace plc4xNamespace = namespaces.get("plc4x");

        // Try to find the target namespace.
        Namespace targetNamespace = null;
        String targetNamespaceUri = dfdlSpecification.getRootElement().attributeValue("targetNamespace");
        for (Namespace namespace : namespaces.values()) {
            if(namespace.getURI().equalsIgnoreCase(targetNamespaceUri)) {
                targetNamespace = namespace;
                break;
            }
        }
        if(targetNamespace == null) {
            throw new RuntimeException("Couldn't find targetNamespace");
        }

        // Parse all simple types (Only Enums are allowed to be included in protocol schemas).
        List<Element> simpleTypes = dfdlSpecification.getRootElement().elements(new QName("simpleType", xsNamespace));
        for (Element simpleTypeElement : simpleTypes) {
            // Get the name of the type
            QName name = new QName(simpleTypeElement.attributeValue("name"), targetNamespace);

            // Get the base-type of the type
            Element restriction = simpleTypeElement.element(new QName("restriction", xsNamespace));
            if(restriction == null) {
                throw new RuntimeException("simpleTypes are only allowed for defining enum types");
            }
            String baseTypeString = restriction.attributeValue("base");
            Namespace baseTypeNamespace = namespaces.get(baseTypeString.substring(0, baseTypeString.indexOf(":")));
            QName baseTypeQName = new QName(baseTypeString.substring(baseTypeString.indexOf(":") + 1), baseTypeNamespace);
            if(!baseTypeQName.getNamespace().equals(plc4xNamespace)) {
                throw new RuntimeException("base is supposed to be coming from the plc4x namespace");
            }
            if(!types.containsKey(baseTypeQName)) {
                throw new RuntimeException("couldn't find base type " + baseTypeQName.getName());
            }
            Type baseType = types.get(baseTypeQName);

            List<Element> enumerations = restriction.elements(new QName("enumeration", xsNamespace));
            List<EnumType.EnumValue> enumValues = new ArrayList<>(enumerations.size());
            for (Element enumeration : enumerations) {
                if(enumeration.attributeValue("value", null) == null) {
                    throw new RuntimeException("enumeration is missing value attribute");
                }
                List<Node> enumNameNodes = enumeration.selectNodes(
                    "xs:annotation/xs:appinfo[@source='http://plc4x.apache.org/plc4x']/plc4x:enumName/text()");
                if(enumNameNodes.size() != 1) {
                    throw new RuntimeException("enumeration value is expecting one enumName");
                }
                String enumName = enumNameNodes.get(0).getText();
                EnumType.EnumValue enumValue = new EnumType.EnumValue(enumName, enumeration.attributeValue("value"));
                enumValues.add(enumValue);
            }

            EnumType enumType = new EnumType(name.getName(), baseType, enumValues);
            types.put(name, enumType);
        }

        // TODO: Parse all the complex types
        // Get the root element (There should only be one)
        // This defines the root of the
        List<Element> rootNodes = dfdlSpecification.getRootElement().elements(new QName("element", xsNamespace));
        if(rootNodes.size() == 0) {
            throw new RuntimeException("Couldn't find root element");
        } else if(rootNodes.size() > 1) {
            throw new RuntimeException("Found multiple root elements");
        }
        Element rootElement = rootNodes.get(0);
        String rootTypeString = rootElement.attributeValue("type");
        Namespace rootTypeNamespace = namespaces.get(rootTypeString.substring(0, rootTypeString.indexOf(":")));
        QName rootTypeName = new QName(rootTypeString.substring(rootTypeString.indexOf(":") + 1), rootTypeNamespace);

        rootType = types.get(rootTypeName);
    }

    private Document parseDFDLSchema(InputStream schemaInputStream) throws MojoExecutionException {
        try {
            SAXReader reader = new SAXReader();
            return reader.read(schemaInputStream);
        } catch (DocumentException e) {
            throw new MojoExecutionException("Unable to parse DFDL schema", e);
        }
    }


}
