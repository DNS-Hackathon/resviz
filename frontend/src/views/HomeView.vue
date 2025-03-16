<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import axios from 'axios';
import qs from 'qs';
import MermaidChart from '@/components/MermaidChart.vue'
import ChartNavigator from '@/components/ChartNavigator.vue'
import SearchBar from '@/components/SearchBar.vue'

const step = ref(0);

const charts = ref(null);

const nextStep = () => {
  if (step.value < charts.length -1) step.value++;
};

const prevStep = () => {
  if (step.value > 0) step.value--;
};

const handleSearch = async (query) => {
  try {

    const response = await axios.post(
      'http://localhost:1999/resviz',
      qs.stringify({ domain: query }),
      { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } }
    );

    console.log(response.data);
    charts.value = [response.data];
    if (response.status !== 200) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
  } catch (err) {
    alert(err);
  }
};

const handleKeydown = (event) => {
  if (event.key === "ArrowRight") {
    nextStep();
  } else if (event.key === "ArrowLeft") {
    prevStep();
  }
};

onMounted(() => {
  window.addEventListener("keydown", handleKeydown);
});

onUnmounted(() => {
  window.removeEventListener("keydown", handleKeydown);
});
</script>

<template>
  <main>
    <SearchBar @search="handleSearch" />
    <div v-if="charts">
      <ChartNavigator :step="step" :maxSteps="charts.length" @next="nextStep" @prev="prevStep" />
      <MermaidChart :step="step" :charts="charts"/>
    </div>
  </main>
</template>
