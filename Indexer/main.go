package main

import (
	"mail-indexer/constants"
	parserhelpers "mail-indexer/parserHelpers"
	fetchzinc "mail-indexer/zincHelpers"
)

func main() {
	fetchzinc.FetchCreateZincIndex()
	parserhelpers.ProcessFiles(constants.FILE_NAME)
}

// CPU profiling
// go tool pprof cpu.prof

// MEMORY profiling
// go tool pprof mem_profile.prof
