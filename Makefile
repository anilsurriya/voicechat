build: clean
	@go build -o output/voicechat

clean:
	@rm -f output/voicechat

run: build
	@output/voicechat