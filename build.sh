#!/bin/bash

echo "Building frontend..."
cd frontend
yarn run build

echo "Building backend..."
cd ..
go build -ldflags="-s -w" -o app ./backend/cmd

echo "Done."
