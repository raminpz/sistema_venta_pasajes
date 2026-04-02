# Sistema de Venta de Pasajes - API REST

## 📚 Documentación Importante
**PRIMERO lee estos documentos:**
- ✅ [`STANDARDS.md`](./STANDARDS.md) - Guía de estándares (Mux + GORM)
- ✅ [`QUICK_REFERENCE.md`](./QUICK_REFERENCE.md) - Referencia rápida para nuevo módulo
- ✅ [`STANDARDS_SUMMARY.md`](./STANDARDS_SUMMARY.md) - Resumen ejecutivo
- ✅ [`API_ROUTES.md`](./API_ROUTES.md) - Documentación de todas las rutas
- ✅ [`VALIDATION_CHECKLIST.md`](./VALIDATION_CHECKLIST.md) - Estado de módulos

## Checklist de lectura
- [x] Arquitectura actual del proyecto
- [x] Estructura de carpetas por características
- [x] Tecnologías base (Go, Gorilla Mux, GORM, MySQL)
- [x] Flujo de arranque de la aplicación
- [x] Manejo centralizado de errores
- [x] Features implementadas (13 módulos)
- [x] Convenciones para siguientes módulos
- [x] Comandos de ejecución y prueba
- [x] **NUEVO**: Estándares Mux y GORM para todos los módulos

## Resumen
Este proyecto implementa un **API REST en Go** para un **sistema de venta de pasajes terrestres**.

La arquitectura ya no está basada en CQRS ni en capas separadas globalmente por `models`, `repositories`, `services` y `handlers`.

Ahora el proyecto usa una **estructura basada en características**:
- cada módulo vive dentro de su propia carpeta
- cada feature contiene su modelo, su repository, su service y su handler
- la infraestructura compartida queda fuera de la feature

## Stack actual
- **Go** 1.25.0
- **Gorilla Mux** v1.8.1 (Router HTTP)
- **GORM** v1.31.1 (ORM)
- **MySQL 8** (Base de datos)
- **Gorilla Handlers** v1.5.2 (CORS)

---

## 1. Objetivo de la arquitectura actual
La idea de esta estructura es que cada funcionalidad sea más fácil de entender y mantener.

En vez de tener código repartido por capas globales, se agrupa por feature.

### Beneficios
- mayor cohesión por módulo
- menos saltos entre carpetas
- más simple de leer para desarrollo incremental
- más fácil de ubicar modelo + lógica + acceso a datos + handler HTTP
- mejor base para seguir el orden real del esquema SQL

---

## 2. Estructura actual del proyecto
```text
proyecto/
├── cmd/
│   └── app/
│       └── main.go
├── config/
│   ├── .env.example
│   ├── README.md
│   └── schema_mysql8.sql
├── internal/
│   ├── bootstrap/
│   │   └── app.go
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── mysql.go
│   ├── http/
│   │   ├── middleware/
│   │   │   ├── recover.go
│   │   │   └── request_id.go
│   │   └── routes/
│   │       ├── router.go
│   │       └── router_test.go
│   ├── proveedor_sistema/
│   │   ├── domain/
│   │   │   └── proveedor_sistema.go
│   │   ├── handler/
│   │   │   ├── handler.go
│   │   │   ├── handler_test.go
│   │   │   └── routes.go
│   │   ├── repository/
│   │   │   └── repository.go
│   │   ├── service/
│   │   │   ├── service.go
│   │   │   └── service_test.go
│   └── shared/
│       ├── apperror.go
│       └── response.go
├── pkg/
├── go.mod
├── go.sum
├── main.go
└── README.md
```

---

## 3. Qué queda como infraestructura compartida
Estas carpetas no pertenecen a una feature específica:

### `internal/bootstrap`
Arranque principal de la aplicación.

### `internal/config`
Carga de variables de entorno y configuración general.

### `internal/database`
Conexión a MySQL mediante GORM.

### `internal/http/middleware`
Middleware compartido:
- `request_id`
- `recover`

### `internal/http/routes`
Router principal de la aplicación y endpoints transversales:
- `/health`
- `/ready`

### `internal/shared`
Utilidades comunes para toda la API:
- respuestas JSON estándar
- errores centralizados

---

## 4. Estructura de una feature
Cada feature debe seguir este patrón:

```text
internal/
└── nombre_feature/
    ├── domain/
    │   └── nombre_feature.go
    ├── input/
    │   └── input.go
    ├── handler/
    │   ├── handler.go
    │   ├── handler_test.go
    │   └── routes.go
    ├── repository/
    │   └── repository.go
    ├── service/
    │   ├── service.go
    │   └── service_test.go
```

### Responsabilidad de cada parte
- `domain/nombre_feature.go`: entidad de dominio de la feature
- `input/input.go`: contratos de entrada usados por `handler`, `service` y `repository`
- `repository/`: acceso a datos con GORM
- `service/`: lógica de negocio y validaciones
- `handler/`: capa HTTP y registro de rutas de la feature

---

## 5. Módulos Implementados (13 módulos)

### ✅ Todos con Mux + GORM

| Módulo | Funcionalidad | Status |
|--------|---------------|--------|
| **Asiento** | Gestión de asientos de vehículos | ✅ Listo |
| **Conductor** | Gestión de conductores | ✅ Listo |
| **Empresa** | Gestión de empresas de transporte | ✅ Listo |
| **Licencia** | Licencias del sistema | ✅ Listo |
| **Pago** | Gestión de pagos | ✅ Listo |
| **Pasajero** | Gestión de pasajeros | ✅ Listo |
| **Programación** | Programaciones de viajes | ✅ Listo |
| **Proveedor** | Gestión de proveedores | ✅ Listo |
| **Ruta** | Gestión de rutas | ✅ Listo |
| **Terminal** | Gestión de terminales | ✅ Listo |
| **Usuario** | Gestión de usuarios | ✅ Listo |
| **Vehiculo** | Gestión de vehículos | ✅ Listo |
| **Venta** | Gestión de ventas | ✅ Listo |

### Endpoints RESTful
Cada módulo implementa:
```
POST   /api/v1/{modulo}           - Crear
GET    /api/v1/{modulo}           - Listar (con paginación)
GET    /api/v1/{modulo}/{id}      - Obtener por ID
PUT    /api/v1/{modulo}/{id}      - Actualizar
DELETE /api/v1/{modulo}/{id}      - Eliminar
```

**Ver documentación completa en** [`API_ROUTES.md`](./API_ROUTES.md)

---

## 6. Flujo de arranque actual
1. `main.go` o `cmd/app/main.go`
2. `internal/bootstrap/app.go`
3. carga de configuración desde `internal/config`
4. conexión GORM desde `internal/database`
5. construcción de router desde `internal/http/routes`
6. registro de features en el router principal
7. arranque del servidor HTTP

---

## 7. Manejo centralizado de errores
El proyecto ya tiene un manejo centralizado y personalizado de errores.

### Archivos clave
- `internal/shared/apperror.go`
- `internal/shared/response.go`

### Formato de error
```json
{
  "code": 422,
  "message": "los datos enviados no son válidos",
  "error": "validation_error"
}
```

### Formato de éxito
```json
{
  "code": 200,
  "message": "operación realizada correctamente",
  "data": {},
  "meta": {}
}
```

### Middlewares activos
- `X-Request-ID`
- recuperación de pánicos

---

## 6. Convenciones para todos los módulos - Mux y GORM

### ⚠️ IMPORTANTE: Estándares Obligatorios

**TODOS los módulos DEBEN usar:**
- ✅ **Gorilla Mux** (`github.com/gorilla/mux`) para enrutamiento HTTP
- ✅ **GORM** (`gorm.io/gorm`) para persistencia de datos

### Estructura estándar de módulo
```text
internal/nombre_modulo/
├── domain/
│   └── entity.go                    # Modelo con GORM tags
├── repository/
│   └── repository.go                # Interfaz + implementación GORM
├── service/
│   └── service.go                   # Lógica de negocio
├── handler/
│   ├── handler.go                   # Handlers HTTP con Mux
│   └── register.go                  # Registro de rutas
├── input/
│   └── input.go                     # DTOs
└── util/
    ├── constants.go
    └── validation.go
```

### Handler (SIEMPRE con firma de Mux)
```go
func (h *YourHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]  // ✓ Usar mux.Vars
	// lógica...
}
```

### Repository (SIEMPRE con GORM)
```go
func (r *repo) Create(entity *domain.Entity) error {
	return r.db.Create(entity).Error  // ✓ Usar r.db (GORM)
}
```

### Registro de rutas (SIEMPRE con RegisterRoutes)
```go
func RegisterRoutes(r *mux.Router, db *gorm.DB) {
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := NewHandler(svc)
	r.HandleFunc("/modulo", h.Create).Methods("POST")
	// ... más rutas
}
```

### ✅ Validación
- [ ] ¿Usa `github.com/gorilla/mux`?
- [ ] ¿Tiene `RegisterRoutes(r *mux.Router, db *gorm.DB)`?
- [ ] ¿Repository usa GORM (`*gorm.DB`)?
- [ ] ¿Las queries usan placeholders (`WHERE "id = ?", value`)?
- [ ] ¿NO tiene concatenación de strings en queries?

**Si respondiste SÍ a todo → ✅ CUMPLE ESTÁNDARES**

### 📖 Referencia Rápida para Nuevo Módulo
Ver [`QUICK_REFERENCE.md`](./QUICK_REFERENCE.md)

---

## 8. Orden de desarrollo recomendado según el esquema
### Tablas independientes
1. ✅ `PROVEEDOR_SISTEMA` - Implementado
2. ✅ `EMPRESA` - Implementado
3. ✅ `TERMINAL` - Implementado
4. `TIPO_VEHICULO` - Por hacer
5. `ROL` - Por hacer
6. ✅ `PASAJERO` - Implementado
7. `TIPO_COMPROBANTE` - Por hacer
8. `METODO_PAGO` - Por hacer

### Tablas con dependencias simples
9. ✅ `LICENCIA_SISTEMA` - Implementado
10. ✅ `RUTA` - Implementado
11. ✅ `VEHICULO` - Implementado
12. ✅ `ASIENTO` - Implementado
13. ✅ `CONDUCTOR` - Implementado
14. ✅ `USUARIO` - Implementado

### Tablas transaccionales
15. `VIAJE` - Por hacer (puede estar como PROGRAMACION)
16. ✅ `VENTA` - Implementado
17. ✅ `PAGO` - Implementado
18. `DETALLE_PASAJE` - Por hacer
19. `ENCOMIENDA` - Por hacer

**Status Actual**: 13/19 módulos implementados (68%)

---

## 9. Variables de entorno sugeridas
```env
APP_PORT=8080
APP_ENV=development

DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=SISTEMA_PASAJES
DB_USER=root
DB_PASS=
DB_PARAMS=parseTime=true&loc=Local&charset=utf8mb4
DB_MAX_OPEN_CONNS=10
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME_MIN=30

HTTP_READ_TIMEOUT=10
HTTP_WRITE_TIMEOUT=10
```

---

## 10. Comandos útiles
### Ejecutar tests
```powershell
Set-Location "C:\Users\Rami\GolandProjects\sistema_venta_pasajes"
go test ./...
```

### Compilar
```powershell
Set-Location "C:\Users\Rami\GolandProjects\sistema_venta_pasajes"
go build ./...
```

### Ejecutar la aplicación
```powershell
Set-Location "C:\Users\Rami\GolandProjects\sistema_venta_pasajes"
go run .
```

O también:

```powershell
Set-Location "C:\Users\Rami\GolandProjects\sistema_venta_pasajes"
go run ./cmd/app
```

---

## 12. Estado actual
Actualmente el proyecto ya tiene:
- bootstrap funcional
- conexión MySQL con GORM
- router base
- middlewares compartidos
- manejo centralizado de errores
- feature `PROVEEDOR_SISTEMA` funcional
- tests unitarios y de rutas

## 13. Siguiente paso sugerido
Continuar con la siguiente feature del esquema:
- `EMPRESA`

Y desarrollarla siguiendo exactamente la misma estructura basada en características.
