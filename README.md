capgo-server
===

### Table of Contents

- [capgo-server](#capgo-server)
    - [Table of Contents](#table-of-contents)
- [System Requirement](#system-requirement)
- [Running in Production](#running-in-production)
    - [Environment Configuration](#environment-configuration)
- [Usage](#usage)
  - [Concepts](#concepts)
    - [Bundle](#bundle)
    - [Release](#release)
  - [Workflow](#workflow)
- [License](#license)


---

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
| CACHE_RESULT_DURATION    | Duration for caching the result of the `POST /updates` API.                                                                                                                                                           | 10 minutes                                                    |
| OAUTH_ISSUER             | OIDC Issuer URL. No tailing slash. Please make sure it's matched with `iss` field in the token.                                                                                                                       | (Optional)                                                    |
| OAUTH_CLIENT_ID          | OAuth 2.0 client ID provided by OIDC issuer.                                                                                                                                                                          | (Optional)                                                    |
| CAPGO_USER_PORT          | Public server listen port for checking bundle update.                                                                                                                                                                 | 8000                                                          |
| CAPGO_MANAGEMENT_PORT    | Management server listen port for managing releases and bundles.                                                                                                                                                      | 8001                                                          |

These environment variables can be used to override the corresponding settings in the `config.yml` file. For more detailed information about the configuration, please refer to the [config/config.go](./config/config.go) file.

# Usage

## Concepts

### Bundle
Bundle is a zip file that contains the compiled Capacitor web assets of the app. 

Usually, when you're going to release your native application, you need to embed a specific bundle version into your native application. This embedded bundle version is called `builtin` bundle.

When you want to perform Over-The-Air update (OTA), you need to build a new bundle version, upload it to the capgo-server, and associate it with a specific release.

### Release
Release is a published (or will be published) native application version. These informations are very crutial and must be known.
- Platform (e.g. iOS, Android)
- App bundle name (e.g. com.example.app)
- App version (e.g. 2.15.1)
- Build number (e.g. 12)

Typically, these information are embed in to the native application binary and maintaned by developers/CI pipeline. For example, a developer or CI pipeline will update these information when preparing for a new release. For iOS, it's in `Info.plist` and for Android, it's in `AndroidManifest.xml`.

Once you determied the release information, you can create a new release in the capgo-server and set a default `builtin` bundle for that release.

After the application is released, you can upload new bundle version to the capgo-server and associate it with the release to perform OTA update.

The Capgo SDK will periodically check for a new bundle by providing the release information to the capgo-server. The capgo-server will identify the release and find the associated bundle for that release. If there's a new bundle available, the SDK will download the new bundle and prompt the user to update the app.

## Workflow

1. **Create a new bundle**
   - Upload a bundle zip file that contains the compiled Capacitor web assets and upload via UI or `POST /api/v1/bundles.upload`.
2. **Create a new release**
   - Provide the release information such as platform, bundle name, app version, and build number and set the default `builtin` bundle for that release via UI or `POST /api/v1/releases.create`.
3. **[Optional] For OTA update.** This step is done after you have released the native application via platform's app store.
   - Upload a new bundle zip version to capgo-server.
   - Associate the new bundle with the release via UI or `POST /api/v1/releases.set-active`.


# License

[MIT License](LICENSE.md)
