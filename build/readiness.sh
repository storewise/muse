#!/bin/bash
test $(curl -sSL http://localhost:9480/synced) == true
