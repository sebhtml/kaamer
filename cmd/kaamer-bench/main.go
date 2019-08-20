package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dgraph-io/badger/options"
	"github.com/pkg/profile"
	"github.com/zorino/kaamer/pkg/kvstore"
	"github.com/zorino/kaamer/pkg/makedb"
)

const (
	MaxInt uint32 = 1<<32 - 1
)

var (
	/*program options*/
	function      = flag.String("f", "", "function to test")
	profilingMode = flag.String("prof", "cpu", "profiling mode")

	/*shared options*/
	dbPath    = flag.String("d", "", "db path argument")
	tableMode = flag.String("tablemode", "memorymap", "table loading mode (fileio, memorymap)")
	valueMode = flag.String("valuemode", "memorymap", "value loading mode (fileio, memorymap)")
	maxSize   = flag.Bool("maxsize", false, "to maximize badger output file size")

	/*opendb options*/
	portNumber = flag.Int("p", 8321, "port argument")
	nbThreads  = flag.Int("t", 2, "number of threads")
	tmpFolder  = flag.String("tmp", "/tmp/", "tmp folder for query import")

	/*makedb options*/
	inputPath    = flag.String("i", "", "input file argument")
	makedbOffset = flag.Uint("offset", 0, "offset to process raw file")
	makedbLenght = flag.Uint("length", uint(MaxInt), "process x number of files")
	noIndex      = flag.Bool("noindex", false, "prevent the indexing of database")

	LoadingMode = map[string]options.FileLoadingMode{"memorymap": options.MemoryMap, "fileio": options.FileIO}
)

func main() {

	usage := `
 kaamer-bench <function> [OPTIONS]

  // Functions to profile

  -opendb (db opened in read-only mode)
    (input)
      -d            database directory
      -p            port (default: 8321)
      -t            number of threads to use (default all)
      -tmp          tmp folder for query import (default /tmp)
      -tableMode    (fileio, memorymap) default memorymap / fileio decreases memory usage
      -valueMode    (fileio, memorymap) default memorymap / fileio decreases memory usage
    (flag)
      -maxsize      will maximize the size of tables (.sst) and vlog (.log) files
                    (to limit the number of open files)

  -makedb
    (input)
      -i            input raw EMBL file
      -d            badger database directory (output)
      -offset       start processing raw uniprot file at protein number x
      -length       process x number of proteins (-1 == infinity)
      -tableMode    (fileio, memorymap) default memorymap / fileio decreases memory usage
      -valueMode    (fileio, memorymap) default memorymap / fileio decreases memory usage
    (flag)
      -maxsize      will maximize the size of tables (.sst) and vlog (.log) files
                    (to limit the number of open files)
      -noindex      will NOT index the database - need to be done afterward with -index

`

	flag.Parse()

	// errHelp := errors.New(usage)

	/* Profiling Mode */
	switch *profilingMode {
	case "cpu":
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	case "mem":
		defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
	default:
		defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	}

	/* Setting values from CLI */
	var tableLoadingMode options.FileLoadingMode
	var valueLoadingMode options.FileLoadingMode
	var ok = false

	if tableLoadingMode, ok = LoadingMode[*tableMode]; !ok {
		fmt.Println("TableMode unrecognized ! use fileio or memorymap!")
		os.Exit(1)
	}
	if valueLoadingMode, ok = LoadingMode[*valueMode]; !ok {
		fmt.Println("ValueMode unrecognized ! use fileio or memorymap!")
		os.Exit(1)
	}

	/* Functions */
	switch *function {
	case "opendb":
		if *dbPath == "" {
			fmt.Println("No output db path !")
			os.Exit(1)
		} else {
			var kvStores *kvstore.KVStores
			kvStores = kvstore.KVStoresNew(*dbPath, *nbThreads, tableLoadingMode, valueLoadingMode, true, false, true)
			defer kvStores.Close()
		}
	case "makedb":
		if *dbPath == "" {
			fmt.Println("No output db path !")
			os.Exit(1)
		} else if *inputPath == "" {
			fmt.Println("No input file !")
			os.Exit(1)
		} else {
			makedb.NewMakedb(*dbPath, *inputPath, *makedbOffset, *makedbLenght, *maxSize, tableLoadingMode, valueLoadingMode, *noIndex)
		}
	default:
		os.Remove("cpu.pprof")
		os.Remove("mem.pprof")
		fmt.Println(usage)
	}

}