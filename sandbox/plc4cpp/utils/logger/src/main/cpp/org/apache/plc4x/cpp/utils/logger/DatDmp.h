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

#ifndef _DATDMP_H_
#define _DATDMP_H_

#include <stdlib.h>
#include <string>
#include <iostream>

using namespace std;

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
					class DatDmp
					{
					public:

						DatDmp();
						DatDmp(unsigned int uiBytesPerLine);

						string Show(unsigned char* pStream, unsigned int uiNumbersBytes, const char* pLinePrefix = NULL);

						string getPrefix() { return _strPrefix; }
						void setPrefix(string strPrefix) { _strPrefix = strPrefix; }

						void setPrefix(const char* pPrefix) {   _strPrefix.clear();
                                                                if (pPrefix != NULL) { _strPrefix = pPrefix; };
						                                    }

						inline unsigned getBytesPerLine() const { return _uiBytesPerLine; }
						void setBytesPerLine(unsigned int uiBytesPerLine) { _uiBytesPerLine = uiBytesPerLine; }

						inline bool getLeadingNewLine() const { return _bLeadingNewLine; }
						void setLeadingNewLine(bool bLeadingNewLine) { _bLeadingNewLine = bLeadingNewLine; }

                        static string HexByte(unsigned char ucByte);
                        static string HexWord(unsigned short usWord);
                        static string HexDWord(unsigned int uiDword);

					protected:

						static const char* Hex(unsigned char ucByte);
                        static char HexNibbleHigh(unsigned char ucByte);
                        static char HexNibbleLow(unsigned char ucByte);

					protected:

						string          _strHexDump;
						string          _strPrefix;
						unsigned int    _uiBytesPerLine;
						bool            _bLeadingNewLine;

					private:

						static char     achHex[3];
					};
				}
			}
		}
	}
}

#endif
