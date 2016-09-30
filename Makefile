BEATNAME=loxybeat
BEAT_DIR=github.com/mod/loxybeat
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS?=./vendor/github.com/elastic/beats
GOPACKAGES=$(shell glide novendor)
PREFIX?=.

APP_NAME = loxy
APP_VERSION = 0.1.0

# Path to the libbeat Makefile
-include $(ES_BEATS)/libbeat/scripts/Makefile

# Initial beat setup
.PHONY: setup
setup: copy-vendor
	make update

# Copy beats into vendor directory
.PHONY: copy-vendor
copy-vendor:
	mkdir -p vendor/github.com/elastic/
	cp -R ${GOPATH}/src/github.com/elastic/beats vendor/github.com/elastic/
	rm -rf vendor/github.com/elastic/beats/.git

.PHONY: git-init
git-init:
	git init
	git add README.md CONTRIBUTING.md
	git commit -m "Initial commit"
	git add LICENSE
	git commit -m "Add the LICENSE"
	git add .gitignore
	git commit -m "Add git settings"
	git add .
	git reset -- .travis.yml
	git commit -m "Add loxybeat"
	git add .travis.yml
	git commit -m "Add Travis CI"

# This is called by the beats packer before building starts
.PHONY: before-build
before-build:


.PHONY: app_build app_tag app_release

app_build:
	docker build -t $(APP_NAME):$(APP_VERSION) --rm .

app_tag:
	docker tag $(APP_NAME):$(APP_VERSION) $(APP_NAME):latest

app_release: test tag_latest
	docker push $(APP_NAME):$(APP_VERSION)
	docker push $(APP_NAME)

app_clean:
	docker rmi $(APP_NAME):$(APP_VERSION)