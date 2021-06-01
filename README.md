# pancast-server
Go HTTP server for PanCast Backend
Structure and helper function adapted from Ivan's CPSC 416 course assignment 5 (key value store).

This is a HTTPS server that intended to service three routes:
1) Device registration
2) Encounter uploading
3) Risk dissemination

Instructions for setting up:
1) Make sure you have Golang installed on your machine. 
2) Clone the repo
3) Navigate to the root directory, and run `make`
4) Now you can run `./app`, and the server will listen to a specified address from the `app_config.json` file in the `/config` folder.
