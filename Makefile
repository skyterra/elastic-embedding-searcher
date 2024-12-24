.PHONY: target local run modelx pb

# build linux binary.
target:
	export CGO_ENABLED=0 && \
	export GOOS=linux && \
	export GOARCH=amd64 && \
	go build -ldflags '-w -s' -o output/bin/searcher

# build local binary.
local:
	cp -r modelx output/ && \
	go build -ldflags '-w -s' -o output/bin/searcher

# execute program in local.
run: local
	cd output && ./bin/searcher -e http://127.0.0.1:9200 -m ./local_models/paraphrase-multilingual-MiniLM-L12-v2

modelx:
	rm -rf output/modelx && \
	cp -r modelx output/ && \
	(cd output && python ./modelx/server.py --model_path ./local_models/paraphrase-multilingual-MiniLM-L12-v2)

pb:	
	python -m grpc_tools.protoc --python_out=modelx --grpc_python_out=modelx --pyi_out=modelx -I ./pb ./pb/modelx.proto && \
    protoc -I=./pb --go_out=./pb ./pb/*.proto && \
    protoc -I=./pb --go-grpc_out=./pb ./pb/*.proto

deps:
	@echo "installing dependencies..."
	@which go > /dev/null || (echo "Go is not installed. Please install Go from https://golang.org/dl/" && exit 1)
	@which python > /dev/null || (echo "Python is not installed. Please install Python from https://www.python.org/downloads/" && exit 1)
	@which protoc > /dev/null || (echo "protoc is not installed. Please install Protoc from https://github.com/protocolbuffers/protobuf/releases" && exit 1)
	@python -m pip install --upgrade pip
	@pip install grpcio-tools
	@pip install -r requirements.txt
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "all dependencies installed."

# TODO make an example for fine-tuning in future.
ft:
	python ./modelx/fine_tuning.py --dataset ./dataset/example.csv --model ./local_models/paraphrase-multilingual-MiniLM-L12-v2 --version v1

# TODO write a command to download model from hugging face.
# Download the model from Hugging Face and save it to the ./output/local_models directory to save startup time and avoid the modelx timeout (10s).
