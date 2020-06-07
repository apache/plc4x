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

import net.objecthunter.exp4j.Expression;
import net.objecthunter.exp4j.ExpressionBuilder;
import org.apache.commons.lang3.math.NumberUtils;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.FreemarkerLanguageTemplateHelper;
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.fields.ArrayField;
import org.apache.plc4x.plugins.codegenerator.types.fields.PropertyField;
import org.apache.plc4x.plugins.codegenerator.types.fields.SwitchField;
import org.apache.plc4x.plugins.codegenerator.types.fields.TypedField;
import org.apache.plc4x.plugins.codegenerator.types.references.ComplexTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.FloatTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.SimpleTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.*;
import sun.tools.tree.ArrayExpression;

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

    /**
     * Check if this is an abstract type ...
     * @param typeDefinition the type definition
     * @return true if this is an abstract type
     */
    public boolean isAbstractType(ComplexTypeDefinition typeDefinition) {
        return typeDefinition.isAbstract();
    }

    public boolean isDiscriminatedType(ComplexTypeDefinition typeDefinition) {
        return typeDefinition instanceof DiscriminatedComplexTypeDefinition;
    }

    public List<DiscriminatedComplexTypeDefinition> getDiscriminatedSubTypes(
        ComplexTypeDefinition complexTypeDefinition) {
        // Sebastian would be proud of me ;-)
        SwitchField switchField =  (SwitchField) complexTypeDefinition.getFields().stream().filter(
            field -> field instanceof SwitchField).findFirst().orElse(null);
        if(switchField != null) {
            return switchField.getCases();
        }
        return Collections.emptyList();
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
        final char[] chars = camelCase.toCharArray();
        for(int i = 0; i < chars.length; i++) {
            String lowerCaseChar = String.valueOf(chars[i]).toLowerCase();
            // If the previous letter is a lowercase letter and this one is uppercase, create a new snake-segment.
            if ((i > 0) && !Character.isUpperCase(chars[i - 1]) && Character.isUpperCase(chars[i])) {
                snakeCase.append('_').append(lowerCaseChar);
            }
            else if((i < (chars.length - 2)) && Character.isUpperCase(chars[i]) && !Character.isUpperCase(chars[i + 1])) {
                snakeCase.append('_').append(lowerCaseChar);
            }
            // If this is uppercase and the previous one is too ... just make this letter lowercase.
            else if ((i > 0) && Character.isUpperCase(chars[i - 1]) && Character.isUpperCase(chars[i])) {
                snakeCase.append(lowerCaseChar);
            } else if (chars[i] == '-') {
                snakeCase.append("_");
            } else {
                snakeCase.append(lowerCaseChar);
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
     * Little wrapper around the actual getLanguageTypeName which handles the case of requiring
     * DataIo type fields.
     *
     * @param field field we want to get the type name for
     * @return type name we should use in C
     */
    public String getLanguageTypeNameForField(TypedField field) {
        // If the referenced type is a DataIo type, the value is of type plc4c_data.
        if (field instanceof PropertyField) {
            PropertyField propertyField = (PropertyField) field;
            if (propertyField.getType() instanceof ComplexTypeReference) {
                ComplexTypeReference complexTypeReference = (ComplexTypeReference) propertyField.getType();
                final TypeDefinition typeDefinition = types.get(complexTypeReference.getName());
                if (typeDefinition instanceof DataIoTypeDefinition) {
                    return "plc4c_data*";
                }
            }
        }
        // If this is an array with variable length, then we have to use our "plc4c_list" to store the data.
        if((field instanceof ArrayField) && (!isFixedValueExpression(((ArrayField) field).getLoopExpression()))) {
            return "plc4c_list";
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
                case INT: {
                    StringBuilder sb = new StringBuilder();
                    if(simpleTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.UINT) {
                        sb.append("u");
                    }
                    if(simpleTypeReference.getSizeInBits() % 64 == 0) {
                        sb.append("int64_t");
                    } else if(simpleTypeReference.getSizeInBits() % 32 == 0) {
                        sb.append("int32_t");
                    } else if(simpleTypeReference.getSizeInBits() % 16 == 0) {
                        sb.append("int16_t");
                    } else if(simpleTypeReference.getSizeInBits() % 8 == 0) {
                        sb.append("int8_t");
                    } else {
                        if(simpleTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.UINT) {
                            // We already have the "u" in there ...
                            sb.append("nsigned ");
                        }
                        sb.append("int");
                    }
                    return sb.toString();
                }
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

    public String getLoopExpressionSuffix(TypedField field) {
        if(field instanceof ArrayField) {
            ArrayField arrayField = (ArrayField) field;
            if(arrayField.getLoopType() == ArrayField.LoopType.COUNT) {
                Term countTerm = arrayField.getLoopExpression();
                if (isFixedValueExpression(countTerm)) {
                    int evaluatedCount = evaluateFixedValueExpression(countTerm);
                    return "[" + evaluatedCount +"]";
                }
            }
        }
        return "";
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
                    // If the bit-size is exactly one of the built-in tpye-sizes, omit the suffix.
                    if((simpleTypeReference.getSizeInBits() == 8) || (simpleTypeReference.getSizeInBits() == 16) ||
                        (simpleTypeReference.getSizeInBits() == 32) || (simpleTypeReference.getSizeInBits() == 64)) {
                        return "";
                    }
                    return " : " + simpleTypeReference.getSizeInBits();
                case FLOAT:
                case UFLOAT:
                    // If the bit-size is exactly one of the built-in tpye-sizes, omit the suffix.
                    if((simpleTypeReference.getSizeInBits() == 32) || (simpleTypeReference.getSizeInBits() == 64)) {
                        return "";
                    }
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

    /**
     * Check if the expression doesn't reference any variables.
     * If this is the case, the expression can be evaluated at code-generation time.
     * @param term term
     * @return true if it doesn't reference any variable literals.
     */
    private boolean isFixedValueExpression(Term term) {
        if(term instanceof VariableLiteral) {
            return false;
        }
        if(term instanceof UnaryTerm) {
            UnaryTerm unaryTerm = (UnaryTerm) term;
            return isFixedValueExpression(unaryTerm.getA());
        }
        if(term instanceof BinaryTerm) {
            BinaryTerm binaryTerm = (BinaryTerm) term;
            return isFixedValueExpression(binaryTerm.getA()) && isFixedValueExpression(binaryTerm.getB());
        }
        if(term instanceof TernaryTerm) {
            TernaryTerm ternaryTerm = (TernaryTerm) term;
            return isFixedValueExpression(ternaryTerm.getA()) && isFixedValueExpression(ternaryTerm.getB()) &&
                isFixedValueExpression(ternaryTerm.getC());
        }
        return true;
    }

    private int evaluateFixedValueExpression(Term term) {
        final Expression expression = new ExpressionBuilder(toString(term)).build();
        return (int) expression.evaluate();
    }

    private String toString(Term term) {
        if(term instanceof NullLiteral) {
            return "null";
        }
        if(term instanceof BooleanLiteral) {
            return Boolean.toString(((BooleanLiteral) term).getValue());
        }
        if(term instanceof NumericLiteral) {
            return ((NumericLiteral) term).getNumber().toString();
        }
        if(term instanceof StringLiteral) {
            return "\"" + ((StringLiteral) term).getValue() + "\"";
        }
        if(term instanceof UnaryTerm) {
            return ((UnaryTerm) term).getOperation() + toString(((UnaryTerm) term).getA());
        }
        if(term instanceof BinaryTerm) {
            return toString(((BinaryTerm) term).getA()) + ((BinaryTerm) term).getOperation() + toString(((BinaryTerm) term).getB());
        }
        if(term instanceof TernaryTerm) {
            return "(" + toString(((TernaryTerm) term).getA()) + ") ? (" + toString(((TernaryTerm) term).getB()) +
                ") : (" + toString(((TernaryTerm) term).getC()) + ")";
        }
        return "";
    }

}
