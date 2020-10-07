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

package org.apache.plc4x.language.go;

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

@SuppressWarnings({"unused", "WeakerAccess"})
public class GoLanguageTemplateHelper extends BaseFreemarkerLanguageTemplateHelper {

    public GoLanguageTemplateHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types) {
        super(thisType, protocolName, flavorName, types);
    }

    public String fileName(String protocolName, String languageName, String languageFlavorName) {
        return "plc4go." + String.join("", languageName.split("\\-")) + "." +
            String.join("", protocolName.split("\\-")) + "." +
            String.join("", languageFlavorName.split("\\-"));
    }

    public String packageName(String languageFlavorName) {
        return String.join("", languageFlavorName.split("\\-"));
    }

    @Override
    public String getLanguageTypeNameForField(Field field) {
        boolean optional = field instanceof OptionalField;
        // If the referenced type is a DataIo type, the value is of type PlcValue.
        if(field instanceof PropertyField) {
            PropertyField propertyField = (PropertyField) field;
            if(propertyField.getType() instanceof ComplexTypeReference) {
                ComplexTypeReference complexTypeReference = (ComplexTypeReference) propertyField.getType();
                final TypeDefinition typeDefinition = getTypeDefinitions().get(complexTypeReference.getName());
                if(typeDefinition instanceof DataIoTypeDefinition) {
                    return "PlcValue";
                }
            }
        }
        return getLanguageTypeNameForTypeReference(((TypedField) field).getType());
    }

    @Override
    public String getLanguageTypeNameForTypeReference(TypeReference typeReference) {
        if(typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT: {
                    return "bool";
                }
                case UINT: {
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 8) {
                        return "uint8";
                    }
                    if (integerTypeReference.getSizeInBits() <= 16) {
                        return "uint16";
                    }
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return "uint32";
                    }
                    if (integerTypeReference.getSizeInBits() <= 64) {
                        return "uint64";
                    }
                    throw new RuntimeException("Unsupported simple type");
                }
                case INT: {
                    IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                    if (integerTypeReference.getSizeInBits() <= 8) {
                        return "int8";
                    }
                    if (integerTypeReference.getSizeInBits() <= 16) {
                        return "int16";
                    }
                    if (integerTypeReference.getSizeInBits() <= 32) {
                        return "int32";
                    }
                    if (integerTypeReference.getSizeInBits() <= 64) {
                        return "int64";
                    }
                    throw new RuntimeException("Unsupported simple type");
                }
                case FLOAT:
                case UFLOAT: {
                    FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                    int sizeInBits = ((floatTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.FLOAT) ? 1 : 0) +
                        floatTypeReference.getExponent() + floatTypeReference.getMantissa();
                    if (sizeInBits <= 32) {
                        return "float32";
                    }
                    if (sizeInBits <= 64) {
                        return "float64";
                    }
                    throw new RuntimeException("Unsupported simple type");
                }
                case STRING: {
                    return "string";
                }
                case TIME: {
                    return "Time";
                }
                case DATE: {
                    return "Date";
                }
                case DATETIME: {
                    return "Date";
                }
            }
            throw new RuntimeException("Unsupported simple type");
        } else {
            return ((ComplexTypeReference) typeReference).getName();
        }
    }

    @Override
    public String getNullValueForTypeReference(TypeReference typeReference) {
        if(typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT: {
                    return "false";
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

    public String getReadBufferReadMethodCall(SimpleTypeReference simpleTypeReference, String valueString) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT: {
                return "io.ReadBit()";
            }
            case UINT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "io.ReadUint8(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "io.ReadUint16(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "io.ReadUint32(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "io.ReadUint64(" + integerTypeReference.getSizeInBits() + ")";
                }
            }
            case INT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "io.ReadInt8(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "io.ReadInt16(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "io.ReadInt32(" + integerTypeReference.getSizeInBits() + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "io.ReadInt64(" + integerTypeReference.getSizeInBits() + ")";
                }
            }
            case FLOAT: {
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                String type = (floatTypeReference.getSizeInBits() <= 32)? "ReadFloat32" : "ReadFloat64";
                String typeCast = (floatTypeReference.getSizeInBits() <= 32) ? "float32" : "float64";
                String defaultNull = (floatTypeReference.getSizeInBits() <= 32) ? "0.0f" : "0.0";
                StringBuilder sb = new StringBuilder("((Supplier<").append(type).append(">) (() -> {");
                sb.append("\n            return (").append(typeCast).append(") toFloat(io, ").append(
                    (floatTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.FLOAT) ? "true" : "false")
                    .append(", ").append(floatTypeReference.getExponent()).append(", ")
                    .append(floatTypeReference.getMantissa()).append(");");
                sb.append("\n        })).get()");
                return sb.toString();
            }
            case STRING: {
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return "io.ReadString(" + stringTypeReference.getSizeInBits() + ", \"" +
                    stringTypeReference.getEncoding() + "\")";
            }
        }
        return "Hurz";
    }

    @Override
    public String getWriteBufferWriteMethodCall(SimpleTypeReference simpleTypeReference, String fieldName) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT: {
                return "io.writeBit((boolean) " + fieldName + ")";
            }
            case UINT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 4) {
                    return "io.writeUnsignedByte(" + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").byteValue())";
                }
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "io.writeUnsignedShort(" + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").shortValue())";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "io.writeUnsignedInt(" + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").intValue())";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "io.writeUnsignedLong(" + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").longValue())";
                }
                return "io.writeUnsignedBigInteger(" + integerTypeReference.getSizeInBits() + ", (BigInteger) " + fieldName + ")";
            }
            case INT: {
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "io.writeByte(" + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").byteValue())";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "io.writeShort(" + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").shortValue())";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "io.writeInt(" + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").intValue())";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "io.writeLong(" + integerTypeReference.getSizeInBits() + ", ((Number) " + fieldName + ").longValue())";
                }
                return "io.writeBigInteger(" + integerTypeReference.getSizeInBits() + ", BigInteger.valueOf( " + fieldName + "))";
            }
            case FLOAT:
            case UFLOAT: {
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                StringBuilder sb = new StringBuilder();
                if(simpleTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.FLOAT) {
                    sb.append("\n        boolean negative = value < 0;");
                    sb.append("\n        io.writeBit(negative);");
                }
                sb.append("\n        final int exponent = Math.getExponent(value);");
                sb.append("\n        final double mantissa = value / Math.pow(2, exponent);");
                sb.append("\n        io.writeInt(").append(floatTypeReference.getExponent()).append(", exponent);");
                sb.append("\n        io.writeDouble(").append(floatTypeReference.getMantissa()).append(", mantissa)");
                return sb.toString().substring(9);
            }
            case STRING: {
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return "io.writeString(" + stringTypeReference.getSizeInBits() + ", \"" +
                    stringTypeReference.getEncoding() + "\", (String) " + fieldName + ")";
            }
        }
        return "Hurz";
    }

    public String getReservedValue(ReservedField reservedField) {
        final String languageTypeName = getLanguageTypeNameForTypeReference(reservedField.getType());
        if("BigInteger".equals(languageTypeName)) {
            return "BigInteger.valueOf(" + reservedField.getReferenceValue() + ")";
        } else {
            return "(" + languageTypeName + ") " + reservedField.getReferenceValue();
        }
    }

    public String toParseExpression(TypedField field, Term term, Argument[] parserArguments) {
        return toExpression(field, term, term1 -> toVariableParseExpression(field, term1, parserArguments));
    }

    public String toSerializationExpression(TypedField field, Term term, Argument[] parserArguments) {
        return toExpression(field, term, term1 -> toVariableSerializationExpression(field, term1, parserArguments));
    }

    private String toExpression(TypedField field, Term term, Function<Term, String> variableExpressionGenerator) {
        if(term == null) {
            return "";
        }
        if(term instanceof Literal) {
            if(term instanceof NullLiteral) {
                return "null";
            } else if(term instanceof BooleanLiteral) {
                return Boolean.toString(((BooleanLiteral) term).getValue());
            } else if(term instanceof NumericLiteral) {
                return ((NumericLiteral) term).getNumber().toString();
            } else if(term instanceof StringLiteral) {
                return "\"" + ((StringLiteral) term).getValue() + "\"";
            } else if(term instanceof VariableLiteral) {
                VariableLiteral variableLiteral = (VariableLiteral) term;
                // If this literal references an Enum type, then we have to output it differently.
                if(getTypeDefinitions().get(variableLiteral.getName()) instanceof EnumTypeDefinition) {
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
            switch(ut.getOperation()) {
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
            if("if".equals(tt.getOperation())) {
                Term a = tt.getA();
                Term b = tt.getB();
                Term c = tt.getC();
                return "((" +  toExpression(field, a, variableExpressionGenerator) + ") ? " + toExpression(field, b, variableExpressionGenerator) + " : " + toExpression(field, c, variableExpressionGenerator) + ")";
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
        if("CAST".equals(vl.getName())) {
            StringBuilder sb = new StringBuilder(vl.getName());
            if((vl.getArgs() == null) || (vl.getArgs().size() != 2)) {
                throw new RuntimeException("A CAST expression expects exactly two arguments.");
            }
            sb.append("(").append(toVariableParseExpression(field, vl.getArgs().get(0), parserArguments))
                .append(", ").append(((VariableLiteral) vl.getArgs().get(1)).getName()).append(".class)");
            return sb.toString() + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "");
        }
        else if("STATIC_CALL".equals(vl.getName())) {
            StringBuilder sb = new StringBuilder();
            if(!(vl.getArgs().get(0) instanceof StringLiteral)) {
                throw new RuntimeException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral");
            }
            // Get the class and method name
            String methodName = ((StringLiteral) vl.getArgs().get(0)).getValue();
            // Cut off the double-quptes
            methodName = methodName.substring(1, methodName.length() - 1);
            sb.append(methodName).append("(");
            for(int i = 1; i < vl.getArgs().size(); i++) {
                Term arg = vl.getArgs().get(i);
                if(i > 1) {
                    sb.append(", ");
                }
                if(arg instanceof VariableLiteral) {
                    VariableLiteral va = (VariableLiteral) arg;
                    // "io" is the default name of the reader argument which is always available.
                    boolean isParserArg = "io".equals(va.getName());
                    boolean isTypeArg = "_type".equals(va.getName());
                    if(!isParserArg && !isTypeArg && parserArguments != null) {
                        for (Argument parserArgument : parserArguments) {
                            if (parserArgument.getName().equals(va.getName())) {
                                isParserArg = true;
                                break;
                            }
                        }
                    }
                    if(isParserArg) {
                        sb.append(va.getName() + ((va.getChild() != null) ? "." + toVariableExpressionRest(va.getChild()) : ""));
                    }
                    // We have to manually evaluate the type information at code-generation time.
                    else if(isTypeArg) {
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
                } else if(arg instanceof StringLiteral) {
                    sb.append(((StringLiteral) arg).getValue());
                }
            }
            sb.append(")");
            return sb.toString();
        }
        // All uppercase names are not fields, but utility methods.
        else if(vl.getName().equals(vl.getName().toUpperCase())) {
            StringBuilder sb = new StringBuilder(vl.getName());
            if(vl.getArgs() != null) {
                sb.append("(");
                boolean firstArg = true;
                for(Term arg : vl.getArgs()) {
                    if(!firstArg) {
                        sb.append(", ");
                    }
                    sb.append(toParseExpression(field, arg, parserArguments));
                    firstArg = false;
                }
                sb.append(")");
            }
            if(vl.getIndex() != VariableLiteral.NO_INDEX) {
                sb.append("[").append(vl.getIndex()).append("]");
            }
            return sb.toString() + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "");
        }
        return vl.getName() + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "");
    }

    private String toVariableSerializationExpression(TypedField field, Term term, Argument[] serialzerArguments) {
        VariableLiteral vl = (VariableLiteral) term;
        if("STATIC_CALL".equals(vl.getName())) {
            StringBuilder sb = new StringBuilder();
            if(!(vl.getArgs().get(0) instanceof StringLiteral)) {
                throw new RuntimeException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral");
            }
            String methodName = ((StringLiteral) vl.getArgs().get(0)).getValue();
            methodName = methodName.substring(1, methodName.length() - 1);
            sb.append(methodName).append("(");
            for(int i = 1; i < vl.getArgs().size(); i++) {
                Term arg = vl.getArgs().get(i);
                if(i > 1) {
                    sb.append(", ");
                }
                if(arg instanceof VariableLiteral) {
                    VariableLiteral va = (VariableLiteral) arg;
                    // "io" and "_value" are always available in every parser.
                    boolean isSerializerArg = "io".equals(va.getName()) || "_value".equals(va.getName()) || "element".equals(va.getName());
                    boolean isTypeArg = "_type".equals(va.getName());
                    if(!isSerializerArg && !isTypeArg && serialzerArguments != null) {
                        for (Argument serializerArgument : serialzerArguments) {
                            if (serializerArgument.getName().equals(va.getName())) {
                                isSerializerArg = true;
                                break;
                            }
                        }
                    }
                    if(isSerializerArg) {
                        sb.append(va.getName() + ((va.getChild() != null) ? "." + toVariableExpressionRest(va.getChild()) : ""));
                    } else if(isTypeArg) {
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
                        sb.append(toVariableSerializationExpression(field, va, null));
                    }
                } else if(arg instanceof StringLiteral) {
                    sb.append(((StringLiteral) arg).getValue());
                }
            }
            sb.append(")");
            return sb.toString();
        }
        // All uppercase names are not fields, but utility methods.
        else if(vl.getName().equals(vl.getName().toUpperCase())) {
            StringBuilder sb = new StringBuilder(vl.getName());
            if(vl.getArgs() != null) {
                sb.append("(");
                boolean firstArg = true;
                for(Term arg : vl.getArgs()) {
                    if(!firstArg) {
                        sb.append(", ");
                    }

                    if(arg instanceof VariableLiteral) {
                        VariableLiteral va = (VariableLiteral) arg;
                        boolean isSerializerArg = "io".equals(va.getName());
                        boolean isTypeArg = "_type".equals(va.getName());
                        if(!isSerializerArg && !isTypeArg && serialzerArguments != null) {
                            for (Argument serializerArgument : serialzerArguments) {
                                if (serializerArgument.getName().equals(va.getName())) {
                                    isSerializerArg = true;
                                    break;
                                }
                            }
                        }
                        if(isSerializerArg) {
                            sb.append(va.getName() + ((va.getChild() != null) ? "." + toVariableExpressionRest(va.getChild()) : ""));
                        } else if(isTypeArg) {
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
                            sb.append(toVariableSerializationExpression(field, va, null));
                        }
                    } else if(arg instanceof StringLiteral) {
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
        if(!isSerializerArg && !isTypeArg && serialzerArguments != null) {
            for (Argument serializerArgument : serialzerArguments) {
                if (serializerArgument.getName().equals(vl.getName())) {
                    isSerializerArg = true;
                    break;
                }
            }
        }
        if(isSerializerArg) {
            return vl.getName() + ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : "");
        } else if(isTypeArg) {
            String part = vl.getChild().getName();
            switch (part) {
                case "name":
                    return"\"" + field.getTypeName() + "\"";
                case "length":
                    return"\"" + ((SimpleTypeReference) field).getSizeInBits() + "\"";
                case "encoding":
                    String encoding = ((StringTypeReference) field.getType()).getEncoding();
                    // Cut off the single quotes.
                    encoding = encoding.substring(1, encoding.length() - 1);
                    return"\"" + encoding + "\"";
                default:
                    return "";
            }
        } else {
            return "_value." + toVariableExpressionRest(vl);
        }
    }

    private String toVariableExpressionRest(VariableLiteral vl) {
        return "get" + WordUtils.capitalize(vl.getName()) + "()" + ((vl.isIndexed() ? "[" + vl.getIndex() + "]" : "") +
            ((vl.getChild() != null) ? "." + toVariableExpressionRest(vl.getChild()) : ""));
    }

    public String getSizeInBits(ComplexTypeDefinition complexTypeDefinition, Argument[] parserArguments) {
        int sizeInBits = 0;
        StringBuilder sb = new StringBuilder("");
        for (Field field : complexTypeDefinition.getFields()) {
            if(field instanceof ArrayField) {
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
            } else if(field instanceof TypedField) {
                TypedField typedField = (TypedField) field;
                final TypeReference type = typedField.getType();
                if(field instanceof ManualField) {
                    ManualField manualField = (ManualField) field;
                    sb.append("(").append(toSerializationExpression(null, manualField.getLengthExpression(), parserArguments)).append(") + ");
                }
                else if(type instanceof SimpleTypeReference) {
                    SimpleTypeReference simpleTypeReference = (SimpleTypeReference) type;
                    sizeInBits += simpleTypeReference.getSizeInBits();
                } else {
                    // No ComplexTypeReference supported.
                }
            }
        }
        return sb.toString() + sizeInBits;
    }

    public String escapeValue(TypeReference typeReference, String valueString) {
        if(valueString == null) {
            return null;
        }
        if(typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case UINT:
                case INT:
                    // If it's a one character string and is numeric, output it as char.
                    if(!NumberUtils.isParsable(valueString) && (valueString.length() == 1)) {
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
            return constantName;
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
