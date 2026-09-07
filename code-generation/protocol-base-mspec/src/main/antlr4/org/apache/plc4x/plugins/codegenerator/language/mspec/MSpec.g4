grammar MSpec;
/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

file
 : contantsDefinition? globalsDefinition? contextDefintion? complexTypeDefinition* EOF
 ;

contantsDefinition
 : LBRACKET CONSTANTS (LBRACKET constField RBRACKET)* RBRACKET
 ;

globalsDefinition
 : LBRACKET GLOBALS fieldDefinition* RBRACKET
 ;

contextDefintion
 : LBRACKET CONTEXT fieldDefinition* RBRACKET
 ;

complexTypeDefinition
 : LBRACKET complexType RBRACKET
 ;

complexType
 : TYPE name=idExpression (LRBRACKET params=argumentList RRBRACKET)? attributes=attributeList (fieldDefinition|batchSetDefinition)*
 | DISCRIMINATEDTYPE name=idExpression (LRBRACKET params=argumentList RRBRACKET)? attributes=attributeList (fieldDefinition|batchSetDefinition)*
 | ENUM (type=dataType)? name=idExpression (LRBRACKET params=argumentList RRBRACKET)? attributes=attributeList enumValues=enumValueDefinition*
 | DATAIO name=idExpression (LRBRACKET params=argumentList RRBRACKET)? (attributes=attributeList) dataIoTypeSwitch=dataIoDefinition
 ;

fieldDefinition
 : LBRACKET field (attributes=attributeList) RBRACKET
 ;

batchSetDefinition
 : LBRACKET BATCHSET attributes=attributeList fieldDefinition+ RBRACKET
 ;

dataIoDefinition
// TODO: possibly alow more fields than just one typeSwitch field.
 : LBRACKET typeSwitchField (LBRACKET params=multipleExpressions RBRACKET)? RBRACKET
 ;

field
 : abstractField
 | arrayField
 | assertField
 | checksumField
 | constField
 | discriminatorField
 | enumField
 | implicitField
 | manualArrayField
 | manualField
 | optionalField
 | paddingField
 | peekField
 | reservedField
 | simpleField
 | stateField
 | typeSwitchField
 | unknownField
 | validationField
 | virtualField
 ;

abstractField
 : ABSTRACT type=typeReference name=idExpression
 ;

arrayField
 : ARRAY type=typeReference name=idExpression loopType=ARRAY_LOOP_TYPE loopExpression=expression
 ;

assertField
 : ASSERT type=typeReference name=idExpression condition=expression
 ;

checksumField
 : CHECKSUM type=dataType name=idExpression checksumExpression=expression
 ;

constField
 : CONST type=typeReference name=idExpression expected=valueLiteral
 ;

discriminatorField
 : DISCRIMINATOR type=typeReference name=idExpression
 ;

enumField
 : ENUM type=typeReference name=idExpression fieldName=idExpression
 ;

implicitField
 : IMPLICIT type=dataType name=idExpression serializeExpression=expression
 ;

manualArrayField
 : MANUALARRAY type=typeReference name=idExpression loopType=ARRAY_LOOP_TYPE loopExpression=expression parseExpression=expression serializeExpression=expression lengthExpression=expression
 ;

manualField
 : MANUAL type=typeReference name=idExpression parseExpression=expression serializeExpression=expression lengthExpression=expression
 ;

optionalField
 : OPTIONAL type=typeReference name=idExpression (condition=expression)?
 ;

paddingField
 : PADDING type=dataType name=idExpression paddingValue=expression timesPadding=expression
 ;

peekField
 : PEEK type=typeReference name=idExpression (offset=expression)?
 ;

reservedField
 : RESERVED type=dataType expected=expression
 ;

simpleField
 : SIMPLE type=typeReference name=idExpression
 ;

stateField
 : STATE name=idExpression
 ;

typeSwitchField
 : TYPESWITCH discriminators=multipleVariableLiterals caseStatement*
 ;

unknownField
 : UNKNOWN type=dataType
 ;

validationField
 : VALIDATION validationExpression=expression (description=STRING_LITERAL)? (SHOULD_FAIL '=' shouldFail=BOOLEAN_LITERAL)?
 ;

virtualField
 : VIRTUAL type=typeReference name=idExpression valueExpression=expression
 ;

enumValueDefinition
 : LBRACKET (valueExpression=expression)? name=IDENTIFIER_LITERAL (LBRACKET constantValueExpressions=multipleExpressions RBRACKET)? RBRACKET
 ;

typeReference
 : complexTypeReference=IDENTIFIER_LITERAL (LRBRACKET params=multipleExpressions RRBRACKET)?
 | simpleTypeReference=dataType
 ;

caseStatement
 : LBRACKET (discriminatorValues=multipleExpressions)? (nameWildcard=ASTERISK)? name=IDENTIFIER_LITERAL (LRBRACKET params=argumentList RRBRACKET)? (fieldDefinition|batchSetDefinition)* RBRACKET
 ;

dataType
 : base=BIT
 | base=BYTE
 | base=INT size=INTEGER_LITERAL
 | base=VINT
 | base=UINT size=INTEGER_LITERAL
 | base=VUINT
 | base=FLOAT size=INTEGER_LITERAL
 | base=UFLOAT size=INTEGER_LITERAL
 | base=STRING size=INTEGER_LITERAL
 | base=VSTRING (length=expression)?
 | base=TIME
 | base=DATE
 | base=DATETIME
 ;

attribute
 : name=IDENTIFIER_LITERAL '=' value=expression
 ;

attributeList
 : attribute*
 ;

argument
 : type=typeReference name=idExpression
 ;

argumentList
 : argument (',' argument)*
 ;

expression
 : TICK expr=innerExpression TICK
 // TODO: check if this is really universal or should be specific to case statement
 | ASTERISK
 ;

multipleExpressions
 : expression (',' expression)*
 ;

multipleVariableLiterals
 : variableLiteral (',' variableLiteral)*
 ;

variableLiteral
 : IDENTIFIER_LITERAL
 | IDENTIFIER_LITERAL '.' variableLiteral // Field Reference or method call
 | variableLiteral '[' + INTEGER_LITERAL + ']' // Array index
 ;

innerExpression
 : valueLiteral
 // Explicitly allow the loop type keywords in expressions
 | ARRAY_LOOP_TYPE
 | idExpression ('(' (innerExpression (',' innerExpression)* )? ')')? ('[' innerExpression ']')?
 | innerExpression '.' innerExpression // Field Reference or method call
 | innerExpression '[' + INTEGER_LITERAL + ']' // Array index
 | innerExpression binaryOperator innerExpression  // Addition
 | innerExpression '?' innerExpression ':' innerExpression
 | '(' innerExpression ')'
 | '"' innerExpression '"'
 | '!' innerExpression
 ;

valueLiteral
 : BOOLEAN_LITERAL
 | HEX_LITERAL
 | INTEGER_LITERAL
 | FLOAT_LITERAL
 | STRING_LITERAL
 ;

idExpression
 : id=IDENTIFIER_LITERAL
 // Explicitly allow keywords in id-expressions
 | id=ARRAY_LOOP_TYPE
 | id=CONSTANTS | id=GLOBALS | id=CONTEXT | id=TYPE | id=DISCRIMINATEDTYPE | id=DATAIO
 | id=ENUM | id=BATCHSET | id=ABSTRACT | id=ARRAY | id=ASSERT | id=CHECKSUM | id=CONST
 | id=DISCRIMINATOR | id=IMPLICIT | id=MANUALARRAY | id=MANUAL | id=OPTIONAL | id=PADDING
 | id=PEEK | id=RESERVED | id=SIMPLE | id=STATE | id=TYPESWITCH | id=UNKNOWN | id=VALIDATION
 | id=VIRTUAL
 | id=BIT | id=BYTE | id=INT | id=VINT | id=UINT | id=VUINT
 | id=FLOAT | id=UFLOAT | id=STRING | id=VSTRING | id=TIME | id=DATE | id=DATETIME
 | id=SHOULD_FAIL
 ;

binaryOperator
 : '+'
 | '-'
 | '/'
 | ASTERISK
 | '^'
 | '=='
 | '!='
 | '>>'
 | '<<'
 | '>'
 | '<'
 | '>='
 | '<='
 | '&&'
 | '||'
 | '&'
 | '|'
 | '%'
 ;

TICK : '\'';
LBRACKET : '[';
RBRACKET : ']';
LRBRACKET : '(';
RRBRACKET : ')';
LCBRACKET : '{';
RCBRACKET : '}';

ASTERISK : '*';

// Keywords used for higher level constructs
CONSTANTS       : 'constants';
GLOBALS         : 'globals';
CONTEXT         : 'context';
TYPE            : 'type';
DISCRIMINATEDTYPE : 'discriminatedType';
DATAIO          : 'dataIo';
ENUM            : 'enum';
BATCHSET        : 'batchSet';

// Keywords used for fields
ABSTRACT        : 'abstract';
ARRAY           : 'array';
ASSERT          : 'assert';
CHECKSUM        : 'checksum';
CONST           : 'const';
DISCRIMINATOR   : 'discriminator';
IMPLICIT        : 'implicit';
MANUALARRAY     : 'manualArray';
MANUAL          : 'manual';
OPTIONAL        : 'optional';
PADDING         : 'padding';
PEEK            : 'peek';
RESERVED        : 'reserved';
SIMPLE          : 'simple';
STATE           : 'state';
TYPESWITCH      : 'typeSwitch';
UNKNOWN         : 'unknown';
VALIDATION      : 'validation';
VIRTUAL         : 'virtual';

// Keywords for types
BIT             : 'bit';
BYTE            : 'byte';
INT             : 'int';
VINT            : 'vint';
UINT            : 'uint';
VUINT           : 'vuint';
FLOAT           : 'float';
UFLOAT          : 'ufloat';
STRING          : 'string';
VSTRING         : 'vstring';
TIME            : 'time';
DATE            : 'date';
DATETIME        : 'dateTime';

// Keywords used elsewhere
SHOULD_FAIL     : 'shouldFail';

ARRAY_LOOP_TYPE
 : 'count'
 | 'length'
 | 'terminated'
 ;

// Integer literals

INTEGER_LITERAL
 : INTEGER_CHARACTERS
 ;

fragment
INTEGER_CHARACTERS
 : INTEGER_CHARACTER+
 ;

fragment
INTEGER_CHARACTER
 : [0-9]
 ;

FLOAT_LITERAL
 : INTEGER_LITERAL.INTEGER_LITERAL
 ;

// Hexadecimal literals

HEX_LITERAL
 : '0' [xX] HEX_CHARACTERS
 ;

fragment
HEX_CHARACTERS
 : HEX_CHARACTER+
 ;

fragment
HEX_CHARACTER
 : [0-9a-fA-F]
 ;

// Boolean literals

BOOLEAN_LITERAL
 : 'true'
 | 'false'
 ;

// String literals

STRING_LITERAL
 : '"' STRING_CHARACTERS? '"'
 ;

// As we're generating property names and class names from these,
// we have to put more restrictions on them.

IDENTIFIER_LITERAL
 : [A-Za-z0-9_-]+
 ;

fragment
STRING_CHARACTERS
 : STRING_CHARACTER+
 ;

fragment
STRING_CHARACTER
 : ~["\\\r\n]
 ;

LINE_COMMENT
 : (('//' ~[\r\n]*)+) -> channel(HIDDEN)
 ;

BLOCK_COMMENT
 : '/*' .*? '*/' -> channel(HIDDEN)
 ;

EmptyLine
 : {getCharPositionInLine() == 0}? [ \t]* '\r'? '\n' -> channel(HIDDEN)
 ;

NEWLINE
 : '\r'? '\n' -> channel(HIDDEN)
 ;

WS
 : [ \t\u000C]+ -> channel(HIDDEN)
 ;
