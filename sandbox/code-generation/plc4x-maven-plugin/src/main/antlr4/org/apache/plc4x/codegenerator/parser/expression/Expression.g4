grammar Expression;
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

expressionString
	: a=expression EOF
	;

expression
	: boolExpr
	;

boolExpr
	: left=sumExpr (op=AND|OR|LT|LTEQ|GT|GTEQ|EQ|NOTEQ right=boolExpr)?
	;

sumExpr
	: left=productExpr (op=SUB|ADD right=sumExpr)?
	;

productExpr
	: left=expExpr (op=DIV|MULT right=productExpr)?
	;

expExpr
	: base=unaryOperation (op=EXP exponent=expExpr)?
	;

unaryOperation
	: a=operand
	| NOT e=expression
	| SUB e=expression
	| LPAREN e=expression RPAREN
	;

operand
	: l=literal
	| v=VARIABLE
	| f=functionExpr
	;

functionExpr
	: f=FUNCNAME LPAREN (a=arguments)? RPAREN
	;

arguments
	: a=expression (COMMA b=arguments)?
	;

literal
	: numeric=NUMBER
	| string=STRING
	| boolan=TRUE|FALSE
	;

STRING
	:
	'"' STRING_EXPRESSION*	'"'
	|
	'\'' STRING_EXPRESSION* '\''
	;

STRING_EXPRESSION
	: ESCAPE_SEQUENCE
	| ~'\\'
	;

TRUE
	: ('t'|'T')('r'|'R')('u'|'U')('e'|'E')
	;

FALSE
	: ('f'|'F')('a'|'A')('l'|'L')('s'|'S')('e'|'E')
	;


NOTEQ   : '!=';
LTEQ    : '<=';
GTEQ    : '>=';
AND		: '&&';
OR      : '||';
NOT	    : '!';
EQ      : '==';
LT      : '<';
GT      : '>';

EXP     : '^';
MULT    : '*';
DIV     : '/';
ADD     : '+';
SUB     : '-';

LPAREN  : '(';
RPAREN  : ')';
COMMA   : ',';

VARIABLE
	: '[' ~('[' | ']')+ ']'
	;
FUNCNAME
	: (LETTER)+
	;
NUMBER
	: (DIGIT)+ ('.' (DIGIT)+)?
	;
WHITESPACE
	: (' ' | '\n' | '\t' | '\r')+ -> skip
	;

fragment
LETTER
	: ('a'..'z') | ('A'..'Z')
	;

fragment
DIGIT
	: ('0'..'9')
	;

fragment
ESCAPE_SEQUENCE
	: '\\' 't'
	| '\\' 'n'
	| '\\' '"'
	| '\\' '\''
	| '\\' '\\'
	;
