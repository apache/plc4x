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

#ifndef _ERROR_INFO_EXCEPTION_H_
#define _ERROR_INFO_EXCEPTION_H_

#include <exception>
#include <boost/exception/exception.hpp>
#include "ErrorCategory.h"

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

                    class ErrorInfoException : public boost::exception, public std::exception
                    {
                    public:
                        ErrorInfoException();

                        virtual ~ErrorInfoException();

                        ErrorInfoException(unsigned int uiUserMessageID, std::string strErrorInfo);

                        ErrorInfoException(unsigned int uiUserMessageID, std::exception& ex);

                        ErrorInfoException(unsigned int uiUserMessageID, boost::exception& ex);

                        ErrorInfoException(unsigned int uiUserMessageID, Error_Category::Category cat, int iErrorValue);

                        ErrorInfoException(unsigned int uiUserMessageID, Error_Category::Category cat, int iErrorValue, std::string strErrorInfo);

                        ErrorInfoException(unsigned int uiUserMessageID, Error_Category::Category cat, int iErrorValue, std::exception& ex);

                        ErrorInfoException(unsigned int uiUserMessageID, Error_Category::Category cat, int iErrorValue, boost::exception& ex);

                        inline unsigned int getUserMessageID() const { return _uiUserMessageID; }

                        virtual const char* what() const throw() { return _strWhat.c_str(); }

                        Error_Category* pCategory;

                        inline int getErrorValue() const { return _iEv; }                        

                    protected:

                        Error_Category* CreateErrorCategory(Error_Category::Category cat) throw();

                    protected:

                        std::string _strWhat;
                        int _iEv;
                        unsigned int _uiUserMessageID;

                    };
                }
            }
        }
    }
}
#endif