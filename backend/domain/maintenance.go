package domain

import "time"

const (
	EstadoOperativa     = "OPERATIVA"
	EstadoMantenimiento = "MANTENIMIENTO"
	EstadoFueraServicio = "FUERA_SERVICIO"
	RegistroCompletado  = "COMPLETADO"
	RegistroPendiente   = "PENDIENTE"
	RegistroCancelado   = "CANCELADO"
)

type Tecnico struct {
	ID             string    `gorm:"primaryKey;size:36" json:"id"`
	NumeroEmpleado string    `gorm:"uniqueIndex;not null;size:64" json:"numero_empleado"`
	Nombre         string    `gorm:"not null;size:160" json:"nombre"`
	Area           string    `gorm:"size:100" json:"area"`
	Activo         bool      `json:"activo"`
	CreatedAt      time.Time `json:"created_at"`
}

func (Tecnico) TableName() string { return "mt_tecnicos" }

type Maquina struct {
	ID                      string    `gorm:"primaryKey;size:36" json:"id"`
	Codigo                  string    `gorm:"uniqueIndex;not null;size:64" json:"codigo"`
	Nombre                  string    `gorm:"not null;size:160" json:"nombre"`
	Ubicacion               string    `gorm:"size:160" json:"ubicacion"`
	Estado                  string    `gorm:"size:32" json:"estado"`
	FrecuenciaMantenimiento string    `gorm:"size:32" json:"frecuencia_mantenimiento"`
	CreatedAt               time.Time `json:"created_at"`
}

func (Maquina) TableName() string { return "mt_maquinas" }

type ItemChecklist struct {
	ID          string    `gorm:"primaryKey;size:36" json:"id"`
	MaquinaID   string    `gorm:"index;size:36" json:"maquina_id"`
	Descripcion string    `gorm:"not null;size:500" json:"descripcion"`
	Orden       int       `json:"orden"`
	CreatedAt   time.Time `json:"created_at"`
}

func (ItemChecklist) TableName() string { return "mt_items_checklist" }

type RegistroMantenimiento struct {
	ID            string                 `gorm:"primaryKey;size:36" json:"id"`
	TecnicoID     string                 `gorm:"index;size:36" json:"tecnico_id"`
	Tecnico       Tecnico                `gorm:"foreignKey:TecnicoID" json:"tecnico,omitempty"`
	MaquinaID     string                 `gorm:"index;size:36" json:"maquina_id"`
	Maquina       Maquina                `gorm:"foreignKey:MaquinaID" json:"maquina,omitempty"`
	FechaHora     time.Time              `gorm:"index" json:"fecha_hora"`
	FotoURL       string                 `gorm:"size:512" json:"foto_url"`
	AudioURL      string                 `gorm:"size:512" json:"audio_url"`
	Observaciones string                 `gorm:"type:text" json:"observaciones"`
	Estado        string                 `gorm:"size:32" json:"estado"`
	Detalles      []DetalleMantenimiento `gorm:"foreignKey:RegistroMantenimientoID" json:"detalles,omitempty"`
}

func (RegistroMantenimiento) TableName() string { return "mt_registros" }

type DetalleMantenimiento struct {
	ID                      string `gorm:"primaryKey;size:36" json:"id"`
	RegistroMantenimientoID string `gorm:"index;size:36" json:"registro_mantenimiento_id"`
	ItemChecklistID         string `gorm:"index;size:36" json:"item_checklist_id"`
	Cumple                  bool   `json:"cumple"`
	ObservacionItem         string `gorm:"type:text" json:"observacion_item"`
}

func (DetalleMantenimiento) TableName() string { return "mt_detalles" }
