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

import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.math.NumberUtils;
import org.apache.commons.text.WordUtils;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.BaseFreemarkerLanguageTemplateHelper;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.FreemarkerException;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.Tracer;
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.references.*;
import org.apache.plc4x.plugins.codegenerator.types.terms.*;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.function.Function;

@SuppressWarnings({"unused", "WeakerAccess"})
public class JavaLanguageTemplateHelper extends BaseFreemarkerLanguageTemplateHelper {

    private final Map<String, String> options;

    public JavaLanguageTemplateHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types,
                                      Map<String, String> options) {
        super(thisType, protocolName, flavorName, types);
        this.options = options;
    }

    public String packageName() {
        return packageName(protocolName, "java", flavorName);
    }

    public String packageName(String protocolName, String languageName, String languageFlavorName) {
        return Optional.ofNullable(options.get("package")).orElseGet(() ->
            "org.apache.plc4x." + String.join("", languageName.split("\\-")) + "." +
                String.join("", protocolName.split("\\-")) + "." +
                String.join("", languageFlavorName.split("\\-")));
    }

    @Override
    public String getLanguageTypeNameForField(Field field) {
        // If the referenced type is a DataIo type, the value is of type PlcValue.
        if (field.isPropertyField()) {
            PropertyField propertyField = field.asPropertyField().orElseThrow(IllegalStateException::new);
            if (propertyField.getType().isComplexTypeReference()) {
                ComplexTypeReference complexTypeReference = propertyField.getType().asComplexTypeReference().orElseThrow(IllegalStateException::new);
                final TypeDefinition typeDefinition = getTypeDefinitions().get(complexTypeReference.getName());
                if (typeDefinition instanceof DataIoTypeDefinition) {
                    return "PlcValue";
                }
            }
        }
        return getLanguageTypeNameForTypeReference(((TypedField) field).getType(), !field.isOptionalField());
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
        if (!(typeReference instanceof SimpleTypeReference)) {
            return ((ComplexTypeReference) typeReference).getName();
        }
        SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return allowPrimitive ? boolean.class.getSimpleName() : Boolean.class.getSimpleName();
            case BYTE:
                return allowPrimitive ? byte.class.getSimpleName() : Byte.class.getSimpleName();
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 4) {
                    return allowPrimitive ? byte.class.getSimpleName() : Byte.class.getSimpleName();
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    return allowPrimitive ? short.class.getSimpleName() : Short.class.getSimpleName();
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    return allowPrimitive ? int.class.getSimpleName() : Integer.class.getSimpleName();
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    return allowPrimitive ? long.class.getSimpleName() : Long.class.getSimpleName();
                }
                return BigInteger.class.getSimpleName();
            case INT:
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
            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                int sizeInBits = floatTypeReference.getSizeInBits();
                if (sizeInBits <= 32) {
                    return allowPrimitive ? float.class.getSimpleName() : Float.class.getSimpleName();
                }
                if (sizeInBits <= 64) {
                    return allowPrimitive ? double.class.getSimpleName() : Double.class.getSimpleName();
                }
                return BigDecimal.class.getSimpleName();
            case STRING:
            case VSTRING:
                return String.class.getSimpleName();
            case TIME:
                return LocalTime.class.getSimpleName();
            case DATE:
                return LocalDate.class.getSimpleName();
            case DATETIME:
                return LocalDateTime.class.getSimpleName();

        }
        throw new RuntimeException("Unsupported simple type");
    }

    public String getPlcValueTypeForTypeReference(TypeReference typeReference) {
        if (!(typeReference instanceof SimpleTypeReference)) {
            return "PlcStruct";
        }
        SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "PlcBOOL";
            case BYTE:
                return "PlcSINT";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 4) {
                    return "PlcUSINT";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    return "PlcUINT";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    return "PlcUDINT";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    return "PlcULINT";
                }
            case INT:
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

            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                int sizeInBits = floatTypeReference.getSizeInBits();
                if (sizeInBits <= 32) {
                    return "PlcREAL";
                }
                if (sizeInBits <= 64) {
                    return "PlcLREAL";
                }
            case STRING:
            case VSTRING:
                return "PlcSTRING";
            case TIME:
            case DATE:
            case DATETIME:
                return "PlcTIME";
        }
        throw new RuntimeException("Unsupported simple type");
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
                    IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                        return "0";
                    }
                    if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                        return "0l";
                    }
                    return "null";
                case INT:
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return "0";
                    }
                    if (integerTypeReference.getSizeInBits() <= 64) {
                        return "0l";
                    }
                    return "null";
                case FLOAT:
                    FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                    int sizeInBits = floatTypeReference.getSizeInBits();
                    if (sizeInBits <= 32) {
                        return "0.0f";
                    }
                    if (sizeInBits <= 64) {
                        return "0.0";
                    }
                    return "null";
                case STRING:
                case VSTRING:
                    return "null";
            }
            throw new FreemarkerException("Unmapped basetype" + simpleTypeReference.getBaseType());
        } else {
            return "null";
        }
    }

    public int getNumBits(SimpleTypeReference simpleTypeReference) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return 1;
            case BYTE:
                return 8;
            case UINT:
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                return integerTypeReference.getSizeInBits();
            case FLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                return floatTypeReference.getSizeInBits();
            case STRING:
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return stringTypeReference.getSizeInBits();
            case VSTRING:
                throw new IllegalArgumentException("getSizeInBits doesn't work for 'vstring' fields");
            default:
                return 0;
        }
    }

    @Deprecated
    @Override
    public String getReadBufferReadMethodCall(SimpleTypeReference simpleTypeReference, String valueString, TypedField field) {
        return "/*TODO: migrate me*/" + getReadBufferReadMethodCall("", simpleTypeReference, valueString, field);
    }

    @Deprecated
    public String getReadBufferReadMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, String valueString, TypedField field) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                String bitType = "Bit";
                return "/*TODO: migrate me*/" + "readBuffer.read" + bitType + "(\"" + logicalName + "\")";
            case BYTE:
                String byteType = "Byte";
                return "/*TODO: migrate me*/" + "readBuffer.read" + byteType + "(\"" + logicalName + "\")";
            case UINT:
                String unsignedIntegerType;
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 4) {
                    unsignedIntegerType = "UnsignedByte";
                } else if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    unsignedIntegerType = "UnsignedShort";
                } else if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    unsignedIntegerType = "UnsignedInt";
                } else if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    unsignedIntegerType = "UnsignedLong";
                } else {
                    unsignedIntegerType = "UnsignedBigInteger";
                }
                return "/*TODO: migrate me*/" + "readBuffer.read" + unsignedIntegerType + "(\"" + logicalName + "\", " + simpleTypeReference.getSizeInBits() + ")";
            case INT:
                String integerType;
                if (simpleTypeReference.getSizeInBits() <= 8) {
                    integerType = "SignedByte";
                } else if (simpleTypeReference.getSizeInBits() <= 16) {
                    integerType = "Short";
                } else if (simpleTypeReference.getSizeInBits() <= 32) {
                    integerType = "Int";
                } else if (simpleTypeReference.getSizeInBits() <= 64) {
                    integerType = "Long";
                } else {
                    integerType = "BigInteger";
                }
                return "/*TODO: migrate me*/" + "readBuffer.read" + integerType + "(\"" + logicalName + "\", " + simpleTypeReference.getSizeInBits() + ")";
            case FLOAT:
                String floatType = (simpleTypeReference.getSizeInBits() <= 32) ? "Float" : "Double";
                return "/*TODO: migrate me*/" + "readBuffer.read" + floatType + "(\"" + logicalName + "\", " + simpleTypeReference.getSizeInBits() + ")";
            case STRING:
            case VSTRING:
                String stringType = "String";
                final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                if (!(encodingTerm instanceof StringLiteral)) {
                    throw new RuntimeException("Encoding must be a quoted string value");
                }
                String encoding = ((StringLiteral) encodingTerm).getValue();
                String length = Integer.toString(simpleTypeReference.getSizeInBits());
                if (simpleTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.VSTRING) {
                    VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                    length = toParseExpression(field, vstringTypeReference.getLengthExpression(), null);
                }
                return "/*TODO: migrate me*/" + "readBuffer.read" + stringType + "(\"" + logicalName + "\", " + length + ", \"" +
                    encoding + "\")";
        }
        return "/*TODO: migrate me*/" + "";
    }

    public String getDataReaderCall(TypeReference typeReference) {
        return getDataReaderCall(typeReference, "enumForValue");
    }

    public String getDataReaderCall(TypeReference typeReference, String resolverMethod) {
        if (isEnumTypeReference(typeReference)) {
            final String languageTypeName = getLanguageTypeNameForTypeReference(typeReference);
            final SimpleTypeReference enumBaseTypeReference = getEnumBaseTypeReference(typeReference);
            return "new DataReaderEnumDefault<>(" + languageTypeName + "::" + resolverMethod + ", " + getDataReaderCall(enumBaseTypeReference) + ")";
        } else if (typeReference.isSimpleTypeReference()) {
            SimpleTypeReference simpleTypeReference = typeReference.asSimpleTypeReference().orElseThrow(IllegalStateException::new);
            return getDataReaderCall(simpleTypeReference);
        } else if (typeReference.isComplexTypeReference()) {
            StringBuilder paramsString = new StringBuilder();
            ComplexTypeReference complexTypeReference = typeReference.asComplexTypeReference().orElseThrow(IllegalStateException::new);
            TypeDefinition typeDefinition = getTypeDefinitionForTypeReference(typeReference);
            String parserCallString = getLanguageTypeNameForTypeReference(typeReference);
            if (typeDefinition.isDiscriminatedChildTypeDefinition()) {
                parserCallString = "(" + getLanguageTypeNameForTypeReference(typeReference) + ") " + typeDefinition.getParentType().getName();
            }
            List<Term> paramTerms = complexTypeReference.getParams().orElse(Collections.emptyList());
            for (int i = 0; i < paramTerms.size(); i++) {
                Term paramTerm = paramTerms.get(i);
                paramsString
                    .append(", (")
                    .append(getLanguageTypeNameForTypeReference(getArgumentType(complexTypeReference, i), true))
                    .append(") (")
                    .append(toParseExpression(null, paramTerm, null))
                    .append(")");
            }
            return "new DataReaderComplexDefault<>(() -> " + parserCallString + "IO.staticParse(readBuffer" + paramsString + "), readBuffer)";
        } else {
            throw new IllegalStateException("What is this type? " + typeReference);
        }
    }

    public String getDataReaderCall(SimpleTypeReference simpleTypeReference) {
        final int sizeInBits = simpleTypeReference.getSizeInBits();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "readBoolean(readBuffer)";
            case BYTE:
                return "readByte(readBuffer, " + sizeInBits + ")";
            case UINT:
                if (sizeInBits <= 4) return "readUnsignedByte(readBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 8) return "readUnsignedShort(readBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 16) return "readUnsignedInt(readBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 32) return "readUnsignedLong(readBuffer, " + sizeInBits + ")";
                return "readUnsignedBigInteger(readBuffer, " + sizeInBits + ")";
            case INT:
                if (sizeInBits <= 8) return "readSignedByte(readBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 16) return "readSignedShort(readBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 32) return "readSignedInt(readBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 64) return "readSignedLong(readBuffer, " + sizeInBits + ")";
                return "readSignedBigInteger(readBuffer, " + sizeInBits + ")";
            case FLOAT:
                if (sizeInBits <= 32) return "readFloat(readBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 64) return "readDouble(readBuffer, " + sizeInBits + ")";
                return "readBigDecimal(readBuffer, " + sizeInBits + ")";
            case STRING:
                return "readString(readBuffer, " + sizeInBits + ")";
            case VSTRING:
                VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                return "readString(readBuffer, " + toParseExpression(null, vstringTypeReference.getLengthExpression(), null) + ")";
            default:
                throw new UnsupportedOperationException("Unsupported type " + simpleTypeReference.getBaseType());
        }
    }

    public String getDataWriterCall(TypeReference typeReference, String fieldName) {
        String resolverMethod = "enumForValue";
        if (isEnumTypeReference(typeReference)) {
            final String languageTypeName = getLanguageTypeNameForTypeReference(typeReference);
            final SimpleTypeReference enumBaseTypeReference = getEnumBaseTypeReference(typeReference);
            return "new DataWriterEnumDefault<>(" + languageTypeName + "::" + resolverMethod + ", " + getDataWriterCall(enumBaseTypeReference, fieldName) + ")";
        } else if (typeReference.isSimpleTypeReference()) {
            SimpleTypeReference simpleTypeReference = typeReference.asSimpleTypeReference().orElseThrow(IllegalStateException::new);
            return getDataWriterCall(simpleTypeReference, fieldName);
        } else if (typeReference.isComplexTypeReference()) {
            StringBuilder paramsString = new StringBuilder();
            ComplexTypeReference complexTypeReference = typeReference.asComplexTypeReference().orElseThrow(IllegalStateException::new);
            TypeDefinition typeDefinition = getTypeDefinitionForTypeReference(typeReference);
            String serializerCallString = getLanguageTypeNameForTypeReference(typeReference);
            if (typeDefinition.isDiscriminatedChildTypeDefinition()) {
                serializerCallString = "(" + getLanguageTypeNameForTypeReference(typeReference) + ") " + typeDefinition.getParentType().getName();
            }
            List<Term> paramTerms = complexTypeReference.getParams().orElse(Collections.emptyList());
            for (int i = 0; i < paramTerms.size(); i++) {
                Term paramTerm = paramTerms.get(i);
                paramsString
                    .append(", (")
                    .append(getLanguageTypeNameForTypeReference(getArgumentType(complexTypeReference, i), true))
                    .append(") (")
                    .append(toParseExpression(null, paramTerm, null))
                    .append(")");
            }
            return "new DataWriterComplexDefault<>(writeBuffer)";
        } else {
            throw new IllegalStateException("What is this type? " + typeReference);
        }
    }

    public String getDataWriterCall(SimpleTypeReference simpleTypeReference, String fieldName) {
        final int sizeInBits = simpleTypeReference.getSizeInBits();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "writeBoolean(writeBuffer)";
            case BYTE:
                return "writeByte(writeBuffer, " + sizeInBits + ")";
            case UINT:
                if (sizeInBits <= 4) return "writeUnsignedByte(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 8) return "writeUnsignedShort(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 16) return "writeUnsignedInt(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 32) return "writeUnsignedLong(writeBuffer, " + sizeInBits + ")";
                return "writeUnsignedBigInteger(writeBuffer, " + sizeInBits + ")";
            case INT:
                if (sizeInBits <= 8) return "writeSignedByte(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 16) return "writeSignedShort(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 32) return "writeSignedInt(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 64) return "writeSignedLong(writeBuffer, " + sizeInBits + ")";
                return "writeSignedBigInteger(writeBuffer, " + sizeInBits + ")";
            case FLOAT:
                if (sizeInBits <= 32) return "writeFloat(writeBuffer, " + sizeInBits + ")";
                if (sizeInBits <= 64) return "writeDouble(writeBuffer, " + sizeInBits + ")";
                return "writeBigDecimal(writeBuffer, " + sizeInBits + ")";
            case STRING:
                return "writeString(writeBuffer, " + sizeInBits + ")";
            case VSTRING:
                VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                return "writeString(writeBuffer, " + toParseExpression(null, vstringTypeReference.getLengthExpression(), null) + ")";
            default:
                throw new UnsupportedOperationException("Unsupported type " + simpleTypeReference.getBaseType());
        }
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
            case BIT:
                return "writeBuffer.writeBit(\"" + logicalName + "\", (boolean) " + fieldName + "" + writerArgsString + ")";
            case BYTE:
                ByteTypeReference byteTypeReference = (ByteTypeReference) simpleTypeReference;
                return "writeBuffer.writeByte(\"" + logicalName + "\", ((Number) " + fieldName + ").byteValue()" + writerArgsString + ")";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 4) {
                    return "writeBuffer.writeUnsignedByte(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").byteValue()" + writerArgsString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    return "writeBuffer.writeUnsignedShort(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").shortValue()" + writerArgsString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    return "writeBuffer.writeUnsignedInt(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").intValue()" + writerArgsString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    return "writeBuffer.writeUnsignedLong(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").longValue()" + writerArgsString + ")";
                }
                return "writeBuffer.writeUnsignedBigInteger(\"" + logicalName + "\", " + unsignedIntegerTypeReference.getSizeInBits() + ", (BigInteger) " + fieldName + "" + writerArgsString + ")";
            case INT:
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
            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                if (floatTypeReference.getSizeInBits() <= 32) {
                    return "writeBuffer.writeFloat(\"" + logicalName + "\", " + floatTypeReference.getSizeInBits() + "," + fieldName + "" + writerArgsString + ")";
                } else if (floatTypeReference.getSizeInBits() <= 64) {
                    return "writeBuffer.writeDouble(\"" + logicalName + "\", " + floatTypeReference.getSizeInBits() + "," + fieldName + "" + writerArgsString + ")";
                } else {
                    throw new RuntimeException("Unsupported float type");
                }
            case STRING:
            case VSTRING:
                final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                if (!(encodingTerm instanceof StringLiteral)) {
                    throw new RuntimeException("Encoding must be a quoted string value");
                }
                String encoding = ((StringLiteral) encodingTerm).getValue();
                String length = Integer.toString(simpleTypeReference.getSizeInBits());
                if (simpleTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.VSTRING) {
                    VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                    length = toSerializationExpression(field, vstringTypeReference.getLengthExpression(), thisType.getParserArguments().orElse(Collections.emptyList()));
                }
                return "writeBuffer.writeString(\"" + logicalName + "\", " + length + ", \"" +
                    encoding + "\", (String) " + fieldName + "" + writerArgsString + ")";
        }
        throw new FreemarkerException("Unmapped basetype" + simpleTypeReference.getBaseType());
    }

    public String getReservedValue(ReservedField reservedField) {
        final String languageTypeName = getLanguageTypeNameForTypeReference(reservedField.getType(), true);
        if ("BigInteger".equals(languageTypeName)) {
            return "BigInteger.valueOf(" + reservedField.getReferenceValue() + ")";
        } else {
            return "(" + languageTypeName + ") " + reservedField.getReferenceValue();
        }
    }

    public String toAccessExpression(TypedField field, Term term, List<Argument> parserArguments) {
        return toExpression(field, term, variableLiteral -> {
            if (isVariableLiteralVirtualField(variableLiteral) || isVariableLiteralDiscriminatorField(variableLiteral)) { // If we are accessing virtual|discriminator fields, we need to call the getter.
                return "get" + StringUtils.capitalize(variableLiteral.getName()) + "()";
            }
            return toVariableParseExpression(field, variableLiteral, parserArguments);
        });
    }

    public String toParseExpression(TypedField field, Term term, List<Argument> parserArguments) {
        Tracer tracer = Tracer.start("toParseExpression");
        return tracer + toExpression(field, term, variableLiteral -> tracer.dive("variableExpressionGenerator") + toVariableParseExpression(field, variableLiteral, parserArguments));
    }

    public String toSerializationExpression(TypedField field, Term term, List<Argument> serializerArguments) {
        Tracer tracer = Tracer.start("toSerializationExpression");
        return tracer + toExpression(field, term, variableLiteral -> tracer.dive("variableExpressionGenerator") + toVariableSerializationExpression(field, variableLiteral, serializerArguments));
    }

    private String toExpression(TypedField field, Term term, Function<VariableLiteral, String> variableExpressionGenerator) {
        Tracer tracer = Tracer.start("toExpression");
        if (term == null) {
            return tracer + "";
        }
        if (term instanceof Literal) {
            return toLiteralTermExpression((Literal) term, variableExpressionGenerator, tracer);
        } else if (term instanceof UnaryTerm) {
            return toUnaryTermExpression(field, (UnaryTerm) term, variableExpressionGenerator, tracer);
        } else if (term instanceof BinaryTerm) {
            return toBinaryTermExpression(field, (BinaryTerm) term, variableExpressionGenerator, tracer);
        } else if (term instanceof TernaryTerm) {
            return toTernaryTermExpression(field, (TernaryTerm) term, variableExpressionGenerator, tracer);
        } else {
            throw new RuntimeException("Unsupported Term type " + term.getClass().getName());
        }
    }

    private String toLiteralTermExpression(Literal literal, Function<VariableLiteral, String> variableExpressionGenerator, Tracer tracer) {
        tracer = tracer.dive("literal term instanceOf");
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
            if ("curPos".equals(((VariableLiteral) literal).getName())) {
                return "(readBuffer.getPos() - startPos)";
            }
            // If this literal references an Enum type, then we have to output it differently.
            if (getTypeDefinitions().get(variableLiteral.getName()) instanceof EnumTypeDefinition) {
                tracer = tracer.dive("enum definition instanceOf");
                VariableLiteral enumDefinitionChild = variableLiteral.getChild()
                    .orElseThrow(() -> new RuntimeException("enum definitions should have childs"));
                return tracer + variableLiteral.getName() + "." + enumDefinitionChild.getName() +
                    enumDefinitionChild.getChild().map(child -> "." + toVariableExpressionRest(child)).orElse("");
            } else {
                return tracer + variableExpressionGenerator.apply(variableLiteral);
            }
        } else {
            throw new RuntimeException("Unsupported Literal type " + literal.getClass().getName());
        }
    }

    private String toUnaryTermExpression(TypedField field, UnaryTerm unaryTerm, Function<VariableLiteral, String> variableExpressionGenerator, Tracer tracer) {
        tracer = tracer.dive("unary term instanceOf");
        Term a = unaryTerm.getA();
        switch (unaryTerm.getOperation()) {
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
                throw new RuntimeException("Unsupported unary operation type " + unaryTerm.getOperation());
        }
    }

    private String toBinaryTermExpression(TypedField field, BinaryTerm binaryTerm, Function<VariableLiteral, String> variableExpressionGenerator, Tracer tracer) {
        tracer = tracer.dive("binary term instanceOf");
        Term a = binaryTerm.getA();
        Term b = binaryTerm.getB();
        String operation = binaryTerm.getOperation();
        switch (operation) {
            case "^":
                tracer = tracer.dive("^");
                return tracer + "Math.pow((" + toExpression(field, a, variableExpressionGenerator) + "), (" + toExpression(field, b, variableExpressionGenerator) + "))";
            default:
                return tracer + "(" + toExpression(field, a, variableExpressionGenerator) + ") " + operation + " (" + toExpression(field, b, variableExpressionGenerator) + ")";
        }
    }

    private String toTernaryTermExpression(TypedField field, TernaryTerm ternaryTerm, Function<VariableLiteral, String> variableExpressionGenerator, Tracer tracer) {
        tracer = tracer.dive("ternary term instanceOf");
        if ("if".equals(ternaryTerm.getOperation())) {
            Term a = ternaryTerm.getA();
            Term b = ternaryTerm.getB();
            Term c = ternaryTerm.getC();
            return tracer + "((" + toExpression(field, a, variableExpressionGenerator) + ") ? " + toExpression(field, b, variableExpressionGenerator) + " : " + toExpression(field, c, variableExpressionGenerator) + ")";
        } else {
            throw new RuntimeException("Unsupported ternary operation type " + ternaryTerm.getOperation());
        }
    }

    public String toVariableEnumAccessExpression(VariableLiteral variableLiteral) {
        return variableLiteral.getName();
    }

    private String toVariableParseExpression(TypedField field, VariableLiteral variableLiteral, List<Argument> parserArguments) {
        Tracer tracer = Tracer.start("toVariableParseExpression");
        // CAST expressions are special as we need to add a ".class" to the second parameter in Java.
        if ("CAST".equals(variableLiteral.getName())) {
            return toCastVariableParseExpression(field, variableLiteral, parserArguments, tracer);
        } else if ("STATIC_CALL".equals(variableLiteral.getName())) {
            return toStaticCallVariableParseExpression(field, variableLiteral, parserArguments, tracer);
        } else if (isVariableLiteralImplicitField(variableLiteral)) { // If we are accessing implicit fields, we need to rely on a local variable instead.
            return toImplictVariableParseExpression(variableLiteral, tracer);
        } else if (variableLiteral.getName().equals(variableLiteral.getName().toUpperCase())) { // All uppercase names are not fields, but utility methods.
            return toUpperCaseVariableParseExpression(field, variableLiteral, parserArguments, tracer);
        }
        return tracer + variableLiteral.getName() + variableLiteral.getChild().map(child -> "." + toVariableExpressionRest(child)).orElse("");
    }

    private String toUpperCaseVariableParseExpression(TypedField field, VariableLiteral variableLiteral, List<Argument> parserArguments, Tracer tracer) {
        tracer = tracer.dive("UPPERCASE");
        StringBuilder sb = new StringBuilder(variableLiteral.getName());
        if (variableLiteral.getArgs().isPresent()) {
            sb.append("(");
            boolean firstArg = true;
            for (Term arg : variableLiteral.getArgs().get()) {
                if (!firstArg) {
                    sb.append(", ");
                }
                sb.append(toParseExpression(field, arg, parserArguments));
                firstArg = false;
            }
            sb.append(")");
        }
        if (variableLiteral.getIndex() != VariableLiteral.NO_INDEX) {
            sb.append(".get(").append(variableLiteral.getIndex()).append(")");
        }
        return tracer + sb.toString() + variableLiteral.getChild().map(child -> "." + toVariableExpressionRest(child)).orElse("");
    }

    private String toImplictVariableParseExpression(VariableLiteral variableLiteral, Tracer tracer) {
        tracer = tracer.dive("implicit");
        return tracer + variableLiteral.getName();
    }

    private String toStaticCallVariableParseExpression(TypedField field, VariableLiteral variableLiteral, List<Argument> parserArguments, Tracer tracer) {
        tracer = tracer.dive("STATIC_CALL");
        StringBuilder sb = new StringBuilder();
        List<Term> arguments = variableLiteral.getArgs().orElseThrow(() -> new RuntimeException("A STATIC_CALL expression needs arguments"));
        if (arguments.size() < 1) {
            throw new RuntimeException("A STATIC_CALL expression expects at least one argument.");
        }
        // TODO: make it as static import with a emitImport so if a static call is present a "utils" package must be present in the import
        sb.append(packageName()).append(".utils.StaticHelper.");
        // Get the class and method name
        String methodName = arguments.get(0).asLiteral()
            .orElseThrow(() -> new RuntimeException("First argument should be a literal"))
            .asStringLiteral()
            .orElseThrow(() -> new RuntimeException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral")).
            getValue();
        sb.append(methodName).append("(");
        for (int i = 1; i < arguments.size(); i++) {
            Term arg = arguments.get(i);
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
                    sb.append(variableLiteralArg.getName()).append(variableLiteralArg.getChild().map(child -> "." + toVariableExpressionRest(child)).orElse(""));
                } else if (isTypeArg) {// We have to manually evaluate the type information at code-generation time.
                    String part = variableLiteralArg.getChild().map(VariableLiteral::getName).orElse("");
                    switch (part) {
                        case "name":
                            sb.append("\"").append(field.getTypeName()).append("\"");
                            break;
                        case "length":
                            sb.append("\"").append(((SimpleTypeReference) field).getSizeInBits()).append("\"");
                            break;
                        case "encoding":
                            String encoding = ((StringLiteral) field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"))).getValue();
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
    }

    private String toCastVariableParseExpression(TypedField field, VariableLiteral variableLiteral, List<Argument> parserArguments, Tracer tracer) {
        tracer = tracer.dive("CAST");
        StringBuilder sb = new StringBuilder(variableLiteral.getName());
        List<Term> arguments = variableLiteral.getArgs().orElseThrow(() -> new RuntimeException("A Cast expression needs arguments"));
        if (arguments.size() != 2) {
            throw new RuntimeException("A CAST expression expects exactly two arguments.");
        }
        VariableLiteral firstArgument = arguments.get(0).asLiteral()
            .orElseThrow(() -> new RuntimeException("First argument should be a literal"))
            .asVariableLiteral()
            .orElseThrow(() -> new RuntimeException("First argument should be a Variable literal"));
        VariableLiteral secondArgument = arguments.get(1).asLiteral().orElseThrow(() -> new RuntimeException("Second argument should be a literal"))
            .asVariableLiteral()
            .orElseThrow(() -> new RuntimeException("Second argument should be a Variable literal"));
        sb.append("(")
            .append(toVariableParseExpression(field, firstArgument, parserArguments))
            .append(", ")
            .append(secondArgument.getName()).append(".class)");
        return tracer + sb.toString() + variableLiteral.getChild().map(child -> "." + toVariableExpressionRest(child)).orElse("");
    }

    private String toVariableSerializationExpression(TypedField field, VariableLiteral variableLiteral, List<Argument> serialzerArguments) {
        Tracer tracer = Tracer.start("variable serialization expression");
        if ("STATIC_CALL".equals(variableLiteral.getName())) {
            return toStaticCallSerializationExpression(field, variableLiteral, serialzerArguments, tracer);
        }
        // All uppercase names are not fields, but utility methods.
        else if (variableLiteral.getName().equals(variableLiteral.getName().toUpperCase())) {
            return toUpperCaseSerializationExpression(field, variableLiteral, serialzerArguments, tracer);
        } else if (isVariableLiteralImplicitField(variableLiteral)) { // If we are accessing implicit fields, we need to rely on a local variable instead.
            tracer = tracer.dive("implicit field");
            return tracer + toSerializationExpression(getReferencedImplicitField(variableLiteral), getReferencedImplicitField(variableLiteral).getSerializeExpression(), serialzerArguments);
        } else if (isVariableLiteralVirtualField(variableLiteral)) {
            tracer = tracer.dive("virtual field");
            return tracer + toVariableExpressionRest(variableLiteral);
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
            return tracer + variableLiteral.getName() + variableLiteral.getChild().map(child -> "." + toVariableExpressionRest(child)).orElse("");
        } else if (isTypeArg) {
            tracer = tracer.dive("type arg");
            String part = variableLiteral.getChild().map(VariableLiteral::getName).orElse("");
            switch (part) {
                case "name":
                    return tracer + "\"" + field.getTypeName() + "\"";
                case "length":
                    return tracer + "\"" + ((SimpleTypeReference) field).getSizeInBits() + "\"";
                case "encoding":
                    String encoding = ((StringLiteral) field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"))).getValue();
                    return tracer + "\"" + encoding + "\"";
                default:
                    return tracer + "";
            }
        } else {
            return tracer + toVariableExpressionRest(variableLiteral);
        }
    }

    private String toUpperCaseSerializationExpression(TypedField field, VariableLiteral variableLiteral, List<Argument> serialzerArguments, Tracer tracer) {
        tracer = tracer.dive("UPPER_CASE");
        StringBuilder sb = new StringBuilder(variableLiteral.getName());
        if (variableLiteral.getArgs().isPresent()) {
            sb.append("(");
            boolean firstArg = true;
            for (Term arg : variableLiteral.getArgs().get()) {
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
                        sb.append(va.getName()).append(va.getChild().map(child -> "." + toVariableExpressionRest(child)).orElse(""));
                    } else if (isTypeArg) {
                        String part = va.getChild().map(VariableLiteral::getName).orElse("");
                        switch (part) {
                            case "name":
                                sb.append("\"").append(field.getTypeName()).append("\"");
                                break;
                            case "length":
                                sb.append("\"").append(((SimpleTypeReference) field).getSizeInBits()).append("\"");
                                break;
                            case "encoding":
                                String encoding = ((StringLiteral) field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"))).getValue();
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
    }

    private String toStaticCallSerializationExpression(TypedField field, VariableLiteral variableLiteral, List<Argument> serialzerArguments, Tracer tracer) {
        tracer = tracer.dive("STATIC_CALL");
        StringBuilder sb = new StringBuilder();
        List<Term> arguments = variableLiteral.getArgs().orElseThrow(() -> new RuntimeException("A STATIC_CALL expression needs arguments"));
        if (arguments.size() < 1) {
            throw new RuntimeException("A STATIC_CALL expression expects at least one argument.");
        }
        // TODO: make it as static import with a emitImport so if a static call is present a "utils" package must be present in the import
        sb.append(packageName()).append(".utils.StaticHelper.");
        // Get the class and method name
        String methodName = arguments.get(0).asLiteral()
            .orElseThrow(() -> new RuntimeException("First argument should be a literal"))
            .asStringLiteral()
            .orElseThrow(() -> new RuntimeException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral")).
            getValue();
        //methodName = methodName.substring(1, methodName.length() - 1);
        sb.append(methodName).append("(");
        for (int i = 1; i < arguments.size(); i++) {
            Term arg = arguments.get(i);
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
                    sb.append(va.getName()).append(va.getChild().map(child -> "." + toVariableExpressionRest(child)).orElse(""));
                } else if (isTypeArg) {
                    String part = va.getChild().map(VariableLiteral::getName).orElse("");
                    switch (part) {
                        case "name":
                            sb.append("\"").append(field.getTypeName()).append("\"");
                            break;
                        case "length":
                            sb.append("\"").append(((SimpleTypeReference) field).getSizeInBits()).append("\"");
                            break;
                        case "encoding":
                            String encoding = ((StringLiteral) field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"))).getValue();
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

    private String toVariableExpressionRest(VariableLiteral variableLiteral) {
        Tracer tracer = Tracer.start("variable expression rest");
        // length is kind of a keyword in mspec, so we shouldn't be naming variables length. if we ask for the length of a object we can just return length().
        // This way we can get the length of a string when serializing
        String variableLiteralName = variableLiteral.getName();
        if (variableLiteralName.equals("length")) {
            tracer = tracer.dive("length");
            return tracer + variableLiteralName + "()" + ((variableLiteral.isIndexed() ? ".get(" + variableLiteral.getIndex() + ")" : "") +
                variableLiteral.getChild().map(child -> "." + toVariableExpressionRest(child)).orElse(""));
        }
        return tracer + "get" + WordUtils.capitalize(variableLiteralName) + "()" + ((variableLiteral.isIndexed() ? ".get(" + variableLiteral.getIndex() + ")" : "") +
            variableLiteral.getChild().map(child -> "." + toVariableExpressionRest(child)).orElse(""));
    }

    public String getSizeInBits(ComplexTypeDefinition complexTypeDefinition, List<Argument> parserArguments) {
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
                    if (simpleTypeReference instanceof VstringTypeReference) {
                        sb.append(toSerializationExpression(null, ((VstringTypeReference) simpleTypeReference).getLengthExpression(), parserArguments)).append(" + ");
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
                case VSTRING:
                    return "\"" + valueString + "\"";
            }
        }
        return valueString;
    }

}
