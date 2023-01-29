# Stelo Finance API
This is the Backend API for the Stelo Finance platform, providing it's core functionality.

## Environment Variables
In development use a .env file to set these, in prod set them in the environment
- `DB_CONNECTION_STRING` Set as the postgres URI connection string
- `JWT_SECRET` Set as the secret to be used for creating JWTs
- `ADMIN_API_KEY` Set as the api key to access admin routes
- `CENTRIFUGO_API_KEY` Set as the api key to post data to Centrifugo
- `CENTRIFUGO_API_ADDR` Set as the api endpoint for Centrifugo
- `CENTRIFUGO_JWT_KEY` Set as the JWT secret for Centrifugo

This doesn't need to be set in a .env file during development, it's absence will default it to false
- `PRODUCTION_ENV` Set to `true` when deployment is running in production

## Database setup
Install the `golang-migrate` software using the `postgres cockroachdb` tags (disregard cockroachdb if you're just using postgres)

Using the `golang-migrate` software, run the following command in the repo root
 - `migrate -database $COCKROACHDB_URL -path migrations/ up`
Where `$COCKROACHDB_URL` is an environment variable. Note, Stelo uses [CockroachDB](https://www.cockroachlabs.com/), if you wish to also then prefix this connection string with `cockroachdb://` not `postgresql://`
