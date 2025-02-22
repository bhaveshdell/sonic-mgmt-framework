#######################################################################
#
# Copyright 2019 Broadcom. All rights reserved.
# The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
#
#######################################################################

TOPDIR := ..
ABS_TOPDIR := $(abspath $(TOPDIR))

BUILD_DIR := $(TOPDIR)/build
CODEGEN_TOOLS_DIR := $(TOPDIR)/tools/swagger_codegen

CODEGEN_VER	:= 2.4.5
CODEGEN_JAR := $(CODEGEN_TOOLS_DIR)/swagger-codegen-cli-$(CODEGEN_VER).jar

SERVER_BUILD_DIR	:= $(BUILD_DIR)/rest_server
SERVER_CODEGEN_DIR	:= $(SERVER_BUILD_DIR)/codegen
SERVER_DIST_DIR		:= $(SERVER_BUILD_DIR)/dist
SERVER_DIST_INIT	:= $(SERVER_DIST_DIR)/.init_done
SERVER_DIST_GO		:= $(SERVER_DIST_DIR)/src/swagger
SERVER_DIST_UI		:= $(SERVER_DIST_DIR)/ui
SERVER_DIST_UI_HOME	:= $(SERVER_DIST_DIR)/ui/index.html

YANGAPI_DIR     := $(TOPDIR)/build/yaml
YANGAPI_SPECS   := $(shell find $(YANGAPI_DIR) -name '*.yaml' | sort)
YANGAPI_NAMES   := $(basename $(notdir $(YANGAPI_SPECS)))
YANGAPI_SERVERS := $(addsuffix /.yangapi_done, $(addprefix $(SERVER_CODEGEN_DIR)/, $(YANGAPI_NAMES)))

OPENAPI_DIR   := openapi
OPENAPI_SPECS := $(shell find $(OPENAPI_DIR) -name '*.yaml' | sort)
OPENAPI_NAMES := $(basename $(notdir $(OPENAPI_SPECS)))
OPENAPI_SERVERS := $(addsuffix /.openapi_done, $(addprefix $(SERVER_CODEGEN_DIR)/, $(OPENAPI_NAMES)))

PY_CLIENT_CODEGEN_DIR  := $(BUILD_DIR)/swagger_client_py
PY_CLIENT_TARGETS := $(addsuffix .yangapi_client, $(addprefix $(PY_CLIENT_CODEGEN_DIR)/, $(YANGAPI_NAMES))) \
                     $(addsuffix .openapi_client, $(addprefix $(PY_CLIENT_CODEGEN_DIR)/, $(OPENAPI_NAMES)))

UIGEN_DIR  = $(TOPDIR)/tools/ui_gen
UIGEN_SRCS = $(shell find $(UIGEN_DIR) -type f)


.PHONY: all clean cleanall go-server py-client

all: go-server py-client

go-server: $(YANGAPI_SERVERS) $(OPENAPI_SERVERS) $(SERVER_DIST_INIT) $(SERVER_DIST_UI_HOME)

$(SERVER_DIST_UI_HOME): $(YANGAPI_SERVERS) $(OPENAPI_SERVERS) $(UIGEN_SRCS)
	@echo "+++ Generating  landing page for Swagger UI +++"
	$(UIGEN_DIR)/src/uigen.py

py-client: $(PY_CLIENT_TARGETS)


.SECONDEXPANSION:

#======================================================================
# Common rule for directories. Use "." suffix, like "xyz/."
#======================================================================
.PRECIOUS: %/. 
%/.:
	mkdir -p $@

#======================================================================
# Download swagger codegen jar from Maven.org repo. It will be saved as
# build/swagger-codegen-cli.jar file.
#======================================================================
$(CODEGEN_JAR): | $$(@D)/.
	cd $(@D) && \
	wget http://central.maven.org/maven2/io/swagger/swagger-codegen-cli/$(CODEGEN_VER)/$(@F)

#======================================================================
# Generate swagger server in GO language for Yang generated OpenAPIs
# specs.
#======================================================================
%/.yangapi_done: $(YANGAPI_DIR)/$$(*F).yaml | $$(@D)/. $(CODEGEN_JAR) $(SERVER_DIST_INIT)
	@echo "+++ Generating GO server for Yang API $$(basename $(@D)).yaml +++"
	java -jar $(CODEGEN_JAR) generate \
		--lang go-server \
		--input-spec $(YANGAPI_DIR)/$$(basename $(@D)).yaml \
		--template-dir $(CODEGEN_TOOLS_DIR)/go-server/templates-yang \
		--output $(@D)
	cp $(@D)/go/api_* $(SERVER_DIST_GO)/
	cp $(@D)/go/routers.go $(SERVER_DIST_GO)/routers_$$(basename $(@D)).go
	cp $(@D)/api/swagger.yaml $(SERVER_DIST_UI)/$$(basename $(@D)).yaml
	touch $@

#======================================================================
# Generate swagger server in GO language for handcoded OpenAPI specs
#======================================================================
%/.openapi_done: $(OPENAPI_DIR)/$$(*F).yaml | $$(@D)/. $(CODEGEN_JAR) $(SERVER_DIST_INIT)
	@echo "+++ Generating GO server for OpenAPI $$(basename $(@D)).yaml +++"
	java -jar $(CODEGEN_JAR) generate \
		--lang go-server \
		--input-spec $(OPENAPI_DIR)/$$(basename $(@D)).yaml \
		--template-dir $(CODEGEN_TOOLS_DIR)/go-server/templates-nonyang \
		--output $(@D)
	cp $(@D)/go/api_* $(@D)/go/model_* $(SERVER_DIST_GO)/
	cp $(@D)/go/routers.go $(SERVER_DIST_GO)/routers_$$(basename $(@D)).go
	cp $(@D)/api/swagger.yaml $(SERVER_DIST_UI)/$$(basename $(@D)).yaml
	touch $@

#======================================================================
# Initialize dist directory for GO server code
#======================================================================
$(SERVER_DIST_INIT): | $$(@D)/.
	cp -r $(CODEGEN_TOOLS_DIR)/ui-dist $(@D)/ui
	cp -r $(CODEGEN_TOOLS_DIR)/go-server/src $(@D)/
	touch $@

#======================================================================
# Generate swagger client in Python for yang generated OpenAPI specs
#======================================================================
%.yangapi_client: $(YANGAPI_DIR)/$$(*F).yaml | $(CODEGEN_JAR) $$(@D)/.
	@echo "+++++ Generating Python client for $(*F).yaml +++++"
	java -jar $(CODEGEN_JAR) generate \
		-DpackageName=$(subst -,_,$(*F))_client \
		--lang python \
		--input-spec $(YANGAPI_DIR)/$(*F).yaml \
		--output $(@D)
	touch $@

#======================================================================
# Generate swagger client in Python for handcoded OpenAPI specs
#======================================================================
%.openapi_client: $(OPENAPI_DIR)/$$(*F).yaml | $(CODEGEN_JAR) $$(@D)/.
	@echo "+++++ Generating Python client for $(*F).yaml +++++"
	java -jar $(CODEGEN_JAR) generate \
		-DpackageName=$(subst -,_,$(*F))_client \
		--lang python \
		--input-spec $(OPENAPI_DIR)/$(*F).yaml \
		--output $(@D)
	touch $@

#======================================================================
# Cleanups
#======================================================================

clean-server:
	rm -rf $(SERVER_DIST_DIR)
	rm -rf $(SERVER_CODEGEN_DIR)

clean-client:
	rm -rf $(PY_CLIENT_CODEGEN_DIR)

clean: clean-server clean-client
	make -C yang clean

cleanall: clean
	rm -f $(CODEGEN_JAR)

