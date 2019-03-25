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

package org.apache.plc4x.plugins.codegenerator.utils;

import com.google.common.collect.Lists;
import org.apache.plc4x.plugins.codegenerator.utils.model.*;
import org.apache.plc4x.plugins.codegenerator.utils.model.types.LengthUnit;
import org.dom4j.*;

import java.util.HashMap;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;

public class TemplateHelper {

    private Namespace targetNamespace;
    private Namespace dfdlNamespace;
    private Namespace xsNamespace;
    private Map<String, Namespace> namespaces;

    private String rootName;
    private Type rootType;
    private Map<QName, Type> types;

    public TemplateHelper(Document dfdlSpecification) {
        types = new HashMap<>();

        // Get all namespaces that are declared.
        namespaces = new HashMap<>();
        for (Namespace declaredNamespace : dfdlSpecification.getRootElement().declaredNamespaces()) {
            namespaces.put(declaredNamespace.getPrefix(), declaredNamespace);
        }

        // Check that some of the expected namespaces are declared:
        if(!namespaces.containsKey("dfdl")) {
            throw new RuntimeException("No namespace declaration for prefix 'dfdl' defined.");
        }
        dfdlNamespace = namespaces.get("dfdl");
        if(!namespaces.containsKey("xs")) {
            throw new RuntimeException("No namespace declaration for prefix 'xs' defined.");
        }
        xsNamespace = namespaces.get("xs");

        // Try to find the target namespace.
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

        // Get the root element (There should only be one)
        // This defines the root of the
        List<Node> rootNodes = dfdlSpecification.getRootElement().selectNodes("xs:element");
        if(rootNodes.size() == 0) {
            throw new RuntimeException("Couldn't find root element");
        } else if(rootNodes.size() > 1) {
            throw new RuntimeException("Found multiple root elements");
        }
        Element rootElement = (Element) rootNodes.get(0);
        rootName = rootElement.attributeValue("name");
        String rootTypeString = rootElement.attributeValue("type");
        Namespace rootTypeNamespace = namespaces.get(rootTypeString.substring(0, rootTypeString.indexOf(":")));
        QName rootTypeName = new QName(rootTypeString.substring(rootTypeString.indexOf(":") + 1), rootTypeNamespace);

        // Parse all simple types.
        List<Node> simpleTypes = dfdlSpecification.getRootElement().selectNodes("xs:simpleType");
        for (Node simpleTypeNode : simpleTypes) {
            Element simpleTypeElement = (Element) simpleTypeNode;
            QName name = new QName(simpleTypeElement.attributeValue("name"), targetNamespace);

            LengthUnit lengthUnit = LengthUnit.valueOf(simpleTypeElement.attributeValue(
                new QName("lengthUnits", dfdlNamespace)).toUpperCase());

            int length = Integer.parseInt(simpleTypeElement.attributeValue(new QName("length", dfdlNamespace)));

            String baseTypeString = simpleTypeElement.element(new QName("restriction", xsNamespace)).attributeValue("base");
            Namespace baseTypeNamespace = namespaces.get(baseTypeString.substring(0, baseTypeString.indexOf(":")));
            QName baseType = new QName(baseTypeString.substring(baseTypeString.indexOf(":") + 1), baseTypeNamespace);

            SimpleType simpleType = new SimpleType(name, lengthUnit, length, baseType);
            types.put(name, simpleType);
        }

        // Parse all complex types (In reverse order as we expect
        // referenced types to come after the ones referencing them).
        List<Node> complexTypes = Lists.reverse(dfdlSpecification.getRootElement().selectNodes("xs:complexType"));
        for (Node complexTypeNode : complexTypes) {
            ComplexType complexType = parseComplexType(complexTypeNode);
            types.put(complexType.getName(), complexType);
        }

        rootType = types.get(rootTypeName);
    }

    private ComplexType parseComplexType(Node complexTypeNode) {
        Element complexTypeElement = (Element) complexTypeNode;
        QName name = new QName(complexTypeElement.attributeValue("name"), targetNamespace);

        List<Segment> segments = new LinkedList<>();
        List<Node> segmentNodes = complexTypeElement.selectNodes("xs:sequence/*[(local-name() = 'element') or (local-name() = 'choice')]");
        for (Node segmentNode : segmentNodes) {
            Element segmentElement = (Element) segmentNode;
            switch(segmentElement.getName()) {
                case "element":
                    String fieldName = segmentElement.attributeValue("name");
                    String typeName = segmentElement.attributeValue("type", null);
                    Type fieldType = null;
                    // If a type is specified, this references another type.
                    if(typeName != null) {
                        String namespacePrefix = typeName.substring(0, typeName.indexOf(":"));
                        if(!namespaces.containsKey(namespacePrefix)) {
                            throw new RuntimeException("Namespace prefix " + namespacePrefix + " not defined in root element");
                        }
                        Namespace namespace = namespaces.get(namespacePrefix);
                        typeName = typeName.substring(typeName.indexOf(":") + 1);
                        QName type = new QName(typeName, namespace);
                        if(!types.containsKey(type)) {
                            throw new RuntimeException("Type " + namespacePrefix + ":" + typeName +
                                " not found (maybe defined before current element)");
                        }
                        fieldType = types.get(type);
                    }
                    // In all other cases, this must be a nested complex type.
                    // In this case we simply add the children instead of the complex element itself.
                    else {
                        // Get the children of the inner "complexType" and continue with them.
                        Node innerComplexTypeNode = segmentElement.selectSingleNode("xs:complexType");
                        parseComplexType(innerComplexTypeNode);
                    }

                    // This is an element that can only occur at most once.
                    if("1".equals(segmentElement.attributeValue("maxOccurs", "1"))) {
                        FieldSegment fieldSegment = new FieldSegment(fieldName, fieldType);
                        segments.add(fieldSegment);
                    }
                    // This is an element that can only occur multiple times.
                    else {
                        String maxOccurs = segmentElement.attributeValue("maxOccurs");
                        RepeaterSegment repeaterSegment = new RepeaterSegment(fieldName, fieldType, maxOccurs);
                        segments.add(repeaterSegment);
                    }
                    break;
                case "choice":
                    String discriminatorRule = segmentElement.attributeValue(new QName("choiceDispatchKey", dfdlNamespace));
                    if(!(discriminatorRule.startsWith("{xs:string(") && discriminatorRule.endsWith(")}"))) {
                        throw new RuntimeException("Discriminator rules are expected to have the form " +
                            "'{xs:string({pattern})}', but was " + discriminatorRule);
                    }
                    // Remove the "{xs:string(" and ")}" around the expression.
                    discriminatorRule = discriminatorRule.substring(
                        "{xs:string(".length(), discriminatorRule.length() - 2);
                    Map<String, FieldSegment> choices = new HashMap<>();
                    List<Node> choiceNodes = segmentElement.selectNodes("xs:element");
                    for (Node choiceNode : choiceNodes) {
                        Element choiceElement = (Element) choiceNode;
                        String discriminatorValue = choiceElement.attributeValue(new QName("choiceBranchKey", dfdlNamespace));
                        // TODO: Do some more parsing here ...
                        FieldSegment choiceField = null;
                        choices.put(discriminatorValue, choiceField);
                    }

                    ChoiceSegment choiceSegment = new ChoiceSegment(discriminatorRule, choices);
                    segments.add(choiceSegment);
                    break;
                default:
                    throw new RuntimeException("Unsupported segment type " + segmentElement.getName());
            }
        }
        return new ComplexType(name, segments);
    }


   /* List<FieldSegment> getFields(Element type) {
        List<Node> nodes = type.selectNodes("xs:sequence/xs:element");
        List<FieldSegment> fields = new LinkedList<>();
        for (Node node : nodes) {
            if(node instanceof Element) {
                Element element = (Element) node;
                if(!"true".equalsIgnoreCase(element.attributeValue("fixed"))) {
                }
            }
        }
    }*/

}
