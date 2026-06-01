package domain

type TipoVehiculo struct {
	IDTipoVehiculo int64  `gorm:"column:ID_TIPO_VEHICULO;primaryKey;autoIncrement" json:"id_tipo_vehiculo"`
	Nombre         string `gorm:"column:NOMBRE" json:"nombre"`
	Descripcion    string `gorm:"column:DESCRIPCION" json:"descripcion"`
}

func (TipoVehiculo) TableName() string {
	return "TIPO_VEHICULO"
}

