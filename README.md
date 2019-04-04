# GoogleAppEngine Empty Project

GoogleAppEngine Gen2 Empty Proejct

# Datastore Emulator

gcloud beta emulators datastore start --poject=Xxx

# Development

cd cmd

export DATASTORE_EMULATOR_HOST=localhost:8081

dev_appserver.py --support_datastore_emulator=true -A=Xxx app.yaml

-A : localhost:8000/datastore ->Projectid

# Deployment

cd cmd
gcloud app deploy --project=Xxx

