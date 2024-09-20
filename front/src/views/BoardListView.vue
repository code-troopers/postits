<script setup lang="ts">
import { useBoardStore } from '@/stores/board';
import { onMounted } from 'vue';

const store = useBoardStore();

onMounted(() => {
  store.getBoards();
})

</script>

<template>
  <main>

    <button @click="store.newBoard()">Add Board</button>

    <div v-for="board in store.boards" :key="board.id">
      <RouterLink :to="`/board/${board.id}`">{{ board.name }}</RouterLink>
      <button @click="store.deleteBoard(board.id)">Delete</button>
    </div>
  </main>
</template>
