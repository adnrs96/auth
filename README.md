### Storyscript Auth Server

The Auth server is used to generate tokens for use with the Storyscript platform.


#### Running the Server

The server requires the following environment variables to be set:
 - `GH_CLIENT_ID` for the Client ID of the GitHub application to authenticate against
 - `GH_CLIENT_SECRET` for the Client Secret of the GitHub application to authenticate against
 - `JWT_SIGNING_KEY` for signing JWTs
 - `DB_CONNECTION_STRING` for connecting to the Storyscript database

#### Running The Tests

The tests can be run using `ginkgo -r` from the root directory.

The server requires the following environment variables to be set:
 - `GH_CLIENT_ID` for the Client ID of the GitHub application to authenticate against
 - `GH_CLIENT_SECRET` for the Client Secret of the GitHub application to authenticate against
 - `JWT_SIGNING_KEY` for signing JWTs
 - `TEST_DB_NAME` `TEST_DB_USER` `TEST_DB_HOST` and `TEST_DB_PASSWORD` for connecting to the test database

The database must already have the Storyscript schema defined from https://github.com/storyscript/database
