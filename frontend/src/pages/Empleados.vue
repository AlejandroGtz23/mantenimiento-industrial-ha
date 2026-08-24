<template>
  <AppFrame>
    <section class="page-heading mb-6"><div><p class="eyebrow mb-2">DIRECTORIO</p><h1 class="text-h4 font-weight-bold">Empleados</h1><p class="text-medium-emphasis mt-2">Administra las credenciales de acceso de la planta.</p></div><v-btn color="primary" prepend-icon="mdi-account-plus-outline" size="large" @click="openCreate">Nuevo empleado</v-btn></section>
    <v-alert v-if="error" type="error" variant="tonal" class="mb-4">{{ error }}</v-alert>
    <v-card class="panel-card" elevation="0" rounded="xl"><v-data-table :headers="headers" :items="employees" :loading="loading" item-value="id" hover><template #item.nombre="{ item }"><div class="d-flex align-center ga-3"><v-avatar color="primary" variant="tonal">{{ initials(item.nombre) }}</v-avatar><div><div class="font-weight-medium">{{ item.nombre }}</div><div class="text-caption text-medium-emphasis">{{ item.numero_empleado }}</div></div></div></template><template #item.area="{ item }"><v-chip size="small" variant="tonal">{{ item.area || 'Sin área' }}</v-chip></template><template #item.actions="{ item }"><v-btn icon="mdi-pencil-outline" variant="text" @click="openEdit(item)"/><v-btn icon="mdi-delete-outline" color="error" variant="text" @click="remove(item.id)"/></template><template #no-data><div class="pa-10 text-center text-medium-emphasis"><v-icon icon="mdi-account-group-outline" size="42" class="mb-3"/><div>No hay empleados registrados.</div></div></template></v-data-table></v-card>
    <v-dialog v-model="dialog" max-width="650" persistent><v-card rounded="xl"><v-card-item :title="editing ? 'Editar empleado' : 'Nuevo empleado'" :subtitle="editing ? 'Actualiza los datos de la credencial.' : 'Registra una nueva credencial de acceso.'"/><v-divider/><v-card-text class="pa-6"><v-row><v-col cols="12" sm="6"><v-text-field v-model="form.numero_empleado" label="Número de empleado" variant="outlined"/></v-col><v-col cols="12" sm="6"><v-text-field v-model="form.nombre" label="Nombre completo" variant="outlined"/></v-col><v-col cols="12" sm="6"><v-text-field v-model="form.area" label="Área" variant="outlined"/></v-col><v-col cols="12" sm="6"><v-text-field v-model="form.rfid_tag" label="Etiqueta RFID" variant="outlined"/></v-col><v-col cols="12"><v-text-field v-model="form.foto_referencia" label="URL de foto de referencia (opcional)" variant="outlined"/></v-col></v-row></v-card-text><v-card-actions class="pa-5"><v-spacer/><v-btn @click="dialog = false">Cancelar</v-btn><v-btn color="primary" :loading="saving" @click="save">Guardar empleado</v-btn></v-card-actions></v-card></v-dialog>
  </AppFrame>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import AppFrame from '@/components/AppFrame.vue'
import api from '@/services/api'
interface Employee { id: string, numero_empleado: string, nombre: string, area: string, rfid_tag: string, foto_referencia: string }
const employees = ref<Employee[]>([]); const loading = ref(false); const saving = ref(false); const dialog = ref(false); const editing = ref(false); const error = ref(''); const form = reactive<Employee>({ id: '', numero_empleado: '', nombre: '', area: '', rfid_tag: '', foto_referencia: '' })
const headers = [{ title: 'Empleado', key: 'nombre' }, { title: 'Área', key: 'area' }, { title: 'RFID', key: 'rfid_tag' }, { title: '', key: 'actions', sortable: false, align: 'end' as const }]
function reset() { Object.assign(form, { id: '', numero_empleado: '', nombre: '', area: '', rfid_tag: '', foto_referencia: '' }) }; function openCreate() { reset(); editing.value = false; dialog.value = true }; function openEdit(item: Employee) { Object.assign(form, item); editing.value = true; dialog.value = true }; function initials(name: string) { return name.split(' ').filter(Boolean).slice(0, 2).map((word) => word[0]).join('').toUpperCase() }
async function load() { loading.value = true; try { employees.value = (await api.get('/admin/empleados')).data.data } finally { loading.value = false } }
async function save() { saving.value = true; error.value = ''; try { if (editing.value) await api.put(`/admin/empleados/${form.id}`, form); else await api.post('/admin/empleados', form); dialog.value = false; await load() } catch (err: any) { error.value = err.response?.data?.error ?? 'No se pudo guardar el empleado.' } finally { saving.value = false } }
async function remove(id: string) { if (!confirm('¿Eliminar este empleado?')) return; try { await api.delete(`/admin/empleados/${id}`); await load() } catch { error.value = 'No se pudo eliminar el empleado.' } }
onMounted(load)
</script>

<style scoped>
.page-heading { align-items: end; display: flex; justify-content: space-between; }.eyebrow { color: rgb(var(--v-theme-primary)); font-size: .75rem; font-weight: 800; letter-spacing: .14em; }.panel-card { border: 1px solid rgba(15, 23, 42, .08); box-shadow: 0 12px 28px rgba(15, 23, 42, .04) !important; }@media (max-width: 600px) { .page-heading { align-items: flex-start; flex-direction: column; } }
</style>
