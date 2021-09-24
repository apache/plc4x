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

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Optional;

import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.builder.ReflectionToStringBuilder;
import org.apache.commons.lang3.math.NumberUtils;
import org.apache.commons.text.WordUtils;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.BaseFreemarkerLanguageTemplateHelper;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.Tracer;
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
                    return allowPrimitive ? boolean.class.getSimpleName() : Boolean.class.getSimpleName();
                }
                case BYTE: {
                    return allowPrimitive ? byte.class.getSimpleName() : Byte.class.getSimpleName();
                }
                case UINT: {
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 4) {
                        return allowPrimitive ? byte.class.getSimpleName() : Byte.class.getSimpleName();
                    }
                    if (integerTypeReference.getSizeInBits() <= 8) {
                        return allowPrimitive ? short.class.getSimpleName() : Short.class.getSimpleName();
                    }
                    if (integerTypeReference.getSizeInBits() <= 16) {
                        return allowPrimitive ? int.class.getSimpleName() : Integer.class.getSimpleName();
                    }
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return allowPrimitive ? long.class.getSimpleName() : Long.class.getSimpleName();
                    }
                    return BigInteger.class.getSimpleName();
                }
                case INT: {
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 8) {
                        return allowPrimitive ? byte.class.getSimpleName() : Byte.class.getSimpleName();
                    }
                    if (integerTypeReference.getSizeInBits() <= 16) {
                        return allowPrimitive ? short.class.getSimpleName() : Short.class.getSimpleName();
                    }
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return allowPrimitive ? int.class.getSimpleName() : Integer.class.getSimpleName();
                    }
                    if (integerTypeReference.getSizeInBits() <= 64) {
                        return allowPrimitive ? long.class.getSimpleName() : Long.class.getSimpleName();
                    }
                    return BigInteger.class.getSimpleName();
                }
                case FLOAT:
                case UFLOAT: {
                    FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                    int sizeInBits = ((floatTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.FLOAT) ? 1 : 0) +
                        floatTypeReference.getExponent() + floatTypeReference.getMantissa();
                    if (sizeInBits <= 32) {
                        return allowPrimitive ? float.class.getSimpleName() : Float.class.getSimpleName();
                    }
                    if (sizeInBits <= 64) {
                        return allowPrimitive ? double.class.getSimpleName() : Double.class.getSimpleName();
                    }
                    return BigDecimal.class.getSimpleName();
                }
                case STRING: {
                    return String.class.getSimpleName();
                }
                case TIME: {
                    return LocalTime.class.getSimpleName();
                }
                case DATE: {
                    return LocalDate.class.getSimpleName();
                }
                case DATETIME: {
                    return LocalDateTime.class.getSimpleName();
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
        return toExpression(field, term, variableLiteral -> {
            if (isVariableLiteralVirtualField(variableLiteral) || isVariableLiteralDiscriminatorField(variableLiteral)) { // If we are accessing virtual|discriminator fields, we need to call the getter.
                return "get" + StringUtils.capitalize(variableLiteral.getName()) + "()";
            }
            return toVariableParseExpression(field, variableLiteral, parserArguments);
        });
    }

    public String toParseExpression(TypedField field, Term term, Argument[] parserArguments) {
        Tracer tracer = Tracer.start("toParseExpression");
        return tracer + toExpression(field, term, variableLiteral -> tracer.dive("variableExpressionGenerator") + toVariableParseExpression(field, variableLiteral, parserArguments));
    }

    public String toSerializationExpression(TypedField field, Term term, Argument[] serializerArgments) {
        Tracer tracer = Tracer.start("toSerializationExpression");
        return tracer + toExpression(field, term, variableLiteral -> tracer.dive("variableExpressionGenerator") + toVariableSerializationExpression(field, variableLiteral, serializerArgments));
    }

    private String toExpression(TypedField field, Term term, Function<VariableLiteral, String> variableExpressionGenerator) {
        Tracer tracer = Tracer.start("toExpression");
        if (term == null) {
            return "";
        }
        if (term instanceof Literal) {
            tracer = tracer.dive("literal term instanceOf");
            Literal literal = (Literal) term;
            if (literal instanceof NullLiteral) {
                tracer = tracer.dive("null literal instanceOf");
                return tracer + "null";
            } else if (literal instanceof BooleanLiteral) {
                tracer = tracer.dive("boolean literal instanceOf");
                return tracer + Boolean.toString(((BooleanLiteral) literal).getValue());
            } else if (literal instanceof NumericLiteral) {
                tracer = tracer.dive("numeric literal instanceOf");
                return tracer + ((NumericLiteral) literal).getNumber().toString();
            } else if (literal instanceof StringLiteral) {
                tracer = tracer.dive("string literal instanceOf");
                return tracer + "\"" + ((StringLiteral) literal).getValue() + "\"";
            } else if (literal instanceof VariableLiteral) {
                tracer = tracer.dive("variable literal instanceOf");
                VariableLiteral variableLiteral = (VariableLiteral) literal;
                // If this literal references an Enum type, then we have to output it differently.
                if (getTypeDefinitions().get(variableLiteral.getName()) instanceof EnumTypeDefinition) {
                    tracer = tracer.dive("enum definition instanceOf");
                    return tracer + variableLiteral.getName() + "." + variableLiteral.getChild().getName() +
                        ((variableLiteral.getChild().getChild() != null) ?
                            "." + toVariableExpressionRest(variableLiteral.getChild().getChild()) : "");
                } else {
                    return tracer + variableExpressionGenerator.apply(variableLiteral);
                }
            } else {
                throw new RuntimeException("Unsupported Literal type " + literal.getClass().getName());
            }
        } else if (term instanceof UnaryTerm) {
            tracer = tracer.dive("unary term instanceOf");
            UnaryTerm ut = (UnaryTerm) term;
            Term a = ut.getA();
            switch (ut.getOperation()) {
                case "!":
                    tracer = tracer.dive("case !");
                    return tracer + "!(" + toExpression(field, a, variableExpressionGenerator) + ")";
                case "-":
                    tracer = tracer.dive("case -");
                    return tracer + "-(" + toExpression(field, a, variableExpressionGenerator) + ")";
                case "()":
                    tracer = tracer.dive("case ()");
                    return tracer + "(" + toExpression(field, a, variableExpressionGenerator) + ")";
                default:
                    throw new RuntimeException("Unsupported unary operation type " + ut.getOperation());
            }
        } else if (term instanceof BinaryTerm) {
            tracer = tracer.dive("binary term instanceOf");
            BinaryTerm bt = (BinaryTerm) term;
            Term a = bt.getA();
            Term b = bt.getB();
            String operation = bt.getOperation();
            switch (operation) {
                case "^":
                    tracer = tracer.dive("^");
                    return tracer + "Math.pow((" + toExpression(field, a, variableExpressionGenerator) + "), (" + toExpression(field, b, variableExpressionGenerator) + "))";
                default:
                    return tracer + "(" + toExpression(field, a, variableExpressionGenerator) + ") " + operation + " (" + toExpression(field, b, variableExpressionGenerator) + ")";
            }
        } else if (term instanceof TernaryTerm) {
            tracer = tracer.dive("ternary term instanceOf");
            TernaryTerm tt = (TernaryTerm) term;
            if ("if".equals(tt.getOperation())) {
                Term a = tt.getA();
                Term b = tt.getB();
                Term c = tt.getC();
                return tracer + "((" + toExpression(field, a, variableExpressionGenerator) + ") ? " + toExpression(field, b, variableExpressionGenerator) + " : " + toExpression(field, c, variableExpressionGenerator) + ")";
            } else {
                throw new RuntimeException("Unsupported ternary operation type " + tt.getOperation());
            }
        } else {
            throw new RuntimeException("Unsupported Term type " + term.getClass().getName());
        }
    }

    private String toVariableParseExpression(TypedField field, VariableLiteral variableLiteral, Argument[] parserArguments) {
        Tracer tracer = Tracer.start("toVariableParseExpression");
        // CAST expressions are special as we need to add a ".class" to the second parameter in Java.
        if ("CAST".equals(variableLiteral.getName())) {
            tracer = tracer.dive("CAST");
            StringBuilder sb = new StringBuilder(variableLiteral.getName());
            if ((variableLiteral.getArgs() == null) || (variableLiteral.getArgs().size() != 2)) {
                throw new RuntimeException("A CAST expression expects exactly two arguments.");
            }
            sb.append("(")
                .append(toVariableParseExpression(field, (VariableLiteral) variableLiteral.getArgs().get(0), parserArguments))
                .append(", ")
                .append(((VariableLiteral) variableLiteral.getArgs().get(1)).getName()).append(".class)");
            return tracer + sb.toString() + ((variableLiteral.getChild() != null) ? "." + toVariableExpressionRest(variableLiteral.getChild()) : "");
        } else if ("STATIC_CALL".equals(variableLiteral.getName())) {
            tracer = tracer.dive("STATIC_CALL");
            StringBuilder sb = new StringBuilder();
            if (!(variableLiteral.getArgs().get(0) instanceof StringLiteral)) {
                throw new RuntimeException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral");
            }
            // Get the class and method name
            String methodName = ((StringLiteral) variableLiteral.getArgs().get(0)).getValue();
            // Cut off the double-quotes
            methodName = methodName.substring(1, methodName.length() - 1);
            sb.append(methodName).append("(");
            for (int i = 1; i < variableLiteral.getArgs().size(); i++) {
                Term arg = variableLiteral.getArgs().get(i);
                if (i > 1) {
                    sb.append(", ");
                }
                if (arg instanceof VariableLiteral) {
                    VariableLiteral variableLiteralArg = (VariableLiteral) arg;
                    // "readBuffer" is the default name of the reader argument which is always available.
                    boolean isParserArg = "readBuffer".equals(variableLiteralArg.getName());
                    boolean isTypeArg = "_type".equals(variableLiteralArg.getName());
                    if (!isParserArg && !isTypeArg && parserArguments != null) {
                        for (Argument parserArgument : parserArguments) {
                            if (parserArgument.getName().equals(variableLiteralArg.getName())) {
                                isParserArg = true;
                                break;
                            }
                        }
                    }
                    if (isParserArg) {
                        sb.append(variableLiteralArg.getName()).append((variableLiteralArg.getChild() != null) ? "." + toVariableExpressionRest(variableLiteralArg.getChild()) : "");
                    } else if (isTypeArg) {// We have to manually evaluate the type information at code-generation time.
                        String part = variableLiteralArg.getChild().getName();
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
                        sb.append(toVariableParseExpression(field, variableLiteralArg, null));
                    }
                } else if (arg instanceof StringLiteral) {
                    sb.append(((StringLiteral) arg).getValue());
                }
            }
            sb.append(")");
            return tracer + sb.toString();
        } else if (isVariableLiteralImplicitField(variableLiteral)) { // If we are accessing implicit fields, we need to rely on a local variable instead.
            tracer = tracer.dive("implicit");
            return tracer + variableLiteral.getName();
        } else if (variableLiteral.getName().equals(variableLiteral.getName().toUpperCase())) { // All uppercase names are not fields, but utility methods.
            tracer = tracer.dive("UPPERCASE");
            StringBuilder sb = new StringBuilder(variableLiteral.getName());
            if (variableLiteral.getArgs() != null) {
                sb.append("(");
                boolean firstArg = true;
                for (Term arg : variableLiteral.getArgs()) {
                    if (!firstArg) {
                        sb.append(", ");
                    }
                    sb.append(toParseExpression(field, arg, parserArguments));
                    firstArg = false;
                }
                sb.append(")");
            }
            if (variableLiteral.getIndex() != VariableLiteral.NO_INDEX) {
                sb.append("[").append(variableLiteral.getIndex()).append("]");
            }
            return tracer + sb.toString() + ((variableLiteral.getChild() != null) ? "." + toVariableExpressionRest(variableLiteral.getChild()) : "");
        }
        return tracer + variableLiteral.getName() + ((variableLiteral.getChild() != null) ? "." + toVariableExpressionRest(variableLiteral.getChild()) : "");
    }

    private String toVariableSerializationExpression(TypedField field, VariableLiteral variableLiteral, Argument[] serialzerArguments) {
        Tracer tracer = Tracer.start("variable serialization expression");
        if ("STATIC_CALL".equals(variableLiteral.getName())) {
            tracer = tracer.dive("STATIC_CALL");
            StringBuilder sb = new StringBuilder();
            if (!(variableLiteral.getArgs().get(0) instanceof StringLiteral)) {
                throw new RuntimeException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral");
            }
            String methodName = ((StringLiteral) variableLiteral.getArgs().get(0)).getValue();
            methodName = methodName.substring(1, methodName.length() - 1);
            sb.append(methodName).append("(");
            for (int i = 1; i < variableLiteral.getArgs().size(); i++) {
                Term arg = variableLiteral.getArgs().get(i);
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
                        sb.append(va.getName()).append((va.getChild() != null) ? "." + toVariableExpressionRest(va.getChild()) : "");
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
            return tracer + sb.toString();
        }
        // All uppercase names are not fields, but utility methods.
        else if (variableLiteral.getName().equals(variableLiteral.getName().toUpperCase())) {
            tracer = tracer.dive("UPPER_CASE");
            StringBuilder sb = new StringBuilder(variableLiteral.getName());
            if (variableLiteral.getArgs() != null) {
                sb.append("(");
                boolean firstArg = true;
                for (Term arg : variableLiteral.getArgs()) {
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
                            sb.append(va.getName()).append((va.getChild() != null) ? "." + toVariableExpressionRest(va.getChild()) : "");
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
            return tracer + sb.toString();
        } else if (isVariableLiteralImplicitField(variableLiteral)) { // If we are accessing implicit fields, we need to rely on a local variable instead.
            tracer = tracer.dive("implicit field");
            return tracer + toSerializationExpression(getReferencedImplicitField(variableLiteral), getReferencedImplicitField(variableLiteral).getSerializeExpression(), serialzerArguments);
        } else if (isVariableLiteralVirtualField(variableLiteral)) {
            tracer = tracer.dive("virtual field");
            return tracer + "_value." + toVariableExpressionRest(variableLiteral);
        }
        // The synthetic checksumRawData is a local field and should not be accessed as bean property.
        boolean isSerializerArg = "checksumRawData".equals(variableLiteral.getName()) || "_value".equals(variableLiteral.getName()) || "element".equals(variableLiteral.getName()) || "size".equals(variableLiteral.getName());
        boolean isTypeArg = "_type".equals(variableLiteral.getName());
        if (!isSerializerArg && !isTypeArg && serialzerArguments != null) {
            for (Argument serializerArgument : serialzerArguments) {
                if (serializerArgument.getName().equals(variableLiteral.getName())) {
                    isSerializerArg = true;
                    break;
                }
            }
        }
        if (isSerializerArg) {
            tracer = tracer.dive("serializer arg");
            return tracer + variableLiteral.getName() + ((variableLiteral.getChild() != null) ? "." + toVariableExpressionRest(variableLiteral.getChild()) : "");
        } else if (isTypeArg) {
            tracer = tracer.dive("type arg");
            String part = variableLiteral.getChild().getName();
            switch (part) {
                case "name":
                    return tracer + "\"" + field.getTypeName() + "\"";
                case "length":
                    return tracer + "\"" + ((SimpleTypeReference) field).getSizeInBits() + "\"";
                case "encoding":
                    String encoding = ((StringTypeReference) field.getType()).getEncoding();
                    // Cut off the single quotes.
                    encoding = encoding.substring(1, encoding.length() - 1);
                    return tracer + "\"" + encoding + "\"";
                default:
                    return tracer + "";
            }
        } else {
            return tracer + "_value." + toVariableExpressionRest(variableLiteral);
        }
    }

    private String toVariableExpressionRest(VariableLiteral variableLiteral) {
        Tracer tracer = Tracer.start("variable expression rest");
        // length is kind of a keyword in mspec, so we shouldn't be naming variables length. if we ask for the length of a object we can just return length().
        // This way we can get the length of a string when serializing
        String variableLiteralName = variableLiteral.getName();
        if (variableLiteralName.equals("length")) {
            tracer = tracer.dive("length");
            return tracer + variableLiteralName + "()" + ((variableLiteral.isIndexed() ? "[" + variableLiteral.getIndex() + "]" : "") +
                ((variableLiteral.getChild() != null) ? "." + toVariableExpressionRest(variableLiteral.getChild()) : ""));
        }
        return tracer + "get" + WordUtils.capitalize(variableLiteralName) + "()" + ((variableLiteral.isIndexed() ? "[" + variableLiteral.getIndex() + "]" : "") +
            ((variableLiteral.getChild() != null) ? "." + toVariableExpressionRest(variableLiteral.getChild()) : ""));
    }

    public String getSizeInBits(ComplexTypeDefinition complexTypeDefinition, Argument[] parserArguments) {
        int sizeInBits = 0;
        StringBuilder sb = new StringBuilder();
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
