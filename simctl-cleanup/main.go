package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// RunCommandReturnCombinedStdoutAndStderr ...
func RunCommandReturnCombinedStdoutAndStderr(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	outBytes, err := cmd.CombinedOutput()
	outStr := string(outBytes)
	return outStr, err
}

// SimInfo ...
type SimInfo struct {
	Name        string
	SimID       string
	Status      string
	StatusOther string
}

// OSVersionSimInfoPair ...
type OSVersionSimInfoPair struct {
	OSVersion     string
	SimulatorInfo SimInfo
}

// SimulatorsGroupedByIOSVersions ...
type SimulatorsGroupedByIOSVersions map[string][]SimInfo

// a simulator info line should look like this:
//  iPhone 5s (EA1C7E48-8137-428C-A0A5-B2C63FF276EB) (Shutdown)
// or
//  iPhone 4s (51B10EBD-C949-49F5-A38B-E658F41640FF) (Shutdown) (unavailable, runtime profile not found)
func getSimInfoFromLine(lineStr string) (SimInfo, error) {
	baseInfosExp := regexp.MustCompile(`(?P<deviceName>[a-zA-Z].*[a-zA-Z0-9 -]*) \((?P<simulatorID>[a-zA-Z0-9-]{36})\) \((?P<status>[a-zA-Z]*)\)`)
	baseInfosRes := baseInfosExp.FindStringSubmatch(lineStr)
	if baseInfosRes == nil {
		return SimInfo{}, fmt.Errorf("No match found")
	}

	simInfo := SimInfo{
		Name:   baseInfosRes[1],
		SimID:  baseInfosRes[2],
		Status: baseInfosRes[3],
	}

	// StatusOther
	restOfTheLine := lineStr[len(baseInfosRes[0]):]
	if len(restOfTheLine) > 0 {
		statusOtherExp := regexp.MustCompile(`\((?P<statusOther>[a-zA-Z ,]*)\)`)
		statusOtherRes := statusOtherExp.FindStringSubmatch(restOfTheLine)
		if statusOtherRes != nil {
			simInfo.StatusOther = statusOtherRes[1]
		}
	}
	return simInfo, nil
}

func collectAllSimIDs(simctlListOutputToScan string) SimulatorsGroupedByIOSVersions {
	simulatorsByIOSVersions := SimulatorsGroupedByIOSVersions{}
	currIOSVersion := ""

	fscanner := bufio.NewScanner(strings.NewReader(simctlListOutputToScan))
	isDevicesSectionFound := false
	for fscanner.Scan() {
		aLine := fscanner.Text()

		if aLine == "== Devices ==" {
			isDevicesSectionFound = true
			continue
		}

		if !isDevicesSectionFound {
			continue
		}
		if strings.HasPrefix(aLine, "==") {
			isDevicesSectionFound = false
			continue
		}
		if strings.HasPrefix(aLine, "--") {
			iosVersionSectionExp := regexp.MustCompile(`-- (?P<iosVersionSection>.*) --`)
			iosVersionSectionRes := iosVersionSectionExp.FindStringSubmatch(aLine)
			if iosVersionSectionRes != nil {
				currIOSVersion = iosVersionSectionRes[1]
			}
			continue
		}

		// fmt.Println("-> ", aLine)
		simInfo, err := getSimInfoFromLine(aLine)
		if err != nil {
			fmt.Println(" [!] Error scanning the line for Simulator info: ", err)
		}

		currIOSVersionSimList := simulatorsByIOSVersions[currIOSVersion]
		currIOSVersionSimList = append(currIOSVersionSimList, simInfo)
		simulatorsByIOSVersions[currIOSVersion] = currIOSVersionSimList
	}

	return simulatorsByIOSVersions
}

func (simsGrouped *SimulatorsGroupedByIOSVersions) flatList() []OSVersionSimInfoPair {
	osVersionSimInfoPairs := []OSVersionSimInfoPair{}

	for osVer, simulatorInfos := range *simsGrouped {
		for _, aSimInfo := range simulatorInfos {
			osVersionSimInfoPairs = append(osVersionSimInfoPairs, OSVersionSimInfoPair{
				OSVersion:     osVer,
				SimulatorInfo: aSimInfo,
			})
		}
	}

	return osVersionSimInfoPairs
}

func deleteIOSSimulator(simInfo SimInfo) error {
	_, err := RunCommandReturnCombinedStdoutAndStderr("xcrun", "simctl", "delete", simInfo.SimID)
	return err
}

func (simsGrouped *SimulatorsGroupedByIOSVersions) duplicates() []OSVersionSimInfoPair {
	duplicates := []OSVersionSimInfoPair{}
	for osVer, simulatorInfos := range *simsGrouped {
		simNameCache := map[string]bool{}
		for _, aSimInfo := range simulatorInfos {
			if _, isFound := simNameCache[aSimInfo.Name]; isFound {
				duplicates = append(duplicates, OSVersionSimInfoPair{
					OSVersion:     osVer,
					SimulatorInfo: aSimInfo,
				})
			}
			simNameCache[aSimInfo.Name] = true
		}
	}
	return duplicates
}

func (osVerSimInfoPair *OSVersionSimInfoPair) String() string {
	return fmt.Sprintf("[OS: %s] %#v", osVerSimInfoPair.OSVersion, osVerSimInfoPair.SimulatorInfo)
}

func main() {
	var (
		isHelp   = flag.Bool("help", false, `Show help`)
		isAll    = flag.Bool("all", false, `List all simulators, not just duplicates`)
		isDelete = flag.Bool("delete", false, `Delete the listed simulators - BE CAREFUL, if you use "--all" it'll delete all simulators!`)
		isIDOnly = flag.Bool("id-only", false, `Will only print the Simulator IDs, one ID per line. Can be used for piping to other tools.`)
	)

	flag.Usage = func() {
		fmt.Println(`# Usage:`)
		fmt.Println()
		fmt.Println(`By default running this tool will list all the duplicate Simulators it can find
by calling: $ xcrun simctl list

Use can use the --delete flag to delete these simulators (same as calling
for every listed id: $ xcrun simctl delete [SIM-ID])

For other flags see the whole list of available flags below.`)
		fmt.Println()
		fmt.Println("# Available parameters / flags:")
		fmt.Println()
		fmt.Printf("Usage: %s [FLAGS]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *isHelp {
		flag.Usage()
		os.Exit(0)
	}

	simctlListOut, err := RunCommandReturnCombinedStdoutAndStderr("xcrun", "simctl", "list")
	if err != nil {
		log.Println(" [!] Failed to get `xcrun simctl list`: ", err)
		log.Println()
		log.Println("Output was:")
		log.Fatalln(simctlListOut)
	}

	allSimIDsGroupedBySimVersion := collectAllSimIDs(simctlListOut)
	var simIDsToWorkWith []OSVersionSimInfoPair
	if !*isAll {
		simIDsToWorkWith = allSimIDsGroupedBySimVersion.duplicates()
	} else {
		simIDsToWorkWith = allSimIDsGroupedBySimVersion.flatList()
	}

	if !*isIDOnly {
		if *isAll {
			fmt.Println("All available Simulators:")
		} else {
			fmt.Println("Duplicated Simulators to remove:")
		}
		fmt.Println()
	}

	for _, itm := range simIDsToWorkWith {
		if *isIDOnly {
			fmt.Println(itm.SimulatorInfo.SimID)
			continue
		}
		fmt.Println("* ", itm.String())
		if *isDelete {
			fmt.Println(" -> deleting ...")
			if err := deleteIOSSimulator(itm.SimulatorInfo); err != nil {
				log.Fatalf("Failed to delete duplicate simulator: %s", err)
			}
			fmt.Println("    deleted [OK]")
		}
	}
}
