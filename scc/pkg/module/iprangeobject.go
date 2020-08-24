/*
 * Copyright 2020 Intel Corporation, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governinog permissions and
 * limitations under the License.
 */

package module

// App contains metadata for Apps
type IPRangeObject struct {
	Metadata ObjectMetaData `json:"metadata"`
	Specification IPRangeObjectSpec `json:"spec"`
}

//IPRangeObjectSpec contains the parameters
type IPRangeObjectSpec struct {
	Subnet    	string 	`json:"subnet" validate:"required"`
	MinIp    	string 	`json:"minIp" validate:"ipv4"`
	Maxip		string 	`json:"maxIp" validate:"ipv4"`
}

func (c *IPRangeObject) GetMetadata() ObjectMetaData {
	return c.Metadata
}