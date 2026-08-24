<template>
  <div class="chart-wrapper"><canvas ref="canvas" /></div>
</template>

<script setup lang="ts">
import { Chart, type ChartConfiguration, LineController, LineElement, PointElement, LinearScale, CategoryScale, Tooltip, Filler } from 'chart.js'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = defineProps<{ points: Array<{ hora: number, total: number }> }>()
const canvas = ref<HTMLCanvasElement>()
let chart: Chart | undefined
Chart.register(LineController, LineElement, PointElement, LinearScale, CategoryScale, Tooltip, Filler)
function draw() {
  if (!canvas.value) return
  chart?.destroy()
  const config: ChartConfiguration<'line'> = { type: 'line', data: { labels: props.points.map((item) => `${String(item.hora).padStart(2, '0')}:00`), datasets: [{ label: 'Ingresos', data: props.points.map((item) => item.total), borderColor: '#1976D2', backgroundColor: 'rgba(25, 118, 210, .14)', fill: true, tension: .3 }] }, options: { responsive: true, maintainAspectRatio: false, plugins: { tooltip: { intersect: false } }, scales: { y: { beginAtZero: true, ticks: { precision: 0 } } } } }
  chart = new Chart(canvas.value, config)
}
onMounted(draw)
watch(() => props.points, draw, { deep: true })
onBeforeUnmount(() => chart?.destroy())
</script>

<style scoped>.chart-wrapper { height: 300px; }</style>
