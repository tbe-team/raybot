export PYTHONPATH=$(shell printenv PYTHONPATH):$(pwd)/src/

.PHONY: run
run:
	uv run src/main.py

.PHONY: test
test:
	uv run pytest

.PHONY: test-cov
test-cov:
	uv run pytest --cov src

.PHONY: lint
lint:
	scripts/lint.sh
