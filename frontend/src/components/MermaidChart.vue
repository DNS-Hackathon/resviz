<script setup>
import { defineProps, onMounted, watch, nextTick } from "vue";
import mermaid from "mermaid";

const props = defineProps({
  step: Number,
  charts: Array,
});

const handleMermaidClick = (id) => {
  alert(`Node Clicked: ${id}`);
};

const renderMermaid = async () => {
  await nextTick();

  const container = document.getElementById("mermaid-container");
  if (!container) return;

  mermaid.initialize({
    startOnLoad: false,
    theme: "dark",
    securityLevel: "loose",
  });

  let modifiedChart = props.charts[props.step]
    .split("\n")
    .map((line) => {
      const match = line.match(/^(\w[\w.-]*)@\{/); // Extract node ID
      if (match) {
        const nodeId = match[1];
        return `${line}\nclick ${nodeId} "javascript:handleMermaidClick('${nodeId}')"`;
      }
      return line;
    })
    .join("\n");

  const { svg } = await mermaid.render("mermaidChart", props.charts[props.step]);

  container.innerHTML = svg;

  setTimeout(() => {
    const nodes = document.querySelectorAll("svg .node");
    nodes.forEach((node) => {
      const nodeId = node.getAttribute("id");
      node.addEventListener("click", () => {
        handleMermaidClick(nodeId);
      });
    });
  }, 0);
};

onMounted(renderMermaid);
watch(() => props.step, renderMermaid);
watch(() => props.charts, renderMermaid);
</script>

<template>
  <div id="mermaid-container" class="mermaid"></div>
</template>

<style scoped>
#mermaid-container {
  min-height: 400px;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  overflow: visible;
  position: relative;
}

svg {
  display: block;
  margin: auto;
  max-wdith: 100%;
  height: auto !important;
}
</style>
