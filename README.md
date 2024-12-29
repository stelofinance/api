> [!NOTE]
> No longer in use, refer to the [stelofinance](https://github.com/stelofinance/stelofinance) mono repo

# Stelo Finance
This is the Backend API for the Stelo Finance platform, providing it's core functionality.

## API Documentation
If you're looking for the API documentation, head over to our [docs](https://docs.stelo.finance)

## Environment Variables
In development use a .env file to set these, in prod set them in the environment
- `DB_CONNECTION_STRING` Set as the postgres URI connection string
- `ADMIN_API_KEY` Set as the api key to access admin routes
- `PUSHER_HOST` Set as the URL for the Pusher instance
- `PUSHER_APP_ID` Set as the app id for the Pusher instance
- `PUSHER_APP_KEY` Set as the app key for the Pusher instance
- `PUSHER_APP_SECRET` Set as your app secret for the Pusher instance

This doesn't need to be set in a .env file during development, it's absence will default it to false
- `PRODUCTION_ENV` Set to `true` when deployment is running in production

## Database setup
Install the `golang-migrate` software using the `postgres cockroachdb` tags (disregard cockroachdb if you're just using postgres)

Using the `golang-migrate` software, run the following command in the repo root
 - `migrate -database $COCKROACHDB_URL -path migrations/ up`
Where `$COCKROACHDB_URL` is an environment variable. Note, Stelo uses [CockroachDB](https://www.cockroachlabs.com/), if you wish to also then prefix this connection string with `cockroach://` not `postgresql://`
