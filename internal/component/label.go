/*
Copyright 2020 The Laputa Cloud Co.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package component

import "github.com/laputacloudco/sevendays-operator/api/v1alpha1"

func standardLabels(sd v1alpha1.SevenDays) map[string]string {
	return map[string]string{
		"app":      "sevendays",
		"instance": sd.Name,
		"managed":  "true",
		"owner":    sd.Name,
	}
}
