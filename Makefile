all: build bench

build:
	python3.7 generatetest.py
	go build find.go
	go build find2.go
	go build find3.go
	
bench:
	./find # Warmup
	time ./find
	time ./find
	time ./find
	time ./find
	time ./find
	./find2
	time ./find2
	time ./find2
	time ./find2
	time ./find2
	time ./find2
	./find3
	time ./find3
	time ./find3
	time ./find3
	time ./find3
	time ./find3