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

package org.apache.plc4x.codegeneration.language.java.generators;

import com.palantir.javapoet.*;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultDataIoTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.*;
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.references.*;
import org.apache.plc4x.plugins.codegenerator.types.terms.*;
import org.apache.plc4x.java.spi.utils.StaticHelper;

import java.io.File;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;
import java.util.function.Function;

public abstract class BaseGenerator<T> {

    public static final TypeReference BOOL_TYPE_REFERENCE = new DefaultBooleanTypeReference();
    public static final TypeReference INT_TYPE_REFERENCE = new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.INT, 32);
    public static final TypeReference STRING_TYPE_REFERENCE = new DefaultVstringTypeReference(SimpleTypeReference.SimpleBaseType.VSTRING, null);
    public static final TypeReference ANY_TYPE_REFERENCE = new DefaultUndefinedTypeReference();

    /**
     * Apache License header to be added to all generated files.
     */
    public static final String APACHE_LICENSE_HEADER =
        """
            Licensed to the Apache Software Foundation (ASF) under one
            or more contributor license agreements.  See the NOTICE file
            distributed with this work for additional information
            regarding copyright ownership.  The ASF licenses this file
            to you under the Apache License, Version 2.0 (the
            "License"); you may not use this file except in compliance
            with the License.  You may obtain a copy of the License at
            
              https://www.apache.org/licenses/LICENSE-2.0
            
            Unless required by applicable law or agreed to in writing,
            software distributed under the License is distributed on an
            "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
            KIND, either express or implied.  See the License for the
            specific language governing permissions and limitations
            under the License.""";

    protected final String targetPackage;
    protected final File outputDirectory;
    protected final File outputRootDir;
    protected final Map<String, TypeDefinition> types;
    protected final Map<String, String> externalTypes;
    protected final Map<String, Object> options;

    public BaseGenerator(String targetPackage, File outputDirectory, File outputRootDir, Map<String, TypeDefinition> types, Map<String, String> externalTypes, Map<String, Object> options) {
        this.targetPackage = targetPackage;
        this.outputDirectory = outputDirectory;
        this.outputRootDir = outputRootDir;
        this.types = types;
        this.externalTypes = externalTypes;
        this.options = options;
    }

    public abstract void generate(T definition) throws BufferException;

    public boolean isGeneratePropertiesForReservedFields() {
        return options.getOrDefault("generate-properties-for-reserved-fields", "false").equals("true");
    }

    public TypeName getLanguageTypeNameForTypeReference(TypeReference typeReference, boolean allowPrimitive) {
        Objects.requireNonNull(typeReference);
        if (typeReference instanceof ArrayTypeReference arrayTypeReference) {
            if (arrayTypeReference.getElementTypeReference().isByteBased()) {
                ClassName list = ClassName.get("java.util", "List");
                TypeName genericTypeByte = ClassName.get(Byte.class);
                return ParameterizedTypeName.get(list, genericTypeByte);
            } else {
                ClassName list = ClassName.get("java.util", "List");
                TypeName genericType = getLanguageTypeNameForTypeReference(arrayTypeReference.getElementTypeReference(), false);
                return ParameterizedTypeName.get(list, genericType);
            }
        }
        // DataIo data-types always have properties of type PlcValue
        if (typeReference.isDataIoTypeReference()) {
            return ClassName.get(PlcValue.class);
        }
        if (typeReference.isNonSimpleTypeReference()) {
            return ClassName.get(targetPackage, typeReference.asNonSimpleTypeReference().orElseThrow().getName());
        }
        SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return allowPrimitive ? TypeName.BOOLEAN : ClassName.get(Boolean.class);
            case BYTE:
                return allowPrimitive ? TypeName.BYTE : ClassName.get(Byte.class);
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 7) {
                    return allowPrimitive ? TypeName.BYTE : ClassName.get(Byte.class);
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 15) {
                    return allowPrimitive ? TypeName.SHORT : ClassName.get(Short.class);
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 31) {
                    return allowPrimitive ? TypeName.INT : ClassName.get(Integer.class);
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 63) {
                    return allowPrimitive ? TypeName.LONG : ClassName.get(Long.class);
                }
                return ClassName.get(BigInteger.class);
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return allowPrimitive ? TypeName.BYTE : ClassName.get(Byte.class);
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return allowPrimitive ? TypeName.SHORT : ClassName.get(Short.class);
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return allowPrimitive ? TypeName.INT : ClassName.get(Integer.class);
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return allowPrimitive ? TypeName.LONG : ClassName.get(Long.class);
                }
                return ClassName.get(BigInteger.class);
            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                int sizeInBits = floatTypeReference.getSizeInBits();
                if (sizeInBits <= 32) {
                    return allowPrimitive ? TypeName.FLOAT : ClassName.get(Float.class);
                }
                if (sizeInBits <= 64) {
                    return allowPrimitive ? TypeName.DOUBLE : ClassName.get(Double.class);
                }
                return ClassName.get(BigDecimal.class);
            case STRING:
            case VSTRING:
                return ClassName.get(String.class);
            case TIME:
                return ClassName.get(LocalTime.class);
            case DATE:
                return ClassName.get(LocalDate.class);
            case DATETIME:
                return ClassName.get(LocalDateTime.class);
            case UNDEFINED:
                return ClassName.get("java.lang", "Object");
        }
        throw new RuntimeException("Unsupported simple type: " + simpleTypeReference.getBaseType());
    }

    public String getNullValueForTypeReference(TypeReference typeReference) {
        if (typeReference instanceof SimpleTypeReference simpleTypeReference) {
            switch (simpleTypeReference.getBaseType()) {
                case BIT:
                    return "false";
                case BYTE:
                    return "0";
                case UINT:
                    IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (unsignedIntegerTypeReference.getSizeInBits() <= 31) {
                        return "0";
                    }
                    if (unsignedIntegerTypeReference.getSizeInBits() <= 63) {
                        return "0L";
                    }
                    return "null";
                case INT:
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return "0";
                    }
                    if (integerTypeReference.getSizeInBits() <= 64) {
                        return "0L";
                    }
                    return "null";
                case FLOAT:
                    FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                    int sizeInBits = floatTypeReference.getSizeInBits();
                    if (sizeInBits <= 32) {
                        return "0.0F";
                    }
                    if (sizeInBits <= 64) {
                        return "0.0";
                    }
                    return "null";
                case STRING:
                case VSTRING:
                    return "null";
            }
            throw new RuntimeException("Unmapped base-type" + simpleTypeReference.getBaseType());
        } else {
            return "null";
        }
    }

    protected TypeName toTypeName(TypeReference typeReference, boolean allowPrimitive) {
        Objects.requireNonNull(typeReference);
        if (typeReference instanceof ArrayTypeReference arrayTypeReference) {
            TypeReference elemRef = arrayTypeReference.getElementTypeReference();
            if (elemRef.isByteBased()) {
                // byte[]
                return ArrayTypeName.of(TypeName.BYTE);
            } else {
                // List<Elem>
                return ParameterizedTypeName.get(
                    ClassName.get("java.util", "List"),
                    toTypeName(elemRef, false)
                );
            }
        }
        if (typeReference.isDataIoTypeReference()) {
            return ClassName.get("org.apache.plc4x.java.api.value", "PlcValue");
        }
        if (typeReference.isNonSimpleTypeReference()) {
            String name = typeReference.asNonSimpleTypeReference().orElseThrow().getName();
            if (externalTypes != null && externalTypes.containsKey(name)) {
                return classNameFromFqn(externalTypes.get(name));
            }
            return ClassName.get(targetPackage, name);
        }
        // Simple types
        SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return allowPrimitive ? TypeName.BOOLEAN : ClassName.get("java.lang", "Boolean");
            case BYTE:
                return allowPrimitive ? TypeName.BYTE : ClassName.get("java.lang", "Byte");
            case UINT: {
                IntegerTypeReference u = (IntegerTypeReference) simpleTypeReference;
                if (u.getSizeInBits() <= 7) {
                    return allowPrimitive ? TypeName.BYTE : ClassName.get("java.lang", "Byte");
                } else if (u.getSizeInBits() <= 15) {
                    return allowPrimitive ? TypeName.SHORT : ClassName.get("java.lang", "Short");
                } else if (u.getSizeInBits() <= 31) {
                    return allowPrimitive ? TypeName.INT : ClassName.get("java.lang", "Integer");
                } else if (u.getSizeInBits() <= 63) {
                    return allowPrimitive ? TypeName.LONG : ClassName.get("java.lang", "Long");
                } else {
                    return ClassName.get("java.math", "BigInteger");
                }
            }
            case INT: {
                IntegerTypeReference i = (IntegerTypeReference) simpleTypeReference;
                if (i.getSizeInBits() <= 8) {
                    return allowPrimitive ? TypeName.BYTE : ClassName.get("java.lang", "Byte");
                } else if (i.getSizeInBits() <= 16) {
                    return allowPrimitive ? TypeName.SHORT : ClassName.get("java.lang", "Short");
                } else if (i.getSizeInBits() <= 32) {
                    return allowPrimitive ? TypeName.INT : ClassName.get("java.lang", "Integer");
                } else if (i.getSizeInBits() <= 64) {
                    return allowPrimitive ? TypeName.LONG : ClassName.get("java.lang", "Long");
                } else {
                    return ClassName.get("java.math", "BigInteger");
                }
            }
            case FLOAT:
            case UFLOAT: {
                FloatTypeReference f = (FloatTypeReference) simpleTypeReference;
                int sizeInBits = f.getSizeInBits();
                if (sizeInBits <= 32) {
                    return allowPrimitive ? TypeName.FLOAT : ClassName.get("java.lang", "Float");
                } else if (sizeInBits <= 64) {
                    return allowPrimitive ? TypeName.DOUBLE : ClassName.get("java.lang", "Double");
                } else {
                    return ClassName.get("java.math", "BigDecimal");
                }
            }
            case STRING:
            case VSTRING:
                return ClassName.get("java.lang", "String");
            case TIME:
                return ClassName.get("java.time", "LocalTime");
            case DATE:
                return ClassName.get("java.time", "LocalDate");
            case DATETIME:
                return ClassName.get("java.time", "LocalDateTime");
        }
        throw new RuntimeException("Unsupported simple type");
    }

    protected ClassName classNameFromFqn(String fqn) {
        int idx = fqn.lastIndexOf('.');
        if (idx < 0) {
            return ClassName.bestGuess(fqn);
        }
        String pkg = fqn.substring(0, idx);
        String simple = fqn.substring(idx + 1);
        return ClassName.get(pkg, simple);
    }

    public String getPlcValueTypeForTypeReference(TypeReference typeReference) {
        if (typeReference.isArrayTypeReference() && typeReference.asArrayTypeReference().orElseThrow().getElementTypeReference().isByteBased()) {
            return "PlcRawByteArray";
        }
        if (!(typeReference instanceof SimpleTypeReference simpleTypeReference)) {
            return "PlcStruct";
        }
        int sizeInBits = simpleTypeReference.getSizeInBits();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "PlcBOOL";
            case BYTE:
                return "PlcSINT";
            case UINT:
                if (sizeInBits <= 8) {
                    return "PlcUSINT";
                }
                if (sizeInBits <= 16) {
                    return "PlcUINT";
                }
                if (sizeInBits <= 32) {
                    return "PlcUDINT";
                }
                if (sizeInBits <= 64) {
                    return "PlcULINT";
                }
                throw new IllegalArgumentException("Unsupported UINT with bit length " + sizeInBits);
            case INT:
                if (sizeInBits <= 8) {
                    return "PlcSINT";
                }
                if (sizeInBits <= 16) {
                    return "PlcINT";
                }
                if (sizeInBits <= 32) {
                    return "PlcDINT";
                }
                if (sizeInBits <= 64) {
                    return "PlcLINT";
                }
                throw new IllegalArgumentException("Unsupported INT with bit length " + sizeInBits);
            case FLOAT:
            case UFLOAT:
                if (sizeInBits <= 32) {
                    return "PlcREAL";
                }
                if (sizeInBits <= 64) {
                    return "PlcLREAL";
                }
                throw new IllegalArgumentException("Unsupported REAL with bit length " + sizeInBits);
            case STRING:
            case VSTRING:
                return "PlcSTRING";
            case TIME:
            case DATE:
            case DATETIME:
                return "PlcTIME";
        }
        throw new IllegalStateException("Unsupported simple type");
    }

    public CodeBlock getDataReaderCall(TypeReference typeReference) throws BufferException {
        return getDataReaderCall(typeReference, "enumForValue");
    }

    public CodeBlock getDataReaderCall(TypeReference typeReference, String resolverMethod) throws BufferException {
        ClassName dataReaderFactory = ClassName.get("org.apache.plc4x.java.spi.fields.data.reader", "DataReaderFactory");
        if (typeReference.isEnumTypeReference()) {
            final TypeName languageTypeName = getLanguageTypeNameForTypeReference(typeReference, false);
            final SimpleTypeReference enumBaseTypeReference = getEnumBaseTypeReference(typeReference);
            return CodeBlock.of("$T.readEnum($T::$L, $T.$L)", dataReaderFactory, languageTypeName, resolverMethod, dataReaderFactory, getDataReaderCall(enumBaseTypeReference));
        } else if (typeReference.isArrayTypeReference()) {
            final ArrayTypeReference arrayTypeReference = typeReference.asArrayTypeReference().orElseThrow();
            return getDataReaderCall(arrayTypeReference.getElementTypeReference(), resolverMethod);
        } else if (typeReference.isSimpleTypeReference()) {
            SimpleTypeReference simpleTypeReference = typeReference.asSimpleTypeReference().orElseThrow(IllegalStateException::new);
            return CodeBlock.of("$T.$L", dataReaderFactory, getDataReaderCall(simpleTypeReference));
        } else if (typeReference.isComplexTypeReference()) {
            StringBuilder paramsString = new StringBuilder();
            ComplexTypeReference complexTypeReference = typeReference.asComplexTypeReference().orElseThrow(IllegalStateException::new);
            ComplexTypeDefinition typeDefinition = complexTypeReference.getTypeDefinition();
            TypeName parserResultTypeString = getLanguageTypeNameForTypeReference(typeReference, false);
            TypeName parserTypeString = complexTypeReference.isDataIoTypeReference() ? ClassName.get(targetPackage, complexTypeReference.getName()) : getLanguageTypeNameForTypeReference(typeReference, false);
            List<Term> paramTerms = complexTypeReference.getParams().orElse(Collections.emptyList());
            for (int i = 0; i < paramTerms.size(); i++) {
                Term paramTerm = paramTerms.get(i);
                final TypeReference argumentType = getArgumentType(complexTypeReference, i);
                paramsString
                    .append(", (")
                    .append(getLanguageTypeNameForTypeReference(argumentType, true))
                    .append(") (")
                    .append(toParseExpression(null, null, argumentType, paramTerm, null))
                    .append(")");
            }
            return CodeBlock.of("$T.readComplex(() -> ($T) $T.staticParse(readBuffer" + paramsString + "), readBuffer)", dataReaderFactory, parserResultTypeString, parserTypeString);
        } else {
            throw new IllegalStateException("What is this type? " + typeReference);
        }
    }

    public CodeBlock getDataReaderCall(SimpleTypeReference simpleTypeReference) throws BufferException {
        final int sizeInBits = simpleTypeReference.getSizeInBits();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return CodeBlock.of("readBoolean(readBuffer)");
            case BYTE:
                return CodeBlock.of("readByte(readBuffer, " + sizeInBits + ")");
            case UINT:
                if (sizeInBits <= 7) return CodeBlock.of("readUnsignedByte(readBuffer, " + sizeInBits + ")");
                if (sizeInBits <= 15) return CodeBlock.of("readUnsignedShort(readBuffer, " + sizeInBits + ")");
                if (sizeInBits <= 31) return CodeBlock.of("readUnsignedInt(readBuffer, " + sizeInBits + ")");
                if (sizeInBits <= 63) return CodeBlock.of("readUnsignedLong(readBuffer, " + sizeInBits + ")");
                return CodeBlock.of("readUnsignedBigInteger(readBuffer, " + sizeInBits + ")");
            case INT:
                if (sizeInBits <= 8) return CodeBlock.of("readSignedByte(readBuffer, " + sizeInBits + ")");
                if (sizeInBits <= 16) return CodeBlock.of("readSignedShort(readBuffer, " + sizeInBits + ")");
                if (sizeInBits <= 32) return CodeBlock.of("readSignedInt(readBuffer, " + sizeInBits + ")");
                if (sizeInBits <= 64) return CodeBlock.of("readSignedLong(readBuffer, " + sizeInBits + ")");
                return CodeBlock.of("readSignedBigInteger(readBuffer, " + sizeInBits + ")");
            case FLOAT:
                if (sizeInBits <= 32) return CodeBlock.of("readFloat(readBuffer, " + sizeInBits + ")");
                if (sizeInBits <= 64) return CodeBlock.of("readDouble(readBuffer, " + sizeInBits + ")");
                return CodeBlock.of("readBigDecimal(readBuffer, " + sizeInBits + ")");
            case STRING:
                return CodeBlock.of("readString(readBuffer, " + sizeInBits + ")");
            case VSTRING:
                VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                return CodeBlock.of("readString(readBuffer, " + toParseExpression(null, null, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), null) + ")");
            case TIME:
                return CodeBlock.of("readTime(readBuffer)");
            case DATE:
                return CodeBlock.of("readDate(readBuffer)");
            case DATETIME:
                return CodeBlock.of("readDateTime(readBuffer)");
            default:
                throw new UnsupportedOperationException("Unsupported type " + simpleTypeReference.getBaseType());
        }
    }

    public CodeBlock getDataWriterCall(TypeReference typeReference) throws BufferException {
        if (typeReference.isSimpleTypeReference()) {
            SimpleTypeReference simpleTypeReference = typeReference.asSimpleTypeReference().orElseThrow(IllegalStateException::new);
            return getDataWriterCall(simpleTypeReference);
        } else if (typeReference.isArrayTypeReference()) {
            final ArrayTypeReference arrayTypeReference = typeReference.asArrayTypeReference().orElseThrow();
            return getDataWriterCall(arrayTypeReference.getElementTypeReference());
        } else if (typeReference.isComplexTypeReference()) {
            ClassName dataWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.data.writer", "DataWriterFactory");
            return CodeBlock.of("$T.writeComplex(writeBuffer)", dataWriterFactory);
        } else if (typeReference.isEnumTypeReference()) {
            final EnumTypeReference enumTypeReference = typeReference.asEnumTypeReference().orElseThrow(IllegalStateException::new);
            SimpleTypeReference enumBaseTypeReference = enumTypeReference.getBaseTypeReference().orElseThrow();
            return getDataWriterCall(enumBaseTypeReference);
        } else {
            throw new IllegalStateException("What is this type? " + typeReference);
        }
    }

    public CodeBlock getDataWriterCall(SimpleTypeReference simpleTypeReference) throws BufferException {
        ClassName dataWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.data.writer", "DataWriterFactory");
        final int sizeInBits = simpleTypeReference.getSizeInBits();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return CodeBlock.of("$T.writeBoolean(writeBuffer)", dataWriterFactory);
            case BYTE:
                return CodeBlock.of("$T.writeByte(writeBuffer, $L)", dataWriterFactory, sizeInBits);
            case UINT:
                if (sizeInBits <= 7)
                    return CodeBlock.of("$T.writeUnsignedByte(writeBuffer, " + sizeInBits + ")", dataWriterFactory);
                if (sizeInBits <= 15)
                    return CodeBlock.of("$T.writeUnsignedShort(writeBuffer, " + sizeInBits + ")", dataWriterFactory);
                if (sizeInBits <= 31)
                    return CodeBlock.of("$T.writeUnsignedInt(writeBuffer, " + sizeInBits + ")", dataWriterFactory);
                if (sizeInBits <= 63)
                    return CodeBlock.of("$T.writeUnsignedLong(writeBuffer, " + sizeInBits + ")", dataWriterFactory);
                return CodeBlock.of("$T.writeUnsignedBigInteger(writeBuffer, " + sizeInBits + ")", dataWriterFactory);
            case INT:
                if (sizeInBits <= 8)
                    return CodeBlock.of("$T.writeSignedByte(writeBuffer, " + sizeInBits + ")", dataWriterFactory);
                if (sizeInBits <= 16)
                    return CodeBlock.of("$T.writeSignedShort(writeBuffer, " + sizeInBits + ")", dataWriterFactory);
                if (sizeInBits <= 32)
                    return CodeBlock.of("$T.writeSignedInt(writeBuffer, " + sizeInBits + ")", dataWriterFactory);
                if (sizeInBits <= 64)
                    return CodeBlock.of("$T.writeSignedLong(writeBuffer, " + sizeInBits + ")", dataWriterFactory);
                return CodeBlock.of("$T.writeSignedBigInteger(writeBuffer, " + sizeInBits + ")", dataWriterFactory);
            case FLOAT:
                if (sizeInBits <= 32)
                    return CodeBlock.of("$T.writeFloat(writeBuffer, " + sizeInBits + ")", dataWriterFactory);
                if (sizeInBits <= 64)
                    return CodeBlock.of("$T.writeDouble(writeBuffer, " + sizeInBits + ")", dataWriterFactory);
                return CodeBlock.of("$T.writeBigDecimal(writeBuffer, " + sizeInBits + ")", dataWriterFactory);
            case STRING:
                return CodeBlock.of("$T.writeString(writeBuffer, " + sizeInBits + ")", dataWriterFactory);
            case VSTRING:
                VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                return CodeBlock.of("$T.writeString(writeBuffer, " + toParseExpression(null, null, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), null) + ")", dataWriterFactory);
            case TIME:
                return CodeBlock.of("$T.writeTime(writeBuffer)", dataWriterFactory);
            case DATE:
                return CodeBlock.of("$T.writeDate(writeBuffer)", dataWriterFactory);
            case DATETIME:
                return CodeBlock.of("$T.writeDateTime(writeBuffer)", dataWriterFactory);
            default:
                throw new UnsupportedOperationException("Unsupported type " + simpleTypeReference.getBaseType());
        }
    }

    public Map<String, TypeReference> getDiscriminatorTypes(TypeDefinition type) {
        // Get the parent type (Which contains the typeSwitch field)
        SwitchField switchField = null;
        Function<String, TypeReference> typeRefRetriever = null;
        if (type.isDiscriminatedComplexTypeDefinition()) {
            DiscriminatedComplexTypeDefinition discriminatedComplexTypeDefinition = type.asDiscriminatedComplexTypeDefinition().orElseThrow();
            switchField = discriminatedComplexTypeDefinition.getSwitchField().orElse(null);
            typeRefRetriever = propertyName -> discriminatedComplexTypeDefinition.getTypeReferenceForProperty(propertyName).orElse(null);
            // Please forgive us, we didn't know what we were doing.
            if (switchField == null) {
                ComplexTypeDefinition parentType = type.asDiscriminatedComplexTypeDefinition().orElseThrow().getParentType().orElseThrow();
                switchField = parentType.getSwitchField().orElse(null);
                typeRefRetriever = propertyName -> parentType.getTypeReferenceForProperty(propertyName).orElse(null);
            }
        } else if (type.isDataIoTypeDefinition()) {
            final DefaultDataIoTypeDefinition dataIoTypeDefinition = (DefaultDataIoTypeDefinition) type;
            switchField = dataIoTypeDefinition.getSwitchField().orElseThrow();
            typeRefRetriever = propertyName -> type.getParserArguments()
                .orElse(Collections.emptyList())
                .stream()
                .filter(argument -> argument.getName().equals(propertyName))
                .findFirst()
                .map(Argument::getType)
                .orElse(null);
        } else if (type.isComplexTypeDefinition()) {
            switchField = ((ComplexTypeDefinition) type).getSwitchField().orElse(null);
            typeRefRetriever = propertyName -> ((ComplexTypeDefinition) type).getTypeReferenceForProperty(propertyName).orElse(null);
        }
        // Get the typeSwitch field from that.
        if (switchField == null) {
            return Collections.emptyMap();
        }
        Map<String, TypeReference> discriminatorTypes = new TreeMap<>();
        for (VariableLiteral variableLiteral : switchField.getDiscriminatorExpressions()) {
            // Get some symbolic name we can use.
            String discriminatorName = variableLiteral.getDiscriminatorName();
            final TypeReference typeReference = typeRefRetriever.apply(variableLiteral.getName());
            Optional<TypeReference> discriminatorType = typeReference.getDiscriminatorType(variableLiteral);
            if (discriminatorType.isEmpty()) {
                throw new RuntimeException("no type for " + discriminatorName);
            }
            discriminatorTypes.put(discriminatorName, discriminatorType.orElse(null));
        }
        return discriminatorTypes;
    }

    public TypeReference getArgumentType(TypeReference typeReference, int index) {
        Objects.requireNonNull(typeReference, "type reference must not be null");
        NonSimpleTypeReference complexTypeReference = typeReference.asNonSimpleTypeReference().orElseThrow(() -> new RuntimeException("Only non simple type references supported here."));
        return complexTypeReference.getArgumentType(index);
    }

    /**
     * @param field           this generally only is needed to access field attributes such as encoding etc.
     * @param resultType      the type the resulting expression should have
     * @param term            the term representing the expression
     * @param parserArguments any parser arguments, which could be referenced in expressions (Needed for getting the type)
     * @return Java code which does the things defined in 'term'
     */
    public CodeBlock toParseExpression(TypeDefinition typeDefinition, Field field, TypeReference resultType, Term term, List<Argument> parserArguments) throws BufferException {
        return toExpression(typeDefinition, field, resultType, term, variableLiteral -> toVariableParseExpression(typeDefinition, field, resultType, variableLiteral, parserArguments), true);
    }

    /**
     * @param field               this generally only is needed to access field attributes such as encoding etc.
     * @param resultType          the type the resulting expression should have
     * @param term                the term representing the expression
     * @param serializerArguments any serializer arguments, which could be referenced in expressions (Needed for getting the type)
     * @return Java code which does the things defined in 'term'
     */
    public CodeBlock toSerializationExpression(TypeDefinition typeDefinition, Field field, TypeReference resultType, Term term, List<Argument> serializerArguments) throws BufferException {
        return toExpression(typeDefinition, field, resultType, term, variableLiteral -> toVariableSerializationExpression(typeDefinition, field, resultType, variableLiteral, serializerArguments), false);
    }

    private CodeBlock toExpression(TypeDefinition typeDefinition, Field field, TypeReference resultType, Term term, ThrowingFunction<VariableLiteral, CodeBlock, BufferException> variableExpressionGenerator, boolean isParse) throws BufferException {
        return switch (term) {
            case null -> CodeBlock.of("");
            case Literal literal ->
                toLiteralTermExpression(typeDefinition, field, resultType, literal, variableExpressionGenerator, isParse);
            case UnaryTerm unaryTerm ->
                toUnaryTermExpression(typeDefinition, field, resultType, unaryTerm, variableExpressionGenerator, isParse);
            case BinaryTerm binaryTerm ->
                toBinaryTermExpression(typeDefinition, field, resultType, binaryTerm, variableExpressionGenerator, isParse);
            case TernaryTerm ternaryTerm ->
                toTernaryTermExpression(typeDefinition, field, resultType, ternaryTerm, variableExpressionGenerator, isParse);
            default ->
                throw new RuntimeException("Unsupported Term type " + term.getClass().getName() + ". Actual type " + resultType);
        };
    }

    private CodeBlock toLiteralTermExpression(TypeDefinition typeDefinition, Field field, TypeReference resultType, Literal literal, ThrowingFunction<VariableLiteral, CodeBlock, BufferException> variableExpressionGenerator, boolean isParse) throws BufferException {
        switch (literal) {
            case NullLiteral ignored -> {
                return CodeBlock.of("null");
            }
            case BooleanLiteral booleanLiteral -> {
                return CodeBlock.of(Boolean.toString(booleanLiteral.getValue()));
            }
            case NumericLiteral numericLiteral -> {
                final String numberString = numericLiteral.getNumber().toString();
                if (resultType.isIntegerTypeReference()) {
                    final IntegerTypeReference integerTypeReference = resultType.asIntegerTypeReference().orElseThrow(RuntimeException::new);
                    if (integerTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.UINT && integerTypeReference.getSizeInBits() >= 32) {
                        return CodeBlock.of("$LL", numberString);
                    } else if (integerTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.INT && integerTypeReference.getSizeInBits() > 32) {
                        return CodeBlock.of("$LL", numberString);
                    }
                } else if (resultType.isFloatTypeReference()) {
                    final FloatTypeReference floatTypeReference = resultType.asFloatTypeReference().orElseThrow(RuntimeException::new);
                    if (floatTypeReference.getSizeInBits() <= 32) {
                        return CodeBlock.of("$LF", numberString);
                    }
                }
                return CodeBlock.of("$L", numberString);
            }
            case HexadecimalLiteral hexadecimalLiteral -> {
                final String hexString = hexadecimalLiteral.getHexString();
                if (resultType.isIntegerTypeReference()) {
                    final IntegerTypeReference integerTypeReference = resultType.asIntegerTypeReference().orElseThrow(RuntimeException::new);
                    if (integerTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.UINT && integerTypeReference.getSizeInBits() >= 32) {
                        return CodeBlock.of("$LL", hexString);
                    } else if (integerTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.INT && integerTypeReference.getSizeInBits() > 32) {
                        return CodeBlock.of("$LL", hexString);
                    }
                }
                return CodeBlock.of("$L", hexString);
            }
            case StringLiteral stringLiteral -> {
                return CodeBlock.of("$S", stringLiteral.getValue());
            }
            case VariableLiteral variableLiteral -> {
                if ("curPos".equals(variableLiteral.getName())) {
                    if (isParse) {
                        return CodeBlock.of("(readBuffer.getPositionInBits() - startPos)");
                    } else {
                        return CodeBlock.of("(writeBuffer.getPositionInBits() - startPos)");
                    }
                }
                // If this literal references an Enum type, then we have to output it differently.
                if (types.get(variableLiteral.getName()) instanceof EnumTypeDefinition enumTypeDefinition) {
                    VariableLiteral enumDefinitionChild = variableLiteral.getChild()
                        .orElseThrow(() -> new RuntimeException("enum definitions should have children"));
                    ClassName enumTypeClassName = ClassName.get(targetPackage, enumTypeDefinition.getName());
                    CodeBlock.Builder codeBlock = CodeBlock.builder().add("$T.$L", enumTypeClassName, enumDefinitionChild.getName());
                    if (enumDefinitionChild.getChild().isPresent()) {
                        codeBlock.add(".$L", toPropertyVariableExpression(field, resultType, enumDefinitionChild.getChild().get()));
                    }
                    return codeBlock.build();
                } else {
                    return variableExpressionGenerator.apply(variableLiteral);
                }
                // If this literal references an Enum type, then we have to output it differently.
            }
            default -> throw new RuntimeException("Unsupported Literal type " + literal.getClass().getName());
        }
    }

    private CodeBlock toUnaryTermExpression(TypeDefinition typeDefinition, Field field, TypeReference resultType, UnaryTerm unaryTerm, ThrowingFunction<VariableLiteral, CodeBlock, BufferException> variableExpressionGenerator, boolean isParse) throws BufferException {
        Term a = unaryTerm.getA();
        return switch (unaryTerm.getOperation()) {
            case "!" -> {
                if ((resultType != ANY_TYPE_REFERENCE) && !resultType.isBooleanTypeReference()) {
                    throw new IllegalArgumentException("'!(...)' expression requires boolean type. Actual type " + resultType);
                }
                yield CodeBlock.of("!($L)",
                    toExpression(typeDefinition, field, resultType, a, variableExpressionGenerator, isParse));
            }
            case "-" -> {
                if ((resultType != ANY_TYPE_REFERENCE) && !resultType.isIntegerTypeReference() && !resultType.isFloatTypeReference()) {
                    throw new IllegalArgumentException("'-(...)' expression requires integer or floating-point type. Actual type " + resultType);
                }
                yield CodeBlock.of("-($L)",
                    toExpression(typeDefinition, field, resultType, a, variableExpressionGenerator, isParse));
            }
            case "()" ->
                CodeBlock.of("($L)", toExpression(typeDefinition, field, resultType, a, variableExpressionGenerator, isParse));
            default ->
                throw new RuntimeException("Unsupported unary operation type " + unaryTerm.getOperation() + ". Actual type " + resultType);
        };
    }

    private CodeBlock toBinaryTermExpression(TypeDefinition typeDefinition, Field field, TypeReference resultType, BinaryTerm binaryTerm, ThrowingFunction<VariableLiteral, CodeBlock, BufferException> variableExpressionGenerator, boolean isParse) throws BufferException {
        Term a = binaryTerm.getA();
        Term b = binaryTerm.getB();
        String operation = binaryTerm.getOperation();
        return switch (operation) {
            case "^" -> {
                if ((resultType != ANY_TYPE_REFERENCE) && !resultType.isIntegerTypeReference() && !resultType.isFloatTypeReference()) {
                    throw new IllegalArgumentException("'A^B' expression requires numeric result type. Actual type " + resultType);
                }
                yield CodeBlock.of("$T.pow($L, $L)", Math.class,
                    toExpression(typeDefinition, field, resultType, a, variableExpressionGenerator, isParse),
                    toExpression(typeDefinition, field, resultType, b, variableExpressionGenerator, isParse));
            }
            case "*", "/", "%", "+", "-" -> {
                if ((resultType != ANY_TYPE_REFERENCE) && !resultType.isIntegerTypeReference() && !resultType.isFloatTypeReference()) {
                    throw new IllegalArgumentException("'A" + operation + "B' expression requires numeric result type. Actual type " + resultType);
                }
                yield CodeBlock.of("($L) $L ($L)",
                    toExpression(typeDefinition, field, resultType, a, variableExpressionGenerator, isParse),
                    operation,
                    toExpression(typeDefinition, field, resultType, b, variableExpressionGenerator, isParse));
            }
            case ">>", "<<" -> CodeBlock.of("($L) $L ($L)",
                toExpression(typeDefinition, field, resultType, a, variableExpressionGenerator, isParse),
                operation,
                toExpression(typeDefinition, field, INT_TYPE_REFERENCE, b, variableExpressionGenerator, isParse));
            case ">=", "<=", ">", "<", "==", "!=" -> {
                if ((resultType != ANY_TYPE_REFERENCE) && !resultType.isBooleanTypeReference()) {
                    throw new IllegalArgumentException("'A" + operation + "B' expression requires boolean result type. Actual type " + resultType);
                }
                // TODO: Try to infer the types of the arguments in this case
                yield CodeBlock.of("($L) $L ($L)",
                    toExpression(typeDefinition, field, ANY_TYPE_REFERENCE, a, variableExpressionGenerator, isParse),
                    operation,
                    toExpression(typeDefinition, field, ANY_TYPE_REFERENCE, b, variableExpressionGenerator, isParse));
                // TODO: Try to infer the types of the arguments in this case
            }
            case "&&", "||" -> {
                if ((resultType != ANY_TYPE_REFERENCE) && !resultType.isBooleanTypeReference()) {
                    throw new IllegalArgumentException("'A" + operation + "B' expression requires boolean result type. Actual type " + resultType);
                }
                yield CodeBlock.of("($L) $L ($L)",
                    toExpression(typeDefinition, field, resultType, a, variableExpressionGenerator, isParse),
                    operation,
                    toExpression(typeDefinition, field, resultType, b, variableExpressionGenerator, isParse));
            }
            case "&", "|" -> {
                if ((resultType != ANY_TYPE_REFERENCE) && !resultType.isIntegerTypeReference() && !resultType.isByteTypeReference()) {
                    throw new IllegalArgumentException("'A" + operation + "B' expression requires byte or integer result type. Actual type " + resultType);
                }
                yield CodeBlock.of("($L) $L ($L)",
                    toExpression(typeDefinition, field, resultType, a, variableExpressionGenerator, isParse),
                    operation,
                    toExpression(typeDefinition, field, resultType, b, variableExpressionGenerator, isParse));
            }
            default -> throw new IllegalArgumentException("Unsupported ternary operation type " + operation);
        };
    }

    private CodeBlock toTernaryTermExpression(TypeDefinition typeDefinition, Field field, TypeReference resultType, TernaryTerm ternaryTerm, ThrowingFunction<VariableLiteral, CodeBlock, BufferException> variableExpressionGenerator, boolean isParse) throws BufferException {
        if ("if".equals(ternaryTerm.getOperation())) {
            Term a = ternaryTerm.getA();
            Term b = ternaryTerm.getB();
            Term c = ternaryTerm.getC();
            return CodeBlock.builder()
                .add("(($L) ? $L : $L)",
                    toExpression(typeDefinition, field, BOOL_TYPE_REFERENCE, a, variableExpressionGenerator, isParse),
                    toExpression(typeDefinition, field, resultType, b, variableExpressionGenerator, isParse),
                    toExpression(typeDefinition, field, resultType, c, variableExpressionGenerator, isParse))
                .build();
        } else {
            throw new IllegalArgumentException("Unsupported ternary operation type " + ternaryTerm.getOperation() + ". Actual type " + resultType);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Parsing

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private CodeBlock toVariableParseExpression(TypeDefinition typeDefinition, Field field, TypeReference resultType, VariableLiteral variableLiteral, List<Argument> parserArguments) throws BufferException {
        // Call a static function in the drivers StaticHelper
        if ("STATIC_CALL".equals(variableLiteral.getName())) {
            return toStaticCallParseExpression(typeDefinition, targetPackage + ".utils", field, resultType, variableLiteral, parserArguments);
        }
        // All uppercase names are not fields, but utility methods.
        else if (variableLiteral.getName().equals(variableLiteral.getName().toUpperCase())) {
            return toStaticCallParseExpression(typeDefinition, "org.apache.plc4x.java.spi.utils", field, resultType, variableLiteral, parserArguments);
        }
        // If we're referencing an implicit field, we need to handle that differently.
        else if (isVariableLiteralImplicitField(typeDefinition, variableLiteral)) {
            return CodeBlock.of(variableLiteral.getName());
        }
        // In contrast to the serialization expression, when parsing, the virtual field should be a local variable.
        /*else if (isVariableLiteralVirtualField(typeDefinition, variableLiteral)) {
            return toGetterVariableExpression(field, resultType, variableLiteral);
        }*/
        // CAST expressions are special as we need to add a ".class" to the second parameter in Java.
        else if ("CAST".equals(variableLiteral.getName())) {
            return toCastVariableParseExpression(typeDefinition, field, resultType, variableLiteral, parserArguments);
        }
        // Special handling for ByteOrder enums (Built in enums)
/*        else if ("BIG_ENDIAN".equals(variableLiteral.getName())) {
            return "ByteOrder.BIG_ENDIAN";
        } else if ("LITTLE_ENDIAN".equals(variableLiteral.getName())) {
            return "ByteOrder.LITTLE_ENDIAN";
        }*/

        // The synthetic checksumRawData is a local field and should not be accessed as bean property.
        boolean isLocalVariable = "readBuffer".equals(variableLiteral.getName());
        boolean isTypeArg = "_type".equals(variableLiteral.getName());
        if (!isLocalVariable && !isTypeArg && parserArguments != null) {
            for (Argument serializerArgument : parserArguments) {
                if (serializerArgument.getName().equals(variableLiteral.getName())) {
                    isLocalVariable = true;
                    break;
                }
            }
        }

        if (isLocalVariable) {
            CodeBlock.Builder builder = CodeBlock.builder().add(variableLiteral.getName());
            if (variableLiteral.getChild().isPresent()) {
                builder.add(".").add(toGetterVariableExpression(field, resultType, variableLiteral.getChild().get()));
            }
            return builder.build();
        }
        // If a "_type" arg is used, then this refers to attributes added to the field (type).
        else if (isTypeArg) {
            String part = variableLiteral.getChild().map(VariableLiteral::getName).orElse("");
            Optional<Term> attribute = field.getAttribute(part);
            return toParseExpression(typeDefinition, field, null, attribute.orElseThrow(() -> new BufferException("Field is mising attribute '" + part + "'")), parserArguments);
            /*return switch (part) {
                case "name" -> CodeBlock.of("\"" + field.getTypeName() + "\"");
                case "length" -> CodeBlock.of("\"" + ((SimpleTypeReference) field).getSizeInBits() + "\"");
                case "encoding" -> {
                    String encoding = ((StringLiteral) field.getEncoding().orElse(new DefaultStringLiteral("UTF8"))).getValue();
                    yield CodeBlock.of("\"" + encoding + "\"");
                }
                default -> CodeBlock.of("");
            };*/
        } else {
            return toPropertyVariableExpression(field, resultType, variableLiteral);
        }
    }

    private CodeBlock toStaticCallParseExpression(TypeDefinition typeDefinition, String packageName, Field field, TypeReference resultType, VariableLiteral variableLiteral, List<Argument> parserArguments) throws BufferException {
        ClassName staticHelper = ClassName.get(packageName, "StaticHelper");

        // Check if we have at least one argument, as this contains the name of the function we want to call.
        List<Term> arguments = variableLiteral.getArgs().orElseThrow(() -> new RuntimeException("A STATIC_CALL expression needs arguments"));
        if (arguments.isEmpty()) {
            throw new RuntimeException("A STATIC_CALL expression expects at least one argument.");
        }

        String methodName;
        if ("STATIC_CALL".equalsIgnoreCase(variableLiteral.getName())) {
            // Get the class and method name
            methodName = arguments.getFirst().asLiteral()
                .orElseThrow(() -> new RuntimeException("First argument should be a literal"))
                .asStringLiteral().orElseThrow(() -> new RuntimeException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral")).
                getValue();
            arguments.removeFirst();
        } else {
            methodName = variableLiteral.getName();
        }

        // Create the call itself.
        CodeBlock.Builder staticCallBuilder = CodeBlock.builder()
            .add("$T.$L(", staticHelper, methodName);

        // TODO: Fix this ugly hack ...
        if ("CAST".equalsIgnoreCase(methodName)) {
            CodeBlock parseExpression = toParseExpression(typeDefinition, field, ANY_TYPE_REFERENCE, arguments.get(0), parserArguments);
            String classNameString = arguments.get(1).asLiteral().orElseThrow().asStringLiteral().orElseThrow().getValue();
            ClassName className = ClassName.get("", classNameString);
            staticCallBuilder.add("$L, $T.class", parseExpression, className);
        } else {
            for (int i = 0; i < arguments.size(); i++) {
                Term arg = arguments.get(i);
                if (i > 0) {
                    staticCallBuilder.add(", ");
                }
                staticCallBuilder.add(toParseExpression(typeDefinition, field, ANY_TYPE_REFERENCE, arg, parserArguments));
           /*if (arg instanceof VariableLiteral) {
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
                            String encoding = ((StringLiteral) field.getEncoding().orElse(new DefaultStringLiteral("UTF8"))).getValue();
                            sb.append("\"").append(encoding).append("\"");
                            break;
                    }
                } else {
                    sb.append(toVariableParseExpression(field, variableLiteralArg, null));
                }
            } else if (arg instanceof StringLiteral) {
                sb.append(((StringLiteral) arg).getValue());
            }*/
            }
        }
        staticCallBuilder.add(")");
        if (variableLiteral.getIndex().isPresent()) {
            // TODO: If this is a byte typed field, this needs to be an array accessor instead.
            staticCallBuilder.add(".get($L)", variableLiteral.getIndex().orElseThrow());
        }
        if (variableLiteral.getChild().isPresent()) {
            staticCallBuilder.add(".$L", toGetterVariableExpression(field, resultType, variableLiteral.getChild().get()));
        }
        return staticCallBuilder.build();
    }

    private CodeBlock toCastVariableParseExpression(TypeDefinition typeDefinition, Field field, TypeReference resultType, VariableLiteral variableLiteral, List<Argument> parserArguments) throws BufferException {
        List<Term> arguments = variableLiteral.getArgs().orElseThrow(() -> new RuntimeException("A Cast expression needs arguments"));
        if (arguments.size() != 2) {
            throw new RuntimeException("A CAST expression expects exactly two arguments.");
        }
        VariableLiteral firstArgument = arguments.get(0).asLiteral()
            .orElseThrow(() -> new RuntimeException("First argument should be a literal"))
            .asVariableLiteral()
            .orElseThrow(() -> new RuntimeException("First argument should be a Variable literal"));
        StringLiteral typeArgument = arguments.get(1).asLiteral().orElseThrow(() -> new RuntimeException("Second argument should be a String literal"))
            .asStringLiteral()
            .orElseThrow(() -> new RuntimeException("Second argument should be a String literal"));
        CodeBlock.Builder codeBlock = CodeBlock.builder()
            .add("CAST($L, $L.class)", toVariableParseExpression(typeDefinition, field, ANY_TYPE_REFERENCE, firstArgument, parserArguments), typeArgument.getValue());
        if (variableLiteral.getChild().isPresent()) {
            codeBlock.add(".$L", toPropertyVariableExpression(field, resultType, variableLiteral.getChild().get()));
        }
        return codeBlock.build();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Serialization
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private CodeBlock toVariableSerializationExpression(TypeDefinition typeDefinition, Field field, TypeReference resultType, VariableLiteral variableLiteral, List<Argument> serializerArguments) throws BufferException {
        if ("STATIC_CALL".equals(variableLiteral.getName())) {
            return toStaticCallSerializationExpression(typeDefinition, targetPackage + ".utils", field, resultType, variableLiteral, serializerArguments);
        }
        // All uppercase names are not fields, but utility methods.
        else if (variableLiteral.getName().equals(variableLiteral.getName().toUpperCase())) {
            return toStaticCallSerializationExpression(typeDefinition, "org.apache.plc4x.java.spi.utils", field, resultType, variableLiteral, serializerArguments);
        }
        // If we are accessing implicit fields, we need to rely on a local variable instead.
        else if (isVariableLiteralImplicitField(typeDefinition, variableLiteral)) {
            final ImplicitField referencedImplicitField = getReferencedImplicitField(typeDefinition, variableLiteral);
            return toSerializationExpression(typeDefinition, referencedImplicitField, referencedImplicitField.getType(),
                getReferencedImplicitField(typeDefinition, variableLiteral).getSerializeExpression(), serializerArguments);
        }
        // If we are accessing a discriminator or virtual field, we need to call a getter method.
        else if (isVariableLiteralDiscriminatorField(typeDefinition, variableLiteral) || isVariableLiteralVirtualField(typeDefinition, variableLiteral)) {
            return toGetterVariableExpression(field, resultType, variableLiteral);
        }
        // CAST expressions are special as we need to add a ".class" to the second parameter in Java.
        /*
        TODO: Possibly implement this.
        if ("CAST".equals(variableLiteral.getName())) {
            return toCastSerializationExpression(field, resultType, variableLiteral, serialzerArguments);
        }*/

        // The synthetic checksumRawData is a local field and should not be accessed as bean property.
        boolean isLocalVariable = "writeBuffer".equals(variableLiteral.getName()) || "checksumRawData".equals(variableLiteral.getName()) || "_value".equals(variableLiteral.getName()) || "element".equals(variableLiteral.getName()) || "size".equals(variableLiteral.getName());
        boolean isTypeArg = "_type".equals(variableLiteral.getName());
        if (!isLocalVariable && !isTypeArg && serializerArguments != null) {
            for (Argument serializerArgument : serializerArguments) {
                if (serializerArgument.getName().equals(variableLiteral.getName())) {
                    isLocalVariable = true;
                    break;
                }
            }
        }

        if (isLocalVariable) {
            CodeBlock.Builder codeBlock = CodeBlock.builder().add(variableLiteral.getName());
            if (variableLiteral.getChild().isPresent()) {
                codeBlock.add(".$L", toGetterVariableExpression(field, resultType, variableLiteral.getChild().get()));
            }
            return codeBlock.build();
        }
/*        else if (isTypeArg) {
            String part = variableLiteral.getChild().map(VariableLiteral::getName).orElse("");
            return switch (part) {
                case "name" -> "\"" + field.getTypeName() + "\"";
                case "length" -> "\"" + ((SimpleTypeReference) field).getSizeInBits() + "\"";
                case "encoding" -> {
                    String encoding = ((StringLiteral) field.getEncoding().orElse(new DefaultStringLiteral("UTF8"))).getValue();
                    yield "\"" + encoding + "\"";
                }
                default -> "";
            };
        }*/
        else {
            return toPropertyVariableExpression(field, resultType, variableLiteral);
        }
    }

    /**
     * Terms containing "STATIC_CALL" or global functions are handled here.
     *
     * @param packageName         the package name for the StaticHelper class to use
     * @param field
     * @param resultType
     * @param variableLiteral
     * @param serializerArguments
     * @return
     */
    private CodeBlock toStaticCallSerializationExpression(TypeDefinition typeDefinition, String packageName, Field field, TypeReference resultType, VariableLiteral variableLiteral, List<Argument> serializerArguments) throws BufferException {
        ClassName staticHelper = ClassName.get(packageName, "StaticHelper");

        // Check if we have at least one argument, as this contains the name of the function we want to call.
        List<Term> arguments = new ArrayList<>(variableLiteral.getArgs().orElseThrow(() -> new RuntimeException("A STATIC_CALL expression needs arguments")));
        if (arguments.isEmpty()) {
            throw new RuntimeException("A STATIC_CALL expression expects at least one argument.");
        }

        String methodName;
        // This is a STATIC_CALL to a StaticHelper method inside the current protocol implementation.
        if ("STATIC_CALL".equalsIgnoreCase(variableLiteral.getName())) {
            // Get the class and method name
            methodName = arguments.getFirst().asLiteral()
                .orElseThrow(() -> new RuntimeException("First argument should be a literal"))
                .asStringLiteral().orElseThrow(() -> new RuntimeException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral")).
                getValue();
            arguments.removeFirst();
        }
        // This is a static call to a method in the SPI.
        else {
            methodName = variableLiteral.getName();
        }

        // Create the call itself.
        CodeBlock.Builder staticCallBuilder = CodeBlock.builder()
            .add("$T.$L(", staticHelper, methodName);

        // Add all arguments.
        for (int i = 0; i < arguments.size(); i++) {
            Term arg = arguments.get(i);
            if (i > 0) {
                staticCallBuilder.add(", ");
            }
            staticCallBuilder.add(toSerializationExpression(typeDefinition, field, ANY_TYPE_REFERENCE, arg, serializerArguments));
            /*if (arg instanceof VariableLiteral) {
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
                            String encoding = ((StringLiteral) field.getEncoding().orElse(new DefaultStringLiteral("UTF8"))).getValue();
                            sb.append("\"").append(encoding).append("\"");
                            break;
                    }
                } else {
                    sb.append(toVariableSerializationExpression(field, va, serialzerArguments));
                }
            } else if (arg instanceof StringLiteral) {
                sb.append(((StringLiteral) arg).getValue());
            }*/
        }
        staticCallBuilder.add(")");
        if (variableLiteral.getIndex().isPresent()) {
            // TODO: If this is a byte typed field, this needs to be an array accessor instead.
            staticCallBuilder.add(".get($L)", variableLiteral.getIndex().orElseThrow());
        }
        return staticCallBuilder.build();
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // General purpose
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    private CodeBlock toPropertyVariableExpression(Field field, TypeReference resultType, VariableLiteral variableLiteral) {
        String variableLiteralName = variableLiteral.getName();
        if ("lengthInBytes".equalsIgnoreCase(variableLiteralName)) {
            variableLiteralName = "getLengthInBytes()";
        } else if ("lengthInBits".equalsIgnoreCase(variableLiteralName)) {
            variableLiteralName = "getLengthInBits()";
        }

        CodeBlock.Builder codeBlock = CodeBlock.builder()
            .add("$L", variableLiteralName);
        if (variableLiteral.getIndex().isPresent()) {
            codeBlock.add(".get($L)", variableLiteral.getIndex().orElseThrow());
        }
        if (variableLiteral.getChild().isPresent()) {
            codeBlock.add("." + toGetterVariableExpression(field, resultType, variableLiteral.getChild().get()));
        }
        return codeBlock.build();
    }

    private CodeBlock toGetterVariableExpression(Field field, TypeReference resultType, VariableLiteral variableLiteral) {
        String variableLiteralName = variableLiteral.getName();
        CodeBlock.Builder codeBlock;
        // TODO: Hack ... generally we should get the type of the variable referenced by the variable literal and use that.
        //  If it's a string, use the first branch.
        if ("length".equalsIgnoreCase(variableLiteralName)) {
            codeBlock = CodeBlock.builder()
                .add("$L()", variableLiteralName);
        } else {
            codeBlock = CodeBlock.builder()
                .add("get$L()", StaticHelper.CAPITALIZE(variableLiteralName));
        }

        if (variableLiteral.getIndex().isPresent()) {
            codeBlock.add(".get($L)", variableLiteral.getIndex().orElseThrow());
        }
        if (variableLiteral.getChild().isPresent()) {
            codeBlock.add("." + toGetterVariableExpression(field, resultType, variableLiteral.getChild().get()));
        }
        return codeBlock.build();
    }

    /**
     * Enums are always based on a main type. This helper accesses this information in a safe manner.
     *
     * @param typeReference type reference
     * @return simple type reference for the enum type referenced by the given type reference
     */
    public SimpleTypeReference getEnumBaseTypeReference(TypeReference typeReference) {
        // Enum types always have simple type references.
        return getEnumTypeDefinition(typeReference).getType().orElseThrow();
    }

    protected EnumTypeDefinition getEnumTypeDefinition(TypeReference typeReference) {
        NonSimpleTypeReference nonSimpleTypeReference = typeReference.asNonSimpleTypeReference().orElseThrow(
            () -> new RuntimeException("type reference for enum types must be of type non simple type"));
        String typeName = nonSimpleTypeReference.getName();
        final TypeDefinition typeDefinition = nonSimpleTypeReference.getTypeDefinition();
        if (typeDefinition == null) {
            throw new RuntimeException("Couldn't find given enum type definition with name " + typeName);
        }
        // TODO: same here. It is named complex type reference but it references a enum...
        if (!typeDefinition.isEnumTypeDefinition()) {
            throw new RuntimeException("Referenced type with name " + typeName + " is not an enum type");
        }
        return (EnumTypeDefinition) typeDefinition;
    }

    protected TypeName getFieldTypeClassName(Field field) {
        if (field.isTypedField()) {
            TypedField typedField = field.asTypedField().orElseThrow();
            TypeReference fieldType = typedField.getType();
            return getLanguageTypeNameForTypeReference(fieldType, !field.isOptionalField());
        } else {
            throw new IllegalArgumentException("Field is not a typed field.");
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Field Generators
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    protected void generateArrayField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getSerializationValueBlock, CodeBlock getSizeBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        ArrayField arrayField = field.asArrayField().orElseThrow();
        String fieldName = arrayField.getName();
        CodeBlock attributesCodeBlock = getAttributes(typeDefinition, field, fieldName, parserArguments);

        TypeReference fieldType = arrayField.getType();
        TypeName fieldTypeClassName = getFieldTypeClassName(field);

        TypeReference elementTypeReference = arrayField.getType().getElementTypeReference();

        ClassName fieldReaderFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.reader", "FieldReaderFactory");
        ClassName fieldWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.writer", "FieldWriterFactory");
        ClassName dataWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.data.writer", "DataWriterFactory");
        CodeBlock readBlock = getDataReaderCall(fieldType);
        CodeBlock writeBlock = getDataWriterCall(fieldType);

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        parseBlockBuilder.add("// Array Field: $L\n", fieldName);
        if (elementTypeReference.isByteBased()) {
            if (!field.isCountArrayField() && !field.isLengthArrayField()) {
                throw new IllegalArgumentException("Array fields of type byte only support 'count' and 'length' loop-types.");
            }
            // In the case of a byte-array field, it doesn't matter if it's count or length.
            CodeBlock loopExpression = toParseExpression(typeDefinition, field, INT_TYPE_REFERENCE, arrayField.getLoopExpression(), parserArguments);
            if(field.isCountArrayField()) {
                parseBlockBuilder.addStatement("byte[] $L = readBuffer.readBits($T.toIntExact(($L) * 8), $L)", arrayField.getName(), Math.class, loopExpression, attributesCodeBlock);
            } else {
                parseBlockBuilder.addStatement("byte[] $L = readBuffer.readBits($T.toIntExact($L), $L)", arrayField.getName(), Math.class, loopExpression, attributesCodeBlock);
            }
        } else if (field.isCountArrayField()) {
            CodeBlock loopExpression = toParseExpression(typeDefinition, field, INT_TYPE_REFERENCE, arrayField.getLoopExpression(), parserArguments);
            parseBlockBuilder.addStatement("$T $L = $T.readCountArrayField($L, $L, $L)", fieldTypeClassName, arrayField.getName(), fieldReaderFactory, readBlock, loopExpression, attributesCodeBlock);
        } else if (field.isLengthArrayField()) {
            CodeBlock loopExpression = toParseExpression(typeDefinition, field, INT_TYPE_REFERENCE, arrayField.getLoopExpression(), parserArguments);
            parseBlockBuilder.addStatement("$T $L = $T.readLengthArrayField($L, $L, $L)", fieldTypeClassName, arrayField.getName(), fieldReaderFactory, readBlock, loopExpression, attributesCodeBlock);
        } else if (field.isTerminatedArrayField()) {
            CodeBlock loopExpression = toParseExpression(typeDefinition, field, INT_TYPE_REFERENCE, arrayField.getLoopExpression(), parserArguments);
            parseBlockBuilder.addStatement("$T $L = $T.readTerminatedArrayField($L, () -> (boolean) $L, $L)", fieldTypeClassName, arrayField.getName(), fieldReaderFactory, readBlock, loopExpression, attributesCodeBlock);
        }
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        serializeCodeBlockBuilder.add("// Array Field: $L\n", fieldName);
        // Shortcut for handling raw byte arrays.
        if (elementTypeReference.isByteBased()) {
            CodeBlock loopExpression = CodeBlock.of("$L.length", getSerializationValueBlock);
            CodeBlock rawByteArrayWriteBlock = CodeBlock.of("$T.writeByteArray(writeBuffer, (int) (($L != null) ? $L : 0))", dataWriterFactory, getSerializationValueBlock, loopExpression);
            serializeCodeBlockBuilder.addStatement("$T.writeByteArrayField($L, $L, $L)", fieldWriterFactory, getSerializationValueBlock, rawByteArrayWriteBlock, attributesCodeBlock);
        }
        // Arrays of simple types need explicit writers.
        else if (elementTypeReference.isSimpleTypeReference()) {
            serializeCodeBlockBuilder.addStatement("$T.writeSimpleTypeArrayField($L, $L, $L)", fieldWriterFactory, getSerializationValueBlock, writeBlock, attributesCodeBlock);
        }
        // Arrays of enum types are treated differently.
        else if (elementTypeReference.isEnumTypeReference()) {
            TypeName enumType = getLanguageTypeNameForTypeReference(elementTypeReference.asEnumTypeReference().orElseThrow(), false);
            serializeCodeBlockBuilder.addStatement("$T.writeSimpleTypeArrayField($L, $T.writeEnum($T::getValue, $T::name, $L), $L)", fieldWriterFactory, getSerializationValueBlock, dataWriterFactory, enumType, enumType, writeBlock, attributesCodeBlock);
        }
        // Arrays of complex types use the Message.serialize method for serialization.
        else {
            serializeCodeBlockBuilder.addStatement("$T.writeComplexTypeArrayField($L, writeBuffer, $L)", fieldWriterFactory, getSerializationValueBlock, attributesCodeBlock);
        }
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Array Field: $L\n", fieldName);
        if (elementTypeReference.isSimpleTypeReference()) {
            SimpleTypeReference simpleTypeReference = elementTypeReference.asSimpleTypeReference().orElseThrow();
            if (simpleTypeReference.isByteBased()) {
                getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += 8 * (($L != null) ? $L.length : 0)", getSizeBlock, getSizeBlock);
            } else if (simpleTypeReference.isVstringTypeReference()) {
                // Every element of the array is a vstring of the same declared length, so the total
                // is that length times the number of elements. The length term is rendered the same
                // way as for a single vstring field; it may reference parser arguments (e.g. a
                // "stringLength" parameter), which is why it cannot be reduced to a constant here.
                VstringTypeReference vstringTypeReference = simpleTypeReference.asVstringTypeReference().orElseThrow();
                CodeBlock elementLength = toSerializationExpression(typeDefinition, field, INT_TYPE_REFERENCE,
                    vstringTypeReference.getLengthExpression(), parserArguments);
                getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += ($L) * $L.size()", elementLength, getSizeBlock);
            } else {
                // TODO: Generate dynamic type length values here.
                Optional<Term> unsignedIntegerEncodingAttribute = field.getAttribute("unsignedIntegerEncoding");
                if(unsignedIntegerEncodingAttribute.isPresent() && unsignedIntegerEncodingAttribute.get().stringRepresentation().startsWith("\"VAR-")) {
                    TypeName languageTypeNameForTypeReference = getLanguageTypeNameForTypeReference(elementTypeReference, false);
                    getLengthInBitsCodeBlockBuilder.beginControlFlow("for ($T _element : $L)", languageTypeNameForTypeReference, fieldName);
                    getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $T.GET_VAR_LENGTH_UINT_IN_BITS(_element)", ClassName.get("org.apache.plc4x.java.spi.utils", "StaticHelper"));
                    getLengthInBitsCodeBlockBuilder.endControlFlow();
                } else {
                    Optional<Term> signedIntegerEncodingAttribute = field.getAttribute("signedIntegerEncoding");
                    if(signedIntegerEncodingAttribute.isPresent() && signedIntegerEncodingAttribute.get().stringRepresentation().startsWith("\"VAR-")) {
                        TypeName languageTypeNameForTypeReference = getLanguageTypeNameForTypeReference(elementTypeReference, false);
                        getLengthInBitsCodeBlockBuilder.beginControlFlow("for ($T _element : $L)", languageTypeNameForTypeReference, fieldName);
                        getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $T.GET_VAR_LENGTH_SINT_IN_BITS(_element)", ClassName.get("org.apache.plc4x.java.spi.utils", "StaticHelper"));
                        getLengthInBitsCodeBlockBuilder.endControlFlow();
                    } else {
                        getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L * $L.size()", simpleTypeReference.getSizeInBits(), getSizeBlock);
                    }
                }
            }
        } else if (elementTypeReference.isEnumTypeReference()) {
            // TODO: Generate dynamic type length values here.
            getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L * $L.size()", elementTypeReference.asEnumTypeReference().orElseThrow().getBaseTypeReference().orElseThrow().getSizeInBits(), getSizeBlock);
        } else if (fieldType.isDataIoTypeReference()) {
            throw new RuntimeException("Array fields of type dataIo are not supported.");
        }
        // In all other cases, we're dealing with complex types.
        else {
            if (arrayField.isCountArrayField()) {
                NonSimpleTypeReference nonSimpleTypeReference = elementTypeReference.asNonSimpleTypeReference().orElseThrow();
                ClassName threadLocalHelper = ClassName.get("org.apache.plc4x.java.spi.fields.utils", "ThreadLocalHelper");
                ClassName nonSimpleTypeClassName = ClassName.get(targetPackage, nonSimpleTypeReference.getName());
                getLengthInBitsCodeBlockBuilder.beginControlFlow("if ($L != null)", fieldName);
                getLengthInBitsCodeBlockBuilder.addStatement("int i = 0");
                getLengthInBitsCodeBlockBuilder.beginControlFlow("for ($T _element : $L)", nonSimpleTypeClassName, fieldName);
                getLengthInBitsCodeBlockBuilder.addStatement("$T.lastItemThreadLocal.set(++i >= $L.size())", threadLocalHelper, fieldName);
                getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += _element.getLengthInBits()");
                getLengthInBitsCodeBlockBuilder.endControlFlow();
                getLengthInBitsCodeBlockBuilder.endControlFlow();
            } else {
                getLengthInBitsCodeBlockBuilder.beginControlFlow("for ($T _element : $L)", ClassName.get("org.apache.plc4x.java.spi.buffers.api", "Message"), fieldName);
                getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += _element.getLengthInBits()");
                getLengthInBitsCodeBlockBuilder.endControlFlow();
            }
        }
        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generateAssertField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        AssertField assertField = field.asAssertField().orElseThrow();
        String fieldName = assertField.getName();
        CodeBlock attributesCodeBlock = getAttributes(typeDefinition, field, fieldName, parserArguments);

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        parseBlockBuilder.add("// Assert Field: $L\n", fieldName);
        // TODO: Implement this.
        parseBlockBuilder.addStatement("fail");
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        serializeCodeBlockBuilder.add("// Assert Field: $L\n", fieldName);
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Assert Field: $L\n", fieldName);
        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generateChecksumField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        ChecksumField checksumField = field.asChecksumField().orElseThrow();
        String fieldName = checksumField.getName();
        CodeBlock attributesCodeBlock = getAttributes(typeDefinition, field, fieldName, parserArguments);

        TypeReference fieldType = checksumField.getType();
        TypeName fieldTypeClassName = getLanguageTypeNameForTypeReference(fieldType, true);

        ClassName fieldReaderFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.reader", "FieldReaderFactory");
        ClassName fieldWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.writer", "FieldWriterFactory");
        CodeBlock readBlock = getDataReaderCall(fieldType);
        CodeBlock writeBlock = getDataWriterCall(fieldType);
        CodeBlock referenceBlock = toSerializationExpression(typeDefinition, checksumField, fieldType, checksumField.getChecksumExpression(), parserArguments);

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        parseBlockBuilder.add("// Checksum Field: $L\n", fieldName);
        parseBlockBuilder.addStatement("$T $L = $T.readChecksumField($L, $L, $L)", fieldTypeClassName, fieldName, fieldReaderFactory, readBlock, referenceBlock, attributesCodeBlock);
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        serializeCodeBlockBuilder.add("// Checksum Field: $L\n", fieldName);
        serializeCodeBlockBuilder.addStatement("$T.writeChecksumField(($T) $L, $L, $L)", fieldWriterFactory, fieldTypeClassName, referenceBlock, writeBlock, attributesCodeBlock);
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Checksum Field: $L\n", fieldName);
        SimpleTypeReference simpleTypeReference = fieldType.asSimpleTypeReference().orElseThrow();
        if (simpleTypeReference.isVstringTypeReference()) {
            VstringTypeReference vstringTypeReference = simpleTypeReference.asVstringTypeReference().orElseThrow();
            getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", toSerializationExpression(typeDefinition, checksumField, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), null));
        } else {
            // TODO: Generate dynamic type length values here.
            getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", simpleTypeReference.getSizeInBits());
        }
        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generateConstField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        ConstField constField = field.asConstField().orElseThrow();
        String fieldName = constField.getName();
        CodeBlock attributesCodeBlock = getAttributes(typeDefinition, field, fieldName, parserArguments);

        TypeReference fieldType = constField.getType();
        TypeName fieldTypeClassName = getLanguageTypeNameForTypeReference(constField.getType(), true);

        ClassName fieldReaderFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.reader", "FieldReaderFactory");
        ClassName fieldWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.writer", "FieldWriterFactory");
        CodeBlock readBlock = getDataReaderCall(fieldType);
        CodeBlock writeBlock = getDataWriterCall(fieldType);

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        String reservedValue = constField.getName().toUpperCase();
        parseBlockBuilder.add("// Const Field: $L\n", fieldName);
        parseBlockBuilder.addStatement("$T $L = $T.readConstField($L, $L, $L)", fieldTypeClassName, constField.getName(), fieldReaderFactory, readBlock, reservedValue, attributesCodeBlock);
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        serializeCodeBlockBuilder.add("// Const Field: $L\n", fieldName);
        serializeCodeBlockBuilder.addStatement("$T.writeConstField(($T) $L, $L, $L)", fieldWriterFactory, fieldTypeClassName, reservedValue, writeBlock, attributesCodeBlock);
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Const Field: $L\n", fieldName);
        SimpleTypeReference simpleTypeReference = fieldType.asSimpleTypeReference().orElseThrow();
        if (simpleTypeReference.isVstringTypeReference()) {
            VstringTypeReference vstringTypeReference = simpleTypeReference.asVstringTypeReference().orElseThrow();
            getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", toSerializationExpression(typeDefinition, constField, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), null));
        } else {
            // TODO: Generate dynamic type length values here.
            getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", simpleTypeReference.getSizeInBits());
        }
        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generateDiscriminatorField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        DiscriminatorField discriminatorField = field.asDiscriminatorField().orElseThrow();
        String fieldName = discriminatorField.getName();
        CodeBlock attributesCodeBlock = getAttributes(typeDefinition, field, fieldName, parserArguments);

        TypeReference fieldType = discriminatorField.getType();
        TypeName fieldTypeClassName = getLanguageTypeNameForTypeReference(discriminatorField.getType(), true);

        ClassName fieldReaderFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.reader", "FieldReaderFactory");
        ClassName fieldWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.writer", "FieldWriterFactory");
        ClassName dataWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.data.writer", "DataWriterFactory");
        CodeBlock readBlock = getDataReaderCall(fieldType);
        CodeBlock writeBlock = getDataWriterCall(fieldType);

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        if (fieldType.isEnumTypeReference()) {
            parseBlockBuilder.add("// Enum Discriminator Field: $L\n", fieldName);
            parseBlockBuilder.addStatement("$T $L = $T.readDiscriminatorEnumField($L, $L)", fieldTypeClassName, fieldName, fieldReaderFactory, readBlock, attributesCodeBlock);
        } else {
            parseBlockBuilder.add("// Discriminator Field: $L\n", fieldName);
            parseBlockBuilder.addStatement("$T $L = $T.readDiscriminatorField($L, $L)", fieldTypeClassName, fieldName, fieldReaderFactory, readBlock, attributesCodeBlock);
        }
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        CodeBlock valueAccessBlock = CodeBlock.of("get$L()", StaticHelper.CAPITALIZE(fieldName));
        serializeCodeBlockBuilder.addStatement("$T $L = ($T) $L", fieldTypeClassName, fieldName, fieldTypeClassName, valueAccessBlock);
        if (fieldType.isEnumTypeReference()) {
            serializeCodeBlockBuilder.add("// Enum Discriminator Field: $L\n", fieldName);
            writeBlock = CodeBlock.of("$T.writeEnum($T::getValue, $T::name, $L)", dataWriterFactory, fieldTypeClassName, fieldTypeClassName, writeBlock);
            serializeCodeBlockBuilder.addStatement("$T.writeDiscriminatorEnumField(($T) $L, $L, $L)", fieldWriterFactory, fieldTypeClassName, fieldName, writeBlock, attributesCodeBlock);
        } else {
            serializeCodeBlockBuilder.add("// Discriminator Field: $L\n", fieldName);
            serializeCodeBlockBuilder.addStatement("$T.writeDiscriminatorField(($T) $L, $L, $L)", fieldWriterFactory, fieldTypeClassName, fieldName, writeBlock, attributesCodeBlock);
        }
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Discriminator Field: $L\n", fieldName);
        if (fieldType.isSimpleTypeReference()) {
            SimpleTypeReference simpleTypeReference = fieldType.asSimpleTypeReference().orElseThrow();
            if (simpleTypeReference.isVstringTypeReference()) {
                VstringTypeReference vstringTypeReference = simpleTypeReference.asVstringTypeReference().orElseThrow();
                getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", toSerializationExpression(typeDefinition, discriminatorField, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), null));
            } else {
                // TODO: Generate dynamic type length values here.
                getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", simpleTypeReference.getSizeInBits());
            }
        } else if (fieldType.isEnumTypeReference()) {
            EnumTypeReference enumTypeReference = fieldType.asEnumTypeReference().orElseThrow();
            // TODO: Generate dynamic type length values here.
            getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", enumTypeReference.getBaseTypeReference().orElseThrow().getSizeInBits());
        } else if (fieldType.isDataIoTypeReference()) {
            throw new RuntimeException("Discriminator fields of type dataIo are not supported.");
        } else {
            throw new RuntimeException("Discriminator fields of type complex are not supported.");
        }
        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generateEnumField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        EnumField enumField = field.asEnumField().orElseThrow();
        String fieldName = enumField.getName();
        CodeBlock attributesCodeBlock = getAttributes(typeDefinition, field, fieldName, parserArguments);

        TypeReference fieldType = enumField.getType();
        EnumTypeReference enumTypeReference = enumField.getType().asEnumTypeReference().orElseThrow();
        TypeName fieldTypeClassName = getLanguageTypeNameForTypeReference(enumTypeReference, true);

        ClassName fieldReaderFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.reader", "FieldReaderFactory");
        ClassName fieldWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.writer", "FieldWriterFactory");
        ClassName dataWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.data.writer", "DataWriterFactory");
        CodeBlock readBlock = getDataReaderCall(enumTypeReference, "firstEnumForField" + StaticHelper.CAPITALIZE(enumField.getFieldName()));
        CodeBlock writeBlock = getDataWriterCall(fieldType);

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        parseBlockBuilder.add("// Enum Field (enum): $L\n", fieldName);
        parseBlockBuilder.addStatement("$T $L = $T.readEnumField($L, $L)", fieldTypeClassName, fieldName, fieldReaderFactory, readBlock, attributesCodeBlock);
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        serializeCodeBlockBuilder.add("// Enum Field (enum): $L\n", fieldName);
        serializeCodeBlockBuilder.addStatement("$T.writeEnumField(($T) $L, $T.writeEnum($T::get$L, $T::name, $L), $L)", fieldWriterFactory, fieldTypeClassName, getValueBlock, dataWriterFactory, fieldTypeClassName, StaticHelper.CAPITALIZE(enumField.getFieldName()), fieldTypeClassName, writeBlock, attributesCodeBlock);
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Enum Field: $L\n", fieldName);
        // TODO: Generate dynamic type length values here.
        getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", getEnumBaseTypeReference(enumTypeReference).getSizeInBits());
        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generateImplicitField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        ImplicitField implicitField = field.asImplicitField().orElseThrow();
        String fieldName = implicitField.getName();
        CodeBlock attributesCodeBlock = getAttributes(typeDefinition, field, fieldName, parserArguments);

        TypeReference fieldType = implicitField.getType();
        TypeName fieldTypeClassName = getFieldTypeClassName(field);

        ClassName fieldReaderFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.reader", "FieldReaderFactory");
        ClassName fieldWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.writer", "FieldWriterFactory");
        CodeBlock readBlock = getDataReaderCall(implicitField.getType());
        CodeBlock writeBlock = getDataWriterCall(fieldType);

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        parseBlockBuilder.add("// Implicit Field: $L\n", fieldName);
        parseBlockBuilder.addStatement("$T $L = $T.readImplicitField($L, $L)", fieldTypeClassName, fieldName, fieldReaderFactory, readBlock, attributesCodeBlock);
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        serializeCodeBlockBuilder.add("// Implicit Field: $L\n", fieldName);
        serializeCodeBlockBuilder.addStatement("$T $L = ($T) ($L)", fieldTypeClassName, fieldName, fieldTypeClassName, toSerializationExpression(typeDefinition, implicitField, implicitField.getType(), implicitField.getSerializeExpression(), parserArguments));
        serializeCodeBlockBuilder.addStatement("$T.writeImplicitField(($T) $L, $L, $L)", fieldWriterFactory, fieldTypeClassName, getValueBlock, writeBlock, attributesCodeBlock);
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Implicit Field: $L\n", fieldName);
        SimpleTypeReference simpleTypeReference = implicitField.getType().asSimpleTypeReference().orElseThrow();
        if (simpleTypeReference.isVstringTypeReference()) {
            VstringTypeReference vstringTypeReference = simpleTypeReference.asVstringTypeReference().orElseThrow();
            getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", toSerializationExpression(typeDefinition, implicitField, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), null));
        } else {
            // TODO: Generate dynamic type length values here.
            getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", simpleTypeReference.getSizeInBits());
        }
        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generateManualArrayField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        ManualArrayField manualArrayField = field.asManualArrayField().orElseThrow();
        String fieldName = manualArrayField.getName();
        CodeBlock attributesCodeBlock = getAttributes(typeDefinition, field, fieldName, parserArguments);

        TypeName fieldTypeClassName = getFieldTypeClassName(field);

        ClassName fieldReaderFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.reader", "FieldReaderFactory");
        ClassName fieldWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.writer", "FieldWriterFactory");

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        parseBlockBuilder.add("// Manual Array Field: $L\n", fieldName);
        CodeBlock loopExpression = toParseExpression(typeDefinition, manualArrayField, BOOL_TYPE_REFERENCE, manualArrayField.getLoopExpression(), parserArguments);
        CodeBlock parseExpression = toParseExpression(typeDefinition, manualArrayField, manualArrayField.getType(), manualArrayField.getParseExpression(), parserArguments);
        if (manualArrayField.getType().asArrayTypeReference().orElseThrow().getElementTypeReference().isByteBased()) {
            parseBlockBuilder.addStatement("byte[] $L = $T.readManualByteArrayField(readBuffer, ($T<Byte> _values) -> (boolean) ($L), () -> (Byte) $L, $L)", fieldName, fieldReaderFactory, List.class, loopExpression, parseExpression, attributesCodeBlock);
        } else {
            TypeName elementTypeClassName = getLanguageTypeNameForTypeReference(manualArrayField.getType().getElementTypeReference(), false);
            parseBlockBuilder.addStatement("$T $L = $T.readManualArrayField(readBuffer, ($T _values) -> (boolean) ($L), () -> ($T) $L, $L)", fieldTypeClassName, fieldName, fieldReaderFactory, fieldTypeClassName, loopExpression, elementTypeClassName, parseExpression, attributesCodeBlock);
        }
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        serializeCodeBlockBuilder.add("// Manual Array Field: $L\n", fieldName);
        CodeBlock serializationExpression = toSerializationExpression(typeDefinition, manualArrayField, manualArrayField.getType(), manualArrayField.getSerializeExpression(), parserArguments);
        if (manualArrayField.getType().asArrayTypeReference().orElseThrow().getElementTypeReference().isByteBased()) {
            serializeCodeBlockBuilder.addStatement("$T.writeManualArrayField($L, (Byte _value) -> $L, writeBuffer, $L)", fieldWriterFactory, fieldName, serializationExpression, attributesCodeBlock);
        } else {
            TypeName elementTypeClassName = getLanguageTypeNameForTypeReference(manualArrayField.getType().getElementTypeReference(), false);
            serializeCodeBlockBuilder.addStatement("$T.writeManualArrayField($L, ($T _value) -> $L, writeBuffer, $L)", fieldWriterFactory, fieldName, elementTypeClassName, serializationExpression, attributesCodeBlock);
        }
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Manual Array Field: $L\n", fieldName);
        CodeBlock lengthExpression = toParseExpression(typeDefinition, manualArrayField, manualArrayField.getType(), manualArrayField.getLengthExpression(), parserArguments);
        getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", lengthExpression);

        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generateManualField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        ManualField manualField = field.asManualField().orElseThrow();
        String fieldName = manualField.getName();
        CodeBlock attributesCodeBlock = getAttributes(typeDefinition, field, fieldName, parserArguments);

        TypeName fieldTypeClassName = getFieldTypeClassName(manualField);

        ClassName fieldReaderFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.reader", "FieldReaderFactory");
        ClassName fieldWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.writer", "FieldWriterFactory");

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        parseBlockBuilder.add("// Manual Field: $L\n", fieldName);
        CodeBlock parseExpression = toParseExpression(typeDefinition, manualField, manualField.getType(), manualField.getParseExpression(), parserArguments);
        parseBlockBuilder.addStatement("$T $L = $T.readManualField(readBuffer, () -> ($T) ($L), $L)", fieldTypeClassName, manualField.getName(), fieldReaderFactory, fieldTypeClassName, parseExpression, attributesCodeBlock);
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        serializeCodeBlockBuilder.add("// Manual Field: $L\n", fieldName);
        CodeBlock serializeExpression = toParseExpression(typeDefinition, manualField, manualField.getType(), manualField.getSerializeExpression(), parserArguments);
        serializeCodeBlockBuilder.addStatement("$T.writeManualField(() -> $L, writeBuffer, $L)", fieldWriterFactory, serializeExpression, attributesCodeBlock);
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Manual Field: $L\n", fieldName);
        CodeBlock lengthExpression = toParseExpression(typeDefinition, manualField, INT_TYPE_REFERENCE, manualField.getLengthExpression(), parserArguments);
        getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", lengthExpression);
        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generateOptionalField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        OptionalField optionalField = field.asOptionalField().orElseThrow();
        String fieldName = optionalField.getName();
        CodeBlock attributesCodeBlock = getAttributes(typeDefinition, field, fieldName, parserArguments);

        Optional<Term> nullBytesHexAttribute = optionalField.getAttribute("nullBytesHex");
        if (nullBytesHexAttribute.isPresent()) {
            String nullBytesHex = nullBytesHexAttribute.get().stringRepresentation();
            attributesCodeBlock = CodeBlock.of("$T.$L($S), $L", ClassName.get("org.apache.plc4x.java.spi.fields.fields", "WithFieldOption"), "WithNullBytesHex", nullBytesHex, attributesCodeBlock);
        }
        TypeReference fieldType = optionalField.getType();
        TypeName fieldTypeClassName = getFieldTypeClassName(field);

        ClassName fieldReaderFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.reader", "FieldReaderFactory");
        ClassName fieldWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.writer", "FieldWriterFactory");
        CodeBlock readBlock = getDataReaderCall(fieldType);
        CodeBlock writeBlock = getDataWriterCall(fieldType);

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        if (optionalField.getType().isEnumTypeReference() && optionalField.getConditionExpression().isEmpty()) {
            parseBlockBuilder.add("// Optional Field (enum): $L\n", fieldName);
            parseBlockBuilder.addStatement("$T $L = $T.readEnumField($L, $L)", fieldTypeClassName, optionalField.getName(), fieldReaderFactory, readBlock, attributesCodeBlock);
        } else if (optionalField.getConditionExpression().isPresent()) {
            parseBlockBuilder.add("// Optional Field (conditional): $L\n", fieldName);
            CodeBlock conditionExpression = toParseExpression(typeDefinition, optionalField, BOOL_TYPE_REFERENCE, optionalField.getConditionExpression().get(), parserArguments);
            if (fieldType.isDataIoTypeReference()) {
                parseBlockBuilder.addStatement("$T $L = $T.readOptionalField($L, $L, $L)", fieldTypeClassName, optionalField.getName(), fieldReaderFactory, readBlock, conditionExpression, attributesCodeBlock);
            } else {
                parseBlockBuilder.addStatement("$T $L = $T.readOptionalField($L, $L, $L)", fieldTypeClassName, optionalField.getName(), fieldReaderFactory, readBlock, conditionExpression, attributesCodeBlock);
            }
        } else {
            parseBlockBuilder.add("// Optional Field: $L\n", fieldName);
            parseBlockBuilder.addStatement("$T $L = $T.readOptionalField($L, $L)", fieldTypeClassName, optionalField.getName(), fieldReaderFactory, readBlock, attributesCodeBlock);
        }
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        if (nullBytesHexAttribute.isEmpty()) {
            serializeCodeBlockBuilder.beginControlFlow("if($L != null)", fieldName);
        }
        if (optionalField.getType().isEnumTypeReference()) {
            EnumTypeReference enumTypeReference = optionalField.getType().asEnumTypeReference().orElseThrow();
            serializeCodeBlockBuilder.add("// Optional Field (enum): $L\n", fieldName);
            TypeName typeClassName = getLanguageTypeNameForTypeReference(enumTypeReference.getBaseTypeReference().orElseThrow(), true);
            serializeCodeBlockBuilder.addStatement("$T.writeOptionalField(($T) $L.getValue(), $L, $L)", fieldWriterFactory, typeClassName, getValueBlock, writeBlock, attributesCodeBlock);
        } else if (fieldType.isDataIoTypeReference()) {
            DataIoTypeReference dataIoTypeReference = fieldType.asDataIoTypeReference().orElseThrow();
            serializeCodeBlockBuilder.add("// Optional Field: $L\n", fieldName);
            // Build parameter string from DataIoTypeReference params
            StringBuilder paramsBuilder = new StringBuilder();
            paramsBuilder.append("writeBuffer, ").append(optionalField.getName());
            List<Term> paramTerms = dataIoTypeReference.getParams().orElse(Collections.emptyList());
            for (int i = 0; i < paramTerms.size(); i++) {
                Term paramTerm = paramTerms.get(i);
                final TypeReference argumentType = getArgumentType(dataIoTypeReference, i);
                paramsBuilder.append(", (").append(getLanguageTypeNameForTypeReference(argumentType, true)).append(") (")
                    .append(toSerializationExpression(typeDefinition, optionalField, argumentType, paramTerm, parserArguments)).append(")");
            }
            serializeCodeBlockBuilder.addStatement("$T.staticSerialize($L)", ClassName.get(targetPackage, dataIoTypeReference.getName()), paramsBuilder.toString());
        } else {
            serializeCodeBlockBuilder.add("// Optional Field: $L\n", fieldName);
            serializeCodeBlockBuilder.addStatement("$T.writeOptionalField(($T) $L, $L, $L)", fieldWriterFactory, fieldTypeClassName, getValueBlock, writeBlock, attributesCodeBlock);
        }
        if (nullBytesHexAttribute.isEmpty()) {
            serializeCodeBlockBuilder.endControlFlow();
        }
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Optional Field: $L\n", fieldName);
        // Optional fields with "nullBytesHex"-attribute always have a size.
        // DataIoTypeReference fields use try-catch instead of if-block.
        if (nullBytesHexAttribute.isEmpty() && !fieldType.isDataIoTypeReference()) {
            getLengthInBitsCodeBlockBuilder.beginControlFlow("if($L != null)", fieldName);
        }
        if (fieldType.isSimpleTypeReference()) {
            SimpleTypeReference simpleTypeReference = fieldType.asSimpleTypeReference().orElseThrow();
            if (simpleTypeReference.isVstringTypeReference()) {
                VstringTypeReference vstringTypeReference = simpleTypeReference.asVstringTypeReference().orElseThrow();
                getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", toSerializationExpression(typeDefinition, optionalField, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), null));
            } else {
                // TODO: Generate dynamic type length values here.
                getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", simpleTypeReference.getSizeInBits());
            }
        } else if (fieldType.isEnumTypeReference()) {
            // TODO: Generate dynamic type length values here.
            getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", getEnumBaseTypeReference(optionalField.getType()).getSizeInBits());
        } else if (fieldType.isDataIoTypeReference()) {
            DataIoTypeReference dataIoTypeReference = fieldType.asDataIoTypeReference().orElseThrow();
            // Build parameter string from DataIoTypeReference params
            StringBuilder paramsBuilder = new StringBuilder();
            paramsBuilder.append(optionalField.getName());
            List<Term> paramTerms = dataIoTypeReference.getParams().orElse(Collections.emptyList());
            for (int i = 0; i < paramTerms.size(); i++) {
                Term paramTerm = paramTerms.get(i);
                final TypeReference argumentType = getArgumentType(dataIoTypeReference, i);
                paramsBuilder.append(", (").append(getLanguageTypeNameForTypeReference(argumentType, true)).append(") (")
                    .append(toSerializationExpression(typeDefinition, optionalField, argumentType, paramTerm, parserArguments)).append(")");
            }
            getLengthInBitsCodeBlockBuilder.beginControlFlow("try");
            getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += ($L != null) ? $T.getLengthInBits($L) : 0", optionalField.getName(), ClassName.get(targetPackage, dataIoTypeReference.getName()), paramsBuilder.toString());
            getLengthInBitsCodeBlockBuilder.nextControlFlow("catch ($T e)", ClassName.get("org.apache.plc4x.java.spi.buffers.api.exceptions", "BufferException"));
            getLengthInBitsCodeBlockBuilder.add("// Ignore\n");
            getLengthInBitsCodeBlockBuilder.endControlFlow();
        } else {
            getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L.getLengthInBits()", optionalField.getName());
        }
        if (nullBytesHexAttribute.isEmpty() && !fieldType.isDataIoTypeReference()) {
            getLengthInBitsCodeBlockBuilder.endControlFlow();
        }

        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generatePaddingField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        PaddingField paddingField = field.asPaddingField().orElseThrow();
        String fieldName = "padding" + fieldIndex;
        CodeBlock attributesCodeBlock = getAttributes(typeDefinition, field, fieldName, parserArguments);

        TypeReference fieldType = paddingField.getType();
        TypeName fieldTypeClassName = getFieldTypeClassName(field);

        ClassName fieldReaderFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.reader", "FieldReaderFactory");
        ClassName fieldWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.writer", "FieldWriterFactory");
        CodeBlock readBlock = getDataReaderCall(fieldType);
        CodeBlock writeBlock = getDataWriterCall(fieldType);

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        parseBlockBuilder.add("// Padding Field: $L\n", fieldName);
        CodeBlock paddingConditionExpression = toParseExpression(typeDefinition, paddingField, paddingField.getType(), paddingField.getPaddingCondition(), parserArguments);
        parseBlockBuilder.addStatement("$T.readPaddingField($L, ($T) $L)", fieldReaderFactory, readBlock, int.class, paddingConditionExpression);
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        serializeCodeBlockBuilder.add("// Padding Field: $L\n", fieldName);
        CodeBlock paddingValueExpression = toParseExpression(typeDefinition, paddingField, paddingField.getType(), paddingField.getPaddingValue(), parserArguments);
        serializeCodeBlockBuilder.addStatement("$T.writePaddingField(($T) $L, ($T) $L, $L)", fieldWriterFactory, int.class, paddingConditionExpression, fieldTypeClassName, paddingValueExpression, writeBlock);
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Padding Field: $L\n", fieldName);
        // TODO: Generate dynamic type length values here.
        getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L * ((int) $L)", paddingField.getType().asSimpleTypeReference().orElseThrow().getSizeInBits(), paddingConditionExpression);
        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generatePeekField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        PeekField peekField = field.asPeekField().orElseThrow();
        String fieldName = peekField.getName();
        CodeBlock attributesCodeBlock = getAttributes(typeDefinition, field, fieldName, parserArguments);

        TypeReference fieldType = peekField.getType();
        TypeName fieldTypeClassName = getFieldTypeClassName(field);

        ClassName fieldReaderFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.reader", "FieldReaderFactory");
        CodeBlock readBlock = getDataReaderCall(fieldType);

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        parseBlockBuilder.add("// Peek Field: $L\n", fieldName);
        parseBlockBuilder.addStatement("$T $L = $T.readPeekField($L, $L)", fieldTypeClassName, peekField.getName(), fieldReaderFactory, readBlock, attributesCodeBlock);
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        serializeCodeBlockBuilder.add("// Peek Field: $L\n", fieldName);
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Peek Field: $L\n", fieldName);
        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generateReservedField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        ReservedField reservedField = field.asReservedField().orElseThrow();

        TypeReference fieldType = reservedField.getType();

        ClassName fieldReaderFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.reader", "FieldReaderFactory");
        ClassName fieldWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.writer", "FieldWriterFactory");
        CodeBlock readBlock = getDataReaderCall(fieldType);
        CodeBlock writeBlock = getDataWriterCall(fieldType);

        TypeName languageTypeName = getLanguageTypeNameForTypeReference(reservedField.getType(), true);
        CodeBlock reservedValue;
        if (ClassName.get(BigInteger.class).equals(languageTypeName)) {
            reservedValue = CodeBlock.of("$T.valueOf($L)", BigInteger.class, reservedField.getReferenceValue());
        } else {
            reservedValue = CodeBlock.of("($L) $L", languageTypeName, reservedField.getReferenceValue());
        }

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        // Outputting type+symbolic field name is extremely helpful when tracking down offset issues.
        CodeBlock withNameBlock = CodeBlock.of("$T.$L($S)", ClassName.get("org.apache.plc4x.java.spi.buffers.api", "WithOption"), "WithName", typeDefinition.getName() + ".reserved" + fieldIndex);
        parseBlockBuilder.add("// Reserved Field\n");
        parseBlockBuilder.addStatement("$T.readReservedField($L, $L, $L)", fieldReaderFactory, readBlock, reservedValue, withNameBlock);
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        serializeCodeBlockBuilder.add("// Reserved Field\n");
        serializeCodeBlockBuilder.addStatement("$T.writeReservedField($L, $L)", fieldWriterFactory, reservedValue, writeBlock);
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Reserved Field\n");
        // TODO: Generate dynamic type length values here.
        getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", fieldType.asSimpleTypeReference().orElseThrow().getSizeInBits());
        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generateSimpleField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        SimpleField simpleField = field.asSimpleField().orElseThrow();
        String fieldName = simpleField.getName();
        CodeBlock attributesCodeBlock = getAttributes(typeDefinition, field, fieldName, parserArguments);

        TypeReference fieldType = simpleField.getType();
        TypeName fieldTypeClassName = getLanguageTypeNameForTypeReference(simpleField.getType(), true);

        ClassName fieldReaderFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.reader", "FieldReaderFactory");
        ClassName fieldWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.writer", "FieldWriterFactory");
        ClassName dataWriterFactory = ClassName.get("org.apache.plc4x.java.spi.fields.data.writer", "DataWriterFactory");
        CodeBlock readBlock = getDataReaderCall(fieldType);
        CodeBlock writeBlock = getDataWriterCall(fieldType);

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        if (simpleField.getType().isEnumTypeReference()) {
            parseBlockBuilder.add("// Simple Field (enum): $L\n", fieldName);
            parseBlockBuilder.addStatement("$T $L = $T.readEnumField($L, $L)", fieldTypeClassName, simpleField.getName(), fieldReaderFactory, readBlock, attributesCodeBlock);
        } else {
            parseBlockBuilder.add("// Simple Field: $L\n", fieldName);
            parseBlockBuilder.addStatement("$T $L = $T.readSimpleField($L, $L)", fieldTypeClassName, simpleField.getName(), fieldReaderFactory, readBlock, attributesCodeBlock);
        }
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        if (simpleField.getType().isEnumTypeReference()) {
            EnumTypeReference enumTypeReference = simpleField.getType().asEnumTypeReference().orElseThrow();
            serializeCodeBlockBuilder.add("// Simple Field (enum): $L\n", fieldName);
            // TODO: Double-check if we can't use readBlock instead?
            serializeCodeBlockBuilder.addStatement("$T.writeSimpleEnumField(($T) $L, $T.writeEnum($T::getValue, $T::name, $L), $L)", fieldWriterFactory, fieldTypeClassName, getValueBlock, dataWriterFactory, fieldTypeClassName, fieldTypeClassName, writeBlock, attributesCodeBlock);
        } else {
            serializeCodeBlockBuilder.add("// Simple Field: $L\n", fieldName);
            serializeCodeBlockBuilder.addStatement("$T.writeSimpleField(($T) $L, $L, $L)", fieldWriterFactory, fieldTypeClassName, getValueBlock, writeBlock, attributesCodeBlock);
        }
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Simple Field: $L\n", fieldName);
        if (fieldType.isSimpleTypeReference()) {
            SimpleTypeReference simpleTypeReference = fieldType.asSimpleTypeReference().orElseThrow();
            if (simpleTypeReference.isVstringTypeReference()) {
                // Calculate the actual length of the string.
                VstringTypeReference vstringTypeReference = simpleTypeReference.asVstringTypeReference().orElseThrow();
                getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", toSerializationExpression(typeDefinition, simpleField, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), null));
            } else {
                // TODO: Generate dynamic type length values here.
                Optional<Term> unsignedIntegerEncodingAttribute = field.getAttribute("unsignedIntegerEncoding");
                if(unsignedIntegerEncodingAttribute.isPresent() && unsignedIntegerEncodingAttribute.get().stringRepresentation().startsWith("\"VAR-")) {
                    getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $T.GET_VAR_LENGTH_UINT_IN_BITS($L)", ClassName.get("org.apache.plc4x.java.spi.utils", "StaticHelper"), fieldName);
                } else {
                    Optional<Term> signedIntegerEncodingAttribute = field.getAttribute("signedIntegerEncoding");
                    if(signedIntegerEncodingAttribute.isPresent() && signedIntegerEncodingAttribute.get().stringRepresentation().startsWith("\"VAR-")) {
                        getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $T.GET_VAR_LENGTH_SINT_IN_BITS($L)", ClassName.get("org.apache.plc4x.java.spi.utils", "StaticHelper"), fieldName);
                    } else {
                        getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", simpleTypeReference.getSizeInBits());
                    }
                }
            }
        } else if (fieldType.isEnumTypeReference()) {
            // TODO: Generate dynamic type length values here.
            getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L", getEnumBaseTypeReference(simpleField.getType()).getSizeInBits());
        } else if (fieldType.isDataIoTypeReference()) {
            throw new RuntimeException("Simple fields of type dataIo are not supported.");
        } else {
            getLengthInBitsCodeBlockBuilder.addStatement("lengthInBits += $L.getLengthInBits()", simpleField.getName());
        }
        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generateStateField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        // TODO: Implement
    }

    protected void generateUnknownField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        UnknownField unknownField = field.asUnknownField().orElseThrow();
        String fieldName = "unknown" + fieldIndex;
        CodeBlock attributesCodeBlock = getAttributes(typeDefinition, field, fieldName, parserArguments);

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        parseBlockBuilder.add("// Unknown Field: $L\n", fieldName);
        parseBlockBuilder.addStatement("fail");
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        serializeCodeBlockBuilder.add("// Unknown Field: $L\n", fieldName);
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Unknown Field: $L\n", fieldName);
        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generateValidationField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) throws BufferException {
        ValidationField validationField = field.asValidationField().orElseThrow();

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        parseBlockBuilder.add("// Validation Field\n");
        parseBlockBuilder.beginControlFlow("if(!($L))", toParseExpression(typeDefinition, validationField, BOOL_TYPE_REFERENCE, validationField.getValidationExpression(), parserArguments));
        parseBlockBuilder.addStatement("throw new $T($L)", ClassName.get("org.apache.plc4x.java.spi.fields.exceptions", "ParseAssertException"), validationField.getDescription().orElse("\"Unknown validation error.\""));
        parseBlockBuilder.endControlFlow();
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        serializeCodeBlockBuilder.add("// Validation Field (Nothing needed here)\n");
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Validation Field (Nothing needed here)\n");
        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    protected void generateVirtualField(TypeDefinition typeDefinition, Field field, int fieldIndex, List<Argument> parserArguments, CodeBlock getValueBlock, List<CodeBlock> caseCodeBlocksParse, List<CodeBlock> caseCodeBlocksSerialize, List<CodeBlock> caseCodeBlocksGetLengthInBits) {
        VirtualField virtualField = field.asVirtualField().orElseThrow();
        String fieldName = virtualField.getName();
        CodeBlock attributesCodeBlock = getAttributes(typeDefinition, field, fieldName, parserArguments);

        TypeName fieldTypeClassName = getFieldTypeClassName(field);

        ClassName fieldReaderFactory = ClassName.get("org.apache.plc4x.java.spi.fields.fields.reader", "FieldReaderFactory");

        CodeBlock.Builder parseBlockBuilder = CodeBlock.builder();
        parseBlockBuilder.add("// Virtual Field: $L (doesn't parse anything, just makes the value available)\n", fieldName);
        parseBlockBuilder.addStatement("$T $L = $T.readVirtualField($T.class, $L, $L)", fieldTypeClassName, fieldName, fieldReaderFactory, fieldTypeClassName, getValueBlock, attributesCodeBlock);
        caseCodeBlocksParse.add(parseBlockBuilder.build());

        CodeBlock.Builder serializeCodeBlockBuilder = CodeBlock.builder();
        serializeCodeBlockBuilder.add("// Virtual Field: $L (doesn't serialize anything, just makes the value available)\n", fieldName);
        serializeCodeBlockBuilder.addStatement("$T $L = ($T) get$L()", fieldTypeClassName, fieldName, fieldTypeClassName, StaticHelper.CAPITALIZE(fieldName));
        caseCodeBlocksSerialize.add(serializeCodeBlockBuilder.build());

        CodeBlock.Builder getLengthInBitsCodeBlockBuilder = CodeBlock.builder();
        getLengthInBitsCodeBlockBuilder.add("// Virtual Field: $L (doesn't produce any output, just makes the value available)\n", fieldName);
        caseCodeBlocksGetLengthInBits.add(getLengthInBitsCodeBlockBuilder.build());
    }

    /**
     * Confirms if a variable references a discriminator variable. These need to be handled differently when serializing and parsing.
     *
     * @param variableLiteral The variable to search for.
     * @return boolean returns true if the variable's name references a discriminator field
     */
    protected boolean isVariableLiteralDiscriminatorField(TypeDefinition typeDefinition, VariableLiteral variableLiteral) {
        if ((typeDefinition == null) || !typeDefinition.isComplexTypeDefinition()) {
            return false;
        }
        return typeDefinition.asComplexTypeDefinition()
            .map(complexTypeDefinition -> complexTypeDefinition.isVariableLiteralDiscriminatorField(variableLiteral))
            .orElse(false);
    }

    /**
     * Confirms if a variable is an implicit variable. These need to be handled differently when serializing and parsing.
     *
     * @param variableLiteral The variable to search for.
     * @return boolean returns true if the variable's name is an implicit field
     */
    protected boolean isVariableLiteralImplicitField(TypeDefinition typeDefinition, VariableLiteral variableLiteral) {
        if ((typeDefinition == null) || !typeDefinition.isComplexTypeDefinition()) {
            return false;
        }
        return typeDefinition.asComplexTypeDefinition()
            .map(complexTypeDefinition -> complexTypeDefinition.isVariableLiteralImplicitField(variableLiteral))
            .orElse(false);
    }

    /**
     * Confirms if a variable is a virtual variable. These need to be handled differently when serializing and parsing.
     *
     * @param variableLiteral The variable to search for.
     * @return boolean returns true if the variable's name is an virtual field
     */
    protected boolean isVariableLiteralVirtualField(TypeDefinition typeDefinition, VariableLiteral variableLiteral) {
        if ((typeDefinition == null) || !typeDefinition.isComplexTypeDefinition()) {
            return false;
        }
        return typeDefinition.asComplexTypeDefinition()
            .map(complexTypeDefinition -> complexTypeDefinition.isVariableLiteralVirtualField(variableLiteral))
            .orElse(false);
    }

    /**
     * Returns the implicit field that has the same name as the variable. These need to be handled differently when serializing and parsing.
     *
     * @param variableLiteral The variable to search for.
     * @return ImplicitField returns the implicit field that corresponds to the variable's name.
     */
    protected ImplicitField getReferencedImplicitField(TypeDefinition typeDefinition, VariableLiteral variableLiteral) {
        if ((typeDefinition == null) || !typeDefinition.isComplexTypeDefinition()) {
            return null;
        }
        return typeDefinition.asComplexTypeDefinition()
            .map(complexTypeDefinition -> complexTypeDefinition.getReferencedImplicitField(variableLiteral))
            .orElse(null);
    }

    protected CodeBlock getAttributes(TypeDefinition typeDefinition, Field field, String fieldName, List<Argument> parserArguments) {
        List<CodeBlock> attributes = new ArrayList<>();
        CodeBlock withNameBlock = CodeBlock.of("$T.$L($S)", ClassName.get("org.apache.plc4x.java.spi.buffers.api", "WithOption"), "WithName", fieldName);
        attributes.add(withNameBlock);
        field.getCurrentAttributeNames().forEach(attributeName -> {
                String methodName = "With" + StaticHelper.CAPITALIZE(attributeName);
                try {
                    CodeBlock parseExpression = toParseExpression(typeDefinition, field, STRING_TYPE_REFERENCE, field.getAttribute(attributeName).orElseThrow(), parserArguments);
                    try {
                        WithOption.class.getDeclaredMethod(methodName, String.class);
                        attributes.add(CodeBlock.of("$T.$L($L)", WithOption.class, methodName, parseExpression));
                    } catch (NoSuchMethodException e) {
                        // TODO: Possibly solve this in a more pluggable way.
                        try {
                            WithByteBasedOption.class.getDeclaredMethod(methodName, String.class);
                            attributes.add(CodeBlock.of("$T.$L($L)", WithByteBasedOption.class, methodName, parseExpression));
                        } catch (NoSuchMethodException e2) {
                            // Just ignore this.
                        }
                    }
                } catch (BufferException e) {
                    // Just ignore this for now.
                }
            }
        );
        return CodeBlock.join(attributes, ", ");
    }

    public boolean hasFieldsWithNames(List<Field> fields, String... names) {
        for (String name : names) {
            boolean foundName = false;
            for (Field field : fields) {
                if (field instanceof NamedField && name.equals(((NamedField) field).getName())) {
                    foundName = true;
                    break;
                }
            }
            if (!foundName) {
                return false;
            }
        }
        // TODO: document why true is returned here.
        return true;
    }

    public boolean isRawByteArray(DiscriminatedComplexTypeDefinition currentCase) {
        Optional<Field> valueFieldOptional = currentCase.getFields().stream().filter(field -> field.isNamedField() && field.asNamedField().orElseThrow().getName().equals("value")).findFirst();
        if (valueFieldOptional.isPresent()) {
            Field valueField = valueFieldOptional.get();
            if (valueField.isTypedField()) {
                TypedField typedField = valueField.asTypedField().orElseThrow();
                return typedField.getType().isArrayTypeReference() && typedField.getType().asArrayTypeReference().orElseThrow().getElementTypeReference().isByteBased();
            }
        }
        return false;
    }

    public boolean isParsableAsNumber(String str) {
        if (str == null || str.isEmpty()) {
            return false;
        }
        try {
            Double.parseDouble(str);
            return true;
        } catch (NumberFormatException ex) {
            return false;
        }
    }

    @FunctionalInterface
    public interface ThrowingFunction<T, R, E extends Exception> {
        R apply(T t) throws E;
    }

}
