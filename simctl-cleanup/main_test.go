// example duplication:

// $ xcrun simctl list
// == Device Types ==
// iPhone 4s (com.apple.CoreSimulator.SimDeviceType.iPhone-4s)
// iPhone 5 (com.apple.CoreSimulator.SimDeviceType.iPhone-5)
// iPhone 5s (com.apple.CoreSimulator.SimDeviceType.iPhone-5s)
// iPhone 6 (com.apple.CoreSimulator.SimDeviceType.iPhone-6)
// iPhone 6 Plus (com.apple.CoreSimulator.SimDeviceType.iPhone-6-Plus)
// iPhone 6s (com.apple.CoreSimulator.SimDeviceType.iPhone-6s)
// iPhone 6s Plus (com.apple.CoreSimulator.SimDeviceType.iPhone-6s-Plus)
// iPad 2 (com.apple.CoreSimulator.SimDeviceType.iPad-2)
// iPad Retina (com.apple.CoreSimulator.SimDeviceType.iPad-Retina)
// iPad Air (com.apple.CoreSimulator.SimDeviceType.iPad-Air)
// iPad Air 2 (com.apple.CoreSimulator.SimDeviceType.iPad-Air-2)
// iPad Pro (com.apple.CoreSimulator.SimDeviceType.iPad-Pro)
// Apple TV 1080p (com.apple.CoreSimulator.SimDeviceType.Apple-TV-1080p)
// Apple Watch - 38mm (com.apple.CoreSimulator.SimDeviceType.Apple-Watch-38mm)
// Apple Watch - 42mm (com.apple.CoreSimulator.SimDeviceType.Apple-Watch-42mm)
// == Runtimes ==
// iOS 8.3 (8.3 - 12F70) (com.apple.CoreSimulator.SimRuntime.iOS-8-3)
// iOS 9.0 (9.0 - 13A344) (com.apple.CoreSimulator.SimRuntime.iOS-9-0)
// iOS 9.2 (9.2 - 13C5055d) (com.apple.CoreSimulator.SimRuntime.iOS-9-2)
// tvOS 9.0 (9.0 - 13T393) (com.apple.CoreSimulator.SimRuntime.tvOS-9-0)
// watchOS 2.0 (2.0 - 13S343) (com.apple.CoreSimulator.SimRuntime.watchOS-2-0)
// == Devices ==
// -- iOS 9.0 --
//     iPhone 4s (5D0C4373-D7C1-4854-9244-CEEE0A467777) (Shutdown)
//     iPhone 4s (A01958BB-6031-4927-AF4F-475B86EBBF1F) (Shutdown)
//     iPhone 4s (B47C3F19-C814-4420-8C86-190CDAC69F6D) (Shutdown)
//     iPhone 5 (7F265FF5-88FF-4A92-9489-D7358F4E6994) (Shutdown)
//     iPhone 5 (2A38DC06-18D6-4227-8430-2F5BDE8AD230) (Shutdown)
//     iPhone 5 (1F9B3122-FC7E-4D40-8E46-D353B075A8AC) (Shutdown)
// ...
// -- Unavailable: com.apple.CoreSimulator.SimRuntime.iOS-8-4 --
//     iPhone 4s (2A6C840F-1A30-481E-A4A8-201B0C48E099) (Shutdown) (unavailable, runtime profile not found)
//     iPhone 4s (55FFD329-916A-4C9D-A8B9-D51BF59BE184) (Shutdown) (unavailable, runtime profile not found)
//     iPhone 5 (EC4D7BC2-5BD7-4FE6-BE22-0FA516FD9C1D) (Shutdown) (unavailable, runtime profile not found)
//     iPhone 5 (DD53BF84-22E5-4865-BAFD-B8E75D023AE8) (Shutdown) (unavailable, runtime profile not found)
// ...

package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getSimInfoFromLine(t *testing.T) {
	{
		simInfo, err := getSimInfoFromLine("   iPhone 5s (EA1C7E48-8137-428C-A0A5-B2C63FF276EB) (Shutdown)")
		require.NoError(t, err)
		require.Equal(t, SimInfo{
			Name:        "iPhone 5s",
			SimID:       "EA1C7E48-8137-428C-A0A5-B2C63FF276EB",
			Status:      "Shutdown",
			StatusOther: "",
		}, simInfo)
	}
	{
		simInfo, err := getSimInfoFromLine("iPhone 4s (51B10EBD-C949-49F5-A38B-E658F41640FF) (Creating) (unavailable, runtime profile not found)")
		require.NoError(t, err)
		require.Equal(t, SimInfo{
			Name:        "iPhone 4s",
			SimID:       "51B10EBD-C949-49F5-A38B-E658F41640FF",
			Status:      "Creating",
			StatusOther: "unavailable, runtime profile not found",
		}, simInfo)
	}

	// no SimInfo
	{
		_, err := getSimInfoFromLine("-- iOS 9.0 --")
		require.Equal(t, errors.New("No match found"), err)
	}
	{
		_, err := getSimInfoFromLine("iPhone 5s (com.apple.CoreSimulator.SimDeviceType.iPhone-5s)")
		require.Equal(t, errors.New("No match found"), err)
	}
}
