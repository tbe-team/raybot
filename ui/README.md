# Raybot UI

The Raybot UI is a UI for Raybot for managing and monitoring the Raybot application.

## Tech Stack

- **Framework**: Vue 3 with Composition API
- **Build Tool**: Vite
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **UI Components**: Custom components based on shadcn/vue
- **Form Validation**: VeeValidate with Zod
- **HTTP Client**: Axios
- **Icons**: Lucide Vue

## Project Structure

```
src/
├── assets/           # Static assets
├── components/       # Reusable components
│   ├── ui/           # UI components (buttons, inputs, etc.)
│   └── shared/       # Shared components across features
├── composables/      # Reusable composition functions
├── layouts/          # Page layouts
├── lib/              # Utility functions and libraries
├── router/           # Vue Router configuration
├── types/            # TypeScript type definitions
└── views/            # Page components
```

## Component Organization

Components are organized following a variant of the Atomic Design methodology:

1. **UI Components**: Low-level, reusable components in `src/components/ui/`
2. **Shared Components**: Common components used across features
3. **Layout Components**: Components that define the structure of pages
4. **View Components**: Page-level components that compose other components

## Getting Started

### Prerequisites

- Node.js 20+
- pnpm 10+

### Installation

```
# Clone the repository
git clone https://github.com/tbe-team/raybot
cd raybot/ui

# Install dependencies
pnpm install
```

### Development

```
# Start the development server
pnpm dev
```

The application will be available at http://localhost:5173 by default.

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Type Support for `.vue` Imports in TS

TypeScript cannot handle type information for `.vue` imports by default, so we replace the `tsc` CLI with `vue-tsc` for type checking. In editors, we need [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) to make the TypeScript language service aware of `.vue` types.

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## Build and Deployment

```
# Type-check, compile and minify for production
pnpm build

# Preview the production build locally
pnpm preview
```

### Lint with [ESLint](https://eslint.org/)

```
pnpm lint
```

## Contributing

We welcome contributions to Raybot UI! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow the existing code style and project structure
- Write clean, maintainable, and testable code
- Use TypeScript for type safety
- Document your code when necessary
- Update the README if you make significant changes

### Commit Guidelines

- Use conventional commit messages
- Keep commits focused on a single change
- Write clear and concise commit messages
