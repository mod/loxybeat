#!/bin/bash

kaigara render loxybeat.yml > $APP_HOME/loxybeat.yml
kaigara render loxybeat.template-es2x.json > $APP_HOME/loxybeat.template-es2x.json
kaigara render loxybeat.template.json > $APP_HOME/loxybeat.template.json
