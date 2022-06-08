/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.protocol.bacnetip;

import org.apache.plc4x.plugins.codegenerator.protocol.TypeContext;
import org.apache.plc4x.plugins.codegenerator.types.definitions.EnumTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.junit.jupiter.api.*;
import org.opentest4j.TestAbortedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import static org.apache.plc4x.protocol.bacnetip.BACnetObjectsDefinitions.*;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

public class ObjectPropertyDeDuplicationTest {
    private static final Logger LOGGER = LoggerFactory.getLogger(ObjectPropertyDeDuplicationTest.class);

    public static final String BACNET_OBJECT_TYPE_TYPE_NAME = "BACnetObjectType";

    public static final String BACNET_PROPERTY_IDENTIFIER_TYPE_NAME = "BACnetPropertyIdentifier";
    Map<String, TypeDefinition> typeDefinitions;

    @BeforeEach
    void setUp() throws Exception {
        TypeContext typeContext = new BacNetIpProtocol().getTypeContext();
        typeDefinitions = typeContext.getTypeDefinitions();
    }

    @TestFactory
    Collection<DynamicNode> testThatEveryObjectIsMapped() {
        List<DynamicNode> tests = new LinkedList<>();
        TypeDefinition baCnetObjectType = typeDefinitions.get(BACNET_OBJECT_TYPE_TYPE_NAME);
        assertNotNull(baCnetObjectType, "We could not find essential type " + BACNET_OBJECT_TYPE_TYPE_NAME);
        EnumTypeDefinition baCnetObjectTypeEnumTypeDefinition = (EnumTypeDefinition) baCnetObjectType;
        for (EnumValue enumValue : baCnetObjectTypeEnumTypeDefinition.getEnumValues()) {
            String objectName = enumValue.getName();
            if ("VENDOR_PROPRIETARY_VALUE".equals(objectName)) continue;
            tests.add(
                DynamicTest.dynamicTest("Test definition for " + objectName,
                    () -> {
                        assertNotNull(objectName);
                        assertNotNull(objectNameToBacNetObjectMap.get(objectName), objectName + " has no definition");
                    })
            );
        }
        return tests;
    }

    @TestFactory
    Collection<DynamicNode> testThatEveryPropertyIsUsed() {
        List<DynamicNode> tests = new LinkedList<>();
        TypeDefinition baCnetPropertyIdentifier = typeDefinitions.get(BACNET_PROPERTY_IDENTIFIER_TYPE_NAME);
        assertNotNull(baCnetPropertyIdentifier, "We could not find essential type " + BACNET_OBJECT_TYPE_TYPE_NAME);
        EnumTypeDefinition baCnetPropertyIdentifierEnumTypeDefinition = (EnumTypeDefinition) baCnetPropertyIdentifier;
        for (EnumValue enumValue : baCnetPropertyIdentifierEnumTypeDefinition.getEnumValues()) {
            String propertyIdentifier = enumValue.getName();
            if ("VENDOR_PROPRIETARY_VALUE".equals(propertyIdentifier)) continue;
            tests.add(
                DynamicTest.dynamicTest("Test definition for " + propertyIdentifier,
                    () -> {
                        switch (propertyIdentifier) {
                            case "ALL":
                                //TODO: check what is up with those properties
                            case "LOG_DEVICE_OBJECT_PROPERTY":
                            case "OPTIONAL":
                            case "PROCESS_IDENTIFIER":
                            case "PROTOCOL_CONFORMANCE_CLASS":
                            case "REQUIRED":
                                throw new TestAbortedException(propertyIdentifier + " not in use");
                        }
                        assertNotNull(propertyIdentifier);
                        assertTrue(propertyIdToPropertyNameMap.containsKey(propertyIdentifier), propertyIdentifier + " has no usage");
                    })
            );
        }
        return tests;
    }

    @TestFactory
    Collection<DynamicNode> testUniqueUsagesAreMappedGeneric() {
        List<DynamicNode> tests = new LinkedList<>();
        new LinkedList<>(propertyToObjectNamesMap.entrySet())
            .stream()
            .filter(propertyToObjectNamesEntry -> propertyToObjectNamesEntry.getValue().size() == 1)
            .sorted(Comparator.comparing(stringListEntry -> stringListEntry.getValue().get(0)))
            .forEach(propertyToObjectNameEntry -> {
                String propertyIdentifier = propertyToObjectNameEntry.getKey();
                String bacNetObjectName = propertyToObjectNameEntry.getValue().get(0);
                tests.add(
                    DynamicTest.dynamicTest(bacNetObjectName + " uses property " + propertyIdentifier + " uniquely",
                        () -> {
                            String searchedTypeName = "BACnetConstructedData" + propertyIdentifier;
                            searchedTypeName = searchedTypeName.replaceAll("_", "");
                            switch (searchedTypeName) {
                                case "BACnetConstructedDataOutofService":
                                    // Global group has a typo in this and writes the "o" lowercase
                                    searchedTypeName = "BACnetConstructedDataOutOfService";
                                    break;
                            }
                            assertNotNull(typeDefinitions.get(searchedTypeName), searchedTypeName + " not found");
                        })
                );
            });
        return tests;
    }

    @Disabled("disabled till all of those issues are resolbed within")
    @TestFactory
    Collection<DynamicNode> testNonUniqueUsagesAreMappedGeneric() {
        List<DynamicNode> tests = new LinkedList<>();
        new LinkedList<>(propertyTypeCombinationToObjectNameMap.entrySet())
            .stream()
            .sorted(Map.Entry.comparingByKey())
            .forEach(propertyTypeCombinationToObjectNameEntry -> {
                PropertyTypeCombination propertyTypeCombination = propertyTypeCombinationToObjectNameEntry.getKey();
                propertyTypeCombinationToObjectNameEntry.getValue().forEach(bacNetObjectName -> {
                        String propertyIdentifier = propertyTypeCombination.propertyIdentifier;
                        Set<String> listOfTypes = propertyToPropertyTypesMaps.get(propertyIdentifier);
                        if (listOfTypes.size() < 2) {
                            tests.add(
                                DynamicTest.dynamicTest(propertyTypeCombination + " is used by " + bacNetObjectName + " uses property shared with all having the same type",
                                    () -> {
                                        String searchedTypeName = "BACnetConstructedData" + propertyIdentifier;
                                        searchedTypeName = searchedTypeName.replaceAll("_", "");
                                        assertNotNull(typeDefinitions.get(searchedTypeName), "shared " + searchedTypeName + " not found (" + propertyTypeCombination + ")");
                                    })
                            );
                        } else {
                            boolean isThisCombinationTheMostCommon = true;
                            Integer numberOfOccurences = propertyTypeCombinationCount.get(propertyTypeCombination);
                            for (String otherType : listOfTypes) {
                                if (otherType.equals(propertyTypeCombination.propertyDataType)) continue;
                                Integer otherOccurence = propertyTypeCombinationCount.get(new PropertyTypeCombination(propertyIdentifier, otherType));
                                if (otherOccurence >= numberOfOccurences) {
                                    isThisCombinationTheMostCommon = false;
                                    break;
                                }
                            }
                            if (isThisCombinationTheMostCommon) {
                                tests.add(DynamicTest.dynamicTest(propertyTypeCombination + " is used by " + bacNetObjectName + " uses property shared with " + numberOfOccurences + " using the same type",
                                    () -> {
                                        // This is the case when there are more than 1 occurrence of this propertyIdentifier with one type and this combination is the one with the most occurrences
                                        String searchedTypeName = "BACnetConstructedData" + propertyIdentifier;
                                        searchedTypeName = searchedTypeName.replaceAll("_", "");
                                        assertNotNull(typeDefinitions.get(searchedTypeName), "shared " + searchedTypeName + " not found (most occurring case with " + numberOfOccurences + " occurrences)");
                                    })
                                );
                            } else {
                                tests.add(DynamicTest.dynamicTest(propertyTypeCombination + " is used by " + bacNetObjectName + " uses property shared with this type in the minority.",
                                    () -> {
                                        // This is the case when there are more than 1 occurrence of this propertyIdentifier with one type
                                        String searchedTypeName = "BACnetConstructedData" + bacNetObjectName + propertyIdentifier;
                                        searchedTypeName = searchedTypeName.replaceAll("[_ ]", "");
                                        Pattern pattern = Pattern.compile("-([a-z])");
                                        Matcher matcher = pattern.matcher(searchedTypeName);
                                        StringBuilder result = new StringBuilder();
                                        while (matcher.find()) {
                                            matcher.appendReplacement(result, matcher.group(1).toUpperCase());
                                        }
                                        matcher.appendTail(result);
                                        searchedTypeName = result.toString();
                                        assertNotNull(typeDefinitions.get(searchedTypeName), "dedicated " + searchedTypeName + " not found (this occurrence: " + numberOfOccurences + ", other variants " + listOfTypes + ").");
                                    })
                                );
                            }
                        }
                    }
                );
            });
        return tests;
    }

    @Nested
    @Tag("just-output")
    class JustOutputs {

        @Test
        void outputObjectChapters() {
            List<Integer> unrelatedTables = List.of(0, 1, 7, 9, 15, 25, 33, 41, 42, 59, 60, 63, 65, 66, 67, 68, 70, 72, 73, 74);
            int tableNo = 0;
            for (BacNetObject bacNetObject : bacNetObjects) {
                while (unrelatedTables.contains(tableNo)) tableNo++;
                LOGGER.info("Table 12-{}. Properties of the {} Object Type", tableNo++, bacNetObject.name);
            }
        }

        @Test
        void outputPropertyUsage() {
            propertyToObjectNamesMap.forEach((propertyIdentifier, bacNetObjectNames) -> LOGGER.info("property {} is used by {}", propertyIdentifier, bacNetObjectNames));
        }

        @Test
        void outputTypeCombinationUsage() {
            LinkedList<Map.Entry<PropertyTypeCombination, List<String>>> listOfCombinationEntries = new LinkedList<>(propertyTypeCombinationToObjectNameMap.entrySet());
            listOfCombinationEntries.sort(Comparator.comparingInt(v -> v.getValue().size()));
            Collections.reverse(listOfCombinationEntries);
            listOfCombinationEntries.forEach(propertyTypeCombinationListEntry -> {
                LOGGER.info("{} appearance of {} in {}", propertyTypeCombinationListEntry.getValue().size(), propertyTypeCombinationListEntry.getKey(), propertyTypeCombinationListEntry.getValue());
            });
        }

        @Test
        void outputTypeCombinationsSorted() {
            Set<PropertyTypeCombination> propertyTypeCombinations = propertyTypeCombinationToObjectNameMap.keySet();
            propertyTypeCombinations.stream().sorted().forEach(propertyTypeCombination -> LOGGER.info("{}", propertyTypeCombination));
        }

        @Test
        void outputTypeCombinationCountSorted() {
            propertyTypeCombinationCount.entrySet().stream().sorted(Map.Entry.comparingByKey()).forEach(propertyTypeCombinationCount -> LOGGER.info("{}", propertyTypeCombinationCount));
        }

        @Test
        void outputUniqueProperties() {
            propertyToObjectNamesMap.forEach((propertyIdentifier, bacNetObjectNames) -> {
                if (bacNetObjectNames.size() > 1) {
                    return;
                }
                LOGGER.info("property {} is used by {} uniquely", propertyIdentifier, bacNetObjectNames.get(0));
            });
        }

        @Test
        void outputNonUniqueProperties() {
            propertyToObjectNamesMap.forEach((propertyIdentifier, bacNetObjectNames) -> {
                if (bacNetObjectNames.size() == 1) {
                    return;
                }
                LOGGER.info("property {} is used by {} non uniquely", propertyIdentifier, bacNetObjectNames);
            });
        }
    }
}
