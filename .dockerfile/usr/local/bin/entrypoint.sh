#!/bin/bash
export PATH=${PATH}:$(dirname "${TARGET}")
exec "$@"