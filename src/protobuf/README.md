# protobuf tutorial

## Build steps
* `mkdir addressbook`
* `protoc -I=./proto --go_out=addressbook ./proto/addressbook.proto`
* `cd add-people`
* `go build`
* `cd list-people`
* `go build`

## Run steps
* `./add-people/add-people addressbook.data`
* `./list-people/list-people addressbook.data`
