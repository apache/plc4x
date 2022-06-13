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

#ifndef _ERRORCATEGORY_H_
#define _ERRORCATEGORY_H_

#include <boost/system/error_code.hpp>
#include <string>

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
                    class IError_Category
                    {
                    public:
                 
                        virtual const char* name() const = 0;
                        virtual inline std::string getName() const = 0;

                        virtual std::string message(int iEv) const = 0;

                    protected:
                        virtual bool IsValidMessageIndex(int iEv) const = 0;
                    };


                    class Error_Category : IError_Category
                    {
                    public:
                        enum Category
                        {
                            System = 0,
                            Application,
                            max_category = Application,
                            num_categories
                        };

                        Error_Category() :
                            _Category(System),
                            _strName("System"),
                            _iMsgOffset(0)				
                        {								
                        }								

                        virtual const char* name() const { return _strName.c_str(); }

                        virtual std::string message(int ev) const;

                        virtual inline std::string getName() const { return _strName; }

                    protected:

                        std::string DefaultMessage(int ev) const;

                        virtual bool IsValidMessageIndex(int ev) const { return false; }

                    protected:
                        std::string _strName;
                        Category _Category;
                        int _iMsgOffset;

                        static const char* _paErrorMessage[];
                    };

                    class Error_Category_Application : public Error_Category
                    {
                    public:

                        enum errc_t						
                        {
                            success = 0,				
                            err1,
                            err2,
                            err3,
                            err4,
                            err5,
                            err6,
                            max_err = err6,	
                            num_codes
                        };

                        Error_Category_Application()
                        {
                            _Category = Application;
                            _strName = "Application";
                            _iMsgOffset = 0;			// Start-Offset to error message index, to be set ..
                        }								// .. individually for each derivation of Error_Category!
                                                        // Provide messages in array Error_Category::_aszErrorMessage[]!
                    protected:

                        virtual bool IsValidMessageIndex(int iEv) const
                        {
                            return iEv >= (int)success && iEv <= _iMsgOffset + (int)max_err;
                        }

                    };
                }
            }
        }
    }
}
#endif