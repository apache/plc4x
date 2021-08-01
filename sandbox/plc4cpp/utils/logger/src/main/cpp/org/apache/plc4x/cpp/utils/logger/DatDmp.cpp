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

#include "DatDmp.h"

namespace org
{
    namespace apache
    {
        namespace plc4x
        {
            namespace cpp
            {
                namespace utils
                {

                    char DatDmp::achHex[3] = { '\0', '\0', '\0' };

                    DatDmp::DatDmp() : 
                        _strHexDump(""),
                        _strPrefix(""),
                        _uiBytesPerLine(16),
                        _bLeadingNewLine(false)
                    {
                    }

                    DatDmp::DatDmp(unsigned int uiBytesPerLine) : 
                        _strHexDump(""),
                        _strPrefix(""),
                        _uiBytesPerLine(uiBytesPerLine),
                        _bLeadingNewLine(false)
                    {
                    }

                    string DatDmp::Show(unsigned char* pStream, unsigned int uiNumbersBytes, const char* pLinePrefix /*= NULL*/)
                    {
                        _strHexDump = "";
                        if (_bLeadingNewLine)
                        {
                            _bLeadingNewLine = false;
                            _strHexDump += '\n';
                        }

                        for (unsigned int uiByteCounter = 0, iStopByte = _uiBytesPerLine;
                            uiByteCounter < uiNumbersBytes;
                            uiByteCounter += _uiBytesPerLine, iStopByte += _uiBytesPerLine
                            )
                        {
                            if (pLinePrefix != NULL) _strHexDump += pLinePrefix;
                            _strHexDump += _strPrefix;

                            int nBoundaryBytes = 0;
                            if (iStopByte > uiNumbersBytes)
                            {
                                nBoundaryBytes = iStopByte - uiNumbersBytes;
                                iStopByte = uiNumbersBytes;
                            }

                            for (unsigned int uiByteLine = uiByteCounter; uiByteLine < iStopByte; uiByteLine++)
                            {
                                unsigned char byte = pStream[uiByteLine];
                                _strHexDump += Hex(byte);
                                _strHexDump += ' ';
                            }
                            if (nBoundaryBytes > 0) _strHexDump += std::string(3 * nBoundaryBytes, ' ');
                            _strHexDump += "| ";

                            for (unsigned int uiByteLine = uiByteCounter; uiByteLine < iStopByte; uiByteLine++)
                            {
                                char ch = (char)pStream[uiByteCounter];
                                if (ch >= 0)
                                {
                                    if (!isprint(ch)) ch = '.';
                                }
                                else
                                {
                                    ch = '.';
                                }
                                _strHexDump += ch;
                            }
                            _strHexDump += "\n";
                        }
                        return _strHexDump;
                    }

                    string DatDmp::HexByte(unsigned char ucByte) 
                    {
                        return string(Hex(ucByte));
                    }
                    
                    string DatDmp::HexWord(unsigned short usWord) 
                    {
                        return string(Hex((unsigned char)(usWord >> 8))) + string(Hex((unsigned char)(usWord & 0xff)));
                    }

                    string DatDmp::HexDWord(unsigned int uiDWord)
                    {
                        return	string(Hex((unsigned char)((uiDWord & 0xff000000) >> 24)))
                            + string(Hex((unsigned char)((uiDWord & 0x00ff0000) >> 16)))
                            + string(Hex((unsigned char)((uiDWord & 0x0000ff00) >> 8)))
                            + string(Hex((unsigned char)((uiDWord & 0x000000ff))));
                    }

                    const char* DatDmp::Hex(unsigned char ucByte)
                    {
                        achHex[0] = HexNibbleHigh(ucByte);
                        achHex[1] = HexNibbleLow(ucByte);
                        achHex[2] = '\0';
                        return (const char*)achHex;
                    }

                    char DatDmp::HexNibbleHigh(unsigned char ucByte)
                    {
                        unsigned char btNibble = (unsigned char)(ucByte >> 4);
                        return btNibble <= (unsigned char)9 ? '0' + btNibble : 'A' + btNibble - (unsigned char)10;
                    }

                    char DatDmp::HexNibbleLow(unsigned char ucByte)
                    {
                        unsigned char btNibble = (unsigned char)(ucByte & 0x0f);
                        return btNibble <= (unsigned char)9 ? '0' + btNibble : 'A' + btNibble - (unsigned char)10;
                    }
                }
            }
        }
    }
}
