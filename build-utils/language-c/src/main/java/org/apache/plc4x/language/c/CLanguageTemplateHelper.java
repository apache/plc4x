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
package org.apache.plc4x.language.c;

import org.apache.commons.lang3.math.NumberUtils;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.FreemarkerLanguageTemplateHelper;
import org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.DataIoTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.EnumTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.fields.PropertyField;
import org.apache.plc4x.plugins.codegenerator.types.fields.TypedField;
import org.apache.plc4x.plugins.codegenerator.types.references.ComplexTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.FloatTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.SimpleTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;

import java.util.*;

public class CLanguageTemplateHelper implements FreemarkerLanguageTemplateHelper {

    private final Map<String, TypeDefinition> types;
    private TypeDefinition thisType;
    private String protocolName;
    private String flavorName;

    public CLanguageTemplateHelper(Map<String, TypeDefinition> types) {
        this.types = types;
    }

    public void setConstants(TypeDefinition thisType, String protocolName, String flavorName) {
        this.thisType = thisType;
        this.protocolName = protocolName;
        this.flavorName = flavorName;
    }

    public String getSourceDirectory() {
        return String.join("", protocolName.split("-")) + ".src";
    }

    public String getIncludesDirectory() {
        return String.join("", protocolName.split("-")) + ".includes";
    }

    public Collection<ComplexTypeReference> getComplexTypeReferencesInFields() {
        List<ComplexTypeReference> complexTypeReferences = new LinkedList<>();
        if (thisType instanceof ComplexTypeDefinition) {
            for (PropertyField propertyField : ((ComplexTypeDefinition) thisType).getAllPropertyFields()) {
                if (propertyField.getType() instanceof ComplexTypeReference) {
                    ComplexTypeReference complexTypeReference = (ComplexTypeReference) propertyField.getType();
                    complexTypeReferences.add(complexTypeReference);
                }
            }
        } else if (thisType instanceof EnumTypeDefinition) {
            for (String constantName : ((EnumTypeDefinition) thisType).getConstantNames()) {
                final TypeReference constantType = ((EnumTypeDefinition) thisType).getConstantType(constantName);
                if (constantType instanceof ComplexTypeReference) {
                    ComplexTypeReference complexTypeReference = (ComplexTypeReference) constantType;
                    complexTypeReferences.add(complexTypeReference);
                }
            }
        }
        return complexTypeReferences;
    }

    /**
     * Little helper that converts a given type name in camel-case into a c-style snake-case expression.
     * In addition it appends a prefix for the protocol name and the output flavor.
     *
     * @param typeName camel-case type name
     * @return c-style type name
     */
    public String getCTypeName(String typeName) {
        return camelCaseToSnakeCase(protocolName).toLowerCase() + "_" + camelCaseToSnakeCase(flavorName).toLowerCase() + "_" +
            camelCaseToSnakeCase(typeName).toLowerCase();
    }

    /**
     * Converts a camel-case string to snake-case.
     *
     * @param camelCase camel-case string
     * @return snake-case string
     */
    public String camelCaseToSnakeCase(String camelCase) {
        StringBuilder snakeCase = new StringBuilder();
        for (char c : camelCase.toCharArray()) {
            if (Character.isUpperCase(c)) {
                snakeCase.append('_').append(String.valueOf(c).toLowerCase());
            } else if (c == '-') {
                snakeCase.append("_");
            } else {
                snakeCase.append(c);
            }
        }
        // If the first letter was a capital letter, the string will start with a "_".
        // In this case we cut that char off.
        if (snakeCase.indexOf("_") == 0) {
            return snakeCase.toString().substring(1);
        }
        return snakeCase.toString();
    }

    /**
     * Little wrapper around the actual getLanguageTypeName which handles the case of reqriting
     * DataIo type fields.
     *
     * @param field field we want to get the type name for
     * @return type name we should use in C
     */
    public String getLanguageTypeNameForField(TypedField field) {
        // If the referenced type is a DataIo type, the value is of type PlcValue.
        if (field instanceof PropertyField) {
            PropertyField propertyField = (PropertyField) field;
            if (propertyField.getType() instanceof ComplexTypeReference) {
                ComplexTypeReference complexTypeReference = (ComplexTypeReference) propertyField.getType();
                final TypeDefinition typeDefinition = types.get(complexTypeReference.getName());
                if (typeDefinition instanceof DataIoTypeDefinition) {
                    return "PlcValue";
                }
            }
        }
        TypeReference typeReference = field.getType();
        return getLanguageTypeName(typeReference);
    }

    /**
     * Converts a given type-reference into a concrete type in C
     * If it's a complex type, this is trivial, as the typename then follows the usual pattern.
     * For simple types it's a little more complex as depending on the parameters the concrete type will be different.
     *
     * @param typeReference type reference
     * @return c type
     */
    public String getLanguageTypeName(TypeReference typeReference) {
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT:
                    return "bool";
                case UINT:
                    return "unsigned int";
                case INT:
                    return "int";
                case FLOAT:
                    FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                    int sizeInBits = ((floatTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.FLOAT) ? 1 : 0) +
                        floatTypeReference.getExponent() + floatTypeReference.getMantissa();
                    if (sizeInBits <= 32) {
                        return "float";
                    }
                    if (sizeInBits <= 64) {
                        return "double";
                    }
                    return "unsupported";
                case UFLOAT:
                    throw new RuntimeException("Unsigned floats are not supported");
                case STRING:
                    return "char *";
                case TIME:
                    return "unsupported";
                case DATE:
                    return "unsupported";
                case DATETIME:
                    return "unsupported";
            }
            return "unsupported";
        } else {
            return "plc4c_" + getCTypeName(((ComplexTypeReference) typeReference).getName());
        }
    }

    /**
     * Ge the type-size suffix in case of simple types.
     *
     * @param field the field we want to get the type-size for
     * @return a type-size string for the given field or an empty string if this does not apply
     */
    public String getTypeSizeForField(TypedField field) {
        TypeReference typeReference = field.getType();
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT:
                    return " : 1";
                case UINT:
                case INT:
                case FLOAT:
                case UFLOAT:
                    return " : " + simpleTypeReference.getSizeInBits();
                case STRING:
                case TIME:
                case DATE:
                case DATETIME:
                    return "";
            }
        }
        return "";
    }

    public String escapeValue(TypeReference typeReference, String valueString) {
        if (valueString == null) {
            return null;
        }
        if ("null".equals(valueString)) {
            // C doesn't like NULL values for enums, so we have to return something else (we'll treat -1 as NULL)
            if (typeReference instanceof ComplexTypeReference) {
                if (types.get(((ComplexTypeReference) typeReference).getName()) instanceof EnumTypeDefinition) {
                    return "-1";
                }
            }
            return "NULL";
        }
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case UINT:
                case INT:
                    // C doesn't like this hex notation, so we have to convert it to a numeric one
                    if (valueString.startsWith("0x")) {
                        return Long.toString(Long.parseLong(valueString.substring(2), 16));
                    }
                    // If it's a one character string and is numeric, output it as char.
                    else if (!NumberUtils.isParsable(valueString) && (valueString.length() == 1)) {
                        return "'" + valueString + "'";
                    }
                    break;
                case STRING:
                    return "\"" + valueString + "\"";
            }
        }
        return valueString;
    }

    public String escapeEnumValue(TypeReference typeReference, String valueString) {
        // Currently the only case in which here complex type references are used are when referencing enum constants.
        if(typeReference instanceof ComplexTypeReference) {
            // C doesn't like NULL values for enums, so we have to return something else (we'll treat -1 as NULL)
            if("null".equals(valueString)) {
                return "-1";
            }
            ComplexTypeReference complexTypeReference = (ComplexTypeReference) typeReference;
            String typeName = valueString.substring(0, valueString.indexOf('.'));
            String constantName = valueString.substring(valueString.indexOf('.') + 1);
            return "plc4c_" + getCTypeName(typeName) + "_" + constantName;
        } else {
            return escapeValue(typeReference, valueString);
        }
    }

    public Collection<EnumValue> getUniqueEnumValues(EnumValue[] enumValues) {
        Map<String, EnumValue> filteredEnumValues = new TreeMap<>();
        for (EnumValue enumValue : enumValues) {
            if (!filteredEnumValues.containsKey(enumValue.getValue())) {
                filteredEnumValues.put(enumValue.getValue(), enumValue);
            }
        }
        return filteredEnumValues.values();
    }

}
