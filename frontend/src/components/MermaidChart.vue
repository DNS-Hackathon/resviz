<script setup>
import { defineProps, onMounted, watch, nextTick } from "vue";
import mermaid from "mermaid";

const props = defineProps({
  step: Number,
  charts: Array,
});

const renderMermaid = async () => {
  await nextTick();

  const container = document.getElementById("mermaid-container");
  if (!container) return;
  
  mermaid.initialize({
    startOnLoad: false,
    theme: "dark",
  });
  const { svg } = await mermaid.render("mermaidChart", props.charts[props.step]);

  container.innerHTML = svg;
};

onMounted(renderMermaid);

watch(() => props.step, renderMermaid)
</script>


<template>
  <div id="mermaid-container" class="mermaid"></div>
</template>

<style scoped>
#mermaid-container {
  min-height: 200px;
  margin: 0 auto;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
}
</style>