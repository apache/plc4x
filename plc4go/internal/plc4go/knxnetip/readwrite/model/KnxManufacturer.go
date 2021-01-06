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

type KnxManufacturer uint16

type IKnxManufacturer interface {
	Text() string
	Serialize(io utils.WriteBuffer) error
}

const (
	KnxManufacturer_SIEMENS                                            KnxManufacturer = 1
	KnxManufacturer_ABB                                                KnxManufacturer = 2
	KnxManufacturer_ALBRECHT_JUNG                                      KnxManufacturer = 4
	KnxManufacturer_BTICINO                                            KnxManufacturer = 5
	KnxManufacturer_BERKER                                             KnxManufacturer = 6
	KnxManufacturer_BUSCH_JAEGER_ELEKTRO                               KnxManufacturer = 7
	KnxManufacturer_GIRA_GIERSIEPEN                                    KnxManufacturer = 8
	KnxManufacturer_HAGER_ELECTRO                                      KnxManufacturer = 9
	KnxManufacturer_INSTA_GMBH                                         KnxManufacturer = 10
	KnxManufacturer_LEGRAND_APPAREILLAGE_ELECTRIQUE                    KnxManufacturer = 11
	KnxManufacturer_MERTEN                                             KnxManufacturer = 12
	KnxManufacturer_ABB_SPA_SACE_DIVISION                              KnxManufacturer = 14
	KnxManufacturer_SIEDLE_AND_SOEHNE                                  KnxManufacturer = 22
	KnxManufacturer_EBERLE                                             KnxManufacturer = 24
	KnxManufacturer_GEWISS                                             KnxManufacturer = 25
	KnxManufacturer_ALBERT_ACKERMANN                                   KnxManufacturer = 27
	KnxManufacturer_SCHUPA_GMBH                                        KnxManufacturer = 28
	KnxManufacturer_ABB_SCHWEIZ                                        KnxManufacturer = 29
	KnxManufacturer_FELLER                                             KnxManufacturer = 30
	KnxManufacturer_GLAMOX_AS                                          KnxManufacturer = 31
	KnxManufacturer_DEHN_AND_SOEHNE                                    KnxManufacturer = 32
	KnxManufacturer_CRABTREE                                           KnxManufacturer = 33
	KnxManufacturer_EVOKNX                                             KnxManufacturer = 34
	KnxManufacturer_PAUL_HOCHKOEPPER                                   KnxManufacturer = 36
	KnxManufacturer_ALTENBURGER_ELECTRONIC                             KnxManufacturer = 37
	KnxManufacturer_GRAESSLIN                                          KnxManufacturer = 41
	KnxManufacturer_SIMON_42                                           KnxManufacturer = 42
	KnxManufacturer_VIMAR                                              KnxManufacturer = 44
	KnxManufacturer_MOELLER_GEBAEUDEAUTOMATION_KG                      KnxManufacturer = 45
	KnxManufacturer_ELTAKO                                             KnxManufacturer = 46
	KnxManufacturer_BOSCH_SIEMENS_HAUSHALTSGERAETE                     KnxManufacturer = 49
	KnxManufacturer_RITTO_GMBHANDCO_KG                                 KnxManufacturer = 52
	KnxManufacturer_POWER_CONTROLS                                     KnxManufacturer = 53
	KnxManufacturer_ZUMTOBEL                                           KnxManufacturer = 55
	KnxManufacturer_PHOENIX_CONTACT                                    KnxManufacturer = 57
	KnxManufacturer_WAGO_KONTAKTTECHNIK                                KnxManufacturer = 61
	KnxManufacturer_KNXPRESSO                                          KnxManufacturer = 62
	KnxManufacturer_WIELAND_ELECTRIC                                   KnxManufacturer = 66
	KnxManufacturer_HERMANN_KLEINHUIS                                  KnxManufacturer = 67
	KnxManufacturer_STIEBEL_ELTRON                                     KnxManufacturer = 69
	KnxManufacturer_TEHALIT                                            KnxManufacturer = 71
	KnxManufacturer_THEBEN_AG                                          KnxManufacturer = 72
	KnxManufacturer_WILHELM_RUTENBECK                                  KnxManufacturer = 73
	KnxManufacturer_WINKHAUS                                           KnxManufacturer = 75
	KnxManufacturer_ROBERT_BOSCH                                       KnxManufacturer = 76
	KnxManufacturer_SOMFY                                              KnxManufacturer = 78
	KnxManufacturer_WOERTZ                                             KnxManufacturer = 80
	KnxManufacturer_VIESSMANN_WERKE                                    KnxManufacturer = 81
	KnxManufacturer_IMI_HYDRONIC_ENGINEERING                           KnxManufacturer = 82
	KnxManufacturer_JOH__VAILLANT                                      KnxManufacturer = 83
	KnxManufacturer_AMP_DEUTSCHLAND                                    KnxManufacturer = 85
	KnxManufacturer_BOSCH_THERMOTECHNIK_GMBH                           KnxManufacturer = 89
	KnxManufacturer_SEF___ECOTEC                                       KnxManufacturer = 90
	KnxManufacturer_DORMA_GMBH_Plus_CO__KG                             KnxManufacturer = 92
	KnxManufacturer_WINDOWMASTER_AS                                    KnxManufacturer = 93
	KnxManufacturer_WALTHER_WERKE                                      KnxManufacturer = 94
	KnxManufacturer_ORAS                                               KnxManufacturer = 95
	KnxManufacturer_DAETWYLER                                          KnxManufacturer = 97
	KnxManufacturer_ELECTRAK                                           KnxManufacturer = 98
	KnxManufacturer_TECHEM                                             KnxManufacturer = 99
	KnxManufacturer_SCHNEIDER_ELECTRIC_INDUSTRIES_SAS                  KnxManufacturer = 100
	KnxManufacturer_WHD_WILHELM_HUBER_Plus_SOEHNE                      KnxManufacturer = 101
	KnxManufacturer_BISCHOFF_ELEKTRONIK                                KnxManufacturer = 102
	KnxManufacturer_JEPAZ                                              KnxManufacturer = 104
	KnxManufacturer_RTS_AUTOMATION                                     KnxManufacturer = 105
	KnxManufacturer_EIBMARKT_GMBH                                      KnxManufacturer = 106
	KnxManufacturer_WAREMA_RENKHOFF_SE                                 KnxManufacturer = 107
	KnxManufacturer_EELECTRON                                          KnxManufacturer = 108
	KnxManufacturer_BELDEN_WIRE_AND_CABLE_B_V_                         KnxManufacturer = 109
	KnxManufacturer_BECKER_ANTRIEBE_GMBH                               KnxManufacturer = 110
	KnxManufacturer_J_STEHLEPlusSOEHNE_GMBH                            KnxManufacturer = 111
	KnxManufacturer_AGFEO                                              KnxManufacturer = 112
	KnxManufacturer_ZENNIO                                             KnxManufacturer = 113
	KnxManufacturer_TAPKO_TECHNOLOGIES                                 KnxManufacturer = 114
	KnxManufacturer_HDL                                                KnxManufacturer = 115
	KnxManufacturer_UPONOR                                             KnxManufacturer = 116
	KnxManufacturer_SE_LIGHTMANAGEMENT_AG                              KnxManufacturer = 117
	KnxManufacturer_ARCUS_EDS                                          KnxManufacturer = 118
	KnxManufacturer_INTESIS                                            KnxManufacturer = 119
	KnxManufacturer_HERHOLDT_CONTROLS_SRL                              KnxManufacturer = 120
	KnxManufacturer_NIKO_ZUBLIN                                        KnxManufacturer = 121
	KnxManufacturer_DURABLE_TECHNOLOGIES                               KnxManufacturer = 122
	KnxManufacturer_INNOTEAM                                           KnxManufacturer = 123
	KnxManufacturer_ISE_GMBH                                           KnxManufacturer = 124
	KnxManufacturer_TEAM_FOR_TRONICS                                   KnxManufacturer = 125
	KnxManufacturer_CIAT                                               KnxManufacturer = 126
	KnxManufacturer_REMEHA_BV                                          KnxManufacturer = 127
	KnxManufacturer_ESYLUX                                             KnxManufacturer = 128
	KnxManufacturer_BASALTE                                            KnxManufacturer = 129
	KnxManufacturer_VESTAMATIC                                         KnxManufacturer = 130
	KnxManufacturer_MDT_TECHNOLOGIES                                   KnxManufacturer = 131
	KnxManufacturer_WARENDORFER_KUECHEN_GMBH                           KnxManufacturer = 132
	KnxManufacturer_VIDEO_STAR                                         KnxManufacturer = 133
	KnxManufacturer_SITEK                                              KnxManufacturer = 134
	KnxManufacturer_CONTROLTRONIC                                      KnxManufacturer = 135
	KnxManufacturer_FUNCTION_TECHNOLOGY                                KnxManufacturer = 136
	KnxManufacturer_AMX                                                KnxManufacturer = 137
	KnxManufacturer_ELDAT                                              KnxManufacturer = 138
	KnxManufacturer_PANASONIC                                          KnxManufacturer = 139
	KnxManufacturer_PULSE_TECHNOLOGIES                                 KnxManufacturer = 140
	KnxManufacturer_CRESTRON                                           KnxManufacturer = 141
	KnxManufacturer_STEINEL_PROFESSIONAL                               KnxManufacturer = 142
	KnxManufacturer_BILTON_LED_LIGHTING                                KnxManufacturer = 143
	KnxManufacturer_DENRO_AG                                           KnxManufacturer = 144
	KnxManufacturer_GEPRO                                              KnxManufacturer = 145
	KnxManufacturer_PREUSSEN_AUTOMATION                                KnxManufacturer = 146
	KnxManufacturer_ZOPPAS_INDUSTRIES                                  KnxManufacturer = 147
	KnxManufacturer_MACTECH                                            KnxManufacturer = 148
	KnxManufacturer_TECHNO_TREND                                       KnxManufacturer = 149
	KnxManufacturer_FS_CABLES                                          KnxManufacturer = 150
	KnxManufacturer_DELTA_DORE                                         KnxManufacturer = 151
	KnxManufacturer_EISSOUND                                           KnxManufacturer = 152
	KnxManufacturer_CISCO                                              KnxManufacturer = 153
	KnxManufacturer_DINUY                                              KnxManufacturer = 154
	KnxManufacturer_IKNIX                                              KnxManufacturer = 155
	KnxManufacturer_RADEMACHER_GERAETE_ELEKTRONIK_GMBH                 KnxManufacturer = 156
	KnxManufacturer_EGI_ELECTROACUSTICA_GENERAL_IBERICA                KnxManufacturer = 157
	KnxManufacturer_BES___INGENIUM                                     KnxManufacturer = 158
	KnxManufacturer_ELABNET                                            KnxManufacturer = 159
	KnxManufacturer_BLUMOTIX                                           KnxManufacturer = 160
	KnxManufacturer_HUNTER_DOUGLAS                                     KnxManufacturer = 161
	KnxManufacturer_APRICUM                                            KnxManufacturer = 162
	KnxManufacturer_TIANSU_AUTOMATION                                  KnxManufacturer = 163
	KnxManufacturer_BUBENDORFF                                         KnxManufacturer = 164
	KnxManufacturer_MBS_GMBH                                           KnxManufacturer = 165
	KnxManufacturer_ENERTEX_BAYERN_GMBH                                KnxManufacturer = 166
	KnxManufacturer_BMS                                                KnxManufacturer = 167
	KnxManufacturer_SINAPSI                                            KnxManufacturer = 168
	KnxManufacturer_EMBEDDED_SYSTEMS_SIA                               KnxManufacturer = 169
	KnxManufacturer_KNX1                                               KnxManufacturer = 170
	KnxManufacturer_TOKKA                                              KnxManufacturer = 171
	KnxManufacturer_NANOSENSE                                          KnxManufacturer = 172
	KnxManufacturer_PEAR_AUTOMATION_GMBH                               KnxManufacturer = 173
	KnxManufacturer_DGA                                                KnxManufacturer = 174
	KnxManufacturer_LUTRON                                             KnxManufacturer = 175
	KnxManufacturer_AIRZONE___ALTRA                                    KnxManufacturer = 176
	KnxManufacturer_LITHOSS_DESIGN_SWITCHES                            KnxManufacturer = 177
	KnxManufacturer_3ATEL                                              KnxManufacturer = 178
	KnxManufacturer_PHILIPS_CONTROLS                                   KnxManufacturer = 179
	KnxManufacturer_VELUX_AS                                           KnxManufacturer = 180
	KnxManufacturer_LOYTEC                                             KnxManufacturer = 181
	KnxManufacturer_EKINEX_S_P_A_                                      KnxManufacturer = 182
	KnxManufacturer_SIRLAN_TECHNOLOGIES                                KnxManufacturer = 183
	KnxManufacturer_PROKNX_SAS                                         KnxManufacturer = 184
	KnxManufacturer_IT_GMBH                                            KnxManufacturer = 185
	KnxManufacturer_RENSON                                             KnxManufacturer = 186
	KnxManufacturer_HEP_GROUP                                          KnxManufacturer = 187
	KnxManufacturer_BALMART                                            KnxManufacturer = 188
	KnxManufacturer_GFS_GMBH                                           KnxManufacturer = 189
	KnxManufacturer_SCHENKER_STOREN_AG                                 KnxManufacturer = 190
	KnxManufacturer_ALGODUE_ELETTRONICA_S_R_L_                         KnxManufacturer = 191
	KnxManufacturer_ABB_FRANCE                                         KnxManufacturer = 192
	KnxManufacturer_MAINTRONIC                                         KnxManufacturer = 193
	KnxManufacturer_VANTAGE                                            KnxManufacturer = 194
	KnxManufacturer_FORESIS                                            KnxManufacturer = 195
	KnxManufacturer_RESEARCH_AND_PRODUCTION_ASSOCIATION_SEM            KnxManufacturer = 196
	KnxManufacturer_WEINZIERL_ENGINEERING_GMBH                         KnxManufacturer = 197
	KnxManufacturer_MOEHLENHOFF_WAERMETECHNIK_GMBH                     KnxManufacturer = 198
	KnxManufacturer_PKC_GROUP_OYJ                                      KnxManufacturer = 199
	KnxManufacturer_B_E_G_                                             KnxManufacturer = 200
	KnxManufacturer_ELSNER_ELEKTRONIK_GMBH                             KnxManufacturer = 201
	KnxManufacturer_SIEMENS_BUILDING_TECHNOLOGIES_HKCHINA_LTD_         KnxManufacturer = 202
	KnxManufacturer_EUTRAC                                             KnxManufacturer = 204
	KnxManufacturer_GUSTAV_HENSEL_GMBH_AND_CO__KG                      KnxManufacturer = 205
	KnxManufacturer_GARO_AB                                            KnxManufacturer = 206
	KnxManufacturer_WALDMANN_LICHTTECHNIK                              KnxManufacturer = 207
	KnxManufacturer_SCHUECO                                            KnxManufacturer = 208
	KnxManufacturer_EMU                                                KnxManufacturer = 209
	KnxManufacturer_JNET_SYSTEMS_AG                                    KnxManufacturer = 210
	KnxManufacturer_TOTAL_SOLUTION_GMBH                                KnxManufacturer = 211
	KnxManufacturer_O_Y_L__ELECTRONICS                                 KnxManufacturer = 214
	KnxManufacturer_GALAX_SYSTEM                                       KnxManufacturer = 215
	KnxManufacturer_DISCH                                              KnxManufacturer = 216
	KnxManufacturer_AUCOTEAM                                           KnxManufacturer = 217
	KnxManufacturer_LUXMATE_CONTROLS                                   KnxManufacturer = 218
	KnxManufacturer_DANFOSS                                            KnxManufacturer = 219
	KnxManufacturer_AST_GMBH                                           KnxManufacturer = 220
	KnxManufacturer_WILA_LEUCHTEN                                      KnxManufacturer = 222
	KnxManufacturer_BPlusB_AUTOMATIONS__UND_STEUERUNGSTECHNIK          KnxManufacturer = 223
	KnxManufacturer_LINGG_AND_JANKE                                    KnxManufacturer = 225
	KnxManufacturer_SAUTER                                             KnxManufacturer = 227
	KnxManufacturer_SIMU                                               KnxManufacturer = 228
	KnxManufacturer_THEBEN_HTS_AG                                      KnxManufacturer = 232
	KnxManufacturer_AMANN_GMBH                                         KnxManufacturer = 233
	KnxManufacturer_BERG_ENERGIEKONTROLLSYSTEME_GMBH                   KnxManufacturer = 234
	KnxManufacturer_HUEPPE_FORM_SONNENSCHUTZSYSTEME_GMBH               KnxManufacturer = 235
	KnxManufacturer_OVENTROP_KG                                        KnxManufacturer = 237
	KnxManufacturer_GRIESSER_AG                                        KnxManufacturer = 238
	KnxManufacturer_IPAS_GMBH                                          KnxManufacturer = 239
	KnxManufacturer_ELERO_GMBH                                         KnxManufacturer = 240
	KnxManufacturer_ARDAN_PRODUCTION_AND_INDUSTRIAL_CONTROLS_LTD_      KnxManufacturer = 241
	KnxManufacturer_METEC_MESSTECHNIK_GMBH                             KnxManufacturer = 242
	KnxManufacturer_ELKA_ELEKTRONIK_GMBH                               KnxManufacturer = 244
	KnxManufacturer_ELEKTROANLAGEN_D__NAGEL                            KnxManufacturer = 245
	KnxManufacturer_TRIDONIC_BAUELEMENTE_GMBH                          KnxManufacturer = 246
	KnxManufacturer_STENGLER_GESELLSCHAFT                              KnxManufacturer = 248
	KnxManufacturer_SCHNEIDER_ELECTRIC_MG                              KnxManufacturer = 249
	KnxManufacturer_KNX_ASSOCIATION                                    KnxManufacturer = 250
	KnxManufacturer_VIVO                                               KnxManufacturer = 251
	KnxManufacturer_HUGO_MUELLER_GMBH_AND_CO_KG                        KnxManufacturer = 252
	KnxManufacturer_SIEMENS_HVAC                                       KnxManufacturer = 253
	KnxManufacturer_APT                                                KnxManufacturer = 254
	KnxManufacturer_HIGHDOM                                            KnxManufacturer = 256
	KnxManufacturer_TOP_SERVICES                                       KnxManufacturer = 257
	KnxManufacturer_AMBIHOME                                           KnxManufacturer = 258
	KnxManufacturer_DATEC_ELECTRONIC_AG                                KnxManufacturer = 259
	KnxManufacturer_ABUS_SECURITY_CENTER                               KnxManufacturer = 260
	KnxManufacturer_LITE_PUTER                                         KnxManufacturer = 261
	KnxManufacturer_TANTRON_ELECTRONIC                                 KnxManufacturer = 262
	KnxManufacturer_INTERRA                                            KnxManufacturer = 263
	KnxManufacturer_DKX_TECH                                           KnxManufacturer = 264
	KnxManufacturer_VIATRON                                            KnxManufacturer = 265
	KnxManufacturer_NAUTIBUS                                           KnxManufacturer = 266
	KnxManufacturer_ON_SEMICONDUCTOR                                   KnxManufacturer = 267
	KnxManufacturer_LONGCHUANG                                         KnxManufacturer = 268
	KnxManufacturer_AIR_ON_AG                                          KnxManufacturer = 269
	KnxManufacturer_IB_COMPANY_GMBH                                    KnxManufacturer = 270
	KnxManufacturer_SATION_FACTORY                                     KnxManufacturer = 271
	KnxManufacturer_AGENTILO_GMBH                                      KnxManufacturer = 272
	KnxManufacturer_MAKEL_ELEKTRIK                                     KnxManufacturer = 273
	KnxManufacturer_HELIOS_VENTILATOREN                                KnxManufacturer = 274
	KnxManufacturer_OTTO_SOLUTIONS_PTE_LTD                             KnxManufacturer = 275
	KnxManufacturer_AIRMASTER                                          KnxManufacturer = 276
	KnxManufacturer_VALLOX_GMBH                                        KnxManufacturer = 277
	KnxManufacturer_DALITEK                                            KnxManufacturer = 278
	KnxManufacturer_ASIN                                               KnxManufacturer = 279
	KnxManufacturer_BRIDGES_INTELLIGENCE_TECHNOLOGY_INC_               KnxManufacturer = 280
	KnxManufacturer_ARBONIA                                            KnxManufacturer = 281
	KnxManufacturer_KERMI                                              KnxManufacturer = 282
	KnxManufacturer_PROLUX                                             KnxManufacturer = 283
	KnxManufacturer_CLICHOME                                           KnxManufacturer = 284
	KnxManufacturer_COMMAX                                             KnxManufacturer = 285
	KnxManufacturer_EAE                                                KnxManufacturer = 286
	KnxManufacturer_TENSE                                              KnxManufacturer = 287
	KnxManufacturer_SEYOUNG_ELECTRONICS                                KnxManufacturer = 288
	KnxManufacturer_LIFEDOMUS                                          KnxManufacturer = 289
	KnxManufacturer_EUROTRONIC_TECHNOLOGY_GMBH                         KnxManufacturer = 290
	KnxManufacturer_TCI                                                KnxManufacturer = 291
	KnxManufacturer_RISHUN_ELECTRONIC                                  KnxManufacturer = 292
	KnxManufacturer_ZIPATO                                             KnxManufacturer = 293
	KnxManufacturer_CM_SECURITY_GMBH_AND_CO_KG                         KnxManufacturer = 294
	KnxManufacturer_QING_CABLES                                        KnxManufacturer = 295
	KnxManufacturer_LABIO                                              KnxManufacturer = 296
	KnxManufacturer_COSTER_TECNOLOGIE_ELETTRONICHE_S_P_A_              KnxManufacturer = 297
	KnxManufacturer_E_G_E                                              KnxManufacturer = 298
	KnxManufacturer_NETXAUTOMATION                                     KnxManufacturer = 299
	KnxManufacturer_TECALOR                                            KnxManufacturer = 300
	KnxManufacturer_URMET_ELECTRONICS_HUIZHOU_LTD_                     KnxManufacturer = 301
	KnxManufacturer_PEIYING_BUILDING_CONTROL                           KnxManufacturer = 302
	KnxManufacturer_BPT_S_P_A__A_SOCIO_UNICO                           KnxManufacturer = 303
	KnxManufacturer_KANONTEC___KANONBUS                                KnxManufacturer = 304
	KnxManufacturer_ISER_TECH                                          KnxManufacturer = 305
	KnxManufacturer_FINELINE                                           KnxManufacturer = 306
	KnxManufacturer_CP_ELECTRONICS_LTD                                 KnxManufacturer = 307
	KnxManufacturer_NIKO_SERVODAN_AS                                   KnxManufacturer = 308
	KnxManufacturer_SIMON_309                                          KnxManufacturer = 309
	KnxManufacturer_GM_MODULAR_PVT__LTD_                               KnxManufacturer = 310
	KnxManufacturer_FU_CHENG_INTELLIGENCE                              KnxManufacturer = 311
	KnxManufacturer_NEXKON                                             KnxManufacturer = 312
	KnxManufacturer_FEEL_S_R_L                                         KnxManufacturer = 313
	KnxManufacturer_NOT_ASSIGNED_314                                   KnxManufacturer = 314
	KnxManufacturer_SHENZHEN_FANHAI_SANJIANG_ELECTRONICS_CO___LTD_     KnxManufacturer = 315
	KnxManufacturer_JIUZHOU_GREEBLE                                    KnxManufacturer = 316
	KnxManufacturer_AUMUELLER_AUMATIC_GMBH                             KnxManufacturer = 317
	KnxManufacturer_ETMAN_ELECTRIC                                     KnxManufacturer = 318
	KnxManufacturer_BLACK_NOVA                                         KnxManufacturer = 319
	KnxManufacturer_ZIDATECH_AG                                        KnxManufacturer = 320
	KnxManufacturer_IDGS_BVBA                                          KnxManufacturer = 321
	KnxManufacturer_DAKANIMO                                           KnxManufacturer = 322
	KnxManufacturer_TREBOR_AUTOMATION_AB                               KnxManufacturer = 323
	KnxManufacturer_SATEL_SP__Z_O_O_                                   KnxManufacturer = 324
	KnxManufacturer_RUSSOUND__INC_                                     KnxManufacturer = 325
	KnxManufacturer_MIDEA_HEATING_AND_VENTILATING_EQUIPMENT_CO_LTD     KnxManufacturer = 326
	KnxManufacturer_CONSORZIO_TERRANUOVA                               KnxManufacturer = 327
	KnxManufacturer_WOLF_HEIZTECHNIK_GMBH                              KnxManufacturer = 328
	KnxManufacturer_SONTEC                                             KnxManufacturer = 329
	KnxManufacturer_BELCOM_CABLES_LTD_                                 KnxManufacturer = 330
	KnxManufacturer_GUANGZHOU_SEAWIN_ELECTRICAL_TECHNOLOGIES_CO___LTD_ KnxManufacturer = 331
	KnxManufacturer_ACREL                                              KnxManufacturer = 332
	KnxManufacturer_FRANKE_AQUAROTTER_GMBH                             KnxManufacturer = 333
	KnxManufacturer_ORION_SYSTEMS                                      KnxManufacturer = 334
	KnxManufacturer_SCHRACK_TECHNIK_GMBH                               KnxManufacturer = 335
	KnxManufacturer_INSPRID                                            KnxManufacturer = 336
	KnxManufacturer_SUNRICHER                                          KnxManufacturer = 337
	KnxManufacturer_MENRED_AUTOMATION_SYSTEMSHANGHAI_CO__LTD_          KnxManufacturer = 338
	KnxManufacturer_AUREX                                              KnxManufacturer = 339
	KnxManufacturer_JOSEF_BARTHELME_GMBH_AND_CO__KG                    KnxManufacturer = 340
	KnxManufacturer_ARCHITECTURE_NUMERIQUE                             KnxManufacturer = 341
	KnxManufacturer_UP_GROUP                                           KnxManufacturer = 342
	KnxManufacturer_TEKNOS_AVINNO                                      KnxManufacturer = 343
	KnxManufacturer_NINGBO_DOOYA_MECHANIC_AND_ELECTRONIC_TECHNOLOGY    KnxManufacturer = 344
	KnxManufacturer_THERMOKON_SENSORTECHNIK_GMBH                       KnxManufacturer = 345
	KnxManufacturer_BELIMO_AUTOMATION_AG                               KnxManufacturer = 346
	KnxManufacturer_ZEHNDER_GROUP_INTERNATIONAL_AG                     KnxManufacturer = 347
	KnxManufacturer_SKS_KINKEL_ELEKTRONIK                              KnxManufacturer = 348
	KnxManufacturer_ECE_WURMITZER_GMBH                                 KnxManufacturer = 349
	KnxManufacturer_LARS                                               KnxManufacturer = 350
	KnxManufacturer_URC                                                KnxManufacturer = 351
	KnxManufacturer_LIGHTCONTROL                                       KnxManufacturer = 352
	KnxManufacturer_SHENZHEN_YM                                        KnxManufacturer = 353
	KnxManufacturer_MEAN_WELL_ENTERPRISES_CO__LTD_                     KnxManufacturer = 354
	KnxManufacturer_OSIX                                               KnxManufacturer = 355
	KnxManufacturer_AYPRO_TECHNOLOGY                                   KnxManufacturer = 356
	KnxManufacturer_HEFEI_ECOLITE_SOFTWARE                             KnxManufacturer = 357
	KnxManufacturer_ENNO                                               KnxManufacturer = 358
	KnxManufacturer_OHOSURE                                            KnxManufacturer = 359
	KnxManufacturer_GAREFOWL                                           KnxManufacturer = 360
	KnxManufacturer_GEZE                                               KnxManufacturer = 361
	KnxManufacturer_LG_ELECTRONICS_INC_                                KnxManufacturer = 362
	KnxManufacturer_SMC_INTERIORS                                      KnxManufacturer = 363
	KnxManufacturer_NOT_ASSIGNED_364                                   KnxManufacturer = 364
	KnxManufacturer_SCS_CABLE                                          KnxManufacturer = 365
	KnxManufacturer_HOVAL                                              KnxManufacturer = 366
	KnxManufacturer_CANST                                              KnxManufacturer = 367
	KnxManufacturer_HANGZHOU_BERLIN                                    KnxManufacturer = 368
	KnxManufacturer_EVN_LICHTTECHNIK                                   KnxManufacturer = 369
	KnxManufacturer_RUTEC                                              KnxManufacturer = 370
	KnxManufacturer_FINDER                                             KnxManufacturer = 371
	KnxManufacturer_FUJITSU_GENERAL_LIMITED                            KnxManufacturer = 372
	KnxManufacturer_ZF_FRIEDRICHSHAFEN_AG                              KnxManufacturer = 373
	KnxManufacturer_CREALED                                            KnxManufacturer = 374
	KnxManufacturer_MILES_MAGIC_AUTOMATION_PRIVATE_LIMITED             KnxManufacturer = 375
	KnxManufacturer_EPlus                                              KnxManufacturer = 376
	KnxManufacturer_ITALCOND                                           KnxManufacturer = 377
	KnxManufacturer_SATION                                             KnxManufacturer = 378
	KnxManufacturer_NEWBEST                                            KnxManufacturer = 379
	KnxManufacturer_GDS_DIGITAL_SYSTEMS                                KnxManufacturer = 380
	KnxManufacturer_IDDERO                                             KnxManufacturer = 381
	KnxManufacturer_MBNLED                                             KnxManufacturer = 382
	KnxManufacturer_VITRUM                                             KnxManufacturer = 383
	KnxManufacturer_EKEY_BIOMETRIC_SYSTEMS_GMBH                        KnxManufacturer = 384
	KnxManufacturer_AMC                                                KnxManufacturer = 385
	KnxManufacturer_TRILUX_GMBH_AND_CO__KG                             KnxManufacturer = 386
	KnxManufacturer_WEXCEDO                                            KnxManufacturer = 387
	KnxManufacturer_VEMER_SPA                                          KnxManufacturer = 388
	KnxManufacturer_ALEXANDER_BUERKLE_GMBH_AND_CO_KG                   KnxManufacturer = 389
	KnxManufacturer_CITRON                                             KnxManufacturer = 390
	KnxManufacturer_SHENZHEN_HEGUANG                                   KnxManufacturer = 391
	KnxManufacturer_NOT_ASSIGNED_392                                   KnxManufacturer = 392
	KnxManufacturer_TRANE_B_V_B_A                                      KnxManufacturer = 393
	KnxManufacturer_CAREL                                              KnxManufacturer = 394
	KnxManufacturer_PROLITE_CONTROLS                                   KnxManufacturer = 395
	KnxManufacturer_BOSMER                                             KnxManufacturer = 396
	KnxManufacturer_EUCHIPS                                            KnxManufacturer = 397
	KnxManufacturer_CONNECT_THINKA_CONNECT                             KnxManufacturer = 398
	KnxManufacturer_PEAKNX_A_DOGAWIST_COMPANY                          KnxManufacturer = 399
	KnxManufacturer_ACEMATIC                                           KnxManufacturer = 400
	KnxManufacturer_ELAUSYS                                            KnxManufacturer = 401
	KnxManufacturer_ITK_ENGINEERING_AG                                 KnxManufacturer = 402
	KnxManufacturer_INTEGRA_METERING_AG                                KnxManufacturer = 403
	KnxManufacturer_FMS_HOSPITALITY_PTE_LTD                            KnxManufacturer = 404
	KnxManufacturer_NUVO                                               KnxManufacturer = 405
	KnxManufacturer_U__LUX_GMBH                                        KnxManufacturer = 406
	KnxManufacturer_BRUMBERG_LEUCHTEN                                  KnxManufacturer = 407
	KnxManufacturer_LIME                                               KnxManufacturer = 408
	KnxManufacturer_GREAT_EMPIRE_INTERNATIONAL_GROUP_CO___LTD_         KnxManufacturer = 409
	KnxManufacturer_KAVOSHPISHRO_ASIA                                  KnxManufacturer = 410
	KnxManufacturer_V2_SPA                                             KnxManufacturer = 411
	KnxManufacturer_JOHNSON_CONTROLS                                   KnxManufacturer = 412
	KnxManufacturer_ARKUD                                              KnxManufacturer = 413
	KnxManufacturer_IRIDIUM_LTD_                                       KnxManufacturer = 414
	KnxManufacturer_BSMART                                             KnxManufacturer = 415
	KnxManufacturer_BAB_TECHNOLOGIE_GMBH                               KnxManufacturer = 416
	KnxManufacturer_NICE_SPA                                           KnxManufacturer = 417
	KnxManufacturer_REDFISH_GROUP_PTY_LTD                              KnxManufacturer = 418
	KnxManufacturer_SABIANA_SPA                                        KnxManufacturer = 419
	KnxManufacturer_UBEE_INTERACTIVE_EUROPE                            KnxManufacturer = 420
	KnxManufacturer_REXEL                                              KnxManufacturer = 421
	KnxManufacturer_GES_TEKNIK_A_S_                                    KnxManufacturer = 422
	KnxManufacturer_AVE_S_P_A_                                         KnxManufacturer = 423
	KnxManufacturer_ZHUHAI_LTECH_TECHNOLOGY_CO___LTD_                  KnxManufacturer = 424
	KnxManufacturer_ARCOM                                              KnxManufacturer = 425
	KnxManufacturer_VIA_TECHNOLOGIES__INC_                             KnxManufacturer = 426
	KnxManufacturer_FEELSMART_                                         KnxManufacturer = 427
	KnxManufacturer_SUPCON                                             KnxManufacturer = 428
	KnxManufacturer_MANIC                                              KnxManufacturer = 429
	KnxManufacturer_TDE_GMBH                                           KnxManufacturer = 430
	KnxManufacturer_NANJING_SHUFAN_INFORMATION_TECHNOLOGY_CO__LTD_     KnxManufacturer = 431
	KnxManufacturer_EWTECH                                             KnxManufacturer = 432
	KnxManufacturer_KLUGER_AUTOMATION_GMBH                             KnxManufacturer = 433
	KnxManufacturer_JOONGANG_CONTROL                                   KnxManufacturer = 434
	KnxManufacturer_GREENCONTROLS_TECHNOLOGY_SDN__BHD_                 KnxManufacturer = 435
	KnxManufacturer_IME_S_P_A_                                         KnxManufacturer = 436
	KnxManufacturer_SICHUAN_HAODING                                    KnxManufacturer = 437
	KnxManufacturer_MINDJAGA_LTD_                                      KnxManufacturer = 438
	KnxManufacturer_RUILI_SMART_CONTROL                                KnxManufacturer = 439
	KnxManufacturer_CODESYS_GMBH                                       KnxManufacturer = 440
	KnxManufacturer_MOORGEN_DEUTSCHLAND_GMBH                           KnxManufacturer = 441
	KnxManufacturer_CULLMANN_TECH                                      KnxManufacturer = 442
	KnxManufacturer_MERCK_WINDOW_TECHNOLOGIES_B_V_                     KnxManufacturer = 443
	KnxManufacturer_ABEGO                                              KnxManufacturer = 444
	KnxManufacturer_MYGEKKO                                            KnxManufacturer = 445
	KnxManufacturer_ERGO3_SARL                                         KnxManufacturer = 446
	KnxManufacturer_STMICROELECTRONICS_INTERNATIONAL_N_V_              KnxManufacturer = 447
	KnxManufacturer_CJC_SYSTEMS                                        KnxManufacturer = 448
	KnxManufacturer_SUDOKU                                             KnxManufacturer = 449
	KnxManufacturer_AZ_E_LITE_PTE_LTD                                  KnxManufacturer = 451
	KnxManufacturer_ARLIGHT                                            KnxManufacturer = 452
	KnxManufacturer_GRUENBECK_WASSERAUFBEREITUNG_GMBH                  KnxManufacturer = 453
	KnxManufacturer_MODULE_ELECTRONIC                                  KnxManufacturer = 454
	KnxManufacturer_KOPLAT                                             KnxManufacturer = 455
	KnxManufacturer_GUANGZHOU_LETOUR_LIFE_TECHNOLOGY_CO___LTD          KnxManufacturer = 456
	KnxManufacturer_ILEVIA                                             KnxManufacturer = 457
	KnxManufacturer_LN_SYSTEMTEQ                                       KnxManufacturer = 458
	KnxManufacturer_HISENSE_SMARTHOME                                  KnxManufacturer = 459
	KnxManufacturer_FLINK_AUTOMATION_SYSTEM                            KnxManufacturer = 460
	KnxManufacturer_XXTER_BV                                           KnxManufacturer = 461
	KnxManufacturer_LYNXUS_TECHNOLOGY                                  KnxManufacturer = 462
	KnxManufacturer_ROBOT_S_A_                                         KnxManufacturer = 463
	KnxManufacturer_SHENZHEN_ATTE_SMART_LIFE_CO__LTD_                  KnxManufacturer = 464
	KnxManufacturer_NOBLESSE                                           KnxManufacturer = 465
	KnxManufacturer_ADVANCED_DEVICES                                   KnxManufacturer = 466
	KnxManufacturer_ATRINA_BUILDING_AUTOMATION_CO__LTD                 KnxManufacturer = 467
	KnxManufacturer_GUANGDONG_DAMING_LAFFEY_ELECTRIC_CO___LTD_         KnxManufacturer = 468
	KnxManufacturer_WESTERSTRAND_URFABRIK_AB                           KnxManufacturer = 469
	KnxManufacturer_CONTROL4_CORPORATE                                 KnxManufacturer = 470
	KnxManufacturer_ONTROL                                             KnxManufacturer = 471
	KnxManufacturer_STARNET                                            KnxManufacturer = 472
	KnxManufacturer_BETA_CAVI                                          KnxManufacturer = 473
	KnxManufacturer_EASEMORE                                           KnxManufacturer = 474
	KnxManufacturer_VIVALDI_SRL                                        KnxManufacturer = 475
	KnxManufacturer_GREE_ELECTRIC_APPLIANCES_INC__OF_ZHUHAI            KnxManufacturer = 476
	KnxManufacturer_HWISCON                                            KnxManufacturer = 477
	KnxManufacturer_SHANGHAI_ELECON_INTELLIGENT_TECHNOLOGY_CO___LTD_   KnxManufacturer = 478
	KnxManufacturer_KAMPMANN                                           KnxManufacturer = 479
	KnxManufacturer_IMPOLUX_GMBH_LEDIMAX                               KnxManufacturer = 480
	KnxManufacturer_EVAUX                                              KnxManufacturer = 481
	KnxManufacturer_WEBRO_CABLES_AND_CONNECTORS_LIMITED                KnxManufacturer = 482
	KnxManufacturer_SHANGHAI_E_TECH_SOLUTION                           KnxManufacturer = 483
	KnxManufacturer_GUANGZHOU_HOKO_ELECTRIC_CO__LTD_                   KnxManufacturer = 484
	KnxManufacturer_LAMMIN_HIGH_TECH_CO__LTD                           KnxManufacturer = 485
	KnxManufacturer_SHENZHEN_MERRYTEK_TECHNOLOGY_CO___LTD              KnxManufacturer = 486
	KnxManufacturer_I_LUXUS                                            KnxManufacturer = 487
	KnxManufacturer_ELMOS_SEMICONDUCTOR_AG                             KnxManufacturer = 488
	KnxManufacturer_EMCOM_TECHNOLOGY_INC                               KnxManufacturer = 489
	KnxManufacturer_PROJECT_INNOVATIONS_GMBH                           KnxManufacturer = 490
	KnxManufacturer_ITC                                                KnxManufacturer = 491
	KnxManufacturer_ABB_LV_INSTALLATION_MATERIALS_COMPANY_LTD__BEIJING KnxManufacturer = 492
	KnxManufacturer_MAICO                                              KnxManufacturer = 493
	KnxManufacturer_ELAN_SRL                                           KnxManufacturer = 495
	KnxManufacturer_MINHHA_TECHNOLOGY_CO__LTD                          KnxManufacturer = 496
	KnxManufacturer_ZHEJIANG_TIANJIE_INDUSTRIAL_CORP_                  KnxManufacturer = 497
	KnxManufacturer_IAUTOMATION_PTY_LIMITED                            KnxManufacturer = 498
	KnxManufacturer_EXTRON                                             KnxManufacturer = 499
	KnxManufacturer_FREEDOMPRO                                         KnxManufacturer = 500
	KnxManufacturer_1HOME                                              KnxManufacturer = 501
	KnxManufacturer_EOS_SAUNATECHNIK_GMBH                              KnxManufacturer = 502
	KnxManufacturer_KUSATEK_GMBH                                       KnxManufacturer = 503
	KnxManufacturer_EISBAER_SCADA                                      KnxManufacturer = 504
	KnxManufacturer_AUTOMATISMI_BENINCA_S_P_A_                         KnxManufacturer = 505
	KnxManufacturer_BLENDOM                                            KnxManufacturer = 506
	KnxManufacturer_MADEL_AIR_TECHNICAL_DIFFUSION                      KnxManufacturer = 507
	KnxManufacturer_NIKO                                               KnxManufacturer = 508
	KnxManufacturer_BOSCH_REXROTH_AG                                   KnxManufacturer = 509
	KnxManufacturer_CANDM_PRODUCTS                                     KnxManufacturer = 512
	KnxManufacturer_HOERMANN_KG_VERKAUFSGESELLSCHAFT                   KnxManufacturer = 513
	KnxManufacturer_SHANGHAI_RAJAYASA_CO__LTD                          KnxManufacturer = 514
	KnxManufacturer_SUZUKI                                             KnxManufacturer = 515
	KnxManufacturer_SILENT_GLISS_INTERNATIONAL_LTD_                    KnxManufacturer = 516
	KnxManufacturer_BEE_CONTROLS_ADGSC_GROUP                           KnxManufacturer = 517
	KnxManufacturer_XDTECGMBH                                          KnxManufacturer = 518
	KnxManufacturer_OSRAM                                              KnxManufacturer = 519
	KnxManufacturer_LEBENOR                                            KnxManufacturer = 520
	KnxManufacturer_AUTOMANENG                                         KnxManufacturer = 521
	KnxManufacturer_HONEYWELL_AUTOMATION_SOLUTION_CONTROLCHINA         KnxManufacturer = 522
	KnxManufacturer_HANGZHOU_BINTHEN_INTELLIGENCE_TECHNOLOGY_CO__LTD   KnxManufacturer = 523
	KnxManufacturer_ETA_HEIZTECHNIK                                    KnxManufacturer = 524
	KnxManufacturer_DIVUS_GMBH                                         KnxManufacturer = 525
	KnxManufacturer_NANJING_TAIJIESAI_INTELLIGENT_TECHNOLOGY_CO__LTD_  KnxManufacturer = 526
	KnxManufacturer_LUNATONE                                           KnxManufacturer = 527
	KnxManufacturer_ZHEJIANG_SCTECH_BUILDING_INTELLIGENT               KnxManufacturer = 528
	KnxManufacturer_FOSHAN_QITE_TECHNOLOGY_CO___LTD_                   KnxManufacturer = 529
	KnxManufacturer_NOKE                                               KnxManufacturer = 530
	KnxManufacturer_LANDCOM                                            KnxManufacturer = 531
	KnxManufacturer_STORK_AS                                           KnxManufacturer = 532
	KnxManufacturer_HANGZHOU_SHENDU_TECHNOLOGY_CO___LTD_               KnxManufacturer = 533
	KnxManufacturer_COOLAUTOMATION                                     KnxManufacturer = 534
	KnxManufacturer_APRSTERN                                           KnxManufacturer = 535
	KnxManufacturer_SONNEN                                             KnxManufacturer = 536
	KnxManufacturer_DNAKE                                              KnxManufacturer = 537
	KnxManufacturer_NEUBERGER_GEBAEUDEAUTOMATION_GMBH                  KnxManufacturer = 538
	KnxManufacturer_STILIGER                                           KnxManufacturer = 539
	KnxManufacturer_BERGHOF_AUTOMATION_GMBH                            KnxManufacturer = 540
	KnxManufacturer_TOTAL_AUTOMATION_AND_CONTROLS_GMBH                 KnxManufacturer = 541
	KnxManufacturer_DOVIT                                              KnxManufacturer = 542
	KnxManufacturer_INSTALIGHTING_GMBH                                 KnxManufacturer = 543
	KnxManufacturer_UNI_TEC                                            KnxManufacturer = 544
	KnxManufacturer_CASATUNES                                          KnxManufacturer = 545
	KnxManufacturer_EMT                                                KnxManufacturer = 546
	KnxManufacturer_SENFFICIENT                                        KnxManufacturer = 547
	KnxManufacturer_AUROLITE_ELECTRICAL_PANYU_GUANGZHOU_LIMITED        KnxManufacturer = 548
	KnxManufacturer_ABB_XIAMEN_SMART_TECHNOLOGY_CO___LTD_              KnxManufacturer = 549
	KnxManufacturer_SAMSON_ELECTRIC_WIRE                               KnxManufacturer = 550
	KnxManufacturer_T_TOUCHING                                         KnxManufacturer = 551
	KnxManufacturer_CORE_SMART_HOME                                    KnxManufacturer = 552
	KnxManufacturer_GREENCONNECT_SOLUTIONS_SA                          KnxManufacturer = 553
	KnxManufacturer_ELETTRONICA_CONDUTTORI                             KnxManufacturer = 554
	KnxManufacturer_MKFC                                               KnxManufacturer = 555
	KnxManufacturer_AUTOMATIONPlus                                     KnxManufacturer = 556
	KnxManufacturer_BLUE_AND_RED                                       KnxManufacturer = 557
	KnxManufacturer_FROGBLUE                                           KnxManufacturer = 558
	KnxManufacturer_SAVESOR                                            KnxManufacturer = 559
	KnxManufacturer_APP_TECH                                           KnxManufacturer = 560
	KnxManufacturer_SENSORTEC_AG                                       KnxManufacturer = 561
	KnxManufacturer_NYSA_TECHNOLOGY_AND_SOLUTIONS                      KnxManufacturer = 562
	KnxManufacturer_FARADITE                                           KnxManufacturer = 563
	KnxManufacturer_OPTIMUS                                            KnxManufacturer = 564
	KnxManufacturer_KTS_S_R_L_                                         KnxManufacturer = 565
	KnxManufacturer_RAMCRO_SPA                                         KnxManufacturer = 566
	KnxManufacturer_WUHAN_WISECREATE_UNIVERSE_TECHNOLOGY_CO___LTD      KnxManufacturer = 567
	KnxManufacturer_BEMI_SMART_HOME_LTD                                KnxManufacturer = 568
	KnxManufacturer_ARDOMUS                                            KnxManufacturer = 569
	KnxManufacturer_CHANGXING                                          KnxManufacturer = 570
	KnxManufacturer_E_CONTROLS                                         KnxManufacturer = 571
	KnxManufacturer_AIB_TECHNOLOGY                                     KnxManufacturer = 572
	KnxManufacturer_NVC                                                KnxManufacturer = 573
	KnxManufacturer_KBOX                                               KnxManufacturer = 574
	KnxManufacturer_CNS                                                KnxManufacturer = 575
	KnxManufacturer_TYBA                                               KnxManufacturer = 576
	KnxManufacturer_ATREL                                              KnxManufacturer = 577
	KnxManufacturer_SIMON_ELECTRIC_CHINA_CO___LTD                      KnxManufacturer = 578
	KnxManufacturer_KORDZ_GROUP                                        KnxManufacturer = 579
	KnxManufacturer_ND_ELECTRIC                                        KnxManufacturer = 580
	KnxManufacturer_CONTROLIUM                                         KnxManufacturer = 581
	KnxManufacturer_FAMO_GMBH_AND_CO__KG                               KnxManufacturer = 582
	KnxManufacturer_CDN_SMART                                          KnxManufacturer = 583
	KnxManufacturer_HESTON                                             KnxManufacturer = 584
	KnxManufacturer_ESLA_CONEXIONES_S_L_                               KnxManufacturer = 585
	KnxManufacturer_WEISHAUPT                                          KnxManufacturer = 586
	KnxManufacturer_ASTRUM_TECHNOLOGY                                  KnxManufacturer = 587
	KnxManufacturer_WUERTH_ELEKTRONIK_STELVIO_KONTEK_S_P_A_            KnxManufacturer = 588
	KnxManufacturer_NANOTECO_CORPORATION                               KnxManufacturer = 589
	KnxManufacturer_NIETIAN                                            KnxManufacturer = 590
	KnxManufacturer_SUMSIR                                             KnxManufacturer = 591
	KnxManufacturer_ORBIS_TECNOLOGIA_ELECTRICA_SA                      KnxManufacturer = 592
	KnxManufacturer_ABB___RESERVED                                     KnxManufacturer = 43954
	KnxManufacturer_BUSCH_JAEGER_ELEKTRO___RESERVED                    KnxManufacturer = 43959
)

func (e KnxManufacturer) Text() string {
	switch e {
	case 1:
		{ /* '1' */
			return "Siemens"
		}
	case 10:
		{ /* '10' */
			return "Insta GmbH"
		}
	case 100:
		{ /* '100' */
			return "Schneider Electric Industries SAS"
		}
	case 101:
		{ /* '101' */
			return "WHD Wilhelm Huber + Söhne"
		}
	case 102:
		{ /* '102' */
			return "Bischoff Elektronik"
		}
	case 104:
		{ /* '104' */
			return "JEPAZ"
		}
	case 105:
		{ /* '105' */
			return "RTS Automation"
		}
	case 106:
		{ /* '106' */
			return "EIBMARKT GmbH"
		}
	case 107:
		{ /* '107' */
			return "WAREMA Renkhoff SE"
		}
	case 108:
		{ /* '108' */
			return "Eelectron"
		}
	case 109:
		{ /* '109' */
			return "Belden Wire & Cable B.V."
		}
	case 11:
		{ /* '11' */
			return "LEGRAND Appareillage électrique"
		}
	case 110:
		{ /* '110' */
			return "Becker-Antriebe GmbH"
		}
	case 111:
		{ /* '111' */
			return "J.Stehle+Söhne GmbH"
		}
	case 112:
		{ /* '112' */
			return "AGFEO"
		}
	case 113:
		{ /* '113' */
			return "Zennio"
		}
	case 114:
		{ /* '114' */
			return "TAPKO Technologies"
		}
	case 115:
		{ /* '115' */
			return "HDL"
		}
	case 116:
		{ /* '116' */
			return "Uponor"
		}
	case 117:
		{ /* '117' */
			return "se Lightmanagement AG"
		}
	case 118:
		{ /* '118' */
			return "Arcus-eds"
		}
	case 119:
		{ /* '119' */
			return "Intesis"
		}
	case 12:
		{ /* '12' */
			return "Merten"
		}
	case 120:
		{ /* '120' */
			return "Herholdt Controls srl"
		}
	case 121:
		{ /* '121' */
			return "Niko-Zublin"
		}
	case 122:
		{ /* '122' */
			return "Durable Technologies"
		}
	case 123:
		{ /* '123' */
			return "Innoteam"
		}
	case 124:
		{ /* '124' */
			return "ise GmbH"
		}
	case 125:
		{ /* '125' */
			return "TEAM FOR TRONICS"
		}
	case 126:
		{ /* '126' */
			return "CIAT"
		}
	case 127:
		{ /* '127' */
			return "Remeha BV"
		}
	case 128:
		{ /* '128' */
			return "ESYLUX"
		}
	case 129:
		{ /* '129' */
			return "BASALTE"
		}
	case 130:
		{ /* '130' */
			return "Vestamatic"
		}
	case 131:
		{ /* '131' */
			return "MDT technologies"
		}
	case 132:
		{ /* '132' */
			return "Warendorfer Küchen GmbH"
		}
	case 133:
		{ /* '133' */
			return "Video-Star"
		}
	case 134:
		{ /* '134' */
			return "Sitek"
		}
	case 135:
		{ /* '135' */
			return "CONTROLtronic"
		}
	case 136:
		{ /* '136' */
			return "function Technology"
		}
	case 137:
		{ /* '137' */
			return "AMX"
		}
	case 138:
		{ /* '138' */
			return "ELDAT"
		}
	case 139:
		{ /* '139' */
			return "Panasonic"
		}
	case 14:
		{ /* '14' */
			return "ABB SpA-SACE Division"
		}
	case 140:
		{ /* '140' */
			return "Pulse Technologies"
		}
	case 141:
		{ /* '141' */
			return "Crestron"
		}
	case 142:
		{ /* '142' */
			return "STEINEL professional"
		}
	case 143:
		{ /* '143' */
			return "BILTON LED Lighting"
		}
	case 144:
		{ /* '144' */
			return "denro AG"
		}
	case 145:
		{ /* '145' */
			return "GePro"
		}
	case 146:
		{ /* '146' */
			return "preussen automation"
		}
	case 147:
		{ /* '147' */
			return "Zoppas Industries"
		}
	case 148:
		{ /* '148' */
			return "MACTECH"
		}
	case 149:
		{ /* '149' */
			return "TECHNO-TREND"
		}
	case 150:
		{ /* '150' */
			return "FS Cables"
		}
	case 151:
		{ /* '151' */
			return "Delta Dore"
		}
	case 152:
		{ /* '152' */
			return "Eissound"
		}
	case 153:
		{ /* '153' */
			return "Cisco"
		}
	case 154:
		{ /* '154' */
			return "Dinuy"
		}
	case 155:
		{ /* '155' */
			return "iKNiX"
		}
	case 156:
		{ /* '156' */
			return "Rademacher Geräte-Elektronik GmbH"
		}
	case 157:
		{ /* '157' */
			return "EGi Electroacustica General Iberica"
		}
	case 158:
		{ /* '158' */
			return "Bes – Ingenium"
		}
	case 159:
		{ /* '159' */
			return "ElabNET"
		}
	case 160:
		{ /* '160' */
			return "Blumotix"
		}
	case 161:
		{ /* '161' */
			return "Hunter Douglas"
		}
	case 162:
		{ /* '162' */
			return "APRICUM"
		}
	case 163:
		{ /* '163' */
			return "TIANSU Automation"
		}
	case 164:
		{ /* '164' */
			return "Bubendorff"
		}
	case 165:
		{ /* '165' */
			return "MBS GmbH"
		}
	case 166:
		{ /* '166' */
			return "Enertex Bayern GmbH"
		}
	case 167:
		{ /* '167' */
			return "BMS"
		}
	case 168:
		{ /* '168' */
			return "Sinapsi"
		}
	case 169:
		{ /* '169' */
			return "Embedded Systems SIA"
		}
	case 170:
		{ /* '170' */
			return "KNX1"
		}
	case 171:
		{ /* '171' */
			return "Tokka"
		}
	case 172:
		{ /* '172' */
			return "NanoSense"
		}
	case 173:
		{ /* '173' */
			return "PEAR Automation GmbH"
		}
	case 174:
		{ /* '174' */
			return "DGA"
		}
	case 175:
		{ /* '175' */
			return "Lutron"
		}
	case 176:
		{ /* '176' */
			return "AIRZONE – ALTRA"
		}
	case 177:
		{ /* '177' */
			return "Lithoss Design Switches"
		}
	case 178:
		{ /* '178' */
			return "3ATEL"
		}
	case 179:
		{ /* '179' */
			return "Philips Controls"
		}
	case 180:
		{ /* '180' */
			return "VELUX A/S"
		}
	case 181:
		{ /* '181' */
			return "LOYTEC"
		}
	case 182:
		{ /* '182' */
			return "Ekinex S.p.A."
		}
	case 183:
		{ /* '183' */
			return "SIRLAN Technologies"
		}
	case 184:
		{ /* '184' */
			return "ProKNX SAS"
		}
	case 185:
		{ /* '185' */
			return "IT GmbH"
		}
	case 186:
		{ /* '186' */
			return "RENSON"
		}
	case 187:
		{ /* '187' */
			return "HEP Group"
		}
	case 188:
		{ /* '188' */
			return "Balmart"
		}
	case 189:
		{ /* '189' */
			return "GFS GmbH"
		}
	case 190:
		{ /* '190' */
			return "Schenker Storen AG"
		}
	case 191:
		{ /* '191' */
			return "Algodue Elettronica S.r.L."
		}
	case 192:
		{ /* '192' */
			return "ABB France"
		}
	case 193:
		{ /* '193' */
			return "maintronic"
		}
	case 194:
		{ /* '194' */
			return "Vantage"
		}
	case 195:
		{ /* '195' */
			return "Foresis"
		}
	case 196:
		{ /* '196' */
			return "Research & Production Association SEM"
		}
	case 197:
		{ /* '197' */
			return "Weinzierl Engineering GmbH"
		}
	case 198:
		{ /* '198' */
			return "Möhlenhoff Wärmetechnik GmbH"
		}
	case 199:
		{ /* '199' */
			return "PKC-GROUP Oyj"
		}
	case 2:
		{ /* '2' */
			return "ABB"
		}
	case 200:
		{ /* '200' */
			return "B.E.G."
		}
	case 201:
		{ /* '201' */
			return "Elsner Elektronik GmbH"
		}
	case 202:
		{ /* '202' */
			return "Siemens Building Technologies (HK/China) Ltd."
		}
	case 204:
		{ /* '204' */
			return "Eutrac"
		}
	case 205:
		{ /* '205' */
			return "Gustav Hensel GmbH & Co. KG"
		}
	case 206:
		{ /* '206' */
			return "GARO AB"
		}
	case 207:
		{ /* '207' */
			return "Waldmann Lichttechnik"
		}
	case 208:
		{ /* '208' */
			return "SCHÜCO"
		}
	case 209:
		{ /* '209' */
			return "EMU"
		}
	case 210:
		{ /* '210' */
			return "JNet Systems AG"
		}
	case 211:
		{ /* '211' */
			return "Total Solution GmbH"
		}
	case 214:
		{ /* '214' */
			return "O.Y.L. Electronics"
		}
	case 215:
		{ /* '215' */
			return "Galax System"
		}
	case 216:
		{ /* '216' */
			return "Disch"
		}
	case 217:
		{ /* '217' */
			return "Aucoteam"
		}
	case 218:
		{ /* '218' */
			return "Luxmate Controls"
		}
	case 219:
		{ /* '219' */
			return "Danfoss"
		}
	case 22:
		{ /* '22' */
			return "Siedle & Söhne"
		}
	case 220:
		{ /* '220' */
			return "AST GmbH"
		}
	case 222:
		{ /* '222' */
			return "WILA Leuchten"
		}
	case 223:
		{ /* '223' */
			return "b+b Automations- und Steuerungstechnik"
		}
	case 225:
		{ /* '225' */
			return "Lingg & Janke"
		}
	case 227:
		{ /* '227' */
			return "Sauter"
		}
	case 228:
		{ /* '228' */
			return "SIMU"
		}
	case 232:
		{ /* '232' */
			return "Theben HTS AG"
		}
	case 233:
		{ /* '233' */
			return "Amann GmbH"
		}
	case 234:
		{ /* '234' */
			return "BERG Energiekontrollsysteme GmbH"
		}
	case 235:
		{ /* '235' */
			return "Hüppe Form Sonnenschutzsysteme GmbH"
		}
	case 237:
		{ /* '237' */
			return "Oventrop KG"
		}
	case 238:
		{ /* '238' */
			return "Griesser AG"
		}
	case 239:
		{ /* '239' */
			return "IPAS GmbH"
		}
	case 24:
		{ /* '24' */
			return "Eberle"
		}
	case 240:
		{ /* '240' */
			return "elero GmbH"
		}
	case 241:
		{ /* '241' */
			return "Ardan Production and Industrial Controls Ltd."
		}
	case 242:
		{ /* '242' */
			return "Metec Meßtechnik GmbH"
		}
	case 244:
		{ /* '244' */
			return "ELKA-Elektronik GmbH"
		}
	case 245:
		{ /* '245' */
			return "ELEKTROANLAGEN D. NAGEL"
		}
	case 246:
		{ /* '246' */
			return "Tridonic Bauelemente GmbH"
		}
	case 248:
		{ /* '248' */
			return "Stengler Gesellschaft"
		}
	case 249:
		{ /* '249' */
			return "Schneider Electric (MG)"
		}
	case 25:
		{ /* '25' */
			return "GEWISS"
		}
	case 250:
		{ /* '250' */
			return "KNX Association"
		}
	case 251:
		{ /* '251' */
			return "VIVO"
		}
	case 252:
		{ /* '252' */
			return "Hugo Müller GmbH & Co KG"
		}
	case 253:
		{ /* '253' */
			return "Siemens HVAC"
		}
	case 254:
		{ /* '254' */
			return "APT"
		}
	case 256:
		{ /* '256' */
			return "HighDom"
		}
	case 257:
		{ /* '257' */
			return "Top Services"
		}
	case 258:
		{ /* '258' */
			return "ambiHome"
		}
	case 259:
		{ /* '259' */
			return "DATEC electronic AG"
		}
	case 260:
		{ /* '260' */
			return "ABUS Security-Center"
		}
	case 261:
		{ /* '261' */
			return "Lite-Puter"
		}
	case 262:
		{ /* '262' */
			return "Tantron Electronic"
		}
	case 263:
		{ /* '263' */
			return "Interra"
		}
	case 264:
		{ /* '264' */
			return "DKX Tech"
		}
	case 265:
		{ /* '265' */
			return "Viatron"
		}
	case 266:
		{ /* '266' */
			return "Nautibus"
		}
	case 267:
		{ /* '267' */
			return "ON Semiconductor"
		}
	case 268:
		{ /* '268' */
			return "Longchuang"
		}
	case 269:
		{ /* '269' */
			return "Air-On AG"
		}
	case 27:
		{ /* '27' */
			return "Albert Ackermann"
		}
	case 270:
		{ /* '270' */
			return "ib-company GmbH"
		}
	case 271:
		{ /* '271' */
			return "Sation Factory"
		}
	case 272:
		{ /* '272' */
			return "Agentilo GmbH"
		}
	case 273:
		{ /* '273' */
			return "Makel Elektrik"
		}
	case 274:
		{ /* '274' */
			return "Helios Ventilatoren"
		}
	case 275:
		{ /* '275' */
			return "Otto Solutions Pte Ltd"
		}
	case 276:
		{ /* '276' */
			return "Airmaster"
		}
	case 277:
		{ /* '277' */
			return "Vallox GmbH"
		}
	case 278:
		{ /* '278' */
			return "Dalitek"
		}
	case 279:
		{ /* '279' */
			return "ASIN"
		}
	case 28:
		{ /* '28' */
			return "Schupa GmbH"
		}
	case 280:
		{ /* '280' */
			return "Bridges Intelligence Technology Inc."
		}
	case 281:
		{ /* '281' */
			return "ARBONIA"
		}
	case 282:
		{ /* '282' */
			return "KERMI"
		}
	case 283:
		{ /* '283' */
			return "PROLUX"
		}
	case 284:
		{ /* '284' */
			return "ClicHome"
		}
	case 285:
		{ /* '285' */
			return "COMMAX"
		}
	case 286:
		{ /* '286' */
			return "EAE"
		}
	case 287:
		{ /* '287' */
			return "Tense"
		}
	case 288:
		{ /* '288' */
			return "Seyoung Electronics"
		}
	case 289:
		{ /* '289' */
			return "Lifedomus"
		}
	case 29:
		{ /* '29' */
			return "ABB SCHWEIZ"
		}
	case 290:
		{ /* '290' */
			return "EUROtronic Technology GmbH"
		}
	case 291:
		{ /* '291' */
			return "tci"
		}
	case 292:
		{ /* '292' */
			return "Rishun Electronic"
		}
	case 293:
		{ /* '293' */
			return "Zipato"
		}
	case 294:
		{ /* '294' */
			return "cm-security GmbH & Co KG"
		}
	case 295:
		{ /* '295' */
			return "Qing Cables"
		}
	case 296:
		{ /* '296' */
			return "LABIO"
		}
	case 297:
		{ /* '297' */
			return "Coster Tecnologie Elettroniche S.p.A."
		}
	case 298:
		{ /* '298' */
			return "E.G.E"
		}
	case 299:
		{ /* '299' */
			return "NETxAutomation"
		}
	case 30:
		{ /* '30' */
			return "Feller"
		}
	case 300:
		{ /* '300' */
			return "tecalor"
		}
	case 301:
		{ /* '301' */
			return "Urmet Electronics (Huizhou) Ltd."
		}
	case 302:
		{ /* '302' */
			return "Peiying Building Control"
		}
	case 303:
		{ /* '303' */
			return "BPT S.p.A. a Socio Unico"
		}
	case 304:
		{ /* '304' */
			return "Kanontec - KanonBUS"
		}
	case 305:
		{ /* '305' */
			return "ISER Tech"
		}
	case 306:
		{ /* '306' */
			return "Fineline"
		}
	case 307:
		{ /* '307' */
			return "CP Electronics Ltd"
		}
	case 308:
		{ /* '308' */
			return "Niko-Servodan A/S"
		}
	case 309:
		{ /* '309' */
			return "Simon"
		}
	case 31:
		{ /* '31' */
			return "Glamox AS"
		}
	case 310:
		{ /* '310' */
			return "GM modular pvt. Ltd."
		}
	case 311:
		{ /* '311' */
			return "FU CHENG Intelligence"
		}
	case 312:
		{ /* '312' */
			return "NexKon"
		}
	case 313:
		{ /* '313' */
			return "FEEL s.r.l"
		}
	case 314:
		{ /* '314' */
			return "Not Assigned"
		}
	case 315:
		{ /* '315' */
			return "Shenzhen Fanhai Sanjiang Electronics Co., Ltd."
		}
	case 316:
		{ /* '316' */
			return "Jiuzhou Greeble"
		}
	case 317:
		{ /* '317' */
			return "Aumüller Aumatic GmbH"
		}
	case 318:
		{ /* '318' */
			return "Etman Electric"
		}
	case 319:
		{ /* '319' */
			return "Black Nova"
		}
	case 32:
		{ /* '32' */
			return "DEHN & SÖHNE"
		}
	case 320:
		{ /* '320' */
			return "ZidaTech AG"
		}
	case 321:
		{ /* '321' */
			return "IDGS bvba"
		}
	case 322:
		{ /* '322' */
			return "dakanimo"
		}
	case 323:
		{ /* '323' */
			return "Trebor Automation AB"
		}
	case 324:
		{ /* '324' */
			return "Satel sp. z o.o."
		}
	case 325:
		{ /* '325' */
			return "Russound, Inc."
		}
	case 326:
		{ /* '326' */
			return "Midea Heating & Ventilating Equipment CO LTD"
		}
	case 327:
		{ /* '327' */
			return "Consorzio Terranuova"
		}
	case 328:
		{ /* '328' */
			return "Wolf Heiztechnik GmbH"
		}
	case 329:
		{ /* '329' */
			return "SONTEC"
		}
	case 33:
		{ /* '33' */
			return "CRABTREE"
		}
	case 330:
		{ /* '330' */
			return "Belcom Cables Ltd."
		}
	case 331:
		{ /* '331' */
			return "Guangzhou SeaWin Electrical Technologies Co., Ltd."
		}
	case 332:
		{ /* '332' */
			return "Acrel"
		}
	case 333:
		{ /* '333' */
			return "Franke Aquarotter GmbH"
		}
	case 334:
		{ /* '334' */
			return "Orion Systems"
		}
	case 335:
		{ /* '335' */
			return "Schrack Technik GmbH"
		}
	case 336:
		{ /* '336' */
			return "INSPRID"
		}
	case 337:
		{ /* '337' */
			return "Sunricher"
		}
	case 338:
		{ /* '338' */
			return "Menred automation system(shanghai) Co.,Ltd."
		}
	case 339:
		{ /* '339' */
			return "Aurex"
		}
	case 34:
		{ /* '34' */
			return "eVoKNX"
		}
	case 340:
		{ /* '340' */
			return "Josef Barthelme GmbH & Co. KG"
		}
	case 341:
		{ /* '341' */
			return "Architecture Numerique"
		}
	case 342:
		{ /* '342' */
			return "UP GROUP"
		}
	case 343:
		{ /* '343' */
			return "Teknos-Avinno"
		}
	case 344:
		{ /* '344' */
			return "Ningbo Dooya Mechanic & Electronic Technology"
		}
	case 345:
		{ /* '345' */
			return "Thermokon Sensortechnik GmbH"
		}
	case 346:
		{ /* '346' */
			return "BELIMO Automation AG"
		}
	case 347:
		{ /* '347' */
			return "Zehnder Group International AG"
		}
	case 348:
		{ /* '348' */
			return "sks Kinkel Elektronik"
		}
	case 349:
		{ /* '349' */
			return "ECE Wurmitzer GmbH"
		}
	case 350:
		{ /* '350' */
			return "LARS"
		}
	case 351:
		{ /* '351' */
			return "URC"
		}
	case 352:
		{ /* '352' */
			return "LightControl"
		}
	case 353:
		{ /* '353' */
			return "ShenZhen YM"
		}
	case 354:
		{ /* '354' */
			return "MEAN WELL Enterprises Co. Ltd."
		}
	case 355:
		{ /* '355' */
			return "OSix"
		}
	case 356:
		{ /* '356' */
			return "AYPRO Technology"
		}
	case 357:
		{ /* '357' */
			return "Hefei Ecolite Software"
		}
	case 358:
		{ /* '358' */
			return "Enno"
		}
	case 359:
		{ /* '359' */
			return "OHOSURE"
		}
	case 36:
		{ /* '36' */
			return "Paul Hochköpper"
		}
	case 360:
		{ /* '360' */
			return "Garefowl"
		}
	case 361:
		{ /* '361' */
			return "GEZE"
		}
	case 362:
		{ /* '362' */
			return "LG Electronics Inc."
		}
	case 363:
		{ /* '363' */
			return "SMC interiors"
		}
	case 364:
		{ /* '364' */
			return "Not Assigned"
		}
	case 365:
		{ /* '365' */
			return "SCS Cable"
		}
	case 366:
		{ /* '366' */
			return "Hoval"
		}
	case 367:
		{ /* '367' */
			return "CANST"
		}
	case 368:
		{ /* '368' */
			return "HangZhou Berlin"
		}
	case 369:
		{ /* '369' */
			return "EVN-Lichttechnik"
		}
	case 37:
		{ /* '37' */
			return "Altenburger Electronic"
		}
	case 370:
		{ /* '370' */
			return "rutec"
		}
	case 371:
		{ /* '371' */
			return "Finder"
		}
	case 372:
		{ /* '372' */
			return "Fujitsu General Limited"
		}
	case 373:
		{ /* '373' */
			return "ZF Friedrichshafen AG"
		}
	case 374:
		{ /* '374' */
			return "Crealed"
		}
	case 375:
		{ /* '375' */
			return "Miles Magic Automation Private Limited"
		}
	case 376:
		{ /* '376' */
			return "E+"
		}
	case 377:
		{ /* '377' */
			return "Italcond"
		}
	case 378:
		{ /* '378' */
			return "SATION"
		}
	case 379:
		{ /* '379' */
			return "NewBest"
		}
	case 380:
		{ /* '380' */
			return "GDS DIGITAL SYSTEMS"
		}
	case 381:
		{ /* '381' */
			return "Iddero"
		}
	case 382:
		{ /* '382' */
			return "MBNLED"
		}
	case 383:
		{ /* '383' */
			return "VITRUM"
		}
	case 384:
		{ /* '384' */
			return "ekey biometric systems GmbH"
		}
	case 385:
		{ /* '385' */
			return "AMC"
		}
	case 386:
		{ /* '386' */
			return "TRILUX GmbH & Co. KG"
		}
	case 387:
		{ /* '387' */
			return "WExcedo"
		}
	case 388:
		{ /* '388' */
			return "VEMER SPA"
		}
	case 389:
		{ /* '389' */
			return "Alexander Bürkle GmbH & Co KG"
		}
	case 390:
		{ /* '390' */
			return "Citron"
		}
	case 391:
		{ /* '391' */
			return "Shenzhen HeGuang"
		}
	case 392:
		{ /* '392' */
			return "Not Assigned"
		}
	case 393:
		{ /* '393' */
			return "TRANE B.V.B.A"
		}
	case 394:
		{ /* '394' */
			return "CAREL"
		}
	case 395:
		{ /* '395' */
			return "Prolite Controls"
		}
	case 396:
		{ /* '396' */
			return "BOSMER"
		}
	case 397:
		{ /* '397' */
			return "EUCHIPS"
		}
	case 398:
		{ /* '398' */
			return "connect (Thinka connect)"
		}
	case 399:
		{ /* '399' */
			return "PEAKnx a DOGAWIST company"
		}
	case 4:
		{ /* '4' */
			return "Albrecht Jung"
		}
	case 400:
		{ /* '400' */
			return "ACEMATIC"
		}
	case 401:
		{ /* '401' */
			return "ELAUSYS"
		}
	case 402:
		{ /* '402' */
			return "ITK Engineering AG"
		}
	case 403:
		{ /* '403' */
			return "INTEGRA METERING AG"
		}
	case 404:
		{ /* '404' */
			return "FMS Hospitality Pte Ltd"
		}
	case 405:
		{ /* '405' */
			return "Nuvo"
		}
	case 406:
		{ /* '406' */
			return "u::Lux GmbH"
		}
	case 407:
		{ /* '407' */
			return "Brumberg Leuchten"
		}
	case 408:
		{ /* '408' */
			return "Lime"
		}
	case 409:
		{ /* '409' */
			return "Great Empire International Group Co., Ltd."
		}
	case 41:
		{ /* '41' */
			return "Grässlin"
		}
	case 410:
		{ /* '410' */
			return "Kavoshpishro Asia"
		}
	case 411:
		{ /* '411' */
			return "V2 SpA"
		}
	case 412:
		{ /* '412' */
			return "Johnson Controls"
		}
	case 413:
		{ /* '413' */
			return "Arkud"
		}
	case 414:
		{ /* '414' */
			return "Iridium Ltd."
		}
	case 415:
		{ /* '415' */
			return "bsmart"
		}
	case 416:
		{ /* '416' */
			return "BAB TECHNOLOGIE GmbH"
		}
	case 417:
		{ /* '417' */
			return "NICE Spa"
		}
	case 418:
		{ /* '418' */
			return "Redfish Group Pty Ltd"
		}
	case 419:
		{ /* '419' */
			return "SABIANA spa"
		}
	case 42:
		{ /* '42' */
			return "Simon"
		}
	case 420:
		{ /* '420' */
			return "Ubee Interactive Europe"
		}
	case 421:
		{ /* '421' */
			return "Rexel"
		}
	case 422:
		{ /* '422' */
			return "Ges Teknik A.S."
		}
	case 423:
		{ /* '423' */
			return "Ave S.p.A."
		}
	case 424:
		{ /* '424' */
			return "Zhuhai Ltech Technology Co., Ltd."
		}
	case 425:
		{ /* '425' */
			return "ARCOM"
		}
	case 426:
		{ /* '426' */
			return "VIA Technologies, Inc."
		}
	case 427:
		{ /* '427' */
			return "FEELSMART."
		}
	case 428:
		{ /* '428' */
			return "SUPCON"
		}
	case 429:
		{ /* '429' */
			return "MANIC"
		}
	case 430:
		{ /* '430' */
			return "TDE GmbH"
		}
	case 431:
		{ /* '431' */
			return "Nanjing Shufan Information technology Co.,Ltd."
		}
	case 432:
		{ /* '432' */
			return "EWTech"
		}
	case 433:
		{ /* '433' */
			return "Kluger Automation GmbH"
		}
	case 434:
		{ /* '434' */
			return "JoongAng Control"
		}
	case 435:
		{ /* '435' */
			return "GreenControls Technology Sdn. Bhd."
		}
	case 436:
		{ /* '436' */
			return "IME S.p.a."
		}
	case 437:
		{ /* '437' */
			return "SiChuan HaoDing"
		}
	case 438:
		{ /* '438' */
			return "Mindjaga Ltd."
		}
	case 439:
		{ /* '439' */
			return "RuiLi Smart Control"
		}
	case 43954:
		{ /* '43954' */
			return "ABB - reserved"
		}
	case 43959:
		{ /* '43959' */
			return "Busch-Jaeger Elektro - reserved"
		}
	case 44:
		{ /* '44' */
			return "VIMAR"
		}
	case 440:
		{ /* '440' */
			return "CODESYS GmbH"
		}
	case 441:
		{ /* '441' */
			return "Moorgen Deutschland GmbH"
		}
	case 442:
		{ /* '442' */
			return "CULLMANN TECH"
		}
	case 443:
		{ /* '443' */
			return "Merck Window Technologies B.V."
		}
	case 444:
		{ /* '444' */
			return "ABEGO"
		}
	case 445:
		{ /* '445' */
			return "myGEKKO"
		}
	case 446:
		{ /* '446' */
			return "Ergo3 Sarl"
		}
	case 447:
		{ /* '447' */
			return "STmicroelectronics International N.V."
		}
	case 448:
		{ /* '448' */
			return "cjc systems"
		}
	case 449:
		{ /* '449' */
			return "Sudoku"
		}
	case 45:
		{ /* '45' */
			return "Moeller Gebäudeautomation KG"
		}
	case 451:
		{ /* '451' */
			return "AZ e-lite Pte Ltd"
		}
	case 452:
		{ /* '452' */
			return "Arlight"
		}
	case 453:
		{ /* '453' */
			return "Grünbeck Wasseraufbereitung GmbH"
		}
	case 454:
		{ /* '454' */
			return "Module Electronic"
		}
	case 455:
		{ /* '455' */
			return "KOPLAT"
		}
	case 456:
		{ /* '456' */
			return "Guangzhou Letour Life Technology Co., Ltd"
		}
	case 457:
		{ /* '457' */
			return "ILEVIA"
		}
	case 458:
		{ /* '458' */
			return "LN SYSTEMTEQ"
		}
	case 459:
		{ /* '459' */
			return "Hisense SmartHome"
		}
	case 46:
		{ /* '46' */
			return "Eltako"
		}
	case 460:
		{ /* '460' */
			return "Flink Automation System"
		}
	case 461:
		{ /* '461' */
			return "xxter bv"
		}
	case 462:
		{ /* '462' */
			return "lynxus technology"
		}
	case 463:
		{ /* '463' */
			return "ROBOT S.A."
		}
	case 464:
		{ /* '464' */
			return "Shenzhen Atte Smart Life Co.,Ltd."
		}
	case 465:
		{ /* '465' */
			return "Noblesse"
		}
	case 466:
		{ /* '466' */
			return "Advanced Devices"
		}
	case 467:
		{ /* '467' */
			return "Atrina Building Automation Co. Ltd"
		}
	case 468:
		{ /* '468' */
			return "Guangdong Daming Laffey electric Co., Ltd."
		}
	case 469:
		{ /* '469' */
			return "Westerstrand Urfabrik AB"
		}
	case 470:
		{ /* '470' */
			return "Control4 Corporate"
		}
	case 471:
		{ /* '471' */
			return "Ontrol"
		}
	case 472:
		{ /* '472' */
			return "Starnet"
		}
	case 473:
		{ /* '473' */
			return "BETA CAVI"
		}
	case 474:
		{ /* '474' */
			return "EaseMore"
		}
	case 475:
		{ /* '475' */
			return "Vivaldi srl"
		}
	case 476:
		{ /* '476' */
			return "Gree Electric Appliances,Inc. of Zhuhai"
		}
	case 477:
		{ /* '477' */
			return "HWISCON"
		}
	case 478:
		{ /* '478' */
			return "Shanghai ELECON Intelligent Technology Co., Ltd."
		}
	case 479:
		{ /* '479' */
			return "Kampmann"
		}
	case 480:
		{ /* '480' */
			return "Impolux GmbH / LEDIMAX"
		}
	case 481:
		{ /* '481' */
			return "Evaux"
		}
	case 482:
		{ /* '482' */
			return "Webro Cables & Connectors Limited"
		}
	case 483:
		{ /* '483' */
			return "Shanghai E-tech Solution"
		}
	case 484:
		{ /* '484' */
			return "Guangzhou HOKO Electric Co.,Ltd."
		}
	case 485:
		{ /* '485' */
			return "LAMMIN HIGH TECH CO.,LTD"
		}
	case 486:
		{ /* '486' */
			return "Shenzhen Merrytek Technology Co., Ltd"
		}
	case 487:
		{ /* '487' */
			return "I-Luxus"
		}
	case 488:
		{ /* '488' */
			return "Elmos Semiconductor AG"
		}
	case 489:
		{ /* '489' */
			return "EmCom Technology Inc"
		}
	case 49:
		{ /* '49' */
			return "Bosch-Siemens Haushaltsgeräte"
		}
	case 490:
		{ /* '490' */
			return "project innovations GmbH"
		}
	case 491:
		{ /* '491' */
			return "Itc"
		}
	case 492:
		{ /* '492' */
			return "ABB LV Installation Materials Company Ltd, Beijing"
		}
	case 493:
		{ /* '493' */
			return "Maico"
		}
	case 495:
		{ /* '495' */
			return "ELAN SRL"
		}
	case 496:
		{ /* '496' */
			return "MinhHa Technology co.,Ltd"
		}
	case 497:
		{ /* '497' */
			return "Zhejiang Tianjie Industrial CORP."
		}
	case 498:
		{ /* '498' */
			return "iAutomation Pty Limited"
		}
	case 499:
		{ /* '499' */
			return "Extron"
		}
	case 5:
		{ /* '5' */
			return "Bticino"
		}
	case 500:
		{ /* '500' */
			return "Freedompro"
		}
	case 501:
		{ /* '501' */
			return "1Home"
		}
	case 502:
		{ /* '502' */
			return "EOS Saunatechnik GmbH"
		}
	case 503:
		{ /* '503' */
			return "KUSATEK GmbH"
		}
	case 504:
		{ /* '504' */
			return "EisBär Scada"
		}
	case 505:
		{ /* '505' */
			return "AUTOMATISMI BENINCA S.P.A."
		}
	case 506:
		{ /* '506' */
			return "Blendom"
		}
	case 507:
		{ /* '507' */
			return "Madel Air Technical diffusion"
		}
	case 508:
		{ /* '508' */
			return "NIKO"
		}
	case 509:
		{ /* '509' */
			return "Bosch Rexroth AG"
		}
	case 512:
		{ /* '512' */
			return "C&M Products"
		}
	case 513:
		{ /* '513' */
			return "Hörmann KG Verkaufsgesellschaft"
		}
	case 514:
		{ /* '514' */
			return "Shanghai Rajayasa co.,LTD"
		}
	case 515:
		{ /* '515' */
			return "SUZUKI"
		}
	case 516:
		{ /* '516' */
			return "Silent Gliss International Ltd."
		}
	case 517:
		{ /* '517' */
			return "BEE Controls (ADGSC Group)"
		}
	case 518:
		{ /* '518' */
			return "xDTecGmbH"
		}
	case 519:
		{ /* '519' */
			return "OSRAM"
		}
	case 52:
		{ /* '52' */
			return "RITTO GmbH&Co.KG"
		}
	case 520:
		{ /* '520' */
			return "Lebenor"
		}
	case 521:
		{ /* '521' */
			return "automaneng"
		}
	case 522:
		{ /* '522' */
			return "Honeywell Automation Solution control(China)"
		}
	case 523:
		{ /* '523' */
			return "Hangzhou binthen Intelligence Technology Co.,Ltd"
		}
	case 524:
		{ /* '524' */
			return "ETA Heiztechnik"
		}
	case 525:
		{ /* '525' */
			return "DIVUS GmbH"
		}
	case 526:
		{ /* '526' */
			return "Nanjing Taijiesai Intelligent Technology Co. Ltd."
		}
	case 527:
		{ /* '527' */
			return "Lunatone"
		}
	case 528:
		{ /* '528' */
			return "ZHEJIANG SCTECH BUILDING INTELLIGENT"
		}
	case 529:
		{ /* '529' */
			return "Foshan Qite Technology Co., Ltd."
		}
	case 53:
		{ /* '53' */
			return "Power Controls"
		}
	case 530:
		{ /* '530' */
			return "NOKE"
		}
	case 531:
		{ /* '531' */
			return "LANDCOM"
		}
	case 532:
		{ /* '532' */
			return "Stork AS"
		}
	case 533:
		{ /* '533' */
			return "Hangzhou Shendu Technology Co., Ltd."
		}
	case 534:
		{ /* '534' */
			return "CoolAutomation"
		}
	case 535:
		{ /* '535' */
			return "Aprstern"
		}
	case 536:
		{ /* '536' */
			return "sonnen"
		}
	case 537:
		{ /* '537' */
			return "DNAKE"
		}
	case 538:
		{ /* '538' */
			return "Neuberger Gebäudeautomation GmbH"
		}
	case 539:
		{ /* '539' */
			return "Stiliger"
		}
	case 540:
		{ /* '540' */
			return "Berghof Automation GmbH"
		}
	case 541:
		{ /* '541' */
			return "Total Automation and controls GmbH"
		}
	case 542:
		{ /* '542' */
			return "dovit"
		}
	case 543:
		{ /* '543' */
			return "Instalighting GmbH"
		}
	case 544:
		{ /* '544' */
			return "UNI-TEC"
		}
	case 545:
		{ /* '545' */
			return "CasaTunes"
		}
	case 546:
		{ /* '546' */
			return "EMT"
		}
	case 547:
		{ /* '547' */
			return "Senfficient"
		}
	case 548:
		{ /* '548' */
			return "Aurolite electrical panyu guangzhou limited"
		}
	case 549:
		{ /* '549' */
			return "ABB Xiamen Smart Technology Co., Ltd."
		}
	case 55:
		{ /* '55' */
			return "ZUMTOBEL"
		}
	case 550:
		{ /* '550' */
			return "Samson Electric Wire"
		}
	case 551:
		{ /* '551' */
			return "T-Touching"
		}
	case 552:
		{ /* '552' */
			return "Core Smart Home"
		}
	case 553:
		{ /* '553' */
			return "GreenConnect Solutions SA"
		}
	case 554:
		{ /* '554' */
			return "ELETTRONICA CONDUTTORI"
		}
	case 555:
		{ /* '555' */
			return "MKFC"
		}
	case 556:
		{ /* '556' */
			return "Automation+"
		}
	case 557:
		{ /* '557' */
			return "blue and red"
		}
	case 558:
		{ /* '558' */
			return "frogblue"
		}
	case 559:
		{ /* '559' */
			return "SAVESOR"
		}
	case 560:
		{ /* '560' */
			return "App Tech"
		}
	case 561:
		{ /* '561' */
			return "sensortec AG"
		}
	case 562:
		{ /* '562' */
			return "nysa technology & solutions"
		}
	case 563:
		{ /* '563' */
			return "FARADITE"
		}
	case 564:
		{ /* '564' */
			return "Optimus"
		}
	case 565:
		{ /* '565' */
			return "KTS s.r.l."
		}
	case 566:
		{ /* '566' */
			return "Ramcro SPA"
		}
	case 567:
		{ /* '567' */
			return "Wuhan WiseCreate Universe Technology Co., Ltd"
		}
	case 568:
		{ /* '568' */
			return "BEMI Smart Home Ltd"
		}
	case 569:
		{ /* '569' */
			return "Ardomus"
		}
	case 57:
		{ /* '57' */
			return "Phoenix Contact"
		}
	case 570:
		{ /* '570' */
			return "ChangXing"
		}
	case 571:
		{ /* '571' */
			return "E-Controls"
		}
	case 572:
		{ /* '572' */
			return "AIB Technology"
		}
	case 573:
		{ /* '573' */
			return "NVC"
		}
	case 574:
		{ /* '574' */
			return "Kbox"
		}
	case 575:
		{ /* '575' */
			return "CNS"
		}
	case 576:
		{ /* '576' */
			return "Tyba"
		}
	case 577:
		{ /* '577' */
			return "Atrel"
		}
	case 578:
		{ /* '578' */
			return "Simon Electric (China) Co., LTD"
		}
	case 579:
		{ /* '579' */
			return "Kordz Group"
		}
	case 580:
		{ /* '580' */
			return "ND Electric"
		}
	case 581:
		{ /* '581' */
			return "Controlium"
		}
	case 582:
		{ /* '582' */
			return "FAMO GmbH & Co. KG"
		}
	case 583:
		{ /* '583' */
			return "CDN Smart"
		}
	case 584:
		{ /* '584' */
			return "Heston"
		}
	case 585:
		{ /* '585' */
			return "ESLA CONEXIONES S.L."
		}
	case 586:
		{ /* '586' */
			return "Weishaupt"
		}
	case 587:
		{ /* '587' */
			return "ASTRUM TECHNOLOGY"
		}
	case 588:
		{ /* '588' */
			return "WUERTH ELEKTRONIK STELVIO KONTEK S.p.A."
		}
	case 589:
		{ /* '589' */
			return "NANOTECO corporation"
		}
	case 590:
		{ /* '590' */
			return "Nietian"
		}
	case 591:
		{ /* '591' */
			return "Sumsir"
		}
	case 592:
		{ /* '592' */
			return "ORBIS TECNOLOGIA ELECTRICA SA"
		}
	case 6:
		{ /* '6' */
			return "Berker"
		}
	case 61:
		{ /* '61' */
			return "WAGO Kontakttechnik"
		}
	case 62:
		{ /* '62' */
			return "knXpresso"
		}
	case 66:
		{ /* '66' */
			return "Wieland Electric"
		}
	case 67:
		{ /* '67' */
			return "Hermann Kleinhuis"
		}
	case 69:
		{ /* '69' */
			return "Stiebel Eltron"
		}
	case 7:
		{ /* '7' */
			return "Busch-Jaeger Elektro"
		}
	case 71:
		{ /* '71' */
			return "Tehalit"
		}
	case 72:
		{ /* '72' */
			return "Theben AG"
		}
	case 73:
		{ /* '73' */
			return "Wilhelm Rutenbeck"
		}
	case 75:
		{ /* '75' */
			return "Winkhaus"
		}
	case 76:
		{ /* '76' */
			return "Robert Bosch"
		}
	case 78:
		{ /* '78' */
			return "Somfy"
		}
	case 8:
		{ /* '8' */
			return "GIRA Giersiepen"
		}
	case 80:
		{ /* '80' */
			return "Woertz"
		}
	case 81:
		{ /* '81' */
			return "Viessmann Werke"
		}
	case 82:
		{ /* '82' */
			return "IMI Hydronic Engineering"
		}
	case 83:
		{ /* '83' */
			return "Joh. Vaillant"
		}
	case 85:
		{ /* '85' */
			return "AMP Deutschland"
		}
	case 89:
		{ /* '89' */
			return "Bosch Thermotechnik GmbH"
		}
	case 9:
		{ /* '9' */
			return "Hager Electro"
		}
	case 90:
		{ /* '90' */
			return "SEF - ECOTEC"
		}
	case 92:
		{ /* '92' */
			return "DORMA GmbH + Co. KG"
		}
	case 93:
		{ /* '93' */
			return "WindowMaster A/S"
		}
	case 94:
		{ /* '94' */
			return "Walther Werke"
		}
	case 95:
		{ /* '95' */
			return "ORAS"
		}
	case 97:
		{ /* '97' */
			return "Dätwyler"
		}
	case 98:
		{ /* '98' */
			return "Electrak"
		}
	case 99:
		{ /* '99' */
			return "Techem"
		}
	default:
		{
			return ""
		}
	}
}
func KnxManufacturerValueOf(value uint16) KnxManufacturer {
	switch value {
	case 1:
		return KnxManufacturer_SIEMENS
	case 10:
		return KnxManufacturer_INSTA_GMBH
	case 100:
		return KnxManufacturer_SCHNEIDER_ELECTRIC_INDUSTRIES_SAS
	case 101:
		return KnxManufacturer_WHD_WILHELM_HUBER_Plus_SOEHNE
	case 102:
		return KnxManufacturer_BISCHOFF_ELEKTRONIK
	case 104:
		return KnxManufacturer_JEPAZ
	case 105:
		return KnxManufacturer_RTS_AUTOMATION
	case 106:
		return KnxManufacturer_EIBMARKT_GMBH
	case 107:
		return KnxManufacturer_WAREMA_RENKHOFF_SE
	case 108:
		return KnxManufacturer_EELECTRON
	case 109:
		return KnxManufacturer_BELDEN_WIRE_AND_CABLE_B_V_
	case 11:
		return KnxManufacturer_LEGRAND_APPAREILLAGE_ELECTRIQUE
	case 110:
		return KnxManufacturer_BECKER_ANTRIEBE_GMBH
	case 111:
		return KnxManufacturer_J_STEHLEPlusSOEHNE_GMBH
	case 112:
		return KnxManufacturer_AGFEO
	case 113:
		return KnxManufacturer_ZENNIO
	case 114:
		return KnxManufacturer_TAPKO_TECHNOLOGIES
	case 115:
		return KnxManufacturer_HDL
	case 116:
		return KnxManufacturer_UPONOR
	case 117:
		return KnxManufacturer_SE_LIGHTMANAGEMENT_AG
	case 118:
		return KnxManufacturer_ARCUS_EDS
	case 119:
		return KnxManufacturer_INTESIS
	case 12:
		return KnxManufacturer_MERTEN
	case 120:
		return KnxManufacturer_HERHOLDT_CONTROLS_SRL
	case 121:
		return KnxManufacturer_NIKO_ZUBLIN
	case 122:
		return KnxManufacturer_DURABLE_TECHNOLOGIES
	case 123:
		return KnxManufacturer_INNOTEAM
	case 124:
		return KnxManufacturer_ISE_GMBH
	case 125:
		return KnxManufacturer_TEAM_FOR_TRONICS
	case 126:
		return KnxManufacturer_CIAT
	case 127:
		return KnxManufacturer_REMEHA_BV
	case 128:
		return KnxManufacturer_ESYLUX
	case 129:
		return KnxManufacturer_BASALTE
	case 130:
		return KnxManufacturer_VESTAMATIC
	case 131:
		return KnxManufacturer_MDT_TECHNOLOGIES
	case 132:
		return KnxManufacturer_WARENDORFER_KUECHEN_GMBH
	case 133:
		return KnxManufacturer_VIDEO_STAR
	case 134:
		return KnxManufacturer_SITEK
	case 135:
		return KnxManufacturer_CONTROLTRONIC
	case 136:
		return KnxManufacturer_FUNCTION_TECHNOLOGY
	case 137:
		return KnxManufacturer_AMX
	case 138:
		return KnxManufacturer_ELDAT
	case 139:
		return KnxManufacturer_PANASONIC
	case 14:
		return KnxManufacturer_ABB_SPA_SACE_DIVISION
	case 140:
		return KnxManufacturer_PULSE_TECHNOLOGIES
	case 141:
		return KnxManufacturer_CRESTRON
	case 142:
		return KnxManufacturer_STEINEL_PROFESSIONAL
	case 143:
		return KnxManufacturer_BILTON_LED_LIGHTING
	case 144:
		return KnxManufacturer_DENRO_AG
	case 145:
		return KnxManufacturer_GEPRO
	case 146:
		return KnxManufacturer_PREUSSEN_AUTOMATION
	case 147:
		return KnxManufacturer_ZOPPAS_INDUSTRIES
	case 148:
		return KnxManufacturer_MACTECH
	case 149:
		return KnxManufacturer_TECHNO_TREND
	case 150:
		return KnxManufacturer_FS_CABLES
	case 151:
		return KnxManufacturer_DELTA_DORE
	case 152:
		return KnxManufacturer_EISSOUND
	case 153:
		return KnxManufacturer_CISCO
	case 154:
		return KnxManufacturer_DINUY
	case 155:
		return KnxManufacturer_IKNIX
	case 156:
		return KnxManufacturer_RADEMACHER_GERAETE_ELEKTRONIK_GMBH
	case 157:
		return KnxManufacturer_EGI_ELECTROACUSTICA_GENERAL_IBERICA
	case 158:
		return KnxManufacturer_BES___INGENIUM
	case 159:
		return KnxManufacturer_ELABNET
	case 160:
		return KnxManufacturer_BLUMOTIX
	case 161:
		return KnxManufacturer_HUNTER_DOUGLAS
	case 162:
		return KnxManufacturer_APRICUM
	case 163:
		return KnxManufacturer_TIANSU_AUTOMATION
	case 164:
		return KnxManufacturer_BUBENDORFF
	case 165:
		return KnxManufacturer_MBS_GMBH
	case 166:
		return KnxManufacturer_ENERTEX_BAYERN_GMBH
	case 167:
		return KnxManufacturer_BMS
	case 168:
		return KnxManufacturer_SINAPSI
	case 169:
		return KnxManufacturer_EMBEDDED_SYSTEMS_SIA
	case 170:
		return KnxManufacturer_KNX1
	case 171:
		return KnxManufacturer_TOKKA
	case 172:
		return KnxManufacturer_NANOSENSE
	case 173:
		return KnxManufacturer_PEAR_AUTOMATION_GMBH
	case 174:
		return KnxManufacturer_DGA
	case 175:
		return KnxManufacturer_LUTRON
	case 176:
		return KnxManufacturer_AIRZONE___ALTRA
	case 177:
		return KnxManufacturer_LITHOSS_DESIGN_SWITCHES
	case 178:
		return KnxManufacturer_3ATEL
	case 179:
		return KnxManufacturer_PHILIPS_CONTROLS
	case 180:
		return KnxManufacturer_VELUX_AS
	case 181:
		return KnxManufacturer_LOYTEC
	case 182:
		return KnxManufacturer_EKINEX_S_P_A_
	case 183:
		return KnxManufacturer_SIRLAN_TECHNOLOGIES
	case 184:
		return KnxManufacturer_PROKNX_SAS
	case 185:
		return KnxManufacturer_IT_GMBH
	case 186:
		return KnxManufacturer_RENSON
	case 187:
		return KnxManufacturer_HEP_GROUP
	case 188:
		return KnxManufacturer_BALMART
	case 189:
		return KnxManufacturer_GFS_GMBH
	case 190:
		return KnxManufacturer_SCHENKER_STOREN_AG
	case 191:
		return KnxManufacturer_ALGODUE_ELETTRONICA_S_R_L_
	case 192:
		return KnxManufacturer_ABB_FRANCE
	case 193:
		return KnxManufacturer_MAINTRONIC
	case 194:
		return KnxManufacturer_VANTAGE
	case 195:
		return KnxManufacturer_FORESIS
	case 196:
		return KnxManufacturer_RESEARCH_AND_PRODUCTION_ASSOCIATION_SEM
	case 197:
		return KnxManufacturer_WEINZIERL_ENGINEERING_GMBH
	case 198:
		return KnxManufacturer_MOEHLENHOFF_WAERMETECHNIK_GMBH
	case 199:
		return KnxManufacturer_PKC_GROUP_OYJ
	case 2:
		return KnxManufacturer_ABB
	case 200:
		return KnxManufacturer_B_E_G_
	case 201:
		return KnxManufacturer_ELSNER_ELEKTRONIK_GMBH
	case 202:
		return KnxManufacturer_SIEMENS_BUILDING_TECHNOLOGIES_HKCHINA_LTD_
	case 204:
		return KnxManufacturer_EUTRAC
	case 205:
		return KnxManufacturer_GUSTAV_HENSEL_GMBH_AND_CO__KG
	case 206:
		return KnxManufacturer_GARO_AB
	case 207:
		return KnxManufacturer_WALDMANN_LICHTTECHNIK
	case 208:
		return KnxManufacturer_SCHUECO
	case 209:
		return KnxManufacturer_EMU
	case 210:
		return KnxManufacturer_JNET_SYSTEMS_AG
	case 211:
		return KnxManufacturer_TOTAL_SOLUTION_GMBH
	case 214:
		return KnxManufacturer_O_Y_L__ELECTRONICS
	case 215:
		return KnxManufacturer_GALAX_SYSTEM
	case 216:
		return KnxManufacturer_DISCH
	case 217:
		return KnxManufacturer_AUCOTEAM
	case 218:
		return KnxManufacturer_LUXMATE_CONTROLS
	case 219:
		return KnxManufacturer_DANFOSS
	case 22:
		return KnxManufacturer_SIEDLE_AND_SOEHNE
	case 220:
		return KnxManufacturer_AST_GMBH
	case 222:
		return KnxManufacturer_WILA_LEUCHTEN
	case 223:
		return KnxManufacturer_BPlusB_AUTOMATIONS__UND_STEUERUNGSTECHNIK
	case 225:
		return KnxManufacturer_LINGG_AND_JANKE
	case 227:
		return KnxManufacturer_SAUTER
	case 228:
		return KnxManufacturer_SIMU
	case 232:
		return KnxManufacturer_THEBEN_HTS_AG
	case 233:
		return KnxManufacturer_AMANN_GMBH
	case 234:
		return KnxManufacturer_BERG_ENERGIEKONTROLLSYSTEME_GMBH
	case 235:
		return KnxManufacturer_HUEPPE_FORM_SONNENSCHUTZSYSTEME_GMBH
	case 237:
		return KnxManufacturer_OVENTROP_KG
	case 238:
		return KnxManufacturer_GRIESSER_AG
	case 239:
		return KnxManufacturer_IPAS_GMBH
	case 24:
		return KnxManufacturer_EBERLE
	case 240:
		return KnxManufacturer_ELERO_GMBH
	case 241:
		return KnxManufacturer_ARDAN_PRODUCTION_AND_INDUSTRIAL_CONTROLS_LTD_
	case 242:
		return KnxManufacturer_METEC_MESSTECHNIK_GMBH
	case 244:
		return KnxManufacturer_ELKA_ELEKTRONIK_GMBH
	case 245:
		return KnxManufacturer_ELEKTROANLAGEN_D__NAGEL
	case 246:
		return KnxManufacturer_TRIDONIC_BAUELEMENTE_GMBH
	case 248:
		return KnxManufacturer_STENGLER_GESELLSCHAFT
	case 249:
		return KnxManufacturer_SCHNEIDER_ELECTRIC_MG
	case 25:
		return KnxManufacturer_GEWISS
	case 250:
		return KnxManufacturer_KNX_ASSOCIATION
	case 251:
		return KnxManufacturer_VIVO
	case 252:
		return KnxManufacturer_HUGO_MUELLER_GMBH_AND_CO_KG
	case 253:
		return KnxManufacturer_SIEMENS_HVAC
	case 254:
		return KnxManufacturer_APT
	case 256:
		return KnxManufacturer_HIGHDOM
	case 257:
		return KnxManufacturer_TOP_SERVICES
	case 258:
		return KnxManufacturer_AMBIHOME
	case 259:
		return KnxManufacturer_DATEC_ELECTRONIC_AG
	case 260:
		return KnxManufacturer_ABUS_SECURITY_CENTER
	case 261:
		return KnxManufacturer_LITE_PUTER
	case 262:
		return KnxManufacturer_TANTRON_ELECTRONIC
	case 263:
		return KnxManufacturer_INTERRA
	case 264:
		return KnxManufacturer_DKX_TECH
	case 265:
		return KnxManufacturer_VIATRON
	case 266:
		return KnxManufacturer_NAUTIBUS
	case 267:
		return KnxManufacturer_ON_SEMICONDUCTOR
	case 268:
		return KnxManufacturer_LONGCHUANG
	case 269:
		return KnxManufacturer_AIR_ON_AG
	case 27:
		return KnxManufacturer_ALBERT_ACKERMANN
	case 270:
		return KnxManufacturer_IB_COMPANY_GMBH
	case 271:
		return KnxManufacturer_SATION_FACTORY
	case 272:
		return KnxManufacturer_AGENTILO_GMBH
	case 273:
		return KnxManufacturer_MAKEL_ELEKTRIK
	case 274:
		return KnxManufacturer_HELIOS_VENTILATOREN
	case 275:
		return KnxManufacturer_OTTO_SOLUTIONS_PTE_LTD
	case 276:
		return KnxManufacturer_AIRMASTER
	case 277:
		return KnxManufacturer_VALLOX_GMBH
	case 278:
		return KnxManufacturer_DALITEK
	case 279:
		return KnxManufacturer_ASIN
	case 28:
		return KnxManufacturer_SCHUPA_GMBH
	case 280:
		return KnxManufacturer_BRIDGES_INTELLIGENCE_TECHNOLOGY_INC_
	case 281:
		return KnxManufacturer_ARBONIA
	case 282:
		return KnxManufacturer_KERMI
	case 283:
		return KnxManufacturer_PROLUX
	case 284:
		return KnxManufacturer_CLICHOME
	case 285:
		return KnxManufacturer_COMMAX
	case 286:
		return KnxManufacturer_EAE
	case 287:
		return KnxManufacturer_TENSE
	case 288:
		return KnxManufacturer_SEYOUNG_ELECTRONICS
	case 289:
		return KnxManufacturer_LIFEDOMUS
	case 29:
		return KnxManufacturer_ABB_SCHWEIZ
	case 290:
		return KnxManufacturer_EUROTRONIC_TECHNOLOGY_GMBH
	case 291:
		return KnxManufacturer_TCI
	case 292:
		return KnxManufacturer_RISHUN_ELECTRONIC
	case 293:
		return KnxManufacturer_ZIPATO
	case 294:
		return KnxManufacturer_CM_SECURITY_GMBH_AND_CO_KG
	case 295:
		return KnxManufacturer_QING_CABLES
	case 296:
		return KnxManufacturer_LABIO
	case 297:
		return KnxManufacturer_COSTER_TECNOLOGIE_ELETTRONICHE_S_P_A_
	case 298:
		return KnxManufacturer_E_G_E
	case 299:
		return KnxManufacturer_NETXAUTOMATION
	case 30:
		return KnxManufacturer_FELLER
	case 300:
		return KnxManufacturer_TECALOR
	case 301:
		return KnxManufacturer_URMET_ELECTRONICS_HUIZHOU_LTD_
	case 302:
		return KnxManufacturer_PEIYING_BUILDING_CONTROL
	case 303:
		return KnxManufacturer_BPT_S_P_A__A_SOCIO_UNICO
	case 304:
		return KnxManufacturer_KANONTEC___KANONBUS
	case 305:
		return KnxManufacturer_ISER_TECH
	case 306:
		return KnxManufacturer_FINELINE
	case 307:
		return KnxManufacturer_CP_ELECTRONICS_LTD
	case 308:
		return KnxManufacturer_NIKO_SERVODAN_AS
	case 309:
		return KnxManufacturer_SIMON_309
	case 31:
		return KnxManufacturer_GLAMOX_AS
	case 310:
		return KnxManufacturer_GM_MODULAR_PVT__LTD_
	case 311:
		return KnxManufacturer_FU_CHENG_INTELLIGENCE
	case 312:
		return KnxManufacturer_NEXKON
	case 313:
		return KnxManufacturer_FEEL_S_R_L
	case 314:
		return KnxManufacturer_NOT_ASSIGNED_314
	case 315:
		return KnxManufacturer_SHENZHEN_FANHAI_SANJIANG_ELECTRONICS_CO___LTD_
	case 316:
		return KnxManufacturer_JIUZHOU_GREEBLE
	case 317:
		return KnxManufacturer_AUMUELLER_AUMATIC_GMBH
	case 318:
		return KnxManufacturer_ETMAN_ELECTRIC
	case 319:
		return KnxManufacturer_BLACK_NOVA
	case 32:
		return KnxManufacturer_DEHN_AND_SOEHNE
	case 320:
		return KnxManufacturer_ZIDATECH_AG
	case 321:
		return KnxManufacturer_IDGS_BVBA
	case 322:
		return KnxManufacturer_DAKANIMO
	case 323:
		return KnxManufacturer_TREBOR_AUTOMATION_AB
	case 324:
		return KnxManufacturer_SATEL_SP__Z_O_O_
	case 325:
		return KnxManufacturer_RUSSOUND__INC_
	case 326:
		return KnxManufacturer_MIDEA_HEATING_AND_VENTILATING_EQUIPMENT_CO_LTD
	case 327:
		return KnxManufacturer_CONSORZIO_TERRANUOVA
	case 328:
		return KnxManufacturer_WOLF_HEIZTECHNIK_GMBH
	case 329:
		return KnxManufacturer_SONTEC
	case 33:
		return KnxManufacturer_CRABTREE
	case 330:
		return KnxManufacturer_BELCOM_CABLES_LTD_
	case 331:
		return KnxManufacturer_GUANGZHOU_SEAWIN_ELECTRICAL_TECHNOLOGIES_CO___LTD_
	case 332:
		return KnxManufacturer_ACREL
	case 333:
		return KnxManufacturer_FRANKE_AQUAROTTER_GMBH
	case 334:
		return KnxManufacturer_ORION_SYSTEMS
	case 335:
		return KnxManufacturer_SCHRACK_TECHNIK_GMBH
	case 336:
		return KnxManufacturer_INSPRID
	case 337:
		return KnxManufacturer_SUNRICHER
	case 338:
		return KnxManufacturer_MENRED_AUTOMATION_SYSTEMSHANGHAI_CO__LTD_
	case 339:
		return KnxManufacturer_AUREX
	case 34:
		return KnxManufacturer_EVOKNX
	case 340:
		return KnxManufacturer_JOSEF_BARTHELME_GMBH_AND_CO__KG
	case 341:
		return KnxManufacturer_ARCHITECTURE_NUMERIQUE
	case 342:
		return KnxManufacturer_UP_GROUP
	case 343:
		return KnxManufacturer_TEKNOS_AVINNO
	case 344:
		return KnxManufacturer_NINGBO_DOOYA_MECHANIC_AND_ELECTRONIC_TECHNOLOGY
	case 345:
		return KnxManufacturer_THERMOKON_SENSORTECHNIK_GMBH
	case 346:
		return KnxManufacturer_BELIMO_AUTOMATION_AG
	case 347:
		return KnxManufacturer_ZEHNDER_GROUP_INTERNATIONAL_AG
	case 348:
		return KnxManufacturer_SKS_KINKEL_ELEKTRONIK
	case 349:
		return KnxManufacturer_ECE_WURMITZER_GMBH
	case 350:
		return KnxManufacturer_LARS
	case 351:
		return KnxManufacturer_URC
	case 352:
		return KnxManufacturer_LIGHTCONTROL
	case 353:
		return KnxManufacturer_SHENZHEN_YM
	case 354:
		return KnxManufacturer_MEAN_WELL_ENTERPRISES_CO__LTD_
	case 355:
		return KnxManufacturer_OSIX
	case 356:
		return KnxManufacturer_AYPRO_TECHNOLOGY
	case 357:
		return KnxManufacturer_HEFEI_ECOLITE_SOFTWARE
	case 358:
		return KnxManufacturer_ENNO
	case 359:
		return KnxManufacturer_OHOSURE
	case 36:
		return KnxManufacturer_PAUL_HOCHKOEPPER
	case 360:
		return KnxManufacturer_GAREFOWL
	case 361:
		return KnxManufacturer_GEZE
	case 362:
		return KnxManufacturer_LG_ELECTRONICS_INC_
	case 363:
		return KnxManufacturer_SMC_INTERIORS
	case 364:
		return KnxManufacturer_NOT_ASSIGNED_364
	case 365:
		return KnxManufacturer_SCS_CABLE
	case 366:
		return KnxManufacturer_HOVAL
	case 367:
		return KnxManufacturer_CANST
	case 368:
		return KnxManufacturer_HANGZHOU_BERLIN
	case 369:
		return KnxManufacturer_EVN_LICHTTECHNIK
	case 37:
		return KnxManufacturer_ALTENBURGER_ELECTRONIC
	case 370:
		return KnxManufacturer_RUTEC
	case 371:
		return KnxManufacturer_FINDER
	case 372:
		return KnxManufacturer_FUJITSU_GENERAL_LIMITED
	case 373:
		return KnxManufacturer_ZF_FRIEDRICHSHAFEN_AG
	case 374:
		return KnxManufacturer_CREALED
	case 375:
		return KnxManufacturer_MILES_MAGIC_AUTOMATION_PRIVATE_LIMITED
	case 376:
		return KnxManufacturer_EPlus
	case 377:
		return KnxManufacturer_ITALCOND
	case 378:
		return KnxManufacturer_SATION
	case 379:
		return KnxManufacturer_NEWBEST
	case 380:
		return KnxManufacturer_GDS_DIGITAL_SYSTEMS
	case 381:
		return KnxManufacturer_IDDERO
	case 382:
		return KnxManufacturer_MBNLED
	case 383:
		return KnxManufacturer_VITRUM
	case 384:
		return KnxManufacturer_EKEY_BIOMETRIC_SYSTEMS_GMBH
	case 385:
		return KnxManufacturer_AMC
	case 386:
		return KnxManufacturer_TRILUX_GMBH_AND_CO__KG
	case 387:
		return KnxManufacturer_WEXCEDO
	case 388:
		return KnxManufacturer_VEMER_SPA
	case 389:
		return KnxManufacturer_ALEXANDER_BUERKLE_GMBH_AND_CO_KG
	case 390:
		return KnxManufacturer_CITRON
	case 391:
		return KnxManufacturer_SHENZHEN_HEGUANG
	case 392:
		return KnxManufacturer_NOT_ASSIGNED_392
	case 393:
		return KnxManufacturer_TRANE_B_V_B_A
	case 394:
		return KnxManufacturer_CAREL
	case 395:
		return KnxManufacturer_PROLITE_CONTROLS
	case 396:
		return KnxManufacturer_BOSMER
	case 397:
		return KnxManufacturer_EUCHIPS
	case 398:
		return KnxManufacturer_CONNECT_THINKA_CONNECT
	case 399:
		return KnxManufacturer_PEAKNX_A_DOGAWIST_COMPANY
	case 4:
		return KnxManufacturer_ALBRECHT_JUNG
	case 400:
		return KnxManufacturer_ACEMATIC
	case 401:
		return KnxManufacturer_ELAUSYS
	case 402:
		return KnxManufacturer_ITK_ENGINEERING_AG
	case 403:
		return KnxManufacturer_INTEGRA_METERING_AG
	case 404:
		return KnxManufacturer_FMS_HOSPITALITY_PTE_LTD
	case 405:
		return KnxManufacturer_NUVO
	case 406:
		return KnxManufacturer_U__LUX_GMBH
	case 407:
		return KnxManufacturer_BRUMBERG_LEUCHTEN
	case 408:
		return KnxManufacturer_LIME
	case 409:
		return KnxManufacturer_GREAT_EMPIRE_INTERNATIONAL_GROUP_CO___LTD_
	case 41:
		return KnxManufacturer_GRAESSLIN
	case 410:
		return KnxManufacturer_KAVOSHPISHRO_ASIA
	case 411:
		return KnxManufacturer_V2_SPA
	case 412:
		return KnxManufacturer_JOHNSON_CONTROLS
	case 413:
		return KnxManufacturer_ARKUD
	case 414:
		return KnxManufacturer_IRIDIUM_LTD_
	case 415:
		return KnxManufacturer_BSMART
	case 416:
		return KnxManufacturer_BAB_TECHNOLOGIE_GMBH
	case 417:
		return KnxManufacturer_NICE_SPA
	case 418:
		return KnxManufacturer_REDFISH_GROUP_PTY_LTD
	case 419:
		return KnxManufacturer_SABIANA_SPA
	case 42:
		return KnxManufacturer_SIMON_42
	case 420:
		return KnxManufacturer_UBEE_INTERACTIVE_EUROPE
	case 421:
		return KnxManufacturer_REXEL
	case 422:
		return KnxManufacturer_GES_TEKNIK_A_S_
	case 423:
		return KnxManufacturer_AVE_S_P_A_
	case 424:
		return KnxManufacturer_ZHUHAI_LTECH_TECHNOLOGY_CO___LTD_
	case 425:
		return KnxManufacturer_ARCOM
	case 426:
		return KnxManufacturer_VIA_TECHNOLOGIES__INC_
	case 427:
		return KnxManufacturer_FEELSMART_
	case 428:
		return KnxManufacturer_SUPCON
	case 429:
		return KnxManufacturer_MANIC
	case 430:
		return KnxManufacturer_TDE_GMBH
	case 431:
		return KnxManufacturer_NANJING_SHUFAN_INFORMATION_TECHNOLOGY_CO__LTD_
	case 432:
		return KnxManufacturer_EWTECH
	case 433:
		return KnxManufacturer_KLUGER_AUTOMATION_GMBH
	case 434:
		return KnxManufacturer_JOONGANG_CONTROL
	case 435:
		return KnxManufacturer_GREENCONTROLS_TECHNOLOGY_SDN__BHD_
	case 436:
		return KnxManufacturer_IME_S_P_A_
	case 437:
		return KnxManufacturer_SICHUAN_HAODING
	case 438:
		return KnxManufacturer_MINDJAGA_LTD_
	case 439:
		return KnxManufacturer_RUILI_SMART_CONTROL
	case 43954:
		return KnxManufacturer_ABB___RESERVED
	case 43959:
		return KnxManufacturer_BUSCH_JAEGER_ELEKTRO___RESERVED
	case 44:
		return KnxManufacturer_VIMAR
	case 440:
		return KnxManufacturer_CODESYS_GMBH
	case 441:
		return KnxManufacturer_MOORGEN_DEUTSCHLAND_GMBH
	case 442:
		return KnxManufacturer_CULLMANN_TECH
	case 443:
		return KnxManufacturer_MERCK_WINDOW_TECHNOLOGIES_B_V_
	case 444:
		return KnxManufacturer_ABEGO
	case 445:
		return KnxManufacturer_MYGEKKO
	case 446:
		return KnxManufacturer_ERGO3_SARL
	case 447:
		return KnxManufacturer_STMICROELECTRONICS_INTERNATIONAL_N_V_
	case 448:
		return KnxManufacturer_CJC_SYSTEMS
	case 449:
		return KnxManufacturer_SUDOKU
	case 45:
		return KnxManufacturer_MOELLER_GEBAEUDEAUTOMATION_KG
	case 451:
		return KnxManufacturer_AZ_E_LITE_PTE_LTD
	case 452:
		return KnxManufacturer_ARLIGHT
	case 453:
		return KnxManufacturer_GRUENBECK_WASSERAUFBEREITUNG_GMBH
	case 454:
		return KnxManufacturer_MODULE_ELECTRONIC
	case 455:
		return KnxManufacturer_KOPLAT
	case 456:
		return KnxManufacturer_GUANGZHOU_LETOUR_LIFE_TECHNOLOGY_CO___LTD
	case 457:
		return KnxManufacturer_ILEVIA
	case 458:
		return KnxManufacturer_LN_SYSTEMTEQ
	case 459:
		return KnxManufacturer_HISENSE_SMARTHOME
	case 46:
		return KnxManufacturer_ELTAKO
	case 460:
		return KnxManufacturer_FLINK_AUTOMATION_SYSTEM
	case 461:
		return KnxManufacturer_XXTER_BV
	case 462:
		return KnxManufacturer_LYNXUS_TECHNOLOGY
	case 463:
		return KnxManufacturer_ROBOT_S_A_
	case 464:
		return KnxManufacturer_SHENZHEN_ATTE_SMART_LIFE_CO__LTD_
	case 465:
		return KnxManufacturer_NOBLESSE
	case 466:
		return KnxManufacturer_ADVANCED_DEVICES
	case 467:
		return KnxManufacturer_ATRINA_BUILDING_AUTOMATION_CO__LTD
	case 468:
		return KnxManufacturer_GUANGDONG_DAMING_LAFFEY_ELECTRIC_CO___LTD_
	case 469:
		return KnxManufacturer_WESTERSTRAND_URFABRIK_AB
	case 470:
		return KnxManufacturer_CONTROL4_CORPORATE
	case 471:
		return KnxManufacturer_ONTROL
	case 472:
		return KnxManufacturer_STARNET
	case 473:
		return KnxManufacturer_BETA_CAVI
	case 474:
		return KnxManufacturer_EASEMORE
	case 475:
		return KnxManufacturer_VIVALDI_SRL
	case 476:
		return KnxManufacturer_GREE_ELECTRIC_APPLIANCES_INC__OF_ZHUHAI
	case 477:
		return KnxManufacturer_HWISCON
	case 478:
		return KnxManufacturer_SHANGHAI_ELECON_INTELLIGENT_TECHNOLOGY_CO___LTD_
	case 479:
		return KnxManufacturer_KAMPMANN
	case 480:
		return KnxManufacturer_IMPOLUX_GMBH_LEDIMAX
	case 481:
		return KnxManufacturer_EVAUX
	case 482:
		return KnxManufacturer_WEBRO_CABLES_AND_CONNECTORS_LIMITED
	case 483:
		return KnxManufacturer_SHANGHAI_E_TECH_SOLUTION
	case 484:
		return KnxManufacturer_GUANGZHOU_HOKO_ELECTRIC_CO__LTD_
	case 485:
		return KnxManufacturer_LAMMIN_HIGH_TECH_CO__LTD
	case 486:
		return KnxManufacturer_SHENZHEN_MERRYTEK_TECHNOLOGY_CO___LTD
	case 487:
		return KnxManufacturer_I_LUXUS
	case 488:
		return KnxManufacturer_ELMOS_SEMICONDUCTOR_AG
	case 489:
		return KnxManufacturer_EMCOM_TECHNOLOGY_INC
	case 49:
		return KnxManufacturer_BOSCH_SIEMENS_HAUSHALTSGERAETE
	case 490:
		return KnxManufacturer_PROJECT_INNOVATIONS_GMBH
	case 491:
		return KnxManufacturer_ITC
	case 492:
		return KnxManufacturer_ABB_LV_INSTALLATION_MATERIALS_COMPANY_LTD__BEIJING
	case 493:
		return KnxManufacturer_MAICO
	case 495:
		return KnxManufacturer_ELAN_SRL
	case 496:
		return KnxManufacturer_MINHHA_TECHNOLOGY_CO__LTD
	case 497:
		return KnxManufacturer_ZHEJIANG_TIANJIE_INDUSTRIAL_CORP_
	case 498:
		return KnxManufacturer_IAUTOMATION_PTY_LIMITED
	case 499:
		return KnxManufacturer_EXTRON
	case 5:
		return KnxManufacturer_BTICINO
	case 500:
		return KnxManufacturer_FREEDOMPRO
	case 501:
		return KnxManufacturer_1HOME
	case 502:
		return KnxManufacturer_EOS_SAUNATECHNIK_GMBH
	case 503:
		return KnxManufacturer_KUSATEK_GMBH
	case 504:
		return KnxManufacturer_EISBAER_SCADA
	case 505:
		return KnxManufacturer_AUTOMATISMI_BENINCA_S_P_A_
	case 506:
		return KnxManufacturer_BLENDOM
	case 507:
		return KnxManufacturer_MADEL_AIR_TECHNICAL_DIFFUSION
	case 508:
		return KnxManufacturer_NIKO
	case 509:
		return KnxManufacturer_BOSCH_REXROTH_AG
	case 512:
		return KnxManufacturer_CANDM_PRODUCTS
	case 513:
		return KnxManufacturer_HOERMANN_KG_VERKAUFSGESELLSCHAFT
	case 514:
		return KnxManufacturer_SHANGHAI_RAJAYASA_CO__LTD
	case 515:
		return KnxManufacturer_SUZUKI
	case 516:
		return KnxManufacturer_SILENT_GLISS_INTERNATIONAL_LTD_
	case 517:
		return KnxManufacturer_BEE_CONTROLS_ADGSC_GROUP
	case 518:
		return KnxManufacturer_XDTECGMBH
	case 519:
		return KnxManufacturer_OSRAM
	case 52:
		return KnxManufacturer_RITTO_GMBHANDCO_KG
	case 520:
		return KnxManufacturer_LEBENOR
	case 521:
		return KnxManufacturer_AUTOMANENG
	case 522:
		return KnxManufacturer_HONEYWELL_AUTOMATION_SOLUTION_CONTROLCHINA
	case 523:
		return KnxManufacturer_HANGZHOU_BINTHEN_INTELLIGENCE_TECHNOLOGY_CO__LTD
	case 524:
		return KnxManufacturer_ETA_HEIZTECHNIK
	case 525:
		return KnxManufacturer_DIVUS_GMBH
	case 526:
		return KnxManufacturer_NANJING_TAIJIESAI_INTELLIGENT_TECHNOLOGY_CO__LTD_
	case 527:
		return KnxManufacturer_LUNATONE
	case 528:
		return KnxManufacturer_ZHEJIANG_SCTECH_BUILDING_INTELLIGENT
	case 529:
		return KnxManufacturer_FOSHAN_QITE_TECHNOLOGY_CO___LTD_
	case 53:
		return KnxManufacturer_POWER_CONTROLS
	case 530:
		return KnxManufacturer_NOKE
	case 531:
		return KnxManufacturer_LANDCOM
	case 532:
		return KnxManufacturer_STORK_AS
	case 533:
		return KnxManufacturer_HANGZHOU_SHENDU_TECHNOLOGY_CO___LTD_
	case 534:
		return KnxManufacturer_COOLAUTOMATION
	case 535:
		return KnxManufacturer_APRSTERN
	case 536:
		return KnxManufacturer_SONNEN
	case 537:
		return KnxManufacturer_DNAKE
	case 538:
		return KnxManufacturer_NEUBERGER_GEBAEUDEAUTOMATION_GMBH
	case 539:
		return KnxManufacturer_STILIGER
	case 540:
		return KnxManufacturer_BERGHOF_AUTOMATION_GMBH
	case 541:
		return KnxManufacturer_TOTAL_AUTOMATION_AND_CONTROLS_GMBH
	case 542:
		return KnxManufacturer_DOVIT
	case 543:
		return KnxManufacturer_INSTALIGHTING_GMBH
	case 544:
		return KnxManufacturer_UNI_TEC
	case 545:
		return KnxManufacturer_CASATUNES
	case 546:
		return KnxManufacturer_EMT
	case 547:
		return KnxManufacturer_SENFFICIENT
	case 548:
		return KnxManufacturer_AUROLITE_ELECTRICAL_PANYU_GUANGZHOU_LIMITED
	case 549:
		return KnxManufacturer_ABB_XIAMEN_SMART_TECHNOLOGY_CO___LTD_
	case 55:
		return KnxManufacturer_ZUMTOBEL
	case 550:
		return KnxManufacturer_SAMSON_ELECTRIC_WIRE
	case 551:
		return KnxManufacturer_T_TOUCHING
	case 552:
		return KnxManufacturer_CORE_SMART_HOME
	case 553:
		return KnxManufacturer_GREENCONNECT_SOLUTIONS_SA
	case 554:
		return KnxManufacturer_ELETTRONICA_CONDUTTORI
	case 555:
		return KnxManufacturer_MKFC
	case 556:
		return KnxManufacturer_AUTOMATIONPlus
	case 557:
		return KnxManufacturer_BLUE_AND_RED
	case 558:
		return KnxManufacturer_FROGBLUE
	case 559:
		return KnxManufacturer_SAVESOR
	case 560:
		return KnxManufacturer_APP_TECH
	case 561:
		return KnxManufacturer_SENSORTEC_AG
	case 562:
		return KnxManufacturer_NYSA_TECHNOLOGY_AND_SOLUTIONS
	case 563:
		return KnxManufacturer_FARADITE
	case 564:
		return KnxManufacturer_OPTIMUS
	case 565:
		return KnxManufacturer_KTS_S_R_L_
	case 566:
		return KnxManufacturer_RAMCRO_SPA
	case 567:
		return KnxManufacturer_WUHAN_WISECREATE_UNIVERSE_TECHNOLOGY_CO___LTD
	case 568:
		return KnxManufacturer_BEMI_SMART_HOME_LTD
	case 569:
		return KnxManufacturer_ARDOMUS
	case 57:
		return KnxManufacturer_PHOENIX_CONTACT
	case 570:
		return KnxManufacturer_CHANGXING
	case 571:
		return KnxManufacturer_E_CONTROLS
	case 572:
		return KnxManufacturer_AIB_TECHNOLOGY
	case 573:
		return KnxManufacturer_NVC
	case 574:
		return KnxManufacturer_KBOX
	case 575:
		return KnxManufacturer_CNS
	case 576:
		return KnxManufacturer_TYBA
	case 577:
		return KnxManufacturer_ATREL
	case 578:
		return KnxManufacturer_SIMON_ELECTRIC_CHINA_CO___LTD
	case 579:
		return KnxManufacturer_KORDZ_GROUP
	case 580:
		return KnxManufacturer_ND_ELECTRIC
	case 581:
		return KnxManufacturer_CONTROLIUM
	case 582:
		return KnxManufacturer_FAMO_GMBH_AND_CO__KG
	case 583:
		return KnxManufacturer_CDN_SMART
	case 584:
		return KnxManufacturer_HESTON
	case 585:
		return KnxManufacturer_ESLA_CONEXIONES_S_L_
	case 586:
		return KnxManufacturer_WEISHAUPT
	case 587:
		return KnxManufacturer_ASTRUM_TECHNOLOGY
	case 588:
		return KnxManufacturer_WUERTH_ELEKTRONIK_STELVIO_KONTEK_S_P_A_
	case 589:
		return KnxManufacturer_NANOTECO_CORPORATION
	case 590:
		return KnxManufacturer_NIETIAN
	case 591:
		return KnxManufacturer_SUMSIR
	case 592:
		return KnxManufacturer_ORBIS_TECNOLOGIA_ELECTRICA_SA
	case 6:
		return KnxManufacturer_BERKER
	case 61:
		return KnxManufacturer_WAGO_KONTAKTTECHNIK
	case 62:
		return KnxManufacturer_KNXPRESSO
	case 66:
		return KnxManufacturer_WIELAND_ELECTRIC
	case 67:
		return KnxManufacturer_HERMANN_KLEINHUIS
	case 69:
		return KnxManufacturer_STIEBEL_ELTRON
	case 7:
		return KnxManufacturer_BUSCH_JAEGER_ELEKTRO
	case 71:
		return KnxManufacturer_TEHALIT
	case 72:
		return KnxManufacturer_THEBEN_AG
	case 73:
		return KnxManufacturer_WILHELM_RUTENBECK
	case 75:
		return KnxManufacturer_WINKHAUS
	case 76:
		return KnxManufacturer_ROBERT_BOSCH
	case 78:
		return KnxManufacturer_SOMFY
	case 8:
		return KnxManufacturer_GIRA_GIERSIEPEN
	case 80:
		return KnxManufacturer_WOERTZ
	case 81:
		return KnxManufacturer_VIESSMANN_WERKE
	case 82:
		return KnxManufacturer_IMI_HYDRONIC_ENGINEERING
	case 83:
		return KnxManufacturer_JOH__VAILLANT
	case 85:
		return KnxManufacturer_AMP_DEUTSCHLAND
	case 89:
		return KnxManufacturer_BOSCH_THERMOTECHNIK_GMBH
	case 9:
		return KnxManufacturer_HAGER_ELECTRO
	case 90:
		return KnxManufacturer_SEF___ECOTEC
	case 92:
		return KnxManufacturer_DORMA_GMBH_Plus_CO__KG
	case 93:
		return KnxManufacturer_WINDOWMASTER_AS
	case 94:
		return KnxManufacturer_WALTHER_WERKE
	case 95:
		return KnxManufacturer_ORAS
	case 97:
		return KnxManufacturer_DAETWYLER
	case 98:
		return KnxManufacturer_ELECTRAK
	case 99:
		return KnxManufacturer_TECHEM
	}
	return 0
}

func CastKnxManufacturer(structType interface{}) KnxManufacturer {
	castFunc := func(typ interface{}) KnxManufacturer {
		if sKnxManufacturer, ok := typ.(KnxManufacturer); ok {
			return sKnxManufacturer
		}
		return 0
	}
	return castFunc(structType)
}

func (m KnxManufacturer) LengthInBits() uint16 {
	return 16
}

func (m KnxManufacturer) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func KnxManufacturerParse(io *utils.ReadBuffer) (KnxManufacturer, error) {
	val, err := io.ReadUint16(16)
	if err != nil {
		return 0, nil
	}
	return KnxManufacturerValueOf(val), nil
}

func (e KnxManufacturer) Serialize(io utils.WriteBuffer) error {
	err := io.WriteUint16(16, uint16(e))
	return err
}

func (e KnxManufacturer) String() string {
	switch e {
	case KnxManufacturer_SIEMENS:
		return "SIEMENS"
	case KnxManufacturer_INSTA_GMBH:
		return "INSTA_GMBH"
	case KnxManufacturer_SCHNEIDER_ELECTRIC_INDUSTRIES_SAS:
		return "SCHNEIDER_ELECTRIC_INDUSTRIES_SAS"
	case KnxManufacturer_WHD_WILHELM_HUBER_Plus_SOEHNE:
		return "WHD_WILHELM_HUBER_Plus_SOEHNE"
	case KnxManufacturer_BISCHOFF_ELEKTRONIK:
		return "BISCHOFF_ELEKTRONIK"
	case KnxManufacturer_JEPAZ:
		return "JEPAZ"
	case KnxManufacturer_RTS_AUTOMATION:
		return "RTS_AUTOMATION"
	case KnxManufacturer_EIBMARKT_GMBH:
		return "EIBMARKT_GMBH"
	case KnxManufacturer_WAREMA_RENKHOFF_SE:
		return "WAREMA_RENKHOFF_SE"
	case KnxManufacturer_EELECTRON:
		return "EELECTRON"
	case KnxManufacturer_BELDEN_WIRE_AND_CABLE_B_V_:
		return "BELDEN_WIRE_AND_CABLE_B_V_"
	case KnxManufacturer_LEGRAND_APPAREILLAGE_ELECTRIQUE:
		return "LEGRAND_APPAREILLAGE_ELECTRIQUE"
	case KnxManufacturer_BECKER_ANTRIEBE_GMBH:
		return "BECKER_ANTRIEBE_GMBH"
	case KnxManufacturer_J_STEHLEPlusSOEHNE_GMBH:
		return "J_STEHLEPlusSOEHNE_GMBH"
	case KnxManufacturer_AGFEO:
		return "AGFEO"
	case KnxManufacturer_ZENNIO:
		return "ZENNIO"
	case KnxManufacturer_TAPKO_TECHNOLOGIES:
		return "TAPKO_TECHNOLOGIES"
	case KnxManufacturer_HDL:
		return "HDL"
	case KnxManufacturer_UPONOR:
		return "UPONOR"
	case KnxManufacturer_SE_LIGHTMANAGEMENT_AG:
		return "SE_LIGHTMANAGEMENT_AG"
	case KnxManufacturer_ARCUS_EDS:
		return "ARCUS_EDS"
	case KnxManufacturer_INTESIS:
		return "INTESIS"
	case KnxManufacturer_MERTEN:
		return "MERTEN"
	case KnxManufacturer_HERHOLDT_CONTROLS_SRL:
		return "HERHOLDT_CONTROLS_SRL"
	case KnxManufacturer_NIKO_ZUBLIN:
		return "NIKO_ZUBLIN"
	case KnxManufacturer_DURABLE_TECHNOLOGIES:
		return "DURABLE_TECHNOLOGIES"
	case KnxManufacturer_INNOTEAM:
		return "INNOTEAM"
	case KnxManufacturer_ISE_GMBH:
		return "ISE_GMBH"
	case KnxManufacturer_TEAM_FOR_TRONICS:
		return "TEAM_FOR_TRONICS"
	case KnxManufacturer_CIAT:
		return "CIAT"
	case KnxManufacturer_REMEHA_BV:
		return "REMEHA_BV"
	case KnxManufacturer_ESYLUX:
		return "ESYLUX"
	case KnxManufacturer_BASALTE:
		return "BASALTE"
	case KnxManufacturer_VESTAMATIC:
		return "VESTAMATIC"
	case KnxManufacturer_MDT_TECHNOLOGIES:
		return "MDT_TECHNOLOGIES"
	case KnxManufacturer_WARENDORFER_KUECHEN_GMBH:
		return "WARENDORFER_KUECHEN_GMBH"
	case KnxManufacturer_VIDEO_STAR:
		return "VIDEO_STAR"
	case KnxManufacturer_SITEK:
		return "SITEK"
	case KnxManufacturer_CONTROLTRONIC:
		return "CONTROLTRONIC"
	case KnxManufacturer_FUNCTION_TECHNOLOGY:
		return "FUNCTION_TECHNOLOGY"
	case KnxManufacturer_AMX:
		return "AMX"
	case KnxManufacturer_ELDAT:
		return "ELDAT"
	case KnxManufacturer_PANASONIC:
		return "PANASONIC"
	case KnxManufacturer_ABB_SPA_SACE_DIVISION:
		return "ABB_SPA_SACE_DIVISION"
	case KnxManufacturer_PULSE_TECHNOLOGIES:
		return "PULSE_TECHNOLOGIES"
	case KnxManufacturer_CRESTRON:
		return "CRESTRON"
	case KnxManufacturer_STEINEL_PROFESSIONAL:
		return "STEINEL_PROFESSIONAL"
	case KnxManufacturer_BILTON_LED_LIGHTING:
		return "BILTON_LED_LIGHTING"
	case KnxManufacturer_DENRO_AG:
		return "DENRO_AG"
	case KnxManufacturer_GEPRO:
		return "GEPRO"
	case KnxManufacturer_PREUSSEN_AUTOMATION:
		return "PREUSSEN_AUTOMATION"
	case KnxManufacturer_ZOPPAS_INDUSTRIES:
		return "ZOPPAS_INDUSTRIES"
	case KnxManufacturer_MACTECH:
		return "MACTECH"
	case KnxManufacturer_TECHNO_TREND:
		return "TECHNO_TREND"
	case KnxManufacturer_FS_CABLES:
		return "FS_CABLES"
	case KnxManufacturer_DELTA_DORE:
		return "DELTA_DORE"
	case KnxManufacturer_EISSOUND:
		return "EISSOUND"
	case KnxManufacturer_CISCO:
		return "CISCO"
	case KnxManufacturer_DINUY:
		return "DINUY"
	case KnxManufacturer_IKNIX:
		return "IKNIX"
	case KnxManufacturer_RADEMACHER_GERAETE_ELEKTRONIK_GMBH:
		return "RADEMACHER_GERAETE_ELEKTRONIK_GMBH"
	case KnxManufacturer_EGI_ELECTROACUSTICA_GENERAL_IBERICA:
		return "EGI_ELECTROACUSTICA_GENERAL_IBERICA"
	case KnxManufacturer_BES___INGENIUM:
		return "BES___INGENIUM"
	case KnxManufacturer_ELABNET:
		return "ELABNET"
	case KnxManufacturer_BLUMOTIX:
		return "BLUMOTIX"
	case KnxManufacturer_HUNTER_DOUGLAS:
		return "HUNTER_DOUGLAS"
	case KnxManufacturer_APRICUM:
		return "APRICUM"
	case KnxManufacturer_TIANSU_AUTOMATION:
		return "TIANSU_AUTOMATION"
	case KnxManufacturer_BUBENDORFF:
		return "BUBENDORFF"
	case KnxManufacturer_MBS_GMBH:
		return "MBS_GMBH"
	case KnxManufacturer_ENERTEX_BAYERN_GMBH:
		return "ENERTEX_BAYERN_GMBH"
	case KnxManufacturer_BMS:
		return "BMS"
	case KnxManufacturer_SINAPSI:
		return "SINAPSI"
	case KnxManufacturer_EMBEDDED_SYSTEMS_SIA:
		return "EMBEDDED_SYSTEMS_SIA"
	case KnxManufacturer_KNX1:
		return "KNX1"
	case KnxManufacturer_TOKKA:
		return "TOKKA"
	case KnxManufacturer_NANOSENSE:
		return "NANOSENSE"
	case KnxManufacturer_PEAR_AUTOMATION_GMBH:
		return "PEAR_AUTOMATION_GMBH"
	case KnxManufacturer_DGA:
		return "DGA"
	case KnxManufacturer_LUTRON:
		return "LUTRON"
	case KnxManufacturer_AIRZONE___ALTRA:
		return "AIRZONE___ALTRA"
	case KnxManufacturer_LITHOSS_DESIGN_SWITCHES:
		return "LITHOSS_DESIGN_SWITCHES"
	case KnxManufacturer_3ATEL:
		return "3ATEL"
	case KnxManufacturer_PHILIPS_CONTROLS:
		return "PHILIPS_CONTROLS"
	case KnxManufacturer_VELUX_AS:
		return "VELUX_AS"
	case KnxManufacturer_LOYTEC:
		return "LOYTEC"
	case KnxManufacturer_EKINEX_S_P_A_:
		return "EKINEX_S_P_A_"
	case KnxManufacturer_SIRLAN_TECHNOLOGIES:
		return "SIRLAN_TECHNOLOGIES"
	case KnxManufacturer_PROKNX_SAS:
		return "PROKNX_SAS"
	case KnxManufacturer_IT_GMBH:
		return "IT_GMBH"
	case KnxManufacturer_RENSON:
		return "RENSON"
	case KnxManufacturer_HEP_GROUP:
		return "HEP_GROUP"
	case KnxManufacturer_BALMART:
		return "BALMART"
	case KnxManufacturer_GFS_GMBH:
		return "GFS_GMBH"
	case KnxManufacturer_SCHENKER_STOREN_AG:
		return "SCHENKER_STOREN_AG"
	case KnxManufacturer_ALGODUE_ELETTRONICA_S_R_L_:
		return "ALGODUE_ELETTRONICA_S_R_L_"
	case KnxManufacturer_ABB_FRANCE:
		return "ABB_FRANCE"
	case KnxManufacturer_MAINTRONIC:
		return "MAINTRONIC"
	case KnxManufacturer_VANTAGE:
		return "VANTAGE"
	case KnxManufacturer_FORESIS:
		return "FORESIS"
	case KnxManufacturer_RESEARCH_AND_PRODUCTION_ASSOCIATION_SEM:
		return "RESEARCH_AND_PRODUCTION_ASSOCIATION_SEM"
	case KnxManufacturer_WEINZIERL_ENGINEERING_GMBH:
		return "WEINZIERL_ENGINEERING_GMBH"
	case KnxManufacturer_MOEHLENHOFF_WAERMETECHNIK_GMBH:
		return "MOEHLENHOFF_WAERMETECHNIK_GMBH"
	case KnxManufacturer_PKC_GROUP_OYJ:
		return "PKC_GROUP_OYJ"
	case KnxManufacturer_ABB:
		return "ABB"
	case KnxManufacturer_B_E_G_:
		return "B_E_G_"
	case KnxManufacturer_ELSNER_ELEKTRONIK_GMBH:
		return "ELSNER_ELEKTRONIK_GMBH"
	case KnxManufacturer_SIEMENS_BUILDING_TECHNOLOGIES_HKCHINA_LTD_:
		return "SIEMENS_BUILDING_TECHNOLOGIES_HKCHINA_LTD_"
	case KnxManufacturer_EUTRAC:
		return "EUTRAC"
	case KnxManufacturer_GUSTAV_HENSEL_GMBH_AND_CO__KG:
		return "GUSTAV_HENSEL_GMBH_AND_CO__KG"
	case KnxManufacturer_GARO_AB:
		return "GARO_AB"
	case KnxManufacturer_WALDMANN_LICHTTECHNIK:
		return "WALDMANN_LICHTTECHNIK"
	case KnxManufacturer_SCHUECO:
		return "SCHUECO"
	case KnxManufacturer_EMU:
		return "EMU"
	case KnxManufacturer_JNET_SYSTEMS_AG:
		return "JNET_SYSTEMS_AG"
	case KnxManufacturer_TOTAL_SOLUTION_GMBH:
		return "TOTAL_SOLUTION_GMBH"
	case KnxManufacturer_O_Y_L__ELECTRONICS:
		return "O_Y_L__ELECTRONICS"
	case KnxManufacturer_GALAX_SYSTEM:
		return "GALAX_SYSTEM"
	case KnxManufacturer_DISCH:
		return "DISCH"
	case KnxManufacturer_AUCOTEAM:
		return "AUCOTEAM"
	case KnxManufacturer_LUXMATE_CONTROLS:
		return "LUXMATE_CONTROLS"
	case KnxManufacturer_DANFOSS:
		return "DANFOSS"
	case KnxManufacturer_SIEDLE_AND_SOEHNE:
		return "SIEDLE_AND_SOEHNE"
	case KnxManufacturer_AST_GMBH:
		return "AST_GMBH"
	case KnxManufacturer_WILA_LEUCHTEN:
		return "WILA_LEUCHTEN"
	case KnxManufacturer_BPlusB_AUTOMATIONS__UND_STEUERUNGSTECHNIK:
		return "BPlusB_AUTOMATIONS__UND_STEUERUNGSTECHNIK"
	case KnxManufacturer_LINGG_AND_JANKE:
		return "LINGG_AND_JANKE"
	case KnxManufacturer_SAUTER:
		return "SAUTER"
	case KnxManufacturer_SIMU:
		return "SIMU"
	case KnxManufacturer_THEBEN_HTS_AG:
		return "THEBEN_HTS_AG"
	case KnxManufacturer_AMANN_GMBH:
		return "AMANN_GMBH"
	case KnxManufacturer_BERG_ENERGIEKONTROLLSYSTEME_GMBH:
		return "BERG_ENERGIEKONTROLLSYSTEME_GMBH"
	case KnxManufacturer_HUEPPE_FORM_SONNENSCHUTZSYSTEME_GMBH:
		return "HUEPPE_FORM_SONNENSCHUTZSYSTEME_GMBH"
	case KnxManufacturer_OVENTROP_KG:
		return "OVENTROP_KG"
	case KnxManufacturer_GRIESSER_AG:
		return "GRIESSER_AG"
	case KnxManufacturer_IPAS_GMBH:
		return "IPAS_GMBH"
	case KnxManufacturer_EBERLE:
		return "EBERLE"
	case KnxManufacturer_ELERO_GMBH:
		return "ELERO_GMBH"
	case KnxManufacturer_ARDAN_PRODUCTION_AND_INDUSTRIAL_CONTROLS_LTD_:
		return "ARDAN_PRODUCTION_AND_INDUSTRIAL_CONTROLS_LTD_"
	case KnxManufacturer_METEC_MESSTECHNIK_GMBH:
		return "METEC_MESSTECHNIK_GMBH"
	case KnxManufacturer_ELKA_ELEKTRONIK_GMBH:
		return "ELKA_ELEKTRONIK_GMBH"
	case KnxManufacturer_ELEKTROANLAGEN_D__NAGEL:
		return "ELEKTROANLAGEN_D__NAGEL"
	case KnxManufacturer_TRIDONIC_BAUELEMENTE_GMBH:
		return "TRIDONIC_BAUELEMENTE_GMBH"
	case KnxManufacturer_STENGLER_GESELLSCHAFT:
		return "STENGLER_GESELLSCHAFT"
	case KnxManufacturer_SCHNEIDER_ELECTRIC_MG:
		return "SCHNEIDER_ELECTRIC_MG"
	case KnxManufacturer_GEWISS:
		return "GEWISS"
	case KnxManufacturer_KNX_ASSOCIATION:
		return "KNX_ASSOCIATION"
	case KnxManufacturer_VIVO:
		return "VIVO"
	case KnxManufacturer_HUGO_MUELLER_GMBH_AND_CO_KG:
		return "HUGO_MUELLER_GMBH_AND_CO_KG"
	case KnxManufacturer_SIEMENS_HVAC:
		return "SIEMENS_HVAC"
	case KnxManufacturer_APT:
		return "APT"
	case KnxManufacturer_HIGHDOM:
		return "HIGHDOM"
	case KnxManufacturer_TOP_SERVICES:
		return "TOP_SERVICES"
	case KnxManufacturer_AMBIHOME:
		return "AMBIHOME"
	case KnxManufacturer_DATEC_ELECTRONIC_AG:
		return "DATEC_ELECTRONIC_AG"
	case KnxManufacturer_ABUS_SECURITY_CENTER:
		return "ABUS_SECURITY_CENTER"
	case KnxManufacturer_LITE_PUTER:
		return "LITE_PUTER"
	case KnxManufacturer_TANTRON_ELECTRONIC:
		return "TANTRON_ELECTRONIC"
	case KnxManufacturer_INTERRA:
		return "INTERRA"
	case KnxManufacturer_DKX_TECH:
		return "DKX_TECH"
	case KnxManufacturer_VIATRON:
		return "VIATRON"
	case KnxManufacturer_NAUTIBUS:
		return "NAUTIBUS"
	case KnxManufacturer_ON_SEMICONDUCTOR:
		return "ON_SEMICONDUCTOR"
	case KnxManufacturer_LONGCHUANG:
		return "LONGCHUANG"
	case KnxManufacturer_AIR_ON_AG:
		return "AIR_ON_AG"
	case KnxManufacturer_ALBERT_ACKERMANN:
		return "ALBERT_ACKERMANN"
	case KnxManufacturer_IB_COMPANY_GMBH:
		return "IB_COMPANY_GMBH"
	case KnxManufacturer_SATION_FACTORY:
		return "SATION_FACTORY"
	case KnxManufacturer_AGENTILO_GMBH:
		return "AGENTILO_GMBH"
	case KnxManufacturer_MAKEL_ELEKTRIK:
		return "MAKEL_ELEKTRIK"
	case KnxManufacturer_HELIOS_VENTILATOREN:
		return "HELIOS_VENTILATOREN"
	case KnxManufacturer_OTTO_SOLUTIONS_PTE_LTD:
		return "OTTO_SOLUTIONS_PTE_LTD"
	case KnxManufacturer_AIRMASTER:
		return "AIRMASTER"
	case KnxManufacturer_VALLOX_GMBH:
		return "VALLOX_GMBH"
	case KnxManufacturer_DALITEK:
		return "DALITEK"
	case KnxManufacturer_ASIN:
		return "ASIN"
	case KnxManufacturer_SCHUPA_GMBH:
		return "SCHUPA_GMBH"
	case KnxManufacturer_BRIDGES_INTELLIGENCE_TECHNOLOGY_INC_:
		return "BRIDGES_INTELLIGENCE_TECHNOLOGY_INC_"
	case KnxManufacturer_ARBONIA:
		return "ARBONIA"
	case KnxManufacturer_KERMI:
		return "KERMI"
	case KnxManufacturer_PROLUX:
		return "PROLUX"
	case KnxManufacturer_CLICHOME:
		return "CLICHOME"
	case KnxManufacturer_COMMAX:
		return "COMMAX"
	case KnxManufacturer_EAE:
		return "EAE"
	case KnxManufacturer_TENSE:
		return "TENSE"
	case KnxManufacturer_SEYOUNG_ELECTRONICS:
		return "SEYOUNG_ELECTRONICS"
	case KnxManufacturer_LIFEDOMUS:
		return "LIFEDOMUS"
	case KnxManufacturer_ABB_SCHWEIZ:
		return "ABB_SCHWEIZ"
	case KnxManufacturer_EUROTRONIC_TECHNOLOGY_GMBH:
		return "EUROTRONIC_TECHNOLOGY_GMBH"
	case KnxManufacturer_TCI:
		return "TCI"
	case KnxManufacturer_RISHUN_ELECTRONIC:
		return "RISHUN_ELECTRONIC"
	case KnxManufacturer_ZIPATO:
		return "ZIPATO"
	case KnxManufacturer_CM_SECURITY_GMBH_AND_CO_KG:
		return "CM_SECURITY_GMBH_AND_CO_KG"
	case KnxManufacturer_QING_CABLES:
		return "QING_CABLES"
	case KnxManufacturer_LABIO:
		return "LABIO"
	case KnxManufacturer_COSTER_TECNOLOGIE_ELETTRONICHE_S_P_A_:
		return "COSTER_TECNOLOGIE_ELETTRONICHE_S_P_A_"
	case KnxManufacturer_E_G_E:
		return "E_G_E"
	case KnxManufacturer_NETXAUTOMATION:
		return "NETXAUTOMATION"
	case KnxManufacturer_FELLER:
		return "FELLER"
	case KnxManufacturer_TECALOR:
		return "TECALOR"
	case KnxManufacturer_URMET_ELECTRONICS_HUIZHOU_LTD_:
		return "URMET_ELECTRONICS_HUIZHOU_LTD_"
	case KnxManufacturer_PEIYING_BUILDING_CONTROL:
		return "PEIYING_BUILDING_CONTROL"
	case KnxManufacturer_BPT_S_P_A__A_SOCIO_UNICO:
		return "BPT_S_P_A__A_SOCIO_UNICO"
	case KnxManufacturer_KANONTEC___KANONBUS:
		return "KANONTEC___KANONBUS"
	case KnxManufacturer_ISER_TECH:
		return "ISER_TECH"
	case KnxManufacturer_FINELINE:
		return "FINELINE"
	case KnxManufacturer_CP_ELECTRONICS_LTD:
		return "CP_ELECTRONICS_LTD"
	case KnxManufacturer_NIKO_SERVODAN_AS:
		return "NIKO_SERVODAN_AS"
	case KnxManufacturer_SIMON_309:
		return "SIMON_309"
	case KnxManufacturer_GLAMOX_AS:
		return "GLAMOX_AS"
	case KnxManufacturer_GM_MODULAR_PVT__LTD_:
		return "GM_MODULAR_PVT__LTD_"
	case KnxManufacturer_FU_CHENG_INTELLIGENCE:
		return "FU_CHENG_INTELLIGENCE"
	case KnxManufacturer_NEXKON:
		return "NEXKON"
	case KnxManufacturer_FEEL_S_R_L:
		return "FEEL_S_R_L"
	case KnxManufacturer_NOT_ASSIGNED_314:
		return "NOT_ASSIGNED_314"
	case KnxManufacturer_SHENZHEN_FANHAI_SANJIANG_ELECTRONICS_CO___LTD_:
		return "SHENZHEN_FANHAI_SANJIANG_ELECTRONICS_CO___LTD_"
	case KnxManufacturer_JIUZHOU_GREEBLE:
		return "JIUZHOU_GREEBLE"
	case KnxManufacturer_AUMUELLER_AUMATIC_GMBH:
		return "AUMUELLER_AUMATIC_GMBH"
	case KnxManufacturer_ETMAN_ELECTRIC:
		return "ETMAN_ELECTRIC"
	case KnxManufacturer_BLACK_NOVA:
		return "BLACK_NOVA"
	case KnxManufacturer_DEHN_AND_SOEHNE:
		return "DEHN_AND_SOEHNE"
	case KnxManufacturer_ZIDATECH_AG:
		return "ZIDATECH_AG"
	case KnxManufacturer_IDGS_BVBA:
		return "IDGS_BVBA"
	case KnxManufacturer_DAKANIMO:
		return "DAKANIMO"
	case KnxManufacturer_TREBOR_AUTOMATION_AB:
		return "TREBOR_AUTOMATION_AB"
	case KnxManufacturer_SATEL_SP__Z_O_O_:
		return "SATEL_SP__Z_O_O_"
	case KnxManufacturer_RUSSOUND__INC_:
		return "RUSSOUND__INC_"
	case KnxManufacturer_MIDEA_HEATING_AND_VENTILATING_EQUIPMENT_CO_LTD:
		return "MIDEA_HEATING_AND_VENTILATING_EQUIPMENT_CO_LTD"
	case KnxManufacturer_CONSORZIO_TERRANUOVA:
		return "CONSORZIO_TERRANUOVA"
	case KnxManufacturer_WOLF_HEIZTECHNIK_GMBH:
		return "WOLF_HEIZTECHNIK_GMBH"
	case KnxManufacturer_SONTEC:
		return "SONTEC"
	case KnxManufacturer_CRABTREE:
		return "CRABTREE"
	case KnxManufacturer_BELCOM_CABLES_LTD_:
		return "BELCOM_CABLES_LTD_"
	case KnxManufacturer_GUANGZHOU_SEAWIN_ELECTRICAL_TECHNOLOGIES_CO___LTD_:
		return "GUANGZHOU_SEAWIN_ELECTRICAL_TECHNOLOGIES_CO___LTD_"
	case KnxManufacturer_ACREL:
		return "ACREL"
	case KnxManufacturer_FRANKE_AQUAROTTER_GMBH:
		return "FRANKE_AQUAROTTER_GMBH"
	case KnxManufacturer_ORION_SYSTEMS:
		return "ORION_SYSTEMS"
	case KnxManufacturer_SCHRACK_TECHNIK_GMBH:
		return "SCHRACK_TECHNIK_GMBH"
	case KnxManufacturer_INSPRID:
		return "INSPRID"
	case KnxManufacturer_SUNRICHER:
		return "SUNRICHER"
	case KnxManufacturer_MENRED_AUTOMATION_SYSTEMSHANGHAI_CO__LTD_:
		return "MENRED_AUTOMATION_SYSTEMSHANGHAI_CO__LTD_"
	case KnxManufacturer_AUREX:
		return "AUREX"
	case KnxManufacturer_EVOKNX:
		return "EVOKNX"
	case KnxManufacturer_JOSEF_BARTHELME_GMBH_AND_CO__KG:
		return "JOSEF_BARTHELME_GMBH_AND_CO__KG"
	case KnxManufacturer_ARCHITECTURE_NUMERIQUE:
		return "ARCHITECTURE_NUMERIQUE"
	case KnxManufacturer_UP_GROUP:
		return "UP_GROUP"
	case KnxManufacturer_TEKNOS_AVINNO:
		return "TEKNOS_AVINNO"
	case KnxManufacturer_NINGBO_DOOYA_MECHANIC_AND_ELECTRONIC_TECHNOLOGY:
		return "NINGBO_DOOYA_MECHANIC_AND_ELECTRONIC_TECHNOLOGY"
	case KnxManufacturer_THERMOKON_SENSORTECHNIK_GMBH:
		return "THERMOKON_SENSORTECHNIK_GMBH"
	case KnxManufacturer_BELIMO_AUTOMATION_AG:
		return "BELIMO_AUTOMATION_AG"
	case KnxManufacturer_ZEHNDER_GROUP_INTERNATIONAL_AG:
		return "ZEHNDER_GROUP_INTERNATIONAL_AG"
	case KnxManufacturer_SKS_KINKEL_ELEKTRONIK:
		return "SKS_KINKEL_ELEKTRONIK"
	case KnxManufacturer_ECE_WURMITZER_GMBH:
		return "ECE_WURMITZER_GMBH"
	case KnxManufacturer_LARS:
		return "LARS"
	case KnxManufacturer_URC:
		return "URC"
	case KnxManufacturer_LIGHTCONTROL:
		return "LIGHTCONTROL"
	case KnxManufacturer_SHENZHEN_YM:
		return "SHENZHEN_YM"
	case KnxManufacturer_MEAN_WELL_ENTERPRISES_CO__LTD_:
		return "MEAN_WELL_ENTERPRISES_CO__LTD_"
	case KnxManufacturer_OSIX:
		return "OSIX"
	case KnxManufacturer_AYPRO_TECHNOLOGY:
		return "AYPRO_TECHNOLOGY"
	case KnxManufacturer_HEFEI_ECOLITE_SOFTWARE:
		return "HEFEI_ECOLITE_SOFTWARE"
	case KnxManufacturer_ENNO:
		return "ENNO"
	case KnxManufacturer_OHOSURE:
		return "OHOSURE"
	case KnxManufacturer_PAUL_HOCHKOEPPER:
		return "PAUL_HOCHKOEPPER"
	case KnxManufacturer_GAREFOWL:
		return "GAREFOWL"
	case KnxManufacturer_GEZE:
		return "GEZE"
	case KnxManufacturer_LG_ELECTRONICS_INC_:
		return "LG_ELECTRONICS_INC_"
	case KnxManufacturer_SMC_INTERIORS:
		return "SMC_INTERIORS"
	case KnxManufacturer_NOT_ASSIGNED_364:
		return "NOT_ASSIGNED_364"
	case KnxManufacturer_SCS_CABLE:
		return "SCS_CABLE"
	case KnxManufacturer_HOVAL:
		return "HOVAL"
	case KnxManufacturer_CANST:
		return "CANST"
	case KnxManufacturer_HANGZHOU_BERLIN:
		return "HANGZHOU_BERLIN"
	case KnxManufacturer_EVN_LICHTTECHNIK:
		return "EVN_LICHTTECHNIK"
	case KnxManufacturer_ALTENBURGER_ELECTRONIC:
		return "ALTENBURGER_ELECTRONIC"
	case KnxManufacturer_RUTEC:
		return "RUTEC"
	case KnxManufacturer_FINDER:
		return "FINDER"
	case KnxManufacturer_FUJITSU_GENERAL_LIMITED:
		return "FUJITSU_GENERAL_LIMITED"
	case KnxManufacturer_ZF_FRIEDRICHSHAFEN_AG:
		return "ZF_FRIEDRICHSHAFEN_AG"
	case KnxManufacturer_CREALED:
		return "CREALED"
	case KnxManufacturer_MILES_MAGIC_AUTOMATION_PRIVATE_LIMITED:
		return "MILES_MAGIC_AUTOMATION_PRIVATE_LIMITED"
	case KnxManufacturer_EPlus:
		return "EPlus"
	case KnxManufacturer_ITALCOND:
		return "ITALCOND"
	case KnxManufacturer_SATION:
		return "SATION"
	case KnxManufacturer_NEWBEST:
		return "NEWBEST"
	case KnxManufacturer_GDS_DIGITAL_SYSTEMS:
		return "GDS_DIGITAL_SYSTEMS"
	case KnxManufacturer_IDDERO:
		return "IDDERO"
	case KnxManufacturer_MBNLED:
		return "MBNLED"
	case KnxManufacturer_VITRUM:
		return "VITRUM"
	case KnxManufacturer_EKEY_BIOMETRIC_SYSTEMS_GMBH:
		return "EKEY_BIOMETRIC_SYSTEMS_GMBH"
	case KnxManufacturer_AMC:
		return "AMC"
	case KnxManufacturer_TRILUX_GMBH_AND_CO__KG:
		return "TRILUX_GMBH_AND_CO__KG"
	case KnxManufacturer_WEXCEDO:
		return "WEXCEDO"
	case KnxManufacturer_VEMER_SPA:
		return "VEMER_SPA"
	case KnxManufacturer_ALEXANDER_BUERKLE_GMBH_AND_CO_KG:
		return "ALEXANDER_BUERKLE_GMBH_AND_CO_KG"
	case KnxManufacturer_CITRON:
		return "CITRON"
	case KnxManufacturer_SHENZHEN_HEGUANG:
		return "SHENZHEN_HEGUANG"
	case KnxManufacturer_NOT_ASSIGNED_392:
		return "NOT_ASSIGNED_392"
	case KnxManufacturer_TRANE_B_V_B_A:
		return "TRANE_B_V_B_A"
	case KnxManufacturer_CAREL:
		return "CAREL"
	case KnxManufacturer_PROLITE_CONTROLS:
		return "PROLITE_CONTROLS"
	case KnxManufacturer_BOSMER:
		return "BOSMER"
	case KnxManufacturer_EUCHIPS:
		return "EUCHIPS"
	case KnxManufacturer_CONNECT_THINKA_CONNECT:
		return "CONNECT_THINKA_CONNECT"
	case KnxManufacturer_PEAKNX_A_DOGAWIST_COMPANY:
		return "PEAKNX_A_DOGAWIST_COMPANY"
	case KnxManufacturer_ALBRECHT_JUNG:
		return "ALBRECHT_JUNG"
	case KnxManufacturer_ACEMATIC:
		return "ACEMATIC"
	case KnxManufacturer_ELAUSYS:
		return "ELAUSYS"
	case KnxManufacturer_ITK_ENGINEERING_AG:
		return "ITK_ENGINEERING_AG"
	case KnxManufacturer_INTEGRA_METERING_AG:
		return "INTEGRA_METERING_AG"
	case KnxManufacturer_FMS_HOSPITALITY_PTE_LTD:
		return "FMS_HOSPITALITY_PTE_LTD"
	case KnxManufacturer_NUVO:
		return "NUVO"
	case KnxManufacturer_U__LUX_GMBH:
		return "U__LUX_GMBH"
	case KnxManufacturer_BRUMBERG_LEUCHTEN:
		return "BRUMBERG_LEUCHTEN"
	case KnxManufacturer_LIME:
		return "LIME"
	case KnxManufacturer_GREAT_EMPIRE_INTERNATIONAL_GROUP_CO___LTD_:
		return "GREAT_EMPIRE_INTERNATIONAL_GROUP_CO___LTD_"
	case KnxManufacturer_GRAESSLIN:
		return "GRAESSLIN"
	case KnxManufacturer_KAVOSHPISHRO_ASIA:
		return "KAVOSHPISHRO_ASIA"
	case KnxManufacturer_V2_SPA:
		return "V2_SPA"
	case KnxManufacturer_JOHNSON_CONTROLS:
		return "JOHNSON_CONTROLS"
	case KnxManufacturer_ARKUD:
		return "ARKUD"
	case KnxManufacturer_IRIDIUM_LTD_:
		return "IRIDIUM_LTD_"
	case KnxManufacturer_BSMART:
		return "BSMART"
	case KnxManufacturer_BAB_TECHNOLOGIE_GMBH:
		return "BAB_TECHNOLOGIE_GMBH"
	case KnxManufacturer_NICE_SPA:
		return "NICE_SPA"
	case KnxManufacturer_REDFISH_GROUP_PTY_LTD:
		return "REDFISH_GROUP_PTY_LTD"
	case KnxManufacturer_SABIANA_SPA:
		return "SABIANA_SPA"
	case KnxManufacturer_SIMON_42:
		return "SIMON_42"
	case KnxManufacturer_UBEE_INTERACTIVE_EUROPE:
		return "UBEE_INTERACTIVE_EUROPE"
	case KnxManufacturer_REXEL:
		return "REXEL"
	case KnxManufacturer_GES_TEKNIK_A_S_:
		return "GES_TEKNIK_A_S_"
	case KnxManufacturer_AVE_S_P_A_:
		return "AVE_S_P_A_"
	case KnxManufacturer_ZHUHAI_LTECH_TECHNOLOGY_CO___LTD_:
		return "ZHUHAI_LTECH_TECHNOLOGY_CO___LTD_"
	case KnxManufacturer_ARCOM:
		return "ARCOM"
	case KnxManufacturer_VIA_TECHNOLOGIES__INC_:
		return "VIA_TECHNOLOGIES__INC_"
	case KnxManufacturer_FEELSMART_:
		return "FEELSMART_"
	case KnxManufacturer_SUPCON:
		return "SUPCON"
	case KnxManufacturer_MANIC:
		return "MANIC"
	case KnxManufacturer_TDE_GMBH:
		return "TDE_GMBH"
	case KnxManufacturer_NANJING_SHUFAN_INFORMATION_TECHNOLOGY_CO__LTD_:
		return "NANJING_SHUFAN_INFORMATION_TECHNOLOGY_CO__LTD_"
	case KnxManufacturer_EWTECH:
		return "EWTECH"
	case KnxManufacturer_KLUGER_AUTOMATION_GMBH:
		return "KLUGER_AUTOMATION_GMBH"
	case KnxManufacturer_JOONGANG_CONTROL:
		return "JOONGANG_CONTROL"
	case KnxManufacturer_GREENCONTROLS_TECHNOLOGY_SDN__BHD_:
		return "GREENCONTROLS_TECHNOLOGY_SDN__BHD_"
	case KnxManufacturer_IME_S_P_A_:
		return "IME_S_P_A_"
	case KnxManufacturer_SICHUAN_HAODING:
		return "SICHUAN_HAODING"
	case KnxManufacturer_MINDJAGA_LTD_:
		return "MINDJAGA_LTD_"
	case KnxManufacturer_RUILI_SMART_CONTROL:
		return "RUILI_SMART_CONTROL"
	case KnxManufacturer_ABB___RESERVED:
		return "ABB___RESERVED"
	case KnxManufacturer_BUSCH_JAEGER_ELEKTRO___RESERVED:
		return "BUSCH_JAEGER_ELEKTRO___RESERVED"
	case KnxManufacturer_VIMAR:
		return "VIMAR"
	case KnxManufacturer_CODESYS_GMBH:
		return "CODESYS_GMBH"
	case KnxManufacturer_MOORGEN_DEUTSCHLAND_GMBH:
		return "MOORGEN_DEUTSCHLAND_GMBH"
	case KnxManufacturer_CULLMANN_TECH:
		return "CULLMANN_TECH"
	case KnxManufacturer_MERCK_WINDOW_TECHNOLOGIES_B_V_:
		return "MERCK_WINDOW_TECHNOLOGIES_B_V_"
	case KnxManufacturer_ABEGO:
		return "ABEGO"
	case KnxManufacturer_MYGEKKO:
		return "MYGEKKO"
	case KnxManufacturer_ERGO3_SARL:
		return "ERGO3_SARL"
	case KnxManufacturer_STMICROELECTRONICS_INTERNATIONAL_N_V_:
		return "STMICROELECTRONICS_INTERNATIONAL_N_V_"
	case KnxManufacturer_CJC_SYSTEMS:
		return "CJC_SYSTEMS"
	case KnxManufacturer_SUDOKU:
		return "SUDOKU"
	case KnxManufacturer_MOELLER_GEBAEUDEAUTOMATION_KG:
		return "MOELLER_GEBAEUDEAUTOMATION_KG"
	case KnxManufacturer_AZ_E_LITE_PTE_LTD:
		return "AZ_E_LITE_PTE_LTD"
	case KnxManufacturer_ARLIGHT:
		return "ARLIGHT"
	case KnxManufacturer_GRUENBECK_WASSERAUFBEREITUNG_GMBH:
		return "GRUENBECK_WASSERAUFBEREITUNG_GMBH"
	case KnxManufacturer_MODULE_ELECTRONIC:
		return "MODULE_ELECTRONIC"
	case KnxManufacturer_KOPLAT:
		return "KOPLAT"
	case KnxManufacturer_GUANGZHOU_LETOUR_LIFE_TECHNOLOGY_CO___LTD:
		return "GUANGZHOU_LETOUR_LIFE_TECHNOLOGY_CO___LTD"
	case KnxManufacturer_ILEVIA:
		return "ILEVIA"
	case KnxManufacturer_LN_SYSTEMTEQ:
		return "LN_SYSTEMTEQ"
	case KnxManufacturer_HISENSE_SMARTHOME:
		return "HISENSE_SMARTHOME"
	case KnxManufacturer_ELTAKO:
		return "ELTAKO"
	case KnxManufacturer_FLINK_AUTOMATION_SYSTEM:
		return "FLINK_AUTOMATION_SYSTEM"
	case KnxManufacturer_XXTER_BV:
		return "XXTER_BV"
	case KnxManufacturer_LYNXUS_TECHNOLOGY:
		return "LYNXUS_TECHNOLOGY"
	case KnxManufacturer_ROBOT_S_A_:
		return "ROBOT_S_A_"
	case KnxManufacturer_SHENZHEN_ATTE_SMART_LIFE_CO__LTD_:
		return "SHENZHEN_ATTE_SMART_LIFE_CO__LTD_"
	case KnxManufacturer_NOBLESSE:
		return "NOBLESSE"
	case KnxManufacturer_ADVANCED_DEVICES:
		return "ADVANCED_DEVICES"
	case KnxManufacturer_ATRINA_BUILDING_AUTOMATION_CO__LTD:
		return "ATRINA_BUILDING_AUTOMATION_CO__LTD"
	case KnxManufacturer_GUANGDONG_DAMING_LAFFEY_ELECTRIC_CO___LTD_:
		return "GUANGDONG_DAMING_LAFFEY_ELECTRIC_CO___LTD_"
	case KnxManufacturer_WESTERSTRAND_URFABRIK_AB:
		return "WESTERSTRAND_URFABRIK_AB"
	case KnxManufacturer_CONTROL4_CORPORATE:
		return "CONTROL4_CORPORATE"
	case KnxManufacturer_ONTROL:
		return "ONTROL"
	case KnxManufacturer_STARNET:
		return "STARNET"
	case KnxManufacturer_BETA_CAVI:
		return "BETA_CAVI"
	case KnxManufacturer_EASEMORE:
		return "EASEMORE"
	case KnxManufacturer_VIVALDI_SRL:
		return "VIVALDI_SRL"
	case KnxManufacturer_GREE_ELECTRIC_APPLIANCES_INC__OF_ZHUHAI:
		return "GREE_ELECTRIC_APPLIANCES_INC__OF_ZHUHAI"
	case KnxManufacturer_HWISCON:
		return "HWISCON"
	case KnxManufacturer_SHANGHAI_ELECON_INTELLIGENT_TECHNOLOGY_CO___LTD_:
		return "SHANGHAI_ELECON_INTELLIGENT_TECHNOLOGY_CO___LTD_"
	case KnxManufacturer_KAMPMANN:
		return "KAMPMANN"
	case KnxManufacturer_IMPOLUX_GMBH_LEDIMAX:
		return "IMPOLUX_GMBH_LEDIMAX"
	case KnxManufacturer_EVAUX:
		return "EVAUX"
	case KnxManufacturer_WEBRO_CABLES_AND_CONNECTORS_LIMITED:
		return "WEBRO_CABLES_AND_CONNECTORS_LIMITED"
	case KnxManufacturer_SHANGHAI_E_TECH_SOLUTION:
		return "SHANGHAI_E_TECH_SOLUTION"
	case KnxManufacturer_GUANGZHOU_HOKO_ELECTRIC_CO__LTD_:
		return "GUANGZHOU_HOKO_ELECTRIC_CO__LTD_"
	case KnxManufacturer_LAMMIN_HIGH_TECH_CO__LTD:
		return "LAMMIN_HIGH_TECH_CO__LTD"
	case KnxManufacturer_SHENZHEN_MERRYTEK_TECHNOLOGY_CO___LTD:
		return "SHENZHEN_MERRYTEK_TECHNOLOGY_CO___LTD"
	case KnxManufacturer_I_LUXUS:
		return "I_LUXUS"
	case KnxManufacturer_ELMOS_SEMICONDUCTOR_AG:
		return "ELMOS_SEMICONDUCTOR_AG"
	case KnxManufacturer_EMCOM_TECHNOLOGY_INC:
		return "EMCOM_TECHNOLOGY_INC"
	case KnxManufacturer_BOSCH_SIEMENS_HAUSHALTSGERAETE:
		return "BOSCH_SIEMENS_HAUSHALTSGERAETE"
	case KnxManufacturer_PROJECT_INNOVATIONS_GMBH:
		return "PROJECT_INNOVATIONS_GMBH"
	case KnxManufacturer_ITC:
		return "ITC"
	case KnxManufacturer_ABB_LV_INSTALLATION_MATERIALS_COMPANY_LTD__BEIJING:
		return "ABB_LV_INSTALLATION_MATERIALS_COMPANY_LTD__BEIJING"
	case KnxManufacturer_MAICO:
		return "MAICO"
	case KnxManufacturer_ELAN_SRL:
		return "ELAN_SRL"
	case KnxManufacturer_MINHHA_TECHNOLOGY_CO__LTD:
		return "MINHHA_TECHNOLOGY_CO__LTD"
	case KnxManufacturer_ZHEJIANG_TIANJIE_INDUSTRIAL_CORP_:
		return "ZHEJIANG_TIANJIE_INDUSTRIAL_CORP_"
	case KnxManufacturer_IAUTOMATION_PTY_LIMITED:
		return "IAUTOMATION_PTY_LIMITED"
	case KnxManufacturer_EXTRON:
		return "EXTRON"
	case KnxManufacturer_BTICINO:
		return "BTICINO"
	case KnxManufacturer_FREEDOMPRO:
		return "FREEDOMPRO"
	case KnxManufacturer_1HOME:
		return "1HOME"
	case KnxManufacturer_EOS_SAUNATECHNIK_GMBH:
		return "EOS_SAUNATECHNIK_GMBH"
	case KnxManufacturer_KUSATEK_GMBH:
		return "KUSATEK_GMBH"
	case KnxManufacturer_EISBAER_SCADA:
		return "EISBAER_SCADA"
	case KnxManufacturer_AUTOMATISMI_BENINCA_S_P_A_:
		return "AUTOMATISMI_BENINCA_S_P_A_"
	case KnxManufacturer_BLENDOM:
		return "BLENDOM"
	case KnxManufacturer_MADEL_AIR_TECHNICAL_DIFFUSION:
		return "MADEL_AIR_TECHNICAL_DIFFUSION"
	case KnxManufacturer_NIKO:
		return "NIKO"
	case KnxManufacturer_BOSCH_REXROTH_AG:
		return "BOSCH_REXROTH_AG"
	case KnxManufacturer_CANDM_PRODUCTS:
		return "CANDM_PRODUCTS"
	case KnxManufacturer_HOERMANN_KG_VERKAUFSGESELLSCHAFT:
		return "HOERMANN_KG_VERKAUFSGESELLSCHAFT"
	case KnxManufacturer_SHANGHAI_RAJAYASA_CO__LTD:
		return "SHANGHAI_RAJAYASA_CO__LTD"
	case KnxManufacturer_SUZUKI:
		return "SUZUKI"
	case KnxManufacturer_SILENT_GLISS_INTERNATIONAL_LTD_:
		return "SILENT_GLISS_INTERNATIONAL_LTD_"
	case KnxManufacturer_BEE_CONTROLS_ADGSC_GROUP:
		return "BEE_CONTROLS_ADGSC_GROUP"
	case KnxManufacturer_XDTECGMBH:
		return "XDTECGMBH"
	case KnxManufacturer_OSRAM:
		return "OSRAM"
	case KnxManufacturer_RITTO_GMBHANDCO_KG:
		return "RITTO_GMBHANDCO_KG"
	case KnxManufacturer_LEBENOR:
		return "LEBENOR"
	case KnxManufacturer_AUTOMANENG:
		return "AUTOMANENG"
	case KnxManufacturer_HONEYWELL_AUTOMATION_SOLUTION_CONTROLCHINA:
		return "HONEYWELL_AUTOMATION_SOLUTION_CONTROLCHINA"
	case KnxManufacturer_HANGZHOU_BINTHEN_INTELLIGENCE_TECHNOLOGY_CO__LTD:
		return "HANGZHOU_BINTHEN_INTELLIGENCE_TECHNOLOGY_CO__LTD"
	case KnxManufacturer_ETA_HEIZTECHNIK:
		return "ETA_HEIZTECHNIK"
	case KnxManufacturer_DIVUS_GMBH:
		return "DIVUS_GMBH"
	case KnxManufacturer_NANJING_TAIJIESAI_INTELLIGENT_TECHNOLOGY_CO__LTD_:
		return "NANJING_TAIJIESAI_INTELLIGENT_TECHNOLOGY_CO__LTD_"
	case KnxManufacturer_LUNATONE:
		return "LUNATONE"
	case KnxManufacturer_ZHEJIANG_SCTECH_BUILDING_INTELLIGENT:
		return "ZHEJIANG_SCTECH_BUILDING_INTELLIGENT"
	case KnxManufacturer_FOSHAN_QITE_TECHNOLOGY_CO___LTD_:
		return "FOSHAN_QITE_TECHNOLOGY_CO___LTD_"
	case KnxManufacturer_POWER_CONTROLS:
		return "POWER_CONTROLS"
	case KnxManufacturer_NOKE:
		return "NOKE"
	case KnxManufacturer_LANDCOM:
		return "LANDCOM"
	case KnxManufacturer_STORK_AS:
		return "STORK_AS"
	case KnxManufacturer_HANGZHOU_SHENDU_TECHNOLOGY_CO___LTD_:
		return "HANGZHOU_SHENDU_TECHNOLOGY_CO___LTD_"
	case KnxManufacturer_COOLAUTOMATION:
		return "COOLAUTOMATION"
	case KnxManufacturer_APRSTERN:
		return "APRSTERN"
	case KnxManufacturer_SONNEN:
		return "SONNEN"
	case KnxManufacturer_DNAKE:
		return "DNAKE"
	case KnxManufacturer_NEUBERGER_GEBAEUDEAUTOMATION_GMBH:
		return "NEUBERGER_GEBAEUDEAUTOMATION_GMBH"
	case KnxManufacturer_STILIGER:
		return "STILIGER"
	case KnxManufacturer_BERGHOF_AUTOMATION_GMBH:
		return "BERGHOF_AUTOMATION_GMBH"
	case KnxManufacturer_TOTAL_AUTOMATION_AND_CONTROLS_GMBH:
		return "TOTAL_AUTOMATION_AND_CONTROLS_GMBH"
	case KnxManufacturer_DOVIT:
		return "DOVIT"
	case KnxManufacturer_INSTALIGHTING_GMBH:
		return "INSTALIGHTING_GMBH"
	case KnxManufacturer_UNI_TEC:
		return "UNI_TEC"
	case KnxManufacturer_CASATUNES:
		return "CASATUNES"
	case KnxManufacturer_EMT:
		return "EMT"
	case KnxManufacturer_SENFFICIENT:
		return "SENFFICIENT"
	case KnxManufacturer_AUROLITE_ELECTRICAL_PANYU_GUANGZHOU_LIMITED:
		return "AUROLITE_ELECTRICAL_PANYU_GUANGZHOU_LIMITED"
	case KnxManufacturer_ABB_XIAMEN_SMART_TECHNOLOGY_CO___LTD_:
		return "ABB_XIAMEN_SMART_TECHNOLOGY_CO___LTD_"
	case KnxManufacturer_ZUMTOBEL:
		return "ZUMTOBEL"
	case KnxManufacturer_SAMSON_ELECTRIC_WIRE:
		return "SAMSON_ELECTRIC_WIRE"
	case KnxManufacturer_T_TOUCHING:
		return "T_TOUCHING"
	case KnxManufacturer_CORE_SMART_HOME:
		return "CORE_SMART_HOME"
	case KnxManufacturer_GREENCONNECT_SOLUTIONS_SA:
		return "GREENCONNECT_SOLUTIONS_SA"
	case KnxManufacturer_ELETTRONICA_CONDUTTORI:
		return "ELETTRONICA_CONDUTTORI"
	case KnxManufacturer_MKFC:
		return "MKFC"
	case KnxManufacturer_AUTOMATIONPlus:
		return "AUTOMATIONPlus"
	case KnxManufacturer_BLUE_AND_RED:
		return "BLUE_AND_RED"
	case KnxManufacturer_FROGBLUE:
		return "FROGBLUE"
	case KnxManufacturer_SAVESOR:
		return "SAVESOR"
	case KnxManufacturer_APP_TECH:
		return "APP_TECH"
	case KnxManufacturer_SENSORTEC_AG:
		return "SENSORTEC_AG"
	case KnxManufacturer_NYSA_TECHNOLOGY_AND_SOLUTIONS:
		return "NYSA_TECHNOLOGY_AND_SOLUTIONS"
	case KnxManufacturer_FARADITE:
		return "FARADITE"
	case KnxManufacturer_OPTIMUS:
		return "OPTIMUS"
	case KnxManufacturer_KTS_S_R_L_:
		return "KTS_S_R_L_"
	case KnxManufacturer_RAMCRO_SPA:
		return "RAMCRO_SPA"
	case KnxManufacturer_WUHAN_WISECREATE_UNIVERSE_TECHNOLOGY_CO___LTD:
		return "WUHAN_WISECREATE_UNIVERSE_TECHNOLOGY_CO___LTD"
	case KnxManufacturer_BEMI_SMART_HOME_LTD:
		return "BEMI_SMART_HOME_LTD"
	case KnxManufacturer_ARDOMUS:
		return "ARDOMUS"
	case KnxManufacturer_PHOENIX_CONTACT:
		return "PHOENIX_CONTACT"
	case KnxManufacturer_CHANGXING:
		return "CHANGXING"
	case KnxManufacturer_E_CONTROLS:
		return "E_CONTROLS"
	case KnxManufacturer_AIB_TECHNOLOGY:
		return "AIB_TECHNOLOGY"
	case KnxManufacturer_NVC:
		return "NVC"
	case KnxManufacturer_KBOX:
		return "KBOX"
	case KnxManufacturer_CNS:
		return "CNS"
	case KnxManufacturer_TYBA:
		return "TYBA"
	case KnxManufacturer_ATREL:
		return "ATREL"
	case KnxManufacturer_SIMON_ELECTRIC_CHINA_CO___LTD:
		return "SIMON_ELECTRIC_CHINA_CO___LTD"
	case KnxManufacturer_KORDZ_GROUP:
		return "KORDZ_GROUP"
	case KnxManufacturer_ND_ELECTRIC:
		return "ND_ELECTRIC"
	case KnxManufacturer_CONTROLIUM:
		return "CONTROLIUM"
	case KnxManufacturer_FAMO_GMBH_AND_CO__KG:
		return "FAMO_GMBH_AND_CO__KG"
	case KnxManufacturer_CDN_SMART:
		return "CDN_SMART"
	case KnxManufacturer_HESTON:
		return "HESTON"
	case KnxManufacturer_ESLA_CONEXIONES_S_L_:
		return "ESLA_CONEXIONES_S_L_"
	case KnxManufacturer_WEISHAUPT:
		return "WEISHAUPT"
	case KnxManufacturer_ASTRUM_TECHNOLOGY:
		return "ASTRUM_TECHNOLOGY"
	case KnxManufacturer_WUERTH_ELEKTRONIK_STELVIO_KONTEK_S_P_A_:
		return "WUERTH_ELEKTRONIK_STELVIO_KONTEK_S_P_A_"
	case KnxManufacturer_NANOTECO_CORPORATION:
		return "NANOTECO_CORPORATION"
	case KnxManufacturer_NIETIAN:
		return "NIETIAN"
	case KnxManufacturer_SUMSIR:
		return "SUMSIR"
	case KnxManufacturer_ORBIS_TECNOLOGIA_ELECTRICA_SA:
		return "ORBIS_TECNOLOGIA_ELECTRICA_SA"
	case KnxManufacturer_BERKER:
		return "BERKER"
	case KnxManufacturer_WAGO_KONTAKTTECHNIK:
		return "WAGO_KONTAKTTECHNIK"
	case KnxManufacturer_KNXPRESSO:
		return "KNXPRESSO"
	case KnxManufacturer_WIELAND_ELECTRIC:
		return "WIELAND_ELECTRIC"
	case KnxManufacturer_HERMANN_KLEINHUIS:
		return "HERMANN_KLEINHUIS"
	case KnxManufacturer_STIEBEL_ELTRON:
		return "STIEBEL_ELTRON"
	case KnxManufacturer_BUSCH_JAEGER_ELEKTRO:
		return "BUSCH_JAEGER_ELEKTRO"
	case KnxManufacturer_TEHALIT:
		return "TEHALIT"
	case KnxManufacturer_THEBEN_AG:
		return "THEBEN_AG"
	case KnxManufacturer_WILHELM_RUTENBECK:
		return "WILHELM_RUTENBECK"
	case KnxManufacturer_WINKHAUS:
		return "WINKHAUS"
	case KnxManufacturer_ROBERT_BOSCH:
		return "ROBERT_BOSCH"
	case KnxManufacturer_SOMFY:
		return "SOMFY"
	case KnxManufacturer_GIRA_GIERSIEPEN:
		return "GIRA_GIERSIEPEN"
	case KnxManufacturer_WOERTZ:
		return "WOERTZ"
	case KnxManufacturer_VIESSMANN_WERKE:
		return "VIESSMANN_WERKE"
	case KnxManufacturer_IMI_HYDRONIC_ENGINEERING:
		return "IMI_HYDRONIC_ENGINEERING"
	case KnxManufacturer_JOH__VAILLANT:
		return "JOH__VAILLANT"
	case KnxManufacturer_AMP_DEUTSCHLAND:
		return "AMP_DEUTSCHLAND"
	case KnxManufacturer_BOSCH_THERMOTECHNIK_GMBH:
		return "BOSCH_THERMOTECHNIK_GMBH"
	case KnxManufacturer_HAGER_ELECTRO:
		return "HAGER_ELECTRO"
	case KnxManufacturer_SEF___ECOTEC:
		return "SEF___ECOTEC"
	case KnxManufacturer_DORMA_GMBH_Plus_CO__KG:
		return "DORMA_GMBH_Plus_CO__KG"
	case KnxManufacturer_WINDOWMASTER_AS:
		return "WINDOWMASTER_AS"
	case KnxManufacturer_WALTHER_WERKE:
		return "WALTHER_WERKE"
	case KnxManufacturer_ORAS:
		return "ORAS"
	case KnxManufacturer_DAETWYLER:
		return "DAETWYLER"
	case KnxManufacturer_ELECTRAK:
		return "ELECTRAK"
	case KnxManufacturer_TECHEM:
		return "TECHEM"
	}
	return ""
}
