validate_version:
ifndef VERSION
	$(error VERSION is undefined)
endif

test:
	go test ./... -v -cover

cover:
	go test -v -coverpkg=./... -coverprofile=profile.cov ./... -count=1
	go tool cover -func profile.cov
	# go tool cover -html coverage.out

release: validate_version
	# linux
	GOOS=linux go build -ldflags "-X main.version=${VERSION}" -o di ;\
	tar -zcvf ./releases/di_${VERSION}_linux.tar.gz ./di ;\

	# macos (arm)
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=${VERSION}" -o di ;\
	tar -zcvf ./releases/di_${VERSION}_macos_arm64.tar.gz ./di ;\

	# macos (amd)
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o di ;\
	tar -zcvf ./releases/di_${VERSION}_macos_amd64.tar.gz ./di ;\

	# windows
	GOOS=windows go build -ldflags "-X main.version=${VERSION}" -o di ;\
	tar -zcvf ./releases/di_${VERSION}_windows.tar.gz ./di ;\

	rm ./di