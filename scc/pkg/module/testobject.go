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
type TestObject struct {
	Metadata ObjectMetaData `json:"metadata"`
	Specification TestObjectSpec `json:"spec"`
}

//TestObjectSpec contains the parameters
type TestObjectSpec struct {
	Ip        	string 	`json:"ip" validate:"required,ipv4"`
	Port 		int 	`json:"port" validate:"gte=0,lte=50000"`
	MyField		string  `json:"myfield" validate:"omitempty,check_field"`
}

func (c *TestObject) GetMetadata() ObjectMetaData {
	return c.Metadata
}