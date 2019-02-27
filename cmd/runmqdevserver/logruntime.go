/*
© Copyright IBM Corporation 2017, 2019

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
package main

import (
	"runtime"
	"strings"

	containerruntime "github.com/ibm-messaging/mq-container/internal/containerruntime"
	"github.com/ibm-messaging/mq-container/internal/user"
)

func logContainerDetails() {
	log.Printf("CPU architecture: %v", runtime.GOARCH)
	kv, err := containerruntime.GetKernelVersion()
	if err == nil {
		log.Printf("Linux kernel version: %v", kv)
	}
	cr, err := containerruntime.GetContainerRuntime()
	if err == nil {
		log.Printf("Container runtime: %v", cr)
	}
	bi, err := containerruntime.GetBaseImage()
	if err == nil {
		log.Printf("Base image: %v", bi)
	}
	u, err := user.GetUser()
	if err == nil {
		if len(u.SupplementalGID) == 0 {
			log.Printf("Running as user ID %v (%v) with primary group %v", u.UID, u.Name, u.PrimaryGID)
		} else {
			log.Printf("Running as user ID %v (%v) with primary group %v, and supplementary groups %v", u.UID, u.Name, u.PrimaryGID, strings.Join(u.SupplementalGID, ","))
		}
	}
	caps, err := containerruntime.GetCapabilities()
	capLogged := false
	if err == nil {
		for k, v := range caps {
			if len(v) > 0 {
				log.Printf("Capabilities (%s set): %v", strings.ToLower(k), strings.Join(v, ","))
				capLogged = true
			}
		}
		if !capLogged {
			log.Print("Capabilities: none")
		}
	} else {
		log.Errorf("Error getting capabilities: %v", err)
	}
	sc, err := containerruntime.GetSeccomp()
	if err == nil {
		log.Printf("seccomp enforcing mode: %v", sc)
	}
	log.Printf("Process security attributes: %v", containerruntime.GetSecurityAttributes())
}
