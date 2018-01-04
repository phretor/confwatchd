.PHONY: build fmt lint run test vet deps install

SRC_PATH=.
TARGET=confwatchd

default: build

build: deps
	@go build $(FLAGS) -o $(TARGET) $(SRC_PATH)

vet:
	@go vet $(SRC_PATH)

fmt:
	@go fmt $(SRC_PATH)/...

lint:
	@golint $(SRC_PATH)

test:
	@go test $(SRC_PATH)/...

clean:
	@rm -rf $(TARGET)

deps:
	@go get github.com/gin-gonic/gin
	@go get github.com/jinzhu/gorm
	@go get github.com/jinzhu/gorm/dialects/sqlite
	@go get gopkg.in/unrolled/secure.v1
	@go get github.com/gosimple/slug
	@go get github.com/gin-gonic/autotls
	@go get github.com/michelloworld/ez-gin-template
	@go get github.com/pariz/gountries


