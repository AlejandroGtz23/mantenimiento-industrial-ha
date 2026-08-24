// Importaciones necesarias
import type { App } from 'vue'
import Swal from 'sweetalert2'

// ============================================================
// 1. Función autoinvocada para mover diálogos (drag & drop)
// ============================================================
interface DragState {
   el?: HTMLElement
   handle?: HTMLElement | null
   mouseStartX?: number
   mouseStartY?: number
   elStartX?: number
   elStartY?: number
   oldTransition?: string
}

(function () {
   const d: DragState = {}

   const isMovable = (targ: EventTarget | null): boolean => {
      if (!(targ instanceof HTMLElement)) return false
      return targ.classList?.contains('vss-movable') ?? false
   }

   document.addEventListener('mousedown', (e: MouseEvent) => {
      const target = e.target as HTMLElement
      const closestDialog = target.closest('.v-overlay__content') as HTMLElement | null
      const title = closestDialog?.querySelector('.v-toolbar-title') as HTMLElement | null

      if (
         e.button === 0 &&
         closestDialog != null &&
         (isMovable(target) || isMovable(target.parentNode))
      ) {
         d.el = closestDialog
         d.handle = title
         d.mouseStartX = e.clientX
         d.mouseStartY = e.clientY
         d.elStartX = d.el.getBoundingClientRect().left
         d.elStartY = d.el.getBoundingClientRect().top
         d.el.style.position = 'fixed'
         d.el.style.margin = '0'
         d.oldTransition = d.el.style.transition
         d.el.style.transition = 'none'
      }
   })

   document.addEventListener('mousemove', (e: MouseEvent) => {
      if (d.el === undefined) return
      if (d.handle === undefined) return // handle necesario para el cálculo de top
      const rect = d.el.getBoundingClientRect()
      d.el.style.left =
         Math.min(
            Math.max(d.elStartX! + e.clientX - d.mouseStartX!, 0),
            window.innerWidth - rect.width
         ) + 'px'
      d.el.style.top =
         Math.min(
            Math.max(d.elStartY! + e.clientY - d.mouseStartY!, 0),
            window.innerHeight - d.handle.getBoundingClientRect().height
         ) + 'px'
   })

   document.addEventListener('mouseup', () => {
      if (d.el === undefined) return
      d.el.style.transition = d.oldTransition || ''
      d.el = undefined
      d.handle = undefined
   })
})()

// ============================================================
// 2. Configuración de propiedades globales de la app Vue
// ============================================================
// Asumimos que 'app' es la instancia de Vue (App) ya creada.
// En tu archivo principal (main.ts) tendrás algo como:
// const app = createApp(App)
// Luego aplicas estas configuraciones.
declare const app: App

// Constante para JWT
app.config.globalProperties.$jwt = ''

// Collator para ordenamiento en español
app.config.globalProperties.$collator = new Intl.Collator('es')

// Fetch genérico (texto o JSON)
app.config.globalProperties.$fetch = async (
   input: RequestInfo | URL,
   init?: RequestInit
): Promise<string | any> => {
   const response = await fetch(input, init)
   if (!response.ok) throw new Error(await response.text())
   const contentType = response.headers.get('content-type') || ''
   return contentType.includes('text/plain')
      ? await response.text()
      : await response.json()
}

// Fetch de imagen como blob y creación de URL
app.config.globalProperties.$fetchImage = async (
   input: RequestInfo | URL,
   init?: RequestInit
): Promise<string> => {
   const response = await fetch(input, init)
   if (!response.ok) throw new Error(await response.text())
   return URL.createObjectURL(await response.blob())
}

// Toast con SweetAlert2 (duración 3500ms)
app.config.globalProperties.$toast = Swal.mixin({
   toast: true,
   position: 'top-end',
   showConfirmButton: false,
   timer: 3500,
   timerProgressBar: true,
   didOpen: (evtToast: HTMLElement) => {
      evtToast.addEventListener('mouseenter', Swal.stopTimer)
      evtToast.addEventListener('mouseleave', Swal.resumeTimer)
   }
})

// Toast con duración 2300ms
app.config.globalProperties.$toast2300 = Swal.mixin({
   toast: true,
   position: 'top-end',
   showConfirmButton: false,
   timer: 2300,
   timerProgressBar: true,
   didOpen: (evtToast: HTMLElement) => {
      evtToast.addEventListener('mouseenter', Swal.stopTimer)
      evtToast.addEventListener('mouseleave', Swal.resumeTimer)
   }
})

// Notificación genérica
app.config.globalProperties.$showNotif = (
   icon: any,
   title: string,
   text: string
) => {
   Swal.fire({ icon, title, text })
}

// Notificación de error
app.config.globalProperties.$showError = (title: string, text: string) => {
   Swal.fire({ icon: 'error', title, text })
}

// Confirmación con SweetAlert2
app.config.globalProperties.$swalConfirm = async (
   title: string,
   icon: any,
   text: string
): Promise<any> => {
   return await Swal.fire({
      title,
      icon,
      text,
      showCancelButton: true,
      confirmButtonColor: '#070047',
      cancelButtonColor: '#d8d7e3',
      confirmButtonText: 'Confirmar',
      cancelButtonText: 'Cancelar',
      reverseButtons: true
   })
}

// Convertir array de objetos a CSV
app.config.globalProperties.$dataToCSV = (data: Record<string, any>[]): string => {
   const csvRows: string[] = []
   const headers = Object.keys(data[0])
   csvRows.push(headers.join(','))
   for (const row of data) {
      const values = headers.map(header => {
         const escaped = ('' + row[header])
            .normalize('NFD')
            .replace(/[\u0300-\u036f]/g, '')
         return `"${escaped}"`
      })
      csvRows.push(values.join(','))
   }
   return csvRows.join('\n')
}

// Descargar archivo CSV
app.config.globalProperties.$downloadCSV = (data: string, file: string) => {
   const blob = new Blob([data], { type: 'text/csv' })
   const url = window.URL.createObjectURL(blob)
   const a = document.createElement('a')
   a.setAttribute('href', url)
   a.setAttribute('download', file)
   document.body.appendChild(a)
   a.click()
   document.body.removeChild(a)
   window.URL.revokeObjectURL(url)
}

// Formatear fecha de YYYY-MM-DD a DD/MM/YYYY
app.config.globalProperties.$ddMMyyFormat = (strDate: string): string => {
   const arr = strDate.split('-')
   return `${arr[2]}/${arr[1]}/${arr[0]}`
}

// Descargar PDF desde una API
app.config.globalProperties.$fetchDownloadPdf = async (
   input: RequestInfo | URL,
   init?: RequestInit
) => {
   const response = await fetch(input, init)
   if (!response.ok) throw new Error(await response.text())
   const url = URL.createObjectURL(await response.blob())
   const a = document.createElement('a')
   a.setAttribute('download', 'Reporte-sistema-dat.pdf')
   a.setAttribute('href', url)
   a.click()
   URL.revokeObjectURL(url)
}

// Manejo de imágenes en CRUD (subida y eliminación)
app.config.globalProperties.$crudImage = function (
   this: any,
   id: string | number,
   result: string,
   folder: string,
   fileList: FileList | File[],
   deletedImages: string[]
) {
   if (fileList.length > 0 && result !== 'failed') {
      const formData = new FormData()
      Array.from(fileList).forEach((file: File) => {
         formData.append('savedImages', file)
      })
      this.$fetch(`/api/image-multiple?tbl=${folder}&fld=mini&id=${id}`, {
         method: 'POST',
         body: formData
      })
         .then((result: any) => console.log(result))
         .catch((error: any) => this.$toast.fire({ icon: 'error', text: error }))
   }

   if (deletedImages.length > 0) {
      const formData = new FormData()
      deletedImages.forEach(imageUrl => {
         const urlParams = new URLSearchParams(imageUrl.split('?')[1])
         formData.append('deletedImages', urlParams.get('url') || '')
      })
      this.$fetch(`/api/delete-images`, {
         method: 'POST',
         body: formData
      })
         .then((result: string) => {
            if (result === 'success') {
               console.log('Images deleted successfully')
            } else {
               this.$toast.fire({ icon: 'error', text: 'No fue posible eliminar las imágenes' })
            }
         })
         .catch((error: any) =>
            this.$toast.fire({ icon: 'error', text: 'Ocurrió un problema al intentar eliminar las imágenes' })
         )
   }
}

// Detección de móvil usando Vuetify
app.config.globalProperties.$isMobile = function (this: any): boolean {
   return this.$vuetify.display.smAndDown
}

// Añadir archivos nuevos a FormData con metadatos
app.config.globalProperties.$appendNewFilesToFormData = function (
   formData: FormData,
   files: File | File[] | null | undefined,
   object: Record<string, any>
) {
   if (!files) return
   const fileArray = Array.isArray(files) ? files : [files]
   fileArray.forEach(file => {
      formData.append('images[]', file)
      formData.append('metadata[]', JSON.stringify(object))
   })
}

// Añadir URLs de imágenes eliminadas a FormData
app.config.globalProperties.$appendDeletedFilesToFormData = function (
   formData: FormData,
   files: string | string[] | null | undefined
) {
   if (!files) return
   const fileArray = Array.isArray(files) ? files : [files]
   fileArray.forEach(imageUrl => {
      const urlParams = new URLSearchParams(imageUrl.split('?')[1])
      formData.append('deletedImages[]', urlParams.get('url') || '')
   })
}

// ============================================================
// (Opcional) Aumento de tipos para Vue
// ============================================================
declare module '@vue/runtime-core' {
   export interface ComponentCustomProperties {
      $jwt: string
      $collator: Intl.Collator
      $fetch: (input: RequestInfo | URL, init?: RequestInit) => Promise<any>
      $fetchImage: (input: RequestInfo | URL, init?: RequestInit) => Promise<string>
      $toast: any // Podría tiparse con SweetAlert2
      $toast2300: any
      $showNotif: (icon: any, title: string, text: string) => void
      $showError: (title: string, text: string) => void
      $swalConfirm: (title: string, icon: any, text: string) => Promise<any>
      $dataToCSV: (data: Record<string, any>[]) => string
      $downloadCSV: (data: string, file: string) => void
      $ddMMyyFormat: (strDate: string) => string
      $fetchDownloadPdf: (input: RequestInfo | URL, init?: RequestInit) => Promise<void>
      $crudImage: (
         id: string | number,
         result: string,
         folder: string,
         fileList: FileList | File[],
         deletedImages: string[]
      ) => void
      $isMobile: () => boolean
      $appendNewFilesToFormData: (
         formData: FormData,
         files: File | File[] | null | undefined,
         object: Record<string, any>
      ) => void
      $appendDeletedFilesToFormData: (
         formData: FormData,
         files: string | string[] | null | undefined
      ) => void
   }
}