#!/bin/bash
git pull && make && make setcap && clear && ./confwatchd -config prod-config.json.json -seed $HOME/confwatch-data && ./confwatchd -config prod-config.json
