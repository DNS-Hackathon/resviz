<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import axios from 'axios';
import qs from 'qs';
import MermaidChart from '@/components/MermaidChart.vue'
import ChartNavigator from '@/components/ChartNavigator.vue'
import SearchBar from '@/components/SearchBar.vue'
import Modal from '@/components/Modal.vue'

const step = ref(0);
const charts = ref([]);
const isLoading = ref(null);
const showModal = ref(false);
const modalContent = ref("This is some fun information about the box you clicked on!");

const fakeData = [
`
graph TD;
i.root-servers.net@{ shape: rect, label: "i.root-servers.net"}
`,
`
graph TD;
i.root-servers.net@{ shape: rect, label: "i.root-servers.net"}
l.gtld-servers.net@{ shape: rect, label: "l.gtld-servers.net"}


i.root-servers.net --> l.gtld-servers.net
`,
`
graph TD;
i.root-servers.net@{ shape: rect, label: "i.root-servers.net"}
l.gtld-servers.net@{ shape: rect, label: "l.gtld-servers.net"}
ns3.afrinic.net@{ shape: rect, label: "ns3.afrinic.net"}


i.root-servers.net --> l.gtld-servers.net
l.gtld-servers.net --> ns3.afrinic.net
`,
];

const nextStep = () => {
  if (step.value < charts.value.length -1) step.value++;
};

const prevStep = () => {
  if (step.value > 0) step.value--;
};

const openModal = (content) => {
  modalContent.value = content;
  showModal.value = true;
};

const closeModal = () => {
  showModal.value = false;
};

const handleSearch = async (query) => {
  try {
    isLoading.value = true;
    if (query === "ripe") {
      charts.value = fakeData;
      return
    }
    const response = await axios.post(
      'http://nerdig.examples.nu/resviz',
      qs.stringify({ domain: query }),
      { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } }
    );

    if (response.status !== 200) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

//    const contentType = response.headers["content-type"];

//    if (contentType && contentType.includes("application/json")) {
//    } else {
//      charts.value = [response.mermaid];
//    } else {
//      console.log("None json recieved");
//      charts.value = [response.data];
//    };

    let parsedData;

    try {
      const rawData = typeof response.data === "string" ? JSON.parse(response.data) : response.data;
      parsedData = rawData?.mermaid ?? ["Invalid data"];
    } catch (error) {
      console.log("JSON Parsing error");
      parsedData = [response.data];
    }

  charts.value = parsedData;

  } catch (err) {
    alert(err);
  } finally {
    isLoading.value = false;
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
    
    <div v-if="isLoading" class="loading-spinner">
      <div class="spinner"></div>
      <p>Loading resviz data...</p>
    </div>
    
    <div v-if="charts.length > 0">
      <ChartNavigator :step="step" :maxSteps="charts.length" @next="nextStep" @prev="prevStep" />
      <MermaidChart :step="step" :charts="charts" @openModal="openModal"/>
      <Modal :isOpen="showModal" @close="closeModal">
        <h2>Deep dive into the world of DNS!</h2>
	<p>{{ modalContent }}</p>
      </Modal>
    </div>
  </main>
</template>

<style scoped>
.loading-spinner {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  height: 100vh; /* Center the spinner vertically */
}

.spinner {
  border: 4px solid #f3f3f3; /* Light background */
  border-top: 4px solid #3498db; /* Blue color */
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 2s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>