package repository

import (
	"fiber3-app/backend/database"
	"fiber3-app/backend/domain"
	"gorm.io/gorm"
	"time"
)

type MaintenanceRepository struct{ db *gorm.DB }

func NewMaintenanceRepository() *MaintenanceRepository { return &MaintenanceRepository{database.DB} }
func (r *MaintenanceRepository) Technician(number string) (*domain.Tecnico, error) {
	var v domain.Tecnico
	err := r.db.Where("UPPER(numero_empleado) = ? AND activo = ?", number, true).First(&v).Error
	return &v, err
}
func (r *MaintenanceRepository) Machines() ([]domain.Maquina, error) {
	var v []domain.Maquina
	return v, r.db.Order("nombre").Find(&v).Error
}
func (r *MaintenanceRepository) Machine(id string) (*domain.Maquina, error) {
	var v domain.Maquina
	err := r.db.First(&v, "id = ?", id).Error
	return &v, err
}
func (r *MaintenanceRepository) Checklist(id string) ([]domain.ItemChecklist, error) {
	var v []domain.ItemChecklist
	return v, r.db.Where("maquina_id = ?", id).Order("orden").Find(&v).Error
}
func (r *MaintenanceRepository) CreateChecklistItems(items []domain.ItemChecklist) error {
	return r.db.Create(&items).Error
}
func (r *MaintenanceRepository) CreateChecklistItem(item *domain.ItemChecklist) error {
	return r.db.Create(item).Error
}
func (r *MaintenanceRepository) UpdateChecklistItem(item *domain.ItemChecklist) error {
	return r.db.Model(&domain.ItemChecklist{}).Where("id = ?", item.ID).Updates(map[string]any{"descripcion": item.Descripcion, "orden": item.Orden}).Error
}
func (r *MaintenanceRepository) DeleteChecklistItem(id string) error {
	return r.db.Delete(&domain.ItemChecklist{}, "id = ?", id).Error
}
func (r *MaintenanceRepository) Save(reg *domain.RegistroMantenimiento, details []domain.DetalleMantenimiento) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(reg).Error; err != nil {
			return err
		}
		for i := range details {
			if err := tx.Create(&details[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
func (r *MaintenanceRepository) History(technician string) ([]domain.RegistroMantenimiento, error) {
	var v []domain.RegistroMantenimiento
	return v, r.db.Preload("Maquina").Where("tecnico_id = ?", technician).Order("fecha_hora desc").Find(&v).Error
}
func (r *MaintenanceRepository) Records(date, machine, technician string) ([]domain.RegistroMantenimiento, error) {
	q := r.db.Preload("Tecnico").Preload("Maquina").Order("fecha_hora desc")
	if date != "" {
		start, e := time.Parse("2006-01-02", date)
		if e != nil {
			return nil, e
		}
		q = q.Where("fecha_hora >= ? AND fecha_hora < ?", start, start.AddDate(0, 0, 1))
	}
	if machine != "" {
		q = q.Where("maquina_id = ?", machine)
	}
	if technician != "" {
		q = q.Where("tecnico_id = ?", technician)
	}
	var v []domain.RegistroMantenimiento
	return v, q.Find(&v).Error
}
func (r *MaintenanceRepository) Record(id string) (*domain.RegistroMantenimiento, error) {
	var v domain.RegistroMantenimiento
	err := r.db.Preload("Tecnico").Preload("Maquina").Preload("Detalles").First(&v, "id = ?", id).Error
	return &v, err
}
func (r *MaintenanceRepository) MachinesPage() ([]domain.Maquina, error) { return r.Machines() }
func (r *MaintenanceRepository) CreateMachine(v *domain.Maquina) error   { return r.db.Create(v).Error }
func (r *MaintenanceRepository) UpdateMachine(v *domain.Maquina) error {
	return r.db.Model(&domain.Maquina{}).Where("id = ?", v.ID).Updates(v).Error
}
func (r *MaintenanceRepository) DeleteMachine(id string) error {
	return r.db.Delete(&domain.Maquina{}, "id = ?", id).Error
}
func (r *MaintenanceRepository) Technicians() ([]domain.Tecnico, error) {
	var v []domain.Tecnico
	return v, r.db.Order("nombre").Find(&v).Error
}
func (r *MaintenanceRepository) CreateTechnician(v *domain.Tecnico) error {
	return r.db.Create(v).Error
}
func (r *MaintenanceRepository) Dashboard() (int64, int64, int64, []map[string]any, error) {
	today := time.Now().Truncate(24 * time.Hour)
	var total, critical, active int64
	if e := r.db.Model(&domain.RegistroMantenimiento{}).Where("fecha_hora >= ?", today).Count(&total).Error; e != nil {
		return 0, 0, 0, nil, e
	}
	r.db.Model(&domain.Maquina{}).Where("estado = ?", domain.EstadoFueraServicio).Count(&critical)
	r.db.Model(&domain.Tecnico{}).Where("activo = ?", true).Count(&active)
	var rows []struct {
		Name  string
		Total int64
	}
	e := r.db.Raw("SELECT m.nombre AS name, COUNT(r.id) AS total FROM mt_maquinas m LEFT JOIN mt_registros r ON r.maquina_id=m.id GROUP BY m.nombre ORDER BY total DESC").Scan(&rows).Error
	chart := []map[string]any{}
	for _, x := range rows {
		chart = append(chart, map[string]any{"maquina": x.Name, "total": x.Total})
	}
	return total, critical, active, chart, e
}
