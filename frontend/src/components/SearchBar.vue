<script setup>
import { ref, defineEmits } from "vue";

const searchQuery = ref("");
const emit = defineEmits(["search"]);

const handleSearch = () => {
  emit("search", searchQuery.value);
};

</script>

<template>
  <div class="search-bar">
    <input
      v-model="searchQuery"
      placeholder="Domain to visualize..."
      type="text"
    />
    <button @click="handleSearch">Visualise</button>
  </div>
</template>

<style scoped>
.search-bar {
  max-wdith: 10%;
  display: flex;
  gap: 10px;
  margin-top: 20px;
  margin-bottom: 10px;
}

input {
  flex: 1;
  padding: 15px;
  background: var(--color-background);
  color: var(--color-text);
  border: 1px solid var(--color-text);
  border-radius: 5px;
}

button {
  padding: 8px 12px;
  border: none;
  background: var(--color-button);
  color: white;
  cursor: pointer;
  border-radius: 5px;
}

button:hover {
  background: var(--color-button-hover);
}
</style>