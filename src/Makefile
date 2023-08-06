# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Main program paths
CONTROLLER_SRC=controller/controller.go controller/client_conn.go controller/storage_conn.go controller/resourceman_conn.go
CONTROLLER_BIN=controllerExec

RESOURCE_MANAGER_SRC=resource_manager/resource_manager.go resource_manager/client_conn.go resource_manager/controller_conn.go resource_manager/nodeman_conn.go
RESOURCE_MANAGER_BIN=resource_managerExec

CLIENT_SRC=client/client_main.go client/dispatch.go client/client.go client/fetch.go client/resourceman_handler.go
CLIENT_BIN=clientExec

SPAWN_SRC=spawn/spawn.go
SPAWN_BIN=spawnExec

NODEMAN_SRC=node_manager/node_manager.go node_manager/resourceman_conn.go node_manager/nodeman_conn.go
NODEMAN_BIN=node_managerExec


all: clean build

build:
	$(GOBUILD) -o $(CONTROLLER_BIN) $(CONTROLLER_SRC)
	$(GOBUILD) -o $(RESOURCE_MANAGER_BIN) $(RESOURCE_MANAGER_SRC)
	$(GOBUILD) -o $(CLIENT_BIN) $(CLIENT_SRC)
	$(GOBUILD) -o $(SPAWN_BIN) $(SPAWN_SRC)
	$(GOBUILD) -o $(NODEMAN_BIN) $(NODEMAN_SRC)


clean:
	$(GOCLEAN)
	rm -f $(CONTROLLER_BIN) $(RESOURCE_MANAGER_BIN) $(CLIENT_BIN) $(SPAWN_BIN) $(NODEMAN_BIN)

run_controller:
	$(GOBUILD) -o $(CONTROLLER_BIN) $(CONTROLLER_SRC)
	./$(CONTROLLER_BIN) 23000 23003 23004

run_resource_manager:
	$(GOBUILD) -o $(RESOURCE_MANAGER_BIN) $(RESOURCE_MANAGER_SRC)
	./$(RESOURCE_MANAGER_BIN) orion01:23004 23009 23010

run_client:
	$(GOBUILD) -o $(CLIENT_BIN) $(CLIENT_SRC)
	./$(CLIENT_BIN)

run_spawn:
	$(GOBUILD) -o $(SPAWN_BIN) $(SPAWN_SRC)
	./$(SPAWN_BIN)

run_nodeman:
	$(GOBUILD) -o $(NODEMAN_BIN) $(NODEMAN_SRC)
	./$(NODEMAN_BIN)