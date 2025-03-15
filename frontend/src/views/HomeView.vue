<script setup lang="ts">
import { ref } from "vue";
import MermaidChart from '@/components/MermaidChart.vue'
import ChartNavigator from '@/components/ChartNavigator.vue'

const step = ref(0);

const charts = [
  `graph TD; A[Start]`,
  `graph TD; A[Start] --> B[Step 1]`,
  `graph TD; A[Start] --> B[Step 1]; B --> C[Step 2]`,
  `graph TD; A[Start] --> B[Step 1]; B --> C[Step 2]; C --> D[Final Step]`,
];

const nextStep = () => {
  if (step.value < charts.length -1) step.value++;
};

const prevStep = () => {
  if (step.value > 0) step.value--;
};
</script>

<template>
  <main>
    <ChartNavigator :step="step" :maxSteps="charts.length" @next="nextStep" @prev="prevStep" />
    <MermaidChart :step="step" :charts="charts"/>
  </main>
</template>
