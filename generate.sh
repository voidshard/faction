#! /bin/bash
#
#go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
#go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Remove existing Go struct files
echo "Removing existing Go struct files"
rm -v pkg/structs/*.pb.go
rm -v pkg/structs/*.encoding.go

set -e

# Generate Go structs from proto files
echo "Generating Go structs"
set -x
protoc -I=./proto --go_out=./ --go-grpc_out=./ ./proto/*.proto
set +x

# Generate marshal, unmarshal, set error functions for core message types.
# These funcs extend the default proto generated funcs to be a bit more helpful
# in meeting some interfaces internally for generic functions.
#
# Because this is all autogenerated code no one needs to read we'll push it all
# into a single file
echo "Generating encoding funcs"
cat > pkg/structs/encoding.go <<EOF
// FILE IS AUTOGENERATED. DO NOT EDIT!
package structs

import (
    "google.golang.org/protobuf/encoding/protojson"

    "buf.build/go/protoyaml"
)
EOF
for file in $(find ./proto -type f -name "*.proto"); do
    grep "message" $file | awk '{print $2}' | while read -r line; do
        lower=$(echo $line | tr '[:upper:]' '[:lower:]')
        cat >> pkg/structs/register_kinds.go <<EOF
Register.Add($line{})
EOF

        cat >> pkg/structs/encoding.go <<EOF

func (x *$line) Kind() string {
    return "$line"
}

func (x *$line) New() *$line {
    return &$line{}
}

func (x *$line) NewSlice() []*$line {
    return []*$line{}
}

func (x *$line) MarshalYAML() ([]byte, error) {
    return protoyaml.Marshal(x)
}

func (x *$line) UnmarshalYAML(data []byte) error {
    return protoyaml.Unmarshal(data, x)
}

func (x *$line) MarshalJSON() ([]byte, error) {
    return protojson.Marshal(x)
}

func (x *$line) UnmarshalJSON(data []byte) error {
    return protojson.Unmarshal(data, x)
}
EOF
        if [[ $line == *"Response"* ]]; then 
            cat >> pkg/structs/encoding.go <<EOF

func (x *$line) SetError(err *Error) {
    x.Error = err
}
EOF
        fi
    echo "+ wrote "$lower" encoding funcs"
    done
done

# sets up bson tag for Id -> _id (for MongoDB)
for file in $(find ./pkg/structs -type f -name "*.pb.go"); do
    sed -i 's/json:"Id,/bson:"_id" &/g' $file 
done
