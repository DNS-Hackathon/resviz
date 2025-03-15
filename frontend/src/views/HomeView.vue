<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import MermaidChart from '@/components/MermaidChart.vue'
import ChartNavigator from '@/components/ChartNavigator.vue'
import SearchBar from '@/components/SearchBar.vue'

const step = ref(0);

const charts = [
  `graph TD;
a.root-servers.net@{ shape: rect, label: "a.root-servers.net"}
  `,
  `graph TD;
a.root-servers.net@{ shape: rect, label: "a.root-servers.net"}
c.gtld-servers.net@{ shape: rect, label: "c.gtld-servers.net"}
a.root-servers.net --> c.gtld-servers.net
  `,
  `graph TD;
a.root-servers.net@{ shape: rect, label: "a.root-servers.net"}
c.gtld-servers.net@{ shape: rect, label: "c.gtld-servers.net"}
ns2.dns-oarc.net@{ shape: rect, label: "ns2.dns-oarc.net"}
a.root-servers.net --> c.gtld-servers.net
c.gtld-servers.net --> ns2.dns-oarc.net

  `,
  `graph TD;
a.root-servers.net@{ shape: rect, label: "a.root-servers.net"}
c.gtld-servers.net@{ shape: rect, label: "c.gtld-servers.net"}
ns2.dns-oarc.net@{ shape: rect, label: "ns2.dns-oarc.net"}
ns.dev.dns-oarc.net@{ shape: rect, label: "ns.dev.dns-oarc.net"}


a.root-servers.net --> c.gtld-servers.net
c.gtld-servers.net --> ns2.dns-oarc.net
ns2.dns-oarc.net --> ns.dev.dns-oarc.net
  `,
  `graph TD;
a.root-servers.net@{ shape: rect, label: "a.root-servers.net"}
c.gtld-servers.net@{ shape: rect, label: "c.gtld-servers.net"}
ns2.dns-oarc.net@{ shape: rect, label: "ns2.dns-oarc.net"}
ns.dev.dns-oarc.net@{ shape: rect, label: "ns.dev.dns-oarc.net"}
ns.cmdns.dev.dns-oarc.net@{ shape: rect, label: "ns.cmdns.dev.dns-oarc.net"}


a.root-servers.net --> c.gtld-servers.net
c.gtld-servers.net --> ns2.dns-oarc.net
ns2.dns-oarc.net --> ns.dev.dns-oarc.net
ns.dev.dns-oarc.net --> ns.cmdns.dev.dns-oarc.net
  `,

];

const nextStep = () => {
  if (step.value < charts.length -1) step.value++;
};

const prevStep = () => {
  if (step.value > 0) step.value--;
};

const handleSearch = (query) => {
  console.log(query);
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
    <ChartNavigator :step="step" :maxSteps="charts.length" @next="nextStep" @prev="prevStep" />
    <MermaidChart :step="step" :charts="charts"/>
  </main>
</template>
