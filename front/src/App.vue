<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <header>
      <nav>
        <RouterLink to="/">Boards list</RouterLink>
      </nav>
  </header>
  <RouterView />
</template>

<style scoped>
header {
  width: 100%;
  background-color: var(--color-background-mute);
  display: flex;
  align-items: center;
  justify-content: start;
  padding: 1rem;
}
</style>