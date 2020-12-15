grammar Expression;
/*

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <http://unlicense.org/>

Borrowed part of this Grammar:
https://github.com/bkiers/tiny-language-antlr4/blob/master/src/main/antlr4/tl/antlr4/TL.g4

*/

expressionString
 : expression EOF
 ;

expression
 : '-' expression                                       #unaryMinusExpression
 | '!' expression                                       #notExpression
 | <assoc=right> expression '^' expression              #powerExpression
 | expression op=( '*' | '/' | '%' ) expression         #multExpression
 | expression op=( '+' | '-' ) expression               #addExpression
 | expression op=( '>>' | '<<' ) expression             #bitShiftExpression
 | expression op=( '>=' | '<=' | '>' | '<' ) expression #compExpression
 | expression op=( '==' | '!=' ) expression             #eqExpression
 | expression '&&' expression                           #andExpression
 | expression '&' expression                            #bitAndExpression
 | expression '||' expression                           #orExpression
 | expression '|' expression                            #bitOrExpression
 | expression '?' expression ':' expression             #ifExpression
 | Number                                               #numberExpression
 | Bool                                                 #boolExpression
 | Null                                                 #nullExpression
 | identifierSegment                                    #identifierExpression
 | String indexes?                                      #stringExpression
 | '(' expression ')' indexes?                          #expressionExpression
 ;

identifierSegment
 : name=Identifier args=identifierSegmentArguments? index=identifierSegmentIndexes? ('.' rest=identifierSegmentRest)?
 ;

identifierSegmentArguments
 : arguments
 ;

identifierSegmentIndexes
 : indexes
 ;

identifierSegmentRest
 : identifierSegment
 ;

arguments
 : '(' (expression (',' expression)*)? ')'
 ;

indexes
 : ( '[' expression ']' )+
 ;

Null     : 'null';

Bool
 : 'true'
 | 'false'
 ;

Number
 : Int ( '.' Digit* )?
 ;

Identifier
 : [a-zA-Z_] [a-zA-Z_0-9]*
 ;

String
 : ["] ( ~["\r\n\\] | '\\' ~[\r\n] )* ["]
 | ['] ( ~['\r\n\\] | '\\' ~[\r\n] )* [']
 ;

Space
 : [ \t\r\n\u000C] -> skip
 ;

fragment Int
 : [1-9] Digit*
 | '0'
 ;

fragment Digit
 : [0-9]
 ;