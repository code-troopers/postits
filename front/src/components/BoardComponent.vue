<template>
  <div>
    <nav>
      <ul>
        <li v-if="!voteModeStatus"><button @click="voteMode">Voter</button></li>
        <li v-if="voteModeStatus"><button @click="voteMode">Édition</button></li>
        <li v-if="showMode"><button @click="showHide">Cacher les tickets</button></li>
        <li v-if="!showMode"><button @click="showHide">Montrer les tickets</button></li>
      </ul>
    </nav>
    <div
      style="width: 100vw; height: 100vh"
      @click="createPostit"
      ref="parentDiv"
      @mousemove="onDrag"
      @mouseup="endDrag"
    >
      <div v-for="postit in postits" :key="postit.id">
        <div
          @mouseover="hovered = postit.id || ''"
          @mouseleave="hovered = ''"
          @mousedown="startDrag(postit, $event)"
          @click.stop="clickOnPostit(postit.id)"
          @contextmenu.prevent="rightClickOnPostit(postit.id)"
          :style="{ left: postit.posX + 'px', top: postit.posY + 'px' }"
          class="postit"
          :class="{ highlight: selectedAuthor === postit.author?.id }"
        >
          <textarea
            class="full-size"
            v-model="postit.text"
            :class="{ 'hidden-font': !notMyPostit(postit.author?.id) && !postit.show }"
            :readonly="voteModeStatus || notMyPostit(postit.author?.id)"
            @change="updatePostit(postit.id, postit.text)"
          ></textarea>
          <button
            v-if="hovered === postit.id"
            class="hover-button"
            @click="deletePostit(postit.id)"
          >
            X
          </button>
          <div class="votes">
            {{ postit.votes }}
          </div>
        </div>
      </div>
    </div>
    <div class="author-list">
      <ul>
        <li v-for="author in authorList" :key="author.id">
          <button
            :class="{ highlight: selectedAuthor === author.id }"
            @click="selectAuthor(author.id || '')"
          >
            {{ author.givenName }}
          </button>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import keycloak from "@/keycloak";
import type { StickyNote } from "@/models/StickyNote";
import type { User } from "@/models/User";
import { useBoardStore } from "@/stores/board";
import { computed, onMounted, ref } from "vue";
import { useRoute } from "vue-router";

const store = useBoardStore();
const route = useRoute();

const boardId = computed(() => route.params.id as string);

const postits = computed(() => store.getPostits(boardId.value));
const mainSelected = ref(true);
const hovered = ref("");
const voteModeStatus = ref(false);
const showMode = ref(false);
let throttleTimeout: number | null = null;
const selectedAuthor = ref("");

const authorList = computed(() => {
  const uniqueAuthors: { [id: string]: User } = {};
  if (!postits.value) {
    return Object.values([]);
  }

  postits.value.forEach((postit) => {
    if (postit.author?.id && postit.author.givenName && !uniqueAuthors[postit.author.id]) {
      uniqueAuthors[postit.author.id] = postit.author;
    }
  });

  return Object.values(uniqueAuthors);
});

function selectAuthor(id: string) {
  if (selectedAuthor.value === id) {
    selectedAuthor.value = "";
  } else {
    selectedAuthor.value = id;
  }
}

const isDragging = ref(false);
const draggedPostit = ref<StickyNote | null>(null); // Index de l'élément actuellement déplacé
const parentDiv = ref<HTMLElement | null>(null);
const initialMouseX = ref(0);
const initialMouseY = ref(0);
const initialX = ref(0);
const initialY = ref(0);

const startDrag = (postit: StickyNote, event: MouseEvent) => {
  isDragging.value = true;
  draggedPostit.value = postit;
  initialMouseX.value = event.clientX;
  initialMouseY.value = event.clientY;
  initialX.value = postit.posX || 0;
  initialY.value = postit.posY || 0;
};

const onDrag = (event: MouseEvent) => {
  if (isDragging.value && draggedPostit.value !== null) {
    const dx = event.clientX - initialMouseX.value;
    const dy = event.clientY - initialMouseY.value;

    draggedPostit.value.posX = initialX.value + dx;
    draggedPostit.value.posY = initialY.value + dy;
    if (throttleTimeout === null) {
      store.movePostit(boardId.value, draggedPostit.value);
      throttleTimeout = window.setTimeout(() => {
        throttleTimeout = null;
      }, 10);
    }
  }
};

const endDrag = () => {
  isDragging.value = false;
  if (
    draggedPostit.value !== null &&
    (draggedPostit.value.posX !== initialX.value ||
      draggedPostit.value.posY !== initialY.value)
  ) {
    store.movePostit(boardId.value, draggedPostit.value);
  }
  draggedPostit.value = null;
};

onMounted(() => {
  store.initPostits(boardId.value);
  const p = postits.value.find(
    (postit) => postit.author?.id === keycloak.tokenParsed?.sub
  );
  if (p) {
    showMode.value = p.author?.id === keycloak.tokenParsed?.sub;
  }
});

function showHide() {
  if (showMode.value) {
    store.hidePostits(boardId.value, keycloak.tokenParsed?.sub || "");
  } else {
    store.showPostits(boardId.value, keycloak.tokenParsed?.sub || "");
  }
  showMode.value = !showMode.value;
}

function notMyPostit(authorId: string | undefined) {
  if (authorId === undefined) {
    return true;
  }
  return keycloak.tokenParsed?.sub !== authorId;
}

function voteMode() {
  voteModeStatus.value = !voteModeStatus.value;
}

function clickOnPostit(id: string | undefined) {
  if (id === undefined) {
    return;
  }
  mainSelected.value = false;
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
  transition: box-shadow 0.3s ease, transform 0.3s ease;
}

.highlight {
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.4), 0 0 20px rgba(43, 70, 222, 0.7),
    inset 0 0 10px rgba(255, 255, 255, 0.6); /* Ombre plus prononcée et effet lumineux accentué */
  border: 1px solid rgba(43, 70, 222, 0.7);
  transform: scale(1.02);
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

.author-list {
  position: absolute;
  top: 5px;
  right: 5px;
}
</style>
