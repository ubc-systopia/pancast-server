# pancast-server
Go HTTP server for PanCast Backend
Structure and helper function adapted from Ivan's CPSC 416 course assignment 5 (key value store).

This is a HTTPS server that intended to service three routes:
1) Device registration
2) Encounter uploading
3) Risk dissemination

Instructions for setting up:
1) Make sure you have Golang installed on your machine.
2) Clone the repo.
3) Create a pancast.key and pancast.cert file. These are your private key and TLS certificate, respectively. For a production environment, we can either register a certificate with a trusted CA, or we can host this server on a cloud service (i.e. Heroku). For now, we can use a self-signed certificate.
4) Create a .env file populated with the DB_USER, DB_PASSWD, DB_NAME, DB_ADDR, DB_PORT fields. These are database parameters.
   Also, add fields for MEAN, SENS, EPSILON and DELTA. These are parameters for the noise generating algorithm. A template can be found at env.template.
5) Create a app_config.json file in the config folder. This is based off of the app_config_template.json file.   
6) Navigate to the root directory, and run `make`.
7) Now you can run `./app`, and the server will listen to a specified address from the `app_config.json` file in the `/config` folder.

Manual Test Log (integration, stress): <br>
/update: <br>
HTTPS request for the cuckoo filter pulls through regardless of how large the cuckoo filter is. Tested for cuckoo filter sizes up to 216MB. Performance can be quite poor though.

/upload: <br>
Can upload an arbitrary number of unique HTTPS requests. Tested for 100,000 concurrent uploads. Upload speed can reach up to 88s. <br>

TODO List index (accurate as of 17-08-2021):
1) Change the logic of how the cuckoo filter grows in size. Currently, it can only support a CF with 2^x entries for some x.
2) Change the way ephemeral IDs are uploaded to be serialized for performance. Either we can create our own way of serializing the data, or we can use Google's protocol buffers which I heard is pretty efficient.
3) Implement ephemeral ID verification as outlined in section 3.4 in the Pancast paper.
4) Remove the telemetry logic in production. We don't want to collect any more data than we need at that stage (right)?
5) Clean the codebase up. Apply good software design practices, unlike what I did :)
6) A whole host of performance tuning if possible.
7) Security and access control. e.g. anyone can upload to the risk database. This needs to be a lot more stringent IMO.

An architectural note: <br>
The backend repo's architecture shares a lot of commonalities with that of JS express servers. <br>
The `app` file invokes the `server`, which creates handlers for HTTPS routes, which are located in the `routes` folder. These routes act as the main controller for the route, invoking functions from the `models` folder (db abstraction), `util` functions etc.
