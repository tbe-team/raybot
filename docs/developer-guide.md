# Developer guide

This guide helps you get started developing Raybot.

## Dependencies

Make sure you have the following dependencies installed:

- [Git](https://git-scm.com/)
- [Go](https://go.dev/) (See [go.mod](../go.mod) for the version used in the project)
- [Node.js](https://nodejs.org/), with corepack enabled. See [.nvmrc](../ui/.nvmrc) for the version used in the project.
We can use a version manager like [nvm](https://github.com/nvm-sh/nvm) to install the correct version of Node.js.
- [GCC](https://gcc.gnu.org/) (required for Cgo dependencies)

## Development

### UI

The UI is built with [Vue.js](https://vuejs.org/) and [Tailwind CSS](https://tailwindcss.com/).

To run the UI, run `pnpm dev` in the `ui` directory.

### Backend

The backend is built with [Go](https://go.dev/).

To run the backend, run `make run` in the root directory.

Run with docker:

```bash
make docker-build-raybot
make docker-run-raybot
```
