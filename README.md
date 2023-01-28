# Stelo Finance API
This is the Backend API for the Stelo Finance platform, providing it's core functionality.

## Environment Variables
In development use a .env file to set these, in prod set them in the environment
- `DB_CONNECTION_STRING` Set as the postgres URI connection string
- `JWT_SECRET` Set as the secret to be used for creating JWTs
- `ADMIN_API_KEY` Set as the api key to access admin routes
- `CENTRIFUGO_API_KEY` Set as the api key to post data to Centrifugo
- `CENTRIFUGO_API_ADDR` Set as the api endpoint for Centrifugo

This doesn't need to be set in a .env file during development, it's absence will default it to false
- `PRODUCTION_ENV` Set to `true` when deployment is running in production
