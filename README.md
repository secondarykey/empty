# GoogleAppEngine Empty Project

GoogleAppEngine Gen2 Empty Proejct

# Datastore Emulator

gcloud beta emulators datastore start --poject=Xxx

# Development

cd cmd
$(gcloud beta eumulators datastore env-init)
dev_appserver.py --support_datastore_emulator=true -A=Xxx .

# Deployment

cd cmd
gcloud app deploy --project=Xxx



