<script setup lang="ts">
import { Board } from '@/models/Board';
import { useBoardStore } from '@/stores/board';
import { ref, onMounted } from 'vue';

const store = useBoardStore();

const b = new Board();
const editingBoard = ref(b);
const isEditing = ref(false);
const isDeleting = ref(false);
const boardToDelete = ref(b);

onMounted(() => {
  store.getBoards();
});

const editBoard = (board: Board) => {
  editingBoard.value = { ...board };
  isEditing.value = true;
};

const saveBoardName = () => {
  if (!editingBoard.value.id || !editingBoard.value.name) {
    return
  }
  store.renameBoard(editingBoard.value.id, editingBoard.value.name);
  isEditing.value = false;
};

const confirmDeleteBoard = (board: Board) => {
  boardToDelete.value = { ...board };
  isDeleting.value = true;
};

const deleteBoard = () => {
  if (boardToDelete.value.id) {
    store.deleteBoard(boardToDelete.value.id);
    isDeleting.value = false;
    boardToDelete.value = b;
  }
};
</script>

<template>
  <main class="container">
    <h1>Boards</h1>
    <button @click="store.newBoard()" class="btn btn-add">Add Board</button>
    <div class="board-grid">
      <div v-for="board in store.boards" :key="board.id" class="board-card">
        <RouterLink :to="`/board/${board.id}`" class="board-name">
          {{ board.name }}
        </RouterLink>
        <div class="btn-group">
          <button @click="editBoard(board)" class="btn btn-edit">Edit</button>
          <button @click="confirmDeleteBoard(board)" class="btn btn-delete">
            Delete
          </button>
        </div>
      </div>
    </div>

    <!-- Edit Board Modal -->
    <div v-if="isEditing" class="modal-overlay">
      <div class="modal">
        <h2>Edit Board Name</h2>
        <input
          v-model="editingBoard.name"
          type="text"
          class="input"
          placeholder="Enter new board name"
        />
        <div class="btn-group">
          <button @click="isEditing = false" class="btn btn-cancel">Cancel</button>
          <button @click="saveBoardName" class="btn btn-save">Save</button>
        </div>
      </div>
    </div>
    <div v-if="isDeleting" class="modal-overlay">
      <div class="modal">
        <h2>Confirm Deletion</h2>
        <p>Are you sure you want to delete the board "{{ boardToDelete.name }}"?</p>
        <div class="btn-group">
          <button @click="isDeleting = false" class="btn btn-cancel">Cancel</button>
          <button @click="deleteBoard" class="btn btn-delete">Delete</button>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
/* Container */
.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 16px;
  text-align: center;
}

/* Heading */
h1 {
  font-size: 2rem;
  font-weight: bold;
  margin-bottom: 16px;
}

/* Button Styles */
.btn {
  cursor: pointer;
  padding: 8px 16px;
  border-radius: 4px;
  border: none;
  color: #fff;
  margin: 4px;
  transition: background-color 0.3s;
}

.btn-add {
  background-color: #3490dc;
}

.btn-add:hover {
  background-color: #2779bd;
}

.btn-edit {
  background-color: #f6ad55;
}

.btn-edit:hover {
  background-color: #e58e26;
}

.btn-delete {
  background-color: #e3342f;
}

.btn-delete:hover {
  background-color: #cc1f1a;
}

.btn-cancel {
  background-color: #606f7b;
}

.btn-cancel:hover {
  background-color: #4a5568;
}

.btn-save {
  background-color: #38a169;
}

.btn-save:hover {
  background-color: #2f855a;
}

/* Board Grid */
.board-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  justify-content: center;
}

/* Board Card */
.board-card {
  background-color: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  padding: 16px;
  width: 100%;
  max-width: 300px;
  transition: transform 0.3s, box-shadow 0.3s;
}

.board-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 12px rgba(0, 0, 0, 0.1);
}

/* Board Name */
.board-name {
  display: block;
  font-size: 1.25rem;
  font-weight: 500;
  color: #3490dc;
  text-decoration: none;
  margin-bottom: 8px;
}

.board-name:hover {
  text-decoration: underline;
}

/* Button Group */
.btn-group {
  display: flex;
  gap: 8px;
  justify-content: center;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal {
  background: #fff;
  padding: 24px;
  border-radius: 8px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
}

.modal h2 {
  margin-bottom: 16px;
  font-size: 1.5rem;
  font-weight: bold;
}

.input {
  width: 100%;
  padding: 8px;
  margin-bottom: 16px;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  outline: none;
  transition: border-color 0.3s;
}

.input:focus {
  border-color: #3490dc;
}
</style>
