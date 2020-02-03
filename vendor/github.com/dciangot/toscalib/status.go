/*
Copyright 2015 - Olivier Wulveryck

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

package toscalib

// Status is used in the PropertyDefinition
type Status string

// Valid values for Status as described in Appendix 5.7.3
const (
	Supported    Status = "supported"
	Unsupported  Status = "unsupported"
	Experimental Status = "experimental"
	Deprecated   Status = "deprecated"
)
