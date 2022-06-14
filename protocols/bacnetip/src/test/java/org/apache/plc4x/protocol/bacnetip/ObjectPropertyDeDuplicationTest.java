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
package org.apache.plc4x.protocol.bacnetip;

import org.apache.plc4x.plugins.codegenerator.protocol.TypeContext;
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.fields.Field;
import org.apache.plc4x.plugins.codegenerator.types.fields.FieldConversions;
import org.apache.plc4x.plugins.codegenerator.types.fields.PropertyField;
import org.apache.plc4x.plugins.codegenerator.types.fields.ValidationField;
import org.junit.jupiter.api.*;
import org.opentest4j.TestAbortedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

import static org.apache.plc4x.protocol.bacnetip.BACnetObjectsDefinitions.*;
import static org.junit.jupiter.api.Assertions.*;

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
                                        switch (searchedTypeName) {
                                            case "BACnetConstructedDataOutofService":
                                                // Typo in spec again
                                                searchedTypeName = "BACnetConstructedDataOutOfService";
                                                break;
                                        }
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

    // Note: if this test fails check testNonUniqueUsagesAreMappedGeneric first
    @TestFactory
    Collection<DynamicContainer> testArrayIndexesAreHandled() {
        List<DynamicContainer> tests = new LinkedList<>();
        String bacnetArrayIdentifierPrefix = "BACnetARRAY";
        String bacnetListIdentifierPrefix = "BACnetLIST";
        new LinkedList<>(propertyTypeCombinationToObjectNameMap.entrySet())
            .stream()
            .filter(propertyTypeCombinationListEntry -> propertyTypeCombinationListEntry.getKey().propertyDataType.startsWith(bacnetArrayIdentifierPrefix) || propertyTypeCombinationListEntry.getKey().propertyDataType.startsWith(bacnetListIdentifierPrefix))
            .sorted(Map.Entry.comparingByKey())
            .forEach(propertyTypeCombinationToObjectNameEntry -> {
                PropertyTypeCombination propertyTypeCombination = propertyTypeCombinationToObjectNameEntry.getKey();
                propertyTypeCombinationToObjectNameEntry.getValue().forEach(bacNetObjectName -> {
                        tests.add(DynamicContainer.dynamicContainer(propertyTypeCombination + " for " + bacNetObjectName, () -> {
                            Collection<DynamicNode> nodes = new LinkedList<>();
                            String propertyIdentifier = propertyTypeCombination.propertyIdentifier;
                            String propertyDataType = propertyTypeCombination.propertyDataType;
                            Set<String> listOfTypes = propertyToPropertyTypesMaps.get(propertyIdentifier);
                            final TypeDefinition typeDefinition;
                            if (listOfTypes.size() < 2) {
                                String searchedTypeName = "BACnetConstructedData" + propertyIdentifier;
                                searchedTypeName = searchedTypeName.replaceAll("_", "");
                                switch (searchedTypeName) {
                                    case "BACnetConstructedDataOutofService":
                                        // Typo in spec again
                                        searchedTypeName = "BACnetConstructedDataOutOfService";
                                        break;
                                }
                                typeDefinition = typeDefinitions.get(searchedTypeName);
                                assertNotNull(typeDefinition, "shared " + searchedTypeName + " not found (" + propertyTypeCombination + ")");
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
                                    // This is the case when there are more than 1 occurrence of this propertyIdentifier with one type and this combination is the one with the most occurrences
                                    String searchedTypeName = "BACnetConstructedData" + propertyIdentifier;
                                    searchedTypeName = searchedTypeName.replaceAll("_", "");
                                    typeDefinition = typeDefinitions.get(searchedTypeName);
                                    assertNotNull(typeDefinition, "shared " + searchedTypeName + " not found (most occurring case with " + numberOfOccurences + " occurrences)");
                                } else {
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
                                    typeDefinition = typeDefinitions.get(searchedTypeName);
                                    assertNotNull(typeDefinitions.get(searchedTypeName), "dedicated " + searchedTypeName + " not found (this occurrence: " + numberOfOccurences + ", other variants " + listOfTypes + ").");
                                }
                            }
                            assertTrue(typeDefinition.isComplexTypeDefinition(), typeDefinition.getName() + " should be complex");
                            ComplexTypeDefinition complexTypeDefinition = typeDefinition.asComplexTypeDefinition().orElseThrow();
                            if (propertyDataType.startsWith(bacnetArrayIdentifierPrefix)) {
                                nodes.add(
                                    DynamicTest.dynamicTest("Check array count for " + propertyTypeCombination + " for " + typeDefinition.getName(),
                                        () -> {
                                            Optional<PropertyField> numberOfDataElements = complexTypeDefinition.getPropertyFieldByName("numberOfDataElements");
                                            assertTrue(numberOfDataElements.isPresent(), "field numberOfDataElements for " + typeDefinition.getName() + " not found");
                                        })
                                );
                                if (propertyDataType.startsWith("BACnetARRAY[") && !propertyDataType.startsWith("BACnetARRAY[N")) {
                                    nodes.add(
                                        DynamicTest.dynamicTest("Check bounds validation for " + propertyTypeCombination + " for " + typeDefinition.getName(),
                                            () -> {
                                                Pattern pattern = Pattern.compile("BACnetARRAY\\[(\\d+)]");
                                                Matcher matcher = pattern.matcher(propertyDataType);
                                                assertTrue(matcher.find(), "we should find the index");
                                                String index = matcher.group(1);
                                                Optional<ValidationField> validationField = complexTypeDefinition.getFields().stream()
                                                    .filter(Field::isValidationField)
                                                    .map(ValidationField.class::cast)
                                                    .filter(foundValidationField ->
                                                        foundValidationField.getValidationExpression().stringRepresentation().contains("COUNT")
                                                            && foundValidationField.getValidationExpression().stringRepresentation().contains("arrayIndexArgument")
                                                            && foundValidationField.getValidationExpression().stringRepresentation().contains(index)
                                                    )
                                                    .findAny();
                                                assertTrue(validationField.isPresent(), "No validation for length of " + index + " found for " + typeDefinition.getName());
                                            })
                                    );

                                }
                            } else if (propertyDataType.startsWith(bacnetListIdentifierPrefix))
                                nodes.add(
                                    DynamicTest.dynamicTest("Check no array count for " + propertyTypeCombination + " for " + typeDefinition.getName(),
                                        () -> {
                                            Optional<PropertyField> numberOfDataElements = complexTypeDefinition.getPropertyFieldByName("numberOfDataElements");
                                            assertFalse(numberOfDataElements.isPresent(), "field numberOfDataElements for " + typeDefinition.getName() + " found");
                                        })
                                );
                            else
                                throw new IllegalStateException("how on earth did we got " + propertyDataType + " here???");
                            return nodes.iterator();
                        }));
                    }
                );
            });
        return tests;
    }

    @TestFactory
    Collection<DynamicNode> singleAttributesHaveAnActualValue() {
        List<DynamicNode> tests = new LinkedList<>();
        ComplexTypeDefinition baCnetConstructedData = typeDefinitions.get("BACnetConstructedData").asComplexTypeDefinition().orElseThrow();
        typeDefinitions.values().stream()
            .filter(TypeDefinitionConversions::isDiscriminatedComplexTypeDefinition)
            .filter(typeDefinition -> !typeDefinition.getName().endsWith("All"))
            .filter(typeDefinition -> !typeDefinition.getName().equals("BACnetConstructedDataOptional"))
            .filter(typeDefinition -> !typeDefinition.getName().equals("BACnetConstructedDataRequired"))
            .map(DiscriminatedComplexTypeDefinition.class::cast)
            .filter(discriminatedComplexTypeDefinition -> discriminatedComplexTypeDefinition.getParentType().isPresent())
            .filter(discriminatedComplexTypeDefinition -> discriminatedComplexTypeDefinition.getParentType().get() == baCnetConstructedData)
            .filter(discriminatedComplexTypeDefinition -> discriminatedComplexTypeDefinition.getPropertyFields().stream().noneMatch(FieldConversions::isArrayField))
            .forEach(discriminatedComplexTypeDefinition -> {
                tests.add(DynamicTest.dynamicTest("Test for actualValue on " + discriminatedComplexTypeDefinition.getName(), () -> {
                    assertTrue(discriminatedComplexTypeDefinition.getNamedFieldByName("actualValue").isPresent());
                }));
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

        @Test
        void outputDataConstructedDataChilds() {
            ComplexTypeDefinition baCnetConstructedData = typeDefinitions.get("BACnetConstructedData").asComplexTypeDefinition().orElseThrow();
            Set<String> childsOfConstructedData = typeDefinitions.values().stream()
                .filter(TypeDefinitionConversions::isDiscriminatedComplexTypeDefinition)
                .map(DiscriminatedComplexTypeDefinition.class::cast)
                .filter(discriminatedComplexTypeDefinition -> discriminatedComplexTypeDefinition.getParentType().isPresent())
                .filter(discriminatedComplexTypeDefinition -> discriminatedComplexTypeDefinition.getParentType().get() == baCnetConstructedData)
                .map(TypeDefinition::getName)
                .collect(Collectors.toSet());
            childsOfConstructedData.stream().sorted().forEach(System.out::println);
        }
    }
}
