package database

import (
	"fmt"
	"strings"
)

func EnsureMaintenanceSchema() error {
	tables := []struct {
		name string
		cols []string
	}{
		{"MT_TECNICOS", []string{"ID VARCHAR(36) NOT NULL", "NUMERO_EMPLEADO VARCHAR(64) NOT NULL", "NOMBRE VARCHAR(160) NOT NULL", "AREA VARCHAR(100)", "ACTIVO BOOLEAN", "CREATED_AT TIMESTAMP"}},
		{"MT_MAQUINAS", []string{"ID VARCHAR(36) NOT NULL", "CODIGO VARCHAR(64) NOT NULL", "NOMBRE VARCHAR(160) NOT NULL", "UBICACION VARCHAR(160)", "ESTADO VARCHAR(32)", "FRECUENCIA_MANTENIMIENTO VARCHAR(32)", "CREATED_AT TIMESTAMP"}},
		{"MT_ITEMS_CHECKLIST", []string{"ID VARCHAR(36) NOT NULL", "MAQUINA_ID VARCHAR(36) NOT NULL", "DESCRIPCION VARCHAR(500) NOT NULL", "ORDEN INTEGER", "CREATED_AT TIMESTAMP"}},
		{"MT_REGISTROS", []string{"ID VARCHAR(36) NOT NULL", "TECNICO_ID VARCHAR(36) NOT NULL", "MAQUINA_ID VARCHAR(36) NOT NULL", "FECHA_HORA TIMESTAMP", "FOTO_URL VARCHAR(512)", "AUDIO_URL VARCHAR(512)", "OBSERVACIONES BLOB SUB_TYPE TEXT", "ESTADO VARCHAR(32)"}},
		{"MT_DETALLES", []string{"ID VARCHAR(36) NOT NULL", "REGISTRO_MANTENIMIENTO_ID VARCHAR(36) NOT NULL", "ITEM_CHECKLIST_ID VARCHAR(36) NOT NULL", "CUMPLE BOOLEAN", "OBSERVACION_ITEM BLOB SUB_TYPE TEXT"}}}
	for _, t := range tables {
		exists, err := tableExists(t.name)
		if err != nil {
			return err
		}
		if !exists {
			if err := DB.Exec(fmt.Sprintf("CREATE TABLE %s (%s)", t.name, strings.Join(t.cols, ", "))).Error; err != nil {
				return err
			}
			continue
		}
		for _, col := range t.cols {
			name := strings.Fields(col)[0]
			found, err := columnExists(t.name, name)
			if err != nil {
				return err
			}
			if !found {
				if err := DB.Exec(fmt.Sprintf("ALTER TABLE %s ADD %s", t.name, col)).Error; err != nil {
					return err
				}
			}
		}
	}
	for name, sql := range map[string]string{"MT_TEC_NUM_UQ": "CREATE UNIQUE INDEX MT_TEC_NUM_UQ ON MT_TECNICOS (NUMERO_EMPLEADO)", "MT_MAQ_COD_UQ": "CREATE UNIQUE INDEX MT_MAQ_COD_UQ ON MT_MAQUINAS (CODIGO)", "MT_REG_TEC_IDX": "CREATE INDEX MT_REG_TEC_IDX ON MT_REGISTROS (TECNICO_ID)", "MT_REG_MAQ_IDX": "CREATE INDEX MT_REG_MAQ_IDX ON MT_REGISTROS (MAQUINA_ID)"} {
		exists, err := indexExists(name)
		if err != nil {
			return err
		}
		if !exists {
			if err := DB.Exec(sql).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
func tableExists(n string) (bool, error) {
	var c int64
	err := DB.Raw("SELECT COUNT(*) FROM RDB$RELATIONS WHERE RDB$RELATION_NAME = ?", n).Scan(&c).Error
	return c > 0, err
}
func columnExists(t, n string) (bool, error) {
	var c int64
	err := DB.Raw("SELECT COUNT(*) FROM RDB$RELATION_FIELDS WHERE RDB$RELATION_NAME = ? AND RDB$FIELD_NAME = ?", t, n).Scan(&c).Error
	return c > 0, err
}
func indexExists(n string) (bool, error) {
	var c int64
	err := DB.Raw("SELECT COUNT(*) FROM RDB$INDICES WHERE RDB$INDEX_NAME = ?", n).Scan(&c).Error
	return c > 0, err
}
