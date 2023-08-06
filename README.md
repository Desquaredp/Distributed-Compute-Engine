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


## Retrospective Questions
### How many extra days did you use for the project?
The project on 11th of April, and will be done by the 30th of April. So, 19 days.
### Given the same goals, how would you complete the project differently if you didn’t have any restrictions imposed by the instructor? This could involve using a particular library, programming language, etc. Be sure to provide sufficient detail to justify your response.

- I think using gRPC for a small part of the project would be a great way to learn gRPC.
- I would have also like liked to use a Key-Value store and Zookeeper to store the metadata of the system. This would have been a great way to learn about Key-Value stores and Zookeeper.
### Let’s imagine that your next project was to improve and extend P2. What are the features/functionality you would add, use cases you would support, etc? Are there any weaknesses in your current implementation that you would like to improve upon? This should include at least three areas you would improve/extend.

#### Features I would add:
- Sort by Value
- Restarting a failed task
- Resource Manager can timeout a task if it takes too long
There's a lot more I would like to add as an extension to this project.


#### Changes I would make to rectify the weaknesses in my current implementation:
##### Controlled parallelism
  I resorted to letting the network decide how many tasks run in parallel. This is not ideal, as there are instances where the program might have several MAP tasks running in parallel on the
  same node. This might result in the node running out of memory, and crashing. I would like to control the parallelism of the tasks, so that I can control how many tasks run in parallel on a node.

##### Better error handling
  There are several instances where I do not handle errors. I would like to handle errors better, and not just print them out to the console.

##### MAP tasks are designed to load the entire processed output in memory
  This is obviously not ideal, and pretty dumb. The default fragment size is capped at 128MBs and each fragment has its independent output which, if loaded in memory, will not exceed 128MBs. This should be fine. However, combine this
  with the network defined parallelism(defined in the first point), and one can see how this can easily overload a node if it houses sever fragments.

There are a couple of ways to fix this:
1. Instead of relying on the network to decide how many tasks run in parallel, I can control the parallelism. This way, I can control how many MAP tasks run in parallel on a node, and thus control how much memory is used.
   I could do that by giving each node manager all the MAP tasks it is responsible for, and letting it decide how many tasks run in parallel. This way, I would not overload a node.
2. I can limit the amount of data ran at a time, in memory, per fragment. Sort the fragment of the fragment processed in memory and write it to disk. This way, I would not overload a node. Although, this could potentially overload a node because
   of the uncontrolled parallelism. There can be several MAP tasks running in parallel on a node, and each task can be writing to disk. This can overload a node. To fix this, I can limit the amount of data written to disk at a time.

##### Better testing
  One can never have enough tests. I would like to write more tests to test the different components of the system. I would also like to write tests to test the system as a whole. Implementing the 'chaos monkey' would be a great way to test the system as a whole.

##### Code cleanliness
  Simply put, the code is awful. I would like to clean up the code, and make it more readable. I would also like to add more comments to the code.


### Give a rough estimate of how long you spent completing this assignment. Additionally, what part of the assignment took the most time?
1. Getting the plugin to work took me ages. The issue was the '-race' flag, which was almost impossible to deduce from the error messages. I spent a lot of time trying to figure out why the plugin was not working. I would say I spent 2 days trying to figure out why the plugin was not working.
2. Although it wasn't a part of the project, I implemented an external sort algorithm using priority queues. This took me a while to implement. I would say I spent a day implementing the external sort algorithm. I really enjoyed implementing the external sort algorithm!

### What did you learn from completing this project? Is there anything you would change about the project?
#### I learned a lot about go and distributed systems. To be specific I learnt about the following:
##### Go plugins
Go plugins are a pain to work w. It took me a while to figure out how to get them to work. This was mainly because if they are ran with the '-race' flag, they do not work as expected.
The reason eludes me. I tried to understand how the '-race' flag affects the plugin, but I could not find any information on it. My guess is since the '-race' creates something called the shadow memory(cool name) space, it might be causing the plugin to not work as expected. 
##### Channels
Channels are a great way to collect the progress of all the MAP, SHUFFLE, and REDUCE tasks. I used channels to collect the progress of the tasks, and to send the progress to the client.
Plus, they are a great way to pause the execution of the program until all the tasks are done.
##### ...

#### Changes I would make to the project:
- External sort. Although I implemented an external sort algorithm, having credits for it would be nice. Plus, it was a lot of fun to implement.
- Job chaining. It took me about 5 hours to implement job chaining, but it was a lot of fun. Cleaning up the previous job's leftovers was paramount to the success of job chaining.
- Implement the 'chaos monkey' to test the system as a whole. Would be pretty cool, but impractical.
