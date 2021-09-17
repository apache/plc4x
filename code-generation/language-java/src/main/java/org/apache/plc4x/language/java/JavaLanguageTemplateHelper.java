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
package org.apache.plc4x.language.java;

import java.util.Optional;

import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.math.NumberUtils;
import org.apache.commons.text.WordUtils;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.BaseFreemarkerLanguageTemplateHelper;
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.references.*;
import org.apache.plc4x.plugins.codegenerator.types.terms.*;

import java.util.Map;
import java.util.function.Function;

@SuppressWarnings({"unused", "WeakerAccess"})
public class JavaLanguageTemplateHelper extends BaseFreemarkerLanguageTemplateHelper {

    private final Map<String, String> options;

    public JavaLanguageTemplateHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types,
                                      Map<String, String> options) {
        super(thisType, protocolName, flavorName, types);
        this.options = options;
    }

    public String packageName(String protocolName, String languageName, String languageFlavorName) {
        return Optional.ofNullable(options.get("package")).orElseGet(() ->
            "org.apache.plc4x." + String.join("", languageName.split("\\-")) + "." +
                String.join("", protocolName.split("\\-")) + "." +
                String.join("", languageFlavorName.split("\\-")));
    }

    @Override
    public String getLanguageTypeNameForField(Field field) {
        boolean optional = field instanceof OptionalField;
        // If the referenced type is a DataIo type, the value is of type PlcValue.
        if (field instanceof PropertyField) {
            PropertyField propertyField = (PropertyField) field;
            if (propertyField.getType() instanceof ComplexTypeReference) {
                ComplexTypeReference complexTypeReference = (ComplexTypeReference) propertyField.getType();
                final TypeDefinition typeDefinition = getTypeDefinitions().get(complexTypeReference.getName());
                if (typeDefinition instanceof DataIoTypeDefinition) {
                    return "PlcValue";
                }
            }
        }
        return getLanguageTypeNameForTypeReference(((TypedField) field).getType(), !optional);
    }

    public String getNonPrimitiveLanguageTypeNameForField(TypedField field) {
        return getLanguageTypeNameForTypeReference(field.getType(), false);
    }

    public String getLanguageTypeNameForSpecType(TypeReference typeReference) {
        return getLanguageTypeNameForTypeReference(typeReference, true);
    }

    @Override
    public String getLanguageTypeNameForTypeReference(TypeReference typeReference) {
        return getLanguageTypeNameForTypeReference(typeReference, false);
    }

    public String adjustLiterals(String javaType, String value) {
        switch (javaType) {
            case "long":
            case "Long":
                return value + "L";
            default:
                return value;
        }
    }

    public String getLanguageTypeNameForTypeReference(TypeReference typeReference, boolean allowPrimitive) {
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT: {
                    return allowPrimitive ? "boolean" : "Boolean";
                }
                case BYTE: {
                    return allowPrimitive ? "byte" : "Byte";
                }
                case UINT: {
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 4) {
                        return allowPrimitive ? "byte" : "Byte";
                    }
                    if (integerTypeReference.getSizeInBits() <= 8) {
                        return allowPrimitive ? "short" : "Short";
                    }
                    if (integerTypeReference.getSizeInBits() <= 16) {
                        return allowPrimitive ? "int" : "Integer";
                    }
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return allowPrimitive ? "long" : "Long";
                    }
                    return "BigInteger";
                }
                case INT: {
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 8) {
                        return allowPrimitive ? "byte" : "Byte";
                    }
                    if (integerTypeReference.getSizeInBits() <= 16) {
                        return allowPrimitive ? "short" : "Short";
                    }
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return allowPrimitive ? "int" : "Integer";
                    }
                    if (integerTypeReference.getSizeInBits() <= 64) {
                        return allowPrimitive ? "long" : "Long";
                    }
                    return "BigInteger";
                }
                case FLOAT:
                case UFLOAT: {
                    FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                    int sizeInBits = ((floatTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.FLOAT) ? 1 : 0) +
                        floatTypeReference.getExponent() + floatTypeReference.getMantissa();
                    if (sizeInBits <= 32) {
                        return allowPrimitive ? "float" : "Float";
                    }
                    if (sizeInBits <= 64) {
                        return allowPrimitive ? "double" : "Double";
                    }
                    return "BigDecimal";
                }
                case STRING: {
                    return "String";
                }
                case TIME: {
                    return "LocalTime";
                }
                case DATE: {
                    return "LocalDate";
                }
                case DATETIME: {
                    return "LocalDateTime";
                }
            }
            throw new RuntimeException("Unsupported simple type");
        } else {
            return ((ComplexTypeReference) typeReference).getName();
        }
    }

    public String getPlcValueTypeForTypeReference(TypeReference typeReference) {
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT: {
                    return "PlcBOOL";
                }
                case BYTE: {
                    return "PlcSINT";
                }
                case UINT: {
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 4) {
                        return "PlcUSINT";
                    }
                    if (integerTypeReference.getSizeInBits() <= 8) {
                        return "PlcUINT";
                    }
                    if (integerTypeReference.getSizeInBits() <= 16) {
                        return "PlcUDINT";
                    }
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return "PlcULINT";
                    }
                }
                case INT: {
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 8) {
                        return "PlcSINT";
                    }
                    if (integerTypeReference.getSizeInBits() <= 16) {
                        return "PlcINT";
                    }
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return "PlcDINT";
                    }
                    if (integerTypeReference.getSizeInBits() <= 64) {
                        return "PlcLINT";
                    }
                }
                case FLOAT:
                case UFLOAT: {
                    FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                    int sizeInBits = ((floatTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.FLOAT) ? 1 : 0) +
                        floatTypeReference.getExponent() + floatTypeReference.getMantissa();
                    if (sizeInBits <= 32) {
                        return "PlcREAL";
                    }
                    if (sizeInBits <= 64) {
                        return "PlcLREAL";
                    }
                }
                case STRING: {
                    return "PlcSTRING";
                }
                case TIME: {
                    return "PlcTIME";
                }
                case DATE: {
                    return "PlcTIME";
                }
                case DATETIME: {
                    return "PlcTIME";
                }
            }
            throw new RuntimeException("Unsupported simple type");
        } else {
            return "PlcStruct";
        }
    }

    @Override
    public String getNullValueForTypeReference(TypeReference typeReference) {
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT: {
                    return "false";
                }
                case BYTE: {
                    return "0";
                }
                case UINT: {
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 16) {
                        return "0";
                    }
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return "0l";
                    }
                    return "null";
                }
                case INT: {
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return "0";
                    }
                    if (integerTypeReference.getSizeInBits() <= 64) {
                        return "0l";
                    }
                    return "null";
                }
                case FLOAT: {
                    FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                    int sizeInBits = ((floatTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.FLOAT) ? 1 : 0) +
                        floatTypeReference.getExponent() + floatTypeReference.getMantissa();
                    if (sizeInBits <= 32) {
                        return "0.0f";
                    }
                    if (sizeInBits <= 64) {
                        return "0.0";
                    }
                    return "null";
                }
                case STRING: {
                    return "null";
                }
            }
            return "Hurz";
        } else {
            return "null";
        }
    }

    /*public String getArgumentType(TypeReference typeReference, int index) {
        if(typeReference instanceof ComplexTypeReference) {
            ComplexTypeReference complexTypeReference = (ComplexTypeReference) typeReference;
            if(!getTypeDefinitions().containsKey(complexTypeReference.getName())) {
                throw new RuntimeException("Could not find definition of complex type " + complexTypeReference.getName());
            }
            TypeDefinition complexTypeDefinition = getTypeDefinitions().get(complexTypeReference.getName());
            if(complexTypeDefinition.getParserArguments().length <= index) {
                throw new RuntimeException("Type " + complexTypeReference.getName() + " specifies too few parser arguments");
            }
            return getLanguageTypeNameForSpecType(complexTypeDefinition.getParserArguments()[index].getType());
        }
        return "Hurz";
    }*/

    public int getNumBits(SimpleTypeReference simpleTypeReference) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT: {
                return 1;
            }
            case BYTE: {
                return 8;
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

    @Override
    public String getReadBufferReadMethodCall(SimpleTypeReference simpleTypeReference, String valueString, TypedField field) {
        return getReadBufferReadMethodCall("", simpleTypeReference, valueString, field);
    }

    public String getReadBufferReadMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, String valueString, TypedField field) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT: {
                return "readBuffer.readBit(\"" + logicalName + "\")";
            }
            case BYTE: {
                ByteTypeReference byteTypeReference = (ByteTypeReference) simpleTypeReference;
                return "readBuffer.readByte(\"" + logicalName + "\")";
            }
            case UINT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 4) {
                    return "readBuffer.readUnsignedByte(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "readBuffer.readUnsignedShort(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "readBuffer.readUnsignedInt(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "readBuffer.readUnsignedLong(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ")";
                }
                return "readBuffer.readUnsignedBigInteger(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ")";
            }
            case INT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "readBuffer.readSignedByte(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "readBuffer.readShort(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "readBuffer.readInt(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "readBuffer.readLong(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ")";
                }
                return "readBuffer.readBigInteger(" + integerTypeReference.getSizeInBits() + ")";
            }
            case FLOAT: {
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                String type = (floatTypeReference.getSizeInBits() <= 32) ? "Float" : "Double";
                String typeCast = (floatTypeReference.getSizeInBits() <= 32) ? "float" : "double";
                String defaultNull = (floatTypeReference.getSizeInBits() <= 32) ? "0.0f" : "0.0";
                return "((Supplier<" + type + ">) (() -> {" +
                    "\n            return (" + typeCast + ") toFloat(readBuffer, \"" + logicalName + "\", " +
                    ((floatTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.FLOAT) ? "true" : "false") +
                    ", " + floatTypeReference.getExponent() + ", " +
                    floatTypeReference.getMantissa() + ");" +
                    "\n        })).get()";
            }
            case STRING: {
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return "readBuffer.readString(\"" + logicalName + "\", " + toParseExpression(field, stringTypeReference.getLengthExpression(), null) + ", \"" +
                    stringTypeReference.getEncoding() + "\")";
            }
        }
        return "Hurz";
    }

    @Override
    public String getWriteBufferWriteMethodCall(SimpleTypeReference simpleTypeReference, String fieldName, TypedField field) {
        return getWriteBufferWriteMethodCall("", simpleTypeReference, fieldName, field);
    }

    public String getWriteBufferWriteMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, String fieldName, TypedField field, String... writerArgs) {
        String writerArgsString = "";
        if (writerArgs.length > 0) {
            writerArgsString += ", " + StringUtils.join(writerArgs, ", ");
        }
        switch (simpleTypeReference.getBaseType()) {
            case BIT: {
                return "writeBuffer.writeBit(\"" + logicalName + "\", (boolean) " + fieldName + "" + writerArgsString + ")";
            }
            case BYTE: {
                ByteTypeReference byteTypeReference = (ByteTypeReference) simpleTypeReference;
                return "writeBuffer.writeByte(\"" + logicalName + "\", ((Number) " + fieldName + ").byteValue()" + writerArgsString + ")";
            }
            case UINT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 4) {
                    return "writeBuffer.writeUnsignedByte(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").byteValue()" + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "writeBuffer.writeUnsignedShort(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").shortValue()" + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "writeBuffer.writeUnsignedInt(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").intValue()" + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "writeBuffer.writeUnsignedLong(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").longValue()" + writerArgsString + ")";
                }
                return "writeBuffer.writeUnsignedBigInteger(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", (BigInteger) " + fieldName + "" + writerArgsString + ")";
            }
            case INT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "writeBuffer.writeSignedByte(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").byteValue()" + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "writeBuffer.writeShort(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").shortValue()" + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "writeBuffer.writeInt(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").intValue()" + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "writeBuffer.writeLong(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").longValue()" + writerArgsString + ")";
                }
                return "writeBuffer.writeBigInteger(\"" + logicalName + "\", " + integerTypeReference.getSizeInBits() + ", BigInteger.valueOf( " + fieldName + ")" + writerArgsString + ")";
            }
            case FLOAT:
            case UFLOAT: {
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                if (floatTypeReference.getSizeInBits() <= 32) {
                    return "writeBuffer.writeFloat(\"" + logicalName + "\", " + fieldName + "," + floatTypeReference.getExponent() + "," + floatTypeReference.getMantissa() + "" + writerArgsString + ")";
                } else if (floatTypeReference.getSizeInBits() <= 64) {
                    return "writeBuffer.writeDouble(\"" + logicalName + "\", " + fieldName + "," + floatTypeReference.getExponent() + "," + floatTypeReference.getMantissa() + "" + writerArgsString + ")";
                } else {
                    throw new RuntimeException("Unsupported float type");
                }
            }
            case STRING: {
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return "writeBuffer.writeString(\"" + logicalName + "\", " + toSerializationExpression(field, stringTypeReference.getLengthExpression(), getThisTypeDefinition().getParserArguments()) + ", \"" +
                    stringTypeReference.getEncoding() + "\", (String) " + fieldName + "" + writerArgsString + ")";
            }
        }
        return "Hurz";
    }

    /*public String getReadMethodName(SimpleTypeReference simpleTypeReference) {
        String languageTypeName = getLanguageTypeNameForSpecType(simpleTypeReference);
        languageTypeName = languageTypeName.substring(0, 1).toUpperCase() + languageTypeName.substring(1);
        if(simpleTypeReference.getBaseType().equals(SimpleTypeReference.SimpleBaseType.UINT)) {
            return "readUnsigned" + languageTypeName;
        } else {
            return "read" + languageTypeName;
        }
    }*/

    public String getReservedValue(ReservedField reservedField) {
        final String languageTypeName = getLanguageTypeNameForTypeReference(reservedField.getType(), true);
        if ("BigInteger".equals(languageTypeName)) {
            return "BigInteger.valueOf(" + reservedField.getReferenceValue() + ")";
        } else {
            return "(" + languageTypeName + ") " + reservedField.getReferenceValue();
        }
    }

    /*public Collection<ComplexTypeReference> getComplexTypes(ComplexTypeDefinition complexTypeDefinition) {
        Map<String, ComplexTypeReference> types = new HashMap<>();
        for (Field field : complexTypeDefinition.getFields()) {
            if(field instanceof TypedField) {
                TypedField typedField = (TypedField) field;
                if(typedField.getType() instanceof ComplexTypeReference) {
                    ComplexTypeReference complexTypeReference = (ComplexTypeReference) typedField.getType();
                    types.put(complexTypeReference.getName(),  complexTypeReference);
                }
            } else if(field instanceof SwitchField) {
                SwitchField switchField = (SwitchField) field;
                for (DiscriminatedComplexTypeDefinition cas : switchField.getCases()) {
                    types.put(cas.getName(), new ComplexTypeReference() {
                        @Override
                        public String getName() {
                            return cas.getName();
                        }
                    });
                }
            }
        }
        return types.values();
    }*/

    /*public Collection<ComplexTypeReference> getEnumTypes(ComplexTypeDefinition complexTypeDefinition) {
        Map<String, ComplexTypeReference> types = new HashMap<>();
        for (Field field : complexTypeDefinition.getFields()) {
            if(field instanceof EnumField) {
                EnumField enumField = (EnumField) field;
                if(enumField.getType() instanceof ComplexTypeReference) {
                    ComplexTypeReference complexTypeReference = (ComplexTypeReference) enumField.getType();
                    types.put(complexTypeReference.getName(),  complexTypeReference);
                }
            }
        }
        for (Field field : complexTypeDefinition.getParentPropertyFields()) {
            if(field instanceof EnumField) {
                EnumField enumField = (EnumField) field;
                if(enumField.getType() instanceof ComplexTypeReference) {
                    ComplexTypeReference complexTypeReference = (ComplexTypeReference) enumField.getType();
                    types.put(complexTypeReference.getName(),  complexTypeReference);
                }
            }
        }
        return types.values();
    }*/

    public String toAccessExpression(TypedField field, Term term, Argument[] parserArguments) {
        return toExpression(field, term, term1 -> {
            VariableLiteral vl = (VariableLiteral) term1;
            if (isVariableLiteralVirtualField(vl) || isVariableLiteralDiscriminatorField(vl)) { // If we are accessing virtual|discriminator fields, we need to call the getter.
                return "get" + StringUtils.capitalize(vl.getName()) + "()";
            }
            return toVariableParseExpression(field, term1, parserArguments);
        });
    }

    public String toParseExpression(TypedField field, Term term, Argument[] parserArguments) {
        return toExpression(field, term, term1 -> toVariableParseExpression(field, term1, parserArguments));
    }

    public String toSerializationExpression(TypedField field, Term term, Argument[] serializerArgments) {
        return toExpression(field, term, term1 -> toVariableSerializationExpression(field, term1, serializerArgments));
    }

    private String toExpression(TypedField field, Term term, Function<Term, String> variableExpressionGenerator) {
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
                    return variableLiteral.getName() + "." + variableLiteral.getChild().getName() +
                        ((variableLiteral.getChild().getChild() != null) ?
                            "." + toVariableExpressionRest(variableLiteral.getChild().getChild()) : "");
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
                    return "!(" + toExpression(field, a, variableExpressionGenerator) + ")";
                case "-":
                    return "-(" + toExpression(field, a, variableExpressionGenerator) + ")";
                case "()":
                    return "(" + toExpression(field, a, variableExpressionGenerator) + ")";
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
                    return "Math.pow((" + toExpression(field, a, variableExpressionGenerator) + "), (" + toExpression(field, b, variableExpressionGenerator) + "))";
                default:
                    return "(" + toExpression(field, a, variableExpressionGenerator) + ") " + operation + " (" + toExpression(field, b, variableExpressionGenerator) + ")";
            }
        } else if (term instanceof TernaryTerm) {
            TernaryTerm tt = (TernaryTerm) term;
            if ("if".equals(tt.getOperation())) {
                Term a = tt.getA();
                Term b = tt.getB();
                Term c = tt.getC();
                return "((" + toExpression(field, a, variableExpressionGenerator) + ") ? " + toExpression(field, b, variableExpressionGenerator) + " : " + toExpression(field, c, variableExpressionGenerator) + ")";
            } else {
                throw new RuntimeException("Unsupported ternary operation type " + tt.getOperation());
            }
        } else {
            throw new RuntimeException("Unsupported Term type " + term.getClass().getName());
        }
    }

    private String toVariableParseExpression(TypedField field, Term term, Argument[] parserArguments) {
        VariableLiteral vl = (VariableLiteral) term;
        // CAST expressions are special as we need to add a ".class" to the second parameter in Java.
        if ("CAST".equals(vl.getName())) {
            StringBuilder sb = new StringBuilder(vl.getName());
            if ((vl.getArgs() == null) || (vl.getArgs().size() != 2)) {
                throw new RuntimeException("A CAST expression expects exactly two arguments.");
            }
            sb.append("(").append(toVariableParseExpression(field, vl.getArgs().get(0), parserArguments))
                .append(", ").append(((VariableLiteral) vl.getArgs().get(1)).getName()).append(".class)");
            return sb + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "");
        } else if ("STATIC_CALL".equals(vl.getName())) {
            StringBuilder sb = new StringBuilder();
            if (!(vl.getArgs().get(0) instanceof StringLiteral)) {
                throw new RuntimeException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral");
            }
            // Get the class and method name
            String methodName = ((StringLiteral) vl.getArgs().get(0)).getValue();
            // Cut off the double-quotes
            methodName = methodName.substring(1, methodName.length() - 1);
            sb.append(methodName).append("(");
            for (int i = 1; i < vl.getArgs().size(); i++) {
                Term arg = vl.getArgs().get(i);
                if (i > 1) {
                    sb.append(", ");
                }
                if (arg instanceof VariableLiteral) {
                    VariableLiteral va = (VariableLiteral) arg;
                    // "readBuffer" is the default name of the reader argument which is always available.
                    boolean isParserArg = "readBuffer".equals(va.getName());
                    boolean isTypeArg = "_type".equals(va.getName());
                    if (!isParserArg && !isTypeArg && parserArguments != null) {
                        for (Argument parserArgument : parserArguments) {
                            if (parserArgument.getName().equals(va.getName())) {
                                isParserArg = true;
                                break;
                            }
                        }
                    }
                    if (isParserArg) {
                        sb.append(va.getName()).append((va.getChild() != null) ? "." + toVariableExpressionRest(va.getChild()) : "");
                    } else if (isTypeArg) {// We have to manually evaluate the type information at code-generation time.
                        String part = va.getChild().getName();
                        switch (part) {
                            case "name":
                                sb.append("\"").append(field.getTypeName()).append("\"");
                                break;
                            case "length":
                                sb.append("\"").append(((SimpleTypeReference) field).getSizeInBits()).append("\"");
                                break;
                            case "encoding":
                                String encoding = ((StringTypeReference) field.getType()).getEncoding();
                                // Cut off the single quotes.
                                encoding = encoding.substring(1, encoding.length() - 1);
                                sb.append("\"").append(encoding).append("\"");
                                break;
                        }
                    } else {
                        sb.append(toVariableParseExpression(field, va, null));
                    }
                } else if (arg instanceof StringLiteral) {
                    sb.append(((StringLiteral) arg).getValue());
                }
            }
            sb.append(")");
            return sb.toString();
        } else if (isVariableLiteralImplicitField(vl)) { // If we are accessing implicit fields, we need to rely on a local variable instead.
            return vl.getName();
        } else if (vl.getName().equals(vl.getName().toUpperCase())) { // All uppercase names are not fields, but utility methods.
            StringBuilder sb = new StringBuilder(vl.getName());
            if (vl.getArgs() != null) {
                sb.append("(");
                boolean firstArg = true;
                for (Term arg : vl.getArgs()) {
                    if (!firstArg) {
                        sb.append(", ");
                    }
                    sb.append(toParseExpression(field, arg, parserArguments));
                    firstArg = false;
                }
                sb.append(")");
            }
            if (vl.getIndex() != VariableLiteral.NO_INDEX) {
                sb.append("[").append(vl.getIndex()).append("]");
            }
            return sb + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "");
        }
        return vl.getName() + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "");
    }

    private String toVariableSerializationExpression(TypedField field, Term term, Argument[] serialzerArguments) {
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
                    // "readBuffer" and "_value" are always available in every parser.
                    boolean isSerializerArg = "readBuffer".equals(va.getName()) || "writeBuffer".equals(va.getName()) || "_value".equals(va.getName()) || "element".equals(va.getName());
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
                                String encoding = ((StringTypeReference) field.getType()).getEncoding();
                                // Cut off the single quotes.
                                encoding = encoding.substring(1, encoding.length() - 1);
                                sb.append("\"").append(encoding).append("\"");
                                break;
                        }
                    } else {
                        sb.append(toVariableSerializationExpression(field, va, serialzerArguments));
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
                        boolean isSerializerArg = "readBuffer".equals(va.getName()) || "writeBuffer".equals(va.getName());
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
                                    String encoding = ((StringTypeReference) field.getType()).getEncoding();
                                    // Cut off the single quotes.
                                    encoding = encoding.substring(1, encoding.length() - 1);
                                    sb.append("\"").append(encoding).append("\"");
                                    break;
                            }
                        } else {
                            sb.append(toVariableSerializationExpression(field, va, serialzerArguments));
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
        else if (isVariableLiteralImplicitField(vl)) {
            return toSerializationExpression(getReferencedImplicitField(vl), getReferencedImplicitField(vl).getSerializeExpression(), serialzerArguments);
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
            return vl.getName() + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "");
        } else if (isTypeArg) {
            String part = vl.getChild().getName();
            switch (part) {
                case "name":
                    return "\"" + field.getTypeName() + "\"";
                case "length":
                    return "\"" + ((SimpleTypeReference) field).getSizeInBits() + "\"";
                case "encoding":
                    String encoding = ((StringTypeReference) field.getType()).getEncoding();
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
        // length is kind of a keyword in mspec, so we shouldn't be naming variables length. if we ask for the length of a object we can just return length().
        // This way we can get the length of a string when serializing
        if (vl.getName().equals("length")) {
            return vl.getName() + "()" + ((vl.isIndexed() ? "[" + vl.getIndex() + "]" : "") +
                ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : ""));
        } else {
            return "get" + WordUtils.capitalize(vl.getName()) + "()" + ((vl.isIndexed() ? "[" + vl.getIndex() + "]" : "") +
                ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : ""));
        }
    }

    public String getSizeInBits(ComplexTypeDefinition complexTypeDefinition, Argument[] parserArguments) {
        int sizeInBits = 0;
        StringBuilder sb = new StringBuilder("");
        for (Field field : complexTypeDefinition.getFields()) {
            if (field instanceof ArrayField) {
                ArrayField arrayField = (ArrayField) field;
                final SimpleTypeReference type = (SimpleTypeReference) arrayField.getType();
                switch (arrayField.getLoopType()) {
                    case COUNT:
                        sb.append("(").append(toSerializationExpression(null, arrayField.getLoopExpression(), parserArguments)).append(" * ").append(type.getSizeInBits()).append(") + ");
                        break;
                    case LENGTH:
                        sb.append("(").append(toSerializationExpression(null, arrayField.getLoopExpression(), parserArguments)).append(" * 8) + ");
                        break;
                    case TERMINATED:
                        // No terminated.
                        break;
                }
            } else if (field instanceof TypedField) {
                TypedField typedField = (TypedField) field;
                final TypeReference type = typedField.getType();
                if (field instanceof ManualField) {
                    ManualField manualField = (ManualField) field;
                    sb.append("(").append(toSerializationExpression(null, manualField.getLengthExpression(), parserArguments)).append(") + ");
                } else if (type instanceof SimpleTypeReference) {
                    SimpleTypeReference simpleTypeReference = (SimpleTypeReference) type;
                    if (simpleTypeReference instanceof StringTypeReference) {
                        sb.append(toSerializationExpression(null, ((StringTypeReference) simpleTypeReference).getLengthExpression(), parserArguments)).append(" + ");
                    } else {
                        sizeInBits += simpleTypeReference.getSizeInBits();
                    }
                } else {
                    // No ComplexTypeReference supported.
                }
            }
        }
        return sb.toString() + sizeInBits;
    }

    public String escapeValue(TypeReference typeReference, String valueString) {
        if (valueString == null) {
            return null;
        }
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case UINT:
                case INT:
                    // If it's a one character string and is numeric, output it as char.
                    if (!NumberUtils.isParsable(valueString) && (valueString.length() == 1)) {
                        return "'" + valueString + "'";
                    }
                    break;
                case STRING:
                    return "\"" + valueString + "\"";
            }
        }
        return valueString;
    }

}
