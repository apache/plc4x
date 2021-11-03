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
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.Tracer;
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
     * In addition, it appends a prefix for the protocol name and the output flavor.
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
        if (baseType != thisType) {
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
            return snakeCase.substring(1);
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
        if (!(field.isTypedField())) {
            throw new FreemarkerException("Field " + field + " is not a TypedField");
        }
        // If this is an array with variable length, then we have to use our "plc4c_list" to store the data.
        if (field.asArrayField().map(ArrayField::getLoopExpression).map(Term::isFixedValueExpression).orElse(false)) {
            return "plc4c_list";
        }
        TypedField typedField = field.asTypedField().orElseThrow(IllegalStateException::new);
        TypeReference typeReference = typedField.getType();
        if (typeReference.isComplexTypeReference()) {
            final TypeDefinition typeDefinition = getTypeDefinitionForTypeReference(typeReference);
            if (typeDefinition instanceof DataIoTypeDefinition) {
                return "plc4c_data*";
            }
        }
        return getLanguageTypeNameForTypeReference(typeReference);
    }

    public Map<ConstField, ComplexTypeDefinition> getAllConstFields() {
        Map<ConstField, ComplexTypeDefinition> constFields = new LinkedHashMap<>();
        ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) this.thisType;
        complexTypeDefinition.getConstFields()
            .forEach(constField -> constFields.put(constField, complexTypeDefinition));
        complexTypeDefinition.getSwitchField()
            .map(SwitchField::getCases)
            .ifPresent(discriminatedComplexTypeDefinitions ->
                discriminatedComplexTypeDefinitions.forEach(switchCase ->
                    switchCase.getConstFields().forEach(constField -> constFields.put(constField, switchCase))
                )
            );
        return constFields;
    }

    /**
     * If a property references a complex type in an argument, we need to pass that as a pointer,
     * same with optional fields.
     *
     * @param typeDefinition type that contains the property or attribute.
     * @param propertyName   name of the property or attribute
     * @return true if the access needs to be using pointers
     */
    public boolean requiresPointerAccess(ComplexTypeDefinition typeDefinition, String propertyName) {
        final Optional<NamedField> namedFieldOptional = typeDefinition.getFields().stream()
            .filter(NamedField.class::isInstance)
            .map(NamedField.class::cast)
            .filter(namedField -> namedField.getName().equals(propertyName))
            .findFirst();
        // If the property name refers to a field, check if it's an optional field.
        // If it is, pointer access is required, if not, it's not.
        if (namedFieldOptional.isPresent()) {
            return namedFieldOptional
                .filter(TypedField.class::isInstance)
                .map(TypedField.class::cast)
                .map(typedField -> !(typedField.isEnumField()) && typedField.getType().isComplexTypeReference())
                .orElse(false);
        }
        return typeDefinition.getParserArguments()
            .orElse(Collections.emptyList())
            .stream()
            .filter(argument -> argument.getName().equals(propertyName)).findFirst()
            // If the property name refers to a parser argument, as soon as it's a complex type,
            // pointer access is required.
            .map(Argument::getType)
            .filter(ComplexTypeReference.class::isInstance)
            .isPresent();
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
                    int sizeInBits = floatTypeReference.getSizeInBits();
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
                case VSTRING:
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
        if (!(field instanceof ArrayField)) {
            return "";
        }
        ArrayField arrayField = (ArrayField) field;
        if (arrayField.getLoopType() != ArrayField.LoopType.COUNT) {
            return "";
        }
        Term countTerm = arrayField.getLoopExpression();
        if (countTerm == null || !countTerm.isFixedValueExpression()) {
            return "";
        }
        int evaluatedCount = evaluateFixedValueExpression(countTerm);
        return "[" + evaluatedCount + "]";
    }

    /**
     * Ge the type-size suffix in case of simple types.
     *
     * @param field the field we want to get the type-size for
     * @return a type-size string for the given field or an empty string if this does not apply
     */
    public String getTypeSizeForField(TypedField field) {
        TypeReference typeReference = field.getType();
        if (!(typeReference instanceof SimpleTypeReference)) {
            return "";
        }
        SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return " : 1";
            case BYTE:
                return " : 8";
            case UINT:
            case INT:
                // If the bit-size is exactly one of the built-in type-sizes, omit the suffix.
                if ((simpleTypeReference.getSizeInBits() == 8) || (simpleTypeReference.getSizeInBits() == 16) ||
                    (simpleTypeReference.getSizeInBits() == 32) || (simpleTypeReference.getSizeInBits() == 64)) {
                    return "";
                }
                return " : " + simpleTypeReference.getSizeInBits();
            case FLOAT:
            case UFLOAT:
                // If the bit-size is exactly one of the built-in type-sizes, omit the suffix.
                if ((simpleTypeReference.getSizeInBits() == 32) || (simpleTypeReference.getSizeInBits() == 64)) {
                    return "";
                }
                return " : " + simpleTypeReference.getSizeInBits();
            case STRING:
            case VSTRING:
            case TIME:
            case DATE:
            case DATETIME:
                return "";
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
                case VSTRING:
                    return "\"" + valueString + "\"";
            }
        }
        return valueString;
    }

    public String escapeEnumValue(TypeReference typeReference, String valueString) {
        // Currently, the only case in which here complex type references are used are when referencing enum constants.
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
                } else if (floatTypeReference.getSizeInBits() <= 64) {
                    return "plc4c_spi_read_double(readBuffer, " + floatTypeReference.getSizeInBits() + ", (double*) " + valueString + ")";
                }
                throw new FreemarkerException("Unsupported float type with " + floatTypeReference.getSizeInBits() + " bits");
            case STRING:
            case VSTRING:
                final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                if (!(encodingTerm instanceof StringLiteral)) {
                    throw new RuntimeException("Encoding must be a quoted string value");
                }
                String encoding = ((StringLiteral) encodingTerm).getValue();
                String length = Integer.toString(simpleTypeReference.getSizeInBits());
                return "plc4c_spi_read_string(readBuffer, " + length + ", \"" +
                    encoding + "\"" + ", (char**) " + valueString + ")";
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
                } else if (floatTypeReference.getSizeInBits() <= 64) {
                    return "plc4c_spi_write_double(writeBuffer, " + floatTypeReference.getSizeInBits() + ", " + fieldName + ")";
                }
                throw new FreemarkerException("Unsupported float type with " + floatTypeReference.getSizeInBits() + " bits");
            case STRING:
            case VSTRING:
                final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                if (!(encodingTerm instanceof StringLiteral)) {
                    throw new RuntimeException("Encoding must be a quoted string value");
                }
                String encoding = ((StringLiteral) encodingTerm).getValue();
                String length = Integer.toString(simpleTypeReference.getSizeInBits());
                return "plc4c_spi_write_string(writeBuffer, " + length + ", \"" +
                    encoding + "\", " + fieldName + ")";
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
                case UINT:
                case INT:
                    return "0";
                case FLOAT:
                case UFLOAT:
                    return "0.0f";
                case STRING:
                case VSTRING:
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
        if ("BigInteger".equals(languageTypeName)) {
            return "BigInteger.valueOf(" + reservedField.getReferenceValue() + ")";
        } else {
            return "(" + languageTypeName + ") " + reservedField.getReferenceValue();
        }
    }

    public String toSpecialParseExpression(TypeDefinition baseType, Field field, Term term, List<Argument> parserArguments) {
        return toParseExpression(baseType, field, term, parserArguments);
    }

    public String toParseExpression(TypeDefinition baseType, Field field, Term term, List<Argument> parserArguments) {
        return toExpression(baseType, field, term, variableLiteral -> toVariableParseExpression(baseType, field, variableLiteral, parserArguments));
    }

    public String toSerializationExpression(TypeDefinition baseType, Field field, Term term, List<Argument> parserArguments) {
        return toExpression(baseType, field, term, variableLiteral -> toVariableSerializationExpression(baseType, field, variableLiteral, parserArguments));
    }

    private String toExpression(TypeDefinition baseType, Field field, Term term, Function<VariableLiteral, String> variableExpressionGenerator) {
        Tracer tracer = Tracer.start("toExpression");
        if (term == null) {
            return tracer + "";
        }
        if (term instanceof Literal) {
            return toLiteralExpression(baseType, term, variableExpressionGenerator, tracer);
        } else if (term instanceof UnaryTerm) {
            return toUnaryTermExpression(baseType, field, (UnaryTerm) term, variableExpressionGenerator, tracer);
        } else if (term instanceof BinaryTerm) {
            return toBinaryTermExpression(baseType, field, (BinaryTerm) term, variableExpressionGenerator, tracer);
        } else if (term instanceof TernaryTerm) {
            return toTernaryTermExpression(baseType, field, (TernaryTerm) term, variableExpressionGenerator, tracer);
        } else {
            throw new FreemarkerException("Unsupported Term type " + term.getClass().getName());
        }
    }

    private String toUnaryTermExpression(TypeDefinition baseType, Field field, UnaryTerm unaryTerm, Function<VariableLiteral, String> variableExpressionGenerator, Tracer tracer) {
        tracer = tracer.dive("unary term instanceOf");
        Term a = unaryTerm.getA();
        switch (unaryTerm.getOperation()) {
            case "!":
                tracer = tracer.dive("case !");
                return tracer + "!(" + toExpression(baseType, field, a, variableExpressionGenerator) + ")";
            case "-":
                tracer = tracer.dive("case -");
                return tracer + "-(" + toExpression(baseType, field, a, variableExpressionGenerator) + ")";
            case "()":
                tracer = tracer.dive("case ()");
                return tracer + "(" + toExpression(baseType, field, a, variableExpressionGenerator) + ")";
            default:
                throw new FreemarkerException("Unsupported unary operation type " + unaryTerm.getOperation());
        }
    }

    private String toBinaryTermExpression(TypeDefinition baseType, Field field, BinaryTerm binaryTerm, Function<VariableLiteral, String> variableExpressionGenerator, Tracer tracer) {
        tracer = tracer.dive("binary term instanceOf");
        Term a = binaryTerm.getA();
        Term b = binaryTerm.getB();
        String operation = binaryTerm.getOperation();
        if ("^".equals(operation)) {
            tracer = tracer.dive("^");
            return tracer + "Math.pow((" + toExpression(baseType, field, a, variableExpressionGenerator) + "), (" + toExpression(baseType, field, b, variableExpressionGenerator) + "))";
        }
        return tracer + "(" + toExpression(baseType, field, a, variableExpressionGenerator) + ") " + operation + " (" + toExpression(baseType, field, b, variableExpressionGenerator) + ")";
    }

    private String toTernaryTermExpression(TypeDefinition baseType, Field field, TernaryTerm ternaryTerm, Function<VariableLiteral, String> variableExpressionGenerator, Tracer tracer) {
        tracer = tracer.dive("ternary term instanceOf");
        if ("if".equals(ternaryTerm.getOperation())) {
            Term a = ternaryTerm.getA();
            Term b = ternaryTerm.getB();
            Term c = ternaryTerm.getC();
            return tracer + "((" + toExpression(baseType, field, a, variableExpressionGenerator) + ") ? " + toExpression(baseType, field, b, variableExpressionGenerator) + " : " + toExpression(baseType, field, c, variableExpressionGenerator) + ")";
        } else {
            throw new FreemarkerException("Unsupported ternary operation type " + ternaryTerm.getOperation());
        }
    }

    private String toLiteralExpression(TypeDefinition baseType, Term term, Function<VariableLiteral, String> variableExpressionGenerator, Tracer tracer) {
        tracer = tracer.dive("literal term instanceOf");
        if (term instanceof NullLiteral) {
            tracer = tracer.dive("null literal instanceOf");
            return tracer + "NULL";
        } else if (term instanceof BooleanLiteral) {
            tracer = tracer.dive("boolean literal instanceOf");
            return tracer + Boolean.toString(((BooleanLiteral) term).getValue());
        } else if (term instanceof NumericLiteral) {
            tracer = tracer.dive("numeric literal instanceOf");
            return tracer + ((NumericLiteral) term).getNumber().toString();
        } else if (term instanceof StringLiteral) {
            tracer = tracer.dive("string literal instanceOf");
            return tracer + toStringLiteralExpression((StringLiteral) term);
        } else if (term instanceof VariableLiteral) {
            tracer = tracer.dive("variable literal instanceOf");
            VariableLiteral variableLiteral = (VariableLiteral) term;
            if ("curPos".equals(variableLiteral.getName())) {
                return "(plc4c_spi_read_get_pos(readBuffer) - startPos)";
            }
            return tracer + toVariableLiteralExpression(baseType, (VariableLiteral) term, variableExpressionGenerator, tracer);
        } else {
            throw new FreemarkerException("Unsupported Literal type " + term.getClass().getName());
        }
    }

    private String toStringLiteralExpression(StringLiteral term) {
        return "\"" + term.getValue() + "\"";
    }

    private String toVariableLiteralExpression(TypeDefinition baseType, VariableLiteral variableLiteral, Function<VariableLiteral, String> variableExpressionGenerator, Tracer tracer) {
        tracer = tracer.dive("toVariableLiteralExpression");
        if (variableLiteral.contains("lengthInBytes")) {
            tracer = tracer.dive("lengthInBytesContained");
            TypeDefinition lengthType;
            String lengthExpression;
            if (variableLiteral.getName().equals("lengthInBytes")) {
                tracer = tracer.dive("lengthInBytes contained in variable name");
                lengthType = (baseType.getParentType() == null) ? baseType : (ComplexTypeDefinition) baseType.getParentType();
                lengthExpression = "_message";
            } else {
                final Optional<TypeReference> typeReferenceForProperty = getTypeReferenceForProperty((ComplexTypeDefinition) baseType, variableLiteral.getName());
                if (!typeReferenceForProperty.isPresent()) {
                    throw new FreemarkerException("Unknown type for property " + variableLiteral.getName());
                }
                lengthType = getTypeDefinitionForTypeReference(typeReferenceForProperty.get());
                lengthExpression = variableExpressionGenerator.apply(variableLiteral);
            }
            return tracer + getCTypeName(lengthType.getName()) + "_length_in_bytes(" + lengthExpression + ")";
        } else if (variableLiteral.getName().equals("lastItem")) {
            tracer = tracer.dive("lastitem");
            return tracer + "lastItem";
            // If this literal references an Enum type, then we have to output it differently.
        } else if (getTypeDefinitions().get(variableLiteral.getName()) instanceof EnumTypeDefinition) {
            tracer = tracer.dive("enum type definition");
            return tracer + getCTypeName(variableLiteral.getName()) + "_" + variableLiteral.getChild().map(VariableLiteral::getName).orElseThrow(() -> new FreemarkerException("child required"));
        } else {
            return tracer + variableExpressionGenerator.apply(variableLiteral);
        }
    }

    public String toVariableParseExpression(TypeDefinition baseType, Field field, VariableLiteral variableLiteral, List<Argument> parserArguments) {
        Tracer tracer = Tracer.start("toVariableParseExpression");
        if ("CAST".equals(variableLiteral.getName())) {
            return toCastVariableParseExpression(baseType, field, variableLiteral, parserArguments, tracer);
        }
        if ("STATIC_CALL".equals(variableLiteral.getName())) {
            return toStaticCallVariableParseExpression(baseType, field, variableLiteral, parserArguments, tracer);
        }
        if (variableLiteral.getName().equals(variableLiteral.getName().toUpperCase())) {
            return toUpperCaseVariableParseExpression(baseType, field, variableLiteral, parserArguments, tracer);
        } else if ("readBuffer".equals(variableLiteral.getName())) {
            return toReadBufferVariableParseExpression(baseType, variableLiteral, tracer);
        } else if ("writeBuffer".equals(variableLiteral.getName())) {
            return toWriteBufferVariableParseExpression(baseType, variableLiteral, tracer);
        } else if ("_type".equals(variableLiteral.getName())) {
            return toTypeVariableParseExpression(field, variableLiteral, tracer);
        }

        final String name = variableLiteral.getName();

        // In case of DataIo types, we'll just check the arguments.
        if (baseType instanceof DataIoTypeDefinition) {
            if (baseType.getParserArguments().isPresent()) {
                for (Argument parserArgument : baseType.getParserArguments().get()) {
                    if (parserArgument.getName().equals(name)) {
                        return name;
                    }
                }
            }
        }

        // Try to find the type of the addressed property.
        Optional<TypeReference> propertyTypeOptional =
            getTypeReferenceForProperty((ComplexTypeDefinition) baseType, name);

        // If we couldn't find the type, we didn't find the property.
        if (!propertyTypeOptional.isPresent()) {
            final List<Argument> arguments = baseType.getAllParserArguments().orElse(Collections.emptyList());
            for (Argument argument : arguments) {
                if(argument.getName().equals(name)) {
                    propertyTypeOptional = Optional.of(argument.getType());
                }
            }
            if(!propertyTypeOptional.isPresent()) {
                throw new FreemarkerException("Could not find property with name '" + name + "' in type " + baseType.getName());
            }
        }

        final TypeReference propertyType = propertyTypeOptional.get();

        // If it's a simple field, there is no sub-type to access.
        if (propertyType instanceof SimpleTypeReference) {
            if (variableLiteral.getChild().isPresent()) {
                throw new FreemarkerException("Simple property '" + name + "' doesn't have child properties.");
            }
            return name;
        }

        // If it references a complex, type we need to get that type's definition first.
        final TypeDefinition propertyTypeDefinition = getTypeDefinitions().get(((ComplexTypeReference) propertyType).getName());
        // If we're not accessing any child property, no need to handle anything special.
        if (!variableLiteral.getChild().isPresent()) {
            return name;
        }
        // If there is a child we need to check if this is a discriminator property.
        // As discriminator properties are not real properties, but saved in the static metadata
        // of a type, we need to generate a different access pattern.
        if (propertyTypeDefinition instanceof ComplexTypeDefinition) {
            final Optional<DiscriminatorField> discriminatorFieldOptional = ((ComplexTypeDefinition) propertyTypeDefinition).getFields().stream().filter(
                curField -> curField instanceof DiscriminatorField).map(curField -> (DiscriminatorField) curField).filter(
                discriminatorField -> discriminatorField.getName().equals(variableLiteral.getChild().get().getName())).findFirst();
            // If child references a discriminator field of the type we found, we have to escape it.
            if (discriminatorFieldOptional.isPresent()) {
                return getCTypeName(propertyTypeDefinition.getName()) + "_get_discriminator(" + name + "->_type)." + variableLiteral.getChild().get().getName();
            }
        }
        // Handling enum properties in C is a little more tricky as we have to use the enum value
        // and pass this to a function that returns the desired property value.
        else if (propertyTypeDefinition instanceof EnumTypeDefinition) {
            return getCTypeName(propertyTypeDefinition.getName()) +
                "_get_" + camelCaseToSnakeCase(variableLiteral.getChild().get().getName()) +
                "(*" + variableLiteral.getName() + ")";
        }
        // Else ... generate a simple access path.
        StringBuilder sb = new StringBuilder(variableLiteral.getName());
        if (variableLiteral.getChild().isPresent()) {
            sb.append(".");
            appendVariableExpressionRest(sb, baseType, variableLiteral.getChild().get());
        }
        return sb.toString();
    }

    private String toTypeVariableParseExpression(Field field, VariableLiteral variableLiteral, Tracer tracer) {
        tracer = tracer.dive("type");
        if (variableLiteral.getChild().isPresent() && "encoding".equals(variableLiteral.getChild().get().getName()) && (field instanceof TypedField) && ((((TypedField) field).getType() instanceof StringTypeReference) || (((TypedField) field).getType() instanceof VstringTypeReference))) {
            // TODO: replace with map join
            final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
            if (!(encodingTerm instanceof StringLiteral)) {
                throw new RuntimeException("Encoding must be a quoted string value");
            }
            String encoding = ((StringLiteral) encodingTerm).getValue();
            return tracer + "\"" + encoding + "\"";
        } else {
            throw new FreemarkerException("_type is currently pretty much hard-coded for some use cases, please check CLanguageTemplateHelper.toVariableParseExpression");
        }
    }

    private String toReadBufferVariableParseExpression(TypeDefinition baseType, VariableLiteral variableLiteral, Tracer tracer) {
        tracer = tracer.dive("readbuffer");
        StringBuilder sb = new StringBuilder("readBuffer");
        if (variableLiteral.getChild().isPresent()) {
            // TODO: replace with map join
            sb.append(".");
            appendVariableExpressionRest(sb, baseType, variableLiteral.getChild().get());
        }
        return tracer + sb.toString();
    }

    /*
     * Any name that is full upper-case is considered a function call.
     * These are generally defined in the spi file evaluation_helper.c.
     * All should have a name prefix "plc4c_spi_evaluation_helper_".
     */
    private String toUpperCaseVariableParseExpression(TypeDefinition baseType, Field field, VariableLiteral variableLiteral, List<Argument> parserArguments, Tracer tracer) {
        tracer = tracer.dive("UPPER");
        StringBuilder sb = new StringBuilder("plc4c_spi_evaluation_helper_" + variableLiteral.getName().toLowerCase());
        if (variableLiteral.getArgs().isPresent()) {
            // TODO: replace with map join
            sb.append("(");
            boolean firstArg = true;
            for (Term arg : variableLiteral.getArgs().get()) {
                if (!firstArg) {
                    sb.append(", ");
                }
                sb.append(toParseExpression(baseType, field, arg, parserArguments));
                firstArg = false;
            }
            sb.append(")");
        }
        if (variableLiteral.getIndex() != VariableLiteral.NO_INDEX) {
            sb.append("[").append(variableLiteral.getIndex()).append("]");
        }
        if (variableLiteral.getChild().isPresent()) {
            sb.append(".");
            appendVariableExpressionRest(sb, baseType, variableLiteral.getChild().get());
        }
        return tracer + sb.toString();
    }

    private String toCastVariableParseExpression(TypeDefinition baseType, Field field, VariableLiteral variableLiteral, List<Argument> parserArguments, Tracer tracer) {
        tracer = tracer.dive("CAST");
        if ((!variableLiteral.getArgs().isPresent()) || (variableLiteral.getArgs().get().size() != 2)) {
            throw new FreemarkerException("A CAST expression expects exactly two arguments.");
        }
        List<Term> args = variableLiteral.getArgs().get();
        final VariableLiteral sourceTerm = (VariableLiteral) args.get(0);
        final VariableLiteral typeTerm = (VariableLiteral) variableLiteral.getArgs().get().get(1);
        final TypeDefinition castType = getTypeDefinitions().get(typeTerm.getName());
        // If we're casting to a sub-type of a discriminated value, we got to cast to the parent
        // type instead and add the name of the sub-type as prefix to the property we're trying to
        // access next.
        StringBuilder sb = new StringBuilder();
        sb.append("((");
        if (castType.getParentType() != null) {
            sb.append(getCTypeName(castType.getParentType().getName()));
        } else {
            sb.append(getCTypeName(castType.getName()));
        }
        sb.append("*) (");
        sb.append(toVariableParseExpression(baseType, field, sourceTerm, parserArguments)).append("))");
        if (variableLiteral.getChild().isPresent()) {
            if (castType.getParentType() != null) {
                // Change the name of the property to contain the sub-type-prefix.
                sb.append("->").append(camelCaseToSnakeCase(castType.getName())).append("_");
                appendVariableExpressionRest(sb, baseType, variableLiteral.getChild().get());
            } else {
                sb.append("->");
                appendVariableExpressionRest(sb, castType, variableLiteral.getChild().get());
            }
        }
        return tracer + sb.toString();
    }

    /*
     * STATIC_CALL implies that driver specific static logic should be called
     */
    private String toStaticCallVariableParseExpression(TypeDefinition baseType, Field field, VariableLiteral variableLiteral, List<Argument> parserArguments, Tracer tracer) {
        tracer = tracer.dive("STATIC_CALL");
        List<Term> terms = variableLiteral.getArgs().orElseThrow(() -> new FreemarkerException("'STATIC_CALL' needs at least one args"));
        String functionName = terms.get(0).asLiteral()
            .orElseThrow(()-> new FreemarkerException("Expecting the first argument of a 'STATIC_CALL' to be a Literal"))
            .asStringLiteral()
            .orElseThrow(()-> new FreemarkerException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral"))
            .getValue();
        // But to make the function name unique, well add the driver prefix to it.
        StringBuilder sb = new StringBuilder(getCTypeName(functionName));
        if (terms.size() > 1) {
            sb.append("(");
            boolean firstArg = true;
            for (int i = 1; i < terms.size(); i++) {
                Term arg = terms.get(i);
                if (!firstArg) {
                    sb.append(", ");
                }
                sb.append(toParseExpression(baseType, field, arg, parserArguments));
                firstArg = false;
            }
            sb.append(")");
        }
        return tracer + sb.toString();
    }

    private String toWriteBufferVariableParseExpression(TypeDefinition baseType, VariableLiteral variableLiteral, Tracer tracer) {
        tracer = tracer.dive("writeBuffer");
        StringBuilder sb = new StringBuilder("writeBuffer");
        if (variableLiteral.getChild().isPresent()) {
            sb.append(".");
            appendVariableExpressionRest(sb, baseType, variableLiteral.getChild().get());
        }
        return tracer + sb.toString();
    }

    private String toVariableSerializationExpression(TypeDefinition baseType, Field field, VariableLiteral variableLiteral, List<Argument> serializerArguments) {
        Tracer tracer = Tracer.start("toVariableSerializationExpression");
        if ("STATIC_CALL".equals(variableLiteral.getName())) {
            return toStaticCallSerializationExpression(baseType, field, serializerArguments, variableLiteral, tracer);
        }
        // All uppercase names are not fields, but utility methods.
        if (variableLiteral.getName().equals(variableLiteral.getName().toUpperCase())) {
            return toUppercaseSerializationExpression(baseType, field, serializerArguments, variableLiteral, tracer);
        }
        // If we are accessing implicit fields, we need to rely on a local variable instead.
        if (isVariableLiteralImplicitField(variableLiteral)) {
            tracer = tracer.dive("is variable implicit field");
            return tracer + toSerializationExpression(thisType, getReferencedImplicitField(variableLiteral), getReferencedImplicitField(variableLiteral).getSerializeExpression(), serializerArguments);
        }
        // The synthetic checksumRawData is a local field and should not be accessed as bean property.
        boolean isSerializerArg = "checksumRawData".equals(variableLiteral.getName()) || "_value".equals(variableLiteral.getName()) || "element".equals(variableLiteral.getName()) || "size".equals(variableLiteral.getName());
        boolean isTypeArg = "_type".equals(variableLiteral.getName());
        if (!isSerializerArg && !isTypeArg && serializerArguments != null) {
            for (Argument serializerArgument : serializerArguments) {
                if (serializerArgument.getName().equals(variableLiteral.getName())) {
                    isSerializerArg = true;
                    break;
                }
            }
        }
        if (isSerializerArg) {
            tracer = tracer.dive("is serializer arg");
            StringBuilder sb = new StringBuilder(variableLiteral.getName());
            if (variableLiteral.getChild().isPresent()) {
                sb.append(".");
                appendVariableExpressionRest(sb, baseType, variableLiteral.getChild().get());
            }
            return tracer + sb.toString();
        }
        if (isTypeArg) {
            tracer = tracer.dive("is type arg");
            String part = variableLiteral.getChild().get().getName();
            switch (part) {
                case "name":
                    tracer = tracer.dive("name");
                    return tracer + "\"" + field.getTypeName() + "\"";
                case "length":
                    tracer = tracer.dive("length");
                    return tracer + "\"" + ((SimpleTypeReference) field).getSizeInBits() + "\"";
                case "encoding":
                    tracer = tracer.dive("encoding");
                    if (!(field instanceof TypedField)) {
                        throw new FreemarkerException("'encoding' only supported for typed fields.");
                    }
                    TypedField typedField = (TypedField) field;
                    if (!(typedField.getType() instanceof StringTypeReference)) {
                        throw new FreemarkerException("Can only access 'encoding' for string types.");
                    }
                    final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                    if (!(encodingTerm instanceof StringLiteral)) {
                        throw new RuntimeException("Encoding must be a quoted string value");
                    }
                    String encoding = ((StringLiteral) encodingTerm).getValue();
                    return tracer + "\"" + encoding + "\"";
                default:
                    return tracer + "";
            }
        }
        // If it wasn't an enum, treat it as a normal property.
        if (variableLiteral.getName().equals("lengthInBits")) {
            tracer = tracer.dive("is length in bits");
            return tracer + getCTypeName(baseType.getName()) + "_length_in_bits(_message)";
        }
        if (variableLiteral.getChild().isPresent() && "length".equals(variableLiteral.getChild().get().getName())) {
            tracer = tracer.dive("is length");
            return tracer + "sizeof(_message->" + camelCaseToSnakeCase(variableLiteral.getName()) + ")";
        }
        StringBuilder sb = new StringBuilder("_message->");
        // If this is a property of a sub-type, add the sub-type name to the property.
        if (baseType != thisType) {
            tracer = tracer.dive("this is not this type");
            sb.append(camelCaseToSnakeCase(baseType.getName())).append("_");
        }

        // If this expression references enum constants we need to do things differently
        final Optional<TypeReference> typeReferenceForProperty =
            getTypeReferenceForProperty((ComplexTypeDefinition) baseType, variableLiteral.getName());
        if (typeReferenceForProperty.isPresent()) {
            final TypeReference typeReference = typeReferenceForProperty.get();
            if (typeReference instanceof ComplexTypeReference) {
                final TypeDefinition typeDefinitionForTypeReference =
                    getTypeDefinitionForTypeReference(typeReference);
                if ((typeDefinitionForTypeReference instanceof EnumTypeDefinition) && (variableLiteral.getChild().isPresent())) {
                    tracer = tracer.dive("is enum type definition");
                    sb.append(camelCaseToSnakeCase(variableLiteral.getName()));
                    return tracer + getCTypeName(typeDefinitionForTypeReference.getName()) + "_get_" + camelCaseToSnakeCase(variableLiteral.getChild().get().getName()) +
                        "(" + sb + ")";
                }
            }
        }
        // If it wasn't an enum, treat it as a normal property.
        appendVariableExpressionRest(sb, baseType, variableLiteral);
        return tracer + sb.toString();
    }

    private String toUppercaseSerializationExpression(TypeDefinition baseType, Field field, List<Argument> serializerArguments, VariableLiteral vl, Tracer tracer) {
        tracer = tracer.dive("UPPER_CASE");
        StringBuilder sb = new StringBuilder("plc4c_spi_evaluation_helper_" + vl.getName().toLowerCase());
        if (vl.getArgs().isPresent()) {
            sb.append("(");
            boolean firstArg = true;
            for (Term arg : vl.getArgs().get()) {
                if (!firstArg) {
                    sb.append(", ");
                }

                if (arg instanceof VariableLiteral) {
                    VariableLiteral va = (VariableLiteral) arg;
                    boolean isSerializerArg = "io".equals(va.getName());
                    boolean isTypeArg = "_type".equals(va.getName());
                    if (!isSerializerArg && !isTypeArg && serializerArguments != null) {
                        for (Argument serializerArgument : serializerArguments) {
                            if (serializerArgument.getName().equals(va.getName())) {
                                isSerializerArg = true;
                                break;
                            }
                        }
                    }
                    if (isSerializerArg) {
                        sb.append(va.getName());
                        if (va.getChild().isPresent()) {
                            sb.append(".");
                            appendVariableExpressionRest(sb, baseType, va.getChild().get());
                        }
                    } else if (isTypeArg) {
                        String part = va.getChild().get().getName();
                        switch (part) {
                            case "name":
                                sb.append("\"").append(field.getTypeName()).append("\"");
                                break;
                            case "length":
                                sb.append("\"").append(((SimpleTypeReference) field).getSizeInBits()).append("\"");
                                break;
                            case "encoding":
                                if (!(field instanceof TypedField)) {
                                    throw new FreemarkerException("'encoding' only supported for typed fields.");
                                }
                                TypedField typedField = (TypedField) field;
                                if (!(typedField.getType() instanceof StringTypeReference)) {
                                    throw new FreemarkerException("Can only access 'encoding' for string types.");
                                }
                                final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                                if (!(encodingTerm instanceof StringLiteral)) {
                                    throw new RuntimeException("Encoding must be a quoted string value");
                                }
                                String encoding = ((StringLiteral) encodingTerm).getValue();
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
        return tracer + sb.toString();
    }

    private String toStaticCallSerializationExpression(TypeDefinition baseType, Field field, List<Argument> serializerArguments, VariableLiteral vl, Tracer tracer) {
        tracer = tracer.dive("toStaticCallSerializationExpression");
        List<Term> args = vl.getArgs().orElseThrow(() -> new FreemarkerException("'STATIC_CALL' needs at least one attribute"));
        String functionName = args.get(0).asLiteral()
            .orElseThrow(()-> new FreemarkerException("Expecting the first argument of a 'STATIC_CALL' to be a Literal"))
            .asStringLiteral()
            .orElseThrow(()-> new FreemarkerException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral"))
            .getValue();
        // But to make the function name unique, well add the driver prefix to it.
        StringBuilder sb = new StringBuilder(getCTypeName(functionName));
        sb.append("(");
        for (int i = 1; i < args.size(); i++) {
            Term arg = args.get(i);
            if (i > 1) {
                sb.append(", ");
            }
            if (arg instanceof VariableLiteral) {
                VariableLiteral va = (VariableLiteral) arg;
                // "io" and "_value" are always available in every parser.
                boolean isSerializerArg = "writeBuffer".equals(va.getName()) || "_value".equals(va.getName()) || "element".equals(va.getName());
                boolean isTypeArg = "_type".equals(va.getName());
                if (!isSerializerArg && !isTypeArg && serializerArguments != null) {
                    for (Argument serializerArgument : serializerArguments) {
                        if (serializerArgument.getName().equals(va.getName())) {
                            isSerializerArg = true;
                            break;
                        }
                    }
                }
                if (isSerializerArg) {
                    if ("_value".equals(va.getName())) {
                        sb.append("_message");
                    } else {
                        sb.append(va.getName());
                    }
                    if (va.getChild().isPresent()) {
                        sb.append("->");
                        appendVariableExpressionRest(sb, baseType, va.getChild().get());
                    }
                } else if (isTypeArg) {
                    String part = va.getChild().get().getName();
                    switch (part) {
                        case "name":
                            sb.append("\"").append(field.getTypeName()).append("\"");
                            break;
                        case "length":
                            sb.append("\"").append(((SimpleTypeReference) field).getSizeInBits()).append("\"");
                            break;
                        case "encoding":
                            if (!(field instanceof TypedField)) {
                                throw new FreemarkerException("'encoding' only supported for typed fields.");
                            }
                            TypedField typedField = (TypedField) field;
                            if (!(typedField.getType() instanceof StringTypeReference)) {
                                throw new FreemarkerException("Can only access 'encoding' for string types.");
                            }
                            final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                            if (!(encodingTerm instanceof StringLiteral)) {
                                throw new RuntimeException("Encoding must be a quoted string value");
                            }
                            String encoding = ((StringLiteral) encodingTerm).getValue();
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
        return tracer + sb.toString();
    }

    private void appendVariableExpressionRest(StringBuilder sb, TypeDefinition baseType, VariableLiteral variableLiteral) {
        Tracer tracer = Tracer.start("appendVariableExpressionRest");
        if (variableLiteral.isIndexed()) {
            tracer = tracer.dive("isindexed");
            sb.insert(0, "plc4c_utils_list_get_value(");
            sb.append(camelCaseToSnakeCase(variableLiteral.getName()));
            sb.append(", ").append(variableLiteral.getIndex()).append(")");
        } else {
            sb.append(camelCaseToSnakeCase(variableLiteral.getName()));
        }
        // Suppress any "lengthInBytes" properties as these are handled differently in C
        variableLiteral.getChild()
            .filter(child -> !"lengthInBytes".equals(child.getName()))
            .ifPresent(child -> {
                sb.append(".");
                appendVariableExpressionRest(sb, baseType, child);
            });
        sb.append(tracer);
    }

    public int getNumBits(SimpleTypeReference simpleTypeReference) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return 1;
            case UINT:
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                return integerTypeReference.getSizeInBits();
            case FLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                return floatTypeReference.getSizeInBits();
            case STRING:
            case VSTRING:
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return stringTypeReference.getSizeInBits();
            default:
                return 0;
        }
    }

    public String getLengthInBitsFunctionNameForComplexTypedField(Field field) {
        return field.asTypedField()
            .map(TypedField::getType)
            .orElseThrow(() -> new FreemarkerException("lengthInBits functions only exist for TypedFields"))
            .asComplexTypeReference()
            .map(complexTypeReference -> getCTypeName(complexTypeReference.getName()) + "_length_in_bits")
            .orElseThrow(() -> new FreemarkerException("lengthInBits functions only exist for complex types"));
    }

    public String getEnumExpression(String expression) {
        String enumName = expression.substring(0, expression.indexOf('.'));
        String enumConstant = expression.substring(expression.indexOf('.') + 1);
        return getCTypeName(enumName) + "_" + enumConstant;
    }

}
