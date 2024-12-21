.PHONY: target local run pb test clean sample run-sample modelx goods

# build linux binary.
target:
	export CGO_ENABLED=0 && \
	export GOOS=linux && \
	export GOARCH=amd64 && \
	go build -ldflags '-w -s' -o output/bin/frigga-finder

# build local binary.
local:
	cp -r modelx output/ && \
	go build -ldflags '-w -s' -o output/bin/searcher

# execute program in local.
run: local
	cd output && ./bin/searcher serve -e 127.0.0.1:9200 -m paraphrase-multilingual-MiniLM-L12-v2

modelx:
	rm -rf output/modelx && \
	cp -r modelx output/ && \
	(cd output && python ./modelx/server.py --model_path paraphrase-multilingual-MiniLM-L12-v2)

run-sample:
	export K8S_SERVICE=clotho-frigga-finder-kr-local && \
	cp -r csv output/ && \
	./output/bin/sample

pb:	
	python -m grpc_tools.protoc --python_out=modelx --grpc_python_out=modelx --pyi_out=modelx -I ./pb ./pb/modelx.proto && \
    protoc -I=./pb --go_out=./pb ./pb/*.proto && \
    protoc -I=./pb --go-grpc_out=./pb ./pb/*.proto

ft:
	python ./modelx/fine_tuning.py --dataset ./dataset/example.csv --model ./local_models/paraphrase-multilingual-MiniLM-L12-v2 --version v1

# convert action format data into importable action library data.
# e.g.
# make goods input=./goods/action_v1118.csv output=./goods/action_v1118_a.csv
goods:
	@awk 'BEGIN { FS=OFS=","; print "id,asset_id,annotation,loop,animation" } \
	     NR>1 { print NR-1, $$1, $$2, $$3, $$4 }' $(input) > $(output)

test:
	go test ./...

clean:
	rm -f ./output/bin/* && \
	rm -f ./output/logs/*
