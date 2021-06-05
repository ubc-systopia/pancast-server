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
4) Create a .env file populated with a DB_USER, DB_PASSWD, DB_NAME field. These are database parameters.
5) Navigate to the root directory, and run `make`.
6) Now you can run `./app`, and the server will listen to a specified address from the `app_config.json` file in the `/config` folder.
