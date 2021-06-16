# ocp_alert_webhook_tester
Simple app for testing OpenShift alert web hooks

Use s2i to deploy this.

Configure the webhook in OpenShift's AlertManager Alert Receiver as URL_OF_APP/api/alerts/set, example: http://someserver.com/api/alert/set

To see the collected alert messages, use /api/alerts/get. Example: curl http://someserver.com/api/alert/set. This will output the message bodies and the timestamps for each message.