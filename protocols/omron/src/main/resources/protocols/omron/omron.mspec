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

[type FinsTCPHeader
    [const uint 32 magicNumber 0x46494E53]      // 固定为 "FINS" 的 ASCII 值
    [simple uint 16 length]                     // 数据帧长度（包括所有字段）
    [simple uint 8 reserved]                    // 保留字段，通常为 0
    [simple uint 8 flags]                       // 标志位，通常为 0
]

[type FinsCommandHeader
    [simple uint 8 ICF]                   // 信息控制字段
    [simple uint 8 RSV]                   // 保留字段
    [simple uint 8 GCT]                   // 网关计数
    [simple uint 8 DNA]                   // 目标网络地址
    [simple uint 8 DA1]                   // 目标节点地址
    [simple uint 8 DA2]                   // 目标单元地址
    [simple uint 8 SNA]                   // 源网络地址
    [simple uint 8 SA1]                   // 源节点地址
    [simple uint 8 SA2]                   // 源单元地址
    [simple uint 8 SID]                   // 服务 ID
]

[type FinsFrame
    [field FinsTCPHeader tcpHeader]         // TCP 头部
    [field FinsCommandHeader cmdHeader]     // FINS 命令头部
    [array uint 8 data 'length - 16'        // 数据段，长度由总长度减去头部长度计算
        'length' 'tcpHeader.length - 16']
]

// 定义读取内存区命令
[type FinsReadMemoryCommand
    [const uint 16 mainFunction 0x0101]   // 主功能码：读取内存区
    [simple uint 16 memoryAreaCode]       // 内存区代码（如 CIO, DM）
    [simple uint 16 address]              // 起始地址
    [simple uint 8 bitOffset]             // 位偏移
    [simple uint 16 readLength]           // 读取长度
]

// 定义读取内存区的响应
[type FinsReadMemoryResponse
    [const uint 16 mainFunction 0x0101]   // 主功能码：读取内存区
    [simple uint 16 statusCode]           // 状态码（如成功或错误）
    [array uint 8 data 'readLength']      // 返回的数据
]

[enum uint 8 OmronDataType(uint 8 dataTypeSize)
    ['1' UINT8 ['1']]
    ['2' INT8 ['1']]
    ['3' INT16 ['2']]
    ['4' UINT16 ['2']]
    ['5' INT32 ['4']]
    ['6' UINT32 ['4']]
    ['7' FLOAT ['4']]
    ['8' INT64 ['8']]
    ['9' UINT64 ['8']]
    ['10' DOUBLE ['8']]
    ['11' BIT ['1']]
    ['12' STRING ['2']]
]