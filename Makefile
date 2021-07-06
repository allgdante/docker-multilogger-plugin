PLUGIN_NAME = allgdante/docker-multilogger-plugin
PLUGIN_TAG := 0.0.2
PLUGIN_DIR = $(CURDIR)/.plugin

clean:
	@echo "### rm $(PLUGIN_DIR)"
	rm -rf $(PLUGIN_DIR) || true

build:
	@echo "### build"
	go build -ldflags '-extldflags "-fno-PIC -static"' -buildmode pie -tags 'osusergo netgo static_build'

docker:
	@echo "### docker build: rootfs image with docker-multilogger-plugin"
	docker build -t $(PLUGIN_NAME):rootfs .

rootfs:
	@echo "### create rootfs directory in $(PLUGIN_DIR)"
	mkdir -p $(PLUGIN_DIR)/rootfs
	docker create --name tmprootfs $(PLUGIN_NAME):rootfs
	docker export tmprootfs | tar -x -C $(PLUGIN_DIR)/rootfs
	@echo "### copy config.json to plugin directory"
	cp config.json $(PLUGIN_DIR)/
	docker rm -vf tmprootfs

create:
	@echo "### remove existing plugin $(PLUGIN_NAME):$(PLUGIN_TAG) if exists"
	docker plugin rm -f $(PLUGIN_NAME):$(PLUGIN_TAG) || true
	@echo "### create new plugin $(PLUGIN_NAME):$(PLUGIN_TAG) from plugin directory"
	docker plugin create $(PLUGIN_NAME):$(PLUGIN_TAG) $(PLUGIN_DIR)

enable:
	@echo "### enable plugin $(PLUGIN_NAME):$(PLUGIN_TAG)"
	docker plugin enable $(PLUGIN_NAME):$(PLUGIN_TAG)

push: clean docker rootfs create enable
	@echo "### push plugin $(PLUGIN_NAME):$(PLUGIN_TAG)"
	docker plugin push $(PLUGIN_NAME):$(PLUGIN_TAG)
