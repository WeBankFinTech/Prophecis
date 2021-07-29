# Prophecis-DI 



DI is the model training execution subsystem of the Prophecis system. It supports users to choose a stand-alone or distributed method to train the model. Users can access the CPU and GPU resources required for training through the UI. The DI subsystem will allocate according to the information configured by the user. Corresponding resources. 



#### Module

- CLI: CLI is a command line tool, which can perform basic command operations on Prophecis-DI.
- Jonmonitor: Jonmonitor is responsible for monitoring the status of task execution.
- LCM: LCM is responsible for various methods of task life cycle management, including start, pause, and terminate. 

- Restapi: Restapi is mainly used to process external service requests and service authentication, and supports load balancing.
- Storage: Storage is responsible for the operation of storage module, such as Minio, ES, Mongo, etc.
- Trainer: Trainer is responsible for managing model training tasks.

