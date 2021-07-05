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
package org.apache.plc4x.plugins.codegenerator.protocol.freemarker;

import net.objecthunter.exp4j.Expression;
import net.objecthunter.exp4j.ExpressionBuilder;
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.references.ComplexTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.SimpleTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.StringTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.*;
import org.w3c.dom.Node;

import java.util.*;
import java.util.stream.Collectors;

public abstract class BaseFreemarkerLanguageTemplateHelper implements FreemarkerLanguageTemplateHelper {

    private final TypeDefinition thisType;
    private final String protocolName;
    private final String flavorName;
    private final Map<String, TypeDefinition> types;

    // In mspec we are using some virtual virtual fields that are useful for code generation.
    // As they should be shared over all language template implementations,
    // these are defined here manually.
    private static final Map<String, SimpleTypeReference> builtInFields;

    public static final String CUR_POS = "curPos";
    public static final String START_POS = "startPos";
    public static final String LAST_ITEM = "lastItem";
    public static final String IMPLICIT = "implicit";

    static {
        builtInFields = new HashMap<>();
        builtInFields.put(CUR_POS, new SimpleTypeReference() {
            @Override
            public SimpleBaseType getBaseType() {
                return SimpleBaseType.UINT;
            }

            @Override
            public int getSizeInBits() {
                return 16;
            }
        });
        builtInFields.put(START_POS, new SimpleTypeReference() {
            @Override
            public SimpleBaseType getBaseType() {
                return SimpleBaseType.UINT;
            }

            @Override
            public int getSizeInBits() {
                return 16;
            }
        });
        builtInFields.put(LAST_ITEM, new SimpleTypeReference() {
            @Override
            public SimpleBaseType getBaseType() {
                return SimpleBaseType.BIT;
            }

            @Override
            public int getSizeInBits() {
                return 1;
            }
        });
    }

    protected BaseFreemarkerLanguageTemplateHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types) {
        this.thisType = thisType;
        this.protocolName = protocolName;
        this.flavorName = flavorName;
        this.types = types;
    }

    protected TypeDefinition getThisTypeDefinition() {
        return thisType;
    }

    public String getProtocolName() {
        return protocolName;
    }

    public String getFlavorName() {
        return flavorName;
    }

    public Map<String, TypeDefinition> getTypeDefinitions() {
        return types;
    }

    public List<TypeDefinition> getComplexTypeRootDefinitions() {
        return types.values().stream().filter(typeDefinition -> (typeDefinition instanceof ComplexTypeDefinition) &&
            !(typeDefinition instanceof DiscriminatedComplexTypeDefinition)).collect(Collectors.toList());
    }

    protected static Map<String, SimpleTypeReference> getBuiltInFieldTypes() {
        return builtInFields;
    }

    /* *********************************************************************************
     * Methods that are language-dependent.
     **********************************************************************************/

    public abstract String getLanguageTypeNameForField(Field field);

    public abstract String getLanguageTypeNameForTypeReference(TypeReference typeReference);

    public String getReadBufferReadMethodCall(SimpleTypeReference simpleTypeReference) {
        return getReadBufferReadMethodCall(simpleTypeReference, null, null);
    }

    public abstract String getReadBufferReadMethodCall(SimpleTypeReference simpleTypeReference, String valueString, TypedField field);

    public abstract String getWriteBufferWriteMethodCall(SimpleTypeReference simpleTypeReference, String fieldName, TypedField field);

    public abstract String getNullValueForTypeReference(TypeReference typeReference);

    /* *********************************************************************************
     * Methods related to type-references.
     **********************************************************************************/

    /**
     * @param typeReference type reference
     * @return true if the given type reference is a simple type reference.
     */
    public boolean isSimpleTypeReference(TypeReference typeReference) {
        return typeReference instanceof SimpleTypeReference;
    }

    /**
     * @param typeReference type reference
     * @return true if the given type reference is a byte based type reference.
     */
    public boolean isByteBased(TypeReference typeReference) {
        if (!isSimpleTypeReference(typeReference)) {
            return false;
        }
        return ((SimpleTypeReference) typeReference).getBaseType() == SimpleTypeReference.SimpleBaseType.BYTE;
    }

    /**
     * @param typeReference type reference
     * @return true if the given type reference is a complex type reference.
     */
    public boolean isComplexTypeReference(TypeReference typeReference) {
        return typeReference instanceof ComplexTypeReference;
    }

    public boolean isEnumTypeReference(TypeReference typeReference) {
        if (!isComplexTypeReference(typeReference)) {
            return false;
        }
        return getTypeDefinitionForTypeReference(typeReference) instanceof EnumTypeDefinition;
    }

    public boolean isStringTypeReference(TypeReference typeReference) {
        return typeReference instanceof StringTypeReference;
    }

    /**
     * Helper for collecting referenced complex types as these usually ned to be
     * imported in some way.
     *
     * @return Collection of all complex type references used in fields or enum constants.
     */
    public Collection<String> getComplexTypeReferences() {
        return getComplexTypeReferences(thisType);
    }

    /**
     * Helper for collecting referenced complex types as these usually need to be
     * imported in some way.
     *
     * @param baseType the base type we want to get the type references from
     * @return collection of complex type references used in the type.
     */
    public Collection<String> getComplexTypeReferences(TypeDefinition baseType) {
        Set<String> complexTypeReferences = new HashSet<>();
        // If this is a subtype of a discriminated type, we have to add a reference to the parent type.
        if (baseType instanceof DiscriminatedComplexTypeDefinition) {
            DiscriminatedComplexTypeDefinition discriminatedComplexTypeDefinition = (DiscriminatedComplexTypeDefinition) baseType;
            if (!discriminatedComplexTypeDefinition.isAbstract()) {
                complexTypeReferences.add(((ComplexTypeReference)
                    discriminatedComplexTypeDefinition.getParentType().getTypeReference()).getName());
            }
        }
        // If it's a complex type definition, add all the types referenced by any property fields
        // (Includes any types referenced by sub-types in case this is a discriminated type parent)
        if (baseType instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) baseType;
            for (Field field : complexTypeDefinition.getFields()) {
                if (field instanceof PropertyField) {
                    PropertyField propertyField = (PropertyField) field;
                    if (propertyField.getType() instanceof ComplexTypeReference) {
                        ComplexTypeReference complexTypeReference = (ComplexTypeReference) propertyField.getType();
                        complexTypeReferences.add(complexTypeReference.getName());
                    }
                } else if (field instanceof SwitchField) {
                    SwitchField switchField = (SwitchField) field;
                    for (DiscriminatedComplexTypeDefinition switchCase : switchField.getCases()) {
                        complexTypeReferences.addAll(getComplexTypeReferences(switchCase));
                    }
                }
            }
        }
        // In case this is a enum type, we have to check all the constant types.
        else if (baseType instanceof EnumTypeDefinition) {
            for (String constantName : ((EnumTypeDefinition) baseType).getConstantNames()) {
                final TypeReference constantType = ((EnumTypeDefinition) thisType).getConstantType(constantName);
                if (constantType instanceof ComplexTypeReference) {
                    ComplexTypeReference complexTypeReference = (ComplexTypeReference) constantType;
                    complexTypeReferences.add(complexTypeReference.getName());
                }
            }
        }
        // If the type has any parser arguments, these have to be checked too.
        if (baseType.getParserArguments() != null) {
            for (Argument parserArgument : baseType.getParserArguments()) {
                if (parserArgument.getType() instanceof ComplexTypeReference) {
                    ComplexTypeReference complexTypeReference = (ComplexTypeReference) parserArgument.getType();
                    complexTypeReferences.add(complexTypeReference.getName());
                }
            }
        }
        return complexTypeReferences;
    }

    /**
     * Little helper to return the type of a given property.
     *
     * @param baseType     base type definition that contains the given property.
     * @param propertyName name of the property
     * @return the type reference of the given property
     */
    public Optional<TypeReference> getTypeReferenceForProperty(ComplexTypeDefinition baseType, String propertyName) {
        // If this is a built-in type, use that.
        if (builtInFields.containsKey(propertyName)) {
            return Optional.of(builtInFields.get(propertyName));
        }
        // Check if the expression root is referencing a field
        final Optional<PropertyField> propertyFieldOptional = baseType.getPropertyFields().stream()
            .filter(propertyField -> propertyField.getName().equals(propertyName))
            .findFirst();
        if (propertyFieldOptional.isPresent()) {
            final PropertyField propertyField = propertyFieldOptional.get();
            return Optional.of(propertyField.getType());
        }
        // Check if the expression is a ImplicitField
        final Optional<ImplicitField> implicitFieldOptional = baseType.getFields().stream()
            .filter(ImplicitField.class::isInstance)
            .map(ImplicitField.class::cast)
            .filter(implicitField -> implicitField.getName().equals(propertyName))
            .findFirst();
        if (implicitFieldOptional.isPresent()) {
            final ImplicitField implicitField = implicitFieldOptional.get();
            return Optional.of(implicitField.getType());
        }
        // Check if the expression is a VirtualField
        final Optional<VirtualField> virtualFieldOptional = baseType.getFields().stream()
            .filter(VirtualField.class::isInstance)
            .map(VirtualField.class::cast)
            .filter(virtualField -> virtualField.getName().equals(propertyName))
            .findFirst();
        if (virtualFieldOptional.isPresent()) {
            final VirtualField virtualField = virtualFieldOptional.get();
            return Optional.of(virtualField.getType());
        }
        // Check if the expression root is referencing an argument
        if (baseType.getParserArguments() != null) {
            final Optional<Argument> argumentOptional = Arrays.stream(baseType.getParserArguments())
                .filter(argument -> argument.getName().equals(propertyName))
                .findFirst();
            if (argumentOptional.isPresent()) {
                final Argument argument = argumentOptional.get();
                return Optional.of(argument.getType());
            }
        }
        // Check if the expression is a DiscriminatorField
        // This is a more theoretical case where the expression is referencing a discriminator value of the current type
        final Optional<DiscriminatorField> discriminatorFieldOptional = baseType.getFields().stream()
            .filter(DiscriminatorField.class::isInstance)
            .map(DiscriminatorField.class::cast)
            .filter(discriminatorField -> discriminatorField.getName().equals(propertyName))
            .findFirst();
        if (discriminatorFieldOptional.isPresent()) {
            final DiscriminatorField discriminatorField = discriminatorFieldOptional.get();
            return Optional.of(discriminatorField.getType());
        }
        return Optional.empty();
    }

    /**
     * Enums are always based on a main type. This helper accesses this information in a safe manner.
     *
     * @param typeReference type reference
     * @return simple type reference for the enum type referenced by the given type reference
     */
    public SimpleTypeReference getEnumBaseTypeReference(TypeReference typeReference) {
        if (!(typeReference instanceof ComplexTypeReference)) {
            throw new FreemarkerException("type reference for enum types must be of type complex type");
        }
        ComplexTypeReference complexTypeReference = (ComplexTypeReference) typeReference;
        final TypeDefinition typeDefinition = types.get(complexTypeReference.getName());
        if (typeDefinition == null) {
            throw new FreemarkerException("Couldn't find given enum type definition with name " + complexTypeReference.getName());
        }
        if (!(typeDefinition instanceof EnumTypeDefinition)) {
            throw new FreemarkerException("Referenced type with name " + complexTypeReference.getName() + " is not an enum type");
        }
        EnumTypeDefinition enumTypeDefinition = (EnumTypeDefinition) typeDefinition;
        // Enum types always have simple type references.
        return (SimpleTypeReference) enumTypeDefinition.getType();
    }

    /* *********************************************************************************
     * Methods related to fields.
     **********************************************************************************/

    public boolean hasFieldOfType(String fieldTypeName) {
        if (getThisTypeDefinition() instanceof ComplexTypeDefinition) {
            return ((ComplexTypeDefinition) getThisTypeDefinition()).getFields().stream().anyMatch(field -> field.getTypeName().equals(fieldTypeName));
        }
        return false;
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

    public Field getFieldForNameFromCurrentOrParent(String fieldName) {
        if (!(getThisTypeDefinition() instanceof ComplexTypeDefinition)) {
            return null;
        }
        return ((ComplexTypeDefinition) getThisTypeDefinition()).getAllPropertyFields()
            .stream()
            .filter(propertyField -> propertyField.getName().equals(fieldName))
            .findFirst()
            .orElse(null);
    }

    public Field getFieldForNameFromCurrent(String fieldName) {
        if (!(getThisTypeDefinition() instanceof ComplexTypeDefinition)) {
            return null;
        }
        return ((ComplexTypeDefinition) getThisTypeDefinition()).getPropertyFields()
            .stream()
            .filter(propertyField -> propertyField.getName().equals(fieldName))
            .findFirst()
            .orElse(null);
    }

    public boolean isAbstractField(Field field) {
        return field instanceof AbstractField;
    }

    public boolean isArrayField(Field field) {
        return field instanceof ArrayField;
    }

    public boolean isChecksumField(Field field) {
        return field instanceof ChecksumField;
    }

    public boolean isConstField(Field field) {
        return field instanceof ConstField;
    }

    public boolean isDiscriminatorField(Field field) {
        return field instanceof DiscriminatorField;
    }

    public boolean isEnumField(Field field) {
        if (!(field instanceof TypedField)) {
            return false;
        }
        TypedField typedField = (TypedField) field;
        TypeReference typeReference = typedField.getType();
        if (isSimpleTypeReference(typeReference)) {
            return false;
        }
        TypeDefinition typeDefinition = getTypeDefinitionForTypeReference(typedField.getType());
        return typeDefinition instanceof EnumTypeDefinition;
    }

    public boolean isImplicitField(Field field) {
        return field instanceof ImplicitField;
    }

    public boolean isManualArrayField(Field field) {
        return field instanceof ManualArrayField;
    }

    public boolean isNamedField(Field field) {
        return field instanceof NamedField;
    }

    public boolean isOptionalField(Field field) {
        return field instanceof OptionalField;
    }

    public boolean isPaddingField(Field field) {
        return field instanceof PaddingField;
    }

    public boolean isPropertyField(Field field) {
        return field instanceof PropertyField;
    }

    public boolean isReservedField(Field field) {
        return field instanceof ReservedField;
    }

    public boolean isSimpleField(Field field) {
        return field instanceof SimpleField;
    }

    public boolean isSwitchField(Field field) {
        return field instanceof SwitchField;
    }

    public boolean isTypedField(Field field) {
        return field instanceof TypedField;
    }

    public boolean isVirtualField(Field field) {
        return field instanceof VirtualField;
    }

    public boolean isCountArrayField(Field field) {
        if (field instanceof ArrayField) {
            ArrayField arrayField = (ArrayField) field;
            return arrayField.getLoopType() == ArrayField.LoopType.COUNT;
        }
        if (field instanceof ManualArrayField) {
            ManualArrayField arrayField = (ManualArrayField) field;
            return arrayField.getLoopType() == ManualArrayField.LoopType.COUNT;
        }
        return false;
    }

    public boolean isLengthArrayField(Field field) {
        if (field instanceof ArrayField) {
            ArrayField arrayField = (ArrayField) field;
            return arrayField.getLoopType() == ArrayField.LoopType.LENGTH;
        }
        if (field instanceof ManualArrayField) {
            ManualArrayField arrayField = (ManualArrayField) field;
            return arrayField.getLoopType() == ManualArrayField.LoopType.LENGTH;
        }
        return false;
    }

    public boolean isTerminatedArrayField(Field field) {
        if (field instanceof ArrayField) {
            ArrayField arrayField = (ArrayField) field;
            return arrayField.getLoopType() == ArrayField.LoopType.TERMINATED;
        }
        if (field instanceof ManualArrayField) {
            ManualArrayField arrayField = (ManualArrayField) field;
            return arrayField.getLoopType() == ManualArrayField.LoopType.TERMINATED;
        }
        return false;
    }

    /**
     * @return switch field of the current base type.
     */
    public SwitchField getSwitchField() {
        return getSwitchField(thisType);
    }

    /**
     * @return switch field of the provided base type.
     */
    protected SwitchField getSwitchField(TypeDefinition typeDefinition) {
        if (!(typeDefinition instanceof ComplexTypeDefinition)) {
            return null;
        }
        ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) typeDefinition;
        // Sebastian would be proud of me ;-)
        return (SwitchField) complexTypeDefinition.getFields().stream()
            .filter(SwitchField.class::isInstance)
            .findFirst()
            .orElse(null);
    }

    public Collection<Field> getPropertyAndSwitchFields() {
        return getPropertyAndSwitchFields(thisType);
    }

    public Collection<Field> getPropertyAndSwitchFields(TypeDefinition typeDefinition) {
        if (!(thisType instanceof ComplexTypeDefinition)) {
            return Collections.emptyList();
        }
        return ((ComplexTypeDefinition) thisType).getFields().stream()
            .filter(field -> (field instanceof PropertyField) || (field instanceof SwitchField))
            .collect(Collectors.toList());
    }

    /* *********************************************************************************
     * Methods related to type-definitions.
     **********************************************************************************/

    public boolean isDiscriminatedParentTypeDefinition() {
        return isDiscriminatedParentTypeDefinition(thisType);
    }

    public boolean isDiscriminatedParentTypeDefinition(TypeDefinition typeDefinition) {
        return (typeDefinition instanceof ComplexTypeDefinition) && ((ComplexTypeDefinition) typeDefinition).isAbstract();
    }

    public boolean isDiscriminatedChildTypeDefinition() {
        return isDiscriminatedChildTypeDefinition(thisType);
    }

    public boolean isDiscriminatedChildTypeDefinition(TypeDefinition typeDefinition) {
        return (typeDefinition instanceof DiscriminatedComplexTypeDefinition) && !((ComplexTypeDefinition) typeDefinition).isAbstract();
    }

    public TypeDefinition getTypeDefinitionForTypeReference(TypeReference typeReference) {
        if (!isComplexTypeReference(typeReference)) {
            throw new FreemarkerException("Type reference must be a complex type reference");
        }
        return getTypeDefinitions().get(((ComplexTypeReference) typeReference).getName());
    }

    /**
     * @return list of sub-types for the current base type or an empty collection, if there are none
     */
    public List<DiscriminatedComplexTypeDefinition> getSubTypeDefinitions() {
        return getSubTypeDefinitions(thisType);
    }

    /**
     * @return list of sub-types for a given type definition or an empty collection, if there are none
     */
    public List<DiscriminatedComplexTypeDefinition> getSubTypeDefinitions(TypeDefinition type) {
        SwitchField switchField = getSwitchField(type);
        if (switchField != null) {
            return switchField.getCases();
        }
        return Collections.emptyList();
    }

    /* *********************************************************************************
     * Methods related to terms and expressions.
     **********************************************************************************/

    /**
     * Check if the expression doesn't reference any variables.
     * If this is the case, the expression can be evaluated at code-generation time.
     *
     * @param term term
     * @return true if it doesn't reference any variable literals.
     */
    protected boolean isFixedValueExpression(Term term) {
        if (term instanceof VariableLiteral) {
            return false;
        }
        if (term instanceof UnaryTerm) {
            UnaryTerm unaryTerm = (UnaryTerm) term;
            return isFixedValueExpression(unaryTerm.getA());
        }
        if (term instanceof BinaryTerm) {
            BinaryTerm binaryTerm = (BinaryTerm) term;
            return isFixedValueExpression(binaryTerm.getA()) && isFixedValueExpression(binaryTerm.getB());
        }
        if (term instanceof TernaryTerm) {
            TernaryTerm ternaryTerm = (TernaryTerm) term;
            return isFixedValueExpression(ternaryTerm.getA()) && isFixedValueExpression(ternaryTerm.getB()) &&
                isFixedValueExpression(ternaryTerm.getC());
        }
        return true;
    }

    protected int evaluateFixedValueExpression(Term term) {
        final Expression expression = new ExpressionBuilder(toString(term)).build();
        return (int) expression.evaluate();
    }

    protected String toString(Term term) {
        if (term instanceof NullLiteral) {
            return "null";
        }
        if (term instanceof BooleanLiteral) {
            return Boolean.toString(((BooleanLiteral) term).getValue());
        }
        if (term instanceof NumericLiteral) {
            return ((NumericLiteral) term).getNumber().toString();
        }
        if (term instanceof StringLiteral) {
            return "\"" + ((StringLiteral) term).getValue() + "\"";
        }
        if (term instanceof UnaryTerm) {
            return ((UnaryTerm) term).getOperation() + toString(((UnaryTerm) term).getA());
        }
        if (term instanceof BinaryTerm) {
            return toString(((BinaryTerm) term).getA()) + ((BinaryTerm) term).getOperation() + toString(((BinaryTerm) term).getB());
        }
        if (term instanceof TernaryTerm) {
            return "(" + toString(((TernaryTerm) term).getA()) + ") ? (" + toString(((TernaryTerm) term).getB()) +
                ") : (" + toString(((TernaryTerm) term).getC()) + ")";
        }
        return "";
    }

    /* *********************************************************************************
     * Methods related to discriminators.
     **********************************************************************************/

    private Optional<TypeReference> getDiscriminatorType(ComplexTypeDefinition parentType, Term disciminatorExpression) {
        if (!(disciminatorExpression instanceof VariableLiteral)) {
            throw new FreemarkerException("Currently no arithmetic expressions are supported as discriminator expressions.");
        }
        VariableLiteral variableLiteral = (VariableLiteral) disciminatorExpression;
        Optional<TypeReference> type = getTypeReferenceForProperty(parentType, variableLiteral.getName());
        // If we found something but there's a "rest" left, we got to use the type we
        // found in this level, get that type's definition and continue from there.
        if (type.isPresent() && (variableLiteral.getChild() != null)) {
            TypeReference typeReference = type.get();
            if (typeReference instanceof ComplexTypeReference) {
                ComplexTypeReference complexTypeReference = (ComplexTypeReference) typeReference;
                final TypeDefinition typeDefinition = this.types.get(complexTypeReference.getName());
                if (typeDefinition instanceof ComplexTypeDefinition) {
                    return getDiscriminatorType((ComplexTypeDefinition) typeDefinition, variableLiteral.getChild());
                }
            }
        }
        return type;
    }

    /**
     * Get an ordered list of generated names for the discriminators.
     * These names can be used to access the type definitions as well as well as the values.
     *
     * @return list of symbolic names for the discriminators.
     */
    public List<String> getDiscriminatorNames() {
        TypeDefinition baseType = thisType;
        if (thisType.getParentType() != null) {
            baseType = thisType.getParentType();
        }
        final SwitchField switchField = getSwitchField(baseType);
        List<String> discriminatorNames = new ArrayList<>();
        if (switchField != null) {
            for (Term discriminatorExpression : switchField.getDiscriminatorExpressions()) {
                discriminatorNames.add(getDiscriminatorName(discriminatorExpression));
            }
        }
        return discriminatorNames;
    }

    /**
     * Check if there's any field with the given name.
     * This is required to suppress the generation of a discriminator field
     * in case a named field is providing the information.
     *
     * @param discriminatorName name of the discriminator name
     * @return true if a field with the given name already exists in the same type.
     */
    public boolean isNonDiscriminatorField(String discriminatorName) {
        return ((ComplexTypeDefinition) thisType).getAllPropertyFields().stream().anyMatch(
            field -> !(field instanceof DiscriminatorField) && field.getName().equals(discriminatorName));
    }

    /**
     * Check if there's any field with the given name.
     * This is required to suppress the generation of a virtual field
     * in case a discriminated field is providing the information.
     *
     * @param discriminatorName name of the virtual name
     * @return true if a field with the given name already exists in the same type.
     */
    public boolean isDiscriminatorField(String discriminatorName) {
        List<String> names = getDiscriminatorNames();
        if (names != null) {
            return getDiscriminatorNames().stream().anyMatch(
                field -> field.equals(discriminatorName));
        }
        return false;
    }

    /**
     * Converts a given discriminator description into a symbolic name.
     *
     * @param discriminatorExpression discriminator expression
     * @return name
     */
    public String getDiscriminatorName(Term discriminatorExpression) {
        if (discriminatorExpression instanceof Literal) {
            Literal literal = (Literal) discriminatorExpression;
            if (literal instanceof NullLiteral) {
                return "null";
            } else if (literal instanceof BooleanLiteral) {
                return Boolean.toString(((BooleanLiteral) literal).getValue());
            } else if (literal instanceof NumericLiteral) {
                return ((NumericLiteral) literal).getNumber().toString();
            } else if (literal instanceof StringLiteral) {
                return ((StringLiteral) literal).getValue();
            } else if (literal instanceof VariableLiteral) {
                VariableLiteral variableLiteral = (VariableLiteral) literal;
                return getVariableLiteralName(variableLiteral);
            }
        } else if (discriminatorExpression instanceof UnaryTerm) {
            UnaryTerm unaryTerm = (UnaryTerm) discriminatorExpression;
            return getDiscriminatorName(unaryTerm.getA());
        } else if (discriminatorExpression instanceof BinaryTerm) {
            BinaryTerm binaryTerm = (BinaryTerm) discriminatorExpression;
            return getDiscriminatorName(binaryTerm.getA()) + "_" + getDiscriminatorName(binaryTerm.getB());
        } else if (discriminatorExpression instanceof TernaryTerm) {
            TernaryTerm ternaryTerm = (TernaryTerm) discriminatorExpression;
            return getDiscriminatorName(ternaryTerm.getA()) + "_" + getDiscriminatorName(ternaryTerm.getB())
                + "_" + getDiscriminatorName(ternaryTerm.getC());
        }
        return "";
    }

    private String getVariableLiteralName(VariableLiteral variableLiteral) {
        String rest = "";
        if (variableLiteral.getChild() != null) {
            rest = getVariableLiteralName(variableLiteral.getChild());
            rest = rest.substring(0, 1).toUpperCase() + rest.substring(1);
        }
        return variableLiteral.getName() + rest;
    }

    /**
     * Get a list of the types for every discriminator name.
     *
     * @return Map mapping discriminator names to types.
     */
    public Map<String, TypeReference> getDiscriminatorTypes() {
        // Get the parent type (Which contains the typeSwitch field)
        ComplexTypeDefinition parentType;
        if (thisType instanceof DiscriminatedComplexTypeDefinition) {
            parentType = (ComplexTypeDefinition) thisType.getParentType();
        } else {
            parentType = (ComplexTypeDefinition) thisType;
        }
        // Get the typeSwitch field from that.
        final SwitchField switchField = getSwitchField(parentType);
        if (switchField == null) {
            return Collections.emptyMap();
        }
        Map<String, TypeReference> discriminatorTypes = new TreeMap<>();
        for (Term discriminatorExpression : switchField.getDiscriminatorExpressions()) {
            // Get some symbolic name we can use.
            String discriminatorName = getDiscriminatorName(discriminatorExpression);
            Optional<TypeReference> discriminatorType = getDiscriminatorType(parentType, discriminatorExpression);
            discriminatorTypes.put(discriminatorName, discriminatorType.orElse(null));
        }
        return discriminatorTypes;
    }

    public Map<String, String> getDiscriminatorValues(TypeDefinition type) {
        if (!(type instanceof DiscriminatedComplexTypeDefinition)) {
            return Collections.emptyMap();
        }
        DiscriminatedComplexTypeDefinition switchType = (DiscriminatedComplexTypeDefinition) type;
        final List<String> discriminatorNames = getDiscriminatorNames();
        final Map<String, String> discriminatorValues = new LinkedHashMap<>();
        for (int i = 0; i < discriminatorNames.size(); i++) {
            String discriminatorValue;
            if (i < switchType.getDiscriminatorValues().length) {
                discriminatorValue = switchType.getDiscriminatorValues()[i];
            } else {
                discriminatorValue = null;
            }
            discriminatorValues.put(discriminatorNames.get(i), discriminatorValue);
        }
        return discriminatorValues;
    }

    /**
     * Get a list of the values for every discriminator name for every discriminated type.
     *
     * @return Map mapping discriminator names to discriminator values for every discriminated type.
     */
    public Map<String, Map<String, String>> getDiscriminatorValues() {
        // Get the parent type (Which contains the typeSwitch field)
        ComplexTypeDefinition parentType;
        if (thisType instanceof DiscriminatedComplexTypeDefinition) {
            parentType = (ComplexTypeDefinition) thisType.getParentType();
        } else {
            parentType = (ComplexTypeDefinition) thisType;
        }
        // Get the typeSwitch field from that.
        final SwitchField switchField = getSwitchField(parentType);
        if (switchField == null) {
            return Collections.emptyMap();
        }
        // Build a map containing the named discriminator values for every case of the typeSwitch.
        Map<String, Map<String, String>> discriminatorTypes = new LinkedHashMap<>();
        for (DiscriminatedComplexTypeDefinition switchCase : switchField.getCases()) {
            discriminatorTypes.put(switchCase.getName(), getDiscriminatorValues(switchCase));
        }
        return discriminatorTypes;
    }

    public TypeReference getArgumentType(TypeReference typeReference, int index) {
        if (!(typeReference instanceof ComplexTypeReference)) {
            throw new FreemarkerException("Only complex type references supported here.");
        }
        ComplexTypeReference complexTypeReference = (ComplexTypeReference) typeReference;
        if (!getTypeDefinitions().containsKey(complexTypeReference.getName())) {
            throw new FreemarkerException("Could not find definition of complex type " + complexTypeReference.getName());
        }
        TypeDefinition complexTypeDefinition = getTypeDefinitions().get(complexTypeReference.getName());
        if (complexTypeDefinition.getParserArguments().length <= index) {
            throw new FreemarkerException("Type " + complexTypeReference.getName() + " specifies too few parser arguments");
        }
        return complexTypeDefinition.getParserArguments()[index].getType();
    }

    /**
     * Filters out the arguments that are user for serializiation.
     *
     * @param arguments list of all arguments.
     * @return list of arguments that are used during serialization.
     */
    public List<Argument> getSerializerArguments(Argument[] arguments) {
        if (arguments == null) {
            return Collections.emptyList();
        }
        List<Argument> serializerArguments = new LinkedList<>();
        for (Argument argument : arguments) {
            if (LAST_ITEM.equals(argument.getName())) {
                serializerArguments.add(argument);
            }
        }
        return serializerArguments;
    }

    public List<Term> getSerializerTerms(Term[] terms) {
        if (terms == null) {
            return Collections.emptyList();
        }
        List<Term> serializerTerms = new LinkedList<>();
        for (Term term : terms) {
            if (term.contains(LAST_ITEM)) {
                serializerTerms.add(term);
            }
        }
        return serializerTerms;
    }

    public boolean hasLastItemTerm(Term[] terms) {
        if (terms == null) {
            return false;
        }
        for (Term term : terms) {
            if (term.contains(LAST_ITEM)) {
                return true;
            }
        }
        return false;
    }

    public boolean discriminatorValueNeedsStringEqualityCheck(Term term) {
        if (!(term instanceof VariableLiteral)) {
            return false;
        }

        VariableLiteral variableLiteral = (VariableLiteral) term;
        // If this literal references an Enum type, then we have to output it differently.
        if (getTypeDefinitions().get(variableLiteral.getName()) instanceof EnumTypeDefinition) {
            return false;
        }

        if (getThisTypeDefinition() instanceof ComplexTypeDefinition) {
            Field referencedField = ((ComplexTypeDefinition) getThisTypeDefinition()).getFields().stream().filter(field -> ((field instanceof NamedField) && ((NamedField) field).getName().equals(variableLiteral.getName()))).findFirst().orElse(null);
            if (referencedField instanceof TypedField
                && ((TypedField) referencedField).getType() instanceof StringTypeReference) {
                return true;
            }
        }
        if (getThisTypeDefinition().getParserArguments() != null) {
            for (Argument parserArgument : getThisTypeDefinition().getParserArguments()) {
                if (parserArgument.getName().equals(variableLiteral.getName())
                    && parserArgument.getType() instanceof StringTypeReference) {
                    return true;
                }
            }
        }
        return false;
    }

    public boolean isEnumExpression(String expression) {
        if (!expression.contains(".")) {
            return false;
        }
        String enumName = expression.substring(0, expression.indexOf('.'));
        TypeDefinition typeDefinition = this.getTypeDefinitions().get(enumName);
        return (typeDefinition instanceof EnumTypeDefinition);
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

    public Collection<EnumValue> getEnumValuesForUniqueConstantValues(EnumValue[] enumValues, String constantName) {
        Map<String, EnumValue> filteredEnumValues = new TreeMap<>();
        for (EnumValue enumValue : enumValues) {
            if (!filteredEnumValues.containsKey(enumValue.getConstant(constantName))) {
                filteredEnumValues.put(enumValue.getConstant(constantName), enumValue);
            }
        }
        return filteredEnumValues.values();
    }

    public Collection<EnumValue> getEnumValuesForConstantValue(EnumValue[] enumValues, String constantName, String constantValue) {
        List<EnumValue> filteredEnumValues = new ArrayList<>();
        for (EnumValue enumValue : enumValues) {
            if (enumValue.getConstant(constantName).equals(constantValue)) {
                filteredEnumValues.add(enumValue);
            }
        }
        return filteredEnumValues;
    }

    public SimpleTypeReference getEnumFieldSimpleTypeReference(TypeReference type, String fieldName) {
        TypeDefinition typeDefinition = getTypeDefinitionForTypeReference(type);

        if (typeDefinition instanceof EnumTypeDefinition
            && ((EnumTypeDefinition) typeDefinition).getConstantType(fieldName) instanceof SimpleTypeReference) {
            return (SimpleTypeReference) ((EnumTypeDefinition) typeDefinition).getConstantType(fieldName);
        }
        return null;
    }

    /**
     * Confirms if a variable is an implicit variable. These need to be handled differently when serializing and parsing.
     *
     * @param vl The variable to search for.
     * @return boolean returns true if the variable's name is an implicit field
     */
    protected boolean isVariableLiteralImplicitField(VariableLiteral vl) {
        List<Field> fields = null;
        if (thisType instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexType = (ComplexTypeDefinition) getThisTypeDefinition();
            fields = complexType.getFields();
        }
        if (fields == null) {
            return false;
        }
        for (Field field : fields) {
            if (field.getTypeName().equals(IMPLICIT)) {
                ImplicitField implicitField = (ImplicitField) field;
                if (vl.getName().equals(implicitField.getName())) {
                    return true;
                }
            }
        }
        return false;
    }

    /**
     * Returns the implicit field that has the same name as the variable. These need to be handled differently when serializing and parsing.
     *
     * @param vl The variable to search for.
     * @return ImplicitField returns the implicit field that corresponds to the variable's name.
     */
    protected ImplicitField getReferencedImplicitField(VariableLiteral vl) {
        List<Field> fields = null;
        if (thisType instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexType = (ComplexTypeDefinition) getThisTypeDefinition();
            fields = complexType.getFields();
        }
        if (fields == null) {
            return null;
        }
        for (Field field : fields) {
            if (field.getTypeName().equals(IMPLICIT)) {
                ImplicitField implicitField = (ImplicitField) field;
                if (vl.getName().equals(implicitField.getName())) {
                    return implicitField;
                }
            }
        }
        return null;
    }

}
