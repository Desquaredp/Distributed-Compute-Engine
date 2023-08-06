# Project 2: Distributed Computation Engine



## Building
To build the project, run the following command from the root of the project:

```make build```

This will build the project and place the binaries in the ```src``` directory.

## Running
Since there are many components in the project, each will have its own set of instructions on how to run it.

### Controller
To run the Controller, run the following command from the ```src``` directory:

```make run_controller```

Note: The controller has default ports it listens on and runs on. The Makefile has these ports baked in. If you want to change the ports, you will have to change them in the Makefile or pass them as arguments to the controller.

To run the controller using the binary, run the following command from the ```src``` directory:

```./controllerExec <Storage Nodes facing port port> <Client facing port> <Resource Manager facing port>```


### Resource Manager
To run the Resource Manager, run the following command from the ```src``` directory:

```make run_resource_manager```

Note: The Resource Manager has default ports it listens on and runs on. The Makefile has these ports baked in. If you want to change the ports, you will have to change them in the Makefile or pass them as arguments to the Resource Manager.

To run the Resource Manager using the binary, run the following command from the ```src``` directory:

```./resourceManagerExec <Controller host:port> <Node Manager Facing Port> <Client facing Port>```


### Node Manager
(Tentative)

To run the Node Manager, it is recommended to use the script provided in the ```scripts``` directory. To run the script, run the following command from the ```scripts``` directory:

```bash startNodeMan```

Note: The script is meant to run on orion. If you want to run it on a different machine, you will have to change the script.

The Node Manager listens on port 23098 for Node Manager to Node Manager communication, and port 23099 for Resource Manager to Node Manager communication. 
There is no way to change these ports, as they are hardcoded in the Node Manager.

### Storage Node
(Tentative)

To run the Storage Node, it is recommended to use the script provided in the ```scripts``` directory. To run the script, run the following command from the ```scripts``` directory:

```bash startStorageNodes```


Note: The script is meant to run on orion. If you want to run it on a different machine, you will have to change the script.


### Client
To run the Client, run the following command from the ```src``` directory:

#### To get list of commands:

```./clientExec -h```

#### To populate the DFS configuration file:

```./clientExec --populate-config dfs <config-file name>```

This will create a config file and populate it with the default values. You can then edit the config file to change the values.


#### To load the config file and run the client:

```./clientExec --load-config dfs <config-file name>```


#### To populate the MapReduce configuration file:

```./clientExec --populate-config mr <config-file name>```

This will create a config file and populate it with the default values. You can then edit the config file to change the values.

#### To load the config file and run a MapReduce job:

```./clientExec --load-config mr <config-file name>```

### To build the plugin:

```cd plug```

```mkdir <plugin name>```

```cd <plugin name>```

Create the plugin code in this directory. Make sure ```var MRPlugin``` is exported, and the plugin follows the interfaces and structs defined in ```src/types```.

```go build -buildmode=plugin -o <plugin name>.so <plugin name>.go```

NOTE: The plugin doesn't run if it is built on a different machine than the one it is being run on. This is because the plugin is built for a specific architecture. To fix this, you will have to build the plugin on the machine you want to run it on.


## Design

The design document can be found [here](./CommsDesign.md).

## Components
### Client
The client is responsible for sending a list of jobs to the execution engine, and sending files to the DFS.
### Controller
The Controller is the entry point of the DFS, and is responsible for indexing the files in the DFS. It also handles the communication between the client and the DFS.
### Resource Manager
The Resource Manager is the entry point of the execution engine. I thought of calling it the Execution Manager, but that sounds weird. It is responsible for assigning tasks to the nodes, and keeping track of the progress of the tasks.
### Storage Node
The Storage Node is responsible for storing the files in the DFS. It is also responsible for sending the files to the nodes that need them. It has to send periodic heartbeats to the Controller to let it know that it is still alive.
### Node Manager
The Node Manager is responsible for executing the tasks assigned to it by the Resource Manager. It sits on the same node, but doesn't have to send periodic heartbeats to the Resource Manager or the Controller. 


#### Changes I would make to the project:
- External sort. Although I implemented an external sort algorithm, having credits for it would be nice. Plus, it was a lot of fun to implement.
- Job chaining. It took me about 5 hours to implement job chaining, but it was a lot of fun. Cleaning up the previous job's leftovers was paramount to the success of job chaining.
- Implement the 'chaos monkey' to test the system as a whole. Would be pretty cool, but impractical.
