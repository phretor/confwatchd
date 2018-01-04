#!/bin/bash
git pull && make && make setcap && clear && ./confwatchd -config prod-config.json
