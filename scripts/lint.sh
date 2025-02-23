#!/bin/bash

set -e
set -x

uv run mypy src
uv run ruff check src
uv run ruff format --check src
