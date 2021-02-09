#!/bin/bash

if [ -z "$DD_API_KEY" ]; then
	echo "DD_API_KEY is not set, proceeding without it"
	exec "$@"
fi

export DD_AGENT_HOST=localhost

if [ -z "$DYNO" ]; then
  export DD_HOSTNAME=localhost
else
  export DD_HOSTNAME="${HEROKU_APP_NAME}-${DYNO}"
fi

export DD_SERVICE="$SERVICE_NAME"
export DD_ENV="$ENV_NAME"
export DD_VERSION="$HEROKU_RELEASE_VERSION"

# Yes, process agent uses SPACES to separate tags and Go uses commas. Using a space and
# a comma appears to work with both.
export DD_TAGS="service:${SERVICE_NAME}, dyno-id:${HEROKU_DYNO_ID}"
/usr/bin/datadog-agent run &
/opt/datadog-agent/embedded/bin/trace-agent &
/opt/datadog-agent/embedded/bin/process-agent &

exec "$@"
