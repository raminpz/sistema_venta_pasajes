# Sistema de Venta de Pasajes - API REST

## Checklist de lectura
- [x] Arquitectura actual del proyecto
- [x] Estructura de carpetas por características
- [x] Tecnologías base
- [x] Flujo de arranque de la aplicación
- [x] Manejo centralizado de errores
- [x] Feature implementada actualmente
- [x] Convenciones para siguientes módulos
- [x] Comandos de ejecución y prueba

## Resumen
Este proyecto implementa un **API REST en Go** para un **sistema de venta de pasajes terrestres**.

La arquitectura ya no está basada en CQRS ni en capas separadas globalmente por `models`, `repositories`, `services` y `handlers`.

Ahora el proyecto usa una **estructura basada en características**:
- cada módulo vive dentro de su propia carpeta
- cada feature contiene su modelo, su repository, su service y su handler
- la infraestructura compartida queda fuera de la feature

## Stack actual
- **Go**
- **Gorilla Mux**
- **GORM**
- **MySQL 8**

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

## 5. Feature implementada actualmente: `PROVEEDOR_SISTEMA`
La primera feature ya implementada es `PROVEEDOR_SISTEMA`.

### Estructura
```text
internal/proveedor_sistema/
├── domain/
│   └── proveedor_sistema.go
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

### Endpoints disponibles
Bajo `/api/v1`:
- `GET /api/v1/proveedor`
- `GET /api/v1/proveedor/{id}`
- `POST /api/v1/proveedor`
- `PUT /api/v1/proveedor/{id}`
- `DELETE /api/v1/proveedor/{id}`

> `DELETE` realiza borrado lógico usando `DELETED_AT`; el registro no se elimina físicamente de la base de datos.

### Validaciones implementadas
- RUC obligatorio y de 11 dígitos
- razón social obligatoria
- validación de longitudes máximas
- email válido si se envía
- web válida si se envía
- ID válido en operaciones por identificador

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

## 8. Convenciones para los siguientes módulos
A partir de ahora, cada nueva tabla o feature debe implementarse siguiendo esta regla:

### Orden por feature
1. crear carpeta `internal/<feature>`
2. crear `<feature>.go`
3. crear `repository/repository.go`
4. crear `service/service.go`
5. crear `handler/handler.go`
6. crear `handler/routes.go`
7. agregar tests mínimos de `service` y `handler`
8. registrar la feature desde `internal/http/routes/router.go`

### Reglas técnicas
- usar **GORM** en repository
- usar `context.Context`
- usar errores centralizados
- usar respuestas JSON estandarizadas
- usar `.Methods("GET")`, `.Methods("POST")`, `.Methods("PUT")`, etc.
- respetar el esquema MySQL ya existente
- no agregar operaciones no definidas funcionalmente para la feature

---

## 9. Orden de desarrollo recomendado según el esquema
### Tablas independientes
1. `PROVEEDOR_SISTEMA`
2. `EMPRESA`
3. `TERMINAL`
4. `TIPO_VEHICULO`
5. `ROL`
6. `PASAJERO`
7. `TIPO_COMPROBANTE`
8. `METODO_PAGO`

### Tablas con dependencias simples
9. `LICENCIA_SISTEMA`
10. `RUTA`
11. `VEHICULO`
12. `ASIENTO`
13. `CONDUCTOR`
14. `USUARIO`

### Tablas transaccionales
15. `VIAJE`
16. `VENTA`
17. `PAGO`
18. `DETALLE_PASAJE`
19. `ENCOMIENDA`

---

## 10. Variables de entorno sugeridas
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

## 11. Comandos útiles
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
