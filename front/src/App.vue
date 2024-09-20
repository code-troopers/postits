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
