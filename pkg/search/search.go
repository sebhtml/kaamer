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

import (
	"bufio"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
	cnt "github.com/zorino/counters"
	"github.com/zorino/kaamer/pkg/kvstore"
)

const (
	NUCLEOTIDE    = 0
	PROTEIN       = 1
	READS         = 2
	KMER_SIZE     = 7
	DNA_QUERY     = "DNA Query"
	PROTEIN_QUERY = "Protein Query"
)

var (
	searchOptions = SearchOptions{}
)

type SearchOptions struct {
	File             string
	InputType        string
	SequenceType     int
	OutFormat        string
	MaxResults       int
	ExtractPositions bool
	Annotations      bool
}

type SearchResults struct {
	Counter      *cnt.CounterBox
	Hits         HitList
	PositionHits map[uint32][]bool
}

type KeyPos struct {
	Key   []byte
	Pos   int
	QSize int
}

type MatchPosition struct {
	HitId uint32
	QPos  int
	QSize int
}

type QueryWriter struct {
	Query
	http.ResponseWriter
}

type QueryResult struct {
	Query         Query
	SearchResults *SearchResults
	HitEntries    map[uint32]kvstore.Protein
}

type Query struct {
	Sequence   string
	Name       string
	SizeInKmer int
	Type       string
	Location   Location
	Contig     string
}

type Hit struct {
	Key    uint32
	Kmatch int64
}

type HitList []Hit

func (p HitList) Len() int           { return len(p) }
func (p HitList) Less(i, j int) bool { return p[i].Kmatch < p[j].Kmatch }
func (p HitList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func syncMapLen(syncMap *sync.Map) int {
	length := 0
	syncMap.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}

func sortMapByValue(hitFrequencies *sync.Map) HitList {
	pl := make(HitList, syncMapLen(hitFrequencies))
	i := 0

	hitFrequencies.Range(func(k, v interface{}) bool {
		key, okKey := k.(string)
		item, okValue := v.(cnt.Counter)
		if okKey && okValue {
			idUint32, err := strconv.Atoi(key)
			if err != nil {
				log.Fatal(err.Error())
			}
			pl[i] = Hit{uint32(idUint32), item.Value()}
			i++
		}
		return true
	})

	sort.Sort(sort.Reverse(pl))
	return pl
}

func NewSearchResult(newSearchOptions SearchOptions, kvStores *kvstore.KVStores, nbOfThreads int, w http.ResponseWriter) []QueryResult {

	// sequence is either file path or the actual sequence (depends on sequenceType)
	queryResults := []QueryResult{}

	searchOptions = newSearchOptions

	switch searchOptions.SequenceType {
	case READS:
		fmt.Println("Searching for Reads file")
		NucleotideSearch(searchOptions, kvStores, nbOfThreads, w, true)
	case NUCLEOTIDE:
		fmt.Println("Searching for Nucleotide file")
		NucleotideSearch(searchOptions, kvStores, nbOfThreads, w, false)
	case PROTEIN:
		fmt.Println("Searching from Protein file")
		ProteinSearch(searchOptions, kvStores, nbOfThreads, w)
	}

	if searchOptions.InputType != "path" {
		os.Remove(searchOptions.File)
	}

	return queryResults

}

func (queryResult *QueryResult) FilterResults(kmerMatchRatio float64) {

	var hitsToDelete []uint32
	var lastGoodHitPosition = len(queryResult.SearchResults.Hits) - 1

	for i, hit := range queryResult.SearchResults.Hits {
		if (float64(hit.Kmatch)/float64(queryResult.Query.SizeInKmer)) < kmerMatchRatio || hit.Kmatch < 10 {
			if lastGoodHitPosition == (len(queryResult.SearchResults.Hits) - 1) {
				lastGoodHitPosition = i - 1
			}
			hitsToDelete = append(hitsToDelete, hit.Key)
		}
	}

	if lastGoodHitPosition >= searchOptions.MaxResults {
		lastGoodHitPosition = searchOptions.MaxResults - 1
		for _, h := range queryResult.SearchResults.Hits[lastGoodHitPosition+1:] {
			hitsToDelete = append(hitsToDelete, h.Key)
		}
	}

	if lastGoodHitPosition < 0 {
		queryResult.SearchResults.Hits = []Hit{}
	} else {
		queryResult.SearchResults.Hits = queryResult.SearchResults.Hits[0 : lastGoodHitPosition+1]
	}

	for _, k := range hitsToDelete {
		delete(queryResult.SearchResults.PositionHits, k)
	}

}

func GetQueriesFasta(fileName string, queryChan chan<- Query, isProtein bool) {

	loc := Location{
		StartPosition:     1,
		EndPosition:       0,
		PlusStrand:        true,
		StartsAlternative: []int{},
	}
	query := Query{
		Sequence:   "",
		Name:       "",
		SizeInKmer: 0,
		Location:   loc,
		Type:       "",
		Contig:     "",
	}

	// queries := []Query{}
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	// check filetype
	buff := make([]byte, 32)
	_, err = file.Read(buff)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	file.Seek(0, 0)

	filetype := http.DetectContentType(buff)

	scanner := new(bufio.Scanner)

	if filetype == "application/x-gzip" {
		gz, err := gzip.NewReader(file)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		scanner = bufio.NewScanner(gz)
		defer gz.Close()
	} else if filetype == "text/plain; charset=utf-8" {
		scanner = bufio.NewScanner(file)
	} else {
		return
	}

	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	l := ""
	queryName := ""

	for scanner.Scan() {
		l = scanner.Text()
		if len(l) < 1 {
			continue
		}
		if l[0] == '>' {
			if query.Sequence != "" {
				query.SizeInKmer = len(query.Sequence) - KMER_SIZE + 1
				if query.Sequence[len(query.Sequence)-1:] == "*" {
					query.SizeInKmer--
				}
				query.Location.EndPosition = len(query.Sequence)
				queryChan <- query
				query = Query{Sequence: "", Name: "", SizeInKmer: 0, Contig: ""}
			}
			queryName = strings.TrimSuffix(l[1:], "\n")
			if isProtein {
				query.Name = queryName
				query.Contig = ""
				query.Location.StartPosition = 1
			} else {
				query.Name = queryName
				query.Contig = queryName
			}
		} else {
			query.Sequence += strings.TrimSpace(l)
		}
	}

	if query.Sequence != "" {
		query.SizeInKmer = len(query.Sequence) - KMER_SIZE + 1
		if query.Sequence[len(query.Sequence)-1:] == "*" {
			query.SizeInKmer--
		}
		query.Location.EndPosition = len(query.Sequence)
		queryChan <- query
	}

}

func GetQueriesFastq(fileName string, queryChan chan<- Query) {

	loc := Location{
		StartPosition:     1,
		EndPosition:       0,
		PlusStrand:        true,
		StartsAlternative: []int{},
	}
	query := Query{
		Sequence:   "",
		Name:       "",
		SizeInKmer: 0,
		Location:   loc,
		Type:       "",
		Contig:     "",
	}

	// queries := []Query{}
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	// check filetype
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filetype := http.DetectContentType(buff)

	scanner := new(bufio.Scanner)

	if filetype == "application/x-gzip" {
		// fmt.Println("Loaded gzip file")
		gz, err := gzip.NewReader(file)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		scanner = bufio.NewScanner(gz)
		// defer gz.Close()
	} else if filetype == "text/plain; charset=utf-8" {
		scanner = bufio.NewScanner(file)
	} else {
		return
	}

	isSequence := regexp.MustCompile(`^[ATGCNatgcn]+$`).MatchString
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	seqNb := 0
	l := ""

	for scanner.Scan() {
		l = scanner.Text()
		if len(l) < 1 {
			continue
		}
		if l[0] == '@' {
			seqNb += 1
			if query.Sequence != "" {
				query.SizeInKmer = len(query.Sequence) - KMER_SIZE + 1
				query.Location.EndPosition = len(query.Sequence)
				queryChan <- query
				query = Query{Sequence: "", Name: "", SizeInKmer: 0}
			}
			query.Name = strings.TrimSuffix(l[1:], "\n")
		} else if isSequence(l) {
			query.Sequence = strings.TrimSuffix(l, "\n")
		}
	}

	if query.Sequence != "" {
		query.SizeInKmer = len(query.Sequence) - KMER_SIZE + 1
		query.Location.EndPosition = len(query.Sequence)
		queryChan <- query
	}

}

func (searchRes *SearchResults) KmerSearch(keyChan <-chan KeyPos, kvStores *kvstore.KVStores, wg *sync.WaitGroup, matchPositionChan chan<- MatchPosition) {

	defer wg.Done()
	for keyPos := range keyChan {

		if kCombId, err := kvStores.KmerStore.GetValueFromBadger(keyPos.Key); err == nil {

			if len(kCombId) < 1 {
				continue
			}

			kCombVal, _ := kvStores.KCombStore.GetValueFromBadger(kCombId)
			kC := &kvstore.KComb{}
			proto.Unmarshal(kCombVal, kC)

			for _, id := range kC.ProteinKeys {
				searchRes.Counter.GetCounter(strconv.Itoa(int(id))).Increment()
				if searchOptions.ExtractPositions {
					matchPositionChan <- MatchPosition{HitId: id, QPos: keyPos.Pos, QSize: keyPos.QSize}
				}
			}

		}
	}

}

func (searchRes *SearchResults) StoreMatchPositions(matchPosition <-chan MatchPosition, wg *sync.WaitGroup) {

	defer wg.Done()
	for mp := range matchPosition {
		if _, ok := searchRes.PositionHits[mp.HitId]; !ok {
			searchRes.PositionHits[mp.HitId] = make([]bool, mp.QSize)
		}
		searchRes.PositionHits[mp.HitId][mp.QPos] = true
	}

}

func (queryResult *QueryResult) FetchHitsInformation(kvStores *kvstore.KVStores) {

	for _, h := range queryResult.SearchResults.Hits {
		proteinId := make([]byte, 4)
		binary.BigEndian.PutUint32(proteinId, h.Key)
		val, err := kvStores.ProteinStore.GetValueFromBadger(proteinId)
		if err != nil {
			return
		}
		prot := &kvstore.Protein{}
		proto.Unmarshal(val, prot)
		queryResult.HitEntries[h.Key] = *prot
	}

}

func QueryResultResponseWriter(queryResult <-chan QueryResult, w http.ResponseWriter, wg *sync.WaitGroup) {

	defer wg.Done()

	if searchOptions.OutFormat == "tsv" {

		// set http response header
		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
		w.WriteHeader(200)

		w.Write([]byte("QueryName\tQueryKSize\tQStart\tQEnd\tKMatch\tHit.Id"))
		if searchOptions.Annotations {
			w.Write([]byte("\tHit.ProteinName\tHit.Organism\tHit.EC\tHit.GO\tHit.HAMAP\tHit.KEGG_id\tHit.KEGG_pathway\tHit.Biocyc_id\tHit.Biocyc_pathway\tHit.Taxonomy"))
		}
		if searchOptions.ExtractPositions {
			w.Write([]byte("\tQueryHit.Positions"))
		}
		w.Write([]byte("\n"))
		output := ""

		for qR := range queryResult {

			for _, h := range qR.SearchResults.Hits {
				output = ""
				output += qR.Query.Name
				output += "\t"
				output += strconv.Itoa(qR.Query.SizeInKmer)
				output += "\t"
				output += strconv.Itoa(qR.Query.Location.StartPosition)
				output += "\t"
				output += strconv.Itoa(qR.Query.Location.EndPosition)
				output += "\t"
				output += strconv.Itoa(int(h.Kmatch))
				output += "\t"
				output += qR.HitEntries[h.Key].Entry
				if searchOptions.Annotations {
					output += "\t"
					output += qR.HitEntries[h.Key].ProteinName
					output += "\t"
					output += qR.HitEntries[h.Key].Organism
					output += "\t"
					output += qR.HitEntries[h.Key].EC
					output += "\t"
					output += strings.Join(qR.HitEntries[h.Key].GO, ";")
					output += "\t"
					output += strings.Join(qR.HitEntries[h.Key].HAMAP, ";")
					output += "\t"
					output += strings.Join(qR.HitEntries[h.Key].KEGG, ";")
					output += "\t"
					output += strings.Join(qR.HitEntries[h.Key].KEGG_Pathways, ";")
					output += "\t"
					output += strings.Join(qR.HitEntries[h.Key].BioCyc, ";")
					output += "\t"
					output += strings.Join(qR.HitEntries[h.Key].Biocyc_Pathways, ";")
					output += "\t"
					output += qR.HitEntries[h.Key].Taxonomy
				}
				if searchOptions.ExtractPositions {
					output += "\t"
					output += FormatPositionsToString(qR.SearchResults.PositionHits[h.Key])
				}
				output += "\n"
				w.Write([]byte(output))
			}

		}

	}

	if searchOptions.OutFormat == "json" {

		// set http response header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		// open results array
		w.Write([]byte("["))

		firstResult := true

		for qR := range queryResult {

			if !firstResult {
				w.Write([]byte(","))
			}
			data, err := json.Marshal(qR)
			if err != nil {
				fmt.Println(err.Error())
			}
			w.Write(data)

			firstResult = false

		}

		// open results array
		w.Write([]byte("]"))

	}

}

func FormatPositionsToString(positions []bool) string {

	currentStart := 0
	inSequence := false

	positionsString := ""

	for pos, match := range positions {
		if match {
			if !inSequence {
				currentStart = pos + 1
				inSequence = true
			}
		} else {
			if inSequence {
				if pos+1 > currentStart {
					if positionsString != "" {
						positionsString += ","
					}
					positionsString += (strconv.Itoa(currentStart) + "-" + (strconv.Itoa(pos + 1)))
					inSequence = false
				} else {
					if positionsString != "" {
						positionsString += ","
					}
					positionsString += strconv.Itoa(currentStart)
					inSequence = false
				}
			}
		}
	}
	if inSequence {
		if positionsString != "" {
			positionsString += ","
		}
		positionsString += (strconv.Itoa(currentStart) + "-" + (strconv.Itoa(len(positions))))
	}

	return positionsString

}
