build:
	mkdir -p out && cd out && \
	go build -o service.hello ../cmd/*
