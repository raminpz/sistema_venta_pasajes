package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sistema_venta_pasajes/pkg"
	"strconv"

	"github.com/gorilla/mux"

	"sistema_venta_pasajes/internal/proveedor/domain"
	providerinput "sistema_venta_pasajes/internal/proveedor/input"
)

type Service interface {
	List(ctx context.Context) ([]domain.ProveedorSistema, error)
	GetByID(ctx context.Context, id int64) (*domain.ProveedorSistema, error)
	Create(ctx context.Context, input providerinput.CreateInput) (*domain.ProveedorSistema, error)
	Update(ctx context.Context, id int64, input providerinput.UpdateInput) (*domain.ProveedorSistema, error)
	Delete(ctx context.Context, id int64) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func toProveedorOutput(p *domain.ProveedorSistema) providerinput.ProveedorOutput {
	return providerinput.ProveedorOutput{
		ID:              p.IDProveedor,
		RUC:             p.RUC,
		RazonSocial:     p.RazonSocial,
		NombreComercial: p.NombreComercial,
		Direccion:       p.Direccion,
		Telefono:        p.Telefono,
		Email:           p.Email,
		Web:             p.Web,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	proveedores, err := h.service.List(r.Context())
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	outputs := make([]providerinput.ProveedorOutput, 0, len(proveedores))
	for _, p := range proveedores {
		outputs = append(outputs, toProveedorOutput(&p))
	}
	// Si no hay resultados, devolver array vacío
	var data any = outputs
	if len(outputs) == 0 {
		data = []interface{}{}
	}
	pkg.WriteSuccess(w, http.StatusOK, "proveedores del sistema obtenidos correctamente", data, map[string]any{
		"count": len(outputs),
	})
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := proveedorSistemaIDFromRequest(r)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	proveedor, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, "proveedor del sistema obtenido correctamente", toProveedorOutput(proveedor), nil)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input providerinput.CreateInput
	if err := decodeJSONBody(r, &input); err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	proveedor, err := h.service.Create(r.Context(), input)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusCreated, "proveedor del sistema creado correctamente", toProveedorOutput(proveedor), nil)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := proveedorSistemaIDFromRequest(r)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	var input providerinput.UpdateInput
	if err := decodeJSONBody(r, &input); err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	proveedor, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}
	pkg.WriteSuccess(w, http.StatusOK, "proveedor del sistema actualizado correctamente", toProveedorOutput(proveedor), nil)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := proveedorSistemaIDFromRequest(r)
	if err != nil {
		pkg.WriteError(w, r, err)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		pkg.WriteError(w, r, err)
		return
	}

	pkg.WriteSuccess(w, http.StatusOK, "proveedor del sistema eliminado correctamente", nil, nil)
}

func proveedorSistemaIDFromRequest(r *http.Request) (int64, error) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		return 0, pkg.BadRequest("invalid_provider_id", "el id del proveedor del sistema no es válido")
	}
	return id, nil
}

func decodeJSONBody(r *http.Request, destination any) error {
	defer func() {
		_ = r.Body.Close()
	}()

	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return mapDecodeError(err)
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return pkg.BadRequest("invalid_json", "el cuerpo JSON debe contener un único objeto")
	}

	return nil
}

func mapDecodeError(err error) error {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.Is(err, io.EOF):
		return pkg.BadRequest("empty_body", "el cuerpo de la solicitud es obligatorio")
	case errors.As(err, &syntaxError):
		return pkg.BadRequest("invalid_json", "el cuerpo JSON no tiene un formato válido")
	case errors.As(err, &unmarshalTypeError):
		return pkg.BadRequest("invalid_json_type", "uno o más campos del JSON tienen un tipo inválido")
	default:
		return pkg.BadRequest("invalid_json", "el cuerpo JSON es inválido")
	}
}
