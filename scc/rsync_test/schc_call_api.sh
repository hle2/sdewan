# /*
#  * Copyright 2020 Intel Corporation, Inc
#  *
#  * Licensed under the Apache License, Version 2.0 (the "License");
#  * you may not use this file except in compliance with the License.
#  * You may obtain a copy of the License at
#  *
#  *     http://www.apache.org/licenses/LICENSE-2.0
#  *
#  * Unless required by applicable law or agreed to in writing, software
#  * distributed under the License is distributed on an "AS IS" BASIS,
#  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  * See the License for the specific language governing permissions and
#  * limitations under the License.
#  */


name="test-name1"
description="test object description"
test_url="http://localhost:9015/v2/tests"


test_object_data="$(cat << EOF
{
 "metadata" : {
    "name": "${name}",
    "description": "${description}",
    "userData1":"<user data 1>",
    "userData2":"<user data 2>"
   },
 "spec" : {
    "ip" : "1.1.1.1",
    "port" : 9000,
    "myfield" : "myfield2"
  }
 }
}
EOF
)"

# Create 
printf "\n\nCreating test object data\n\n"
curl -d "${test_object_data}" -X POST ${test_url}

# Get logical cloud data
printf "\n\nGetting test object\n\n"
curl -X GET "${test_url}/${name}"

printf "\n\nGetting clusters info for logical cloud\n\n"
curl -X GET ${test_url}