# Aries

As opposed to a monolithic repository, the UVA digital repository solution involves a great number of independent systems.  Experienced and knowledgeable staff make finding, accessing and integrating the pieces possible, but it's time to bring the information together in a more machine-actionable way.

The goal of Aries is to provide a method for seeing the full picture in a coherent way.

### System Requirements
* GO version 1.11 or greater
* DEP (https://golang.github.io/dep/) version 0.5 or greater
* Node version 8.11.1 or higher (https://nodejs.org/en/)
* Yarn version 1.9.3 or greater
* vue-cli 3 version 3.0.3 or greater
* Vue 2.5 or greater

### Build Instructions

1. After clone, `cd frontend` and execute `yarn install` to prepare the front end
2. Run the Makefile target `build all` to generate binaries for linux, darwin and the front end.  All results will be in the bin directory.

### Current API

* GET /version : return service version info
* GET /healthcheck : test health of system components; results returned as JSON.
* GET /resources/:ID : Get a block of JSON data containing details about the specified resource.
* GET /services : Get a JSON list services that are a part of the aries constellation.

### Notes

To run in a development mode for the frontend only, build and launch aries, and supply a port param: `./aries.linux -port 8085`. Set an ENV variable:

* `ARIES_API` - set to localhost at the port used to launch the server (8085 from the example)

Run the frontend in development mode with `yarn serve` from the frontend directory. All API requests from the frontend will be redirected to the local instance of the backend services.
