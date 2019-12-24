all: build bench

build:
	python3.7 generatetest.py
	go build find.go
	go build find2.go
	
bench:
	./find # Warmup
	time ./find
	time ./find
	time ./find
	time ./find
	time ./find
	time ./find2
	time ./find2
	time ./find2
	time ./find2
	time ./find2