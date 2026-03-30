package domain

type Ruta struct {
	IDRuta            int     `json:"id_ruta" gorm:"column:ID_RUTA;primaryKey;autoIncrement"`
	IDOrigenTerminal  int     `json:"id_origen_terminal" gorm:"column:ID_ORIGEN_TERMINAL"`
	IDDestinoTerminal int     `json:"id_destino_terminal" gorm:"column:ID_DESTINO_TERMINAL"`
	DuracionHoras     float64 `json:"duracion_horas" gorm:"column:DURACION_HORAS"`
}

func (Ruta) TableName() string {
	return "RUTA"
}

