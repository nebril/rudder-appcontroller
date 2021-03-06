/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package releaseutil

import (
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
)

// SimpleHead defines what the structure of the head of a manifest file
type SimpleHead struct {
	Version  string `json:"apiVersion"`
	Kind     string `json:"kind,omitempty"`
	Metadata *struct {
		Name        string            `json:"name"`
		Annotations map[string]string `json:"annotations"`
	} `json:"metadata,omitempty"`
}

// SplitManifests takes a string of manifest and returns a map contains individual manifests
func SplitManifests(bigfile string) map[string]string {
	// This is not the best way of doing things, but it's how k8s itself does it.
	// Basically, we're quickly splitting a stream of YAML documents into an
	// array of YAML docs. In the current implementation, the file name is just
	// a place holder, and doesn't have any further meaning.
	sep := "\n---\n"
	cutset := " \n\t"
	tpl := "manifest-%d"
	res := map[string]string{}
	tmp := strings.Split(bigfile, sep)
	for i, d := range tmp {
		if len(strings.Trim(d, cutset)) > 0 {
			res[fmt.Sprintf(tpl, i)] = d
		}
	}
	return res
}

// Manifest reperestens a single manifest content with SimpleHead added for additional metadata
type Manifest struct {
	SimpleHead
	Content string
}

// SplitManifestsWithHeads
func SplitManifestsWithHeads(bigfile string) ([]Manifest, error) {
	raws := SplitManifests(bigfile)

	result := make([]Manifest, 0, len(raws))
	var err error

	for _, raw := range raws {
		var head SimpleHead
		err = yaml.Unmarshal([]byte(raw), &head)

		result = append(result, Manifest{
			Content:    raw,
			SimpleHead: head,
		})
	}
	return result, err
}
