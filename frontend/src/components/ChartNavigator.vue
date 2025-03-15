<script setup>
import { defineProps, defineEmits } from "vue";

const props = defineProps({
  step: Number,
  maxSteps: Number,
});

const emit = defineEmits(["next", "prev"]);
</script>

<template>
  <div class="buttons">
    <button @click="emit('prev')" :disabled="step === 0">Back</button>
    <button @click="emit('next')" :disabled="step === maxSteps - 1">Next</button>
  </div>
</template>

<style scoped>
.buttons {
  display: flex;
  gap: 10px;
}

button {
  padding: 10px;
  border: none;
  background: var(--color-button);
  color: white;
  cursor: pointer;
  border-radius: 5px;
}

button:hover {
  background: var(--color-button-hover);
}

button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

</style>