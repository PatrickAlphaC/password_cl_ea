# Password chainlink EA

You need a working MongoDB database for this to work. 

# Testing

```
go test
```

# Upload to GCP

```
gcloud functions deploy password-api --runtime go113 --trigger-http --allow-unauthenticated --source password_cl_ea --entry-point MakeRequest
```
