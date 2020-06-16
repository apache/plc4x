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
import org.apache.commons.text.WordUtils;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.BaseFreemarkerLanguageTemplateHelper;
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
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
        return String.join("", getProtocolName().split("-")) + ".includes";
    }

    /**
     * Little helper that converts a given type name in camel-case into a c-style snake-case expression.
     * In addition it appends a prefix for the protocol name and the output flavor.
     *
     * @param typeName camel-case type name
     * @return c-style type name
     */
    public String getCTypeName(String typeName) {
        return camelCaseToSnakeCase(getProtocolName()).toLowerCase() +
            "_" + camelCaseToSnakeCase(getFlavorName()).toLowerCase() +
            "_" + camelCaseToSnakeCase(typeName).toLowerCase();
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
            throw new RuntimeException("Field " + field + " is not a TypedField");
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
        final Optional<NamedField> typeField = typeDefinition.getFields().stream().filter(field -> field instanceof NamedField).map(field -> (NamedField) field).filter(namedField -> namedField.getName().equals(propertyName)).findFirst();
        // If the property name refers to a field, check if it's an optional field.
        // If it is, pointer access is required, if not, it's not.
        if(typeField.isPresent()) {
            final NamedField namedField = typeField.get();
            return namedField instanceof OptionalField;
        }
        final Optional<Argument> parserArgument = Arrays.stream(typeDefinition.getParserArguments()).filter(argument -> argument.getName().equals(propertyName)).findFirst();
        // If the property name refers to a parser argument, as soon as it's a complex type,
        // pointer access is required.
        if(parserArgument.isPresent()) {
            return parserArgument.get().getType() instanceof ComplexTypeReference;
        }
        // In all other cases, the property might be a built-in constant, so we don't need pointer
        // access for any of these.
        return false;
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
                case UINT:
                case INT: {
                    StringBuilder sb = new StringBuilder();
                    if (simpleTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.UINT) {
                        sb.append("u");
                    }
                    if (simpleTypeReference.getSizeInBits() % 64 == 0) {
                        sb.append("int64_t");
                    } else if (simpleTypeReference.getSizeInBits() % 32 == 0) {
                        sb.append("int32_t");
                    } else if (simpleTypeReference.getSizeInBits() % 16 == 0) {
                        sb.append("int16_t");
                    } else if (simpleTypeReference.getSizeInBits() % 8 == 0) {
                        sb.append("int8_t");
                    } else {
                        if (simpleTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.UINT) {
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
            return null;
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

    @Override
    public String getReadBufferReadMethodCall(SimpleTypeReference simpleTypeReference) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT: {
                return "plc4c_spi_read_bit(buf)";
            }
            case UINT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 4) {
                    return "plc4c_spi_read_unsigned_byte(buf, " + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "plc4c_spi_read_unsigned_short(buf, " + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "plc4c_spi_read_unsigned_int(buf, " + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "plc4c_spi_read_unsigned_long(buf, " + integerTypeReference.getSizeInBits() + ")";
                }
                return "plc4c_spi_read_unsigned_big_integer(buf, " + integerTypeReference.getSizeInBits() + ")";
            }
            case INT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "plc4c_spi_read_byte(buf, " + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "plc4c_spi_read_short(buf, " + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "plc4c_spi_read_int(buf, " + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "plc4c_spi_read_long(buf, " + integerTypeReference.getSizeInBits() + ")";
                }
                return "plc4c_spi_read_big_integer(buf, " + integerTypeReference.getSizeInBits() + ")";
            }
            case FLOAT: {
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                if (floatTypeReference.getSizeInBits() <= 32) {
                    return "plc4c_spi_read_float(buf, " + floatTypeReference.getSizeInBits() + ")";
                } else {
                    return "plc4c_spi_read_double(buf, " + floatTypeReference.getSizeInBits() + ")";
                }
            }
            case STRING: {
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return "plc4c_spi_read_string(buf, " + stringTypeReference.getSizeInBits() + ", \"" +
                    stringTypeReference.getEncoding() + "\")";
            }
        }
        return "Hurz";
    }

    @Override
    public String getWriteBufferReadMethodCall(SimpleTypeReference simpleTypeReference, String fieldName) {
        return null;
    }

    @Override
    public String getNullValueForTypeReference(TypeReference typeReference) {
        return null;
    }






    public String toParseExpression(ComplexTypeDefinition baseType, Field field, Term term, Argument[] parserArguments) {
        return toExpression(baseType, field, term, term1 -> toVariableParseExpression(baseType, field, term1, parserArguments));
    }

    public String toSerializationExpression(ComplexTypeDefinition baseType, Field field, Term term, Argument[] parserArguments) {
        return toExpression(baseType, field, term, term1 -> toVariableSerializationExpression(baseType, field, term1, parserArguments));
    }

    private String toExpression(ComplexTypeDefinition baseType, Field field, Term term, Function<Term, String> variableExpressionGenerator) {
        if (term == null) {
            return "";
        }
        if (term instanceof Literal) {
            if (term instanceof NullLiteral) {
                return "null";
            } else if (term instanceof BooleanLiteral) {
                return Boolean.toString(((BooleanLiteral) term).getValue());
            } else if (term instanceof NumericLiteral) {
                return ((NumericLiteral) term).getNumber().toString();
            } else if (term instanceof StringLiteral) {
                return "\"" + ((StringLiteral) term).getValue() + "\"";
            } else if (term instanceof VariableLiteral) {
                VariableLiteral variableLiteral = (VariableLiteral) term;
                // If this literal references an Enum type, then we have to output it differently.
                if (getTypeDefinitions().get(variableLiteral.getName()) instanceof EnumTypeDefinition) {
                    return variableLiteral.getName() + "." + variableLiteral.getChild().getName();
                } else {
                    return variableExpressionGenerator.apply(term);
                }
            } else {
                throw new RuntimeException("Unsupported Literal type " + term.getClass().getName());
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
                    throw new RuntimeException("Unsupported unary operation type " + ut.getOperation());
            }
        } else if (term instanceof BinaryTerm) {
            BinaryTerm bt = (BinaryTerm) term;
            Term a = bt.getA();
            Term b = bt.getB();
            String operation = bt.getOperation();
            switch (operation) {
                case "^":
                    return "Math.pow((" + toExpression(baseType, field, a, variableExpressionGenerator) + "), (" + toExpression(baseType, field, b, variableExpressionGenerator) + "))";
                default:
                    return "(" + toExpression(baseType, field, a, variableExpressionGenerator) + ") " + operation + " (" + toExpression(baseType, field, b, variableExpressionGenerator) + ")";
            }
        } else if (term instanceof TernaryTerm) {
            TernaryTerm tt = (TernaryTerm) term;
            if ("if".equals(tt.getOperation())) {
                Term a = tt.getA();
                Term b = tt.getB();
                Term c = tt.getC();
                return "((" + toExpression(baseType, field, a, variableExpressionGenerator) + ") ? " + toExpression(baseType, field, b, variableExpressionGenerator) + " : " + toExpression(baseType, field, c, variableExpressionGenerator) + ")";
            } else {
                throw new RuntimeException("Unsupported ternary operation type " + tt.getOperation());
            }
        } else {
            throw new RuntimeException("Unsupported Term type " + term.getClass().getName());
        }
    }

    public String toVariableParseExpression(ComplexTypeDefinition baseType, Field field, Term term, Argument[] parserArguments) {
        VariableLiteral vl = (VariableLiteral) term;
        if("CAST".equals(vl.getName())) {

            StringBuilder sb = new StringBuilder();
            if((vl.getArgs() == null) || (vl.getArgs().size() != 2)) {
                throw new RuntimeException("A CAST expression expects exactly two arguments.");
            }
            final VariableLiteral sourceTerm = (VariableLiteral) vl.getArgs().get(0);
            final VariableLiteral typeTerm = (VariableLiteral) vl.getArgs().get(1);
            ComplexTypeReference castTypeReference = typeTerm::getName;
            final TypeDefinition castType = getTypeDefinitionForTypeReference(castTypeReference);
            // If we're casting to a sub-type of a discriminated value, we got to cast to the parent
            // type instead and add the name of the sub-type as prefix to the property we're tryging to
            // access next.
            String castToType;
            String restExpression;
            if(castType.getParentType() != null) {
                castToType = castType.getParentType().getName();
                if(vl.getChild() != null) {
                    // Change the name of the property to contain the sub-type-prefix.
                    restExpression = "." + camelCaseToSnakeCase(castType.getName()) + "_" + toVariableExpressionRest(vl.getChild());
                } else {
                    restExpression = "";
                }
            } else {
                castToType = castType.getName();
                if(vl.getChild() != null) {
                    restExpression = "." + toVariableExpressionRest(vl.getChild());
                } else {
                    restExpression = "";
                }
            }
            sb.append("((plc4c_").append(getCTypeName(castToType)).append(") (")
                .append(requiresPointerAccess(baseType, sourceTerm.getName()) ? "*" : "")
                .append(toVariableParseExpression(baseType, field, sourceTerm, parserArguments)).append("))")
                .append(restExpression);
            return sb.toString();
        }
        // Any name that is full upper-case is considered a function call.
        // These are generally defined in the spi file evaluation_helper.c.
        // All should have a name prefix "plc4c_spi_evaluation_helper_".
        if (vl.getName().equals(vl.getName().toUpperCase())) {
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
            return sb.toString() + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "");
        }

        final String name = vl.getName();

        // Try to find the type of the addressed property.
        final Optional<TypeReference> propertyTypeOptional =
            getTypeReferenceForProperty(baseType, name);

        // If we couldn't find the type, we didn't find the property.
        if(!propertyTypeOptional.isPresent()) {
            throw new RuntimeException("Could not find property with name '" + name + "' in type " + baseType.getName());
        }

        final TypeReference propertyType = propertyTypeOptional.get();

        // If it's a simple field, there is no sub-type to access.
        if(propertyType instanceof SimpleTypeReference) {
            if(vl.getChild() != null) {
                throw new RuntimeException("Simple property '" + name + "' doesn't have child properties.");
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
                return "plc4c_" + getCTypeName(propertyTypeDefinition.getName()) + "_get_discriminator(" + name + "->_type)." + vl.getChild().getName();
            }
        }
        // Handling enum properties in C is a little more tricky as we have to use the enum value
        // and pass this to a function that then returns the desired property value.
        else if(propertyTypeDefinition instanceof EnumTypeDefinition) {
            EnumTypeDefinition enumTypeDefinition = (EnumTypeDefinition) propertyTypeDefinition;
            StringBuilder sb = new StringBuilder("plc4c_")
                .append(getCTypeName(propertyTypeDefinition.getName()))
                .append("_get_").append(camelCaseToSnakeCase(vl.getChild().getName()));
            return sb.toString();
        }
        // Else ... generate a simple access path.
        return vl.getName() + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "");
    }

    private String toVariableSerializationExpression(ComplexTypeDefinition baseType, Field field, Term term, Argument[] serialzerArguments) {
        VariableLiteral vl = (VariableLiteral) term;
        if ("STATIC_CALL".equals(vl.getName())) {
            StringBuilder sb = new StringBuilder();
            if (!(vl.getArgs().get(0) instanceof StringLiteral)) {
                throw new RuntimeException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral");
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
                        sb.append(va.getName() + ((va.getChild() != null) ? "." + toVariableExpressionRest(va.getChild()) : ""));
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
                                    throw new RuntimeException("'encoding' only supported for typed fields.");
                                }
                                TypedField typedField = (TypedField) field;
                                if(!(typedField.getType() instanceof StringTypeReference)) {
                                    throw new RuntimeException("Can only access 'encoding' for string types.");
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
        // Discriminator values have to be handled a little differently.
        /*else if(vl.getName().equals("DISCRIMINATOR_VALUES")) {
            final String typeName = getLanguageTypeNameForSpecType(field.getType());
            switch (typeName) {
                case "byte":
                    return "((Number) _value.getDiscriminatorValues()[" + vl.getIndex() + "]).byteValue()";
                case "short":
                    return "((Number) _value.getDiscriminatorValues()[" + vl.getIndex() + "]).shortValue()";
                case "int":
                    return "((Number) _value.getDiscriminatorValues()[" + vl.getIndex() + "]).intValue()";
                case "long":
                    return "((Number) _value.getDiscriminatorValues()[" + vl.getIndex() + "]).longValue()";
                case "float":
                    return "((Number) _value.getDiscriminatorValues()[" + vl.getIndex() + "]).floatValue()";
                case "double":
                    return "((Number) _value.getDiscriminatorValues()[" + vl.getIndex() + "]).doubleValue()";
                default:
                    return "_value.getDiscriminatorValues()[" + vl.getIndex() + "]";
            }
        }*/
        // All uppercase names are not fields, but utility methods.
        else if (vl.getName().equals(vl.getName().toUpperCase())) {
            StringBuilder sb = new StringBuilder(vl.getName());
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
                            sb.append(va.getName() + ((va.getChild() != null) ? "." + toVariableExpressionRest(va.getChild()) : ""));
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
                                        throw new RuntimeException("'encoding' only supported for typed fields.");
                                    }
                                    TypedField typedField = (TypedField) field;
                                    if(!(typedField.getType() instanceof StringTypeReference)) {
                                        throw new RuntimeException("Can only access 'encoding' for string types.");
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
        // The synthetic checksumRawData is a local field and should not be accessed as bean property.
        boolean isSerializerArg = "checksumRawData".equals(vl.getName()) || "_value".equals(vl.getName()) || "element".equals(vl.getName());
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
            return vl.getName() + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "");
        } else if (isTypeArg) {
            String part = vl.getChild().getName();
            switch (part) {
                case "name":
                    return "\"" + field.getTypeName() + "\"";
                case "length":
                    return "\"" + ((SimpleTypeReference) field).getSizeInBits() + "\"";
                case "encoding":
                    if(!(field instanceof TypedField)) {
                        throw new RuntimeException("'encoding' only supported for typed fields.");
                    }
                    TypedField typedField = (TypedField) field;
                    if(!(typedField.getType() instanceof StringTypeReference)) {
                        throw new RuntimeException("Can only access 'encoding' for string types.");
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
            return "_value." + toVariableExpressionRest(vl);
        }
    }

    private String toVariableExpressionRest(VariableLiteral vl) {
        return camelCaseToSnakeCase(vl.getName()) + ((vl.isIndexed() ? "[" + vl.getIndex() + "]" : "") +
            ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : ""));
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
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                return integerTypeReference.getSizeInBits();
            }
            default: {
                return 0;
            }
        }
    }

}
