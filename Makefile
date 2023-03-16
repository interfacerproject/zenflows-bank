all: build

generate:
	cd fabcoin && rm -f Fabcoin.go && abigen --bin=Fabcoin.bin --abi=Fabcoin.abi --pkg=fabcoin --out=Fabcoin.go && cd -

build:
	go build -o bank .

