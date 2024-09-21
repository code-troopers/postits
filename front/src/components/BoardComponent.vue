<template>
  <div>
  <nav>
    <ul>
      <li v-if="!voteModeStatus"><button @click="voteMode">Voter</button></li>
      <li v-if="voteModeStatus"><button @click="voteMode">Édition</button></li>
    </ul>
  </nav>
  <div style="width: 100vw; height: 100vh" @click="createPostit">
    <div v-for="postit in postits" :key="postit.id">
      <div
        @mouseover="hovered = true"
        @mouseleave="hovered = false"
        @click.stop="clickOnPostit(postit.id)"
         @contextmenu.prevent="rightClickOnPostit(postit.id)"
        :style="{ left: postit.posX + 'px', top: postit.posY + 'px' }"
        class="postit"
      >
        <textarea
          class="full-size"
          v-model="postit.text"
          :readonly="voteModeStatus"
          @change="updatePostit(postit.id, postit.text)"
        ></textarea>
      <button v-if="hovered" class="hover-button" @click="deletePostit(postit.id)">
        X
      </button>
      <div class="votes">
        {{ postit.votes }}
      </div>
      </div>
    </div>
  </div>
  </div>
</template>

<script setup lang="ts">
import { useBoardStore } from "@/stores/board";
import { computed, onMounted, ref } from "vue";
import { useRoute } from "vue-router";

const store = useBoardStore();
const route = useRoute();

const boardId = computed(() => route.params.id as string);

const postits = computed(() => store.getPostits(boardId.value));
const mainSelected = ref(true);
const hovered = ref(false);
const voteModeStatus = ref(false);

onMounted(() => {
  store.initPostits(boardId.value);
});

function voteMode() {
  voteModeStatus.value = !voteModeStatus.value;
}

function clickOnPostit(id: string | undefined) {
  if (id === undefined) {
    return;
  }
  mainSelected.value = false
  if (voteModeStatus.value) {
    store.addVote(id, boardId.value);
  }
}
function rightClickOnPostit(id: string | undefined) {
  if (id === undefined) {
    return;
  }
  if (voteModeStatus.value) {
    store.removeVote(id, boardId.value);
  }
}

function createPostit(event: MouseEvent) {
  if (!mainSelected.value) {
    mainSelected.value = true;
  } else {
    store.newPostit(boardId.value, event.clientX, event.clientY);
  }
}

function updatePostit(id: string | undefined, text: string | undefined) {
  if (id === undefined || text === undefined) {
    return;
  }
  store.updateContent(boardId.value, id, text);
}

function deletePostit(id: string | undefined) {
  if (id === undefined) {
    return;
  }
  store.deletePostit(boardId.value, id);
}
</script>

<style scoped>
.postit {
  position: absolute;
  width: 120px;
  height: 120px;
  background-color: #ffeb3b;
  border: 0.5px solid #000000;
}

.full-size {
  position: absolute; /* Permet au textarea de se positionner par rapport à .container */
  top: 0; /* Se caler sur le bord supérieur de .container */
  left: 0; /* Se caler sur le bord gauche de .container */
  right: 0; /* Se caler sur le bord droit de .container */
  bottom: 0; /* Se caler sur le bord inférieur de .container */
  border: none; /* Supprime les bordures par défaut */
  padding: 0; /* Supprime les marges internes */
  margin: 0; /* Supprime les marges externes */
  resize: none; /* Empêche le redimensionnement par l'utilisateur */
  box-sizing: border-box; /* Prend en compte la bordure et le padding dans la taille */
  font: inherit; /* Hérite de la police du parent */
  background: none; /* Supprime le background par défaut */
  outline: none; /* Supprime le contour bleu par défaut lors du focus */
}

.hover-button {
  position: absolute;
  top: 5px;
  right: 5px;
  display: block;
  background-color: #f56c6c;
  border: none;
  color: white;
  padding: 5px 10px;
  cursor: pointer;
  border-radius: 3px;
}

.hover-button:hover {
  background-color: #c0392b;
}

.votes {
  position: absolute;
  bottom: 5px;
  right: 5px;
}
</style>
