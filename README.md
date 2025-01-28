# Map-Reduce

This repository contains an implementation of the **Map-Reduce** algorithm in Go. Map-Reduce is a programming model used for processing and generating large datasets that can be parallelized across a distributed cluster of computers.

## Features
- **Map Function**: Applies a function to each element of a dataset and outputs a collection of intermediate key-value pairs.
- **Reduce Function**: Processes the intermediate key-value pairs and combines them into a final result.
- **Fault Tolerance**: Designed to handle failures during the execution of tasks.
- **Parallel Processing**: Can process large datasets by distributing tasks across multiple workers.

## File Overview
- `coordinator.go`: Manages the coordination of the Map-Reduce job, including task distribution and communication between workers.
- `rpc.go`: Contains the remote procedure calls (RPC) that facilitate communication between the coordinator and workers.
- `worker.go`: Implements the worker functionality to execute the Map and Reduce tasks.
- `README.md`: Documentation of the repository.
