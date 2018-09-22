#!/bin/bash
(CHAMBER_SECRET_BACKEND=ASM aws-vault exec home -- ./chamber read devops-dashboard testsecret)
