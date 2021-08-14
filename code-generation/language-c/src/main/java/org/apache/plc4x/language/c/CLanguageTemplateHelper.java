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
package org.apache.plc4x.language.c;

import org.apache.commons.lang3.math.NumberUtils;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.BaseFreemarkerLanguageTemplateHelper;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.FreemarkerException;
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.references.*;
import org.apache.plc4x.plugins.codegenerator.types.terms.*;

import java.util.*;
import java.util.function.Function;

public class CLanguageTemplateHelper extends BaseFreemarkerLanguageTemplateHelper {

    public CLanguageTemplateHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types) {
        super(thisType, protocolName, flavorName, types);
    }

    public String getSourceDirectory() {
        return String.join("", getProtocolName().split("-")) + ".src";
    }

    public String getIncludesDirectory() {
        return String.join("", getProtocolName().split("-")) + ".include";
    }

    /**
     * Little helper that converts a given type name in camel-case into a c-style snake-case expression.
     * In addition it appends a prefix for the protocol name and the output flavor.
     *
     * @param typeName camel-case type name
     * @return c-style type name
     */
    public String getCTypeName(String typeName) {
        return "plc4c_" + camelCaseToSnakeCase(getProtocolName()).toLowerCase() +
            "_" + camelCaseToSnakeCase(getFlavorName()).toLowerCase() +
            "_" + camelCaseToSnakeCase(typeName).toLowerCase();
    }

    public String getFieldName(ComplexTypeDefinition baseType, NamedField field) {
        StringBuilder sb = new StringBuilder();
        if (baseType != getThisTypeDefinition()) {
            sb.append(camelCaseToSnakeCase(baseType.getName())).append("_");
        }
        sb.append(camelCaseToSnakeCase(field.getName()));
        return sb.toString();
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
        for (int i = 0; i < chars.length; i++) {
            String lowerCaseChar = String.valueOf(chars[i]).toLowerCase();
            // If the previous letter is a lowercase letter and this one is uppercase, create a new snake-segment.
            if ((i > 0) && !Character.isUpperCase(chars[i - 1]) && Character.isUpperCase(chars[i])) {
                snakeCase.append('_').append(lowerCaseChar);
            } else if ((i < (chars.length - 2)) && Character.isUpperCase(chars[i]) && !Character.isUpperCase(chars[i + 1])) {
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
    @Override
    public String getLanguageTypeNameForField(Field field) {
        if(!(field instanceof TypedField)) {
            throw new FreemarkerException("Field " + field + " is not a TypedField");
        }
        // If this is an array with variable length, then we have to use our "plc4c_list" to store the data.
        if ((field instanceof ArrayField) && (!isFixedValueExpression(((ArrayField) field).getLoopExpression()))) {
            return "plc4c_list";
        }
        TypedField typedField = (TypedField) field;
        TypeReference typeReference = typedField.getType();
        if (typeReference instanceof ComplexTypeReference) {
            final TypeDefinition typeDefinition = getTypeDefinitionForTypeReference(typeReference);
            if (typeDefinition instanceof DataIoTypeDefinition) {
                return "plc4c_data*";
            }
        }
        return getLanguageTypeNameForTypeReference(typeReference);
    }

    public Map<ComplexTypeDefinition, ConstField> getAllConstFields() {
        Map<ComplexTypeDefinition, ConstField> constFields = new HashMap<>();
        ((ComplexTypeDefinition) getThisTypeDefinition()).getConstFields().forEach(
            constField -> constFields.put((ComplexTypeDefinition) getThisTypeDefinition(), constField));
        if(getSwitchField() != null) {
            for (DiscriminatedComplexTypeDefinition switchCase : getSwitchField().getCases()) {
                switchCase.getConstFields().forEach(
                    constField -> constFields.put(switchCase, constField));
            }
        }
        return constFields;
    }

    /**
     * If a property references a complex type in an argument, we need to pass that as a pointer,
     * same with optional fields.
     *
     * @param typeDefinition type that contains the property or attribute.
     * @param propertyName name of the property or attribute
     * @return true if the access needs to be using pointers
     */
    public boolean requiresPointerAccess(ComplexTypeDefinition typeDefinition, String propertyName) {
        final Optional<NamedField> namedFieldOptional = typeDefinition.getFields().stream().filter(field -> field instanceof NamedField).map(field -> (NamedField) field).filter(namedField -> namedField.getName().equals(propertyName)).findFirst();
        // If the property name refers to a field, check if it's an optional field.
        // If it is, pointer access is required, if not, it's not.
        if(namedFieldOptional.isPresent()) {
            final NamedField namedField = namedFieldOptional.get();
            if(namedField instanceof TypedField) {
                TypedField typedField = (TypedField) namedField;
                return !(namedField instanceof EnumField) && (isComplexTypeReference(typedField.getType()));
            }
            return false;
        }
        final Optional<Argument> parserArgument = Arrays.stream(typeDefinition.getParserArguments()).filter(argument -> argument.getName().equals(propertyName)).findFirst();
        // If the property name refers to a parser argument, as soon as it's a complex type,
        // pointer access is required.
        return parserArgument.filter(argument -> argument.getType() instanceof ComplexTypeReference).isPresent();
        // In all other cases, the property might be a built-in constant, so we don't need pointer
        // access for any of these.
    }

    /**
     * Converts a given type-reference into a concrete type in C
     * If it's a complex type, this is trivial, as the typename then follows the usual pattern.
     * For simple types it's a little more complex as depending on the parameters the concrete type will be different.
     *
     * @param typeReference type reference
     * @return c type
     */
    @Override
    public String getLanguageTypeNameForTypeReference(TypeReference typeReference) {
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT:
                    return "bool";
                case BYTE:
                    return "char";
                case UINT: {
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 8) {
                        return "uint8_t";
                    }
                    if (integerTypeReference.getSizeInBits() <= 16) {
                        return "uint16_t";
                    }
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return "uint32_t";
                    }
                    if (integerTypeReference.getSizeInBits() <= 64) {
                        return "uint64_t";
                    }
                    throw new RuntimeException("Unsupported simple type");
                }
                case INT: {
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 8) {
                        return "int8_t";
                    }
                    if (integerTypeReference.getSizeInBits() <= 16) {
                        return "int16_t";
                    }
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return "int32_t";
                    }
                    if (integerTypeReference.getSizeInBits() <= 64) {
                        return "int64_t";
                    }
                    throw new RuntimeException("Unsupported simple type");
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
                    throw new FreemarkerException("Unsupported float type with " + sizeInBits + " bits.");
                case UFLOAT:
                    throw new FreemarkerException("Unsupported unsigned float type.");
                case STRING:
                    return "char*";
                case TIME:
                    return "time_t";//throw new FreemarkerException("Unsupported time type.");
                case DATE:
                    return "time_t";//throw new FreemarkerException("Unsupported date type.");
                case DATETIME:
                    return "time_t";//throw new FreemarkerException("Unsupported date-time type.");
            }
            throw new FreemarkerException("Unsupported simple type. " + simpleTypeReference.getBaseType());
        } else {
            return getCTypeName(((ComplexTypeReference) typeReference).getName());
        }
    }

    public String getLoopExpressionSuffix(TypedField field) {
        if (field instanceof ArrayField) {
            ArrayField arrayField = (ArrayField) field;
            if (arrayField.getLoopType() == ArrayField.LoopType.COUNT) {
                Term countTerm = arrayField.getLoopExpression();
                if (isFixedValueExpression(countTerm)) {
                    int evaluatedCount = evaluateFixedValueExpression(countTerm);
                    return "[" + evaluatedCount + "]";
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
                case BYTE:
                    return " : 8";
                case UINT:
                case INT:
                    // If the bit-size is exactly one of the built-in tpye-sizes, omit the suffix.
                    if ((simpleTypeReference.getSizeInBits() == 8) || (simpleTypeReference.getSizeInBits() == 16) ||
                        (simpleTypeReference.getSizeInBits() == 32) || (simpleTypeReference.getSizeInBits() == 64)) {
                        return "";
                    }
                    return " : " + simpleTypeReference.getSizeInBits();
                case FLOAT:
                case UFLOAT:
                    // If the bit-size is exactly one of the built-in tpye-sizes, omit the suffix.
                    if ((simpleTypeReference.getSizeInBits() == 32) || (simpleTypeReference.getSizeInBits() == 64)) {
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
            return "NULL";
        }
        if ("null".equals(valueString)) {
            // C doesn't like NULL values for enums, so we have to return something else (we'll treat -1 as NULL)
            if (typeReference instanceof ComplexTypeReference) {
                if (getTypeDefinitionForTypeReference(typeReference) instanceof EnumTypeDefinition) {
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
        if (typeReference instanceof ComplexTypeReference) {
            // C doesn't like NULL values for enums, so we have to return something else (we'll treat -1 as NULL)
            if ("null".equals(valueString)) {
                return "-1";
            }
            String typeName = getLanguageTypeNameForTypeReference(typeReference);
            return typeName + "_" + valueString;
        } else {
            return escapeValue(typeReference, valueString);
        }
    }

    @Override
    public String getReadBufferReadMethodCall(SimpleTypeReference simpleTypeReference, String valueString, TypedField field) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "plc4c_spi_read_bit(readBuffer, (bool*) " + valueString + ")";
            case BYTE:
                return "plc4c_spi_read_char(readBuffer, (char*) " + valueString + ")";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    return "plc4c_spi_read_unsigned_byte(readBuffer, " + unsignedIntegerTypeReference.getSizeInBits() + ", (uint8_t*) " + valueString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    return "plc4c_spi_read_unsigned_short(readBuffer, " + unsignedIntegerTypeReference.getSizeInBits() + ", (uint16_t*) " + valueString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    return "plc4c_spi_read_unsigned_int(readBuffer, " + unsignedIntegerTypeReference.getSizeInBits() + ", (uint32_t*) " + valueString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 64) {
                    return "plc4c_spi_read_unsigned_long(readBuffer, " + unsignedIntegerTypeReference.getSizeInBits() + ", (uint64_t*) " + valueString + ")";
                }
                throw new FreemarkerException("Unsupported unsigned integer type with " + unsignedIntegerTypeReference.getSizeInBits() + " bits");
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "plc4c_spi_read_signed_byte(readBuffer, " + integerTypeReference.getSizeInBits() + ", (int8_t*) " + valueString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "plc4c_spi_read_signed_short(readBuffer, " + integerTypeReference.getSizeInBits() + ", (int16_t*) " + valueString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "plc4c_spi_read_signed_int(readBuffer, " + integerTypeReference.getSizeInBits() + ", (int32_t*) " + valueString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "plc4c_spi_read_signed_long(readBuffer, " + integerTypeReference.getSizeInBits() + ", (int64_t*) " + valueString + ")";
                }
                throw new FreemarkerException("Unsupported signed integer type with " + integerTypeReference.getSizeInBits() + " bits");
            case FLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                if (floatTypeReference.getSizeInBits() <= 32) {
                    return "plc4c_spi_read_float(readBuffer, " + floatTypeReference.getSizeInBits() + ", (float*) " + valueString + ")";
                } else if(floatTypeReference.getSizeInBits() <= 64) {
                    return "plc4c_spi_read_double(readBuffer, " + floatTypeReference.getSizeInBits() + ", (double*) " + valueString + ")";
                }
                throw new FreemarkerException("Unsupported float type with " + floatTypeReference.getSizeInBits() + " bits");
            case STRING:
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return "plc4c_spi_read_string(readBuffer, " + toParseExpression(getThisTypeDefinition(), field, stringTypeReference.getLengthExpression(), null) + ", \"" +
                    stringTypeReference.getEncoding() + "\"" + ", (char**) " + valueString + ")";
            default:
                throw new FreemarkerException("Unsupported type " + simpleTypeReference.getBaseType().name());
        }
    }

    @Override
    public String getWriteBufferWriteMethodCall(SimpleTypeReference simpleTypeReference, String fieldName, TypedField field) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "plc4c_spi_write_bit(writeBuffer, " + fieldName + ")";
            case BYTE:
                return "plc4c_spi_write_char(writeBuffer, " + fieldName + ")";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    return "plc4c_spi_write_unsigned_byte(writeBuffer, " + unsignedIntegerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    return "plc4c_spi_write_unsigned_short(writeBuffer, " + unsignedIntegerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    return "plc4c_spi_write_unsigned_int(writeBuffer, " + unsignedIntegerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 64) {
                    return "plc4c_spi_write_unsigned_long(writeBuffer, " + unsignedIntegerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                throw new FreemarkerException("Unsupported unsigned integer type with " + unsignedIntegerTypeReference.getSizeInBits() + " bits");
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "plc4c_spi_write_signed_byte(writeBuffer, " + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "plc4c_spi_write_signed_short(writeBuffer, " + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "plc4c_spi_write_signed_int(writeBuffer, " + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "plc4c_spi_write_signed_long(writeBuffer, " + integerTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                throw new FreemarkerException("Unsupported signed integer type with " + integerTypeReference.getSizeInBits() + " bits");
            case FLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                if (floatTypeReference.getSizeInBits() <= 32) {
                    return "plc4c_spi_write_float(writeBuffer, " + floatTypeReference.getSizeInBits() + ", " + fieldName + ")";
                } else if(floatTypeReference.getSizeInBits() <= 64) {
                    return "plc4c_spi_write_double(writeBuffer, " + floatTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                throw new FreemarkerException("Unsupported float type with " + floatTypeReference.getSizeInBits() + " bits");
            case STRING:
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return "plc4c_spi_write_string(writeBuffer, " + toSerializationExpression(getThisTypeDefinition(), field, stringTypeReference.getLengthExpression(), null) + ", \"" +
                    stringTypeReference.getEncoding() + "\", " + fieldName + ")";
            default:
                throw new FreemarkerException("Unsupported type " + simpleTypeReference.getBaseType().name());
        }
    }

    @Override
    public String getNullValueForTypeReference(TypeReference typeReference) {
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT:
                    return "false";
                case BYTE:
                    return "0";
                case UINT:
                case INT:
                    return "0";
                case FLOAT:
                case UFLOAT:
                    return "0.0f";
                case STRING:
                    return "\"\"";
                case TIME:
                    throw new FreemarkerException("Unsupported time type.");
                case DATE:
                    throw new FreemarkerException("Unsupported date type.");
                case DATETIME:
                    throw new FreemarkerException("Unsupported date-time type.");
                default:
                    throw new FreemarkerException("Unsupported type.");
            }
        } else {
            ComplexTypeReference complexTypeReference = (ComplexTypeReference) typeReference;
            return getCTypeName(complexTypeReference.getName()) + "_null()";
        }
    }

    public String getReservedValue(ReservedField reservedField) {
        final String languageTypeName = getLanguageTypeNameForTypeReference(reservedField.getType());
        if("BigInteger".equals(languageTypeName)) {
            return "BigInteger.valueOf(" + reservedField.getReferenceValue() + ")";
        } else {
            return "(" + languageTypeName + ") " + reservedField.getReferenceValue();
        }
    }

    public String toParseExpression(TypeDefinition baseType, Field field, Term term, Argument[] parserArguments) {
        return toExpression(baseType, field, term, term1 -> toVariableParseExpression(baseType, field, term1, parserArguments));
    }

    public String toSerializationExpression(TypeDefinition baseType, Field field, Term term, Argument[] parserArguments) {
        return toExpression(baseType, field, term, term1 -> toVariableSerializationExpression(baseType, field, term1, parserArguments));
    }

    private String toExpression(TypeDefinition baseType, Field field, Term term, Function<Term, String> variableExpressionGenerator) {
        if (term == null) {
            return "";
        }
        if (term instanceof Literal) {
            if (term instanceof NullLiteral) {
                return "NULL";
            } else if (term instanceof BooleanLiteral) {
                return Boolean.toString(((BooleanLiteral) term).getValue());
            } else if (term instanceof NumericLiteral) {
                return ((NumericLiteral) term).getNumber().toString();
            } else if (term instanceof StringLiteral) {
                return "\"" + ((StringLiteral) term).getValue() + "\"";
            } else if (term instanceof VariableLiteral) {
                VariableLiteral variableLiteral = (VariableLiteral) term;
                if(variableLiteral.contains("lengthInBytes")) {
                    TypeDefinition lengthType;
                    String lengthExpression;
                    if(variableLiteral.getName().equals("lengthInBytes")) {
                        lengthType = (baseType.getParentType() == null) ? baseType : (ComplexTypeDefinition) baseType.getParentType();
                        lengthExpression = "_message";
                    } else {
                        final Optional<TypeReference> typeReferenceForProperty = getTypeReferenceForProperty( (ComplexTypeDefinition) baseType, variableLiteral.getName());
                        if(!typeReferenceForProperty.isPresent()) {
                            throw new FreemarkerException("Unknown type for property " + variableLiteral.getName());
                        }
                        lengthType = getTypeDefinitionForTypeReference(typeReferenceForProperty.get());
                        lengthExpression = variableExpressionGenerator.apply(term);
                    }
                    return getCTypeName(lengthType.getName()) + "_length_in_bytes(" + lengthExpression + ")";
                } else if (variableLiteral.getName().equals("lastItem")) {
                    return "lastItem";
                // If this literal references an Enum type, then we have to output it differently.
                } else if (getTypeDefinitions().get(variableLiteral.getName()) instanceof EnumTypeDefinition) {
                    return getCTypeName(variableLiteral.getName()) + "_" + variableLiteral.getChild().getName();
                } else {
                    return variableExpressionGenerator.apply(term);
                }
            } else {
                throw new FreemarkerException("Unsupported Literal type " + term.getClass().getName());
            }
        } else if (term instanceof UnaryTerm) {
            UnaryTerm ut = (UnaryTerm) term;
            Term a = ut.getA();
            switch (ut.getOperation()) {
                case "!":
                    return "!(" + toExpression(baseType, field, a, variableExpressionGenerator) + ")";
                case "-":
                    return "-(" + toExpression(baseType, field, a, variableExpressionGenerator) + ")";
                case "()":
                    return "(" + toExpression(baseType, field, a, variableExpressionGenerator) + ")";
                default:
                    throw new FreemarkerException("Unsupported unary operation type " + ut.getOperation());
            }
        } else if (term instanceof BinaryTerm) {
            BinaryTerm bt = (BinaryTerm) term;
            Term a = bt.getA();
            Term b = bt.getB();
            String operation = bt.getOperation();
            if ("^".equals(operation)) {
                return "Math.pow((" + toExpression(baseType, field, a, variableExpressionGenerator) + "), (" + toExpression(baseType, field, b, variableExpressionGenerator) + "))";
            }
            return "(" + toExpression(baseType, field, a, variableExpressionGenerator) + ") " + operation + " (" + toExpression(baseType, field, b, variableExpressionGenerator) + ")";
        } else if (term instanceof TernaryTerm) {
            TernaryTerm tt = (TernaryTerm) term;
            if ("if".equals(tt.getOperation())) {
                Term a = tt.getA();
                Term b = tt.getB();
                Term c = tt.getC();
                return "((" + toExpression(baseType, field, a, variableExpressionGenerator) + ") ? " + toExpression(baseType, field, b, variableExpressionGenerator) + " : " + toExpression(baseType, field, c, variableExpressionGenerator) + ")";
            } else {
                throw new FreemarkerException("Unsupported ternary operation type " + tt.getOperation());
            }
        } else {
            throw new FreemarkerException("Unsupported Term type " + term.getClass().getName());
        }
    }

    public String toVariableParseExpression(TypeDefinition baseType, Field field, Term term, Argument[] parserArguments) {
        VariableLiteral vl = (VariableLiteral) term;
        if("CAST".equals(vl.getName())) {

            if((vl.getArgs() == null) || (vl.getArgs().size() != 2)) {
                throw new FreemarkerException("A CAST expression expects exactly two arguments.");
            }
            final VariableLiteral sourceTerm = (VariableLiteral) vl.getArgs().get(0);
            final VariableLiteral typeTerm = (VariableLiteral) vl.getArgs().get(1);
            ComplexTypeReference castTypeReference = typeTerm::getName;
            final TypeDefinition castType = getTypeDefinitionForTypeReference(castTypeReference);
            // If we're casting to a sub-type of a discriminated value, we got to cast to the parent
            // type instead and add the name of the sub-type as prefix to the property we're tryging to
            // access next.
            StringBuilder sb = new StringBuilder();
            sb.append("((");
            if(castType.getParentType() != null) {
                sb.append(getCTypeName(castType.getParentType().getName()));
            } else {
                sb.append(getCTypeName(castType.getName()));
            }
            sb.append("*) (");
            sb.append(toVariableParseExpression(baseType, field, sourceTerm, parserArguments)).append("))");
            if(vl.getChild() != null) {
                if(castType.getParentType() != null) {
                    // Change the name of the property to contain the sub-type-prefix.
                    sb.append("->").append(camelCaseToSnakeCase(castType.getName())).append("_");
                    appendVariableExpressionRest(sb, baseType, vl.getChild());
                } else {
                    sb.append("->");
                    appendVariableExpressionRest(sb, castType, vl.getChild());
                }
            }
            return sb.toString();
        }
        // STATIC_CALL implies that driver specific static logic should be called
        if ("STATIC_CALL".equals(vl.getName())) {
            String functionName = ((StringLiteral) vl.getArgs().get(0)).getValue();
            // We'll cut off the java package structure and just take the segment after the last "."
            functionName = functionName.substring(functionName.lastIndexOf('.') + 1, functionName.length() -1);
            // But to make the function name unique, well add the driver prefix to it.
            StringBuilder sb = new StringBuilder(getCTypeName(functionName));
            if (vl.getArgs().size() > 1) {
                sb.append("(");
                boolean firstArg = true;
                for (int i = 1; i < vl.getArgs().size(); i++) {
                    Term arg = vl.getArgs().get(i);
                    if (!firstArg) {
                        sb.append(", ");
                    }
                    sb.append(toParseExpression(baseType, field, arg, parserArguments));
                    firstArg = false;
                }
                sb.append(")");
            }
            return sb.toString();
        }
        // Any name that is full upper-case is considered a function call.
        // These are generally defined in the spi file evaluation_helper.c.
        // All should have a name prefix "plc4c_spi_evaluation_helper_".
        else if (vl.getName().equals(vl.getName().toUpperCase())) {
            StringBuilder sb = new StringBuilder("plc4c_spi_evaluation_helper_" + vl.getName().toLowerCase());
            if (vl.getArgs() != null) {
                sb.append("(");
                boolean firstArg = true;
                for (Term arg : vl.getArgs()) {
                    if (!firstArg) {
                        sb.append(", ");
                    }
                    sb.append(toParseExpression(baseType, field, arg, parserArguments));
                    firstArg = false;
                }
                sb.append(")");
            }
            if (vl.getIndex() != VariableLiteral.NO_INDEX) {
                sb.append("[").append(vl.getIndex()).append("]");
            }
            if(vl.getChild() != null) {
                sb.append(".");
                appendVariableExpressionRest(sb, baseType, vl.getChild());
            }
            return sb.toString();
        } else if("readBuffer".equals(vl.getName())) {
            StringBuilder sb = new StringBuilder("readBuffer");
            if(vl.getChild() != null) {
                sb.append(".");
                appendVariableExpressionRest(sb, baseType, vl.getChild());
            }
            return sb.toString();
        } else if("writeBuffer".equals(vl.getName())) {
            StringBuilder sb = new StringBuilder("writeBuffer");
            if(vl.getChild() != null) {
                sb.append(".");
                appendVariableExpressionRest(sb, baseType, vl.getChild());
            }
            return sb.toString();
        } else if("_type".equals(vl.getName())) {
            if((vl.getChild() != null) && "encoding".equals(vl.getChild().getName()) && (field instanceof TypedField) && (((TypedField) field).getType() instanceof StringTypeReference)) {
                TypedField typedField = (TypedField) field;
                StringTypeReference stringTypeReference = (StringTypeReference) typedField.getType();
                return "\"" + stringTypeReference.getEncoding().substring(1, stringTypeReference.getEncoding().length() - 1) + "\"";
            } else {
                throw new FreemarkerException("_type is currently pretty much hard-coded for some usecases, please check CLanguageTemplateHelper.toVariableParseExpression");
            }
        }

        final String name = vl.getName();

        // In case of DataIo types, we'll just check the arguments.
        if(baseType instanceof DataIoTypeDefinition) {
            if(baseType.getParserArguments() != null) {
                for (Argument parserArgument : baseType.getParserArguments()) {
                    if(parserArgument.getName().equals(name)) {
                        return name;
                    }
                }
            }
        }

        // Try to find the type of the addressed property.
        final Optional<TypeReference> propertyTypeOptional =
            getTypeReferenceForProperty((ComplexTypeDefinition) baseType, name);

        // If we couldn't find the type, we didn't find the property.
        if(!propertyTypeOptional.isPresent()) {
            throw new FreemarkerException("Could not find property with name '" + name + "' in type " + baseType.getName());
        }

        final TypeReference propertyType = propertyTypeOptional.get();

        // If it's a simple field, there is no sub-type to access.
        if(propertyType instanceof SimpleTypeReference) {
            if(vl.getChild() != null) {
                throw new FreemarkerException("Simple property '" + name + "' doesn't have child properties.");
            }
            return name;
        }

        // If it references a complex, type we need to get that type's definition first.
        final TypeDefinition propertyTypeDefinition = getTypeDefinitions().get(((ComplexTypeReference) propertyType).getName());
        // If we're not accessing any child property, no need to handle anything special.
        if(vl.getChild() == null) {
            return name;
        }
        // If there is a child we need to check if this is a discriminator property.
        // As discriminator properties are not real properties, but saved in the static metadata
        // of a type, we need to generate a different access pattern.
        if(propertyTypeDefinition instanceof ComplexTypeDefinition) {
            final Optional<DiscriminatorField> discriminatorFieldOptional = ((ComplexTypeDefinition) propertyTypeDefinition).getFields().stream().filter(
                curField -> curField instanceof DiscriminatorField).map(curField -> (DiscriminatorField) curField).filter(
                discriminatorField -> discriminatorField.getName().equals(vl.getChild().getName())).findFirst();
            // If child references a discriminator field of the type we found, we have to escape it.
            if (discriminatorFieldOptional.isPresent()) {
                return getCTypeName(propertyTypeDefinition.getName()) + "_get_discriminator(" + name + "->_type)." + vl.getChild().getName();
            }
        }
        // Handling enum properties in C is a little more tricky as we have to use the enum value
        // and pass this to a function that then returns the desired property value.
        else if(propertyTypeDefinition instanceof EnumTypeDefinition) {
            return getCTypeName(propertyTypeDefinition.getName()) +
                "_get_" + camelCaseToSnakeCase(vl.getChild().getName()) +
                "(" + vl.getName() + ")";
        }
        // Else ... generate a simple access path.
        StringBuilder sb = new StringBuilder(vl.getName());
        if(vl.getChild() != null) {
            sb.append(".");
            appendVariableExpressionRest(sb, baseType, vl.getChild());
        }
        return sb.toString();
    }

    private String toVariableSerializationExpression(TypeDefinition baseType, Field field, Term term, Argument[] serialzerArguments) {
        VariableLiteral vl = (VariableLiteral) term;
        if ("STATIC_CALL".equals(vl.getName())) {
            StringBuilder sb = new StringBuilder();
            if (!(vl.getArgs().get(0) instanceof StringLiteral)) {
                throw new FreemarkerException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral");
            }
            String methodName = ((StringLiteral) vl.getArgs().get(0)).getValue();
            methodName = methodName.substring(1, methodName.length() - 1);
            sb.append(methodName).append("(");
            for (int i = 1; i < vl.getArgs().size(); i++) {
                Term arg = vl.getArgs().get(i);
                if (i > 1) {
                    sb.append(", ");
                }
                if (arg instanceof VariableLiteral) {
                    VariableLiteral va = (VariableLiteral) arg;
                    // "io" and "_value" are always available in every parser.
                    boolean isSerializerArg = "io".equals(va.getName()) || "_value".equals(va.getName()) || "element".equals(va.getName());
                    boolean isTypeArg = "_type".equals(va.getName());
                    if (!isSerializerArg && !isTypeArg && serialzerArguments != null) {
                        for (Argument serializerArgument : serialzerArguments) {
                            if (serializerArgument.getName().equals(va.getName())) {
                                isSerializerArg = true;
                                break;
                            }
                        }
                    }
                    if (isSerializerArg) {
                        sb.append(va.getName());
                        if(va.getChild() != null) {
                            sb.append(".");
                            appendVariableExpressionRest(sb, baseType, va.getChild());
                        }
                    } else if (isTypeArg) {
                        String part = va.getChild().getName();
                        switch (part) {
                            case "name":
                                sb.append("\"").append(field.getTypeName()).append("\"");
                                break;
                            case "length":
                                sb.append("\"").append(((SimpleTypeReference) field).getSizeInBits()).append("\"");
                                break;
                            case "encoding":
                                if(!(field instanceof TypedField)) {
                                    throw new FreemarkerException("'encoding' only supported for typed fields.");
                                }
                                TypedField typedField = (TypedField) field;
                                if(!(typedField.getType() instanceof StringTypeReference)) {
                                    throw new FreemarkerException("Can only access 'encoding' for string types.");
                                }
                                StringTypeReference stringTypeReference = (StringTypeReference) typedField.getType();
                                String encoding = stringTypeReference.getEncoding();
                                // Cut off the single quotes.
                                encoding = encoding.substring(1, encoding.length() - 1);
                                sb.append("\"").append(encoding).append("\"");
                                break;
                        }
                    } else {
                        sb.append(toVariableSerializationExpression(baseType, field, va, null));
                    }
                } else if (arg instanceof StringLiteral) {
                    sb.append(((StringLiteral) arg).getValue());
                }
            }
            sb.append(")");
            return sb.toString();
        }
        // All uppercase names are not fields, but utility methods.
        else if (vl.getName().equals(vl.getName().toUpperCase())) {
            StringBuilder sb = new StringBuilder("plc4c_spi_evaluation_helper_" + vl.getName().toLowerCase());
            if (vl.getArgs() != null) {
                sb.append("(");
                boolean firstArg = true;
                for (Term arg : vl.getArgs()) {
                    if (!firstArg) {
                        sb.append(", ");
                    }

                    if (arg instanceof VariableLiteral) {
                        VariableLiteral va = (VariableLiteral) arg;
                        boolean isSerializerArg = "io".equals(va.getName());
                        boolean isTypeArg = "_type".equals(va.getName());
                        if (!isSerializerArg && !isTypeArg && serialzerArguments != null) {
                            for (Argument serializerArgument : serialzerArguments) {
                                if (serializerArgument.getName().equals(va.getName())) {
                                    isSerializerArg = true;
                                    break;
                                }
                            }
                        }
                        if (isSerializerArg) {
                            sb.append(va.getName());
                            if(va.getChild() != null) {
                                sb.append(".");
                                appendVariableExpressionRest(sb, baseType, va.getChild());
                            }
                        } else if (isTypeArg) {
                            String part = va.getChild().getName();
                            switch (part) {
                                case "name":
                                    sb.append("\"").append(field.getTypeName()).append("\"");
                                    break;
                                case "length":
                                    sb.append("\"").append(((SimpleTypeReference) field).getSizeInBits()).append("\"");
                                    break;
                                case "encoding":
                                    if(!(field instanceof TypedField)) {
                                        throw new FreemarkerException("'encoding' only supported for typed fields.");
                                    }
                                    TypedField typedField = (TypedField) field;
                                    if(!(typedField.getType() instanceof StringTypeReference)) {
                                        throw new FreemarkerException("Can only access 'encoding' for string types.");
                                    }
                                    StringTypeReference stringTypeReference = (StringTypeReference) typedField.getType();
                                    String encoding = stringTypeReference.getEncoding();
                                    // Cut off the single quotes.
                                    encoding = encoding.substring(1, encoding.length() - 1);
                                    sb.append("\"").append(encoding).append("\"");
                                    break;
                            }
                        } else {
                            sb.append(toVariableSerializationExpression(baseType, field, va, null));
                        }
                    } else if (arg instanceof StringLiteral) {
                        sb.append(((StringLiteral) arg).getValue());
                    }
                    firstArg = false;
                }
                sb.append(")");
            }
            return sb.toString();
        }
        // If we are accessing implicit fields, we need to rely on a local variable instead.
        if (isVariableLiteralImplicitField(vl)) {
            return toSerializationExpression(getThisTypeDefinition(), getReferencedImplicitField(vl), getReferencedImplicitField(vl).getSerializeExpression(), serialzerArguments);
        }
        // The synthetic checksumRawData is a local field and should not be accessed as bean property.
        boolean isSerializerArg = "checksumRawData".equals(vl.getName()) || "_value".equals(vl.getName()) || "element".equals(vl.getName()) || "size".equals(vl.getName());
        boolean isTypeArg = "_type".equals(vl.getName());
        if (!isSerializerArg && !isTypeArg && serialzerArguments != null) {
            for (Argument serializerArgument : serialzerArguments) {
                if (serializerArgument.getName().equals(vl.getName())) {
                    isSerializerArg = true;
                    break;
                }
            }
        }
        if (isSerializerArg) {
            StringBuilder sb = new StringBuilder(vl.getName());
            if(vl.getChild() != null) {
                sb.append(".");
                appendVariableExpressionRest(sb, baseType, vl.getChild());
            }
            return sb.toString();
        } else if (isTypeArg) {
            String part = vl.getChild().getName();
            switch (part) {
                case "name":
                    return "\"" + field.getTypeName() + "\"";
                case "length":
                    return "\"" + ((SimpleTypeReference) field).getSizeInBits() + "\"";
                case "encoding":
                    if(!(field instanceof TypedField)) {
                        throw new FreemarkerException("'encoding' only supported for typed fields.");
                    }
                    TypedField typedField = (TypedField) field;
                    if(!(typedField.getType() instanceof StringTypeReference)) {
                        throw new FreemarkerException("Can only access 'encoding' for string types.");
                    }
                    StringTypeReference stringTypeReference = (StringTypeReference) typedField.getType();
                    String encoding = stringTypeReference.getEncoding();
                    // Cut off the single quotes.
                    encoding = encoding.substring(1, encoding.length() - 1);
                    return "\"" + encoding + "\"";
                default:
                    return "";
            }
        } else {
            // If it wasn't an enum, treat it as a normal property.
            if (vl.getName().equals("lengthInBits")) {
                StringBuilder sb = new StringBuilder(getCTypeName(baseType.getName()));
                sb.append("_length_in_bits(_message)");
                return sb.toString();
            } else if (vl.getChild() != null && "length".equals(vl.getChild().getName())) {
                StringBuilder sb = new StringBuilder("");
                sb.append("sizeof(_message->" + camelCaseToSnakeCase(vl.getName()) + ")");
                return sb.toString();
            } else {
                StringBuilder sb = new StringBuilder("_message->");
                // If this is a property of a sub-type, add the sub-type name to the property.
                if(baseType != getThisTypeDefinition()) {
                    sb.append(camelCaseToSnakeCase(baseType.getName())).append("_");
                }

                // If this expression references enum constants we need to do things differently
                final Optional<TypeReference> typeReferenceForProperty =
                    getTypeReferenceForProperty((ComplexTypeDefinition) baseType, vl.getName());
                if(typeReferenceForProperty.isPresent()) {
                    final TypeReference typeReference = typeReferenceForProperty.get();
                    if(typeReference instanceof ComplexTypeReference) {
                        final TypeDefinition typeDefinitionForTypeReference =
                            getTypeDefinitionForTypeReference(typeReference);
                        if ((typeDefinitionForTypeReference instanceof EnumTypeDefinition) && (vl.getChild() != null)){
                            sb.append(camelCaseToSnakeCase(vl.getName()));
                            return getCTypeName(typeDefinitionForTypeReference.getName()) +
                                "_get_" + camelCaseToSnakeCase(vl.getChild().getName()) +
                                "(" + sb.toString() + ")";
                        }
                    }
                }
                // If it wasn't an enum, treat it as a normal property.
                appendVariableExpressionRest(sb, baseType, vl);
                return sb.toString();
            }
        }
    }

    private void appendVariableExpressionRest(StringBuilder sb, TypeDefinition baseType, VariableLiteral vl) {
        if(vl.isIndexed()) {
            sb.insert(0, "plc4c_utils_list_get_value(");
            sb.append(camelCaseToSnakeCase(vl.getName()));
            sb.append(", ").append(vl.getIndex()).append(")");
        } else {
            sb.append(camelCaseToSnakeCase(vl.getName()));
        }
        // Suppress any "lengthInBytes" properties as these are handled differently in C
        if((vl.getChild() != null) && !vl.getChild().getName().equals("lengthInBytes")) {
            sb.append(".");
            appendVariableExpressionRest(sb, baseType, vl.getChild());
        }
    }

    public int getNumBits(SimpleTypeReference simpleTypeReference) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT: {
                return 1;
            }
            case UINT:
            case INT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                return integerTypeReference.getSizeInBits();
            }
            case FLOAT: {
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                return floatTypeReference.getSizeInBits();
            }
            case STRING: {
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return stringTypeReference.getSizeInBits();
            }
            default: {
                return 0;
            }
        }
    }

    public String getLengthInBitsFunctionNameForComplexTypedField(Field field) {
        if(field instanceof TypedField) {
            TypedField typedField = (TypedField) field;
            final TypeReference typeReference = typedField.getType();
            if(typeReference instanceof ComplexTypeReference) {
                ComplexTypeReference complexTypeReference = (ComplexTypeReference) typeReference;
                return getCTypeName(complexTypeReference.getName()) + "_length_in_bits";
            } else {
                throw new FreemarkerException("lengthInBits functions only exist for complex types");
            }
        } else {
            throw new FreemarkerException("lengthInBits functions only exist for TypedFields");
        }
    }

    public String getEnumExpression(String expression) {
        String enumName = expression.substring(0, expression.indexOf('.'));
        String enumConstant = expression.substring(expression.indexOf('.') + 1);
        return getCTypeName(enumName) + "_" + enumConstant;
    }

}
