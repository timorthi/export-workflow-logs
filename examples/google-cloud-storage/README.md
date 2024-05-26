The Google Cloud Storage exporter uses the [Google Cloud Client Libraries for Go](https://pkg.go.dev/cloud.google.com/go#hdr-Authentication_and_Authorization) under the hood.

This exporter does not accept Action-level inputs for credentials or path to credentials files. Authenticating to GCP via [google-github-actions/auth](https://github.com/google-github-actions/auth) is recommended.
