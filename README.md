# COLMENA Agent Repository

This GitHub repository contains all the files and software necessary to deploy a COLMENA agent. COLMENA (COLaboración entre dispositivos Mediante tecnología de ENjAmbre) aims to ease the development, deployment, operation and maintenance of extremely-high available, reliable and intelligent services running seamlessly across the device-edge-cloud continuum. It leverages a swarm approach organising a dynamic group of autonomous, collaborative nodes following an agile, fully-decentralised, robust, secure and trustworthy open architecture.

## Table of Contents
- [Repository Structure](#repository-structure)
- [Getting Started](#getting-started)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgments](#acknowledgments)



## Repository Structure
The repository is organized into the following directories and files:
### Directories
- **agent**: Contains all the code necessary to deploy a COLMENA Agent.
- **dcp**: Contains all the necessary code to deploy a Centralized Colmena Platform to coordinate the COLMENA Agents.
- **library**: Contains all the necessary code to test the code.
### Files
- **.gitignore**: Specifies files and directories to be ignored by Git.
- **CODE_OF_CONDUCT.md**: Outlines the expected behavior and guidelines for participants within the project's community. 
- **CONTRIBUTING.md**: Overview of the repository, setup instructions, and basic usage examples.
- **LICENSE**: License information for the repository.
- **README.md**: Overview of the repository, setup instructions, and basic usage examples.
- **setting.gradle**: Configuration of gradle to build the project


## Getting Started
To get started with deploying a COLMENA agents, follow these steps:
1. Build the DCP and the agent
    ```bash
    gradle clean assemble
    ```
    The project can also be build using gradle's docker image:
    ```bash
    docker run --rm -v .:/home/gradle gradle:8.9.0-jdk17 gradle clean assemble
    ```
2. Starting DCP
    In this initial prototype, COLMENA agents require a centralized version of the Colmena platform. Before deploying the Agents, the DCP (Distributed Colmena Platform) needs to be deployed. The following command deploys the DCP on the local node on port 5555:
    ```bash
    java -classpath <path_to_agent_repository_root>/dcp/build/libs/dcp-0.0.1-all.jar es.bsc.colmena.dcp.Application
    ```

3. Starting the Agents    
    Once DCP is running, then Agents can be started. Since they do not have neighbor discovery implemented, they need an endpoint to enter the Colmena platform. Likewise, Colmena does not provide Agents with mechanisms for self-discovery to identify its capabilities. Hence, the capabilities of the Agents need to be indicated at booting time. Currently, two capabilities can be indicates CPU - for high processing capacity - and CAMERA - to capture images -.
    On the other hand, Agents are not able to read any configuration file to indicate which policy to use; therefore, the policy needs to be indicated. Agents can be started in Eager or Lazy mode. In Eager mode, Agent will start all possible roles. In Lazy mode, Agent will only start roles with broken KPIs.

    The following command starts an agent capable of collecting images driven by the EAGER policy and contacting the COLMENA DCP deployed in localhost:5555.
    ``` bash
    java -classpath <path_to_agent_repository_root>/agent/build/libs/agent-0.0.1-all.jar es.bsc.colmena.Application localhost 5555 EAGER CAMERA
    ```

## Contributing
Please read our [contribution guidelines](CONTRIBUTING.md) before making a pull request.

## License
The COLMENA programming model is released under the Apache 2.0 license.
Copyright © 2022-2024 Barcelona Supercomputing Center - Centro Nacional de Supercomputación. All rights reserved.
See the [LICENSE](LICENSE) file for more information.


<sub>
	This work is co-financed by the COLMENA project of the UNICO I+D Cloud program that has the Ministry for Digital Transformation and of Civil Service and the EU-Next Generation EU as financing entities, within the framework of the PRTR and the MRR. It has also been supported by the Spanish Government (PID2019-107255GB-C21), MCIN/AEI /10.13039/501100011033 (CEX2021-001148-S), and Generalitat de Catalunya (2021-SGR-00412).
</sub>
<p align="center">
	<img src="https://github.com/colmena-swarm/.github/blob/assets/images/funding_logos/Logos_entidades_OK.png?raw=true" width="600">
</p>

