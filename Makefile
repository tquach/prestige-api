MOD_NAME		:= github.com/tquach/prestige-api
REPO			?= tquach
APP_NAME		?= prestige-api
TAG				?= latest

all: $(APP_NAME)

$(APP_NAME): main.go test
	# go get -u github.com/tools/godep
	go build $(MOD_NAME)

dist:
	docker build -t $(REPO)/$(APP_NAME):$(TAG) .

start: dist
	docker run -it --rm --name $(APP_NAME) -p 9000:9000 $(REPO)/$(APP_NAME) $(APP_NAME) --configFile _config/settings/local.yml

test: 
	go test ./...

clean:
	@rm -f $(APP_NAME)

deploy: dist
	docker push $(REPO)/$(APP_NAME):$(TAG) 

.PHONY:
	start clean test
