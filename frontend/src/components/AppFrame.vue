<template>
  <v-app class="access-app">
    <v-navigation-drawer v-model="drawer" class="access-drawer" color="#101827" elevation="0" theme="dark">
      <div class="brand pa-6"><div class="brand-mark"><v-icon icon="mdi-shield-check" size="26" /></div><div><div class="text-subtitle-1 font-weight-bold">Acceso seguro</div><div class="text-caption text-medium-emphasis">Planta industrial</div></div></div>
      <v-divider class="mx-4" />
      <v-list nav class="pa-4"><v-list-subheader class="px-3 text-caption">OPERACIÓN</v-list-subheader><v-list-item v-for="item in menuItems" :key="item.to" :prepend-icon="item.icon" :title="item.text" :to="item.to" rounded="lg" class="mb-2" /></v-list>
      <template #append><div class="pa-4"><v-card color="rgba(255,255,255,.07)" elevation="0" rounded="lg"><v-card-text class="text-caption"><v-icon icon="mdi-shield-lock-outline" class="mr-2" />Sesión protegida</v-card-text></v-card></div></template>
    </v-navigation-drawer>
    <v-app-bar class="access-bar" color="transparent" elevation="0"><v-app-bar-nav-icon class="d-md-none" @click="drawer = !drawer" /><v-spacer /><v-chip color="success" prepend-icon="mdi-circle" size="small" variant="tonal" class="mr-3">Sistema activo</v-chip><v-btn icon="mdi-logout" variant="text" aria-label="Cerrar sesión" @click="logout" /></v-app-bar>
    <v-main><v-container class="pa-4 pa-md-8" fluid><slot /></v-container></v-main>
  </v-app>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter(); const drawer = ref(true)
const menuItems = [{ icon: 'mdi-view-dashboard-outline', text: 'Resumen', to: '/' }, { icon: 'mdi-clipboard-check-outline', text: 'Mantenimientos', to: '/mantenimientos' }, { icon: 'mdi-cog-outline', text: 'Máquinas', to: '/maquinas' }, { icon: 'mdi-account-hard-hat-outline', text: 'Técnicos', to: '/tecnicos' }]
function logout() { localStorage.removeItem('access-admin-token'); router.push('/login') }
</script>

<style scoped>
.access-drawer { border: 0 !important; }.brand { align-items: center; display: flex; gap: 12px; }.brand-mark { align-items: center; background: linear-gradient(135deg, #38bdf8, #2563eb); border-radius: 12px; display: flex; height: 46px; justify-content: center; width: 46px; }
</style>
