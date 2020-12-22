//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package model

import (
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
)

type KnxManufacturers uint16

type IKnxManufacturers interface {
    Text() string
    Serialize(io utils.WriteBuffer) error
}

const(
    KnxManufacturers_M_0001 KnxManufacturers = 1
    KnxManufacturers_M_0002 KnxManufacturers = 2
    KnxManufacturers_M_0004 KnxManufacturers = 4
    KnxManufacturers_M_0005 KnxManufacturers = 5
    KnxManufacturers_M_0006 KnxManufacturers = 6
    KnxManufacturers_M_0007 KnxManufacturers = 7
    KnxManufacturers_M_0008 KnxManufacturers = 8
    KnxManufacturers_M_0009 KnxManufacturers = 9
    KnxManufacturers_M_000A KnxManufacturers = 10
    KnxManufacturers_M_000B KnxManufacturers = 11
    KnxManufacturers_M_000C KnxManufacturers = 12
    KnxManufacturers_M_000E KnxManufacturers = 14
    KnxManufacturers_M_0016 KnxManufacturers = 22
    KnxManufacturers_M_0018 KnxManufacturers = 24
    KnxManufacturers_M_0019 KnxManufacturers = 25
    KnxManufacturers_M_001B KnxManufacturers = 27
    KnxManufacturers_M_001C KnxManufacturers = 28
    KnxManufacturers_M_001D KnxManufacturers = 29
    KnxManufacturers_M_001E KnxManufacturers = 30
    KnxManufacturers_M_001F KnxManufacturers = 31
    KnxManufacturers_M_0020 KnxManufacturers = 32
    KnxManufacturers_M_0021 KnxManufacturers = 33
    KnxManufacturers_M_0022 KnxManufacturers = 34
    KnxManufacturers_M_0024 KnxManufacturers = 36
    KnxManufacturers_M_0025 KnxManufacturers = 37
    KnxManufacturers_M_0029 KnxManufacturers = 41
    KnxManufacturers_M_002A KnxManufacturers = 42
    KnxManufacturers_M_002C KnxManufacturers = 44
    KnxManufacturers_M_002D KnxManufacturers = 45
    KnxManufacturers_M_002E KnxManufacturers = 46
    KnxManufacturers_M_0031 KnxManufacturers = 49
    KnxManufacturers_M_0034 KnxManufacturers = 52
    KnxManufacturers_M_0035 KnxManufacturers = 53
    KnxManufacturers_M_0037 KnxManufacturers = 55
    KnxManufacturers_M_0039 KnxManufacturers = 57
    KnxManufacturers_M_003D KnxManufacturers = 61
    KnxManufacturers_M_003E KnxManufacturers = 62
    KnxManufacturers_M_0042 KnxManufacturers = 66
    KnxManufacturers_M_0043 KnxManufacturers = 67
    KnxManufacturers_M_0045 KnxManufacturers = 69
    KnxManufacturers_M_0047 KnxManufacturers = 71
    KnxManufacturers_M_0048 KnxManufacturers = 72
    KnxManufacturers_M_0049 KnxManufacturers = 73
    KnxManufacturers_M_004B KnxManufacturers = 75
    KnxManufacturers_M_004C KnxManufacturers = 76
    KnxManufacturers_M_004E KnxManufacturers = 78
    KnxManufacturers_M_0050 KnxManufacturers = 80
    KnxManufacturers_M_0051 KnxManufacturers = 81
    KnxManufacturers_M_0052 KnxManufacturers = 82
    KnxManufacturers_M_0053 KnxManufacturers = 83
    KnxManufacturers_M_0055 KnxManufacturers = 85
    KnxManufacturers_M_0059 KnxManufacturers = 89
    KnxManufacturers_M_005A KnxManufacturers = 90
    KnxManufacturers_M_005C KnxManufacturers = 92
    KnxManufacturers_M_005D KnxManufacturers = 93
    KnxManufacturers_M_005E KnxManufacturers = 94
    KnxManufacturers_M_005F KnxManufacturers = 95
    KnxManufacturers_M_0061 KnxManufacturers = 97
    KnxManufacturers_M_0062 KnxManufacturers = 98
    KnxManufacturers_M_0063 KnxManufacturers = 99
    KnxManufacturers_M_0064 KnxManufacturers = 100
    KnxManufacturers_M_0065 KnxManufacturers = 101
    KnxManufacturers_M_0066 KnxManufacturers = 102
    KnxManufacturers_M_0068 KnxManufacturers = 104
    KnxManufacturers_M_0069 KnxManufacturers = 105
    KnxManufacturers_M_006A KnxManufacturers = 106
    KnxManufacturers_M_006B KnxManufacturers = 107
    KnxManufacturers_M_006C KnxManufacturers = 108
    KnxManufacturers_M_006D KnxManufacturers = 109
    KnxManufacturers_M_006E KnxManufacturers = 110
    KnxManufacturers_M_006F KnxManufacturers = 111
    KnxManufacturers_M_0070 KnxManufacturers = 112
    KnxManufacturers_M_0071 KnxManufacturers = 113
    KnxManufacturers_M_0072 KnxManufacturers = 114
    KnxManufacturers_M_0073 KnxManufacturers = 115
    KnxManufacturers_M_0074 KnxManufacturers = 116
    KnxManufacturers_M_0075 KnxManufacturers = 117
    KnxManufacturers_M_0076 KnxManufacturers = 118
    KnxManufacturers_M_0077 KnxManufacturers = 119
    KnxManufacturers_M_0078 KnxManufacturers = 120
    KnxManufacturers_M_0079 KnxManufacturers = 121
    KnxManufacturers_M_007A KnxManufacturers = 122
    KnxManufacturers_M_007B KnxManufacturers = 123
    KnxManufacturers_M_007C KnxManufacturers = 124
    KnxManufacturers_M_007D KnxManufacturers = 125
    KnxManufacturers_M_007E KnxManufacturers = 126
    KnxManufacturers_M_007F KnxManufacturers = 127
    KnxManufacturers_M_0080 KnxManufacturers = 128
    KnxManufacturers_M_0081 KnxManufacturers = 129
    KnxManufacturers_M_0082 KnxManufacturers = 130
    KnxManufacturers_M_0083 KnxManufacturers = 131
    KnxManufacturers_M_0084 KnxManufacturers = 132
    KnxManufacturers_M_0085 KnxManufacturers = 133
    KnxManufacturers_M_0086 KnxManufacturers = 134
    KnxManufacturers_M_0087 KnxManufacturers = 135
    KnxManufacturers_M_0088 KnxManufacturers = 136
    KnxManufacturers_M_0089 KnxManufacturers = 137
    KnxManufacturers_M_008A KnxManufacturers = 138
    KnxManufacturers_M_008B KnxManufacturers = 139
    KnxManufacturers_M_008C KnxManufacturers = 140
    KnxManufacturers_M_008D KnxManufacturers = 141
    KnxManufacturers_M_008E KnxManufacturers = 142
    KnxManufacturers_M_008F KnxManufacturers = 143
    KnxManufacturers_M_0090 KnxManufacturers = 144
    KnxManufacturers_M_0091 KnxManufacturers = 145
    KnxManufacturers_M_0092 KnxManufacturers = 146
    KnxManufacturers_M_0093 KnxManufacturers = 147
    KnxManufacturers_M_0094 KnxManufacturers = 148
    KnxManufacturers_M_0095 KnxManufacturers = 149
    KnxManufacturers_M_0096 KnxManufacturers = 150
    KnxManufacturers_M_0097 KnxManufacturers = 151
    KnxManufacturers_M_0098 KnxManufacturers = 152
    KnxManufacturers_M_0099 KnxManufacturers = 153
    KnxManufacturers_M_009A KnxManufacturers = 154
    KnxManufacturers_M_009B KnxManufacturers = 155
    KnxManufacturers_M_009C KnxManufacturers = 156
    KnxManufacturers_M_009D KnxManufacturers = 157
    KnxManufacturers_M_009E KnxManufacturers = 158
    KnxManufacturers_M_009F KnxManufacturers = 159
    KnxManufacturers_M_00A0 KnxManufacturers = 160
    KnxManufacturers_M_00A1 KnxManufacturers = 161
    KnxManufacturers_M_00A2 KnxManufacturers = 162
    KnxManufacturers_M_00A3 KnxManufacturers = 163
    KnxManufacturers_M_00A4 KnxManufacturers = 164
    KnxManufacturers_M_00A5 KnxManufacturers = 165
    KnxManufacturers_M_00A6 KnxManufacturers = 166
    KnxManufacturers_M_00A7 KnxManufacturers = 167
    KnxManufacturers_M_00A8 KnxManufacturers = 168
    KnxManufacturers_M_00A9 KnxManufacturers = 169
    KnxManufacturers_M_00AA KnxManufacturers = 170
    KnxManufacturers_M_00AB KnxManufacturers = 171
    KnxManufacturers_M_00AC KnxManufacturers = 172
    KnxManufacturers_M_00AD KnxManufacturers = 173
    KnxManufacturers_M_00AE KnxManufacturers = 174
    KnxManufacturers_M_00AF KnxManufacturers = 175
    KnxManufacturers_M_00B0 KnxManufacturers = 176
    KnxManufacturers_M_00B1 KnxManufacturers = 177
    KnxManufacturers_M_00B2 KnxManufacturers = 178
    KnxManufacturers_M_00B3 KnxManufacturers = 179
    KnxManufacturers_M_00B4 KnxManufacturers = 180
    KnxManufacturers_M_00B5 KnxManufacturers = 181
    KnxManufacturers_M_00B6 KnxManufacturers = 182
    KnxManufacturers_M_00B7 KnxManufacturers = 183
    KnxManufacturers_M_00B8 KnxManufacturers = 184
    KnxManufacturers_M_00B9 KnxManufacturers = 185
    KnxManufacturers_M_00BA KnxManufacturers = 186
    KnxManufacturers_M_00BB KnxManufacturers = 187
    KnxManufacturers_M_00BC KnxManufacturers = 188
    KnxManufacturers_M_00BD KnxManufacturers = 189
    KnxManufacturers_M_00BE KnxManufacturers = 190
    KnxManufacturers_M_00BF KnxManufacturers = 191
    KnxManufacturers_M_00C0 KnxManufacturers = 192
    KnxManufacturers_M_00C1 KnxManufacturers = 193
    KnxManufacturers_M_00C2 KnxManufacturers = 194
    KnxManufacturers_M_00C3 KnxManufacturers = 195
    KnxManufacturers_M_00C4 KnxManufacturers = 196
    KnxManufacturers_M_00C5 KnxManufacturers = 197
    KnxManufacturers_M_00C6 KnxManufacturers = 198
    KnxManufacturers_M_00C7 KnxManufacturers = 199
    KnxManufacturers_M_00C8 KnxManufacturers = 200
    KnxManufacturers_M_00C9 KnxManufacturers = 201
    KnxManufacturers_M_00CA KnxManufacturers = 202
    KnxManufacturers_M_00CC KnxManufacturers = 204
    KnxManufacturers_M_00CD KnxManufacturers = 205
    KnxManufacturers_M_00CE KnxManufacturers = 206
    KnxManufacturers_M_00CF KnxManufacturers = 207
    KnxManufacturers_M_00D0 KnxManufacturers = 208
    KnxManufacturers_M_00D1 KnxManufacturers = 209
    KnxManufacturers_M_00D2 KnxManufacturers = 210
    KnxManufacturers_M_00D3 KnxManufacturers = 211
    KnxManufacturers_M_00D6 KnxManufacturers = 214
    KnxManufacturers_M_00D7 KnxManufacturers = 215
    KnxManufacturers_M_00D8 KnxManufacturers = 216
    KnxManufacturers_M_00D9 KnxManufacturers = 217
    KnxManufacturers_M_00DA KnxManufacturers = 218
    KnxManufacturers_M_00DB KnxManufacturers = 219
    KnxManufacturers_M_00DC KnxManufacturers = 220
    KnxManufacturers_M_00DE KnxManufacturers = 222
    KnxManufacturers_M_00DF KnxManufacturers = 223
    KnxManufacturers_M_00E1 KnxManufacturers = 225
    KnxManufacturers_M_00E3 KnxManufacturers = 227
    KnxManufacturers_M_00E4 KnxManufacturers = 228
    KnxManufacturers_M_00E8 KnxManufacturers = 232
    KnxManufacturers_M_00E9 KnxManufacturers = 233
    KnxManufacturers_M_00EA KnxManufacturers = 234
    KnxManufacturers_M_00EB KnxManufacturers = 235
    KnxManufacturers_M_00ED KnxManufacturers = 237
    KnxManufacturers_M_00EE KnxManufacturers = 238
    KnxManufacturers_M_00EF KnxManufacturers = 239
    KnxManufacturers_M_00F0 KnxManufacturers = 240
    KnxManufacturers_M_00F1 KnxManufacturers = 241
    KnxManufacturers_M_00F2 KnxManufacturers = 242
    KnxManufacturers_M_00F4 KnxManufacturers = 244
    KnxManufacturers_M_00F5 KnxManufacturers = 245
    KnxManufacturers_M_00F6 KnxManufacturers = 246
    KnxManufacturers_M_00F8 KnxManufacturers = 248
    KnxManufacturers_M_00F9 KnxManufacturers = 249
    KnxManufacturers_M_00FA KnxManufacturers = 250
    KnxManufacturers_M_00FB KnxManufacturers = 251
    KnxManufacturers_M_00FC KnxManufacturers = 252
    KnxManufacturers_M_00FD KnxManufacturers = 253
    KnxManufacturers_M_00FE KnxManufacturers = 254
    KnxManufacturers_M_0100 KnxManufacturers = 256
    KnxManufacturers_M_0101 KnxManufacturers = 257
    KnxManufacturers_M_0102 KnxManufacturers = 258
    KnxManufacturers_M_0103 KnxManufacturers = 259
    KnxManufacturers_M_0104 KnxManufacturers = 260
    KnxManufacturers_M_0105 KnxManufacturers = 261
    KnxManufacturers_M_0106 KnxManufacturers = 262
    KnxManufacturers_M_0107 KnxManufacturers = 263
    KnxManufacturers_M_0108 KnxManufacturers = 264
    KnxManufacturers_M_0109 KnxManufacturers = 265
    KnxManufacturers_M_010A KnxManufacturers = 266
    KnxManufacturers_M_010B KnxManufacturers = 267
    KnxManufacturers_M_010C KnxManufacturers = 268
    KnxManufacturers_M_010D KnxManufacturers = 269
    KnxManufacturers_M_010E KnxManufacturers = 270
    KnxManufacturers_M_010F KnxManufacturers = 271
    KnxManufacturers_M_0110 KnxManufacturers = 272
    KnxManufacturers_M_0111 KnxManufacturers = 273
    KnxManufacturers_M_0112 KnxManufacturers = 274
    KnxManufacturers_M_0113 KnxManufacturers = 275
    KnxManufacturers_M_0114 KnxManufacturers = 276
    KnxManufacturers_M_0115 KnxManufacturers = 277
    KnxManufacturers_M_0116 KnxManufacturers = 278
    KnxManufacturers_M_0117 KnxManufacturers = 279
    KnxManufacturers_M_0118 KnxManufacturers = 280
    KnxManufacturers_M_0119 KnxManufacturers = 281
    KnxManufacturers_M_011A KnxManufacturers = 282
    KnxManufacturers_M_011B KnxManufacturers = 283
    KnxManufacturers_M_011C KnxManufacturers = 284
    KnxManufacturers_M_011D KnxManufacturers = 285
    KnxManufacturers_M_011E KnxManufacturers = 286
    KnxManufacturers_M_011F KnxManufacturers = 287
    KnxManufacturers_M_0120 KnxManufacturers = 288
    KnxManufacturers_M_0121 KnxManufacturers = 289
    KnxManufacturers_M_0122 KnxManufacturers = 290
    KnxManufacturers_M_0123 KnxManufacturers = 291
    KnxManufacturers_M_0124 KnxManufacturers = 292
    KnxManufacturers_M_0125 KnxManufacturers = 293
    KnxManufacturers_M_0126 KnxManufacturers = 294
    KnxManufacturers_M_0127 KnxManufacturers = 295
    KnxManufacturers_M_0128 KnxManufacturers = 296
    KnxManufacturers_M_0129 KnxManufacturers = 297
    KnxManufacturers_M_012A KnxManufacturers = 298
    KnxManufacturers_M_012B KnxManufacturers = 299
    KnxManufacturers_M_012C KnxManufacturers = 300
    KnxManufacturers_M_012D KnxManufacturers = 301
    KnxManufacturers_M_012E KnxManufacturers = 302
    KnxManufacturers_M_012F KnxManufacturers = 303
    KnxManufacturers_M_0130 KnxManufacturers = 304
    KnxManufacturers_M_0131 KnxManufacturers = 305
    KnxManufacturers_M_0132 KnxManufacturers = 306
    KnxManufacturers_M_0133 KnxManufacturers = 307
    KnxManufacturers_M_0134 KnxManufacturers = 308
    KnxManufacturers_M_0135 KnxManufacturers = 309
    KnxManufacturers_M_0136 KnxManufacturers = 310
    KnxManufacturers_M_0137 KnxManufacturers = 311
    KnxManufacturers_M_0138 KnxManufacturers = 312
    KnxManufacturers_M_0139 KnxManufacturers = 313
    KnxManufacturers_M_013A KnxManufacturers = 314
    KnxManufacturers_M_013B KnxManufacturers = 315
    KnxManufacturers_M_013C KnxManufacturers = 316
    KnxManufacturers_M_013D KnxManufacturers = 317
    KnxManufacturers_M_013E KnxManufacturers = 318
    KnxManufacturers_M_013F KnxManufacturers = 319
    KnxManufacturers_M_0140 KnxManufacturers = 320
    KnxManufacturers_M_0141 KnxManufacturers = 321
    KnxManufacturers_M_0142 KnxManufacturers = 322
    KnxManufacturers_M_0143 KnxManufacturers = 323
    KnxManufacturers_M_0144 KnxManufacturers = 324
    KnxManufacturers_M_0145 KnxManufacturers = 325
    KnxManufacturers_M_0146 KnxManufacturers = 326
    KnxManufacturers_M_0147 KnxManufacturers = 327
    KnxManufacturers_M_0148 KnxManufacturers = 328
    KnxManufacturers_M_0149 KnxManufacturers = 329
    KnxManufacturers_M_014A KnxManufacturers = 330
    KnxManufacturers_M_014B KnxManufacturers = 331
    KnxManufacturers_M_014C KnxManufacturers = 332
    KnxManufacturers_M_014D KnxManufacturers = 333
    KnxManufacturers_M_014E KnxManufacturers = 334
    KnxManufacturers_M_014F KnxManufacturers = 335
    KnxManufacturers_M_0150 KnxManufacturers = 336
    KnxManufacturers_M_0151 KnxManufacturers = 337
    KnxManufacturers_M_0152 KnxManufacturers = 338
    KnxManufacturers_M_0153 KnxManufacturers = 339
    KnxManufacturers_M_0154 KnxManufacturers = 340
    KnxManufacturers_M_0155 KnxManufacturers = 341
    KnxManufacturers_M_0156 KnxManufacturers = 342
    KnxManufacturers_M_0157 KnxManufacturers = 343
    KnxManufacturers_M_0158 KnxManufacturers = 344
    KnxManufacturers_M_0159 KnxManufacturers = 345
    KnxManufacturers_M_015A KnxManufacturers = 346
    KnxManufacturers_M_015B KnxManufacturers = 347
    KnxManufacturers_M_015C KnxManufacturers = 348
    KnxManufacturers_M_015D KnxManufacturers = 349
    KnxManufacturers_M_015E KnxManufacturers = 350
    KnxManufacturers_M_015F KnxManufacturers = 351
    KnxManufacturers_M_0160 KnxManufacturers = 352
    KnxManufacturers_M_0161 KnxManufacturers = 353
    KnxManufacturers_M_0162 KnxManufacturers = 354
    KnxManufacturers_M_0163 KnxManufacturers = 355
    KnxManufacturers_M_0164 KnxManufacturers = 356
    KnxManufacturers_M_0165 KnxManufacturers = 357
    KnxManufacturers_M_0166 KnxManufacturers = 358
    KnxManufacturers_M_0167 KnxManufacturers = 359
    KnxManufacturers_M_0168 KnxManufacturers = 360
    KnxManufacturers_M_0169 KnxManufacturers = 361
    KnxManufacturers_M_016A KnxManufacturers = 362
    KnxManufacturers_M_016B KnxManufacturers = 363
    KnxManufacturers_M_016C KnxManufacturers = 364
    KnxManufacturers_M_016D KnxManufacturers = 365
    KnxManufacturers_M_016E KnxManufacturers = 366
    KnxManufacturers_M_016F KnxManufacturers = 367
    KnxManufacturers_M_0170 KnxManufacturers = 368
    KnxManufacturers_M_0171 KnxManufacturers = 369
    KnxManufacturers_M_0172 KnxManufacturers = 370
    KnxManufacturers_M_0173 KnxManufacturers = 371
    KnxManufacturers_M_0174 KnxManufacturers = 372
    KnxManufacturers_M_0175 KnxManufacturers = 373
    KnxManufacturers_M_0176 KnxManufacturers = 374
    KnxManufacturers_M_0177 KnxManufacturers = 375
    KnxManufacturers_M_0178 KnxManufacturers = 376
    KnxManufacturers_M_0179 KnxManufacturers = 377
    KnxManufacturers_M_017A KnxManufacturers = 378
    KnxManufacturers_M_017B KnxManufacturers = 379
    KnxManufacturers_M_017C KnxManufacturers = 380
    KnxManufacturers_M_017D KnxManufacturers = 381
    KnxManufacturers_M_017E KnxManufacturers = 382
    KnxManufacturers_M_017F KnxManufacturers = 383
    KnxManufacturers_M_0180 KnxManufacturers = 384
    KnxManufacturers_M_0181 KnxManufacturers = 385
    KnxManufacturers_M_0182 KnxManufacturers = 386
    KnxManufacturers_M_0183 KnxManufacturers = 387
    KnxManufacturers_M_0184 KnxManufacturers = 388
    KnxManufacturers_M_0185 KnxManufacturers = 389
    KnxManufacturers_M_0186 KnxManufacturers = 390
    KnxManufacturers_M_0187 KnxManufacturers = 391
    KnxManufacturers_M_0188 KnxManufacturers = 392
    KnxManufacturers_M_0189 KnxManufacturers = 393
    KnxManufacturers_M_018A KnxManufacturers = 394
    KnxManufacturers_M_018B KnxManufacturers = 395
    KnxManufacturers_M_018C KnxManufacturers = 396
    KnxManufacturers_M_018D KnxManufacturers = 397
    KnxManufacturers_M_018E KnxManufacturers = 398
    KnxManufacturers_M_018F KnxManufacturers = 399
    KnxManufacturers_M_0190 KnxManufacturers = 400
    KnxManufacturers_M_0191 KnxManufacturers = 401
    KnxManufacturers_M_0192 KnxManufacturers = 402
    KnxManufacturers_M_0193 KnxManufacturers = 403
    KnxManufacturers_M_0194 KnxManufacturers = 404
    KnxManufacturers_M_0195 KnxManufacturers = 405
    KnxManufacturers_M_0196 KnxManufacturers = 406
    KnxManufacturers_M_0197 KnxManufacturers = 407
    KnxManufacturers_M_0198 KnxManufacturers = 408
    KnxManufacturers_M_0199 KnxManufacturers = 409
    KnxManufacturers_M_019A KnxManufacturers = 410
    KnxManufacturers_M_019B KnxManufacturers = 411
    KnxManufacturers_M_019C KnxManufacturers = 412
    KnxManufacturers_M_019D KnxManufacturers = 413
    KnxManufacturers_M_019E KnxManufacturers = 414
    KnxManufacturers_M_019F KnxManufacturers = 415
    KnxManufacturers_M_01A0 KnxManufacturers = 416
    KnxManufacturers_M_01A1 KnxManufacturers = 417
    KnxManufacturers_M_01A2 KnxManufacturers = 418
    KnxManufacturers_M_01A3 KnxManufacturers = 419
    KnxManufacturers_M_01A4 KnxManufacturers = 420
    KnxManufacturers_M_01A5 KnxManufacturers = 421
    KnxManufacturers_M_01A6 KnxManufacturers = 422
    KnxManufacturers_M_01A7 KnxManufacturers = 423
    KnxManufacturers_M_01A8 KnxManufacturers = 424
    KnxManufacturers_M_01A9 KnxManufacturers = 425
    KnxManufacturers_M_01AA KnxManufacturers = 426
    KnxManufacturers_M_01AB KnxManufacturers = 427
    KnxManufacturers_M_01AC KnxManufacturers = 428
    KnxManufacturers_M_01AD KnxManufacturers = 429
    KnxManufacturers_M_01AE KnxManufacturers = 430
    KnxManufacturers_M_01AF KnxManufacturers = 431
    KnxManufacturers_M_01B0 KnxManufacturers = 432
    KnxManufacturers_M_01B1 KnxManufacturers = 433
    KnxManufacturers_M_01B2 KnxManufacturers = 434
    KnxManufacturers_M_01B3 KnxManufacturers = 435
    KnxManufacturers_M_01B4 KnxManufacturers = 436
    KnxManufacturers_M_01B5 KnxManufacturers = 437
    KnxManufacturers_M_01B6 KnxManufacturers = 438
    KnxManufacturers_M_01B7 KnxManufacturers = 439
    KnxManufacturers_M_01B8 KnxManufacturers = 440
    KnxManufacturers_M_01B9 KnxManufacturers = 441
    KnxManufacturers_M_01BA KnxManufacturers = 442
    KnxManufacturers_M_01BB KnxManufacturers = 443
    KnxManufacturers_M_01BC KnxManufacturers = 444
    KnxManufacturers_M_01BD KnxManufacturers = 445
    KnxManufacturers_M_01BE KnxManufacturers = 446
    KnxManufacturers_M_01BF KnxManufacturers = 447
    KnxManufacturers_M_01C0 KnxManufacturers = 448
    KnxManufacturers_M_01C1 KnxManufacturers = 449
    KnxManufacturers_M_01C3 KnxManufacturers = 451
    KnxManufacturers_M_01C4 KnxManufacturers = 452
    KnxManufacturers_M_01C5 KnxManufacturers = 453
    KnxManufacturers_M_01C6 KnxManufacturers = 454
    KnxManufacturers_M_01C7 KnxManufacturers = 455
    KnxManufacturers_M_01C8 KnxManufacturers = 456
    KnxManufacturers_M_01C9 KnxManufacturers = 457
    KnxManufacturers_M_01CA KnxManufacturers = 458
    KnxManufacturers_M_01CB KnxManufacturers = 459
    KnxManufacturers_M_01CC KnxManufacturers = 460
    KnxManufacturers_M_01CD KnxManufacturers = 461
    KnxManufacturers_M_01CE KnxManufacturers = 462
    KnxManufacturers_M_01CF KnxManufacturers = 463
    KnxManufacturers_M_01D0 KnxManufacturers = 464
    KnxManufacturers_M_01D1 KnxManufacturers = 465
    KnxManufacturers_M_01D2 KnxManufacturers = 466
    KnxManufacturers_M_01D3 KnxManufacturers = 467
    KnxManufacturers_M_01D4 KnxManufacturers = 468
    KnxManufacturers_M_01D5 KnxManufacturers = 469
    KnxManufacturers_M_01D6 KnxManufacturers = 470
    KnxManufacturers_M_01D7 KnxManufacturers = 471
    KnxManufacturers_M_01D8 KnxManufacturers = 472
    KnxManufacturers_M_01D9 KnxManufacturers = 473
    KnxManufacturers_M_01DA KnxManufacturers = 474
    KnxManufacturers_M_01DB KnxManufacturers = 475
    KnxManufacturers_M_01DC KnxManufacturers = 476
    KnxManufacturers_M_01DD KnxManufacturers = 477
    KnxManufacturers_M_01DE KnxManufacturers = 478
    KnxManufacturers_M_01DF KnxManufacturers = 479
    KnxManufacturers_M_01E0 KnxManufacturers = 480
    KnxManufacturers_M_01E1 KnxManufacturers = 481
    KnxManufacturers_M_01E2 KnxManufacturers = 482
    KnxManufacturers_M_01E3 KnxManufacturers = 483
    KnxManufacturers_M_01E4 KnxManufacturers = 484
    KnxManufacturers_M_01E5 KnxManufacturers = 485
    KnxManufacturers_M_01E6 KnxManufacturers = 486
    KnxManufacturers_M_01E7 KnxManufacturers = 487
    KnxManufacturers_M_01E8 KnxManufacturers = 488
    KnxManufacturers_M_01E9 KnxManufacturers = 489
    KnxManufacturers_M_01EA KnxManufacturers = 490
    KnxManufacturers_M_01EB KnxManufacturers = 491
    KnxManufacturers_M_01EC KnxManufacturers = 492
    KnxManufacturers_M_01ED KnxManufacturers = 493
    KnxManufacturers_M_01EF KnxManufacturers = 495
    KnxManufacturers_M_01F0 KnxManufacturers = 496
    KnxManufacturers_M_01F1 KnxManufacturers = 497
    KnxManufacturers_M_01F2 KnxManufacturers = 498
    KnxManufacturers_M_01F3 KnxManufacturers = 499
    KnxManufacturers_M_01F4 KnxManufacturers = 500
    KnxManufacturers_M_01F5 KnxManufacturers = 501
    KnxManufacturers_M_01F6 KnxManufacturers = 502
    KnxManufacturers_M_01F7 KnxManufacturers = 503
    KnxManufacturers_M_01F8 KnxManufacturers = 504
    KnxManufacturers_M_01F9 KnxManufacturers = 505
    KnxManufacturers_M_01FA KnxManufacturers = 506
    KnxManufacturers_M_01FB KnxManufacturers = 507
    KnxManufacturers_M_01FC KnxManufacturers = 508
    KnxManufacturers_M_01FD KnxManufacturers = 509
    KnxManufacturers_M_0200 KnxManufacturers = 512
    KnxManufacturers_M_0201 KnxManufacturers = 513
    KnxManufacturers_M_0202 KnxManufacturers = 514
    KnxManufacturers_M_0203 KnxManufacturers = 515
    KnxManufacturers_M_0204 KnxManufacturers = 516
    KnxManufacturers_M_0205 KnxManufacturers = 517
    KnxManufacturers_M_0206 KnxManufacturers = 518
    KnxManufacturers_M_0207 KnxManufacturers = 519
    KnxManufacturers_M_0208 KnxManufacturers = 520
    KnxManufacturers_M_0209 KnxManufacturers = 521
    KnxManufacturers_M_020A KnxManufacturers = 522
    KnxManufacturers_M_020B KnxManufacturers = 523
    KnxManufacturers_M_020C KnxManufacturers = 524
    KnxManufacturers_M_020D KnxManufacturers = 525
    KnxManufacturers_M_020E KnxManufacturers = 526
    KnxManufacturers_M_020F KnxManufacturers = 527
    KnxManufacturers_M_0210 KnxManufacturers = 528
    KnxManufacturers_M_0211 KnxManufacturers = 529
    KnxManufacturers_M_0212 KnxManufacturers = 530
    KnxManufacturers_M_0213 KnxManufacturers = 531
    KnxManufacturers_M_0214 KnxManufacturers = 532
    KnxManufacturers_M_0215 KnxManufacturers = 533
    KnxManufacturers_M_0216 KnxManufacturers = 534
    KnxManufacturers_M_0217 KnxManufacturers = 535
    KnxManufacturers_M_0218 KnxManufacturers = 536
    KnxManufacturers_M_0219 KnxManufacturers = 537
    KnxManufacturers_M_021A KnxManufacturers = 538
    KnxManufacturers_M_021B KnxManufacturers = 539
    KnxManufacturers_M_021C KnxManufacturers = 540
    KnxManufacturers_M_021D KnxManufacturers = 541
    KnxManufacturers_M_021E KnxManufacturers = 542
    KnxManufacturers_M_021F KnxManufacturers = 543
    KnxManufacturers_M_0220 KnxManufacturers = 544
    KnxManufacturers_M_0221 KnxManufacturers = 545
    KnxManufacturers_M_0222 KnxManufacturers = 546
    KnxManufacturers_M_0223 KnxManufacturers = 547
    KnxManufacturers_M_0224 KnxManufacturers = 548
    KnxManufacturers_M_0225 KnxManufacturers = 549
    KnxManufacturers_M_0226 KnxManufacturers = 550
    KnxManufacturers_M_0227 KnxManufacturers = 551
    KnxManufacturers_M_0228 KnxManufacturers = 552
    KnxManufacturers_M_0229 KnxManufacturers = 553
    KnxManufacturers_M_022A KnxManufacturers = 554
    KnxManufacturers_M_022B KnxManufacturers = 555
    KnxManufacturers_M_022C KnxManufacturers = 556
    KnxManufacturers_M_022D KnxManufacturers = 557
    KnxManufacturers_M_022E KnxManufacturers = 558
    KnxManufacturers_M_022F KnxManufacturers = 559
    KnxManufacturers_M_0230 KnxManufacturers = 560
    KnxManufacturers_M_0231 KnxManufacturers = 561
    KnxManufacturers_M_0232 KnxManufacturers = 562
    KnxManufacturers_M_0233 KnxManufacturers = 563
    KnxManufacturers_M_0234 KnxManufacturers = 564
    KnxManufacturers_M_0235 KnxManufacturers = 565
    KnxManufacturers_M_0236 KnxManufacturers = 566
    KnxManufacturers_M_0237 KnxManufacturers = 567
    KnxManufacturers_M_0238 KnxManufacturers = 568
    KnxManufacturers_M_0239 KnxManufacturers = 569
    KnxManufacturers_M_023A KnxManufacturers = 570
    KnxManufacturers_M_023B KnxManufacturers = 571
    KnxManufacturers_M_023C KnxManufacturers = 572
    KnxManufacturers_M_023D KnxManufacturers = 573
    KnxManufacturers_M_023E KnxManufacturers = 574
    KnxManufacturers_M_023F KnxManufacturers = 575
    KnxManufacturers_M_0240 KnxManufacturers = 576
    KnxManufacturers_M_0241 KnxManufacturers = 577
    KnxManufacturers_M_0242 KnxManufacturers = 578
    KnxManufacturers_M_0243 KnxManufacturers = 579
    KnxManufacturers_M_0244 KnxManufacturers = 580
    KnxManufacturers_M_0245 KnxManufacturers = 581
    KnxManufacturers_M_0246 KnxManufacturers = 582
    KnxManufacturers_M_0247 KnxManufacturers = 583
    KnxManufacturers_M_0248 KnxManufacturers = 584
    KnxManufacturers_M_0249 KnxManufacturers = 585
    KnxManufacturers_M_024A KnxManufacturers = 586
    KnxManufacturers_M_024B KnxManufacturers = 587
    KnxManufacturers_M_024C KnxManufacturers = 588
    KnxManufacturers_M_024D KnxManufacturers = 589
    KnxManufacturers_M_024E KnxManufacturers = 590
    KnxManufacturers_M_024F KnxManufacturers = 591
    KnxManufacturers_M_0250 KnxManufacturers = 592
    KnxManufacturers_M_ABB2 KnxManufacturers = 43954
    KnxManufacturers_M_ABB7 KnxManufacturers = 43959
)


func (e KnxManufacturers) Text() string {
    switch e  {
        case 1: { /* '1' */
            return "Siemens"
        }
        case 10: { /* '10' */
            return "Insta GmbH"
        }
        case 100: { /* '100' */
            return "Schneider Electric Industries SAS"
        }
        case 101: { /* '101' */
            return "WHD Wilhelm Huber + Söhne"
        }
        case 102: { /* '102' */
            return "Bischoff Elektronik"
        }
        case 104: { /* '104' */
            return "JEPAZ"
        }
        case 105: { /* '105' */
            return "RTS Automation"
        }
        case 106: { /* '106' */
            return "EIBMARKT GmbH"
        }
        case 107: { /* '107' */
            return "WAREMA Renkhoff SE"
        }
        case 108: { /* '108' */
            return "Eelectron"
        }
        case 109: { /* '109' */
            return "Belden Wire & Cable B.V."
        }
        case 11: { /* '11' */
            return "LEGRAND Appareillage électrique"
        }
        case 110: { /* '110' */
            return "Becker-Antriebe GmbH"
        }
        case 111: { /* '111' */
            return "J.Stehle+Söhne GmbH"
        }
        case 112: { /* '112' */
            return "AGFEO"
        }
        case 113: { /* '113' */
            return "Zennio"
        }
        case 114: { /* '114' */
            return "TAPKO Technologies"
        }
        case 115: { /* '115' */
            return "HDL"
        }
        case 116: { /* '116' */
            return "Uponor"
        }
        case 117: { /* '117' */
            return "se Lightmanagement AG"
        }
        case 118: { /* '118' */
            return "Arcus-eds"
        }
        case 119: { /* '119' */
            return "Intesis"
        }
        case 12: { /* '12' */
            return "Merten"
        }
        case 120: { /* '120' */
            return "Herholdt Controls srl"
        }
        case 121: { /* '121' */
            return "Niko-Zublin"
        }
        case 122: { /* '122' */
            return "Durable Technologies"
        }
        case 123: { /* '123' */
            return "Innoteam"
        }
        case 124: { /* '124' */
            return "ise GmbH"
        }
        case 125: { /* '125' */
            return "TEAM FOR TRONICS"
        }
        case 126: { /* '126' */
            return "CIAT"
        }
        case 127: { /* '127' */
            return "Remeha BV"
        }
        case 128: { /* '128' */
            return "ESYLUX"
        }
        case 129: { /* '129' */
            return "BASALTE"
        }
        case 130: { /* '130' */
            return "Vestamatic"
        }
        case 131: { /* '131' */
            return "MDT technologies"
        }
        case 132: { /* '132' */
            return "Warendorfer Küchen GmbH"
        }
        case 133: { /* '133' */
            return "Video-Star"
        }
        case 134: { /* '134' */
            return "Sitek"
        }
        case 135: { /* '135' */
            return "CONTROLtronic"
        }
        case 136: { /* '136' */
            return "function Technology"
        }
        case 137: { /* '137' */
            return "AMX"
        }
        case 138: { /* '138' */
            return "ELDAT"
        }
        case 139: { /* '139' */
            return "Panasonic"
        }
        case 14: { /* '14' */
            return "ABB SpA-SACE Division"
        }
        case 140: { /* '140' */
            return "Pulse Technologies"
        }
        case 141: { /* '141' */
            return "Crestron"
        }
        case 142: { /* '142' */
            return "STEINEL professional"
        }
        case 143: { /* '143' */
            return "BILTON LED Lighting"
        }
        case 144: { /* '144' */
            return "denro AG"
        }
        case 145: { /* '145' */
            return "GePro"
        }
        case 146: { /* '146' */
            return "preussen automation"
        }
        case 147: { /* '147' */
            return "Zoppas Industries"
        }
        case 148: { /* '148' */
            return "MACTECH"
        }
        case 149: { /* '149' */
            return "TECHNO-TREND"
        }
        case 150: { /* '150' */
            return "FS Cables"
        }
        case 151: { /* '151' */
            return "Delta Dore"
        }
        case 152: { /* '152' */
            return "Eissound"
        }
        case 153: { /* '153' */
            return "Cisco"
        }
        case 154: { /* '154' */
            return "Dinuy"
        }
        case 155: { /* '155' */
            return "iKNiX"
        }
        case 156: { /* '156' */
            return "Rademacher Geräte-Elektronik GmbH"
        }
        case 157: { /* '157' */
            return "EGi Electroacustica General Iberica"
        }
        case 158: { /* '158' */
            return "Bes – Ingenium"
        }
        case 159: { /* '159' */
            return "ElabNET"
        }
        case 160: { /* '160' */
            return "Blumotix"
        }
        case 161: { /* '161' */
            return "Hunter Douglas"
        }
        case 162: { /* '162' */
            return "APRICUM"
        }
        case 163: { /* '163' */
            return "TIANSU Automation"
        }
        case 164: { /* '164' */
            return "Bubendorff"
        }
        case 165: { /* '165' */
            return "MBS GmbH"
        }
        case 166: { /* '166' */
            return "Enertex Bayern GmbH"
        }
        case 167: { /* '167' */
            return "BMS"
        }
        case 168: { /* '168' */
            return "Sinapsi"
        }
        case 169: { /* '169' */
            return "Embedded Systems SIA"
        }
        case 170: { /* '170' */
            return "KNX1"
        }
        case 171: { /* '171' */
            return "Tokka"
        }
        case 172: { /* '172' */
            return "NanoSense"
        }
        case 173: { /* '173' */
            return "PEAR Automation GmbH"
        }
        case 174: { /* '174' */
            return "DGA"
        }
        case 175: { /* '175' */
            return "Lutron"
        }
        case 176: { /* '176' */
            return "AIRZONE – ALTRA"
        }
        case 177: { /* '177' */
            return "Lithoss Design Switches"
        }
        case 178: { /* '178' */
            return "3ATEL"
        }
        case 179: { /* '179' */
            return "Philips Controls"
        }
        case 180: { /* '180' */
            return "VELUX A/S"
        }
        case 181: { /* '181' */
            return "LOYTEC"
        }
        case 182: { /* '182' */
            return "Ekinex S.p.A."
        }
        case 183: { /* '183' */
            return "SIRLAN Technologies"
        }
        case 184: { /* '184' */
            return "ProKNX SAS"
        }
        case 185: { /* '185' */
            return "IT GmbH"
        }
        case 186: { /* '186' */
            return "RENSON"
        }
        case 187: { /* '187' */
            return "HEP Group"
        }
        case 188: { /* '188' */
            return "Balmart"
        }
        case 189: { /* '189' */
            return "GFS GmbH"
        }
        case 190: { /* '190' */
            return "Schenker Storen AG"
        }
        case 191: { /* '191' */
            return "Algodue Elettronica S.r.L."
        }
        case 192: { /* '192' */
            return "ABB France"
        }
        case 193: { /* '193' */
            return "maintronic"
        }
        case 194: { /* '194' */
            return "Vantage"
        }
        case 195: { /* '195' */
            return "Foresis"
        }
        case 196: { /* '196' */
            return "Research & Production Association SEM"
        }
        case 197: { /* '197' */
            return "Weinzierl Engineering GmbH"
        }
        case 198: { /* '198' */
            return "Möhlenhoff Wärmetechnik GmbH"
        }
        case 199: { /* '199' */
            return "PKC-GROUP Oyj"
        }
        case 2: { /* '2' */
            return "ABB"
        }
        case 200: { /* '200' */
            return "B.E.G."
        }
        case 201: { /* '201' */
            return "Elsner Elektronik GmbH"
        }
        case 202: { /* '202' */
            return "Siemens Building Technologies (HK/China) Ltd."
        }
        case 204: { /* '204' */
            return "Eutrac"
        }
        case 205: { /* '205' */
            return "Gustav Hensel GmbH & Co. KG"
        }
        case 206: { /* '206' */
            return "GARO AB"
        }
        case 207: { /* '207' */
            return "Waldmann Lichttechnik"
        }
        case 208: { /* '208' */
            return "SCHÜCO"
        }
        case 209: { /* '209' */
            return "EMU"
        }
        case 210: { /* '210' */
            return "JNet Systems AG"
        }
        case 211: { /* '211' */
            return "Total Solution GmbH"
        }
        case 214: { /* '214' */
            return "O.Y.L. Electronics"
        }
        case 215: { /* '215' */
            return "Galax System"
        }
        case 216: { /* '216' */
            return "Disch"
        }
        case 217: { /* '217' */
            return "Aucoteam"
        }
        case 218: { /* '218' */
            return "Luxmate Controls"
        }
        case 219: { /* '219' */
            return "Danfoss"
        }
        case 22: { /* '22' */
            return "Siedle & Söhne"
        }
        case 220: { /* '220' */
            return "AST GmbH"
        }
        case 222: { /* '222' */
            return "WILA Leuchten"
        }
        case 223: { /* '223' */
            return "b+b Automations- und Steuerungstechnik"
        }
        case 225: { /* '225' */
            return "Lingg & Janke"
        }
        case 227: { /* '227' */
            return "Sauter"
        }
        case 228: { /* '228' */
            return "SIMU"
        }
        case 232: { /* '232' */
            return "Theben HTS AG"
        }
        case 233: { /* '233' */
            return "Amann GmbH"
        }
        case 234: { /* '234' */
            return "BERG Energiekontrollsysteme GmbH"
        }
        case 235: { /* '235' */
            return "Hüppe Form Sonnenschutzsysteme GmbH"
        }
        case 237: { /* '237' */
            return "Oventrop KG"
        }
        case 238: { /* '238' */
            return "Griesser AG"
        }
        case 239: { /* '239' */
            return "IPAS GmbH"
        }
        case 24: { /* '24' */
            return "Eberle"
        }
        case 240: { /* '240' */
            return "elero GmbH"
        }
        case 241: { /* '241' */
            return "Ardan Production and Industrial Controls Ltd."
        }
        case 242: { /* '242' */
            return "Metec Meßtechnik GmbH"
        }
        case 244: { /* '244' */
            return "ELKA-Elektronik GmbH"
        }
        case 245: { /* '245' */
            return "ELEKTROANLAGEN D. NAGEL"
        }
        case 246: { /* '246' */
            return "Tridonic Bauelemente GmbH"
        }
        case 248: { /* '248' */
            return "Stengler Gesellschaft"
        }
        case 249: { /* '249' */
            return "Schneider Electric (MG)"
        }
        case 25: { /* '25' */
            return "GEWISS"
        }
        case 250: { /* '250' */
            return "KNX Association"
        }
        case 251: { /* '251' */
            return "VIVO"
        }
        case 252: { /* '252' */
            return "Hugo Müller GmbH & Co KG"
        }
        case 253: { /* '253' */
            return "Siemens HVAC"
        }
        case 254: { /* '254' */
            return "APT"
        }
        case 256: { /* '256' */
            return "HighDom"
        }
        case 257: { /* '257' */
            return "Top Services"
        }
        case 258: { /* '258' */
            return "ambiHome"
        }
        case 259: { /* '259' */
            return "DATEC electronic AG"
        }
        case 260: { /* '260' */
            return "ABUS Security-Center"
        }
        case 261: { /* '261' */
            return "Lite-Puter"
        }
        case 262: { /* '262' */
            return "Tantron Electronic"
        }
        case 263: { /* '263' */
            return "Interra"
        }
        case 264: { /* '264' */
            return "DKX Tech"
        }
        case 265: { /* '265' */
            return "Viatron"
        }
        case 266: { /* '266' */
            return "Nautibus"
        }
        case 267: { /* '267' */
            return "ON Semiconductor"
        }
        case 268: { /* '268' */
            return "Longchuang"
        }
        case 269: { /* '269' */
            return "Air-On AG"
        }
        case 27: { /* '27' */
            return "Albert Ackermann"
        }
        case 270: { /* '270' */
            return "ib-company GmbH"
        }
        case 271: { /* '271' */
            return "Sation Factory"
        }
        case 272: { /* '272' */
            return "Agentilo GmbH"
        }
        case 273: { /* '273' */
            return "Makel Elektrik"
        }
        case 274: { /* '274' */
            return "Helios Ventilatoren"
        }
        case 275: { /* '275' */
            return "Otto Solutions Pte Ltd"
        }
        case 276: { /* '276' */
            return "Airmaster"
        }
        case 277: { /* '277' */
            return "Vallox GmbH"
        }
        case 278: { /* '278' */
            return "Dalitek"
        }
        case 279: { /* '279' */
            return "ASIN"
        }
        case 28: { /* '28' */
            return "Schupa GmbH"
        }
        case 280: { /* '280' */
            return "Bridges Intelligence Technology Inc."
        }
        case 281: { /* '281' */
            return "ARBONIA"
        }
        case 282: { /* '282' */
            return "KERMI"
        }
        case 283: { /* '283' */
            return "PROLUX"
        }
        case 284: { /* '284' */
            return "ClicHome"
        }
        case 285: { /* '285' */
            return "COMMAX"
        }
        case 286: { /* '286' */
            return "EAE"
        }
        case 287: { /* '287' */
            return "Tense"
        }
        case 288: { /* '288' */
            return "Seyoung Electronics"
        }
        case 289: { /* '289' */
            return "Lifedomus"
        }
        case 29: { /* '29' */
            return "ABB SCHWEIZ"
        }
        case 290: { /* '290' */
            return "EUROtronic Technology GmbH"
        }
        case 291: { /* '291' */
            return "tci"
        }
        case 292: { /* '292' */
            return "Rishun Electronic"
        }
        case 293: { /* '293' */
            return "Zipato"
        }
        case 294: { /* '294' */
            return "cm-security GmbH & Co KG"
        }
        case 295: { /* '295' */
            return "Qing Cables"
        }
        case 296: { /* '296' */
            return "LABIO"
        }
        case 297: { /* '297' */
            return "Coster Tecnologie Elettroniche S.p.A."
        }
        case 298: { /* '298' */
            return "E.G.E"
        }
        case 299: { /* '299' */
            return "NETxAutomation"
        }
        case 30: { /* '30' */
            return "Feller"
        }
        case 300: { /* '300' */
            return "tecalor"
        }
        case 301: { /* '301' */
            return "Urmet Electronics (Huizhou) Ltd."
        }
        case 302: { /* '302' */
            return "Peiying Building Control"
        }
        case 303: { /* '303' */
            return "BPT S.p.A. a Socio Unico"
        }
        case 304: { /* '304' */
            return "Kanontec - KanonBUS"
        }
        case 305: { /* '305' */
            return "ISER Tech"
        }
        case 306: { /* '306' */
            return "Fineline"
        }
        case 307: { /* '307' */
            return "CP Electronics Ltd"
        }
        case 308: { /* '308' */
            return "Niko-Servodan A/S"
        }
        case 309: { /* '309' */
            return "Simon"
        }
        case 31: { /* '31' */
            return "Glamox AS"
        }
        case 310: { /* '310' */
            return "GM modular pvt. Ltd."
        }
        case 311: { /* '311' */
            return "FU CHENG Intelligence"
        }
        case 312: { /* '312' */
            return "NexKon"
        }
        case 313: { /* '313' */
            return "FEEL s.r.l"
        }
        case 314: { /* '314' */
            return "Not Assigned"
        }
        case 315: { /* '315' */
            return "Shenzhen Fanhai Sanjiang Electronics Co., Ltd."
        }
        case 316: { /* '316' */
            return "Jiuzhou Greeble"
        }
        case 317: { /* '317' */
            return "Aumüller Aumatic GmbH"
        }
        case 318: { /* '318' */
            return "Etman Electric"
        }
        case 319: { /* '319' */
            return "Black Nova"
        }
        case 32: { /* '32' */
            return "DEHN & SÖHNE"
        }
        case 320: { /* '320' */
            return "ZidaTech AG"
        }
        case 321: { /* '321' */
            return "IDGS bvba"
        }
        case 322: { /* '322' */
            return "dakanimo"
        }
        case 323: { /* '323' */
            return "Trebor Automation AB"
        }
        case 324: { /* '324' */
            return "Satel sp. z o.o."
        }
        case 325: { /* '325' */
            return "Russound, Inc."
        }
        case 326: { /* '326' */
            return "Midea Heating & Ventilating Equipment CO LTD"
        }
        case 327: { /* '327' */
            return "Consorzio Terranuova"
        }
        case 328: { /* '328' */
            return "Wolf Heiztechnik GmbH"
        }
        case 329: { /* '329' */
            return "SONTEC"
        }
        case 33: { /* '33' */
            return "CRABTREE"
        }
        case 330: { /* '330' */
            return "Belcom Cables Ltd."
        }
        case 331: { /* '331' */
            return "Guangzhou SeaWin Electrical Technologies Co., Ltd."
        }
        case 332: { /* '332' */
            return "Acrel"
        }
        case 333: { /* '333' */
            return "Franke Aquarotter GmbH"
        }
        case 334: { /* '334' */
            return "Orion Systems"
        }
        case 335: { /* '335' */
            return "Schrack Technik GmbH"
        }
        case 336: { /* '336' */
            return "INSPRID"
        }
        case 337: { /* '337' */
            return "Sunricher"
        }
        case 338: { /* '338' */
            return "Menred automation system(shanghai) Co.,Ltd."
        }
        case 339: { /* '339' */
            return "Aurex"
        }
        case 34: { /* '34' */
            return "eVoKNX"
        }
        case 340: { /* '340' */
            return "Josef Barthelme GmbH & Co. KG"
        }
        case 341: { /* '341' */
            return "Architecture Numerique"
        }
        case 342: { /* '342' */
            return "UP GROUP"
        }
        case 343: { /* '343' */
            return "Teknos-Avinno"
        }
        case 344: { /* '344' */
            return "Ningbo Dooya Mechanic & Electronic Technology"
        }
        case 345: { /* '345' */
            return "Thermokon Sensortechnik GmbH"
        }
        case 346: { /* '346' */
            return "BELIMO Automation AG"
        }
        case 347: { /* '347' */
            return "Zehnder Group International AG"
        }
        case 348: { /* '348' */
            return "sks Kinkel Elektronik"
        }
        case 349: { /* '349' */
            return "ECE Wurmitzer GmbH"
        }
        case 350: { /* '350' */
            return "LARS"
        }
        case 351: { /* '351' */
            return "URC"
        }
        case 352: { /* '352' */
            return "LightControl"
        }
        case 353: { /* '353' */
            return "ShenZhen YM"
        }
        case 354: { /* '354' */
            return "MEAN WELL Enterprises Co. Ltd."
        }
        case 355: { /* '355' */
            return "OSix"
        }
        case 356: { /* '356' */
            return "AYPRO Technology"
        }
        case 357: { /* '357' */
            return "Hefei Ecolite Software"
        }
        case 358: { /* '358' */
            return "Enno"
        }
        case 359: { /* '359' */
            return "OHOSURE"
        }
        case 36: { /* '36' */
            return "Paul Hochköpper"
        }
        case 360: { /* '360' */
            return "Garefowl"
        }
        case 361: { /* '361' */
            return "GEZE"
        }
        case 362: { /* '362' */
            return "LG Electronics Inc."
        }
        case 363: { /* '363' */
            return "SMC interiors"
        }
        case 364: { /* '364' */
            return "Not Assigned"
        }
        case 365: { /* '365' */
            return "SCS Cable"
        }
        case 366: { /* '366' */
            return "Hoval"
        }
        case 367: { /* '367' */
            return "CANST"
        }
        case 368: { /* '368' */
            return "HangZhou Berlin"
        }
        case 369: { /* '369' */
            return "EVN-Lichttechnik"
        }
        case 37: { /* '37' */
            return "Altenburger Electronic"
        }
        case 370: { /* '370' */
            return "rutec"
        }
        case 371: { /* '371' */
            return "Finder"
        }
        case 372: { /* '372' */
            return "Fujitsu General Limited"
        }
        case 373: { /* '373' */
            return "ZF Friedrichshafen AG"
        }
        case 374: { /* '374' */
            return "Crealed"
        }
        case 375: { /* '375' */
            return "Miles Magic Automation Private Limited"
        }
        case 376: { /* '376' */
            return "E+"
        }
        case 377: { /* '377' */
            return "Italcond"
        }
        case 378: { /* '378' */
            return "SATION"
        }
        case 379: { /* '379' */
            return "NewBest"
        }
        case 380: { /* '380' */
            return "GDS DIGITAL SYSTEMS"
        }
        case 381: { /* '381' */
            return "Iddero"
        }
        case 382: { /* '382' */
            return "MBNLED"
        }
        case 383: { /* '383' */
            return "VITRUM"
        }
        case 384: { /* '384' */
            return "ekey biometric systems GmbH"
        }
        case 385: { /* '385' */
            return "AMC"
        }
        case 386: { /* '386' */
            return "TRILUX GmbH & Co. KG"
        }
        case 387: { /* '387' */
            return "WExcedo"
        }
        case 388: { /* '388' */
            return "VEMER SPA"
        }
        case 389: { /* '389' */
            return "Alexander Bürkle GmbH & Co KG"
        }
        case 390: { /* '390' */
            return "Citron"
        }
        case 391: { /* '391' */
            return "Shenzhen HeGuang"
        }
        case 392: { /* '392' */
            return "Not Assigned"
        }
        case 393: { /* '393' */
            return "TRANE B.V.B.A"
        }
        case 394: { /* '394' */
            return "CAREL"
        }
        case 395: { /* '395' */
            return "Prolite Controls"
        }
        case 396: { /* '396' */
            return "BOSMER"
        }
        case 397: { /* '397' */
            return "EUCHIPS"
        }
        case 398: { /* '398' */
            return "connect (Thinka connect)"
        }
        case 399: { /* '399' */
            return "PEAKnx a DOGAWIST company "
        }
        case 4: { /* '4' */
            return "Albrecht Jung"
        }
        case 400: { /* '400' */
            return "ACEMATIC"
        }
        case 401: { /* '401' */
            return "ELAUSYS"
        }
        case 402: { /* '402' */
            return "ITK Engineering AG"
        }
        case 403: { /* '403' */
            return "INTEGRA METERING AG"
        }
        case 404: { /* '404' */
            return "FMS Hospitality Pte Ltd"
        }
        case 405: { /* '405' */
            return "Nuvo"
        }
        case 406: { /* '406' */
            return "u::Lux GmbH"
        }
        case 407: { /* '407' */
            return "Brumberg Leuchten"
        }
        case 408: { /* '408' */
            return "Lime"
        }
        case 409: { /* '409' */
            return "Great Empire International Group Co., Ltd."
        }
        case 41: { /* '41' */
            return "Grässlin"
        }
        case 410: { /* '410' */
            return "Kavoshpishro Asia"
        }
        case 411: { /* '411' */
            return "V2 SpA"
        }
        case 412: { /* '412' */
            return "Johnson Controls"
        }
        case 413: { /* '413' */
            return "Arkud"
        }
        case 414: { /* '414' */
            return "Iridium Ltd."
        }
        case 415: { /* '415' */
            return "bsmart"
        }
        case 416: { /* '416' */
            return "BAB TECHNOLOGIE GmbH"
        }
        case 417: { /* '417' */
            return "NICE Spa"
        }
        case 418: { /* '418' */
            return "Redfish Group Pty Ltd"
        }
        case 419: { /* '419' */
            return "SABIANA spa"
        }
        case 42: { /* '42' */
            return "Simon"
        }
        case 420: { /* '420' */
            return "Ubee Interactive Europe"
        }
        case 421: { /* '421' */
            return "Rexel"
        }
        case 422: { /* '422' */
            return "Ges Teknik A.S."
        }
        case 423: { /* '423' */
            return "Ave S.p.A. "
        }
        case 424: { /* '424' */
            return "Zhuhai Ltech Technology Co., Ltd. "
        }
        case 425: { /* '425' */
            return "ARCOM"
        }
        case 426: { /* '426' */
            return "VIA Technologies, Inc."
        }
        case 427: { /* '427' */
            return "FEELSMART."
        }
        case 428: { /* '428' */
            return "SUPCON"
        }
        case 429: { /* '429' */
            return "MANIC"
        }
        case 430: { /* '430' */
            return "TDE GmbH"
        }
        case 431: { /* '431' */
            return "Nanjing Shufan Information technology Co.,Ltd."
        }
        case 432: { /* '432' */
            return "EWTech"
        }
        case 433: { /* '433' */
            return "Kluger Automation GmbH"
        }
        case 434: { /* '434' */
            return "JoongAng Control"
        }
        case 435: { /* '435' */
            return "GreenControls Technology Sdn. Bhd."
        }
        case 436: { /* '436' */
            return "IME S.p.a."
        }
        case 437: { /* '437' */
            return "SiChuan HaoDing"
        }
        case 438: { /* '438' */
            return "Mindjaga Ltd."
        }
        case 439: { /* '439' */
            return "RuiLi Smart Control"
        }
        case 43954: { /* '43954' */
            return "ABB - reserved"
        }
        case 43959: { /* '43959' */
            return "Busch-Jaeger Elektro - reserved"
        }
        case 44: { /* '44' */
            return "VIMAR"
        }
        case 440: { /* '440' */
            return "CODESYS GmbH"
        }
        case 441: { /* '441' */
            return "Moorgen Deutschland GmbH"
        }
        case 442: { /* '442' */
            return "CULLMANN TECH"
        }
        case 443: { /* '443' */
            return "Merck Window Technologies B.V. "
        }
        case 444: { /* '444' */
            return "ABEGO"
        }
        case 445: { /* '445' */
            return "myGEKKO"
        }
        case 446: { /* '446' */
            return "Ergo3 Sarl"
        }
        case 447: { /* '447' */
            return "STmicroelectronics International N.V."
        }
        case 448: { /* '448' */
            return "cjc systems"
        }
        case 449: { /* '449' */
            return "Sudoku"
        }
        case 45: { /* '45' */
            return "Moeller Gebäudeautomation KG"
        }
        case 451: { /* '451' */
            return "AZ e-lite Pte Ltd"
        }
        case 452: { /* '452' */
            return "Arlight"
        }
        case 453: { /* '453' */
            return "Grünbeck Wasseraufbereitung GmbH"
        }
        case 454: { /* '454' */
            return "Module Electronic"
        }
        case 455: { /* '455' */
            return "KOPLAT"
        }
        case 456: { /* '456' */
            return "Guangzhou Letour Life Technology Co., Ltd"
        }
        case 457: { /* '457' */
            return "ILEVIA"
        }
        case 458: { /* '458' */
            return "LN SYSTEMTEQ"
        }
        case 459: { /* '459' */
            return "Hisense SmartHome"
        }
        case 46: { /* '46' */
            return "Eltako"
        }
        case 460: { /* '460' */
            return "Flink Automation System"
        }
        case 461: { /* '461' */
            return "xxter bv"
        }
        case 462: { /* '462' */
            return "lynxus technology"
        }
        case 463: { /* '463' */
            return "ROBOT S.A."
        }
        case 464: { /* '464' */
            return "Shenzhen Atte Smart Life Co.,Ltd."
        }
        case 465: { /* '465' */
            return "Noblesse"
        }
        case 466: { /* '466' */
            return "Advanced Devices"
        }
        case 467: { /* '467' */
            return "Atrina Building Automation Co. Ltd"
        }
        case 468: { /* '468' */
            return "Guangdong Daming Laffey electric Co., Ltd."
        }
        case 469: { /* '469' */
            return "Westerstrand Urfabrik AB"
        }
        case 470: { /* '470' */
            return "Control4 Corporate"
        }
        case 471: { /* '471' */
            return "Ontrol"
        }
        case 472: { /* '472' */
            return "Starnet"
        }
        case 473: { /* '473' */
            return "BETA CAVI"
        }
        case 474: { /* '474' */
            return "EaseMore"
        }
        case 475: { /* '475' */
            return "Vivaldi srl"
        }
        case 476: { /* '476' */
            return "Gree Electric Appliances,Inc. of Zhuhai"
        }
        case 477: { /* '477' */
            return "HWISCON"
        }
        case 478: { /* '478' */
            return "Shanghai ELECON Intelligent Technology Co., Ltd."
        }
        case 479: { /* '479' */
            return "Kampmann"
        }
        case 480: { /* '480' */
            return "Impolux GmbH / LEDIMAX"
        }
        case 481: { /* '481' */
            return "Evaux"
        }
        case 482: { /* '482' */
            return "Webro Cables & Connectors Limited"
        }
        case 483: { /* '483' */
            return "Shanghai E-tech Solution"
        }
        case 484: { /* '484' */
            return "Guangzhou HOKO Electric Co.,Ltd."
        }
        case 485: { /* '485' */
            return "LAMMIN HIGH TECH CO.,LTD"
        }
        case 486: { /* '486' */
            return "Shenzhen Merrytek Technology Co., Ltd"
        }
        case 487: { /* '487' */
            return "I-Luxus"
        }
        case 488: { /* '488' */
            return "Elmos Semiconductor AG"
        }
        case 489: { /* '489' */
            return "EmCom Technology Inc"
        }
        case 49: { /* '49' */
            return "Bosch-Siemens Haushaltsgeräte"
        }
        case 490: { /* '490' */
            return "project innovations GmbH"
        }
        case 491: { /* '491' */
            return "Itc"
        }
        case 492: { /* '492' */
            return "ABB LV Installation Materials Company Ltd, Beijing"
        }
        case 493: { /* '493' */
            return "Maico "
        }
        case 495: { /* '495' */
            return "ELAN SRL"
        }
        case 496: { /* '496' */
            return "MinhHa Technology co.,Ltd"
        }
        case 497: { /* '497' */
            return "Zhejiang Tianjie Industrial CORP."
        }
        case 498: { /* '498' */
            return "iAutomation Pty Limited"
        }
        case 499: { /* '499' */
            return "Extron"
        }
        case 5: { /* '5' */
            return "Bticino"
        }
        case 500: { /* '500' */
            return "Freedompro"
        }
        case 501: { /* '501' */
            return "1Home"
        }
        case 502: { /* '502' */
            return "EOS Saunatechnik GmbH"
        }
        case 503: { /* '503' */
            return "KUSATEK GmbH"
        }
        case 504: { /* '504' */
            return "EisBär Scada"
        }
        case 505: { /* '505' */
            return "AUTOMATISMI BENINCA S.P.A."
        }
        case 506: { /* '506' */
            return "Blendom"
        }
        case 507: { /* '507' */
            return "Madel Air Technical diffusion"
        }
        case 508: { /* '508' */
            return "NIKO"
        }
        case 509: { /* '509' */
            return "Bosch Rexroth AG"
        }
        case 512: { /* '512' */
            return "C&M Products"
        }
        case 513: { /* '513' */
            return "Hörmann KG Verkaufsgesellschaft"
        }
        case 514: { /* '514' */
            return "Shanghai Rajayasa co.,LTD"
        }
        case 515: { /* '515' */
            return "SUZUKI"
        }
        case 516: { /* '516' */
            return "Silent Gliss International Ltd."
        }
        case 517: { /* '517' */
            return "BEE Controls (ADGSC Group)"
        }
        case 518: { /* '518' */
            return "xDTecGmbH"
        }
        case 519: { /* '519' */
            return "OSRAM"
        }
        case 52: { /* '52' */
            return "RITTO GmbH&Co.KG"
        }
        case 520: { /* '520' */
            return "Lebenor"
        }
        case 521: { /* '521' */
            return "automaneng"
        }
        case 522: { /* '522' */
            return "Honeywell Automation Solution control(China)"
        }
        case 523: { /* '523' */
            return "Hangzhou binthen Intelligence Technology Co.,Ltd"
        }
        case 524: { /* '524' */
            return "ETA Heiztechnik"
        }
        case 525: { /* '525' */
            return "DIVUS GmbH"
        }
        case 526: { /* '526' */
            return "Nanjing Taijiesai Intelligent Technology Co. Ltd."
        }
        case 527: { /* '527' */
            return "Lunatone"
        }
        case 528: { /* '528' */
            return "ZHEJIANG SCTECH BUILDING INTELLIGENT"
        }
        case 529: { /* '529' */
            return "Foshan Qite Technology Co., Ltd."
        }
        case 53: { /* '53' */
            return "Power Controls"
        }
        case 530: { /* '530' */
            return "NOKE"
        }
        case 531: { /* '531' */
            return "LANDCOM"
        }
        case 532: { /* '532' */
            return "Stork AS"
        }
        case 533: { /* '533' */
            return "Hangzhou Shendu Technology Co., Ltd."
        }
        case 534: { /* '534' */
            return "CoolAutomation"
        }
        case 535: { /* '535' */
            return "Aprstern"
        }
        case 536: { /* '536' */
            return "sonnen"
        }
        case 537: { /* '537' */
            return "DNAKE"
        }
        case 538: { /* '538' */
            return "Neuberger Gebäudeautomation GmbH"
        }
        case 539: { /* '539' */
            return "Stiliger"
        }
        case 540: { /* '540' */
            return "Berghof Automation GmbH"
        }
        case 541: { /* '541' */
            return "Total Automation and controls GmbH"
        }
        case 542: { /* '542' */
            return "dovit"
        }
        case 543: { /* '543' */
            return "Instalighting GmbH"
        }
        case 544: { /* '544' */
            return "UNI-TEC"
        }
        case 545: { /* '545' */
            return "CasaTunes"
        }
        case 546: { /* '546' */
            return "EMT"
        }
        case 547: { /* '547' */
            return "Senfficient"
        }
        case 548: { /* '548' */
            return "Aurolite electrical panyu guangzhou limited"
        }
        case 549: { /* '549' */
            return "ABB Xiamen Smart Technology Co., Ltd."
        }
        case 55: { /* '55' */
            return "ZUMTOBEL"
        }
        case 550: { /* '550' */
            return "Samson Electric Wire"
        }
        case 551: { /* '551' */
            return "T-Touching"
        }
        case 552: { /* '552' */
            return "Core Smart Home"
        }
        case 553: { /* '553' */
            return "GreenConnect Solutions SA"
        }
        case 554: { /* '554' */
            return "ELETTRONICA CONDUTTORI"
        }
        case 555: { /* '555' */
            return "MKFC"
        }
        case 556: { /* '556' */
            return "Automation+"
        }
        case 557: { /* '557' */
            return "blue and red"
        }
        case 558: { /* '558' */
            return "frogblue"
        }
        case 559: { /* '559' */
            return "SAVESOR"
        }
        case 560: { /* '560' */
            return "App Tech"
        }
        case 561: { /* '561' */
            return "sensortec AG"
        }
        case 562: { /* '562' */
            return "nysa technology & solutions"
        }
        case 563: { /* '563' */
            return "FARADITE"
        }
        case 564: { /* '564' */
            return "Optimus"
        }
        case 565: { /* '565' */
            return "KTS s.r.l."
        }
        case 566: { /* '566' */
            return "Ramcro SPA"
        }
        case 567: { /* '567' */
            return "Wuhan WiseCreate Universe Technology Co., Ltd"
        }
        case 568: { /* '568' */
            return "BEMI Smart Home Ltd"
        }
        case 569: { /* '569' */
            return "Ardomus"
        }
        case 57: { /* '57' */
            return "Phoenix Contact"
        }
        case 570: { /* '570' */
            return "ChangXing"
        }
        case 571: { /* '571' */
            return "E-Controls"
        }
        case 572: { /* '572' */
            return "AIB Technology"
        }
        case 573: { /* '573' */
            return "NVC"
        }
        case 574: { /* '574' */
            return "Kbox"
        }
        case 575: { /* '575' */
            return "CNS"
        }
        case 576: { /* '576' */
            return "Tyba"
        }
        case 577: { /* '577' */
            return "Atrel"
        }
        case 578: { /* '578' */
            return "Simon Electric (China) Co., LTD"
        }
        case 579: { /* '579' */
            return "Kordz Group"
        }
        case 580: { /* '580' */
            return "ND Electric"
        }
        case 581: { /* '581' */
            return "Controlium"
        }
        case 582: { /* '582' */
            return "FAMO GmbH & Co. KG"
        }
        case 583: { /* '583' */
            return "CDN Smart"
        }
        case 584: { /* '584' */
            return "Heston"
        }
        case 585: { /* '585' */
            return "ESLA CONEXIONES S.L."
        }
        case 586: { /* '586' */
            return "Weishaupt"
        }
        case 587: { /* '587' */
            return "ASTRUM TECHNOLOGY"
        }
        case 588: { /* '588' */
            return "WUERTH ELEKTRONIK STELVIO KONTEK S.p.A."
        }
        case 589: { /* '589' */
            return "NANOTECO corporation"
        }
        case 590: { /* '590' */
            return "Nietian"
        }
        case 591: { /* '591' */
            return "Sumsir"
        }
        case 592: { /* '592' */
            return "ORBIS TECNOLOGIA ELECTRICA SA"
        }
        case 6: { /* '6' */
            return "Berker"
        }
        case 61: { /* '61' */
            return "WAGO Kontakttechnik"
        }
        case 62: { /* '62' */
            return "knXpresso"
        }
        case 66: { /* '66' */
            return "Wieland Electric"
        }
        case 67: { /* '67' */
            return "Hermann Kleinhuis"
        }
        case 69: { /* '69' */
            return "Stiebel Eltron"
        }
        case 7: { /* '7' */
            return "Busch-Jaeger Elektro"
        }
        case 71: { /* '71' */
            return "Tehalit"
        }
        case 72: { /* '72' */
            return "Theben AG"
        }
        case 73: { /* '73' */
            return "Wilhelm Rutenbeck"
        }
        case 75: { /* '75' */
            return "Winkhaus"
        }
        case 76: { /* '76' */
            return "Robert Bosch"
        }
        case 78: { /* '78' */
            return "Somfy"
        }
        case 8: { /* '8' */
            return "GIRA Giersiepen"
        }
        case 80: { /* '80' */
            return "Woertz"
        }
        case 81: { /* '81' */
            return "Viessmann Werke"
        }
        case 82: { /* '82' */
            return "IMI Hydronic Engineering"
        }
        case 83: { /* '83' */
            return "Joh. Vaillant"
        }
        case 85: { /* '85' */
            return "AMP Deutschland"
        }
        case 89: { /* '89' */
            return "Bosch Thermotechnik GmbH"
        }
        case 9: { /* '9' */
            return "Hager Electro"
        }
        case 90: { /* '90' */
            return "SEF - ECOTEC"
        }
        case 92: { /* '92' */
            return "DORMA GmbH + Co. KG"
        }
        case 93: { /* '93' */
            return "WindowMaster A/S"
        }
        case 94: { /* '94' */
            return "Walther Werke"
        }
        case 95: { /* '95' */
            return "ORAS"
        }
        case 97: { /* '97' */
            return "Dätwyler"
        }
        case 98: { /* '98' */
            return "Electrak"
        }
        case 99: { /* '99' */
            return "Techem"
        }
        default: {
            return ""
        }
    }
}
func KnxManufacturersValueOf(value uint16) KnxManufacturers {
    switch value {
        case 1:
            return KnxManufacturers_M_0001
        case 10:
            return KnxManufacturers_M_000A
        case 100:
            return KnxManufacturers_M_0064
        case 101:
            return KnxManufacturers_M_0065
        case 102:
            return KnxManufacturers_M_0066
        case 104:
            return KnxManufacturers_M_0068
        case 105:
            return KnxManufacturers_M_0069
        case 106:
            return KnxManufacturers_M_006A
        case 107:
            return KnxManufacturers_M_006B
        case 108:
            return KnxManufacturers_M_006C
        case 109:
            return KnxManufacturers_M_006D
        case 11:
            return KnxManufacturers_M_000B
        case 110:
            return KnxManufacturers_M_006E
        case 111:
            return KnxManufacturers_M_006F
        case 112:
            return KnxManufacturers_M_0070
        case 113:
            return KnxManufacturers_M_0071
        case 114:
            return KnxManufacturers_M_0072
        case 115:
            return KnxManufacturers_M_0073
        case 116:
            return KnxManufacturers_M_0074
        case 117:
            return KnxManufacturers_M_0075
        case 118:
            return KnxManufacturers_M_0076
        case 119:
            return KnxManufacturers_M_0077
        case 12:
            return KnxManufacturers_M_000C
        case 120:
            return KnxManufacturers_M_0078
        case 121:
            return KnxManufacturers_M_0079
        case 122:
            return KnxManufacturers_M_007A
        case 123:
            return KnxManufacturers_M_007B
        case 124:
            return KnxManufacturers_M_007C
        case 125:
            return KnxManufacturers_M_007D
        case 126:
            return KnxManufacturers_M_007E
        case 127:
            return KnxManufacturers_M_007F
        case 128:
            return KnxManufacturers_M_0080
        case 129:
            return KnxManufacturers_M_0081
        case 130:
            return KnxManufacturers_M_0082
        case 131:
            return KnxManufacturers_M_0083
        case 132:
            return KnxManufacturers_M_0084
        case 133:
            return KnxManufacturers_M_0085
        case 134:
            return KnxManufacturers_M_0086
        case 135:
            return KnxManufacturers_M_0087
        case 136:
            return KnxManufacturers_M_0088
        case 137:
            return KnxManufacturers_M_0089
        case 138:
            return KnxManufacturers_M_008A
        case 139:
            return KnxManufacturers_M_008B
        case 14:
            return KnxManufacturers_M_000E
        case 140:
            return KnxManufacturers_M_008C
        case 141:
            return KnxManufacturers_M_008D
        case 142:
            return KnxManufacturers_M_008E
        case 143:
            return KnxManufacturers_M_008F
        case 144:
            return KnxManufacturers_M_0090
        case 145:
            return KnxManufacturers_M_0091
        case 146:
            return KnxManufacturers_M_0092
        case 147:
            return KnxManufacturers_M_0093
        case 148:
            return KnxManufacturers_M_0094
        case 149:
            return KnxManufacturers_M_0095
        case 150:
            return KnxManufacturers_M_0096
        case 151:
            return KnxManufacturers_M_0097
        case 152:
            return KnxManufacturers_M_0098
        case 153:
            return KnxManufacturers_M_0099
        case 154:
            return KnxManufacturers_M_009A
        case 155:
            return KnxManufacturers_M_009B
        case 156:
            return KnxManufacturers_M_009C
        case 157:
            return KnxManufacturers_M_009D
        case 158:
            return KnxManufacturers_M_009E
        case 159:
            return KnxManufacturers_M_009F
        case 160:
            return KnxManufacturers_M_00A0
        case 161:
            return KnxManufacturers_M_00A1
        case 162:
            return KnxManufacturers_M_00A2
        case 163:
            return KnxManufacturers_M_00A3
        case 164:
            return KnxManufacturers_M_00A4
        case 165:
            return KnxManufacturers_M_00A5
        case 166:
            return KnxManufacturers_M_00A6
        case 167:
            return KnxManufacturers_M_00A7
        case 168:
            return KnxManufacturers_M_00A8
        case 169:
            return KnxManufacturers_M_00A9
        case 170:
            return KnxManufacturers_M_00AA
        case 171:
            return KnxManufacturers_M_00AB
        case 172:
            return KnxManufacturers_M_00AC
        case 173:
            return KnxManufacturers_M_00AD
        case 174:
            return KnxManufacturers_M_00AE
        case 175:
            return KnxManufacturers_M_00AF
        case 176:
            return KnxManufacturers_M_00B0
        case 177:
            return KnxManufacturers_M_00B1
        case 178:
            return KnxManufacturers_M_00B2
        case 179:
            return KnxManufacturers_M_00B3
        case 180:
            return KnxManufacturers_M_00B4
        case 181:
            return KnxManufacturers_M_00B5
        case 182:
            return KnxManufacturers_M_00B6
        case 183:
            return KnxManufacturers_M_00B7
        case 184:
            return KnxManufacturers_M_00B8
        case 185:
            return KnxManufacturers_M_00B9
        case 186:
            return KnxManufacturers_M_00BA
        case 187:
            return KnxManufacturers_M_00BB
        case 188:
            return KnxManufacturers_M_00BC
        case 189:
            return KnxManufacturers_M_00BD
        case 190:
            return KnxManufacturers_M_00BE
        case 191:
            return KnxManufacturers_M_00BF
        case 192:
            return KnxManufacturers_M_00C0
        case 193:
            return KnxManufacturers_M_00C1
        case 194:
            return KnxManufacturers_M_00C2
        case 195:
            return KnxManufacturers_M_00C3
        case 196:
            return KnxManufacturers_M_00C4
        case 197:
            return KnxManufacturers_M_00C5
        case 198:
            return KnxManufacturers_M_00C6
        case 199:
            return KnxManufacturers_M_00C7
        case 2:
            return KnxManufacturers_M_0002
        case 200:
            return KnxManufacturers_M_00C8
        case 201:
            return KnxManufacturers_M_00C9
        case 202:
            return KnxManufacturers_M_00CA
        case 204:
            return KnxManufacturers_M_00CC
        case 205:
            return KnxManufacturers_M_00CD
        case 206:
            return KnxManufacturers_M_00CE
        case 207:
            return KnxManufacturers_M_00CF
        case 208:
            return KnxManufacturers_M_00D0
        case 209:
            return KnxManufacturers_M_00D1
        case 210:
            return KnxManufacturers_M_00D2
        case 211:
            return KnxManufacturers_M_00D3
        case 214:
            return KnxManufacturers_M_00D6
        case 215:
            return KnxManufacturers_M_00D7
        case 216:
            return KnxManufacturers_M_00D8
        case 217:
            return KnxManufacturers_M_00D9
        case 218:
            return KnxManufacturers_M_00DA
        case 219:
            return KnxManufacturers_M_00DB
        case 22:
            return KnxManufacturers_M_0016
        case 220:
            return KnxManufacturers_M_00DC
        case 222:
            return KnxManufacturers_M_00DE
        case 223:
            return KnxManufacturers_M_00DF
        case 225:
            return KnxManufacturers_M_00E1
        case 227:
            return KnxManufacturers_M_00E3
        case 228:
            return KnxManufacturers_M_00E4
        case 232:
            return KnxManufacturers_M_00E8
        case 233:
            return KnxManufacturers_M_00E9
        case 234:
            return KnxManufacturers_M_00EA
        case 235:
            return KnxManufacturers_M_00EB
        case 237:
            return KnxManufacturers_M_00ED
        case 238:
            return KnxManufacturers_M_00EE
        case 239:
            return KnxManufacturers_M_00EF
        case 24:
            return KnxManufacturers_M_0018
        case 240:
            return KnxManufacturers_M_00F0
        case 241:
            return KnxManufacturers_M_00F1
        case 242:
            return KnxManufacturers_M_00F2
        case 244:
            return KnxManufacturers_M_00F4
        case 245:
            return KnxManufacturers_M_00F5
        case 246:
            return KnxManufacturers_M_00F6
        case 248:
            return KnxManufacturers_M_00F8
        case 249:
            return KnxManufacturers_M_00F9
        case 25:
            return KnxManufacturers_M_0019
        case 250:
            return KnxManufacturers_M_00FA
        case 251:
            return KnxManufacturers_M_00FB
        case 252:
            return KnxManufacturers_M_00FC
        case 253:
            return KnxManufacturers_M_00FD
        case 254:
            return KnxManufacturers_M_00FE
        case 256:
            return KnxManufacturers_M_0100
        case 257:
            return KnxManufacturers_M_0101
        case 258:
            return KnxManufacturers_M_0102
        case 259:
            return KnxManufacturers_M_0103
        case 260:
            return KnxManufacturers_M_0104
        case 261:
            return KnxManufacturers_M_0105
        case 262:
            return KnxManufacturers_M_0106
        case 263:
            return KnxManufacturers_M_0107
        case 264:
            return KnxManufacturers_M_0108
        case 265:
            return KnxManufacturers_M_0109
        case 266:
            return KnxManufacturers_M_010A
        case 267:
            return KnxManufacturers_M_010B
        case 268:
            return KnxManufacturers_M_010C
        case 269:
            return KnxManufacturers_M_010D
        case 27:
            return KnxManufacturers_M_001B
        case 270:
            return KnxManufacturers_M_010E
        case 271:
            return KnxManufacturers_M_010F
        case 272:
            return KnxManufacturers_M_0110
        case 273:
            return KnxManufacturers_M_0111
        case 274:
            return KnxManufacturers_M_0112
        case 275:
            return KnxManufacturers_M_0113
        case 276:
            return KnxManufacturers_M_0114
        case 277:
            return KnxManufacturers_M_0115
        case 278:
            return KnxManufacturers_M_0116
        case 279:
            return KnxManufacturers_M_0117
        case 28:
            return KnxManufacturers_M_001C
        case 280:
            return KnxManufacturers_M_0118
        case 281:
            return KnxManufacturers_M_0119
        case 282:
            return KnxManufacturers_M_011A
        case 283:
            return KnxManufacturers_M_011B
        case 284:
            return KnxManufacturers_M_011C
        case 285:
            return KnxManufacturers_M_011D
        case 286:
            return KnxManufacturers_M_011E
        case 287:
            return KnxManufacturers_M_011F
        case 288:
            return KnxManufacturers_M_0120
        case 289:
            return KnxManufacturers_M_0121
        case 29:
            return KnxManufacturers_M_001D
        case 290:
            return KnxManufacturers_M_0122
        case 291:
            return KnxManufacturers_M_0123
        case 292:
            return KnxManufacturers_M_0124
        case 293:
            return KnxManufacturers_M_0125
        case 294:
            return KnxManufacturers_M_0126
        case 295:
            return KnxManufacturers_M_0127
        case 296:
            return KnxManufacturers_M_0128
        case 297:
            return KnxManufacturers_M_0129
        case 298:
            return KnxManufacturers_M_012A
        case 299:
            return KnxManufacturers_M_012B
        case 30:
            return KnxManufacturers_M_001E
        case 300:
            return KnxManufacturers_M_012C
        case 301:
            return KnxManufacturers_M_012D
        case 302:
            return KnxManufacturers_M_012E
        case 303:
            return KnxManufacturers_M_012F
        case 304:
            return KnxManufacturers_M_0130
        case 305:
            return KnxManufacturers_M_0131
        case 306:
            return KnxManufacturers_M_0132
        case 307:
            return KnxManufacturers_M_0133
        case 308:
            return KnxManufacturers_M_0134
        case 309:
            return KnxManufacturers_M_0135
        case 31:
            return KnxManufacturers_M_001F
        case 310:
            return KnxManufacturers_M_0136
        case 311:
            return KnxManufacturers_M_0137
        case 312:
            return KnxManufacturers_M_0138
        case 313:
            return KnxManufacturers_M_0139
        case 314:
            return KnxManufacturers_M_013A
        case 315:
            return KnxManufacturers_M_013B
        case 316:
            return KnxManufacturers_M_013C
        case 317:
            return KnxManufacturers_M_013D
        case 318:
            return KnxManufacturers_M_013E
        case 319:
            return KnxManufacturers_M_013F
        case 32:
            return KnxManufacturers_M_0020
        case 320:
            return KnxManufacturers_M_0140
        case 321:
            return KnxManufacturers_M_0141
        case 322:
            return KnxManufacturers_M_0142
        case 323:
            return KnxManufacturers_M_0143
        case 324:
            return KnxManufacturers_M_0144
        case 325:
            return KnxManufacturers_M_0145
        case 326:
            return KnxManufacturers_M_0146
        case 327:
            return KnxManufacturers_M_0147
        case 328:
            return KnxManufacturers_M_0148
        case 329:
            return KnxManufacturers_M_0149
        case 33:
            return KnxManufacturers_M_0021
        case 330:
            return KnxManufacturers_M_014A
        case 331:
            return KnxManufacturers_M_014B
        case 332:
            return KnxManufacturers_M_014C
        case 333:
            return KnxManufacturers_M_014D
        case 334:
            return KnxManufacturers_M_014E
        case 335:
            return KnxManufacturers_M_014F
        case 336:
            return KnxManufacturers_M_0150
        case 337:
            return KnxManufacturers_M_0151
        case 338:
            return KnxManufacturers_M_0152
        case 339:
            return KnxManufacturers_M_0153
        case 34:
            return KnxManufacturers_M_0022
        case 340:
            return KnxManufacturers_M_0154
        case 341:
            return KnxManufacturers_M_0155
        case 342:
            return KnxManufacturers_M_0156
        case 343:
            return KnxManufacturers_M_0157
        case 344:
            return KnxManufacturers_M_0158
        case 345:
            return KnxManufacturers_M_0159
        case 346:
            return KnxManufacturers_M_015A
        case 347:
            return KnxManufacturers_M_015B
        case 348:
            return KnxManufacturers_M_015C
        case 349:
            return KnxManufacturers_M_015D
        case 350:
            return KnxManufacturers_M_015E
        case 351:
            return KnxManufacturers_M_015F
        case 352:
            return KnxManufacturers_M_0160
        case 353:
            return KnxManufacturers_M_0161
        case 354:
            return KnxManufacturers_M_0162
        case 355:
            return KnxManufacturers_M_0163
        case 356:
            return KnxManufacturers_M_0164
        case 357:
            return KnxManufacturers_M_0165
        case 358:
            return KnxManufacturers_M_0166
        case 359:
            return KnxManufacturers_M_0167
        case 36:
            return KnxManufacturers_M_0024
        case 360:
            return KnxManufacturers_M_0168
        case 361:
            return KnxManufacturers_M_0169
        case 362:
            return KnxManufacturers_M_016A
        case 363:
            return KnxManufacturers_M_016B
        case 364:
            return KnxManufacturers_M_016C
        case 365:
            return KnxManufacturers_M_016D
        case 366:
            return KnxManufacturers_M_016E
        case 367:
            return KnxManufacturers_M_016F
        case 368:
            return KnxManufacturers_M_0170
        case 369:
            return KnxManufacturers_M_0171
        case 37:
            return KnxManufacturers_M_0025
        case 370:
            return KnxManufacturers_M_0172
        case 371:
            return KnxManufacturers_M_0173
        case 372:
            return KnxManufacturers_M_0174
        case 373:
            return KnxManufacturers_M_0175
        case 374:
            return KnxManufacturers_M_0176
        case 375:
            return KnxManufacturers_M_0177
        case 376:
            return KnxManufacturers_M_0178
        case 377:
            return KnxManufacturers_M_0179
        case 378:
            return KnxManufacturers_M_017A
        case 379:
            return KnxManufacturers_M_017B
        case 380:
            return KnxManufacturers_M_017C
        case 381:
            return KnxManufacturers_M_017D
        case 382:
            return KnxManufacturers_M_017E
        case 383:
            return KnxManufacturers_M_017F
        case 384:
            return KnxManufacturers_M_0180
        case 385:
            return KnxManufacturers_M_0181
        case 386:
            return KnxManufacturers_M_0182
        case 387:
            return KnxManufacturers_M_0183
        case 388:
            return KnxManufacturers_M_0184
        case 389:
            return KnxManufacturers_M_0185
        case 390:
            return KnxManufacturers_M_0186
        case 391:
            return KnxManufacturers_M_0187
        case 392:
            return KnxManufacturers_M_0188
        case 393:
            return KnxManufacturers_M_0189
        case 394:
            return KnxManufacturers_M_018A
        case 395:
            return KnxManufacturers_M_018B
        case 396:
            return KnxManufacturers_M_018C
        case 397:
            return KnxManufacturers_M_018D
        case 398:
            return KnxManufacturers_M_018E
        case 399:
            return KnxManufacturers_M_018F
        case 4:
            return KnxManufacturers_M_0004
        case 400:
            return KnxManufacturers_M_0190
        case 401:
            return KnxManufacturers_M_0191
        case 402:
            return KnxManufacturers_M_0192
        case 403:
            return KnxManufacturers_M_0193
        case 404:
            return KnxManufacturers_M_0194
        case 405:
            return KnxManufacturers_M_0195
        case 406:
            return KnxManufacturers_M_0196
        case 407:
            return KnxManufacturers_M_0197
        case 408:
            return KnxManufacturers_M_0198
        case 409:
            return KnxManufacturers_M_0199
        case 41:
            return KnxManufacturers_M_0029
        case 410:
            return KnxManufacturers_M_019A
        case 411:
            return KnxManufacturers_M_019B
        case 412:
            return KnxManufacturers_M_019C
        case 413:
            return KnxManufacturers_M_019D
        case 414:
            return KnxManufacturers_M_019E
        case 415:
            return KnxManufacturers_M_019F
        case 416:
            return KnxManufacturers_M_01A0
        case 417:
            return KnxManufacturers_M_01A1
        case 418:
            return KnxManufacturers_M_01A2
        case 419:
            return KnxManufacturers_M_01A3
        case 42:
            return KnxManufacturers_M_002A
        case 420:
            return KnxManufacturers_M_01A4
        case 421:
            return KnxManufacturers_M_01A5
        case 422:
            return KnxManufacturers_M_01A6
        case 423:
            return KnxManufacturers_M_01A7
        case 424:
            return KnxManufacturers_M_01A8
        case 425:
            return KnxManufacturers_M_01A9
        case 426:
            return KnxManufacturers_M_01AA
        case 427:
            return KnxManufacturers_M_01AB
        case 428:
            return KnxManufacturers_M_01AC
        case 429:
            return KnxManufacturers_M_01AD
        case 430:
            return KnxManufacturers_M_01AE
        case 431:
            return KnxManufacturers_M_01AF
        case 432:
            return KnxManufacturers_M_01B0
        case 433:
            return KnxManufacturers_M_01B1
        case 434:
            return KnxManufacturers_M_01B2
        case 435:
            return KnxManufacturers_M_01B3
        case 436:
            return KnxManufacturers_M_01B4
        case 437:
            return KnxManufacturers_M_01B5
        case 438:
            return KnxManufacturers_M_01B6
        case 439:
            return KnxManufacturers_M_01B7
        case 43954:
            return KnxManufacturers_M_ABB2
        case 43959:
            return KnxManufacturers_M_ABB7
        case 44:
            return KnxManufacturers_M_002C
        case 440:
            return KnxManufacturers_M_01B8
        case 441:
            return KnxManufacturers_M_01B9
        case 442:
            return KnxManufacturers_M_01BA
        case 443:
            return KnxManufacturers_M_01BB
        case 444:
            return KnxManufacturers_M_01BC
        case 445:
            return KnxManufacturers_M_01BD
        case 446:
            return KnxManufacturers_M_01BE
        case 447:
            return KnxManufacturers_M_01BF
        case 448:
            return KnxManufacturers_M_01C0
        case 449:
            return KnxManufacturers_M_01C1
        case 45:
            return KnxManufacturers_M_002D
        case 451:
            return KnxManufacturers_M_01C3
        case 452:
            return KnxManufacturers_M_01C4
        case 453:
            return KnxManufacturers_M_01C5
        case 454:
            return KnxManufacturers_M_01C6
        case 455:
            return KnxManufacturers_M_01C7
        case 456:
            return KnxManufacturers_M_01C8
        case 457:
            return KnxManufacturers_M_01C9
        case 458:
            return KnxManufacturers_M_01CA
        case 459:
            return KnxManufacturers_M_01CB
        case 46:
            return KnxManufacturers_M_002E
        case 460:
            return KnxManufacturers_M_01CC
        case 461:
            return KnxManufacturers_M_01CD
        case 462:
            return KnxManufacturers_M_01CE
        case 463:
            return KnxManufacturers_M_01CF
        case 464:
            return KnxManufacturers_M_01D0
        case 465:
            return KnxManufacturers_M_01D1
        case 466:
            return KnxManufacturers_M_01D2
        case 467:
            return KnxManufacturers_M_01D3
        case 468:
            return KnxManufacturers_M_01D4
        case 469:
            return KnxManufacturers_M_01D5
        case 470:
            return KnxManufacturers_M_01D6
        case 471:
            return KnxManufacturers_M_01D7
        case 472:
            return KnxManufacturers_M_01D8
        case 473:
            return KnxManufacturers_M_01D9
        case 474:
            return KnxManufacturers_M_01DA
        case 475:
            return KnxManufacturers_M_01DB
        case 476:
            return KnxManufacturers_M_01DC
        case 477:
            return KnxManufacturers_M_01DD
        case 478:
            return KnxManufacturers_M_01DE
        case 479:
            return KnxManufacturers_M_01DF
        case 480:
            return KnxManufacturers_M_01E0
        case 481:
            return KnxManufacturers_M_01E1
        case 482:
            return KnxManufacturers_M_01E2
        case 483:
            return KnxManufacturers_M_01E3
        case 484:
            return KnxManufacturers_M_01E4
        case 485:
            return KnxManufacturers_M_01E5
        case 486:
            return KnxManufacturers_M_01E6
        case 487:
            return KnxManufacturers_M_01E7
        case 488:
            return KnxManufacturers_M_01E8
        case 489:
            return KnxManufacturers_M_01E9
        case 49:
            return KnxManufacturers_M_0031
        case 490:
            return KnxManufacturers_M_01EA
        case 491:
            return KnxManufacturers_M_01EB
        case 492:
            return KnxManufacturers_M_01EC
        case 493:
            return KnxManufacturers_M_01ED
        case 495:
            return KnxManufacturers_M_01EF
        case 496:
            return KnxManufacturers_M_01F0
        case 497:
            return KnxManufacturers_M_01F1
        case 498:
            return KnxManufacturers_M_01F2
        case 499:
            return KnxManufacturers_M_01F3
        case 5:
            return KnxManufacturers_M_0005
        case 500:
            return KnxManufacturers_M_01F4
        case 501:
            return KnxManufacturers_M_01F5
        case 502:
            return KnxManufacturers_M_01F6
        case 503:
            return KnxManufacturers_M_01F7
        case 504:
            return KnxManufacturers_M_01F8
        case 505:
            return KnxManufacturers_M_01F9
        case 506:
            return KnxManufacturers_M_01FA
        case 507:
            return KnxManufacturers_M_01FB
        case 508:
            return KnxManufacturers_M_01FC
        case 509:
            return KnxManufacturers_M_01FD
        case 512:
            return KnxManufacturers_M_0200
        case 513:
            return KnxManufacturers_M_0201
        case 514:
            return KnxManufacturers_M_0202
        case 515:
            return KnxManufacturers_M_0203
        case 516:
            return KnxManufacturers_M_0204
        case 517:
            return KnxManufacturers_M_0205
        case 518:
            return KnxManufacturers_M_0206
        case 519:
            return KnxManufacturers_M_0207
        case 52:
            return KnxManufacturers_M_0034
        case 520:
            return KnxManufacturers_M_0208
        case 521:
            return KnxManufacturers_M_0209
        case 522:
            return KnxManufacturers_M_020A
        case 523:
            return KnxManufacturers_M_020B
        case 524:
            return KnxManufacturers_M_020C
        case 525:
            return KnxManufacturers_M_020D
        case 526:
            return KnxManufacturers_M_020E
        case 527:
            return KnxManufacturers_M_020F
        case 528:
            return KnxManufacturers_M_0210
        case 529:
            return KnxManufacturers_M_0211
        case 53:
            return KnxManufacturers_M_0035
        case 530:
            return KnxManufacturers_M_0212
        case 531:
            return KnxManufacturers_M_0213
        case 532:
            return KnxManufacturers_M_0214
        case 533:
            return KnxManufacturers_M_0215
        case 534:
            return KnxManufacturers_M_0216
        case 535:
            return KnxManufacturers_M_0217
        case 536:
            return KnxManufacturers_M_0218
        case 537:
            return KnxManufacturers_M_0219
        case 538:
            return KnxManufacturers_M_021A
        case 539:
            return KnxManufacturers_M_021B
        case 540:
            return KnxManufacturers_M_021C
        case 541:
            return KnxManufacturers_M_021D
        case 542:
            return KnxManufacturers_M_021E
        case 543:
            return KnxManufacturers_M_021F
        case 544:
            return KnxManufacturers_M_0220
        case 545:
            return KnxManufacturers_M_0221
        case 546:
            return KnxManufacturers_M_0222
        case 547:
            return KnxManufacturers_M_0223
        case 548:
            return KnxManufacturers_M_0224
        case 549:
            return KnxManufacturers_M_0225
        case 55:
            return KnxManufacturers_M_0037
        case 550:
            return KnxManufacturers_M_0226
        case 551:
            return KnxManufacturers_M_0227
        case 552:
            return KnxManufacturers_M_0228
        case 553:
            return KnxManufacturers_M_0229
        case 554:
            return KnxManufacturers_M_022A
        case 555:
            return KnxManufacturers_M_022B
        case 556:
            return KnxManufacturers_M_022C
        case 557:
            return KnxManufacturers_M_022D
        case 558:
            return KnxManufacturers_M_022E
        case 559:
            return KnxManufacturers_M_022F
        case 560:
            return KnxManufacturers_M_0230
        case 561:
            return KnxManufacturers_M_0231
        case 562:
            return KnxManufacturers_M_0232
        case 563:
            return KnxManufacturers_M_0233
        case 564:
            return KnxManufacturers_M_0234
        case 565:
            return KnxManufacturers_M_0235
        case 566:
            return KnxManufacturers_M_0236
        case 567:
            return KnxManufacturers_M_0237
        case 568:
            return KnxManufacturers_M_0238
        case 569:
            return KnxManufacturers_M_0239
        case 57:
            return KnxManufacturers_M_0039
        case 570:
            return KnxManufacturers_M_023A
        case 571:
            return KnxManufacturers_M_023B
        case 572:
            return KnxManufacturers_M_023C
        case 573:
            return KnxManufacturers_M_023D
        case 574:
            return KnxManufacturers_M_023E
        case 575:
            return KnxManufacturers_M_023F
        case 576:
            return KnxManufacturers_M_0240
        case 577:
            return KnxManufacturers_M_0241
        case 578:
            return KnxManufacturers_M_0242
        case 579:
            return KnxManufacturers_M_0243
        case 580:
            return KnxManufacturers_M_0244
        case 581:
            return KnxManufacturers_M_0245
        case 582:
            return KnxManufacturers_M_0246
        case 583:
            return KnxManufacturers_M_0247
        case 584:
            return KnxManufacturers_M_0248
        case 585:
            return KnxManufacturers_M_0249
        case 586:
            return KnxManufacturers_M_024A
        case 587:
            return KnxManufacturers_M_024B
        case 588:
            return KnxManufacturers_M_024C
        case 589:
            return KnxManufacturers_M_024D
        case 590:
            return KnxManufacturers_M_024E
        case 591:
            return KnxManufacturers_M_024F
        case 592:
            return KnxManufacturers_M_0250
        case 6:
            return KnxManufacturers_M_0006
        case 61:
            return KnxManufacturers_M_003D
        case 62:
            return KnxManufacturers_M_003E
        case 66:
            return KnxManufacturers_M_0042
        case 67:
            return KnxManufacturers_M_0043
        case 69:
            return KnxManufacturers_M_0045
        case 7:
            return KnxManufacturers_M_0007
        case 71:
            return KnxManufacturers_M_0047
        case 72:
            return KnxManufacturers_M_0048
        case 73:
            return KnxManufacturers_M_0049
        case 75:
            return KnxManufacturers_M_004B
        case 76:
            return KnxManufacturers_M_004C
        case 78:
            return KnxManufacturers_M_004E
        case 8:
            return KnxManufacturers_M_0008
        case 80:
            return KnxManufacturers_M_0050
        case 81:
            return KnxManufacturers_M_0051
        case 82:
            return KnxManufacturers_M_0052
        case 83:
            return KnxManufacturers_M_0053
        case 85:
            return KnxManufacturers_M_0055
        case 89:
            return KnxManufacturers_M_0059
        case 9:
            return KnxManufacturers_M_0009
        case 90:
            return KnxManufacturers_M_005A
        case 92:
            return KnxManufacturers_M_005C
        case 93:
            return KnxManufacturers_M_005D
        case 94:
            return KnxManufacturers_M_005E
        case 95:
            return KnxManufacturers_M_005F
        case 97:
            return KnxManufacturers_M_0061
        case 98:
            return KnxManufacturers_M_0062
        case 99:
            return KnxManufacturers_M_0063
    }
    return 0
}

func CastKnxManufacturers(structType interface{}) KnxManufacturers {
    castFunc := func(typ interface{}) KnxManufacturers {
        if sKnxManufacturers, ok := typ.(KnxManufacturers); ok {
            return sKnxManufacturers
        }
        return 0
    }
    return castFunc(structType)
}

func (m KnxManufacturers) LengthInBits() uint16 {
    return 16
}

func (m KnxManufacturers) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KnxManufacturersParse(io *utils.ReadBuffer) (KnxManufacturers, error) {
    val, err := io.ReadUint16(16)
    if err != nil {
        return 0, nil
    }
    return KnxManufacturersValueOf(val), nil
}

func (e KnxManufacturers) Serialize(io utils.WriteBuffer) error {
    err := io.WriteUint16(16, uint16(e))
    return err
}

func (e KnxManufacturers) String() string {
    switch e {
    case KnxManufacturers_M_0001:
        return "M_0001"
    case KnxManufacturers_M_000A:
        return "M_000A"
    case KnxManufacturers_M_0064:
        return "M_0064"
    case KnxManufacturers_M_0065:
        return "M_0065"
    case KnxManufacturers_M_0066:
        return "M_0066"
    case KnxManufacturers_M_0068:
        return "M_0068"
    case KnxManufacturers_M_0069:
        return "M_0069"
    case KnxManufacturers_M_006A:
        return "M_006A"
    case KnxManufacturers_M_006B:
        return "M_006B"
    case KnxManufacturers_M_006C:
        return "M_006C"
    case KnxManufacturers_M_006D:
        return "M_006D"
    case KnxManufacturers_M_000B:
        return "M_000B"
    case KnxManufacturers_M_006E:
        return "M_006E"
    case KnxManufacturers_M_006F:
        return "M_006F"
    case KnxManufacturers_M_0070:
        return "M_0070"
    case KnxManufacturers_M_0071:
        return "M_0071"
    case KnxManufacturers_M_0072:
        return "M_0072"
    case KnxManufacturers_M_0073:
        return "M_0073"
    case KnxManufacturers_M_0074:
        return "M_0074"
    case KnxManufacturers_M_0075:
        return "M_0075"
    case KnxManufacturers_M_0076:
        return "M_0076"
    case KnxManufacturers_M_0077:
        return "M_0077"
    case KnxManufacturers_M_000C:
        return "M_000C"
    case KnxManufacturers_M_0078:
        return "M_0078"
    case KnxManufacturers_M_0079:
        return "M_0079"
    case KnxManufacturers_M_007A:
        return "M_007A"
    case KnxManufacturers_M_007B:
        return "M_007B"
    case KnxManufacturers_M_007C:
        return "M_007C"
    case KnxManufacturers_M_007D:
        return "M_007D"
    case KnxManufacturers_M_007E:
        return "M_007E"
    case KnxManufacturers_M_007F:
        return "M_007F"
    case KnxManufacturers_M_0080:
        return "M_0080"
    case KnxManufacturers_M_0081:
        return "M_0081"
    case KnxManufacturers_M_0082:
        return "M_0082"
    case KnxManufacturers_M_0083:
        return "M_0083"
    case KnxManufacturers_M_0084:
        return "M_0084"
    case KnxManufacturers_M_0085:
        return "M_0085"
    case KnxManufacturers_M_0086:
        return "M_0086"
    case KnxManufacturers_M_0087:
        return "M_0087"
    case KnxManufacturers_M_0088:
        return "M_0088"
    case KnxManufacturers_M_0089:
        return "M_0089"
    case KnxManufacturers_M_008A:
        return "M_008A"
    case KnxManufacturers_M_008B:
        return "M_008B"
    case KnxManufacturers_M_000E:
        return "M_000E"
    case KnxManufacturers_M_008C:
        return "M_008C"
    case KnxManufacturers_M_008D:
        return "M_008D"
    case KnxManufacturers_M_008E:
        return "M_008E"
    case KnxManufacturers_M_008F:
        return "M_008F"
    case KnxManufacturers_M_0090:
        return "M_0090"
    case KnxManufacturers_M_0091:
        return "M_0091"
    case KnxManufacturers_M_0092:
        return "M_0092"
    case KnxManufacturers_M_0093:
        return "M_0093"
    case KnxManufacturers_M_0094:
        return "M_0094"
    case KnxManufacturers_M_0095:
        return "M_0095"
    case KnxManufacturers_M_0096:
        return "M_0096"
    case KnxManufacturers_M_0097:
        return "M_0097"
    case KnxManufacturers_M_0098:
        return "M_0098"
    case KnxManufacturers_M_0099:
        return "M_0099"
    case KnxManufacturers_M_009A:
        return "M_009A"
    case KnxManufacturers_M_009B:
        return "M_009B"
    case KnxManufacturers_M_009C:
        return "M_009C"
    case KnxManufacturers_M_009D:
        return "M_009D"
    case KnxManufacturers_M_009E:
        return "M_009E"
    case KnxManufacturers_M_009F:
        return "M_009F"
    case KnxManufacturers_M_00A0:
        return "M_00A0"
    case KnxManufacturers_M_00A1:
        return "M_00A1"
    case KnxManufacturers_M_00A2:
        return "M_00A2"
    case KnxManufacturers_M_00A3:
        return "M_00A3"
    case KnxManufacturers_M_00A4:
        return "M_00A4"
    case KnxManufacturers_M_00A5:
        return "M_00A5"
    case KnxManufacturers_M_00A6:
        return "M_00A6"
    case KnxManufacturers_M_00A7:
        return "M_00A7"
    case KnxManufacturers_M_00A8:
        return "M_00A8"
    case KnxManufacturers_M_00A9:
        return "M_00A9"
    case KnxManufacturers_M_00AA:
        return "M_00AA"
    case KnxManufacturers_M_00AB:
        return "M_00AB"
    case KnxManufacturers_M_00AC:
        return "M_00AC"
    case KnxManufacturers_M_00AD:
        return "M_00AD"
    case KnxManufacturers_M_00AE:
        return "M_00AE"
    case KnxManufacturers_M_00AF:
        return "M_00AF"
    case KnxManufacturers_M_00B0:
        return "M_00B0"
    case KnxManufacturers_M_00B1:
        return "M_00B1"
    case KnxManufacturers_M_00B2:
        return "M_00B2"
    case KnxManufacturers_M_00B3:
        return "M_00B3"
    case KnxManufacturers_M_00B4:
        return "M_00B4"
    case KnxManufacturers_M_00B5:
        return "M_00B5"
    case KnxManufacturers_M_00B6:
        return "M_00B6"
    case KnxManufacturers_M_00B7:
        return "M_00B7"
    case KnxManufacturers_M_00B8:
        return "M_00B8"
    case KnxManufacturers_M_00B9:
        return "M_00B9"
    case KnxManufacturers_M_00BA:
        return "M_00BA"
    case KnxManufacturers_M_00BB:
        return "M_00BB"
    case KnxManufacturers_M_00BC:
        return "M_00BC"
    case KnxManufacturers_M_00BD:
        return "M_00BD"
    case KnxManufacturers_M_00BE:
        return "M_00BE"
    case KnxManufacturers_M_00BF:
        return "M_00BF"
    case KnxManufacturers_M_00C0:
        return "M_00C0"
    case KnxManufacturers_M_00C1:
        return "M_00C1"
    case KnxManufacturers_M_00C2:
        return "M_00C2"
    case KnxManufacturers_M_00C3:
        return "M_00C3"
    case KnxManufacturers_M_00C4:
        return "M_00C4"
    case KnxManufacturers_M_00C5:
        return "M_00C5"
    case KnxManufacturers_M_00C6:
        return "M_00C6"
    case KnxManufacturers_M_00C7:
        return "M_00C7"
    case KnxManufacturers_M_0002:
        return "M_0002"
    case KnxManufacturers_M_00C8:
        return "M_00C8"
    case KnxManufacturers_M_00C9:
        return "M_00C9"
    case KnxManufacturers_M_00CA:
        return "M_00CA"
    case KnxManufacturers_M_00CC:
        return "M_00CC"
    case KnxManufacturers_M_00CD:
        return "M_00CD"
    case KnxManufacturers_M_00CE:
        return "M_00CE"
    case KnxManufacturers_M_00CF:
        return "M_00CF"
    case KnxManufacturers_M_00D0:
        return "M_00D0"
    case KnxManufacturers_M_00D1:
        return "M_00D1"
    case KnxManufacturers_M_00D2:
        return "M_00D2"
    case KnxManufacturers_M_00D3:
        return "M_00D3"
    case KnxManufacturers_M_00D6:
        return "M_00D6"
    case KnxManufacturers_M_00D7:
        return "M_00D7"
    case KnxManufacturers_M_00D8:
        return "M_00D8"
    case KnxManufacturers_M_00D9:
        return "M_00D9"
    case KnxManufacturers_M_00DA:
        return "M_00DA"
    case KnxManufacturers_M_00DB:
        return "M_00DB"
    case KnxManufacturers_M_0016:
        return "M_0016"
    case KnxManufacturers_M_00DC:
        return "M_00DC"
    case KnxManufacturers_M_00DE:
        return "M_00DE"
    case KnxManufacturers_M_00DF:
        return "M_00DF"
    case KnxManufacturers_M_00E1:
        return "M_00E1"
    case KnxManufacturers_M_00E3:
        return "M_00E3"
    case KnxManufacturers_M_00E4:
        return "M_00E4"
    case KnxManufacturers_M_00E8:
        return "M_00E8"
    case KnxManufacturers_M_00E9:
        return "M_00E9"
    case KnxManufacturers_M_00EA:
        return "M_00EA"
    case KnxManufacturers_M_00EB:
        return "M_00EB"
    case KnxManufacturers_M_00ED:
        return "M_00ED"
    case KnxManufacturers_M_00EE:
        return "M_00EE"
    case KnxManufacturers_M_00EF:
        return "M_00EF"
    case KnxManufacturers_M_0018:
        return "M_0018"
    case KnxManufacturers_M_00F0:
        return "M_00F0"
    case KnxManufacturers_M_00F1:
        return "M_00F1"
    case KnxManufacturers_M_00F2:
        return "M_00F2"
    case KnxManufacturers_M_00F4:
        return "M_00F4"
    case KnxManufacturers_M_00F5:
        return "M_00F5"
    case KnxManufacturers_M_00F6:
        return "M_00F6"
    case KnxManufacturers_M_00F8:
        return "M_00F8"
    case KnxManufacturers_M_00F9:
        return "M_00F9"
    case KnxManufacturers_M_0019:
        return "M_0019"
    case KnxManufacturers_M_00FA:
        return "M_00FA"
    case KnxManufacturers_M_00FB:
        return "M_00FB"
    case KnxManufacturers_M_00FC:
        return "M_00FC"
    case KnxManufacturers_M_00FD:
        return "M_00FD"
    case KnxManufacturers_M_00FE:
        return "M_00FE"
    case KnxManufacturers_M_0100:
        return "M_0100"
    case KnxManufacturers_M_0101:
        return "M_0101"
    case KnxManufacturers_M_0102:
        return "M_0102"
    case KnxManufacturers_M_0103:
        return "M_0103"
    case KnxManufacturers_M_0104:
        return "M_0104"
    case KnxManufacturers_M_0105:
        return "M_0105"
    case KnxManufacturers_M_0106:
        return "M_0106"
    case KnxManufacturers_M_0107:
        return "M_0107"
    case KnxManufacturers_M_0108:
        return "M_0108"
    case KnxManufacturers_M_0109:
        return "M_0109"
    case KnxManufacturers_M_010A:
        return "M_010A"
    case KnxManufacturers_M_010B:
        return "M_010B"
    case KnxManufacturers_M_010C:
        return "M_010C"
    case KnxManufacturers_M_010D:
        return "M_010D"
    case KnxManufacturers_M_001B:
        return "M_001B"
    case KnxManufacturers_M_010E:
        return "M_010E"
    case KnxManufacturers_M_010F:
        return "M_010F"
    case KnxManufacturers_M_0110:
        return "M_0110"
    case KnxManufacturers_M_0111:
        return "M_0111"
    case KnxManufacturers_M_0112:
        return "M_0112"
    case KnxManufacturers_M_0113:
        return "M_0113"
    case KnxManufacturers_M_0114:
        return "M_0114"
    case KnxManufacturers_M_0115:
        return "M_0115"
    case KnxManufacturers_M_0116:
        return "M_0116"
    case KnxManufacturers_M_0117:
        return "M_0117"
    case KnxManufacturers_M_001C:
        return "M_001C"
    case KnxManufacturers_M_0118:
        return "M_0118"
    case KnxManufacturers_M_0119:
        return "M_0119"
    case KnxManufacturers_M_011A:
        return "M_011A"
    case KnxManufacturers_M_011B:
        return "M_011B"
    case KnxManufacturers_M_011C:
        return "M_011C"
    case KnxManufacturers_M_011D:
        return "M_011D"
    case KnxManufacturers_M_011E:
        return "M_011E"
    case KnxManufacturers_M_011F:
        return "M_011F"
    case KnxManufacturers_M_0120:
        return "M_0120"
    case KnxManufacturers_M_0121:
        return "M_0121"
    case KnxManufacturers_M_001D:
        return "M_001D"
    case KnxManufacturers_M_0122:
        return "M_0122"
    case KnxManufacturers_M_0123:
        return "M_0123"
    case KnxManufacturers_M_0124:
        return "M_0124"
    case KnxManufacturers_M_0125:
        return "M_0125"
    case KnxManufacturers_M_0126:
        return "M_0126"
    case KnxManufacturers_M_0127:
        return "M_0127"
    case KnxManufacturers_M_0128:
        return "M_0128"
    case KnxManufacturers_M_0129:
        return "M_0129"
    case KnxManufacturers_M_012A:
        return "M_012A"
    case KnxManufacturers_M_012B:
        return "M_012B"
    case KnxManufacturers_M_001E:
        return "M_001E"
    case KnxManufacturers_M_012C:
        return "M_012C"
    case KnxManufacturers_M_012D:
        return "M_012D"
    case KnxManufacturers_M_012E:
        return "M_012E"
    case KnxManufacturers_M_012F:
        return "M_012F"
    case KnxManufacturers_M_0130:
        return "M_0130"
    case KnxManufacturers_M_0131:
        return "M_0131"
    case KnxManufacturers_M_0132:
        return "M_0132"
    case KnxManufacturers_M_0133:
        return "M_0133"
    case KnxManufacturers_M_0134:
        return "M_0134"
    case KnxManufacturers_M_0135:
        return "M_0135"
    case KnxManufacturers_M_001F:
        return "M_001F"
    case KnxManufacturers_M_0136:
        return "M_0136"
    case KnxManufacturers_M_0137:
        return "M_0137"
    case KnxManufacturers_M_0138:
        return "M_0138"
    case KnxManufacturers_M_0139:
        return "M_0139"
    case KnxManufacturers_M_013A:
        return "M_013A"
    case KnxManufacturers_M_013B:
        return "M_013B"
    case KnxManufacturers_M_013C:
        return "M_013C"
    case KnxManufacturers_M_013D:
        return "M_013D"
    case KnxManufacturers_M_013E:
        return "M_013E"
    case KnxManufacturers_M_013F:
        return "M_013F"
    case KnxManufacturers_M_0020:
        return "M_0020"
    case KnxManufacturers_M_0140:
        return "M_0140"
    case KnxManufacturers_M_0141:
        return "M_0141"
    case KnxManufacturers_M_0142:
        return "M_0142"
    case KnxManufacturers_M_0143:
        return "M_0143"
    case KnxManufacturers_M_0144:
        return "M_0144"
    case KnxManufacturers_M_0145:
        return "M_0145"
    case KnxManufacturers_M_0146:
        return "M_0146"
    case KnxManufacturers_M_0147:
        return "M_0147"
    case KnxManufacturers_M_0148:
        return "M_0148"
    case KnxManufacturers_M_0149:
        return "M_0149"
    case KnxManufacturers_M_0021:
        return "M_0021"
    case KnxManufacturers_M_014A:
        return "M_014A"
    case KnxManufacturers_M_014B:
        return "M_014B"
    case KnxManufacturers_M_014C:
        return "M_014C"
    case KnxManufacturers_M_014D:
        return "M_014D"
    case KnxManufacturers_M_014E:
        return "M_014E"
    case KnxManufacturers_M_014F:
        return "M_014F"
    case KnxManufacturers_M_0150:
        return "M_0150"
    case KnxManufacturers_M_0151:
        return "M_0151"
    case KnxManufacturers_M_0152:
        return "M_0152"
    case KnxManufacturers_M_0153:
        return "M_0153"
    case KnxManufacturers_M_0022:
        return "M_0022"
    case KnxManufacturers_M_0154:
        return "M_0154"
    case KnxManufacturers_M_0155:
        return "M_0155"
    case KnxManufacturers_M_0156:
        return "M_0156"
    case KnxManufacturers_M_0157:
        return "M_0157"
    case KnxManufacturers_M_0158:
        return "M_0158"
    case KnxManufacturers_M_0159:
        return "M_0159"
    case KnxManufacturers_M_015A:
        return "M_015A"
    case KnxManufacturers_M_015B:
        return "M_015B"
    case KnxManufacturers_M_015C:
        return "M_015C"
    case KnxManufacturers_M_015D:
        return "M_015D"
    case KnxManufacturers_M_015E:
        return "M_015E"
    case KnxManufacturers_M_015F:
        return "M_015F"
    case KnxManufacturers_M_0160:
        return "M_0160"
    case KnxManufacturers_M_0161:
        return "M_0161"
    case KnxManufacturers_M_0162:
        return "M_0162"
    case KnxManufacturers_M_0163:
        return "M_0163"
    case KnxManufacturers_M_0164:
        return "M_0164"
    case KnxManufacturers_M_0165:
        return "M_0165"
    case KnxManufacturers_M_0166:
        return "M_0166"
    case KnxManufacturers_M_0167:
        return "M_0167"
    case KnxManufacturers_M_0024:
        return "M_0024"
    case KnxManufacturers_M_0168:
        return "M_0168"
    case KnxManufacturers_M_0169:
        return "M_0169"
    case KnxManufacturers_M_016A:
        return "M_016A"
    case KnxManufacturers_M_016B:
        return "M_016B"
    case KnxManufacturers_M_016C:
        return "M_016C"
    case KnxManufacturers_M_016D:
        return "M_016D"
    case KnxManufacturers_M_016E:
        return "M_016E"
    case KnxManufacturers_M_016F:
        return "M_016F"
    case KnxManufacturers_M_0170:
        return "M_0170"
    case KnxManufacturers_M_0171:
        return "M_0171"
    case KnxManufacturers_M_0025:
        return "M_0025"
    case KnxManufacturers_M_0172:
        return "M_0172"
    case KnxManufacturers_M_0173:
        return "M_0173"
    case KnxManufacturers_M_0174:
        return "M_0174"
    case KnxManufacturers_M_0175:
        return "M_0175"
    case KnxManufacturers_M_0176:
        return "M_0176"
    case KnxManufacturers_M_0177:
        return "M_0177"
    case KnxManufacturers_M_0178:
        return "M_0178"
    case KnxManufacturers_M_0179:
        return "M_0179"
    case KnxManufacturers_M_017A:
        return "M_017A"
    case KnxManufacturers_M_017B:
        return "M_017B"
    case KnxManufacturers_M_017C:
        return "M_017C"
    case KnxManufacturers_M_017D:
        return "M_017D"
    case KnxManufacturers_M_017E:
        return "M_017E"
    case KnxManufacturers_M_017F:
        return "M_017F"
    case KnxManufacturers_M_0180:
        return "M_0180"
    case KnxManufacturers_M_0181:
        return "M_0181"
    case KnxManufacturers_M_0182:
        return "M_0182"
    case KnxManufacturers_M_0183:
        return "M_0183"
    case KnxManufacturers_M_0184:
        return "M_0184"
    case KnxManufacturers_M_0185:
        return "M_0185"
    case KnxManufacturers_M_0186:
        return "M_0186"
    case KnxManufacturers_M_0187:
        return "M_0187"
    case KnxManufacturers_M_0188:
        return "M_0188"
    case KnxManufacturers_M_0189:
        return "M_0189"
    case KnxManufacturers_M_018A:
        return "M_018A"
    case KnxManufacturers_M_018B:
        return "M_018B"
    case KnxManufacturers_M_018C:
        return "M_018C"
    case KnxManufacturers_M_018D:
        return "M_018D"
    case KnxManufacturers_M_018E:
        return "M_018E"
    case KnxManufacturers_M_018F:
        return "M_018F"
    case KnxManufacturers_M_0004:
        return "M_0004"
    case KnxManufacturers_M_0190:
        return "M_0190"
    case KnxManufacturers_M_0191:
        return "M_0191"
    case KnxManufacturers_M_0192:
        return "M_0192"
    case KnxManufacturers_M_0193:
        return "M_0193"
    case KnxManufacturers_M_0194:
        return "M_0194"
    case KnxManufacturers_M_0195:
        return "M_0195"
    case KnxManufacturers_M_0196:
        return "M_0196"
    case KnxManufacturers_M_0197:
        return "M_0197"
    case KnxManufacturers_M_0198:
        return "M_0198"
    case KnxManufacturers_M_0199:
        return "M_0199"
    case KnxManufacturers_M_0029:
        return "M_0029"
    case KnxManufacturers_M_019A:
        return "M_019A"
    case KnxManufacturers_M_019B:
        return "M_019B"
    case KnxManufacturers_M_019C:
        return "M_019C"
    case KnxManufacturers_M_019D:
        return "M_019D"
    case KnxManufacturers_M_019E:
        return "M_019E"
    case KnxManufacturers_M_019F:
        return "M_019F"
    case KnxManufacturers_M_01A0:
        return "M_01A0"
    case KnxManufacturers_M_01A1:
        return "M_01A1"
    case KnxManufacturers_M_01A2:
        return "M_01A2"
    case KnxManufacturers_M_01A3:
        return "M_01A3"
    case KnxManufacturers_M_002A:
        return "M_002A"
    case KnxManufacturers_M_01A4:
        return "M_01A4"
    case KnxManufacturers_M_01A5:
        return "M_01A5"
    case KnxManufacturers_M_01A6:
        return "M_01A6"
    case KnxManufacturers_M_01A7:
        return "M_01A7"
    case KnxManufacturers_M_01A8:
        return "M_01A8"
    case KnxManufacturers_M_01A9:
        return "M_01A9"
    case KnxManufacturers_M_01AA:
        return "M_01AA"
    case KnxManufacturers_M_01AB:
        return "M_01AB"
    case KnxManufacturers_M_01AC:
        return "M_01AC"
    case KnxManufacturers_M_01AD:
        return "M_01AD"
    case KnxManufacturers_M_01AE:
        return "M_01AE"
    case KnxManufacturers_M_01AF:
        return "M_01AF"
    case KnxManufacturers_M_01B0:
        return "M_01B0"
    case KnxManufacturers_M_01B1:
        return "M_01B1"
    case KnxManufacturers_M_01B2:
        return "M_01B2"
    case KnxManufacturers_M_01B3:
        return "M_01B3"
    case KnxManufacturers_M_01B4:
        return "M_01B4"
    case KnxManufacturers_M_01B5:
        return "M_01B5"
    case KnxManufacturers_M_01B6:
        return "M_01B6"
    case KnxManufacturers_M_01B7:
        return "M_01B7"
    case KnxManufacturers_M_ABB2:
        return "M_ABB2"
    case KnxManufacturers_M_ABB7:
        return "M_ABB7"
    case KnxManufacturers_M_002C:
        return "M_002C"
    case KnxManufacturers_M_01B8:
        return "M_01B8"
    case KnxManufacturers_M_01B9:
        return "M_01B9"
    case KnxManufacturers_M_01BA:
        return "M_01BA"
    case KnxManufacturers_M_01BB:
        return "M_01BB"
    case KnxManufacturers_M_01BC:
        return "M_01BC"
    case KnxManufacturers_M_01BD:
        return "M_01BD"
    case KnxManufacturers_M_01BE:
        return "M_01BE"
    case KnxManufacturers_M_01BF:
        return "M_01BF"
    case KnxManufacturers_M_01C0:
        return "M_01C0"
    case KnxManufacturers_M_01C1:
        return "M_01C1"
    case KnxManufacturers_M_002D:
        return "M_002D"
    case KnxManufacturers_M_01C3:
        return "M_01C3"
    case KnxManufacturers_M_01C4:
        return "M_01C4"
    case KnxManufacturers_M_01C5:
        return "M_01C5"
    case KnxManufacturers_M_01C6:
        return "M_01C6"
    case KnxManufacturers_M_01C7:
        return "M_01C7"
    case KnxManufacturers_M_01C8:
        return "M_01C8"
    case KnxManufacturers_M_01C9:
        return "M_01C9"
    case KnxManufacturers_M_01CA:
        return "M_01CA"
    case KnxManufacturers_M_01CB:
        return "M_01CB"
    case KnxManufacturers_M_002E:
        return "M_002E"
    case KnxManufacturers_M_01CC:
        return "M_01CC"
    case KnxManufacturers_M_01CD:
        return "M_01CD"
    case KnxManufacturers_M_01CE:
        return "M_01CE"
    case KnxManufacturers_M_01CF:
        return "M_01CF"
    case KnxManufacturers_M_01D0:
        return "M_01D0"
    case KnxManufacturers_M_01D1:
        return "M_01D1"
    case KnxManufacturers_M_01D2:
        return "M_01D2"
    case KnxManufacturers_M_01D3:
        return "M_01D3"
    case KnxManufacturers_M_01D4:
        return "M_01D4"
    case KnxManufacturers_M_01D5:
        return "M_01D5"
    case KnxManufacturers_M_01D6:
        return "M_01D6"
    case KnxManufacturers_M_01D7:
        return "M_01D7"
    case KnxManufacturers_M_01D8:
        return "M_01D8"
    case KnxManufacturers_M_01D9:
        return "M_01D9"
    case KnxManufacturers_M_01DA:
        return "M_01DA"
    case KnxManufacturers_M_01DB:
        return "M_01DB"
    case KnxManufacturers_M_01DC:
        return "M_01DC"
    case KnxManufacturers_M_01DD:
        return "M_01DD"
    case KnxManufacturers_M_01DE:
        return "M_01DE"
    case KnxManufacturers_M_01DF:
        return "M_01DF"
    case KnxManufacturers_M_01E0:
        return "M_01E0"
    case KnxManufacturers_M_01E1:
        return "M_01E1"
    case KnxManufacturers_M_01E2:
        return "M_01E2"
    case KnxManufacturers_M_01E3:
        return "M_01E3"
    case KnxManufacturers_M_01E4:
        return "M_01E4"
    case KnxManufacturers_M_01E5:
        return "M_01E5"
    case KnxManufacturers_M_01E6:
        return "M_01E6"
    case KnxManufacturers_M_01E7:
        return "M_01E7"
    case KnxManufacturers_M_01E8:
        return "M_01E8"
    case KnxManufacturers_M_01E9:
        return "M_01E9"
    case KnxManufacturers_M_0031:
        return "M_0031"
    case KnxManufacturers_M_01EA:
        return "M_01EA"
    case KnxManufacturers_M_01EB:
        return "M_01EB"
    case KnxManufacturers_M_01EC:
        return "M_01EC"
    case KnxManufacturers_M_01ED:
        return "M_01ED"
    case KnxManufacturers_M_01EF:
        return "M_01EF"
    case KnxManufacturers_M_01F0:
        return "M_01F0"
    case KnxManufacturers_M_01F1:
        return "M_01F1"
    case KnxManufacturers_M_01F2:
        return "M_01F2"
    case KnxManufacturers_M_01F3:
        return "M_01F3"
    case KnxManufacturers_M_0005:
        return "M_0005"
    case KnxManufacturers_M_01F4:
        return "M_01F4"
    case KnxManufacturers_M_01F5:
        return "M_01F5"
    case KnxManufacturers_M_01F6:
        return "M_01F6"
    case KnxManufacturers_M_01F7:
        return "M_01F7"
    case KnxManufacturers_M_01F8:
        return "M_01F8"
    case KnxManufacturers_M_01F9:
        return "M_01F9"
    case KnxManufacturers_M_01FA:
        return "M_01FA"
    case KnxManufacturers_M_01FB:
        return "M_01FB"
    case KnxManufacturers_M_01FC:
        return "M_01FC"
    case KnxManufacturers_M_01FD:
        return "M_01FD"
    case KnxManufacturers_M_0200:
        return "M_0200"
    case KnxManufacturers_M_0201:
        return "M_0201"
    case KnxManufacturers_M_0202:
        return "M_0202"
    case KnxManufacturers_M_0203:
        return "M_0203"
    case KnxManufacturers_M_0204:
        return "M_0204"
    case KnxManufacturers_M_0205:
        return "M_0205"
    case KnxManufacturers_M_0206:
        return "M_0206"
    case KnxManufacturers_M_0207:
        return "M_0207"
    case KnxManufacturers_M_0034:
        return "M_0034"
    case KnxManufacturers_M_0208:
        return "M_0208"
    case KnxManufacturers_M_0209:
        return "M_0209"
    case KnxManufacturers_M_020A:
        return "M_020A"
    case KnxManufacturers_M_020B:
        return "M_020B"
    case KnxManufacturers_M_020C:
        return "M_020C"
    case KnxManufacturers_M_020D:
        return "M_020D"
    case KnxManufacturers_M_020E:
        return "M_020E"
    case KnxManufacturers_M_020F:
        return "M_020F"
    case KnxManufacturers_M_0210:
        return "M_0210"
    case KnxManufacturers_M_0211:
        return "M_0211"
    case KnxManufacturers_M_0035:
        return "M_0035"
    case KnxManufacturers_M_0212:
        return "M_0212"
    case KnxManufacturers_M_0213:
        return "M_0213"
    case KnxManufacturers_M_0214:
        return "M_0214"
    case KnxManufacturers_M_0215:
        return "M_0215"
    case KnxManufacturers_M_0216:
        return "M_0216"
    case KnxManufacturers_M_0217:
        return "M_0217"
    case KnxManufacturers_M_0218:
        return "M_0218"
    case KnxManufacturers_M_0219:
        return "M_0219"
    case KnxManufacturers_M_021A:
        return "M_021A"
    case KnxManufacturers_M_021B:
        return "M_021B"
    case KnxManufacturers_M_021C:
        return "M_021C"
    case KnxManufacturers_M_021D:
        return "M_021D"
    case KnxManufacturers_M_021E:
        return "M_021E"
    case KnxManufacturers_M_021F:
        return "M_021F"
    case KnxManufacturers_M_0220:
        return "M_0220"
    case KnxManufacturers_M_0221:
        return "M_0221"
    case KnxManufacturers_M_0222:
        return "M_0222"
    case KnxManufacturers_M_0223:
        return "M_0223"
    case KnxManufacturers_M_0224:
        return "M_0224"
    case KnxManufacturers_M_0225:
        return "M_0225"
    case KnxManufacturers_M_0037:
        return "M_0037"
    case KnxManufacturers_M_0226:
        return "M_0226"
    case KnxManufacturers_M_0227:
        return "M_0227"
    case KnxManufacturers_M_0228:
        return "M_0228"
    case KnxManufacturers_M_0229:
        return "M_0229"
    case KnxManufacturers_M_022A:
        return "M_022A"
    case KnxManufacturers_M_022B:
        return "M_022B"
    case KnxManufacturers_M_022C:
        return "M_022C"
    case KnxManufacturers_M_022D:
        return "M_022D"
    case KnxManufacturers_M_022E:
        return "M_022E"
    case KnxManufacturers_M_022F:
        return "M_022F"
    case KnxManufacturers_M_0230:
        return "M_0230"
    case KnxManufacturers_M_0231:
        return "M_0231"
    case KnxManufacturers_M_0232:
        return "M_0232"
    case KnxManufacturers_M_0233:
        return "M_0233"
    case KnxManufacturers_M_0234:
        return "M_0234"
    case KnxManufacturers_M_0235:
        return "M_0235"
    case KnxManufacturers_M_0236:
        return "M_0236"
    case KnxManufacturers_M_0237:
        return "M_0237"
    case KnxManufacturers_M_0238:
        return "M_0238"
    case KnxManufacturers_M_0239:
        return "M_0239"
    case KnxManufacturers_M_0039:
        return "M_0039"
    case KnxManufacturers_M_023A:
        return "M_023A"
    case KnxManufacturers_M_023B:
        return "M_023B"
    case KnxManufacturers_M_023C:
        return "M_023C"
    case KnxManufacturers_M_023D:
        return "M_023D"
    case KnxManufacturers_M_023E:
        return "M_023E"
    case KnxManufacturers_M_023F:
        return "M_023F"
    case KnxManufacturers_M_0240:
        return "M_0240"
    case KnxManufacturers_M_0241:
        return "M_0241"
    case KnxManufacturers_M_0242:
        return "M_0242"
    case KnxManufacturers_M_0243:
        return "M_0243"
    case KnxManufacturers_M_0244:
        return "M_0244"
    case KnxManufacturers_M_0245:
        return "M_0245"
    case KnxManufacturers_M_0246:
        return "M_0246"
    case KnxManufacturers_M_0247:
        return "M_0247"
    case KnxManufacturers_M_0248:
        return "M_0248"
    case KnxManufacturers_M_0249:
        return "M_0249"
    case KnxManufacturers_M_024A:
        return "M_024A"
    case KnxManufacturers_M_024B:
        return "M_024B"
    case KnxManufacturers_M_024C:
        return "M_024C"
    case KnxManufacturers_M_024D:
        return "M_024D"
    case KnxManufacturers_M_024E:
        return "M_024E"
    case KnxManufacturers_M_024F:
        return "M_024F"
    case KnxManufacturers_M_0250:
        return "M_0250"
    case KnxManufacturers_M_0006:
        return "M_0006"
    case KnxManufacturers_M_003D:
        return "M_003D"
    case KnxManufacturers_M_003E:
        return "M_003E"
    case KnxManufacturers_M_0042:
        return "M_0042"
    case KnxManufacturers_M_0043:
        return "M_0043"
    case KnxManufacturers_M_0045:
        return "M_0045"
    case KnxManufacturers_M_0007:
        return "M_0007"
    case KnxManufacturers_M_0047:
        return "M_0047"
    case KnxManufacturers_M_0048:
        return "M_0048"
    case KnxManufacturers_M_0049:
        return "M_0049"
    case KnxManufacturers_M_004B:
        return "M_004B"
    case KnxManufacturers_M_004C:
        return "M_004C"
    case KnxManufacturers_M_004E:
        return "M_004E"
    case KnxManufacturers_M_0008:
        return "M_0008"
    case KnxManufacturers_M_0050:
        return "M_0050"
    case KnxManufacturers_M_0051:
        return "M_0051"
    case KnxManufacturers_M_0052:
        return "M_0052"
    case KnxManufacturers_M_0053:
        return "M_0053"
    case KnxManufacturers_M_0055:
        return "M_0055"
    case KnxManufacturers_M_0059:
        return "M_0059"
    case KnxManufacturers_M_0009:
        return "M_0009"
    case KnxManufacturers_M_005A:
        return "M_005A"
    case KnxManufacturers_M_005C:
        return "M_005C"
    case KnxManufacturers_M_005D:
        return "M_005D"
    case KnxManufacturers_M_005E:
        return "M_005E"
    case KnxManufacturers_M_005F:
        return "M_005F"
    case KnxManufacturers_M_0061:
        return "M_0061"
    case KnxManufacturers_M_0062:
        return "M_0062"
    case KnxManufacturers_M_0063:
        return "M_0063"
    }
    return ""
}
