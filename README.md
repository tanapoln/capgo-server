capgo-server
===

This repo provide server for [Capgo](https://capgo.app/) (live update for Capacitor).

To getting started with docker-compose:
```bash
docker-compose up -d
```
Default web management UI will be served at http://localhost:8080/ui

# System Requirement

capgo-server requires MongoDB version 4 or later but it's tested with MongoDB 5.

S3 or a compatible storage is required for storing app bundles.

# Running in Production

You can run capgo-server using docker image from docker hub: `tanapolsh/capgo-server:latest`

### Environment Configuration

The following table describes the environment variables that can be used to configure the capgo-server:

| Environment Variable     | Description                                                                                                                                                                                                           | Default Value                                                 |
| ------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------- |
| MONGO_CONNECTION_STRING  | MongoDB connection string                                                                                                                                                                                             | mongodb://user:pass@localhost:27017                           |
| MONGO_DATABASE           | Name of the MongoDB database                                                                                                                                                                                          | capgo                                                         |
| S3_BASE_ENDPOINT         | For overriding endpoint with s3-compatible storage                                                                                                                                                                    | (Optional)                                                    |
| S3_BUCKET                | Name of the S3 bucket for storing app bundles                                                                                                                                                                         | (Required)                                                    |
| MANAGEMENT_API_TOKENS    | Comma-separated list of API tokens for management access                                                                                                                                                              | (Required)                                                    |
| LIMIT_REQUEST_PER_MINUTE | Rate limit for API requests per minute                                                                                                                                                                                | 100                                                           |
| TRUSTED_PROXIES          | Comma-separated list of trusted proxy IP addresses or CIDR ranges. capgo-server will use this to determine if the forwarded client IP address is trusted. Then, the client IP address will be used for rate limiting. | (Optional)                                                    |
| AWS_ACCESS_KEY_ID        | AWS access key ID for S3 authentication                                                                                                                                                                               | Automatically resolve using AWS SDK Credential Provider Chain |
| AWS_SECRET_ACCESS_KEY    | AWS secret access key for S3 authentication                                                                                                                                                                           | Automatically resolve using AWS SDK Credential Provider Chain |
| AWS_REGION               | AWS region for S3 bucket                                                                                                                                                                                              | Automatically resolve using AWS SDK configuration resolution. |
| CACHE_RESULT_DURATION    | Duration for caching the result of the `POST /updates` API.                                                                                                                                                               | 10 minutes                                                    |

These environment variables can be used to override the corresponding settings in the `config.yml` file. For more detailed information about the configuration, please refer to the [config/config.go](./config/config.go) file.


# License

[MIT License](LICENSE.md)
