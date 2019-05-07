package search

import (
	"runtime"

	"github.com/zorino/metaprot/pkg/kvstore"
	"github.com/zorino/metaprot/pkg/search"
)

const (
	FILE              = 0
	STRING            = 1
	NUCLEOTIDE_STRING = 0
	NUCLEOTIDE_FILE   = 1
	PROTEIN_STRING    = 2
	PROTEIN_FILE      = 3
)

func NewSearch(dbPath string, input string, inputType int) {

	nbOfThreads := runtime.NumCPU()
	kvStores := kvstore.KVStoresNew(dbPath, nbOfThreads)


	search.NewSearchResult(input, PROTEIN_STRING, kvStores, nbOfThreads)
	// searchRes.G_hits

}
