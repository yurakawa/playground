- インストール

brew install protobuf


- Go言語でのコード生成用プラグインを入れておきます
go get -u github.com/golang/protobuf/protoc-gen-go

- protocコマンド実行

protoc \
    -Iproto \
    --go_out=plugins=grpc:api \
    proto/*.proto
    
### Usage

- 起動

`go run server.go`


### CLIツールで動作確認する
- grpc-cliのインストール

brew tap grpc/grpc
brew install grpc


### CLIツールを使う

- サービス確認
grpc_cli ls localhost:50051 pancake.baker.PancakeBakerService -l filename: grpc_reflection_v1alpha/reflection.proto

- RPCを実行
grpc_cli call localhost:50051 pancake.baker.PancakeBakerService/Bake 'menu: 1'
