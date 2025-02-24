.PHONY: run
run:
	uv run src/main.py

.PHONY: test
test:
	uv run pytest tests

.PHONY: test-cov
test-cov:
	uv run pytest tests --cov src

.PHONY: lint
lint:
	scripts/lint.sh
