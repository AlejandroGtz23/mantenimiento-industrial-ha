<template>
  <AppFrame>
    <section class="page-heading mb-6"><div><p class="eyebrow mb-2">TRAZABILIDAD</p><h1 class="text-h4 font-weight-bold">Mantenimientos</h1><p class="text-medium-emphasis mt-2">Consulta evidencias y reportes realizados por los técnicos.</p></div><v-chip prepend-icon="mdi-database-outline" variant="tonal">{{ records.length }} resultados</v-chip></section>
    <v-card class="panel-card mb-5" elevation="0" rounded="xl"><v-card-text class="pa-4"><v-row align="center"><v-col cols="12" md="6"><v-text-field v-model="filters.fecha" type="date" label="Fecha" prepend-inner-icon="mdi-calendar" density="comfortable" hide-details variant="outlined"/></v-col><v-col cols="12" md="6"><v-btn color="primary" prepend-icon="mdi-magnify" size="large" block :loading="loading" @click="load">Buscar mantenimientos</v-btn></v-col></v-row></v-card-text></v-card>
    <v-card class="panel-card" elevation="0" rounded="xl"><v-data-table :headers="headers" :items="records" :loading="loading" item-value="id" hover><template #item.tecnico.nombre="{ item }"><div class="font-weight-medium">{{ item.tecnico?.nombre ?? '—' }}</div></template><template #item.maquina.nombre="{ item }">{{ item.maquina?.nombre ?? '—' }}</template><template #item.fecha_hora="{ item }">{{ new Date(item.fecha_hora).toLocaleString() }}</template><template #item.estado="{ item }"><v-chip :color="statusColor(item.estado)" size="small" variant="tonal">{{ item.estado }}</v-chip></template><template #item.foto_url="{ item }"><v-btn v-if="item.foto_url" color="primary" prepend-icon="mdi-image-outline" size="small" variant="text" @click="selectedPhoto = item.foto_url">Ver evidencia</v-btn><span v-else>—</span></template><template #no-data><div class="pa-10 text-center text-medium-emphasis"><v-icon icon="mdi-clipboard-text-search-outline" size="42" class="mb-3"/><div>Aún no hay registros para estos filtros.</div></div></template></v-data-table></v-card>
    <v-dialog v-model="showPhoto" max-width="650"><v-card rounded="xl" title="Evidencia fotográfica"><v-img :src="selectedPhoto" max-height="580" contain/><v-card-actions class="pa-4"><v-spacer/><v-btn @click="showPhoto = false">Cerrar</v-btn></v-card-actions></v-card></v-dialog>
  </AppFrame>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import AppFrame from '@/components/AppFrame.vue'
import api from '@/services/api'
interface RecordItem { id: string, fecha_hora: string, estado: string, foto_url: string, audio_url: string, tecnico?: { nombre: string }, maquina?: { nombre: string } }
const filters = reactive({ fecha: new Date().toISOString().slice(0, 10) }); const records = ref<RecordItem[]>([]); const loading = ref(false); const selectedPhoto = ref('')
const showPhoto = computed({ get: () => selectedPhoto.value !== '', set: (value: boolean) => { if (!value) selectedPhoto.value = '' } })
const headers = [{ title: 'Técnico', key: 'tecnico.nombre' }, { title: 'Máquina', key: 'maquina.nombre' }, { title: 'Fecha', key: 'fecha_hora' }, { title: 'Estado', key: 'estado' }, { title: '', key: 'foto_url', sortable: false, align: 'end' as const }]
function statusColor(status: string) { return status === 'COMPLETADO' ? 'success' : status === 'CANCELADO' ? 'error' : 'warning' }
async function load() { loading.value = true; try { records.value = (await api.get('/admin/mantenimientos', { params: filters })).data } finally { loading.value = false } }
onMounted(load)
</script>

<style scoped>
.page-heading { align-items: end; display: flex; justify-content: space-between; }.eyebrow { color: rgb(var(--v-theme-primary)); font-size: .75rem; font-weight: 800; letter-spacing: .14em; }.panel-card { border: 1px solid rgba(15, 23, 42, .08); box-shadow: 0 12px 28px rgba(15, 23, 42, .04) !important; }@media (max-width: 600px) { .page-heading { align-items: flex-start; flex-direction: column; } }
</style>
