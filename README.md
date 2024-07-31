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
- **agent**: Contains the code of a COLMENA Agent.
- **agent-zenoh-client**: Rust interface service to Zenoh.
- **agent-zenoh-router**: Contains files to start up the zenoh router as a Docker container.

### Files
- **.gitignore**: Specifies files and directories to be ignored by Git.
- **changeLog**: Change highlights associated with official releases.
- **CODE_OF_CONDUCT.md**: Outlines the expected behavior and guidelines for participants within the project's community. 
- **compose.yaml**: Configuration file for the multi-container application.
- **CONTRIBUTING.md**: Overview of the repository, setup instructions, and basic usage examples.
- **LICENSE**: License information for the repository.
- **README.md**: Overview of the repository, setup instructions, and basic usage examples.
- **setting.gradle**: Configuration of gradle to build the project.


## Getting Started
1. Start Zenoh Router
`docker compose -f compose-zenoh.yaml up --abort-on-container-exit`

2. Start COLMENA agent
`DEVICE_STRATEGY=${DEVICE_STRATEGY} DEVICE_HARDWARE=${DEVICE_HARDWARE} DISCOVERY_INTERFACE={NETWORK_INTERFACE} DEVICE_NAME={DEVICE_NAME} docker compose -f compose.yaml up --abort-on-container-exit`
where DEVICE_STRATEGY can be EAGER or LAZY and DEVICE_HARDWARE is defined by the role.
DISCOVERY_INTERFACE is the network interface for peer discovery, if not defined then a default interface will be picked.

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

