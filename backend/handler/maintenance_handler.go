package handler

import (
	"encoding/json"
	"errors"
	"fiber3-app/backend/config"
	"fiber3-app/backend/domain"
	"fiber3-app/backend/middleware"
	"fiber3-app/backend/repository"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MaintenanceHandler struct {
	repo *repository.MaintenanceRepository
	auth *repository.AuthRepository
}

func NewMaintenanceHandler() *MaintenanceHandler {
	return &MaintenanceHandler{repository.NewMaintenanceRepository(), repository.NewAuthRepository()}
}
func (h *MaintenanceHandler) MobileLogin(c fiber.Ctx) error {
	var req struct {
		NumeroEmpleado string `json:"numero_empleado"`
	}
	if e := c.Bind().Body(&req); e != nil || strings.TrimSpace(req.NumeroEmpleado) == "" {
		return bad(c, "Número de empleado obligatorio")
	}
	t, e := h.repo.Technician(strings.ToUpper(strings.TrimSpace(req.NumeroEmpleado)))
	if e != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Técnico no encontrado o inactivo"})
	}
	token, e := middleware.CreateJWT(t.ID, t.NumeroEmpleado, "TECNICO")
	if e != nil {
		return internal(c)
	}
	return c.JSON(fiber.Map{"token": token, "tecnico": t})
}
func (h *MaintenanceHandler) Machines(c fiber.Ctx) error {
	v, e := h.repo.Machines()
	if e != nil {
		return internal(c)
	}
	return c.JSON(v)
}
func (h *MaintenanceHandler) Checklist(c fiber.Ctx) error {
	v, e := h.repo.Checklist(c.Params("id"))
	if e != nil {
		return internal(c)
	}
	if len(v) == 0 {
		defaults := []string{"Verificar estado general de la máquina", "Verificar ruidos o vibraciones anormales", "Verificar limpieza y seguridad del área"}
		v = make([]domain.ItemChecklist, len(defaults))
		for i, description := range defaults {
			v[i] = domain.ItemChecklist{ID: uuid.NewString(), MaquinaID: c.Params("id"), Descripcion: description, Orden: i + 1, CreatedAt: time.Now()}
		}
	}
	return c.JSON(v)
}

type detailInput struct {
	ItemChecklistID string `json:"item_checklist_id"`
	Cumple          bool   `json:"cumple"`
	ObservacionItem string `json:"observacion_item"`
}

func (h *MaintenanceHandler) Register(c fiber.Ctx) error {
	form, e := c.MultipartForm()
	if e != nil {
		return bad(c, "Se requiere formulario multipart")
	}
	claims := c.Locals("claims").(middleware.Claims)
	machine := formValue(form, "maquina_id")
	if machine == "" {
		return bad(c, "maquina_id es obligatorio")
	}
	if _, e = h.repo.Machine(machine); e != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Máquina no encontrada"})
	}
	var input []detailInput
	if raw := formValue(form, "detalles"); raw != "" {
		if json.Unmarshal([]byte(raw), &input) != nil {
			return bad(c, "detalles inválidos")
		}
	}
	if len(input) == 0 {
		return bad(c, "Debe registrar al menos un ítem")
	}
	photo, e := saveEvidence(c, form, "foto", []string{".jpg", ".jpeg", ".png"})
	if e != nil {
		return bad(c, e.Error())
	}
	audio, e := saveEvidence(c, form, "audio", []string{".mp3", ".3gp", ".m4a", ".aac"})
	if e != nil {
		return bad(c, e.Error())
	}
	reg := domain.RegistroMantenimiento{ID: uuid.NewString(), TecnicoID: claims.Subject, MaquinaID: machine, FechaHora: time.Now(), FotoURL: photo, AudioURL: audio, Observaciones: formValue(form, "observaciones"), Estado: domain.RegistroCompletado}
	ds := make([]domain.DetalleMantenimiento, len(input))
	for i, x := range input {
		ds[i] = domain.DetalleMantenimiento{ID: uuid.NewString(), RegistroMantenimientoID: reg.ID, ItemChecklistID: x.ItemChecklistID, Cumple: x.Cumple, ObservacionItem: x.ObservacionItem}
	}
	if e = h.repo.Save(&reg, ds); e != nil {
		return internal(c)
	}
	return c.Status(201).JSON(reg)
}
func (h *MaintenanceHandler) History(c fiber.Ctx) error {
	claims := c.Locals("claims").(middleware.Claims)
	v, e := h.repo.History(claims.Subject)
	if e != nil {
		return internal(c)
	}
	return c.JSON(v)
}
func (h *MaintenanceHandler) AdminLogin(c fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if c.Bind().Body(&req) != nil || req.Username == "" || req.Password == "" {
		return bad(c, "Usuario y contraseña obligatorios")
	}
	a := h.auth.UserAccount(req.Username)
	// Permite iniciar una instalación Docker vacía sin transportar la antigua
	// base PROFILE. En cuanto exista un usuario PROFILE, éste conserva prioridad.
	if a.ID == "" && config.Env.BootstrapAdminUser != "" && config.Env.BootstrapAdminPassword != "" &&
		req.Username == config.Env.BootstrapAdminUser && req.Password == config.Env.BootstrapAdminPassword {
		token, e := middleware.CreateJWT("bootstrap-admin", req.Username, "ADMIN")
		if e != nil {
			return internal(c)
		}
		return c.JSON(fiber.Map{"token": token})
	}
	if a.ID == "" || bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(req.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Credenciales inválidas"})
	}
	token, e := middleware.CreateJWT(a.ID, req.Username, "ADMIN")
	if e != nil {
		return internal(c)
	}
	return c.JSON(fiber.Map{"token": token})
}
func (h *MaintenanceHandler) Dashboard(c fiber.Ctx) error {
	a, b, d, chart, e := h.repo.Dashboard()
	if e != nil {
		return internal(c)
	}
	recent, e := h.repo.Records("", "", "")
	if e != nil {
		return internal(c)
	}
	if len(recent) > 5 {
		recent = recent[:5]
	}
	return c.JSON(fiber.Map{"mantenimientos_hoy": a, "maquinas_criticas": b, "tecnicos_activos": d, "por_maquina": chart, "ultimos": recent})
}
func (h *MaintenanceHandler) Records(c fiber.Ctx) error {
	v, e := h.repo.Records(c.Query("fecha"), c.Query("maquina"), c.Query("tecnico"))
	if e != nil {
		return bad(c, "Filtros inválidos")
	}
	return c.JSON(v)
}
func (h *MaintenanceHandler) Record(c fiber.Ctx) error {
	v, e := h.repo.Record(c.Params("id"))
	if errors.Is(e, gorm.ErrRecordNotFound) {
		return c.SendStatus(404)
	}
	if e != nil {
		return internal(c)
	}
	return c.JSON(v)
}
func (h *MaintenanceHandler) MachinesAdmin(c fiber.Ctx) error {
	v, e := h.repo.MachinesPage()
	if e != nil {
		return internal(c)
	}
	return c.JSON(v)
}
func (h *MaintenanceHandler) CreateMachine(c fiber.Ctx) error {
	var v domain.Maquina
	if c.Bind().Body(&v) != nil || v.Codigo == "" || v.Nombre == "" {
		return bad(c, "Código y nombre obligatorios")
	}
	v.ID = uuid.NewString()
	v.CreatedAt = time.Now()
	if v.Estado == "" {
		v.Estado = domain.EstadoOperativa
	}
	if e := h.repo.CreateMachine(&v); e != nil {
		return c.Status(409).JSON(fiber.Map{"error": "Código de máquina duplicado"})
	}
	return c.Status(201).JSON(v)
}
func (h *MaintenanceHandler) UpdateMachine(c fiber.Ctx) error {
	var v domain.Maquina
	if c.Bind().Body(&v) != nil {
		return bad(c, "Datos inválidos")
	}
	v.ID = c.Params("id")
	if e := h.repo.UpdateMachine(&v); e != nil {
		return internal(c)
	}
	return c.JSON(v)
}
func (h *MaintenanceHandler) DeleteMachine(c fiber.Ctx) error {
	if e := h.repo.DeleteMachine(c.Params("id")); e != nil {
		return internal(c)
	}
	return c.SendStatus(204)
}
func (h *MaintenanceHandler) ChecklistAdmin(c fiber.Ctx) error {
	v, e := h.repo.Checklist(c.Params("id"))
	if e != nil {
		return internal(c)
	}
	return c.JSON(v)
}
func (h *MaintenanceHandler) CreateChecklistItem(c fiber.Ctx) error {
	var item domain.ItemChecklist
	if c.Bind().Body(&item) != nil || strings.TrimSpace(item.Descripcion) == "" {
		return bad(c, "La descripción es obligatoria")
	}
	if _, e := h.repo.Machine(c.Params("id")); e != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Máquina no encontrada"})
	}
	item.ID = uuid.NewString()
	item.MaquinaID = c.Params("id")
	item.CreatedAt = time.Now()
	if e := h.repo.CreateChecklistItem(&item); e != nil {
		return internal(c)
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}
func (h *MaintenanceHandler) UpdateChecklistItem(c fiber.Ctx) error {
	var item domain.ItemChecklist
	if c.Bind().Body(&item) != nil || strings.TrimSpace(item.Descripcion) == "" {
		return bad(c, "La descripción es obligatoria")
	}
	item.ID = c.Params("itemId")
	if e := h.repo.UpdateChecklistItem(&item); e != nil {
		return internal(c)
	}
	return c.JSON(item)
}
func (h *MaintenanceHandler) DeleteChecklistItem(c fiber.Ctx) error {
	if e := h.repo.DeleteChecklistItem(c.Params("itemId")); e != nil {
		return internal(c)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (h *MaintenanceHandler) Technicians(c fiber.Ctx) error {
	v, e := h.repo.Technicians()
	if e != nil {
		return internal(c)
	}
	return c.JSON(v)
}
func (h *MaintenanceHandler) CreateTechnician(c fiber.Ctx) error {
	var v domain.Tecnico
	if c.Bind().Body(&v) != nil || v.NumeroEmpleado == "" || v.Nombre == "" {
		return bad(c, "Número y nombre obligatorios")
	}
	v.ID = uuid.NewString()
	v.NumeroEmpleado = strings.ToUpper(strings.TrimSpace(v.NumeroEmpleado))
	v.Activo = true
	v.CreatedAt = time.Now()
	if e := h.repo.CreateTechnician(&v); e != nil {
		return c.Status(409).JSON(fiber.Map{"error": "Número de empleado duplicado"})
	}
	return c.Status(201).JSON(v)
}
func formValue(f *multipart.Form, n string) string {
	if v := f.Value[n]; len(v) > 0 {
		return v[0]
	}
	return ""
}
func saveEvidence(c fiber.Ctx, f *multipart.Form, n string, exts []string) (string, error) {
	files := f.File[n]
	if len(files) == 0 {
		return "", nil
	}
	file := files[0]
	if file.Size > 15*1024*1024 {
		return "", errors.New("archivo demasiado grande")
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	ok := false
	for _, x := range exts {
		if ext == x {
			ok = true
		}
	}
	if !ok {
		return "", errors.New("tipo de archivo no permitido")
	}
	folder := filepath.Join(config.Env.EvidenceDirectory, n)
	if e := os.MkdirAll(folder, 0750); e != nil {
		return "", e
	}
	name := uuid.NewString() + ext
	if e := c.SaveFile(file, filepath.Join(folder, name)); e != nil {
		return "", e
	}
	return "/uploads/mantenimiento/" + n + "/" + name, nil
}
func bad(c fiber.Ctx, m string) error { return c.Status(400).JSON(fiber.Map{"error": m}) }
func internal(c fiber.Ctx) error      { return c.Status(500).JSON(fiber.Map{"error": "Error interno"}) }
