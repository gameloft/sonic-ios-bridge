/*
 *  Copyright (C) [SonicCloudOrg] Sonic Project
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
package cmd

import (
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	giDevice "github.com/electricbubble/gidevice"
	"os"

	"github.com/spf13/cobra"
)

// locationCmd represents the location command
var locationCmd = &cobra.Command{
	Use:   "location",
	Short: "Simulate location to your device.",
	Long:  "Simulate location to your device.",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		err := device.SimulateLocationUpdate(long, lat, giDevice.CoordinateSystemBD09)
		if err != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "location", err)
		}
		return nil
	},
}

var long, lat float64

func init() {
	rootCmd.AddCommand(locationCmd)
	locationCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	locationCmd.Flags().Float64Var(&long, "long", 0, "longitude")
	locationCmd.Flags().Float64Var(&lat, "lat", 0, "latitude")
	locationCmd.MarkFlagRequired("long")
	locationCmd.MarkFlagRequired("lat")
}
