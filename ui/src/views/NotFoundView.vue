<template>
  <div class="flex flex-col items-center justify-center min-h-screen p-4">
    <div class="space-y-6 text-center">
      <h1 class="text-6xl font-bold text-primary">
        404
      </h1>
      <h2 class="text-2xl font-medium">
        Page not found
      </h2>
      <p class="max-w-md mx-auto text-muted-foreground">
        The page you're looking for doesn't exist or has been moved.
      </p>
      <RouterLink
        to="/"
        class="inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium transition-colors rounded-md shadow whitespace-nowrap focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-9"
      >
        Go back home
      </RouterLink>
    </div>
  </div>
</template>
