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
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
)

type DeviceInformation uint16

type IDeviceInformation interface {
	ComObjectTableAddress() uint16
	DeviceDescriptor() uint16
	Serialize(io utils.WriteBuffer) error
}

const (
	DeviceInformation_DEV0001914201 DeviceInformation = 1
	DeviceInformation_DEV0001140C13 DeviceInformation = 2
	DeviceInformation_DEV0001140B11 DeviceInformation = 3
	DeviceInformation_DEV0001803002 DeviceInformation = 4
	DeviceInformation_DEV00641BD610 DeviceInformation = 5
	DeviceInformation_DEV0064760210 DeviceInformation = 6
	DeviceInformation_DEV0064182410 DeviceInformation = 7
	DeviceInformation_DEV0064182310 DeviceInformation = 8
	DeviceInformation_DEV0064705C01 DeviceInformation = 9
	DeviceInformation_DEV0064181910 DeviceInformation = 10
	DeviceInformation_DEV0064181810 DeviceInformation = 11
	DeviceInformation_DEV0064181710 DeviceInformation = 12
	DeviceInformation_DEV0064181610 DeviceInformation = 13
	DeviceInformation_DEV006420C011 DeviceInformation = 14
	DeviceInformation_DEV006420BA11 DeviceInformation = 15
	DeviceInformation_DEV0064182010 DeviceInformation = 16
	DeviceInformation_DEV0064182510 DeviceInformation = 17
	DeviceInformation_DEV0064182610 DeviceInformation = 18
	DeviceInformation_DEV0064182910 DeviceInformation = 19
	DeviceInformation_DEV0064130610 DeviceInformation = 20
	DeviceInformation_DEV0064130710 DeviceInformation = 21
	DeviceInformation_DEV0064133510 DeviceInformation = 22
	DeviceInformation_DEV0064133310 DeviceInformation = 23
	DeviceInformation_DEV0064133410 DeviceInformation = 24
	DeviceInformation_DEV0064133610 DeviceInformation = 25
	DeviceInformation_DEV0064130510 DeviceInformation = 26
	DeviceInformation_DEV0064480611 DeviceInformation = 27
	DeviceInformation_DEV0064482011 DeviceInformation = 28
	DeviceInformation_DEV0064182210 DeviceInformation = 29
	DeviceInformation_DEV0064182710 DeviceInformation = 30
	DeviceInformation_DEV0064183010 DeviceInformation = 31
	DeviceInformation_DEV0064B00812 DeviceInformation = 32
	DeviceInformation_DEV0064B00A01 DeviceInformation = 33
	DeviceInformation_DEV0064760110 DeviceInformation = 34
	DeviceInformation_DEV0064242313 DeviceInformation = 35
	DeviceInformation_DEV0064FF2111 DeviceInformation = 36
	DeviceInformation_DEV0064FF2112 DeviceInformation = 37
	DeviceInformation_DEV0064648B10 DeviceInformation = 38
	DeviceInformation_DEV0064724010 DeviceInformation = 39
	DeviceInformation_DEV006420BD11 DeviceInformation = 40
	DeviceInformation_DEV0064570011 DeviceInformation = 41
	DeviceInformation_DEV0064570310 DeviceInformation = 42
	DeviceInformation_DEV0064570211 DeviceInformation = 43
	DeviceInformation_DEV0064570411 DeviceInformation = 44
	DeviceInformation_DEV0064570110 DeviceInformation = 45
	DeviceInformation_DEV0064615022 DeviceInformation = 46
	DeviceInformation_DEV0064182810 DeviceInformation = 47
	DeviceInformation_DEV0064183110 DeviceInformation = 48
	DeviceInformation_DEV0064133611 DeviceInformation = 49
	DeviceInformation_DEV006A000122 DeviceInformation = 50
	DeviceInformation_DEV006A000222 DeviceInformation = 51
	DeviceInformation_DEV006A070210 DeviceInformation = 52
	DeviceInformation_DEV006BFFF713 DeviceInformation = 53
	DeviceInformation_DEV006BFF2111 DeviceInformation = 54
	DeviceInformation_DEV006BFFF820 DeviceInformation = 55
	DeviceInformation_DEV006B106D10 DeviceInformation = 56
	DeviceInformation_DEV0071123130 DeviceInformation = 57
	DeviceInformation_DEV0071413133 DeviceInformation = 58
	DeviceInformation_DEV0071114019 DeviceInformation = 59
	DeviceInformation_DEV007111306C DeviceInformation = 60
	DeviceInformation_DEV0071231112 DeviceInformation = 61
	DeviceInformation_DEV0071113080 DeviceInformation = 62
	DeviceInformation_DEV0071321212 DeviceInformation = 63
	DeviceInformation_DEV0071321113 DeviceInformation = 64
	DeviceInformation_DEV0071322212 DeviceInformation = 65
	DeviceInformation_DEV0071322112 DeviceInformation = 66
	DeviceInformation_DEV0071322312 DeviceInformation = 67
	DeviceInformation_DEV0071122124 DeviceInformation = 68
	DeviceInformation_DEV007112221E DeviceInformation = 69
	DeviceInformation_DEV0071413314 DeviceInformation = 70
	DeviceInformation_DEV0072300110 DeviceInformation = 71
	DeviceInformation_DEV0076002101 DeviceInformation = 72
	DeviceInformation_DEV0076002001 DeviceInformation = 73
	DeviceInformation_DEV0076002A15 DeviceInformation = 74
	DeviceInformation_DEV0076002815 DeviceInformation = 75
	DeviceInformation_DEV0076002215 DeviceInformation = 76
	DeviceInformation_DEV0076002B15 DeviceInformation = 77
	DeviceInformation_DEV0076002715 DeviceInformation = 78
	DeviceInformation_DEV0076002315 DeviceInformation = 79
	DeviceInformation_DEV0076002415 DeviceInformation = 80
	DeviceInformation_DEV0076002615 DeviceInformation = 81
	DeviceInformation_DEV0076002515 DeviceInformation = 82
	DeviceInformation_DEV0076000201 DeviceInformation = 83
	DeviceInformation_DEV0076000101 DeviceInformation = 84
	DeviceInformation_DEV0076000301 DeviceInformation = 85
	DeviceInformation_DEV0076000401 DeviceInformation = 86
	DeviceInformation_DEV0076002903 DeviceInformation = 87
	DeviceInformation_DEV0076002901 DeviceInformation = 88
	DeviceInformation_DEV007600E503 DeviceInformation = 89
	DeviceInformation_DEV0076004002 DeviceInformation = 90
	DeviceInformation_DEV0076004003 DeviceInformation = 91
	DeviceInformation_DEV0076003402 DeviceInformation = 92
	DeviceInformation_DEV0076003401 DeviceInformation = 93
	DeviceInformation_DEV007600E908 DeviceInformation = 94
	DeviceInformation_DEV007600E907 DeviceInformation = 95
	DeviceInformation_DEV000C181710 DeviceInformation = 96
	DeviceInformation_DEV000C130510 DeviceInformation = 97
	DeviceInformation_DEV000C130610 DeviceInformation = 98
	DeviceInformation_DEV000C133610 DeviceInformation = 99
	DeviceInformation_DEV000C133410 DeviceInformation = 100
	DeviceInformation_DEV000C133310 DeviceInformation = 101
	DeviceInformation_DEV000C133611 DeviceInformation = 102
	DeviceInformation_DEV000C133510 DeviceInformation = 103
	DeviceInformation_DEV000C130710 DeviceInformation = 104
	DeviceInformation_DEV000C760210 DeviceInformation = 105
	DeviceInformation_DEV000C1BD610 DeviceInformation = 106
	DeviceInformation_DEV000C181610 DeviceInformation = 107
	DeviceInformation_DEV000C648B10 DeviceInformation = 108
	DeviceInformation_DEV000C480611 DeviceInformation = 109
	DeviceInformation_DEV000C482011 DeviceInformation = 110
	DeviceInformation_DEV000C724010 DeviceInformation = 111
	DeviceInformation_DEV000C570211 DeviceInformation = 112
	DeviceInformation_DEV000C570310 DeviceInformation = 113
	DeviceInformation_DEV000C570411 DeviceInformation = 114
	DeviceInformation_DEV000C570110 DeviceInformation = 115
	DeviceInformation_DEV000C570011 DeviceInformation = 116
	DeviceInformation_DEV000C20BD11 DeviceInformation = 117
	DeviceInformation_DEV000C20BA11 DeviceInformation = 118
	DeviceInformation_DEV000C760110 DeviceInformation = 119
	DeviceInformation_DEV000C705C01 DeviceInformation = 120
	DeviceInformation_DEV000CFF2112 DeviceInformation = 121
	DeviceInformation_DEV000CB00812 DeviceInformation = 122
	DeviceInformation_DEV000CB00713 DeviceInformation = 123
	DeviceInformation_DEV000C181910 DeviceInformation = 124
	DeviceInformation_DEV000C181810 DeviceInformation = 125
	DeviceInformation_DEV000C20C011 DeviceInformation = 126
	DeviceInformation_DEV0079002527 DeviceInformation = 127
	DeviceInformation_DEV0079004027 DeviceInformation = 128
	DeviceInformation_DEV0079000223 DeviceInformation = 129
	DeviceInformation_DEV0079000123 DeviceInformation = 130
	DeviceInformation_DEV0079001427 DeviceInformation = 131
	DeviceInformation_DEV0079003027 DeviceInformation = 132
	DeviceInformation_DEV0079100C13 DeviceInformation = 133
	DeviceInformation_DEV0079101C11 DeviceInformation = 134
	DeviceInformation_DEV0080709010 DeviceInformation = 135
	DeviceInformation_DEV0080707010 DeviceInformation = 136
	DeviceInformation_DEV0080706010 DeviceInformation = 137
	DeviceInformation_DEV0080706810 DeviceInformation = 138
	DeviceInformation_DEV0080705010 DeviceInformation = 139
	DeviceInformation_DEV0080703013 DeviceInformation = 140
	DeviceInformation_DEV0080704021 DeviceInformation = 141
	DeviceInformation_DEV0080704022 DeviceInformation = 142
	DeviceInformation_DEV0080704020 DeviceInformation = 143
	DeviceInformation_DEV0080701111 DeviceInformation = 144
	DeviceInformation_DEV0080701811 DeviceInformation = 145
	DeviceInformation_DEV008020A110 DeviceInformation = 146
	DeviceInformation_DEV008020A210 DeviceInformation = 147
	DeviceInformation_DEV008020A010 DeviceInformation = 148
	DeviceInformation_DEV0080207212 DeviceInformation = 149
	DeviceInformation_DEV0080209111 DeviceInformation = 150
	DeviceInformation_DEV0080204310 DeviceInformation = 151
	DeviceInformation_DEV008020B612 DeviceInformation = 152
	DeviceInformation_DEV008020B412 DeviceInformation = 153
	DeviceInformation_DEV008020B512 DeviceInformation = 154
	DeviceInformation_DEV0080208310 DeviceInformation = 155
	DeviceInformation_DEV0080702111 DeviceInformation = 156
	DeviceInformation_DEV0081FE0111 DeviceInformation = 157
	DeviceInformation_DEV0081FF3131 DeviceInformation = 158
	DeviceInformation_DEV0081F01313 DeviceInformation = 159
	DeviceInformation_DEV0083002C16 DeviceInformation = 160
	DeviceInformation_DEV0083002E16 DeviceInformation = 161
	DeviceInformation_DEV0083002F16 DeviceInformation = 162
	DeviceInformation_DEV0083012F16 DeviceInformation = 163
	DeviceInformation_DEV0083003210 DeviceInformation = 164
	DeviceInformation_DEV0083001D13 DeviceInformation = 165
	DeviceInformation_DEV0083001E13 DeviceInformation = 166
	DeviceInformation_DEV0083001B13 DeviceInformation = 167
	DeviceInformation_DEV0083001C13 DeviceInformation = 168
	DeviceInformation_DEV0083001F11 DeviceInformation = 169
	DeviceInformation_DEV0083003C10 DeviceInformation = 170
	DeviceInformation_DEV0083001C20 DeviceInformation = 171
	DeviceInformation_DEV0083001B22 DeviceInformation = 172
	DeviceInformation_DEV0083003A14 DeviceInformation = 173
	DeviceInformation_DEV0083003B14 DeviceInformation = 174
	DeviceInformation_DEV0083003B24 DeviceInformation = 175
	DeviceInformation_DEV0083003A24 DeviceInformation = 176
	DeviceInformation_DEV0083005824 DeviceInformation = 177
	DeviceInformation_DEV0083002828 DeviceInformation = 178
	DeviceInformation_DEV0083002928 DeviceInformation = 179
	DeviceInformation_DEV0083002A18 DeviceInformation = 180
	DeviceInformation_DEV0083002B18 DeviceInformation = 181
	DeviceInformation_DEV0083002337 DeviceInformation = 182
	DeviceInformation_DEV0083002838 DeviceInformation = 183
	DeviceInformation_DEV0083002938 DeviceInformation = 184
	DeviceInformation_DEV0083002A38 DeviceInformation = 185
	DeviceInformation_DEV0083002B38 DeviceInformation = 186
	DeviceInformation_DEV0083001321 DeviceInformation = 187
	DeviceInformation_DEV0083001421 DeviceInformation = 188
	DeviceInformation_DEV0083001521 DeviceInformation = 189
	DeviceInformation_DEV0083001621 DeviceInformation = 190
	DeviceInformation_DEV0083000921 DeviceInformation = 191
	DeviceInformation_DEV0083000D11 DeviceInformation = 192
	DeviceInformation_DEV0083000C11 DeviceInformation = 193
	DeviceInformation_DEV0083000E11 DeviceInformation = 194
	DeviceInformation_DEV0083000B11 DeviceInformation = 195
	DeviceInformation_DEV0083000A11 DeviceInformation = 196
	DeviceInformation_DEV0083000A21 DeviceInformation = 197
	DeviceInformation_DEV0083000B21 DeviceInformation = 198
	DeviceInformation_DEV0083000C21 DeviceInformation = 199
	DeviceInformation_DEV0083000D21 DeviceInformation = 200
	DeviceInformation_DEV0083000821 DeviceInformation = 201
	DeviceInformation_DEV0083000E21 DeviceInformation = 202
	DeviceInformation_DEV0083001812 DeviceInformation = 203
	DeviceInformation_DEV0083001712 DeviceInformation = 204
	DeviceInformation_DEV0083001816 DeviceInformation = 205
	DeviceInformation_DEV0083001916 DeviceInformation = 206
	DeviceInformation_DEV0083001716 DeviceInformation = 207
	DeviceInformation_DEV0083001921 DeviceInformation = 208
	DeviceInformation_DEV0083001721 DeviceInformation = 209
	DeviceInformation_DEV0083001821 DeviceInformation = 210
	DeviceInformation_DEV0083001A20 DeviceInformation = 211
	DeviceInformation_DEV0083002320 DeviceInformation = 212
	DeviceInformation_DEV0083001010 DeviceInformation = 213
	DeviceInformation_DEV0083000F10 DeviceInformation = 214
	DeviceInformation_DEV0083003D14 DeviceInformation = 215
	DeviceInformation_DEV0083003E14 DeviceInformation = 216
	DeviceInformation_DEV0083003F14 DeviceInformation = 217
	DeviceInformation_DEV0083004014 DeviceInformation = 218
	DeviceInformation_DEV0083004024 DeviceInformation = 219
	DeviceInformation_DEV0083003D24 DeviceInformation = 220
	DeviceInformation_DEV0083003E24 DeviceInformation = 221
	DeviceInformation_DEV0083003F24 DeviceInformation = 222
	DeviceInformation_DEV0083001112 DeviceInformation = 223
	DeviceInformation_DEV0083001212 DeviceInformation = 224
	DeviceInformation_DEV0083005B12 DeviceInformation = 225
	DeviceInformation_DEV0083005A12 DeviceInformation = 226
	DeviceInformation_DEV0083008410 DeviceInformation = 227
	DeviceInformation_DEV0083008510 DeviceInformation = 228
	DeviceInformation_DEV0083008610 DeviceInformation = 229
	DeviceInformation_DEV0083008710 DeviceInformation = 230
	DeviceInformation_DEV0083002515 DeviceInformation = 231
	DeviceInformation_DEV0083002115 DeviceInformation = 232
	DeviceInformation_DEV0083002015 DeviceInformation = 233
	DeviceInformation_DEV0083002415 DeviceInformation = 234
	DeviceInformation_DEV0083002615 DeviceInformation = 235
	DeviceInformation_DEV0083002215 DeviceInformation = 236
	DeviceInformation_DEV0083002715 DeviceInformation = 237
	DeviceInformation_DEV0083002315 DeviceInformation = 238
	DeviceInformation_DEV0083008B25 DeviceInformation = 239
	DeviceInformation_DEV0083008A25 DeviceInformation = 240
	DeviceInformation_DEV0083008B28 DeviceInformation = 241
	DeviceInformation_DEV0083008A28 DeviceInformation = 242
	DeviceInformation_DEV0083009013 DeviceInformation = 243
	DeviceInformation_DEV0083009213 DeviceInformation = 244
	DeviceInformation_DEV0083009113 DeviceInformation = 245
	DeviceInformation_DEV0083009313 DeviceInformation = 246
	DeviceInformation_DEV0083009413 DeviceInformation = 247
	DeviceInformation_DEV0083009513 DeviceInformation = 248
	DeviceInformation_DEV0083009613 DeviceInformation = 249
	DeviceInformation_DEV0083009713 DeviceInformation = 250
	DeviceInformation_DEV0083009A13 DeviceInformation = 251
	DeviceInformation_DEV0083009B13 DeviceInformation = 252
	DeviceInformation_DEV0083004B11 DeviceInformation = 253
	DeviceInformation_DEV0083004B20 DeviceInformation = 254
	DeviceInformation_DEV0083005514 DeviceInformation = 255
	DeviceInformation_DEV0083006824 DeviceInformation = 256
	DeviceInformation_DEV0083006624 DeviceInformation = 257
	DeviceInformation_DEV0083006524 DeviceInformation = 258
	DeviceInformation_DEV0083006424 DeviceInformation = 259
	DeviceInformation_DEV0083006734 DeviceInformation = 260
	DeviceInformation_DEV0083006434 DeviceInformation = 261
	DeviceInformation_DEV0083006634 DeviceInformation = 262
	DeviceInformation_DEV0083006534 DeviceInformation = 263
	DeviceInformation_DEV0083006A34 DeviceInformation = 264
	DeviceInformation_DEV0083006B34 DeviceInformation = 265
	DeviceInformation_DEV0083006934 DeviceInformation = 266
	DeviceInformation_DEV0083004F11 DeviceInformation = 267
	DeviceInformation_DEV0083004E10 DeviceInformation = 268
	DeviceInformation_DEV0083004D13 DeviceInformation = 269
	DeviceInformation_DEV0083004414 DeviceInformation = 270
	DeviceInformation_DEV0083004114 DeviceInformation = 271
	DeviceInformation_DEV0083004514 DeviceInformation = 272
	DeviceInformation_DEV0083004213 DeviceInformation = 273
	DeviceInformation_DEV0083004313 DeviceInformation = 274
	DeviceInformation_DEV0083004C11 DeviceInformation = 275
	DeviceInformation_DEV0083004913 DeviceInformation = 276
	DeviceInformation_DEV0083004A13 DeviceInformation = 277
	DeviceInformation_DEV0083004712 DeviceInformation = 278
	DeviceInformation_DEV0083004610 DeviceInformation = 279
	DeviceInformation_DEV0083008E12 DeviceInformation = 280
	DeviceInformation_DEV0083004813 DeviceInformation = 281
	DeviceInformation_DEV0083005611 DeviceInformation = 282
	DeviceInformation_DEV0083005710 DeviceInformation = 283
	DeviceInformation_DEV0083005010 DeviceInformation = 284
	DeviceInformation_DEV0083001A10 DeviceInformation = 285
	DeviceInformation_DEV0083002918 DeviceInformation = 286
	DeviceInformation_DEV0083002818 DeviceInformation = 287
	DeviceInformation_DEV0083006724 DeviceInformation = 288
	DeviceInformation_DEV0083006D41 DeviceInformation = 289
	DeviceInformation_DEV0083006E41 DeviceInformation = 290
	DeviceInformation_DEV0083007342 DeviceInformation = 291
	DeviceInformation_DEV0083007242 DeviceInformation = 292
	DeviceInformation_DEV0083006C42 DeviceInformation = 293
	DeviceInformation_DEV0083007542 DeviceInformation = 294
	DeviceInformation_DEV0083007442 DeviceInformation = 295
	DeviceInformation_DEV0083007742 DeviceInformation = 296
	DeviceInformation_DEV0083007642 DeviceInformation = 297
	DeviceInformation_DEV008300B030 DeviceInformation = 298
	DeviceInformation_DEV008300B130 DeviceInformation = 299
	DeviceInformation_DEV008300B230 DeviceInformation = 300
	DeviceInformation_DEV008300B330 DeviceInformation = 301
	DeviceInformation_DEV008300B430 DeviceInformation = 302
	DeviceInformation_DEV008300B530 DeviceInformation = 303
	DeviceInformation_DEV008300B630 DeviceInformation = 304
	DeviceInformation_DEV008300B730 DeviceInformation = 305
	DeviceInformation_DEV0083012843 DeviceInformation = 306
	DeviceInformation_DEV0083012943 DeviceInformation = 307
	DeviceInformation_DEV008300A421 DeviceInformation = 308
	DeviceInformation_DEV008300A521 DeviceInformation = 309
	DeviceInformation_DEV008300A621 DeviceInformation = 310
	DeviceInformation_DEV0083001332 DeviceInformation = 311
	DeviceInformation_DEV0083000932 DeviceInformation = 312
	DeviceInformation_DEV0083001432 DeviceInformation = 313
	DeviceInformation_DEV0083001532 DeviceInformation = 314
	DeviceInformation_DEV0083001632 DeviceInformation = 315
	DeviceInformation_DEV008300A432 DeviceInformation = 316
	DeviceInformation_DEV008300A532 DeviceInformation = 317
	DeviceInformation_DEV008300A632 DeviceInformation = 318
	DeviceInformation_DEV0083000F32 DeviceInformation = 319
	DeviceInformation_DEV0083001032 DeviceInformation = 320
	DeviceInformation_DEV0083000632 DeviceInformation = 321
	DeviceInformation_DEV0083009810 DeviceInformation = 322
	DeviceInformation_DEV0083009910 DeviceInformation = 323
	DeviceInformation_DEV0083005C11 DeviceInformation = 324
	DeviceInformation_DEV0083005D11 DeviceInformation = 325
	DeviceInformation_DEV0083005E11 DeviceInformation = 326
	DeviceInformation_DEV0083005F11 DeviceInformation = 327
	DeviceInformation_DEV0083005413 DeviceInformation = 328
	DeviceInformation_DEV0085000520 DeviceInformation = 329
	DeviceInformation_DEV0085000620 DeviceInformation = 330
	DeviceInformation_DEV0085000720 DeviceInformation = 331
	DeviceInformation_DEV0085012210 DeviceInformation = 332
	DeviceInformation_DEV0085011210 DeviceInformation = 333
	DeviceInformation_DEV0085013220 DeviceInformation = 334
	DeviceInformation_DEV0085010210 DeviceInformation = 335
	DeviceInformation_DEV0085000A10 DeviceInformation = 336
	DeviceInformation_DEV0085000B10 DeviceInformation = 337
	DeviceInformation_DEV0085071010 DeviceInformation = 338
	DeviceInformation_DEV008500FB10 DeviceInformation = 339
	DeviceInformation_DEV0085060210 DeviceInformation = 340
	DeviceInformation_DEV0085060110 DeviceInformation = 341
	DeviceInformation_DEV0085000D20 DeviceInformation = 342
	DeviceInformation_DEV008500C810 DeviceInformation = 343
	DeviceInformation_DEV0085040111 DeviceInformation = 344
	DeviceInformation_DEV008500C910 DeviceInformation = 345
	DeviceInformation_DEV0085045020 DeviceInformation = 346
	DeviceInformation_DEV0085070210 DeviceInformation = 347
	DeviceInformation_DEV0085070110 DeviceInformation = 348
	DeviceInformation_DEV0085070310 DeviceInformation = 349
	DeviceInformation_DEV0085000E20 DeviceInformation = 350
	DeviceInformation_DEV008E596010 DeviceInformation = 351
	DeviceInformation_DEV008E593710 DeviceInformation = 352
	DeviceInformation_DEV008E597710 DeviceInformation = 353
	DeviceInformation_DEV008E598310 DeviceInformation = 354
	DeviceInformation_DEV008E598910 DeviceInformation = 355
	DeviceInformation_DEV008E593720 DeviceInformation = 356
	DeviceInformation_DEV008E598920 DeviceInformation = 357
	DeviceInformation_DEV008E598320 DeviceInformation = 358
	DeviceInformation_DEV008E596021 DeviceInformation = 359
	DeviceInformation_DEV008E597721 DeviceInformation = 360
	DeviceInformation_DEV008E587320 DeviceInformation = 361
	DeviceInformation_DEV008E587020 DeviceInformation = 362
	DeviceInformation_DEV008E587220 DeviceInformation = 363
	DeviceInformation_DEV008E587120 DeviceInformation = 364
	DeviceInformation_DEV008E679910 DeviceInformation = 365
	DeviceInformation_DEV008E618310 DeviceInformation = 366
	DeviceInformation_DEV008E707910 DeviceInformation = 367
	DeviceInformation_DEV008E004010 DeviceInformation = 368
	DeviceInformation_DEV008E570910 DeviceInformation = 369
	DeviceInformation_DEV008E558810 DeviceInformation = 370
	DeviceInformation_DEV008E683410 DeviceInformation = 371
	DeviceInformation_DEV008E707710 DeviceInformation = 372
	DeviceInformation_DEV008E707810 DeviceInformation = 373
	DeviceInformation_DEV0091100013 DeviceInformation = 374
	DeviceInformation_DEV0091100110 DeviceInformation = 375
	DeviceInformation_DEV009E670101 DeviceInformation = 376
	DeviceInformation_DEV009E119311 DeviceInformation = 377
	DeviceInformation_DEV00A2100C13 DeviceInformation = 378
	DeviceInformation_DEV00A2101C11 DeviceInformation = 379
	DeviceInformation_DEV00A2300110 DeviceInformation = 380
	DeviceInformation_DEV0002A05814 DeviceInformation = 381
	DeviceInformation_DEV0002A07114 DeviceInformation = 382
	DeviceInformation_DEV0002134A10 DeviceInformation = 383
	DeviceInformation_DEV0002A03422 DeviceInformation = 384
	DeviceInformation_DEV0002A03321 DeviceInformation = 385
	DeviceInformation_DEV0002648B10 DeviceInformation = 386
	DeviceInformation_DEV0002A09013 DeviceInformation = 387
	DeviceInformation_DEV0002A08F13 DeviceInformation = 388
	DeviceInformation_DEV0002A05510 DeviceInformation = 389
	DeviceInformation_DEV0002A05910 DeviceInformation = 390
	DeviceInformation_DEV0002A05326 DeviceInformation = 391
	DeviceInformation_DEV0002A05428 DeviceInformation = 392
	DeviceInformation_DEV0002A08411 DeviceInformation = 393
	DeviceInformation_DEV0002A08511 DeviceInformation = 394
	DeviceInformation_DEV0002A00F11 DeviceInformation = 395
	DeviceInformation_DEV0002A07310 DeviceInformation = 396
	DeviceInformation_DEV0002A04110 DeviceInformation = 397
	DeviceInformation_DEV0002A03813 DeviceInformation = 398
	DeviceInformation_DEV0002A07F13 DeviceInformation = 399
	DeviceInformation_DEV0002A08832 DeviceInformation = 400
	DeviceInformation_DEV0002A06E32 DeviceInformation = 401
	DeviceInformation_DEV0002A08132 DeviceInformation = 402
	DeviceInformation_DEV0002A01D20 DeviceInformation = 403
	DeviceInformation_DEV0002A02120 DeviceInformation = 404
	DeviceInformation_DEV0002A02520 DeviceInformation = 405
	DeviceInformation_DEV0002A02920 DeviceInformation = 406
	DeviceInformation_DEV0002A03A20 DeviceInformation = 407
	DeviceInformation_DEV0002A05C32 DeviceInformation = 408
	DeviceInformation_DEV0002A06A32 DeviceInformation = 409
	DeviceInformation_DEV0002A09632 DeviceInformation = 410
	DeviceInformation_DEV0002A08932 DeviceInformation = 411
	DeviceInformation_DEV0002A06F32 DeviceInformation = 412
	DeviceInformation_DEV0002A08232 DeviceInformation = 413
	DeviceInformation_DEV0002A01E20 DeviceInformation = 414
	DeviceInformation_DEV0002A02220 DeviceInformation = 415
	DeviceInformation_DEV0002A02620 DeviceInformation = 416
	DeviceInformation_DEV0002A02A20 DeviceInformation = 417
	DeviceInformation_DEV0002A03B20 DeviceInformation = 418
	DeviceInformation_DEV0002A05D32 DeviceInformation = 419
	DeviceInformation_DEV0002A06B32 DeviceInformation = 420
	DeviceInformation_DEV0002A09732 DeviceInformation = 421
	DeviceInformation_DEV0002A08A32 DeviceInformation = 422
	DeviceInformation_DEV0002A07032 DeviceInformation = 423
	DeviceInformation_DEV0002A08332 DeviceInformation = 424
	DeviceInformation_DEV0002A01F20 DeviceInformation = 425
	DeviceInformation_DEV0002A02320 DeviceInformation = 426
	DeviceInformation_DEV0002A02720 DeviceInformation = 427
	DeviceInformation_DEV0002A02B20 DeviceInformation = 428
	DeviceInformation_DEV0002A04820 DeviceInformation = 429
	DeviceInformation_DEV0002A06C32 DeviceInformation = 430
	DeviceInformation_DEV0002A05E32 DeviceInformation = 431
	DeviceInformation_DEV0002A09832 DeviceInformation = 432
	DeviceInformation_DEV0002A06D32 DeviceInformation = 433
	DeviceInformation_DEV0002A08032 DeviceInformation = 434
	DeviceInformation_DEV0002A02020 DeviceInformation = 435
	DeviceInformation_DEV0002A02420 DeviceInformation = 436
	DeviceInformation_DEV0002A02820 DeviceInformation = 437
	DeviceInformation_DEV0002A03920 DeviceInformation = 438
	DeviceInformation_DEV0002A05B32 DeviceInformation = 439
	DeviceInformation_DEV0002A06932 DeviceInformation = 440
	DeviceInformation_DEV0002A09532 DeviceInformation = 441
	DeviceInformation_DEV0002B00813 DeviceInformation = 442
	DeviceInformation_DEV0002A0A610 DeviceInformation = 443
	DeviceInformation_DEV0002A0A611 DeviceInformation = 444
	DeviceInformation_DEV0002A0A510 DeviceInformation = 445
	DeviceInformation_DEV0002A0A511 DeviceInformation = 446
	DeviceInformation_DEV0002A00510 DeviceInformation = 447
	DeviceInformation_DEV0002A00610 DeviceInformation = 448
	DeviceInformation_DEV0002A01511 DeviceInformation = 449
	DeviceInformation_DEV0002A03D11 DeviceInformation = 450
	DeviceInformation_DEV000220C011 DeviceInformation = 451
	DeviceInformation_DEV0002A05213 DeviceInformation = 452
	DeviceInformation_DEV0002A08B10 DeviceInformation = 453
	DeviceInformation_DEV0002A0A210 DeviceInformation = 454
	DeviceInformation_DEV0002A0A010 DeviceInformation = 455
	DeviceInformation_DEV0002A09F10 DeviceInformation = 456
	DeviceInformation_DEV0002A0A110 DeviceInformation = 457
	DeviceInformation_DEV0002A0A013 DeviceInformation = 458
	DeviceInformation_DEV0002A09F13 DeviceInformation = 459
	DeviceInformation_DEV0002A0A213 DeviceInformation = 460
	DeviceInformation_DEV0002A0A113 DeviceInformation = 461
	DeviceInformation_DEV0002A03F11 DeviceInformation = 462
	DeviceInformation_DEV0002A04011 DeviceInformation = 463
	DeviceInformation_DEV0002FF2112 DeviceInformation = 464
	DeviceInformation_DEV0002FF2111 DeviceInformation = 465
	DeviceInformation_DEV0002720111 DeviceInformation = 466
	DeviceInformation_DEV0002613812 DeviceInformation = 467
	DeviceInformation_DEV0002A05713 DeviceInformation = 468
	DeviceInformation_DEV0002A07610 DeviceInformation = 469
	DeviceInformation_DEV0002A01911 DeviceInformation = 470
	DeviceInformation_DEV0002A07611 DeviceInformation = 471
	DeviceInformation_DEV0002A04B10 DeviceInformation = 472
	DeviceInformation_DEV0002A01B14 DeviceInformation = 473
	DeviceInformation_DEV0002A09B11 DeviceInformation = 474
	DeviceInformation_DEV0002A09B12 DeviceInformation = 475
	DeviceInformation_DEV0002A03C10 DeviceInformation = 476
	DeviceInformation_DEV0002A00213 DeviceInformation = 477
	DeviceInformation_DEV0002A00113 DeviceInformation = 478
	DeviceInformation_DEV0002A02C12 DeviceInformation = 479
	DeviceInformation_DEV0002A02D12 DeviceInformation = 480
	DeviceInformation_DEV0002A02E12 DeviceInformation = 481
	DeviceInformation_DEV0002A04C13 DeviceInformation = 482
	DeviceInformation_DEV0002A04D13 DeviceInformation = 483
	DeviceInformation_DEV0002A02F12 DeviceInformation = 484
	DeviceInformation_DEV0002A03012 DeviceInformation = 485
	DeviceInformation_DEV0002A03112 DeviceInformation = 486
	DeviceInformation_DEV0002A04E13 DeviceInformation = 487
	DeviceInformation_DEV0002A04F13 DeviceInformation = 488
	DeviceInformation_DEV0002A01A13 DeviceInformation = 489
	DeviceInformation_DEV0002A09C11 DeviceInformation = 490
	DeviceInformation_DEV0002A09C12 DeviceInformation = 491
	DeviceInformation_DEV0002A01C20 DeviceInformation = 492
	DeviceInformation_DEV0002A09A10 DeviceInformation = 493
	DeviceInformation_DEV0002A09A12 DeviceInformation = 494
	DeviceInformation_DEV000220BA11 DeviceInformation = 495
	DeviceInformation_DEV0002A03D12 DeviceInformation = 496
	DeviceInformation_DEV0002A09110 DeviceInformation = 497
	DeviceInformation_DEV0002A09210 DeviceInformation = 498
	DeviceInformation_DEV0002A09111 DeviceInformation = 499
	DeviceInformation_DEV0002A09211 DeviceInformation = 500
	DeviceInformation_DEV0002A00E21 DeviceInformation = 501
	DeviceInformation_DEV0002A03710 DeviceInformation = 502
	DeviceInformation_DEV0002A01112 DeviceInformation = 503
	DeviceInformation_DEV0002A01216 DeviceInformation = 504
	DeviceInformation_DEV0002A01217 DeviceInformation = 505
	DeviceInformation_DEV000220BD11 DeviceInformation = 506
	DeviceInformation_DEV0002A07F12 DeviceInformation = 507
	DeviceInformation_DEV0002A06613 DeviceInformation = 508
	DeviceInformation_DEV0002A06713 DeviceInformation = 509
	DeviceInformation_DEV0002A06013 DeviceInformation = 510
	DeviceInformation_DEV0002A06113 DeviceInformation = 511
	DeviceInformation_DEV0002A06213 DeviceInformation = 512
	DeviceInformation_DEV0002A06413 DeviceInformation = 513
	DeviceInformation_DEV0002A07713 DeviceInformation = 514
	DeviceInformation_DEV0002A07813 DeviceInformation = 515
	DeviceInformation_DEV0002A07913 DeviceInformation = 516
	DeviceInformation_DEV0002A07914 DeviceInformation = 517
	DeviceInformation_DEV0002A06114 DeviceInformation = 518
	DeviceInformation_DEV0002A06714 DeviceInformation = 519
	DeviceInformation_DEV0002A06414 DeviceInformation = 520
	DeviceInformation_DEV0002A06214 DeviceInformation = 521
	DeviceInformation_DEV0002A06514 DeviceInformation = 522
	DeviceInformation_DEV0002A07714 DeviceInformation = 523
	DeviceInformation_DEV0002A06014 DeviceInformation = 524
	DeviceInformation_DEV0002A06614 DeviceInformation = 525
	DeviceInformation_DEV0002A07814 DeviceInformation = 526
	DeviceInformation_DEV0002A0C310 DeviceInformation = 527
	DeviceInformation_DEV0002632010 DeviceInformation = 528
	DeviceInformation_DEV0002632020 DeviceInformation = 529
	DeviceInformation_DEV0002632040 DeviceInformation = 530
	DeviceInformation_DEV0002632180 DeviceInformation = 531
	DeviceInformation_DEV0002632170 DeviceInformation = 532
	DeviceInformation_DEV0002FF1140 DeviceInformation = 533
	DeviceInformation_DEV0002A07E10 DeviceInformation = 534
	DeviceInformation_DEV0002A07213 DeviceInformation = 535
	DeviceInformation_DEV0002A04A35 DeviceInformation = 536
	DeviceInformation_DEV0002A07420 DeviceInformation = 537
	DeviceInformation_DEV0002A07520 DeviceInformation = 538
	DeviceInformation_DEV0002A07B12 DeviceInformation = 539
	DeviceInformation_DEV0002A07C12 DeviceInformation = 540
	DeviceInformation_DEV0002A04312 DeviceInformation = 541
	DeviceInformation_DEV0002A04412 DeviceInformation = 542
	DeviceInformation_DEV0002A04512 DeviceInformation = 543
	DeviceInformation_DEV0002A04912 DeviceInformation = 544
	DeviceInformation_DEV0002A05012 DeviceInformation = 545
	DeviceInformation_DEV0002A01811 DeviceInformation = 546
	DeviceInformation_DEV0002A03E11 DeviceInformation = 547
	DeviceInformation_DEV0002A08711 DeviceInformation = 548
	DeviceInformation_DEV0002A09311 DeviceInformation = 549
	DeviceInformation_DEV0002A01011 DeviceInformation = 550
	DeviceInformation_DEV0002A01622 DeviceInformation = 551
	DeviceInformation_DEV0002A04210 DeviceInformation = 552
	DeviceInformation_DEV0002A09A13 DeviceInformation = 553
	DeviceInformation_DEV00C8272040 DeviceInformation = 554
	DeviceInformation_DEV00C8272260 DeviceInformation = 555
	DeviceInformation_DEV00C8272060 DeviceInformation = 556
	DeviceInformation_DEV00C8272160 DeviceInformation = 557
	DeviceInformation_DEV00C8272050 DeviceInformation = 558
	DeviceInformation_DEV00C9106D10 DeviceInformation = 559
	DeviceInformation_DEV00C9107C20 DeviceInformation = 560
	DeviceInformation_DEV00C9108511 DeviceInformation = 561
	DeviceInformation_DEV00C9106210 DeviceInformation = 562
	DeviceInformation_DEV00C9109310 DeviceInformation = 563
	DeviceInformation_DEV00C9109210 DeviceInformation = 564
	DeviceInformation_DEV00C9109810 DeviceInformation = 565
	DeviceInformation_DEV00C9109A10 DeviceInformation = 566
	DeviceInformation_DEV00C910A420 DeviceInformation = 567
	DeviceInformation_DEV00C910A110 DeviceInformation = 568
	DeviceInformation_DEV00C910A010 DeviceInformation = 569
	DeviceInformation_DEV00C910A310 DeviceInformation = 570
	DeviceInformation_DEV00C910A210 DeviceInformation = 571
	DeviceInformation_DEV00C9109B10 DeviceInformation = 572
	DeviceInformation_DEV00C9106110 DeviceInformation = 573
	DeviceInformation_DEV00C9109110 DeviceInformation = 574
	DeviceInformation_DEV00C9109610 DeviceInformation = 575
	DeviceInformation_DEV00C9109710 DeviceInformation = 576
	DeviceInformation_DEV00C9109510 DeviceInformation = 577
	DeviceInformation_DEV00C9109910 DeviceInformation = 578
	DeviceInformation_DEV00C9109C10 DeviceInformation = 579
	DeviceInformation_DEV00C910AB10 DeviceInformation = 580
	DeviceInformation_DEV00C910AC10 DeviceInformation = 581
	DeviceInformation_DEV00C910AD10 DeviceInformation = 582
	DeviceInformation_DEV00C910A810 DeviceInformation = 583
	DeviceInformation_DEV00C9106510 DeviceInformation = 584
	DeviceInformation_DEV00C910A710 DeviceInformation = 585
	DeviceInformation_DEV00C9107610 DeviceInformation = 586
	DeviceInformation_DEV00C910890A DeviceInformation = 587
	DeviceInformation_DEV00C9FF1012 DeviceInformation = 588
	DeviceInformation_DEV00C9FF0913 DeviceInformation = 589
	DeviceInformation_DEV00C9FF1112 DeviceInformation = 590
	DeviceInformation_DEV00C9100310 DeviceInformation = 591
	DeviceInformation_DEV00C9101110 DeviceInformation = 592
	DeviceInformation_DEV00C9101010 DeviceInformation = 593
	DeviceInformation_DEV00C9103710 DeviceInformation = 594
	DeviceInformation_DEV00C9101310 DeviceInformation = 595
	DeviceInformation_DEV00C9FF0D12 DeviceInformation = 596
	DeviceInformation_DEV00C9100E10 DeviceInformation = 597
	DeviceInformation_DEV00C9100610 DeviceInformation = 598
	DeviceInformation_DEV00C9100510 DeviceInformation = 599
	DeviceInformation_DEV00C9100710 DeviceInformation = 600
	DeviceInformation_DEV00C9FF1D20 DeviceInformation = 601
	DeviceInformation_DEV00C9FF1C10 DeviceInformation = 602
	DeviceInformation_DEV00C9100810 DeviceInformation = 603
	DeviceInformation_DEV00C9FF1420 DeviceInformation = 604
	DeviceInformation_DEV00C9100D10 DeviceInformation = 605
	DeviceInformation_DEV00C9101220 DeviceInformation = 606
	DeviceInformation_DEV00C9102330 DeviceInformation = 607
	DeviceInformation_DEV00C9102130 DeviceInformation = 608
	DeviceInformation_DEV00C9102430 DeviceInformation = 609
	DeviceInformation_DEV00C9100831 DeviceInformation = 610
	DeviceInformation_DEV00C9102530 DeviceInformation = 611
	DeviceInformation_DEV00C9100531 DeviceInformation = 612
	DeviceInformation_DEV00C9102030 DeviceInformation = 613
	DeviceInformation_DEV00C9100731 DeviceInformation = 614
	DeviceInformation_DEV00C9100631 DeviceInformation = 615
	DeviceInformation_DEV00C9102230 DeviceInformation = 616
	DeviceInformation_DEV00C9100632 DeviceInformation = 617
	DeviceInformation_DEV00C9100532 DeviceInformation = 618
	DeviceInformation_DEV00C9100732 DeviceInformation = 619
	DeviceInformation_DEV00C9100832 DeviceInformation = 620
	DeviceInformation_DEV00C9102532 DeviceInformation = 621
	DeviceInformation_DEV00C9102132 DeviceInformation = 622
	DeviceInformation_DEV00C9102332 DeviceInformation = 623
	DeviceInformation_DEV00C9102432 DeviceInformation = 624
	DeviceInformation_DEV00C9102032 DeviceInformation = 625
	DeviceInformation_DEV00C9102232 DeviceInformation = 626
	DeviceInformation_DEV00C9104432 DeviceInformation = 627
	DeviceInformation_DEV00C9100D11 DeviceInformation = 628
	DeviceInformation_DEV00C9100633 DeviceInformation = 629
	DeviceInformation_DEV00C9100533 DeviceInformation = 630
	DeviceInformation_DEV00C9100733 DeviceInformation = 631
	DeviceInformation_DEV00C9100833 DeviceInformation = 632
	DeviceInformation_DEV00C9102533 DeviceInformation = 633
	DeviceInformation_DEV00C9102133 DeviceInformation = 634
	DeviceInformation_DEV00C9102333 DeviceInformation = 635
	DeviceInformation_DEV00C9102433 DeviceInformation = 636
	DeviceInformation_DEV00C9102033 DeviceInformation = 637
	DeviceInformation_DEV00C9102233 DeviceInformation = 638
	DeviceInformation_DEV00C9104810 DeviceInformation = 639
	DeviceInformation_DEV00C9FF1A11 DeviceInformation = 640
	DeviceInformation_DEV00C9100212 DeviceInformation = 641
	DeviceInformation_DEV00C9FF0A11 DeviceInformation = 642
	DeviceInformation_DEV00C9FF0C12 DeviceInformation = 643
	DeviceInformation_DEV00C9100112 DeviceInformation = 644
	DeviceInformation_DEV00C9FF1911 DeviceInformation = 645
	DeviceInformation_DEV00C9FF0B12 DeviceInformation = 646
	DeviceInformation_DEV00C9FF0715 DeviceInformation = 647
	DeviceInformation_DEV00C9FF1B10 DeviceInformation = 648
	DeviceInformation_DEV00C9101610 DeviceInformation = 649
	DeviceInformation_DEV00C9FF1B11 DeviceInformation = 650
	DeviceInformation_DEV00C9101611 DeviceInformation = 651
	DeviceInformation_DEV00C9101612 DeviceInformation = 652
	DeviceInformation_DEV00C9FF0F11 DeviceInformation = 653
	DeviceInformation_DEV00C9FF1E30 DeviceInformation = 654
	DeviceInformation_DEV00C9100410 DeviceInformation = 655
	DeviceInformation_DEV00C9106410 DeviceInformation = 656
	DeviceInformation_DEV00C9106710 DeviceInformation = 657
	DeviceInformation_DEV00C9106810 DeviceInformation = 658
	DeviceInformation_DEV00C9106010 DeviceInformation = 659
	DeviceInformation_DEV00C9106310 DeviceInformation = 660
	DeviceInformation_DEV00C9107110 DeviceInformation = 661
	DeviceInformation_DEV00C9107210 DeviceInformation = 662
	DeviceInformation_DEV00C9107310 DeviceInformation = 663
	DeviceInformation_DEV00C9107010 DeviceInformation = 664
	DeviceInformation_DEV00C9107A20 DeviceInformation = 665
	DeviceInformation_DEV00C9107B20 DeviceInformation = 666
	DeviceInformation_DEV00C9107820 DeviceInformation = 667
	DeviceInformation_DEV00C9107920 DeviceInformation = 668
	DeviceInformation_DEV00C9104433 DeviceInformation = 669
	DeviceInformation_DEV00C9107C11 DeviceInformation = 670
	DeviceInformation_DEV00C9107711 DeviceInformation = 671
	DeviceInformation_DEV00C9108310 DeviceInformation = 672
	DeviceInformation_DEV00C9108210 DeviceInformation = 673
	DeviceInformation_DEV00C9108610 DeviceInformation = 674
	DeviceInformation_DEV00C9107D10 DeviceInformation = 675
	DeviceInformation_DEV00CE648B10 DeviceInformation = 676
	DeviceInformation_DEV00CE494513 DeviceInformation = 677
	DeviceInformation_DEV00CE494311 DeviceInformation = 678
	DeviceInformation_DEV00CE494810 DeviceInformation = 679
	DeviceInformation_DEV00CE494712 DeviceInformation = 680
	DeviceInformation_DEV00CE494012 DeviceInformation = 681
	DeviceInformation_DEV00CE494111 DeviceInformation = 682
	DeviceInformation_DEV00CE494210 DeviceInformation = 683
	DeviceInformation_DEV00CE494610 DeviceInformation = 684
	DeviceInformation_DEV00CE494412 DeviceInformation = 685
	DeviceInformation_DEV00D0660212 DeviceInformation = 686
	DeviceInformation_DEV00E8000A10 DeviceInformation = 687
	DeviceInformation_DEV00E8000B10 DeviceInformation = 688
	DeviceInformation_DEV00E8000910 DeviceInformation = 689
	DeviceInformation_DEV00E8001112 DeviceInformation = 690
	DeviceInformation_DEV00E8000C14 DeviceInformation = 691
	DeviceInformation_DEV00E8000D13 DeviceInformation = 692
	DeviceInformation_DEV00E8000E12 DeviceInformation = 693
	DeviceInformation_DEV00E8001310 DeviceInformation = 694
	DeviceInformation_DEV00E8001410 DeviceInformation = 695
	DeviceInformation_DEV00E8001510 DeviceInformation = 696
	DeviceInformation_DEV00E8000F10 DeviceInformation = 697
	DeviceInformation_DEV00E8001010 DeviceInformation = 698
	DeviceInformation_DEV00E8000612 DeviceInformation = 699
	DeviceInformation_DEV00E8000812 DeviceInformation = 700
	DeviceInformation_DEV00E8000712 DeviceInformation = 701
	DeviceInformation_DEV00F4501311 DeviceInformation = 702
	DeviceInformation_DEV00F4B00911 DeviceInformation = 703
	DeviceInformation_DEV0019512710 DeviceInformation = 704
	DeviceInformation_DEV0019512810 DeviceInformation = 705
	DeviceInformation_DEV0019512910 DeviceInformation = 706
	DeviceInformation_DEV0019E30D10 DeviceInformation = 707
	DeviceInformation_DEV0019512211 DeviceInformation = 708
	DeviceInformation_DEV0019512311 DeviceInformation = 709
	DeviceInformation_DEV0019512111 DeviceInformation = 710
	DeviceInformation_DEV0019520D11 DeviceInformation = 711
	DeviceInformation_DEV0019E30B12 DeviceInformation = 712
	DeviceInformation_DEV0019530812 DeviceInformation = 713
	DeviceInformation_DEV0019530912 DeviceInformation = 714
	DeviceInformation_DEV0019530612 DeviceInformation = 715
	DeviceInformation_DEV0019530711 DeviceInformation = 716
	DeviceInformation_DEV0019E30A11 DeviceInformation = 717
	DeviceInformation_DEV0019E20111 DeviceInformation = 718
	DeviceInformation_DEV0019E20210 DeviceInformation = 719
	DeviceInformation_DEV0019E30C11 DeviceInformation = 720
	DeviceInformation_DEV0019E11310 DeviceInformation = 721
	DeviceInformation_DEV0019E11210 DeviceInformation = 722
	DeviceInformation_DEV0019E30610 DeviceInformation = 723
	DeviceInformation_DEV0019E30710 DeviceInformation = 724
	DeviceInformation_DEV0019E30910 DeviceInformation = 725
	DeviceInformation_DEV0019E30810 DeviceInformation = 726
	DeviceInformation_DEV0019E25510 DeviceInformation = 727
	DeviceInformation_DEV0019E20410 DeviceInformation = 728
	DeviceInformation_DEV0019E20310 DeviceInformation = 729
	DeviceInformation_DEV0019E25610 DeviceInformation = 730
	DeviceInformation_DEV0019512010 DeviceInformation = 731
	DeviceInformation_DEV0019520C10 DeviceInformation = 732
	DeviceInformation_DEV0019520710 DeviceInformation = 733
	DeviceInformation_DEV0019520210 DeviceInformation = 734
	DeviceInformation_DEV0019E25010 DeviceInformation = 735
	DeviceInformation_DEV0019E25110 DeviceInformation = 736
	DeviceInformation_DEV0019130710 DeviceInformation = 737
	DeviceInformation_DEV0019272050 DeviceInformation = 738
	DeviceInformation_DEV0019520910 DeviceInformation = 739
	DeviceInformation_DEV0019520A10 DeviceInformation = 740
	DeviceInformation_DEV0019520B10 DeviceInformation = 741
	DeviceInformation_DEV0019520412 DeviceInformation = 742
	DeviceInformation_DEV0019520812 DeviceInformation = 743
	DeviceInformation_DEV0019512510 DeviceInformation = 744
	DeviceInformation_DEV0019512410 DeviceInformation = 745
	DeviceInformation_DEV0019512610 DeviceInformation = 746
	DeviceInformation_DEV0019511711 DeviceInformation = 747
	DeviceInformation_DEV0019511811 DeviceInformation = 748
	DeviceInformation_DEV0019522212 DeviceInformation = 749
	DeviceInformation_DEV0019FF0716 DeviceInformation = 750
	DeviceInformation_DEV0019FF1420 DeviceInformation = 751
	DeviceInformation_DEV0019522112 DeviceInformation = 752
	DeviceInformation_DEV0019522011 DeviceInformation = 753
	DeviceInformation_DEV0019522311 DeviceInformation = 754
	DeviceInformation_DEV0019E12410 DeviceInformation = 755
	DeviceInformation_DEV0019000311 DeviceInformation = 756
	DeviceInformation_DEV0019000411 DeviceInformation = 757
	DeviceInformation_DEV0019070210 DeviceInformation = 758
	DeviceInformation_DEV0019070E11 DeviceInformation = 759
	DeviceInformation_DEV0019724010 DeviceInformation = 760
	DeviceInformation_DEV0019520610 DeviceInformation = 761
	DeviceInformation_DEV0019520510 DeviceInformation = 762
	DeviceInformation_DEV00FB101111 DeviceInformation = 763
	DeviceInformation_DEV00FB103001 DeviceInformation = 764
	DeviceInformation_DEV00FB104401 DeviceInformation = 765
	DeviceInformation_DEV00FB124002 DeviceInformation = 766
	DeviceInformation_DEV00FB104102 DeviceInformation = 767
	DeviceInformation_DEV00FB104201 DeviceInformation = 768
	DeviceInformation_DEV00FBF77603 DeviceInformation = 769
	DeviceInformation_DEV00FB104301 DeviceInformation = 770
	DeviceInformation_DEV00FB104601 DeviceInformation = 771
	DeviceInformation_DEV00FB104701 DeviceInformation = 772
	DeviceInformation_DEV00FB105101 DeviceInformation = 773
	DeviceInformation_DEV0103030110 DeviceInformation = 774
	DeviceInformation_DEV0103010113 DeviceInformation = 775
	DeviceInformation_DEV0103090110 DeviceInformation = 776
	DeviceInformation_DEV0103020111 DeviceInformation = 777
	DeviceInformation_DEV0103020112 DeviceInformation = 778
	DeviceInformation_DEV0103040110 DeviceInformation = 779
	DeviceInformation_DEV0103050111 DeviceInformation = 780
	DeviceInformation_DEV0107000301 DeviceInformation = 781
	DeviceInformation_DEV0107000101 DeviceInformation = 782
	DeviceInformation_DEV0107000201 DeviceInformation = 783
	DeviceInformation_DEV0107020801 DeviceInformation = 784
	DeviceInformation_DEV0107020401 DeviceInformation = 785
	DeviceInformation_DEV0107020001 DeviceInformation = 786
	DeviceInformation_DEV010701F801 DeviceInformation = 787
	DeviceInformation_DEV010701FC01 DeviceInformation = 788
	DeviceInformation_DEV0107020C01 DeviceInformation = 789
	DeviceInformation_DEV010F100801 DeviceInformation = 790
	DeviceInformation_DEV010F100601 DeviceInformation = 791
	DeviceInformation_DEV010F100401 DeviceInformation = 792
	DeviceInformation_DEV010F030601 DeviceInformation = 793
	DeviceInformation_DEV010F010301 DeviceInformation = 794
	DeviceInformation_DEV010F010101 DeviceInformation = 795
	DeviceInformation_DEV010F010201 DeviceInformation = 796
	DeviceInformation_DEV010F000302 DeviceInformation = 797
	DeviceInformation_DEV010F000402 DeviceInformation = 798
	DeviceInformation_DEV010F000102 DeviceInformation = 799
	DeviceInformation_DEV011EBB8211 DeviceInformation = 800
	DeviceInformation_DEV011E108111 DeviceInformation = 801
	DeviceInformation_DEV0123010010 DeviceInformation = 802
	DeviceInformation_DEV001E478010 DeviceInformation = 803
	DeviceInformation_DEV001E706611 DeviceInformation = 804
	DeviceInformation_DEV001E706811 DeviceInformation = 805
	DeviceInformation_DEV001E473012 DeviceInformation = 806
	DeviceInformation_DEV001E20A011 DeviceInformation = 807
	DeviceInformation_DEV001E209011 DeviceInformation = 808
	DeviceInformation_DEV001E209811 DeviceInformation = 809
	DeviceInformation_DEV001E208811 DeviceInformation = 810
	DeviceInformation_DEV001E208011 DeviceInformation = 811
	DeviceInformation_DEV001E207821 DeviceInformation = 812
	DeviceInformation_DEV001E20CA12 DeviceInformation = 813
	DeviceInformation_DEV001E20B312 DeviceInformation = 814
	DeviceInformation_DEV001E20B012 DeviceInformation = 815
	DeviceInformation_DEV001E302612 DeviceInformation = 816
	DeviceInformation_DEV001E302312 DeviceInformation = 817
	DeviceInformation_DEV001E302012 DeviceInformation = 818
	DeviceInformation_DEV001E20A811 DeviceInformation = 819
	DeviceInformation_DEV001E20C412 DeviceInformation = 820
	DeviceInformation_DEV001E20C712 DeviceInformation = 821
	DeviceInformation_DEV001E20AD12 DeviceInformation = 822
	DeviceInformation_DEV001E443720 DeviceInformation = 823
	DeviceInformation_DEV001E441821 DeviceInformation = 824
	DeviceInformation_DEV001E443810 DeviceInformation = 825
	DeviceInformation_DEV001E140C12 DeviceInformation = 826
	DeviceInformation_DEV001E471611 DeviceInformation = 827
	DeviceInformation_DEV001E479024 DeviceInformation = 828
	DeviceInformation_DEV001E471A11 DeviceInformation = 829
	DeviceInformation_DEV001E477A10 DeviceInformation = 830
	DeviceInformation_DEV001E470A11 DeviceInformation = 831
	DeviceInformation_DEV001E480B11 DeviceInformation = 832
	DeviceInformation_DEV001E487B10 DeviceInformation = 833
	DeviceInformation_DEV001E440411 DeviceInformation = 834
	DeviceInformation_DEV001E447211 DeviceInformation = 835
	DeviceInformation_DEV0142010010 DeviceInformation = 836
	DeviceInformation_DEV0142010011 DeviceInformation = 837
	DeviceInformation_DEV017A130401 DeviceInformation = 838
	DeviceInformation_DEV017A130201 DeviceInformation = 839
	DeviceInformation_DEV017A130801 DeviceInformation = 840
	DeviceInformation_DEV017A130601 DeviceInformation = 841
	DeviceInformation_DEV017A300102 DeviceInformation = 842
	DeviceInformation_DEV0193323C11 DeviceInformation = 843
	DeviceInformation_DEV0198101110 DeviceInformation = 844
	DeviceInformation_DEV01C4030110 DeviceInformation = 845
	DeviceInformation_DEV01C4030210 DeviceInformation = 846
	DeviceInformation_DEV01C4021210 DeviceInformation = 847
	DeviceInformation_DEV01C4001010 DeviceInformation = 848
	DeviceInformation_DEV01C4020610 DeviceInformation = 849
	DeviceInformation_DEV01C4020910 DeviceInformation = 850
	DeviceInformation_DEV01C4020810 DeviceInformation = 851
	DeviceInformation_DEV01C4010710 DeviceInformation = 852
	DeviceInformation_DEV01C4050210 DeviceInformation = 853
	DeviceInformation_DEV01C4010810 DeviceInformation = 854
	DeviceInformation_DEV01C4020510 DeviceInformation = 855
	DeviceInformation_DEV01C4040110 DeviceInformation = 856
	DeviceInformation_DEV01C4040310 DeviceInformation = 857
	DeviceInformation_DEV01C4040210 DeviceInformation = 858
	DeviceInformation_DEV01C4101210 DeviceInformation = 859
	DeviceInformation_DEV003D020109 DeviceInformation = 860
	DeviceInformation_DEV01DB000301 DeviceInformation = 861
	DeviceInformation_DEV01DB000201 DeviceInformation = 862
	DeviceInformation_DEV01DB000401 DeviceInformation = 863
	DeviceInformation_DEV01DB000801 DeviceInformation = 864
	DeviceInformation_DEV01DB001201 DeviceInformation = 865
	DeviceInformation_DEV009A000400 DeviceInformation = 866
	DeviceInformation_DEV009A100400 DeviceInformation = 867
	DeviceInformation_DEV009A200C00 DeviceInformation = 868
	DeviceInformation_DEV009A200E00 DeviceInformation = 869
	DeviceInformation_DEV009A000201 DeviceInformation = 870
	DeviceInformation_DEV009A000300 DeviceInformation = 871
	DeviceInformation_DEV009A00C000 DeviceInformation = 872
	DeviceInformation_DEV009A00B000 DeviceInformation = 873
	DeviceInformation_DEV009A00C002 DeviceInformation = 874
	DeviceInformation_DEV009A200100 DeviceInformation = 875
	DeviceInformation_DEV004E400010 DeviceInformation = 876
	DeviceInformation_DEV004E030031 DeviceInformation = 877
	DeviceInformation_DEV012B010110 DeviceInformation = 878
	DeviceInformation_DEV01F6E0E110 DeviceInformation = 879
	DeviceInformation_DEV0088100010 DeviceInformation = 880
	DeviceInformation_DEV0088100210 DeviceInformation = 881
	DeviceInformation_DEV0088100110 DeviceInformation = 882
	DeviceInformation_DEV0088110010 DeviceInformation = 883
	DeviceInformation_DEV0088120412 DeviceInformation = 884
	DeviceInformation_DEV0088120113 DeviceInformation = 885
	DeviceInformation_DEV011A4B5201 DeviceInformation = 886
	DeviceInformation_DEV008B020301 DeviceInformation = 887
	DeviceInformation_DEV008B010610 DeviceInformation = 888
	DeviceInformation_DEV008B030110 DeviceInformation = 889
	DeviceInformation_DEV008B030310 DeviceInformation = 890
	DeviceInformation_DEV008B030210 DeviceInformation = 891
	DeviceInformation_DEV008B031512 DeviceInformation = 892
	DeviceInformation_DEV008B031412 DeviceInformation = 893
	DeviceInformation_DEV008B031312 DeviceInformation = 894
	DeviceInformation_DEV008B031212 DeviceInformation = 895
	DeviceInformation_DEV008B031112 DeviceInformation = 896
	DeviceInformation_DEV008B031012 DeviceInformation = 897
	DeviceInformation_DEV008B030510 DeviceInformation = 898
	DeviceInformation_DEV008B030410 DeviceInformation = 899
	DeviceInformation_DEV008B020310 DeviceInformation = 900
	DeviceInformation_DEV008B020210 DeviceInformation = 901
	DeviceInformation_DEV008B020110 DeviceInformation = 902
	DeviceInformation_DEV008B010110 DeviceInformation = 903
	DeviceInformation_DEV008B010210 DeviceInformation = 904
	DeviceInformation_DEV008B010310 DeviceInformation = 905
	DeviceInformation_DEV008B010410 DeviceInformation = 906
	DeviceInformation_DEV008B040110 DeviceInformation = 907
	DeviceInformation_DEV008B040210 DeviceInformation = 908
	DeviceInformation_DEV008B010910 DeviceInformation = 909
	DeviceInformation_DEV008B010710 DeviceInformation = 910
	DeviceInformation_DEV008B010810 DeviceInformation = 911
	DeviceInformation_DEV008B041111 DeviceInformation = 912
	DeviceInformation_DEV008B041211 DeviceInformation = 913
	DeviceInformation_DEV008B041311 DeviceInformation = 914
	DeviceInformation_DEV00A600020A DeviceInformation = 915
	DeviceInformation_DEV00A6000B10 DeviceInformation = 916
	DeviceInformation_DEV00A6000300 DeviceInformation = 917
	DeviceInformation_DEV00A6000705 DeviceInformation = 918
	DeviceInformation_DEV00A6000605 DeviceInformation = 919
	DeviceInformation_DEV00A6000500 DeviceInformation = 920
	DeviceInformation_DEV00A6000C10 DeviceInformation = 921
	DeviceInformation_DEV002CA01811 DeviceInformation = 922
	DeviceInformation_DEV002CA07033 DeviceInformation = 923
	DeviceInformation_DEV002C555020 DeviceInformation = 924
	DeviceInformation_DEV002C556421 DeviceInformation = 925
	DeviceInformation_DEV002C05F211 DeviceInformation = 926
	DeviceInformation_DEV002C05F411 DeviceInformation = 927
	DeviceInformation_DEV002C05E613 DeviceInformation = 928
	DeviceInformation_DEV002CA07914 DeviceInformation = 929
	DeviceInformation_DEV002C060A13 DeviceInformation = 930
	DeviceInformation_DEV002C3A0212 DeviceInformation = 931
	DeviceInformation_DEV002C060210 DeviceInformation = 932
	DeviceInformation_DEV002CA00213 DeviceInformation = 933
	DeviceInformation_DEV002CA0A611 DeviceInformation = 934
	DeviceInformation_DEV002CA07B11 DeviceInformation = 935
	DeviceInformation_DEV002C063210 DeviceInformation = 936
	DeviceInformation_DEV002C063110 DeviceInformation = 937
	DeviceInformation_DEV002C062E10 DeviceInformation = 938
	DeviceInformation_DEV002C062C10 DeviceInformation = 939
	DeviceInformation_DEV002C062D10 DeviceInformation = 940
	DeviceInformation_DEV002C063310 DeviceInformation = 941
	DeviceInformation_DEV002C05EB10 DeviceInformation = 942
	DeviceInformation_DEV002C05F110 DeviceInformation = 943
	DeviceInformation_DEV002C0B8830 DeviceInformation = 944
	DeviceInformation_DEV00A0B07101 DeviceInformation = 945
	DeviceInformation_DEV00A0B07001 DeviceInformation = 946
	DeviceInformation_DEV00A0B07203 DeviceInformation = 947
	DeviceInformation_DEV00A0B02101 DeviceInformation = 948
	DeviceInformation_DEV00A0B02401 DeviceInformation = 949
	DeviceInformation_DEV00A0B02301 DeviceInformation = 950
	DeviceInformation_DEV00A0B02601 DeviceInformation = 951
	DeviceInformation_DEV00A0B02201 DeviceInformation = 952
	DeviceInformation_DEV00A0B01902 DeviceInformation = 953
	DeviceInformation_DEV0004147112 DeviceInformation = 954
	DeviceInformation_DEV000410A411 DeviceInformation = 955
	DeviceInformation_DEV0004109911 DeviceInformation = 956
	DeviceInformation_DEV0004109912 DeviceInformation = 957
	DeviceInformation_DEV0004109913 DeviceInformation = 958
	DeviceInformation_DEV0004109914 DeviceInformation = 959
	DeviceInformation_DEV000410A211 DeviceInformation = 960
	DeviceInformation_DEV000410FC12 DeviceInformation = 961
	DeviceInformation_DEV000410FD12 DeviceInformation = 962
	DeviceInformation_DEV000410B212 DeviceInformation = 963
	DeviceInformation_DEV0004110B11 DeviceInformation = 964
	DeviceInformation_DEV0004110711 DeviceInformation = 965
	DeviceInformation_DEV000410B213 DeviceInformation = 966
	DeviceInformation_DEV0004109811 DeviceInformation = 967
	DeviceInformation_DEV0004109812 DeviceInformation = 968
	DeviceInformation_DEV0004109813 DeviceInformation = 969
	DeviceInformation_DEV0004109814 DeviceInformation = 970
	DeviceInformation_DEV000410A011 DeviceInformation = 971
	DeviceInformation_DEV000410A111 DeviceInformation = 972
	DeviceInformation_DEV000410FA12 DeviceInformation = 973
	DeviceInformation_DEV000410FB12 DeviceInformation = 974
	DeviceInformation_DEV000410B112 DeviceInformation = 975
	DeviceInformation_DEV0004110A11 DeviceInformation = 976
	DeviceInformation_DEV0004110611 DeviceInformation = 977
	DeviceInformation_DEV000410B113 DeviceInformation = 978
	DeviceInformation_DEV0004109A11 DeviceInformation = 979
	DeviceInformation_DEV0004109A12 DeviceInformation = 980
	DeviceInformation_DEV0004109A13 DeviceInformation = 981
	DeviceInformation_DEV0004109A14 DeviceInformation = 982
	DeviceInformation_DEV000410A311 DeviceInformation = 983
	DeviceInformation_DEV000410B312 DeviceInformation = 984
	DeviceInformation_DEV0004110C11 DeviceInformation = 985
	DeviceInformation_DEV0004110811 DeviceInformation = 986
	DeviceInformation_DEV000410B313 DeviceInformation = 987
	DeviceInformation_DEV0004109B11 DeviceInformation = 988
	DeviceInformation_DEV0004109B12 DeviceInformation = 989
	DeviceInformation_DEV0004109B13 DeviceInformation = 990
	DeviceInformation_DEV0004109B14 DeviceInformation = 991
	DeviceInformation_DEV000410A511 DeviceInformation = 992
	DeviceInformation_DEV000410B412 DeviceInformation = 993
	DeviceInformation_DEV0004110D11 DeviceInformation = 994
	DeviceInformation_DEV0004110911 DeviceInformation = 995
	DeviceInformation_DEV000410B413 DeviceInformation = 996
	DeviceInformation_DEV0004109C11 DeviceInformation = 997
	DeviceInformation_DEV0004109C12 DeviceInformation = 998
	DeviceInformation_DEV0004109C13 DeviceInformation = 999
	DeviceInformation_DEV0004109C14 DeviceInformation = 1000
	DeviceInformation_DEV000410A611 DeviceInformation = 1001
	DeviceInformation_DEV0004146B13 DeviceInformation = 1002
	DeviceInformation_DEV0004147211 DeviceInformation = 1003
	DeviceInformation_DEV000410FE12 DeviceInformation = 1004
	DeviceInformation_DEV0004209016 DeviceInformation = 1005
	DeviceInformation_DEV000420A011 DeviceInformation = 1006
	DeviceInformation_DEV0004209011 DeviceInformation = 1007
	DeviceInformation_DEV000420CA11 DeviceInformation = 1008
	DeviceInformation_DEV0004208012 DeviceInformation = 1009
	DeviceInformation_DEV0004207812 DeviceInformation = 1010
	DeviceInformation_DEV000420BA11 DeviceInformation = 1011
	DeviceInformation_DEV000420B311 DeviceInformation = 1012
	DeviceInformation_DEV0004209811 DeviceInformation = 1013
	DeviceInformation_DEV0004208811 DeviceInformation = 1014
	DeviceInformation_DEV0004B00812 DeviceInformation = 1015
	DeviceInformation_DEV0004302613 DeviceInformation = 1016
	DeviceInformation_DEV0004302313 DeviceInformation = 1017
	DeviceInformation_DEV0004302013 DeviceInformation = 1018
	DeviceInformation_DEV0004302B12 DeviceInformation = 1019
	DeviceInformation_DEV0004706811 DeviceInformation = 1020
	DeviceInformation_DEV0004705D11 DeviceInformation = 1021
	DeviceInformation_DEV0004705C11 DeviceInformation = 1022
	DeviceInformation_DEV0004B00713 DeviceInformation = 1023
	DeviceInformation_DEV0004B00A01 DeviceInformation = 1024
	DeviceInformation_DEV0004706611 DeviceInformation = 1025
	DeviceInformation_DEV0004C01410 DeviceInformation = 1026
	DeviceInformation_DEV0004C00102 DeviceInformation = 1027
	DeviceInformation_DEV0004705E11 DeviceInformation = 1028
	DeviceInformation_DEV0004706211 DeviceInformation = 1029
	DeviceInformation_DEV0004706411 DeviceInformation = 1030
	DeviceInformation_DEV0004706412 DeviceInformation = 1031
	DeviceInformation_DEV000420C011 DeviceInformation = 1032
	DeviceInformation_DEV000420B011 DeviceInformation = 1033
	DeviceInformation_DEV0004B00911 DeviceInformation = 1034
	DeviceInformation_DEV0004A01211 DeviceInformation = 1035
	DeviceInformation_DEV0004A01112 DeviceInformation = 1036
	DeviceInformation_DEV0004A01111 DeviceInformation = 1037
	DeviceInformation_DEV0004A01212 DeviceInformation = 1038
	DeviceInformation_DEV0004A03312 DeviceInformation = 1039
	DeviceInformation_DEV0004A03212 DeviceInformation = 1040
	DeviceInformation_DEV0004A01113 DeviceInformation = 1041
	DeviceInformation_DEV0004A01711 DeviceInformation = 1042
	DeviceInformation_DEV000420C711 DeviceInformation = 1043
	DeviceInformation_DEV000420BD11 DeviceInformation = 1044
	DeviceInformation_DEV000420C411 DeviceInformation = 1045
	DeviceInformation_DEV000420A812 DeviceInformation = 1046
	DeviceInformation_DEV000420CD11 DeviceInformation = 1047
	DeviceInformation_DEV000420AD11 DeviceInformation = 1048
	DeviceInformation_DEV000420B611 DeviceInformation = 1049
	DeviceInformation_DEV000420A811 DeviceInformation = 1050
	DeviceInformation_DEV0004501311 DeviceInformation = 1051
	DeviceInformation_DEV0004501411 DeviceInformation = 1052
	DeviceInformation_DEV0004B01002 DeviceInformation = 1053
	DeviceInformation_DEV0006D00610 DeviceInformation = 1054
	DeviceInformation_DEV0006D01510 DeviceInformation = 1055
	DeviceInformation_DEV0006D00110 DeviceInformation = 1056
	DeviceInformation_DEV0006D00310 DeviceInformation = 1057
	DeviceInformation_DEV0006D03210 DeviceInformation = 1058
	DeviceInformation_DEV0006D03310 DeviceInformation = 1059
	DeviceInformation_DEV0006D02E20 DeviceInformation = 1060
	DeviceInformation_DEV0006D02F20 DeviceInformation = 1061
	DeviceInformation_DEV0006D03020 DeviceInformation = 1062
	DeviceInformation_DEV0006D03120 DeviceInformation = 1063
	DeviceInformation_DEV0006D02110 DeviceInformation = 1064
	DeviceInformation_DEV0006D00010 DeviceInformation = 1065
	DeviceInformation_DEV0006D01810 DeviceInformation = 1066
	DeviceInformation_DEV0006D00910 DeviceInformation = 1067
	DeviceInformation_DEV0006D01110 DeviceInformation = 1068
	DeviceInformation_DEV0006D03510 DeviceInformation = 1069
	DeviceInformation_DEV0006D03410 DeviceInformation = 1070
	DeviceInformation_DEV0006D02410 DeviceInformation = 1071
	DeviceInformation_DEV0006D02510 DeviceInformation = 1072
	DeviceInformation_DEV0006D00810 DeviceInformation = 1073
	DeviceInformation_DEV0006D00710 DeviceInformation = 1074
	DeviceInformation_DEV0006D01310 DeviceInformation = 1075
	DeviceInformation_DEV0006D01410 DeviceInformation = 1076
	DeviceInformation_DEV0006D00210 DeviceInformation = 1077
	DeviceInformation_DEV0006D00510 DeviceInformation = 1078
	DeviceInformation_DEV0006D00410 DeviceInformation = 1079
	DeviceInformation_DEV0006D02210 DeviceInformation = 1080
	DeviceInformation_DEV0006D02310 DeviceInformation = 1081
	DeviceInformation_DEV0006D01710 DeviceInformation = 1082
	DeviceInformation_DEV0006D01610 DeviceInformation = 1083
	DeviceInformation_DEV0006D01010 DeviceInformation = 1084
	DeviceInformation_DEV0006D01210 DeviceInformation = 1085
	DeviceInformation_DEV0006D04820 DeviceInformation = 1086
	DeviceInformation_DEV0006D04C11 DeviceInformation = 1087
	DeviceInformation_DEV0006D05610 DeviceInformation = 1088
	DeviceInformation_DEV0006D02910 DeviceInformation = 1089
	DeviceInformation_DEV0006D02A10 DeviceInformation = 1090
	DeviceInformation_DEV0006D02B10 DeviceInformation = 1091
	DeviceInformation_DEV0006D02C10 DeviceInformation = 1092
	DeviceInformation_DEV0006D02810 DeviceInformation = 1093
	DeviceInformation_DEV0006D02610 DeviceInformation = 1094
	DeviceInformation_DEV0006D02710 DeviceInformation = 1095
	DeviceInformation_DEV0006D03610 DeviceInformation = 1096
	DeviceInformation_DEV0006D03710 DeviceInformation = 1097
	DeviceInformation_DEV0006D02D11 DeviceInformation = 1098
	DeviceInformation_DEV0006D03C10 DeviceInformation = 1099
	DeviceInformation_DEV0006D03B10 DeviceInformation = 1100
	DeviceInformation_DEV0006D03910 DeviceInformation = 1101
	DeviceInformation_DEV0006D03A10 DeviceInformation = 1102
	DeviceInformation_DEV0006D03D11 DeviceInformation = 1103
	DeviceInformation_DEV0006D03E10 DeviceInformation = 1104
	DeviceInformation_DEV0006C00102 DeviceInformation = 1105
	DeviceInformation_DEV0006E05611 DeviceInformation = 1106
	DeviceInformation_DEV0006E05212 DeviceInformation = 1107
	DeviceInformation_DEV000620B011 DeviceInformation = 1108
	DeviceInformation_DEV000620B311 DeviceInformation = 1109
	DeviceInformation_DEV000620C011 DeviceInformation = 1110
	DeviceInformation_DEV000620BA11 DeviceInformation = 1111
	DeviceInformation_DEV0006705C11 DeviceInformation = 1112
	DeviceInformation_DEV0006705D11 DeviceInformation = 1113
	DeviceInformation_DEV0006E07710 DeviceInformation = 1114
	DeviceInformation_DEV0006E07712 DeviceInformation = 1115
	DeviceInformation_DEV0006706210 DeviceInformation = 1116
	DeviceInformation_DEV0006302611 DeviceInformation = 1117
	DeviceInformation_DEV0006302612 DeviceInformation = 1118
	DeviceInformation_DEV0006E00810 DeviceInformation = 1119
	DeviceInformation_DEV0006E01F01 DeviceInformation = 1120
	DeviceInformation_DEV0006302311 DeviceInformation = 1121
	DeviceInformation_DEV0006302312 DeviceInformation = 1122
	DeviceInformation_DEV0006E00910 DeviceInformation = 1123
	DeviceInformation_DEV0006E02001 DeviceInformation = 1124
	DeviceInformation_DEV0006302011 DeviceInformation = 1125
	DeviceInformation_DEV0006302012 DeviceInformation = 1126
	DeviceInformation_DEV0006C00C13 DeviceInformation = 1127
	DeviceInformation_DEV0006E00811 DeviceInformation = 1128
	DeviceInformation_DEV0006E00911 DeviceInformation = 1129
	DeviceInformation_DEV0006E01F20 DeviceInformation = 1130
	DeviceInformation_DEV0006E03410 DeviceInformation = 1131
	DeviceInformation_DEV0006E03110 DeviceInformation = 1132
	DeviceInformation_DEV0006E0A210 DeviceInformation = 1133
	DeviceInformation_DEV0006E0CE10 DeviceInformation = 1134
	DeviceInformation_DEV0006E0A111 DeviceInformation = 1135
	DeviceInformation_DEV0006E0CD11 DeviceInformation = 1136
	DeviceInformation_DEV0006E02020 DeviceInformation = 1137
	DeviceInformation_DEV0006E02D11 DeviceInformation = 1138
	DeviceInformation_DEV0006E03011 DeviceInformation = 1139
	DeviceInformation_DEV0006E0C110 DeviceInformation = 1140
	DeviceInformation_DEV0006E0C510 DeviceInformation = 1141
	DeviceInformation_DEV0006B00A01 DeviceInformation = 1142
	DeviceInformation_DEV0006B00602 DeviceInformation = 1143
	DeviceInformation_DEV0006E0C410 DeviceInformation = 1144
	DeviceInformation_DEV0006E0C312 DeviceInformation = 1145
	DeviceInformation_DEV0006E0C210 DeviceInformation = 1146
	DeviceInformation_DEV0006209016 DeviceInformation = 1147
	DeviceInformation_DEV0006E01A01 DeviceInformation = 1148
	DeviceInformation_DEV0006E09910 DeviceInformation = 1149
	DeviceInformation_DEV0006E03710 DeviceInformation = 1150
	DeviceInformation_DEV0006209011 DeviceInformation = 1151
	DeviceInformation_DEV000620A011 DeviceInformation = 1152
	DeviceInformation_DEV0006E02410 DeviceInformation = 1153
	DeviceInformation_DEV0006E02301 DeviceInformation = 1154
	DeviceInformation_DEV0006E02510 DeviceInformation = 1155
	DeviceInformation_DEV0006E01B01 DeviceInformation = 1156
	DeviceInformation_DEV0006E01C01 DeviceInformation = 1157
	DeviceInformation_DEV0006E01D01 DeviceInformation = 1158
	DeviceInformation_DEV0006E01E01 DeviceInformation = 1159
	DeviceInformation_DEV0006207812 DeviceInformation = 1160
	DeviceInformation_DEV0006B00811 DeviceInformation = 1161
	DeviceInformation_DEV0006E01001 DeviceInformation = 1162
	DeviceInformation_DEV0006E03610 DeviceInformation = 1163
	DeviceInformation_DEV0006E09810 DeviceInformation = 1164
	DeviceInformation_DEV0006208811 DeviceInformation = 1165
	DeviceInformation_DEV0006209811 DeviceInformation = 1166
	DeviceInformation_DEV0006E02610 DeviceInformation = 1167
	DeviceInformation_DEV0006E02710 DeviceInformation = 1168
	DeviceInformation_DEV0006E02A10 DeviceInformation = 1169
	DeviceInformation_DEV0006E02B10 DeviceInformation = 1170
	DeviceInformation_DEV0006E00C10 DeviceInformation = 1171
	DeviceInformation_DEV0006010110 DeviceInformation = 1172
	DeviceInformation_DEV0006010210 DeviceInformation = 1173
	DeviceInformation_DEV0006E00B10 DeviceInformation = 1174
	DeviceInformation_DEV0006E09C10 DeviceInformation = 1175
	DeviceInformation_DEV0006E09B10 DeviceInformation = 1176
	DeviceInformation_DEV0006E03510 DeviceInformation = 1177
	DeviceInformation_DEV0006FF1B11 DeviceInformation = 1178
	DeviceInformation_DEV0006E0CF10 DeviceInformation = 1179
	DeviceInformation_DEV000620A812 DeviceInformation = 1180
	DeviceInformation_DEV000620CD11 DeviceInformation = 1181
	DeviceInformation_DEV0006E00E01 DeviceInformation = 1182
	DeviceInformation_DEV0006E02201 DeviceInformation = 1183
	DeviceInformation_DEV000620AD11 DeviceInformation = 1184
	DeviceInformation_DEV0006E00F01 DeviceInformation = 1185
	DeviceInformation_DEV0006E02101 DeviceInformation = 1186
	DeviceInformation_DEV000620BD11 DeviceInformation = 1187
	DeviceInformation_DEV0006E00D01 DeviceInformation = 1188
	DeviceInformation_DEV0006E03910 DeviceInformation = 1189
	DeviceInformation_DEV0006E02810 DeviceInformation = 1190
	DeviceInformation_DEV0006E02910 DeviceInformation = 1191
	DeviceInformation_DEV0006E02C10 DeviceInformation = 1192
	DeviceInformation_DEV0006C00403 DeviceInformation = 1193
	DeviceInformation_DEV0006590101 DeviceInformation = 1194
	DeviceInformation_DEV0006E0CC11 DeviceInformation = 1195
	DeviceInformation_DEV0006E09A10 DeviceInformation = 1196
	DeviceInformation_DEV0006E03811 DeviceInformation = 1197
	DeviceInformation_DEV0006E0C710 DeviceInformation = 1198
	DeviceInformation_DEV0006E0C610 DeviceInformation = 1199
	DeviceInformation_DEV0006E05A10 DeviceInformation = 1200
	DeviceInformation_DEV0048493A1C DeviceInformation = 1201
	DeviceInformation_DEV0048494712 DeviceInformation = 1202
	DeviceInformation_DEV0048494810 DeviceInformation = 1203
	DeviceInformation_DEV0048855A10 DeviceInformation = 1204
	DeviceInformation_DEV0048855B10 DeviceInformation = 1205
	DeviceInformation_DEV0048A05713 DeviceInformation = 1206
	DeviceInformation_DEV0048494414 DeviceInformation = 1207
	DeviceInformation_DEV0048824A11 DeviceInformation = 1208
	DeviceInformation_DEV0048824A12 DeviceInformation = 1209
	DeviceInformation_DEV0048770A10 DeviceInformation = 1210
	DeviceInformation_DEV0048494311 DeviceInformation = 1211
	DeviceInformation_DEV0048494513 DeviceInformation = 1212
	DeviceInformation_DEV0048494012 DeviceInformation = 1213
	DeviceInformation_DEV0048494111 DeviceInformation = 1214
	DeviceInformation_DEV0048494210 DeviceInformation = 1215
	DeviceInformation_DEV0048494610 DeviceInformation = 1216
	DeviceInformation_DEV0048494910 DeviceInformation = 1217
	DeviceInformation_DEV0048134A10 DeviceInformation = 1218
	DeviceInformation_DEV0048107E12 DeviceInformation = 1219
	DeviceInformation_DEV0048FF2112 DeviceInformation = 1220
	DeviceInformation_DEV0048140A11 DeviceInformation = 1221
	DeviceInformation_DEV0048140B12 DeviceInformation = 1222
	DeviceInformation_DEV0048140C13 DeviceInformation = 1223
	DeviceInformation_DEV0048139A10 DeviceInformation = 1224
	DeviceInformation_DEV0048648B10 DeviceInformation = 1225
	DeviceInformation_DEV0008A01111 DeviceInformation = 1226
	DeviceInformation_DEV0008A01211 DeviceInformation = 1227
	DeviceInformation_DEV0008A01212 DeviceInformation = 1228
	DeviceInformation_DEV0008A01112 DeviceInformation = 1229
	DeviceInformation_DEV0008A03213 DeviceInformation = 1230
	DeviceInformation_DEV0008A03313 DeviceInformation = 1231
	DeviceInformation_DEV0008A01113 DeviceInformation = 1232
	DeviceInformation_DEV0008A01711 DeviceInformation = 1233
	DeviceInformation_DEV0008B00911 DeviceInformation = 1234
	DeviceInformation_DEV0008C00102 DeviceInformation = 1235
	DeviceInformation_DEV0008C00101 DeviceInformation = 1236
	DeviceInformation_DEV0008901501 DeviceInformation = 1237
	DeviceInformation_DEV0008901310 DeviceInformation = 1238
	DeviceInformation_DEV000820B011 DeviceInformation = 1239
	DeviceInformation_DEV0008705C11 DeviceInformation = 1240
	DeviceInformation_DEV0008705D11 DeviceInformation = 1241
	DeviceInformation_DEV0008706211 DeviceInformation = 1242
	DeviceInformation_DEV000820BA11 DeviceInformation = 1243
	DeviceInformation_DEV000820C011 DeviceInformation = 1244
	DeviceInformation_DEV000820B311 DeviceInformation = 1245
	DeviceInformation_DEV0008301A11 DeviceInformation = 1246
	DeviceInformation_DEV0008C00C13 DeviceInformation = 1247
	DeviceInformation_DEV0008302611 DeviceInformation = 1248
	DeviceInformation_DEV0008302311 DeviceInformation = 1249
	DeviceInformation_DEV0008302011 DeviceInformation = 1250
	DeviceInformation_DEV0008C00C11 DeviceInformation = 1251
	DeviceInformation_DEV0008302612 DeviceInformation = 1252
	DeviceInformation_DEV0008302312 DeviceInformation = 1253
	DeviceInformation_DEV0008302012 DeviceInformation = 1254
	DeviceInformation_DEV0008C00C15 DeviceInformation = 1255
	DeviceInformation_DEV0008C00C14 DeviceInformation = 1256
	DeviceInformation_DEV0008B00713 DeviceInformation = 1257
	DeviceInformation_DEV0008706611 DeviceInformation = 1258
	DeviceInformation_DEV0008706811 DeviceInformation = 1259
	DeviceInformation_DEV0008B00812 DeviceInformation = 1260
	DeviceInformation_DEV0008209016 DeviceInformation = 1261
	DeviceInformation_DEV0008209011 DeviceInformation = 1262
	DeviceInformation_DEV000820A011 DeviceInformation = 1263
	DeviceInformation_DEV0008208811 DeviceInformation = 1264
	DeviceInformation_DEV0008209811 DeviceInformation = 1265
	DeviceInformation_DEV000820CA11 DeviceInformation = 1266
	DeviceInformation_DEV0008208012 DeviceInformation = 1267
	DeviceInformation_DEV0008207812 DeviceInformation = 1268
	DeviceInformation_DEV0008207811 DeviceInformation = 1269
	DeviceInformation_DEV0008208011 DeviceInformation = 1270
	DeviceInformation_DEV000810D111 DeviceInformation = 1271
	DeviceInformation_DEV000810D511 DeviceInformation = 1272
	DeviceInformation_DEV000810FA12 DeviceInformation = 1273
	DeviceInformation_DEV000810FB12 DeviceInformation = 1274
	DeviceInformation_DEV000810F211 DeviceInformation = 1275
	DeviceInformation_DEV000810D211 DeviceInformation = 1276
	DeviceInformation_DEV000810E211 DeviceInformation = 1277
	DeviceInformation_DEV000810D611 DeviceInformation = 1278
	DeviceInformation_DEV000810F212 DeviceInformation = 1279
	DeviceInformation_DEV000810E212 DeviceInformation = 1280
	DeviceInformation_DEV000810FC13 DeviceInformation = 1281
	DeviceInformation_DEV000810FD13 DeviceInformation = 1282
	DeviceInformation_DEV000810F311 DeviceInformation = 1283
	DeviceInformation_DEV000810D311 DeviceInformation = 1284
	DeviceInformation_DEV000810D711 DeviceInformation = 1285
	DeviceInformation_DEV000810F312 DeviceInformation = 1286
	DeviceInformation_DEV000810D811 DeviceInformation = 1287
	DeviceInformation_DEV000810E511 DeviceInformation = 1288
	DeviceInformation_DEV000810E512 DeviceInformation = 1289
	DeviceInformation_DEV000810F611 DeviceInformation = 1290
	DeviceInformation_DEV000810D911 DeviceInformation = 1291
	DeviceInformation_DEV000810F612 DeviceInformation = 1292
	DeviceInformation_DEV000820A812 DeviceInformation = 1293
	DeviceInformation_DEV000820AD11 DeviceInformation = 1294
	DeviceInformation_DEV000820BD11 DeviceInformation = 1295
	DeviceInformation_DEV000820C711 DeviceInformation = 1296
	DeviceInformation_DEV000820CD11 DeviceInformation = 1297
	DeviceInformation_DEV000820C411 DeviceInformation = 1298
	DeviceInformation_DEV000820A811 DeviceInformation = 1299
	DeviceInformation_DEV0008501411 DeviceInformation = 1300
	DeviceInformation_DEV0008C01602 DeviceInformation = 1301
	DeviceInformation_DEV0008302613 DeviceInformation = 1302
	DeviceInformation_DEV0008302313 DeviceInformation = 1303
	DeviceInformation_DEV0008302013 DeviceInformation = 1304
	DeviceInformation_DEV0009207730 DeviceInformation = 1305
	DeviceInformation_DEV0009208F10 DeviceInformation = 1306
	DeviceInformation_DEV0009C00C13 DeviceInformation = 1307
	DeviceInformation_DEV0009209910 DeviceInformation = 1308
	DeviceInformation_DEV0009209A10 DeviceInformation = 1309
	DeviceInformation_DEV0009207930 DeviceInformation = 1310
	DeviceInformation_DEV0009201720 DeviceInformation = 1311
	DeviceInformation_DEV0009500D01 DeviceInformation = 1312
	DeviceInformation_DEV0009500E01 DeviceInformation = 1313
	DeviceInformation_DEV0009209911 DeviceInformation = 1314
	DeviceInformation_DEV0009209A11 DeviceInformation = 1315
	DeviceInformation_DEV0009C00C12 DeviceInformation = 1316
	DeviceInformation_DEV0009C00C11 DeviceInformation = 1317
	DeviceInformation_DEV0009500D20 DeviceInformation = 1318
	DeviceInformation_DEV0009500E20 DeviceInformation = 1319
	DeviceInformation_DEV000920B910 DeviceInformation = 1320
	DeviceInformation_DEV0009E0CE10 DeviceInformation = 1321
	DeviceInformation_DEV0009E0A210 DeviceInformation = 1322
	DeviceInformation_DEV0009501410 DeviceInformation = 1323
	DeviceInformation_DEV0009207830 DeviceInformation = 1324
	DeviceInformation_DEV0009201620 DeviceInformation = 1325
	DeviceInformation_DEV0009E0A111 DeviceInformation = 1326
	DeviceInformation_DEV0009E0CD11 DeviceInformation = 1327
	DeviceInformation_DEV000920B811 DeviceInformation = 1328
	DeviceInformation_DEV000920B611 DeviceInformation = 1329
	DeviceInformation_DEV0009207E10 DeviceInformation = 1330
	DeviceInformation_DEV0009207630 DeviceInformation = 1331
	DeviceInformation_DEV0009205910 DeviceInformation = 1332
	DeviceInformation_DEV0009500B01 DeviceInformation = 1333
	DeviceInformation_DEV000920AC10 DeviceInformation = 1334
	DeviceInformation_DEV0009207430 DeviceInformation = 1335
	DeviceInformation_DEV0009204521 DeviceInformation = 1336
	DeviceInformation_DEV0009500A01 DeviceInformation = 1337
	DeviceInformation_DEV0009500001 DeviceInformation = 1338
	DeviceInformation_DEV000920AB10 DeviceInformation = 1339
	DeviceInformation_DEV000920BF11 DeviceInformation = 1340
	DeviceInformation_DEV0009203510 DeviceInformation = 1341
	DeviceInformation_DEV0009207A30 DeviceInformation = 1342
	DeviceInformation_DEV0009500701 DeviceInformation = 1343
	DeviceInformation_DEV0009501710 DeviceInformation = 1344
	DeviceInformation_DEV000920B310 DeviceInformation = 1345
	DeviceInformation_DEV0009207530 DeviceInformation = 1346
	DeviceInformation_DEV0009203321 DeviceInformation = 1347
	DeviceInformation_DEV0009500C01 DeviceInformation = 1348
	DeviceInformation_DEV000920AD10 DeviceInformation = 1349
	DeviceInformation_DEV0009207230 DeviceInformation = 1350
	DeviceInformation_DEV0009500801 DeviceInformation = 1351
	DeviceInformation_DEV0009501810 DeviceInformation = 1352
	DeviceInformation_DEV000920B410 DeviceInformation = 1353
	DeviceInformation_DEV0009207330 DeviceInformation = 1354
	DeviceInformation_DEV0009204421 DeviceInformation = 1355
	DeviceInformation_DEV0009500901 DeviceInformation = 1356
	DeviceInformation_DEV000920AA10 DeviceInformation = 1357
	DeviceInformation_DEV0009209D01 DeviceInformation = 1358
	DeviceInformation_DEV000920B010 DeviceInformation = 1359
	DeviceInformation_DEV0009E0BE01 DeviceInformation = 1360
	DeviceInformation_DEV000920B110 DeviceInformation = 1361
	DeviceInformation_DEV0009E0BD01 DeviceInformation = 1362
	DeviceInformation_DEV0009D03F10 DeviceInformation = 1363
	DeviceInformation_DEV0009305F10 DeviceInformation = 1364
	DeviceInformation_DEV0009305610 DeviceInformation = 1365
	DeviceInformation_DEV0009D03E10 DeviceInformation = 1366
	DeviceInformation_DEV0009306010 DeviceInformation = 1367
	DeviceInformation_DEV0009306110 DeviceInformation = 1368
	DeviceInformation_DEV0009306310 DeviceInformation = 1369
	DeviceInformation_DEV0009D03B10 DeviceInformation = 1370
	DeviceInformation_DEV0009D03C10 DeviceInformation = 1371
	DeviceInformation_DEV0009D03910 DeviceInformation = 1372
	DeviceInformation_DEV0009D03A10 DeviceInformation = 1373
	DeviceInformation_DEV0009305411 DeviceInformation = 1374
	DeviceInformation_DEV0009D03D11 DeviceInformation = 1375
	DeviceInformation_DEV0009304B11 DeviceInformation = 1376
	DeviceInformation_DEV0009304C11 DeviceInformation = 1377
	DeviceInformation_DEV0009306220 DeviceInformation = 1378
	DeviceInformation_DEV0009302E10 DeviceInformation = 1379
	DeviceInformation_DEV0009302F10 DeviceInformation = 1380
	DeviceInformation_DEV0009303010 DeviceInformation = 1381
	DeviceInformation_DEV0009303110 DeviceInformation = 1382
	DeviceInformation_DEV0009306510 DeviceInformation = 1383
	DeviceInformation_DEV0009306610 DeviceInformation = 1384
	DeviceInformation_DEV0009306410 DeviceInformation = 1385
	DeviceInformation_DEV0009401110 DeviceInformation = 1386
	DeviceInformation_DEV0009400610 DeviceInformation = 1387
	DeviceInformation_DEV0009401510 DeviceInformation = 1388
	DeviceInformation_DEV0009402110 DeviceInformation = 1389
	DeviceInformation_DEV0009400110 DeviceInformation = 1390
	DeviceInformation_DEV0009400910 DeviceInformation = 1391
	DeviceInformation_DEV0009400010 DeviceInformation = 1392
	DeviceInformation_DEV0009401810 DeviceInformation = 1393
	DeviceInformation_DEV0009400310 DeviceInformation = 1394
	DeviceInformation_DEV0009301810 DeviceInformation = 1395
	DeviceInformation_DEV0009301910 DeviceInformation = 1396
	DeviceInformation_DEV0009301A10 DeviceInformation = 1397
	DeviceInformation_DEV0009401210 DeviceInformation = 1398
	DeviceInformation_DEV0009400810 DeviceInformation = 1399
	DeviceInformation_DEV0009400710 DeviceInformation = 1400
	DeviceInformation_DEV0009401310 DeviceInformation = 1401
	DeviceInformation_DEV0009401410 DeviceInformation = 1402
	DeviceInformation_DEV0009402210 DeviceInformation = 1403
	DeviceInformation_DEV0009402310 DeviceInformation = 1404
	DeviceInformation_DEV0009401710 DeviceInformation = 1405
	DeviceInformation_DEV0009401610 DeviceInformation = 1406
	DeviceInformation_DEV0009400210 DeviceInformation = 1407
	DeviceInformation_DEV0009401010 DeviceInformation = 1408
	DeviceInformation_DEV0009400510 DeviceInformation = 1409
	DeviceInformation_DEV0009400410 DeviceInformation = 1410
	DeviceInformation_DEV0009D04B20 DeviceInformation = 1411
	DeviceInformation_DEV0009D04920 DeviceInformation = 1412
	DeviceInformation_DEV0009D04A20 DeviceInformation = 1413
	DeviceInformation_DEV0009D04820 DeviceInformation = 1414
	DeviceInformation_DEV0009D04C11 DeviceInformation = 1415
	DeviceInformation_DEV0009D05610 DeviceInformation = 1416
	DeviceInformation_DEV0009305510 DeviceInformation = 1417
	DeviceInformation_DEV0009209810 DeviceInformation = 1418
	DeviceInformation_DEV0009202A10 DeviceInformation = 1419
	DeviceInformation_DEV0009209510 DeviceInformation = 1420
	DeviceInformation_DEV0009501110 DeviceInformation = 1421
	DeviceInformation_DEV0009209310 DeviceInformation = 1422
	DeviceInformation_DEV0009209410 DeviceInformation = 1423
	DeviceInformation_DEV0009209210 DeviceInformation = 1424
	DeviceInformation_DEV0009501210 DeviceInformation = 1425
	DeviceInformation_DEV0009205411 DeviceInformation = 1426
	DeviceInformation_DEV000920A111 DeviceInformation = 1427
	DeviceInformation_DEV000920A311 DeviceInformation = 1428
	DeviceInformation_DEV0009205112 DeviceInformation = 1429
	DeviceInformation_DEV0009204110 DeviceInformation = 1430
	DeviceInformation_DEV0009E07710 DeviceInformation = 1431
	DeviceInformation_DEV0009E07712 DeviceInformation = 1432
	DeviceInformation_DEV0009205212 DeviceInformation = 1433
	DeviceInformation_DEV0009205211 DeviceInformation = 1434
	DeviceInformation_DEV0009205311 DeviceInformation = 1435
	DeviceInformation_DEV0009206B10 DeviceInformation = 1436
	DeviceInformation_DEV0009208010 DeviceInformation = 1437
	DeviceInformation_DEV0009206A12 DeviceInformation = 1438
	DeviceInformation_DEV0009206810 DeviceInformation = 1439
	DeviceInformation_DEV0009208110 DeviceInformation = 1440
	DeviceInformation_DEV0009205511 DeviceInformation = 1441
	DeviceInformation_DEV0009209F01 DeviceInformation = 1442
	DeviceInformation_DEV0009208C10 DeviceInformation = 1443
	DeviceInformation_DEV0009208E10 DeviceInformation = 1444
	DeviceInformation_DEV000920B511 DeviceInformation = 1445
	DeviceInformation_DEV0009501910 DeviceInformation = 1446
	DeviceInformation_DEV000920BE11 DeviceInformation = 1447
	DeviceInformation_DEV0009209710 DeviceInformation = 1448
	DeviceInformation_DEV0009208510 DeviceInformation = 1449
	DeviceInformation_DEV0009208610 DeviceInformation = 1450
	DeviceInformation_DEV000920BD10 DeviceInformation = 1451
	DeviceInformation_DEV0009500210 DeviceInformation = 1452
	DeviceInformation_DEV0009500310 DeviceInformation = 1453
	DeviceInformation_DEV0009E0BF10 DeviceInformation = 1454
	DeviceInformation_DEV0009E0C010 DeviceInformation = 1455
	DeviceInformation_DEV0009500110 DeviceInformation = 1456
	DeviceInformation_DEV0009209B10 DeviceInformation = 1457
	DeviceInformation_DEV0009207D10 DeviceInformation = 1458
	DeviceInformation_DEV0009202F11 DeviceInformation = 1459
	DeviceInformation_DEV0009203011 DeviceInformation = 1460
	DeviceInformation_DEV0009207C10 DeviceInformation = 1461
	DeviceInformation_DEV0009207B10 DeviceInformation = 1462
	DeviceInformation_DEV0009208710 DeviceInformation = 1463
	DeviceInformation_DEV0009E06610 DeviceInformation = 1464
	DeviceInformation_DEV0009E06611 DeviceInformation = 1465
	DeviceInformation_DEV0009E06410 DeviceInformation = 1466
	DeviceInformation_DEV0009E06411 DeviceInformation = 1467
	DeviceInformation_DEV0009E06210 DeviceInformation = 1468
	DeviceInformation_DEV0009E0E910 DeviceInformation = 1469
	DeviceInformation_DEV0009E0EB10 DeviceInformation = 1470
	DeviceInformation_DEV000920BB10 DeviceInformation = 1471
	DeviceInformation_DEV0009FF1B11 DeviceInformation = 1472
	DeviceInformation_DEV0009E0CF10 DeviceInformation = 1473
	DeviceInformation_DEV0009206C30 DeviceInformation = 1474
	DeviceInformation_DEV0009206D30 DeviceInformation = 1475
	DeviceInformation_DEV0009206E30 DeviceInformation = 1476
	DeviceInformation_DEV0009206F30 DeviceInformation = 1477
	DeviceInformation_DEV0009207130 DeviceInformation = 1478
	DeviceInformation_DEV0009204720 DeviceInformation = 1479
	DeviceInformation_DEV0009204820 DeviceInformation = 1480
	DeviceInformation_DEV0009204920 DeviceInformation = 1481
	DeviceInformation_DEV0009204A20 DeviceInformation = 1482
	DeviceInformation_DEV0009205A10 DeviceInformation = 1483
	DeviceInformation_DEV0009207030 DeviceInformation = 1484
	DeviceInformation_DEV0009205B10 DeviceInformation = 1485
	DeviceInformation_DEV0009500501 DeviceInformation = 1486
	DeviceInformation_DEV0009501001 DeviceInformation = 1487
	DeviceInformation_DEV0009500601 DeviceInformation = 1488
	DeviceInformation_DEV0009500F01 DeviceInformation = 1489
	DeviceInformation_DEV0009500401 DeviceInformation = 1490
	DeviceInformation_DEV000920B210 DeviceInformation = 1491
	DeviceInformation_DEV000920AE10 DeviceInformation = 1492
	DeviceInformation_DEV000920BC10 DeviceInformation = 1493
	DeviceInformation_DEV000920AF10 DeviceInformation = 1494
	DeviceInformation_DEV0009207F10 DeviceInformation = 1495
	DeviceInformation_DEV0009208910 DeviceInformation = 1496
	DeviceInformation_DEV0009205710 DeviceInformation = 1497
	DeviceInformation_DEV0009205810 DeviceInformation = 1498
	DeviceInformation_DEV0009203810 DeviceInformation = 1499
	DeviceInformation_DEV0009203910 DeviceInformation = 1500
	DeviceInformation_DEV0009203E10 DeviceInformation = 1501
	DeviceInformation_DEV0009204B10 DeviceInformation = 1502
	DeviceInformation_DEV0009203F10 DeviceInformation = 1503
	DeviceInformation_DEV0009204C10 DeviceInformation = 1504
	DeviceInformation_DEV0009204010 DeviceInformation = 1505
	DeviceInformation_DEV0009206411 DeviceInformation = 1506
	DeviceInformation_DEV0009205E10 DeviceInformation = 1507
	DeviceInformation_DEV0009206711 DeviceInformation = 1508
	DeviceInformation_DEV000920A710 DeviceInformation = 1509
	DeviceInformation_DEV000920A610 DeviceInformation = 1510
	DeviceInformation_DEV0009203A10 DeviceInformation = 1511
	DeviceInformation_DEV0009203B10 DeviceInformation = 1512
	DeviceInformation_DEV0009203C10 DeviceInformation = 1513
	DeviceInformation_DEV0009203D10 DeviceInformation = 1514
	DeviceInformation_DEV0009E05E12 DeviceInformation = 1515
	DeviceInformation_DEV0009E0B711 DeviceInformation = 1516
	DeviceInformation_DEV0009E06A12 DeviceInformation = 1517
	DeviceInformation_DEV0009E06E12 DeviceInformation = 1518
	DeviceInformation_DEV0009E0B720 DeviceInformation = 1519
	DeviceInformation_DEV0009E0E611 DeviceInformation = 1520
	DeviceInformation_DEV0009E0B321 DeviceInformation = 1521
	DeviceInformation_DEV0009E0E512 DeviceInformation = 1522
	DeviceInformation_DEV0009204210 DeviceInformation = 1523
	DeviceInformation_DEV0009208210 DeviceInformation = 1524
	DeviceInformation_DEV0009E07211 DeviceInformation = 1525
	DeviceInformation_DEV0009E0CC11 DeviceInformation = 1526
	DeviceInformation_DEV0009110111 DeviceInformation = 1527
	DeviceInformation_DEV0009110211 DeviceInformation = 1528
	DeviceInformation_DEV000916B212 DeviceInformation = 1529
	DeviceInformation_DEV0009110212 DeviceInformation = 1530
	DeviceInformation_DEV0009110311 DeviceInformation = 1531
	DeviceInformation_DEV000916B312 DeviceInformation = 1532
	DeviceInformation_DEV0009110312 DeviceInformation = 1533
	DeviceInformation_DEV0009110411 DeviceInformation = 1534
	DeviceInformation_DEV0009110412 DeviceInformation = 1535
	DeviceInformation_DEV0009501615 DeviceInformation = 1536
	DeviceInformation_DEV0009E0ED10 DeviceInformation = 1537
	DeviceInformation_DEV014F030110 DeviceInformation = 1538
	DeviceInformation_DEV014F030310 DeviceInformation = 1539
	DeviceInformation_DEV014F030210 DeviceInformation = 1540
	DeviceInformation_DEV00EE7FFF10 DeviceInformation = 1541
	DeviceInformation_DEV00B6464101 DeviceInformation = 1542
	DeviceInformation_DEV00B6464201 DeviceInformation = 1543
	DeviceInformation_DEV00B6464501 DeviceInformation = 1544
	DeviceInformation_DEV00B6434101 DeviceInformation = 1545
	DeviceInformation_DEV00B6434201 DeviceInformation = 1546
	DeviceInformation_DEV00B6434202 DeviceInformation = 1547
	DeviceInformation_DEV00B6454101 DeviceInformation = 1548
	DeviceInformation_DEV00B6454201 DeviceInformation = 1549
	DeviceInformation_DEV00B6455001 DeviceInformation = 1550
	DeviceInformation_DEV00B6453101 DeviceInformation = 1551
	DeviceInformation_DEV00B6453102 DeviceInformation = 1552
	DeviceInformation_DEV00B6454102 DeviceInformation = 1553
	DeviceInformation_DEV00B6454401 DeviceInformation = 1554
	DeviceInformation_DEV00B6454402 DeviceInformation = 1555
	DeviceInformation_DEV00B6454202 DeviceInformation = 1556
	DeviceInformation_DEV00B6453103 DeviceInformation = 1557
	DeviceInformation_DEV00B6453201 DeviceInformation = 1558
	DeviceInformation_DEV00B6453301 DeviceInformation = 1559
	DeviceInformation_DEV00B6453104 DeviceInformation = 1560
	DeviceInformation_DEV00B6454403 DeviceInformation = 1561
	DeviceInformation_DEV00B6454801 DeviceInformation = 1562
	DeviceInformation_DEV00B6414701 DeviceInformation = 1563
	DeviceInformation_DEV00B6414201 DeviceInformation = 1564
	DeviceInformation_DEV00B6474101 DeviceInformation = 1565
	DeviceInformation_DEV00B6474302 DeviceInformation = 1566
	DeviceInformation_DEV00B6474602 DeviceInformation = 1567
	DeviceInformation_DEV00B6534D01 DeviceInformation = 1568
	DeviceInformation_DEV00B6535001 DeviceInformation = 1569
	DeviceInformation_DEV00B6455002 DeviceInformation = 1570
	DeviceInformation_DEV00B6453701 DeviceInformation = 1571
	DeviceInformation_DEV00B6484101 DeviceInformation = 1572
	DeviceInformation_DEV00B6484201 DeviceInformation = 1573
	DeviceInformation_DEV00B6484202 DeviceInformation = 1574
	DeviceInformation_DEV00B6484301 DeviceInformation = 1575
	DeviceInformation_DEV00B6484102 DeviceInformation = 1576
	DeviceInformation_DEV00B6455101 DeviceInformation = 1577
	DeviceInformation_DEV00B6455003 DeviceInformation = 1578
	DeviceInformation_DEV00B6455102 DeviceInformation = 1579
	DeviceInformation_DEV00B6453702 DeviceInformation = 1580
	DeviceInformation_DEV00B6453703 DeviceInformation = 1581
	DeviceInformation_DEV00B6484302 DeviceInformation = 1582
	DeviceInformation_DEV00B6484801 DeviceInformation = 1583
	DeviceInformation_DEV00B6484501 DeviceInformation = 1584
	DeviceInformation_DEV00B6484203 DeviceInformation = 1585
	DeviceInformation_DEV00B6484103 DeviceInformation = 1586
	DeviceInformation_DEV00B6455004 DeviceInformation = 1587
	DeviceInformation_DEV00B6455103 DeviceInformation = 1588
	DeviceInformation_DEV00B6455401 DeviceInformation = 1589
	DeviceInformation_DEV00B6455201 DeviceInformation = 1590
	DeviceInformation_DEV00B6455402 DeviceInformation = 1591
	DeviceInformation_DEV00B6455403 DeviceInformation = 1592
	DeviceInformation_DEV00B603430A DeviceInformation = 1593
	DeviceInformation_DEV00B600010A DeviceInformation = 1594
	DeviceInformation_DEV00B6FF110A DeviceInformation = 1595
	DeviceInformation_DEV00B6434601 DeviceInformation = 1596
	DeviceInformation_DEV00B6434602 DeviceInformation = 1597
	DeviceInformation_DEV00B6455301 DeviceInformation = 1598
	DeviceInformation_DEV00C5070410 DeviceInformation = 1599
	DeviceInformation_DEV00C5070210 DeviceInformation = 1600
	DeviceInformation_DEV00C5070610 DeviceInformation = 1601
	DeviceInformation_DEV00C5070E11 DeviceInformation = 1602
	DeviceInformation_DEV00C5060240 DeviceInformation = 1603
	DeviceInformation_DEV00C5062010 DeviceInformation = 1604
	DeviceInformation_DEV00C5080230 DeviceInformation = 1605
	DeviceInformation_DEV00C5060310 DeviceInformation = 1606
	DeviceInformation_DEV006C070E11 DeviceInformation = 1607
	DeviceInformation_DEV006C050002 DeviceInformation = 1608
	DeviceInformation_DEV006C011311 DeviceInformation = 1609
	DeviceInformation_DEV006C011411 DeviceInformation = 1610
	DeviceInformation_DEV0007632010 DeviceInformation = 1611
	DeviceInformation_DEV0007632020 DeviceInformation = 1612
	DeviceInformation_DEV0007632180 DeviceInformation = 1613
	DeviceInformation_DEV0007632040 DeviceInformation = 1614
	DeviceInformation_DEV0007613812 DeviceInformation = 1615
	DeviceInformation_DEV0007613810 DeviceInformation = 1616
	DeviceInformation_DEV000720C011 DeviceInformation = 1617
	DeviceInformation_DEV0007A05210 DeviceInformation = 1618
	DeviceInformation_DEV0007A08B10 DeviceInformation = 1619
	DeviceInformation_DEV0007A05B32 DeviceInformation = 1620
	DeviceInformation_DEV0007A06932 DeviceInformation = 1621
	DeviceInformation_DEV0007A06D32 DeviceInformation = 1622
	DeviceInformation_DEV0007A08032 DeviceInformation = 1623
	DeviceInformation_DEV0007A00213 DeviceInformation = 1624
	DeviceInformation_DEV0007A09532 DeviceInformation = 1625
	DeviceInformation_DEV0007A06C32 DeviceInformation = 1626
	DeviceInformation_DEV0007A05E32 DeviceInformation = 1627
	DeviceInformation_DEV0007A08A32 DeviceInformation = 1628
	DeviceInformation_DEV0007A07032 DeviceInformation = 1629
	DeviceInformation_DEV0007A08332 DeviceInformation = 1630
	DeviceInformation_DEV0007A09832 DeviceInformation = 1631
	DeviceInformation_DEV0007A05C32 DeviceInformation = 1632
	DeviceInformation_DEV0007A06A32 DeviceInformation = 1633
	DeviceInformation_DEV0007A08832 DeviceInformation = 1634
	DeviceInformation_DEV0007A06E32 DeviceInformation = 1635
	DeviceInformation_DEV0007A08132 DeviceInformation = 1636
	DeviceInformation_DEV0007A00113 DeviceInformation = 1637
	DeviceInformation_DEV0007A09632 DeviceInformation = 1638
	DeviceInformation_DEV0007A05D32 DeviceInformation = 1639
	DeviceInformation_DEV0007A06B32 DeviceInformation = 1640
	DeviceInformation_DEV0007A08932 DeviceInformation = 1641
	DeviceInformation_DEV0007A06F32 DeviceInformation = 1642
	DeviceInformation_DEV0007A08232 DeviceInformation = 1643
	DeviceInformation_DEV0007A09732 DeviceInformation = 1644
	DeviceInformation_DEV0007A05713 DeviceInformation = 1645
	DeviceInformation_DEV0007A01811 DeviceInformation = 1646
	DeviceInformation_DEV0007A01911 DeviceInformation = 1647
	DeviceInformation_DEV0007A04912 DeviceInformation = 1648
	DeviceInformation_DEV0007A05814 DeviceInformation = 1649
	DeviceInformation_DEV0007A07114 DeviceInformation = 1650
	DeviceInformation_DEV0007A05810 DeviceInformation = 1651
	DeviceInformation_DEV0007A04312 DeviceInformation = 1652
	DeviceInformation_DEV0007A04412 DeviceInformation = 1653
	DeviceInformation_DEV0007A04512 DeviceInformation = 1654
	DeviceInformation_DEV000720BD11 DeviceInformation = 1655
	DeviceInformation_DEV0007A04C13 DeviceInformation = 1656
	DeviceInformation_DEV0007A04D13 DeviceInformation = 1657
	DeviceInformation_DEV0007A04B10 DeviceInformation = 1658
	DeviceInformation_DEV0007A04E13 DeviceInformation = 1659
	DeviceInformation_DEV0007A04F13 DeviceInformation = 1660
	DeviceInformation_DEV000720BA11 DeviceInformation = 1661
	DeviceInformation_DEV0007A03D11 DeviceInformation = 1662
	DeviceInformation_DEV0007A09211 DeviceInformation = 1663
	DeviceInformation_DEV0007A09111 DeviceInformation = 1664
	DeviceInformation_DEV0007FF1115 DeviceInformation = 1665
	DeviceInformation_DEV0007A01511 DeviceInformation = 1666
	DeviceInformation_DEV0007A08411 DeviceInformation = 1667
	DeviceInformation_DEV0007A08511 DeviceInformation = 1668
	DeviceInformation_DEV0007A03422 DeviceInformation = 1669
	DeviceInformation_DEV0007A07213 DeviceInformation = 1670
	DeviceInformation_DEV0007A07420 DeviceInformation = 1671
	DeviceInformation_DEV0007A07520 DeviceInformation = 1672
	DeviceInformation_DEV0007A07B12 DeviceInformation = 1673
	DeviceInformation_DEV0007A07C12 DeviceInformation = 1674
	DeviceInformation_DEV0007A09311 DeviceInformation = 1675
	DeviceInformation_DEV0007A09013 DeviceInformation = 1676
	DeviceInformation_DEV0007A08F13 DeviceInformation = 1677
	DeviceInformation_DEV0007A07E10 DeviceInformation = 1678
	DeviceInformation_DEV0007A05510 DeviceInformation = 1679
	DeviceInformation_DEV0007A05910 DeviceInformation = 1680
	DeviceInformation_DEV0007A08711 DeviceInformation = 1681
	DeviceInformation_DEV0007A03D12 DeviceInformation = 1682
	DeviceInformation_DEV0007A09A12 DeviceInformation = 1683
	DeviceInformation_DEV0007A09B12 DeviceInformation = 1684
	DeviceInformation_DEV0007A06614 DeviceInformation = 1685
	DeviceInformation_DEV0007A06514 DeviceInformation = 1686
	DeviceInformation_DEV0007A06014 DeviceInformation = 1687
	DeviceInformation_DEV0007A07714 DeviceInformation = 1688
	DeviceInformation_DEV0007A06414 DeviceInformation = 1689
	DeviceInformation_DEV0007A06114 DeviceInformation = 1690
	DeviceInformation_DEV0007A07814 DeviceInformation = 1691
	DeviceInformation_DEV0007A06714 DeviceInformation = 1692
	DeviceInformation_DEV0007A06214 DeviceInformation = 1693
	DeviceInformation_DEV0007A07914 DeviceInformation = 1694
	DeviceInformation_DEV000B0A8410 DeviceInformation = 1695
	DeviceInformation_DEV000B0A7E10 DeviceInformation = 1696
	DeviceInformation_DEV000B0A7F10 DeviceInformation = 1697
	DeviceInformation_DEV000B0A8010 DeviceInformation = 1698
	DeviceInformation_DEV000BBF9111 DeviceInformation = 1699
	DeviceInformation_DEV000B0A7810 DeviceInformation = 1700
	DeviceInformation_DEV000B0A7910 DeviceInformation = 1701
	DeviceInformation_DEV000B0A7A10 DeviceInformation = 1702
	DeviceInformation_DEV000B0A8910 DeviceInformation = 1703
	DeviceInformation_DEV000B0A8310 DeviceInformation = 1704
	DeviceInformation_DEV000B0A8510 DeviceInformation = 1705
	DeviceInformation_DEV000B0A6319 DeviceInformation = 1706
)

func (e DeviceInformation) ComObjectTableAddress() uint16 {
	switch e {
	case 1:
		{ /* '1' */
			return 0x43FE
		}
	case 10:
		{ /* '10' */
			return 0x43FE
		}
	case 100:
		{ /* '100' */
			return 0x4400
		}
	case 1000:
		{ /* '1000' */
			return 0x4195
		}
	case 1001:
		{ /* '1001' */
			return 0x41E5
		}
	case 1002:
		{ /* '1002' */
			return 0x43FF
		}
	case 1003:
		{ /* '1003' */
			return 0x43FF
		}
	case 1004:
		{ /* '1004' */
			return 0x43FF
		}
	case 1005:
		{ /* '1005' */
			return 0x43FF
		}
	case 1006:
		{ /* '1006' */
			return 0x43FF
		}
	case 1007:
		{ /* '1007' */
			return 0x43FF
		}
	case 1008:
		{ /* '1008' */
			return 0x43FF
		}
	case 1009:
		{ /* '1009' */
			return 0x43FF
		}
	case 101:
		{ /* '101' */
			return 0x4400
		}
	case 1010:
		{ /* '1010' */
			return 0x43FF
		}
	case 1011:
		{ /* '1011' */
			return 0x43FF
		}
	case 1012:
		{ /* '1012' */
			return 0x43FF
		}
	case 1013:
		{ /* '1013' */
			return 0x43FF
		}
	case 1014:
		{ /* '1014' */
			return 0x43FF
		}
	case 1015:
		{ /* '1015' */
			return 0x4324
		}
	case 1016:
		{ /* '1016' */
			return 0x43FF
		}
	case 1017:
		{ /* '1017' */
			return 0x43FF
		}
	case 1018:
		{ /* '1018' */
			return 0x43FF
		}
	case 1019:
		{ /* '1019' */
			return 0x4194
		}
	case 102:
		{ /* '102' */
			return 0x4400
		}
	case 1020:
		{ /* '1020' */
			return 0x43FE
		}
	case 1021:
		{ /* '1021' */
			return 0x43FE
		}
	case 1022:
		{ /* '1022' */
			return 0x43FE
		}
	case 1023:
		{ /* '1023' */
			return 0x4324
		}
	case 1024:
		{ /* '1024' */
			return 0x4324
		}
	case 1025:
		{ /* '1025' */
			return 0x43FE
		}
	case 1026:
		{ /* '1026' */
			return 0x43EC
		}
	case 1027:
		{ /* '1027' */
			return 0x41C8
		}
	case 1028:
		{ /* '1028' */
			return 0x43FE
		}
	case 1029:
		{ /* '1029' */
			return 0x43FF
		}
	case 103:
		{ /* '103' */
			return 0x4400
		}
	case 1030:
		{ /* '1030' */
			return 0x43FF
		}
	case 1031:
		{ /* '1031' */
			return 0x43FF
		}
	case 1032:
		{ /* '1032' */
			return 0x43FF
		}
	case 1033:
		{ /* '1033' */
			return 0x43FF
		}
	case 1034:
		{ /* '1034' */
			return 0x4324
		}
	case 1035:
		{ /* '1035' */
			return 0x43FF
		}
	case 1036:
		{ /* '1036' */
			return 0x43FF
		}
	case 1037:
		{ /* '1037' */
			return 0x43FF
		}
	case 1038:
		{ /* '1038' */
			return 0x43FF
		}
	case 1039:
		{ /* '1039' */
			return 0x43FF
		}
	case 104:
		{ /* '104' */
			return 0x4400
		}
	case 1040:
		{ /* '1040' */
			return 0x43FF
		}
	case 1041:
		{ /* '1041' */
			return 0x43FF
		}
	case 1042:
		{ /* '1042' */
			return 0x43FF
		}
	case 1043:
		{ /* '1043' */
			return 0x43FF
		}
	case 1044:
		{ /* '1044' */
			return 0x43FF
		}
	case 1045:
		{ /* '1045' */
			return 0x43FF
		}
	case 1046:
		{ /* '1046' */
			return 0x43FF
		}
	case 1047:
		{ /* '1047' */
			return 0x43FF
		}
	case 1048:
		{ /* '1048' */
			return 0x43FF
		}
	case 1049:
		{ /* '1049' */
			return 0x43FF
		}
	case 105:
		{ /* '105' */
			return 0x43FE
		}
	case 1050:
		{ /* '1050' */
			return 0x43FF
		}
	case 1051:
		{ /* '1051' */
			return 0x43FF
		}
	case 1052:
		{ /* '1052' */
			return 0x8700
		}
	case 1053:
		{ /* '1053' */
			return 0x4324
		}
	case 1054:
		{ /* '1054' */
			return 0x4000
		}
	case 1055:
		{ /* '1055' */
			return 0x4000
		}
	case 1056:
		{ /* '1056' */
			return 0x4000
		}
	case 1057:
		{ /* '1057' */
			return 0x4000
		}
	case 1058:
		{ /* '1058' */
			return 0x4000
		}
	case 1059:
		{ /* '1059' */
			return 0x4000
		}
	case 106:
		{ /* '106' */
			return 0x402C
		}
	case 1060:
		{ /* '1060' */
			return 0x4000
		}
	case 1061:
		{ /* '1061' */
			return 0x4000
		}
	case 1062:
		{ /* '1062' */
			return 0x4000
		}
	case 1063:
		{ /* '1063' */
			return 0x4000
		}
	case 1064:
		{ /* '1064' */
			return 0x4000
		}
	case 1065:
		{ /* '1065' */
			return 0x4000
		}
	case 1066:
		{ /* '1066' */
			return 0x4000
		}
	case 1067:
		{ /* '1067' */
			return 0x4000
		}
	case 1068:
		{ /* '1068' */
			return 0x4000
		}
	case 1069:
		{ /* '1069' */
			return 0x4000
		}
	case 107:
		{ /* '107' */
			return 0x43FE
		}
	case 1070:
		{ /* '1070' */
			return 0x4000
		}
	case 1071:
		{ /* '1071' */
			return 0x4000
		}
	case 1072:
		{ /* '1072' */
			return 0x4000
		}
	case 1073:
		{ /* '1073' */
			return 0x4000
		}
	case 1074:
		{ /* '1074' */
			return 0x4000
		}
	case 1075:
		{ /* '1075' */
			return 0x4000
		}
	case 1076:
		{ /* '1076' */
			return 0x4000
		}
	case 1077:
		{ /* '1077' */
			return 0x4000
		}
	case 1078:
		{ /* '1078' */
			return 0x4000
		}
	case 1079:
		{ /* '1079' */
			return 0x4000
		}
	case 108:
		{ /* '108' */
			return 0x4400
		}
	case 1080:
		{ /* '1080' */
			return 0x4000
		}
	case 1081:
		{ /* '1081' */
			return 0x4000
		}
	case 1082:
		{ /* '1082' */
			return 0x4000
		}
	case 1083:
		{ /* '1083' */
			return 0x4000
		}
	case 1084:
		{ /* '1084' */
			return 0x4000
		}
	case 1085:
		{ /* '1085' */
			return 0x4000
		}
	case 1086:
		{ /* '1086' */
			return 0x4000
		}
	case 1087:
		{ /* '1087' */
			return 0x4000
		}
	case 1088:
		{ /* '1088' */
			return 0x4000
		}
	case 1089:
		{ /* '1089' */
			return 0x4000
		}
	case 109:
		{ /* '109' */
			return 0x43FE
		}
	case 1090:
		{ /* '1090' */
			return 0x4000
		}
	case 1091:
		{ /* '1091' */
			return 0x4000
		}
	case 1092:
		{ /* '1092' */
			return 0x4000
		}
	case 1093:
		{ /* '1093' */
			return 0x4000
		}
	case 1094:
		{ /* '1094' */
			return 0x4000
		}
	case 1095:
		{ /* '1095' */
			return 0x4000
		}
	case 1096:
		{ /* '1096' */
			return 0x4000
		}
	case 1097:
		{ /* '1097' */
			return 0x4000
		}
	case 1098:
		{ /* '1098' */
			return 0x4000
		}
	case 1099:
		{ /* '1099' */
			return 0x4000
		}
	case 11:
		{ /* '11' */
			return 0x43FE
		}
	case 110:
		{ /* '110' */
			return 0x4400
		}
	case 1100:
		{ /* '1100' */
			return 0x4000
		}
	case 1101:
		{ /* '1101' */
			return 0x7000
		}
	case 1102:
		{ /* '1102' */
			return 0x7000
		}
	case 1103:
		{ /* '1103' */
			return 0x4000
		}
	case 1104:
		{ /* '1104' */
			return 0x4000
		}
	case 1105:
		{ /* '1105' */
			return 0x41C8
		}
	case 1106:
		{ /* '1106' */
			return 0x43FE
		}
	case 1107:
		{ /* '1107' */
			return 0x43FE
		}
	case 1108:
		{ /* '1108' */
			return 0x43FF
		}
	case 1109:
		{ /* '1109' */
			return 0x43FF
		}
	case 111:
		{ /* '111' */
			return 0x43CE
		}
	case 1110:
		{ /* '1110' */
			return 0x43FF
		}
	case 1111:
		{ /* '1111' */
			return 0x43FF
		}
	case 1112:
		{ /* '1112' */
			return 0x43FE
		}
	case 1113:
		{ /* '1113' */
			return 0x43FE
		}
	case 1114:
		{ /* '1114' */
			return 0x43FE
		}
	case 1115:
		{ /* '1115' */
			return 0x43FE
		}
	case 1116:
		{ /* '1116' */
			return 0x43FF
		}
	case 1117:
		{ /* '1117' */
			return 0x43FF
		}
	case 1118:
		{ /* '1118' */
			return 0x43FF
		}
	case 1119:
		{ /* '1119' */
			return 0x4400
		}
	case 112:
		{ /* '112' */
			return 0x4400
		}
	case 1120:
		{ /* '1120' */
			return 0x4400
		}
	case 1121:
		{ /* '1121' */
			return 0x43FF
		}
	case 1122:
		{ /* '1122' */
			return 0x43FF
		}
	case 1123:
		{ /* '1123' */
			return 0x4400
		}
	case 1124:
		{ /* '1124' */
			return 0x4400
		}
	case 1125:
		{ /* '1125' */
			return 0x43FF
		}
	case 1126:
		{ /* '1126' */
			return 0x43FF
		}
	case 1127:
		{ /* '1127' */
			return 0x43FF
		}
	case 1128:
		{ /* '1128' */
			return 0x4400
		}
	case 1129:
		{ /* '1129' */
			return 0x4400
		}
	case 113:
		{ /* '113' */
			return 0x4400
		}
	case 1130:
		{ /* '1130' */
			return 0x4400
		}
	case 1131:
		{ /* '1131' */
			return 0x4400
		}
	case 1132:
		{ /* '1132' */
			return 0x4400
		}
	case 1133:
		{ /* '1133' */
			return 0x4400
		}
	case 1134:
		{ /* '1134' */
			return 0x4400
		}
	case 1135:
		{ /* '1135' */
			return 0x4400
		}
	case 1136:
		{ /* '1136' */
			return 0x4400
		}
	case 1137:
		{ /* '1137' */
			return 0x4400
		}
	case 1138:
		{ /* '1138' */
			return 0x4400
		}
	case 1139:
		{ /* '1139' */
			return 0x4400
		}
	case 114:
		{ /* '114' */
			return 0x4400
		}
	case 1140:
		{ /* '1140' */
			return 0x43FE
		}
	case 1141:
		{ /* '1141' */
			return 0x43FE
		}
	case 1142:
		{ /* '1142' */
			return 0x4324
		}
	case 1143:
		{ /* '1143' */
			return 0x4193
		}
	case 1144:
		{ /* '1144' */
			return 0x43FE
		}
	case 1145:
		{ /* '1145' */
			return 0x43FE
		}
	case 1146:
		{ /* '1146' */
			return 0x43FE
		}
	case 1147:
		{ /* '1147' */
			return 0x43FF
		}
	case 1148:
		{ /* '1148' */
			return 0x4400
		}
	case 1149:
		{ /* '1149' */
			return 0x4400
		}
	case 115:
		{ /* '115' */
			return 0x4400
		}
	case 1150:
		{ /* '1150' */
			return 0x4400
		}
	case 1151:
		{ /* '1151' */
			return 0x43FF
		}
	case 1152:
		{ /* '1152' */
			return 0x43FF
		}
	case 1153:
		{ /* '1153' */
			return 0x4400
		}
	case 1154:
		{ /* '1154' */
			return 0x4400
		}
	case 1155:
		{ /* '1155' */
			return 0x4400
		}
	case 1156:
		{ /* '1156' */
			return 0x4400
		}
	case 1157:
		{ /* '1157' */
			return 0x4400
		}
	case 1158:
		{ /* '1158' */
			return 0x4400
		}
	case 1159:
		{ /* '1159' */
			return 0x4400
		}
	case 116:
		{ /* '116' */
			return 0x4400
		}
	case 1160:
		{ /* '1160' */
			return 0x43FF
		}
	case 1161:
		{ /* '1161' */
			return 0x4324
		}
	case 1162:
		{ /* '1162' */
			return 0x4400
		}
	case 1163:
		{ /* '1163' */
			return 0x4400
		}
	case 1164:
		{ /* '1164' */
			return 0x4400
		}
	case 1165:
		{ /* '1165' */
			return 0x43FF
		}
	case 1166:
		{ /* '1166' */
			return 0x43FF
		}
	case 1167:
		{ /* '1167' */
			return 0x4400
		}
	case 1168:
		{ /* '1168' */
			return 0x4400
		}
	case 1169:
		{ /* '1169' */
			return 0x4400
		}
	case 117:
		{ /* '117' */
			return 0x43FF
		}
	case 1170:
		{ /* '1170' */
			return 0x4400
		}
	case 1171:
		{ /* '1171' */
			return 0x43FE
		}
	case 1172:
		{ /* '1172' */
			return 0x43FE
		}
	case 1173:
		{ /* '1173' */
			return 0x43FE
		}
	case 1174:
		{ /* '1174' */
			return 0x43FE
		}
	case 1175:
		{ /* '1175' */
			return 0x43FE
		}
	case 1176:
		{ /* '1176' */
			return 0x43FE
		}
	case 1177:
		{ /* '1177' */
			return 0x4400
		}
	case 1178:
		{ /* '1178' */
			return 0x43FC
		}
	case 1179:
		{ /* '1179' */
			return 0x4400
		}
	case 118:
		{ /* '118' */
			return 0x43FF
		}
	case 1180:
		{ /* '1180' */
			return 0x43FF
		}
	case 1181:
		{ /* '1181' */
			return 0x43FF
		}
	case 1182:
		{ /* '1182' */
			return 0x4400
		}
	case 1183:
		{ /* '1183' */
			return 0x4400
		}
	case 1184:
		{ /* '1184' */
			return 0x43FF
		}
	case 1185:
		{ /* '1185' */
			return 0x4400
		}
	case 1186:
		{ /* '1186' */
			return 0x4400
		}
	case 1187:
		{ /* '1187' */
			return 0x43FF
		}
	case 1188:
		{ /* '1188' */
			return 0x4400
		}
	case 1189:
		{ /* '1189' */
			return 0x4400
		}
	case 119:
		{ /* '119' */
			return 0x4400
		}
	case 1190:
		{ /* '1190' */
			return 0x4400
		}
	case 1191:
		{ /* '1191' */
			return 0x4400
		}
	case 1192:
		{ /* '1192' */
			return 0x4400
		}
	case 1193:
		{ /* '1193' */
			return 0x43FC
		}
	case 1194:
		{ /* '1194' */
			return 0x4000
		}
	case 1195:
		{ /* '1195' */
			return 0x43FE
		}
	case 1196:
		{ /* '1196' */
			return 0x4400
		}
	case 1197:
		{ /* '1197' */
			return 0x4400
		}
	case 1198:
		{ /* '1198' */
			return 0x43FE
		}
	case 1199:
		{ /* '1199' */
			return 0x43FE
		}
	case 12:
		{ /* '12' */
			return 0x43FE
		}
	case 120:
		{ /* '120' */
			return 0x43FE
		}
	case 1200:
		{ /* '1200' */
			return 0x43FE
		}
	case 1201:
		{ /* '1201' */
			return 0x43FE
		}
	case 1202:
		{ /* '1202' */
			return 0x43FE
		}
	case 1203:
		{ /* '1203' */
			return 0x43FE
		}
	case 1204:
		{ /* '1204' */
			return 0x419C
		}
	case 1205:
		{ /* '1205' */
			return 0x419C
		}
	case 1206:
		{ /* '1206' */
			return 0x48D6
		}
	case 1207:
		{ /* '1207' */
			return 0x43FE
		}
	case 1208:
		{ /* '1208' */
			return 0x426C
		}
	case 1209:
		{ /* '1209' */
			return 0x426C
		}
	case 121:
		{ /* '121' */
			return 0x4204
		}
	case 1210:
		{ /* '1210' */
			return 0x4204
		}
	case 1211:
		{ /* '1211' */
			return 0x43FE
		}
	case 1212:
		{ /* '1212' */
			return 0x43FE
		}
	case 1213:
		{ /* '1213' */
			return 0x43FE
		}
	case 1214:
		{ /* '1214' */
			return 0x43FE
		}
	case 1215:
		{ /* '1215' */
			return 0x43FE
		}
	case 1216:
		{ /* '1216' */
			return 0x43FE
		}
	case 1217:
		{ /* '1217' */
			return 0x43FE
		}
	case 1218:
		{ /* '1218' */
			return 0x4400
		}
	case 1219:
		{ /* '1219' */
			return 0x43FC
		}
	case 122:
		{ /* '122' */
			return 0x4324
		}
	case 1220:
		{ /* '1220' */
			return 0x4204
		}
	case 1221:
		{ /* '1221' */
			return 0x4400
		}
	case 1222:
		{ /* '1222' */
			return 0x4400
		}
	case 1223:
		{ /* '1223' */
			return 0x4400
		}
	case 1224:
		{ /* '1224' */
			return 0x4324
		}
	case 1225:
		{ /* '1225' */
			return 0x4400
		}
	case 1226:
		{ /* '1226' */
			return 0x43FF
		}
	case 1227:
		{ /* '1227' */
			return 0x43FF
		}
	case 1228:
		{ /* '1228' */
			return 0x43FF
		}
	case 1229:
		{ /* '1229' */
			return 0x43FF
		}
	case 123:
		{ /* '123' */
			return 0x4324
		}
	case 1230:
		{ /* '1230' */
			return 0x43FF
		}
	case 1231:
		{ /* '1231' */
			return 0x43FF
		}
	case 1232:
		{ /* '1232' */
			return 0x43FF
		}
	case 1233:
		{ /* '1233' */
			return 0x43FF
		}
	case 1234:
		{ /* '1234' */
			return 0x4324
		}
	case 1235:
		{ /* '1235' */
			return 0x41C8
		}
	case 1236:
		{ /* '1236' */
			return 0x41C8
		}
	case 1237:
		{ /* '1237' */
			return 0x40F4
		}
	case 1238:
		{ /* '1238' */
			return 0x40F4
		}
	case 1239:
		{ /* '1239' */
			return 0x43FF
		}
	case 124:
		{ /* '124' */
			return 0x43FE
		}
	case 1240:
		{ /* '1240' */
			return 0x43FE
		}
	case 1241:
		{ /* '1241' */
			return 0x43FE
		}
	case 1242:
		{ /* '1242' */
			return 0x43FF
		}
	case 1243:
		{ /* '1243' */
			return 0x43FF
		}
	case 1244:
		{ /* '1244' */
			return 0x43FF
		}
	case 1245:
		{ /* '1245' */
			return 0x43FF
		}
	case 1246:
		{ /* '1246' */
			return 0x43FF
		}
	case 1247:
		{ /* '1247' */
			return 0x43FF
		}
	case 1248:
		{ /* '1248' */
			return 0x43FF
		}
	case 1249:
		{ /* '1249' */
			return 0x43FF
		}
	case 125:
		{ /* '125' */
			return 0x43FE
		}
	case 1250:
		{ /* '1250' */
			return 0x43FF
		}
	case 1251:
		{ /* '1251' */
			return 0x43FF
		}
	case 1252:
		{ /* '1252' */
			return 0x43FF
		}
	case 1253:
		{ /* '1253' */
			return 0x43FF
		}
	case 1254:
		{ /* '1254' */
			return 0x43FF
		}
	case 1255:
		{ /* '1255' */
			return 0x43FF
		}
	case 1256:
		{ /* '1256' */
			return 0x43FF
		}
	case 1257:
		{ /* '1257' */
			return 0x4324
		}
	case 1258:
		{ /* '1258' */
			return 0x43FE
		}
	case 1259:
		{ /* '1259' */
			return 0x43FE
		}
	case 126:
		{ /* '126' */
			return 0x43FF
		}
	case 1260:
		{ /* '1260' */
			return 0x4324
		}
	case 1261:
		{ /* '1261' */
			return 0x43FF
		}
	case 1262:
		{ /* '1262' */
			return 0x43FF
		}
	case 1263:
		{ /* '1263' */
			return 0x43FF
		}
	case 1264:
		{ /* '1264' */
			return 0x43FF
		}
	case 1265:
		{ /* '1265' */
			return 0x43FF
		}
	case 1266:
		{ /* '1266' */
			return 0x43FF
		}
	case 1267:
		{ /* '1267' */
			return 0x43FF
		}
	case 1268:
		{ /* '1268' */
			return 0x43FF
		}
	case 1269:
		{ /* '1269' */
			return 0x43FF
		}
	case 127:
		{ /* '127' */
			return 0x4400
		}
	case 1270:
		{ /* '1270' */
			return 0x43FF
		}
	case 1271:
		{ /* '1271' */
			return 0x4195
		}
	case 1272:
		{ /* '1272' */
			return 0x41E5
		}
	case 1273:
		{ /* '1273' */
			return 0x4195
		}
	case 1274:
		{ /* '1274' */
			return 0x4195
		}
	case 1275:
		{ /* '1275' */
			return 0x43FF
		}
	case 1276:
		{ /* '1276' */
			return 0x4195
		}
	case 1277:
		{ /* '1277' */
			return 0x43FF
		}
	case 1278:
		{ /* '1278' */
			return 0x41E5
		}
	case 1279:
		{ /* '1279' */
			return 0x43FF
		}
	case 128:
		{ /* '128' */
			return 0x4400
		}
	case 1280:
		{ /* '1280' */
			return 0x43FF
		}
	case 1281:
		{ /* '1281' */
			return 0x4195
		}
	case 1282:
		{ /* '1282' */
			return 0x4195
		}
	case 1283:
		{ /* '1283' */
			return 0x43FF
		}
	case 1284:
		{ /* '1284' */
			return 0x4195
		}
	case 1285:
		{ /* '1285' */
			return 0x41E5
		}
	case 1286:
		{ /* '1286' */
			return 0x43FF
		}
	case 1287:
		{ /* '1287' */
			return 0x41E5
		}
	case 1288:
		{ /* '1288' */
			return 0x43FF
		}
	case 1289:
		{ /* '1289' */
			return 0x43FF
		}
	case 129:
		{ /* '129' */
			return 0x40F4
		}
	case 1290:
		{ /* '1290' */
			return 0x43FF
		}
	case 1291:
		{ /* '1291' */
			return 0x41E5
		}
	case 1292:
		{ /* '1292' */
			return 0x43FF
		}
	case 1293:
		{ /* '1293' */
			return 0x43FF
		}
	case 1294:
		{ /* '1294' */
			return 0x43FF
		}
	case 1295:
		{ /* '1295' */
			return 0x43FF
		}
	case 1296:
		{ /* '1296' */
			return 0x43FF
		}
	case 1297:
		{ /* '1297' */
			return 0x43FF
		}
	case 1298:
		{ /* '1298' */
			return 0x43FF
		}
	case 1299:
		{ /* '1299' */
			return 0x43FF
		}
	case 13:
		{ /* '13' */
			return 0x43FE
		}
	case 130:
		{ /* '130' */
			return 0x40F4
		}
	case 1300:
		{ /* '1300' */
			return 0x8700
		}
	case 1301:
		{ /* '1301' */
			return 0x4094
		}
	case 1302:
		{ /* '1302' */
			return 0x43FF
		}
	case 1303:
		{ /* '1303' */
			return 0x43FF
		}
	case 1304:
		{ /* '1304' */
			return 0x43FF
		}
	case 1305:
		{ /* '1305' */
			return 0x43FE
		}
	case 1306:
		{ /* '1306' */
			return 0x43FE
		}
	case 1307:
		{ /* '1307' */
			return 0x43FF
		}
	case 1308:
		{ /* '1308' */
			return 0x4400
		}
	case 1309:
		{ /* '1309' */
			return 0x4400
		}
	case 131:
		{ /* '131' */
			return 0x4400
		}
	case 1310:
		{ /* '1310' */
			return 0x43FE
		}
	case 1311:
		{ /* '1311' */
			return 0x43F8
		}
	case 1312:
		{ /* '1312' */
			return 0x4400
		}
	case 1313:
		{ /* '1313' */
			return 0x4400
		}
	case 1314:
		{ /* '1314' */
			return 0x4400
		}
	case 1315:
		{ /* '1315' */
			return 0x4400
		}
	case 1316:
		{ /* '1316' */
			return 0x43FF
		}
	case 1317:
		{ /* '1317' */
			return 0x43FF
		}
	case 1318:
		{ /* '1318' */
			return 0x4400
		}
	case 1319:
		{ /* '1319' */
			return 0x4400
		}
	case 132:
		{ /* '132' */
			return 0x4400
		}
	case 1320:
		{ /* '1320' */
			return 0x4400
		}
	case 1321:
		{ /* '1321' */
			return 0x4400
		}
	case 1322:
		{ /* '1322' */
			return 0x4400
		}
	case 1323:
		{ /* '1323' */
			return 0x4400
		}
	case 1324:
		{ /* '1324' */
			return 0x43FE
		}
	case 1325:
		{ /* '1325' */
			return 0x43F8
		}
	case 1326:
		{ /* '1326' */
			return 0x4400
		}
	case 1327:
		{ /* '1327' */
			return 0x4400
		}
	case 1328:
		{ /* '1328' */
			return 0x4400
		}
	case 1329:
		{ /* '1329' */
			return 0x4400
		}
	case 133:
		{ /* '133' */
			return 0x4400
		}
	case 1330:
		{ /* '1330' */
			return 0x43FE
		}
	case 1331:
		{ /* '1331' */
			return 0x43FE
		}
	case 1332:
		{ /* '1332' */
			return 0x43FE
		}
	case 1333:
		{ /* '1333' */
			return 0x4400
		}
	case 1334:
		{ /* '1334' */
			return 0x4400
		}
	case 1335:
		{ /* '1335' */
			return 0x43FE
		}
	case 1336:
		{ /* '1336' */
			return 0x43FE
		}
	case 1337:
		{ /* '1337' */
			return 0x4400
		}
	case 1338:
		{ /* '1338' */
			return 0x4400
		}
	case 1339:
		{ /* '1339' */
			return 0x4400
		}
	case 134:
		{ /* '134' */
			return 0x4400
		}
	case 1340:
		{ /* '1340' */
			return 0x43FF
		}
	case 1341:
		{ /* '1341' */
			return 0x43FE
		}
	case 1342:
		{ /* '1342' */
			return 0x43FE
		}
	case 1343:
		{ /* '1343' */
			return 0x4400
		}
	case 1344:
		{ /* '1344' */
			return 0x4400
		}
	case 1345:
		{ /* '1345' */
			return 0x4400
		}
	case 1346:
		{ /* '1346' */
			return 0x43FE
		}
	case 1347:
		{ /* '1347' */
			return 0x43FE
		}
	case 1348:
		{ /* '1348' */
			return 0x4400
		}
	case 1349:
		{ /* '1349' */
			return 0x4400
		}
	case 135:
		{ /* '135' */
			return 0x425C
		}
	case 1350:
		{ /* '1350' */
			return 0x43FE
		}
	case 1351:
		{ /* '1351' */
			return 0x4400
		}
	case 1352:
		{ /* '1352' */
			return 0x4400
		}
	case 1353:
		{ /* '1353' */
			return 0x4400
		}
	case 1354:
		{ /* '1354' */
			return 0x43FE
		}
	case 1355:
		{ /* '1355' */
			return 0x43FE
		}
	case 1356:
		{ /* '1356' */
			return 0x4400
		}
	case 1357:
		{ /* '1357' */
			return 0x4400
		}
	case 1358:
		{ /* '1358' */
			return 0x4324
		}
	case 1359:
		{ /* '1359' */
			return 0x4400
		}
	case 136:
		{ /* '136' */
			return 0x40F0
		}
	case 1360:
		{ /* '1360' */
			return 0x43F4
		}
	case 1361:
		{ /* '1361' */
			return 0x4400
		}
	case 1362:
		{ /* '1362' */
			return 0x43F4
		}
	case 1363:
		{ /* '1363' */
			return 0x4000
		}
	case 1364:
		{ /* '1364' */
			return 0x4000
		}
	case 1365:
		{ /* '1365' */
			return 0x4000
		}
	case 1366:
		{ /* '1366' */
			return 0x4000
		}
	case 1367:
		{ /* '1367' */
			return 0x4000
		}
	case 1368:
		{ /* '1368' */
			return 0x4000
		}
	case 1369:
		{ /* '1369' */
			return 0x4000
		}
	case 137:
		{ /* '137' */
			return 0x4324
		}
	case 1370:
		{ /* '1370' */
			return 0x4000
		}
	case 1371:
		{ /* '1371' */
			return 0x4000
		}
	case 1372:
		{ /* '1372' */
			return 0x7000
		}
	case 1373:
		{ /* '1373' */
			return 0x7000
		}
	case 1374:
		{ /* '1374' */
			return 0x4000
		}
	case 1375:
		{ /* '1375' */
			return 0x4000
		}
	case 1376:
		{ /* '1376' */
			return 0x4000
		}
	case 1377:
		{ /* '1377' */
			return 0x4000
		}
	case 1378:
		{ /* '1378' */
			return 0x4000
		}
	case 1379:
		{ /* '1379' */
			return 0x4000
		}
	case 138:
		{ /* '138' */
			return 0x4324
		}
	case 1380:
		{ /* '1380' */
			return 0x4000
		}
	case 1381:
		{ /* '1381' */
			return 0x4000
		}
	case 1382:
		{ /* '1382' */
			return 0x4000
		}
	case 1383:
		{ /* '1383' */
			return 0x4000
		}
	case 1384:
		{ /* '1384' */
			return 0x4000
		}
	case 1385:
		{ /* '1385' */
			return 0x4000
		}
	case 1386:
		{ /* '1386' */
			return 0x4000
		}
	case 1387:
		{ /* '1387' */
			return 0x4000
		}
	case 1388:
		{ /* '1388' */
			return 0x4000
		}
	case 1389:
		{ /* '1389' */
			return 0x4000
		}
	case 139:
		{ /* '139' */
			return 0x4324
		}
	case 1390:
		{ /* '1390' */
			return 0x4000
		}
	case 1391:
		{ /* '1391' */
			return 0x4000
		}
	case 1392:
		{ /* '1392' */
			return 0x4000
		}
	case 1393:
		{ /* '1393' */
			return 0x4000
		}
	case 1394:
		{ /* '1394' */
			return 0x4000
		}
	case 1395:
		{ /* '1395' */
			return 0x4000
		}
	case 1396:
		{ /* '1396' */
			return 0x4000
		}
	case 1397:
		{ /* '1397' */
			return 0x4000
		}
	case 1398:
		{ /* '1398' */
			return 0x4000
		}
	case 1399:
		{ /* '1399' */
			return 0x4000
		}
	case 14:
		{ /* '14' */
			return 0x43FF
		}
	case 140:
		{ /* '140' */
			return 0x4324
		}
	case 1400:
		{ /* '1400' */
			return 0x4000
		}
	case 1401:
		{ /* '1401' */
			return 0x4000
		}
	case 1402:
		{ /* '1402' */
			return 0x4000
		}
	case 1403:
		{ /* '1403' */
			return 0x4000
		}
	case 1404:
		{ /* '1404' */
			return 0x4000
		}
	case 1405:
		{ /* '1405' */
			return 0x4000
		}
	case 1406:
		{ /* '1406' */
			return 0x4000
		}
	case 1407:
		{ /* '1407' */
			return 0x4000
		}
	case 1408:
		{ /* '1408' */
			return 0x4000
		}
	case 1409:
		{ /* '1409' */
			return 0x4000
		}
	case 141:
		{ /* '141' */
			return 0x4284
		}
	case 1410:
		{ /* '1410' */
			return 0x4000
		}
	case 1411:
		{ /* '1411' */
			return 0x4000
		}
	case 1412:
		{ /* '1412' */
			return 0x4000
		}
	case 1413:
		{ /* '1413' */
			return 0x4000
		}
	case 1414:
		{ /* '1414' */
			return 0x4000
		}
	case 1415:
		{ /* '1415' */
			return 0x4000
		}
	case 1416:
		{ /* '1416' */
			return 0x4000
		}
	case 1417:
		{ /* '1417' */
			return 0x4000
		}
	case 1418:
		{ /* '1418' */
			return 0x43FE
		}
	case 1419:
		{ /* '1419' */
			return 0x43FE
		}
	case 142:
		{ /* '142' */
			return 0x4284
		}
	case 1420:
		{ /* '1420' */
			return 0x43FE
		}
	case 1421:
		{ /* '1421' */
			return 0x43FE
		}
	case 1422:
		{ /* '1422' */
			return 0x4400
		}
	case 1423:
		{ /* '1423' */
			return 0x4400
		}
	case 1424:
		{ /* '1424' */
			return 0x43FE
		}
	case 1425:
		{ /* '1425' */
			return 0x43FE
		}
	case 1426:
		{ /* '1426' */
			return 0x43FE
		}
	case 1427:
		{ /* '1427' */
			return 0x43FF
		}
	case 1428:
		{ /* '1428' */
			return 0x43FF
		}
	case 1429:
		{ /* '1429' */
			return 0x43FE
		}
	case 143:
		{ /* '143' */
			return 0x4284
		}
	case 1430:
		{ /* '1430' */
			return 0x43FE
		}
	case 1431:
		{ /* '1431' */
			return 0x43FE
		}
	case 1432:
		{ /* '1432' */
			return 0x43FE
		}
	case 1433:
		{ /* '1433' */
			return 0x43FE
		}
	case 1434:
		{ /* '1434' */
			return 0x43FE
		}
	case 1435:
		{ /* '1435' */
			return 0x43FE
		}
	case 1436:
		{ /* '1436' */
			return 0x43FE
		}
	case 1437:
		{ /* '1437' */
			return 0x43FE
		}
	case 1438:
		{ /* '1438' */
			return 0x43FE
		}
	case 1439:
		{ /* '1439' */
			return 0x43FE
		}
	case 144:
		{ /* '144' */
			return 0x4324
		}
	case 1440:
		{ /* '1440' */
			return 0x43FE
		}
	case 1441:
		{ /* '1441' */
			return 0x43FE
		}
	case 1442:
		{ /* '1442' */
			return 0x4324
		}
	case 1443:
		{ /* '1443' */
			return 0x43FE
		}
	case 1444:
		{ /* '1444' */
			return 0x43FE
		}
	case 1445:
		{ /* '1445' */
			return 0x4400
		}
	case 1446:
		{ /* '1446' */
			return 0x4400
		}
	case 1447:
		{ /* '1447' */
			return 0x43FF
		}
	case 1448:
		{ /* '1448' */
			return 0x4400
		}
	case 1449:
		{ /* '1449' */
			return 0x43FC
		}
	case 145:
		{ /* '145' */
			return 0x4324
		}
	case 1450:
		{ /* '1450' */
			return 0x43FC
		}
	case 1451:
		{ /* '1451' */
			return 0x4400
		}
	case 1452:
		{ /* '1452' */
			return 0x407A
		}
	case 1453:
		{ /* '1453' */
			return 0x407A
		}
	case 1454:
		{ /* '1454' */
			return 0x4400
		}
	case 1455:
		{ /* '1455' */
			return 0x4400
		}
	case 1456:
		{ /* '1456' */
			return 0x43FE
		}
	case 1457:
		{ /* '1457' */
			return 0x43FE
		}
	case 1458:
		{ /* '1458' */
			return 0x43FE
		}
	case 1459:
		{ /* '1459' */
			return 0x43FE
		}
	case 146:
		{ /* '146' */
			return 0x43FC
		}
	case 1460:
		{ /* '1460' */
			return 0x43FE
		}
	case 1461:
		{ /* '1461' */
			return 0x43FE
		}
	case 1462:
		{ /* '1462' */
			return 0x43FE
		}
	case 1463:
		{ /* '1463' */
			return 0x43FE
		}
	case 1464:
		{ /* '1464' */
			return 0x43FE
		}
	case 1465:
		{ /* '1465' */
			return 0x43FE
		}
	case 1466:
		{ /* '1466' */
			return 0x43FE
		}
	case 1467:
		{ /* '1467' */
			return 0x43FE
		}
	case 1468:
		{ /* '1468' */
			return 0x43FE
		}
	case 1469:
		{ /* '1469' */
			return 0x43FE
		}
	case 147:
		{ /* '147' */
			return 0x43FC
		}
	case 1470:
		{ /* '1470' */
			return 0x43FE
		}
	case 1471:
		{ /* '1471' */
			return 0x4400
		}
	case 1472:
		{ /* '1472' */
			return 0x43FC
		}
	case 1473:
		{ /* '1473' */
			return 0x4400
		}
	case 1474:
		{ /* '1474' */
			return 0x43FE
		}
	case 1475:
		{ /* '1475' */
			return 0x43FE
		}
	case 1476:
		{ /* '1476' */
			return 0x43FE
		}
	case 1477:
		{ /* '1477' */
			return 0x43FE
		}
	case 1478:
		{ /* '1478' */
			return 0x43FE
		}
	case 1479:
		{ /* '1479' */
			return 0x43FE
		}
	case 148:
		{ /* '148' */
			return 0x43FC
		}
	case 1480:
		{ /* '1480' */
			return 0x43FE
		}
	case 1481:
		{ /* '1481' */
			return 0x43FE
		}
	case 1482:
		{ /* '1482' */
			return 0x43FE
		}
	case 1483:
		{ /* '1483' */
			return 0x43FE
		}
	case 1484:
		{ /* '1484' */
			return 0x43FE
		}
	case 1485:
		{ /* '1485' */
			return 0x43FE
		}
	case 1486:
		{ /* '1486' */
			return 0x4400
		}
	case 1487:
		{ /* '1487' */
			return 0x4400
		}
	case 1488:
		{ /* '1488' */
			return 0x4400
		}
	case 1489:
		{ /* '1489' */
			return 0x4400
		}
	case 149:
		{ /* '149' */
			return 0x43FC
		}
	case 1490:
		{ /* '1490' */
			return 0x4400
		}
	case 1491:
		{ /* '1491' */
			return 0x4400
		}
	case 1492:
		{ /* '1492' */
			return 0x4400
		}
	case 1493:
		{ /* '1493' */
			return 0x4400
		}
	case 1494:
		{ /* '1494' */
			return 0x4400
		}
	case 1495:
		{ /* '1495' */
			return 0x43FE
		}
	case 1496:
		{ /* '1496' */
			return 0x43FE
		}
	case 1497:
		{ /* '1497' */
			return 0x43FE
		}
	case 1498:
		{ /* '1498' */
			return 0x43FE
		}
	case 1499:
		{ /* '1499' */
			return 0x43FE
		}
	case 15:
		{ /* '15' */
			return 0x43FF
		}
	case 150:
		{ /* '150' */
			return 0x43FC
		}
	case 1500:
		{ /* '1500' */
			return 0x43FE
		}
	case 1501:
		{ /* '1501' */
			return 0x43FE
		}
	case 1502:
		{ /* '1502' */
			return 0x43FE
		}
	case 1503:
		{ /* '1503' */
			return 0x43FE
		}
	case 1504:
		{ /* '1504' */
			return 0x43FE
		}
	case 1505:
		{ /* '1505' */
			return 0x43FE
		}
	case 1506:
		{ /* '1506' */
			return 0x43FE
		}
	case 1507:
		{ /* '1507' */
			return 0x43FE
		}
	case 1508:
		{ /* '1508' */
			return 0x43FE
		}
	case 1509:
		{ /* '1509' */
			return 0x43FE
		}
	case 151:
		{ /* '151' */
			return 0x43FC
		}
	case 1510:
		{ /* '1510' */
			return 0x43FE
		}
	case 1511:
		{ /* '1511' */
			return 0x43FE
		}
	case 1512:
		{ /* '1512' */
			return 0x43FE
		}
	case 1513:
		{ /* '1513' */
			return 0x43FE
		}
	case 1514:
		{ /* '1514' */
			return 0x43FE
		}
	case 1515:
		{ /* '1515' */
			return 0x43FE
		}
	case 1516:
		{ /* '1516' */
			return 0x43FE
		}
	case 1517:
		{ /* '1517' */
			return 0x43FE
		}
	case 1518:
		{ /* '1518' */
			return 0x43FE
		}
	case 1519:
		{ /* '1519' */
			return 0x43FE
		}
	case 152:
		{ /* '152' */
			return 0x43FC
		}
	case 1520:
		{ /* '1520' */
			return 0x43FE
		}
	case 1521:
		{ /* '1521' */
			return 0x43FE
		}
	case 1522:
		{ /* '1522' */
			return 0x43FE
		}
	case 1523:
		{ /* '1523' */
			return 0x43FE
		}
	case 1524:
		{ /* '1524' */
			return 0x43FE
		}
	case 1525:
		{ /* '1525' */
			return 0x43FE
		}
	case 1526:
		{ /* '1526' */
			return 0x43FE
		}
	case 1527:
		{ /* '1527' */
			return 0x41E5
		}
	case 1528:
		{ /* '1528' */
			return 0x41E5
		}
	case 1529:
		{ /* '1529' */
			return 0x43FF
		}
	case 153:
		{ /* '153' */
			return 0x43FC
		}
	case 1530:
		{ /* '1530' */
			return 0x41E5
		}
	case 1531:
		{ /* '1531' */
			return 0x41E5
		}
	case 1532:
		{ /* '1532' */
			return 0x43FF
		}
	case 1533:
		{ /* '1533' */
			return 0x41E5
		}
	case 1534:
		{ /* '1534' */
			return 0x41E5
		}
	case 1535:
		{ /* '1535' */
			return 0x41E5
		}
	case 1536:
		{ /* '1536' */
			return 0x43FE
		}
	case 1537:
		{ /* '1537' */
			return 0x4400
		}
	case 1538:
		{ /* '1538' */
			return 0x43FE
		}
	case 1539:
		{ /* '1539' */
			return 0x43FE
		}
	case 154:
		{ /* '154' */
			return 0x43FC
		}
	case 1540:
		{ /* '1540' */
			return 0x43FE
		}
	case 1541:
		{ /* '1541' */
			return 0x4400
		}
	case 1542:
		{ /* '1542' */
			return 0x4136
		}
	case 1543:
		{ /* '1543' */
			return 0x4266
		}
	case 1544:
		{ /* '1544' */
			return 0x437E
		}
	case 1545:
		{ /* '1545' */
			return 0x4276
		}
	case 1546:
		{ /* '1546' */
			return 0x41DE
		}
	case 1547:
		{ /* '1547' */
			return 0x41DE
		}
	case 1548:
		{ /* '1548' */
			return 0x4276
		}
	case 1549:
		{ /* '1549' */
			return 0x43A6
		}
	case 155:
		{ /* '155' */
			return 0x43FC
		}
	case 1550:
		{ /* '1550' */
			return 0x4304
		}
	case 1551:
		{ /* '1551' */
			return 0x437E
		}
	case 1552:
		{ /* '1552' */
			return 0x437E
		}
	case 1553:
		{ /* '1553' */
			return 0x4276
		}
	case 1554:
		{ /* '1554' */
			return 0x437E
		}
	case 1555:
		{ /* '1555' */
			return 0x439A
		}
	case 1556:
		{ /* '1556' */
			return 0x43A6
		}
	case 1557:
		{ /* '1557' */
			return 0x439A
		}
	case 1558:
		{ /* '1558' */
			return 0x439A
		}
	case 1559:
		{ /* '1559' */
			return 0x439A
		}
	case 156:
		{ /* '156' */
			return 0x4324
		}
	case 1560:
		{ /* '1560' */
			return 0x439A
		}
	case 1561:
		{ /* '1561' */
			return 0x439A
		}
	case 1562:
		{ /* '1562' */
			return 0x43A6
		}
	case 1563:
		{ /* '1563' */
			return 0x4136
		}
	case 1564:
		{ /* '1564' */
			return 0x4136
		}
	case 1565:
		{ /* '1565' */
			return 0x4324
		}
	case 1566:
		{ /* '1566' */
			return 0x43FC
		}
	case 1567:
		{ /* '1567' */
			return 0x43FC
		}
	case 1568:
		{ /* '1568' */
			return 0x4400
		}
	case 1569:
		{ /* '1569' */
			return 0x4400
		}
	case 157:
		{ /* '157' */
			return 0x405C
		}
	case 1570:
		{ /* '1570' */
			return 0x4314
		}
	case 1571:
		{ /* '1571' */
			return 0x41CC
		}
	case 1572:
		{ /* '1572' */
			return 0x43FC
		}
	case 1573:
		{ /* '1573' */
			return 0x43FC
		}
	case 1574:
		{ /* '1574' */
			return 0x43FC
		}
	case 1575:
		{ /* '1575' */
			return 0x43FC
		}
	case 1576:
		{ /* '1576' */
			return 0x43FC
		}
	case 1577:
		{ /* '1577' */
			return 0x431C
		}
	case 1578:
		{ /* '1578' */
			return 0x4238
		}
	case 1579:
		{ /* '1579' */
			return 0x4238
		}
	case 158:
		{ /* '158' */
			return 0x42B0
		}
	case 1580:
		{ /* '1580' */
			return 0x41CC
		}
	case 1581:
		{ /* '1581' */
			return 0x41D8
		}
	case 1582:
		{ /* '1582' */
			return 0x43FC
		}
	case 1583:
		{ /* '1583' */
			return 0x43FC
		}
	case 1584:
		{ /* '1584' */
			return 0x43FC
		}
	case 1585:
		{ /* '1585' */
			return 0x43FC
		}
	case 1586:
		{ /* '1586' */
			return 0x43FC
		}
	case 1587:
		{ /* '1587' */
			return 0x4244
		}
	case 1588:
		{ /* '1588' */
			return 0x4244
		}
	case 1589:
		{ /* '1589' */
			return 0x43FC
		}
	case 159:
		{ /* '159' */
			return 0x4328
		}
	case 1590:
		{ /* '1590' */
			return 0x4244
		}
	case 1591:
		{ /* '1591' */
			return 0x43FC
		}
	case 1592:
		{ /* '1592' */
			return 0x43FC
		}
	case 1593:
		{ /* '1593' */
			return 0x4404
		}
	case 1594:
		{ /* '1594' */
			return 0x4404
		}
	case 1595:
		{ /* '1595' */
			return 0x4404
		}
	case 1596:
		{ /* '1596' */
			return 0x43FC
		}
	case 1597:
		{ /* '1597' */
			return 0x43FC
		}
	case 1598:
		{ /* '1598' */
			return 0x43FC
		}
	case 1599:
		{ /* '1599' */
			return 0x4400
		}
	case 16:
		{ /* '16' */
			return 0x43FE
		}
	case 160:
		{ /* '160' */
			return 0x4400
		}
	case 1600:
		{ /* '1600' */
			return 0x402C
		}
	case 1601:
		{ /* '1601' */
			return 0x4054
		}
	case 1602:
		{ /* '1602' */
			return 0x402C
		}
	case 1603:
		{ /* '1603' */
			return 0x43EC
		}
	case 1604:
		{ /* '1604' */
			return 0x43EC
		}
	case 1605:
		{ /* '1605' */
			return 0x4400
		}
	case 1606:
		{ /* '1606' */
			return 0x4400
		}
	case 1607:
		{ /* '1607' */
			return 0x402C
		}
	case 1608:
		{ /* '1608' */
			return 0x43FE
		}
	case 1609:
		{ /* '1609' */
			return 0x43FC
		}
	case 161:
		{ /* '161' */
			return 0x4400
		}
	case 1610:
		{ /* '1610' */
			return 0x43FC
		}
	case 1611:
		{ /* '1611' */
			return 0x43FE
		}
	case 1612:
		{ /* '1612' */
			return 0x43FE
		}
	case 1613:
		{ /* '1613' */
			return 0x43FE
		}
	case 1614:
		{ /* '1614' */
			return 0x43FE
		}
	case 1615:
		{ /* '1615' */
			return 0x43FE
		}
	case 1616:
		{ /* '1616' */
			return 0x43FE
		}
	case 1617:
		{ /* '1617' */
			return 0x43FF
		}
	case 1618:
		{ /* '1618' */
			return 0x4400
		}
	case 1619:
		{ /* '1619' */
			return 0x4400
		}
	case 162:
		{ /* '162' */
			return 0x4400
		}
	case 1620:
		{ /* '1620' */
			return 0x4400
		}
	case 1621:
		{ /* '1621' */
			return 0x4400
		}
	case 1622:
		{ /* '1622' */
			return 0x4400
		}
	case 1623:
		{ /* '1623' */
			return 0x4400
		}
	case 1624:
		{ /* '1624' */
			return 0x43FE
		}
	case 1625:
		{ /* '1625' */
			return 0x4400
		}
	case 1626:
		{ /* '1626' */
			return 0x4400
		}
	case 1627:
		{ /* '1627' */
			return 0x4400
		}
	case 1628:
		{ /* '1628' */
			return 0x4400
		}
	case 1629:
		{ /* '1629' */
			return 0x4400
		}
	case 163:
		{ /* '163' */
			return 0x4400
		}
	case 1630:
		{ /* '1630' */
			return 0x4400
		}
	case 1631:
		{ /* '1631' */
			return 0x4400
		}
	case 1632:
		{ /* '1632' */
			return 0x4400
		}
	case 1633:
		{ /* '1633' */
			return 0x4400
		}
	case 1634:
		{ /* '1634' */
			return 0x4400
		}
	case 1635:
		{ /* '1635' */
			return 0x4400
		}
	case 1636:
		{ /* '1636' */
			return 0x4400
		}
	case 1637:
		{ /* '1637' */
			return 0x43FE
		}
	case 1638:
		{ /* '1638' */
			return 0x4400
		}
	case 1639:
		{ /* '1639' */
			return 0x4400
		}
	case 164:
		{ /* '164' */
			return 0x402C
		}
	case 1640:
		{ /* '1640' */
			return 0x4400
		}
	case 1641:
		{ /* '1641' */
			return 0x4400
		}
	case 1642:
		{ /* '1642' */
			return 0x4400
		}
	case 1643:
		{ /* '1643' */
			return 0x4400
		}
	case 1644:
		{ /* '1644' */
			return 0x4400
		}
	case 1645:
		{ /* '1645' */
			return 0x48D6
		}
	case 1646:
		{ /* '1646' */
			return 0x46A0
		}
	case 1647:
		{ /* '1647' */
			return 0x5220
		}
	case 1648:
		{ /* '1648' */
			return 0x4400
		}
	case 1649:
		{ /* '1649' */
			return 0x4400
		}
	case 165:
		{ /* '165' */
			return 0x4400
		}
	case 1650:
		{ /* '1650' */
			return 0x4400
		}
	case 1651:
		{ /* '1651' */
			return 0x4400
		}
	case 1652:
		{ /* '1652' */
			return 0x4400
		}
	case 1653:
		{ /* '1653' */
			return 0x4400
		}
	case 1654:
		{ /* '1654' */
			return 0x4400
		}
	case 1655:
		{ /* '1655' */
			return 0x43FF
		}
	case 1656:
		{ /* '1656' */
			return 0x4400
		}
	case 1657:
		{ /* '1657' */
			return 0x4400
		}
	case 1658:
		{ /* '1658' */
			return 0x4100
		}
	case 1659:
		{ /* '1659' */
			return 0x4400
		}
	case 166:
		{ /* '166' */
			return 0x4400
		}
	case 1660:
		{ /* '1660' */
			return 0x4400
		}
	case 1661:
		{ /* '1661' */
			return 0x43FF
		}
	case 1662:
		{ /* '1662' */
			return 0x43FE
		}
	case 1663:
		{ /* '1663' */
			return 0x4400
		}
	case 1664:
		{ /* '1664' */
			return 0x4400
		}
	case 1665:
		{ /* '1665' */
			return 0x43FE
		}
	case 1666:
		{ /* '1666' */
			return 0x45E0
		}
	case 1667:
		{ /* '1667' */
			return 0x4400
		}
	case 1668:
		{ /* '1668' */
			return 0x4400
		}
	case 1669:
		{ /* '1669' */
			return 0x4400
		}
	case 167:
		{ /* '167' */
			return 0x4400
		}
	case 1670:
		{ /* '1670' */
			return 0x4400
		}
	case 1671:
		{ /* '1671' */
			return 0x4400
		}
	case 1672:
		{ /* '1672' */
			return 0x4400
		}
	case 1673:
		{ /* '1673' */
			return 0x4400
		}
	case 1674:
		{ /* '1674' */
			return 0x4400
		}
	case 1675:
		{ /* '1675' */
			return 0x4400
		}
	case 1676:
		{ /* '1676' */
			return 0x4400
		}
	case 1677:
		{ /* '1677' */
			return 0x4400
		}
	case 1678:
		{ /* '1678' */
			return 0x4400
		}
	case 1679:
		{ /* '1679' */
			return 0x4400
		}
	case 168:
		{ /* '168' */
			return 0x4400
		}
	case 1680:
		{ /* '1680' */
			return 0x4400
		}
	case 1681:
		{ /* '1681' */
			return 0x4400
		}
	case 1682:
		{ /* '1682' */
			return 0x43FE
		}
	case 1683:
		{ /* '1683' */
			return 0x4400
		}
	case 1684:
		{ /* '1684' */
			return 0x4200
		}
	case 1685:
		{ /* '1685' */
			return 0x4400
		}
	case 1686:
		{ /* '1686' */
			return 0x4400
		}
	case 1687:
		{ /* '1687' */
			return 0x4400
		}
	case 1688:
		{ /* '1688' */
			return 0x4400
		}
	case 1689:
		{ /* '1689' */
			return 0x4400
		}
	case 169:
		{ /* '169' */
			return 0x4400
		}
	case 1690:
		{ /* '1690' */
			return 0x4400
		}
	case 1691:
		{ /* '1691' */
			return 0x4400
		}
	case 1692:
		{ /* '1692' */
			return 0x4400
		}
	case 1693:
		{ /* '1693' */
			return 0x4400
		}
	case 1694:
		{ /* '1694' */
			return 0x4400
		}
	case 1695:
		{ /* '1695' */
			return 0x4144
		}
	case 1696:
		{ /* '1696' */
			return 0x4100
		}
	case 1697:
		{ /* '1697' */
			return 0x4100
		}
	case 1698:
		{ /* '1698' */
			return 0x4100
		}
	case 1699:
		{ /* '1699' */
			return 0x43FC
		}
	case 17:
		{ /* '17' */
			return 0x43FE
		}
	case 170:
		{ /* '170' */
			return 0x4400
		}
	case 1700:
		{ /* '1700' */
			return 0x43C4
		}
	case 1701:
		{ /* '1701' */
			return 0x42E8
		}
	case 1702:
		{ /* '1702' */
			return 0x4760
		}
	case 1703:
		{ /* '1703' */
			return 0x40CC
		}
	case 1704:
		{ /* '1704' */
			return 0x4324
		}
	case 1705:
		{ /* '1705' */
			return 0x41A8
		}
	case 1706:
		{ /* '1706' */
			return 0x43FE
		}
	case 171:
		{ /* '171' */
			return 0x4400
		}
	case 172:
		{ /* '172' */
			return 0x4400
		}
	case 173:
		{ /* '173' */
			return 0x4400
		}
	case 174:
		{ /* '174' */
			return 0x4400
		}
	case 175:
		{ /* '175' */
			return 0x4400
		}
	case 176:
		{ /* '176' */
			return 0x4400
		}
	case 177:
		{ /* '177' */
			return 0x4400
		}
	case 178:
		{ /* '178' */
			return 0x4400
		}
	case 179:
		{ /* '179' */
			return 0x4400
		}
	case 18:
		{ /* '18' */
			return 0x43FE
		}
	case 180:
		{ /* '180' */
			return 0x4400
		}
	case 181:
		{ /* '181' */
			return 0x4400
		}
	case 182:
		{ /* '182' */
			return 0x4400
		}
	case 183:
		{ /* '183' */
			return 0x4400
		}
	case 184:
		{ /* '184' */
			return 0x4400
		}
	case 185:
		{ /* '185' */
			return 0x4400
		}
	case 186:
		{ /* '186' */
			return 0x4400
		}
	case 187:
		{ /* '187' */
			return 0x4400
		}
	case 188:
		{ /* '188' */
			return 0x4400
		}
	case 189:
		{ /* '189' */
			return 0x4400
		}
	case 19:
		{ /* '19' */
			return 0x43FE
		}
	case 190:
		{ /* '190' */
			return 0x4400
		}
	case 191:
		{ /* '191' */
			return 0x4400
		}
	case 192:
		{ /* '192' */
			return 0x4400
		}
	case 193:
		{ /* '193' */
			return 0x4400
		}
	case 194:
		{ /* '194' */
			return 0x4400
		}
	case 195:
		{ /* '195' */
			return 0x4400
		}
	case 196:
		{ /* '196' */
			return 0x4400
		}
	case 197:
		{ /* '197' */
			return 0x4400
		}
	case 198:
		{ /* '198' */
			return 0x4400
		}
	case 199:
		{ /* '199' */
			return 0x4400
		}
	case 2:
		{ /* '2' */
			return 0x4400
		}
	case 20:
		{ /* '20' */
			return 0x4400
		}
	case 200:
		{ /* '200' */
			return 0x4400
		}
	case 201:
		{ /* '201' */
			return 0x4400
		}
	case 202:
		{ /* '202' */
			return 0x4400
		}
	case 203:
		{ /* '203' */
			return 0x4400
		}
	case 204:
		{ /* '204' */
			return 0x4400
		}
	case 205:
		{ /* '205' */
			return 0x4400
		}
	case 206:
		{ /* '206' */
			return 0x4400
		}
	case 207:
		{ /* '207' */
			return 0x4400
		}
	case 208:
		{ /* '208' */
			return 0x4400
		}
	case 209:
		{ /* '209' */
			return 0x4400
		}
	case 21:
		{ /* '21' */
			return 0x4400
		}
	case 210:
		{ /* '210' */
			return 0x4400
		}
	case 211:
		{ /* '211' */
			return 0x4400
		}
	case 212:
		{ /* '212' */
			return 0x4400
		}
	case 213:
		{ /* '213' */
			return 0x4400
		}
	case 214:
		{ /* '214' */
			return 0x4400
		}
	case 215:
		{ /* '215' */
			return 0x4400
		}
	case 216:
		{ /* '216' */
			return 0x4400
		}
	case 217:
		{ /* '217' */
			return 0x4400
		}
	case 218:
		{ /* '218' */
			return 0x4400
		}
	case 219:
		{ /* '219' */
			return 0x4400
		}
	case 22:
		{ /* '22' */
			return 0x4400
		}
	case 220:
		{ /* '220' */
			return 0x4400
		}
	case 221:
		{ /* '221' */
			return 0x4400
		}
	case 222:
		{ /* '222' */
			return 0x4400
		}
	case 223:
		{ /* '223' */
			return 0x4400
		}
	case 224:
		{ /* '224' */
			return 0x4400
		}
	case 225:
		{ /* '225' */
			return 0x4400
		}
	case 226:
		{ /* '226' */
			return 0x4400
		}
	case 227:
		{ /* '227' */
			return 0x4400
		}
	case 228:
		{ /* '228' */
			return 0x4400
		}
	case 229:
		{ /* '229' */
			return 0x4400
		}
	case 23:
		{ /* '23' */
			return 0x4400
		}
	case 230:
		{ /* '230' */
			return 0x4400
		}
	case 231:
		{ /* '231' */
			return 0x4400
		}
	case 232:
		{ /* '232' */
			return 0x4400
		}
	case 233:
		{ /* '233' */
			return 0x4400
		}
	case 234:
		{ /* '234' */
			return 0x4400
		}
	case 235:
		{ /* '235' */
			return 0x4400
		}
	case 236:
		{ /* '236' */
			return 0x4400
		}
	case 237:
		{ /* '237' */
			return 0x4400
		}
	case 238:
		{ /* '238' */
			return 0x4400
		}
	case 239:
		{ /* '239' */
			return 0x4400
		}
	case 24:
		{ /* '24' */
			return 0x4400
		}
	case 240:
		{ /* '240' */
			return 0x4400
		}
	case 241:
		{ /* '241' */
			return 0x4400
		}
	case 242:
		{ /* '242' */
			return 0x4400
		}
	case 243:
		{ /* '243' */
			return 0x4400
		}
	case 244:
		{ /* '244' */
			return 0x4400
		}
	case 245:
		{ /* '245' */
			return 0x4400
		}
	case 246:
		{ /* '246' */
			return 0x4400
		}
	case 247:
		{ /* '247' */
			return 0x4400
		}
	case 248:
		{ /* '248' */
			return 0x4400
		}
	case 249:
		{ /* '249' */
			return 0x4400
		}
	case 25:
		{ /* '25' */
			return 0x4400
		}
	case 250:
		{ /* '250' */
			return 0x4400
		}
	case 251:
		{ /* '251' */
			return 0x4400
		}
	case 252:
		{ /* '252' */
			return 0x4400
		}
	case 253:
		{ /* '253' */
			return 0x4400
		}
	case 254:
		{ /* '254' */
			return 0x4400
		}
	case 255:
		{ /* '255' */
			return 0x4400
		}
	case 256:
		{ /* '256' */
			return 0x4400
		}
	case 257:
		{ /* '257' */
			return 0x4400
		}
	case 258:
		{ /* '258' */
			return 0x4400
		}
	case 259:
		{ /* '259' */
			return 0x4400
		}
	case 26:
		{ /* '26' */
			return 0x4400
		}
	case 260:
		{ /* '260' */
			return 0x4400
		}
	case 261:
		{ /* '261' */
			return 0x4400
		}
	case 262:
		{ /* '262' */
			return 0x4400
		}
	case 263:
		{ /* '263' */
			return 0x4400
		}
	case 264:
		{ /* '264' */
			return 0x4400
		}
	case 265:
		{ /* '265' */
			return 0x4400
		}
	case 266:
		{ /* '266' */
			return 0x4400
		}
	case 267:
		{ /* '267' */
			return 0x4400
		}
	case 268:
		{ /* '268' */
			return 0x4400
		}
	case 269:
		{ /* '269' */
			return 0x4400
		}
	case 27:
		{ /* '27' */
			return 0x43FE
		}
	case 270:
		{ /* '270' */
			return 0x4080
		}
	case 271:
		{ /* '271' */
			return 0x4080
		}
	case 272:
		{ /* '272' */
			return 0x4080
		}
	case 273:
		{ /* '273' */
			return 0x4400
		}
	case 274:
		{ /* '274' */
			return 0x4400
		}
	case 275:
		{ /* '275' */
			return 0x4400
		}
	case 276:
		{ /* '276' */
			return 0x4400
		}
	case 277:
		{ /* '277' */
			return 0x4400
		}
	case 278:
		{ /* '278' */
			return 0x4400
		}
	case 279:
		{ /* '279' */
			return 0x4194
		}
	case 28:
		{ /* '28' */
			return 0x4400
		}
	case 280:
		{ /* '280' */
			return 0x4400
		}
	case 281:
		{ /* '281' */
			return 0x4400
		}
	case 282:
		{ /* '282' */
			return 0x4400
		}
	case 283:
		{ /* '283' */
			return 0x4400
		}
	case 284:
		{ /* '284' */
			return 0x4400
		}
	case 285:
		{ /* '285' */
			return 0x4400
		}
	case 286:
		{ /* '286' */
			return 0x4400
		}
	case 287:
		{ /* '287' */
			return 0x4400
		}
	case 288:
		{ /* '288' */
			return 0x4400
		}
	case 289:
		{ /* '289' */
			return 0x4400
		}
	case 29:
		{ /* '29' */
			return 0x43FE
		}
	case 290:
		{ /* '290' */
			return 0x4400
		}
	case 291:
		{ /* '291' */
			return 0x4400
		}
	case 292:
		{ /* '292' */
			return 0x4400
		}
	case 293:
		{ /* '293' */
			return 0x4400
		}
	case 294:
		{ /* '294' */
			return 0x4400
		}
	case 295:
		{ /* '295' */
			return 0x4400
		}
	case 296:
		{ /* '296' */
			return 0x4400
		}
	case 297:
		{ /* '297' */
			return 0x4400
		}
	case 298:
		{ /* '298' */
			return 0x4400
		}
	case 299:
		{ /* '299' */
			return 0x4400
		}
	case 3:
		{ /* '3' */
			return 0x4400
		}
	case 30:
		{ /* '30' */
			return 0x43FE
		}
	case 300:
		{ /* '300' */
			return 0x4400
		}
	case 301:
		{ /* '301' */
			return 0x4400
		}
	case 302:
		{ /* '302' */
			return 0x4400
		}
	case 303:
		{ /* '303' */
			return 0x4400
		}
	case 304:
		{ /* '304' */
			return 0x4400
		}
	case 305:
		{ /* '305' */
			return 0x4400
		}
	case 306:
		{ /* '306' */
			return 0x4400
		}
	case 307:
		{ /* '307' */
			return 0x4400
		}
	case 308:
		{ /* '308' */
			return 0x4400
		}
	case 309:
		{ /* '309' */
			return 0x4400
		}
	case 31:
		{ /* '31' */
			return 0x43FE
		}
	case 310:
		{ /* '310' */
			return 0x4400
		}
	case 311:
		{ /* '311' */
			return 0x4400
		}
	case 312:
		{ /* '312' */
			return 0x4400
		}
	case 313:
		{ /* '313' */
			return 0x4400
		}
	case 314:
		{ /* '314' */
			return 0x4400
		}
	case 315:
		{ /* '315' */
			return 0x4400
		}
	case 316:
		{ /* '316' */
			return 0x4400
		}
	case 317:
		{ /* '317' */
			return 0x4400
		}
	case 318:
		{ /* '318' */
			return 0x4400
		}
	case 319:
		{ /* '319' */
			return 0x4400
		}
	case 32:
		{ /* '32' */
			return 0x4324
		}
	case 320:
		{ /* '320' */
			return 0x4400
		}
	case 321:
		{ /* '321' */
			return 0x4400
		}
	case 322:
		{ /* '322' */
			return 0x4400
		}
	case 323:
		{ /* '323' */
			return 0x4400
		}
	case 324:
		{ /* '324' */
			return 0x4400
		}
	case 325:
		{ /* '325' */
			return 0x4400
		}
	case 326:
		{ /* '326' */
			return 0x4400
		}
	case 327:
		{ /* '327' */
			return 0x4400
		}
	case 328:
		{ /* '328' */
			return 0x4400
		}
	case 329:
		{ /* '329' */
			return 0x43C4
		}
	case 33:
		{ /* '33' */
			return 0x4324
		}
	case 330:
		{ /* '330' */
			return 0x42E8
		}
	case 331:
		{ /* '331' */
			return 0x4760
		}
	case 332:
		{ /* '332' */
			return 0x4100
		}
	case 333:
		{ /* '333' */
			return 0x4100
		}
	case 334:
		{ /* '334' */
			return 0x4100
		}
	case 335:
		{ /* '335' */
			return 0x4100
		}
	case 336:
		{ /* '336' */
			return 0x4100
		}
	case 337:
		{ /* '337' */
			return 0x4100
		}
	case 338:
		{ /* '338' */
			return 0x43EC
		}
	case 339:
		{ /* '339' */
			return 0x40CC
		}
	case 34:
		{ /* '34' */
			return 0x4400
		}
	case 340:
		{ /* '340' */
			return 0x41A8
		}
	case 341:
		{ /* '341' */
			return 0x4144
		}
	case 342:
		{ /* '342' */
			return 0x42D4
		}
	case 343:
		{ /* '343' */
			return 0x43FC
		}
	case 344:
		{ /* '344' */
			return 0x4144
		}
	case 345:
		{ /* '345' */
			return 0x43FC
		}
	case 346:
		{ /* '346' */
			return 0x4202
		}
	case 347:
		{ /* '347' */
			return 0x40F4
		}
	case 348:
		{ /* '348' */
			return 0x4064
		}
	case 349:
		{ /* '349' */
			return 0x4124
		}
	case 35:
		{ /* '35' */
			return 0x4400
		}
	case 350:
		{ /* '350' */
			return 0x40CC
		}
	case 351:
		{ /* '351' */
			return 0x43EC
		}
	case 352:
		{ /* '352' */
			return 0x43EC
		}
	case 353:
		{ /* '353' */
			return 0x43EC
		}
	case 354:
		{ /* '354' */
			return 0x43EC
		}
	case 355:
		{ /* '355' */
			return 0x43EC
		}
	case 356:
		{ /* '356' */
			return 0x43EC
		}
	case 357:
		{ /* '357' */
			return 0x43EC
		}
	case 358:
		{ /* '358' */
			return 0x43EC
		}
	case 359:
		{ /* '359' */
			return 0x43EC
		}
	case 36:
		{ /* '36' */
			return 0x4204
		}
	case 360:
		{ /* '360' */
			return 0x43EC
		}
	case 361:
		{ /* '361' */
			return 0x43EC
		}
	case 362:
		{ /* '362' */
			return 0x43EC
		}
	case 363:
		{ /* '363' */
			return 0x43EC
		}
	case 364:
		{ /* '364' */
			return 0x43EC
		}
	case 365:
		{ /* '365' */
			return 0x43EC
		}
	case 366:
		{ /* '366' */
			return 0x43EC
		}
	case 367:
		{ /* '367' */
			return 0x43EC
		}
	case 368:
		{ /* '368' */
			return 0x43EC
		}
	case 369:
		{ /* '369' */
			return 0x43EC
		}
	case 37:
		{ /* '37' */
			return 0x4204
		}
	case 370:
		{ /* '370' */
			return 0x43EC
		}
	case 371:
		{ /* '371' */
			return 0x43EC
		}
	case 372:
		{ /* '372' */
			return 0x43EC
		}
	case 373:
		{ /* '373' */
			return 0x43EC
		}
	case 374:
		{ /* '374' */
			return 0x43EC
		}
	case 375:
		{ /* '375' */
			return 0x43EC
		}
	case 376:
		{ /* '376' */
			return 0x0000
		}
	case 377:
		{ /* '377' */
			return 0x0000
		}
	case 378:
		{ /* '378' */
			return 0x4400
		}
	case 379:
		{ /* '379' */
			return 0x4400
		}
	case 38:
		{ /* '38' */
			return 0x4400
		}
	case 380:
		{ /* '380' */
			return 0x4400
		}
	case 381:
		{ /* '381' */
			return 0x4400
		}
	case 382:
		{ /* '382' */
			return 0x4400
		}
	case 383:
		{ /* '383' */
			return 0x4400
		}
	case 384:
		{ /* '384' */
			return 0x4400
		}
	case 385:
		{ /* '385' */
			return 0x4400
		}
	case 386:
		{ /* '386' */
			return 0x4400
		}
	case 387:
		{ /* '387' */
			return 0x4400
		}
	case 388:
		{ /* '388' */
			return 0x4400
		}
	case 389:
		{ /* '389' */
			return 0x4400
		}
	case 39:
		{ /* '39' */
			return 0x43CE
		}
	case 390:
		{ /* '390' */
			return 0x4400
		}
	case 391:
		{ /* '391' */
			return 0x4400
		}
	case 392:
		{ /* '392' */
			return 0x4400
		}
	case 393:
		{ /* '393' */
			return 0x4400
		}
	case 394:
		{ /* '394' */
			return 0x4400
		}
	case 395:
		{ /* '395' */
			return 0x4400
		}
	case 396:
		{ /* '396' */
			return 0x4326
		}
	case 397:
		{ /* '397' */
			return 0x4326
		}
	case 398:
		{ /* '398' */
			return 0x4400
		}
	case 399:
		{ /* '399' */
			return 0x4400
		}
	case 4:
		{ /* '4' */
			return 0x43EC
		}
	case 40:
		{ /* '40' */
			return 0x43FF
		}
	case 400:
		{ /* '400' */
			return 0x4400
		}
	case 401:
		{ /* '401' */
			return 0x4400
		}
	case 402:
		{ /* '402' */
			return 0x4400
		}
	case 403:
		{ /* '403' */
			return 0x4400
		}
	case 404:
		{ /* '404' */
			return 0x4400
		}
	case 405:
		{ /* '405' */
			return 0x4400
		}
	case 406:
		{ /* '406' */
			return 0x4400
		}
	case 407:
		{ /* '407' */
			return 0x4400
		}
	case 408:
		{ /* '408' */
			return 0x4400
		}
	case 409:
		{ /* '409' */
			return 0x4400
		}
	case 41:
		{ /* '41' */
			return 0x4400
		}
	case 410:
		{ /* '410' */
			return 0x4400
		}
	case 411:
		{ /* '411' */
			return 0x4400
		}
	case 412:
		{ /* '412' */
			return 0x4400
		}
	case 413:
		{ /* '413' */
			return 0x4400
		}
	case 414:
		{ /* '414' */
			return 0x4400
		}
	case 415:
		{ /* '415' */
			return 0x4400
		}
	case 416:
		{ /* '416' */
			return 0x4400
		}
	case 417:
		{ /* '417' */
			return 0x4400
		}
	case 418:
		{ /* '418' */
			return 0x4400
		}
	case 419:
		{ /* '419' */
			return 0x4400
		}
	case 42:
		{ /* '42' */
			return 0x4400
		}
	case 420:
		{ /* '420' */
			return 0x4400
		}
	case 421:
		{ /* '421' */
			return 0x4400
		}
	case 422:
		{ /* '422' */
			return 0x4400
		}
	case 423:
		{ /* '423' */
			return 0x4400
		}
	case 424:
		{ /* '424' */
			return 0x4400
		}
	case 425:
		{ /* '425' */
			return 0x4400
		}
	case 426:
		{ /* '426' */
			return 0x4400
		}
	case 427:
		{ /* '427' */
			return 0x4400
		}
	case 428:
		{ /* '428' */
			return 0x4400
		}
	case 429:
		{ /* '429' */
			return 0x4400
		}
	case 43:
		{ /* '43' */
			return 0x4400
		}
	case 430:
		{ /* '430' */
			return 0x4400
		}
	case 431:
		{ /* '431' */
			return 0x4400
		}
	case 432:
		{ /* '432' */
			return 0x4400
		}
	case 433:
		{ /* '433' */
			return 0x4400
		}
	case 434:
		{ /* '434' */
			return 0x4400
		}
	case 435:
		{ /* '435' */
			return 0x4400
		}
	case 436:
		{ /* '436' */
			return 0x4400
		}
	case 437:
		{ /* '437' */
			return 0x4400
		}
	case 438:
		{ /* '438' */
			return 0x4400
		}
	case 439:
		{ /* '439' */
			return 0x4400
		}
	case 44:
		{ /* '44' */
			return 0x4400
		}
	case 440:
		{ /* '440' */
			return 0x4400
		}
	case 441:
		{ /* '441' */
			return 0x4400
		}
	case 442:
		{ /* '442' */
			return 0x4324
		}
	case 443:
		{ /* '443' */
			return 0x4400
		}
	case 444:
		{ /* '444' */
			return 0x4400
		}
	case 445:
		{ /* '445' */
			return 0x4400
		}
	case 446:
		{ /* '446' */
			return 0x4400
		}
	case 447:
		{ /* '447' */
			return 0x4400
		}
	case 448:
		{ /* '448' */
			return 0x4400
		}
	case 449:
		{ /* '449' */
			return 0x45E0
		}
	case 45:
		{ /* '45' */
			return 0x4400
		}
	case 450:
		{ /* '450' */
			return 0x43FE
		}
	case 451:
		{ /* '451' */
			return 0x43FF
		}
	case 452:
		{ /* '452' */
			return 0x4400
		}
	case 453:
		{ /* '453' */
			return 0x4400
		}
	case 454:
		{ /* '454' */
			return 0x4400
		}
	case 455:
		{ /* '455' */
			return 0x4400
		}
	case 456:
		{ /* '456' */
			return 0x4400
		}
	case 457:
		{ /* '457' */
			return 0x4400
		}
	case 458:
		{ /* '458' */
			return 0x4400
		}
	case 459:
		{ /* '459' */
			return 0x4400
		}
	case 46:
		{ /* '46' */
			return 0x4000
		}
	case 460:
		{ /* '460' */
			return 0x4400
		}
	case 461:
		{ /* '461' */
			return 0x4400
		}
	case 462:
		{ /* '462' */
			return 0x4400
		}
	case 463:
		{ /* '463' */
			return 0x4400
		}
	case 464:
		{ /* '464' */
			return 0x4204
		}
	case 465:
		{ /* '465' */
			return 0x4204
		}
	case 466:
		{ /* '466' */
			return 0x4144
		}
	case 467:
		{ /* '467' */
			return 0x43FE
		}
	case 468:
		{ /* '468' */
			return 0x48D6
		}
	case 469:
		{ /* '469' */
			return 0x4CAC
		}
	case 47:
		{ /* '47' */
			return 0x43FE
		}
	case 470:
		{ /* '470' */
			return 0x5220
		}
	case 471:
		{ /* '471' */
			return 0x4CAC
		}
	case 472:
		{ /* '472' */
			return 0x4100
		}
	case 473:
		{ /* '473' */
			return 0x4198
		}
	case 474:
		{ /* '474' */
			return 0x4200
		}
	case 475:
		{ /* '475' */
			return 0x4200
		}
	case 476:
		{ /* '476' */
			return 0x43FE
		}
	case 477:
		{ /* '477' */
			return 0x43FE
		}
	case 478:
		{ /* '478' */
			return 0x43FE
		}
	case 479:
		{ /* '479' */
			return 0x43FE
		}
	case 48:
		{ /* '48' */
			return 0x43FE
		}
	case 480:
		{ /* '480' */
			return 0x43FE
		}
	case 481:
		{ /* '481' */
			return 0x43FE
		}
	case 482:
		{ /* '482' */
			return 0x4400
		}
	case 483:
		{ /* '483' */
			return 0x4400
		}
	case 484:
		{ /* '484' */
			return 0x43FE
		}
	case 485:
		{ /* '485' */
			return 0x43FE
		}
	case 486:
		{ /* '486' */
			return 0x43FE
		}
	case 487:
		{ /* '487' */
			return 0x4400
		}
	case 488:
		{ /* '488' */
			return 0x4400
		}
	case 489:
		{ /* '489' */
			return 0x4198
		}
	case 49:
		{ /* '49' */
			return 0x4400
		}
	case 490:
		{ /* '490' */
			return 0x4200
		}
	case 491:
		{ /* '491' */
			return 0x4200
		}
	case 492:
		{ /* '492' */
			return 0x4196
		}
	case 493:
		{ /* '493' */
			return 0x4400
		}
	case 494:
		{ /* '494' */
			return 0x4400
		}
	case 495:
		{ /* '495' */
			return 0x43FF
		}
	case 496:
		{ /* '496' */
			return 0x43FE
		}
	case 497:
		{ /* '497' */
			return 0x4400
		}
	case 498:
		{ /* '498' */
			return 0x4400
		}
	case 499:
		{ /* '499' */
			return 0x4400
		}
	case 5:
		{ /* '5' */
			return 0x402C
		}
	case 50:
		{ /* '50' */
			return 0x40F4
		}
	case 500:
		{ /* '500' */
			return 0x4400
		}
	case 501:
		{ /* '501' */
			return 0x4400
		}
	case 502:
		{ /* '502' */
			return 0x4400
		}
	case 503:
		{ /* '503' */
			return 0x4402
		}
	case 504:
		{ /* '504' */
			return 0x4EE0
		}
	case 505:
		{ /* '505' */
			return 0x4EE0
		}
	case 506:
		{ /* '506' */
			return 0x43FF
		}
	case 507:
		{ /* '507' */
			return 0x4400
		}
	case 508:
		{ /* '508' */
			return 0x4400
		}
	case 509:
		{ /* '509' */
			return 0x4400
		}
	case 51:
		{ /* '51' */
			return 0x40F4
		}
	case 510:
		{ /* '510' */
			return 0x4400
		}
	case 511:
		{ /* '511' */
			return 0x4400
		}
	case 512:
		{ /* '512' */
			return 0x4400
		}
	case 513:
		{ /* '513' */
			return 0x4400
		}
	case 514:
		{ /* '514' */
			return 0x4400
		}
	case 515:
		{ /* '515' */
			return 0x4400
		}
	case 516:
		{ /* '516' */
			return 0x4400
		}
	case 517:
		{ /* '517' */
			return 0x4400
		}
	case 518:
		{ /* '518' */
			return 0x4400
		}
	case 519:
		{ /* '519' */
			return 0x4400
		}
	case 52:
		{ /* '52' */
			return 0x402C
		}
	case 520:
		{ /* '520' */
			return 0x4400
		}
	case 521:
		{ /* '521' */
			return 0x4400
		}
	case 522:
		{ /* '522' */
			return 0x4400
		}
	case 523:
		{ /* '523' */
			return 0x4400
		}
	case 524:
		{ /* '524' */
			return 0x4400
		}
	case 525:
		{ /* '525' */
			return 0x4400
		}
	case 526:
		{ /* '526' */
			return 0x4400
		}
	case 527:
		{ /* '527' */
			return 0x4400
		}
	case 528:
		{ /* '528' */
			return 0x43FE
		}
	case 529:
		{ /* '529' */
			return 0x43FE
		}
	case 53:
		{ /* '53' */
			return 0x43F0
		}
	case 530:
		{ /* '530' */
			return 0x43FE
		}
	case 531:
		{ /* '531' */
			return 0x43FE
		}
	case 532:
		{ /* '532' */
			return 0x43FE
		}
	case 533:
		{ /* '533' */
			return 0x43F2
		}
	case 534:
		{ /* '534' */
			return 0x4400
		}
	case 535:
		{ /* '535' */
			return 0x4400
		}
	case 536:
		{ /* '536' */
			return 0x4400
		}
	case 537:
		{ /* '537' */
			return 0x4400
		}
	case 538:
		{ /* '538' */
			return 0x4400
		}
	case 539:
		{ /* '539' */
			return 0x4400
		}
	case 54:
		{ /* '54' */
			return 0x43EC
		}
	case 540:
		{ /* '540' */
			return 0x4400
		}
	case 541:
		{ /* '541' */
			return 0x4400
		}
	case 542:
		{ /* '542' */
			return 0x4400
		}
	case 543:
		{ /* '543' */
			return 0x4400
		}
	case 544:
		{ /* '544' */
			return 0x4400
		}
	case 545:
		{ /* '545' */
			return 0x4400
		}
	case 546:
		{ /* '546' */
			return 0x46A0
		}
	case 547:
		{ /* '547' */
			return 0x4400
		}
	case 548:
		{ /* '548' */
			return 0x4400
		}
	case 549:
		{ /* '549' */
			return 0x4400
		}
	case 55:
		{ /* '55' */
			return 0x43EC
		}
	case 550:
		{ /* '550' */
			return 0x4400
		}
	case 551:
		{ /* '551' */
			return 0x4400
		}
	case 552:
		{ /* '552' */
			return 0x4326
		}
	case 553:
		{ /* '553' */
			return 0x4400
		}
	case 554:
		{ /* '554' */
			return 0x4194
		}
	case 555:
		{ /* '555' */
			return 0x43FE
		}
	case 556:
		{ /* '556' */
			return 0x43FE
		}
	case 557:
		{ /* '557' */
			return 0x43FE
		}
	case 558:
		{ /* '558' */
			return 0x43FE
		}
	case 559:
		{ /* '559' */
			return 0x43FC
		}
	case 56:
		{ /* '56' */
			return 0x43FC
		}
	case 560:
		{ /* '560' */
			return 0x43FC
		}
	case 561:
		{ /* '561' */
			return 0x4338
		}
	case 562:
		{ /* '562' */
			return 0x43FC
		}
	case 563:
		{ /* '563' */
			return 0x42E0
		}
	case 564:
		{ /* '564' */
			return 0x42E0
		}
	case 565:
		{ /* '565' */
			return 0x43FC
		}
	case 566:
		{ /* '566' */
			return 0x42E0
		}
	case 567:
		{ /* '567' */
			return 0x43FC
		}
	case 568:
		{ /* '568' */
			return 0x42A0
		}
	case 569:
		{ /* '569' */
			return 0x42A0
		}
	case 57:
		{ /* '57' */
			return 0x4400
		}
	case 570:
		{ /* '570' */
			return 0x42A0
		}
	case 571:
		{ /* '571' */
			return 0x42A0
		}
	case 572:
		{ /* '572' */
			return 0x42E0
		}
	case 573:
		{ /* '573' */
			return 0x43FC
		}
	case 574:
		{ /* '574' */
			return 0x42E0
		}
	case 575:
		{ /* '575' */
			return 0x43FC
		}
	case 576:
		{ /* '576' */
			return 0x43FC
		}
	case 577:
		{ /* '577' */
			return 0x43FC
		}
	case 578:
		{ /* '578' */
			return 0x42E0
		}
	case 579:
		{ /* '579' */
			return 0x43FC
		}
	case 58:
		{ /* '58' */
			return 0x4400
		}
	case 580:
		{ /* '580' */
			return 0x42A0
		}
	case 581:
		{ /* '581' */
			return 0x42A0
		}
	case 582:
		{ /* '582' */
			return 0x42A0
		}
	case 583:
		{ /* '583' */
			return 0x43FC
		}
	case 584:
		{ /* '584' */
			return 0x43FC
		}
	case 585:
		{ /* '585' */
			return 0x43F4
		}
	case 586:
		{ /* '586' */
			return 0x43FC
		}
	case 587:
		{ /* '587' */
			return 0x43FC
		}
	case 588:
		{ /* '588' */
			return 0x4370
		}
	case 589:
		{ /* '589' */
			return 0x4370
		}
	case 59:
		{ /* '59' */
			return 0x4400
		}
	case 590:
		{ /* '590' */
			return 0x4370
		}
	case 591:
		{ /* '591' */
			return 0x439C
		}
	case 592:
		{ /* '592' */
			return 0x43FC
		}
	case 593:
		{ /* '593' */
			return 0x43FC
		}
	case 594:
		{ /* '594' */
			return 0x4324
		}
	case 595:
		{ /* '595' */
			return 0x4374
		}
	case 596:
		{ /* '596' */
			return 0x4374
		}
	case 597:
		{ /* '597' */
			return 0x4274
		}
	case 598:
		{ /* '598' */
			return 0x43FC
		}
	case 599:
		{ /* '599' */
			return 0x43FC
		}
	case 6:
		{ /* '6' */
			return 0x43FE
		}
	case 60:
		{ /* '60' */
			return 0x4400
		}
	case 600:
		{ /* '600' */
			return 0x43FC
		}
	case 601:
		{ /* '601' */
			return 0x43FC
		}
	case 602:
		{ /* '602' */
			return 0x43FC
		}
	case 603:
		{ /* '603' */
			return 0x43FC
		}
	case 604:
		{ /* '604' */
			return 0x42A0
		}
	case 605:
		{ /* '605' */
			return 0x43FC
		}
	case 606:
		{ /* '606' */
			return 0x42A0
		}
	case 607:
		{ /* '607' */
			return 0x43FC
		}
	case 608:
		{ /* '608' */
			return 0x43FC
		}
	case 609:
		{ /* '609' */
			return 0x43FC
		}
	case 61:
		{ /* '61' */
			return 0x4400
		}
	case 610:
		{ /* '610' */
			return 0x43FC
		}
	case 611:
		{ /* '611' */
			return 0x43FC
		}
	case 612:
		{ /* '612' */
			return 0x43FC
		}
	case 613:
		{ /* '613' */
			return 0x43FC
		}
	case 614:
		{ /* '614' */
			return 0x43FC
		}
	case 615:
		{ /* '615' */
			return 0x43FC
		}
	case 616:
		{ /* '616' */
			return 0x43FC
		}
	case 617:
		{ /* '617' */
			return 0x43FC
		}
	case 618:
		{ /* '618' */
			return 0x43FC
		}
	case 619:
		{ /* '619' */
			return 0x43FC
		}
	case 62:
		{ /* '62' */
			return 0x4400
		}
	case 620:
		{ /* '620' */
			return 0x43FC
		}
	case 621:
		{ /* '621' */
			return 0x43FC
		}
	case 622:
		{ /* '622' */
			return 0x43FC
		}
	case 623:
		{ /* '623' */
			return 0x43FC
		}
	case 624:
		{ /* '624' */
			return 0x43FC
		}
	case 625:
		{ /* '625' */
			return 0x43FC
		}
	case 626:
		{ /* '626' */
			return 0x43FC
		}
	case 627:
		{ /* '627' */
			return 0x43FC
		}
	case 628:
		{ /* '628' */
			return 0x43FC
		}
	case 629:
		{ /* '629' */
			return 0x43FC
		}
	case 63:
		{ /* '63' */
			return 0x4400
		}
	case 630:
		{ /* '630' */
			return 0x43FC
		}
	case 631:
		{ /* '631' */
			return 0x43FC
		}
	case 632:
		{ /* '632' */
			return 0x43FC
		}
	case 633:
		{ /* '633' */
			return 0x43FC
		}
	case 634:
		{ /* '634' */
			return 0x43FC
		}
	case 635:
		{ /* '635' */
			return 0x43FC
		}
	case 636:
		{ /* '636' */
			return 0x43FC
		}
	case 637:
		{ /* '637' */
			return 0x43FC
		}
	case 638:
		{ /* '638' */
			return 0x43FC
		}
	case 639:
		{ /* '639' */
			return 0x43FC
		}
	case 64:
		{ /* '64' */
			return 0x4400
		}
	case 640:
		{ /* '640' */
			return 0x43FE
		}
	case 641:
		{ /* '641' */
			return 0x43FE
		}
	case 642:
		{ /* '642' */
			return 0x43FE
		}
	case 643:
		{ /* '643' */
			return 0x43FE
		}
	case 644:
		{ /* '644' */
			return 0x43FE
		}
	case 645:
		{ /* '645' */
			return 0x4324
		}
	case 646:
		{ /* '646' */
			return 0x43FE
		}
	case 647:
		{ /* '647' */
			return 0x43FE
		}
	case 648:
		{ /* '648' */
			return 0x43FC
		}
	case 649:
		{ /* '649' */
			return 0x43FC
		}
	case 65:
		{ /* '65' */
			return 0x4400
		}
	case 650:
		{ /* '650' */
			return 0x43FC
		}
	case 651:
		{ /* '651' */
			return 0x43FC
		}
	case 652:
		{ /* '652' */
			return 0x43FC
		}
	case 653:
		{ /* '653' */
			return 0x4388
		}
	case 654:
		{ /* '654' */
			return 0x43FC
		}
	case 655:
		{ /* '655' */
			return 0x4400
		}
	case 656:
		{ /* '656' */
			return 0x43FC
		}
	case 657:
		{ /* '657' */
			return 0x43FC
		}
	case 658:
		{ /* '658' */
			return 0x43FC
		}
	case 659:
		{ /* '659' */
			return 0x43FC
		}
	case 66:
		{ /* '66' */
			return 0x4400
		}
	case 660:
		{ /* '660' */
			return 0x43FC
		}
	case 661:
		{ /* '661' */
			return 0x43FC
		}
	case 662:
		{ /* '662' */
			return 0x43FC
		}
	case 663:
		{ /* '663' */
			return 0x43FC
		}
	case 664:
		{ /* '664' */
			return 0x43FC
		}
	case 665:
		{ /* '665' */
			return 0x43FC
		}
	case 666:
		{ /* '666' */
			return 0x43FC
		}
	case 667:
		{ /* '667' */
			return 0x43FC
		}
	case 668:
		{ /* '668' */
			return 0x43FC
		}
	case 669:
		{ /* '669' */
			return 0x43FC
		}
	case 67:
		{ /* '67' */
			return 0x4400
		}
	case 670:
		{ /* '670' */
			return 0x43FC
		}
	case 671:
		{ /* '671' */
			return 0x43FC
		}
	case 672:
		{ /* '672' */
			return 0x43FC
		}
	case 673:
		{ /* '673' */
			return 0x43FC
		}
	case 674:
		{ /* '674' */
			return 0x43FC
		}
	case 675:
		{ /* '675' */
			return 0x43FC
		}
	case 676:
		{ /* '676' */
			return 0x4400
		}
	case 677:
		{ /* '677' */
			return 0x43FE
		}
	case 678:
		{ /* '678' */
			return 0x43FE
		}
	case 679:
		{ /* '679' */
			return 0x43FE
		}
	case 68:
		{ /* '68' */
			return 0x4400
		}
	case 680:
		{ /* '680' */
			return 0x43FE
		}
	case 681:
		{ /* '681' */
			return 0x43FE
		}
	case 682:
		{ /* '682' */
			return 0x43FE
		}
	case 683:
		{ /* '683' */
			return 0x43FE
		}
	case 684:
		{ /* '684' */
			return 0x43FE
		}
	case 685:
		{ /* '685' */
			return 0x43FE
		}
	case 686:
		{ /* '686' */
			return 0x43FE
		}
	case 687:
		{ /* '687' */
			return 0x43FE
		}
	case 688:
		{ /* '688' */
			return 0x43FE
		}
	case 689:
		{ /* '689' */
			return 0x43FE
		}
	case 69:
		{ /* '69' */
			return 0x4400
		}
	case 690:
		{ /* '690' */
			return 0x43FE
		}
	case 691:
		{ /* '691' */
			return 0x43FE
		}
	case 692:
		{ /* '692' */
			return 0x43FE
		}
	case 693:
		{ /* '693' */
			return 0x43FE
		}
	case 694:
		{ /* '694' */
			return 0x43FE
		}
	case 695:
		{ /* '695' */
			return 0x43FE
		}
	case 696:
		{ /* '696' */
			return 0x43FE
		}
	case 697:
		{ /* '697' */
			return 0x43FE
		}
	case 698:
		{ /* '698' */
			return 0x43FE
		}
	case 699:
		{ /* '699' */
			return 0x416C
		}
	case 7:
		{ /* '7' */
			return 0x43FE
		}
	case 70:
		{ /* '70' */
			return 0x4400
		}
	case 700:
		{ /* '700' */
			return 0x416C
		}
	case 701:
		{ /* '701' */
			return 0x416C
		}
	case 702:
		{ /* '702' */
			return 0x43FF
		}
	case 703:
		{ /* '703' */
			return 0x4324
		}
	case 704:
		{ /* '704' */
			return 0x43FC
		}
	case 705:
		{ /* '705' */
			return 0x43FC
		}
	case 706:
		{ /* '706' */
			return 0x43FC
		}
	case 707:
		{ /* '707' */
			return 0x43FC
		}
	case 708:
		{ /* '708' */
			return 0x4354
		}
	case 709:
		{ /* '709' */
			return 0x4334
		}
	case 71:
		{ /* '71' */
			return 0x4400
		}
	case 710:
		{ /* '710' */
			return 0x4354
		}
	case 711:
		{ /* '711' */
			return 0x43FC
		}
	case 712:
		{ /* '712' */
			return 0x43FC
		}
	case 713:
		{ /* '713' */
			return 0x43FC
		}
	case 714:
		{ /* '714' */
			return 0x43FC
		}
	case 715:
		{ /* '715' */
			return 0x43FC
		}
	case 716:
		{ /* '716' */
			return 0x43FC
		}
	case 717:
		{ /* '717' */
			return 0x43FC
		}
	case 718:
		{ /* '718' */
			return 0x43FC
		}
	case 719:
		{ /* '719' */
			return 0x43FC
		}
	case 72:
		{ /* '72' */
			return 0x4052
		}
	case 720:
		{ /* '720' */
			return 0x43FC
		}
	case 721:
		{ /* '721' */
			return 0x43FC
		}
	case 722:
		{ /* '722' */
			return 0x43FC
		}
	case 723:
		{ /* '723' */
			return 0x43FC
		}
	case 724:
		{ /* '724' */
			return 0x43FC
		}
	case 725:
		{ /* '725' */
			return 0x43FC
		}
	case 726:
		{ /* '726' */
			return 0x43FC
		}
	case 727:
		{ /* '727' */
			return 0x43FC
		}
	case 728:
		{ /* '728' */
			return 0x43FC
		}
	case 729:
		{ /* '729' */
			return 0x43FC
		}
	case 73:
		{ /* '73' */
			return 0x402A
		}
	case 730:
		{ /* '730' */
			return 0x43FC
		}
	case 731:
		{ /* '731' */
			return 0x4354
		}
	case 732:
		{ /* '732' */
			return 0x43FC
		}
	case 733:
		{ /* '733' */
			return 0x43FC
		}
	case 734:
		{ /* '734' */
			return 0x43FC
		}
	case 735:
		{ /* '735' */
			return 0x43FC
		}
	case 736:
		{ /* '736' */
			return 0x43FC
		}
	case 737:
		{ /* '737' */
			return 0x4400
		}
	case 738:
		{ /* '738' */
			return 0x43FE
		}
	case 739:
		{ /* '739' */
			return 0x43FC
		}
	case 74:
		{ /* '74' */
			return 0x4106
		}
	case 740:
		{ /* '740' */
			return 0x43FC
		}
	case 741:
		{ /* '741' */
			return 0x43FC
		}
	case 742:
		{ /* '742' */
			return 0x43FC
		}
	case 743:
		{ /* '743' */
			return 0x43FC
		}
	case 744:
		{ /* '744' */
			return 0x4400
		}
	case 745:
		{ /* '745' */
			return 0x4400
		}
	case 746:
		{ /* '746' */
			return 0x4400
		}
	case 747:
		{ /* '747' */
			return 0x43FE
		}
	case 748:
		{ /* '748' */
			return 0x43FE
		}
	case 749:
		{ /* '749' */
			return 0x43FE
		}
	case 75:
		{ /* '75' */
			return 0x4106
		}
	case 750:
		{ /* '750' */
			return 0x43FE
		}
	case 751:
		{ /* '751' */
			return 0x42A0
		}
	case 752:
		{ /* '752' */
			return 0x43FE
		}
	case 753:
		{ /* '753' */
			return 0x43FE
		}
	case 754:
		{ /* '754' */
			return 0x43FE
		}
	case 755:
		{ /* '755' */
			return 0x43FC
		}
	case 756:
		{ /* '756' */
			return 0x43FC
		}
	case 757:
		{ /* '757' */
			return 0x43FC
		}
	case 758:
		{ /* '758' */
			return 0x402C
		}
	case 759:
		{ /* '759' */
			return 0x402C
		}
	case 76:
		{ /* '76' */
			return 0x4106
		}
	case 760:
		{ /* '760' */
			return 0x43CE
		}
	case 761:
		{ /* '761' */
			return 0x43FC
		}
	case 762:
		{ /* '762' */
			return 0x43FC
		}
	case 763:
		{ /* '763' */
			return 0x4400
		}
	case 764:
		{ /* '764' */
			return 0x4136
		}
	case 765:
		{ /* '765' */
			return 0x437E
		}
	case 766:
		{ /* '766' */
			return 0x439A
		}
	case 767:
		{ /* '767' */
			return 0x4314
		}
	case 768:
		{ /* '768' */
			return 0x431C
		}
	case 769:
		{ /* '769' */
			return 0x41D8
		}
	case 77:
		{ /* '77' */
			return 0x4106
		}
	case 770:
		{ /* '770' */
			return 0x4276
		}
	case 771:
		{ /* '771' */
			return 0x4136
		}
	case 772:
		{ /* '772' */
			return 0x4266
		}
	case 773:
		{ /* '773' */
			return 0x4324
		}
	case 774:
		{ /* '774' */
			return 0x43FC
		}
	case 775:
		{ /* '775' */
			return 0x43FC
		}
	case 776:
		{ /* '776' */
			return 0x43FC
		}
	case 777:
		{ /* '777' */
			return 0x43FC
		}
	case 778:
		{ /* '778' */
			return 0x43FC
		}
	case 779:
		{ /* '779' */
			return 0x43FC
		}
	case 78:
		{ /* '78' */
			return 0x4106
		}
	case 780:
		{ /* '780' */
			return 0x43FC
		}
	case 781:
		{ /* '781' */
			return 0x43F4
		}
	case 782:
		{ /* '782' */
			return 0x43F4
		}
	case 783:
		{ /* '783' */
			return 0x43F4
		}
	case 784:
		{ /* '784' */
			return 0x43F4
		}
	case 785:
		{ /* '785' */
			return 0x43F4
		}
	case 786:
		{ /* '786' */
			return 0x43F4
		}
	case 787:
		{ /* '787' */
			return 0x43F4
		}
	case 788:
		{ /* '788' */
			return 0x43F4
		}
	case 789:
		{ /* '789' */
			return 0x43F4
		}
	case 79:
		{ /* '79' */
			return 0x4106
		}
	case 790:
		{ /* '790' */
			return 0x4400
		}
	case 791:
		{ /* '791' */
			return 0x4400
		}
	case 792:
		{ /* '792' */
			return 0x4400
		}
	case 793:
		{ /* '793' */
			return 0x4400
		}
	case 794:
		{ /* '794' */
			return 0x4400
		}
	case 795:
		{ /* '795' */
			return 0x4400
		}
	case 796:
		{ /* '796' */
			return 0x4400
		}
	case 797:
		{ /* '797' */
			return 0x4400
		}
	case 798:
		{ /* '798' */
			return 0x4400
		}
	case 799:
		{ /* '799' */
			return 0x4400
		}
	case 8:
		{ /* '8' */
			return 0x43FE
		}
	case 80:
		{ /* '80' */
			return 0x4106
		}
	case 800:
		{ /* '800' */
			return 0x43F8
		}
	case 801:
		{ /* '801' */
			return 0x43F8
		}
	case 802:
		{ /* '802' */
			return 0x4258
		}
	case 803:
		{ /* '803' */
			return 0x4400
		}
	case 804:
		{ /* '804' */
			return 0x43FE
		}
	case 805:
		{ /* '805' */
			return 0x43FE
		}
	case 806:
		{ /* '806' */
			return 0x43FF
		}
	case 807:
		{ /* '807' */
			return 0x43FF
		}
	case 808:
		{ /* '808' */
			return 0x43FF
		}
	case 809:
		{ /* '809' */
			return 0x43FF
		}
	case 81:
		{ /* '81' */
			return 0x4106
		}
	case 810:
		{ /* '810' */
			return 0x43FF
		}
	case 811:
		{ /* '811' */
			return 0x43FF
		}
	case 812:
		{ /* '812' */
			return 0x43FF
		}
	case 813:
		{ /* '813' */
			return 0x43FF
		}
	case 814:
		{ /* '814' */
			return 0x43FF
		}
	case 815:
		{ /* '815' */
			return 0x43FF
		}
	case 816:
		{ /* '816' */
			return 0x43FF
		}
	case 817:
		{ /* '817' */
			return 0x43FF
		}
	case 818:
		{ /* '818' */
			return 0x43FF
		}
	case 819:
		{ /* '819' */
			return 0x43FF
		}
	case 82:
		{ /* '82' */
			return 0x4106
		}
	case 820:
		{ /* '820' */
			return 0x43FF
		}
	case 821:
		{ /* '821' */
			return 0x43FF
		}
	case 822:
		{ /* '822' */
			return 0x43FF
		}
	case 823:
		{ /* '823' */
			return 0x4400
		}
	case 824:
		{ /* '824' */
			return 0x4400
		}
	case 825:
		{ /* '825' */
			return 0x4400
		}
	case 826:
		{ /* '826' */
			return 0x4400
		}
	case 827:
		{ /* '827' */
			return 0x4400
		}
	case 828:
		{ /* '828' */
			return 0x4400
		}
	case 829:
		{ /* '829' */
			return 0x4400
		}
	case 83:
		{ /* '83' */
			return 0x40A6
		}
	case 830:
		{ /* '830' */
			return 0x4400
		}
	case 831:
		{ /* '831' */
			return 0x4400
		}
	case 832:
		{ /* '832' */
			return 0x4400
		}
	case 833:
		{ /* '833' */
			return 0x4400
		}
	case 834:
		{ /* '834' */
			return 0x4400
		}
	case 835:
		{ /* '835' */
			return 0x4400
		}
	case 836:
		{ /* '836' */
			return 0x4258
		}
	case 837:
		{ /* '837' */
			return 0x4258
		}
	case 838:
		{ /* '838' */
			return 0x4400
		}
	case 839:
		{ /* '839' */
			return 0x4400
		}
	case 84:
		{ /* '84' */
			return 0x40A6
		}
	case 840:
		{ /* '840' */
			return 0x4400
		}
	case 841:
		{ /* '841' */
			return 0x4400
		}
	case 842:
		{ /* '842' */
			return 0x4400
		}
	case 843:
		{ /* '843' */
			return 0x43FC
		}
	case 844:
		{ /* '844' */
			return 0x4200
		}
	case 845:
		{ /* '845' */
			return 0x4100
		}
	case 846:
		{ /* '846' */
			return 0x4100
		}
	case 847:
		{ /* '847' */
			return 0x4100
		}
	case 848:
		{ /* '848' */
			return 0x41E4
		}
	case 849:
		{ /* '849' */
			return 0x43C4
		}
	case 85:
		{ /* '85' */
			return 0x40A6
		}
	case 850:
		{ /* '850' */
			return 0x4760
		}
	case 851:
		{ /* '851' */
			return 0x42E8
		}
	case 852:
		{ /* '852' */
			return 0x4324
		}
	case 853:
		{ /* '853' */
			return 0x42D4
		}
	case 854:
		{ /* '854' */
			return 0x4144
		}
	case 855:
		{ /* '855' */
			return 0x40CC
		}
	case 856:
		{ /* '856' */
			return 0x41A8
		}
	case 857:
		{ /* '857' */
			return 0x4144
		}
	case 858:
		{ /* '858' */
			return 0x4202
		}
	case 859:
		{ /* '859' */
			return 0x4200
		}
	case 86:
		{ /* '86' */
			return 0x40A6
		}
	case 860:
		{ /* '860' */
			return 0x43FC
		}
	case 861:
		{ /* '861' */
			return 0x4194
		}
	case 862:
		{ /* '862' */
			return 0x4324
		}
	case 863:
		{ /* '863' */
			return 0x4324
		}
	case 864:
		{ /* '864' */
			return 0x4324
		}
	case 865:
		{ /* '865' */
			return 0x4324
		}
	case 866:
		{ /* '866' */
			return 0x4400
		}
	case 867:
		{ /* '867' */
			return 0x4400
		}
	case 868:
		{ /* '868' */
			return 0x4400
		}
	case 869:
		{ /* '869' */
			return 0x4400
		}
	case 87:
		{ /* '87' */
			return 0x4086
		}
	case 870:
		{ /* '870' */
			return 0x4400
		}
	case 871:
		{ /* '871' */
			return 0x4400
		}
	case 872:
		{ /* '872' */
			return 0x4400
		}
	case 873:
		{ /* '873' */
			return 0x4400
		}
	case 874:
		{ /* '874' */
			return 0x4400
		}
	case 875:
		{ /* '875' */
			return 0x4400
		}
	case 876:
		{ /* '876' */
			return 0x43EC
		}
	case 877:
		{ /* '877' */
			return 0x43EC
		}
	case 878:
		{ /* '878' */
			return 0x402C
		}
	case 879:
		{ /* '879' */
			return 0x4400
		}
	case 88:
		{ /* '88' */
			return 0x4086
		}
	case 880:
		{ /* '880' */
			return 0x43EC
		}
	case 881:
		{ /* '881' */
			return 0x43EC
		}
	case 882:
		{ /* '882' */
			return 0x43EC
		}
	case 883:
		{ /* '883' */
			return 0x43EC
		}
	case 884:
		{ /* '884' */
			return 0x43EC
		}
	case 885:
		{ /* '885' */
			return 0x43EC
		}
	case 886:
		{ /* '886' */
			return 0x43FC
		}
	case 887:
		{ /* '887' */
			return 0x4204
		}
	case 888:
		{ /* '888' */
			return 0x4204
		}
	case 889:
		{ /* '889' */
			return 0x43FE
		}
	case 89:
		{ /* '89' */
			return 0x40C2
		}
	case 890:
		{ /* '890' */
			return 0x43FE
		}
	case 891:
		{ /* '891' */
			return 0x43FE
		}
	case 892:
		{ /* '892' */
			return 0x43EC
		}
	case 893:
		{ /* '893' */
			return 0x43EC
		}
	case 894:
		{ /* '894' */
			return 0x43EC
		}
	case 895:
		{ /* '895' */
			return 0x43EC
		}
	case 896:
		{ /* '896' */
			return 0x43EC
		}
	case 897:
		{ /* '897' */
			return 0x43EC
		}
	case 898:
		{ /* '898' */
			return 0x43FE
		}
	case 899:
		{ /* '899' */
			return 0x43FE
		}
	case 9:
		{ /* '9' */
			return 0x43FE
		}
	case 90:
		{ /* '90' */
			return 0x40C4
		}
	case 900:
		{ /* '900' */
			return 0x40F4
		}
	case 901:
		{ /* '901' */
			return 0x40F4
		}
	case 902:
		{ /* '902' */
			return 0x40F4
		}
	case 903:
		{ /* '903' */
			return 0x40F4
		}
	case 904:
		{ /* '904' */
			return 0x40F4
		}
	case 905:
		{ /* '905' */
			return 0x40F4
		}
	case 906:
		{ /* '906' */
			return 0x40F4
		}
	case 907:
		{ /* '907' */
			return 0x43FE
		}
	case 908:
		{ /* '908' */
			return 0x43FE
		}
	case 909:
		{ /* '909' */
			return 0x40F4
		}
	case 91:
		{ /* '91' */
			return 0x40C4
		}
	case 910:
		{ /* '910' */
			return 0x40F4
		}
	case 911:
		{ /* '911' */
			return 0x40F4
		}
	case 912:
		{ /* '912' */
			return 0x43EC
		}
	case 913:
		{ /* '913' */
			return 0x43EC
		}
	case 914:
		{ /* '914' */
			return 0x43EC
		}
	case 915:
		{ /* '915' */
			return 0x4324
		}
	case 916:
		{ /* '916' */
			return 0x4400
		}
	case 917:
		{ /* '917' */
			return 0x4324
		}
	case 918:
		{ /* '918' */
			return 0x4194
		}
	case 919:
		{ /* '919' */
			return 0x4324
		}
	case 92:
		{ /* '92' */
			return 0x43FE
		}
	case 920:
		{ /* '920' */
			return 0x4324
		}
	case 921:
		{ /* '921' */
			return 0x42BC
		}
	case 922:
		{ /* '922' */
			return 0x46A0
		}
	case 923:
		{ /* '923' */
			return 0x4400
		}
	case 924:
		{ /* '924' */
			return 0x41E4
		}
	case 925:
		{ /* '925' */
			return 0x41E4
		}
	case 926:
		{ /* '926' */
			return 0x4400
		}
	case 927:
		{ /* '927' */
			return 0x4400
		}
	case 928:
		{ /* '928' */
			return 0x4400
		}
	case 929:
		{ /* '929' */
			return 0x4400
		}
	case 93:
		{ /* '93' */
			return 0x43FE
		}
	case 930:
		{ /* '930' */
			return 0x43FE
		}
	case 931:
		{ /* '931' */
			return 0x4324
		}
	case 932:
		{ /* '932' */
			return 0x4400
		}
	case 933:
		{ /* '933' */
			return 0x43FE
		}
	case 934:
		{ /* '934' */
			return 0x4400
		}
	case 935:
		{ /* '935' */
			return 0x4400
		}
	case 936:
		{ /* '936' */
			return 0x4400
		}
	case 937:
		{ /* '937' */
			return 0x4400
		}
	case 938:
		{ /* '938' */
			return 0x4400
		}
	case 939:
		{ /* '939' */
			return 0x4400
		}
	case 94:
		{ /* '94' */
			return 0x414E
		}
	case 940:
		{ /* '940' */
			return 0x4400
		}
	case 941:
		{ /* '941' */
			return 0x4400
		}
	case 942:
		{ /* '942' */
			return 0x4144
		}
	case 943:
		{ /* '943' */
			return 0x4760
		}
	case 944:
		{ /* '944' */
			return 0x4400
		}
	case 945:
		{ /* '945' */
			return 0x4194
		}
	case 946:
		{ /* '946' */
			return 0x4194
		}
	case 947:
		{ /* '947' */
			return 0x4324
		}
	case 948:
		{ /* '948' */
			return 0x4324
		}
	case 949:
		{ /* '949' */
			return 0x4324
		}
	case 95:
		{ /* '95' */
			return 0x414E
		}
	case 950:
		{ /* '950' */
			return 0x4324
		}
	case 951:
		{ /* '951' */
			return 0x4324
		}
	case 952:
		{ /* '952' */
			return 0x4324
		}
	case 953:
		{ /* '953' */
			return 0x4194
		}
	case 954:
		{ /* '954' */
			return 0x43FF
		}
	case 955:
		{ /* '955' */
			return 0x4145
		}
	case 956:
		{ /* '956' */
			return 0x4195
		}
	case 957:
		{ /* '957' */
			return 0x4195
		}
	case 958:
		{ /* '958' */
			return 0x4195
		}
	case 959:
		{ /* '959' */
			return 0x4195
		}
	case 96:
		{ /* '96' */
			return 0x43FE
		}
	case 960:
		{ /* '960' */
			return 0x4145
		}
	case 961:
		{ /* '961' */
			return 0x4195
		}
	case 962:
		{ /* '962' */
			return 0x4195
		}
	case 963:
		{ /* '963' */
			return 0x41E5
		}
	case 964:
		{ /* '964' */
			return 0x4195
		}
	case 965:
		{ /* '965' */
			return 0x43FF
		}
	case 966:
		{ /* '966' */
			return 0x41E5
		}
	case 967:
		{ /* '967' */
			return 0x4195
		}
	case 968:
		{ /* '968' */
			return 0x4195
		}
	case 969:
		{ /* '969' */
			return 0x4195
		}
	case 97:
		{ /* '97' */
			return 0x4400
		}
	case 970:
		{ /* '970' */
			return 0x4195
		}
	case 971:
		{ /* '971' */
			return 0x4145
		}
	case 972:
		{ /* '972' */
			return 0x4145
		}
	case 973:
		{ /* '973' */
			return 0x4195
		}
	case 974:
		{ /* '974' */
			return 0x4195
		}
	case 975:
		{ /* '975' */
			return 0x41E5
		}
	case 976:
		{ /* '976' */
			return 0x4195
		}
	case 977:
		{ /* '977' */
			return 0x43FF
		}
	case 978:
		{ /* '978' */
			return 0x41E5
		}
	case 979:
		{ /* '979' */
			return 0x4195
		}
	case 98:
		{ /* '98' */
			return 0x4400
		}
	case 980:
		{ /* '980' */
			return 0x4195
		}
	case 981:
		{ /* '981' */
			return 0x4195
		}
	case 982:
		{ /* '982' */
			return 0x4195
		}
	case 983:
		{ /* '983' */
			return 0x4145
		}
	case 984:
		{ /* '984' */
			return 0x41E5
		}
	case 985:
		{ /* '985' */
			return 0x4195
		}
	case 986:
		{ /* '986' */
			return 0x43FF
		}
	case 987:
		{ /* '987' */
			return 0x41E5
		}
	case 988:
		{ /* '988' */
			return 0x4195
		}
	case 989:
		{ /* '989' */
			return 0x4195
		}
	case 99:
		{ /* '99' */
			return 0x4400
		}
	case 990:
		{ /* '990' */
			return 0x4195
		}
	case 991:
		{ /* '991' */
			return 0x4195
		}
	case 992:
		{ /* '992' */
			return 0x4145
		}
	case 993:
		{ /* '993' */
			return 0x41E5
		}
	case 994:
		{ /* '994' */
			return 0x4195
		}
	case 995:
		{ /* '995' */
			return 0x43FF
		}
	case 996:
		{ /* '996' */
			return 0x41E5
		}
	case 997:
		{ /* '997' */
			return 0x4195
		}
	case 998:
		{ /* '998' */
			return 0x4195
		}
	case 999:
		{ /* '999' */
			return 0x4195
		}
	default:
		{
			return 0
		}
	}
}

func (e DeviceInformation) DeviceDescriptor() uint16 {
	switch e {
	case 1:
		{ /* '1' */
			return 0x0701
		}
	case 10:
		{ /* '10' */
			return 0x0705
		}
	case 100:
		{ /* '100' */
			return 0x0701
		}
	case 1000:
		{ /* '1000' */
			return 0x0705
		}
	case 1001:
		{ /* '1001' */
			return 0x0705
		}
	case 1002:
		{ /* '1002' */
			return 0x0705
		}
	case 1003:
		{ /* '1003' */
			return 0x0705
		}
	case 1004:
		{ /* '1004' */
			return 0x0705
		}
	case 1005:
		{ /* '1005' */
			return 0x0705
		}
	case 1006:
		{ /* '1006' */
			return 0x0705
		}
	case 1007:
		{ /* '1007' */
			return 0x0705
		}
	case 1008:
		{ /* '1008' */
			return 0x0705
		}
	case 1009:
		{ /* '1009' */
			return 0x0705
		}
	case 101:
		{ /* '101' */
			return 0x0701
		}
	case 1010:
		{ /* '1010' */
			return 0x0705
		}
	case 1011:
		{ /* '1011' */
			return 0x0705
		}
	case 1012:
		{ /* '1012' */
			return 0x0705
		}
	case 1013:
		{ /* '1013' */
			return 0x0705
		}
	case 1014:
		{ /* '1014' */
			return 0x0705
		}
	case 1015:
		{ /* '1015' */
			return 0x0701
		}
	case 1016:
		{ /* '1016' */
			return 0x0705
		}
	case 1017:
		{ /* '1017' */
			return 0x0705
		}
	case 1018:
		{ /* '1018' */
			return 0x0705
		}
	case 1019:
		{ /* '1019' */
			return 0x0705
		}
	case 102:
		{ /* '102' */
			return 0x0701
		}
	case 1020:
		{ /* '1020' */
			return 0x0705
		}
	case 1021:
		{ /* '1021' */
			return 0x0705
		}
	case 1022:
		{ /* '1022' */
			return 0x0705
		}
	case 1023:
		{ /* '1023' */
			return 0x0701
		}
	case 1024:
		{ /* '1024' */
			return 0x0701
		}
	case 1025:
		{ /* '1025' */
			return 0x0705
		}
	case 1026:
		{ /* '1026' */
			return 0x0701
		}
	case 1027:
		{ /* '1027' */
			return 0x0701
		}
	case 1028:
		{ /* '1028' */
			return 0x0705
		}
	case 1029:
		{ /* '1029' */
			return 0x0705
		}
	case 103:
		{ /* '103' */
			return 0x0701
		}
	case 1030:
		{ /* '1030' */
			return 0x0705
		}
	case 1031:
		{ /* '1031' */
			return 0x0705
		}
	case 1032:
		{ /* '1032' */
			return 0x0705
		}
	case 1033:
		{ /* '1033' */
			return 0x0705
		}
	case 1034:
		{ /* '1034' */
			return 0x0705
		}
	case 1035:
		{ /* '1035' */
			return 0x0705
		}
	case 1036:
		{ /* '1036' */
			return 0x0705
		}
	case 1037:
		{ /* '1037' */
			return 0x0705
		}
	case 1038:
		{ /* '1038' */
			return 0x0705
		}
	case 1039:
		{ /* '1039' */
			return 0x0705
		}
	case 104:
		{ /* '104' */
			return 0x0701
		}
	case 1040:
		{ /* '1040' */
			return 0x0705
		}
	case 1041:
		{ /* '1041' */
			return 0x0705
		}
	case 1042:
		{ /* '1042' */
			return 0x0705
		}
	case 1043:
		{ /* '1043' */
			return 0x0705
		}
	case 1044:
		{ /* '1044' */
			return 0x0705
		}
	case 1045:
		{ /* '1045' */
			return 0x0705
		}
	case 1046:
		{ /* '1046' */
			return 0x0705
		}
	case 1047:
		{ /* '1047' */
			return 0x0705
		}
	case 1048:
		{ /* '1048' */
			return 0x0705
		}
	case 1049:
		{ /* '1049' */
			return 0x0705
		}
	case 105:
		{ /* '105' */
			return 0x0701
		}
	case 1050:
		{ /* '1050' */
			return 0x0705
		}
	case 1051:
		{ /* '1051' */
			return 0x0705
		}
	case 1052:
		{ /* '1052' */
			return 0x0701
		}
	case 1053:
		{ /* '1053' */
			return 0x0705
		}
	case 1054:
		{ /* '1054' */
			return 0x0701
		}
	case 1055:
		{ /* '1055' */
			return 0x0701
		}
	case 1056:
		{ /* '1056' */
			return 0x0701
		}
	case 1057:
		{ /* '1057' */
			return 0x0701
		}
	case 1058:
		{ /* '1058' */
			return 0x0701
		}
	case 1059:
		{ /* '1059' */
			return 0x0701
		}
	case 106:
		{ /* '106' */
			return 0x0701
		}
	case 1060:
		{ /* '1060' */
			return 0x0701
		}
	case 1061:
		{ /* '1061' */
			return 0x0701
		}
	case 1062:
		{ /* '1062' */
			return 0x0701
		}
	case 1063:
		{ /* '1063' */
			return 0x0701
		}
	case 1064:
		{ /* '1064' */
			return 0x0701
		}
	case 1065:
		{ /* '1065' */
			return 0x0701
		}
	case 1066:
		{ /* '1066' */
			return 0x0701
		}
	case 1067:
		{ /* '1067' */
			return 0x0701
		}
	case 1068:
		{ /* '1068' */
			return 0x0701
		}
	case 1069:
		{ /* '1069' */
			return 0x0701
		}
	case 107:
		{ /* '107' */
			return 0x0705
		}
	case 1070:
		{ /* '1070' */
			return 0x0701
		}
	case 1071:
		{ /* '1071' */
			return 0x0701
		}
	case 1072:
		{ /* '1072' */
			return 0x0701
		}
	case 1073:
		{ /* '1073' */
			return 0x0701
		}
	case 1074:
		{ /* '1074' */
			return 0x0701
		}
	case 1075:
		{ /* '1075' */
			return 0x0701
		}
	case 1076:
		{ /* '1076' */
			return 0x0701
		}
	case 1077:
		{ /* '1077' */
			return 0x0701
		}
	case 1078:
		{ /* '1078' */
			return 0x0701
		}
	case 1079:
		{ /* '1079' */
			return 0x0701
		}
	case 108:
		{ /* '108' */
			return 0x0701
		}
	case 1080:
		{ /* '1080' */
			return 0x0701
		}
	case 1081:
		{ /* '1081' */
			return 0x0701
		}
	case 1082:
		{ /* '1082' */
			return 0x0701
		}
	case 1083:
		{ /* '1083' */
			return 0x0701
		}
	case 1084:
		{ /* '1084' */
			return 0x0701
		}
	case 1085:
		{ /* '1085' */
			return 0x0701
		}
	case 1086:
		{ /* '1086' */
			return 0x0701
		}
	case 1087:
		{ /* '1087' */
			return 0x0701
		}
	case 1088:
		{ /* '1088' */
			return 0x0701
		}
	case 1089:
		{ /* '1089' */
			return 0x0701
		}
	case 109:
		{ /* '109' */
			return 0x0705
		}
	case 1090:
		{ /* '1090' */
			return 0x0701
		}
	case 1091:
		{ /* '1091' */
			return 0x0701
		}
	case 1092:
		{ /* '1092' */
			return 0x0701
		}
	case 1093:
		{ /* '1093' */
			return 0x0701
		}
	case 1094:
		{ /* '1094' */
			return 0x0701
		}
	case 1095:
		{ /* '1095' */
			return 0x0701
		}
	case 1096:
		{ /* '1096' */
			return 0x0701
		}
	case 1097:
		{ /* '1097' */
			return 0x0701
		}
	case 1098:
		{ /* '1098' */
			return 0x0701
		}
	case 1099:
		{ /* '1099' */
			return 0x0701
		}
	case 11:
		{ /* '11' */
			return 0x0705
		}
	case 110:
		{ /* '110' */
			return 0x0701
		}
	case 1100:
		{ /* '1100' */
			return 0x0701
		}
	case 1101:
		{ /* '1101' */
			return 0x0701
		}
	case 1102:
		{ /* '1102' */
			return 0x0701
		}
	case 1103:
		{ /* '1103' */
			return 0x0701
		}
	case 1104:
		{ /* '1104' */
			return 0x0701
		}
	case 1105:
		{ /* '1105' */
			return 0x0701
		}
	case 1106:
		{ /* '1106' */
			return 0x0701
		}
	case 1107:
		{ /* '1107' */
			return 0x0701
		}
	case 1108:
		{ /* '1108' */
			return 0x0705
		}
	case 1109:
		{ /* '1109' */
			return 0x0705
		}
	case 111:
		{ /* '111' */
			return 0x0701
		}
	case 1110:
		{ /* '1110' */
			return 0x0705
		}
	case 1111:
		{ /* '1111' */
			return 0x0705
		}
	case 1112:
		{ /* '1112' */
			return 0x0705
		}
	case 1113:
		{ /* '1113' */
			return 0x0705
		}
	case 1114:
		{ /* '1114' */
			return 0x0701
		}
	case 1115:
		{ /* '1115' */
			return 0x0701
		}
	case 1116:
		{ /* '1116' */
			return 0x0705
		}
	case 1117:
		{ /* '1117' */
			return 0x0705
		}
	case 1118:
		{ /* '1118' */
			return 0x0705
		}
	case 1119:
		{ /* '1119' */
			return 0x0701
		}
	case 112:
		{ /* '112' */
			return 0x0701
		}
	case 1120:
		{ /* '1120' */
			return 0x0701
		}
	case 1121:
		{ /* '1121' */
			return 0x0705
		}
	case 1122:
		{ /* '1122' */
			return 0x0705
		}
	case 1123:
		{ /* '1123' */
			return 0x0701
		}
	case 1124:
		{ /* '1124' */
			return 0x0701
		}
	case 1125:
		{ /* '1125' */
			return 0x0705
		}
	case 1126:
		{ /* '1126' */
			return 0x0705
		}
	case 1127:
		{ /* '1127' */
			return 0x0705
		}
	case 1128:
		{ /* '1128' */
			return 0x0701
		}
	case 1129:
		{ /* '1129' */
			return 0x0701
		}
	case 113:
		{ /* '113' */
			return 0x0701
		}
	case 1130:
		{ /* '1130' */
			return 0x0701
		}
	case 1131:
		{ /* '1131' */
			return 0x0701
		}
	case 1132:
		{ /* '1132' */
			return 0x0701
		}
	case 1133:
		{ /* '1133' */
			return 0x0701
		}
	case 1134:
		{ /* '1134' */
			return 0x0701
		}
	case 1135:
		{ /* '1135' */
			return 0x0701
		}
	case 1136:
		{ /* '1136' */
			return 0x0701
		}
	case 1137:
		{ /* '1137' */
			return 0x0701
		}
	case 1138:
		{ /* '1138' */
			return 0x0701
		}
	case 1139:
		{ /* '1139' */
			return 0x0701
		}
	case 114:
		{ /* '114' */
			return 0x0701
		}
	case 1140:
		{ /* '1140' */
			return 0x0701
		}
	case 1141:
		{ /* '1141' */
			return 0x0701
		}
	case 1142:
		{ /* '1142' */
			return 0x0701
		}
	case 1143:
		{ /* '1143' */
			return 0x0701
		}
	case 1144:
		{ /* '1144' */
			return 0x0701
		}
	case 1145:
		{ /* '1145' */
			return 0x0701
		}
	case 1146:
		{ /* '1146' */
			return 0x0701
		}
	case 1147:
		{ /* '1147' */
			return 0x0705
		}
	case 1148:
		{ /* '1148' */
			return 0x0701
		}
	case 1149:
		{ /* '1149' */
			return 0x0701
		}
	case 115:
		{ /* '115' */
			return 0x0701
		}
	case 1150:
		{ /* '1150' */
			return 0x0701
		}
	case 1151:
		{ /* '1151' */
			return 0x0705
		}
	case 1152:
		{ /* '1152' */
			return 0x0705
		}
	case 1153:
		{ /* '1153' */
			return 0x0701
		}
	case 1154:
		{ /* '1154' */
			return 0x0701
		}
	case 1155:
		{ /* '1155' */
			return 0x0701
		}
	case 1156:
		{ /* '1156' */
			return 0x0701
		}
	case 1157:
		{ /* '1157' */
			return 0x0701
		}
	case 1158:
		{ /* '1158' */
			return 0x0701
		}
	case 1159:
		{ /* '1159' */
			return 0x0701
		}
	case 116:
		{ /* '116' */
			return 0x0701
		}
	case 1160:
		{ /* '1160' */
			return 0x0705
		}
	case 1161:
		{ /* '1161' */
			return 0x0701
		}
	case 1162:
		{ /* '1162' */
			return 0x0701
		}
	case 1163:
		{ /* '1163' */
			return 0x0701
		}
	case 1164:
		{ /* '1164' */
			return 0x0701
		}
	case 1165:
		{ /* '1165' */
			return 0x0705
		}
	case 1166:
		{ /* '1166' */
			return 0x0705
		}
	case 1167:
		{ /* '1167' */
			return 0x0701
		}
	case 1168:
		{ /* '1168' */
			return 0x0701
		}
	case 1169:
		{ /* '1169' */
			return 0x0701
		}
	case 117:
		{ /* '117' */
			return 0x0705
		}
	case 1170:
		{ /* '1170' */
			return 0x0701
		}
	case 1171:
		{ /* '1171' */
			return 0x0701
		}
	case 1172:
		{ /* '1172' */
			return 0x0701
		}
	case 1173:
		{ /* '1173' */
			return 0x0701
		}
	case 1174:
		{ /* '1174' */
			return 0x0701
		}
	case 1175:
		{ /* '1175' */
			return 0x0701
		}
	case 1176:
		{ /* '1176' */
			return 0x0701
		}
	case 1177:
		{ /* '1177' */
			return 0x0701
		}
	case 1178:
		{ /* '1178' */
			return 0x0701
		}
	case 1179:
		{ /* '1179' */
			return 0x0701
		}
	case 118:
		{ /* '118' */
			return 0x0705
		}
	case 1180:
		{ /* '1180' */
			return 0x0705
		}
	case 1181:
		{ /* '1181' */
			return 0x0705
		}
	case 1182:
		{ /* '1182' */
			return 0x0701
		}
	case 1183:
		{ /* '1183' */
			return 0x0701
		}
	case 1184:
		{ /* '1184' */
			return 0x0705
		}
	case 1185:
		{ /* '1185' */
			return 0x0701
		}
	case 1186:
		{ /* '1186' */
			return 0x0701
		}
	case 1187:
		{ /* '1187' */
			return 0x0705
		}
	case 1188:
		{ /* '1188' */
			return 0x0701
		}
	case 1189:
		{ /* '1189' */
			return 0x0701
		}
	case 119:
		{ /* '119' */
			return 0x0701
		}
	case 1190:
		{ /* '1190' */
			return 0x0701
		}
	case 1191:
		{ /* '1191' */
			return 0x0701
		}
	case 1192:
		{ /* '1192' */
			return 0x0701
		}
	case 1193:
		{ /* '1193' */
			return 0x0701
		}
	case 1194:
		{ /* '1194' */
			return 0x0701
		}
	case 1195:
		{ /* '1195' */
			return 0x0701
		}
	case 1196:
		{ /* '1196' */
			return 0x0701
		}
	case 1197:
		{ /* '1197' */
			return 0x0701
		}
	case 1198:
		{ /* '1198' */
			return 0x0701
		}
	case 1199:
		{ /* '1199' */
			return 0x0701
		}
	case 12:
		{ /* '12' */
			return 0x0705
		}
	case 120:
		{ /* '120' */
			return 0x0705
		}
	case 1200:
		{ /* '1200' */
			return 0x0701
		}
	case 1201:
		{ /* '1201' */
			return 0x0701
		}
	case 1202:
		{ /* '1202' */
			return 0x0701
		}
	case 1203:
		{ /* '1203' */
			return 0x0705
		}
	case 1204:
		{ /* '1204' */
			return 0x0701
		}
	case 1205:
		{ /* '1205' */
			return 0x0705
		}
	case 1206:
		{ /* '1206' */
			return 0x0701
		}
	case 1207:
		{ /* '1207' */
			return 0x0701
		}
	case 1208:
		{ /* '1208' */
			return 0x0701
		}
	case 1209:
		{ /* '1209' */
			return 0x0701
		}
	case 121:
		{ /* '121' */
			return 0x0701
		}
	case 1210:
		{ /* '1210' */
			return 0x0701
		}
	case 1211:
		{ /* '1211' */
			return 0x0701
		}
	case 1212:
		{ /* '1212' */
			return 0x0701
		}
	case 1213:
		{ /* '1213' */
			return 0x0701
		}
	case 1214:
		{ /* '1214' */
			return 0x0701
		}
	case 1215:
		{ /* '1215' */
			return 0x0701
		}
	case 1216:
		{ /* '1216' */
			return 0x0701
		}
	case 1217:
		{ /* '1217' */
			return 0x0705
		}
	case 1218:
		{ /* '1218' */
			return 0x0701
		}
	case 1219:
		{ /* '1219' */
			return 0x0705
		}
	case 122:
		{ /* '122' */
			return 0x0701
		}
	case 1220:
		{ /* '1220' */
			return 0x0701
		}
	case 1221:
		{ /* '1221' */
			return 0x0701
		}
	case 1222:
		{ /* '1222' */
			return 0x0701
		}
	case 1223:
		{ /* '1223' */
			return 0x0701
		}
	case 1224:
		{ /* '1224' */
			return 0x0701
		}
	case 1225:
		{ /* '1225' */
			return 0x0701
		}
	case 1226:
		{ /* '1226' */
			return 0x0705
		}
	case 1227:
		{ /* '1227' */
			return 0x0705
		}
	case 1228:
		{ /* '1228' */
			return 0x0705
		}
	case 1229:
		{ /* '1229' */
			return 0x0705
		}
	case 123:
		{ /* '123' */
			return 0x0701
		}
	case 1230:
		{ /* '1230' */
			return 0x0705
		}
	case 1231:
		{ /* '1231' */
			return 0x0705
		}
	case 1232:
		{ /* '1232' */
			return 0x0705
		}
	case 1233:
		{ /* '1233' */
			return 0x0705
		}
	case 1234:
		{ /* '1234' */
			return 0x0705
		}
	case 1235:
		{ /* '1235' */
			return 0x0701
		}
	case 1236:
		{ /* '1236' */
			return 0x0701
		}
	case 1237:
		{ /* '1237' */
			return 0x0705
		}
	case 1238:
		{ /* '1238' */
			return 0x0705
		}
	case 1239:
		{ /* '1239' */
			return 0x0705
		}
	case 124:
		{ /* '124' */
			return 0x0705
		}
	case 1240:
		{ /* '1240' */
			return 0x0705
		}
	case 1241:
		{ /* '1241' */
			return 0x0705
		}
	case 1242:
		{ /* '1242' */
			return 0x0705
		}
	case 1243:
		{ /* '1243' */
			return 0x0705
		}
	case 1244:
		{ /* '1244' */
			return 0x0705
		}
	case 1245:
		{ /* '1245' */
			return 0x0705
		}
	case 1246:
		{ /* '1246' */
			return 0x0705
		}
	case 1247:
		{ /* '1247' */
			return 0x0705
		}
	case 1248:
		{ /* '1248' */
			return 0x0705
		}
	case 1249:
		{ /* '1249' */
			return 0x0705
		}
	case 125:
		{ /* '125' */
			return 0x0705
		}
	case 1250:
		{ /* '1250' */
			return 0x0705
		}
	case 1251:
		{ /* '1251' */
			return 0x0705
		}
	case 1252:
		{ /* '1252' */
			return 0x0705
		}
	case 1253:
		{ /* '1253' */
			return 0x0705
		}
	case 1254:
		{ /* '1254' */
			return 0x0705
		}
	case 1255:
		{ /* '1255' */
			return 0x0705
		}
	case 1256:
		{ /* '1256' */
			return 0x0705
		}
	case 1257:
		{ /* '1257' */
			return 0x0701
		}
	case 1258:
		{ /* '1258' */
			return 0x0705
		}
	case 1259:
		{ /* '1259' */
			return 0x0705
		}
	case 126:
		{ /* '126' */
			return 0x0705
		}
	case 1260:
		{ /* '1260' */
			return 0x0701
		}
	case 1261:
		{ /* '1261' */
			return 0x0705
		}
	case 1262:
		{ /* '1262' */
			return 0x0705
		}
	case 1263:
		{ /* '1263' */
			return 0x0705
		}
	case 1264:
		{ /* '1264' */
			return 0x0705
		}
	case 1265:
		{ /* '1265' */
			return 0x0705
		}
	case 1266:
		{ /* '1266' */
			return 0x0705
		}
	case 1267:
		{ /* '1267' */
			return 0x0705
		}
	case 1268:
		{ /* '1268' */
			return 0x0705
		}
	case 1269:
		{ /* '1269' */
			return 0x0705
		}
	case 127:
		{ /* '127' */
			return 0x0705
		}
	case 1270:
		{ /* '1270' */
			return 0x0705
		}
	case 1271:
		{ /* '1271' */
			return 0x0705
		}
	case 1272:
		{ /* '1272' */
			return 0x0705
		}
	case 1273:
		{ /* '1273' */
			return 0x0705
		}
	case 1274:
		{ /* '1274' */
			return 0x0705
		}
	case 1275:
		{ /* '1275' */
			return 0x0705
		}
	case 1276:
		{ /* '1276' */
			return 0x0705
		}
	case 1277:
		{ /* '1277' */
			return 0x0705
		}
	case 1278:
		{ /* '1278' */
			return 0x0705
		}
	case 1279:
		{ /* '1279' */
			return 0x0705
		}
	case 128:
		{ /* '128' */
			return 0x0705
		}
	case 1280:
		{ /* '1280' */
			return 0x0705
		}
	case 1281:
		{ /* '1281' */
			return 0x0705
		}
	case 1282:
		{ /* '1282' */
			return 0x0705
		}
	case 1283:
		{ /* '1283' */
			return 0x0705
		}
	case 1284:
		{ /* '1284' */
			return 0x0705
		}
	case 1285:
		{ /* '1285' */
			return 0x0705
		}
	case 1286:
		{ /* '1286' */
			return 0x0705
		}
	case 1287:
		{ /* '1287' */
			return 0x0705
		}
	case 1288:
		{ /* '1288' */
			return 0x0705
		}
	case 1289:
		{ /* '1289' */
			return 0x0705
		}
	case 129:
		{ /* '129' */
			return 0x0701
		}
	case 1290:
		{ /* '1290' */
			return 0x0705
		}
	case 1291:
		{ /* '1291' */
			return 0x0705
		}
	case 1292:
		{ /* '1292' */
			return 0x0705
		}
	case 1293:
		{ /* '1293' */
			return 0x0705
		}
	case 1294:
		{ /* '1294' */
			return 0x0705
		}
	case 1295:
		{ /* '1295' */
			return 0x0705
		}
	case 1296:
		{ /* '1296' */
			return 0x0705
		}
	case 1297:
		{ /* '1297' */
			return 0x0705
		}
	case 1298:
		{ /* '1298' */
			return 0x0705
		}
	case 1299:
		{ /* '1299' */
			return 0x0705
		}
	case 13:
		{ /* '13' */
			return 0x0705
		}
	case 130:
		{ /* '130' */
			return 0x0701
		}
	case 1300:
		{ /* '1300' */
			return 0x0701
		}
	case 1301:
		{ /* '1301' */
			return 0x0705
		}
	case 1302:
		{ /* '1302' */
			return 0x0705
		}
	case 1303:
		{ /* '1303' */
			return 0x0705
		}
	case 1304:
		{ /* '1304' */
			return 0x0705
		}
	case 1305:
		{ /* '1305' */
			return 0x0701
		}
	case 1306:
		{ /* '1306' */
			return 0x0701
		}
	case 1307:
		{ /* '1307' */
			return 0x0705
		}
	case 1308:
		{ /* '1308' */
			return 0x0701
		}
	case 1309:
		{ /* '1309' */
			return 0x0701
		}
	case 131:
		{ /* '131' */
			return 0x0705
		}
	case 1310:
		{ /* '1310' */
			return 0x0701
		}
	case 1311:
		{ /* '1311' */
			return 0x0701
		}
	case 1312:
		{ /* '1312' */
			return 0x0701
		}
	case 1313:
		{ /* '1313' */
			return 0x0701
		}
	case 1314:
		{ /* '1314' */
			return 0x0701
		}
	case 1315:
		{ /* '1315' */
			return 0x0701
		}
	case 1316:
		{ /* '1316' */
			return 0x0705
		}
	case 1317:
		{ /* '1317' */
			return 0x0705
		}
	case 1318:
		{ /* '1318' */
			return 0x0701
		}
	case 1319:
		{ /* '1319' */
			return 0x0701
		}
	case 132:
		{ /* '132' */
			return 0x0705
		}
	case 1320:
		{ /* '1320' */
			return 0x0701
		}
	case 1321:
		{ /* '1321' */
			return 0x0701
		}
	case 1322:
		{ /* '1322' */
			return 0x0701
		}
	case 1323:
		{ /* '1323' */
			return 0x0701
		}
	case 1324:
		{ /* '1324' */
			return 0x0701
		}
	case 1325:
		{ /* '1325' */
			return 0x0701
		}
	case 1326:
		{ /* '1326' */
			return 0x0701
		}
	case 1327:
		{ /* '1327' */
			return 0x0701
		}
	case 1328:
		{ /* '1328' */
			return 0x0701
		}
	case 1329:
		{ /* '1329' */
			return 0x0701
		}
	case 133:
		{ /* '133' */
			return 0x0705
		}
	case 1330:
		{ /* '1330' */
			return 0x0701
		}
	case 1331:
		{ /* '1331' */
			return 0x0701
		}
	case 1332:
		{ /* '1332' */
			return 0x0701
		}
	case 1333:
		{ /* '1333' */
			return 0x0701
		}
	case 1334:
		{ /* '1334' */
			return 0x0701
		}
	case 1335:
		{ /* '1335' */
			return 0x0701
		}
	case 1336:
		{ /* '1336' */
			return 0x0701
		}
	case 1337:
		{ /* '1337' */
			return 0x0701
		}
	case 1338:
		{ /* '1338' */
			return 0x0701
		}
	case 1339:
		{ /* '1339' */
			return 0x0701
		}
	case 134:
		{ /* '134' */
			return 0x0705
		}
	case 1340:
		{ /* '1340' */
			return 0x0705
		}
	case 1341:
		{ /* '1341' */
			return 0x0701
		}
	case 1342:
		{ /* '1342' */
			return 0x0701
		}
	case 1343:
		{ /* '1343' */
			return 0x0701
		}
	case 1344:
		{ /* '1344' */
			return 0x0701
		}
	case 1345:
		{ /* '1345' */
			return 0x0701
		}
	case 1346:
		{ /* '1346' */
			return 0x0701
		}
	case 1347:
		{ /* '1347' */
			return 0x0701
		}
	case 1348:
		{ /* '1348' */
			return 0x0701
		}
	case 1349:
		{ /* '1349' */
			return 0x0701
		}
	case 135:
		{ /* '135' */
			return 0x0705
		}
	case 1350:
		{ /* '1350' */
			return 0x0701
		}
	case 1351:
		{ /* '1351' */
			return 0x0701
		}
	case 1352:
		{ /* '1352' */
			return 0x0701
		}
	case 1353:
		{ /* '1353' */
			return 0x0701
		}
	case 1354:
		{ /* '1354' */
			return 0x0701
		}
	case 1355:
		{ /* '1355' */
			return 0x0701
		}
	case 1356:
		{ /* '1356' */
			return 0x0701
		}
	case 1357:
		{ /* '1357' */
			return 0x0701
		}
	case 1358:
		{ /* '1358' */
			return 0x0701
		}
	case 1359:
		{ /* '1359' */
			return 0x0701
		}
	case 136:
		{ /* '136' */
			return 0x0705
		}
	case 1360:
		{ /* '1360' */
			return 0x0705
		}
	case 1361:
		{ /* '1361' */
			return 0x0701
		}
	case 1362:
		{ /* '1362' */
			return 0x0705
		}
	case 1363:
		{ /* '1363' */
			return 0x0701
		}
	case 1364:
		{ /* '1364' */
			return 0x0701
		}
	case 1365:
		{ /* '1365' */
			return 0x0701
		}
	case 1366:
		{ /* '1366' */
			return 0x0701
		}
	case 1367:
		{ /* '1367' */
			return 0x0701
		}
	case 1368:
		{ /* '1368' */
			return 0x0701
		}
	case 1369:
		{ /* '1369' */
			return 0x0701
		}
	case 137:
		{ /* '137' */
			return 0x0705
		}
	case 1370:
		{ /* '1370' */
			return 0x0701
		}
	case 1371:
		{ /* '1371' */
			return 0x0701
		}
	case 1372:
		{ /* '1372' */
			return 0x0701
		}
	case 1373:
		{ /* '1373' */
			return 0x0701
		}
	case 1374:
		{ /* '1374' */
			return 0x0701
		}
	case 1375:
		{ /* '1375' */
			return 0x0701
		}
	case 1376:
		{ /* '1376' */
			return 0x0701
		}
	case 1377:
		{ /* '1377' */
			return 0x0701
		}
	case 1378:
		{ /* '1378' */
			return 0x0701
		}
	case 1379:
		{ /* '1379' */
			return 0x0701
		}
	case 138:
		{ /* '138' */
			return 0x0705
		}
	case 1380:
		{ /* '1380' */
			return 0x0701
		}
	case 1381:
		{ /* '1381' */
			return 0x0701
		}
	case 1382:
		{ /* '1382' */
			return 0x0701
		}
	case 1383:
		{ /* '1383' */
			return 0x0701
		}
	case 1384:
		{ /* '1384' */
			return 0x0701
		}
	case 1385:
		{ /* '1385' */
			return 0x0701
		}
	case 1386:
		{ /* '1386' */
			return 0x0701
		}
	case 1387:
		{ /* '1387' */
			return 0x0701
		}
	case 1388:
		{ /* '1388' */
			return 0x0701
		}
	case 1389:
		{ /* '1389' */
			return 0x0701
		}
	case 139:
		{ /* '139' */
			return 0x0705
		}
	case 1390:
		{ /* '1390' */
			return 0x0701
		}
	case 1391:
		{ /* '1391' */
			return 0x0701
		}
	case 1392:
		{ /* '1392' */
			return 0x0701
		}
	case 1393:
		{ /* '1393' */
			return 0x0701
		}
	case 1394:
		{ /* '1394' */
			return 0x0701
		}
	case 1395:
		{ /* '1395' */
			return 0x0701
		}
	case 1396:
		{ /* '1396' */
			return 0x0701
		}
	case 1397:
		{ /* '1397' */
			return 0x0701
		}
	case 1398:
		{ /* '1398' */
			return 0x0701
		}
	case 1399:
		{ /* '1399' */
			return 0x0701
		}
	case 14:
		{ /* '14' */
			return 0x0705
		}
	case 140:
		{ /* '140' */
			return 0x0705
		}
	case 1400:
		{ /* '1400' */
			return 0x0701
		}
	case 1401:
		{ /* '1401' */
			return 0x0701
		}
	case 1402:
		{ /* '1402' */
			return 0x0701
		}
	case 1403:
		{ /* '1403' */
			return 0x0701
		}
	case 1404:
		{ /* '1404' */
			return 0x0701
		}
	case 1405:
		{ /* '1405' */
			return 0x0701
		}
	case 1406:
		{ /* '1406' */
			return 0x0701
		}
	case 1407:
		{ /* '1407' */
			return 0x0701
		}
	case 1408:
		{ /* '1408' */
			return 0x0701
		}
	case 1409:
		{ /* '1409' */
			return 0x0701
		}
	case 141:
		{ /* '141' */
			return 0x0705
		}
	case 1410:
		{ /* '1410' */
			return 0x0701
		}
	case 1411:
		{ /* '1411' */
			return 0x0701
		}
	case 1412:
		{ /* '1412' */
			return 0x0701
		}
	case 1413:
		{ /* '1413' */
			return 0x0701
		}
	case 1414:
		{ /* '1414' */
			return 0x0701
		}
	case 1415:
		{ /* '1415' */
			return 0x0701
		}
	case 1416:
		{ /* '1416' */
			return 0x0701
		}
	case 1417:
		{ /* '1417' */
			return 0x0701
		}
	case 1418:
		{ /* '1418' */
			return 0x0701
		}
	case 1419:
		{ /* '1419' */
			return 0x0701
		}
	case 142:
		{ /* '142' */
			return 0x0705
		}
	case 1420:
		{ /* '1420' */
			return 0x0701
		}
	case 1421:
		{ /* '1421' */
			return 0x0701
		}
	case 1422:
		{ /* '1422' */
			return 0x0701
		}
	case 1423:
		{ /* '1423' */
			return 0x0701
		}
	case 1424:
		{ /* '1424' */
			return 0x0701
		}
	case 1425:
		{ /* '1425' */
			return 0x0701
		}
	case 1426:
		{ /* '1426' */
			return 0x0701
		}
	case 1427:
		{ /* '1427' */
			return 0x0705
		}
	case 1428:
		{ /* '1428' */
			return 0x0705
		}
	case 1429:
		{ /* '1429' */
			return 0x0701
		}
	case 143:
		{ /* '143' */
			return 0x0705
		}
	case 1430:
		{ /* '1430' */
			return 0x0701
		}
	case 1431:
		{ /* '1431' */
			return 0x0701
		}
	case 1432:
		{ /* '1432' */
			return 0x0701
		}
	case 1433:
		{ /* '1433' */
			return 0x0701
		}
	case 1434:
		{ /* '1434' */
			return 0x0701
		}
	case 1435:
		{ /* '1435' */
			return 0x0701
		}
	case 1436:
		{ /* '1436' */
			return 0x0701
		}
	case 1437:
		{ /* '1437' */
			return 0x0701
		}
	case 1438:
		{ /* '1438' */
			return 0x0701
		}
	case 1439:
		{ /* '1439' */
			return 0x0701
		}
	case 144:
		{ /* '144' */
			return 0x0705
		}
	case 1440:
		{ /* '1440' */
			return 0x0701
		}
	case 1441:
		{ /* '1441' */
			return 0x0701
		}
	case 1442:
		{ /* '1442' */
			return 0x0701
		}
	case 1443:
		{ /* '1443' */
			return 0x0701
		}
	case 1444:
		{ /* '1444' */
			return 0x0701
		}
	case 1445:
		{ /* '1445' */
			return 0x0701
		}
	case 1446:
		{ /* '1446' */
			return 0x0701
		}
	case 1447:
		{ /* '1447' */
			return 0x0705
		}
	case 1448:
		{ /* '1448' */
			return 0x0701
		}
	case 1449:
		{ /* '1449' */
			return 0x0701
		}
	case 145:
		{ /* '145' */
			return 0x0705
		}
	case 1450:
		{ /* '1450' */
			return 0x0701
		}
	case 1451:
		{ /* '1451' */
			return 0x0701
		}
	case 1452:
		{ /* '1452' */
			return 0x0701
		}
	case 1453:
		{ /* '1453' */
			return 0x0701
		}
	case 1454:
		{ /* '1454' */
			return 0x0701
		}
	case 1455:
		{ /* '1455' */
			return 0x0701
		}
	case 1456:
		{ /* '1456' */
			return 0x0701
		}
	case 1457:
		{ /* '1457' */
			return 0x0701
		}
	case 1458:
		{ /* '1458' */
			return 0x0701
		}
	case 1459:
		{ /* '1459' */
			return 0x0701
		}
	case 146:
		{ /* '146' */
			return 0x0701
		}
	case 1460:
		{ /* '1460' */
			return 0x0701
		}
	case 1461:
		{ /* '1461' */
			return 0x0701
		}
	case 1462:
		{ /* '1462' */
			return 0x0701
		}
	case 1463:
		{ /* '1463' */
			return 0x0701
		}
	case 1464:
		{ /* '1464' */
			return 0x0701
		}
	case 1465:
		{ /* '1465' */
			return 0x0701
		}
	case 1466:
		{ /* '1466' */
			return 0x0701
		}
	case 1467:
		{ /* '1467' */
			return 0x0701
		}
	case 1468:
		{ /* '1468' */
			return 0x0701
		}
	case 1469:
		{ /* '1469' */
			return 0x0701
		}
	case 147:
		{ /* '147' */
			return 0x0701
		}
	case 1470:
		{ /* '1470' */
			return 0x0701
		}
	case 1471:
		{ /* '1471' */
			return 0x0701
		}
	case 1472:
		{ /* '1472' */
			return 0x0701
		}
	case 1473:
		{ /* '1473' */
			return 0x0701
		}
	case 1474:
		{ /* '1474' */
			return 0x0701
		}
	case 1475:
		{ /* '1475' */
			return 0x0701
		}
	case 1476:
		{ /* '1476' */
			return 0x0701
		}
	case 1477:
		{ /* '1477' */
			return 0x0701
		}
	case 1478:
		{ /* '1478' */
			return 0x0701
		}
	case 1479:
		{ /* '1479' */
			return 0x0701
		}
	case 148:
		{ /* '148' */
			return 0x0701
		}
	case 1480:
		{ /* '1480' */
			return 0x0701
		}
	case 1481:
		{ /* '1481' */
			return 0x0701
		}
	case 1482:
		{ /* '1482' */
			return 0x0701
		}
	case 1483:
		{ /* '1483' */
			return 0x0701
		}
	case 1484:
		{ /* '1484' */
			return 0x0701
		}
	case 1485:
		{ /* '1485' */
			return 0x0701
		}
	case 1486:
		{ /* '1486' */
			return 0x0701
		}
	case 1487:
		{ /* '1487' */
			return 0x0701
		}
	case 1488:
		{ /* '1488' */
			return 0x0701
		}
	case 1489:
		{ /* '1489' */
			return 0x0701
		}
	case 149:
		{ /* '149' */
			return 0x0701
		}
	case 1490:
		{ /* '1490' */
			return 0x0701
		}
	case 1491:
		{ /* '1491' */
			return 0x0701
		}
	case 1492:
		{ /* '1492' */
			return 0x0701
		}
	case 1493:
		{ /* '1493' */
			return 0x0701
		}
	case 1494:
		{ /* '1494' */
			return 0x0701
		}
	case 1495:
		{ /* '1495' */
			return 0x0701
		}
	case 1496:
		{ /* '1496' */
			return 0x0701
		}
	case 1497:
		{ /* '1497' */
			return 0x0701
		}
	case 1498:
		{ /* '1498' */
			return 0x0701
		}
	case 1499:
		{ /* '1499' */
			return 0x0701
		}
	case 15:
		{ /* '15' */
			return 0x0705
		}
	case 150:
		{ /* '150' */
			return 0x0701
		}
	case 1500:
		{ /* '1500' */
			return 0x0701
		}
	case 1501:
		{ /* '1501' */
			return 0x0701
		}
	case 1502:
		{ /* '1502' */
			return 0x0701
		}
	case 1503:
		{ /* '1503' */
			return 0x0701
		}
	case 1504:
		{ /* '1504' */
			return 0x0701
		}
	case 1505:
		{ /* '1505' */
			return 0x0701
		}
	case 1506:
		{ /* '1506' */
			return 0x0701
		}
	case 1507:
		{ /* '1507' */
			return 0x0701
		}
	case 1508:
		{ /* '1508' */
			return 0x0701
		}
	case 1509:
		{ /* '1509' */
			return 0x0701
		}
	case 151:
		{ /* '151' */
			return 0x0701
		}
	case 1510:
		{ /* '1510' */
			return 0x0701
		}
	case 1511:
		{ /* '1511' */
			return 0x0701
		}
	case 1512:
		{ /* '1512' */
			return 0x0701
		}
	case 1513:
		{ /* '1513' */
			return 0x0701
		}
	case 1514:
		{ /* '1514' */
			return 0x0701
		}
	case 1515:
		{ /* '1515' */
			return 0x0701
		}
	case 1516:
		{ /* '1516' */
			return 0x0701
		}
	case 1517:
		{ /* '1517' */
			return 0x0701
		}
	case 1518:
		{ /* '1518' */
			return 0x0701
		}
	case 1519:
		{ /* '1519' */
			return 0x0701
		}
	case 152:
		{ /* '152' */
			return 0x0701
		}
	case 1520:
		{ /* '1520' */
			return 0x0701
		}
	case 1521:
		{ /* '1521' */
			return 0x0701
		}
	case 1522:
		{ /* '1522' */
			return 0x0701
		}
	case 1523:
		{ /* '1523' */
			return 0x0701
		}
	case 1524:
		{ /* '1524' */
			return 0x0701
		}
	case 1525:
		{ /* '1525' */
			return 0x0701
		}
	case 1526:
		{ /* '1526' */
			return 0x0701
		}
	case 1527:
		{ /* '1527' */
			return 0x0705
		}
	case 1528:
		{ /* '1528' */
			return 0x0705
		}
	case 1529:
		{ /* '1529' */
			return 0x0705
		}
	case 153:
		{ /* '153' */
			return 0x0701
		}
	case 1530:
		{ /* '1530' */
			return 0x0705
		}
	case 1531:
		{ /* '1531' */
			return 0x0705
		}
	case 1532:
		{ /* '1532' */
			return 0x0705
		}
	case 1533:
		{ /* '1533' */
			return 0x0705
		}
	case 1534:
		{ /* '1534' */
			return 0x0705
		}
	case 1535:
		{ /* '1535' */
			return 0x0705
		}
	case 1536:
		{ /* '1536' */
			return 0x0705
		}
	case 1537:
		{ /* '1537' */
			return 0x0701
		}
	case 1538:
		{ /* '1538' */
			return 0x0705
		}
	case 1539:
		{ /* '1539' */
			return 0x0705
		}
	case 154:
		{ /* '154' */
			return 0x0701
		}
	case 1540:
		{ /* '1540' */
			return 0x0705
		}
	case 1541:
		{ /* '1541' */
			return 0x0705
		}
	case 1542:
		{ /* '1542' */
			return 0x0705
		}
	case 1543:
		{ /* '1543' */
			return 0x0705
		}
	case 1544:
		{ /* '1544' */
			return 0x0705
		}
	case 1545:
		{ /* '1545' */
			return 0x0705
		}
	case 1546:
		{ /* '1546' */
			return 0x0705
		}
	case 1547:
		{ /* '1547' */
			return 0x0705
		}
	case 1548:
		{ /* '1548' */
			return 0x0705
		}
	case 1549:
		{ /* '1549' */
			return 0x0705
		}
	case 155:
		{ /* '155' */
			return 0x0701
		}
	case 1550:
		{ /* '1550' */
			return 0x0705
		}
	case 1551:
		{ /* '1551' */
			return 0x0705
		}
	case 1552:
		{ /* '1552' */
			return 0x0705
		}
	case 1553:
		{ /* '1553' */
			return 0x0705
		}
	case 1554:
		{ /* '1554' */
			return 0x0705
		}
	case 1555:
		{ /* '1555' */
			return 0x0705
		}
	case 1556:
		{ /* '1556' */
			return 0x0705
		}
	case 1557:
		{ /* '1557' */
			return 0x0705
		}
	case 1558:
		{ /* '1558' */
			return 0x0705
		}
	case 1559:
		{ /* '1559' */
			return 0x0705
		}
	case 156:
		{ /* '156' */
			return 0x0705
		}
	case 1560:
		{ /* '1560' */
			return 0x0705
		}
	case 1561:
		{ /* '1561' */
			return 0x0705
		}
	case 1562:
		{ /* '1562' */
			return 0x0705
		}
	case 1563:
		{ /* '1563' */
			return 0x0705
		}
	case 1564:
		{ /* '1564' */
			return 0x0705
		}
	case 1565:
		{ /* '1565' */
			return 0x0705
		}
	case 1566:
		{ /* '1566' */
			return 0x0705
		}
	case 1567:
		{ /* '1567' */
			return 0x0705
		}
	case 1568:
		{ /* '1568' */
			return 0x0705
		}
	case 1569:
		{ /* '1569' */
			return 0x0705
		}
	case 157:
		{ /* '157' */
			return 0x0705
		}
	case 1570:
		{ /* '1570' */
			return 0x0705
		}
	case 1571:
		{ /* '1571' */
			return 0x0705
		}
	case 1572:
		{ /* '1572' */
			return 0x0705
		}
	case 1573:
		{ /* '1573' */
			return 0x0705
		}
	case 1574:
		{ /* '1574' */
			return 0x0705
		}
	case 1575:
		{ /* '1575' */
			return 0x0705
		}
	case 1576:
		{ /* '1576' */
			return 0x0705
		}
	case 1577:
		{ /* '1577' */
			return 0x0705
		}
	case 1578:
		{ /* '1578' */
			return 0x0705
		}
	case 1579:
		{ /* '1579' */
			return 0x0705
		}
	case 158:
		{ /* '158' */
			return 0x0705
		}
	case 1580:
		{ /* '1580' */
			return 0x0705
		}
	case 1581:
		{ /* '1581' */
			return 0x0705
		}
	case 1582:
		{ /* '1582' */
			return 0x0705
		}
	case 1583:
		{ /* '1583' */
			return 0x0705
		}
	case 1584:
		{ /* '1584' */
			return 0x0705
		}
	case 1585:
		{ /* '1585' */
			return 0x0705
		}
	case 1586:
		{ /* '1586' */
			return 0x0705
		}
	case 1587:
		{ /* '1587' */
			return 0x0705
		}
	case 1588:
		{ /* '1588' */
			return 0x0705
		}
	case 1589:
		{ /* '1589' */
			return 0x0705
		}
	case 159:
		{ /* '159' */
			return 0x0705
		}
	case 1590:
		{ /* '1590' */
			return 0x0705
		}
	case 1591:
		{ /* '1591' */
			return 0x0705
		}
	case 1592:
		{ /* '1592' */
			return 0x0705
		}
	case 1593:
		{ /* '1593' */
			return 0x0705
		}
	case 1594:
		{ /* '1594' */
			return 0x0705
		}
	case 1595:
		{ /* '1595' */
			return 0x0705
		}
	case 1596:
		{ /* '1596' */
			return 0x0705
		}
	case 1597:
		{ /* '1597' */
			return 0x0705
		}
	case 1598:
		{ /* '1598' */
			return 0x0705
		}
	case 1599:
		{ /* '1599' */
			return 0x0701
		}
	case 16:
		{ /* '16' */
			return 0x0705
		}
	case 160:
		{ /* '160' */
			return 0x0701
		}
	case 1600:
		{ /* '1600' */
			return 0x0701
		}
	case 1601:
		{ /* '1601' */
			return 0x0701
		}
	case 1602:
		{ /* '1602' */
			return 0x0701
		}
	case 1603:
		{ /* '1603' */
			return 0x0701
		}
	case 1604:
		{ /* '1604' */
			return 0x0705
		}
	case 1605:
		{ /* '1605' */
			return 0x0701
		}
	case 1606:
		{ /* '1606' */
			return 0x0705
		}
	case 1607:
		{ /* '1607' */
			return 0x0701
		}
	case 1608:
		{ /* '1608' */
			return 0x0701
		}
	case 1609:
		{ /* '1609' */
			return 0x0701
		}
	case 161:
		{ /* '161' */
			return 0x0701
		}
	case 1610:
		{ /* '1610' */
			return 0x0701
		}
	case 1611:
		{ /* '1611' */
			return 0x0701
		}
	case 1612:
		{ /* '1612' */
			return 0x0701
		}
	case 1613:
		{ /* '1613' */
			return 0x0701
		}
	case 1614:
		{ /* '1614' */
			return 0x0701
		}
	case 1615:
		{ /* '1615' */
			return 0x0701
		}
	case 1616:
		{ /* '1616' */
			return 0x0701
		}
	case 1617:
		{ /* '1617' */
			return 0x0705
		}
	case 1618:
		{ /* '1618' */
			return 0x0701
		}
	case 1619:
		{ /* '1619' */
			return 0x0701
		}
	case 162:
		{ /* '162' */
			return 0x0701
		}
	case 1620:
		{ /* '1620' */
			return 0x0701
		}
	case 1621:
		{ /* '1621' */
			return 0x0701
		}
	case 1622:
		{ /* '1622' */
			return 0x0701
		}
	case 1623:
		{ /* '1623' */
			return 0x0701
		}
	case 1624:
		{ /* '1624' */
			return 0x0701
		}
	case 1625:
		{ /* '1625' */
			return 0x0701
		}
	case 1626:
		{ /* '1626' */
			return 0x0701
		}
	case 1627:
		{ /* '1627' */
			return 0x0701
		}
	case 1628:
		{ /* '1628' */
			return 0x0701
		}
	case 1629:
		{ /* '1629' */
			return 0x0701
		}
	case 163:
		{ /* '163' */
			return 0x0701
		}
	case 1630:
		{ /* '1630' */
			return 0x0701
		}
	case 1631:
		{ /* '1631' */
			return 0x0701
		}
	case 1632:
		{ /* '1632' */
			return 0x0701
		}
	case 1633:
		{ /* '1633' */
			return 0x0701
		}
	case 1634:
		{ /* '1634' */
			return 0x0701
		}
	case 1635:
		{ /* '1635' */
			return 0x0701
		}
	case 1636:
		{ /* '1636' */
			return 0x0701
		}
	case 1637:
		{ /* '1637' */
			return 0x0701
		}
	case 1638:
		{ /* '1638' */
			return 0x0701
		}
	case 1639:
		{ /* '1639' */
			return 0x0701
		}
	case 164:
		{ /* '164' */
			return 0x0701
		}
	case 1640:
		{ /* '1640' */
			return 0x0701
		}
	case 1641:
		{ /* '1641' */
			return 0x0701
		}
	case 1642:
		{ /* '1642' */
			return 0x0701
		}
	case 1643:
		{ /* '1643' */
			return 0x0701
		}
	case 1644:
		{ /* '1644' */
			return 0x0701
		}
	case 1645:
		{ /* '1645' */
			return 0x0701
		}
	case 1646:
		{ /* '1646' */
			return 0x0701
		}
	case 1647:
		{ /* '1647' */
			return 0x0701
		}
	case 1648:
		{ /* '1648' */
			return 0x0701
		}
	case 1649:
		{ /* '1649' */
			return 0x0701
		}
	case 165:
		{ /* '165' */
			return 0x0701
		}
	case 1650:
		{ /* '1650' */
			return 0x0701
		}
	case 1651:
		{ /* '1651' */
			return 0x0701
		}
	case 1652:
		{ /* '1652' */
			return 0x0701
		}
	case 1653:
		{ /* '1653' */
			return 0x0701
		}
	case 1654:
		{ /* '1654' */
			return 0x0701
		}
	case 1655:
		{ /* '1655' */
			return 0x0705
		}
	case 1656:
		{ /* '1656' */
			return 0x0701
		}
	case 1657:
		{ /* '1657' */
			return 0x0701
		}
	case 1658:
		{ /* '1658' */
			return 0x0701
		}
	case 1659:
		{ /* '1659' */
			return 0x0701
		}
	case 166:
		{ /* '166' */
			return 0x0701
		}
	case 1660:
		{ /* '1660' */
			return 0x0701
		}
	case 1661:
		{ /* '1661' */
			return 0x0705
		}
	case 1662:
		{ /* '1662' */
			return 0x0701
		}
	case 1663:
		{ /* '1663' */
			return 0x0701
		}
	case 1664:
		{ /* '1664' */
			return 0x0701
		}
	case 1665:
		{ /* '1665' */
			return 0x0705
		}
	case 1666:
		{ /* '1666' */
			return 0x0701
		}
	case 1667:
		{ /* '1667' */
			return 0x0701
		}
	case 1668:
		{ /* '1668' */
			return 0x0701
		}
	case 1669:
		{ /* '1669' */
			return 0x0701
		}
	case 167:
		{ /* '167' */
			return 0x0701
		}
	case 1670:
		{ /* '1670' */
			return 0x0701
		}
	case 1671:
		{ /* '1671' */
			return 0x0701
		}
	case 1672:
		{ /* '1672' */
			return 0x0701
		}
	case 1673:
		{ /* '1673' */
			return 0x0701
		}
	case 1674:
		{ /* '1674' */
			return 0x0701
		}
	case 1675:
		{ /* '1675' */
			return 0x0701
		}
	case 1676:
		{ /* '1676' */
			return 0x0701
		}
	case 1677:
		{ /* '1677' */
			return 0x0701
		}
	case 1678:
		{ /* '1678' */
			return 0x0701
		}
	case 1679:
		{ /* '1679' */
			return 0x0701
		}
	case 168:
		{ /* '168' */
			return 0x0701
		}
	case 1680:
		{ /* '1680' */
			return 0x0701
		}
	case 1681:
		{ /* '1681' */
			return 0x0701
		}
	case 1682:
		{ /* '1682' */
			return 0x0701
		}
	case 1683:
		{ /* '1683' */
			return 0x0701
		}
	case 1684:
		{ /* '1684' */
			return 0x0701
		}
	case 1685:
		{ /* '1685' */
			return 0x0701
		}
	case 1686:
		{ /* '1686' */
			return 0x0701
		}
	case 1687:
		{ /* '1687' */
			return 0x0701
		}
	case 1688:
		{ /* '1688' */
			return 0x0701
		}
	case 1689:
		{ /* '1689' */
			return 0x0701
		}
	case 169:
		{ /* '169' */
			return 0x0701
		}
	case 1690:
		{ /* '1690' */
			return 0x0701
		}
	case 1691:
		{ /* '1691' */
			return 0x0701
		}
	case 1692:
		{ /* '1692' */
			return 0x0701
		}
	case 1693:
		{ /* '1693' */
			return 0x0701
		}
	case 1694:
		{ /* '1694' */
			return 0x0701
		}
	case 1695:
		{ /* '1695' */
			return 0x0705
		}
	case 1696:
		{ /* '1696' */
			return 0x0705
		}
	case 1697:
		{ /* '1697' */
			return 0x0705
		}
	case 1698:
		{ /* '1698' */
			return 0x0705
		}
	case 1699:
		{ /* '1699' */
			return 0x0701
		}
	case 17:
		{ /* '17' */
			return 0x0705
		}
	case 170:
		{ /* '170' */
			return 0x0701
		}
	case 1700:
		{ /* '1700' */
			return 0x0705
		}
	case 1701:
		{ /* '1701' */
			return 0x0705
		}
	case 1702:
		{ /* '1702' */
			return 0x0705
		}
	case 1703:
		{ /* '1703' */
			return 0x0705
		}
	case 1704:
		{ /* '1704' */
			return 0x0705
		}
	case 1705:
		{ /* '1705' */
			return 0x0705
		}
	case 1706:
		{ /* '1706' */
			return 0x0701
		}
	case 171:
		{ /* '171' */
			return 0x0705
		}
	case 172:
		{ /* '172' */
			return 0x0705
		}
	case 173:
		{ /* '173' */
			return 0x0701
		}
	case 174:
		{ /* '174' */
			return 0x0701
		}
	case 175:
		{ /* '175' */
			return 0x0705
		}
	case 176:
		{ /* '176' */
			return 0x0705
		}
	case 177:
		{ /* '177' */
			return 0x0705
		}
	case 178:
		{ /* '178' */
			return 0x0701
		}
	case 179:
		{ /* '179' */
			return 0x0701
		}
	case 18:
		{ /* '18' */
			return 0x0705
		}
	case 180:
		{ /* '180' */
			return 0x0701
		}
	case 181:
		{ /* '181' */
			return 0x0701
		}
	case 182:
		{ /* '182' */
			return 0x0705
		}
	case 183:
		{ /* '183' */
			return 0x0705
		}
	case 184:
		{ /* '184' */
			return 0x0705
		}
	case 185:
		{ /* '185' */
			return 0x0705
		}
	case 186:
		{ /* '186' */
			return 0x0705
		}
	case 187:
		{ /* '187' */
			return 0x0705
		}
	case 188:
		{ /* '188' */
			return 0x0705
		}
	case 189:
		{ /* '189' */
			return 0x0705
		}
	case 19:
		{ /* '19' */
			return 0x0705
		}
	case 190:
		{ /* '190' */
			return 0x0705
		}
	case 191:
		{ /* '191' */
			return 0x0705
		}
	case 192:
		{ /* '192' */
			return 0x0701
		}
	case 193:
		{ /* '193' */
			return 0x0701
		}
	case 194:
		{ /* '194' */
			return 0x0701
		}
	case 195:
		{ /* '195' */
			return 0x0701
		}
	case 196:
		{ /* '196' */
			return 0x0701
		}
	case 197:
		{ /* '197' */
			return 0x0705
		}
	case 198:
		{ /* '198' */
			return 0x0705
		}
	case 199:
		{ /* '199' */
			return 0x0705
		}
	case 2:
		{ /* '2' */
			return 0x0701
		}
	case 20:
		{ /* '20' */
			return 0x0701
		}
	case 200:
		{ /* '200' */
			return 0x0705
		}
	case 201:
		{ /* '201' */
			return 0x0705
		}
	case 202:
		{ /* '202' */
			return 0x0705
		}
	case 203:
		{ /* '203' */
			return 0x0701
		}
	case 204:
		{ /* '204' */
			return 0x0701
		}
	case 205:
		{ /* '205' */
			return 0x0701
		}
	case 206:
		{ /* '206' */
			return 0x0701
		}
	case 207:
		{ /* '207' */
			return 0x0701
		}
	case 208:
		{ /* '208' */
			return 0x0701
		}
	case 209:
		{ /* '209' */
			return 0x0701
		}
	case 21:
		{ /* '21' */
			return 0x0701
		}
	case 210:
		{ /* '210' */
			return 0x0701
		}
	case 211:
		{ /* '211' */
			return 0x0705
		}
	case 212:
		{ /* '212' */
			return 0x0705
		}
	case 213:
		{ /* '213' */
			return 0x0701
		}
	case 214:
		{ /* '214' */
			return 0x0701
		}
	case 215:
		{ /* '215' */
			return 0x0701
		}
	case 216:
		{ /* '216' */
			return 0x0701
		}
	case 217:
		{ /* '217' */
			return 0x0701
		}
	case 218:
		{ /* '218' */
			return 0x0701
		}
	case 219:
		{ /* '219' */
			return 0x0705
		}
	case 22:
		{ /* '22' */
			return 0x0701
		}
	case 220:
		{ /* '220' */
			return 0x0705
		}
	case 221:
		{ /* '221' */
			return 0x0705
		}
	case 222:
		{ /* '222' */
			return 0x0705
		}
	case 223:
		{ /* '223' */
			return 0x0705
		}
	case 224:
		{ /* '224' */
			return 0x0705
		}
	case 225:
		{ /* '225' */
			return 0x0705
		}
	case 226:
		{ /* '226' */
			return 0x0705
		}
	case 227:
		{ /* '227' */
			return 0x0701
		}
	case 228:
		{ /* '228' */
			return 0x0701
		}
	case 229:
		{ /* '229' */
			return 0x0701
		}
	case 23:
		{ /* '23' */
			return 0x0701
		}
	case 230:
		{ /* '230' */
			return 0x0701
		}
	case 231:
		{ /* '231' */
			return 0x0701
		}
	case 232:
		{ /* '232' */
			return 0x0701
		}
	case 233:
		{ /* '233' */
			return 0x0701
		}
	case 234:
		{ /* '234' */
			return 0x0701
		}
	case 235:
		{ /* '235' */
			return 0x0701
		}
	case 236:
		{ /* '236' */
			return 0x0701
		}
	case 237:
		{ /* '237' */
			return 0x0701
		}
	case 238:
		{ /* '238' */
			return 0x0701
		}
	case 239:
		{ /* '239' */
			return 0x0705
		}
	case 24:
		{ /* '24' */
			return 0x0701
		}
	case 240:
		{ /* '240' */
			return 0x0705
		}
	case 241:
		{ /* '241' */
			return 0x0705
		}
	case 242:
		{ /* '242' */
			return 0x0705
		}
	case 243:
		{ /* '243' */
			return 0x0705
		}
	case 244:
		{ /* '244' */
			return 0x0705
		}
	case 245:
		{ /* '245' */
			return 0x0705
		}
	case 246:
		{ /* '246' */
			return 0x0705
		}
	case 247:
		{ /* '247' */
			return 0x0705
		}
	case 248:
		{ /* '248' */
			return 0x0705
		}
	case 249:
		{ /* '249' */
			return 0x0705
		}
	case 25:
		{ /* '25' */
			return 0x0701
		}
	case 250:
		{ /* '250' */
			return 0x0705
		}
	case 251:
		{ /* '251' */
			return 0x0705
		}
	case 252:
		{ /* '252' */
			return 0x0705
		}
	case 253:
		{ /* '253' */
			return 0x0701
		}
	case 254:
		{ /* '254' */
			return 0x0705
		}
	case 255:
		{ /* '255' */
			return 0x0701
		}
	case 256:
		{ /* '256' */
			return 0x0705
		}
	case 257:
		{ /* '257' */
			return 0x0705
		}
	case 258:
		{ /* '258' */
			return 0x0705
		}
	case 259:
		{ /* '259' */
			return 0x0705
		}
	case 26:
		{ /* '26' */
			return 0x0701
		}
	case 260:
		{ /* '260' */
			return 0x0705
		}
	case 261:
		{ /* '261' */
			return 0x0705
		}
	case 262:
		{ /* '262' */
			return 0x0705
		}
	case 263:
		{ /* '263' */
			return 0x0705
		}
	case 264:
		{ /* '264' */
			return 0x0705
		}
	case 265:
		{ /* '265' */
			return 0x0705
		}
	case 266:
		{ /* '266' */
			return 0x0705
		}
	case 267:
		{ /* '267' */
			return 0x0701
		}
	case 268:
		{ /* '268' */
			return 0x0701
		}
	case 269:
		{ /* '269' */
			return 0x0701
		}
	case 27:
		{ /* '27' */
			return 0x0705
		}
	case 270:
		{ /* '270' */
			return 0x0701
		}
	case 271:
		{ /* '271' */
			return 0x0701
		}
	case 272:
		{ /* '272' */
			return 0x0701
		}
	case 273:
		{ /* '273' */
			return 0x0701
		}
	case 274:
		{ /* '274' */
			return 0x0701
		}
	case 275:
		{ /* '275' */
			return 0x0705
		}
	case 276:
		{ /* '276' */
			return 0x0701
		}
	case 277:
		{ /* '277' */
			return 0x0701
		}
	case 278:
		{ /* '278' */
			return 0x0701
		}
	case 279:
		{ /* '279' */
			return 0x0701
		}
	case 28:
		{ /* '28' */
			return 0x0701
		}
	case 280:
		{ /* '280' */
			return 0x0701
		}
	case 281:
		{ /* '281' */
			return 0x0701
		}
	case 282:
		{ /* '282' */
			return 0x0705
		}
	case 283:
		{ /* '283' */
			return 0x0705
		}
	case 284:
		{ /* '284' */
			return 0x0705
		}
	case 285:
		{ /* '285' */
			return 0x0701
		}
	case 286:
		{ /* '286' */
			return 0x0701
		}
	case 287:
		{ /* '287' */
			return 0x0701
		}
	case 288:
		{ /* '288' */
			return 0x0705
		}
	case 289:
		{ /* '289' */
			return 0x0705
		}
	case 29:
		{ /* '29' */
			return 0x0705
		}
	case 290:
		{ /* '290' */
			return 0x0705
		}
	case 291:
		{ /* '291' */
			return 0x0705
		}
	case 292:
		{ /* '292' */
			return 0x0705
		}
	case 293:
		{ /* '293' */
			return 0x0705
		}
	case 294:
		{ /* '294' */
			return 0x0705
		}
	case 295:
		{ /* '295' */
			return 0x0705
		}
	case 296:
		{ /* '296' */
			return 0x0705
		}
	case 297:
		{ /* '297' */
			return 0x0705
		}
	case 298:
		{ /* '298' */
			return 0x0705
		}
	case 299:
		{ /* '299' */
			return 0x0705
		}
	case 3:
		{ /* '3' */
			return 0x0701
		}
	case 30:
		{ /* '30' */
			return 0x0705
		}
	case 300:
		{ /* '300' */
			return 0x0705
		}
	case 301:
		{ /* '301' */
			return 0x0705
		}
	case 302:
		{ /* '302' */
			return 0x0705
		}
	case 303:
		{ /* '303' */
			return 0x0705
		}
	case 304:
		{ /* '304' */
			return 0x0705
		}
	case 305:
		{ /* '305' */
			return 0x0705
		}
	case 306:
		{ /* '306' */
			return 0x0705
		}
	case 307:
		{ /* '307' */
			return 0x0705
		}
	case 308:
		{ /* '308' */
			return 0x0705
		}
	case 309:
		{ /* '309' */
			return 0x0705
		}
	case 31:
		{ /* '31' */
			return 0x0705
		}
	case 310:
		{ /* '310' */
			return 0x0705
		}
	case 311:
		{ /* '311' */
			return 0x0705
		}
	case 312:
		{ /* '312' */
			return 0x0705
		}
	case 313:
		{ /* '313' */
			return 0x0705
		}
	case 314:
		{ /* '314' */
			return 0x0705
		}
	case 315:
		{ /* '315' */
			return 0x0705
		}
	case 316:
		{ /* '316' */
			return 0x0705
		}
	case 317:
		{ /* '317' */
			return 0x0705
		}
	case 318:
		{ /* '318' */
			return 0x0705
		}
	case 319:
		{ /* '319' */
			return 0x0705
		}
	case 32:
		{ /* '32' */
			return 0x0701
		}
	case 320:
		{ /* '320' */
			return 0x0705
		}
	case 321:
		{ /* '321' */
			return 0x0705
		}
	case 322:
		{ /* '322' */
			return 0x0705
		}
	case 323:
		{ /* '323' */
			return 0x0705
		}
	case 324:
		{ /* '324' */
			return 0x0705
		}
	case 325:
		{ /* '325' */
			return 0x0705
		}
	case 326:
		{ /* '326' */
			return 0x0705
		}
	case 327:
		{ /* '327' */
			return 0x0705
		}
	case 328:
		{ /* '328' */
			return 0x0705
		}
	case 329:
		{ /* '329' */
			return 0x0705
		}
	case 33:
		{ /* '33' */
			return 0x0701
		}
	case 330:
		{ /* '330' */
			return 0x0705
		}
	case 331:
		{ /* '331' */
			return 0x0705
		}
	case 332:
		{ /* '332' */
			return 0x0705
		}
	case 333:
		{ /* '333' */
			return 0x0705
		}
	case 334:
		{ /* '334' */
			return 0x0705
		}
	case 335:
		{ /* '335' */
			return 0x0705
		}
	case 336:
		{ /* '336' */
			return 0x0705
		}
	case 337:
		{ /* '337' */
			return 0x0705
		}
	case 338:
		{ /* '338' */
			return 0x0705
		}
	case 339:
		{ /* '339' */
			return 0x0705
		}
	case 34:
		{ /* '34' */
			return 0x0701
		}
	case 340:
		{ /* '340' */
			return 0x0705
		}
	case 341:
		{ /* '341' */
			return 0x0705
		}
	case 342:
		{ /* '342' */
			return 0x0705
		}
	case 343:
		{ /* '343' */
			return 0x0705
		}
	case 344:
		{ /* '344' */
			return 0x0705
		}
	case 345:
		{ /* '345' */
			return 0x0705
		}
	case 346:
		{ /* '346' */
			return 0x0705
		}
	case 347:
		{ /* '347' */
			return 0x0705
		}
	case 348:
		{ /* '348' */
			return 0x0705
		}
	case 349:
		{ /* '349' */
			return 0x0705
		}
	case 35:
		{ /* '35' */
			return 0x0701
		}
	case 350:
		{ /* '350' */
			return 0x0705
		}
	case 351:
		{ /* '351' */
			return 0x0705
		}
	case 352:
		{ /* '352' */
			return 0x0705
		}
	case 353:
		{ /* '353' */
			return 0x0705
		}
	case 354:
		{ /* '354' */
			return 0x0705
		}
	case 355:
		{ /* '355' */
			return 0x0705
		}
	case 356:
		{ /* '356' */
			return 0x0705
		}
	case 357:
		{ /* '357' */
			return 0x0705
		}
	case 358:
		{ /* '358' */
			return 0x0705
		}
	case 359:
		{ /* '359' */
			return 0x0705
		}
	case 36:
		{ /* '36' */
			return 0x0701
		}
	case 360:
		{ /* '360' */
			return 0x0705
		}
	case 361:
		{ /* '361' */
			return 0x0705
		}
	case 362:
		{ /* '362' */
			return 0x0705
		}
	case 363:
		{ /* '363' */
			return 0x0705
		}
	case 364:
		{ /* '364' */
			return 0x0705
		}
	case 365:
		{ /* '365' */
			return 0x0705
		}
	case 366:
		{ /* '366' */
			return 0x0705
		}
	case 367:
		{ /* '367' */
			return 0x0705
		}
	case 368:
		{ /* '368' */
			return 0x0705
		}
	case 369:
		{ /* '369' */
			return 0x0705
		}
	case 37:
		{ /* '37' */
			return 0x0701
		}
	case 370:
		{ /* '370' */
			return 0x0705
		}
	case 371:
		{ /* '371' */
			return 0x0705
		}
	case 372:
		{ /* '372' */
			return 0x0705
		}
	case 373:
		{ /* '373' */
			return 0x0705
		}
	case 374:
		{ /* '374' */
			return 0x0705
		}
	case 375:
		{ /* '375' */
			return 0x0705
		}
	case 376:
		{ /* '376' */
			return 0x0701
		}
	case 377:
		{ /* '377' */
			return 0x0701
		}
	case 378:
		{ /* '378' */
			return 0x0705
		}
	case 379:
		{ /* '379' */
			return 0x0705
		}
	case 38:
		{ /* '38' */
			return 0x0701
		}
	case 380:
		{ /* '380' */
			return 0x0705
		}
	case 381:
		{ /* '381' */
			return 0x0701
		}
	case 382:
		{ /* '382' */
			return 0x0701
		}
	case 383:
		{ /* '383' */
			return 0x0701
		}
	case 384:
		{ /* '384' */
			return 0x0701
		}
	case 385:
		{ /* '385' */
			return 0x0701
		}
	case 386:
		{ /* '386' */
			return 0x0701
		}
	case 387:
		{ /* '387' */
			return 0x0701
		}
	case 388:
		{ /* '388' */
			return 0x0701
		}
	case 389:
		{ /* '389' */
			return 0x0701
		}
	case 39:
		{ /* '39' */
			return 0x0701
		}
	case 390:
		{ /* '390' */
			return 0x0701
		}
	case 391:
		{ /* '391' */
			return 0x0701
		}
	case 392:
		{ /* '392' */
			return 0x0701
		}
	case 393:
		{ /* '393' */
			return 0x0701
		}
	case 394:
		{ /* '394' */
			return 0x0701
		}
	case 395:
		{ /* '395' */
			return 0x0701
		}
	case 396:
		{ /* '396' */
			return 0x0701
		}
	case 397:
		{ /* '397' */
			return 0x0701
		}
	case 398:
		{ /* '398' */
			return 0x0701
		}
	case 399:
		{ /* '399' */
			return 0x0701
		}
	case 4:
		{ /* '4' */
			return 0x0705
		}
	case 40:
		{ /* '40' */
			return 0x0705
		}
	case 400:
		{ /* '400' */
			return 0x0701
		}
	case 401:
		{ /* '401' */
			return 0x0701
		}
	case 402:
		{ /* '402' */
			return 0x0701
		}
	case 403:
		{ /* '403' */
			return 0x0701
		}
	case 404:
		{ /* '404' */
			return 0x0701
		}
	case 405:
		{ /* '405' */
			return 0x0701
		}
	case 406:
		{ /* '406' */
			return 0x0701
		}
	case 407:
		{ /* '407' */
			return 0x0701
		}
	case 408:
		{ /* '408' */
			return 0x0701
		}
	case 409:
		{ /* '409' */
			return 0x0701
		}
	case 41:
		{ /* '41' */
			return 0x0701
		}
	case 410:
		{ /* '410' */
			return 0x0701
		}
	case 411:
		{ /* '411' */
			return 0x0701
		}
	case 412:
		{ /* '412' */
			return 0x0701
		}
	case 413:
		{ /* '413' */
			return 0x0701
		}
	case 414:
		{ /* '414' */
			return 0x0701
		}
	case 415:
		{ /* '415' */
			return 0x0701
		}
	case 416:
		{ /* '416' */
			return 0x0701
		}
	case 417:
		{ /* '417' */
			return 0x0701
		}
	case 418:
		{ /* '418' */
			return 0x0701
		}
	case 419:
		{ /* '419' */
			return 0x0701
		}
	case 42:
		{ /* '42' */
			return 0x0701
		}
	case 420:
		{ /* '420' */
			return 0x0701
		}
	case 421:
		{ /* '421' */
			return 0x0701
		}
	case 422:
		{ /* '422' */
			return 0x0701
		}
	case 423:
		{ /* '423' */
			return 0x0701
		}
	case 424:
		{ /* '424' */
			return 0x0701
		}
	case 425:
		{ /* '425' */
			return 0x0701
		}
	case 426:
		{ /* '426' */
			return 0x0701
		}
	case 427:
		{ /* '427' */
			return 0x0701
		}
	case 428:
		{ /* '428' */
			return 0x0701
		}
	case 429:
		{ /* '429' */
			return 0x0701
		}
	case 43:
		{ /* '43' */
			return 0x0701
		}
	case 430:
		{ /* '430' */
			return 0x0701
		}
	case 431:
		{ /* '431' */
			return 0x0701
		}
	case 432:
		{ /* '432' */
			return 0x0701
		}
	case 433:
		{ /* '433' */
			return 0x0701
		}
	case 434:
		{ /* '434' */
			return 0x0701
		}
	case 435:
		{ /* '435' */
			return 0x0701
		}
	case 436:
		{ /* '436' */
			return 0x0701
		}
	case 437:
		{ /* '437' */
			return 0x0701
		}
	case 438:
		{ /* '438' */
			return 0x0701
		}
	case 439:
		{ /* '439' */
			return 0x0701
		}
	case 44:
		{ /* '44' */
			return 0x0701
		}
	case 440:
		{ /* '440' */
			return 0x0701
		}
	case 441:
		{ /* '441' */
			return 0x0701
		}
	case 442:
		{ /* '442' */
			return 0x0701
		}
	case 443:
		{ /* '443' */
			return 0x0701
		}
	case 444:
		{ /* '444' */
			return 0x0701
		}
	case 445:
		{ /* '445' */
			return 0x0701
		}
	case 446:
		{ /* '446' */
			return 0x0701
		}
	case 447:
		{ /* '447' */
			return 0x0701
		}
	case 448:
		{ /* '448' */
			return 0x0701
		}
	case 449:
		{ /* '449' */
			return 0x0701
		}
	case 45:
		{ /* '45' */
			return 0x0701
		}
	case 450:
		{ /* '450' */
			return 0x0701
		}
	case 451:
		{ /* '451' */
			return 0x0705
		}
	case 452:
		{ /* '452' */
			return 0x0701
		}
	case 453:
		{ /* '453' */
			return 0x0701
		}
	case 454:
		{ /* '454' */
			return 0x0701
		}
	case 455:
		{ /* '455' */
			return 0x0701
		}
	case 456:
		{ /* '456' */
			return 0x0701
		}
	case 457:
		{ /* '457' */
			return 0x0701
		}
	case 458:
		{ /* '458' */
			return 0x0701
		}
	case 459:
		{ /* '459' */
			return 0x0701
		}
	case 46:
		{ /* '46' */
			return 0x0701
		}
	case 460:
		{ /* '460' */
			return 0x0701
		}
	case 461:
		{ /* '461' */
			return 0x0701
		}
	case 462:
		{ /* '462' */
			return 0x0701
		}
	case 463:
		{ /* '463' */
			return 0x0701
		}
	case 464:
		{ /* '464' */
			return 0x0701
		}
	case 465:
		{ /* '465' */
			return 0x0701
		}
	case 466:
		{ /* '466' */
			return 0x0701
		}
	case 467:
		{ /* '467' */
			return 0x0701
		}
	case 468:
		{ /* '468' */
			return 0x0701
		}
	case 469:
		{ /* '469' */
			return 0x0701
		}
	case 47:
		{ /* '47' */
			return 0x0705
		}
	case 470:
		{ /* '470' */
			return 0x0701
		}
	case 471:
		{ /* '471' */
			return 0x0701
		}
	case 472:
		{ /* '472' */
			return 0x0701
		}
	case 473:
		{ /* '473' */
			return 0x0701
		}
	case 474:
		{ /* '474' */
			return 0x0701
		}
	case 475:
		{ /* '475' */
			return 0x0701
		}
	case 476:
		{ /* '476' */
			return 0x0701
		}
	case 477:
		{ /* '477' */
			return 0x0701
		}
	case 478:
		{ /* '478' */
			return 0x0701
		}
	case 479:
		{ /* '479' */
			return 0x0701
		}
	case 48:
		{ /* '48' */
			return 0x0705
		}
	case 480:
		{ /* '480' */
			return 0x0701
		}
	case 481:
		{ /* '481' */
			return 0x0701
		}
	case 482:
		{ /* '482' */
			return 0x0701
		}
	case 483:
		{ /* '483' */
			return 0x0701
		}
	case 484:
		{ /* '484' */
			return 0x0701
		}
	case 485:
		{ /* '485' */
			return 0x0701
		}
	case 486:
		{ /* '486' */
			return 0x0701
		}
	case 487:
		{ /* '487' */
			return 0x0701
		}
	case 488:
		{ /* '488' */
			return 0x0701
		}
	case 489:
		{ /* '489' */
			return 0x0701
		}
	case 49:
		{ /* '49' */
			return 0x0701
		}
	case 490:
		{ /* '490' */
			return 0x0701
		}
	case 491:
		{ /* '491' */
			return 0x0701
		}
	case 492:
		{ /* '492' */
			return 0x0701
		}
	case 493:
		{ /* '493' */
			return 0x0701
		}
	case 494:
		{ /* '494' */
			return 0x0701
		}
	case 495:
		{ /* '495' */
			return 0x0705
		}
	case 496:
		{ /* '496' */
			return 0x0701
		}
	case 497:
		{ /* '497' */
			return 0x0701
		}
	case 498:
		{ /* '498' */
			return 0x0701
		}
	case 499:
		{ /* '499' */
			return 0x0701
		}
	case 5:
		{ /* '5' */
			return 0x0701
		}
	case 50:
		{ /* '50' */
			return 0x0701
		}
	case 500:
		{ /* '500' */
			return 0x0701
		}
	case 501:
		{ /* '501' */
			return 0x0701
		}
	case 502:
		{ /* '502' */
			return 0x0701
		}
	case 503:
		{ /* '503' */
			return 0x0701
		}
	case 504:
		{ /* '504' */
			return 0x0701
		}
	case 505:
		{ /* '505' */
			return 0x0701
		}
	case 506:
		{ /* '506' */
			return 0x0705
		}
	case 507:
		{ /* '507' */
			return 0x0701
		}
	case 508:
		{ /* '508' */
			return 0x0701
		}
	case 509:
		{ /* '509' */
			return 0x0701
		}
	case 51:
		{ /* '51' */
			return 0x0701
		}
	case 510:
		{ /* '510' */
			return 0x0701
		}
	case 511:
		{ /* '511' */
			return 0x0701
		}
	case 512:
		{ /* '512' */
			return 0x0701
		}
	case 513:
		{ /* '513' */
			return 0x0701
		}
	case 514:
		{ /* '514' */
			return 0x0701
		}
	case 515:
		{ /* '515' */
			return 0x0701
		}
	case 516:
		{ /* '516' */
			return 0x0701
		}
	case 517:
		{ /* '517' */
			return 0x0701
		}
	case 518:
		{ /* '518' */
			return 0x0701
		}
	case 519:
		{ /* '519' */
			return 0x0701
		}
	case 52:
		{ /* '52' */
			return 0x0701
		}
	case 520:
		{ /* '520' */
			return 0x0701
		}
	case 521:
		{ /* '521' */
			return 0x0701
		}
	case 522:
		{ /* '522' */
			return 0x0701
		}
	case 523:
		{ /* '523' */
			return 0x0701
		}
	case 524:
		{ /* '524' */
			return 0x0701
		}
	case 525:
		{ /* '525' */
			return 0x0701
		}
	case 526:
		{ /* '526' */
			return 0x0701
		}
	case 527:
		{ /* '527' */
			return 0x0701
		}
	case 528:
		{ /* '528' */
			return 0x0701
		}
	case 529:
		{ /* '529' */
			return 0x0701
		}
	case 53:
		{ /* '53' */
			return 0x0701
		}
	case 530:
		{ /* '530' */
			return 0x0701
		}
	case 531:
		{ /* '531' */
			return 0x0701
		}
	case 532:
		{ /* '532' */
			return 0x0701
		}
	case 533:
		{ /* '533' */
			return 0x0705
		}
	case 534:
		{ /* '534' */
			return 0x0701
		}
	case 535:
		{ /* '535' */
			return 0x0701
		}
	case 536:
		{ /* '536' */
			return 0x0701
		}
	case 537:
		{ /* '537' */
			return 0x0701
		}
	case 538:
		{ /* '538' */
			return 0x0701
		}
	case 539:
		{ /* '539' */
			return 0x0701
		}
	case 54:
		{ /* '54' */
			return 0x0705
		}
	case 540:
		{ /* '540' */
			return 0x0701
		}
	case 541:
		{ /* '541' */
			return 0x0701
		}
	case 542:
		{ /* '542' */
			return 0x0701
		}
	case 543:
		{ /* '543' */
			return 0x0701
		}
	case 544:
		{ /* '544' */
			return 0x0701
		}
	case 545:
		{ /* '545' */
			return 0x0701
		}
	case 546:
		{ /* '546' */
			return 0x0701
		}
	case 547:
		{ /* '547' */
			return 0x0701
		}
	case 548:
		{ /* '548' */
			return 0x0701
		}
	case 549:
		{ /* '549' */
			return 0x0701
		}
	case 55:
		{ /* '55' */
			return 0x0701
		}
	case 550:
		{ /* '550' */
			return 0x0701
		}
	case 551:
		{ /* '551' */
			return 0x0701
		}
	case 552:
		{ /* '552' */
			return 0x0701
		}
	case 553:
		{ /* '553' */
			return 0x0701
		}
	case 554:
		{ /* '554' */
			return 0x0705
		}
	case 555:
		{ /* '555' */
			return 0x0705
		}
	case 556:
		{ /* '556' */
			return 0x0705
		}
	case 557:
		{ /* '557' */
			return 0x0705
		}
	case 558:
		{ /* '558' */
			return 0x0705
		}
	case 559:
		{ /* '559' */
			return 0x0701
		}
	case 56:
		{ /* '56' */
			return 0x0701
		}
	case 560:
		{ /* '560' */
			return 0x0701
		}
	case 561:
		{ /* '561' */
			return 0x0701
		}
	case 562:
		{ /* '562' */
			return 0x0701
		}
	case 563:
		{ /* '563' */
			return 0x0701
		}
	case 564:
		{ /* '564' */
			return 0x0701
		}
	case 565:
		{ /* '565' */
			return 0x0701
		}
	case 566:
		{ /* '566' */
			return 0x0701
		}
	case 567:
		{ /* '567' */
			return 0x0701
		}
	case 568:
		{ /* '568' */
			return 0x0701
		}
	case 569:
		{ /* '569' */
			return 0x0701
		}
	case 57:
		{ /* '57' */
			return 0x0701
		}
	case 570:
		{ /* '570' */
			return 0x0701
		}
	case 571:
		{ /* '571' */
			return 0x0701
		}
	case 572:
		{ /* '572' */
			return 0x0701
		}
	case 573:
		{ /* '573' */
			return 0x0701
		}
	case 574:
		{ /* '574' */
			return 0x0701
		}
	case 575:
		{ /* '575' */
			return 0x0701
		}
	case 576:
		{ /* '576' */
			return 0x0701
		}
	case 577:
		{ /* '577' */
			return 0x0701
		}
	case 578:
		{ /* '578' */
			return 0x0701
		}
	case 579:
		{ /* '579' */
			return 0x0701
		}
	case 58:
		{ /* '58' */
			return 0x0701
		}
	case 580:
		{ /* '580' */
			return 0x0701
		}
	case 581:
		{ /* '581' */
			return 0x0701
		}
	case 582:
		{ /* '582' */
			return 0x0701
		}
	case 583:
		{ /* '583' */
			return 0x0701
		}
	case 584:
		{ /* '584' */
			return 0x0701
		}
	case 585:
		{ /* '585' */
			return 0x0701
		}
	case 586:
		{ /* '586' */
			return 0x0701
		}
	case 587:
		{ /* '587' */
			return 0x0701
		}
	case 588:
		{ /* '588' */
			return 0x0701
		}
	case 589:
		{ /* '589' */
			return 0x0701
		}
	case 59:
		{ /* '59' */
			return 0x0701
		}
	case 590:
		{ /* '590' */
			return 0x0701
		}
	case 591:
		{ /* '591' */
			return 0x0701
		}
	case 592:
		{ /* '592' */
			return 0x0701
		}
	case 593:
		{ /* '593' */
			return 0x0701
		}
	case 594:
		{ /* '594' */
			return 0x0701
		}
	case 595:
		{ /* '595' */
			return 0x0701
		}
	case 596:
		{ /* '596' */
			return 0x0701
		}
	case 597:
		{ /* '597' */
			return 0x0701
		}
	case 598:
		{ /* '598' */
			return 0x0701
		}
	case 599:
		{ /* '599' */
			return 0x0701
		}
	case 6:
		{ /* '6' */
			return 0x0701
		}
	case 60:
		{ /* '60' */
			return 0x0701
		}
	case 600:
		{ /* '600' */
			return 0x0701
		}
	case 601:
		{ /* '601' */
			return 0x0701
		}
	case 602:
		{ /* '602' */
			return 0x0701
		}
	case 603:
		{ /* '603' */
			return 0x0701
		}
	case 604:
		{ /* '604' */
			return 0x0701
		}
	case 605:
		{ /* '605' */
			return 0x0701
		}
	case 606:
		{ /* '606' */
			return 0x0701
		}
	case 607:
		{ /* '607' */
			return 0x0701
		}
	case 608:
		{ /* '608' */
			return 0x0701
		}
	case 609:
		{ /* '609' */
			return 0x0701
		}
	case 61:
		{ /* '61' */
			return 0x0701
		}
	case 610:
		{ /* '610' */
			return 0x0701
		}
	case 611:
		{ /* '611' */
			return 0x0701
		}
	case 612:
		{ /* '612' */
			return 0x0701
		}
	case 613:
		{ /* '613' */
			return 0x0701
		}
	case 614:
		{ /* '614' */
			return 0x0701
		}
	case 615:
		{ /* '615' */
			return 0x0701
		}
	case 616:
		{ /* '616' */
			return 0x0701
		}
	case 617:
		{ /* '617' */
			return 0x0701
		}
	case 618:
		{ /* '618' */
			return 0x0701
		}
	case 619:
		{ /* '619' */
			return 0x0701
		}
	case 62:
		{ /* '62' */
			return 0x0701
		}
	case 620:
		{ /* '620' */
			return 0x0701
		}
	case 621:
		{ /* '621' */
			return 0x0701
		}
	case 622:
		{ /* '622' */
			return 0x0701
		}
	case 623:
		{ /* '623' */
			return 0x0701
		}
	case 624:
		{ /* '624' */
			return 0x0701
		}
	case 625:
		{ /* '625' */
			return 0x0701
		}
	case 626:
		{ /* '626' */
			return 0x0701
		}
	case 627:
		{ /* '627' */
			return 0x0701
		}
	case 628:
		{ /* '628' */
			return 0x0701
		}
	case 629:
		{ /* '629' */
			return 0x0701
		}
	case 63:
		{ /* '63' */
			return 0x0701
		}
	case 630:
		{ /* '630' */
			return 0x0701
		}
	case 631:
		{ /* '631' */
			return 0x0701
		}
	case 632:
		{ /* '632' */
			return 0x0701
		}
	case 633:
		{ /* '633' */
			return 0x0701
		}
	case 634:
		{ /* '634' */
			return 0x0701
		}
	case 635:
		{ /* '635' */
			return 0x0701
		}
	case 636:
		{ /* '636' */
			return 0x0701
		}
	case 637:
		{ /* '637' */
			return 0x0701
		}
	case 638:
		{ /* '638' */
			return 0x0701
		}
	case 639:
		{ /* '639' */
			return 0x0701
		}
	case 64:
		{ /* '64' */
			return 0x0701
		}
	case 640:
		{ /* '640' */
			return 0x0701
		}
	case 641:
		{ /* '641' */
			return 0x0701
		}
	case 642:
		{ /* '642' */
			return 0x0701
		}
	case 643:
		{ /* '643' */
			return 0x0701
		}
	case 644:
		{ /* '644' */
			return 0x0701
		}
	case 645:
		{ /* '645' */
			return 0x0701
		}
	case 646:
		{ /* '646' */
			return 0x0701
		}
	case 647:
		{ /* '647' */
			return 0x0701
		}
	case 648:
		{ /* '648' */
			return 0x0701
		}
	case 649:
		{ /* '649' */
			return 0x0701
		}
	case 65:
		{ /* '65' */
			return 0x0701
		}
	case 650:
		{ /* '650' */
			return 0x0701
		}
	case 651:
		{ /* '651' */
			return 0x0701
		}
	case 652:
		{ /* '652' */
			return 0x0701
		}
	case 653:
		{ /* '653' */
			return 0x0701
		}
	case 654:
		{ /* '654' */
			return 0x0701
		}
	case 655:
		{ /* '655' */
			return 0x0701
		}
	case 656:
		{ /* '656' */
			return 0x0701
		}
	case 657:
		{ /* '657' */
			return 0x0701
		}
	case 658:
		{ /* '658' */
			return 0x0701
		}
	case 659:
		{ /* '659' */
			return 0x0701
		}
	case 66:
		{ /* '66' */
			return 0x0701
		}
	case 660:
		{ /* '660' */
			return 0x0701
		}
	case 661:
		{ /* '661' */
			return 0x0701
		}
	case 662:
		{ /* '662' */
			return 0x0701
		}
	case 663:
		{ /* '663' */
			return 0x0701
		}
	case 664:
		{ /* '664' */
			return 0x0701
		}
	case 665:
		{ /* '665' */
			return 0x0701
		}
	case 666:
		{ /* '666' */
			return 0x0701
		}
	case 667:
		{ /* '667' */
			return 0x0701
		}
	case 668:
		{ /* '668' */
			return 0x0701
		}
	case 669:
		{ /* '669' */
			return 0x0701
		}
	case 67:
		{ /* '67' */
			return 0x0701
		}
	case 670:
		{ /* '670' */
			return 0x0701
		}
	case 671:
		{ /* '671' */
			return 0x0701
		}
	case 672:
		{ /* '672' */
			return 0x0701
		}
	case 673:
		{ /* '673' */
			return 0x0701
		}
	case 674:
		{ /* '674' */
			return 0x0701
		}
	case 675:
		{ /* '675' */
			return 0x0701
		}
	case 676:
		{ /* '676' */
			return 0x0701
		}
	case 677:
		{ /* '677' */
			return 0x0701
		}
	case 678:
		{ /* '678' */
			return 0x0701
		}
	case 679:
		{ /* '679' */
			return 0x0705
		}
	case 68:
		{ /* '68' */
			return 0x0701
		}
	case 680:
		{ /* '680' */
			return 0x0701
		}
	case 681:
		{ /* '681' */
			return 0x0701
		}
	case 682:
		{ /* '682' */
			return 0x0701
		}
	case 683:
		{ /* '683' */
			return 0x0701
		}
	case 684:
		{ /* '684' */
			return 0x0701
		}
	case 685:
		{ /* '685' */
			return 0x0701
		}
	case 686:
		{ /* '686' */
			return 0x0705
		}
	case 687:
		{ /* '687' */
			return 0x0701
		}
	case 688:
		{ /* '688' */
			return 0x0701
		}
	case 689:
		{ /* '689' */
			return 0x0701
		}
	case 69:
		{ /* '69' */
			return 0x0701
		}
	case 690:
		{ /* '690' */
			return 0x0705
		}
	case 691:
		{ /* '691' */
			return 0x0705
		}
	case 692:
		{ /* '692' */
			return 0x0705
		}
	case 693:
		{ /* '693' */
			return 0x0705
		}
	case 694:
		{ /* '694' */
			return 0x0705
		}
	case 695:
		{ /* '695' */
			return 0x0705
		}
	case 696:
		{ /* '696' */
			return 0x0705
		}
	case 697:
		{ /* '697' */
			return 0x0705
		}
	case 698:
		{ /* '698' */
			return 0x0705
		}
	case 699:
		{ /* '699' */
			return 0x0701
		}
	case 7:
		{ /* '7' */
			return 0x0705
		}
	case 70:
		{ /* '70' */
			return 0x0701
		}
	case 700:
		{ /* '700' */
			return 0x0701
		}
	case 701:
		{ /* '701' */
			return 0x0701
		}
	case 702:
		{ /* '702' */
			return 0x0705
		}
	case 703:
		{ /* '703' */
			return 0x0705
		}
	case 704:
		{ /* '704' */
			return 0x0705
		}
	case 705:
		{ /* '705' */
			return 0x0705
		}
	case 706:
		{ /* '706' */
			return 0x0705
		}
	case 707:
		{ /* '707' */
			return 0x0705
		}
	case 708:
		{ /* '708' */
			return 0x0705
		}
	case 709:
		{ /* '709' */
			return 0x0705
		}
	case 71:
		{ /* '71' */
			return 0x0705
		}
	case 710:
		{ /* '710' */
			return 0x0705
		}
	case 711:
		{ /* '711' */
			return 0x0705
		}
	case 712:
		{ /* '712' */
			return 0x0705
		}
	case 713:
		{ /* '713' */
			return 0x0705
		}
	case 714:
		{ /* '714' */
			return 0x0705
		}
	case 715:
		{ /* '715' */
			return 0x0705
		}
	case 716:
		{ /* '716' */
			return 0x0705
		}
	case 717:
		{ /* '717' */
			return 0x0705
		}
	case 718:
		{ /* '718' */
			return 0x0705
		}
	case 719:
		{ /* '719' */
			return 0x0705
		}
	case 72:
		{ /* '72' */
			return 0x0701
		}
	case 720:
		{ /* '720' */
			return 0x0705
		}
	case 721:
		{ /* '721' */
			return 0x0705
		}
	case 722:
		{ /* '722' */
			return 0x0705
		}
	case 723:
		{ /* '723' */
			return 0x0705
		}
	case 724:
		{ /* '724' */
			return 0x0705
		}
	case 725:
		{ /* '725' */
			return 0x0705
		}
	case 726:
		{ /* '726' */
			return 0x0705
		}
	case 727:
		{ /* '727' */
			return 0x0705
		}
	case 728:
		{ /* '728' */
			return 0x0705
		}
	case 729:
		{ /* '729' */
			return 0x0705
		}
	case 73:
		{ /* '73' */
			return 0x0701
		}
	case 730:
		{ /* '730' */
			return 0x0705
		}
	case 731:
		{ /* '731' */
			return 0x0705
		}
	case 732:
		{ /* '732' */
			return 0x0705
		}
	case 733:
		{ /* '733' */
			return 0x0705
		}
	case 734:
		{ /* '734' */
			return 0x0705
		}
	case 735:
		{ /* '735' */
			return 0x0705
		}
	case 736:
		{ /* '736' */
			return 0x0705
		}
	case 737:
		{ /* '737' */
			return 0x0701
		}
	case 738:
		{ /* '738' */
			return 0x0705
		}
	case 739:
		{ /* '739' */
			return 0x0705
		}
	case 74:
		{ /* '74' */
			return 0x0701
		}
	case 740:
		{ /* '740' */
			return 0x0705
		}
	case 741:
		{ /* '741' */
			return 0x0705
		}
	case 742:
		{ /* '742' */
			return 0x0705
		}
	case 743:
		{ /* '743' */
			return 0x0705
		}
	case 744:
		{ /* '744' */
			return 0x0705
		}
	case 745:
		{ /* '745' */
			return 0x0705
		}
	case 746:
		{ /* '746' */
			return 0x0705
		}
	case 747:
		{ /* '747' */
			return 0x0705
		}
	case 748:
		{ /* '748' */
			return 0x0705
		}
	case 749:
		{ /* '749' */
			return 0x0701
		}
	case 75:
		{ /* '75' */
			return 0x0701
		}
	case 750:
		{ /* '750' */
			return 0x0701
		}
	case 751:
		{ /* '751' */
			return 0x0701
		}
	case 752:
		{ /* '752' */
			return 0x0701
		}
	case 753:
		{ /* '753' */
			return 0x0701
		}
	case 754:
		{ /* '754' */
			return 0x0701
		}
	case 755:
		{ /* '755' */
			return 0x0705
		}
	case 756:
		{ /* '756' */
			return 0x0701
		}
	case 757:
		{ /* '757' */
			return 0x0701
		}
	case 758:
		{ /* '758' */
			return 0x0701
		}
	case 759:
		{ /* '759' */
			return 0x0701
		}
	case 76:
		{ /* '76' */
			return 0x0701
		}
	case 760:
		{ /* '760' */
			return 0x0701
		}
	case 761:
		{ /* '761' */
			return 0x0705
		}
	case 762:
		{ /* '762' */
			return 0x0705
		}
	case 763:
		{ /* '763' */
			return 0x0705
		}
	case 764:
		{ /* '764' */
			return 0x0705
		}
	case 765:
		{ /* '765' */
			return 0x0705
		}
	case 766:
		{ /* '766' */
			return 0x0705
		}
	case 767:
		{ /* '767' */
			return 0x0705
		}
	case 768:
		{ /* '768' */
			return 0x0705
		}
	case 769:
		{ /* '769' */
			return 0x0705
		}
	case 77:
		{ /* '77' */
			return 0x0701
		}
	case 770:
		{ /* '770' */
			return 0x0705
		}
	case 771:
		{ /* '771' */
			return 0x0705
		}
	case 772:
		{ /* '772' */
			return 0x0705
		}
	case 773:
		{ /* '773' */
			return 0x0705
		}
	case 774:
		{ /* '774' */
			return 0x0705
		}
	case 775:
		{ /* '775' */
			return 0x0705
		}
	case 776:
		{ /* '776' */
			return 0x0705
		}
	case 777:
		{ /* '777' */
			return 0x0705
		}
	case 778:
		{ /* '778' */
			return 0x0705
		}
	case 779:
		{ /* '779' */
			return 0x0705
		}
	case 78:
		{ /* '78' */
			return 0x0701
		}
	case 780:
		{ /* '780' */
			return 0x0705
		}
	case 781:
		{ /* '781' */
			return 0x0705
		}
	case 782:
		{ /* '782' */
			return 0x0705
		}
	case 783:
		{ /* '783' */
			return 0x0705
		}
	case 784:
		{ /* '784' */
			return 0x0705
		}
	case 785:
		{ /* '785' */
			return 0x0705
		}
	case 786:
		{ /* '786' */
			return 0x0705
		}
	case 787:
		{ /* '787' */
			return 0x0705
		}
	case 788:
		{ /* '788' */
			return 0x0705
		}
	case 789:
		{ /* '789' */
			return 0x0705
		}
	case 79:
		{ /* '79' */
			return 0x0701
		}
	case 790:
		{ /* '790' */
			return 0x0701
		}
	case 791:
		{ /* '791' */
			return 0x0701
		}
	case 792:
		{ /* '792' */
			return 0x0701
		}
	case 793:
		{ /* '793' */
			return 0x0701
		}
	case 794:
		{ /* '794' */
			return 0x0701
		}
	case 795:
		{ /* '795' */
			return 0x0701
		}
	case 796:
		{ /* '796' */
			return 0x0701
		}
	case 797:
		{ /* '797' */
			return 0x0701
		}
	case 798:
		{ /* '798' */
			return 0x0701
		}
	case 799:
		{ /* '799' */
			return 0x0701
		}
	case 8:
		{ /* '8' */
			return 0x0705
		}
	case 80:
		{ /* '80' */
			return 0x0701
		}
	case 800:
		{ /* '800' */
			return 0x0705
		}
	case 801:
		{ /* '801' */
			return 0x0705
		}
	case 802:
		{ /* '802' */
			return 0x0705
		}
	case 803:
		{ /* '803' */
			return 0x0701
		}
	case 804:
		{ /* '804' */
			return 0x0705
		}
	case 805:
		{ /* '805' */
			return 0x0705
		}
	case 806:
		{ /* '806' */
			return 0x0705
		}
	case 807:
		{ /* '807' */
			return 0x0705
		}
	case 808:
		{ /* '808' */
			return 0x0705
		}
	case 809:
		{ /* '809' */
			return 0x0705
		}
	case 81:
		{ /* '81' */
			return 0x0701
		}
	case 810:
		{ /* '810' */
			return 0x0705
		}
	case 811:
		{ /* '811' */
			return 0x0705
		}
	case 812:
		{ /* '812' */
			return 0x0705
		}
	case 813:
		{ /* '813' */
			return 0x0705
		}
	case 814:
		{ /* '814' */
			return 0x0705
		}
	case 815:
		{ /* '815' */
			return 0x0705
		}
	case 816:
		{ /* '816' */
			return 0x0705
		}
	case 817:
		{ /* '817' */
			return 0x0705
		}
	case 818:
		{ /* '818' */
			return 0x0705
		}
	case 819:
		{ /* '819' */
			return 0x0705
		}
	case 82:
		{ /* '82' */
			return 0x0701
		}
	case 820:
		{ /* '820' */
			return 0x0705
		}
	case 821:
		{ /* '821' */
			return 0x0705
		}
	case 822:
		{ /* '822' */
			return 0x0705
		}
	case 823:
		{ /* '823' */
			return 0x0701
		}
	case 824:
		{ /* '824' */
			return 0x0701
		}
	case 825:
		{ /* '825' */
			return 0x0701
		}
	case 826:
		{ /* '826' */
			return 0x0701
		}
	case 827:
		{ /* '827' */
			return 0x0701
		}
	case 828:
		{ /* '828' */
			return 0x0701
		}
	case 829:
		{ /* '829' */
			return 0x0705
		}
	case 83:
		{ /* '83' */
			return 0x0701
		}
	case 830:
		{ /* '830' */
			return 0x0705
		}
	case 831:
		{ /* '831' */
			return 0x0705
		}
	case 832:
		{ /* '832' */
			return 0x0705
		}
	case 833:
		{ /* '833' */
			return 0x0705
		}
	case 834:
		{ /* '834' */
			return 0x0705
		}
	case 835:
		{ /* '835' */
			return 0x0705
		}
	case 836:
		{ /* '836' */
			return 0x0705
		}
	case 837:
		{ /* '837' */
			return 0x0705
		}
	case 838:
		{ /* '838' */
			return 0x0701
		}
	case 839:
		{ /* '839' */
			return 0x0701
		}
	case 84:
		{ /* '84' */
			return 0x0701
		}
	case 840:
		{ /* '840' */
			return 0x0701
		}
	case 841:
		{ /* '841' */
			return 0x0701
		}
	case 842:
		{ /* '842' */
			return 0x0701
		}
	case 843:
		{ /* '843' */
			return 0x0705
		}
	case 844:
		{ /* '844' */
			return 0x0705
		}
	case 845:
		{ /* '845' */
			return 0x0705
		}
	case 846:
		{ /* '846' */
			return 0x0705
		}
	case 847:
		{ /* '847' */
			return 0x0705
		}
	case 848:
		{ /* '848' */
			return 0x0705
		}
	case 849:
		{ /* '849' */
			return 0x0705
		}
	case 85:
		{ /* '85' */
			return 0x0701
		}
	case 850:
		{ /* '850' */
			return 0x0705
		}
	case 851:
		{ /* '851' */
			return 0x0705
		}
	case 852:
		{ /* '852' */
			return 0x0705
		}
	case 853:
		{ /* '853' */
			return 0x0705
		}
	case 854:
		{ /* '854' */
			return 0x0705
		}
	case 855:
		{ /* '855' */
			return 0x0705
		}
	case 856:
		{ /* '856' */
			return 0x0705
		}
	case 857:
		{ /* '857' */
			return 0x0705
		}
	case 858:
		{ /* '858' */
			return 0x0705
		}
	case 859:
		{ /* '859' */
			return 0x0705
		}
	case 86:
		{ /* '86' */
			return 0x0701
		}
	case 860:
		{ /* '860' */
			return 0x0701
		}
	case 861:
		{ /* '861' */
			return 0x0705
		}
	case 862:
		{ /* '862' */
			return 0x0705
		}
	case 863:
		{ /* '863' */
			return 0x0705
		}
	case 864:
		{ /* '864' */
			return 0x0705
		}
	case 865:
		{ /* '865' */
			return 0x0705
		}
	case 866:
		{ /* '866' */
			return 0x0705
		}
	case 867:
		{ /* '867' */
			return 0x0705
		}
	case 868:
		{ /* '868' */
			return 0x0705
		}
	case 869:
		{ /* '869' */
			return 0x0705
		}
	case 87:
		{ /* '87' */
			return 0x0701
		}
	case 870:
		{ /* '870' */
			return 0x0705
		}
	case 871:
		{ /* '871' */
			return 0x0705
		}
	case 872:
		{ /* '872' */
			return 0x0705
		}
	case 873:
		{ /* '873' */
			return 0x0705
		}
	case 874:
		{ /* '874' */
			return 0x0705
		}
	case 875:
		{ /* '875' */
			return 0x0705
		}
	case 876:
		{ /* '876' */
			return 0x0705
		}
	case 877:
		{ /* '877' */
			return 0x0705
		}
	case 878:
		{ /* '878' */
			return 0x0701
		}
	case 879:
		{ /* '879' */
			return 0x0705
		}
	case 88:
		{ /* '88' */
			return 0x0701
		}
	case 880:
		{ /* '880' */
			return 0x0705
		}
	case 881:
		{ /* '881' */
			return 0x0705
		}
	case 882:
		{ /* '882' */
			return 0x0705
		}
	case 883:
		{ /* '883' */
			return 0x0705
		}
	case 884:
		{ /* '884' */
			return 0x0705
		}
	case 885:
		{ /* '885' */
			return 0x0705
		}
	case 886:
		{ /* '886' */
			return 0x0705
		}
	case 887:
		{ /* '887' */
			return 0x0701
		}
	case 888:
		{ /* '888' */
			return 0x0701
		}
	case 889:
		{ /* '889' */
			return 0x0705
		}
	case 89:
		{ /* '89' */
			return 0x0701
		}
	case 890:
		{ /* '890' */
			return 0x0705
		}
	case 891:
		{ /* '891' */
			return 0x0705
		}
	case 892:
		{ /* '892' */
			return 0x0705
		}
	case 893:
		{ /* '893' */
			return 0x0705
		}
	case 894:
		{ /* '894' */
			return 0x0705
		}
	case 895:
		{ /* '895' */
			return 0x0705
		}
	case 896:
		{ /* '896' */
			return 0x0705
		}
	case 897:
		{ /* '897' */
			return 0x0705
		}
	case 898:
		{ /* '898' */
			return 0x0705
		}
	case 899:
		{ /* '899' */
			return 0x0705
		}
	case 9:
		{ /* '9' */
			return 0x0705
		}
	case 90:
		{ /* '90' */
			return 0x0701
		}
	case 900:
		{ /* '900' */
			return 0x0701
		}
	case 901:
		{ /* '901' */
			return 0x0701
		}
	case 902:
		{ /* '902' */
			return 0x0701
		}
	case 903:
		{ /* '903' */
			return 0x0701
		}
	case 904:
		{ /* '904' */
			return 0x0701
		}
	case 905:
		{ /* '905' */
			return 0x0701
		}
	case 906:
		{ /* '906' */
			return 0x0701
		}
	case 907:
		{ /* '907' */
			return 0x0705
		}
	case 908:
		{ /* '908' */
			return 0x0705
		}
	case 909:
		{ /* '909' */
			return 0x0701
		}
	case 91:
		{ /* '91' */
			return 0x0701
		}
	case 910:
		{ /* '910' */
			return 0x0701
		}
	case 911:
		{ /* '911' */
			return 0x0701
		}
	case 912:
		{ /* '912' */
			return 0x0705
		}
	case 913:
		{ /* '913' */
			return 0x0705
		}
	case 914:
		{ /* '914' */
			return 0x0705
		}
	case 915:
		{ /* '915' */
			return 0x0705
		}
	case 916:
		{ /* '916' */
			return 0x0705
		}
	case 917:
		{ /* '917' */
			return 0x0705
		}
	case 918:
		{ /* '918' */
			return 0x0705
		}
	case 919:
		{ /* '919' */
			return 0x0705
		}
	case 92:
		{ /* '92' */
			return 0x0701
		}
	case 920:
		{ /* '920' */
			return 0x0705
		}
	case 921:
		{ /* '921' */
			return 0x0705
		}
	case 922:
		{ /* '922' */
			return 0x0701
		}
	case 923:
		{ /* '923' */
			return 0x0701
		}
	case 924:
		{ /* '924' */
			return 0x0705
		}
	case 925:
		{ /* '925' */
			return 0x0705
		}
	case 926:
		{ /* '926' */
			return 0x0701
		}
	case 927:
		{ /* '927' */
			return 0x0701
		}
	case 928:
		{ /* '928' */
			return 0x0701
		}
	case 929:
		{ /* '929' */
			return 0x0701
		}
	case 93:
		{ /* '93' */
			return 0x0701
		}
	case 930:
		{ /* '930' */
			return 0x0701
		}
	case 931:
		{ /* '931' */
			return 0x0701
		}
	case 932:
		{ /* '932' */
			return 0x0705
		}
	case 933:
		{ /* '933' */
			return 0x0701
		}
	case 934:
		{ /* '934' */
			return 0x0701
		}
	case 935:
		{ /* '935' */
			return 0x0701
		}
	case 936:
		{ /* '936' */
			return 0x0705
		}
	case 937:
		{ /* '937' */
			return 0x0705
		}
	case 938:
		{ /* '938' */
			return 0x0705
		}
	case 939:
		{ /* '939' */
			return 0x0705
		}
	case 94:
		{ /* '94' */
			return 0x0701
		}
	case 940:
		{ /* '940' */
			return 0x0705
		}
	case 941:
		{ /* '941' */
			return 0x0705
		}
	case 942:
		{ /* '942' */
			return 0x0705
		}
	case 943:
		{ /* '943' */
			return 0x0705
		}
	case 944:
		{ /* '944' */
			return 0x0705
		}
	case 945:
		{ /* '945' */
			return 0x0705
		}
	case 946:
		{ /* '946' */
			return 0x0705
		}
	case 947:
		{ /* '947' */
			return 0x0705
		}
	case 948:
		{ /* '948' */
			return 0x0705
		}
	case 949:
		{ /* '949' */
			return 0x0705
		}
	case 95:
		{ /* '95' */
			return 0x0701
		}
	case 950:
		{ /* '950' */
			return 0x0705
		}
	case 951:
		{ /* '951' */
			return 0x0705
		}
	case 952:
		{ /* '952' */
			return 0x0705
		}
	case 953:
		{ /* '953' */
			return 0x0705
		}
	case 954:
		{ /* '954' */
			return 0x0705
		}
	case 955:
		{ /* '955' */
			return 0x0705
		}
	case 956:
		{ /* '956' */
			return 0x0705
		}
	case 957:
		{ /* '957' */
			return 0x0705
		}
	case 958:
		{ /* '958' */
			return 0x0705
		}
	case 959:
		{ /* '959' */
			return 0x0705
		}
	case 96:
		{ /* '96' */
			return 0x0705
		}
	case 960:
		{ /* '960' */
			return 0x0705
		}
	case 961:
		{ /* '961' */
			return 0x0705
		}
	case 962:
		{ /* '962' */
			return 0x0705
		}
	case 963:
		{ /* '963' */
			return 0x0705
		}
	case 964:
		{ /* '964' */
			return 0x0705
		}
	case 965:
		{ /* '965' */
			return 0x0705
		}
	case 966:
		{ /* '966' */
			return 0x0705
		}
	case 967:
		{ /* '967' */
			return 0x0705
		}
	case 968:
		{ /* '968' */
			return 0x0705
		}
	case 969:
		{ /* '969' */
			return 0x0705
		}
	case 97:
		{ /* '97' */
			return 0x0701
		}
	case 970:
		{ /* '970' */
			return 0x0705
		}
	case 971:
		{ /* '971' */
			return 0x0705
		}
	case 972:
		{ /* '972' */
			return 0x0705
		}
	case 973:
		{ /* '973' */
			return 0x0705
		}
	case 974:
		{ /* '974' */
			return 0x0705
		}
	case 975:
		{ /* '975' */
			return 0x0705
		}
	case 976:
		{ /* '976' */
			return 0x0705
		}
	case 977:
		{ /* '977' */
			return 0x0705
		}
	case 978:
		{ /* '978' */
			return 0x0705
		}
	case 979:
		{ /* '979' */
			return 0x0705
		}
	case 98:
		{ /* '98' */
			return 0x0701
		}
	case 980:
		{ /* '980' */
			return 0x0705
		}
	case 981:
		{ /* '981' */
			return 0x0705
		}
	case 982:
		{ /* '982' */
			return 0x0705
		}
	case 983:
		{ /* '983' */
			return 0x0705
		}
	case 984:
		{ /* '984' */
			return 0x0705
		}
	case 985:
		{ /* '985' */
			return 0x0705
		}
	case 986:
		{ /* '986' */
			return 0x0705
		}
	case 987:
		{ /* '987' */
			return 0x0705
		}
	case 988:
		{ /* '988' */
			return 0x0705
		}
	case 989:
		{ /* '989' */
			return 0x0705
		}
	case 99:
		{ /* '99' */
			return 0x0701
		}
	case 990:
		{ /* '990' */
			return 0x0705
		}
	case 991:
		{ /* '991' */
			return 0x0705
		}
	case 992:
		{ /* '992' */
			return 0x0705
		}
	case 993:
		{ /* '993' */
			return 0x0705
		}
	case 994:
		{ /* '994' */
			return 0x0705
		}
	case 995:
		{ /* '995' */
			return 0x0705
		}
	case 996:
		{ /* '996' */
			return 0x0705
		}
	case 997:
		{ /* '997' */
			return 0x0705
		}
	case 998:
		{ /* '998' */
			return 0x0705
		}
	case 999:
		{ /* '999' */
			return 0x0705
		}
	default:
		{
			return 0
		}
	}
}
func DeviceInformationByValue(value uint16) DeviceInformation {
	switch value {
	case 1:
		return DeviceInformation_DEV0001914201
	case 10:
		return DeviceInformation_DEV0064181910
	case 100:
		return DeviceInformation_DEV000C133410
	case 1000:
		return DeviceInformation_DEV0004109C14
	case 1001:
		return DeviceInformation_DEV000410A611
	case 1002:
		return DeviceInformation_DEV0004146B13
	case 1003:
		return DeviceInformation_DEV0004147211
	case 1004:
		return DeviceInformation_DEV000410FE12
	case 1005:
		return DeviceInformation_DEV0004209016
	case 1006:
		return DeviceInformation_DEV000420A011
	case 1007:
		return DeviceInformation_DEV0004209011
	case 1008:
		return DeviceInformation_DEV000420CA11
	case 1009:
		return DeviceInformation_DEV0004208012
	case 101:
		return DeviceInformation_DEV000C133310
	case 1010:
		return DeviceInformation_DEV0004207812
	case 1011:
		return DeviceInformation_DEV000420BA11
	case 1012:
		return DeviceInformation_DEV000420B311
	case 1013:
		return DeviceInformation_DEV0004209811
	case 1014:
		return DeviceInformation_DEV0004208811
	case 1015:
		return DeviceInformation_DEV0004B00812
	case 1016:
		return DeviceInformation_DEV0004302613
	case 1017:
		return DeviceInformation_DEV0004302313
	case 1018:
		return DeviceInformation_DEV0004302013
	case 1019:
		return DeviceInformation_DEV0004302B12
	case 102:
		return DeviceInformation_DEV000C133611
	case 1020:
		return DeviceInformation_DEV0004706811
	case 1021:
		return DeviceInformation_DEV0004705D11
	case 1022:
		return DeviceInformation_DEV0004705C11
	case 1023:
		return DeviceInformation_DEV0004B00713
	case 1024:
		return DeviceInformation_DEV0004B00A01
	case 1025:
		return DeviceInformation_DEV0004706611
	case 1026:
		return DeviceInformation_DEV0004C01410
	case 1027:
		return DeviceInformation_DEV0004C00102
	case 1028:
		return DeviceInformation_DEV0004705E11
	case 1029:
		return DeviceInformation_DEV0004706211
	case 103:
		return DeviceInformation_DEV000C133510
	case 1030:
		return DeviceInformation_DEV0004706411
	case 1031:
		return DeviceInformation_DEV0004706412
	case 1032:
		return DeviceInformation_DEV000420C011
	case 1033:
		return DeviceInformation_DEV000420B011
	case 1034:
		return DeviceInformation_DEV0004B00911
	case 1035:
		return DeviceInformation_DEV0004A01211
	case 1036:
		return DeviceInformation_DEV0004A01112
	case 1037:
		return DeviceInformation_DEV0004A01111
	case 1038:
		return DeviceInformation_DEV0004A01212
	case 1039:
		return DeviceInformation_DEV0004A03312
	case 104:
		return DeviceInformation_DEV000C130710
	case 1040:
		return DeviceInformation_DEV0004A03212
	case 1041:
		return DeviceInformation_DEV0004A01113
	case 1042:
		return DeviceInformation_DEV0004A01711
	case 1043:
		return DeviceInformation_DEV000420C711
	case 1044:
		return DeviceInformation_DEV000420BD11
	case 1045:
		return DeviceInformation_DEV000420C411
	case 1046:
		return DeviceInformation_DEV000420A812
	case 1047:
		return DeviceInformation_DEV000420CD11
	case 1048:
		return DeviceInformation_DEV000420AD11
	case 1049:
		return DeviceInformation_DEV000420B611
	case 105:
		return DeviceInformation_DEV000C760210
	case 1050:
		return DeviceInformation_DEV000420A811
	case 1051:
		return DeviceInformation_DEV0004501311
	case 1052:
		return DeviceInformation_DEV0004501411
	case 1053:
		return DeviceInformation_DEV0004B01002
	case 1054:
		return DeviceInformation_DEV0006D00610
	case 1055:
		return DeviceInformation_DEV0006D01510
	case 1056:
		return DeviceInformation_DEV0006D00110
	case 1057:
		return DeviceInformation_DEV0006D00310
	case 1058:
		return DeviceInformation_DEV0006D03210
	case 1059:
		return DeviceInformation_DEV0006D03310
	case 106:
		return DeviceInformation_DEV000C1BD610
	case 1060:
		return DeviceInformation_DEV0006D02E20
	case 1061:
		return DeviceInformation_DEV0006D02F20
	case 1062:
		return DeviceInformation_DEV0006D03020
	case 1063:
		return DeviceInformation_DEV0006D03120
	case 1064:
		return DeviceInformation_DEV0006D02110
	case 1065:
		return DeviceInformation_DEV0006D00010
	case 1066:
		return DeviceInformation_DEV0006D01810
	case 1067:
		return DeviceInformation_DEV0006D00910
	case 1068:
		return DeviceInformation_DEV0006D01110
	case 1069:
		return DeviceInformation_DEV0006D03510
	case 107:
		return DeviceInformation_DEV000C181610
	case 1070:
		return DeviceInformation_DEV0006D03410
	case 1071:
		return DeviceInformation_DEV0006D02410
	case 1072:
		return DeviceInformation_DEV0006D02510
	case 1073:
		return DeviceInformation_DEV0006D00810
	case 1074:
		return DeviceInformation_DEV0006D00710
	case 1075:
		return DeviceInformation_DEV0006D01310
	case 1076:
		return DeviceInformation_DEV0006D01410
	case 1077:
		return DeviceInformation_DEV0006D00210
	case 1078:
		return DeviceInformation_DEV0006D00510
	case 1079:
		return DeviceInformation_DEV0006D00410
	case 108:
		return DeviceInformation_DEV000C648B10
	case 1080:
		return DeviceInformation_DEV0006D02210
	case 1081:
		return DeviceInformation_DEV0006D02310
	case 1082:
		return DeviceInformation_DEV0006D01710
	case 1083:
		return DeviceInformation_DEV0006D01610
	case 1084:
		return DeviceInformation_DEV0006D01010
	case 1085:
		return DeviceInformation_DEV0006D01210
	case 1086:
		return DeviceInformation_DEV0006D04820
	case 1087:
		return DeviceInformation_DEV0006D04C11
	case 1088:
		return DeviceInformation_DEV0006D05610
	case 1089:
		return DeviceInformation_DEV0006D02910
	case 109:
		return DeviceInformation_DEV000C480611
	case 1090:
		return DeviceInformation_DEV0006D02A10
	case 1091:
		return DeviceInformation_DEV0006D02B10
	case 1092:
		return DeviceInformation_DEV0006D02C10
	case 1093:
		return DeviceInformation_DEV0006D02810
	case 1094:
		return DeviceInformation_DEV0006D02610
	case 1095:
		return DeviceInformation_DEV0006D02710
	case 1096:
		return DeviceInformation_DEV0006D03610
	case 1097:
		return DeviceInformation_DEV0006D03710
	case 1098:
		return DeviceInformation_DEV0006D02D11
	case 1099:
		return DeviceInformation_DEV0006D03C10
	case 11:
		return DeviceInformation_DEV0064181810
	case 110:
		return DeviceInformation_DEV000C482011
	case 1100:
		return DeviceInformation_DEV0006D03B10
	case 1101:
		return DeviceInformation_DEV0006D03910
	case 1102:
		return DeviceInformation_DEV0006D03A10
	case 1103:
		return DeviceInformation_DEV0006D03D11
	case 1104:
		return DeviceInformation_DEV0006D03E10
	case 1105:
		return DeviceInformation_DEV0006C00102
	case 1106:
		return DeviceInformation_DEV0006E05611
	case 1107:
		return DeviceInformation_DEV0006E05212
	case 1108:
		return DeviceInformation_DEV000620B011
	case 1109:
		return DeviceInformation_DEV000620B311
	case 111:
		return DeviceInformation_DEV000C724010
	case 1110:
		return DeviceInformation_DEV000620C011
	case 1111:
		return DeviceInformation_DEV000620BA11
	case 1112:
		return DeviceInformation_DEV0006705C11
	case 1113:
		return DeviceInformation_DEV0006705D11
	case 1114:
		return DeviceInformation_DEV0006E07710
	case 1115:
		return DeviceInformation_DEV0006E07712
	case 1116:
		return DeviceInformation_DEV0006706210
	case 1117:
		return DeviceInformation_DEV0006302611
	case 1118:
		return DeviceInformation_DEV0006302612
	case 1119:
		return DeviceInformation_DEV0006E00810
	case 112:
		return DeviceInformation_DEV000C570211
	case 1120:
		return DeviceInformation_DEV0006E01F01
	case 1121:
		return DeviceInformation_DEV0006302311
	case 1122:
		return DeviceInformation_DEV0006302312
	case 1123:
		return DeviceInformation_DEV0006E00910
	case 1124:
		return DeviceInformation_DEV0006E02001
	case 1125:
		return DeviceInformation_DEV0006302011
	case 1126:
		return DeviceInformation_DEV0006302012
	case 1127:
		return DeviceInformation_DEV0006C00C13
	case 1128:
		return DeviceInformation_DEV0006E00811
	case 1129:
		return DeviceInformation_DEV0006E00911
	case 113:
		return DeviceInformation_DEV000C570310
	case 1130:
		return DeviceInformation_DEV0006E01F20
	case 1131:
		return DeviceInformation_DEV0006E03410
	case 1132:
		return DeviceInformation_DEV0006E03110
	case 1133:
		return DeviceInformation_DEV0006E0A210
	case 1134:
		return DeviceInformation_DEV0006E0CE10
	case 1135:
		return DeviceInformation_DEV0006E0A111
	case 1136:
		return DeviceInformation_DEV0006E0CD11
	case 1137:
		return DeviceInformation_DEV0006E02020
	case 1138:
		return DeviceInformation_DEV0006E02D11
	case 1139:
		return DeviceInformation_DEV0006E03011
	case 114:
		return DeviceInformation_DEV000C570411
	case 1140:
		return DeviceInformation_DEV0006E0C110
	case 1141:
		return DeviceInformation_DEV0006E0C510
	case 1142:
		return DeviceInformation_DEV0006B00A01
	case 1143:
		return DeviceInformation_DEV0006B00602
	case 1144:
		return DeviceInformation_DEV0006E0C410
	case 1145:
		return DeviceInformation_DEV0006E0C312
	case 1146:
		return DeviceInformation_DEV0006E0C210
	case 1147:
		return DeviceInformation_DEV0006209016
	case 1148:
		return DeviceInformation_DEV0006E01A01
	case 1149:
		return DeviceInformation_DEV0006E09910
	case 115:
		return DeviceInformation_DEV000C570110
	case 1150:
		return DeviceInformation_DEV0006E03710
	case 1151:
		return DeviceInformation_DEV0006209011
	case 1152:
		return DeviceInformation_DEV000620A011
	case 1153:
		return DeviceInformation_DEV0006E02410
	case 1154:
		return DeviceInformation_DEV0006E02301
	case 1155:
		return DeviceInformation_DEV0006E02510
	case 1156:
		return DeviceInformation_DEV0006E01B01
	case 1157:
		return DeviceInformation_DEV0006E01C01
	case 1158:
		return DeviceInformation_DEV0006E01D01
	case 1159:
		return DeviceInformation_DEV0006E01E01
	case 116:
		return DeviceInformation_DEV000C570011
	case 1160:
		return DeviceInformation_DEV0006207812
	case 1161:
		return DeviceInformation_DEV0006B00811
	case 1162:
		return DeviceInformation_DEV0006E01001
	case 1163:
		return DeviceInformation_DEV0006E03610
	case 1164:
		return DeviceInformation_DEV0006E09810
	case 1165:
		return DeviceInformation_DEV0006208811
	case 1166:
		return DeviceInformation_DEV0006209811
	case 1167:
		return DeviceInformation_DEV0006E02610
	case 1168:
		return DeviceInformation_DEV0006E02710
	case 1169:
		return DeviceInformation_DEV0006E02A10
	case 117:
		return DeviceInformation_DEV000C20BD11
	case 1170:
		return DeviceInformation_DEV0006E02B10
	case 1171:
		return DeviceInformation_DEV0006E00C10
	case 1172:
		return DeviceInformation_DEV0006010110
	case 1173:
		return DeviceInformation_DEV0006010210
	case 1174:
		return DeviceInformation_DEV0006E00B10
	case 1175:
		return DeviceInformation_DEV0006E09C10
	case 1176:
		return DeviceInformation_DEV0006E09B10
	case 1177:
		return DeviceInformation_DEV0006E03510
	case 1178:
		return DeviceInformation_DEV0006FF1B11
	case 1179:
		return DeviceInformation_DEV0006E0CF10
	case 118:
		return DeviceInformation_DEV000C20BA11
	case 1180:
		return DeviceInformation_DEV000620A812
	case 1181:
		return DeviceInformation_DEV000620CD11
	case 1182:
		return DeviceInformation_DEV0006E00E01
	case 1183:
		return DeviceInformation_DEV0006E02201
	case 1184:
		return DeviceInformation_DEV000620AD11
	case 1185:
		return DeviceInformation_DEV0006E00F01
	case 1186:
		return DeviceInformation_DEV0006E02101
	case 1187:
		return DeviceInformation_DEV000620BD11
	case 1188:
		return DeviceInformation_DEV0006E00D01
	case 1189:
		return DeviceInformation_DEV0006E03910
	case 119:
		return DeviceInformation_DEV000C760110
	case 1190:
		return DeviceInformation_DEV0006E02810
	case 1191:
		return DeviceInformation_DEV0006E02910
	case 1192:
		return DeviceInformation_DEV0006E02C10
	case 1193:
		return DeviceInformation_DEV0006C00403
	case 1194:
		return DeviceInformation_DEV0006590101
	case 1195:
		return DeviceInformation_DEV0006E0CC11
	case 1196:
		return DeviceInformation_DEV0006E09A10
	case 1197:
		return DeviceInformation_DEV0006E03811
	case 1198:
		return DeviceInformation_DEV0006E0C710
	case 1199:
		return DeviceInformation_DEV0006E0C610
	case 12:
		return DeviceInformation_DEV0064181710
	case 120:
		return DeviceInformation_DEV000C705C01
	case 1200:
		return DeviceInformation_DEV0006E05A10
	case 1201:
		return DeviceInformation_DEV0048493A1C
	case 1202:
		return DeviceInformation_DEV0048494712
	case 1203:
		return DeviceInformation_DEV0048494810
	case 1204:
		return DeviceInformation_DEV0048855A10
	case 1205:
		return DeviceInformation_DEV0048855B10
	case 1206:
		return DeviceInformation_DEV0048A05713
	case 1207:
		return DeviceInformation_DEV0048494414
	case 1208:
		return DeviceInformation_DEV0048824A11
	case 1209:
		return DeviceInformation_DEV0048824A12
	case 121:
		return DeviceInformation_DEV000CFF2112
	case 1210:
		return DeviceInformation_DEV0048770A10
	case 1211:
		return DeviceInformation_DEV0048494311
	case 1212:
		return DeviceInformation_DEV0048494513
	case 1213:
		return DeviceInformation_DEV0048494012
	case 1214:
		return DeviceInformation_DEV0048494111
	case 1215:
		return DeviceInformation_DEV0048494210
	case 1216:
		return DeviceInformation_DEV0048494610
	case 1217:
		return DeviceInformation_DEV0048494910
	case 1218:
		return DeviceInformation_DEV0048134A10
	case 1219:
		return DeviceInformation_DEV0048107E12
	case 122:
		return DeviceInformation_DEV000CB00812
	case 1220:
		return DeviceInformation_DEV0048FF2112
	case 1221:
		return DeviceInformation_DEV0048140A11
	case 1222:
		return DeviceInformation_DEV0048140B12
	case 1223:
		return DeviceInformation_DEV0048140C13
	case 1224:
		return DeviceInformation_DEV0048139A10
	case 1225:
		return DeviceInformation_DEV0048648B10
	case 1226:
		return DeviceInformation_DEV0008A01111
	case 1227:
		return DeviceInformation_DEV0008A01211
	case 1228:
		return DeviceInformation_DEV0008A01212
	case 1229:
		return DeviceInformation_DEV0008A01112
	case 123:
		return DeviceInformation_DEV000CB00713
	case 1230:
		return DeviceInformation_DEV0008A03213
	case 1231:
		return DeviceInformation_DEV0008A03313
	case 1232:
		return DeviceInformation_DEV0008A01113
	case 1233:
		return DeviceInformation_DEV0008A01711
	case 1234:
		return DeviceInformation_DEV0008B00911
	case 1235:
		return DeviceInformation_DEV0008C00102
	case 1236:
		return DeviceInformation_DEV0008C00101
	case 1237:
		return DeviceInformation_DEV0008901501
	case 1238:
		return DeviceInformation_DEV0008901310
	case 1239:
		return DeviceInformation_DEV000820B011
	case 124:
		return DeviceInformation_DEV000C181910
	case 1240:
		return DeviceInformation_DEV0008705C11
	case 1241:
		return DeviceInformation_DEV0008705D11
	case 1242:
		return DeviceInformation_DEV0008706211
	case 1243:
		return DeviceInformation_DEV000820BA11
	case 1244:
		return DeviceInformation_DEV000820C011
	case 1245:
		return DeviceInformation_DEV000820B311
	case 1246:
		return DeviceInformation_DEV0008301A11
	case 1247:
		return DeviceInformation_DEV0008C00C13
	case 1248:
		return DeviceInformation_DEV0008302611
	case 1249:
		return DeviceInformation_DEV0008302311
	case 125:
		return DeviceInformation_DEV000C181810
	case 1250:
		return DeviceInformation_DEV0008302011
	case 1251:
		return DeviceInformation_DEV0008C00C11
	case 1252:
		return DeviceInformation_DEV0008302612
	case 1253:
		return DeviceInformation_DEV0008302312
	case 1254:
		return DeviceInformation_DEV0008302012
	case 1255:
		return DeviceInformation_DEV0008C00C15
	case 1256:
		return DeviceInformation_DEV0008C00C14
	case 1257:
		return DeviceInformation_DEV0008B00713
	case 1258:
		return DeviceInformation_DEV0008706611
	case 1259:
		return DeviceInformation_DEV0008706811
	case 126:
		return DeviceInformation_DEV000C20C011
	case 1260:
		return DeviceInformation_DEV0008B00812
	case 1261:
		return DeviceInformation_DEV0008209016
	case 1262:
		return DeviceInformation_DEV0008209011
	case 1263:
		return DeviceInformation_DEV000820A011
	case 1264:
		return DeviceInformation_DEV0008208811
	case 1265:
		return DeviceInformation_DEV0008209811
	case 1266:
		return DeviceInformation_DEV000820CA11
	case 1267:
		return DeviceInformation_DEV0008208012
	case 1268:
		return DeviceInformation_DEV0008207812
	case 1269:
		return DeviceInformation_DEV0008207811
	case 127:
		return DeviceInformation_DEV0079002527
	case 1270:
		return DeviceInformation_DEV0008208011
	case 1271:
		return DeviceInformation_DEV000810D111
	case 1272:
		return DeviceInformation_DEV000810D511
	case 1273:
		return DeviceInformation_DEV000810FA12
	case 1274:
		return DeviceInformation_DEV000810FB12
	case 1275:
		return DeviceInformation_DEV000810F211
	case 1276:
		return DeviceInformation_DEV000810D211
	case 1277:
		return DeviceInformation_DEV000810E211
	case 1278:
		return DeviceInformation_DEV000810D611
	case 1279:
		return DeviceInformation_DEV000810F212
	case 128:
		return DeviceInformation_DEV0079004027
	case 1280:
		return DeviceInformation_DEV000810E212
	case 1281:
		return DeviceInformation_DEV000810FC13
	case 1282:
		return DeviceInformation_DEV000810FD13
	case 1283:
		return DeviceInformation_DEV000810F311
	case 1284:
		return DeviceInformation_DEV000810D311
	case 1285:
		return DeviceInformation_DEV000810D711
	case 1286:
		return DeviceInformation_DEV000810F312
	case 1287:
		return DeviceInformation_DEV000810D811
	case 1288:
		return DeviceInformation_DEV000810E511
	case 1289:
		return DeviceInformation_DEV000810E512
	case 129:
		return DeviceInformation_DEV0079000223
	case 1290:
		return DeviceInformation_DEV000810F611
	case 1291:
		return DeviceInformation_DEV000810D911
	case 1292:
		return DeviceInformation_DEV000810F612
	case 1293:
		return DeviceInformation_DEV000820A812
	case 1294:
		return DeviceInformation_DEV000820AD11
	case 1295:
		return DeviceInformation_DEV000820BD11
	case 1296:
		return DeviceInformation_DEV000820C711
	case 1297:
		return DeviceInformation_DEV000820CD11
	case 1298:
		return DeviceInformation_DEV000820C411
	case 1299:
		return DeviceInformation_DEV000820A811
	case 13:
		return DeviceInformation_DEV0064181610
	case 130:
		return DeviceInformation_DEV0079000123
	case 1300:
		return DeviceInformation_DEV0008501411
	case 1301:
		return DeviceInformation_DEV0008C01602
	case 1302:
		return DeviceInformation_DEV0008302613
	case 1303:
		return DeviceInformation_DEV0008302313
	case 1304:
		return DeviceInformation_DEV0008302013
	case 1305:
		return DeviceInformation_DEV0009207730
	case 1306:
		return DeviceInformation_DEV0009208F10
	case 1307:
		return DeviceInformation_DEV0009C00C13
	case 1308:
		return DeviceInformation_DEV0009209910
	case 1309:
		return DeviceInformation_DEV0009209A10
	case 131:
		return DeviceInformation_DEV0079001427
	case 1310:
		return DeviceInformation_DEV0009207930
	case 1311:
		return DeviceInformation_DEV0009201720
	case 1312:
		return DeviceInformation_DEV0009500D01
	case 1313:
		return DeviceInformation_DEV0009500E01
	case 1314:
		return DeviceInformation_DEV0009209911
	case 1315:
		return DeviceInformation_DEV0009209A11
	case 1316:
		return DeviceInformation_DEV0009C00C12
	case 1317:
		return DeviceInformation_DEV0009C00C11
	case 1318:
		return DeviceInformation_DEV0009500D20
	case 1319:
		return DeviceInformation_DEV0009500E20
	case 132:
		return DeviceInformation_DEV0079003027
	case 1320:
		return DeviceInformation_DEV000920B910
	case 1321:
		return DeviceInformation_DEV0009E0CE10
	case 1322:
		return DeviceInformation_DEV0009E0A210
	case 1323:
		return DeviceInformation_DEV0009501410
	case 1324:
		return DeviceInformation_DEV0009207830
	case 1325:
		return DeviceInformation_DEV0009201620
	case 1326:
		return DeviceInformation_DEV0009E0A111
	case 1327:
		return DeviceInformation_DEV0009E0CD11
	case 1328:
		return DeviceInformation_DEV000920B811
	case 1329:
		return DeviceInformation_DEV000920B611
	case 133:
		return DeviceInformation_DEV0079100C13
	case 1330:
		return DeviceInformation_DEV0009207E10
	case 1331:
		return DeviceInformation_DEV0009207630
	case 1332:
		return DeviceInformation_DEV0009205910
	case 1333:
		return DeviceInformation_DEV0009500B01
	case 1334:
		return DeviceInformation_DEV000920AC10
	case 1335:
		return DeviceInformation_DEV0009207430
	case 1336:
		return DeviceInformation_DEV0009204521
	case 1337:
		return DeviceInformation_DEV0009500A01
	case 1338:
		return DeviceInformation_DEV0009500001
	case 1339:
		return DeviceInformation_DEV000920AB10
	case 134:
		return DeviceInformation_DEV0079101C11
	case 1340:
		return DeviceInformation_DEV000920BF11
	case 1341:
		return DeviceInformation_DEV0009203510
	case 1342:
		return DeviceInformation_DEV0009207A30
	case 1343:
		return DeviceInformation_DEV0009500701
	case 1344:
		return DeviceInformation_DEV0009501710
	case 1345:
		return DeviceInformation_DEV000920B310
	case 1346:
		return DeviceInformation_DEV0009207530
	case 1347:
		return DeviceInformation_DEV0009203321
	case 1348:
		return DeviceInformation_DEV0009500C01
	case 1349:
		return DeviceInformation_DEV000920AD10
	case 135:
		return DeviceInformation_DEV0080709010
	case 1350:
		return DeviceInformation_DEV0009207230
	case 1351:
		return DeviceInformation_DEV0009500801
	case 1352:
		return DeviceInformation_DEV0009501810
	case 1353:
		return DeviceInformation_DEV000920B410
	case 1354:
		return DeviceInformation_DEV0009207330
	case 1355:
		return DeviceInformation_DEV0009204421
	case 1356:
		return DeviceInformation_DEV0009500901
	case 1357:
		return DeviceInformation_DEV000920AA10
	case 1358:
		return DeviceInformation_DEV0009209D01
	case 1359:
		return DeviceInformation_DEV000920B010
	case 136:
		return DeviceInformation_DEV0080707010
	case 1360:
		return DeviceInformation_DEV0009E0BE01
	case 1361:
		return DeviceInformation_DEV000920B110
	case 1362:
		return DeviceInformation_DEV0009E0BD01
	case 1363:
		return DeviceInformation_DEV0009D03F10
	case 1364:
		return DeviceInformation_DEV0009305F10
	case 1365:
		return DeviceInformation_DEV0009305610
	case 1366:
		return DeviceInformation_DEV0009D03E10
	case 1367:
		return DeviceInformation_DEV0009306010
	case 1368:
		return DeviceInformation_DEV0009306110
	case 1369:
		return DeviceInformation_DEV0009306310
	case 137:
		return DeviceInformation_DEV0080706010
	case 1370:
		return DeviceInformation_DEV0009D03B10
	case 1371:
		return DeviceInformation_DEV0009D03C10
	case 1372:
		return DeviceInformation_DEV0009D03910
	case 1373:
		return DeviceInformation_DEV0009D03A10
	case 1374:
		return DeviceInformation_DEV0009305411
	case 1375:
		return DeviceInformation_DEV0009D03D11
	case 1376:
		return DeviceInformation_DEV0009304B11
	case 1377:
		return DeviceInformation_DEV0009304C11
	case 1378:
		return DeviceInformation_DEV0009306220
	case 1379:
		return DeviceInformation_DEV0009302E10
	case 138:
		return DeviceInformation_DEV0080706810
	case 1380:
		return DeviceInformation_DEV0009302F10
	case 1381:
		return DeviceInformation_DEV0009303010
	case 1382:
		return DeviceInformation_DEV0009303110
	case 1383:
		return DeviceInformation_DEV0009306510
	case 1384:
		return DeviceInformation_DEV0009306610
	case 1385:
		return DeviceInformation_DEV0009306410
	case 1386:
		return DeviceInformation_DEV0009401110
	case 1387:
		return DeviceInformation_DEV0009400610
	case 1388:
		return DeviceInformation_DEV0009401510
	case 1389:
		return DeviceInformation_DEV0009402110
	case 139:
		return DeviceInformation_DEV0080705010
	case 1390:
		return DeviceInformation_DEV0009400110
	case 1391:
		return DeviceInformation_DEV0009400910
	case 1392:
		return DeviceInformation_DEV0009400010
	case 1393:
		return DeviceInformation_DEV0009401810
	case 1394:
		return DeviceInformation_DEV0009400310
	case 1395:
		return DeviceInformation_DEV0009301810
	case 1396:
		return DeviceInformation_DEV0009301910
	case 1397:
		return DeviceInformation_DEV0009301A10
	case 1398:
		return DeviceInformation_DEV0009401210
	case 1399:
		return DeviceInformation_DEV0009400810
	case 14:
		return DeviceInformation_DEV006420C011
	case 140:
		return DeviceInformation_DEV0080703013
	case 1400:
		return DeviceInformation_DEV0009400710
	case 1401:
		return DeviceInformation_DEV0009401310
	case 1402:
		return DeviceInformation_DEV0009401410
	case 1403:
		return DeviceInformation_DEV0009402210
	case 1404:
		return DeviceInformation_DEV0009402310
	case 1405:
		return DeviceInformation_DEV0009401710
	case 1406:
		return DeviceInformation_DEV0009401610
	case 1407:
		return DeviceInformation_DEV0009400210
	case 1408:
		return DeviceInformation_DEV0009401010
	case 1409:
		return DeviceInformation_DEV0009400510
	case 141:
		return DeviceInformation_DEV0080704021
	case 1410:
		return DeviceInformation_DEV0009400410
	case 1411:
		return DeviceInformation_DEV0009D04B20
	case 1412:
		return DeviceInformation_DEV0009D04920
	case 1413:
		return DeviceInformation_DEV0009D04A20
	case 1414:
		return DeviceInformation_DEV0009D04820
	case 1415:
		return DeviceInformation_DEV0009D04C11
	case 1416:
		return DeviceInformation_DEV0009D05610
	case 1417:
		return DeviceInformation_DEV0009305510
	case 1418:
		return DeviceInformation_DEV0009209810
	case 1419:
		return DeviceInformation_DEV0009202A10
	case 142:
		return DeviceInformation_DEV0080704022
	case 1420:
		return DeviceInformation_DEV0009209510
	case 1421:
		return DeviceInformation_DEV0009501110
	case 1422:
		return DeviceInformation_DEV0009209310
	case 1423:
		return DeviceInformation_DEV0009209410
	case 1424:
		return DeviceInformation_DEV0009209210
	case 1425:
		return DeviceInformation_DEV0009501210
	case 1426:
		return DeviceInformation_DEV0009205411
	case 1427:
		return DeviceInformation_DEV000920A111
	case 1428:
		return DeviceInformation_DEV000920A311
	case 1429:
		return DeviceInformation_DEV0009205112
	case 143:
		return DeviceInformation_DEV0080704020
	case 1430:
		return DeviceInformation_DEV0009204110
	case 1431:
		return DeviceInformation_DEV0009E07710
	case 1432:
		return DeviceInformation_DEV0009E07712
	case 1433:
		return DeviceInformation_DEV0009205212
	case 1434:
		return DeviceInformation_DEV0009205211
	case 1435:
		return DeviceInformation_DEV0009205311
	case 1436:
		return DeviceInformation_DEV0009206B10
	case 1437:
		return DeviceInformation_DEV0009208010
	case 1438:
		return DeviceInformation_DEV0009206A12
	case 1439:
		return DeviceInformation_DEV0009206810
	case 144:
		return DeviceInformation_DEV0080701111
	case 1440:
		return DeviceInformation_DEV0009208110
	case 1441:
		return DeviceInformation_DEV0009205511
	case 1442:
		return DeviceInformation_DEV0009209F01
	case 1443:
		return DeviceInformation_DEV0009208C10
	case 1444:
		return DeviceInformation_DEV0009208E10
	case 1445:
		return DeviceInformation_DEV000920B511
	case 1446:
		return DeviceInformation_DEV0009501910
	case 1447:
		return DeviceInformation_DEV000920BE11
	case 1448:
		return DeviceInformation_DEV0009209710
	case 1449:
		return DeviceInformation_DEV0009208510
	case 145:
		return DeviceInformation_DEV0080701811
	case 1450:
		return DeviceInformation_DEV0009208610
	case 1451:
		return DeviceInformation_DEV000920BD10
	case 1452:
		return DeviceInformation_DEV0009500210
	case 1453:
		return DeviceInformation_DEV0009500310
	case 1454:
		return DeviceInformation_DEV0009E0BF10
	case 1455:
		return DeviceInformation_DEV0009E0C010
	case 1456:
		return DeviceInformation_DEV0009500110
	case 1457:
		return DeviceInformation_DEV0009209B10
	case 1458:
		return DeviceInformation_DEV0009207D10
	case 1459:
		return DeviceInformation_DEV0009202F11
	case 146:
		return DeviceInformation_DEV008020A110
	case 1460:
		return DeviceInformation_DEV0009203011
	case 1461:
		return DeviceInformation_DEV0009207C10
	case 1462:
		return DeviceInformation_DEV0009207B10
	case 1463:
		return DeviceInformation_DEV0009208710
	case 1464:
		return DeviceInformation_DEV0009E06610
	case 1465:
		return DeviceInformation_DEV0009E06611
	case 1466:
		return DeviceInformation_DEV0009E06410
	case 1467:
		return DeviceInformation_DEV0009E06411
	case 1468:
		return DeviceInformation_DEV0009E06210
	case 1469:
		return DeviceInformation_DEV0009E0E910
	case 147:
		return DeviceInformation_DEV008020A210
	case 1470:
		return DeviceInformation_DEV0009E0EB10
	case 1471:
		return DeviceInformation_DEV000920BB10
	case 1472:
		return DeviceInformation_DEV0009FF1B11
	case 1473:
		return DeviceInformation_DEV0009E0CF10
	case 1474:
		return DeviceInformation_DEV0009206C30
	case 1475:
		return DeviceInformation_DEV0009206D30
	case 1476:
		return DeviceInformation_DEV0009206E30
	case 1477:
		return DeviceInformation_DEV0009206F30
	case 1478:
		return DeviceInformation_DEV0009207130
	case 1479:
		return DeviceInformation_DEV0009204720
	case 148:
		return DeviceInformation_DEV008020A010
	case 1480:
		return DeviceInformation_DEV0009204820
	case 1481:
		return DeviceInformation_DEV0009204920
	case 1482:
		return DeviceInformation_DEV0009204A20
	case 1483:
		return DeviceInformation_DEV0009205A10
	case 1484:
		return DeviceInformation_DEV0009207030
	case 1485:
		return DeviceInformation_DEV0009205B10
	case 1486:
		return DeviceInformation_DEV0009500501
	case 1487:
		return DeviceInformation_DEV0009501001
	case 1488:
		return DeviceInformation_DEV0009500601
	case 1489:
		return DeviceInformation_DEV0009500F01
	case 149:
		return DeviceInformation_DEV0080207212
	case 1490:
		return DeviceInformation_DEV0009500401
	case 1491:
		return DeviceInformation_DEV000920B210
	case 1492:
		return DeviceInformation_DEV000920AE10
	case 1493:
		return DeviceInformation_DEV000920BC10
	case 1494:
		return DeviceInformation_DEV000920AF10
	case 1495:
		return DeviceInformation_DEV0009207F10
	case 1496:
		return DeviceInformation_DEV0009208910
	case 1497:
		return DeviceInformation_DEV0009205710
	case 1498:
		return DeviceInformation_DEV0009205810
	case 1499:
		return DeviceInformation_DEV0009203810
	case 15:
		return DeviceInformation_DEV006420BA11
	case 150:
		return DeviceInformation_DEV0080209111
	case 1500:
		return DeviceInformation_DEV0009203910
	case 1501:
		return DeviceInformation_DEV0009203E10
	case 1502:
		return DeviceInformation_DEV0009204B10
	case 1503:
		return DeviceInformation_DEV0009203F10
	case 1504:
		return DeviceInformation_DEV0009204C10
	case 1505:
		return DeviceInformation_DEV0009204010
	case 1506:
		return DeviceInformation_DEV0009206411
	case 1507:
		return DeviceInformation_DEV0009205E10
	case 1508:
		return DeviceInformation_DEV0009206711
	case 1509:
		return DeviceInformation_DEV000920A710
	case 151:
		return DeviceInformation_DEV0080204310
	case 1510:
		return DeviceInformation_DEV000920A610
	case 1511:
		return DeviceInformation_DEV0009203A10
	case 1512:
		return DeviceInformation_DEV0009203B10
	case 1513:
		return DeviceInformation_DEV0009203C10
	case 1514:
		return DeviceInformation_DEV0009203D10
	case 1515:
		return DeviceInformation_DEV0009E05E12
	case 1516:
		return DeviceInformation_DEV0009E0B711
	case 1517:
		return DeviceInformation_DEV0009E06A12
	case 1518:
		return DeviceInformation_DEV0009E06E12
	case 1519:
		return DeviceInformation_DEV0009E0B720
	case 152:
		return DeviceInformation_DEV008020B612
	case 1520:
		return DeviceInformation_DEV0009E0E611
	case 1521:
		return DeviceInformation_DEV0009E0B321
	case 1522:
		return DeviceInformation_DEV0009E0E512
	case 1523:
		return DeviceInformation_DEV0009204210
	case 1524:
		return DeviceInformation_DEV0009208210
	case 1525:
		return DeviceInformation_DEV0009E07211
	case 1526:
		return DeviceInformation_DEV0009E0CC11
	case 1527:
		return DeviceInformation_DEV0009110111
	case 1528:
		return DeviceInformation_DEV0009110211
	case 1529:
		return DeviceInformation_DEV000916B212
	case 153:
		return DeviceInformation_DEV008020B412
	case 1530:
		return DeviceInformation_DEV0009110212
	case 1531:
		return DeviceInformation_DEV0009110311
	case 1532:
		return DeviceInformation_DEV000916B312
	case 1533:
		return DeviceInformation_DEV0009110312
	case 1534:
		return DeviceInformation_DEV0009110411
	case 1535:
		return DeviceInformation_DEV0009110412
	case 1536:
		return DeviceInformation_DEV0009501615
	case 1537:
		return DeviceInformation_DEV0009E0ED10
	case 1538:
		return DeviceInformation_DEV014F030110
	case 1539:
		return DeviceInformation_DEV014F030310
	case 154:
		return DeviceInformation_DEV008020B512
	case 1540:
		return DeviceInformation_DEV014F030210
	case 1541:
		return DeviceInformation_DEV00EE7FFF10
	case 1542:
		return DeviceInformation_DEV00B6464101
	case 1543:
		return DeviceInformation_DEV00B6464201
	case 1544:
		return DeviceInformation_DEV00B6464501
	case 1545:
		return DeviceInformation_DEV00B6434101
	case 1546:
		return DeviceInformation_DEV00B6434201
	case 1547:
		return DeviceInformation_DEV00B6434202
	case 1548:
		return DeviceInformation_DEV00B6454101
	case 1549:
		return DeviceInformation_DEV00B6454201
	case 155:
		return DeviceInformation_DEV0080208310
	case 1550:
		return DeviceInformation_DEV00B6455001
	case 1551:
		return DeviceInformation_DEV00B6453101
	case 1552:
		return DeviceInformation_DEV00B6453102
	case 1553:
		return DeviceInformation_DEV00B6454102
	case 1554:
		return DeviceInformation_DEV00B6454401
	case 1555:
		return DeviceInformation_DEV00B6454402
	case 1556:
		return DeviceInformation_DEV00B6454202
	case 1557:
		return DeviceInformation_DEV00B6453103
	case 1558:
		return DeviceInformation_DEV00B6453201
	case 1559:
		return DeviceInformation_DEV00B6453301
	case 156:
		return DeviceInformation_DEV0080702111
	case 1560:
		return DeviceInformation_DEV00B6453104
	case 1561:
		return DeviceInformation_DEV00B6454403
	case 1562:
		return DeviceInformation_DEV00B6454801
	case 1563:
		return DeviceInformation_DEV00B6414701
	case 1564:
		return DeviceInformation_DEV00B6414201
	case 1565:
		return DeviceInformation_DEV00B6474101
	case 1566:
		return DeviceInformation_DEV00B6474302
	case 1567:
		return DeviceInformation_DEV00B6474602
	case 1568:
		return DeviceInformation_DEV00B6534D01
	case 1569:
		return DeviceInformation_DEV00B6535001
	case 157:
		return DeviceInformation_DEV0081FE0111
	case 1570:
		return DeviceInformation_DEV00B6455002
	case 1571:
		return DeviceInformation_DEV00B6453701
	case 1572:
		return DeviceInformation_DEV00B6484101
	case 1573:
		return DeviceInformation_DEV00B6484201
	case 1574:
		return DeviceInformation_DEV00B6484202
	case 1575:
		return DeviceInformation_DEV00B6484301
	case 1576:
		return DeviceInformation_DEV00B6484102
	case 1577:
		return DeviceInformation_DEV00B6455101
	case 1578:
		return DeviceInformation_DEV00B6455003
	case 1579:
		return DeviceInformation_DEV00B6455102
	case 158:
		return DeviceInformation_DEV0081FF3131
	case 1580:
		return DeviceInformation_DEV00B6453702
	case 1581:
		return DeviceInformation_DEV00B6453703
	case 1582:
		return DeviceInformation_DEV00B6484302
	case 1583:
		return DeviceInformation_DEV00B6484801
	case 1584:
		return DeviceInformation_DEV00B6484501
	case 1585:
		return DeviceInformation_DEV00B6484203
	case 1586:
		return DeviceInformation_DEV00B6484103
	case 1587:
		return DeviceInformation_DEV00B6455004
	case 1588:
		return DeviceInformation_DEV00B6455103
	case 1589:
		return DeviceInformation_DEV00B6455401
	case 159:
		return DeviceInformation_DEV0081F01313
	case 1590:
		return DeviceInformation_DEV00B6455201
	case 1591:
		return DeviceInformation_DEV00B6455402
	case 1592:
		return DeviceInformation_DEV00B6455403
	case 1593:
		return DeviceInformation_DEV00B603430A
	case 1594:
		return DeviceInformation_DEV00B600010A
	case 1595:
		return DeviceInformation_DEV00B6FF110A
	case 1596:
		return DeviceInformation_DEV00B6434601
	case 1597:
		return DeviceInformation_DEV00B6434602
	case 1598:
		return DeviceInformation_DEV00B6455301
	case 1599:
		return DeviceInformation_DEV00C5070410
	case 16:
		return DeviceInformation_DEV0064182010
	case 160:
		return DeviceInformation_DEV0083002C16
	case 1600:
		return DeviceInformation_DEV00C5070210
	case 1601:
		return DeviceInformation_DEV00C5070610
	case 1602:
		return DeviceInformation_DEV00C5070E11
	case 1603:
		return DeviceInformation_DEV00C5060240
	case 1604:
		return DeviceInformation_DEV00C5062010
	case 1605:
		return DeviceInformation_DEV00C5080230
	case 1606:
		return DeviceInformation_DEV00C5060310
	case 1607:
		return DeviceInformation_DEV006C070E11
	case 1608:
		return DeviceInformation_DEV006C050002
	case 1609:
		return DeviceInformation_DEV006C011311
	case 161:
		return DeviceInformation_DEV0083002E16
	case 1610:
		return DeviceInformation_DEV006C011411
	case 1611:
		return DeviceInformation_DEV0007632010
	case 1612:
		return DeviceInformation_DEV0007632020
	case 1613:
		return DeviceInformation_DEV0007632180
	case 1614:
		return DeviceInformation_DEV0007632040
	case 1615:
		return DeviceInformation_DEV0007613812
	case 1616:
		return DeviceInformation_DEV0007613810
	case 1617:
		return DeviceInformation_DEV000720C011
	case 1618:
		return DeviceInformation_DEV0007A05210
	case 1619:
		return DeviceInformation_DEV0007A08B10
	case 162:
		return DeviceInformation_DEV0083002F16
	case 1620:
		return DeviceInformation_DEV0007A05B32
	case 1621:
		return DeviceInformation_DEV0007A06932
	case 1622:
		return DeviceInformation_DEV0007A06D32
	case 1623:
		return DeviceInformation_DEV0007A08032
	case 1624:
		return DeviceInformation_DEV0007A00213
	case 1625:
		return DeviceInformation_DEV0007A09532
	case 1626:
		return DeviceInformation_DEV0007A06C32
	case 1627:
		return DeviceInformation_DEV0007A05E32
	case 1628:
		return DeviceInformation_DEV0007A08A32
	case 1629:
		return DeviceInformation_DEV0007A07032
	case 163:
		return DeviceInformation_DEV0083012F16
	case 1630:
		return DeviceInformation_DEV0007A08332
	case 1631:
		return DeviceInformation_DEV0007A09832
	case 1632:
		return DeviceInformation_DEV0007A05C32
	case 1633:
		return DeviceInformation_DEV0007A06A32
	case 1634:
		return DeviceInformation_DEV0007A08832
	case 1635:
		return DeviceInformation_DEV0007A06E32
	case 1636:
		return DeviceInformation_DEV0007A08132
	case 1637:
		return DeviceInformation_DEV0007A00113
	case 1638:
		return DeviceInformation_DEV0007A09632
	case 1639:
		return DeviceInformation_DEV0007A05D32
	case 164:
		return DeviceInformation_DEV0083003210
	case 1640:
		return DeviceInformation_DEV0007A06B32
	case 1641:
		return DeviceInformation_DEV0007A08932
	case 1642:
		return DeviceInformation_DEV0007A06F32
	case 1643:
		return DeviceInformation_DEV0007A08232
	case 1644:
		return DeviceInformation_DEV0007A09732
	case 1645:
		return DeviceInformation_DEV0007A05713
	case 1646:
		return DeviceInformation_DEV0007A01811
	case 1647:
		return DeviceInformation_DEV0007A01911
	case 1648:
		return DeviceInformation_DEV0007A04912
	case 1649:
		return DeviceInformation_DEV0007A05814
	case 165:
		return DeviceInformation_DEV0083001D13
	case 1650:
		return DeviceInformation_DEV0007A07114
	case 1651:
		return DeviceInformation_DEV0007A05810
	case 1652:
		return DeviceInformation_DEV0007A04312
	case 1653:
		return DeviceInformation_DEV0007A04412
	case 1654:
		return DeviceInformation_DEV0007A04512
	case 1655:
		return DeviceInformation_DEV000720BD11
	case 1656:
		return DeviceInformation_DEV0007A04C13
	case 1657:
		return DeviceInformation_DEV0007A04D13
	case 1658:
		return DeviceInformation_DEV0007A04B10
	case 1659:
		return DeviceInformation_DEV0007A04E13
	case 166:
		return DeviceInformation_DEV0083001E13
	case 1660:
		return DeviceInformation_DEV0007A04F13
	case 1661:
		return DeviceInformation_DEV000720BA11
	case 1662:
		return DeviceInformation_DEV0007A03D11
	case 1663:
		return DeviceInformation_DEV0007A09211
	case 1664:
		return DeviceInformation_DEV0007A09111
	case 1665:
		return DeviceInformation_DEV0007FF1115
	case 1666:
		return DeviceInformation_DEV0007A01511
	case 1667:
		return DeviceInformation_DEV0007A08411
	case 1668:
		return DeviceInformation_DEV0007A08511
	case 1669:
		return DeviceInformation_DEV0007A03422
	case 167:
		return DeviceInformation_DEV0083001B13
	case 1670:
		return DeviceInformation_DEV0007A07213
	case 1671:
		return DeviceInformation_DEV0007A07420
	case 1672:
		return DeviceInformation_DEV0007A07520
	case 1673:
		return DeviceInformation_DEV0007A07B12
	case 1674:
		return DeviceInformation_DEV0007A07C12
	case 1675:
		return DeviceInformation_DEV0007A09311
	case 1676:
		return DeviceInformation_DEV0007A09013
	case 1677:
		return DeviceInformation_DEV0007A08F13
	case 1678:
		return DeviceInformation_DEV0007A07E10
	case 1679:
		return DeviceInformation_DEV0007A05510
	case 168:
		return DeviceInformation_DEV0083001C13
	case 1680:
		return DeviceInformation_DEV0007A05910
	case 1681:
		return DeviceInformation_DEV0007A08711
	case 1682:
		return DeviceInformation_DEV0007A03D12
	case 1683:
		return DeviceInformation_DEV0007A09A12
	case 1684:
		return DeviceInformation_DEV0007A09B12
	case 1685:
		return DeviceInformation_DEV0007A06614
	case 1686:
		return DeviceInformation_DEV0007A06514
	case 1687:
		return DeviceInformation_DEV0007A06014
	case 1688:
		return DeviceInformation_DEV0007A07714
	case 1689:
		return DeviceInformation_DEV0007A06414
	case 169:
		return DeviceInformation_DEV0083001F11
	case 1690:
		return DeviceInformation_DEV0007A06114
	case 1691:
		return DeviceInformation_DEV0007A07814
	case 1692:
		return DeviceInformation_DEV0007A06714
	case 1693:
		return DeviceInformation_DEV0007A06214
	case 1694:
		return DeviceInformation_DEV0007A07914
	case 1695:
		return DeviceInformation_DEV000B0A8410
	case 1696:
		return DeviceInformation_DEV000B0A7E10
	case 1697:
		return DeviceInformation_DEV000B0A7F10
	case 1698:
		return DeviceInformation_DEV000B0A8010
	case 1699:
		return DeviceInformation_DEV000BBF9111
	case 17:
		return DeviceInformation_DEV0064182510
	case 170:
		return DeviceInformation_DEV0083003C10
	case 1700:
		return DeviceInformation_DEV000B0A7810
	case 1701:
		return DeviceInformation_DEV000B0A7910
	case 1702:
		return DeviceInformation_DEV000B0A7A10
	case 1703:
		return DeviceInformation_DEV000B0A8910
	case 1704:
		return DeviceInformation_DEV000B0A8310
	case 1705:
		return DeviceInformation_DEV000B0A8510
	case 1706:
		return DeviceInformation_DEV000B0A6319
	case 171:
		return DeviceInformation_DEV0083001C20
	case 172:
		return DeviceInformation_DEV0083001B22
	case 173:
		return DeviceInformation_DEV0083003A14
	case 174:
		return DeviceInformation_DEV0083003B14
	case 175:
		return DeviceInformation_DEV0083003B24
	case 176:
		return DeviceInformation_DEV0083003A24
	case 177:
		return DeviceInformation_DEV0083005824
	case 178:
		return DeviceInformation_DEV0083002828
	case 179:
		return DeviceInformation_DEV0083002928
	case 18:
		return DeviceInformation_DEV0064182610
	case 180:
		return DeviceInformation_DEV0083002A18
	case 181:
		return DeviceInformation_DEV0083002B18
	case 182:
		return DeviceInformation_DEV0083002337
	case 183:
		return DeviceInformation_DEV0083002838
	case 184:
		return DeviceInformation_DEV0083002938
	case 185:
		return DeviceInformation_DEV0083002A38
	case 186:
		return DeviceInformation_DEV0083002B38
	case 187:
		return DeviceInformation_DEV0083001321
	case 188:
		return DeviceInformation_DEV0083001421
	case 189:
		return DeviceInformation_DEV0083001521
	case 19:
		return DeviceInformation_DEV0064182910
	case 190:
		return DeviceInformation_DEV0083001621
	case 191:
		return DeviceInformation_DEV0083000921
	case 192:
		return DeviceInformation_DEV0083000D11
	case 193:
		return DeviceInformation_DEV0083000C11
	case 194:
		return DeviceInformation_DEV0083000E11
	case 195:
		return DeviceInformation_DEV0083000B11
	case 196:
		return DeviceInformation_DEV0083000A11
	case 197:
		return DeviceInformation_DEV0083000A21
	case 198:
		return DeviceInformation_DEV0083000B21
	case 199:
		return DeviceInformation_DEV0083000C21
	case 2:
		return DeviceInformation_DEV0001140C13
	case 20:
		return DeviceInformation_DEV0064130610
	case 200:
		return DeviceInformation_DEV0083000D21
	case 201:
		return DeviceInformation_DEV0083000821
	case 202:
		return DeviceInformation_DEV0083000E21
	case 203:
		return DeviceInformation_DEV0083001812
	case 204:
		return DeviceInformation_DEV0083001712
	case 205:
		return DeviceInformation_DEV0083001816
	case 206:
		return DeviceInformation_DEV0083001916
	case 207:
		return DeviceInformation_DEV0083001716
	case 208:
		return DeviceInformation_DEV0083001921
	case 209:
		return DeviceInformation_DEV0083001721
	case 21:
		return DeviceInformation_DEV0064130710
	case 210:
		return DeviceInformation_DEV0083001821
	case 211:
		return DeviceInformation_DEV0083001A20
	case 212:
		return DeviceInformation_DEV0083002320
	case 213:
		return DeviceInformation_DEV0083001010
	case 214:
		return DeviceInformation_DEV0083000F10
	case 215:
		return DeviceInformation_DEV0083003D14
	case 216:
		return DeviceInformation_DEV0083003E14
	case 217:
		return DeviceInformation_DEV0083003F14
	case 218:
		return DeviceInformation_DEV0083004014
	case 219:
		return DeviceInformation_DEV0083004024
	case 22:
		return DeviceInformation_DEV0064133510
	case 220:
		return DeviceInformation_DEV0083003D24
	case 221:
		return DeviceInformation_DEV0083003E24
	case 222:
		return DeviceInformation_DEV0083003F24
	case 223:
		return DeviceInformation_DEV0083001112
	case 224:
		return DeviceInformation_DEV0083001212
	case 225:
		return DeviceInformation_DEV0083005B12
	case 226:
		return DeviceInformation_DEV0083005A12
	case 227:
		return DeviceInformation_DEV0083008410
	case 228:
		return DeviceInformation_DEV0083008510
	case 229:
		return DeviceInformation_DEV0083008610
	case 23:
		return DeviceInformation_DEV0064133310
	case 230:
		return DeviceInformation_DEV0083008710
	case 231:
		return DeviceInformation_DEV0083002515
	case 232:
		return DeviceInformation_DEV0083002115
	case 233:
		return DeviceInformation_DEV0083002015
	case 234:
		return DeviceInformation_DEV0083002415
	case 235:
		return DeviceInformation_DEV0083002615
	case 236:
		return DeviceInformation_DEV0083002215
	case 237:
		return DeviceInformation_DEV0083002715
	case 238:
		return DeviceInformation_DEV0083002315
	case 239:
		return DeviceInformation_DEV0083008B25
	case 24:
		return DeviceInformation_DEV0064133410
	case 240:
		return DeviceInformation_DEV0083008A25
	case 241:
		return DeviceInformation_DEV0083008B28
	case 242:
		return DeviceInformation_DEV0083008A28
	case 243:
		return DeviceInformation_DEV0083009013
	case 244:
		return DeviceInformation_DEV0083009213
	case 245:
		return DeviceInformation_DEV0083009113
	case 246:
		return DeviceInformation_DEV0083009313
	case 247:
		return DeviceInformation_DEV0083009413
	case 248:
		return DeviceInformation_DEV0083009513
	case 249:
		return DeviceInformation_DEV0083009613
	case 25:
		return DeviceInformation_DEV0064133610
	case 250:
		return DeviceInformation_DEV0083009713
	case 251:
		return DeviceInformation_DEV0083009A13
	case 252:
		return DeviceInformation_DEV0083009B13
	case 253:
		return DeviceInformation_DEV0083004B11
	case 254:
		return DeviceInformation_DEV0083004B20
	case 255:
		return DeviceInformation_DEV0083005514
	case 256:
		return DeviceInformation_DEV0083006824
	case 257:
		return DeviceInformation_DEV0083006624
	case 258:
		return DeviceInformation_DEV0083006524
	case 259:
		return DeviceInformation_DEV0083006424
	case 26:
		return DeviceInformation_DEV0064130510
	case 260:
		return DeviceInformation_DEV0083006734
	case 261:
		return DeviceInformation_DEV0083006434
	case 262:
		return DeviceInformation_DEV0083006634
	case 263:
		return DeviceInformation_DEV0083006534
	case 264:
		return DeviceInformation_DEV0083006A34
	case 265:
		return DeviceInformation_DEV0083006B34
	case 266:
		return DeviceInformation_DEV0083006934
	case 267:
		return DeviceInformation_DEV0083004F11
	case 268:
		return DeviceInformation_DEV0083004E10
	case 269:
		return DeviceInformation_DEV0083004D13
	case 27:
		return DeviceInformation_DEV0064480611
	case 270:
		return DeviceInformation_DEV0083004414
	case 271:
		return DeviceInformation_DEV0083004114
	case 272:
		return DeviceInformation_DEV0083004514
	case 273:
		return DeviceInformation_DEV0083004213
	case 274:
		return DeviceInformation_DEV0083004313
	case 275:
		return DeviceInformation_DEV0083004C11
	case 276:
		return DeviceInformation_DEV0083004913
	case 277:
		return DeviceInformation_DEV0083004A13
	case 278:
		return DeviceInformation_DEV0083004712
	case 279:
		return DeviceInformation_DEV0083004610
	case 28:
		return DeviceInformation_DEV0064482011
	case 280:
		return DeviceInformation_DEV0083008E12
	case 281:
		return DeviceInformation_DEV0083004813
	case 282:
		return DeviceInformation_DEV0083005611
	case 283:
		return DeviceInformation_DEV0083005710
	case 284:
		return DeviceInformation_DEV0083005010
	case 285:
		return DeviceInformation_DEV0083001A10
	case 286:
		return DeviceInformation_DEV0083002918
	case 287:
		return DeviceInformation_DEV0083002818
	case 288:
		return DeviceInformation_DEV0083006724
	case 289:
		return DeviceInformation_DEV0083006D41
	case 29:
		return DeviceInformation_DEV0064182210
	case 290:
		return DeviceInformation_DEV0083006E41
	case 291:
		return DeviceInformation_DEV0083007342
	case 292:
		return DeviceInformation_DEV0083007242
	case 293:
		return DeviceInformation_DEV0083006C42
	case 294:
		return DeviceInformation_DEV0083007542
	case 295:
		return DeviceInformation_DEV0083007442
	case 296:
		return DeviceInformation_DEV0083007742
	case 297:
		return DeviceInformation_DEV0083007642
	case 298:
		return DeviceInformation_DEV008300B030
	case 299:
		return DeviceInformation_DEV008300B130
	case 3:
		return DeviceInformation_DEV0001140B11
	case 30:
		return DeviceInformation_DEV0064182710
	case 300:
		return DeviceInformation_DEV008300B230
	case 301:
		return DeviceInformation_DEV008300B330
	case 302:
		return DeviceInformation_DEV008300B430
	case 303:
		return DeviceInformation_DEV008300B530
	case 304:
		return DeviceInformation_DEV008300B630
	case 305:
		return DeviceInformation_DEV008300B730
	case 306:
		return DeviceInformation_DEV0083012843
	case 307:
		return DeviceInformation_DEV0083012943
	case 308:
		return DeviceInformation_DEV008300A421
	case 309:
		return DeviceInformation_DEV008300A521
	case 31:
		return DeviceInformation_DEV0064183010
	case 310:
		return DeviceInformation_DEV008300A621
	case 311:
		return DeviceInformation_DEV0083001332
	case 312:
		return DeviceInformation_DEV0083000932
	case 313:
		return DeviceInformation_DEV0083001432
	case 314:
		return DeviceInformation_DEV0083001532
	case 315:
		return DeviceInformation_DEV0083001632
	case 316:
		return DeviceInformation_DEV008300A432
	case 317:
		return DeviceInformation_DEV008300A532
	case 318:
		return DeviceInformation_DEV008300A632
	case 319:
		return DeviceInformation_DEV0083000F32
	case 32:
		return DeviceInformation_DEV0064B00812
	case 320:
		return DeviceInformation_DEV0083001032
	case 321:
		return DeviceInformation_DEV0083000632
	case 322:
		return DeviceInformation_DEV0083009810
	case 323:
		return DeviceInformation_DEV0083009910
	case 324:
		return DeviceInformation_DEV0083005C11
	case 325:
		return DeviceInformation_DEV0083005D11
	case 326:
		return DeviceInformation_DEV0083005E11
	case 327:
		return DeviceInformation_DEV0083005F11
	case 328:
		return DeviceInformation_DEV0083005413
	case 329:
		return DeviceInformation_DEV0085000520
	case 33:
		return DeviceInformation_DEV0064B00A01
	case 330:
		return DeviceInformation_DEV0085000620
	case 331:
		return DeviceInformation_DEV0085000720
	case 332:
		return DeviceInformation_DEV0085012210
	case 333:
		return DeviceInformation_DEV0085011210
	case 334:
		return DeviceInformation_DEV0085013220
	case 335:
		return DeviceInformation_DEV0085010210
	case 336:
		return DeviceInformation_DEV0085000A10
	case 337:
		return DeviceInformation_DEV0085000B10
	case 338:
		return DeviceInformation_DEV0085071010
	case 339:
		return DeviceInformation_DEV008500FB10
	case 34:
		return DeviceInformation_DEV0064760110
	case 340:
		return DeviceInformation_DEV0085060210
	case 341:
		return DeviceInformation_DEV0085060110
	case 342:
		return DeviceInformation_DEV0085000D20
	case 343:
		return DeviceInformation_DEV008500C810
	case 344:
		return DeviceInformation_DEV0085040111
	case 345:
		return DeviceInformation_DEV008500C910
	case 346:
		return DeviceInformation_DEV0085045020
	case 347:
		return DeviceInformation_DEV0085070210
	case 348:
		return DeviceInformation_DEV0085070110
	case 349:
		return DeviceInformation_DEV0085070310
	case 35:
		return DeviceInformation_DEV0064242313
	case 350:
		return DeviceInformation_DEV0085000E20
	case 351:
		return DeviceInformation_DEV008E596010
	case 352:
		return DeviceInformation_DEV008E593710
	case 353:
		return DeviceInformation_DEV008E597710
	case 354:
		return DeviceInformation_DEV008E598310
	case 355:
		return DeviceInformation_DEV008E598910
	case 356:
		return DeviceInformation_DEV008E593720
	case 357:
		return DeviceInformation_DEV008E598920
	case 358:
		return DeviceInformation_DEV008E598320
	case 359:
		return DeviceInformation_DEV008E596021
	case 36:
		return DeviceInformation_DEV0064FF2111
	case 360:
		return DeviceInformation_DEV008E597721
	case 361:
		return DeviceInformation_DEV008E587320
	case 362:
		return DeviceInformation_DEV008E587020
	case 363:
		return DeviceInformation_DEV008E587220
	case 364:
		return DeviceInformation_DEV008E587120
	case 365:
		return DeviceInformation_DEV008E679910
	case 366:
		return DeviceInformation_DEV008E618310
	case 367:
		return DeviceInformation_DEV008E707910
	case 368:
		return DeviceInformation_DEV008E004010
	case 369:
		return DeviceInformation_DEV008E570910
	case 37:
		return DeviceInformation_DEV0064FF2112
	case 370:
		return DeviceInformation_DEV008E558810
	case 371:
		return DeviceInformation_DEV008E683410
	case 372:
		return DeviceInformation_DEV008E707710
	case 373:
		return DeviceInformation_DEV008E707810
	case 374:
		return DeviceInformation_DEV0091100013
	case 375:
		return DeviceInformation_DEV0091100110
	case 376:
		return DeviceInformation_DEV009E670101
	case 377:
		return DeviceInformation_DEV009E119311
	case 378:
		return DeviceInformation_DEV00A2100C13
	case 379:
		return DeviceInformation_DEV00A2101C11
	case 38:
		return DeviceInformation_DEV0064648B10
	case 380:
		return DeviceInformation_DEV00A2300110
	case 381:
		return DeviceInformation_DEV0002A05814
	case 382:
		return DeviceInformation_DEV0002A07114
	case 383:
		return DeviceInformation_DEV0002134A10
	case 384:
		return DeviceInformation_DEV0002A03422
	case 385:
		return DeviceInformation_DEV0002A03321
	case 386:
		return DeviceInformation_DEV0002648B10
	case 387:
		return DeviceInformation_DEV0002A09013
	case 388:
		return DeviceInformation_DEV0002A08F13
	case 389:
		return DeviceInformation_DEV0002A05510
	case 39:
		return DeviceInformation_DEV0064724010
	case 390:
		return DeviceInformation_DEV0002A05910
	case 391:
		return DeviceInformation_DEV0002A05326
	case 392:
		return DeviceInformation_DEV0002A05428
	case 393:
		return DeviceInformation_DEV0002A08411
	case 394:
		return DeviceInformation_DEV0002A08511
	case 395:
		return DeviceInformation_DEV0002A00F11
	case 396:
		return DeviceInformation_DEV0002A07310
	case 397:
		return DeviceInformation_DEV0002A04110
	case 398:
		return DeviceInformation_DEV0002A03813
	case 399:
		return DeviceInformation_DEV0002A07F13
	case 4:
		return DeviceInformation_DEV0001803002
	case 40:
		return DeviceInformation_DEV006420BD11
	case 400:
		return DeviceInformation_DEV0002A08832
	case 401:
		return DeviceInformation_DEV0002A06E32
	case 402:
		return DeviceInformation_DEV0002A08132
	case 403:
		return DeviceInformation_DEV0002A01D20
	case 404:
		return DeviceInformation_DEV0002A02120
	case 405:
		return DeviceInformation_DEV0002A02520
	case 406:
		return DeviceInformation_DEV0002A02920
	case 407:
		return DeviceInformation_DEV0002A03A20
	case 408:
		return DeviceInformation_DEV0002A05C32
	case 409:
		return DeviceInformation_DEV0002A06A32
	case 41:
		return DeviceInformation_DEV0064570011
	case 410:
		return DeviceInformation_DEV0002A09632
	case 411:
		return DeviceInformation_DEV0002A08932
	case 412:
		return DeviceInformation_DEV0002A06F32
	case 413:
		return DeviceInformation_DEV0002A08232
	case 414:
		return DeviceInformation_DEV0002A01E20
	case 415:
		return DeviceInformation_DEV0002A02220
	case 416:
		return DeviceInformation_DEV0002A02620
	case 417:
		return DeviceInformation_DEV0002A02A20
	case 418:
		return DeviceInformation_DEV0002A03B20
	case 419:
		return DeviceInformation_DEV0002A05D32
	case 42:
		return DeviceInformation_DEV0064570310
	case 420:
		return DeviceInformation_DEV0002A06B32
	case 421:
		return DeviceInformation_DEV0002A09732
	case 422:
		return DeviceInformation_DEV0002A08A32
	case 423:
		return DeviceInformation_DEV0002A07032
	case 424:
		return DeviceInformation_DEV0002A08332
	case 425:
		return DeviceInformation_DEV0002A01F20
	case 426:
		return DeviceInformation_DEV0002A02320
	case 427:
		return DeviceInformation_DEV0002A02720
	case 428:
		return DeviceInformation_DEV0002A02B20
	case 429:
		return DeviceInformation_DEV0002A04820
	case 43:
		return DeviceInformation_DEV0064570211
	case 430:
		return DeviceInformation_DEV0002A06C32
	case 431:
		return DeviceInformation_DEV0002A05E32
	case 432:
		return DeviceInformation_DEV0002A09832
	case 433:
		return DeviceInformation_DEV0002A06D32
	case 434:
		return DeviceInformation_DEV0002A08032
	case 435:
		return DeviceInformation_DEV0002A02020
	case 436:
		return DeviceInformation_DEV0002A02420
	case 437:
		return DeviceInformation_DEV0002A02820
	case 438:
		return DeviceInformation_DEV0002A03920
	case 439:
		return DeviceInformation_DEV0002A05B32
	case 44:
		return DeviceInformation_DEV0064570411
	case 440:
		return DeviceInformation_DEV0002A06932
	case 441:
		return DeviceInformation_DEV0002A09532
	case 442:
		return DeviceInformation_DEV0002B00813
	case 443:
		return DeviceInformation_DEV0002A0A610
	case 444:
		return DeviceInformation_DEV0002A0A611
	case 445:
		return DeviceInformation_DEV0002A0A510
	case 446:
		return DeviceInformation_DEV0002A0A511
	case 447:
		return DeviceInformation_DEV0002A00510
	case 448:
		return DeviceInformation_DEV0002A00610
	case 449:
		return DeviceInformation_DEV0002A01511
	case 45:
		return DeviceInformation_DEV0064570110
	case 450:
		return DeviceInformation_DEV0002A03D11
	case 451:
		return DeviceInformation_DEV000220C011
	case 452:
		return DeviceInformation_DEV0002A05213
	case 453:
		return DeviceInformation_DEV0002A08B10
	case 454:
		return DeviceInformation_DEV0002A0A210
	case 455:
		return DeviceInformation_DEV0002A0A010
	case 456:
		return DeviceInformation_DEV0002A09F10
	case 457:
		return DeviceInformation_DEV0002A0A110
	case 458:
		return DeviceInformation_DEV0002A0A013
	case 459:
		return DeviceInformation_DEV0002A09F13
	case 46:
		return DeviceInformation_DEV0064615022
	case 460:
		return DeviceInformation_DEV0002A0A213
	case 461:
		return DeviceInformation_DEV0002A0A113
	case 462:
		return DeviceInformation_DEV0002A03F11
	case 463:
		return DeviceInformation_DEV0002A04011
	case 464:
		return DeviceInformation_DEV0002FF2112
	case 465:
		return DeviceInformation_DEV0002FF2111
	case 466:
		return DeviceInformation_DEV0002720111
	case 467:
		return DeviceInformation_DEV0002613812
	case 468:
		return DeviceInformation_DEV0002A05713
	case 469:
		return DeviceInformation_DEV0002A07610
	case 47:
		return DeviceInformation_DEV0064182810
	case 470:
		return DeviceInformation_DEV0002A01911
	case 471:
		return DeviceInformation_DEV0002A07611
	case 472:
		return DeviceInformation_DEV0002A04B10
	case 473:
		return DeviceInformation_DEV0002A01B14
	case 474:
		return DeviceInformation_DEV0002A09B11
	case 475:
		return DeviceInformation_DEV0002A09B12
	case 476:
		return DeviceInformation_DEV0002A03C10
	case 477:
		return DeviceInformation_DEV0002A00213
	case 478:
		return DeviceInformation_DEV0002A00113
	case 479:
		return DeviceInformation_DEV0002A02C12
	case 48:
		return DeviceInformation_DEV0064183110
	case 480:
		return DeviceInformation_DEV0002A02D12
	case 481:
		return DeviceInformation_DEV0002A02E12
	case 482:
		return DeviceInformation_DEV0002A04C13
	case 483:
		return DeviceInformation_DEV0002A04D13
	case 484:
		return DeviceInformation_DEV0002A02F12
	case 485:
		return DeviceInformation_DEV0002A03012
	case 486:
		return DeviceInformation_DEV0002A03112
	case 487:
		return DeviceInformation_DEV0002A04E13
	case 488:
		return DeviceInformation_DEV0002A04F13
	case 489:
		return DeviceInformation_DEV0002A01A13
	case 49:
		return DeviceInformation_DEV0064133611
	case 490:
		return DeviceInformation_DEV0002A09C11
	case 491:
		return DeviceInformation_DEV0002A09C12
	case 492:
		return DeviceInformation_DEV0002A01C20
	case 493:
		return DeviceInformation_DEV0002A09A10
	case 494:
		return DeviceInformation_DEV0002A09A12
	case 495:
		return DeviceInformation_DEV000220BA11
	case 496:
		return DeviceInformation_DEV0002A03D12
	case 497:
		return DeviceInformation_DEV0002A09110
	case 498:
		return DeviceInformation_DEV0002A09210
	case 499:
		return DeviceInformation_DEV0002A09111
	case 5:
		return DeviceInformation_DEV00641BD610
	case 50:
		return DeviceInformation_DEV006A000122
	case 500:
		return DeviceInformation_DEV0002A09211
	case 501:
		return DeviceInformation_DEV0002A00E21
	case 502:
		return DeviceInformation_DEV0002A03710
	case 503:
		return DeviceInformation_DEV0002A01112
	case 504:
		return DeviceInformation_DEV0002A01216
	case 505:
		return DeviceInformation_DEV0002A01217
	case 506:
		return DeviceInformation_DEV000220BD11
	case 507:
		return DeviceInformation_DEV0002A07F12
	case 508:
		return DeviceInformation_DEV0002A06613
	case 509:
		return DeviceInformation_DEV0002A06713
	case 51:
		return DeviceInformation_DEV006A000222
	case 510:
		return DeviceInformation_DEV0002A06013
	case 511:
		return DeviceInformation_DEV0002A06113
	case 512:
		return DeviceInformation_DEV0002A06213
	case 513:
		return DeviceInformation_DEV0002A06413
	case 514:
		return DeviceInformation_DEV0002A07713
	case 515:
		return DeviceInformation_DEV0002A07813
	case 516:
		return DeviceInformation_DEV0002A07913
	case 517:
		return DeviceInformation_DEV0002A07914
	case 518:
		return DeviceInformation_DEV0002A06114
	case 519:
		return DeviceInformation_DEV0002A06714
	case 52:
		return DeviceInformation_DEV006A070210
	case 520:
		return DeviceInformation_DEV0002A06414
	case 521:
		return DeviceInformation_DEV0002A06214
	case 522:
		return DeviceInformation_DEV0002A06514
	case 523:
		return DeviceInformation_DEV0002A07714
	case 524:
		return DeviceInformation_DEV0002A06014
	case 525:
		return DeviceInformation_DEV0002A06614
	case 526:
		return DeviceInformation_DEV0002A07814
	case 527:
		return DeviceInformation_DEV0002A0C310
	case 528:
		return DeviceInformation_DEV0002632010
	case 529:
		return DeviceInformation_DEV0002632020
	case 53:
		return DeviceInformation_DEV006BFFF713
	case 530:
		return DeviceInformation_DEV0002632040
	case 531:
		return DeviceInformation_DEV0002632180
	case 532:
		return DeviceInformation_DEV0002632170
	case 533:
		return DeviceInformation_DEV0002FF1140
	case 534:
		return DeviceInformation_DEV0002A07E10
	case 535:
		return DeviceInformation_DEV0002A07213
	case 536:
		return DeviceInformation_DEV0002A04A35
	case 537:
		return DeviceInformation_DEV0002A07420
	case 538:
		return DeviceInformation_DEV0002A07520
	case 539:
		return DeviceInformation_DEV0002A07B12
	case 54:
		return DeviceInformation_DEV006BFF2111
	case 540:
		return DeviceInformation_DEV0002A07C12
	case 541:
		return DeviceInformation_DEV0002A04312
	case 542:
		return DeviceInformation_DEV0002A04412
	case 543:
		return DeviceInformation_DEV0002A04512
	case 544:
		return DeviceInformation_DEV0002A04912
	case 545:
		return DeviceInformation_DEV0002A05012
	case 546:
		return DeviceInformation_DEV0002A01811
	case 547:
		return DeviceInformation_DEV0002A03E11
	case 548:
		return DeviceInformation_DEV0002A08711
	case 549:
		return DeviceInformation_DEV0002A09311
	case 55:
		return DeviceInformation_DEV006BFFF820
	case 550:
		return DeviceInformation_DEV0002A01011
	case 551:
		return DeviceInformation_DEV0002A01622
	case 552:
		return DeviceInformation_DEV0002A04210
	case 553:
		return DeviceInformation_DEV0002A09A13
	case 554:
		return DeviceInformation_DEV00C8272040
	case 555:
		return DeviceInformation_DEV00C8272260
	case 556:
		return DeviceInformation_DEV00C8272060
	case 557:
		return DeviceInformation_DEV00C8272160
	case 558:
		return DeviceInformation_DEV00C8272050
	case 559:
		return DeviceInformation_DEV00C9106D10
	case 56:
		return DeviceInformation_DEV006B106D10
	case 560:
		return DeviceInformation_DEV00C9107C20
	case 561:
		return DeviceInformation_DEV00C9108511
	case 562:
		return DeviceInformation_DEV00C9106210
	case 563:
		return DeviceInformation_DEV00C9109310
	case 564:
		return DeviceInformation_DEV00C9109210
	case 565:
		return DeviceInformation_DEV00C9109810
	case 566:
		return DeviceInformation_DEV00C9109A10
	case 567:
		return DeviceInformation_DEV00C910A420
	case 568:
		return DeviceInformation_DEV00C910A110
	case 569:
		return DeviceInformation_DEV00C910A010
	case 57:
		return DeviceInformation_DEV0071123130
	case 570:
		return DeviceInformation_DEV00C910A310
	case 571:
		return DeviceInformation_DEV00C910A210
	case 572:
		return DeviceInformation_DEV00C9109B10
	case 573:
		return DeviceInformation_DEV00C9106110
	case 574:
		return DeviceInformation_DEV00C9109110
	case 575:
		return DeviceInformation_DEV00C9109610
	case 576:
		return DeviceInformation_DEV00C9109710
	case 577:
		return DeviceInformation_DEV00C9109510
	case 578:
		return DeviceInformation_DEV00C9109910
	case 579:
		return DeviceInformation_DEV00C9109C10
	case 58:
		return DeviceInformation_DEV0071413133
	case 580:
		return DeviceInformation_DEV00C910AB10
	case 581:
		return DeviceInformation_DEV00C910AC10
	case 582:
		return DeviceInformation_DEV00C910AD10
	case 583:
		return DeviceInformation_DEV00C910A810
	case 584:
		return DeviceInformation_DEV00C9106510
	case 585:
		return DeviceInformation_DEV00C910A710
	case 586:
		return DeviceInformation_DEV00C9107610
	case 587:
		return DeviceInformation_DEV00C910890A
	case 588:
		return DeviceInformation_DEV00C9FF1012
	case 589:
		return DeviceInformation_DEV00C9FF0913
	case 59:
		return DeviceInformation_DEV0071114019
	case 590:
		return DeviceInformation_DEV00C9FF1112
	case 591:
		return DeviceInformation_DEV00C9100310
	case 592:
		return DeviceInformation_DEV00C9101110
	case 593:
		return DeviceInformation_DEV00C9101010
	case 594:
		return DeviceInformation_DEV00C9103710
	case 595:
		return DeviceInformation_DEV00C9101310
	case 596:
		return DeviceInformation_DEV00C9FF0D12
	case 597:
		return DeviceInformation_DEV00C9100E10
	case 598:
		return DeviceInformation_DEV00C9100610
	case 599:
		return DeviceInformation_DEV00C9100510
	case 6:
		return DeviceInformation_DEV0064760210
	case 60:
		return DeviceInformation_DEV007111306C
	case 600:
		return DeviceInformation_DEV00C9100710
	case 601:
		return DeviceInformation_DEV00C9FF1D20
	case 602:
		return DeviceInformation_DEV00C9FF1C10
	case 603:
		return DeviceInformation_DEV00C9100810
	case 604:
		return DeviceInformation_DEV00C9FF1420
	case 605:
		return DeviceInformation_DEV00C9100D10
	case 606:
		return DeviceInformation_DEV00C9101220
	case 607:
		return DeviceInformation_DEV00C9102330
	case 608:
		return DeviceInformation_DEV00C9102130
	case 609:
		return DeviceInformation_DEV00C9102430
	case 61:
		return DeviceInformation_DEV0071231112
	case 610:
		return DeviceInformation_DEV00C9100831
	case 611:
		return DeviceInformation_DEV00C9102530
	case 612:
		return DeviceInformation_DEV00C9100531
	case 613:
		return DeviceInformation_DEV00C9102030
	case 614:
		return DeviceInformation_DEV00C9100731
	case 615:
		return DeviceInformation_DEV00C9100631
	case 616:
		return DeviceInformation_DEV00C9102230
	case 617:
		return DeviceInformation_DEV00C9100632
	case 618:
		return DeviceInformation_DEV00C9100532
	case 619:
		return DeviceInformation_DEV00C9100732
	case 62:
		return DeviceInformation_DEV0071113080
	case 620:
		return DeviceInformation_DEV00C9100832
	case 621:
		return DeviceInformation_DEV00C9102532
	case 622:
		return DeviceInformation_DEV00C9102132
	case 623:
		return DeviceInformation_DEV00C9102332
	case 624:
		return DeviceInformation_DEV00C9102432
	case 625:
		return DeviceInformation_DEV00C9102032
	case 626:
		return DeviceInformation_DEV00C9102232
	case 627:
		return DeviceInformation_DEV00C9104432
	case 628:
		return DeviceInformation_DEV00C9100D11
	case 629:
		return DeviceInformation_DEV00C9100633
	case 63:
		return DeviceInformation_DEV0071321212
	case 630:
		return DeviceInformation_DEV00C9100533
	case 631:
		return DeviceInformation_DEV00C9100733
	case 632:
		return DeviceInformation_DEV00C9100833
	case 633:
		return DeviceInformation_DEV00C9102533
	case 634:
		return DeviceInformation_DEV00C9102133
	case 635:
		return DeviceInformation_DEV00C9102333
	case 636:
		return DeviceInformation_DEV00C9102433
	case 637:
		return DeviceInformation_DEV00C9102033
	case 638:
		return DeviceInformation_DEV00C9102233
	case 639:
		return DeviceInformation_DEV00C9104810
	case 64:
		return DeviceInformation_DEV0071321113
	case 640:
		return DeviceInformation_DEV00C9FF1A11
	case 641:
		return DeviceInformation_DEV00C9100212
	case 642:
		return DeviceInformation_DEV00C9FF0A11
	case 643:
		return DeviceInformation_DEV00C9FF0C12
	case 644:
		return DeviceInformation_DEV00C9100112
	case 645:
		return DeviceInformation_DEV00C9FF1911
	case 646:
		return DeviceInformation_DEV00C9FF0B12
	case 647:
		return DeviceInformation_DEV00C9FF0715
	case 648:
		return DeviceInformation_DEV00C9FF1B10
	case 649:
		return DeviceInformation_DEV00C9101610
	case 65:
		return DeviceInformation_DEV0071322212
	case 650:
		return DeviceInformation_DEV00C9FF1B11
	case 651:
		return DeviceInformation_DEV00C9101611
	case 652:
		return DeviceInformation_DEV00C9101612
	case 653:
		return DeviceInformation_DEV00C9FF0F11
	case 654:
		return DeviceInformation_DEV00C9FF1E30
	case 655:
		return DeviceInformation_DEV00C9100410
	case 656:
		return DeviceInformation_DEV00C9106410
	case 657:
		return DeviceInformation_DEV00C9106710
	case 658:
		return DeviceInformation_DEV00C9106810
	case 659:
		return DeviceInformation_DEV00C9106010
	case 66:
		return DeviceInformation_DEV0071322112
	case 660:
		return DeviceInformation_DEV00C9106310
	case 661:
		return DeviceInformation_DEV00C9107110
	case 662:
		return DeviceInformation_DEV00C9107210
	case 663:
		return DeviceInformation_DEV00C9107310
	case 664:
		return DeviceInformation_DEV00C9107010
	case 665:
		return DeviceInformation_DEV00C9107A20
	case 666:
		return DeviceInformation_DEV00C9107B20
	case 667:
		return DeviceInformation_DEV00C9107820
	case 668:
		return DeviceInformation_DEV00C9107920
	case 669:
		return DeviceInformation_DEV00C9104433
	case 67:
		return DeviceInformation_DEV0071322312
	case 670:
		return DeviceInformation_DEV00C9107C11
	case 671:
		return DeviceInformation_DEV00C9107711
	case 672:
		return DeviceInformation_DEV00C9108310
	case 673:
		return DeviceInformation_DEV00C9108210
	case 674:
		return DeviceInformation_DEV00C9108610
	case 675:
		return DeviceInformation_DEV00C9107D10
	case 676:
		return DeviceInformation_DEV00CE648B10
	case 677:
		return DeviceInformation_DEV00CE494513
	case 678:
		return DeviceInformation_DEV00CE494311
	case 679:
		return DeviceInformation_DEV00CE494810
	case 68:
		return DeviceInformation_DEV0071122124
	case 680:
		return DeviceInformation_DEV00CE494712
	case 681:
		return DeviceInformation_DEV00CE494012
	case 682:
		return DeviceInformation_DEV00CE494111
	case 683:
		return DeviceInformation_DEV00CE494210
	case 684:
		return DeviceInformation_DEV00CE494610
	case 685:
		return DeviceInformation_DEV00CE494412
	case 686:
		return DeviceInformation_DEV00D0660212
	case 687:
		return DeviceInformation_DEV00E8000A10
	case 688:
		return DeviceInformation_DEV00E8000B10
	case 689:
		return DeviceInformation_DEV00E8000910
	case 69:
		return DeviceInformation_DEV007112221E
	case 690:
		return DeviceInformation_DEV00E8001112
	case 691:
		return DeviceInformation_DEV00E8000C14
	case 692:
		return DeviceInformation_DEV00E8000D13
	case 693:
		return DeviceInformation_DEV00E8000E12
	case 694:
		return DeviceInformation_DEV00E8001310
	case 695:
		return DeviceInformation_DEV00E8001410
	case 696:
		return DeviceInformation_DEV00E8001510
	case 697:
		return DeviceInformation_DEV00E8000F10
	case 698:
		return DeviceInformation_DEV00E8001010
	case 699:
		return DeviceInformation_DEV00E8000612
	case 7:
		return DeviceInformation_DEV0064182410
	case 70:
		return DeviceInformation_DEV0071413314
	case 700:
		return DeviceInformation_DEV00E8000812
	case 701:
		return DeviceInformation_DEV00E8000712
	case 702:
		return DeviceInformation_DEV00F4501311
	case 703:
		return DeviceInformation_DEV00F4B00911
	case 704:
		return DeviceInformation_DEV0019512710
	case 705:
		return DeviceInformation_DEV0019512810
	case 706:
		return DeviceInformation_DEV0019512910
	case 707:
		return DeviceInformation_DEV0019E30D10
	case 708:
		return DeviceInformation_DEV0019512211
	case 709:
		return DeviceInformation_DEV0019512311
	case 71:
		return DeviceInformation_DEV0072300110
	case 710:
		return DeviceInformation_DEV0019512111
	case 711:
		return DeviceInformation_DEV0019520D11
	case 712:
		return DeviceInformation_DEV0019E30B12
	case 713:
		return DeviceInformation_DEV0019530812
	case 714:
		return DeviceInformation_DEV0019530912
	case 715:
		return DeviceInformation_DEV0019530612
	case 716:
		return DeviceInformation_DEV0019530711
	case 717:
		return DeviceInformation_DEV0019E30A11
	case 718:
		return DeviceInformation_DEV0019E20111
	case 719:
		return DeviceInformation_DEV0019E20210
	case 72:
		return DeviceInformation_DEV0076002101
	case 720:
		return DeviceInformation_DEV0019E30C11
	case 721:
		return DeviceInformation_DEV0019E11310
	case 722:
		return DeviceInformation_DEV0019E11210
	case 723:
		return DeviceInformation_DEV0019E30610
	case 724:
		return DeviceInformation_DEV0019E30710
	case 725:
		return DeviceInformation_DEV0019E30910
	case 726:
		return DeviceInformation_DEV0019E30810
	case 727:
		return DeviceInformation_DEV0019E25510
	case 728:
		return DeviceInformation_DEV0019E20410
	case 729:
		return DeviceInformation_DEV0019E20310
	case 73:
		return DeviceInformation_DEV0076002001
	case 730:
		return DeviceInformation_DEV0019E25610
	case 731:
		return DeviceInformation_DEV0019512010
	case 732:
		return DeviceInformation_DEV0019520C10
	case 733:
		return DeviceInformation_DEV0019520710
	case 734:
		return DeviceInformation_DEV0019520210
	case 735:
		return DeviceInformation_DEV0019E25010
	case 736:
		return DeviceInformation_DEV0019E25110
	case 737:
		return DeviceInformation_DEV0019130710
	case 738:
		return DeviceInformation_DEV0019272050
	case 739:
		return DeviceInformation_DEV0019520910
	case 74:
		return DeviceInformation_DEV0076002A15
	case 740:
		return DeviceInformation_DEV0019520A10
	case 741:
		return DeviceInformation_DEV0019520B10
	case 742:
		return DeviceInformation_DEV0019520412
	case 743:
		return DeviceInformation_DEV0019520812
	case 744:
		return DeviceInformation_DEV0019512510
	case 745:
		return DeviceInformation_DEV0019512410
	case 746:
		return DeviceInformation_DEV0019512610
	case 747:
		return DeviceInformation_DEV0019511711
	case 748:
		return DeviceInformation_DEV0019511811
	case 749:
		return DeviceInformation_DEV0019522212
	case 75:
		return DeviceInformation_DEV0076002815
	case 750:
		return DeviceInformation_DEV0019FF0716
	case 751:
		return DeviceInformation_DEV0019FF1420
	case 752:
		return DeviceInformation_DEV0019522112
	case 753:
		return DeviceInformation_DEV0019522011
	case 754:
		return DeviceInformation_DEV0019522311
	case 755:
		return DeviceInformation_DEV0019E12410
	case 756:
		return DeviceInformation_DEV0019000311
	case 757:
		return DeviceInformation_DEV0019000411
	case 758:
		return DeviceInformation_DEV0019070210
	case 759:
		return DeviceInformation_DEV0019070E11
	case 76:
		return DeviceInformation_DEV0076002215
	case 760:
		return DeviceInformation_DEV0019724010
	case 761:
		return DeviceInformation_DEV0019520610
	case 762:
		return DeviceInformation_DEV0019520510
	case 763:
		return DeviceInformation_DEV00FB101111
	case 764:
		return DeviceInformation_DEV00FB103001
	case 765:
		return DeviceInformation_DEV00FB104401
	case 766:
		return DeviceInformation_DEV00FB124002
	case 767:
		return DeviceInformation_DEV00FB104102
	case 768:
		return DeviceInformation_DEV00FB104201
	case 769:
		return DeviceInformation_DEV00FBF77603
	case 77:
		return DeviceInformation_DEV0076002B15
	case 770:
		return DeviceInformation_DEV00FB104301
	case 771:
		return DeviceInformation_DEV00FB104601
	case 772:
		return DeviceInformation_DEV00FB104701
	case 773:
		return DeviceInformation_DEV00FB105101
	case 774:
		return DeviceInformation_DEV0103030110
	case 775:
		return DeviceInformation_DEV0103010113
	case 776:
		return DeviceInformation_DEV0103090110
	case 777:
		return DeviceInformation_DEV0103020111
	case 778:
		return DeviceInformation_DEV0103020112
	case 779:
		return DeviceInformation_DEV0103040110
	case 78:
		return DeviceInformation_DEV0076002715
	case 780:
		return DeviceInformation_DEV0103050111
	case 781:
		return DeviceInformation_DEV0107000301
	case 782:
		return DeviceInformation_DEV0107000101
	case 783:
		return DeviceInformation_DEV0107000201
	case 784:
		return DeviceInformation_DEV0107020801
	case 785:
		return DeviceInformation_DEV0107020401
	case 786:
		return DeviceInformation_DEV0107020001
	case 787:
		return DeviceInformation_DEV010701F801
	case 788:
		return DeviceInformation_DEV010701FC01
	case 789:
		return DeviceInformation_DEV0107020C01
	case 79:
		return DeviceInformation_DEV0076002315
	case 790:
		return DeviceInformation_DEV010F100801
	case 791:
		return DeviceInformation_DEV010F100601
	case 792:
		return DeviceInformation_DEV010F100401
	case 793:
		return DeviceInformation_DEV010F030601
	case 794:
		return DeviceInformation_DEV010F010301
	case 795:
		return DeviceInformation_DEV010F010101
	case 796:
		return DeviceInformation_DEV010F010201
	case 797:
		return DeviceInformation_DEV010F000302
	case 798:
		return DeviceInformation_DEV010F000402
	case 799:
		return DeviceInformation_DEV010F000102
	case 8:
		return DeviceInformation_DEV0064182310
	case 80:
		return DeviceInformation_DEV0076002415
	case 800:
		return DeviceInformation_DEV011EBB8211
	case 801:
		return DeviceInformation_DEV011E108111
	case 802:
		return DeviceInformation_DEV0123010010
	case 803:
		return DeviceInformation_DEV001E478010
	case 804:
		return DeviceInformation_DEV001E706611
	case 805:
		return DeviceInformation_DEV001E706811
	case 806:
		return DeviceInformation_DEV001E473012
	case 807:
		return DeviceInformation_DEV001E20A011
	case 808:
		return DeviceInformation_DEV001E209011
	case 809:
		return DeviceInformation_DEV001E209811
	case 81:
		return DeviceInformation_DEV0076002615
	case 810:
		return DeviceInformation_DEV001E208811
	case 811:
		return DeviceInformation_DEV001E208011
	case 812:
		return DeviceInformation_DEV001E207821
	case 813:
		return DeviceInformation_DEV001E20CA12
	case 814:
		return DeviceInformation_DEV001E20B312
	case 815:
		return DeviceInformation_DEV001E20B012
	case 816:
		return DeviceInformation_DEV001E302612
	case 817:
		return DeviceInformation_DEV001E302312
	case 818:
		return DeviceInformation_DEV001E302012
	case 819:
		return DeviceInformation_DEV001E20A811
	case 82:
		return DeviceInformation_DEV0076002515
	case 820:
		return DeviceInformation_DEV001E20C412
	case 821:
		return DeviceInformation_DEV001E20C712
	case 822:
		return DeviceInformation_DEV001E20AD12
	case 823:
		return DeviceInformation_DEV001E443720
	case 824:
		return DeviceInformation_DEV001E441821
	case 825:
		return DeviceInformation_DEV001E443810
	case 826:
		return DeviceInformation_DEV001E140C12
	case 827:
		return DeviceInformation_DEV001E471611
	case 828:
		return DeviceInformation_DEV001E479024
	case 829:
		return DeviceInformation_DEV001E471A11
	case 83:
		return DeviceInformation_DEV0076000201
	case 830:
		return DeviceInformation_DEV001E477A10
	case 831:
		return DeviceInformation_DEV001E470A11
	case 832:
		return DeviceInformation_DEV001E480B11
	case 833:
		return DeviceInformation_DEV001E487B10
	case 834:
		return DeviceInformation_DEV001E440411
	case 835:
		return DeviceInformation_DEV001E447211
	case 836:
		return DeviceInformation_DEV0142010010
	case 837:
		return DeviceInformation_DEV0142010011
	case 838:
		return DeviceInformation_DEV017A130401
	case 839:
		return DeviceInformation_DEV017A130201
	case 84:
		return DeviceInformation_DEV0076000101
	case 840:
		return DeviceInformation_DEV017A130801
	case 841:
		return DeviceInformation_DEV017A130601
	case 842:
		return DeviceInformation_DEV017A300102
	case 843:
		return DeviceInformation_DEV0193323C11
	case 844:
		return DeviceInformation_DEV0198101110
	case 845:
		return DeviceInformation_DEV01C4030110
	case 846:
		return DeviceInformation_DEV01C4030210
	case 847:
		return DeviceInformation_DEV01C4021210
	case 848:
		return DeviceInformation_DEV01C4001010
	case 849:
		return DeviceInformation_DEV01C4020610
	case 85:
		return DeviceInformation_DEV0076000301
	case 850:
		return DeviceInformation_DEV01C4020910
	case 851:
		return DeviceInformation_DEV01C4020810
	case 852:
		return DeviceInformation_DEV01C4010710
	case 853:
		return DeviceInformation_DEV01C4050210
	case 854:
		return DeviceInformation_DEV01C4010810
	case 855:
		return DeviceInformation_DEV01C4020510
	case 856:
		return DeviceInformation_DEV01C4040110
	case 857:
		return DeviceInformation_DEV01C4040310
	case 858:
		return DeviceInformation_DEV01C4040210
	case 859:
		return DeviceInformation_DEV01C4101210
	case 86:
		return DeviceInformation_DEV0076000401
	case 860:
		return DeviceInformation_DEV003D020109
	case 861:
		return DeviceInformation_DEV01DB000301
	case 862:
		return DeviceInformation_DEV01DB000201
	case 863:
		return DeviceInformation_DEV01DB000401
	case 864:
		return DeviceInformation_DEV01DB000801
	case 865:
		return DeviceInformation_DEV01DB001201
	case 866:
		return DeviceInformation_DEV009A000400
	case 867:
		return DeviceInformation_DEV009A100400
	case 868:
		return DeviceInformation_DEV009A200C00
	case 869:
		return DeviceInformation_DEV009A200E00
	case 87:
		return DeviceInformation_DEV0076002903
	case 870:
		return DeviceInformation_DEV009A000201
	case 871:
		return DeviceInformation_DEV009A000300
	case 872:
		return DeviceInformation_DEV009A00C000
	case 873:
		return DeviceInformation_DEV009A00B000
	case 874:
		return DeviceInformation_DEV009A00C002
	case 875:
		return DeviceInformation_DEV009A200100
	case 876:
		return DeviceInformation_DEV004E400010
	case 877:
		return DeviceInformation_DEV004E030031
	case 878:
		return DeviceInformation_DEV012B010110
	case 879:
		return DeviceInformation_DEV01F6E0E110
	case 88:
		return DeviceInformation_DEV0076002901
	case 880:
		return DeviceInformation_DEV0088100010
	case 881:
		return DeviceInformation_DEV0088100210
	case 882:
		return DeviceInformation_DEV0088100110
	case 883:
		return DeviceInformation_DEV0088110010
	case 884:
		return DeviceInformation_DEV0088120412
	case 885:
		return DeviceInformation_DEV0088120113
	case 886:
		return DeviceInformation_DEV011A4B5201
	case 887:
		return DeviceInformation_DEV008B020301
	case 888:
		return DeviceInformation_DEV008B010610
	case 889:
		return DeviceInformation_DEV008B030110
	case 89:
		return DeviceInformation_DEV007600E503
	case 890:
		return DeviceInformation_DEV008B030310
	case 891:
		return DeviceInformation_DEV008B030210
	case 892:
		return DeviceInformation_DEV008B031512
	case 893:
		return DeviceInformation_DEV008B031412
	case 894:
		return DeviceInformation_DEV008B031312
	case 895:
		return DeviceInformation_DEV008B031212
	case 896:
		return DeviceInformation_DEV008B031112
	case 897:
		return DeviceInformation_DEV008B031012
	case 898:
		return DeviceInformation_DEV008B030510
	case 899:
		return DeviceInformation_DEV008B030410
	case 9:
		return DeviceInformation_DEV0064705C01
	case 90:
		return DeviceInformation_DEV0076004002
	case 900:
		return DeviceInformation_DEV008B020310
	case 901:
		return DeviceInformation_DEV008B020210
	case 902:
		return DeviceInformation_DEV008B020110
	case 903:
		return DeviceInformation_DEV008B010110
	case 904:
		return DeviceInformation_DEV008B010210
	case 905:
		return DeviceInformation_DEV008B010310
	case 906:
		return DeviceInformation_DEV008B010410
	case 907:
		return DeviceInformation_DEV008B040110
	case 908:
		return DeviceInformation_DEV008B040210
	case 909:
		return DeviceInformation_DEV008B010910
	case 91:
		return DeviceInformation_DEV0076004003
	case 910:
		return DeviceInformation_DEV008B010710
	case 911:
		return DeviceInformation_DEV008B010810
	case 912:
		return DeviceInformation_DEV008B041111
	case 913:
		return DeviceInformation_DEV008B041211
	case 914:
		return DeviceInformation_DEV008B041311
	case 915:
		return DeviceInformation_DEV00A600020A
	case 916:
		return DeviceInformation_DEV00A6000B10
	case 917:
		return DeviceInformation_DEV00A6000300
	case 918:
		return DeviceInformation_DEV00A6000705
	case 919:
		return DeviceInformation_DEV00A6000605
	case 92:
		return DeviceInformation_DEV0076003402
	case 920:
		return DeviceInformation_DEV00A6000500
	case 921:
		return DeviceInformation_DEV00A6000C10
	case 922:
		return DeviceInformation_DEV002CA01811
	case 923:
		return DeviceInformation_DEV002CA07033
	case 924:
		return DeviceInformation_DEV002C555020
	case 925:
		return DeviceInformation_DEV002C556421
	case 926:
		return DeviceInformation_DEV002C05F211
	case 927:
		return DeviceInformation_DEV002C05F411
	case 928:
		return DeviceInformation_DEV002C05E613
	case 929:
		return DeviceInformation_DEV002CA07914
	case 93:
		return DeviceInformation_DEV0076003401
	case 930:
		return DeviceInformation_DEV002C060A13
	case 931:
		return DeviceInformation_DEV002C3A0212
	case 932:
		return DeviceInformation_DEV002C060210
	case 933:
		return DeviceInformation_DEV002CA00213
	case 934:
		return DeviceInformation_DEV002CA0A611
	case 935:
		return DeviceInformation_DEV002CA07B11
	case 936:
		return DeviceInformation_DEV002C063210
	case 937:
		return DeviceInformation_DEV002C063110
	case 938:
		return DeviceInformation_DEV002C062E10
	case 939:
		return DeviceInformation_DEV002C062C10
	case 94:
		return DeviceInformation_DEV007600E908
	case 940:
		return DeviceInformation_DEV002C062D10
	case 941:
		return DeviceInformation_DEV002C063310
	case 942:
		return DeviceInformation_DEV002C05EB10
	case 943:
		return DeviceInformation_DEV002C05F110
	case 944:
		return DeviceInformation_DEV002C0B8830
	case 945:
		return DeviceInformation_DEV00A0B07101
	case 946:
		return DeviceInformation_DEV00A0B07001
	case 947:
		return DeviceInformation_DEV00A0B07203
	case 948:
		return DeviceInformation_DEV00A0B02101
	case 949:
		return DeviceInformation_DEV00A0B02401
	case 95:
		return DeviceInformation_DEV007600E907
	case 950:
		return DeviceInformation_DEV00A0B02301
	case 951:
		return DeviceInformation_DEV00A0B02601
	case 952:
		return DeviceInformation_DEV00A0B02201
	case 953:
		return DeviceInformation_DEV00A0B01902
	case 954:
		return DeviceInformation_DEV0004147112
	case 955:
		return DeviceInformation_DEV000410A411
	case 956:
		return DeviceInformation_DEV0004109911
	case 957:
		return DeviceInformation_DEV0004109912
	case 958:
		return DeviceInformation_DEV0004109913
	case 959:
		return DeviceInformation_DEV0004109914
	case 96:
		return DeviceInformation_DEV000C181710
	case 960:
		return DeviceInformation_DEV000410A211
	case 961:
		return DeviceInformation_DEV000410FC12
	case 962:
		return DeviceInformation_DEV000410FD12
	case 963:
		return DeviceInformation_DEV000410B212
	case 964:
		return DeviceInformation_DEV0004110B11
	case 965:
		return DeviceInformation_DEV0004110711
	case 966:
		return DeviceInformation_DEV000410B213
	case 967:
		return DeviceInformation_DEV0004109811
	case 968:
		return DeviceInformation_DEV0004109812
	case 969:
		return DeviceInformation_DEV0004109813
	case 97:
		return DeviceInformation_DEV000C130510
	case 970:
		return DeviceInformation_DEV0004109814
	case 971:
		return DeviceInformation_DEV000410A011
	case 972:
		return DeviceInformation_DEV000410A111
	case 973:
		return DeviceInformation_DEV000410FA12
	case 974:
		return DeviceInformation_DEV000410FB12
	case 975:
		return DeviceInformation_DEV000410B112
	case 976:
		return DeviceInformation_DEV0004110A11
	case 977:
		return DeviceInformation_DEV0004110611
	case 978:
		return DeviceInformation_DEV000410B113
	case 979:
		return DeviceInformation_DEV0004109A11
	case 98:
		return DeviceInformation_DEV000C130610
	case 980:
		return DeviceInformation_DEV0004109A12
	case 981:
		return DeviceInformation_DEV0004109A13
	case 982:
		return DeviceInformation_DEV0004109A14
	case 983:
		return DeviceInformation_DEV000410A311
	case 984:
		return DeviceInformation_DEV000410B312
	case 985:
		return DeviceInformation_DEV0004110C11
	case 986:
		return DeviceInformation_DEV0004110811
	case 987:
		return DeviceInformation_DEV000410B313
	case 988:
		return DeviceInformation_DEV0004109B11
	case 989:
		return DeviceInformation_DEV0004109B12
	case 99:
		return DeviceInformation_DEV000C133610
	case 990:
		return DeviceInformation_DEV0004109B13
	case 991:
		return DeviceInformation_DEV0004109B14
	case 992:
		return DeviceInformation_DEV000410A511
	case 993:
		return DeviceInformation_DEV000410B412
	case 994:
		return DeviceInformation_DEV0004110D11
	case 995:
		return DeviceInformation_DEV0004110911
	case 996:
		return DeviceInformation_DEV000410B413
	case 997:
		return DeviceInformation_DEV0004109C11
	case 998:
		return DeviceInformation_DEV0004109C12
	case 999:
		return DeviceInformation_DEV0004109C13
	}
	return 0
}

func DeviceInformationByName(value string) DeviceInformation {
	switch value {
	case "DEV0001914201":
		return DeviceInformation_DEV0001914201
	case "DEV0064181910":
		return DeviceInformation_DEV0064181910
	case "DEV000C133410":
		return DeviceInformation_DEV000C133410
	case "DEV0004109C14":
		return DeviceInformation_DEV0004109C14
	case "DEV000410A611":
		return DeviceInformation_DEV000410A611
	case "DEV0004146B13":
		return DeviceInformation_DEV0004146B13
	case "DEV0004147211":
		return DeviceInformation_DEV0004147211
	case "DEV000410FE12":
		return DeviceInformation_DEV000410FE12
	case "DEV0004209016":
		return DeviceInformation_DEV0004209016
	case "DEV000420A011":
		return DeviceInformation_DEV000420A011
	case "DEV0004209011":
		return DeviceInformation_DEV0004209011
	case "DEV000420CA11":
		return DeviceInformation_DEV000420CA11
	case "DEV0004208012":
		return DeviceInformation_DEV0004208012
	case "DEV000C133310":
		return DeviceInformation_DEV000C133310
	case "DEV0004207812":
		return DeviceInformation_DEV0004207812
	case "DEV000420BA11":
		return DeviceInformation_DEV000420BA11
	case "DEV000420B311":
		return DeviceInformation_DEV000420B311
	case "DEV0004209811":
		return DeviceInformation_DEV0004209811
	case "DEV0004208811":
		return DeviceInformation_DEV0004208811
	case "DEV0004B00812":
		return DeviceInformation_DEV0004B00812
	case "DEV0004302613":
		return DeviceInformation_DEV0004302613
	case "DEV0004302313":
		return DeviceInformation_DEV0004302313
	case "DEV0004302013":
		return DeviceInformation_DEV0004302013
	case "DEV0004302B12":
		return DeviceInformation_DEV0004302B12
	case "DEV000C133611":
		return DeviceInformation_DEV000C133611
	case "DEV0004706811":
		return DeviceInformation_DEV0004706811
	case "DEV0004705D11":
		return DeviceInformation_DEV0004705D11
	case "DEV0004705C11":
		return DeviceInformation_DEV0004705C11
	case "DEV0004B00713":
		return DeviceInformation_DEV0004B00713
	case "DEV0004B00A01":
		return DeviceInformation_DEV0004B00A01
	case "DEV0004706611":
		return DeviceInformation_DEV0004706611
	case "DEV0004C01410":
		return DeviceInformation_DEV0004C01410
	case "DEV0004C00102":
		return DeviceInformation_DEV0004C00102
	case "DEV0004705E11":
		return DeviceInformation_DEV0004705E11
	case "DEV0004706211":
		return DeviceInformation_DEV0004706211
	case "DEV000C133510":
		return DeviceInformation_DEV000C133510
	case "DEV0004706411":
		return DeviceInformation_DEV0004706411
	case "DEV0004706412":
		return DeviceInformation_DEV0004706412
	case "DEV000420C011":
		return DeviceInformation_DEV000420C011
	case "DEV000420B011":
		return DeviceInformation_DEV000420B011
	case "DEV0004B00911":
		return DeviceInformation_DEV0004B00911
	case "DEV0004A01211":
		return DeviceInformation_DEV0004A01211
	case "DEV0004A01112":
		return DeviceInformation_DEV0004A01112
	case "DEV0004A01111":
		return DeviceInformation_DEV0004A01111
	case "DEV0004A01212":
		return DeviceInformation_DEV0004A01212
	case "DEV0004A03312":
		return DeviceInformation_DEV0004A03312
	case "DEV000C130710":
		return DeviceInformation_DEV000C130710
	case "DEV0004A03212":
		return DeviceInformation_DEV0004A03212
	case "DEV0004A01113":
		return DeviceInformation_DEV0004A01113
	case "DEV0004A01711":
		return DeviceInformation_DEV0004A01711
	case "DEV000420C711":
		return DeviceInformation_DEV000420C711
	case "DEV000420BD11":
		return DeviceInformation_DEV000420BD11
	case "DEV000420C411":
		return DeviceInformation_DEV000420C411
	case "DEV000420A812":
		return DeviceInformation_DEV000420A812
	case "DEV000420CD11":
		return DeviceInformation_DEV000420CD11
	case "DEV000420AD11":
		return DeviceInformation_DEV000420AD11
	case "DEV000420B611":
		return DeviceInformation_DEV000420B611
	case "DEV000C760210":
		return DeviceInformation_DEV000C760210
	case "DEV000420A811":
		return DeviceInformation_DEV000420A811
	case "DEV0004501311":
		return DeviceInformation_DEV0004501311
	case "DEV0004501411":
		return DeviceInformation_DEV0004501411
	case "DEV0004B01002":
		return DeviceInformation_DEV0004B01002
	case "DEV0006D00610":
		return DeviceInformation_DEV0006D00610
	case "DEV0006D01510":
		return DeviceInformation_DEV0006D01510
	case "DEV0006D00110":
		return DeviceInformation_DEV0006D00110
	case "DEV0006D00310":
		return DeviceInformation_DEV0006D00310
	case "DEV0006D03210":
		return DeviceInformation_DEV0006D03210
	case "DEV0006D03310":
		return DeviceInformation_DEV0006D03310
	case "DEV000C1BD610":
		return DeviceInformation_DEV000C1BD610
	case "DEV0006D02E20":
		return DeviceInformation_DEV0006D02E20
	case "DEV0006D02F20":
		return DeviceInformation_DEV0006D02F20
	case "DEV0006D03020":
		return DeviceInformation_DEV0006D03020
	case "DEV0006D03120":
		return DeviceInformation_DEV0006D03120
	case "DEV0006D02110":
		return DeviceInformation_DEV0006D02110
	case "DEV0006D00010":
		return DeviceInformation_DEV0006D00010
	case "DEV0006D01810":
		return DeviceInformation_DEV0006D01810
	case "DEV0006D00910":
		return DeviceInformation_DEV0006D00910
	case "DEV0006D01110":
		return DeviceInformation_DEV0006D01110
	case "DEV0006D03510":
		return DeviceInformation_DEV0006D03510
	case "DEV000C181610":
		return DeviceInformation_DEV000C181610
	case "DEV0006D03410":
		return DeviceInformation_DEV0006D03410
	case "DEV0006D02410":
		return DeviceInformation_DEV0006D02410
	case "DEV0006D02510":
		return DeviceInformation_DEV0006D02510
	case "DEV0006D00810":
		return DeviceInformation_DEV0006D00810
	case "DEV0006D00710":
		return DeviceInformation_DEV0006D00710
	case "DEV0006D01310":
		return DeviceInformation_DEV0006D01310
	case "DEV0006D01410":
		return DeviceInformation_DEV0006D01410
	case "DEV0006D00210":
		return DeviceInformation_DEV0006D00210
	case "DEV0006D00510":
		return DeviceInformation_DEV0006D00510
	case "DEV0006D00410":
		return DeviceInformation_DEV0006D00410
	case "DEV000C648B10":
		return DeviceInformation_DEV000C648B10
	case "DEV0006D02210":
		return DeviceInformation_DEV0006D02210
	case "DEV0006D02310":
		return DeviceInformation_DEV0006D02310
	case "DEV0006D01710":
		return DeviceInformation_DEV0006D01710
	case "DEV0006D01610":
		return DeviceInformation_DEV0006D01610
	case "DEV0006D01010":
		return DeviceInformation_DEV0006D01010
	case "DEV0006D01210":
		return DeviceInformation_DEV0006D01210
	case "DEV0006D04820":
		return DeviceInformation_DEV0006D04820
	case "DEV0006D04C11":
		return DeviceInformation_DEV0006D04C11
	case "DEV0006D05610":
		return DeviceInformation_DEV0006D05610
	case "DEV0006D02910":
		return DeviceInformation_DEV0006D02910
	case "DEV000C480611":
		return DeviceInformation_DEV000C480611
	case "DEV0006D02A10":
		return DeviceInformation_DEV0006D02A10
	case "DEV0006D02B10":
		return DeviceInformation_DEV0006D02B10
	case "DEV0006D02C10":
		return DeviceInformation_DEV0006D02C10
	case "DEV0006D02810":
		return DeviceInformation_DEV0006D02810
	case "DEV0006D02610":
		return DeviceInformation_DEV0006D02610
	case "DEV0006D02710":
		return DeviceInformation_DEV0006D02710
	case "DEV0006D03610":
		return DeviceInformation_DEV0006D03610
	case "DEV0006D03710":
		return DeviceInformation_DEV0006D03710
	case "DEV0006D02D11":
		return DeviceInformation_DEV0006D02D11
	case "DEV0006D03C10":
		return DeviceInformation_DEV0006D03C10
	case "DEV0064181810":
		return DeviceInformation_DEV0064181810
	case "DEV000C482011":
		return DeviceInformation_DEV000C482011
	case "DEV0006D03B10":
		return DeviceInformation_DEV0006D03B10
	case "DEV0006D03910":
		return DeviceInformation_DEV0006D03910
	case "DEV0006D03A10":
		return DeviceInformation_DEV0006D03A10
	case "DEV0006D03D11":
		return DeviceInformation_DEV0006D03D11
	case "DEV0006D03E10":
		return DeviceInformation_DEV0006D03E10
	case "DEV0006C00102":
		return DeviceInformation_DEV0006C00102
	case "DEV0006E05611":
		return DeviceInformation_DEV0006E05611
	case "DEV0006E05212":
		return DeviceInformation_DEV0006E05212
	case "DEV000620B011":
		return DeviceInformation_DEV000620B011
	case "DEV000620B311":
		return DeviceInformation_DEV000620B311
	case "DEV000C724010":
		return DeviceInformation_DEV000C724010
	case "DEV000620C011":
		return DeviceInformation_DEV000620C011
	case "DEV000620BA11":
		return DeviceInformation_DEV000620BA11
	case "DEV0006705C11":
		return DeviceInformation_DEV0006705C11
	case "DEV0006705D11":
		return DeviceInformation_DEV0006705D11
	case "DEV0006E07710":
		return DeviceInformation_DEV0006E07710
	case "DEV0006E07712":
		return DeviceInformation_DEV0006E07712
	case "DEV0006706210":
		return DeviceInformation_DEV0006706210
	case "DEV0006302611":
		return DeviceInformation_DEV0006302611
	case "DEV0006302612":
		return DeviceInformation_DEV0006302612
	case "DEV0006E00810":
		return DeviceInformation_DEV0006E00810
	case "DEV000C570211":
		return DeviceInformation_DEV000C570211
	case "DEV0006E01F01":
		return DeviceInformation_DEV0006E01F01
	case "DEV0006302311":
		return DeviceInformation_DEV0006302311
	case "DEV0006302312":
		return DeviceInformation_DEV0006302312
	case "DEV0006E00910":
		return DeviceInformation_DEV0006E00910
	case "DEV0006E02001":
		return DeviceInformation_DEV0006E02001
	case "DEV0006302011":
		return DeviceInformation_DEV0006302011
	case "DEV0006302012":
		return DeviceInformation_DEV0006302012
	case "DEV0006C00C13":
		return DeviceInformation_DEV0006C00C13
	case "DEV0006E00811":
		return DeviceInformation_DEV0006E00811
	case "DEV0006E00911":
		return DeviceInformation_DEV0006E00911
	case "DEV000C570310":
		return DeviceInformation_DEV000C570310
	case "DEV0006E01F20":
		return DeviceInformation_DEV0006E01F20
	case "DEV0006E03410":
		return DeviceInformation_DEV0006E03410
	case "DEV0006E03110":
		return DeviceInformation_DEV0006E03110
	case "DEV0006E0A210":
		return DeviceInformation_DEV0006E0A210
	case "DEV0006E0CE10":
		return DeviceInformation_DEV0006E0CE10
	case "DEV0006E0A111":
		return DeviceInformation_DEV0006E0A111
	case "DEV0006E0CD11":
		return DeviceInformation_DEV0006E0CD11
	case "DEV0006E02020":
		return DeviceInformation_DEV0006E02020
	case "DEV0006E02D11":
		return DeviceInformation_DEV0006E02D11
	case "DEV0006E03011":
		return DeviceInformation_DEV0006E03011
	case "DEV000C570411":
		return DeviceInformation_DEV000C570411
	case "DEV0006E0C110":
		return DeviceInformation_DEV0006E0C110
	case "DEV0006E0C510":
		return DeviceInformation_DEV0006E0C510
	case "DEV0006B00A01":
		return DeviceInformation_DEV0006B00A01
	case "DEV0006B00602":
		return DeviceInformation_DEV0006B00602
	case "DEV0006E0C410":
		return DeviceInformation_DEV0006E0C410
	case "DEV0006E0C312":
		return DeviceInformation_DEV0006E0C312
	case "DEV0006E0C210":
		return DeviceInformation_DEV0006E0C210
	case "DEV0006209016":
		return DeviceInformation_DEV0006209016
	case "DEV0006E01A01":
		return DeviceInformation_DEV0006E01A01
	case "DEV0006E09910":
		return DeviceInformation_DEV0006E09910
	case "DEV000C570110":
		return DeviceInformation_DEV000C570110
	case "DEV0006E03710":
		return DeviceInformation_DEV0006E03710
	case "DEV0006209011":
		return DeviceInformation_DEV0006209011
	case "DEV000620A011":
		return DeviceInformation_DEV000620A011
	case "DEV0006E02410":
		return DeviceInformation_DEV0006E02410
	case "DEV0006E02301":
		return DeviceInformation_DEV0006E02301
	case "DEV0006E02510":
		return DeviceInformation_DEV0006E02510
	case "DEV0006E01B01":
		return DeviceInformation_DEV0006E01B01
	case "DEV0006E01C01":
		return DeviceInformation_DEV0006E01C01
	case "DEV0006E01D01":
		return DeviceInformation_DEV0006E01D01
	case "DEV0006E01E01":
		return DeviceInformation_DEV0006E01E01
	case "DEV000C570011":
		return DeviceInformation_DEV000C570011
	case "DEV0006207812":
		return DeviceInformation_DEV0006207812
	case "DEV0006B00811":
		return DeviceInformation_DEV0006B00811
	case "DEV0006E01001":
		return DeviceInformation_DEV0006E01001
	case "DEV0006E03610":
		return DeviceInformation_DEV0006E03610
	case "DEV0006E09810":
		return DeviceInformation_DEV0006E09810
	case "DEV0006208811":
		return DeviceInformation_DEV0006208811
	case "DEV0006209811":
		return DeviceInformation_DEV0006209811
	case "DEV0006E02610":
		return DeviceInformation_DEV0006E02610
	case "DEV0006E02710":
		return DeviceInformation_DEV0006E02710
	case "DEV0006E02A10":
		return DeviceInformation_DEV0006E02A10
	case "DEV000C20BD11":
		return DeviceInformation_DEV000C20BD11
	case "DEV0006E02B10":
		return DeviceInformation_DEV0006E02B10
	case "DEV0006E00C10":
		return DeviceInformation_DEV0006E00C10
	case "DEV0006010110":
		return DeviceInformation_DEV0006010110
	case "DEV0006010210":
		return DeviceInformation_DEV0006010210
	case "DEV0006E00B10":
		return DeviceInformation_DEV0006E00B10
	case "DEV0006E09C10":
		return DeviceInformation_DEV0006E09C10
	case "DEV0006E09B10":
		return DeviceInformation_DEV0006E09B10
	case "DEV0006E03510":
		return DeviceInformation_DEV0006E03510
	case "DEV0006FF1B11":
		return DeviceInformation_DEV0006FF1B11
	case "DEV0006E0CF10":
		return DeviceInformation_DEV0006E0CF10
	case "DEV000C20BA11":
		return DeviceInformation_DEV000C20BA11
	case "DEV000620A812":
		return DeviceInformation_DEV000620A812
	case "DEV000620CD11":
		return DeviceInformation_DEV000620CD11
	case "DEV0006E00E01":
		return DeviceInformation_DEV0006E00E01
	case "DEV0006E02201":
		return DeviceInformation_DEV0006E02201
	case "DEV000620AD11":
		return DeviceInformation_DEV000620AD11
	case "DEV0006E00F01":
		return DeviceInformation_DEV0006E00F01
	case "DEV0006E02101":
		return DeviceInformation_DEV0006E02101
	case "DEV000620BD11":
		return DeviceInformation_DEV000620BD11
	case "DEV0006E00D01":
		return DeviceInformation_DEV0006E00D01
	case "DEV0006E03910":
		return DeviceInformation_DEV0006E03910
	case "DEV000C760110":
		return DeviceInformation_DEV000C760110
	case "DEV0006E02810":
		return DeviceInformation_DEV0006E02810
	case "DEV0006E02910":
		return DeviceInformation_DEV0006E02910
	case "DEV0006E02C10":
		return DeviceInformation_DEV0006E02C10
	case "DEV0006C00403":
		return DeviceInformation_DEV0006C00403
	case "DEV0006590101":
		return DeviceInformation_DEV0006590101
	case "DEV0006E0CC11":
		return DeviceInformation_DEV0006E0CC11
	case "DEV0006E09A10":
		return DeviceInformation_DEV0006E09A10
	case "DEV0006E03811":
		return DeviceInformation_DEV0006E03811
	case "DEV0006E0C710":
		return DeviceInformation_DEV0006E0C710
	case "DEV0006E0C610":
		return DeviceInformation_DEV0006E0C610
	case "DEV0064181710":
		return DeviceInformation_DEV0064181710
	case "DEV000C705C01":
		return DeviceInformation_DEV000C705C01
	case "DEV0006E05A10":
		return DeviceInformation_DEV0006E05A10
	case "DEV0048493A1C":
		return DeviceInformation_DEV0048493A1C
	case "DEV0048494712":
		return DeviceInformation_DEV0048494712
	case "DEV0048494810":
		return DeviceInformation_DEV0048494810
	case "DEV0048855A10":
		return DeviceInformation_DEV0048855A10
	case "DEV0048855B10":
		return DeviceInformation_DEV0048855B10
	case "DEV0048A05713":
		return DeviceInformation_DEV0048A05713
	case "DEV0048494414":
		return DeviceInformation_DEV0048494414
	case "DEV0048824A11":
		return DeviceInformation_DEV0048824A11
	case "DEV0048824A12":
		return DeviceInformation_DEV0048824A12
	case "DEV000CFF2112":
		return DeviceInformation_DEV000CFF2112
	case "DEV0048770A10":
		return DeviceInformation_DEV0048770A10
	case "DEV0048494311":
		return DeviceInformation_DEV0048494311
	case "DEV0048494513":
		return DeviceInformation_DEV0048494513
	case "DEV0048494012":
		return DeviceInformation_DEV0048494012
	case "DEV0048494111":
		return DeviceInformation_DEV0048494111
	case "DEV0048494210":
		return DeviceInformation_DEV0048494210
	case "DEV0048494610":
		return DeviceInformation_DEV0048494610
	case "DEV0048494910":
		return DeviceInformation_DEV0048494910
	case "DEV0048134A10":
		return DeviceInformation_DEV0048134A10
	case "DEV0048107E12":
		return DeviceInformation_DEV0048107E12
	case "DEV000CB00812":
		return DeviceInformation_DEV000CB00812
	case "DEV0048FF2112":
		return DeviceInformation_DEV0048FF2112
	case "DEV0048140A11":
		return DeviceInformation_DEV0048140A11
	case "DEV0048140B12":
		return DeviceInformation_DEV0048140B12
	case "DEV0048140C13":
		return DeviceInformation_DEV0048140C13
	case "DEV0048139A10":
		return DeviceInformation_DEV0048139A10
	case "DEV0048648B10":
		return DeviceInformation_DEV0048648B10
	case "DEV0008A01111":
		return DeviceInformation_DEV0008A01111
	case "DEV0008A01211":
		return DeviceInformation_DEV0008A01211
	case "DEV0008A01212":
		return DeviceInformation_DEV0008A01212
	case "DEV0008A01112":
		return DeviceInformation_DEV0008A01112
	case "DEV000CB00713":
		return DeviceInformation_DEV000CB00713
	case "DEV0008A03213":
		return DeviceInformation_DEV0008A03213
	case "DEV0008A03313":
		return DeviceInformation_DEV0008A03313
	case "DEV0008A01113":
		return DeviceInformation_DEV0008A01113
	case "DEV0008A01711":
		return DeviceInformation_DEV0008A01711
	case "DEV0008B00911":
		return DeviceInformation_DEV0008B00911
	case "DEV0008C00102":
		return DeviceInformation_DEV0008C00102
	case "DEV0008C00101":
		return DeviceInformation_DEV0008C00101
	case "DEV0008901501":
		return DeviceInformation_DEV0008901501
	case "DEV0008901310":
		return DeviceInformation_DEV0008901310
	case "DEV000820B011":
		return DeviceInformation_DEV000820B011
	case "DEV000C181910":
		return DeviceInformation_DEV000C181910
	case "DEV0008705C11":
		return DeviceInformation_DEV0008705C11
	case "DEV0008705D11":
		return DeviceInformation_DEV0008705D11
	case "DEV0008706211":
		return DeviceInformation_DEV0008706211
	case "DEV000820BA11":
		return DeviceInformation_DEV000820BA11
	case "DEV000820C011":
		return DeviceInformation_DEV000820C011
	case "DEV000820B311":
		return DeviceInformation_DEV000820B311
	case "DEV0008301A11":
		return DeviceInformation_DEV0008301A11
	case "DEV0008C00C13":
		return DeviceInformation_DEV0008C00C13
	case "DEV0008302611":
		return DeviceInformation_DEV0008302611
	case "DEV0008302311":
		return DeviceInformation_DEV0008302311
	case "DEV000C181810":
		return DeviceInformation_DEV000C181810
	case "DEV0008302011":
		return DeviceInformation_DEV0008302011
	case "DEV0008C00C11":
		return DeviceInformation_DEV0008C00C11
	case "DEV0008302612":
		return DeviceInformation_DEV0008302612
	case "DEV0008302312":
		return DeviceInformation_DEV0008302312
	case "DEV0008302012":
		return DeviceInformation_DEV0008302012
	case "DEV0008C00C15":
		return DeviceInformation_DEV0008C00C15
	case "DEV0008C00C14":
		return DeviceInformation_DEV0008C00C14
	case "DEV0008B00713":
		return DeviceInformation_DEV0008B00713
	case "DEV0008706611":
		return DeviceInformation_DEV0008706611
	case "DEV0008706811":
		return DeviceInformation_DEV0008706811
	case "DEV000C20C011":
		return DeviceInformation_DEV000C20C011
	case "DEV0008B00812":
		return DeviceInformation_DEV0008B00812
	case "DEV0008209016":
		return DeviceInformation_DEV0008209016
	case "DEV0008209011":
		return DeviceInformation_DEV0008209011
	case "DEV000820A011":
		return DeviceInformation_DEV000820A011
	case "DEV0008208811":
		return DeviceInformation_DEV0008208811
	case "DEV0008209811":
		return DeviceInformation_DEV0008209811
	case "DEV000820CA11":
		return DeviceInformation_DEV000820CA11
	case "DEV0008208012":
		return DeviceInformation_DEV0008208012
	case "DEV0008207812":
		return DeviceInformation_DEV0008207812
	case "DEV0008207811":
		return DeviceInformation_DEV0008207811
	case "DEV0079002527":
		return DeviceInformation_DEV0079002527
	case "DEV0008208011":
		return DeviceInformation_DEV0008208011
	case "DEV000810D111":
		return DeviceInformation_DEV000810D111
	case "DEV000810D511":
		return DeviceInformation_DEV000810D511
	case "DEV000810FA12":
		return DeviceInformation_DEV000810FA12
	case "DEV000810FB12":
		return DeviceInformation_DEV000810FB12
	case "DEV000810F211":
		return DeviceInformation_DEV000810F211
	case "DEV000810D211":
		return DeviceInformation_DEV000810D211
	case "DEV000810E211":
		return DeviceInformation_DEV000810E211
	case "DEV000810D611":
		return DeviceInformation_DEV000810D611
	case "DEV000810F212":
		return DeviceInformation_DEV000810F212
	case "DEV0079004027":
		return DeviceInformation_DEV0079004027
	case "DEV000810E212":
		return DeviceInformation_DEV000810E212
	case "DEV000810FC13":
		return DeviceInformation_DEV000810FC13
	case "DEV000810FD13":
		return DeviceInformation_DEV000810FD13
	case "DEV000810F311":
		return DeviceInformation_DEV000810F311
	case "DEV000810D311":
		return DeviceInformation_DEV000810D311
	case "DEV000810D711":
		return DeviceInformation_DEV000810D711
	case "DEV000810F312":
		return DeviceInformation_DEV000810F312
	case "DEV000810D811":
		return DeviceInformation_DEV000810D811
	case "DEV000810E511":
		return DeviceInformation_DEV000810E511
	case "DEV000810E512":
		return DeviceInformation_DEV000810E512
	case "DEV0079000223":
		return DeviceInformation_DEV0079000223
	case "DEV000810F611":
		return DeviceInformation_DEV000810F611
	case "DEV000810D911":
		return DeviceInformation_DEV000810D911
	case "DEV000810F612":
		return DeviceInformation_DEV000810F612
	case "DEV000820A812":
		return DeviceInformation_DEV000820A812
	case "DEV000820AD11":
		return DeviceInformation_DEV000820AD11
	case "DEV000820BD11":
		return DeviceInformation_DEV000820BD11
	case "DEV000820C711":
		return DeviceInformation_DEV000820C711
	case "DEV000820CD11":
		return DeviceInformation_DEV000820CD11
	case "DEV000820C411":
		return DeviceInformation_DEV000820C411
	case "DEV000820A811":
		return DeviceInformation_DEV000820A811
	case "DEV0064181610":
		return DeviceInformation_DEV0064181610
	case "DEV0079000123":
		return DeviceInformation_DEV0079000123
	case "DEV0008501411":
		return DeviceInformation_DEV0008501411
	case "DEV0008C01602":
		return DeviceInformation_DEV0008C01602
	case "DEV0008302613":
		return DeviceInformation_DEV0008302613
	case "DEV0008302313":
		return DeviceInformation_DEV0008302313
	case "DEV0008302013":
		return DeviceInformation_DEV0008302013
	case "DEV0009207730":
		return DeviceInformation_DEV0009207730
	case "DEV0009208F10":
		return DeviceInformation_DEV0009208F10
	case "DEV0009C00C13":
		return DeviceInformation_DEV0009C00C13
	case "DEV0009209910":
		return DeviceInformation_DEV0009209910
	case "DEV0009209A10":
		return DeviceInformation_DEV0009209A10
	case "DEV0079001427":
		return DeviceInformation_DEV0079001427
	case "DEV0009207930":
		return DeviceInformation_DEV0009207930
	case "DEV0009201720":
		return DeviceInformation_DEV0009201720
	case "DEV0009500D01":
		return DeviceInformation_DEV0009500D01
	case "DEV0009500E01":
		return DeviceInformation_DEV0009500E01
	case "DEV0009209911":
		return DeviceInformation_DEV0009209911
	case "DEV0009209A11":
		return DeviceInformation_DEV0009209A11
	case "DEV0009C00C12":
		return DeviceInformation_DEV0009C00C12
	case "DEV0009C00C11":
		return DeviceInformation_DEV0009C00C11
	case "DEV0009500D20":
		return DeviceInformation_DEV0009500D20
	case "DEV0009500E20":
		return DeviceInformation_DEV0009500E20
	case "DEV0079003027":
		return DeviceInformation_DEV0079003027
	case "DEV000920B910":
		return DeviceInformation_DEV000920B910
	case "DEV0009E0CE10":
		return DeviceInformation_DEV0009E0CE10
	case "DEV0009E0A210":
		return DeviceInformation_DEV0009E0A210
	case "DEV0009501410":
		return DeviceInformation_DEV0009501410
	case "DEV0009207830":
		return DeviceInformation_DEV0009207830
	case "DEV0009201620":
		return DeviceInformation_DEV0009201620
	case "DEV0009E0A111":
		return DeviceInformation_DEV0009E0A111
	case "DEV0009E0CD11":
		return DeviceInformation_DEV0009E0CD11
	case "DEV000920B811":
		return DeviceInformation_DEV000920B811
	case "DEV000920B611":
		return DeviceInformation_DEV000920B611
	case "DEV0079100C13":
		return DeviceInformation_DEV0079100C13
	case "DEV0009207E10":
		return DeviceInformation_DEV0009207E10
	case "DEV0009207630":
		return DeviceInformation_DEV0009207630
	case "DEV0009205910":
		return DeviceInformation_DEV0009205910
	case "DEV0009500B01":
		return DeviceInformation_DEV0009500B01
	case "DEV000920AC10":
		return DeviceInformation_DEV000920AC10
	case "DEV0009207430":
		return DeviceInformation_DEV0009207430
	case "DEV0009204521":
		return DeviceInformation_DEV0009204521
	case "DEV0009500A01":
		return DeviceInformation_DEV0009500A01
	case "DEV0009500001":
		return DeviceInformation_DEV0009500001
	case "DEV000920AB10":
		return DeviceInformation_DEV000920AB10
	case "DEV0079101C11":
		return DeviceInformation_DEV0079101C11
	case "DEV000920BF11":
		return DeviceInformation_DEV000920BF11
	case "DEV0009203510":
		return DeviceInformation_DEV0009203510
	case "DEV0009207A30":
		return DeviceInformation_DEV0009207A30
	case "DEV0009500701":
		return DeviceInformation_DEV0009500701
	case "DEV0009501710":
		return DeviceInformation_DEV0009501710
	case "DEV000920B310":
		return DeviceInformation_DEV000920B310
	case "DEV0009207530":
		return DeviceInformation_DEV0009207530
	case "DEV0009203321":
		return DeviceInformation_DEV0009203321
	case "DEV0009500C01":
		return DeviceInformation_DEV0009500C01
	case "DEV000920AD10":
		return DeviceInformation_DEV000920AD10
	case "DEV0080709010":
		return DeviceInformation_DEV0080709010
	case "DEV0009207230":
		return DeviceInformation_DEV0009207230
	case "DEV0009500801":
		return DeviceInformation_DEV0009500801
	case "DEV0009501810":
		return DeviceInformation_DEV0009501810
	case "DEV000920B410":
		return DeviceInformation_DEV000920B410
	case "DEV0009207330":
		return DeviceInformation_DEV0009207330
	case "DEV0009204421":
		return DeviceInformation_DEV0009204421
	case "DEV0009500901":
		return DeviceInformation_DEV0009500901
	case "DEV000920AA10":
		return DeviceInformation_DEV000920AA10
	case "DEV0009209D01":
		return DeviceInformation_DEV0009209D01
	case "DEV000920B010":
		return DeviceInformation_DEV000920B010
	case "DEV0080707010":
		return DeviceInformation_DEV0080707010
	case "DEV0009E0BE01":
		return DeviceInformation_DEV0009E0BE01
	case "DEV000920B110":
		return DeviceInformation_DEV000920B110
	case "DEV0009E0BD01":
		return DeviceInformation_DEV0009E0BD01
	case "DEV0009D03F10":
		return DeviceInformation_DEV0009D03F10
	case "DEV0009305F10":
		return DeviceInformation_DEV0009305F10
	case "DEV0009305610":
		return DeviceInformation_DEV0009305610
	case "DEV0009D03E10":
		return DeviceInformation_DEV0009D03E10
	case "DEV0009306010":
		return DeviceInformation_DEV0009306010
	case "DEV0009306110":
		return DeviceInformation_DEV0009306110
	case "DEV0009306310":
		return DeviceInformation_DEV0009306310
	case "DEV0080706010":
		return DeviceInformation_DEV0080706010
	case "DEV0009D03B10":
		return DeviceInformation_DEV0009D03B10
	case "DEV0009D03C10":
		return DeviceInformation_DEV0009D03C10
	case "DEV0009D03910":
		return DeviceInformation_DEV0009D03910
	case "DEV0009D03A10":
		return DeviceInformation_DEV0009D03A10
	case "DEV0009305411":
		return DeviceInformation_DEV0009305411
	case "DEV0009D03D11":
		return DeviceInformation_DEV0009D03D11
	case "DEV0009304B11":
		return DeviceInformation_DEV0009304B11
	case "DEV0009304C11":
		return DeviceInformation_DEV0009304C11
	case "DEV0009306220":
		return DeviceInformation_DEV0009306220
	case "DEV0009302E10":
		return DeviceInformation_DEV0009302E10
	case "DEV0080706810":
		return DeviceInformation_DEV0080706810
	case "DEV0009302F10":
		return DeviceInformation_DEV0009302F10
	case "DEV0009303010":
		return DeviceInformation_DEV0009303010
	case "DEV0009303110":
		return DeviceInformation_DEV0009303110
	case "DEV0009306510":
		return DeviceInformation_DEV0009306510
	case "DEV0009306610":
		return DeviceInformation_DEV0009306610
	case "DEV0009306410":
		return DeviceInformation_DEV0009306410
	case "DEV0009401110":
		return DeviceInformation_DEV0009401110
	case "DEV0009400610":
		return DeviceInformation_DEV0009400610
	case "DEV0009401510":
		return DeviceInformation_DEV0009401510
	case "DEV0009402110":
		return DeviceInformation_DEV0009402110
	case "DEV0080705010":
		return DeviceInformation_DEV0080705010
	case "DEV0009400110":
		return DeviceInformation_DEV0009400110
	case "DEV0009400910":
		return DeviceInformation_DEV0009400910
	case "DEV0009400010":
		return DeviceInformation_DEV0009400010
	case "DEV0009401810":
		return DeviceInformation_DEV0009401810
	case "DEV0009400310":
		return DeviceInformation_DEV0009400310
	case "DEV0009301810":
		return DeviceInformation_DEV0009301810
	case "DEV0009301910":
		return DeviceInformation_DEV0009301910
	case "DEV0009301A10":
		return DeviceInformation_DEV0009301A10
	case "DEV0009401210":
		return DeviceInformation_DEV0009401210
	case "DEV0009400810":
		return DeviceInformation_DEV0009400810
	case "DEV006420C011":
		return DeviceInformation_DEV006420C011
	case "DEV0080703013":
		return DeviceInformation_DEV0080703013
	case "DEV0009400710":
		return DeviceInformation_DEV0009400710
	case "DEV0009401310":
		return DeviceInformation_DEV0009401310
	case "DEV0009401410":
		return DeviceInformation_DEV0009401410
	case "DEV0009402210":
		return DeviceInformation_DEV0009402210
	case "DEV0009402310":
		return DeviceInformation_DEV0009402310
	case "DEV0009401710":
		return DeviceInformation_DEV0009401710
	case "DEV0009401610":
		return DeviceInformation_DEV0009401610
	case "DEV0009400210":
		return DeviceInformation_DEV0009400210
	case "DEV0009401010":
		return DeviceInformation_DEV0009401010
	case "DEV0009400510":
		return DeviceInformation_DEV0009400510
	case "DEV0080704021":
		return DeviceInformation_DEV0080704021
	case "DEV0009400410":
		return DeviceInformation_DEV0009400410
	case "DEV0009D04B20":
		return DeviceInformation_DEV0009D04B20
	case "DEV0009D04920":
		return DeviceInformation_DEV0009D04920
	case "DEV0009D04A20":
		return DeviceInformation_DEV0009D04A20
	case "DEV0009D04820":
		return DeviceInformation_DEV0009D04820
	case "DEV0009D04C11":
		return DeviceInformation_DEV0009D04C11
	case "DEV0009D05610":
		return DeviceInformation_DEV0009D05610
	case "DEV0009305510":
		return DeviceInformation_DEV0009305510
	case "DEV0009209810":
		return DeviceInformation_DEV0009209810
	case "DEV0009202A10":
		return DeviceInformation_DEV0009202A10
	case "DEV0080704022":
		return DeviceInformation_DEV0080704022
	case "DEV0009209510":
		return DeviceInformation_DEV0009209510
	case "DEV0009501110":
		return DeviceInformation_DEV0009501110
	case "DEV0009209310":
		return DeviceInformation_DEV0009209310
	case "DEV0009209410":
		return DeviceInformation_DEV0009209410
	case "DEV0009209210":
		return DeviceInformation_DEV0009209210
	case "DEV0009501210":
		return DeviceInformation_DEV0009501210
	case "DEV0009205411":
		return DeviceInformation_DEV0009205411
	case "DEV000920A111":
		return DeviceInformation_DEV000920A111
	case "DEV000920A311":
		return DeviceInformation_DEV000920A311
	case "DEV0009205112":
		return DeviceInformation_DEV0009205112
	case "DEV0080704020":
		return DeviceInformation_DEV0080704020
	case "DEV0009204110":
		return DeviceInformation_DEV0009204110
	case "DEV0009E07710":
		return DeviceInformation_DEV0009E07710
	case "DEV0009E07712":
		return DeviceInformation_DEV0009E07712
	case "DEV0009205212":
		return DeviceInformation_DEV0009205212
	case "DEV0009205211":
		return DeviceInformation_DEV0009205211
	case "DEV0009205311":
		return DeviceInformation_DEV0009205311
	case "DEV0009206B10":
		return DeviceInformation_DEV0009206B10
	case "DEV0009208010":
		return DeviceInformation_DEV0009208010
	case "DEV0009206A12":
		return DeviceInformation_DEV0009206A12
	case "DEV0009206810":
		return DeviceInformation_DEV0009206810
	case "DEV0080701111":
		return DeviceInformation_DEV0080701111
	case "DEV0009208110":
		return DeviceInformation_DEV0009208110
	case "DEV0009205511":
		return DeviceInformation_DEV0009205511
	case "DEV0009209F01":
		return DeviceInformation_DEV0009209F01
	case "DEV0009208C10":
		return DeviceInformation_DEV0009208C10
	case "DEV0009208E10":
		return DeviceInformation_DEV0009208E10
	case "DEV000920B511":
		return DeviceInformation_DEV000920B511
	case "DEV0009501910":
		return DeviceInformation_DEV0009501910
	case "DEV000920BE11":
		return DeviceInformation_DEV000920BE11
	case "DEV0009209710":
		return DeviceInformation_DEV0009209710
	case "DEV0009208510":
		return DeviceInformation_DEV0009208510
	case "DEV0080701811":
		return DeviceInformation_DEV0080701811
	case "DEV0009208610":
		return DeviceInformation_DEV0009208610
	case "DEV000920BD10":
		return DeviceInformation_DEV000920BD10
	case "DEV0009500210":
		return DeviceInformation_DEV0009500210
	case "DEV0009500310":
		return DeviceInformation_DEV0009500310
	case "DEV0009E0BF10":
		return DeviceInformation_DEV0009E0BF10
	case "DEV0009E0C010":
		return DeviceInformation_DEV0009E0C010
	case "DEV0009500110":
		return DeviceInformation_DEV0009500110
	case "DEV0009209B10":
		return DeviceInformation_DEV0009209B10
	case "DEV0009207D10":
		return DeviceInformation_DEV0009207D10
	case "DEV0009202F11":
		return DeviceInformation_DEV0009202F11
	case "DEV008020A110":
		return DeviceInformation_DEV008020A110
	case "DEV0009203011":
		return DeviceInformation_DEV0009203011
	case "DEV0009207C10":
		return DeviceInformation_DEV0009207C10
	case "DEV0009207B10":
		return DeviceInformation_DEV0009207B10
	case "DEV0009208710":
		return DeviceInformation_DEV0009208710
	case "DEV0009E06610":
		return DeviceInformation_DEV0009E06610
	case "DEV0009E06611":
		return DeviceInformation_DEV0009E06611
	case "DEV0009E06410":
		return DeviceInformation_DEV0009E06410
	case "DEV0009E06411":
		return DeviceInformation_DEV0009E06411
	case "DEV0009E06210":
		return DeviceInformation_DEV0009E06210
	case "DEV0009E0E910":
		return DeviceInformation_DEV0009E0E910
	case "DEV008020A210":
		return DeviceInformation_DEV008020A210
	case "DEV0009E0EB10":
		return DeviceInformation_DEV0009E0EB10
	case "DEV000920BB10":
		return DeviceInformation_DEV000920BB10
	case "DEV0009FF1B11":
		return DeviceInformation_DEV0009FF1B11
	case "DEV0009E0CF10":
		return DeviceInformation_DEV0009E0CF10
	case "DEV0009206C30":
		return DeviceInformation_DEV0009206C30
	case "DEV0009206D30":
		return DeviceInformation_DEV0009206D30
	case "DEV0009206E30":
		return DeviceInformation_DEV0009206E30
	case "DEV0009206F30":
		return DeviceInformation_DEV0009206F30
	case "DEV0009207130":
		return DeviceInformation_DEV0009207130
	case "DEV0009204720":
		return DeviceInformation_DEV0009204720
	case "DEV008020A010":
		return DeviceInformation_DEV008020A010
	case "DEV0009204820":
		return DeviceInformation_DEV0009204820
	case "DEV0009204920":
		return DeviceInformation_DEV0009204920
	case "DEV0009204A20":
		return DeviceInformation_DEV0009204A20
	case "DEV0009205A10":
		return DeviceInformation_DEV0009205A10
	case "DEV0009207030":
		return DeviceInformation_DEV0009207030
	case "DEV0009205B10":
		return DeviceInformation_DEV0009205B10
	case "DEV0009500501":
		return DeviceInformation_DEV0009500501
	case "DEV0009501001":
		return DeviceInformation_DEV0009501001
	case "DEV0009500601":
		return DeviceInformation_DEV0009500601
	case "DEV0009500F01":
		return DeviceInformation_DEV0009500F01
	case "DEV0080207212":
		return DeviceInformation_DEV0080207212
	case "DEV0009500401":
		return DeviceInformation_DEV0009500401
	case "DEV000920B210":
		return DeviceInformation_DEV000920B210
	case "DEV000920AE10":
		return DeviceInformation_DEV000920AE10
	case "DEV000920BC10":
		return DeviceInformation_DEV000920BC10
	case "DEV000920AF10":
		return DeviceInformation_DEV000920AF10
	case "DEV0009207F10":
		return DeviceInformation_DEV0009207F10
	case "DEV0009208910":
		return DeviceInformation_DEV0009208910
	case "DEV0009205710":
		return DeviceInformation_DEV0009205710
	case "DEV0009205810":
		return DeviceInformation_DEV0009205810
	case "DEV0009203810":
		return DeviceInformation_DEV0009203810
	case "DEV006420BA11":
		return DeviceInformation_DEV006420BA11
	case "DEV0080209111":
		return DeviceInformation_DEV0080209111
	case "DEV0009203910":
		return DeviceInformation_DEV0009203910
	case "DEV0009203E10":
		return DeviceInformation_DEV0009203E10
	case "DEV0009204B10":
		return DeviceInformation_DEV0009204B10
	case "DEV0009203F10":
		return DeviceInformation_DEV0009203F10
	case "DEV0009204C10":
		return DeviceInformation_DEV0009204C10
	case "DEV0009204010":
		return DeviceInformation_DEV0009204010
	case "DEV0009206411":
		return DeviceInformation_DEV0009206411
	case "DEV0009205E10":
		return DeviceInformation_DEV0009205E10
	case "DEV0009206711":
		return DeviceInformation_DEV0009206711
	case "DEV000920A710":
		return DeviceInformation_DEV000920A710
	case "DEV0080204310":
		return DeviceInformation_DEV0080204310
	case "DEV000920A610":
		return DeviceInformation_DEV000920A610
	case "DEV0009203A10":
		return DeviceInformation_DEV0009203A10
	case "DEV0009203B10":
		return DeviceInformation_DEV0009203B10
	case "DEV0009203C10":
		return DeviceInformation_DEV0009203C10
	case "DEV0009203D10":
		return DeviceInformation_DEV0009203D10
	case "DEV0009E05E12":
		return DeviceInformation_DEV0009E05E12
	case "DEV0009E0B711":
		return DeviceInformation_DEV0009E0B711
	case "DEV0009E06A12":
		return DeviceInformation_DEV0009E06A12
	case "DEV0009E06E12":
		return DeviceInformation_DEV0009E06E12
	case "DEV0009E0B720":
		return DeviceInformation_DEV0009E0B720
	case "DEV008020B612":
		return DeviceInformation_DEV008020B612
	case "DEV0009E0E611":
		return DeviceInformation_DEV0009E0E611
	case "DEV0009E0B321":
		return DeviceInformation_DEV0009E0B321
	case "DEV0009E0E512":
		return DeviceInformation_DEV0009E0E512
	case "DEV0009204210":
		return DeviceInformation_DEV0009204210
	case "DEV0009208210":
		return DeviceInformation_DEV0009208210
	case "DEV0009E07211":
		return DeviceInformation_DEV0009E07211
	case "DEV0009E0CC11":
		return DeviceInformation_DEV0009E0CC11
	case "DEV0009110111":
		return DeviceInformation_DEV0009110111
	case "DEV0009110211":
		return DeviceInformation_DEV0009110211
	case "DEV000916B212":
		return DeviceInformation_DEV000916B212
	case "DEV008020B412":
		return DeviceInformation_DEV008020B412
	case "DEV0009110212":
		return DeviceInformation_DEV0009110212
	case "DEV0009110311":
		return DeviceInformation_DEV0009110311
	case "DEV000916B312":
		return DeviceInformation_DEV000916B312
	case "DEV0009110312":
		return DeviceInformation_DEV0009110312
	case "DEV0009110411":
		return DeviceInformation_DEV0009110411
	case "DEV0009110412":
		return DeviceInformation_DEV0009110412
	case "DEV0009501615":
		return DeviceInformation_DEV0009501615
	case "DEV0009E0ED10":
		return DeviceInformation_DEV0009E0ED10
	case "DEV014F030110":
		return DeviceInformation_DEV014F030110
	case "DEV014F030310":
		return DeviceInformation_DEV014F030310
	case "DEV008020B512":
		return DeviceInformation_DEV008020B512
	case "DEV014F030210":
		return DeviceInformation_DEV014F030210
	case "DEV00EE7FFF10":
		return DeviceInformation_DEV00EE7FFF10
	case "DEV00B6464101":
		return DeviceInformation_DEV00B6464101
	case "DEV00B6464201":
		return DeviceInformation_DEV00B6464201
	case "DEV00B6464501":
		return DeviceInformation_DEV00B6464501
	case "DEV00B6434101":
		return DeviceInformation_DEV00B6434101
	case "DEV00B6434201":
		return DeviceInformation_DEV00B6434201
	case "DEV00B6434202":
		return DeviceInformation_DEV00B6434202
	case "DEV00B6454101":
		return DeviceInformation_DEV00B6454101
	case "DEV00B6454201":
		return DeviceInformation_DEV00B6454201
	case "DEV0080208310":
		return DeviceInformation_DEV0080208310
	case "DEV00B6455001":
		return DeviceInformation_DEV00B6455001
	case "DEV00B6453101":
		return DeviceInformation_DEV00B6453101
	case "DEV00B6453102":
		return DeviceInformation_DEV00B6453102
	case "DEV00B6454102":
		return DeviceInformation_DEV00B6454102
	case "DEV00B6454401":
		return DeviceInformation_DEV00B6454401
	case "DEV00B6454402":
		return DeviceInformation_DEV00B6454402
	case "DEV00B6454202":
		return DeviceInformation_DEV00B6454202
	case "DEV00B6453103":
		return DeviceInformation_DEV00B6453103
	case "DEV00B6453201":
		return DeviceInformation_DEV00B6453201
	case "DEV00B6453301":
		return DeviceInformation_DEV00B6453301
	case "DEV0080702111":
		return DeviceInformation_DEV0080702111
	case "DEV00B6453104":
		return DeviceInformation_DEV00B6453104
	case "DEV00B6454403":
		return DeviceInformation_DEV00B6454403
	case "DEV00B6454801":
		return DeviceInformation_DEV00B6454801
	case "DEV00B6414701":
		return DeviceInformation_DEV00B6414701
	case "DEV00B6414201":
		return DeviceInformation_DEV00B6414201
	case "DEV00B6474101":
		return DeviceInformation_DEV00B6474101
	case "DEV00B6474302":
		return DeviceInformation_DEV00B6474302
	case "DEV00B6474602":
		return DeviceInformation_DEV00B6474602
	case "DEV00B6534D01":
		return DeviceInformation_DEV00B6534D01
	case "DEV00B6535001":
		return DeviceInformation_DEV00B6535001
	case "DEV0081FE0111":
		return DeviceInformation_DEV0081FE0111
	case "DEV00B6455002":
		return DeviceInformation_DEV00B6455002
	case "DEV00B6453701":
		return DeviceInformation_DEV00B6453701
	case "DEV00B6484101":
		return DeviceInformation_DEV00B6484101
	case "DEV00B6484201":
		return DeviceInformation_DEV00B6484201
	case "DEV00B6484202":
		return DeviceInformation_DEV00B6484202
	case "DEV00B6484301":
		return DeviceInformation_DEV00B6484301
	case "DEV00B6484102":
		return DeviceInformation_DEV00B6484102
	case "DEV00B6455101":
		return DeviceInformation_DEV00B6455101
	case "DEV00B6455003":
		return DeviceInformation_DEV00B6455003
	case "DEV00B6455102":
		return DeviceInformation_DEV00B6455102
	case "DEV0081FF3131":
		return DeviceInformation_DEV0081FF3131
	case "DEV00B6453702":
		return DeviceInformation_DEV00B6453702
	case "DEV00B6453703":
		return DeviceInformation_DEV00B6453703
	case "DEV00B6484302":
		return DeviceInformation_DEV00B6484302
	case "DEV00B6484801":
		return DeviceInformation_DEV00B6484801
	case "DEV00B6484501":
		return DeviceInformation_DEV00B6484501
	case "DEV00B6484203":
		return DeviceInformation_DEV00B6484203
	case "DEV00B6484103":
		return DeviceInformation_DEV00B6484103
	case "DEV00B6455004":
		return DeviceInformation_DEV00B6455004
	case "DEV00B6455103":
		return DeviceInformation_DEV00B6455103
	case "DEV00B6455401":
		return DeviceInformation_DEV00B6455401
	case "DEV0081F01313":
		return DeviceInformation_DEV0081F01313
	case "DEV00B6455201":
		return DeviceInformation_DEV00B6455201
	case "DEV00B6455402":
		return DeviceInformation_DEV00B6455402
	case "DEV00B6455403":
		return DeviceInformation_DEV00B6455403
	case "DEV00B603430A":
		return DeviceInformation_DEV00B603430A
	case "DEV00B600010A":
		return DeviceInformation_DEV00B600010A
	case "DEV00B6FF110A":
		return DeviceInformation_DEV00B6FF110A
	case "DEV00B6434601":
		return DeviceInformation_DEV00B6434601
	case "DEV00B6434602":
		return DeviceInformation_DEV00B6434602
	case "DEV00B6455301":
		return DeviceInformation_DEV00B6455301
	case "DEV00C5070410":
		return DeviceInformation_DEV00C5070410
	case "DEV0064182010":
		return DeviceInformation_DEV0064182010
	case "DEV0083002C16":
		return DeviceInformation_DEV0083002C16
	case "DEV00C5070210":
		return DeviceInformation_DEV00C5070210
	case "DEV00C5070610":
		return DeviceInformation_DEV00C5070610
	case "DEV00C5070E11":
		return DeviceInformation_DEV00C5070E11
	case "DEV00C5060240":
		return DeviceInformation_DEV00C5060240
	case "DEV00C5062010":
		return DeviceInformation_DEV00C5062010
	case "DEV00C5080230":
		return DeviceInformation_DEV00C5080230
	case "DEV00C5060310":
		return DeviceInformation_DEV00C5060310
	case "DEV006C070E11":
		return DeviceInformation_DEV006C070E11
	case "DEV006C050002":
		return DeviceInformation_DEV006C050002
	case "DEV006C011311":
		return DeviceInformation_DEV006C011311
	case "DEV0083002E16":
		return DeviceInformation_DEV0083002E16
	case "DEV006C011411":
		return DeviceInformation_DEV006C011411
	case "DEV0007632010":
		return DeviceInformation_DEV0007632010
	case "DEV0007632020":
		return DeviceInformation_DEV0007632020
	case "DEV0007632180":
		return DeviceInformation_DEV0007632180
	case "DEV0007632040":
		return DeviceInformation_DEV0007632040
	case "DEV0007613812":
		return DeviceInformation_DEV0007613812
	case "DEV0007613810":
		return DeviceInformation_DEV0007613810
	case "DEV000720C011":
		return DeviceInformation_DEV000720C011
	case "DEV0007A05210":
		return DeviceInformation_DEV0007A05210
	case "DEV0007A08B10":
		return DeviceInformation_DEV0007A08B10
	case "DEV0083002F16":
		return DeviceInformation_DEV0083002F16
	case "DEV0007A05B32":
		return DeviceInformation_DEV0007A05B32
	case "DEV0007A06932":
		return DeviceInformation_DEV0007A06932
	case "DEV0007A06D32":
		return DeviceInformation_DEV0007A06D32
	case "DEV0007A08032":
		return DeviceInformation_DEV0007A08032
	case "DEV0007A00213":
		return DeviceInformation_DEV0007A00213
	case "DEV0007A09532":
		return DeviceInformation_DEV0007A09532
	case "DEV0007A06C32":
		return DeviceInformation_DEV0007A06C32
	case "DEV0007A05E32":
		return DeviceInformation_DEV0007A05E32
	case "DEV0007A08A32":
		return DeviceInformation_DEV0007A08A32
	case "DEV0007A07032":
		return DeviceInformation_DEV0007A07032
	case "DEV0083012F16":
		return DeviceInformation_DEV0083012F16
	case "DEV0007A08332":
		return DeviceInformation_DEV0007A08332
	case "DEV0007A09832":
		return DeviceInformation_DEV0007A09832
	case "DEV0007A05C32":
		return DeviceInformation_DEV0007A05C32
	case "DEV0007A06A32":
		return DeviceInformation_DEV0007A06A32
	case "DEV0007A08832":
		return DeviceInformation_DEV0007A08832
	case "DEV0007A06E32":
		return DeviceInformation_DEV0007A06E32
	case "DEV0007A08132":
		return DeviceInformation_DEV0007A08132
	case "DEV0007A00113":
		return DeviceInformation_DEV0007A00113
	case "DEV0007A09632":
		return DeviceInformation_DEV0007A09632
	case "DEV0007A05D32":
		return DeviceInformation_DEV0007A05D32
	case "DEV0083003210":
		return DeviceInformation_DEV0083003210
	case "DEV0007A06B32":
		return DeviceInformation_DEV0007A06B32
	case "DEV0007A08932":
		return DeviceInformation_DEV0007A08932
	case "DEV0007A06F32":
		return DeviceInformation_DEV0007A06F32
	case "DEV0007A08232":
		return DeviceInformation_DEV0007A08232
	case "DEV0007A09732":
		return DeviceInformation_DEV0007A09732
	case "DEV0007A05713":
		return DeviceInformation_DEV0007A05713
	case "DEV0007A01811":
		return DeviceInformation_DEV0007A01811
	case "DEV0007A01911":
		return DeviceInformation_DEV0007A01911
	case "DEV0007A04912":
		return DeviceInformation_DEV0007A04912
	case "DEV0007A05814":
		return DeviceInformation_DEV0007A05814
	case "DEV0083001D13":
		return DeviceInformation_DEV0083001D13
	case "DEV0007A07114":
		return DeviceInformation_DEV0007A07114
	case "DEV0007A05810":
		return DeviceInformation_DEV0007A05810
	case "DEV0007A04312":
		return DeviceInformation_DEV0007A04312
	case "DEV0007A04412":
		return DeviceInformation_DEV0007A04412
	case "DEV0007A04512":
		return DeviceInformation_DEV0007A04512
	case "DEV000720BD11":
		return DeviceInformation_DEV000720BD11
	case "DEV0007A04C13":
		return DeviceInformation_DEV0007A04C13
	case "DEV0007A04D13":
		return DeviceInformation_DEV0007A04D13
	case "DEV0007A04B10":
		return DeviceInformation_DEV0007A04B10
	case "DEV0007A04E13":
		return DeviceInformation_DEV0007A04E13
	case "DEV0083001E13":
		return DeviceInformation_DEV0083001E13
	case "DEV0007A04F13":
		return DeviceInformation_DEV0007A04F13
	case "DEV000720BA11":
		return DeviceInformation_DEV000720BA11
	case "DEV0007A03D11":
		return DeviceInformation_DEV0007A03D11
	case "DEV0007A09211":
		return DeviceInformation_DEV0007A09211
	case "DEV0007A09111":
		return DeviceInformation_DEV0007A09111
	case "DEV0007FF1115":
		return DeviceInformation_DEV0007FF1115
	case "DEV0007A01511":
		return DeviceInformation_DEV0007A01511
	case "DEV0007A08411":
		return DeviceInformation_DEV0007A08411
	case "DEV0007A08511":
		return DeviceInformation_DEV0007A08511
	case "DEV0007A03422":
		return DeviceInformation_DEV0007A03422
	case "DEV0083001B13":
		return DeviceInformation_DEV0083001B13
	case "DEV0007A07213":
		return DeviceInformation_DEV0007A07213
	case "DEV0007A07420":
		return DeviceInformation_DEV0007A07420
	case "DEV0007A07520":
		return DeviceInformation_DEV0007A07520
	case "DEV0007A07B12":
		return DeviceInformation_DEV0007A07B12
	case "DEV0007A07C12":
		return DeviceInformation_DEV0007A07C12
	case "DEV0007A09311":
		return DeviceInformation_DEV0007A09311
	case "DEV0007A09013":
		return DeviceInformation_DEV0007A09013
	case "DEV0007A08F13":
		return DeviceInformation_DEV0007A08F13
	case "DEV0007A07E10":
		return DeviceInformation_DEV0007A07E10
	case "DEV0007A05510":
		return DeviceInformation_DEV0007A05510
	case "DEV0083001C13":
		return DeviceInformation_DEV0083001C13
	case "DEV0007A05910":
		return DeviceInformation_DEV0007A05910
	case "DEV0007A08711":
		return DeviceInformation_DEV0007A08711
	case "DEV0007A03D12":
		return DeviceInformation_DEV0007A03D12
	case "DEV0007A09A12":
		return DeviceInformation_DEV0007A09A12
	case "DEV0007A09B12":
		return DeviceInformation_DEV0007A09B12
	case "DEV0007A06614":
		return DeviceInformation_DEV0007A06614
	case "DEV0007A06514":
		return DeviceInformation_DEV0007A06514
	case "DEV0007A06014":
		return DeviceInformation_DEV0007A06014
	case "DEV0007A07714":
		return DeviceInformation_DEV0007A07714
	case "DEV0007A06414":
		return DeviceInformation_DEV0007A06414
	case "DEV0083001F11":
		return DeviceInformation_DEV0083001F11
	case "DEV0007A06114":
		return DeviceInformation_DEV0007A06114
	case "DEV0007A07814":
		return DeviceInformation_DEV0007A07814
	case "DEV0007A06714":
		return DeviceInformation_DEV0007A06714
	case "DEV0007A06214":
		return DeviceInformation_DEV0007A06214
	case "DEV0007A07914":
		return DeviceInformation_DEV0007A07914
	case "DEV000B0A8410":
		return DeviceInformation_DEV000B0A8410
	case "DEV000B0A7E10":
		return DeviceInformation_DEV000B0A7E10
	case "DEV000B0A7F10":
		return DeviceInformation_DEV000B0A7F10
	case "DEV000B0A8010":
		return DeviceInformation_DEV000B0A8010
	case "DEV000BBF9111":
		return DeviceInformation_DEV000BBF9111
	case "DEV0064182510":
		return DeviceInformation_DEV0064182510
	case "DEV0083003C10":
		return DeviceInformation_DEV0083003C10
	case "DEV000B0A7810":
		return DeviceInformation_DEV000B0A7810
	case "DEV000B0A7910":
		return DeviceInformation_DEV000B0A7910
	case "DEV000B0A7A10":
		return DeviceInformation_DEV000B0A7A10
	case "DEV000B0A8910":
		return DeviceInformation_DEV000B0A8910
	case "DEV000B0A8310":
		return DeviceInformation_DEV000B0A8310
	case "DEV000B0A8510":
		return DeviceInformation_DEV000B0A8510
	case "DEV000B0A6319":
		return DeviceInformation_DEV000B0A6319
	case "DEV0083001C20":
		return DeviceInformation_DEV0083001C20
	case "DEV0083001B22":
		return DeviceInformation_DEV0083001B22
	case "DEV0083003A14":
		return DeviceInformation_DEV0083003A14
	case "DEV0083003B14":
		return DeviceInformation_DEV0083003B14
	case "DEV0083003B24":
		return DeviceInformation_DEV0083003B24
	case "DEV0083003A24":
		return DeviceInformation_DEV0083003A24
	case "DEV0083005824":
		return DeviceInformation_DEV0083005824
	case "DEV0083002828":
		return DeviceInformation_DEV0083002828
	case "DEV0083002928":
		return DeviceInformation_DEV0083002928
	case "DEV0064182610":
		return DeviceInformation_DEV0064182610
	case "DEV0083002A18":
		return DeviceInformation_DEV0083002A18
	case "DEV0083002B18":
		return DeviceInformation_DEV0083002B18
	case "DEV0083002337":
		return DeviceInformation_DEV0083002337
	case "DEV0083002838":
		return DeviceInformation_DEV0083002838
	case "DEV0083002938":
		return DeviceInformation_DEV0083002938
	case "DEV0083002A38":
		return DeviceInformation_DEV0083002A38
	case "DEV0083002B38":
		return DeviceInformation_DEV0083002B38
	case "DEV0083001321":
		return DeviceInformation_DEV0083001321
	case "DEV0083001421":
		return DeviceInformation_DEV0083001421
	case "DEV0083001521":
		return DeviceInformation_DEV0083001521
	case "DEV0064182910":
		return DeviceInformation_DEV0064182910
	case "DEV0083001621":
		return DeviceInformation_DEV0083001621
	case "DEV0083000921":
		return DeviceInformation_DEV0083000921
	case "DEV0083000D11":
		return DeviceInformation_DEV0083000D11
	case "DEV0083000C11":
		return DeviceInformation_DEV0083000C11
	case "DEV0083000E11":
		return DeviceInformation_DEV0083000E11
	case "DEV0083000B11":
		return DeviceInformation_DEV0083000B11
	case "DEV0083000A11":
		return DeviceInformation_DEV0083000A11
	case "DEV0083000A21":
		return DeviceInformation_DEV0083000A21
	case "DEV0083000B21":
		return DeviceInformation_DEV0083000B21
	case "DEV0083000C21":
		return DeviceInformation_DEV0083000C21
	case "DEV0001140C13":
		return DeviceInformation_DEV0001140C13
	case "DEV0064130610":
		return DeviceInformation_DEV0064130610
	case "DEV0083000D21":
		return DeviceInformation_DEV0083000D21
	case "DEV0083000821":
		return DeviceInformation_DEV0083000821
	case "DEV0083000E21":
		return DeviceInformation_DEV0083000E21
	case "DEV0083001812":
		return DeviceInformation_DEV0083001812
	case "DEV0083001712":
		return DeviceInformation_DEV0083001712
	case "DEV0083001816":
		return DeviceInformation_DEV0083001816
	case "DEV0083001916":
		return DeviceInformation_DEV0083001916
	case "DEV0083001716":
		return DeviceInformation_DEV0083001716
	case "DEV0083001921":
		return DeviceInformation_DEV0083001921
	case "DEV0083001721":
		return DeviceInformation_DEV0083001721
	case "DEV0064130710":
		return DeviceInformation_DEV0064130710
	case "DEV0083001821":
		return DeviceInformation_DEV0083001821
	case "DEV0083001A20":
		return DeviceInformation_DEV0083001A20
	case "DEV0083002320":
		return DeviceInformation_DEV0083002320
	case "DEV0083001010":
		return DeviceInformation_DEV0083001010
	case "DEV0083000F10":
		return DeviceInformation_DEV0083000F10
	case "DEV0083003D14":
		return DeviceInformation_DEV0083003D14
	case "DEV0083003E14":
		return DeviceInformation_DEV0083003E14
	case "DEV0083003F14":
		return DeviceInformation_DEV0083003F14
	case "DEV0083004014":
		return DeviceInformation_DEV0083004014
	case "DEV0083004024":
		return DeviceInformation_DEV0083004024
	case "DEV0064133510":
		return DeviceInformation_DEV0064133510
	case "DEV0083003D24":
		return DeviceInformation_DEV0083003D24
	case "DEV0083003E24":
		return DeviceInformation_DEV0083003E24
	case "DEV0083003F24":
		return DeviceInformation_DEV0083003F24
	case "DEV0083001112":
		return DeviceInformation_DEV0083001112
	case "DEV0083001212":
		return DeviceInformation_DEV0083001212
	case "DEV0083005B12":
		return DeviceInformation_DEV0083005B12
	case "DEV0083005A12":
		return DeviceInformation_DEV0083005A12
	case "DEV0083008410":
		return DeviceInformation_DEV0083008410
	case "DEV0083008510":
		return DeviceInformation_DEV0083008510
	case "DEV0083008610":
		return DeviceInformation_DEV0083008610
	case "DEV0064133310":
		return DeviceInformation_DEV0064133310
	case "DEV0083008710":
		return DeviceInformation_DEV0083008710
	case "DEV0083002515":
		return DeviceInformation_DEV0083002515
	case "DEV0083002115":
		return DeviceInformation_DEV0083002115
	case "DEV0083002015":
		return DeviceInformation_DEV0083002015
	case "DEV0083002415":
		return DeviceInformation_DEV0083002415
	case "DEV0083002615":
		return DeviceInformation_DEV0083002615
	case "DEV0083002215":
		return DeviceInformation_DEV0083002215
	case "DEV0083002715":
		return DeviceInformation_DEV0083002715
	case "DEV0083002315":
		return DeviceInformation_DEV0083002315
	case "DEV0083008B25":
		return DeviceInformation_DEV0083008B25
	case "DEV0064133410":
		return DeviceInformation_DEV0064133410
	case "DEV0083008A25":
		return DeviceInformation_DEV0083008A25
	case "DEV0083008B28":
		return DeviceInformation_DEV0083008B28
	case "DEV0083008A28":
		return DeviceInformation_DEV0083008A28
	case "DEV0083009013":
		return DeviceInformation_DEV0083009013
	case "DEV0083009213":
		return DeviceInformation_DEV0083009213
	case "DEV0083009113":
		return DeviceInformation_DEV0083009113
	case "DEV0083009313":
		return DeviceInformation_DEV0083009313
	case "DEV0083009413":
		return DeviceInformation_DEV0083009413
	case "DEV0083009513":
		return DeviceInformation_DEV0083009513
	case "DEV0083009613":
		return DeviceInformation_DEV0083009613
	case "DEV0064133610":
		return DeviceInformation_DEV0064133610
	case "DEV0083009713":
		return DeviceInformation_DEV0083009713
	case "DEV0083009A13":
		return DeviceInformation_DEV0083009A13
	case "DEV0083009B13":
		return DeviceInformation_DEV0083009B13
	case "DEV0083004B11":
		return DeviceInformation_DEV0083004B11
	case "DEV0083004B20":
		return DeviceInformation_DEV0083004B20
	case "DEV0083005514":
		return DeviceInformation_DEV0083005514
	case "DEV0083006824":
		return DeviceInformation_DEV0083006824
	case "DEV0083006624":
		return DeviceInformation_DEV0083006624
	case "DEV0083006524":
		return DeviceInformation_DEV0083006524
	case "DEV0083006424":
		return DeviceInformation_DEV0083006424
	case "DEV0064130510":
		return DeviceInformation_DEV0064130510
	case "DEV0083006734":
		return DeviceInformation_DEV0083006734
	case "DEV0083006434":
		return DeviceInformation_DEV0083006434
	case "DEV0083006634":
		return DeviceInformation_DEV0083006634
	case "DEV0083006534":
		return DeviceInformation_DEV0083006534
	case "DEV0083006A34":
		return DeviceInformation_DEV0083006A34
	case "DEV0083006B34":
		return DeviceInformation_DEV0083006B34
	case "DEV0083006934":
		return DeviceInformation_DEV0083006934
	case "DEV0083004F11":
		return DeviceInformation_DEV0083004F11
	case "DEV0083004E10":
		return DeviceInformation_DEV0083004E10
	case "DEV0083004D13":
		return DeviceInformation_DEV0083004D13
	case "DEV0064480611":
		return DeviceInformation_DEV0064480611
	case "DEV0083004414":
		return DeviceInformation_DEV0083004414
	case "DEV0083004114":
		return DeviceInformation_DEV0083004114
	case "DEV0083004514":
		return DeviceInformation_DEV0083004514
	case "DEV0083004213":
		return DeviceInformation_DEV0083004213
	case "DEV0083004313":
		return DeviceInformation_DEV0083004313
	case "DEV0083004C11":
		return DeviceInformation_DEV0083004C11
	case "DEV0083004913":
		return DeviceInformation_DEV0083004913
	case "DEV0083004A13":
		return DeviceInformation_DEV0083004A13
	case "DEV0083004712":
		return DeviceInformation_DEV0083004712
	case "DEV0083004610":
		return DeviceInformation_DEV0083004610
	case "DEV0064482011":
		return DeviceInformation_DEV0064482011
	case "DEV0083008E12":
		return DeviceInformation_DEV0083008E12
	case "DEV0083004813":
		return DeviceInformation_DEV0083004813
	case "DEV0083005611":
		return DeviceInformation_DEV0083005611
	case "DEV0083005710":
		return DeviceInformation_DEV0083005710
	case "DEV0083005010":
		return DeviceInformation_DEV0083005010
	case "DEV0083001A10":
		return DeviceInformation_DEV0083001A10
	case "DEV0083002918":
		return DeviceInformation_DEV0083002918
	case "DEV0083002818":
		return DeviceInformation_DEV0083002818
	case "DEV0083006724":
		return DeviceInformation_DEV0083006724
	case "DEV0083006D41":
		return DeviceInformation_DEV0083006D41
	case "DEV0064182210":
		return DeviceInformation_DEV0064182210
	case "DEV0083006E41":
		return DeviceInformation_DEV0083006E41
	case "DEV0083007342":
		return DeviceInformation_DEV0083007342
	case "DEV0083007242":
		return DeviceInformation_DEV0083007242
	case "DEV0083006C42":
		return DeviceInformation_DEV0083006C42
	case "DEV0083007542":
		return DeviceInformation_DEV0083007542
	case "DEV0083007442":
		return DeviceInformation_DEV0083007442
	case "DEV0083007742":
		return DeviceInformation_DEV0083007742
	case "DEV0083007642":
		return DeviceInformation_DEV0083007642
	case "DEV008300B030":
		return DeviceInformation_DEV008300B030
	case "DEV008300B130":
		return DeviceInformation_DEV008300B130
	case "DEV0001140B11":
		return DeviceInformation_DEV0001140B11
	case "DEV0064182710":
		return DeviceInformation_DEV0064182710
	case "DEV008300B230":
		return DeviceInformation_DEV008300B230
	case "DEV008300B330":
		return DeviceInformation_DEV008300B330
	case "DEV008300B430":
		return DeviceInformation_DEV008300B430
	case "DEV008300B530":
		return DeviceInformation_DEV008300B530
	case "DEV008300B630":
		return DeviceInformation_DEV008300B630
	case "DEV008300B730":
		return DeviceInformation_DEV008300B730
	case "DEV0083012843":
		return DeviceInformation_DEV0083012843
	case "DEV0083012943":
		return DeviceInformation_DEV0083012943
	case "DEV008300A421":
		return DeviceInformation_DEV008300A421
	case "DEV008300A521":
		return DeviceInformation_DEV008300A521
	case "DEV0064183010":
		return DeviceInformation_DEV0064183010
	case "DEV008300A621":
		return DeviceInformation_DEV008300A621
	case "DEV0083001332":
		return DeviceInformation_DEV0083001332
	case "DEV0083000932":
		return DeviceInformation_DEV0083000932
	case "DEV0083001432":
		return DeviceInformation_DEV0083001432
	case "DEV0083001532":
		return DeviceInformation_DEV0083001532
	case "DEV0083001632":
		return DeviceInformation_DEV0083001632
	case "DEV008300A432":
		return DeviceInformation_DEV008300A432
	case "DEV008300A532":
		return DeviceInformation_DEV008300A532
	case "DEV008300A632":
		return DeviceInformation_DEV008300A632
	case "DEV0083000F32":
		return DeviceInformation_DEV0083000F32
	case "DEV0064B00812":
		return DeviceInformation_DEV0064B00812
	case "DEV0083001032":
		return DeviceInformation_DEV0083001032
	case "DEV0083000632":
		return DeviceInformation_DEV0083000632
	case "DEV0083009810":
		return DeviceInformation_DEV0083009810
	case "DEV0083009910":
		return DeviceInformation_DEV0083009910
	case "DEV0083005C11":
		return DeviceInformation_DEV0083005C11
	case "DEV0083005D11":
		return DeviceInformation_DEV0083005D11
	case "DEV0083005E11":
		return DeviceInformation_DEV0083005E11
	case "DEV0083005F11":
		return DeviceInformation_DEV0083005F11
	case "DEV0083005413":
		return DeviceInformation_DEV0083005413
	case "DEV0085000520":
		return DeviceInformation_DEV0085000520
	case "DEV0064B00A01":
		return DeviceInformation_DEV0064B00A01
	case "DEV0085000620":
		return DeviceInformation_DEV0085000620
	case "DEV0085000720":
		return DeviceInformation_DEV0085000720
	case "DEV0085012210":
		return DeviceInformation_DEV0085012210
	case "DEV0085011210":
		return DeviceInformation_DEV0085011210
	case "DEV0085013220":
		return DeviceInformation_DEV0085013220
	case "DEV0085010210":
		return DeviceInformation_DEV0085010210
	case "DEV0085000A10":
		return DeviceInformation_DEV0085000A10
	case "DEV0085000B10":
		return DeviceInformation_DEV0085000B10
	case "DEV0085071010":
		return DeviceInformation_DEV0085071010
	case "DEV008500FB10":
		return DeviceInformation_DEV008500FB10
	case "DEV0064760110":
		return DeviceInformation_DEV0064760110
	case "DEV0085060210":
		return DeviceInformation_DEV0085060210
	case "DEV0085060110":
		return DeviceInformation_DEV0085060110
	case "DEV0085000D20":
		return DeviceInformation_DEV0085000D20
	case "DEV008500C810":
		return DeviceInformation_DEV008500C810
	case "DEV0085040111":
		return DeviceInformation_DEV0085040111
	case "DEV008500C910":
		return DeviceInformation_DEV008500C910
	case "DEV0085045020":
		return DeviceInformation_DEV0085045020
	case "DEV0085070210":
		return DeviceInformation_DEV0085070210
	case "DEV0085070110":
		return DeviceInformation_DEV0085070110
	case "DEV0085070310":
		return DeviceInformation_DEV0085070310
	case "DEV0064242313":
		return DeviceInformation_DEV0064242313
	case "DEV0085000E20":
		return DeviceInformation_DEV0085000E20
	case "DEV008E596010":
		return DeviceInformation_DEV008E596010
	case "DEV008E593710":
		return DeviceInformation_DEV008E593710
	case "DEV008E597710":
		return DeviceInformation_DEV008E597710
	case "DEV008E598310":
		return DeviceInformation_DEV008E598310
	case "DEV008E598910":
		return DeviceInformation_DEV008E598910
	case "DEV008E593720":
		return DeviceInformation_DEV008E593720
	case "DEV008E598920":
		return DeviceInformation_DEV008E598920
	case "DEV008E598320":
		return DeviceInformation_DEV008E598320
	case "DEV008E596021":
		return DeviceInformation_DEV008E596021
	case "DEV0064FF2111":
		return DeviceInformation_DEV0064FF2111
	case "DEV008E597721":
		return DeviceInformation_DEV008E597721
	case "DEV008E587320":
		return DeviceInformation_DEV008E587320
	case "DEV008E587020":
		return DeviceInformation_DEV008E587020
	case "DEV008E587220":
		return DeviceInformation_DEV008E587220
	case "DEV008E587120":
		return DeviceInformation_DEV008E587120
	case "DEV008E679910":
		return DeviceInformation_DEV008E679910
	case "DEV008E618310":
		return DeviceInformation_DEV008E618310
	case "DEV008E707910":
		return DeviceInformation_DEV008E707910
	case "DEV008E004010":
		return DeviceInformation_DEV008E004010
	case "DEV008E570910":
		return DeviceInformation_DEV008E570910
	case "DEV0064FF2112":
		return DeviceInformation_DEV0064FF2112
	case "DEV008E558810":
		return DeviceInformation_DEV008E558810
	case "DEV008E683410":
		return DeviceInformation_DEV008E683410
	case "DEV008E707710":
		return DeviceInformation_DEV008E707710
	case "DEV008E707810":
		return DeviceInformation_DEV008E707810
	case "DEV0091100013":
		return DeviceInformation_DEV0091100013
	case "DEV0091100110":
		return DeviceInformation_DEV0091100110
	case "DEV009E670101":
		return DeviceInformation_DEV009E670101
	case "DEV009E119311":
		return DeviceInformation_DEV009E119311
	case "DEV00A2100C13":
		return DeviceInformation_DEV00A2100C13
	case "DEV00A2101C11":
		return DeviceInformation_DEV00A2101C11
	case "DEV0064648B10":
		return DeviceInformation_DEV0064648B10
	case "DEV00A2300110":
		return DeviceInformation_DEV00A2300110
	case "DEV0002A05814":
		return DeviceInformation_DEV0002A05814
	case "DEV0002A07114":
		return DeviceInformation_DEV0002A07114
	case "DEV0002134A10":
		return DeviceInformation_DEV0002134A10
	case "DEV0002A03422":
		return DeviceInformation_DEV0002A03422
	case "DEV0002A03321":
		return DeviceInformation_DEV0002A03321
	case "DEV0002648B10":
		return DeviceInformation_DEV0002648B10
	case "DEV0002A09013":
		return DeviceInformation_DEV0002A09013
	case "DEV0002A08F13":
		return DeviceInformation_DEV0002A08F13
	case "DEV0002A05510":
		return DeviceInformation_DEV0002A05510
	case "DEV0064724010":
		return DeviceInformation_DEV0064724010
	case "DEV0002A05910":
		return DeviceInformation_DEV0002A05910
	case "DEV0002A05326":
		return DeviceInformation_DEV0002A05326
	case "DEV0002A05428":
		return DeviceInformation_DEV0002A05428
	case "DEV0002A08411":
		return DeviceInformation_DEV0002A08411
	case "DEV0002A08511":
		return DeviceInformation_DEV0002A08511
	case "DEV0002A00F11":
		return DeviceInformation_DEV0002A00F11
	case "DEV0002A07310":
		return DeviceInformation_DEV0002A07310
	case "DEV0002A04110":
		return DeviceInformation_DEV0002A04110
	case "DEV0002A03813":
		return DeviceInformation_DEV0002A03813
	case "DEV0002A07F13":
		return DeviceInformation_DEV0002A07F13
	case "DEV0001803002":
		return DeviceInformation_DEV0001803002
	case "DEV006420BD11":
		return DeviceInformation_DEV006420BD11
	case "DEV0002A08832":
		return DeviceInformation_DEV0002A08832
	case "DEV0002A06E32":
		return DeviceInformation_DEV0002A06E32
	case "DEV0002A08132":
		return DeviceInformation_DEV0002A08132
	case "DEV0002A01D20":
		return DeviceInformation_DEV0002A01D20
	case "DEV0002A02120":
		return DeviceInformation_DEV0002A02120
	case "DEV0002A02520":
		return DeviceInformation_DEV0002A02520
	case "DEV0002A02920":
		return DeviceInformation_DEV0002A02920
	case "DEV0002A03A20":
		return DeviceInformation_DEV0002A03A20
	case "DEV0002A05C32":
		return DeviceInformation_DEV0002A05C32
	case "DEV0002A06A32":
		return DeviceInformation_DEV0002A06A32
	case "DEV0064570011":
		return DeviceInformation_DEV0064570011
	case "DEV0002A09632":
		return DeviceInformation_DEV0002A09632
	case "DEV0002A08932":
		return DeviceInformation_DEV0002A08932
	case "DEV0002A06F32":
		return DeviceInformation_DEV0002A06F32
	case "DEV0002A08232":
		return DeviceInformation_DEV0002A08232
	case "DEV0002A01E20":
		return DeviceInformation_DEV0002A01E20
	case "DEV0002A02220":
		return DeviceInformation_DEV0002A02220
	case "DEV0002A02620":
		return DeviceInformation_DEV0002A02620
	case "DEV0002A02A20":
		return DeviceInformation_DEV0002A02A20
	case "DEV0002A03B20":
		return DeviceInformation_DEV0002A03B20
	case "DEV0002A05D32":
		return DeviceInformation_DEV0002A05D32
	case "DEV0064570310":
		return DeviceInformation_DEV0064570310
	case "DEV0002A06B32":
		return DeviceInformation_DEV0002A06B32
	case "DEV0002A09732":
		return DeviceInformation_DEV0002A09732
	case "DEV0002A08A32":
		return DeviceInformation_DEV0002A08A32
	case "DEV0002A07032":
		return DeviceInformation_DEV0002A07032
	case "DEV0002A08332":
		return DeviceInformation_DEV0002A08332
	case "DEV0002A01F20":
		return DeviceInformation_DEV0002A01F20
	case "DEV0002A02320":
		return DeviceInformation_DEV0002A02320
	case "DEV0002A02720":
		return DeviceInformation_DEV0002A02720
	case "DEV0002A02B20":
		return DeviceInformation_DEV0002A02B20
	case "DEV0002A04820":
		return DeviceInformation_DEV0002A04820
	case "DEV0064570211":
		return DeviceInformation_DEV0064570211
	case "DEV0002A06C32":
		return DeviceInformation_DEV0002A06C32
	case "DEV0002A05E32":
		return DeviceInformation_DEV0002A05E32
	case "DEV0002A09832":
		return DeviceInformation_DEV0002A09832
	case "DEV0002A06D32":
		return DeviceInformation_DEV0002A06D32
	case "DEV0002A08032":
		return DeviceInformation_DEV0002A08032
	case "DEV0002A02020":
		return DeviceInformation_DEV0002A02020
	case "DEV0002A02420":
		return DeviceInformation_DEV0002A02420
	case "DEV0002A02820":
		return DeviceInformation_DEV0002A02820
	case "DEV0002A03920":
		return DeviceInformation_DEV0002A03920
	case "DEV0002A05B32":
		return DeviceInformation_DEV0002A05B32
	case "DEV0064570411":
		return DeviceInformation_DEV0064570411
	case "DEV0002A06932":
		return DeviceInformation_DEV0002A06932
	case "DEV0002A09532":
		return DeviceInformation_DEV0002A09532
	case "DEV0002B00813":
		return DeviceInformation_DEV0002B00813
	case "DEV0002A0A610":
		return DeviceInformation_DEV0002A0A610
	case "DEV0002A0A611":
		return DeviceInformation_DEV0002A0A611
	case "DEV0002A0A510":
		return DeviceInformation_DEV0002A0A510
	case "DEV0002A0A511":
		return DeviceInformation_DEV0002A0A511
	case "DEV0002A00510":
		return DeviceInformation_DEV0002A00510
	case "DEV0002A00610":
		return DeviceInformation_DEV0002A00610
	case "DEV0002A01511":
		return DeviceInformation_DEV0002A01511
	case "DEV0064570110":
		return DeviceInformation_DEV0064570110
	case "DEV0002A03D11":
		return DeviceInformation_DEV0002A03D11
	case "DEV000220C011":
		return DeviceInformation_DEV000220C011
	case "DEV0002A05213":
		return DeviceInformation_DEV0002A05213
	case "DEV0002A08B10":
		return DeviceInformation_DEV0002A08B10
	case "DEV0002A0A210":
		return DeviceInformation_DEV0002A0A210
	case "DEV0002A0A010":
		return DeviceInformation_DEV0002A0A010
	case "DEV0002A09F10":
		return DeviceInformation_DEV0002A09F10
	case "DEV0002A0A110":
		return DeviceInformation_DEV0002A0A110
	case "DEV0002A0A013":
		return DeviceInformation_DEV0002A0A013
	case "DEV0002A09F13":
		return DeviceInformation_DEV0002A09F13
	case "DEV0064615022":
		return DeviceInformation_DEV0064615022
	case "DEV0002A0A213":
		return DeviceInformation_DEV0002A0A213
	case "DEV0002A0A113":
		return DeviceInformation_DEV0002A0A113
	case "DEV0002A03F11":
		return DeviceInformation_DEV0002A03F11
	case "DEV0002A04011":
		return DeviceInformation_DEV0002A04011
	case "DEV0002FF2112":
		return DeviceInformation_DEV0002FF2112
	case "DEV0002FF2111":
		return DeviceInformation_DEV0002FF2111
	case "DEV0002720111":
		return DeviceInformation_DEV0002720111
	case "DEV0002613812":
		return DeviceInformation_DEV0002613812
	case "DEV0002A05713":
		return DeviceInformation_DEV0002A05713
	case "DEV0002A07610":
		return DeviceInformation_DEV0002A07610
	case "DEV0064182810":
		return DeviceInformation_DEV0064182810
	case "DEV0002A01911":
		return DeviceInformation_DEV0002A01911
	case "DEV0002A07611":
		return DeviceInformation_DEV0002A07611
	case "DEV0002A04B10":
		return DeviceInformation_DEV0002A04B10
	case "DEV0002A01B14":
		return DeviceInformation_DEV0002A01B14
	case "DEV0002A09B11":
		return DeviceInformation_DEV0002A09B11
	case "DEV0002A09B12":
		return DeviceInformation_DEV0002A09B12
	case "DEV0002A03C10":
		return DeviceInformation_DEV0002A03C10
	case "DEV0002A00213":
		return DeviceInformation_DEV0002A00213
	case "DEV0002A00113":
		return DeviceInformation_DEV0002A00113
	case "DEV0002A02C12":
		return DeviceInformation_DEV0002A02C12
	case "DEV0064183110":
		return DeviceInformation_DEV0064183110
	case "DEV0002A02D12":
		return DeviceInformation_DEV0002A02D12
	case "DEV0002A02E12":
		return DeviceInformation_DEV0002A02E12
	case "DEV0002A04C13":
		return DeviceInformation_DEV0002A04C13
	case "DEV0002A04D13":
		return DeviceInformation_DEV0002A04D13
	case "DEV0002A02F12":
		return DeviceInformation_DEV0002A02F12
	case "DEV0002A03012":
		return DeviceInformation_DEV0002A03012
	case "DEV0002A03112":
		return DeviceInformation_DEV0002A03112
	case "DEV0002A04E13":
		return DeviceInformation_DEV0002A04E13
	case "DEV0002A04F13":
		return DeviceInformation_DEV0002A04F13
	case "DEV0002A01A13":
		return DeviceInformation_DEV0002A01A13
	case "DEV0064133611":
		return DeviceInformation_DEV0064133611
	case "DEV0002A09C11":
		return DeviceInformation_DEV0002A09C11
	case "DEV0002A09C12":
		return DeviceInformation_DEV0002A09C12
	case "DEV0002A01C20":
		return DeviceInformation_DEV0002A01C20
	case "DEV0002A09A10":
		return DeviceInformation_DEV0002A09A10
	case "DEV0002A09A12":
		return DeviceInformation_DEV0002A09A12
	case "DEV000220BA11":
		return DeviceInformation_DEV000220BA11
	case "DEV0002A03D12":
		return DeviceInformation_DEV0002A03D12
	case "DEV0002A09110":
		return DeviceInformation_DEV0002A09110
	case "DEV0002A09210":
		return DeviceInformation_DEV0002A09210
	case "DEV0002A09111":
		return DeviceInformation_DEV0002A09111
	case "DEV00641BD610":
		return DeviceInformation_DEV00641BD610
	case "DEV006A000122":
		return DeviceInformation_DEV006A000122
	case "DEV0002A09211":
		return DeviceInformation_DEV0002A09211
	case "DEV0002A00E21":
		return DeviceInformation_DEV0002A00E21
	case "DEV0002A03710":
		return DeviceInformation_DEV0002A03710
	case "DEV0002A01112":
		return DeviceInformation_DEV0002A01112
	case "DEV0002A01216":
		return DeviceInformation_DEV0002A01216
	case "DEV0002A01217":
		return DeviceInformation_DEV0002A01217
	case "DEV000220BD11":
		return DeviceInformation_DEV000220BD11
	case "DEV0002A07F12":
		return DeviceInformation_DEV0002A07F12
	case "DEV0002A06613":
		return DeviceInformation_DEV0002A06613
	case "DEV0002A06713":
		return DeviceInformation_DEV0002A06713
	case "DEV006A000222":
		return DeviceInformation_DEV006A000222
	case "DEV0002A06013":
		return DeviceInformation_DEV0002A06013
	case "DEV0002A06113":
		return DeviceInformation_DEV0002A06113
	case "DEV0002A06213":
		return DeviceInformation_DEV0002A06213
	case "DEV0002A06413":
		return DeviceInformation_DEV0002A06413
	case "DEV0002A07713":
		return DeviceInformation_DEV0002A07713
	case "DEV0002A07813":
		return DeviceInformation_DEV0002A07813
	case "DEV0002A07913":
		return DeviceInformation_DEV0002A07913
	case "DEV0002A07914":
		return DeviceInformation_DEV0002A07914
	case "DEV0002A06114":
		return DeviceInformation_DEV0002A06114
	case "DEV0002A06714":
		return DeviceInformation_DEV0002A06714
	case "DEV006A070210":
		return DeviceInformation_DEV006A070210
	case "DEV0002A06414":
		return DeviceInformation_DEV0002A06414
	case "DEV0002A06214":
		return DeviceInformation_DEV0002A06214
	case "DEV0002A06514":
		return DeviceInformation_DEV0002A06514
	case "DEV0002A07714":
		return DeviceInformation_DEV0002A07714
	case "DEV0002A06014":
		return DeviceInformation_DEV0002A06014
	case "DEV0002A06614":
		return DeviceInformation_DEV0002A06614
	case "DEV0002A07814":
		return DeviceInformation_DEV0002A07814
	case "DEV0002A0C310":
		return DeviceInformation_DEV0002A0C310
	case "DEV0002632010":
		return DeviceInformation_DEV0002632010
	case "DEV0002632020":
		return DeviceInformation_DEV0002632020
	case "DEV006BFFF713":
		return DeviceInformation_DEV006BFFF713
	case "DEV0002632040":
		return DeviceInformation_DEV0002632040
	case "DEV0002632180":
		return DeviceInformation_DEV0002632180
	case "DEV0002632170":
		return DeviceInformation_DEV0002632170
	case "DEV0002FF1140":
		return DeviceInformation_DEV0002FF1140
	case "DEV0002A07E10":
		return DeviceInformation_DEV0002A07E10
	case "DEV0002A07213":
		return DeviceInformation_DEV0002A07213
	case "DEV0002A04A35":
		return DeviceInformation_DEV0002A04A35
	case "DEV0002A07420":
		return DeviceInformation_DEV0002A07420
	case "DEV0002A07520":
		return DeviceInformation_DEV0002A07520
	case "DEV0002A07B12":
		return DeviceInformation_DEV0002A07B12
	case "DEV006BFF2111":
		return DeviceInformation_DEV006BFF2111
	case "DEV0002A07C12":
		return DeviceInformation_DEV0002A07C12
	case "DEV0002A04312":
		return DeviceInformation_DEV0002A04312
	case "DEV0002A04412":
		return DeviceInformation_DEV0002A04412
	case "DEV0002A04512":
		return DeviceInformation_DEV0002A04512
	case "DEV0002A04912":
		return DeviceInformation_DEV0002A04912
	case "DEV0002A05012":
		return DeviceInformation_DEV0002A05012
	case "DEV0002A01811":
		return DeviceInformation_DEV0002A01811
	case "DEV0002A03E11":
		return DeviceInformation_DEV0002A03E11
	case "DEV0002A08711":
		return DeviceInformation_DEV0002A08711
	case "DEV0002A09311":
		return DeviceInformation_DEV0002A09311
	case "DEV006BFFF820":
		return DeviceInformation_DEV006BFFF820
	case "DEV0002A01011":
		return DeviceInformation_DEV0002A01011
	case "DEV0002A01622":
		return DeviceInformation_DEV0002A01622
	case "DEV0002A04210":
		return DeviceInformation_DEV0002A04210
	case "DEV0002A09A13":
		return DeviceInformation_DEV0002A09A13
	case "DEV00C8272040":
		return DeviceInformation_DEV00C8272040
	case "DEV00C8272260":
		return DeviceInformation_DEV00C8272260
	case "DEV00C8272060":
		return DeviceInformation_DEV00C8272060
	case "DEV00C8272160":
		return DeviceInformation_DEV00C8272160
	case "DEV00C8272050":
		return DeviceInformation_DEV00C8272050
	case "DEV00C9106D10":
		return DeviceInformation_DEV00C9106D10
	case "DEV006B106D10":
		return DeviceInformation_DEV006B106D10
	case "DEV00C9107C20":
		return DeviceInformation_DEV00C9107C20
	case "DEV00C9108511":
		return DeviceInformation_DEV00C9108511
	case "DEV00C9106210":
		return DeviceInformation_DEV00C9106210
	case "DEV00C9109310":
		return DeviceInformation_DEV00C9109310
	case "DEV00C9109210":
		return DeviceInformation_DEV00C9109210
	case "DEV00C9109810":
		return DeviceInformation_DEV00C9109810
	case "DEV00C9109A10":
		return DeviceInformation_DEV00C9109A10
	case "DEV00C910A420":
		return DeviceInformation_DEV00C910A420
	case "DEV00C910A110":
		return DeviceInformation_DEV00C910A110
	case "DEV00C910A010":
		return DeviceInformation_DEV00C910A010
	case "DEV0071123130":
		return DeviceInformation_DEV0071123130
	case "DEV00C910A310":
		return DeviceInformation_DEV00C910A310
	case "DEV00C910A210":
		return DeviceInformation_DEV00C910A210
	case "DEV00C9109B10":
		return DeviceInformation_DEV00C9109B10
	case "DEV00C9106110":
		return DeviceInformation_DEV00C9106110
	case "DEV00C9109110":
		return DeviceInformation_DEV00C9109110
	case "DEV00C9109610":
		return DeviceInformation_DEV00C9109610
	case "DEV00C9109710":
		return DeviceInformation_DEV00C9109710
	case "DEV00C9109510":
		return DeviceInformation_DEV00C9109510
	case "DEV00C9109910":
		return DeviceInformation_DEV00C9109910
	case "DEV00C9109C10":
		return DeviceInformation_DEV00C9109C10
	case "DEV0071413133":
		return DeviceInformation_DEV0071413133
	case "DEV00C910AB10":
		return DeviceInformation_DEV00C910AB10
	case "DEV00C910AC10":
		return DeviceInformation_DEV00C910AC10
	case "DEV00C910AD10":
		return DeviceInformation_DEV00C910AD10
	case "DEV00C910A810":
		return DeviceInformation_DEV00C910A810
	case "DEV00C9106510":
		return DeviceInformation_DEV00C9106510
	case "DEV00C910A710":
		return DeviceInformation_DEV00C910A710
	case "DEV00C9107610":
		return DeviceInformation_DEV00C9107610
	case "DEV00C910890A":
		return DeviceInformation_DEV00C910890A
	case "DEV00C9FF1012":
		return DeviceInformation_DEV00C9FF1012
	case "DEV00C9FF0913":
		return DeviceInformation_DEV00C9FF0913
	case "DEV0071114019":
		return DeviceInformation_DEV0071114019
	case "DEV00C9FF1112":
		return DeviceInformation_DEV00C9FF1112
	case "DEV00C9100310":
		return DeviceInformation_DEV00C9100310
	case "DEV00C9101110":
		return DeviceInformation_DEV00C9101110
	case "DEV00C9101010":
		return DeviceInformation_DEV00C9101010
	case "DEV00C9103710":
		return DeviceInformation_DEV00C9103710
	case "DEV00C9101310":
		return DeviceInformation_DEV00C9101310
	case "DEV00C9FF0D12":
		return DeviceInformation_DEV00C9FF0D12
	case "DEV00C9100E10":
		return DeviceInformation_DEV00C9100E10
	case "DEV00C9100610":
		return DeviceInformation_DEV00C9100610
	case "DEV00C9100510":
		return DeviceInformation_DEV00C9100510
	case "DEV0064760210":
		return DeviceInformation_DEV0064760210
	case "DEV007111306C":
		return DeviceInformation_DEV007111306C
	case "DEV00C9100710":
		return DeviceInformation_DEV00C9100710
	case "DEV00C9FF1D20":
		return DeviceInformation_DEV00C9FF1D20
	case "DEV00C9FF1C10":
		return DeviceInformation_DEV00C9FF1C10
	case "DEV00C9100810":
		return DeviceInformation_DEV00C9100810
	case "DEV00C9FF1420":
		return DeviceInformation_DEV00C9FF1420
	case "DEV00C9100D10":
		return DeviceInformation_DEV00C9100D10
	case "DEV00C9101220":
		return DeviceInformation_DEV00C9101220
	case "DEV00C9102330":
		return DeviceInformation_DEV00C9102330
	case "DEV00C9102130":
		return DeviceInformation_DEV00C9102130
	case "DEV00C9102430":
		return DeviceInformation_DEV00C9102430
	case "DEV0071231112":
		return DeviceInformation_DEV0071231112
	case "DEV00C9100831":
		return DeviceInformation_DEV00C9100831
	case "DEV00C9102530":
		return DeviceInformation_DEV00C9102530
	case "DEV00C9100531":
		return DeviceInformation_DEV00C9100531
	case "DEV00C9102030":
		return DeviceInformation_DEV00C9102030
	case "DEV00C9100731":
		return DeviceInformation_DEV00C9100731
	case "DEV00C9100631":
		return DeviceInformation_DEV00C9100631
	case "DEV00C9102230":
		return DeviceInformation_DEV00C9102230
	case "DEV00C9100632":
		return DeviceInformation_DEV00C9100632
	case "DEV00C9100532":
		return DeviceInformation_DEV00C9100532
	case "DEV00C9100732":
		return DeviceInformation_DEV00C9100732
	case "DEV0071113080":
		return DeviceInformation_DEV0071113080
	case "DEV00C9100832":
		return DeviceInformation_DEV00C9100832
	case "DEV00C9102532":
		return DeviceInformation_DEV00C9102532
	case "DEV00C9102132":
		return DeviceInformation_DEV00C9102132
	case "DEV00C9102332":
		return DeviceInformation_DEV00C9102332
	case "DEV00C9102432":
		return DeviceInformation_DEV00C9102432
	case "DEV00C9102032":
		return DeviceInformation_DEV00C9102032
	case "DEV00C9102232":
		return DeviceInformation_DEV00C9102232
	case "DEV00C9104432":
		return DeviceInformation_DEV00C9104432
	case "DEV00C9100D11":
		return DeviceInformation_DEV00C9100D11
	case "DEV00C9100633":
		return DeviceInformation_DEV00C9100633
	case "DEV0071321212":
		return DeviceInformation_DEV0071321212
	case "DEV00C9100533":
		return DeviceInformation_DEV00C9100533
	case "DEV00C9100733":
		return DeviceInformation_DEV00C9100733
	case "DEV00C9100833":
		return DeviceInformation_DEV00C9100833
	case "DEV00C9102533":
		return DeviceInformation_DEV00C9102533
	case "DEV00C9102133":
		return DeviceInformation_DEV00C9102133
	case "DEV00C9102333":
		return DeviceInformation_DEV00C9102333
	case "DEV00C9102433":
		return DeviceInformation_DEV00C9102433
	case "DEV00C9102033":
		return DeviceInformation_DEV00C9102033
	case "DEV00C9102233":
		return DeviceInformation_DEV00C9102233
	case "DEV00C9104810":
		return DeviceInformation_DEV00C9104810
	case "DEV0071321113":
		return DeviceInformation_DEV0071321113
	case "DEV00C9FF1A11":
		return DeviceInformation_DEV00C9FF1A11
	case "DEV00C9100212":
		return DeviceInformation_DEV00C9100212
	case "DEV00C9FF0A11":
		return DeviceInformation_DEV00C9FF0A11
	case "DEV00C9FF0C12":
		return DeviceInformation_DEV00C9FF0C12
	case "DEV00C9100112":
		return DeviceInformation_DEV00C9100112
	case "DEV00C9FF1911":
		return DeviceInformation_DEV00C9FF1911
	case "DEV00C9FF0B12":
		return DeviceInformation_DEV00C9FF0B12
	case "DEV00C9FF0715":
		return DeviceInformation_DEV00C9FF0715
	case "DEV00C9FF1B10":
		return DeviceInformation_DEV00C9FF1B10
	case "DEV00C9101610":
		return DeviceInformation_DEV00C9101610
	case "DEV0071322212":
		return DeviceInformation_DEV0071322212
	case "DEV00C9FF1B11":
		return DeviceInformation_DEV00C9FF1B11
	case "DEV00C9101611":
		return DeviceInformation_DEV00C9101611
	case "DEV00C9101612":
		return DeviceInformation_DEV00C9101612
	case "DEV00C9FF0F11":
		return DeviceInformation_DEV00C9FF0F11
	case "DEV00C9FF1E30":
		return DeviceInformation_DEV00C9FF1E30
	case "DEV00C9100410":
		return DeviceInformation_DEV00C9100410
	case "DEV00C9106410":
		return DeviceInformation_DEV00C9106410
	case "DEV00C9106710":
		return DeviceInformation_DEV00C9106710
	case "DEV00C9106810":
		return DeviceInformation_DEV00C9106810
	case "DEV00C9106010":
		return DeviceInformation_DEV00C9106010
	case "DEV0071322112":
		return DeviceInformation_DEV0071322112
	case "DEV00C9106310":
		return DeviceInformation_DEV00C9106310
	case "DEV00C9107110":
		return DeviceInformation_DEV00C9107110
	case "DEV00C9107210":
		return DeviceInformation_DEV00C9107210
	case "DEV00C9107310":
		return DeviceInformation_DEV00C9107310
	case "DEV00C9107010":
		return DeviceInformation_DEV00C9107010
	case "DEV00C9107A20":
		return DeviceInformation_DEV00C9107A20
	case "DEV00C9107B20":
		return DeviceInformation_DEV00C9107B20
	case "DEV00C9107820":
		return DeviceInformation_DEV00C9107820
	case "DEV00C9107920":
		return DeviceInformation_DEV00C9107920
	case "DEV00C9104433":
		return DeviceInformation_DEV00C9104433
	case "DEV0071322312":
		return DeviceInformation_DEV0071322312
	case "DEV00C9107C11":
		return DeviceInformation_DEV00C9107C11
	case "DEV00C9107711":
		return DeviceInformation_DEV00C9107711
	case "DEV00C9108310":
		return DeviceInformation_DEV00C9108310
	case "DEV00C9108210":
		return DeviceInformation_DEV00C9108210
	case "DEV00C9108610":
		return DeviceInformation_DEV00C9108610
	case "DEV00C9107D10":
		return DeviceInformation_DEV00C9107D10
	case "DEV00CE648B10":
		return DeviceInformation_DEV00CE648B10
	case "DEV00CE494513":
		return DeviceInformation_DEV00CE494513
	case "DEV00CE494311":
		return DeviceInformation_DEV00CE494311
	case "DEV00CE494810":
		return DeviceInformation_DEV00CE494810
	case "DEV0071122124":
		return DeviceInformation_DEV0071122124
	case "DEV00CE494712":
		return DeviceInformation_DEV00CE494712
	case "DEV00CE494012":
		return DeviceInformation_DEV00CE494012
	case "DEV00CE494111":
		return DeviceInformation_DEV00CE494111
	case "DEV00CE494210":
		return DeviceInformation_DEV00CE494210
	case "DEV00CE494610":
		return DeviceInformation_DEV00CE494610
	case "DEV00CE494412":
		return DeviceInformation_DEV00CE494412
	case "DEV00D0660212":
		return DeviceInformation_DEV00D0660212
	case "DEV00E8000A10":
		return DeviceInformation_DEV00E8000A10
	case "DEV00E8000B10":
		return DeviceInformation_DEV00E8000B10
	case "DEV00E8000910":
		return DeviceInformation_DEV00E8000910
	case "DEV007112221E":
		return DeviceInformation_DEV007112221E
	case "DEV00E8001112":
		return DeviceInformation_DEV00E8001112
	case "DEV00E8000C14":
		return DeviceInformation_DEV00E8000C14
	case "DEV00E8000D13":
		return DeviceInformation_DEV00E8000D13
	case "DEV00E8000E12":
		return DeviceInformation_DEV00E8000E12
	case "DEV00E8001310":
		return DeviceInformation_DEV00E8001310
	case "DEV00E8001410":
		return DeviceInformation_DEV00E8001410
	case "DEV00E8001510":
		return DeviceInformation_DEV00E8001510
	case "DEV00E8000F10":
		return DeviceInformation_DEV00E8000F10
	case "DEV00E8001010":
		return DeviceInformation_DEV00E8001010
	case "DEV00E8000612":
		return DeviceInformation_DEV00E8000612
	case "DEV0064182410":
		return DeviceInformation_DEV0064182410
	case "DEV0071413314":
		return DeviceInformation_DEV0071413314
	case "DEV00E8000812":
		return DeviceInformation_DEV00E8000812
	case "DEV00E8000712":
		return DeviceInformation_DEV00E8000712
	case "DEV00F4501311":
		return DeviceInformation_DEV00F4501311
	case "DEV00F4B00911":
		return DeviceInformation_DEV00F4B00911
	case "DEV0019512710":
		return DeviceInformation_DEV0019512710
	case "DEV0019512810":
		return DeviceInformation_DEV0019512810
	case "DEV0019512910":
		return DeviceInformation_DEV0019512910
	case "DEV0019E30D10":
		return DeviceInformation_DEV0019E30D10
	case "DEV0019512211":
		return DeviceInformation_DEV0019512211
	case "DEV0019512311":
		return DeviceInformation_DEV0019512311
	case "DEV0072300110":
		return DeviceInformation_DEV0072300110
	case "DEV0019512111":
		return DeviceInformation_DEV0019512111
	case "DEV0019520D11":
		return DeviceInformation_DEV0019520D11
	case "DEV0019E30B12":
		return DeviceInformation_DEV0019E30B12
	case "DEV0019530812":
		return DeviceInformation_DEV0019530812
	case "DEV0019530912":
		return DeviceInformation_DEV0019530912
	case "DEV0019530612":
		return DeviceInformation_DEV0019530612
	case "DEV0019530711":
		return DeviceInformation_DEV0019530711
	case "DEV0019E30A11":
		return DeviceInformation_DEV0019E30A11
	case "DEV0019E20111":
		return DeviceInformation_DEV0019E20111
	case "DEV0019E20210":
		return DeviceInformation_DEV0019E20210
	case "DEV0076002101":
		return DeviceInformation_DEV0076002101
	case "DEV0019E30C11":
		return DeviceInformation_DEV0019E30C11
	case "DEV0019E11310":
		return DeviceInformation_DEV0019E11310
	case "DEV0019E11210":
		return DeviceInformation_DEV0019E11210
	case "DEV0019E30610":
		return DeviceInformation_DEV0019E30610
	case "DEV0019E30710":
		return DeviceInformation_DEV0019E30710
	case "DEV0019E30910":
		return DeviceInformation_DEV0019E30910
	case "DEV0019E30810":
		return DeviceInformation_DEV0019E30810
	case "DEV0019E25510":
		return DeviceInformation_DEV0019E25510
	case "DEV0019E20410":
		return DeviceInformation_DEV0019E20410
	case "DEV0019E20310":
		return DeviceInformation_DEV0019E20310
	case "DEV0076002001":
		return DeviceInformation_DEV0076002001
	case "DEV0019E25610":
		return DeviceInformation_DEV0019E25610
	case "DEV0019512010":
		return DeviceInformation_DEV0019512010
	case "DEV0019520C10":
		return DeviceInformation_DEV0019520C10
	case "DEV0019520710":
		return DeviceInformation_DEV0019520710
	case "DEV0019520210":
		return DeviceInformation_DEV0019520210
	case "DEV0019E25010":
		return DeviceInformation_DEV0019E25010
	case "DEV0019E25110":
		return DeviceInformation_DEV0019E25110
	case "DEV0019130710":
		return DeviceInformation_DEV0019130710
	case "DEV0019272050":
		return DeviceInformation_DEV0019272050
	case "DEV0019520910":
		return DeviceInformation_DEV0019520910
	case "DEV0076002A15":
		return DeviceInformation_DEV0076002A15
	case "DEV0019520A10":
		return DeviceInformation_DEV0019520A10
	case "DEV0019520B10":
		return DeviceInformation_DEV0019520B10
	case "DEV0019520412":
		return DeviceInformation_DEV0019520412
	case "DEV0019520812":
		return DeviceInformation_DEV0019520812
	case "DEV0019512510":
		return DeviceInformation_DEV0019512510
	case "DEV0019512410":
		return DeviceInformation_DEV0019512410
	case "DEV0019512610":
		return DeviceInformation_DEV0019512610
	case "DEV0019511711":
		return DeviceInformation_DEV0019511711
	case "DEV0019511811":
		return DeviceInformation_DEV0019511811
	case "DEV0019522212":
		return DeviceInformation_DEV0019522212
	case "DEV0076002815":
		return DeviceInformation_DEV0076002815
	case "DEV0019FF0716":
		return DeviceInformation_DEV0019FF0716
	case "DEV0019FF1420":
		return DeviceInformation_DEV0019FF1420
	case "DEV0019522112":
		return DeviceInformation_DEV0019522112
	case "DEV0019522011":
		return DeviceInformation_DEV0019522011
	case "DEV0019522311":
		return DeviceInformation_DEV0019522311
	case "DEV0019E12410":
		return DeviceInformation_DEV0019E12410
	case "DEV0019000311":
		return DeviceInformation_DEV0019000311
	case "DEV0019000411":
		return DeviceInformation_DEV0019000411
	case "DEV0019070210":
		return DeviceInformation_DEV0019070210
	case "DEV0019070E11":
		return DeviceInformation_DEV0019070E11
	case "DEV0076002215":
		return DeviceInformation_DEV0076002215
	case "DEV0019724010":
		return DeviceInformation_DEV0019724010
	case "DEV0019520610":
		return DeviceInformation_DEV0019520610
	case "DEV0019520510":
		return DeviceInformation_DEV0019520510
	case "DEV00FB101111":
		return DeviceInformation_DEV00FB101111
	case "DEV00FB103001":
		return DeviceInformation_DEV00FB103001
	case "DEV00FB104401":
		return DeviceInformation_DEV00FB104401
	case "DEV00FB124002":
		return DeviceInformation_DEV00FB124002
	case "DEV00FB104102":
		return DeviceInformation_DEV00FB104102
	case "DEV00FB104201":
		return DeviceInformation_DEV00FB104201
	case "DEV00FBF77603":
		return DeviceInformation_DEV00FBF77603
	case "DEV0076002B15":
		return DeviceInformation_DEV0076002B15
	case "DEV00FB104301":
		return DeviceInformation_DEV00FB104301
	case "DEV00FB104601":
		return DeviceInformation_DEV00FB104601
	case "DEV00FB104701":
		return DeviceInformation_DEV00FB104701
	case "DEV00FB105101":
		return DeviceInformation_DEV00FB105101
	case "DEV0103030110":
		return DeviceInformation_DEV0103030110
	case "DEV0103010113":
		return DeviceInformation_DEV0103010113
	case "DEV0103090110":
		return DeviceInformation_DEV0103090110
	case "DEV0103020111":
		return DeviceInformation_DEV0103020111
	case "DEV0103020112":
		return DeviceInformation_DEV0103020112
	case "DEV0103040110":
		return DeviceInformation_DEV0103040110
	case "DEV0076002715":
		return DeviceInformation_DEV0076002715
	case "DEV0103050111":
		return DeviceInformation_DEV0103050111
	case "DEV0107000301":
		return DeviceInformation_DEV0107000301
	case "DEV0107000101":
		return DeviceInformation_DEV0107000101
	case "DEV0107000201":
		return DeviceInformation_DEV0107000201
	case "DEV0107020801":
		return DeviceInformation_DEV0107020801
	case "DEV0107020401":
		return DeviceInformation_DEV0107020401
	case "DEV0107020001":
		return DeviceInformation_DEV0107020001
	case "DEV010701F801":
		return DeviceInformation_DEV010701F801
	case "DEV010701FC01":
		return DeviceInformation_DEV010701FC01
	case "DEV0107020C01":
		return DeviceInformation_DEV0107020C01
	case "DEV0076002315":
		return DeviceInformation_DEV0076002315
	case "DEV010F100801":
		return DeviceInformation_DEV010F100801
	case "DEV010F100601":
		return DeviceInformation_DEV010F100601
	case "DEV010F100401":
		return DeviceInformation_DEV010F100401
	case "DEV010F030601":
		return DeviceInformation_DEV010F030601
	case "DEV010F010301":
		return DeviceInformation_DEV010F010301
	case "DEV010F010101":
		return DeviceInformation_DEV010F010101
	case "DEV010F010201":
		return DeviceInformation_DEV010F010201
	case "DEV010F000302":
		return DeviceInformation_DEV010F000302
	case "DEV010F000402":
		return DeviceInformation_DEV010F000402
	case "DEV010F000102":
		return DeviceInformation_DEV010F000102
	case "DEV0064182310":
		return DeviceInformation_DEV0064182310
	case "DEV0076002415":
		return DeviceInformation_DEV0076002415
	case "DEV011EBB8211":
		return DeviceInformation_DEV011EBB8211
	case "DEV011E108111":
		return DeviceInformation_DEV011E108111
	case "DEV0123010010":
		return DeviceInformation_DEV0123010010
	case "DEV001E478010":
		return DeviceInformation_DEV001E478010
	case "DEV001E706611":
		return DeviceInformation_DEV001E706611
	case "DEV001E706811":
		return DeviceInformation_DEV001E706811
	case "DEV001E473012":
		return DeviceInformation_DEV001E473012
	case "DEV001E20A011":
		return DeviceInformation_DEV001E20A011
	case "DEV001E209011":
		return DeviceInformation_DEV001E209011
	case "DEV001E209811":
		return DeviceInformation_DEV001E209811
	case "DEV0076002615":
		return DeviceInformation_DEV0076002615
	case "DEV001E208811":
		return DeviceInformation_DEV001E208811
	case "DEV001E208011":
		return DeviceInformation_DEV001E208011
	case "DEV001E207821":
		return DeviceInformation_DEV001E207821
	case "DEV001E20CA12":
		return DeviceInformation_DEV001E20CA12
	case "DEV001E20B312":
		return DeviceInformation_DEV001E20B312
	case "DEV001E20B012":
		return DeviceInformation_DEV001E20B012
	case "DEV001E302612":
		return DeviceInformation_DEV001E302612
	case "DEV001E302312":
		return DeviceInformation_DEV001E302312
	case "DEV001E302012":
		return DeviceInformation_DEV001E302012
	case "DEV001E20A811":
		return DeviceInformation_DEV001E20A811
	case "DEV0076002515":
		return DeviceInformation_DEV0076002515
	case "DEV001E20C412":
		return DeviceInformation_DEV001E20C412
	case "DEV001E20C712":
		return DeviceInformation_DEV001E20C712
	case "DEV001E20AD12":
		return DeviceInformation_DEV001E20AD12
	case "DEV001E443720":
		return DeviceInformation_DEV001E443720
	case "DEV001E441821":
		return DeviceInformation_DEV001E441821
	case "DEV001E443810":
		return DeviceInformation_DEV001E443810
	case "DEV001E140C12":
		return DeviceInformation_DEV001E140C12
	case "DEV001E471611":
		return DeviceInformation_DEV001E471611
	case "DEV001E479024":
		return DeviceInformation_DEV001E479024
	case "DEV001E471A11":
		return DeviceInformation_DEV001E471A11
	case "DEV0076000201":
		return DeviceInformation_DEV0076000201
	case "DEV001E477A10":
		return DeviceInformation_DEV001E477A10
	case "DEV001E470A11":
		return DeviceInformation_DEV001E470A11
	case "DEV001E480B11":
		return DeviceInformation_DEV001E480B11
	case "DEV001E487B10":
		return DeviceInformation_DEV001E487B10
	case "DEV001E440411":
		return DeviceInformation_DEV001E440411
	case "DEV001E447211":
		return DeviceInformation_DEV001E447211
	case "DEV0142010010":
		return DeviceInformation_DEV0142010010
	case "DEV0142010011":
		return DeviceInformation_DEV0142010011
	case "DEV017A130401":
		return DeviceInformation_DEV017A130401
	case "DEV017A130201":
		return DeviceInformation_DEV017A130201
	case "DEV0076000101":
		return DeviceInformation_DEV0076000101
	case "DEV017A130801":
		return DeviceInformation_DEV017A130801
	case "DEV017A130601":
		return DeviceInformation_DEV017A130601
	case "DEV017A300102":
		return DeviceInformation_DEV017A300102
	case "DEV0193323C11":
		return DeviceInformation_DEV0193323C11
	case "DEV0198101110":
		return DeviceInformation_DEV0198101110
	case "DEV01C4030110":
		return DeviceInformation_DEV01C4030110
	case "DEV01C4030210":
		return DeviceInformation_DEV01C4030210
	case "DEV01C4021210":
		return DeviceInformation_DEV01C4021210
	case "DEV01C4001010":
		return DeviceInformation_DEV01C4001010
	case "DEV01C4020610":
		return DeviceInformation_DEV01C4020610
	case "DEV0076000301":
		return DeviceInformation_DEV0076000301
	case "DEV01C4020910":
		return DeviceInformation_DEV01C4020910
	case "DEV01C4020810":
		return DeviceInformation_DEV01C4020810
	case "DEV01C4010710":
		return DeviceInformation_DEV01C4010710
	case "DEV01C4050210":
		return DeviceInformation_DEV01C4050210
	case "DEV01C4010810":
		return DeviceInformation_DEV01C4010810
	case "DEV01C4020510":
		return DeviceInformation_DEV01C4020510
	case "DEV01C4040110":
		return DeviceInformation_DEV01C4040110
	case "DEV01C4040310":
		return DeviceInformation_DEV01C4040310
	case "DEV01C4040210":
		return DeviceInformation_DEV01C4040210
	case "DEV01C4101210":
		return DeviceInformation_DEV01C4101210
	case "DEV0076000401":
		return DeviceInformation_DEV0076000401
	case "DEV003D020109":
		return DeviceInformation_DEV003D020109
	case "DEV01DB000301":
		return DeviceInformation_DEV01DB000301
	case "DEV01DB000201":
		return DeviceInformation_DEV01DB000201
	case "DEV01DB000401":
		return DeviceInformation_DEV01DB000401
	case "DEV01DB000801":
		return DeviceInformation_DEV01DB000801
	case "DEV01DB001201":
		return DeviceInformation_DEV01DB001201
	case "DEV009A000400":
		return DeviceInformation_DEV009A000400
	case "DEV009A100400":
		return DeviceInformation_DEV009A100400
	case "DEV009A200C00":
		return DeviceInformation_DEV009A200C00
	case "DEV009A200E00":
		return DeviceInformation_DEV009A200E00
	case "DEV0076002903":
		return DeviceInformation_DEV0076002903
	case "DEV009A000201":
		return DeviceInformation_DEV009A000201
	case "DEV009A000300":
		return DeviceInformation_DEV009A000300
	case "DEV009A00C000":
		return DeviceInformation_DEV009A00C000
	case "DEV009A00B000":
		return DeviceInformation_DEV009A00B000
	case "DEV009A00C002":
		return DeviceInformation_DEV009A00C002
	case "DEV009A200100":
		return DeviceInformation_DEV009A200100
	case "DEV004E400010":
		return DeviceInformation_DEV004E400010
	case "DEV004E030031":
		return DeviceInformation_DEV004E030031
	case "DEV012B010110":
		return DeviceInformation_DEV012B010110
	case "DEV01F6E0E110":
		return DeviceInformation_DEV01F6E0E110
	case "DEV0076002901":
		return DeviceInformation_DEV0076002901
	case "DEV0088100010":
		return DeviceInformation_DEV0088100010
	case "DEV0088100210":
		return DeviceInformation_DEV0088100210
	case "DEV0088100110":
		return DeviceInformation_DEV0088100110
	case "DEV0088110010":
		return DeviceInformation_DEV0088110010
	case "DEV0088120412":
		return DeviceInformation_DEV0088120412
	case "DEV0088120113":
		return DeviceInformation_DEV0088120113
	case "DEV011A4B5201":
		return DeviceInformation_DEV011A4B5201
	case "DEV008B020301":
		return DeviceInformation_DEV008B020301
	case "DEV008B010610":
		return DeviceInformation_DEV008B010610
	case "DEV008B030110":
		return DeviceInformation_DEV008B030110
	case "DEV007600E503":
		return DeviceInformation_DEV007600E503
	case "DEV008B030310":
		return DeviceInformation_DEV008B030310
	case "DEV008B030210":
		return DeviceInformation_DEV008B030210
	case "DEV008B031512":
		return DeviceInformation_DEV008B031512
	case "DEV008B031412":
		return DeviceInformation_DEV008B031412
	case "DEV008B031312":
		return DeviceInformation_DEV008B031312
	case "DEV008B031212":
		return DeviceInformation_DEV008B031212
	case "DEV008B031112":
		return DeviceInformation_DEV008B031112
	case "DEV008B031012":
		return DeviceInformation_DEV008B031012
	case "DEV008B030510":
		return DeviceInformation_DEV008B030510
	case "DEV008B030410":
		return DeviceInformation_DEV008B030410
	case "DEV0064705C01":
		return DeviceInformation_DEV0064705C01
	case "DEV0076004002":
		return DeviceInformation_DEV0076004002
	case "DEV008B020310":
		return DeviceInformation_DEV008B020310
	case "DEV008B020210":
		return DeviceInformation_DEV008B020210
	case "DEV008B020110":
		return DeviceInformation_DEV008B020110
	case "DEV008B010110":
		return DeviceInformation_DEV008B010110
	case "DEV008B010210":
		return DeviceInformation_DEV008B010210
	case "DEV008B010310":
		return DeviceInformation_DEV008B010310
	case "DEV008B010410":
		return DeviceInformation_DEV008B010410
	case "DEV008B040110":
		return DeviceInformation_DEV008B040110
	case "DEV008B040210":
		return DeviceInformation_DEV008B040210
	case "DEV008B010910":
		return DeviceInformation_DEV008B010910
	case "DEV0076004003":
		return DeviceInformation_DEV0076004003
	case "DEV008B010710":
		return DeviceInformation_DEV008B010710
	case "DEV008B010810":
		return DeviceInformation_DEV008B010810
	case "DEV008B041111":
		return DeviceInformation_DEV008B041111
	case "DEV008B041211":
		return DeviceInformation_DEV008B041211
	case "DEV008B041311":
		return DeviceInformation_DEV008B041311
	case "DEV00A600020A":
		return DeviceInformation_DEV00A600020A
	case "DEV00A6000B10":
		return DeviceInformation_DEV00A6000B10
	case "DEV00A6000300":
		return DeviceInformation_DEV00A6000300
	case "DEV00A6000705":
		return DeviceInformation_DEV00A6000705
	case "DEV00A6000605":
		return DeviceInformation_DEV00A6000605
	case "DEV0076003402":
		return DeviceInformation_DEV0076003402
	case "DEV00A6000500":
		return DeviceInformation_DEV00A6000500
	case "DEV00A6000C10":
		return DeviceInformation_DEV00A6000C10
	case "DEV002CA01811":
		return DeviceInformation_DEV002CA01811
	case "DEV002CA07033":
		return DeviceInformation_DEV002CA07033
	case "DEV002C555020":
		return DeviceInformation_DEV002C555020
	case "DEV002C556421":
		return DeviceInformation_DEV002C556421
	case "DEV002C05F211":
		return DeviceInformation_DEV002C05F211
	case "DEV002C05F411":
		return DeviceInformation_DEV002C05F411
	case "DEV002C05E613":
		return DeviceInformation_DEV002C05E613
	case "DEV002CA07914":
		return DeviceInformation_DEV002CA07914
	case "DEV0076003401":
		return DeviceInformation_DEV0076003401
	case "DEV002C060A13":
		return DeviceInformation_DEV002C060A13
	case "DEV002C3A0212":
		return DeviceInformation_DEV002C3A0212
	case "DEV002C060210":
		return DeviceInformation_DEV002C060210
	case "DEV002CA00213":
		return DeviceInformation_DEV002CA00213
	case "DEV002CA0A611":
		return DeviceInformation_DEV002CA0A611
	case "DEV002CA07B11":
		return DeviceInformation_DEV002CA07B11
	case "DEV002C063210":
		return DeviceInformation_DEV002C063210
	case "DEV002C063110":
		return DeviceInformation_DEV002C063110
	case "DEV002C062E10":
		return DeviceInformation_DEV002C062E10
	case "DEV002C062C10":
		return DeviceInformation_DEV002C062C10
	case "DEV007600E908":
		return DeviceInformation_DEV007600E908
	case "DEV002C062D10":
		return DeviceInformation_DEV002C062D10
	case "DEV002C063310":
		return DeviceInformation_DEV002C063310
	case "DEV002C05EB10":
		return DeviceInformation_DEV002C05EB10
	case "DEV002C05F110":
		return DeviceInformation_DEV002C05F110
	case "DEV002C0B8830":
		return DeviceInformation_DEV002C0B8830
	case "DEV00A0B07101":
		return DeviceInformation_DEV00A0B07101
	case "DEV00A0B07001":
		return DeviceInformation_DEV00A0B07001
	case "DEV00A0B07203":
		return DeviceInformation_DEV00A0B07203
	case "DEV00A0B02101":
		return DeviceInformation_DEV00A0B02101
	case "DEV00A0B02401":
		return DeviceInformation_DEV00A0B02401
	case "DEV007600E907":
		return DeviceInformation_DEV007600E907
	case "DEV00A0B02301":
		return DeviceInformation_DEV00A0B02301
	case "DEV00A0B02601":
		return DeviceInformation_DEV00A0B02601
	case "DEV00A0B02201":
		return DeviceInformation_DEV00A0B02201
	case "DEV00A0B01902":
		return DeviceInformation_DEV00A0B01902
	case "DEV0004147112":
		return DeviceInformation_DEV0004147112
	case "DEV000410A411":
		return DeviceInformation_DEV000410A411
	case "DEV0004109911":
		return DeviceInformation_DEV0004109911
	case "DEV0004109912":
		return DeviceInformation_DEV0004109912
	case "DEV0004109913":
		return DeviceInformation_DEV0004109913
	case "DEV0004109914":
		return DeviceInformation_DEV0004109914
	case "DEV000C181710":
		return DeviceInformation_DEV000C181710
	case "DEV000410A211":
		return DeviceInformation_DEV000410A211
	case "DEV000410FC12":
		return DeviceInformation_DEV000410FC12
	case "DEV000410FD12":
		return DeviceInformation_DEV000410FD12
	case "DEV000410B212":
		return DeviceInformation_DEV000410B212
	case "DEV0004110B11":
		return DeviceInformation_DEV0004110B11
	case "DEV0004110711":
		return DeviceInformation_DEV0004110711
	case "DEV000410B213":
		return DeviceInformation_DEV000410B213
	case "DEV0004109811":
		return DeviceInformation_DEV0004109811
	case "DEV0004109812":
		return DeviceInformation_DEV0004109812
	case "DEV0004109813":
		return DeviceInformation_DEV0004109813
	case "DEV000C130510":
		return DeviceInformation_DEV000C130510
	case "DEV0004109814":
		return DeviceInformation_DEV0004109814
	case "DEV000410A011":
		return DeviceInformation_DEV000410A011
	case "DEV000410A111":
		return DeviceInformation_DEV000410A111
	case "DEV000410FA12":
		return DeviceInformation_DEV000410FA12
	case "DEV000410FB12":
		return DeviceInformation_DEV000410FB12
	case "DEV000410B112":
		return DeviceInformation_DEV000410B112
	case "DEV0004110A11":
		return DeviceInformation_DEV0004110A11
	case "DEV0004110611":
		return DeviceInformation_DEV0004110611
	case "DEV000410B113":
		return DeviceInformation_DEV000410B113
	case "DEV0004109A11":
		return DeviceInformation_DEV0004109A11
	case "DEV000C130610":
		return DeviceInformation_DEV000C130610
	case "DEV0004109A12":
		return DeviceInformation_DEV0004109A12
	case "DEV0004109A13":
		return DeviceInformation_DEV0004109A13
	case "DEV0004109A14":
		return DeviceInformation_DEV0004109A14
	case "DEV000410A311":
		return DeviceInformation_DEV000410A311
	case "DEV000410B312":
		return DeviceInformation_DEV000410B312
	case "DEV0004110C11":
		return DeviceInformation_DEV0004110C11
	case "DEV0004110811":
		return DeviceInformation_DEV0004110811
	case "DEV000410B313":
		return DeviceInformation_DEV000410B313
	case "DEV0004109B11":
		return DeviceInformation_DEV0004109B11
	case "DEV0004109B12":
		return DeviceInformation_DEV0004109B12
	case "DEV000C133610":
		return DeviceInformation_DEV000C133610
	case "DEV0004109B13":
		return DeviceInformation_DEV0004109B13
	case "DEV0004109B14":
		return DeviceInformation_DEV0004109B14
	case "DEV000410A511":
		return DeviceInformation_DEV000410A511
	case "DEV000410B412":
		return DeviceInformation_DEV000410B412
	case "DEV0004110D11":
		return DeviceInformation_DEV0004110D11
	case "DEV0004110911":
		return DeviceInformation_DEV0004110911
	case "DEV000410B413":
		return DeviceInformation_DEV000410B413
	case "DEV0004109C11":
		return DeviceInformation_DEV0004109C11
	case "DEV0004109C12":
		return DeviceInformation_DEV0004109C12
	case "DEV0004109C13":
		return DeviceInformation_DEV0004109C13
	}
	return 0
}

func CastDeviceInformation(structType interface{}) DeviceInformation {
	castFunc := func(typ interface{}) DeviceInformation {
		if sDeviceInformation, ok := typ.(DeviceInformation); ok {
			return sDeviceInformation
		}
		return 0
	}
	return castFunc(structType)
}

func (m DeviceInformation) LengthInBits() uint16 {
	return 16
}

func (m DeviceInformation) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func DeviceInformationParse(io *utils.ReadBuffer) (DeviceInformation, error) {
	val, err := io.ReadUint16(16)
	if err != nil {
		return 0, nil
	}
	return DeviceInformationByValue(val), nil
}

func (e DeviceInformation) Serialize(io utils.WriteBuffer) error {
	err := io.WriteUint16(16, uint16(e))
	return err
}

func (e DeviceInformation) String() string {
	switch e {
	case DeviceInformation_DEV0001914201:
		return "DEV0001914201"
	case DeviceInformation_DEV0064181910:
		return "DEV0064181910"
	case DeviceInformation_DEV000C133410:
		return "DEV000C133410"
	case DeviceInformation_DEV0004109C14:
		return "DEV0004109C14"
	case DeviceInformation_DEV000410A611:
		return "DEV000410A611"
	case DeviceInformation_DEV0004146B13:
		return "DEV0004146B13"
	case DeviceInformation_DEV0004147211:
		return "DEV0004147211"
	case DeviceInformation_DEV000410FE12:
		return "DEV000410FE12"
	case DeviceInformation_DEV0004209016:
		return "DEV0004209016"
	case DeviceInformation_DEV000420A011:
		return "DEV000420A011"
	case DeviceInformation_DEV0004209011:
		return "DEV0004209011"
	case DeviceInformation_DEV000420CA11:
		return "DEV000420CA11"
	case DeviceInformation_DEV0004208012:
		return "DEV0004208012"
	case DeviceInformation_DEV000C133310:
		return "DEV000C133310"
	case DeviceInformation_DEV0004207812:
		return "DEV0004207812"
	case DeviceInformation_DEV000420BA11:
		return "DEV000420BA11"
	case DeviceInformation_DEV000420B311:
		return "DEV000420B311"
	case DeviceInformation_DEV0004209811:
		return "DEV0004209811"
	case DeviceInformation_DEV0004208811:
		return "DEV0004208811"
	case DeviceInformation_DEV0004B00812:
		return "DEV0004B00812"
	case DeviceInformation_DEV0004302613:
		return "DEV0004302613"
	case DeviceInformation_DEV0004302313:
		return "DEV0004302313"
	case DeviceInformation_DEV0004302013:
		return "DEV0004302013"
	case DeviceInformation_DEV0004302B12:
		return "DEV0004302B12"
	case DeviceInformation_DEV000C133611:
		return "DEV000C133611"
	case DeviceInformation_DEV0004706811:
		return "DEV0004706811"
	case DeviceInformation_DEV0004705D11:
		return "DEV0004705D11"
	case DeviceInformation_DEV0004705C11:
		return "DEV0004705C11"
	case DeviceInformation_DEV0004B00713:
		return "DEV0004B00713"
	case DeviceInformation_DEV0004B00A01:
		return "DEV0004B00A01"
	case DeviceInformation_DEV0004706611:
		return "DEV0004706611"
	case DeviceInformation_DEV0004C01410:
		return "DEV0004C01410"
	case DeviceInformation_DEV0004C00102:
		return "DEV0004C00102"
	case DeviceInformation_DEV0004705E11:
		return "DEV0004705E11"
	case DeviceInformation_DEV0004706211:
		return "DEV0004706211"
	case DeviceInformation_DEV000C133510:
		return "DEV000C133510"
	case DeviceInformation_DEV0004706411:
		return "DEV0004706411"
	case DeviceInformation_DEV0004706412:
		return "DEV0004706412"
	case DeviceInformation_DEV000420C011:
		return "DEV000420C011"
	case DeviceInformation_DEV000420B011:
		return "DEV000420B011"
	case DeviceInformation_DEV0004B00911:
		return "DEV0004B00911"
	case DeviceInformation_DEV0004A01211:
		return "DEV0004A01211"
	case DeviceInformation_DEV0004A01112:
		return "DEV0004A01112"
	case DeviceInformation_DEV0004A01111:
		return "DEV0004A01111"
	case DeviceInformation_DEV0004A01212:
		return "DEV0004A01212"
	case DeviceInformation_DEV0004A03312:
		return "DEV0004A03312"
	case DeviceInformation_DEV000C130710:
		return "DEV000C130710"
	case DeviceInformation_DEV0004A03212:
		return "DEV0004A03212"
	case DeviceInformation_DEV0004A01113:
		return "DEV0004A01113"
	case DeviceInformation_DEV0004A01711:
		return "DEV0004A01711"
	case DeviceInformation_DEV000420C711:
		return "DEV000420C711"
	case DeviceInformation_DEV000420BD11:
		return "DEV000420BD11"
	case DeviceInformation_DEV000420C411:
		return "DEV000420C411"
	case DeviceInformation_DEV000420A812:
		return "DEV000420A812"
	case DeviceInformation_DEV000420CD11:
		return "DEV000420CD11"
	case DeviceInformation_DEV000420AD11:
		return "DEV000420AD11"
	case DeviceInformation_DEV000420B611:
		return "DEV000420B611"
	case DeviceInformation_DEV000C760210:
		return "DEV000C760210"
	case DeviceInformation_DEV000420A811:
		return "DEV000420A811"
	case DeviceInformation_DEV0004501311:
		return "DEV0004501311"
	case DeviceInformation_DEV0004501411:
		return "DEV0004501411"
	case DeviceInformation_DEV0004B01002:
		return "DEV0004B01002"
	case DeviceInformation_DEV0006D00610:
		return "DEV0006D00610"
	case DeviceInformation_DEV0006D01510:
		return "DEV0006D01510"
	case DeviceInformation_DEV0006D00110:
		return "DEV0006D00110"
	case DeviceInformation_DEV0006D00310:
		return "DEV0006D00310"
	case DeviceInformation_DEV0006D03210:
		return "DEV0006D03210"
	case DeviceInformation_DEV0006D03310:
		return "DEV0006D03310"
	case DeviceInformation_DEV000C1BD610:
		return "DEV000C1BD610"
	case DeviceInformation_DEV0006D02E20:
		return "DEV0006D02E20"
	case DeviceInformation_DEV0006D02F20:
		return "DEV0006D02F20"
	case DeviceInformation_DEV0006D03020:
		return "DEV0006D03020"
	case DeviceInformation_DEV0006D03120:
		return "DEV0006D03120"
	case DeviceInformation_DEV0006D02110:
		return "DEV0006D02110"
	case DeviceInformation_DEV0006D00010:
		return "DEV0006D00010"
	case DeviceInformation_DEV0006D01810:
		return "DEV0006D01810"
	case DeviceInformation_DEV0006D00910:
		return "DEV0006D00910"
	case DeviceInformation_DEV0006D01110:
		return "DEV0006D01110"
	case DeviceInformation_DEV0006D03510:
		return "DEV0006D03510"
	case DeviceInformation_DEV000C181610:
		return "DEV000C181610"
	case DeviceInformation_DEV0006D03410:
		return "DEV0006D03410"
	case DeviceInformation_DEV0006D02410:
		return "DEV0006D02410"
	case DeviceInformation_DEV0006D02510:
		return "DEV0006D02510"
	case DeviceInformation_DEV0006D00810:
		return "DEV0006D00810"
	case DeviceInformation_DEV0006D00710:
		return "DEV0006D00710"
	case DeviceInformation_DEV0006D01310:
		return "DEV0006D01310"
	case DeviceInformation_DEV0006D01410:
		return "DEV0006D01410"
	case DeviceInformation_DEV0006D00210:
		return "DEV0006D00210"
	case DeviceInformation_DEV0006D00510:
		return "DEV0006D00510"
	case DeviceInformation_DEV0006D00410:
		return "DEV0006D00410"
	case DeviceInformation_DEV000C648B10:
		return "DEV000C648B10"
	case DeviceInformation_DEV0006D02210:
		return "DEV0006D02210"
	case DeviceInformation_DEV0006D02310:
		return "DEV0006D02310"
	case DeviceInformation_DEV0006D01710:
		return "DEV0006D01710"
	case DeviceInformation_DEV0006D01610:
		return "DEV0006D01610"
	case DeviceInformation_DEV0006D01010:
		return "DEV0006D01010"
	case DeviceInformation_DEV0006D01210:
		return "DEV0006D01210"
	case DeviceInformation_DEV0006D04820:
		return "DEV0006D04820"
	case DeviceInformation_DEV0006D04C11:
		return "DEV0006D04C11"
	case DeviceInformation_DEV0006D05610:
		return "DEV0006D05610"
	case DeviceInformation_DEV0006D02910:
		return "DEV0006D02910"
	case DeviceInformation_DEV000C480611:
		return "DEV000C480611"
	case DeviceInformation_DEV0006D02A10:
		return "DEV0006D02A10"
	case DeviceInformation_DEV0006D02B10:
		return "DEV0006D02B10"
	case DeviceInformation_DEV0006D02C10:
		return "DEV0006D02C10"
	case DeviceInformation_DEV0006D02810:
		return "DEV0006D02810"
	case DeviceInformation_DEV0006D02610:
		return "DEV0006D02610"
	case DeviceInformation_DEV0006D02710:
		return "DEV0006D02710"
	case DeviceInformation_DEV0006D03610:
		return "DEV0006D03610"
	case DeviceInformation_DEV0006D03710:
		return "DEV0006D03710"
	case DeviceInformation_DEV0006D02D11:
		return "DEV0006D02D11"
	case DeviceInformation_DEV0006D03C10:
		return "DEV0006D03C10"
	case DeviceInformation_DEV0064181810:
		return "DEV0064181810"
	case DeviceInformation_DEV000C482011:
		return "DEV000C482011"
	case DeviceInformation_DEV0006D03B10:
		return "DEV0006D03B10"
	case DeviceInformation_DEV0006D03910:
		return "DEV0006D03910"
	case DeviceInformation_DEV0006D03A10:
		return "DEV0006D03A10"
	case DeviceInformation_DEV0006D03D11:
		return "DEV0006D03D11"
	case DeviceInformation_DEV0006D03E10:
		return "DEV0006D03E10"
	case DeviceInformation_DEV0006C00102:
		return "DEV0006C00102"
	case DeviceInformation_DEV0006E05611:
		return "DEV0006E05611"
	case DeviceInformation_DEV0006E05212:
		return "DEV0006E05212"
	case DeviceInformation_DEV000620B011:
		return "DEV000620B011"
	case DeviceInformation_DEV000620B311:
		return "DEV000620B311"
	case DeviceInformation_DEV000C724010:
		return "DEV000C724010"
	case DeviceInformation_DEV000620C011:
		return "DEV000620C011"
	case DeviceInformation_DEV000620BA11:
		return "DEV000620BA11"
	case DeviceInformation_DEV0006705C11:
		return "DEV0006705C11"
	case DeviceInformation_DEV0006705D11:
		return "DEV0006705D11"
	case DeviceInformation_DEV0006E07710:
		return "DEV0006E07710"
	case DeviceInformation_DEV0006E07712:
		return "DEV0006E07712"
	case DeviceInformation_DEV0006706210:
		return "DEV0006706210"
	case DeviceInformation_DEV0006302611:
		return "DEV0006302611"
	case DeviceInformation_DEV0006302612:
		return "DEV0006302612"
	case DeviceInformation_DEV0006E00810:
		return "DEV0006E00810"
	case DeviceInformation_DEV000C570211:
		return "DEV000C570211"
	case DeviceInformation_DEV0006E01F01:
		return "DEV0006E01F01"
	case DeviceInformation_DEV0006302311:
		return "DEV0006302311"
	case DeviceInformation_DEV0006302312:
		return "DEV0006302312"
	case DeviceInformation_DEV0006E00910:
		return "DEV0006E00910"
	case DeviceInformation_DEV0006E02001:
		return "DEV0006E02001"
	case DeviceInformation_DEV0006302011:
		return "DEV0006302011"
	case DeviceInformation_DEV0006302012:
		return "DEV0006302012"
	case DeviceInformation_DEV0006C00C13:
		return "DEV0006C00C13"
	case DeviceInformation_DEV0006E00811:
		return "DEV0006E00811"
	case DeviceInformation_DEV0006E00911:
		return "DEV0006E00911"
	case DeviceInformation_DEV000C570310:
		return "DEV000C570310"
	case DeviceInformation_DEV0006E01F20:
		return "DEV0006E01F20"
	case DeviceInformation_DEV0006E03410:
		return "DEV0006E03410"
	case DeviceInformation_DEV0006E03110:
		return "DEV0006E03110"
	case DeviceInformation_DEV0006E0A210:
		return "DEV0006E0A210"
	case DeviceInformation_DEV0006E0CE10:
		return "DEV0006E0CE10"
	case DeviceInformation_DEV0006E0A111:
		return "DEV0006E0A111"
	case DeviceInformation_DEV0006E0CD11:
		return "DEV0006E0CD11"
	case DeviceInformation_DEV0006E02020:
		return "DEV0006E02020"
	case DeviceInformation_DEV0006E02D11:
		return "DEV0006E02D11"
	case DeviceInformation_DEV0006E03011:
		return "DEV0006E03011"
	case DeviceInformation_DEV000C570411:
		return "DEV000C570411"
	case DeviceInformation_DEV0006E0C110:
		return "DEV0006E0C110"
	case DeviceInformation_DEV0006E0C510:
		return "DEV0006E0C510"
	case DeviceInformation_DEV0006B00A01:
		return "DEV0006B00A01"
	case DeviceInformation_DEV0006B00602:
		return "DEV0006B00602"
	case DeviceInformation_DEV0006E0C410:
		return "DEV0006E0C410"
	case DeviceInformation_DEV0006E0C312:
		return "DEV0006E0C312"
	case DeviceInformation_DEV0006E0C210:
		return "DEV0006E0C210"
	case DeviceInformation_DEV0006209016:
		return "DEV0006209016"
	case DeviceInformation_DEV0006E01A01:
		return "DEV0006E01A01"
	case DeviceInformation_DEV0006E09910:
		return "DEV0006E09910"
	case DeviceInformation_DEV000C570110:
		return "DEV000C570110"
	case DeviceInformation_DEV0006E03710:
		return "DEV0006E03710"
	case DeviceInformation_DEV0006209011:
		return "DEV0006209011"
	case DeviceInformation_DEV000620A011:
		return "DEV000620A011"
	case DeviceInformation_DEV0006E02410:
		return "DEV0006E02410"
	case DeviceInformation_DEV0006E02301:
		return "DEV0006E02301"
	case DeviceInformation_DEV0006E02510:
		return "DEV0006E02510"
	case DeviceInformation_DEV0006E01B01:
		return "DEV0006E01B01"
	case DeviceInformation_DEV0006E01C01:
		return "DEV0006E01C01"
	case DeviceInformation_DEV0006E01D01:
		return "DEV0006E01D01"
	case DeviceInformation_DEV0006E01E01:
		return "DEV0006E01E01"
	case DeviceInformation_DEV000C570011:
		return "DEV000C570011"
	case DeviceInformation_DEV0006207812:
		return "DEV0006207812"
	case DeviceInformation_DEV0006B00811:
		return "DEV0006B00811"
	case DeviceInformation_DEV0006E01001:
		return "DEV0006E01001"
	case DeviceInformation_DEV0006E03610:
		return "DEV0006E03610"
	case DeviceInformation_DEV0006E09810:
		return "DEV0006E09810"
	case DeviceInformation_DEV0006208811:
		return "DEV0006208811"
	case DeviceInformation_DEV0006209811:
		return "DEV0006209811"
	case DeviceInformation_DEV0006E02610:
		return "DEV0006E02610"
	case DeviceInformation_DEV0006E02710:
		return "DEV0006E02710"
	case DeviceInformation_DEV0006E02A10:
		return "DEV0006E02A10"
	case DeviceInformation_DEV000C20BD11:
		return "DEV000C20BD11"
	case DeviceInformation_DEV0006E02B10:
		return "DEV0006E02B10"
	case DeviceInformation_DEV0006E00C10:
		return "DEV0006E00C10"
	case DeviceInformation_DEV0006010110:
		return "DEV0006010110"
	case DeviceInformation_DEV0006010210:
		return "DEV0006010210"
	case DeviceInformation_DEV0006E00B10:
		return "DEV0006E00B10"
	case DeviceInformation_DEV0006E09C10:
		return "DEV0006E09C10"
	case DeviceInformation_DEV0006E09B10:
		return "DEV0006E09B10"
	case DeviceInformation_DEV0006E03510:
		return "DEV0006E03510"
	case DeviceInformation_DEV0006FF1B11:
		return "DEV0006FF1B11"
	case DeviceInformation_DEV0006E0CF10:
		return "DEV0006E0CF10"
	case DeviceInformation_DEV000C20BA11:
		return "DEV000C20BA11"
	case DeviceInformation_DEV000620A812:
		return "DEV000620A812"
	case DeviceInformation_DEV000620CD11:
		return "DEV000620CD11"
	case DeviceInformation_DEV0006E00E01:
		return "DEV0006E00E01"
	case DeviceInformation_DEV0006E02201:
		return "DEV0006E02201"
	case DeviceInformation_DEV000620AD11:
		return "DEV000620AD11"
	case DeviceInformation_DEV0006E00F01:
		return "DEV0006E00F01"
	case DeviceInformation_DEV0006E02101:
		return "DEV0006E02101"
	case DeviceInformation_DEV000620BD11:
		return "DEV000620BD11"
	case DeviceInformation_DEV0006E00D01:
		return "DEV0006E00D01"
	case DeviceInformation_DEV0006E03910:
		return "DEV0006E03910"
	case DeviceInformation_DEV000C760110:
		return "DEV000C760110"
	case DeviceInformation_DEV0006E02810:
		return "DEV0006E02810"
	case DeviceInformation_DEV0006E02910:
		return "DEV0006E02910"
	case DeviceInformation_DEV0006E02C10:
		return "DEV0006E02C10"
	case DeviceInformation_DEV0006C00403:
		return "DEV0006C00403"
	case DeviceInformation_DEV0006590101:
		return "DEV0006590101"
	case DeviceInformation_DEV0006E0CC11:
		return "DEV0006E0CC11"
	case DeviceInformation_DEV0006E09A10:
		return "DEV0006E09A10"
	case DeviceInformation_DEV0006E03811:
		return "DEV0006E03811"
	case DeviceInformation_DEV0006E0C710:
		return "DEV0006E0C710"
	case DeviceInformation_DEV0006E0C610:
		return "DEV0006E0C610"
	case DeviceInformation_DEV0064181710:
		return "DEV0064181710"
	case DeviceInformation_DEV000C705C01:
		return "DEV000C705C01"
	case DeviceInformation_DEV0006E05A10:
		return "DEV0006E05A10"
	case DeviceInformation_DEV0048493A1C:
		return "DEV0048493A1C"
	case DeviceInformation_DEV0048494712:
		return "DEV0048494712"
	case DeviceInformation_DEV0048494810:
		return "DEV0048494810"
	case DeviceInformation_DEV0048855A10:
		return "DEV0048855A10"
	case DeviceInformation_DEV0048855B10:
		return "DEV0048855B10"
	case DeviceInformation_DEV0048A05713:
		return "DEV0048A05713"
	case DeviceInformation_DEV0048494414:
		return "DEV0048494414"
	case DeviceInformation_DEV0048824A11:
		return "DEV0048824A11"
	case DeviceInformation_DEV0048824A12:
		return "DEV0048824A12"
	case DeviceInformation_DEV000CFF2112:
		return "DEV000CFF2112"
	case DeviceInformation_DEV0048770A10:
		return "DEV0048770A10"
	case DeviceInformation_DEV0048494311:
		return "DEV0048494311"
	case DeviceInformation_DEV0048494513:
		return "DEV0048494513"
	case DeviceInformation_DEV0048494012:
		return "DEV0048494012"
	case DeviceInformation_DEV0048494111:
		return "DEV0048494111"
	case DeviceInformation_DEV0048494210:
		return "DEV0048494210"
	case DeviceInformation_DEV0048494610:
		return "DEV0048494610"
	case DeviceInformation_DEV0048494910:
		return "DEV0048494910"
	case DeviceInformation_DEV0048134A10:
		return "DEV0048134A10"
	case DeviceInformation_DEV0048107E12:
		return "DEV0048107E12"
	case DeviceInformation_DEV000CB00812:
		return "DEV000CB00812"
	case DeviceInformation_DEV0048FF2112:
		return "DEV0048FF2112"
	case DeviceInformation_DEV0048140A11:
		return "DEV0048140A11"
	case DeviceInformation_DEV0048140B12:
		return "DEV0048140B12"
	case DeviceInformation_DEV0048140C13:
		return "DEV0048140C13"
	case DeviceInformation_DEV0048139A10:
		return "DEV0048139A10"
	case DeviceInformation_DEV0048648B10:
		return "DEV0048648B10"
	case DeviceInformation_DEV0008A01111:
		return "DEV0008A01111"
	case DeviceInformation_DEV0008A01211:
		return "DEV0008A01211"
	case DeviceInformation_DEV0008A01212:
		return "DEV0008A01212"
	case DeviceInformation_DEV0008A01112:
		return "DEV0008A01112"
	case DeviceInformation_DEV000CB00713:
		return "DEV000CB00713"
	case DeviceInformation_DEV0008A03213:
		return "DEV0008A03213"
	case DeviceInformation_DEV0008A03313:
		return "DEV0008A03313"
	case DeviceInformation_DEV0008A01113:
		return "DEV0008A01113"
	case DeviceInformation_DEV0008A01711:
		return "DEV0008A01711"
	case DeviceInformation_DEV0008B00911:
		return "DEV0008B00911"
	case DeviceInformation_DEV0008C00102:
		return "DEV0008C00102"
	case DeviceInformation_DEV0008C00101:
		return "DEV0008C00101"
	case DeviceInformation_DEV0008901501:
		return "DEV0008901501"
	case DeviceInformation_DEV0008901310:
		return "DEV0008901310"
	case DeviceInformation_DEV000820B011:
		return "DEV000820B011"
	case DeviceInformation_DEV000C181910:
		return "DEV000C181910"
	case DeviceInformation_DEV0008705C11:
		return "DEV0008705C11"
	case DeviceInformation_DEV0008705D11:
		return "DEV0008705D11"
	case DeviceInformation_DEV0008706211:
		return "DEV0008706211"
	case DeviceInformation_DEV000820BA11:
		return "DEV000820BA11"
	case DeviceInformation_DEV000820C011:
		return "DEV000820C011"
	case DeviceInformation_DEV000820B311:
		return "DEV000820B311"
	case DeviceInformation_DEV0008301A11:
		return "DEV0008301A11"
	case DeviceInformation_DEV0008C00C13:
		return "DEV0008C00C13"
	case DeviceInformation_DEV0008302611:
		return "DEV0008302611"
	case DeviceInformation_DEV0008302311:
		return "DEV0008302311"
	case DeviceInformation_DEV000C181810:
		return "DEV000C181810"
	case DeviceInformation_DEV0008302011:
		return "DEV0008302011"
	case DeviceInformation_DEV0008C00C11:
		return "DEV0008C00C11"
	case DeviceInformation_DEV0008302612:
		return "DEV0008302612"
	case DeviceInformation_DEV0008302312:
		return "DEV0008302312"
	case DeviceInformation_DEV0008302012:
		return "DEV0008302012"
	case DeviceInformation_DEV0008C00C15:
		return "DEV0008C00C15"
	case DeviceInformation_DEV0008C00C14:
		return "DEV0008C00C14"
	case DeviceInformation_DEV0008B00713:
		return "DEV0008B00713"
	case DeviceInformation_DEV0008706611:
		return "DEV0008706611"
	case DeviceInformation_DEV0008706811:
		return "DEV0008706811"
	case DeviceInformation_DEV000C20C011:
		return "DEV000C20C011"
	case DeviceInformation_DEV0008B00812:
		return "DEV0008B00812"
	case DeviceInformation_DEV0008209016:
		return "DEV0008209016"
	case DeviceInformation_DEV0008209011:
		return "DEV0008209011"
	case DeviceInformation_DEV000820A011:
		return "DEV000820A011"
	case DeviceInformation_DEV0008208811:
		return "DEV0008208811"
	case DeviceInformation_DEV0008209811:
		return "DEV0008209811"
	case DeviceInformation_DEV000820CA11:
		return "DEV000820CA11"
	case DeviceInformation_DEV0008208012:
		return "DEV0008208012"
	case DeviceInformation_DEV0008207812:
		return "DEV0008207812"
	case DeviceInformation_DEV0008207811:
		return "DEV0008207811"
	case DeviceInformation_DEV0079002527:
		return "DEV0079002527"
	case DeviceInformation_DEV0008208011:
		return "DEV0008208011"
	case DeviceInformation_DEV000810D111:
		return "DEV000810D111"
	case DeviceInformation_DEV000810D511:
		return "DEV000810D511"
	case DeviceInformation_DEV000810FA12:
		return "DEV000810FA12"
	case DeviceInformation_DEV000810FB12:
		return "DEV000810FB12"
	case DeviceInformation_DEV000810F211:
		return "DEV000810F211"
	case DeviceInformation_DEV000810D211:
		return "DEV000810D211"
	case DeviceInformation_DEV000810E211:
		return "DEV000810E211"
	case DeviceInformation_DEV000810D611:
		return "DEV000810D611"
	case DeviceInformation_DEV000810F212:
		return "DEV000810F212"
	case DeviceInformation_DEV0079004027:
		return "DEV0079004027"
	case DeviceInformation_DEV000810E212:
		return "DEV000810E212"
	case DeviceInformation_DEV000810FC13:
		return "DEV000810FC13"
	case DeviceInformation_DEV000810FD13:
		return "DEV000810FD13"
	case DeviceInformation_DEV000810F311:
		return "DEV000810F311"
	case DeviceInformation_DEV000810D311:
		return "DEV000810D311"
	case DeviceInformation_DEV000810D711:
		return "DEV000810D711"
	case DeviceInformation_DEV000810F312:
		return "DEV000810F312"
	case DeviceInformation_DEV000810D811:
		return "DEV000810D811"
	case DeviceInformation_DEV000810E511:
		return "DEV000810E511"
	case DeviceInformation_DEV000810E512:
		return "DEV000810E512"
	case DeviceInformation_DEV0079000223:
		return "DEV0079000223"
	case DeviceInformation_DEV000810F611:
		return "DEV000810F611"
	case DeviceInformation_DEV000810D911:
		return "DEV000810D911"
	case DeviceInformation_DEV000810F612:
		return "DEV000810F612"
	case DeviceInformation_DEV000820A812:
		return "DEV000820A812"
	case DeviceInformation_DEV000820AD11:
		return "DEV000820AD11"
	case DeviceInformation_DEV000820BD11:
		return "DEV000820BD11"
	case DeviceInformation_DEV000820C711:
		return "DEV000820C711"
	case DeviceInformation_DEV000820CD11:
		return "DEV000820CD11"
	case DeviceInformation_DEV000820C411:
		return "DEV000820C411"
	case DeviceInformation_DEV000820A811:
		return "DEV000820A811"
	case DeviceInformation_DEV0064181610:
		return "DEV0064181610"
	case DeviceInformation_DEV0079000123:
		return "DEV0079000123"
	case DeviceInformation_DEV0008501411:
		return "DEV0008501411"
	case DeviceInformation_DEV0008C01602:
		return "DEV0008C01602"
	case DeviceInformation_DEV0008302613:
		return "DEV0008302613"
	case DeviceInformation_DEV0008302313:
		return "DEV0008302313"
	case DeviceInformation_DEV0008302013:
		return "DEV0008302013"
	case DeviceInformation_DEV0009207730:
		return "DEV0009207730"
	case DeviceInformation_DEV0009208F10:
		return "DEV0009208F10"
	case DeviceInformation_DEV0009C00C13:
		return "DEV0009C00C13"
	case DeviceInformation_DEV0009209910:
		return "DEV0009209910"
	case DeviceInformation_DEV0009209A10:
		return "DEV0009209A10"
	case DeviceInformation_DEV0079001427:
		return "DEV0079001427"
	case DeviceInformation_DEV0009207930:
		return "DEV0009207930"
	case DeviceInformation_DEV0009201720:
		return "DEV0009201720"
	case DeviceInformation_DEV0009500D01:
		return "DEV0009500D01"
	case DeviceInformation_DEV0009500E01:
		return "DEV0009500E01"
	case DeviceInformation_DEV0009209911:
		return "DEV0009209911"
	case DeviceInformation_DEV0009209A11:
		return "DEV0009209A11"
	case DeviceInformation_DEV0009C00C12:
		return "DEV0009C00C12"
	case DeviceInformation_DEV0009C00C11:
		return "DEV0009C00C11"
	case DeviceInformation_DEV0009500D20:
		return "DEV0009500D20"
	case DeviceInformation_DEV0009500E20:
		return "DEV0009500E20"
	case DeviceInformation_DEV0079003027:
		return "DEV0079003027"
	case DeviceInformation_DEV000920B910:
		return "DEV000920B910"
	case DeviceInformation_DEV0009E0CE10:
		return "DEV0009E0CE10"
	case DeviceInformation_DEV0009E0A210:
		return "DEV0009E0A210"
	case DeviceInformation_DEV0009501410:
		return "DEV0009501410"
	case DeviceInformation_DEV0009207830:
		return "DEV0009207830"
	case DeviceInformation_DEV0009201620:
		return "DEV0009201620"
	case DeviceInformation_DEV0009E0A111:
		return "DEV0009E0A111"
	case DeviceInformation_DEV0009E0CD11:
		return "DEV0009E0CD11"
	case DeviceInformation_DEV000920B811:
		return "DEV000920B811"
	case DeviceInformation_DEV000920B611:
		return "DEV000920B611"
	case DeviceInformation_DEV0079100C13:
		return "DEV0079100C13"
	case DeviceInformation_DEV0009207E10:
		return "DEV0009207E10"
	case DeviceInformation_DEV0009207630:
		return "DEV0009207630"
	case DeviceInformation_DEV0009205910:
		return "DEV0009205910"
	case DeviceInformation_DEV0009500B01:
		return "DEV0009500B01"
	case DeviceInformation_DEV000920AC10:
		return "DEV000920AC10"
	case DeviceInformation_DEV0009207430:
		return "DEV0009207430"
	case DeviceInformation_DEV0009204521:
		return "DEV0009204521"
	case DeviceInformation_DEV0009500A01:
		return "DEV0009500A01"
	case DeviceInformation_DEV0009500001:
		return "DEV0009500001"
	case DeviceInformation_DEV000920AB10:
		return "DEV000920AB10"
	case DeviceInformation_DEV0079101C11:
		return "DEV0079101C11"
	case DeviceInformation_DEV000920BF11:
		return "DEV000920BF11"
	case DeviceInformation_DEV0009203510:
		return "DEV0009203510"
	case DeviceInformation_DEV0009207A30:
		return "DEV0009207A30"
	case DeviceInformation_DEV0009500701:
		return "DEV0009500701"
	case DeviceInformation_DEV0009501710:
		return "DEV0009501710"
	case DeviceInformation_DEV000920B310:
		return "DEV000920B310"
	case DeviceInformation_DEV0009207530:
		return "DEV0009207530"
	case DeviceInformation_DEV0009203321:
		return "DEV0009203321"
	case DeviceInformation_DEV0009500C01:
		return "DEV0009500C01"
	case DeviceInformation_DEV000920AD10:
		return "DEV000920AD10"
	case DeviceInformation_DEV0080709010:
		return "DEV0080709010"
	case DeviceInformation_DEV0009207230:
		return "DEV0009207230"
	case DeviceInformation_DEV0009500801:
		return "DEV0009500801"
	case DeviceInformation_DEV0009501810:
		return "DEV0009501810"
	case DeviceInformation_DEV000920B410:
		return "DEV000920B410"
	case DeviceInformation_DEV0009207330:
		return "DEV0009207330"
	case DeviceInformation_DEV0009204421:
		return "DEV0009204421"
	case DeviceInformation_DEV0009500901:
		return "DEV0009500901"
	case DeviceInformation_DEV000920AA10:
		return "DEV000920AA10"
	case DeviceInformation_DEV0009209D01:
		return "DEV0009209D01"
	case DeviceInformation_DEV000920B010:
		return "DEV000920B010"
	case DeviceInformation_DEV0080707010:
		return "DEV0080707010"
	case DeviceInformation_DEV0009E0BE01:
		return "DEV0009E0BE01"
	case DeviceInformation_DEV000920B110:
		return "DEV000920B110"
	case DeviceInformation_DEV0009E0BD01:
		return "DEV0009E0BD01"
	case DeviceInformation_DEV0009D03F10:
		return "DEV0009D03F10"
	case DeviceInformation_DEV0009305F10:
		return "DEV0009305F10"
	case DeviceInformation_DEV0009305610:
		return "DEV0009305610"
	case DeviceInformation_DEV0009D03E10:
		return "DEV0009D03E10"
	case DeviceInformation_DEV0009306010:
		return "DEV0009306010"
	case DeviceInformation_DEV0009306110:
		return "DEV0009306110"
	case DeviceInformation_DEV0009306310:
		return "DEV0009306310"
	case DeviceInformation_DEV0080706010:
		return "DEV0080706010"
	case DeviceInformation_DEV0009D03B10:
		return "DEV0009D03B10"
	case DeviceInformation_DEV0009D03C10:
		return "DEV0009D03C10"
	case DeviceInformation_DEV0009D03910:
		return "DEV0009D03910"
	case DeviceInformation_DEV0009D03A10:
		return "DEV0009D03A10"
	case DeviceInformation_DEV0009305411:
		return "DEV0009305411"
	case DeviceInformation_DEV0009D03D11:
		return "DEV0009D03D11"
	case DeviceInformation_DEV0009304B11:
		return "DEV0009304B11"
	case DeviceInformation_DEV0009304C11:
		return "DEV0009304C11"
	case DeviceInformation_DEV0009306220:
		return "DEV0009306220"
	case DeviceInformation_DEV0009302E10:
		return "DEV0009302E10"
	case DeviceInformation_DEV0080706810:
		return "DEV0080706810"
	case DeviceInformation_DEV0009302F10:
		return "DEV0009302F10"
	case DeviceInformation_DEV0009303010:
		return "DEV0009303010"
	case DeviceInformation_DEV0009303110:
		return "DEV0009303110"
	case DeviceInformation_DEV0009306510:
		return "DEV0009306510"
	case DeviceInformation_DEV0009306610:
		return "DEV0009306610"
	case DeviceInformation_DEV0009306410:
		return "DEV0009306410"
	case DeviceInformation_DEV0009401110:
		return "DEV0009401110"
	case DeviceInformation_DEV0009400610:
		return "DEV0009400610"
	case DeviceInformation_DEV0009401510:
		return "DEV0009401510"
	case DeviceInformation_DEV0009402110:
		return "DEV0009402110"
	case DeviceInformation_DEV0080705010:
		return "DEV0080705010"
	case DeviceInformation_DEV0009400110:
		return "DEV0009400110"
	case DeviceInformation_DEV0009400910:
		return "DEV0009400910"
	case DeviceInformation_DEV0009400010:
		return "DEV0009400010"
	case DeviceInformation_DEV0009401810:
		return "DEV0009401810"
	case DeviceInformation_DEV0009400310:
		return "DEV0009400310"
	case DeviceInformation_DEV0009301810:
		return "DEV0009301810"
	case DeviceInformation_DEV0009301910:
		return "DEV0009301910"
	case DeviceInformation_DEV0009301A10:
		return "DEV0009301A10"
	case DeviceInformation_DEV0009401210:
		return "DEV0009401210"
	case DeviceInformation_DEV0009400810:
		return "DEV0009400810"
	case DeviceInformation_DEV006420C011:
		return "DEV006420C011"
	case DeviceInformation_DEV0080703013:
		return "DEV0080703013"
	case DeviceInformation_DEV0009400710:
		return "DEV0009400710"
	case DeviceInformation_DEV0009401310:
		return "DEV0009401310"
	case DeviceInformation_DEV0009401410:
		return "DEV0009401410"
	case DeviceInformation_DEV0009402210:
		return "DEV0009402210"
	case DeviceInformation_DEV0009402310:
		return "DEV0009402310"
	case DeviceInformation_DEV0009401710:
		return "DEV0009401710"
	case DeviceInformation_DEV0009401610:
		return "DEV0009401610"
	case DeviceInformation_DEV0009400210:
		return "DEV0009400210"
	case DeviceInformation_DEV0009401010:
		return "DEV0009401010"
	case DeviceInformation_DEV0009400510:
		return "DEV0009400510"
	case DeviceInformation_DEV0080704021:
		return "DEV0080704021"
	case DeviceInformation_DEV0009400410:
		return "DEV0009400410"
	case DeviceInformation_DEV0009D04B20:
		return "DEV0009D04B20"
	case DeviceInformation_DEV0009D04920:
		return "DEV0009D04920"
	case DeviceInformation_DEV0009D04A20:
		return "DEV0009D04A20"
	case DeviceInformation_DEV0009D04820:
		return "DEV0009D04820"
	case DeviceInformation_DEV0009D04C11:
		return "DEV0009D04C11"
	case DeviceInformation_DEV0009D05610:
		return "DEV0009D05610"
	case DeviceInformation_DEV0009305510:
		return "DEV0009305510"
	case DeviceInformation_DEV0009209810:
		return "DEV0009209810"
	case DeviceInformation_DEV0009202A10:
		return "DEV0009202A10"
	case DeviceInformation_DEV0080704022:
		return "DEV0080704022"
	case DeviceInformation_DEV0009209510:
		return "DEV0009209510"
	case DeviceInformation_DEV0009501110:
		return "DEV0009501110"
	case DeviceInformation_DEV0009209310:
		return "DEV0009209310"
	case DeviceInformation_DEV0009209410:
		return "DEV0009209410"
	case DeviceInformation_DEV0009209210:
		return "DEV0009209210"
	case DeviceInformation_DEV0009501210:
		return "DEV0009501210"
	case DeviceInformation_DEV0009205411:
		return "DEV0009205411"
	case DeviceInformation_DEV000920A111:
		return "DEV000920A111"
	case DeviceInformation_DEV000920A311:
		return "DEV000920A311"
	case DeviceInformation_DEV0009205112:
		return "DEV0009205112"
	case DeviceInformation_DEV0080704020:
		return "DEV0080704020"
	case DeviceInformation_DEV0009204110:
		return "DEV0009204110"
	case DeviceInformation_DEV0009E07710:
		return "DEV0009E07710"
	case DeviceInformation_DEV0009E07712:
		return "DEV0009E07712"
	case DeviceInformation_DEV0009205212:
		return "DEV0009205212"
	case DeviceInformation_DEV0009205211:
		return "DEV0009205211"
	case DeviceInformation_DEV0009205311:
		return "DEV0009205311"
	case DeviceInformation_DEV0009206B10:
		return "DEV0009206B10"
	case DeviceInformation_DEV0009208010:
		return "DEV0009208010"
	case DeviceInformation_DEV0009206A12:
		return "DEV0009206A12"
	case DeviceInformation_DEV0009206810:
		return "DEV0009206810"
	case DeviceInformation_DEV0080701111:
		return "DEV0080701111"
	case DeviceInformation_DEV0009208110:
		return "DEV0009208110"
	case DeviceInformation_DEV0009205511:
		return "DEV0009205511"
	case DeviceInformation_DEV0009209F01:
		return "DEV0009209F01"
	case DeviceInformation_DEV0009208C10:
		return "DEV0009208C10"
	case DeviceInformation_DEV0009208E10:
		return "DEV0009208E10"
	case DeviceInformation_DEV000920B511:
		return "DEV000920B511"
	case DeviceInformation_DEV0009501910:
		return "DEV0009501910"
	case DeviceInformation_DEV000920BE11:
		return "DEV000920BE11"
	case DeviceInformation_DEV0009209710:
		return "DEV0009209710"
	case DeviceInformation_DEV0009208510:
		return "DEV0009208510"
	case DeviceInformation_DEV0080701811:
		return "DEV0080701811"
	case DeviceInformation_DEV0009208610:
		return "DEV0009208610"
	case DeviceInformation_DEV000920BD10:
		return "DEV000920BD10"
	case DeviceInformation_DEV0009500210:
		return "DEV0009500210"
	case DeviceInformation_DEV0009500310:
		return "DEV0009500310"
	case DeviceInformation_DEV0009E0BF10:
		return "DEV0009E0BF10"
	case DeviceInformation_DEV0009E0C010:
		return "DEV0009E0C010"
	case DeviceInformation_DEV0009500110:
		return "DEV0009500110"
	case DeviceInformation_DEV0009209B10:
		return "DEV0009209B10"
	case DeviceInformation_DEV0009207D10:
		return "DEV0009207D10"
	case DeviceInformation_DEV0009202F11:
		return "DEV0009202F11"
	case DeviceInformation_DEV008020A110:
		return "DEV008020A110"
	case DeviceInformation_DEV0009203011:
		return "DEV0009203011"
	case DeviceInformation_DEV0009207C10:
		return "DEV0009207C10"
	case DeviceInformation_DEV0009207B10:
		return "DEV0009207B10"
	case DeviceInformation_DEV0009208710:
		return "DEV0009208710"
	case DeviceInformation_DEV0009E06610:
		return "DEV0009E06610"
	case DeviceInformation_DEV0009E06611:
		return "DEV0009E06611"
	case DeviceInformation_DEV0009E06410:
		return "DEV0009E06410"
	case DeviceInformation_DEV0009E06411:
		return "DEV0009E06411"
	case DeviceInformation_DEV0009E06210:
		return "DEV0009E06210"
	case DeviceInformation_DEV0009E0E910:
		return "DEV0009E0E910"
	case DeviceInformation_DEV008020A210:
		return "DEV008020A210"
	case DeviceInformation_DEV0009E0EB10:
		return "DEV0009E0EB10"
	case DeviceInformation_DEV000920BB10:
		return "DEV000920BB10"
	case DeviceInformation_DEV0009FF1B11:
		return "DEV0009FF1B11"
	case DeviceInformation_DEV0009E0CF10:
		return "DEV0009E0CF10"
	case DeviceInformation_DEV0009206C30:
		return "DEV0009206C30"
	case DeviceInformation_DEV0009206D30:
		return "DEV0009206D30"
	case DeviceInformation_DEV0009206E30:
		return "DEV0009206E30"
	case DeviceInformation_DEV0009206F30:
		return "DEV0009206F30"
	case DeviceInformation_DEV0009207130:
		return "DEV0009207130"
	case DeviceInformation_DEV0009204720:
		return "DEV0009204720"
	case DeviceInformation_DEV008020A010:
		return "DEV008020A010"
	case DeviceInformation_DEV0009204820:
		return "DEV0009204820"
	case DeviceInformation_DEV0009204920:
		return "DEV0009204920"
	case DeviceInformation_DEV0009204A20:
		return "DEV0009204A20"
	case DeviceInformation_DEV0009205A10:
		return "DEV0009205A10"
	case DeviceInformation_DEV0009207030:
		return "DEV0009207030"
	case DeviceInformation_DEV0009205B10:
		return "DEV0009205B10"
	case DeviceInformation_DEV0009500501:
		return "DEV0009500501"
	case DeviceInformation_DEV0009501001:
		return "DEV0009501001"
	case DeviceInformation_DEV0009500601:
		return "DEV0009500601"
	case DeviceInformation_DEV0009500F01:
		return "DEV0009500F01"
	case DeviceInformation_DEV0080207212:
		return "DEV0080207212"
	case DeviceInformation_DEV0009500401:
		return "DEV0009500401"
	case DeviceInformation_DEV000920B210:
		return "DEV000920B210"
	case DeviceInformation_DEV000920AE10:
		return "DEV000920AE10"
	case DeviceInformation_DEV000920BC10:
		return "DEV000920BC10"
	case DeviceInformation_DEV000920AF10:
		return "DEV000920AF10"
	case DeviceInformation_DEV0009207F10:
		return "DEV0009207F10"
	case DeviceInformation_DEV0009208910:
		return "DEV0009208910"
	case DeviceInformation_DEV0009205710:
		return "DEV0009205710"
	case DeviceInformation_DEV0009205810:
		return "DEV0009205810"
	case DeviceInformation_DEV0009203810:
		return "DEV0009203810"
	case DeviceInformation_DEV006420BA11:
		return "DEV006420BA11"
	case DeviceInformation_DEV0080209111:
		return "DEV0080209111"
	case DeviceInformation_DEV0009203910:
		return "DEV0009203910"
	case DeviceInformation_DEV0009203E10:
		return "DEV0009203E10"
	case DeviceInformation_DEV0009204B10:
		return "DEV0009204B10"
	case DeviceInformation_DEV0009203F10:
		return "DEV0009203F10"
	case DeviceInformation_DEV0009204C10:
		return "DEV0009204C10"
	case DeviceInformation_DEV0009204010:
		return "DEV0009204010"
	case DeviceInformation_DEV0009206411:
		return "DEV0009206411"
	case DeviceInformation_DEV0009205E10:
		return "DEV0009205E10"
	case DeviceInformation_DEV0009206711:
		return "DEV0009206711"
	case DeviceInformation_DEV000920A710:
		return "DEV000920A710"
	case DeviceInformation_DEV0080204310:
		return "DEV0080204310"
	case DeviceInformation_DEV000920A610:
		return "DEV000920A610"
	case DeviceInformation_DEV0009203A10:
		return "DEV0009203A10"
	case DeviceInformation_DEV0009203B10:
		return "DEV0009203B10"
	case DeviceInformation_DEV0009203C10:
		return "DEV0009203C10"
	case DeviceInformation_DEV0009203D10:
		return "DEV0009203D10"
	case DeviceInformation_DEV0009E05E12:
		return "DEV0009E05E12"
	case DeviceInformation_DEV0009E0B711:
		return "DEV0009E0B711"
	case DeviceInformation_DEV0009E06A12:
		return "DEV0009E06A12"
	case DeviceInformation_DEV0009E06E12:
		return "DEV0009E06E12"
	case DeviceInformation_DEV0009E0B720:
		return "DEV0009E0B720"
	case DeviceInformation_DEV008020B612:
		return "DEV008020B612"
	case DeviceInformation_DEV0009E0E611:
		return "DEV0009E0E611"
	case DeviceInformation_DEV0009E0B321:
		return "DEV0009E0B321"
	case DeviceInformation_DEV0009E0E512:
		return "DEV0009E0E512"
	case DeviceInformation_DEV0009204210:
		return "DEV0009204210"
	case DeviceInformation_DEV0009208210:
		return "DEV0009208210"
	case DeviceInformation_DEV0009E07211:
		return "DEV0009E07211"
	case DeviceInformation_DEV0009E0CC11:
		return "DEV0009E0CC11"
	case DeviceInformation_DEV0009110111:
		return "DEV0009110111"
	case DeviceInformation_DEV0009110211:
		return "DEV0009110211"
	case DeviceInformation_DEV000916B212:
		return "DEV000916B212"
	case DeviceInformation_DEV008020B412:
		return "DEV008020B412"
	case DeviceInformation_DEV0009110212:
		return "DEV0009110212"
	case DeviceInformation_DEV0009110311:
		return "DEV0009110311"
	case DeviceInformation_DEV000916B312:
		return "DEV000916B312"
	case DeviceInformation_DEV0009110312:
		return "DEV0009110312"
	case DeviceInformation_DEV0009110411:
		return "DEV0009110411"
	case DeviceInformation_DEV0009110412:
		return "DEV0009110412"
	case DeviceInformation_DEV0009501615:
		return "DEV0009501615"
	case DeviceInformation_DEV0009E0ED10:
		return "DEV0009E0ED10"
	case DeviceInformation_DEV014F030110:
		return "DEV014F030110"
	case DeviceInformation_DEV014F030310:
		return "DEV014F030310"
	case DeviceInformation_DEV008020B512:
		return "DEV008020B512"
	case DeviceInformation_DEV014F030210:
		return "DEV014F030210"
	case DeviceInformation_DEV00EE7FFF10:
		return "DEV00EE7FFF10"
	case DeviceInformation_DEV00B6464101:
		return "DEV00B6464101"
	case DeviceInformation_DEV00B6464201:
		return "DEV00B6464201"
	case DeviceInformation_DEV00B6464501:
		return "DEV00B6464501"
	case DeviceInformation_DEV00B6434101:
		return "DEV00B6434101"
	case DeviceInformation_DEV00B6434201:
		return "DEV00B6434201"
	case DeviceInformation_DEV00B6434202:
		return "DEV00B6434202"
	case DeviceInformation_DEV00B6454101:
		return "DEV00B6454101"
	case DeviceInformation_DEV00B6454201:
		return "DEV00B6454201"
	case DeviceInformation_DEV0080208310:
		return "DEV0080208310"
	case DeviceInformation_DEV00B6455001:
		return "DEV00B6455001"
	case DeviceInformation_DEV00B6453101:
		return "DEV00B6453101"
	case DeviceInformation_DEV00B6453102:
		return "DEV00B6453102"
	case DeviceInformation_DEV00B6454102:
		return "DEV00B6454102"
	case DeviceInformation_DEV00B6454401:
		return "DEV00B6454401"
	case DeviceInformation_DEV00B6454402:
		return "DEV00B6454402"
	case DeviceInformation_DEV00B6454202:
		return "DEV00B6454202"
	case DeviceInformation_DEV00B6453103:
		return "DEV00B6453103"
	case DeviceInformation_DEV00B6453201:
		return "DEV00B6453201"
	case DeviceInformation_DEV00B6453301:
		return "DEV00B6453301"
	case DeviceInformation_DEV0080702111:
		return "DEV0080702111"
	case DeviceInformation_DEV00B6453104:
		return "DEV00B6453104"
	case DeviceInformation_DEV00B6454403:
		return "DEV00B6454403"
	case DeviceInformation_DEV00B6454801:
		return "DEV00B6454801"
	case DeviceInformation_DEV00B6414701:
		return "DEV00B6414701"
	case DeviceInformation_DEV00B6414201:
		return "DEV00B6414201"
	case DeviceInformation_DEV00B6474101:
		return "DEV00B6474101"
	case DeviceInformation_DEV00B6474302:
		return "DEV00B6474302"
	case DeviceInformation_DEV00B6474602:
		return "DEV00B6474602"
	case DeviceInformation_DEV00B6534D01:
		return "DEV00B6534D01"
	case DeviceInformation_DEV00B6535001:
		return "DEV00B6535001"
	case DeviceInformation_DEV0081FE0111:
		return "DEV0081FE0111"
	case DeviceInformation_DEV00B6455002:
		return "DEV00B6455002"
	case DeviceInformation_DEV00B6453701:
		return "DEV00B6453701"
	case DeviceInformation_DEV00B6484101:
		return "DEV00B6484101"
	case DeviceInformation_DEV00B6484201:
		return "DEV00B6484201"
	case DeviceInformation_DEV00B6484202:
		return "DEV00B6484202"
	case DeviceInformation_DEV00B6484301:
		return "DEV00B6484301"
	case DeviceInformation_DEV00B6484102:
		return "DEV00B6484102"
	case DeviceInformation_DEV00B6455101:
		return "DEV00B6455101"
	case DeviceInformation_DEV00B6455003:
		return "DEV00B6455003"
	case DeviceInformation_DEV00B6455102:
		return "DEV00B6455102"
	case DeviceInformation_DEV0081FF3131:
		return "DEV0081FF3131"
	case DeviceInformation_DEV00B6453702:
		return "DEV00B6453702"
	case DeviceInformation_DEV00B6453703:
		return "DEV00B6453703"
	case DeviceInformation_DEV00B6484302:
		return "DEV00B6484302"
	case DeviceInformation_DEV00B6484801:
		return "DEV00B6484801"
	case DeviceInformation_DEV00B6484501:
		return "DEV00B6484501"
	case DeviceInformation_DEV00B6484203:
		return "DEV00B6484203"
	case DeviceInformation_DEV00B6484103:
		return "DEV00B6484103"
	case DeviceInformation_DEV00B6455004:
		return "DEV00B6455004"
	case DeviceInformation_DEV00B6455103:
		return "DEV00B6455103"
	case DeviceInformation_DEV00B6455401:
		return "DEV00B6455401"
	case DeviceInformation_DEV0081F01313:
		return "DEV0081F01313"
	case DeviceInformation_DEV00B6455201:
		return "DEV00B6455201"
	case DeviceInformation_DEV00B6455402:
		return "DEV00B6455402"
	case DeviceInformation_DEV00B6455403:
		return "DEV00B6455403"
	case DeviceInformation_DEV00B603430A:
		return "DEV00B603430A"
	case DeviceInformation_DEV00B600010A:
		return "DEV00B600010A"
	case DeviceInformation_DEV00B6FF110A:
		return "DEV00B6FF110A"
	case DeviceInformation_DEV00B6434601:
		return "DEV00B6434601"
	case DeviceInformation_DEV00B6434602:
		return "DEV00B6434602"
	case DeviceInformation_DEV00B6455301:
		return "DEV00B6455301"
	case DeviceInformation_DEV00C5070410:
		return "DEV00C5070410"
	case DeviceInformation_DEV0064182010:
		return "DEV0064182010"
	case DeviceInformation_DEV0083002C16:
		return "DEV0083002C16"
	case DeviceInformation_DEV00C5070210:
		return "DEV00C5070210"
	case DeviceInformation_DEV00C5070610:
		return "DEV00C5070610"
	case DeviceInformation_DEV00C5070E11:
		return "DEV00C5070E11"
	case DeviceInformation_DEV00C5060240:
		return "DEV00C5060240"
	case DeviceInformation_DEV00C5062010:
		return "DEV00C5062010"
	case DeviceInformation_DEV00C5080230:
		return "DEV00C5080230"
	case DeviceInformation_DEV00C5060310:
		return "DEV00C5060310"
	case DeviceInformation_DEV006C070E11:
		return "DEV006C070E11"
	case DeviceInformation_DEV006C050002:
		return "DEV006C050002"
	case DeviceInformation_DEV006C011311:
		return "DEV006C011311"
	case DeviceInformation_DEV0083002E16:
		return "DEV0083002E16"
	case DeviceInformation_DEV006C011411:
		return "DEV006C011411"
	case DeviceInformation_DEV0007632010:
		return "DEV0007632010"
	case DeviceInformation_DEV0007632020:
		return "DEV0007632020"
	case DeviceInformation_DEV0007632180:
		return "DEV0007632180"
	case DeviceInformation_DEV0007632040:
		return "DEV0007632040"
	case DeviceInformation_DEV0007613812:
		return "DEV0007613812"
	case DeviceInformation_DEV0007613810:
		return "DEV0007613810"
	case DeviceInformation_DEV000720C011:
		return "DEV000720C011"
	case DeviceInformation_DEV0007A05210:
		return "DEV0007A05210"
	case DeviceInformation_DEV0007A08B10:
		return "DEV0007A08B10"
	case DeviceInformation_DEV0083002F16:
		return "DEV0083002F16"
	case DeviceInformation_DEV0007A05B32:
		return "DEV0007A05B32"
	case DeviceInformation_DEV0007A06932:
		return "DEV0007A06932"
	case DeviceInformation_DEV0007A06D32:
		return "DEV0007A06D32"
	case DeviceInformation_DEV0007A08032:
		return "DEV0007A08032"
	case DeviceInformation_DEV0007A00213:
		return "DEV0007A00213"
	case DeviceInformation_DEV0007A09532:
		return "DEV0007A09532"
	case DeviceInformation_DEV0007A06C32:
		return "DEV0007A06C32"
	case DeviceInformation_DEV0007A05E32:
		return "DEV0007A05E32"
	case DeviceInformation_DEV0007A08A32:
		return "DEV0007A08A32"
	case DeviceInformation_DEV0007A07032:
		return "DEV0007A07032"
	case DeviceInformation_DEV0083012F16:
		return "DEV0083012F16"
	case DeviceInformation_DEV0007A08332:
		return "DEV0007A08332"
	case DeviceInformation_DEV0007A09832:
		return "DEV0007A09832"
	case DeviceInformation_DEV0007A05C32:
		return "DEV0007A05C32"
	case DeviceInformation_DEV0007A06A32:
		return "DEV0007A06A32"
	case DeviceInformation_DEV0007A08832:
		return "DEV0007A08832"
	case DeviceInformation_DEV0007A06E32:
		return "DEV0007A06E32"
	case DeviceInformation_DEV0007A08132:
		return "DEV0007A08132"
	case DeviceInformation_DEV0007A00113:
		return "DEV0007A00113"
	case DeviceInformation_DEV0007A09632:
		return "DEV0007A09632"
	case DeviceInformation_DEV0007A05D32:
		return "DEV0007A05D32"
	case DeviceInformation_DEV0083003210:
		return "DEV0083003210"
	case DeviceInformation_DEV0007A06B32:
		return "DEV0007A06B32"
	case DeviceInformation_DEV0007A08932:
		return "DEV0007A08932"
	case DeviceInformation_DEV0007A06F32:
		return "DEV0007A06F32"
	case DeviceInformation_DEV0007A08232:
		return "DEV0007A08232"
	case DeviceInformation_DEV0007A09732:
		return "DEV0007A09732"
	case DeviceInformation_DEV0007A05713:
		return "DEV0007A05713"
	case DeviceInformation_DEV0007A01811:
		return "DEV0007A01811"
	case DeviceInformation_DEV0007A01911:
		return "DEV0007A01911"
	case DeviceInformation_DEV0007A04912:
		return "DEV0007A04912"
	case DeviceInformation_DEV0007A05814:
		return "DEV0007A05814"
	case DeviceInformation_DEV0083001D13:
		return "DEV0083001D13"
	case DeviceInformation_DEV0007A07114:
		return "DEV0007A07114"
	case DeviceInformation_DEV0007A05810:
		return "DEV0007A05810"
	case DeviceInformation_DEV0007A04312:
		return "DEV0007A04312"
	case DeviceInformation_DEV0007A04412:
		return "DEV0007A04412"
	case DeviceInformation_DEV0007A04512:
		return "DEV0007A04512"
	case DeviceInformation_DEV000720BD11:
		return "DEV000720BD11"
	case DeviceInformation_DEV0007A04C13:
		return "DEV0007A04C13"
	case DeviceInformation_DEV0007A04D13:
		return "DEV0007A04D13"
	case DeviceInformation_DEV0007A04B10:
		return "DEV0007A04B10"
	case DeviceInformation_DEV0007A04E13:
		return "DEV0007A04E13"
	case DeviceInformation_DEV0083001E13:
		return "DEV0083001E13"
	case DeviceInformation_DEV0007A04F13:
		return "DEV0007A04F13"
	case DeviceInformation_DEV000720BA11:
		return "DEV000720BA11"
	case DeviceInformation_DEV0007A03D11:
		return "DEV0007A03D11"
	case DeviceInformation_DEV0007A09211:
		return "DEV0007A09211"
	case DeviceInformation_DEV0007A09111:
		return "DEV0007A09111"
	case DeviceInformation_DEV0007FF1115:
		return "DEV0007FF1115"
	case DeviceInformation_DEV0007A01511:
		return "DEV0007A01511"
	case DeviceInformation_DEV0007A08411:
		return "DEV0007A08411"
	case DeviceInformation_DEV0007A08511:
		return "DEV0007A08511"
	case DeviceInformation_DEV0007A03422:
		return "DEV0007A03422"
	case DeviceInformation_DEV0083001B13:
		return "DEV0083001B13"
	case DeviceInformation_DEV0007A07213:
		return "DEV0007A07213"
	case DeviceInformation_DEV0007A07420:
		return "DEV0007A07420"
	case DeviceInformation_DEV0007A07520:
		return "DEV0007A07520"
	case DeviceInformation_DEV0007A07B12:
		return "DEV0007A07B12"
	case DeviceInformation_DEV0007A07C12:
		return "DEV0007A07C12"
	case DeviceInformation_DEV0007A09311:
		return "DEV0007A09311"
	case DeviceInformation_DEV0007A09013:
		return "DEV0007A09013"
	case DeviceInformation_DEV0007A08F13:
		return "DEV0007A08F13"
	case DeviceInformation_DEV0007A07E10:
		return "DEV0007A07E10"
	case DeviceInformation_DEV0007A05510:
		return "DEV0007A05510"
	case DeviceInformation_DEV0083001C13:
		return "DEV0083001C13"
	case DeviceInformation_DEV0007A05910:
		return "DEV0007A05910"
	case DeviceInformation_DEV0007A08711:
		return "DEV0007A08711"
	case DeviceInformation_DEV0007A03D12:
		return "DEV0007A03D12"
	case DeviceInformation_DEV0007A09A12:
		return "DEV0007A09A12"
	case DeviceInformation_DEV0007A09B12:
		return "DEV0007A09B12"
	case DeviceInformation_DEV0007A06614:
		return "DEV0007A06614"
	case DeviceInformation_DEV0007A06514:
		return "DEV0007A06514"
	case DeviceInformation_DEV0007A06014:
		return "DEV0007A06014"
	case DeviceInformation_DEV0007A07714:
		return "DEV0007A07714"
	case DeviceInformation_DEV0007A06414:
		return "DEV0007A06414"
	case DeviceInformation_DEV0083001F11:
		return "DEV0083001F11"
	case DeviceInformation_DEV0007A06114:
		return "DEV0007A06114"
	case DeviceInformation_DEV0007A07814:
		return "DEV0007A07814"
	case DeviceInformation_DEV0007A06714:
		return "DEV0007A06714"
	case DeviceInformation_DEV0007A06214:
		return "DEV0007A06214"
	case DeviceInformation_DEV0007A07914:
		return "DEV0007A07914"
	case DeviceInformation_DEV000B0A8410:
		return "DEV000B0A8410"
	case DeviceInformation_DEV000B0A7E10:
		return "DEV000B0A7E10"
	case DeviceInformation_DEV000B0A7F10:
		return "DEV000B0A7F10"
	case DeviceInformation_DEV000B0A8010:
		return "DEV000B0A8010"
	case DeviceInformation_DEV000BBF9111:
		return "DEV000BBF9111"
	case DeviceInformation_DEV0064182510:
		return "DEV0064182510"
	case DeviceInformation_DEV0083003C10:
		return "DEV0083003C10"
	case DeviceInformation_DEV000B0A7810:
		return "DEV000B0A7810"
	case DeviceInformation_DEV000B0A7910:
		return "DEV000B0A7910"
	case DeviceInformation_DEV000B0A7A10:
		return "DEV000B0A7A10"
	case DeviceInformation_DEV000B0A8910:
		return "DEV000B0A8910"
	case DeviceInformation_DEV000B0A8310:
		return "DEV000B0A8310"
	case DeviceInformation_DEV000B0A8510:
		return "DEV000B0A8510"
	case DeviceInformation_DEV000B0A6319:
		return "DEV000B0A6319"
	case DeviceInformation_DEV0083001C20:
		return "DEV0083001C20"
	case DeviceInformation_DEV0083001B22:
		return "DEV0083001B22"
	case DeviceInformation_DEV0083003A14:
		return "DEV0083003A14"
	case DeviceInformation_DEV0083003B14:
		return "DEV0083003B14"
	case DeviceInformation_DEV0083003B24:
		return "DEV0083003B24"
	case DeviceInformation_DEV0083003A24:
		return "DEV0083003A24"
	case DeviceInformation_DEV0083005824:
		return "DEV0083005824"
	case DeviceInformation_DEV0083002828:
		return "DEV0083002828"
	case DeviceInformation_DEV0083002928:
		return "DEV0083002928"
	case DeviceInformation_DEV0064182610:
		return "DEV0064182610"
	case DeviceInformation_DEV0083002A18:
		return "DEV0083002A18"
	case DeviceInformation_DEV0083002B18:
		return "DEV0083002B18"
	case DeviceInformation_DEV0083002337:
		return "DEV0083002337"
	case DeviceInformation_DEV0083002838:
		return "DEV0083002838"
	case DeviceInformation_DEV0083002938:
		return "DEV0083002938"
	case DeviceInformation_DEV0083002A38:
		return "DEV0083002A38"
	case DeviceInformation_DEV0083002B38:
		return "DEV0083002B38"
	case DeviceInformation_DEV0083001321:
		return "DEV0083001321"
	case DeviceInformation_DEV0083001421:
		return "DEV0083001421"
	case DeviceInformation_DEV0083001521:
		return "DEV0083001521"
	case DeviceInformation_DEV0064182910:
		return "DEV0064182910"
	case DeviceInformation_DEV0083001621:
		return "DEV0083001621"
	case DeviceInformation_DEV0083000921:
		return "DEV0083000921"
	case DeviceInformation_DEV0083000D11:
		return "DEV0083000D11"
	case DeviceInformation_DEV0083000C11:
		return "DEV0083000C11"
	case DeviceInformation_DEV0083000E11:
		return "DEV0083000E11"
	case DeviceInformation_DEV0083000B11:
		return "DEV0083000B11"
	case DeviceInformation_DEV0083000A11:
		return "DEV0083000A11"
	case DeviceInformation_DEV0083000A21:
		return "DEV0083000A21"
	case DeviceInformation_DEV0083000B21:
		return "DEV0083000B21"
	case DeviceInformation_DEV0083000C21:
		return "DEV0083000C21"
	case DeviceInformation_DEV0001140C13:
		return "DEV0001140C13"
	case DeviceInformation_DEV0064130610:
		return "DEV0064130610"
	case DeviceInformation_DEV0083000D21:
		return "DEV0083000D21"
	case DeviceInformation_DEV0083000821:
		return "DEV0083000821"
	case DeviceInformation_DEV0083000E21:
		return "DEV0083000E21"
	case DeviceInformation_DEV0083001812:
		return "DEV0083001812"
	case DeviceInformation_DEV0083001712:
		return "DEV0083001712"
	case DeviceInformation_DEV0083001816:
		return "DEV0083001816"
	case DeviceInformation_DEV0083001916:
		return "DEV0083001916"
	case DeviceInformation_DEV0083001716:
		return "DEV0083001716"
	case DeviceInformation_DEV0083001921:
		return "DEV0083001921"
	case DeviceInformation_DEV0083001721:
		return "DEV0083001721"
	case DeviceInformation_DEV0064130710:
		return "DEV0064130710"
	case DeviceInformation_DEV0083001821:
		return "DEV0083001821"
	case DeviceInformation_DEV0083001A20:
		return "DEV0083001A20"
	case DeviceInformation_DEV0083002320:
		return "DEV0083002320"
	case DeviceInformation_DEV0083001010:
		return "DEV0083001010"
	case DeviceInformation_DEV0083000F10:
		return "DEV0083000F10"
	case DeviceInformation_DEV0083003D14:
		return "DEV0083003D14"
	case DeviceInformation_DEV0083003E14:
		return "DEV0083003E14"
	case DeviceInformation_DEV0083003F14:
		return "DEV0083003F14"
	case DeviceInformation_DEV0083004014:
		return "DEV0083004014"
	case DeviceInformation_DEV0083004024:
		return "DEV0083004024"
	case DeviceInformation_DEV0064133510:
		return "DEV0064133510"
	case DeviceInformation_DEV0083003D24:
		return "DEV0083003D24"
	case DeviceInformation_DEV0083003E24:
		return "DEV0083003E24"
	case DeviceInformation_DEV0083003F24:
		return "DEV0083003F24"
	case DeviceInformation_DEV0083001112:
		return "DEV0083001112"
	case DeviceInformation_DEV0083001212:
		return "DEV0083001212"
	case DeviceInformation_DEV0083005B12:
		return "DEV0083005B12"
	case DeviceInformation_DEV0083005A12:
		return "DEV0083005A12"
	case DeviceInformation_DEV0083008410:
		return "DEV0083008410"
	case DeviceInformation_DEV0083008510:
		return "DEV0083008510"
	case DeviceInformation_DEV0083008610:
		return "DEV0083008610"
	case DeviceInformation_DEV0064133310:
		return "DEV0064133310"
	case DeviceInformation_DEV0083008710:
		return "DEV0083008710"
	case DeviceInformation_DEV0083002515:
		return "DEV0083002515"
	case DeviceInformation_DEV0083002115:
		return "DEV0083002115"
	case DeviceInformation_DEV0083002015:
		return "DEV0083002015"
	case DeviceInformation_DEV0083002415:
		return "DEV0083002415"
	case DeviceInformation_DEV0083002615:
		return "DEV0083002615"
	case DeviceInformation_DEV0083002215:
		return "DEV0083002215"
	case DeviceInformation_DEV0083002715:
		return "DEV0083002715"
	case DeviceInformation_DEV0083002315:
		return "DEV0083002315"
	case DeviceInformation_DEV0083008B25:
		return "DEV0083008B25"
	case DeviceInformation_DEV0064133410:
		return "DEV0064133410"
	case DeviceInformation_DEV0083008A25:
		return "DEV0083008A25"
	case DeviceInformation_DEV0083008B28:
		return "DEV0083008B28"
	case DeviceInformation_DEV0083008A28:
		return "DEV0083008A28"
	case DeviceInformation_DEV0083009013:
		return "DEV0083009013"
	case DeviceInformation_DEV0083009213:
		return "DEV0083009213"
	case DeviceInformation_DEV0083009113:
		return "DEV0083009113"
	case DeviceInformation_DEV0083009313:
		return "DEV0083009313"
	case DeviceInformation_DEV0083009413:
		return "DEV0083009413"
	case DeviceInformation_DEV0083009513:
		return "DEV0083009513"
	case DeviceInformation_DEV0083009613:
		return "DEV0083009613"
	case DeviceInformation_DEV0064133610:
		return "DEV0064133610"
	case DeviceInformation_DEV0083009713:
		return "DEV0083009713"
	case DeviceInformation_DEV0083009A13:
		return "DEV0083009A13"
	case DeviceInformation_DEV0083009B13:
		return "DEV0083009B13"
	case DeviceInformation_DEV0083004B11:
		return "DEV0083004B11"
	case DeviceInformation_DEV0083004B20:
		return "DEV0083004B20"
	case DeviceInformation_DEV0083005514:
		return "DEV0083005514"
	case DeviceInformation_DEV0083006824:
		return "DEV0083006824"
	case DeviceInformation_DEV0083006624:
		return "DEV0083006624"
	case DeviceInformation_DEV0083006524:
		return "DEV0083006524"
	case DeviceInformation_DEV0083006424:
		return "DEV0083006424"
	case DeviceInformation_DEV0064130510:
		return "DEV0064130510"
	case DeviceInformation_DEV0083006734:
		return "DEV0083006734"
	case DeviceInformation_DEV0083006434:
		return "DEV0083006434"
	case DeviceInformation_DEV0083006634:
		return "DEV0083006634"
	case DeviceInformation_DEV0083006534:
		return "DEV0083006534"
	case DeviceInformation_DEV0083006A34:
		return "DEV0083006A34"
	case DeviceInformation_DEV0083006B34:
		return "DEV0083006B34"
	case DeviceInformation_DEV0083006934:
		return "DEV0083006934"
	case DeviceInformation_DEV0083004F11:
		return "DEV0083004F11"
	case DeviceInformation_DEV0083004E10:
		return "DEV0083004E10"
	case DeviceInformation_DEV0083004D13:
		return "DEV0083004D13"
	case DeviceInformation_DEV0064480611:
		return "DEV0064480611"
	case DeviceInformation_DEV0083004414:
		return "DEV0083004414"
	case DeviceInformation_DEV0083004114:
		return "DEV0083004114"
	case DeviceInformation_DEV0083004514:
		return "DEV0083004514"
	case DeviceInformation_DEV0083004213:
		return "DEV0083004213"
	case DeviceInformation_DEV0083004313:
		return "DEV0083004313"
	case DeviceInformation_DEV0083004C11:
		return "DEV0083004C11"
	case DeviceInformation_DEV0083004913:
		return "DEV0083004913"
	case DeviceInformation_DEV0083004A13:
		return "DEV0083004A13"
	case DeviceInformation_DEV0083004712:
		return "DEV0083004712"
	case DeviceInformation_DEV0083004610:
		return "DEV0083004610"
	case DeviceInformation_DEV0064482011:
		return "DEV0064482011"
	case DeviceInformation_DEV0083008E12:
		return "DEV0083008E12"
	case DeviceInformation_DEV0083004813:
		return "DEV0083004813"
	case DeviceInformation_DEV0083005611:
		return "DEV0083005611"
	case DeviceInformation_DEV0083005710:
		return "DEV0083005710"
	case DeviceInformation_DEV0083005010:
		return "DEV0083005010"
	case DeviceInformation_DEV0083001A10:
		return "DEV0083001A10"
	case DeviceInformation_DEV0083002918:
		return "DEV0083002918"
	case DeviceInformation_DEV0083002818:
		return "DEV0083002818"
	case DeviceInformation_DEV0083006724:
		return "DEV0083006724"
	case DeviceInformation_DEV0083006D41:
		return "DEV0083006D41"
	case DeviceInformation_DEV0064182210:
		return "DEV0064182210"
	case DeviceInformation_DEV0083006E41:
		return "DEV0083006E41"
	case DeviceInformation_DEV0083007342:
		return "DEV0083007342"
	case DeviceInformation_DEV0083007242:
		return "DEV0083007242"
	case DeviceInformation_DEV0083006C42:
		return "DEV0083006C42"
	case DeviceInformation_DEV0083007542:
		return "DEV0083007542"
	case DeviceInformation_DEV0083007442:
		return "DEV0083007442"
	case DeviceInformation_DEV0083007742:
		return "DEV0083007742"
	case DeviceInformation_DEV0083007642:
		return "DEV0083007642"
	case DeviceInformation_DEV008300B030:
		return "DEV008300B030"
	case DeviceInformation_DEV008300B130:
		return "DEV008300B130"
	case DeviceInformation_DEV0001140B11:
		return "DEV0001140B11"
	case DeviceInformation_DEV0064182710:
		return "DEV0064182710"
	case DeviceInformation_DEV008300B230:
		return "DEV008300B230"
	case DeviceInformation_DEV008300B330:
		return "DEV008300B330"
	case DeviceInformation_DEV008300B430:
		return "DEV008300B430"
	case DeviceInformation_DEV008300B530:
		return "DEV008300B530"
	case DeviceInformation_DEV008300B630:
		return "DEV008300B630"
	case DeviceInformation_DEV008300B730:
		return "DEV008300B730"
	case DeviceInformation_DEV0083012843:
		return "DEV0083012843"
	case DeviceInformation_DEV0083012943:
		return "DEV0083012943"
	case DeviceInformation_DEV008300A421:
		return "DEV008300A421"
	case DeviceInformation_DEV008300A521:
		return "DEV008300A521"
	case DeviceInformation_DEV0064183010:
		return "DEV0064183010"
	case DeviceInformation_DEV008300A621:
		return "DEV008300A621"
	case DeviceInformation_DEV0083001332:
		return "DEV0083001332"
	case DeviceInformation_DEV0083000932:
		return "DEV0083000932"
	case DeviceInformation_DEV0083001432:
		return "DEV0083001432"
	case DeviceInformation_DEV0083001532:
		return "DEV0083001532"
	case DeviceInformation_DEV0083001632:
		return "DEV0083001632"
	case DeviceInformation_DEV008300A432:
		return "DEV008300A432"
	case DeviceInformation_DEV008300A532:
		return "DEV008300A532"
	case DeviceInformation_DEV008300A632:
		return "DEV008300A632"
	case DeviceInformation_DEV0083000F32:
		return "DEV0083000F32"
	case DeviceInformation_DEV0064B00812:
		return "DEV0064B00812"
	case DeviceInformation_DEV0083001032:
		return "DEV0083001032"
	case DeviceInformation_DEV0083000632:
		return "DEV0083000632"
	case DeviceInformation_DEV0083009810:
		return "DEV0083009810"
	case DeviceInformation_DEV0083009910:
		return "DEV0083009910"
	case DeviceInformation_DEV0083005C11:
		return "DEV0083005C11"
	case DeviceInformation_DEV0083005D11:
		return "DEV0083005D11"
	case DeviceInformation_DEV0083005E11:
		return "DEV0083005E11"
	case DeviceInformation_DEV0083005F11:
		return "DEV0083005F11"
	case DeviceInformation_DEV0083005413:
		return "DEV0083005413"
	case DeviceInformation_DEV0085000520:
		return "DEV0085000520"
	case DeviceInformation_DEV0064B00A01:
		return "DEV0064B00A01"
	case DeviceInformation_DEV0085000620:
		return "DEV0085000620"
	case DeviceInformation_DEV0085000720:
		return "DEV0085000720"
	case DeviceInformation_DEV0085012210:
		return "DEV0085012210"
	case DeviceInformation_DEV0085011210:
		return "DEV0085011210"
	case DeviceInformation_DEV0085013220:
		return "DEV0085013220"
	case DeviceInformation_DEV0085010210:
		return "DEV0085010210"
	case DeviceInformation_DEV0085000A10:
		return "DEV0085000A10"
	case DeviceInformation_DEV0085000B10:
		return "DEV0085000B10"
	case DeviceInformation_DEV0085071010:
		return "DEV0085071010"
	case DeviceInformation_DEV008500FB10:
		return "DEV008500FB10"
	case DeviceInformation_DEV0064760110:
		return "DEV0064760110"
	case DeviceInformation_DEV0085060210:
		return "DEV0085060210"
	case DeviceInformation_DEV0085060110:
		return "DEV0085060110"
	case DeviceInformation_DEV0085000D20:
		return "DEV0085000D20"
	case DeviceInformation_DEV008500C810:
		return "DEV008500C810"
	case DeviceInformation_DEV0085040111:
		return "DEV0085040111"
	case DeviceInformation_DEV008500C910:
		return "DEV008500C910"
	case DeviceInformation_DEV0085045020:
		return "DEV0085045020"
	case DeviceInformation_DEV0085070210:
		return "DEV0085070210"
	case DeviceInformation_DEV0085070110:
		return "DEV0085070110"
	case DeviceInformation_DEV0085070310:
		return "DEV0085070310"
	case DeviceInformation_DEV0064242313:
		return "DEV0064242313"
	case DeviceInformation_DEV0085000E20:
		return "DEV0085000E20"
	case DeviceInformation_DEV008E596010:
		return "DEV008E596010"
	case DeviceInformation_DEV008E593710:
		return "DEV008E593710"
	case DeviceInformation_DEV008E597710:
		return "DEV008E597710"
	case DeviceInformation_DEV008E598310:
		return "DEV008E598310"
	case DeviceInformation_DEV008E598910:
		return "DEV008E598910"
	case DeviceInformation_DEV008E593720:
		return "DEV008E593720"
	case DeviceInformation_DEV008E598920:
		return "DEV008E598920"
	case DeviceInformation_DEV008E598320:
		return "DEV008E598320"
	case DeviceInformation_DEV008E596021:
		return "DEV008E596021"
	case DeviceInformation_DEV0064FF2111:
		return "DEV0064FF2111"
	case DeviceInformation_DEV008E597721:
		return "DEV008E597721"
	case DeviceInformation_DEV008E587320:
		return "DEV008E587320"
	case DeviceInformation_DEV008E587020:
		return "DEV008E587020"
	case DeviceInformation_DEV008E587220:
		return "DEV008E587220"
	case DeviceInformation_DEV008E587120:
		return "DEV008E587120"
	case DeviceInformation_DEV008E679910:
		return "DEV008E679910"
	case DeviceInformation_DEV008E618310:
		return "DEV008E618310"
	case DeviceInformation_DEV008E707910:
		return "DEV008E707910"
	case DeviceInformation_DEV008E004010:
		return "DEV008E004010"
	case DeviceInformation_DEV008E570910:
		return "DEV008E570910"
	case DeviceInformation_DEV0064FF2112:
		return "DEV0064FF2112"
	case DeviceInformation_DEV008E558810:
		return "DEV008E558810"
	case DeviceInformation_DEV008E683410:
		return "DEV008E683410"
	case DeviceInformation_DEV008E707710:
		return "DEV008E707710"
	case DeviceInformation_DEV008E707810:
		return "DEV008E707810"
	case DeviceInformation_DEV0091100013:
		return "DEV0091100013"
	case DeviceInformation_DEV0091100110:
		return "DEV0091100110"
	case DeviceInformation_DEV009E670101:
		return "DEV009E670101"
	case DeviceInformation_DEV009E119311:
		return "DEV009E119311"
	case DeviceInformation_DEV00A2100C13:
		return "DEV00A2100C13"
	case DeviceInformation_DEV00A2101C11:
		return "DEV00A2101C11"
	case DeviceInformation_DEV0064648B10:
		return "DEV0064648B10"
	case DeviceInformation_DEV00A2300110:
		return "DEV00A2300110"
	case DeviceInformation_DEV0002A05814:
		return "DEV0002A05814"
	case DeviceInformation_DEV0002A07114:
		return "DEV0002A07114"
	case DeviceInformation_DEV0002134A10:
		return "DEV0002134A10"
	case DeviceInformation_DEV0002A03422:
		return "DEV0002A03422"
	case DeviceInformation_DEV0002A03321:
		return "DEV0002A03321"
	case DeviceInformation_DEV0002648B10:
		return "DEV0002648B10"
	case DeviceInformation_DEV0002A09013:
		return "DEV0002A09013"
	case DeviceInformation_DEV0002A08F13:
		return "DEV0002A08F13"
	case DeviceInformation_DEV0002A05510:
		return "DEV0002A05510"
	case DeviceInformation_DEV0064724010:
		return "DEV0064724010"
	case DeviceInformation_DEV0002A05910:
		return "DEV0002A05910"
	case DeviceInformation_DEV0002A05326:
		return "DEV0002A05326"
	case DeviceInformation_DEV0002A05428:
		return "DEV0002A05428"
	case DeviceInformation_DEV0002A08411:
		return "DEV0002A08411"
	case DeviceInformation_DEV0002A08511:
		return "DEV0002A08511"
	case DeviceInformation_DEV0002A00F11:
		return "DEV0002A00F11"
	case DeviceInformation_DEV0002A07310:
		return "DEV0002A07310"
	case DeviceInformation_DEV0002A04110:
		return "DEV0002A04110"
	case DeviceInformation_DEV0002A03813:
		return "DEV0002A03813"
	case DeviceInformation_DEV0002A07F13:
		return "DEV0002A07F13"
	case DeviceInformation_DEV0001803002:
		return "DEV0001803002"
	case DeviceInformation_DEV006420BD11:
		return "DEV006420BD11"
	case DeviceInformation_DEV0002A08832:
		return "DEV0002A08832"
	case DeviceInformation_DEV0002A06E32:
		return "DEV0002A06E32"
	case DeviceInformation_DEV0002A08132:
		return "DEV0002A08132"
	case DeviceInformation_DEV0002A01D20:
		return "DEV0002A01D20"
	case DeviceInformation_DEV0002A02120:
		return "DEV0002A02120"
	case DeviceInformation_DEV0002A02520:
		return "DEV0002A02520"
	case DeviceInformation_DEV0002A02920:
		return "DEV0002A02920"
	case DeviceInformation_DEV0002A03A20:
		return "DEV0002A03A20"
	case DeviceInformation_DEV0002A05C32:
		return "DEV0002A05C32"
	case DeviceInformation_DEV0002A06A32:
		return "DEV0002A06A32"
	case DeviceInformation_DEV0064570011:
		return "DEV0064570011"
	case DeviceInformation_DEV0002A09632:
		return "DEV0002A09632"
	case DeviceInformation_DEV0002A08932:
		return "DEV0002A08932"
	case DeviceInformation_DEV0002A06F32:
		return "DEV0002A06F32"
	case DeviceInformation_DEV0002A08232:
		return "DEV0002A08232"
	case DeviceInformation_DEV0002A01E20:
		return "DEV0002A01E20"
	case DeviceInformation_DEV0002A02220:
		return "DEV0002A02220"
	case DeviceInformation_DEV0002A02620:
		return "DEV0002A02620"
	case DeviceInformation_DEV0002A02A20:
		return "DEV0002A02A20"
	case DeviceInformation_DEV0002A03B20:
		return "DEV0002A03B20"
	case DeviceInformation_DEV0002A05D32:
		return "DEV0002A05D32"
	case DeviceInformation_DEV0064570310:
		return "DEV0064570310"
	case DeviceInformation_DEV0002A06B32:
		return "DEV0002A06B32"
	case DeviceInformation_DEV0002A09732:
		return "DEV0002A09732"
	case DeviceInformation_DEV0002A08A32:
		return "DEV0002A08A32"
	case DeviceInformation_DEV0002A07032:
		return "DEV0002A07032"
	case DeviceInformation_DEV0002A08332:
		return "DEV0002A08332"
	case DeviceInformation_DEV0002A01F20:
		return "DEV0002A01F20"
	case DeviceInformation_DEV0002A02320:
		return "DEV0002A02320"
	case DeviceInformation_DEV0002A02720:
		return "DEV0002A02720"
	case DeviceInformation_DEV0002A02B20:
		return "DEV0002A02B20"
	case DeviceInformation_DEV0002A04820:
		return "DEV0002A04820"
	case DeviceInformation_DEV0064570211:
		return "DEV0064570211"
	case DeviceInformation_DEV0002A06C32:
		return "DEV0002A06C32"
	case DeviceInformation_DEV0002A05E32:
		return "DEV0002A05E32"
	case DeviceInformation_DEV0002A09832:
		return "DEV0002A09832"
	case DeviceInformation_DEV0002A06D32:
		return "DEV0002A06D32"
	case DeviceInformation_DEV0002A08032:
		return "DEV0002A08032"
	case DeviceInformation_DEV0002A02020:
		return "DEV0002A02020"
	case DeviceInformation_DEV0002A02420:
		return "DEV0002A02420"
	case DeviceInformation_DEV0002A02820:
		return "DEV0002A02820"
	case DeviceInformation_DEV0002A03920:
		return "DEV0002A03920"
	case DeviceInformation_DEV0002A05B32:
		return "DEV0002A05B32"
	case DeviceInformation_DEV0064570411:
		return "DEV0064570411"
	case DeviceInformation_DEV0002A06932:
		return "DEV0002A06932"
	case DeviceInformation_DEV0002A09532:
		return "DEV0002A09532"
	case DeviceInformation_DEV0002B00813:
		return "DEV0002B00813"
	case DeviceInformation_DEV0002A0A610:
		return "DEV0002A0A610"
	case DeviceInformation_DEV0002A0A611:
		return "DEV0002A0A611"
	case DeviceInformation_DEV0002A0A510:
		return "DEV0002A0A510"
	case DeviceInformation_DEV0002A0A511:
		return "DEV0002A0A511"
	case DeviceInformation_DEV0002A00510:
		return "DEV0002A00510"
	case DeviceInformation_DEV0002A00610:
		return "DEV0002A00610"
	case DeviceInformation_DEV0002A01511:
		return "DEV0002A01511"
	case DeviceInformation_DEV0064570110:
		return "DEV0064570110"
	case DeviceInformation_DEV0002A03D11:
		return "DEV0002A03D11"
	case DeviceInformation_DEV000220C011:
		return "DEV000220C011"
	case DeviceInformation_DEV0002A05213:
		return "DEV0002A05213"
	case DeviceInformation_DEV0002A08B10:
		return "DEV0002A08B10"
	case DeviceInformation_DEV0002A0A210:
		return "DEV0002A0A210"
	case DeviceInformation_DEV0002A0A010:
		return "DEV0002A0A010"
	case DeviceInformation_DEV0002A09F10:
		return "DEV0002A09F10"
	case DeviceInformation_DEV0002A0A110:
		return "DEV0002A0A110"
	case DeviceInformation_DEV0002A0A013:
		return "DEV0002A0A013"
	case DeviceInformation_DEV0002A09F13:
		return "DEV0002A09F13"
	case DeviceInformation_DEV0064615022:
		return "DEV0064615022"
	case DeviceInformation_DEV0002A0A213:
		return "DEV0002A0A213"
	case DeviceInformation_DEV0002A0A113:
		return "DEV0002A0A113"
	case DeviceInformation_DEV0002A03F11:
		return "DEV0002A03F11"
	case DeviceInformation_DEV0002A04011:
		return "DEV0002A04011"
	case DeviceInformation_DEV0002FF2112:
		return "DEV0002FF2112"
	case DeviceInformation_DEV0002FF2111:
		return "DEV0002FF2111"
	case DeviceInformation_DEV0002720111:
		return "DEV0002720111"
	case DeviceInformation_DEV0002613812:
		return "DEV0002613812"
	case DeviceInformation_DEV0002A05713:
		return "DEV0002A05713"
	case DeviceInformation_DEV0002A07610:
		return "DEV0002A07610"
	case DeviceInformation_DEV0064182810:
		return "DEV0064182810"
	case DeviceInformation_DEV0002A01911:
		return "DEV0002A01911"
	case DeviceInformation_DEV0002A07611:
		return "DEV0002A07611"
	case DeviceInformation_DEV0002A04B10:
		return "DEV0002A04B10"
	case DeviceInformation_DEV0002A01B14:
		return "DEV0002A01B14"
	case DeviceInformation_DEV0002A09B11:
		return "DEV0002A09B11"
	case DeviceInformation_DEV0002A09B12:
		return "DEV0002A09B12"
	case DeviceInformation_DEV0002A03C10:
		return "DEV0002A03C10"
	case DeviceInformation_DEV0002A00213:
		return "DEV0002A00213"
	case DeviceInformation_DEV0002A00113:
		return "DEV0002A00113"
	case DeviceInformation_DEV0002A02C12:
		return "DEV0002A02C12"
	case DeviceInformation_DEV0064183110:
		return "DEV0064183110"
	case DeviceInformation_DEV0002A02D12:
		return "DEV0002A02D12"
	case DeviceInformation_DEV0002A02E12:
		return "DEV0002A02E12"
	case DeviceInformation_DEV0002A04C13:
		return "DEV0002A04C13"
	case DeviceInformation_DEV0002A04D13:
		return "DEV0002A04D13"
	case DeviceInformation_DEV0002A02F12:
		return "DEV0002A02F12"
	case DeviceInformation_DEV0002A03012:
		return "DEV0002A03012"
	case DeviceInformation_DEV0002A03112:
		return "DEV0002A03112"
	case DeviceInformation_DEV0002A04E13:
		return "DEV0002A04E13"
	case DeviceInformation_DEV0002A04F13:
		return "DEV0002A04F13"
	case DeviceInformation_DEV0002A01A13:
		return "DEV0002A01A13"
	case DeviceInformation_DEV0064133611:
		return "DEV0064133611"
	case DeviceInformation_DEV0002A09C11:
		return "DEV0002A09C11"
	case DeviceInformation_DEV0002A09C12:
		return "DEV0002A09C12"
	case DeviceInformation_DEV0002A01C20:
		return "DEV0002A01C20"
	case DeviceInformation_DEV0002A09A10:
		return "DEV0002A09A10"
	case DeviceInformation_DEV0002A09A12:
		return "DEV0002A09A12"
	case DeviceInformation_DEV000220BA11:
		return "DEV000220BA11"
	case DeviceInformation_DEV0002A03D12:
		return "DEV0002A03D12"
	case DeviceInformation_DEV0002A09110:
		return "DEV0002A09110"
	case DeviceInformation_DEV0002A09210:
		return "DEV0002A09210"
	case DeviceInformation_DEV0002A09111:
		return "DEV0002A09111"
	case DeviceInformation_DEV00641BD610:
		return "DEV00641BD610"
	case DeviceInformation_DEV006A000122:
		return "DEV006A000122"
	case DeviceInformation_DEV0002A09211:
		return "DEV0002A09211"
	case DeviceInformation_DEV0002A00E21:
		return "DEV0002A00E21"
	case DeviceInformation_DEV0002A03710:
		return "DEV0002A03710"
	case DeviceInformation_DEV0002A01112:
		return "DEV0002A01112"
	case DeviceInformation_DEV0002A01216:
		return "DEV0002A01216"
	case DeviceInformation_DEV0002A01217:
		return "DEV0002A01217"
	case DeviceInformation_DEV000220BD11:
		return "DEV000220BD11"
	case DeviceInformation_DEV0002A07F12:
		return "DEV0002A07F12"
	case DeviceInformation_DEV0002A06613:
		return "DEV0002A06613"
	case DeviceInformation_DEV0002A06713:
		return "DEV0002A06713"
	case DeviceInformation_DEV006A000222:
		return "DEV006A000222"
	case DeviceInformation_DEV0002A06013:
		return "DEV0002A06013"
	case DeviceInformation_DEV0002A06113:
		return "DEV0002A06113"
	case DeviceInformation_DEV0002A06213:
		return "DEV0002A06213"
	case DeviceInformation_DEV0002A06413:
		return "DEV0002A06413"
	case DeviceInformation_DEV0002A07713:
		return "DEV0002A07713"
	case DeviceInformation_DEV0002A07813:
		return "DEV0002A07813"
	case DeviceInformation_DEV0002A07913:
		return "DEV0002A07913"
	case DeviceInformation_DEV0002A07914:
		return "DEV0002A07914"
	case DeviceInformation_DEV0002A06114:
		return "DEV0002A06114"
	case DeviceInformation_DEV0002A06714:
		return "DEV0002A06714"
	case DeviceInformation_DEV006A070210:
		return "DEV006A070210"
	case DeviceInformation_DEV0002A06414:
		return "DEV0002A06414"
	case DeviceInformation_DEV0002A06214:
		return "DEV0002A06214"
	case DeviceInformation_DEV0002A06514:
		return "DEV0002A06514"
	case DeviceInformation_DEV0002A07714:
		return "DEV0002A07714"
	case DeviceInformation_DEV0002A06014:
		return "DEV0002A06014"
	case DeviceInformation_DEV0002A06614:
		return "DEV0002A06614"
	case DeviceInformation_DEV0002A07814:
		return "DEV0002A07814"
	case DeviceInformation_DEV0002A0C310:
		return "DEV0002A0C310"
	case DeviceInformation_DEV0002632010:
		return "DEV0002632010"
	case DeviceInformation_DEV0002632020:
		return "DEV0002632020"
	case DeviceInformation_DEV006BFFF713:
		return "DEV006BFFF713"
	case DeviceInformation_DEV0002632040:
		return "DEV0002632040"
	case DeviceInformation_DEV0002632180:
		return "DEV0002632180"
	case DeviceInformation_DEV0002632170:
		return "DEV0002632170"
	case DeviceInformation_DEV0002FF1140:
		return "DEV0002FF1140"
	case DeviceInformation_DEV0002A07E10:
		return "DEV0002A07E10"
	case DeviceInformation_DEV0002A07213:
		return "DEV0002A07213"
	case DeviceInformation_DEV0002A04A35:
		return "DEV0002A04A35"
	case DeviceInformation_DEV0002A07420:
		return "DEV0002A07420"
	case DeviceInformation_DEV0002A07520:
		return "DEV0002A07520"
	case DeviceInformation_DEV0002A07B12:
		return "DEV0002A07B12"
	case DeviceInformation_DEV006BFF2111:
		return "DEV006BFF2111"
	case DeviceInformation_DEV0002A07C12:
		return "DEV0002A07C12"
	case DeviceInformation_DEV0002A04312:
		return "DEV0002A04312"
	case DeviceInformation_DEV0002A04412:
		return "DEV0002A04412"
	case DeviceInformation_DEV0002A04512:
		return "DEV0002A04512"
	case DeviceInformation_DEV0002A04912:
		return "DEV0002A04912"
	case DeviceInformation_DEV0002A05012:
		return "DEV0002A05012"
	case DeviceInformation_DEV0002A01811:
		return "DEV0002A01811"
	case DeviceInformation_DEV0002A03E11:
		return "DEV0002A03E11"
	case DeviceInformation_DEV0002A08711:
		return "DEV0002A08711"
	case DeviceInformation_DEV0002A09311:
		return "DEV0002A09311"
	case DeviceInformation_DEV006BFFF820:
		return "DEV006BFFF820"
	case DeviceInformation_DEV0002A01011:
		return "DEV0002A01011"
	case DeviceInformation_DEV0002A01622:
		return "DEV0002A01622"
	case DeviceInformation_DEV0002A04210:
		return "DEV0002A04210"
	case DeviceInformation_DEV0002A09A13:
		return "DEV0002A09A13"
	case DeviceInformation_DEV00C8272040:
		return "DEV00C8272040"
	case DeviceInformation_DEV00C8272260:
		return "DEV00C8272260"
	case DeviceInformation_DEV00C8272060:
		return "DEV00C8272060"
	case DeviceInformation_DEV00C8272160:
		return "DEV00C8272160"
	case DeviceInformation_DEV00C8272050:
		return "DEV00C8272050"
	case DeviceInformation_DEV00C9106D10:
		return "DEV00C9106D10"
	case DeviceInformation_DEV006B106D10:
		return "DEV006B106D10"
	case DeviceInformation_DEV00C9107C20:
		return "DEV00C9107C20"
	case DeviceInformation_DEV00C9108511:
		return "DEV00C9108511"
	case DeviceInformation_DEV00C9106210:
		return "DEV00C9106210"
	case DeviceInformation_DEV00C9109310:
		return "DEV00C9109310"
	case DeviceInformation_DEV00C9109210:
		return "DEV00C9109210"
	case DeviceInformation_DEV00C9109810:
		return "DEV00C9109810"
	case DeviceInformation_DEV00C9109A10:
		return "DEV00C9109A10"
	case DeviceInformation_DEV00C910A420:
		return "DEV00C910A420"
	case DeviceInformation_DEV00C910A110:
		return "DEV00C910A110"
	case DeviceInformation_DEV00C910A010:
		return "DEV00C910A010"
	case DeviceInformation_DEV0071123130:
		return "DEV0071123130"
	case DeviceInformation_DEV00C910A310:
		return "DEV00C910A310"
	case DeviceInformation_DEV00C910A210:
		return "DEV00C910A210"
	case DeviceInformation_DEV00C9109B10:
		return "DEV00C9109B10"
	case DeviceInformation_DEV00C9106110:
		return "DEV00C9106110"
	case DeviceInformation_DEV00C9109110:
		return "DEV00C9109110"
	case DeviceInformation_DEV00C9109610:
		return "DEV00C9109610"
	case DeviceInformation_DEV00C9109710:
		return "DEV00C9109710"
	case DeviceInformation_DEV00C9109510:
		return "DEV00C9109510"
	case DeviceInformation_DEV00C9109910:
		return "DEV00C9109910"
	case DeviceInformation_DEV00C9109C10:
		return "DEV00C9109C10"
	case DeviceInformation_DEV0071413133:
		return "DEV0071413133"
	case DeviceInformation_DEV00C910AB10:
		return "DEV00C910AB10"
	case DeviceInformation_DEV00C910AC10:
		return "DEV00C910AC10"
	case DeviceInformation_DEV00C910AD10:
		return "DEV00C910AD10"
	case DeviceInformation_DEV00C910A810:
		return "DEV00C910A810"
	case DeviceInformation_DEV00C9106510:
		return "DEV00C9106510"
	case DeviceInformation_DEV00C910A710:
		return "DEV00C910A710"
	case DeviceInformation_DEV00C9107610:
		return "DEV00C9107610"
	case DeviceInformation_DEV00C910890A:
		return "DEV00C910890A"
	case DeviceInformation_DEV00C9FF1012:
		return "DEV00C9FF1012"
	case DeviceInformation_DEV00C9FF0913:
		return "DEV00C9FF0913"
	case DeviceInformation_DEV0071114019:
		return "DEV0071114019"
	case DeviceInformation_DEV00C9FF1112:
		return "DEV00C9FF1112"
	case DeviceInformation_DEV00C9100310:
		return "DEV00C9100310"
	case DeviceInformation_DEV00C9101110:
		return "DEV00C9101110"
	case DeviceInformation_DEV00C9101010:
		return "DEV00C9101010"
	case DeviceInformation_DEV00C9103710:
		return "DEV00C9103710"
	case DeviceInformation_DEV00C9101310:
		return "DEV00C9101310"
	case DeviceInformation_DEV00C9FF0D12:
		return "DEV00C9FF0D12"
	case DeviceInformation_DEV00C9100E10:
		return "DEV00C9100E10"
	case DeviceInformation_DEV00C9100610:
		return "DEV00C9100610"
	case DeviceInformation_DEV00C9100510:
		return "DEV00C9100510"
	case DeviceInformation_DEV0064760210:
		return "DEV0064760210"
	case DeviceInformation_DEV007111306C:
		return "DEV007111306C"
	case DeviceInformation_DEV00C9100710:
		return "DEV00C9100710"
	case DeviceInformation_DEV00C9FF1D20:
		return "DEV00C9FF1D20"
	case DeviceInformation_DEV00C9FF1C10:
		return "DEV00C9FF1C10"
	case DeviceInformation_DEV00C9100810:
		return "DEV00C9100810"
	case DeviceInformation_DEV00C9FF1420:
		return "DEV00C9FF1420"
	case DeviceInformation_DEV00C9100D10:
		return "DEV00C9100D10"
	case DeviceInformation_DEV00C9101220:
		return "DEV00C9101220"
	case DeviceInformation_DEV00C9102330:
		return "DEV00C9102330"
	case DeviceInformation_DEV00C9102130:
		return "DEV00C9102130"
	case DeviceInformation_DEV00C9102430:
		return "DEV00C9102430"
	case DeviceInformation_DEV0071231112:
		return "DEV0071231112"
	case DeviceInformation_DEV00C9100831:
		return "DEV00C9100831"
	case DeviceInformation_DEV00C9102530:
		return "DEV00C9102530"
	case DeviceInformation_DEV00C9100531:
		return "DEV00C9100531"
	case DeviceInformation_DEV00C9102030:
		return "DEV00C9102030"
	case DeviceInformation_DEV00C9100731:
		return "DEV00C9100731"
	case DeviceInformation_DEV00C9100631:
		return "DEV00C9100631"
	case DeviceInformation_DEV00C9102230:
		return "DEV00C9102230"
	case DeviceInformation_DEV00C9100632:
		return "DEV00C9100632"
	case DeviceInformation_DEV00C9100532:
		return "DEV00C9100532"
	case DeviceInformation_DEV00C9100732:
		return "DEV00C9100732"
	case DeviceInformation_DEV0071113080:
		return "DEV0071113080"
	case DeviceInformation_DEV00C9100832:
		return "DEV00C9100832"
	case DeviceInformation_DEV00C9102532:
		return "DEV00C9102532"
	case DeviceInformation_DEV00C9102132:
		return "DEV00C9102132"
	case DeviceInformation_DEV00C9102332:
		return "DEV00C9102332"
	case DeviceInformation_DEV00C9102432:
		return "DEV00C9102432"
	case DeviceInformation_DEV00C9102032:
		return "DEV00C9102032"
	case DeviceInformation_DEV00C9102232:
		return "DEV00C9102232"
	case DeviceInformation_DEV00C9104432:
		return "DEV00C9104432"
	case DeviceInformation_DEV00C9100D11:
		return "DEV00C9100D11"
	case DeviceInformation_DEV00C9100633:
		return "DEV00C9100633"
	case DeviceInformation_DEV0071321212:
		return "DEV0071321212"
	case DeviceInformation_DEV00C9100533:
		return "DEV00C9100533"
	case DeviceInformation_DEV00C9100733:
		return "DEV00C9100733"
	case DeviceInformation_DEV00C9100833:
		return "DEV00C9100833"
	case DeviceInformation_DEV00C9102533:
		return "DEV00C9102533"
	case DeviceInformation_DEV00C9102133:
		return "DEV00C9102133"
	case DeviceInformation_DEV00C9102333:
		return "DEV00C9102333"
	case DeviceInformation_DEV00C9102433:
		return "DEV00C9102433"
	case DeviceInformation_DEV00C9102033:
		return "DEV00C9102033"
	case DeviceInformation_DEV00C9102233:
		return "DEV00C9102233"
	case DeviceInformation_DEV00C9104810:
		return "DEV00C9104810"
	case DeviceInformation_DEV0071321113:
		return "DEV0071321113"
	case DeviceInformation_DEV00C9FF1A11:
		return "DEV00C9FF1A11"
	case DeviceInformation_DEV00C9100212:
		return "DEV00C9100212"
	case DeviceInformation_DEV00C9FF0A11:
		return "DEV00C9FF0A11"
	case DeviceInformation_DEV00C9FF0C12:
		return "DEV00C9FF0C12"
	case DeviceInformation_DEV00C9100112:
		return "DEV00C9100112"
	case DeviceInformation_DEV00C9FF1911:
		return "DEV00C9FF1911"
	case DeviceInformation_DEV00C9FF0B12:
		return "DEV00C9FF0B12"
	case DeviceInformation_DEV00C9FF0715:
		return "DEV00C9FF0715"
	case DeviceInformation_DEV00C9FF1B10:
		return "DEV00C9FF1B10"
	case DeviceInformation_DEV00C9101610:
		return "DEV00C9101610"
	case DeviceInformation_DEV0071322212:
		return "DEV0071322212"
	case DeviceInformation_DEV00C9FF1B11:
		return "DEV00C9FF1B11"
	case DeviceInformation_DEV00C9101611:
		return "DEV00C9101611"
	case DeviceInformation_DEV00C9101612:
		return "DEV00C9101612"
	case DeviceInformation_DEV00C9FF0F11:
		return "DEV00C9FF0F11"
	case DeviceInformation_DEV00C9FF1E30:
		return "DEV00C9FF1E30"
	case DeviceInformation_DEV00C9100410:
		return "DEV00C9100410"
	case DeviceInformation_DEV00C9106410:
		return "DEV00C9106410"
	case DeviceInformation_DEV00C9106710:
		return "DEV00C9106710"
	case DeviceInformation_DEV00C9106810:
		return "DEV00C9106810"
	case DeviceInformation_DEV00C9106010:
		return "DEV00C9106010"
	case DeviceInformation_DEV0071322112:
		return "DEV0071322112"
	case DeviceInformation_DEV00C9106310:
		return "DEV00C9106310"
	case DeviceInformation_DEV00C9107110:
		return "DEV00C9107110"
	case DeviceInformation_DEV00C9107210:
		return "DEV00C9107210"
	case DeviceInformation_DEV00C9107310:
		return "DEV00C9107310"
	case DeviceInformation_DEV00C9107010:
		return "DEV00C9107010"
	case DeviceInformation_DEV00C9107A20:
		return "DEV00C9107A20"
	case DeviceInformation_DEV00C9107B20:
		return "DEV00C9107B20"
	case DeviceInformation_DEV00C9107820:
		return "DEV00C9107820"
	case DeviceInformation_DEV00C9107920:
		return "DEV00C9107920"
	case DeviceInformation_DEV00C9104433:
		return "DEV00C9104433"
	case DeviceInformation_DEV0071322312:
		return "DEV0071322312"
	case DeviceInformation_DEV00C9107C11:
		return "DEV00C9107C11"
	case DeviceInformation_DEV00C9107711:
		return "DEV00C9107711"
	case DeviceInformation_DEV00C9108310:
		return "DEV00C9108310"
	case DeviceInformation_DEV00C9108210:
		return "DEV00C9108210"
	case DeviceInformation_DEV00C9108610:
		return "DEV00C9108610"
	case DeviceInformation_DEV00C9107D10:
		return "DEV00C9107D10"
	case DeviceInformation_DEV00CE648B10:
		return "DEV00CE648B10"
	case DeviceInformation_DEV00CE494513:
		return "DEV00CE494513"
	case DeviceInformation_DEV00CE494311:
		return "DEV00CE494311"
	case DeviceInformation_DEV00CE494810:
		return "DEV00CE494810"
	case DeviceInformation_DEV0071122124:
		return "DEV0071122124"
	case DeviceInformation_DEV00CE494712:
		return "DEV00CE494712"
	case DeviceInformation_DEV00CE494012:
		return "DEV00CE494012"
	case DeviceInformation_DEV00CE494111:
		return "DEV00CE494111"
	case DeviceInformation_DEV00CE494210:
		return "DEV00CE494210"
	case DeviceInformation_DEV00CE494610:
		return "DEV00CE494610"
	case DeviceInformation_DEV00CE494412:
		return "DEV00CE494412"
	case DeviceInformation_DEV00D0660212:
		return "DEV00D0660212"
	case DeviceInformation_DEV00E8000A10:
		return "DEV00E8000A10"
	case DeviceInformation_DEV00E8000B10:
		return "DEV00E8000B10"
	case DeviceInformation_DEV00E8000910:
		return "DEV00E8000910"
	case DeviceInformation_DEV007112221E:
		return "DEV007112221E"
	case DeviceInformation_DEV00E8001112:
		return "DEV00E8001112"
	case DeviceInformation_DEV00E8000C14:
		return "DEV00E8000C14"
	case DeviceInformation_DEV00E8000D13:
		return "DEV00E8000D13"
	case DeviceInformation_DEV00E8000E12:
		return "DEV00E8000E12"
	case DeviceInformation_DEV00E8001310:
		return "DEV00E8001310"
	case DeviceInformation_DEV00E8001410:
		return "DEV00E8001410"
	case DeviceInformation_DEV00E8001510:
		return "DEV00E8001510"
	case DeviceInformation_DEV00E8000F10:
		return "DEV00E8000F10"
	case DeviceInformation_DEV00E8001010:
		return "DEV00E8001010"
	case DeviceInformation_DEV00E8000612:
		return "DEV00E8000612"
	case DeviceInformation_DEV0064182410:
		return "DEV0064182410"
	case DeviceInformation_DEV0071413314:
		return "DEV0071413314"
	case DeviceInformation_DEV00E8000812:
		return "DEV00E8000812"
	case DeviceInformation_DEV00E8000712:
		return "DEV00E8000712"
	case DeviceInformation_DEV00F4501311:
		return "DEV00F4501311"
	case DeviceInformation_DEV00F4B00911:
		return "DEV00F4B00911"
	case DeviceInformation_DEV0019512710:
		return "DEV0019512710"
	case DeviceInformation_DEV0019512810:
		return "DEV0019512810"
	case DeviceInformation_DEV0019512910:
		return "DEV0019512910"
	case DeviceInformation_DEV0019E30D10:
		return "DEV0019E30D10"
	case DeviceInformation_DEV0019512211:
		return "DEV0019512211"
	case DeviceInformation_DEV0019512311:
		return "DEV0019512311"
	case DeviceInformation_DEV0072300110:
		return "DEV0072300110"
	case DeviceInformation_DEV0019512111:
		return "DEV0019512111"
	case DeviceInformation_DEV0019520D11:
		return "DEV0019520D11"
	case DeviceInformation_DEV0019E30B12:
		return "DEV0019E30B12"
	case DeviceInformation_DEV0019530812:
		return "DEV0019530812"
	case DeviceInformation_DEV0019530912:
		return "DEV0019530912"
	case DeviceInformation_DEV0019530612:
		return "DEV0019530612"
	case DeviceInformation_DEV0019530711:
		return "DEV0019530711"
	case DeviceInformation_DEV0019E30A11:
		return "DEV0019E30A11"
	case DeviceInformation_DEV0019E20111:
		return "DEV0019E20111"
	case DeviceInformation_DEV0019E20210:
		return "DEV0019E20210"
	case DeviceInformation_DEV0076002101:
		return "DEV0076002101"
	case DeviceInformation_DEV0019E30C11:
		return "DEV0019E30C11"
	case DeviceInformation_DEV0019E11310:
		return "DEV0019E11310"
	case DeviceInformation_DEV0019E11210:
		return "DEV0019E11210"
	case DeviceInformation_DEV0019E30610:
		return "DEV0019E30610"
	case DeviceInformation_DEV0019E30710:
		return "DEV0019E30710"
	case DeviceInformation_DEV0019E30910:
		return "DEV0019E30910"
	case DeviceInformation_DEV0019E30810:
		return "DEV0019E30810"
	case DeviceInformation_DEV0019E25510:
		return "DEV0019E25510"
	case DeviceInformation_DEV0019E20410:
		return "DEV0019E20410"
	case DeviceInformation_DEV0019E20310:
		return "DEV0019E20310"
	case DeviceInformation_DEV0076002001:
		return "DEV0076002001"
	case DeviceInformation_DEV0019E25610:
		return "DEV0019E25610"
	case DeviceInformation_DEV0019512010:
		return "DEV0019512010"
	case DeviceInformation_DEV0019520C10:
		return "DEV0019520C10"
	case DeviceInformation_DEV0019520710:
		return "DEV0019520710"
	case DeviceInformation_DEV0019520210:
		return "DEV0019520210"
	case DeviceInformation_DEV0019E25010:
		return "DEV0019E25010"
	case DeviceInformation_DEV0019E25110:
		return "DEV0019E25110"
	case DeviceInformation_DEV0019130710:
		return "DEV0019130710"
	case DeviceInformation_DEV0019272050:
		return "DEV0019272050"
	case DeviceInformation_DEV0019520910:
		return "DEV0019520910"
	case DeviceInformation_DEV0076002A15:
		return "DEV0076002A15"
	case DeviceInformation_DEV0019520A10:
		return "DEV0019520A10"
	case DeviceInformation_DEV0019520B10:
		return "DEV0019520B10"
	case DeviceInformation_DEV0019520412:
		return "DEV0019520412"
	case DeviceInformation_DEV0019520812:
		return "DEV0019520812"
	case DeviceInformation_DEV0019512510:
		return "DEV0019512510"
	case DeviceInformation_DEV0019512410:
		return "DEV0019512410"
	case DeviceInformation_DEV0019512610:
		return "DEV0019512610"
	case DeviceInformation_DEV0019511711:
		return "DEV0019511711"
	case DeviceInformation_DEV0019511811:
		return "DEV0019511811"
	case DeviceInformation_DEV0019522212:
		return "DEV0019522212"
	case DeviceInformation_DEV0076002815:
		return "DEV0076002815"
	case DeviceInformation_DEV0019FF0716:
		return "DEV0019FF0716"
	case DeviceInformation_DEV0019FF1420:
		return "DEV0019FF1420"
	case DeviceInformation_DEV0019522112:
		return "DEV0019522112"
	case DeviceInformation_DEV0019522011:
		return "DEV0019522011"
	case DeviceInformation_DEV0019522311:
		return "DEV0019522311"
	case DeviceInformation_DEV0019E12410:
		return "DEV0019E12410"
	case DeviceInformation_DEV0019000311:
		return "DEV0019000311"
	case DeviceInformation_DEV0019000411:
		return "DEV0019000411"
	case DeviceInformation_DEV0019070210:
		return "DEV0019070210"
	case DeviceInformation_DEV0019070E11:
		return "DEV0019070E11"
	case DeviceInformation_DEV0076002215:
		return "DEV0076002215"
	case DeviceInformation_DEV0019724010:
		return "DEV0019724010"
	case DeviceInformation_DEV0019520610:
		return "DEV0019520610"
	case DeviceInformation_DEV0019520510:
		return "DEV0019520510"
	case DeviceInformation_DEV00FB101111:
		return "DEV00FB101111"
	case DeviceInformation_DEV00FB103001:
		return "DEV00FB103001"
	case DeviceInformation_DEV00FB104401:
		return "DEV00FB104401"
	case DeviceInformation_DEV00FB124002:
		return "DEV00FB124002"
	case DeviceInformation_DEV00FB104102:
		return "DEV00FB104102"
	case DeviceInformation_DEV00FB104201:
		return "DEV00FB104201"
	case DeviceInformation_DEV00FBF77603:
		return "DEV00FBF77603"
	case DeviceInformation_DEV0076002B15:
		return "DEV0076002B15"
	case DeviceInformation_DEV00FB104301:
		return "DEV00FB104301"
	case DeviceInformation_DEV00FB104601:
		return "DEV00FB104601"
	case DeviceInformation_DEV00FB104701:
		return "DEV00FB104701"
	case DeviceInformation_DEV00FB105101:
		return "DEV00FB105101"
	case DeviceInformation_DEV0103030110:
		return "DEV0103030110"
	case DeviceInformation_DEV0103010113:
		return "DEV0103010113"
	case DeviceInformation_DEV0103090110:
		return "DEV0103090110"
	case DeviceInformation_DEV0103020111:
		return "DEV0103020111"
	case DeviceInformation_DEV0103020112:
		return "DEV0103020112"
	case DeviceInformation_DEV0103040110:
		return "DEV0103040110"
	case DeviceInformation_DEV0076002715:
		return "DEV0076002715"
	case DeviceInformation_DEV0103050111:
		return "DEV0103050111"
	case DeviceInformation_DEV0107000301:
		return "DEV0107000301"
	case DeviceInformation_DEV0107000101:
		return "DEV0107000101"
	case DeviceInformation_DEV0107000201:
		return "DEV0107000201"
	case DeviceInformation_DEV0107020801:
		return "DEV0107020801"
	case DeviceInformation_DEV0107020401:
		return "DEV0107020401"
	case DeviceInformation_DEV0107020001:
		return "DEV0107020001"
	case DeviceInformation_DEV010701F801:
		return "DEV010701F801"
	case DeviceInformation_DEV010701FC01:
		return "DEV010701FC01"
	case DeviceInformation_DEV0107020C01:
		return "DEV0107020C01"
	case DeviceInformation_DEV0076002315:
		return "DEV0076002315"
	case DeviceInformation_DEV010F100801:
		return "DEV010F100801"
	case DeviceInformation_DEV010F100601:
		return "DEV010F100601"
	case DeviceInformation_DEV010F100401:
		return "DEV010F100401"
	case DeviceInformation_DEV010F030601:
		return "DEV010F030601"
	case DeviceInformation_DEV010F010301:
		return "DEV010F010301"
	case DeviceInformation_DEV010F010101:
		return "DEV010F010101"
	case DeviceInformation_DEV010F010201:
		return "DEV010F010201"
	case DeviceInformation_DEV010F000302:
		return "DEV010F000302"
	case DeviceInformation_DEV010F000402:
		return "DEV010F000402"
	case DeviceInformation_DEV010F000102:
		return "DEV010F000102"
	case DeviceInformation_DEV0064182310:
		return "DEV0064182310"
	case DeviceInformation_DEV0076002415:
		return "DEV0076002415"
	case DeviceInformation_DEV011EBB8211:
		return "DEV011EBB8211"
	case DeviceInformation_DEV011E108111:
		return "DEV011E108111"
	case DeviceInformation_DEV0123010010:
		return "DEV0123010010"
	case DeviceInformation_DEV001E478010:
		return "DEV001E478010"
	case DeviceInformation_DEV001E706611:
		return "DEV001E706611"
	case DeviceInformation_DEV001E706811:
		return "DEV001E706811"
	case DeviceInformation_DEV001E473012:
		return "DEV001E473012"
	case DeviceInformation_DEV001E20A011:
		return "DEV001E20A011"
	case DeviceInformation_DEV001E209011:
		return "DEV001E209011"
	case DeviceInformation_DEV001E209811:
		return "DEV001E209811"
	case DeviceInformation_DEV0076002615:
		return "DEV0076002615"
	case DeviceInformation_DEV001E208811:
		return "DEV001E208811"
	case DeviceInformation_DEV001E208011:
		return "DEV001E208011"
	case DeviceInformation_DEV001E207821:
		return "DEV001E207821"
	case DeviceInformation_DEV001E20CA12:
		return "DEV001E20CA12"
	case DeviceInformation_DEV001E20B312:
		return "DEV001E20B312"
	case DeviceInformation_DEV001E20B012:
		return "DEV001E20B012"
	case DeviceInformation_DEV001E302612:
		return "DEV001E302612"
	case DeviceInformation_DEV001E302312:
		return "DEV001E302312"
	case DeviceInformation_DEV001E302012:
		return "DEV001E302012"
	case DeviceInformation_DEV001E20A811:
		return "DEV001E20A811"
	case DeviceInformation_DEV0076002515:
		return "DEV0076002515"
	case DeviceInformation_DEV001E20C412:
		return "DEV001E20C412"
	case DeviceInformation_DEV001E20C712:
		return "DEV001E20C712"
	case DeviceInformation_DEV001E20AD12:
		return "DEV001E20AD12"
	case DeviceInformation_DEV001E443720:
		return "DEV001E443720"
	case DeviceInformation_DEV001E441821:
		return "DEV001E441821"
	case DeviceInformation_DEV001E443810:
		return "DEV001E443810"
	case DeviceInformation_DEV001E140C12:
		return "DEV001E140C12"
	case DeviceInformation_DEV001E471611:
		return "DEV001E471611"
	case DeviceInformation_DEV001E479024:
		return "DEV001E479024"
	case DeviceInformation_DEV001E471A11:
		return "DEV001E471A11"
	case DeviceInformation_DEV0076000201:
		return "DEV0076000201"
	case DeviceInformation_DEV001E477A10:
		return "DEV001E477A10"
	case DeviceInformation_DEV001E470A11:
		return "DEV001E470A11"
	case DeviceInformation_DEV001E480B11:
		return "DEV001E480B11"
	case DeviceInformation_DEV001E487B10:
		return "DEV001E487B10"
	case DeviceInformation_DEV001E440411:
		return "DEV001E440411"
	case DeviceInformation_DEV001E447211:
		return "DEV001E447211"
	case DeviceInformation_DEV0142010010:
		return "DEV0142010010"
	case DeviceInformation_DEV0142010011:
		return "DEV0142010011"
	case DeviceInformation_DEV017A130401:
		return "DEV017A130401"
	case DeviceInformation_DEV017A130201:
		return "DEV017A130201"
	case DeviceInformation_DEV0076000101:
		return "DEV0076000101"
	case DeviceInformation_DEV017A130801:
		return "DEV017A130801"
	case DeviceInformation_DEV017A130601:
		return "DEV017A130601"
	case DeviceInformation_DEV017A300102:
		return "DEV017A300102"
	case DeviceInformation_DEV0193323C11:
		return "DEV0193323C11"
	case DeviceInformation_DEV0198101110:
		return "DEV0198101110"
	case DeviceInformation_DEV01C4030110:
		return "DEV01C4030110"
	case DeviceInformation_DEV01C4030210:
		return "DEV01C4030210"
	case DeviceInformation_DEV01C4021210:
		return "DEV01C4021210"
	case DeviceInformation_DEV01C4001010:
		return "DEV01C4001010"
	case DeviceInformation_DEV01C4020610:
		return "DEV01C4020610"
	case DeviceInformation_DEV0076000301:
		return "DEV0076000301"
	case DeviceInformation_DEV01C4020910:
		return "DEV01C4020910"
	case DeviceInformation_DEV01C4020810:
		return "DEV01C4020810"
	case DeviceInformation_DEV01C4010710:
		return "DEV01C4010710"
	case DeviceInformation_DEV01C4050210:
		return "DEV01C4050210"
	case DeviceInformation_DEV01C4010810:
		return "DEV01C4010810"
	case DeviceInformation_DEV01C4020510:
		return "DEV01C4020510"
	case DeviceInformation_DEV01C4040110:
		return "DEV01C4040110"
	case DeviceInformation_DEV01C4040310:
		return "DEV01C4040310"
	case DeviceInformation_DEV01C4040210:
		return "DEV01C4040210"
	case DeviceInformation_DEV01C4101210:
		return "DEV01C4101210"
	case DeviceInformation_DEV0076000401:
		return "DEV0076000401"
	case DeviceInformation_DEV003D020109:
		return "DEV003D020109"
	case DeviceInformation_DEV01DB000301:
		return "DEV01DB000301"
	case DeviceInformation_DEV01DB000201:
		return "DEV01DB000201"
	case DeviceInformation_DEV01DB000401:
		return "DEV01DB000401"
	case DeviceInformation_DEV01DB000801:
		return "DEV01DB000801"
	case DeviceInformation_DEV01DB001201:
		return "DEV01DB001201"
	case DeviceInformation_DEV009A000400:
		return "DEV009A000400"
	case DeviceInformation_DEV009A100400:
		return "DEV009A100400"
	case DeviceInformation_DEV009A200C00:
		return "DEV009A200C00"
	case DeviceInformation_DEV009A200E00:
		return "DEV009A200E00"
	case DeviceInformation_DEV0076002903:
		return "DEV0076002903"
	case DeviceInformation_DEV009A000201:
		return "DEV009A000201"
	case DeviceInformation_DEV009A000300:
		return "DEV009A000300"
	case DeviceInformation_DEV009A00C000:
		return "DEV009A00C000"
	case DeviceInformation_DEV009A00B000:
		return "DEV009A00B000"
	case DeviceInformation_DEV009A00C002:
		return "DEV009A00C002"
	case DeviceInformation_DEV009A200100:
		return "DEV009A200100"
	case DeviceInformation_DEV004E400010:
		return "DEV004E400010"
	case DeviceInformation_DEV004E030031:
		return "DEV004E030031"
	case DeviceInformation_DEV012B010110:
		return "DEV012B010110"
	case DeviceInformation_DEV01F6E0E110:
		return "DEV01F6E0E110"
	case DeviceInformation_DEV0076002901:
		return "DEV0076002901"
	case DeviceInformation_DEV0088100010:
		return "DEV0088100010"
	case DeviceInformation_DEV0088100210:
		return "DEV0088100210"
	case DeviceInformation_DEV0088100110:
		return "DEV0088100110"
	case DeviceInformation_DEV0088110010:
		return "DEV0088110010"
	case DeviceInformation_DEV0088120412:
		return "DEV0088120412"
	case DeviceInformation_DEV0088120113:
		return "DEV0088120113"
	case DeviceInformation_DEV011A4B5201:
		return "DEV011A4B5201"
	case DeviceInformation_DEV008B020301:
		return "DEV008B020301"
	case DeviceInformation_DEV008B010610:
		return "DEV008B010610"
	case DeviceInformation_DEV008B030110:
		return "DEV008B030110"
	case DeviceInformation_DEV007600E503:
		return "DEV007600E503"
	case DeviceInformation_DEV008B030310:
		return "DEV008B030310"
	case DeviceInformation_DEV008B030210:
		return "DEV008B030210"
	case DeviceInformation_DEV008B031512:
		return "DEV008B031512"
	case DeviceInformation_DEV008B031412:
		return "DEV008B031412"
	case DeviceInformation_DEV008B031312:
		return "DEV008B031312"
	case DeviceInformation_DEV008B031212:
		return "DEV008B031212"
	case DeviceInformation_DEV008B031112:
		return "DEV008B031112"
	case DeviceInformation_DEV008B031012:
		return "DEV008B031012"
	case DeviceInformation_DEV008B030510:
		return "DEV008B030510"
	case DeviceInformation_DEV008B030410:
		return "DEV008B030410"
	case DeviceInformation_DEV0064705C01:
		return "DEV0064705C01"
	case DeviceInformation_DEV0076004002:
		return "DEV0076004002"
	case DeviceInformation_DEV008B020310:
		return "DEV008B020310"
	case DeviceInformation_DEV008B020210:
		return "DEV008B020210"
	case DeviceInformation_DEV008B020110:
		return "DEV008B020110"
	case DeviceInformation_DEV008B010110:
		return "DEV008B010110"
	case DeviceInformation_DEV008B010210:
		return "DEV008B010210"
	case DeviceInformation_DEV008B010310:
		return "DEV008B010310"
	case DeviceInformation_DEV008B010410:
		return "DEV008B010410"
	case DeviceInformation_DEV008B040110:
		return "DEV008B040110"
	case DeviceInformation_DEV008B040210:
		return "DEV008B040210"
	case DeviceInformation_DEV008B010910:
		return "DEV008B010910"
	case DeviceInformation_DEV0076004003:
		return "DEV0076004003"
	case DeviceInformation_DEV008B010710:
		return "DEV008B010710"
	case DeviceInformation_DEV008B010810:
		return "DEV008B010810"
	case DeviceInformation_DEV008B041111:
		return "DEV008B041111"
	case DeviceInformation_DEV008B041211:
		return "DEV008B041211"
	case DeviceInformation_DEV008B041311:
		return "DEV008B041311"
	case DeviceInformation_DEV00A600020A:
		return "DEV00A600020A"
	case DeviceInformation_DEV00A6000B10:
		return "DEV00A6000B10"
	case DeviceInformation_DEV00A6000300:
		return "DEV00A6000300"
	case DeviceInformation_DEV00A6000705:
		return "DEV00A6000705"
	case DeviceInformation_DEV00A6000605:
		return "DEV00A6000605"
	case DeviceInformation_DEV0076003402:
		return "DEV0076003402"
	case DeviceInformation_DEV00A6000500:
		return "DEV00A6000500"
	case DeviceInformation_DEV00A6000C10:
		return "DEV00A6000C10"
	case DeviceInformation_DEV002CA01811:
		return "DEV002CA01811"
	case DeviceInformation_DEV002CA07033:
		return "DEV002CA07033"
	case DeviceInformation_DEV002C555020:
		return "DEV002C555020"
	case DeviceInformation_DEV002C556421:
		return "DEV002C556421"
	case DeviceInformation_DEV002C05F211:
		return "DEV002C05F211"
	case DeviceInformation_DEV002C05F411:
		return "DEV002C05F411"
	case DeviceInformation_DEV002C05E613:
		return "DEV002C05E613"
	case DeviceInformation_DEV002CA07914:
		return "DEV002CA07914"
	case DeviceInformation_DEV0076003401:
		return "DEV0076003401"
	case DeviceInformation_DEV002C060A13:
		return "DEV002C060A13"
	case DeviceInformation_DEV002C3A0212:
		return "DEV002C3A0212"
	case DeviceInformation_DEV002C060210:
		return "DEV002C060210"
	case DeviceInformation_DEV002CA00213:
		return "DEV002CA00213"
	case DeviceInformation_DEV002CA0A611:
		return "DEV002CA0A611"
	case DeviceInformation_DEV002CA07B11:
		return "DEV002CA07B11"
	case DeviceInformation_DEV002C063210:
		return "DEV002C063210"
	case DeviceInformation_DEV002C063110:
		return "DEV002C063110"
	case DeviceInformation_DEV002C062E10:
		return "DEV002C062E10"
	case DeviceInformation_DEV002C062C10:
		return "DEV002C062C10"
	case DeviceInformation_DEV007600E908:
		return "DEV007600E908"
	case DeviceInformation_DEV002C062D10:
		return "DEV002C062D10"
	case DeviceInformation_DEV002C063310:
		return "DEV002C063310"
	case DeviceInformation_DEV002C05EB10:
		return "DEV002C05EB10"
	case DeviceInformation_DEV002C05F110:
		return "DEV002C05F110"
	case DeviceInformation_DEV002C0B8830:
		return "DEV002C0B8830"
	case DeviceInformation_DEV00A0B07101:
		return "DEV00A0B07101"
	case DeviceInformation_DEV00A0B07001:
		return "DEV00A0B07001"
	case DeviceInformation_DEV00A0B07203:
		return "DEV00A0B07203"
	case DeviceInformation_DEV00A0B02101:
		return "DEV00A0B02101"
	case DeviceInformation_DEV00A0B02401:
		return "DEV00A0B02401"
	case DeviceInformation_DEV007600E907:
		return "DEV007600E907"
	case DeviceInformation_DEV00A0B02301:
		return "DEV00A0B02301"
	case DeviceInformation_DEV00A0B02601:
		return "DEV00A0B02601"
	case DeviceInformation_DEV00A0B02201:
		return "DEV00A0B02201"
	case DeviceInformation_DEV00A0B01902:
		return "DEV00A0B01902"
	case DeviceInformation_DEV0004147112:
		return "DEV0004147112"
	case DeviceInformation_DEV000410A411:
		return "DEV000410A411"
	case DeviceInformation_DEV0004109911:
		return "DEV0004109911"
	case DeviceInformation_DEV0004109912:
		return "DEV0004109912"
	case DeviceInformation_DEV0004109913:
		return "DEV0004109913"
	case DeviceInformation_DEV0004109914:
		return "DEV0004109914"
	case DeviceInformation_DEV000C181710:
		return "DEV000C181710"
	case DeviceInformation_DEV000410A211:
		return "DEV000410A211"
	case DeviceInformation_DEV000410FC12:
		return "DEV000410FC12"
	case DeviceInformation_DEV000410FD12:
		return "DEV000410FD12"
	case DeviceInformation_DEV000410B212:
		return "DEV000410B212"
	case DeviceInformation_DEV0004110B11:
		return "DEV0004110B11"
	case DeviceInformation_DEV0004110711:
		return "DEV0004110711"
	case DeviceInformation_DEV000410B213:
		return "DEV000410B213"
	case DeviceInformation_DEV0004109811:
		return "DEV0004109811"
	case DeviceInformation_DEV0004109812:
		return "DEV0004109812"
	case DeviceInformation_DEV0004109813:
		return "DEV0004109813"
	case DeviceInformation_DEV000C130510:
		return "DEV000C130510"
	case DeviceInformation_DEV0004109814:
		return "DEV0004109814"
	case DeviceInformation_DEV000410A011:
		return "DEV000410A011"
	case DeviceInformation_DEV000410A111:
		return "DEV000410A111"
	case DeviceInformation_DEV000410FA12:
		return "DEV000410FA12"
	case DeviceInformation_DEV000410FB12:
		return "DEV000410FB12"
	case DeviceInformation_DEV000410B112:
		return "DEV000410B112"
	case DeviceInformation_DEV0004110A11:
		return "DEV0004110A11"
	case DeviceInformation_DEV0004110611:
		return "DEV0004110611"
	case DeviceInformation_DEV000410B113:
		return "DEV000410B113"
	case DeviceInformation_DEV0004109A11:
		return "DEV0004109A11"
	case DeviceInformation_DEV000C130610:
		return "DEV000C130610"
	case DeviceInformation_DEV0004109A12:
		return "DEV0004109A12"
	case DeviceInformation_DEV0004109A13:
		return "DEV0004109A13"
	case DeviceInformation_DEV0004109A14:
		return "DEV0004109A14"
	case DeviceInformation_DEV000410A311:
		return "DEV000410A311"
	case DeviceInformation_DEV000410B312:
		return "DEV000410B312"
	case DeviceInformation_DEV0004110C11:
		return "DEV0004110C11"
	case DeviceInformation_DEV0004110811:
		return "DEV0004110811"
	case DeviceInformation_DEV000410B313:
		return "DEV000410B313"
	case DeviceInformation_DEV0004109B11:
		return "DEV0004109B11"
	case DeviceInformation_DEV0004109B12:
		return "DEV0004109B12"
	case DeviceInformation_DEV000C133610:
		return "DEV000C133610"
	case DeviceInformation_DEV0004109B13:
		return "DEV0004109B13"
	case DeviceInformation_DEV0004109B14:
		return "DEV0004109B14"
	case DeviceInformation_DEV000410A511:
		return "DEV000410A511"
	case DeviceInformation_DEV000410B412:
		return "DEV000410B412"
	case DeviceInformation_DEV0004110D11:
		return "DEV0004110D11"
	case DeviceInformation_DEV0004110911:
		return "DEV0004110911"
	case DeviceInformation_DEV000410B413:
		return "DEV000410B413"
	case DeviceInformation_DEV0004109C11:
		return "DEV0004109C11"
	case DeviceInformation_DEV0004109C12:
		return "DEV0004109C12"
	case DeviceInformation_DEV0004109C13:
		return "DEV0004109C13"
	}
	return ""
}
