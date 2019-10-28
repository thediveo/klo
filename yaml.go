// Copyright 2019 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package klo

import (
	"io"

	"sigs.k8s.io/yaml"
)

// YAMLPrinter prints values in JSON format.
type YAMLPrinter struct{}

// NewYAMLPrinter returns a printer for outputting values in YAML format.
func NewYAMLPrinter() (ValuePrinter, error) {
	return &YAMLPrinter{}, nil
}

func (p *YAMLPrinter) Fprint(w io.Writer, v interface{}) error {
	txt, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(txt)
	return err
}
