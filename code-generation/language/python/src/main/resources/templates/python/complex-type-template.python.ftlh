<#--
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License") you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
-->
<#-- Prevent freemarker from escaping stuff -->
<#outputformat "undefined">
<#-- Declare the name and type of variables passed in to the template -->
<#-- @ftlvariable name="languageName" type="java.lang.String" -->
<#-- @ftlvariable name="protocolName" type="java.lang.String" -->
<#-- @ftlvariable name="outputFlavor" type="java.lang.String" -->
<#-- @ftlvariable name="helper" type="org.apache.plc4x.language.java.JavaLanguageTemplateHelper" -->
<#-- @ftlvariable name="tracer" type="org.apache.plc4x.plugins.codegenerator.protocol.freemarker.Tracer" -->
<#-- @ftlvariable name="type" type="org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition" -->
<#-- Declare the name and type of variables declared locally inside the template -->
${helper.packageName(protocolName, languageName, outputFlavor)?replace(".", "/")}/${type.name}.py
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

from dataclasses import dataclass
<#macro importSectionWithContentBelow>
    <#local sectionContent><#nested></#local>
    <#if helper.getRequiredImports()?has_content>
        <#list helper.getRequiredImports() as import>
${import}
        </#list>
    </#if>
    ${sectionContent}
</#macro>

<#macro emitImport import>${helper.emitRequiredImport(import)}</#macro>

<@importSectionWithContentBelow>
<#-- TODO: the code below implies that parserArguments will be null if not present... not pretty  -->
<#if type.parserArguments.isPresent()><#assign parserArguments=type.allParserArguments.orElseThrow()></#if>

<@emitImport import="from plc4py.api.messages.PlcMessage import PlcMessage" />
<@emitImport import="from plc4py.spi.generation.WriteBuffer import WriteBuffer" />
@dataclass
class ${type.name}<#if type.isDiscriminatedParentTypeDefinition()></#if>(<#if type.isDiscriminatedParentTypeDefinition()>ABC, PlcMessage<@emitImport import="from abc import ABC" /></#if><#if type.parentType.isPresent()>${type.parentType.orElseThrow().name}<@emitImport import="from plc4py.protocols.${protocolName}.${outputFlavor?replace('-', '')}.${type.parentType.orElseThrow().name} import ${type.parentType.orElseThrow().name}" /></#if>):
    <#if parserArguments?has_content>
        <#assign filteredParserArguments=parserArguments?filter(arg -> !type.isDiscriminatorField(arg.name) && !type.getPropertyFieldFromThisOrParentByName(arg.name).isPresent())>
    </#if>
    <#-- Property fields are fields that require a property in the pojo -->
    <#if type.propertyFields?has_content>
        <#list type.propertyFields as field>
    ${helper.camelCaseToSnakeCase(field.name)}: ${helper.getLanguageTypeNameForTypeReference(field.type, !field.isOptionalField())}
        </#list>
    </#if>
    <#if filteredParserArguments?has_content>
    # Arguments.
        <#list filteredParserArguments as parserArgument>
    ${helper.camelCaseToSnakeCase(parserArgument.name)}: ${helper.getLanguageTypeNameForTypeReference(parserArgument.type)}
        </#list>
    </#if>
    <#assign reservedFields=type.getFields()?filter(f->f.isReservedField())>
    <#if reservedFields?has_content>
    # Reserved Fields
        <#list reservedFields as reservedField>
    reserved_field_${reservedField?index}: ${helper.getLanguageTypeNameForTypeReference(reservedField.asReservedField().orElseThrow().type, false)}
        </#list>
    </#if>
    <#-- If the current type contains "const" fields, generate some Python constants for holing their values -->
    <#if type.constFields?has_content>
        <#list type.constFields as field>
    ${helper.camelCaseToSnakeCase(field.name)?upper_case}: ${helper.getLanguageTypeNameForTypeReference(field.type)} = ${helper.toParseExpression(field, field.type, field.referenceValue, parserArguments)}
        </#list>
    </#if>
<#--
    If this is a discriminated child type, we need to generate methods for accessing it's discriminator
    values, as if they were normal java properties.
-->
<#if type.isDiscriminatedChildTypeDefinition()>
    <#assign discriminatedChildType = type.asDiscriminatedComplexTypeDefinition().orElseThrow()>
    # Accessors for discriminator values.
    <#list discriminatedChildType.getDiscriminatorMap() as discriminatorName, discriminatorValue>
        <#-- If the discriminator name matches that of another field, suppress the methods generation -->
        <#if !discriminatedChildType.isNonDiscriminatorField(discriminatorName)><#--&& !discriminatedChildType.isParserArgument(discriminatorName)-->
            <#assign discriminatorType = helper.getDiscriminatorTypes()[discriminatorName]>
                <@emitImport import="from typing import ClassVar" />
    ${helper.camelCaseToSnakeCase(discriminatorName)}: ClassVar[${helper.getLanguageTypeNameForTypeReference(discriminatorType)}] = <@compress single_line=true>
            <#if discriminatorValue?? && !helper.isWildcard(discriminatorValue)>
                <#if discriminatorType.isEnumTypeReference()>
        ${helper.getLanguageTypeNameForTypeReference(discriminatorType)}.${discriminatorValue.name}
                <#else>
        ${helper.toParseExpression(null, discriminatorType, discriminatorValue, parserArguments)}
                </#if>
            <#else>
        ${helper.getNullValueForTypeReference(discriminatorType)}
            </#if>
            </@compress>

        </#if>
    </#list>
</#if>

<#--
    If this is a discriminated parent type, we need to generate the abstract methods for accessing it's
    discriminator values instead.
-->
    <#if type.isDiscriminatedParentTypeDefinition()>
        <#assign discriminatedParentType = type>
    <#-- @ftlvariable name="discriminatedParentType" type="org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition" -->
    # Abstract accessors for discriminator values.
        <#list helper.discriminatorTypes as discriminatorName, discriminatorType>
        <#-- If the discriminator name matches that of another field, suppress the methods generation -->
            <#if !type.isNonDiscriminatorField(discriminatorName)><#-- && !type.isParserArgument(discriminatorName)-->
    <@emitImport import="from abc import abstractmethod" />
    @property
    def ${helper.camelCaseToSnakeCase(discriminatorName)}(self) -> ${helper.getLanguageTypeNameForTypeReference(discriminatorType)}:
        pass
            </#if>
        </#list>
    </#if>

<#list type.abstractFields as field>
    <@emitImport import="from abc import abstractmethod" />
    @property
    def ${helper.camelCaseToSnakeCase(field.asNamedField().orElseThrow().name)}(self) -> ${helper.getLanguageTypeNameForField(field)}:
        pass

</#list>
<#list type.virtualFields as field>
    @property
    def ${helper.camelCaseToSnakeCase(field.name)}(self) -> ${helper.getLanguageTypeNameForField(field)}:
        <#if helper.getLanguageTypeNameForField(field) = 'String'>
        return ${helper.getLanguageTypeNameForField(field)}(${helper.toSerializationExpression(field, field.type, field.valueExpression, parserArguments)})
        <#--elseif helper.getLanguageTypeNameForField(field) = 'BigInteger' && !helper.isBigIntegerSource(field.valueExpression)-->
        <#elseif helper.getLanguageTypeNameForField(field) = 'BigInteger'>
        o = ${helper.toSerializationExpression(field, field.type, field.valueExpression, parserArguments)}
        if isinstance(o, bignum):
            return o
        return bignum(o)
        <#else>
        return ${helper.getLanguageTypeNameForField(field)}(${helper.toSerializationExpression(field, field.type, field.valueExpression, parserArguments)})
        </#if>


</#list>

    <#if outputFlavor != "passive">
<#if type.isDiscriminatedChildTypeDefinition()>
    def serialize_${helper.camelCaseToSnakeCase(type.parentType.orElseThrow().name)}_child(self, write_buffer: WriteBuffer):
<#else>
    <#if type.isDiscriminatedParentTypeDefinition()>
    <@emitImport import="from abc import abstractmethod" />
    @abstractmethod
    def serialize_${helper.camelCaseToSnakeCase(type.name)}_child(self, write_buffer: WriteBuffer) -> None:
        pass

    </#if>
    def serialize(self, write_buffer: WriteBuffer):
</#if>
        <#if helper.hasFieldOfType("unknown")>
        <@emitImport import="from plc4py.api.exceptions.exceptions import SerializationException" />
        raise SerializationException("Unknown field not serializable")
        <#else>
        <#if helper.requiresStartPos() || helper.requiresCurPos()>
            start_pos: int = write_buffer.get_pos()
        </#if>
        write_buffer.push_context("${type.name}")
            <#assign reservedFieldIndex=0>
            <#list type.fields as field>
                <#switch field.typeName>
                    <#case "array">
                        <#assign arrayField = field.asArrayField().orElseThrow()>
                        <#assign typedField = field.asTypedField().orElseThrow()>
                        <#assign namedField = field.asNamedField().orElseThrow()>

        # Array Field (${arrayField.name})
                        <#if arrayField.type.elementTypeReference.isByteBased()>
        write_buffer.write_byte_array(self.${helper.camelCaseToSnakeCase(namedField.name)}, logical_name="${helper.camelCaseToSnakeCase(namedField.name)}")
                        <#elseif arrayField.type.elementTypeReference.isSimpleTypeReference()>
        write_buffer.write_simple_array(self.${helper.camelCaseToSnakeCase(namedField.name)}, write_buffer.${helper.getDataWriterCall(arrayField.type.elementTypeReference, namedField.name)}, logical_name="${helper.camelCaseToSnakeCase(namedField.name)}")
                        <#else>
        write_buffer.write_complex_array(self.${helper.camelCaseToSnakeCase(namedField.name)}, logical_name="${helper.camelCaseToSnakeCase(namedField.name)}")
                        </#if>
                        <#break>
                    <#case "checksum">
                        <#assign checksumField = field.asChecksumField().orElseThrow()>
                        <#assign typedField = field.asTypedField().orElseThrow()>
                        <#assign namedField = field.asNamedField().orElseThrow()>
                        <@emitImport import="from plc4py.protocols.${protocolName} import StaticHelper" />

        # Checksum Field (checksum) (Calculated)
        write_buffer.${helper.getDataWriterCall(typedField.type, namedField.name)}(${helper.getLanguageTypeNameForField(field)}(StaticHelper.${helper.toParseExpression(checksumField, checksumField.type, checksumField.checksumExpression, parserArguments)}), logical_name="${helper.camelCaseToSnakeCase(namedField.name)}")
                        <#break>
                    <#case "const">
                        <#assign constField = field.asConstField().orElseThrow()>
                        <#assign typedField = field.asTypedField().orElseThrow()>
                        <#assign namedField = field.asNamedField().orElseThrow()>

        # Const Field (${constField.name})
        <#if typedField.type.isEnumTypeReference()>write_buffer.${helper.getDataWriterCall(helper.getEnumBaseTypeReference(typedField.type), namedField.name)}(self.${helper.camelCaseToSnakeCase(namedField.name)?upper_case}.value, logical_name="${constField.name}")<#else>write_buffer.${helper.getDataWriterCall(typedField.type, namedField.name)}(self.${helper.camelCaseToSnakeCase(namedField.name)?upper_case}, logical_name="${constField.name}")</#if>
                        <#break>
                    <#case "discriminator">
                        <#assign discriminatorField = field.asDiscriminatorField().orElseThrow()>
                        <#assign typedField = field.asTypedField().orElseThrow()>
                        <#assign namedField = field.asNamedField().orElseThrow()>

        # Discriminator Field (${discriminatorField.name}) (Used as input to a switch field)
                        <#if typedField.type.isEnumTypeReference()>
        write_buffer.${helper.getEnumDataWriterCall(typedField.type, namedField.name, "value")}(self.${helper.camelCaseToSnakeCase(discriminatorField.name)}, logical_name="${helper.camelCaseToSnakeCase(namedField.name)}" ${helper.getFieldOptions(typedField, parserArguments)})
                        <#else>
        <#assign simpleTypeReference = discriminatorField.type.asSimpleTypeReference().orElseThrow()>
        write_buffer.${helper.getDataWriterCall(typedField.type, namedField.name)}(self.${helper.camelCaseToSnakeCase(discriminatorField.name)}, logical_name="${helper.camelCaseToSnakeCase(namedField.name)}", bit_length=${simpleTypeReference.sizeInBits} ${helper.getFieldOptions(typedField, parserArguments)})
                        </#if>
                        <#break>
                    <#case "enum">
                        <#assign enumField = field.asEnumField().orElseThrow()>
                        <#assign typedField = field.asTypedField().orElseThrow()>
                        <#assign namedField = field.asNamedField().orElseThrow()>

        # Enum field (${namedField.name})
        write_buffer.${helper.getEnumDataWriterCall(typedField.type, namedField.name, enumField.fieldName)}(self.${helper.camelCaseToSnakeCase(namedField.name)}, logical_name="${helper.camelCaseToSnakeCase(namedField.name)}")
                        <#break>
                    <#case "implicit">
                        <#assign implicitField = field.asImplicitField().orElseThrow()>
                        <#assign typedField = field.asTypedField().orElseThrow()>
                        <#assign namedField = field.asNamedField().orElseThrow()>

        # Implicit Field (${helper.camelCaseToSnakeCase(implicitField.name)}) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
                        <#-- Implicit field values might be used in expressions, in order to avoid problems, we generate a temporary variable with the given name. -->
        ${helper.camelCaseToSnakeCase(implicitField.name)}: ${helper.getLanguageTypeNameForField(field)} = (${helper.toSerializationExpression(implicitField, implicitField.type, implicitField.serializeExpression, parserArguments)})
        write_buffer.${helper.getDataWriterCall(typedField.type, namedField.name)}(${helper.camelCaseToSnakeCase(implicitField.name)}, logical_name="${helper.camelCaseToSnakeCase(namedField.name)}")
                        <#break>
                    <#case "manualArray">
                        <#assign manualArrayField = field.asManualArrayField().orElseThrow()>
                        <#assign typedField = field.asTypedField().orElseThrow()>
                        <#assign namedField = field.asNamedField().orElseThrow()>

        # Manual Array Field (${manualArrayField.name})
        write_buffer.${helper.getDataWriterCall(typedField.type, namedField.name)}((${helper.getLanguageTypeNameForTypeReference(manualArrayField.type.elementTypeReference)} _value) -> ${helper.toParseExpression(manualArrayField, manualArrayField.type.elementTypeReference, manualArrayField.serializeExpression, parserArguments)}, logical_name="${helper.camelCaseToSnakeCase(namedField.name)}")
                        <#break>
                    <#case "manual">
                        <#assign manualField = field.asManualField().orElseThrow()>
                        <#assign typedField = field.asTypedField().orElseThrow()>
                        <#assign namedField = field.asNamedField().orElseThrow()>
                        <@emitImport import="from plc4py.protocols.${protocolName} import StaticHelper" />
        # Manual Field (${manualField.name})
        write_buffer.write_manual(write_function=lambda : StaticHelper.${helper.toSerializationExpression(manualField, manualField.type, manualField.serializeExpression, parserArguments)}, logical_name="${helper.camelCaseToSnakeCase(namedField.name)}")
                        <#break>
                    <#case "optional">
                        <#assign optionalField = field.asOptionalField().orElseThrow()>
                        <#assign typedField = field.asTypedField().orElseThrow()>
                        <#assign namedField = field.asNamedField().orElseThrow()>

        # Optional Field (${optionalField.name}) (Can be skipped, if the value is null)
                        <#if optionalField.conditionExpression.present>
        if ${helper.toSerializationExpression(optionalField, helper.boolTypeReference, optionalField.conditionExpression.get(), parserArguments)}:
                            <#if optionalField.type.isEnumTypeReference()>
            write_buffer.${helper.getEnumDataWriterCall(optionalField.type, optionalField.name, "value")}(self.${helper.camelCaseToSnakeCase(optionalField.name)}, logical_name="${optionalField.name}")
                            <#elseif optionalField.type.isDataIoTypeReference()>
            write_buffer.${helper.getDataWriterCall(typedField.type, optionalField.name)}(self.${helper.camelCaseToSnakeCase(optionalField.name)}, logical_name="${optionalField.name}", <#list optionalField.type.asComplexTypeReference().orElseThrow().params.orElseThrow() as param>${helper.toParseExpression(optionalField, helper.anyTypeReference, param, parserArguments)}<#sep>, </#sep></#list>)))
                            <#else>
            write_buffer.${helper.getDataWriterCall(typedField.type, optionalField.name)}(self.${helper.camelCaseToSnakeCase(optionalField.name)}, logical_name="${optionalField.name}"<#if optionalField.conditionExpression.present></#if>)
                            </#if>
                        </#if>
                        <#break>
                    <#case "padding">
                        <#assign paddingField = field.asPaddingField().orElseThrow()>
                        <#assign typedField = field.asTypedField().orElseThrow()>

        # Padding Field (padding)
        write_buffer.${helper.getDataWriterCall(typedField.type, "padding")}(int(${helper.toParseExpression(paddingField, helper.intTypeReference, paddingField.paddingCondition, parserArguments)}), logical_name="padding")
                        <#break>
                    <#case "reserved">
                        <#assign reservedField = field.asReservedField().orElseThrow()>
                        <#assign typedField = field.asTypedField().orElseThrow()>

        # Reserved Field (reserved)
        write_buffer.${helper.getDataWriterCall(typedField.type, "reserved")}(self.reserved_field_${reservedFieldIndex} if self.reserved_field_${reservedFieldIndex} is not None else ${helper.getReservedValue(reservedField)}, logical_name="reserved")<#assign reservedFieldIndex=reservedFieldIndex+1>
                        <#break>
                    <#case "simple">
                        <#assign simpleField = field.asSimpleField().orElseThrow()>
                        <#assign typedField = field.asTypedField().orElseThrow()>
                        <#assign namedField = field.asNamedField().orElseThrow()>

        # Simple Field (${namedField.name})
                        <#if typedField.type.isEnumTypeReference()>
        write_buffer.${helper.getEnumDataWriterCall(simpleField.type, simpleField.name, "value")}(self.${helper.camelCaseToSnakeCase(simpleField.name)}, logical_name="${simpleField.name}")
                        <#elseif simpleField.type.isDataIoTypeReference()>
        write_buffer.${helper.getEnumDataWriterCall(simpleField.type, simpleField.name, "value")}(DataWriterDataIoDefault(write_buffer, (wb, val) -> ${simpleField.type.asComplexTypeReference().orElseThrow().name}.staticSerialize(wb, val<#if simpleField.type.asComplexTypeReference().orElseThrow().params?has_content>, <#list simpleField.type.asComplexTypeReference().orElseThrow().params.orElseThrow() as param>${helper.toParseExpression(simpleField, helper.anyTypeReference, param, parserArguments)}<#sep>, </#sep></#list></#if>, logical_name="${simpleField.name}")
                        <#else>
                            <#if simpleField.type.isSimpleTypeReference()>
                                <#assign simpleTypeReference = simpleField.type.asSimpleTypeReference().orElseThrow()>
        write_buffer.${helper.getDataWriterCall(typedField.type, simpleField.name)}(self.${helper.camelCaseToSnakeCase(simpleField.name)},bit_length=${simpleTypeReference.sizeInBits},logical_name="${simpleField.name}")
                            <#else>
        write_buffer.${helper.getDataWriterCall(typedField.type, simpleField.name)}(self.${helper.camelCaseToSnakeCase(simpleField.name)},logical_name="${simpleField.name}")
                            </#if>
                        </#if>
                        <#break>
                    <#case "switch">
                        <#assign switchField = field.asSwitchField().orElseThrow()>

        # Switch field (Serialize the sub-type)
        self.serialize_${helper.camelCaseToSnakeCase(type.name)}_child(write_buffer)
                        <#break>
                    <#case "virtual">
                        <#assign virtualField = field.asVirtualField().orElseThrow()>
                        <#assign typedField = field.asTypedField().orElseThrow()>
                        <#assign namedField = field.asNamedField().orElseThrow()>

        # Virtual field (doesn't actually serialize anything, just makes the value available)
        ${namedField.name}: ${helper.getLanguageTypeNameForField(field)}  = self.${helper.camelCaseToSnakeCase(namedField.name)}()
                        <#break>
                </#switch>
            </#list>

        write_buffer.pop_context("${type.name}")
        </#if>

    </#if>

    <@emitImport import="import math" />
    def length_in_bytes(self) -> int:
        return int(math.ceil(float(self.length_in_bits() / 8.0)))

    def length_in_bits(self) -> int:
        length_in_bits: int = <#if type.parentType.isPresent()>super().length_in_bits()<#else>0</#if>
        _value: ${type.name} = self
<#list type.fields as field>
<#switch field.typeName>
    <#case "array">
        <#assign arrayField = field.asArrayField().orElseThrow()>
        <#assign arrayElementTypeReference = arrayField.type.asArrayTypeReference().orElseThrow().getElementTypeReference()>

        # Array field
        if self.${helper.camelCaseToSnakeCase(arrayField.name)} is not None:
        <#if arrayElementTypeReference.isSimpleTypeReference()>
            <#assign simpleTypeReference = arrayElementTypeReference.asSimpleTypeReference().orElseThrow()>
            length_in_bits += ${simpleTypeReference.sizeInBits} * len(self.${helper.camelCaseToSnakeCase(arrayField.name)})
        <#elseif arrayField.isCountArrayField()>
            <#assign nonSimpleTypeReference = arrayElementTypeReference.asNonSimpleTypeReference().orElseThrow()>
            for element in self.${helper.camelCaseToSnakeCase(arrayField.name)}:
                length_in_bits += element.length_in_bits()

        <#else>
            for element in self.${helper.camelCaseToSnakeCase(arrayField.name)}:
                length_in_bits += element.length_in_bits()

        </#if>

        <#break>
    <#case "checksum">
        <#assign checksumField = field.asChecksumField().orElseThrow()>
        <#assign typedField = field.asTypedField().orElseThrow()>
        <#assign simpleTypeReference = typedField.type.asSimpleTypeReference().orElseThrow()>

        # Checksum Field (checksum)
        length_in_bits += ${simpleTypeReference.sizeInBits}
        <#break>
    <#case "const">
        <#assign constField = field.asConstField().orElseThrow()>
        <#assign typedField = field.asTypedField().orElseThrow()>

        # Const Field (${constField.name})
        <#if typedField.type.isSimpleTypeReference()>
        <#assign simpleTypeReference = typedField.type.asSimpleTypeReference().orElseThrow()>
        length_in_bits += ${simpleTypeReference.sizeInBits}
        <#else>
        length_in_bits += ${helper.getEnumBaseTypeReference(typedField.type).sizeInBits}
        </#if>
        <#break>
    <#case "discriminator">
        <#assign discriminatorField = field.asDiscriminatorField().orElseThrow()>

        # Discriminator Field (${discriminatorField.name})
        <#if discriminatorField.type.isSimpleTypeReference()>
            <#assign simpleTypeReference = discriminatorField.type.asSimpleTypeReference().orElseThrow()>
            <#if simpleTypeReference.isVstringTypeReference()>
                <#assign vstringTypeReference = simpleTypeReference.asVstringTypeReference().orElseThrow()>
        length_in_bits += ${helper.toSerializationExpression(discriminatorField, helper.intTypeReference, vstringTypeReference.getLengthExpression(), parserArguments)}
            <#else>
        length_in_bits += ${simpleTypeReference.sizeInBits}
            </#if>
        <#elseif helper.isEnumField(field)>
        length_in_bits += ${helper.getEnumBaseTypeReference(discriminatorField.type).sizeInBits}
        <#else>
        length_in_bits += self.${helper.camelCaseToSnakeCase(discriminatorField.name)}.length_in_bits()
        </#if>
        <#break>
    <#case "enum">
        <#assign enumField = field.asEnumField().orElseThrow()>

        # Enum Field (${enumField.name})
        length_in_bits += ${helper.getEnumBaseTypeReference(enumField.type).sizeInBits}
        <#break>
    <#case "implicit">
        <#assign implicitField = field.asImplicitField().orElseThrow()>
        <#assign simpleTypeReference = implicitField.type.asSimpleTypeReference().orElseThrow()>

        # Implicit Field (${implicitField.name})
        length_in_bits += ${simpleTypeReference.sizeInBits}
        <#break>
    <#case "manualArray">
        <#assign manualArrayField = field.asManualArrayField().orElseThrow()>
        <#assign arrayElementTypeReference = manualArrayField.type.asArrayTypeReference().orElseThrow().getElementTypeReference()>

        # Manual Array Field (${manualArrayField.name})
        length_in_bits += ${helper.toParseExpression(manualArrayField, helper.intTypeReference, manualArrayField.lengthExpression, parserArguments)} * 8
        <#break>
    <#case "manual">
        <#assign manualField = field.asManualField().orElseThrow()>

        # Manual Field (${manualField.name})
        length_in_bits += ${helper.toSerializationExpression(manualField, helper.intTypeReference, manualField.lengthExpression, parserArguments)}
        <#break>
    <#case "optional">
        <#assign optionalField = field.asOptionalField().orElseThrow()>

        # Optional Field (${optionalField.name})
        <#if optionalField.type.isSimpleTypeReference()>
        if ${helper.toSerializationExpression(optionalField, helper.boolTypeReference, optionalField.conditionExpression.get(), parserArguments)}:
            <#assign simpleTypeReference = optionalField.type.asSimpleTypeReference().orElseThrow()>
            <#if simpleTypeReference.isVstringTypeReference()>
                <#assign vstringTypeReference = simpleTypeReference.asVstringTypeReference().orElseThrow()>
            length_in_bits += ${helper.toSerializationExpression(optionalField, helper.intTypeReference, vstringTypeReference.getLengthExpression(), parserArguments)}
            <#else>
            length_in_bits += ${simpleTypeReference.sizeInBits}
            </#if>
        <#elseif helper.isEnumField(field)>
            length_in_bits += ${helper.getEnumBaseTypeReference(optionalField.type).sizeInBits}
        <#elseif optionalField.type.isDataIoTypeReference()>
            length_in_bits += ${optionalField.type.asComplexTypeReference().orElseThrow().name}.length_in_bits(self.${helper.camelCaseToSnakeCase(optionalField.name)}<#if optionalField.type.asComplexTypeReference().orElseThrow().params?has_content>, <#list optionalField.type.asComplexTypeReference().orElseThrow().params.orElseThrow() as param>${helper.toParseExpression(optionalField, helper.anyTypeReference, param, parserArguments)}<#sep>, </#sep></#list></#if>)
        <#else>
            length_in_bits += self.${helper.camelCaseToSnakeCase(optionalField.name)}.length_in_bits()
        </#if>

        <#break>
    <#case "padding">
        <#assign paddingField = field.asPaddingField().orElseThrow()>
        <#assign simpleTypeReference = paddingField.type.asSimpleTypeReference().orElseThrow()>

        # Padding Field (padding)
        <#-- We're replacing the "lastItem" with 'false' here as the item itself can't know if it is the last -->
        _timesPadding: int = int(${helper.toSerializationExpression(paddingField, helper.intTypeReference, paddingField.paddingCondition, parserArguments)})
        length_in_bits = ${simpleTypeReference.sizeInBits} * _timesPadding

        <#break>
    <#case "reserved">
        <#assign reservedField = field.asReservedField().orElseThrow()>
        <#assign simpleTypeReference = reservedField.type.asSimpleTypeReference().orElseThrow()>

        # Reserved Field (reserved)
        length_in_bits += ${simpleTypeReference.sizeInBits}
        <#break>
    <#case "simple">
        <#assign simpleField = field.asSimpleField().orElseThrow()>

        # Simple field (${simpleField.name})
        <#if simpleField.type.isSimpleTypeReference()>
            <#assign simpleTypeReference = simpleField.type.asSimpleTypeReference().orElseThrow()>
            <#if simpleTypeReference.isVstringTypeReference()>
                <#assign vstringTypeReference = simpleTypeReference.asVstringTypeReference().orElseThrow()>
        length_in_bits += ${helper.toSerializationExpression(simpleField, helper.intTypeReference, vstringTypeReference.getLengthExpression(), parserArguments)}
            <#else>
        length_in_bits += ${simpleTypeReference.sizeInBits}
            </#if>
        <#elseif helper.isEnumField(field)>
        length_in_bits += ${helper.getEnumBaseTypeReference(simpleField.type).sizeInBits}
        <#elseif simpleField.type.isDataIoTypeReference()>
        length_in_bits += self.${helper.camelCaseToSnakeCase(simpleField.type.asComplexTypeReference().orElseThrow().name)}.length_in_bits(${simpleField.name}<#if simpleField.type.asComplexTypeReference().orElseThrow().params?has_content>, <#list simpleField.type.asComplexTypeReference().orElseThrow().params.orElseThrow() as param>${helper.toParseExpression(simpleField, helper.anyTypeReference, param, parserArguments)}<#sep>, </#sep></#list></#if>)
        <#else>
        length_in_bits += self.${helper.camelCaseToSnakeCase(simpleField.name)}.length_in_bits()
        </#if>
        <#break>
    <#case "switch">
        <#assign switchField = field.asSwitchField().orElseThrow()>

        # Length of subtype elements will be added by sub-type...
        <#break>
    <#case "unknown">
        <#assign unknownField = field.asUnknownField().orElseThrow()>
        <#assign simpleTypeReference = unknownField.type.asSimpleTypeReference().orElseThrow()>

        # Unknown field
        length_in_bits += ${simpleTypeReference.sizeInBits}
        <#break>
    <#case "virtual">
        <#assign virtualField = field.asVirtualField().orElseThrow()>

        # A virtual field doesn't have any in- or output.
        <#break>
</#switch>
</#list>

        return length_in_bits


<#-- The parse and serialize methods here are just proxies for forwardning the requests to static counterparts -->
    <#if !type.isDiscriminatedChildTypeDefinition()>
        <@emitImport import="from plc4py.spi.generation.ReadBuffer import ReadBuffer" />
    @staticmethod
    def static_parse(read_buffer: ReadBuffer, **kwargs):
        <#if parserArguments?has_content>
            <@emitImport import="from plc4py.api.exceptions.exceptions import PlcRuntimeException" />

        if kwargs is None:
            raise PlcRuntimeException("Wrong number of arguments, expected ${parserArguments?size}, but got None")

            <#list parserArguments as parserArgument>
                <#assign languageName=helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)>
        ${helper.camelCaseToSnakeCase(parserArgument.name)}: ${languageName} = ${helper.getNullValueForTypeReference(parserArgument.type)}
        if isinstance(kwargs.get("${helper.camelCaseToSnakeCase(parserArgument.name)}"), ${languageName}):
            ${helper.camelCaseToSnakeCase(parserArgument.name)} = ${languageName}(kwargs.get("${helper.camelCaseToSnakeCase(parserArgument.name)}"))
                <#if parserArgument.type.isSimpleTypeReference() || parserArgument.type.isEnumTypeReference()>
        elif isinstance(kwargs.get("${helper.camelCaseToSnakeCase(parserArgument.name)}"), str):
            ${helper.camelCaseToSnakeCase(parserArgument.name)} = ${languageName}(str(kwargs.get("${helper.camelCaseToSnakeCase(parserArgument.name)}")))
                </#if>
        else:
            raise PlcRuntimeException("Argument ${parserArgument?index} expected to be of type ${languageName} or a string which is parseable but was " + kwargs.get("${helper.camelCaseToSnakeCase(parserArgument.name)}").getClass().getName())

            </#list>
        </#if>
        return ${type.name}.static_parse_context(read_buffer<#if parserArguments?has_content>, <#list parserArguments as parserArgument>${helper.camelCaseToSnakeCase(parserArgument.name)}<#sep>, </#sep></#list></#if>)


    </#if>
<#-- Here come the actual parse and serialize methods that actually do the parsing and serlaizing -->
    <#assign hasParserArguments=parserArguments?has_content/>
    <#assign parserArgumentList><#if hasParserArguments><#list parserArguments as parserArgument>${helper.camelCaseToSnakeCase(parserArgument.name)}: ${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)}<#sep>, </#sep></#list></#if></#assign>
    <#assign parsedList><#if hasParserArguments><#list parserArguments as parserArgument>${helper.camelCaseToSnakeCase(parserArgument.name)}=${helper.camelCaseToSnakeCase(parserArgument.name)}<#sep>, </#sep></#list></#if></#assign>
    <@emitImport import="from plc4py.spi.generation.ReadBuffer import ReadBuffer" />
    @staticmethod
    def <#if type.isDiscriminatedChildTypeDefinition()>static_parse_builder<#else>static_parse_context</#if>(read_buffer: ReadBuffer<#if hasParserArguments>, ${parserArgumentList}</#if>):
        read_buffer.push_context("${type.name}")
        <#if helper.requiresStartPos() || helper.requiresCurPos()>
        start_pos: int = read_buffer.get_pos()
        cur_pos: int = 0
        </#if>

        <#if hasParserArguments>
            <#list parserArguments as parserArgument>
        if isinstance(${helper.camelCaseToSnakeCase(parserArgument.name)}, str):
            <#if parserArgument.type.isEnumTypeReference()>
            ${helper.camelCaseToSnakeCase(parserArgument.name)} = ${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)}[${helper.camelCaseToSnakeCase(parserArgument.name)}]
            <#elseif helper.getLanguageTypeNameForTypeReference(parserArgument.type, false) = "bool">
            <@emitImport import="from plc4py.utils.ConnectionStringHandling import strtobool" />
            ${helper.camelCaseToSnakeCase(parserArgument.name)} = bool(strtobool(${helper.camelCaseToSnakeCase(parserArgument.name)}))
            <#else>
            ${helper.camelCaseToSnakeCase(parserArgument.name)} = ${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)}(${helper.camelCaseToSnakeCase(parserArgument.name)})
            </#if>
            </#list>
        </#if>

    <#assign reservedFieldIndex=0>
    <#list type.fields as field>
        <#switch field.typeName>
            <#case "array">
                <#assign arrayField = field.asArrayField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>
                <#assign namedField = field.asNamedField().orElseThrow()>
                <#assign arrayElementTypeReference = arrayField.type.asArrayTypeReference().orElseThrow().getElementTypeReference()>

                <#if arrayElementTypeReference.isByteBased()>
                    <#if !field.isCountArrayField() && !field.isLengthArrayField()>
        raise ParseException("array fields of type byte only support 'count' and 'length' loop-types.")
                    </#if>
                    <@emitImport import="from typing import List" />
                    <@emitImport import="from typing import Any" />
        ${helper.camelCaseToSnakeCase(arrayField.name)}: List[Any] = read_buffer.read_array_field(logical_name="${arrayField.name}", read_function=<#if arrayElementTypeReference.isComplexTypeReference()><#else>read_buffer.</#if>${helper.getDataReaderCall(arrayField.type)}<#if arrayElementTypeReference.isComplexTypeReference()><#else>, </#if>count=${helper.toParseExpression(arrayField, helper.intTypeReference, arrayField.loopExpression, parserArguments)}${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
                <#else>
                <#-- If this is a count array, we can directly initialize an array with the given size -->
                    <#if field.isCountArrayField()>
                    <@emitImport import="from typing import List" />
                    <@emitImport import="from typing import Any" />
        ${helper.camelCaseToSnakeCase(arrayField.name)}: List[Any] = read_buffer.read_array_field(logical_name="${arrayField.name}", read_function=<#if arrayElementTypeReference.isComplexTypeReference()><#else>read_buffer.</#if>${helper.getDataReaderCall(arrayField.type)}<#if arrayElementTypeReference.isComplexTypeReference()><#else>, </#if>count=${helper.toParseExpression(arrayField, helper.intTypeReference, arrayField.loopExpression, parserArguments)}${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
                    <#-- In all other cases do we have to work with a list, that is later converted to an array -->
                    <#else>
                    <#-- For a length array, we read data till the read position of the buffer reaches a given position -->
                        <@emitImport import="from typing import List" />
                        <@emitImport import="from typing import Any" />
                        <#if field.isLengthArrayField()>
        ${helper.camelCaseToSnakeCase(arrayField.name)}: List[Any] = read_buffer.read_array_field(logical_name="${arrayField.name}", read_function=<#if arrayElementTypeReference.isComplexTypeReference()><#else>read_buffer.</#if>${helper.getDataReaderCall(arrayField.type)}<#if arrayElementTypeReference.isComplexTypeReference()><#else>, </#if>length=${helper.camelCaseToSnakeCase(helper.toParseExpression(arrayField, helper.intTypeReference, arrayField.loopExpression, parserArguments))}${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
                        <#-- A terminated array keeps on reading data as long as the termination expression evaluates to false -->
                        <#elseif field.isTerminatedArrayField()>
        ${helper.camelCaseToSnakeCase(arrayField.name)}: List[Any] = read_buffer.read_array_field(logical_name="${arrayField.name}", read_function=<#if arrayElementTypeReference.isComplexTypeReference()><#else>read_buffer.</#if>${helper.getDataReaderCall(arrayField.type)}<#if arrayElementTypeReference.isComplexTypeReference()><#else>, </#if>terminated=lambda : (c_bool(${helper.toParseExpression(arrayField, helper.intTypeReference, arrayField.loopExpression, parserArguments)}))${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
                        </#if>
                    </#if>
                </#if>
                <#break>
            <#case "assert">
                <#assign assertField = field.asAssertField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>
                <#assign namedField = field.asNamedField().orElseThrow()>

        ${helper.camelCaseToSnakeCase(namedField.name)}: ${helper.getLanguageTypeNameForField(field)} = read_buffer.${helper.getDataReaderCall(assertField.type)}(logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
                <#break>
            <#case "checksum">
                <#assign checksumField = field.asChecksumField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>
                <#assign namedField = field.asNamedField().orElseThrow()>
                <@emitImport import="from plc4py.protocols.${protocolName} import StaticHelper" />

        ${helper.camelCaseToSnakeCase(namedField.name)}: ${helper.getLanguageTypeNameForField(field)} = read_buffer.${helper.getDataReaderCall(checksumField.type)}(logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
                <#break>
            <#case "const">
                <#assign constField = field.asConstField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>
                <#assign namedField = field.asNamedField().orElseThrow()>

        ${helper.camelCaseToSnakeCase(namedField.name)?upper_case}: ${helper.getLanguageTypeNameForField(field)} = read_buffer.${helper.getDataReaderCall(constField.type)}(logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
                <#break>
            <#case "discriminator">
                <#assign discriminatorField = field.asDiscriminatorField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>
                <#assign namedField = field.asNamedField().orElseThrow()>
            <#if typedField.type.isSimpleTypeReference()>
                <#assign simpleTypeReference = typedField.type.asSimpleTypeReference().orElseThrow()>
        ${helper.camelCaseToSnakeCase(namedField.name)}: ${helper.getLanguageTypeNameForField(field)} = read_buffer.${helper.getDataReaderCall(discriminatorField.type)}(logical_name="${helper.camelCaseToSnakeCase(namedField.name)}", bit_length=${simpleTypeReference.sizeInBits}${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
            <#else>
        ${helper.camelCaseToSnakeCase(namedField.name)}: ${helper.getLanguageTypeNameForField(field)} = read_buffer.${helper.getDataReaderCall(discriminatorField.type)}logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
            </#if>


                <#break>
            <#case "enum">
                <#assign enumField = field.asEnumField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>
                <#assign namedField = field.asNamedField().orElseThrow()>

        ${helper.camelCaseToSnakeCase(namedField.name)} = read_buffer.${helper.getDataReaderCall(enumField.type)}logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
                <#break>
            <#case "implicit">
                <#assign implicitField = field.asImplicitField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>
                <#assign namedField = field.asNamedField().orElseThrow()>

        ${helper.camelCaseToSnakeCase(namedField.name)}: ${helper.getLanguageTypeNameForField(field)} = read_buffer.${helper.getDataReaderCall(implicitField.type)}(logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
                <#break>
            <#case "manualArray">
                <#assign manualArrayField = field.asManualArrayField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>
                <#assign namedField = field.asNamedField().orElseThrow()>
                <#assign arrayElementTypeReference = manualArrayField.type.asArrayTypeReference().orElseThrow().getElementTypeReference()>

                <#if arrayElementTypeReference.isByteBased()>
                    <@emitImport import="from typing import List" />
                    <@emitImport import="from typing import Any" />
        ${helper.camelCaseToSnakeCase(arrayField.name)}: List[Any] = read_buffer.read_array_field(logical_name="${manualArrayField.name}", read_function=<#if manualArrayField.type.isComplexTypeReference()><#else>read_buffer.</#if>${helper.getDataReaderCall(manualArrayField.type)}manual_start = lambda _values: bool(${helper.toParseExpression(manualArrayField, helper.boolTypeReference, manualArrayField.loopExpression, parserArguments)}), manual_stop = lambda : byte(${helper.toParseExpression(manualArrayField, manualArrayField.type, manualArrayField.parseExpression, parserArguments)})${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
                <#else>
                    <@emitImport import="from typing import List" />
                    <@emitImport import="from typing import Any" />
        ${helper.camelCaseToSnakeCase(arrayField.name)}: List[Any] = read_buffer.read_array_field(logical_name="${manualArrayField.name}", read_function=<#if manualArrayField.type.isComplexTypeReference()><#else>read_buffer.</#if>${helper.getDataReaderCall(manualArrayField.type)}manual_start = lambda _values: ${helper.getNonPrimitiveLanguageTypeNameForField(manualArrayField)} : bool(${helper.toParseExpression(manualArrayField, helper.boolTypeReference, manualArrayField.loopExpression, parserArguments)}), lambda : ${helper.getLanguageTypeNameForTypeReference(manualArrayField.type.elementTypeReference)}(${helper.toParseExpression(manualArrayField, manualArrayField.type, manualArrayField.parseExpression, parserArguments)})${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
                </#if>
                <#break>
            <#case "manual">
                <#assign manualField = field.asManualField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>
                <#assign namedField = field.asNamedField().orElseThrow()>

        ${helper.camelCaseToSnakeCase(manualField.name)} = read_buffer.read_manual(read_function=lambda : StaticHelper.${helper.toParseExpression(manualField, manualField.type, manualField.parseExpression, parserArguments)}, logical_name="${helper.camelCaseToSnakeCase(namedField.name)}")
                <#break>
            <#case "optional">
                <#assign optionalField = field.asOptionalField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>
                <#assign namedField = field.asNamedField().orElseThrow()>

        ${helper.camelCaseToSnakeCase(namedField.name)}: ${helper.getLanguageTypeNameForTypeReference(optionalField.type, false)} = None
        <#if optionalField.conditionExpression.present>
        if ${helper.toParseExpression(optionalField, helper.boolTypeReference, optionalField.conditionExpression.get(), parserArguments)}:
            ${helper.camelCaseToSnakeCase(namedField.name)} = <#if typedField.type.isEnumTypeReference()>read_buffer.${helper.getDataReaderCall(optionalField.type)}bit_length=${helper.getEnumBaseTypeReference(typedField.type).sizeInBits}, logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)<#else><#if optionalField.type.isComplexTypeReference()>read_buffer.${helper.getDataReaderCall(optionalField.type)}logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)<#else>read_buffer.${helper.getDataReaderCall(optionalField.type)}(logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)  </#if></#if>
        </#if>
                <#break>
            <#case "padding">
                <#assign paddingField = field.asPaddingField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>
                <#assign simpleTypeReference = paddingField.type.asSimpleTypeReference().orElseThrow()>

        for _ in range(${helper.toSerializationExpression(paddingField, paddingField.type, paddingField.paddingCondition, parserArguments)}):
            read_buffer.${helper.getDataReaderCall(paddingField.type)}(${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
                <#break>
            <#case "reserved">
                <#assign reservedField = field.asReservedField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>

        reserved_field_${reservedFieldIndex}: ${helper.getLanguageTypeNameForTypeReference(reservedField.type, false)}<#assign reservedFieldIndex=reservedFieldIndex+1> = read_buffer.${helper.getDataReaderCall(reservedField.type)}(logical_name="reserved"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
        if reserved_field_${reservedFieldIndex} != ${helper.getReservedValue(reservedField)}:
            raise ParseValidationException("reserved_field_${reservedFieldIndex} != ${helper.getReservedValue(reservedField)}")
                <#break>
            <#case "simple">
                <#assign simpleField = field.asSimpleField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>
                <#assign namedField = field.asNamedField().orElseThrow()>

                <#if typedField.type.isSimpleTypeReference()>
                    <#assign simpleTypeReference = typedField.type.asSimpleTypeReference().orElseThrow()>
                        <#if simpleTypeReference.isVstringTypeReference()>
                            <#assign vstringTypeReference = simpleTypeReference.asVstringTypeReference().orElseThrow()>
        ${helper.camelCaseToSnakeCase(namedField.name)}: ${helper.getLanguageTypeNameForTypeReference(simpleField.type, false)} = <#if typedField.type.isEnumTypeReference()>read_buffer.${helper.getDataReaderCall(simpleField.type)}bit_length=${helper.toParseExpression(simpleField, helper.intTypeReference, vstringTypeReference.getLengthExpression(), parserArguments)}, logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)<#else><#if simpleField.type.isComplexTypeReference()>read_buffer.${helper.getDataReaderCall(simpleField.type)}logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)<#else>read_buffer.${helper.getDataReaderCall(simpleField.type)}(logical_name="${helper.camelCaseToSnakeCase(namedField.name)}", bit_length=${helper.toParseExpression(simpleField, helper.intTypeReference, vstringTypeReference.getLengthExpression(), parserArguments)}${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)  </#if></#if>
                        <#else>
        ${helper.camelCaseToSnakeCase(namedField.name)}: ${helper.getLanguageTypeNameForTypeReference(simpleField.type, false)} = <#if typedField.type.isEnumTypeReference()>read_buffer.${helper.getDataReaderCall(simpleField.type)}bit_length=${simpleTypeReference.sizeInBits}, logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)<#else><#if simpleField.type.isComplexTypeReference()>read_buffer.${helper.getDataReaderCall(simpleField.type)}logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)<#else>read_buffer.${helper.getDataReaderCall(simpleField.type)}(logical_name="${helper.camelCaseToSnakeCase(namedField.name)}", bit_length=${simpleTypeReference.sizeInBits}${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)  </#if></#if>
                        </#if>
                <#else>
        ${helper.camelCaseToSnakeCase(namedField.name)}: ${helper.getLanguageTypeNameForTypeReference(simpleField.type, false)} = <#if typedField.type.isEnumTypeReference()>read_buffer.${helper.getDataReaderCall(simpleField.type)}bit_length=${helper.getEnumBaseTypeReference(typedField.type).sizeInBits}, logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)<#else><#if simpleField.type.isComplexTypeReference()>read_buffer.${helper.getDataReaderCall(simpleField.type)}logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)<#else>read_buffer.${helper.getDataReaderCall(simpleField.type)}(logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)  </#if></#if>
                </#if>


                <#break>
            <#case "switch">
                <#assign switchField = field.asSwitchField().orElseThrow()>

        # Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
        builder: ${type.name}Builder = None
                <#list switchField.cases as case>
                    <#if case.discriminatorValueTerms?has_content>
        from plc4py.protocols.${protocolName}.${outputFlavor?replace('-', '')}.${case.name} import ${case.name}
        if <@compress single_line=true>
                            <#list case.discriminatorValueTerms as discriminatorValueTerm>
                                <#if helper.isWildcard(discriminatorValueTerm)>
                                    True
                                <#else>
                                    <#assign discriminatorExpression=switchField.discriminatorExpressions[discriminatorValueTerm?index].asLiteral().orElseThrow().asVariableLiteral().orElseThrow()>
                                    <#assign discriminatorType=helper.getDiscriminatorTypes()[discriminatorExpression.discriminatorName]>
                                    ${helper.toParseExpression(switchField, discriminatorType, discriminatorExpression, parserArguments)} ==
                                    <#if discriminatorType.isEnumTypeReference()>
                                        ${helper.getLanguageTypeNameForTypeReference(discriminatorType)}.${helper.toParseExpression(switchField, discriminatorType, discriminatorValueTerm, parserArguments)}
                                    <#else>
                                        ${helper.getLanguageTypeNameForTypeReference(discriminatorType, true)}(${helper.toParseExpression(switchField, discriminatorType, discriminatorValueTerm, parserArguments)})
                                    </#if>
                                </#if>
                                <#sep> and </#sep>
                            </#list>
                            :</@compress>
                        </#if>

            builder = <@compress single_line=true>
                        <#assign hasCaseParseArguments=case.allParserArguments.isPresent() && case.allParserArguments.orElseThrow()?has_content>
                        <#assign caseParseArguments><#if hasCaseParseArguments><#list case.allParserArguments.orElseThrow() as parserArgument>${helper.camelCaseToSnakeCase(parserArgument.name)}<#sep>, </#sep></#list></#if></#assign>
                        ${case.name}.static_parse_builder(read_buffer<#if hasCaseParseArguments>, ${tracer.dive("case parse arguments")} ${caseParseArguments}</#if>)
                    </@compress>

        </#list>
        <@emitImport import="from plc4py.api.exceptions.exceptions import ParseException" />
        if builder is None:
            raise ParseException("Unsupported case for discriminated type"<#if switchField.getDiscriminatorExpressions()?has_content>+" parameters ["<#list switchField.getDiscriminatorExpressions() as discriminatorExpression>+"${discriminatorExpression.stringRepresentation()}="+str(${helper.toParseExpression(null, null, discriminatorExpression, parserArguments)})<#sep>+" "</#sep></#list>+"]"</#if>)

                <#break>
            <#case "unknown">
                <#assign unknownField = field.asUnknownField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>

        read_buffer.${helper.getDataReaderCall(unknownField.type)}(logical_name="unknown"${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
                <#break>
            <#case "virtual">
                <#assign virtualField = field.asVirtualField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>
                <#assign namedField = field.asNamedField().orElseThrow()>
        ${helper.camelCaseToSnakeCase(namedField.name)}: ${helper.getLanguageTypeNameForField(field)} = ${helper.toParseExpression(virtualField, virtualField.type, virtualField.valueExpression, parserArguments)}
                <#break>
            <#case "validation">
                <#assign validationField = field.asValidationField().orElseThrow()>

        # Validation
        if not ${helper.toParseExpression(validationField, helper.boolTypeReference, validationField.getValidationExpression(), null)}:
            <#assign errorType="ParseValidationException">
            <#if !validationField.shouldFail()><#assign errorType="ParseAssertException"></#if>
            raise ${errorType}(${validationField.getDescription().orElse("\"Validation failed\"")})

                <#break>
            <#case "peek">
                <#assign peekField = field.asPeekField().orElseThrow()>
                <#assign typedField = field.asTypedField().orElseThrow()>
                <#assign namedField = field.asNamedField().orElseThrow()>

        ${helper.camelCaseToSnakeCase(namedField.name)}: ${helper.getLanguageTypeNameForField(field)} = read_buffer.${helper.getDataReaderCall(peekField.type)}(logical_name="${helper.camelCaseToSnakeCase(namedField.name)}"<#if peekField.offsetExpression.present>, ${helper.toParseExpression(peekField, helper.boolTypeReference, peekField.offsetExpression.get(), parserArguments)}</#if>${helper.getFieldOptions(typedField, parserArguments)}<#if hasParserArguments>, ${parsedList}</#if>)
                <#break>
        </#switch>
    </#list>

        read_buffer.pop_context("${type.name}")
        # Create the instance
    <#if type.isDiscriminatedChildTypeDefinition()>
        return ${type.name}Builder(<@compress single_line=true>
        <#list type.propertyFields as field>
            ${helper.camelCaseToSnakeCase(field.name)}<#sep>, </#sep>
        </#list>
        <#list reservedFields as reservedField>
            reserved_field_${reservedField?index}<#sep>, </#sep>
        </#list>
        )</@compress>
    <#elseif type.isDiscriminatedParentTypeDefinition()>
        _${helper.camelCaseToSnakeCase(type.name)}: ${type.name} = builder.build(<@compress single_line=true>
        <#list type.propertyFields as field>
            ${helper.camelCaseToSnakeCase(field.name)}<#sep>, </#sep>
        </#list>
        <#if filteredParserArguments?has_content>
            <#if type.propertyFields?has_content>, </#if>
            <#list filteredParserArguments as arg>
                ${helper.camelCaseToSnakeCase(arg.name)}<#sep>, </#sep>
            </#list>
        </#if>
        )</@compress>
        <#list reservedFields as reservedField>
        _${helper.camelCaseToSnakeCase(type.name)}.reserved_field${reservedField?index} = reserved_field_${reservedField?index}
        </#list>
        return _${helper.camelCaseToSnakeCase(type.name)}
    <#else>
        _${helper.camelCaseToSnakeCase(type.name)}: ${type.name} = ${type.name}(<@compress single_line=true>
        <#list type.propertyFields as field>
            ${helper.camelCaseToSnakeCase(field.name)}<#sep>, </#sep>
        </#list>
        <#if filteredParserArguments?has_content>
            <#if type.propertyFields?has_content>, </#if>
            <#list filteredParserArguments as arg>
                ${helper.camelCaseToSnakeCase(arg.name)}<#sep>, </#sep>
            </#list>
        </#if>
        )</@compress>
        <#list reservedFields as reservedField>
        _${helper.camelCaseToSnakeCase(type.name)}.reserved_field_${reservedField?index} = reserved_field_${reservedField?index}
        </#list>
        return _${helper.camelCaseToSnakeCase(type.name)}
    </#if>


    def equals(self, o: object) -> bool:
        if self == o:
            return True

        if not isinstance(o, ${type.name}):
            return False

        that: ${type.name} = ${type.name}(o)
        return <@compress single_line=true>
            <#if type.propertyFields?has_content>
            <#list type.propertyFields as field>
            (self.${helper.camelCaseToSnakeCase(field.name)} == that.${helper.camelCaseToSnakeCase(field.name)}) and
            </#list>
            </#if>
            <#if type.parentType.isPresent()>
            super().equals(that) and
            </#if>
            True</@compress>

    def hash_code(self) -> int:
        return hash(self)

    <@emitImport import="from plc4py.api.exceptions.exceptions import SerializationException" />
    <@emitImport import="from plc4py.api.exceptions.exceptions import PlcRuntimeException" />
    def __str__(self) -> str:
        # TODO:- Implement a generic python object to probably json convertor or something.
        return ""


    <#if type.isDiscriminatedParentTypeDefinition()>
@dataclass
class ${type.name}Builder:
    <#list type.propertyFields as field>
    ${helper.camelCaseToSnakeCase(field.name)}: ${helper.getLanguageTypeNameForField(field)}
    </#list>
    <@compress single_line=true>def build(self,
<#if filteredParserArguments?has_content>
    <#list filteredParserArguments as arg>
        ${helper.camelCaseToSnakeCase(arg.name)}: ${helper.getLanguageTypeNameForTypeReference(arg.type)}<#sep>, </#sep>
    </#list>
</#if>
) -> ${type.name}:</@compress>
        pass

    </#if>

    <#if type.isDiscriminatedChildTypeDefinition()>
@dataclass
class ${type.name}Builder:
<#if type.propertyFields?has_content>
    <#list type.propertyFields as field>
    ${helper.camelCaseToSnakeCase(field.name)}: ${helper.getLanguageTypeNameForField(field)}
    </#list>
</#if>
<#list reservedFields as reservedField>
    reservedField${reservedField?index}: ${helper.getLanguageTypeNameForTypeReference(reservedField.type, false)}
</#list>

    def build(self,<@compress single_line=true>
    <#if filteredParserArguments?has_content>
        <#list filteredParserArguments as arg>
            ${helper.camelCaseToSnakeCase(arg.name)}: ${helper.getLanguageTypeNameForTypeReference(arg.type)}<#sep>, </#sep>
        </#list>
        ,
    </#if>
    <#if type.getParentPropertyFields()?has_content>
        <#list type.getParentPropertyFields() as field>
            ${helper.camelCaseToSnakeCase(field.name)}<#sep>, </#sep>
        </#list>
    </#if>
    <#if filteredParentParserArguments?has_content>
        <#if type.parentType.orElseThrow().asComplexTypeDefinition().orElseThrow().propertyFields?has_content>, </#if>
        <#list filteredParentParserArguments as arg>
            ${helper.camelCaseToSnakeCase(arg.name)}: ${helper.getLanguageTypeNameForTypeReference(arg.type)}<#sep>, </#sep>
        </#list>
    </#if>
    ) -> ${type.name}:</@compress>
        ${helper.camelCaseToSnakeCase(type.name)}: ${type.name} = ${type.name}(<@compress single_line=true>
    <#if filteredParserArguments?has_content>
        <#list filteredParserArguments as arg>
            ${helper.camelCaseToSnakeCase(arg.name)}<#sep>, </#sep>
        </#list>
    </#if>
    <#if filteredParserArguments?has_content>, </#if>
    <#list type.getParentPropertyFields() as field>
        ${helper.camelCaseToSnakeCase(field.name)}<#sep>, </#sep>
    </#list>
    <#if type.getParentPropertyFields()?has_content>, </#if>
    <#list type.propertyFields as field>
        self.${helper.camelCaseToSnakeCase(field.name)}<#sep>, </#sep>
    </#list>
    )</@compress>
    <#list reservedFields as reservedField>
        ${helper.camelCaseToSnakeCase(type.name)}.reservedField${reservedField?index} = reservedField${reservedField?index}
    </#list>
        return ${helper.camelCaseToSnakeCase(type.name)}

    </#if>
</@importSectionWithContentBelow>

</#outputformat>