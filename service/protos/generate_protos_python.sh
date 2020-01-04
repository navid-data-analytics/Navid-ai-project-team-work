python3 -m grpc_tools.protoc -I ./service/protos \
  --python_out=./service/gen/protos \
  --grpc_python_out=./service/gen/protos \
  ./service/protos/*.proto
