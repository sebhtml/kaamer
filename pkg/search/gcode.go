/*
Copyright 2019 The kaamer Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package search

var gcodeBacteria = map[string]AminoAcid{
	"ttt": AminoAcid{AA: "F", Start: false, Stop: false},
	"ttc": AminoAcid{AA: "F", Start: false, Stop: false},
	"tta": AminoAcid{AA: "L", Start: false, Stop: false},
	"ttg": AminoAcid{AA: "L", Start: true, Stop: false},
	"tct": AminoAcid{AA: "S", Start: false, Stop: false},
	"tcc": AminoAcid{AA: "S", Start: false, Stop: false},
	"tca": AminoAcid{AA: "S", Start: false, Stop: false},
	"tcg": AminoAcid{AA: "S", Start: false, Stop: false},
	"tat": AminoAcid{AA: "Y", Start: false, Stop: false},
	"tac": AminoAcid{AA: "Y", Start: false, Stop: false},
	"taa": AminoAcid{AA: "*", Start: false, Stop: true},
	"tag": AminoAcid{AA: "*", Start: false, Stop: true},
	"tgt": AminoAcid{AA: "C", Start: false, Stop: false},
	"tgc": AminoAcid{AA: "C", Start: false, Stop: false},
	"tga": AminoAcid{AA: "*", Start: false, Stop: true},
	"tgg": AminoAcid{AA: "W", Start: false, Stop: false},
	"ctt": AminoAcid{AA: "L", Start: false, Stop: false},
	"ctc": AminoAcid{AA: "L", Start: false, Stop: false},
	"cta": AminoAcid{AA: "L", Start: false, Stop: false},
	"ctg": AminoAcid{AA: "L", Start: true, Stop: false},
	"cct": AminoAcid{AA: "P", Start: false, Stop: false},
	"ccc": AminoAcid{AA: "P", Start: false, Stop: false},
	"cca": AminoAcid{AA: "P", Start: false, Stop: false},
	"ccg": AminoAcid{AA: "P", Start: false, Stop: false},
	"cat": AminoAcid{AA: "H", Start: false, Stop: false},
	"cac": AminoAcid{AA: "H", Start: false, Stop: false},
	"caa": AminoAcid{AA: "Q", Start: false, Stop: false},
	"cag": AminoAcid{AA: "Q", Start: false, Stop: false},
	"cgt": AminoAcid{AA: "R", Start: false, Stop: false},
	"cgc": AminoAcid{AA: "R", Start: false, Stop: false},
	"cga": AminoAcid{AA: "R", Start: false, Stop: false},
	"cgg": AminoAcid{AA: "R", Start: false, Stop: false},
	"att": AminoAcid{AA: "I", Start: true, Stop: false},
	"atc": AminoAcid{AA: "I", Start: true, Stop: false},
	"ata": AminoAcid{AA: "I", Start: true, Stop: false},
	"atg": AminoAcid{AA: "M", Start: true, Stop: false},
	"act": AminoAcid{AA: "T", Start: false, Stop: false},
	"acc": AminoAcid{AA: "T", Start: false, Stop: false},
	"aca": AminoAcid{AA: "T", Start: false, Stop: false},
	"acg": AminoAcid{AA: "T", Start: false, Stop: false},
	"aat": AminoAcid{AA: "N", Start: false, Stop: false},
	"aac": AminoAcid{AA: "N", Start: false, Stop: false},
	"aaa": AminoAcid{AA: "K", Start: false, Stop: false},
	"aag": AminoAcid{AA: "K", Start: false, Stop: false},
	"agt": AminoAcid{AA: "S", Start: false, Stop: false},
	"agc": AminoAcid{AA: "S", Start: false, Stop: false},
	"aga": AminoAcid{AA: "R", Start: false, Stop: false},
	"agg": AminoAcid{AA: "R", Start: false, Stop: false},
	"gtt": AminoAcid{AA: "V", Start: false, Stop: false},
	"gtc": AminoAcid{AA: "V", Start: false, Stop: false},
	"gta": AminoAcid{AA: "V", Start: false, Stop: false},
	"gtg": AminoAcid{AA: "V", Start: true, Stop: false},
	"gct": AminoAcid{AA: "A", Start: false, Stop: false},
	"gcc": AminoAcid{AA: "A", Start: false, Stop: false},
	"gca": AminoAcid{AA: "A", Start: false, Stop: false},
	"gcg": AminoAcid{AA: "A", Start: false, Stop: false},
	"gat": AminoAcid{AA: "D", Start: false, Stop: false},
	"gac": AminoAcid{AA: "D", Start: false, Stop: false},
	"gaa": AminoAcid{AA: "E", Start: false, Stop: false},
	"gag": AminoAcid{AA: "E", Start: false, Stop: false},
	"ggt": AminoAcid{AA: "G", Start: false, Stop: false},
	"ggc": AminoAcid{AA: "G", Start: false, Stop: false},
	"gga": AminoAcid{AA: "G", Start: false, Stop: false},
	"ggg": AminoAcid{AA: "G", Start: false, Stop: false},
}

var startCodonWeight = map[string]int{
	"atg": 85,
	"gtg": 15,
	"ttg": 10,
	"att": 1,
	"ata": 1,
	"atc": 1,
}
