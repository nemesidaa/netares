package main

import (
	"flag"
	"fmt"
	httpclient "netares/internal/httpc"
	"netares/internal/parser"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	maskFile   string // * path to mask file.
	outputType string // * type of output, "raw" only now.
	targetName string // * target name, can be everything, which we can use instead of * in the sources of the masks.
	timeout    uint   // * in milliseconds.
	retries    uint   // * number of retries.

	// ! Not the flags.
	masks []string // * found masks.
)

func ParseFlags() {
	flag.StringVar(&maskFile, "mask", "./...", "path to mask file")
	flag.StringVar(&outputType, "type", "raw", "type of output")
	flag.StringVar(&targetName, "target", "username", "target name")
	flag.UintVar(&timeout, "timeout", 1000, "timeout in milliseconds")
	flag.UintVar(&retries, "retries", 3, "number of retries")
	flag.Parse()

	if maskFile == "" {
		fmt.Println("Error: mask file path cannot be empty")
		flag.Usage()
		os.Exit(1)
	}

	// * Only replace "..." at the end of the string.
	maskFile = strings.TrimSuffix(maskFile, "...")
}

// ? searchMasksRecursive walks through all directories and collects mask files.
func searchMasksRecursive(dir string) error {
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// ? Check if it's a file and contains "_mask.".
		if !d.IsDir() && strings.Contains(d.Name(), "_mask.") {
			masks = append(masks, path)
		}

		return nil
	})
}

func main() {
	ParseFlags()

	// ? Parsing mask.
	if maskFile[len(maskFile)-1:] == "/" {
		// * Recursively search for mask files if "..." is in the path.
		err := searchMasksRecursive(maskFile)
		if err != nil {
			fmt.Printf("Error searching mask files: %v\n", err)
			os.Exit(1)
		}

		if len(masks) == 0 {
			fmt.Println("No mask files found")
			os.Exit(1)
		}
	} else {
		// * If it's a single file, just add it to masks.
		masks = append(masks, maskFile)
	}

	parsedMasks := make([]*parser.ParsedMask, len(masks))
	// * Output the found masks.
	for idx, mask := range masks {
		pm := new(parser.ParsedMask)

		// * Unmarshalling the mask.
		data, err := os.ReadFile(mask)
		if err != nil {
			fmt.Printf("Error reading mask file %s: %v\n", mask, err)
			continue
		}
		if err := pm.UnmarshalJSON(data); err != nil {
			fmt.Printf("Error unmarshalling mask %s: %v\n", mask, err)
			continue
		}
		parsedMasks[idx] = pm
	}

	// ? Create the http client.
	httpclient := httpclient.NewHTTPClient(parsedMasks, parser.NewOutputForm(parser.Watchable), targetName, int(retries), time.Duration(timeout)*time.Millisecond)
	fmt.Println(httpclient.Do())
}
