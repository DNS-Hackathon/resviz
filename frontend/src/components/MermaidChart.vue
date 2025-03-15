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
  
  mermaid.initialize({ startOnLoad: false });
  console.log(props.step);
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
  margin: 20px 0;
  min-height: 200px;
}
</style>