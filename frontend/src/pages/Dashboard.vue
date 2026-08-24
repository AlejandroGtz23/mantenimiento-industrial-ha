<template>
  <AppFrame>
    <section class="hero mb-7"><div><p class="eyebrow mb-2">CENTRO DE MANTENIMIENTO</p><h1 class="text-h4 text-md-h3 font-weight-bold">Operación de la planta.</h1><p class="text-body-1 text-medium-emphasis mt-2">Monitorea recorridos, activos y técnicos en tiempo real.</p></div><v-btn prepend-icon="mdi-refresh" variant="tonal" color="primary" :loading="loading" @click="load">Actualizar</v-btn></section>
    <v-alert v-if="error" type="error" variant="tonal" class="mb-5">{{ error }}</v-alert>
    <v-row><v-col v-for="card in cards" :key="card.title" cols="12" sm="6" lg="4"><v-card class="stat-card" elevation="0" rounded="xl"><v-card-text class="pa-5"><div class="d-flex justify-space-between align-start"><div><p class="text-body-2 text-medium-emphasis mb-2">{{ card.title }}</p><p class="text-h3 font-weight-bold mb-1">{{ card.value }}</p><span class="text-caption text-medium-emphasis">{{ card.caption }}</span></div><v-avatar :color="card.color" size="46" variant="tonal"><v-icon :icon="card.icon" /></v-avatar></div></v-card-text></v-card></v-col></v-row>
    <v-row class="mt-2"><v-col cols="12" lg="7"><v-card class="panel-card" elevation="0" rounded="xl"><v-card-item title="Mantenimientos por máquina" subtitle="Total histórico por activo"/><v-card-text><v-list><v-list-item v-for="row in dashboard.por_maquina" :key="row.maquina" :title="row.maquina"><template #append><v-chip color="primary" variant="tonal">{{ row.total }}</v-chip></template></v-list-item></v-list></v-card-text></v-card></v-col><v-col cols="12" lg="5"><v-card class="panel-card" elevation="0" rounded="xl"><v-card-item title="Últimos mantenimientos"/><v-card-text><v-list><v-list-item v-for="item in dashboard.ultimos" :key="item.id" :title="item.maquina?.nombre" :subtitle="item.tecnico?.nombre"><template #append><v-chip size="small" color="success">{{ item.estado }}</v-chip></template></v-list-item></v-list></v-card-text></v-card></v-col></v-row>
  </AppFrame>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import AppFrame from '@/components/AppFrame.vue'
import api from '@/services/api'
const dashboard = reactive({ mantenimientos_hoy: 0, maquinas_criticas: 0, tecnicos_activos: 0, por_maquina: [] as Array<{ maquina: string, total: number }>, ultimos: [] as any[] }); const error = ref(''); const loading = ref(false)
const cards = computed(() => [{ title: 'Mantenimientos hoy', value: dashboard.mantenimientos_hoy, caption: 'Reportes completados', icon: 'mdi-clipboard-check-outline', color: 'primary' }, { title: 'Máquinas críticas', value: dashboard.maquinas_criticas, caption: 'Fuera de servicio', icon: 'mdi-alert-octagon-outline', color: 'error' }, { title: 'Técnicos activos', value: dashboard.tecnicos_activos, caption: 'Disponibles para recorrido', icon: 'mdi-account-hard-hat-outline', color: 'success' }])
async function load() { loading.value = true; error.value = ''; try { Object.assign(dashboard, (await api.get('/admin/dashboard')).data) } catch { error.value = 'No fue posible cargar las estadísticas.' } finally { loading.value = false } }; onMounted(load)
</script>

<style scoped>
.hero { align-items: end; display: flex; gap: 20px; justify-content: space-between; }.eyebrow { color: rgb(var(--v-theme-primary)); font-size: .75rem; font-weight: 800; letter-spacing: .14em; }.stat-card, .panel-card { border: 1px solid rgba(15, 23, 42, .08); box-shadow: 0 12px 28px rgba(15, 23, 42, .04) !important; }.status-row { align-items: center; display: flex; gap: 12px; padding: 13px 0; }.status-row + .status-row { border-top: 1px solid rgba(15, 23, 42, .08); }@media (max-width: 600px) { .hero { align-items: flex-start; flex-direction: column; } }
</style>
