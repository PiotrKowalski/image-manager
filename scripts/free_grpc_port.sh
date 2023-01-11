#!/bin/bash

comm -23 <(seq 5301 5305 | sort) <(ss -Htan | awk '{print $4}' | cut -d':' -f2 | sort -u) | sort -n | head -n 1
